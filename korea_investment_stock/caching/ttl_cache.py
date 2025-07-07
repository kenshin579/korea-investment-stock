"""
TTL (Time To Live) 캐시 구현

한국투자증권 API 응답을 캐싱하여 API 호출 횟수를 줄이고 성능을 향상시킵니다.
"""

import time
import threading
import sys
import zlib
import pickle
import logging
from typing import Dict, Any, Optional, Callable, Tuple
from functools import wraps
from collections import defaultdict
from datetime import datetime, timedelta
from functools import lru_cache

logger = logging.getLogger(__name__)


class CacheEntry:
    """캐시 항목 데이터 구조"""
    
    def __init__(self, value: Any, ttl: int):
        """
        Args:
            value: 캐시할 값
            ttl: Time To Live (초)
        """
        self.value = value
        self.expires_at = time.time() + ttl
        self.created_at = time.time()
        self.access_count = 0
        self.last_accessed = time.time()
        self.size = sys.getsizeof(value)
    
    def is_expired(self) -> bool:
        """만료 여부 확인"""
        return time.time() > self.expires_at
    
    def access(self) -> Any:
        """값 접근 시 통계 업데이트"""
        self.access_count += 1
        self.last_accessed = time.time()
        return self.value
    
    def get_age(self) -> float:
        """항목 나이 (초)"""
        return time.time() - self.created_at


class TTLCache:
    """Thread-safe TTL 캐시"""
    
    def __init__(self, default_ttl: int = 300, max_size: int = 10000):
        """
        Args:
            default_ttl: 기본 TTL (초, 기본값: 300초=5분)
            max_size: 최대 캐시 항목 수
        """
        self._cache: Dict[str, CacheEntry] = {}
        self._default_ttl = default_ttl
        self._max_size = max_size
        self._lock = threading.RLock()
        
        # 통계
        self._hit_count = 0
        self._miss_count = 0
        self._eviction_count = 0
        self._access_count = defaultdict(int)
        
        # 메모리 관리
        self._total_memory = 0
        self._cleanup_interval = 60  # 1분마다 정리
        self._last_cleanup = time.time()
        self._expired_count = 0
        
        # 백그라운드 정리 스레드
        self._stop_cleanup = False
        self._cleanup_thread = None
        self.start_cleanup_thread()
        
        # 압축 설정
        self._compression_threshold = 1024  # 1KB 이상 압축
        
        logger.info(f"TTLCache 초기화: default_ttl={default_ttl}s, max_size={max_size}")
    
    def get(self, key: str) -> Optional[Any]:
        """캐시에서 값 조회"""
        with self._lock:
            # 주기적 정리
            if time.time() - self._last_cleanup > self._cleanup_interval:
                self._cleanup_expired()
            
            # 원본 키와 압축된 키 모두 확인
            entry = None
            actual_key = key
            
            if key in self._cache:
                entry = self._cache[key]
            elif f"{key}:compressed" in self._cache:
                actual_key = f"{key}:compressed"
                entry = self._cache[actual_key]
            
            if entry is None:
                self._miss_count += 1
                logger.debug(f"Cache MISS for {key}")
                return None
            
            # 만료 확인
            if entry.is_expired():
                self._miss_count += 1
                logger.debug(f"Cache MISS for {key} (expired)")
                del self._cache[actual_key]
                self._total_memory -= entry.size
                return None
            
            # 히트
            self._hit_count += 1
            self._access_count[key] += 1
            logger.debug(f"Cache HIT for {key}")
            
            # 압축된 데이터 복원
            value = entry.access()
            if isinstance(value, bytes) and actual_key.endswith(":compressed"):
                value = pickle.loads(zlib.decompress(value))
            
            return value
    
    def set(self, key: str, value: Any, ttl: Optional[int] = None) -> None:
        """캐시에 값 저장"""
        with self._lock:
            # TTL 설정
            if ttl is None:
                ttl = self._default_ttl
            
            # 크기 제한 확인
            if len(self._cache) >= self._max_size:
                self._evict_lru()
            
            # 기존 항목 제거
            if key in self._cache:
                old_entry = self._cache[key]
                self._total_memory -= old_entry.size
            
            # 압축 처리
            compressed = False
            original_size = sys.getsizeof(value)
            if original_size > self._compression_threshold:
                compressed_value = zlib.compress(pickle.dumps(value))
                compressed_size = sys.getsizeof(compressed_value)
                if compressed_size < original_size * 0.8:  # 20% 이상 압축 효과
                    value = compressed_value
                    key = f"{key}:compressed"
                    compressed = True
                    logger.debug(f"Compressed {original_size} -> {compressed_size} bytes ({compressed_size/original_size:.1%})")
            
            # 새 항목 추가
            entry = CacheEntry(value, ttl)
            self._cache[key] = entry
            self._total_memory += entry.size
            
            logger.debug(f"Cache SET for {key}, TTL={ttl}s, compressed={compressed}")
    
    def delete(self, key: str) -> bool:
        """캐시에서 항목 삭제"""
        with self._lock:
            # 압축된 키도 확인
            keys_to_check = [key, f"{key}:compressed"]
            
            for k in keys_to_check:
                if k in self._cache:
                    entry = self._cache[k]
                    self._total_memory -= entry.size
                    del self._cache[k]
                    logger.debug(f"Cache DELETE for {k}")
                    return True
            
            return False
    
    def clear(self, pattern: Optional[str] = None) -> int:
        """캐시 비우기
        
        Args:
            pattern: 삭제할 키 패턴 (None이면 전체 삭제)
            
        Returns:
            삭제된 항목 수
        """
        with self._lock:
            if pattern is None:
                # 전체 삭제
                count = len(self._cache)
                self._cache.clear()
                self._total_memory = 0
                logger.info(f"Cache CLEAR: {count} items removed")
                return count
            
            # 패턴 매칭 삭제
            import fnmatch
            keys_to_delete = []
            
            for key in self._cache:
                if fnmatch.fnmatch(key, pattern):
                    keys_to_delete.append(key)
            
            for key in keys_to_delete:
                entry = self._cache[key]
                self._total_memory -= entry.size
                del self._cache[key]
            
            logger.info(f"Cache CLEAR pattern '{pattern}': {len(keys_to_delete)} items removed")
            return len(keys_to_delete)
    
    def _cleanup_expired(self) -> None:
        """만료된 항목 정리"""
        with self._lock:
            expired_keys = []
            for key, entry in list(self._cache.items()):
                if entry.is_expired():
                    expired_keys.append(key)
            
            for key in expired_keys:
                self._remove_entry(key)
                self._expired_count += 1
            
            self._last_cleanup = time.time()
            if expired_keys:
                logger.debug(f"{len(expired_keys)}개의 만료된 항목 제거")
    
    def _evict_lru(self) -> None:
        """LRU (Least Recently Used) 정책으로 항목 제거"""
        if not self._cache:
            return
        
        # 가장 오래 사용하지 않은 항목 찾기
        oldest_key = min(self._cache.keys(), 
                        key=lambda k: self._cache[k].last_accessed)
        self._remove_entry(oldest_key)
        self._eviction_count += 1
        logger.debug(f"LRU 제거: {oldest_key}")
    
    def _evict_lfu(self) -> None:
        """LFU (Least Frequently Used) 정책으로 항목 제거"""
        if not self._cache:
            return
        
        # 가장 적게 사용된 항목 찾기
        least_used_key = min(self._cache.keys(), 
                            key=lambda k: self._cache[k].access_count)
        self._remove_entry(least_used_key)
        self._eviction_count += 1
        logger.debug(f"LFU 제거: {least_used_key}")
    
    def get_stats(self) -> dict:
        """캐시 통계 조회"""
        with self._lock:
            total_requests = self._hit_count + self._miss_count
            hit_rate = self._hit_count / total_requests if total_requests > 0 else 0
            
            # 평균 항목 나이 계산
            ages = [entry.get_age() for entry in self._cache.values()]
            avg_age = sum(ages) / len(ages) if ages else 0
            
            return {
                'hit_count': self._hit_count,
                'miss_count': self._miss_count,
                'hit_rate': hit_rate,
                'miss_rate': 1 - hit_rate,
                'total_entries': len(self._cache),
                'eviction_count': self._eviction_count,
                'avg_entry_age': avg_age,
                'memory_usage_mb': self._total_memory / (1024 * 1024),
                'api_calls_saved': self._hit_count,
            }
    
    def print_stats(self) -> None:
        """통계 출력"""
        stats = self.get_stats()
        print(f"\n📊 캐시 통계:")
        print(f"- 캐시 적중률: {stats['hit_rate']:.1%}")
        print(f"- 총 항목 수: {stats['total_entries']:,}")
        print(f"- 메모리 사용량: {stats['memory_usage_mb']:.2f} MB")
        print(f"- API 호출 절감: {stats['api_calls_saved']:,}회")
        print(f"- 평균 항목 나이: {stats['avg_entry_age']:.1f}초")
        print(f"- 제거된 항목: {stats['eviction_count']:,}개")
    
    @property
    def hit_rate(self) -> float:
        """캐시 적중률"""
        total = self._hit_count + self._miss_count
        return self._hit_count / total if total > 0 else 0
    
    @property
    def memory_usage(self) -> float:
        """메모리 사용량 (MB)"""
        return self._total_memory / (1024 * 1024)
    
    @property
    def expired_count(self) -> int:
        """만료된 항목 수"""
        with self._lock:
            return sum(1 for entry in self._cache.values() if entry.is_expired())

    def start_cleanup_thread(self):
        """백그라운드 정리 스레드 시작"""
        if not self._cleanup_thread or not self._cleanup_thread.is_alive():
            self._cleanup_thread = threading.Thread(target=self._cleanup_worker, daemon=True)
            self._cleanup_thread.start()
            logger.debug("캐시 정리 스레드 시작")
    
    def stop_cleanup_thread(self):
        """백그라운드 정리 스레드 정지"""
        self._stop_cleanup = True
        if self._cleanup_thread and self._cleanup_thread.is_alive():
            self._cleanup_thread.join(timeout=5)
            logger.debug("캐시 정리 스레드 정지")
    
    def _cleanup_worker(self):
        """백그라운드 정리 워커"""
        while not self._stop_cleanup:
            time.sleep(self._cleanup_interval)
            if not self._stop_cleanup:
                self._cleanup_expired()
                logger.debug(f"캐시 정리 완료 - 만료된 항목: {self._expired_count}개")

    def _remove_entry(self, key: str):
        """캐시 항목 제거 (내부 사용)"""
        if key in self._cache:
            entry = self._cache[key]
            self._total_memory -= entry.size
            del self._cache[key]


# API별 기본 TTL 설정 (실제 메서드 기준)
CACHE_TTL_CONFIG = {
    'fetch_domestic_price': 300,            # 5분
    'fetch_etf_domestic_price': 300,        # 5분
    'fetch_price_list': 300,                # 5분
    'fetch_price_detail_oversea_list': 300, # 5분
    'fetch_stock_info_list': 18000,         # 5시간
    'fetch_search_stock_info_list': 18000,  # 5시간
    'fetch_kospi_symbols': 259200,          # 3일
    'fetch_kosdaq_symbols': 259200,         # 3일
    'fetch_symbols': 259200,                # 3일
}


def generate_cache_key(method_name: str, *args, **kwargs) -> str:
    """
    메서드명과 파라미터를 조합하여 유니크한 캐시 키 생성
    
    예시:
    - fetch_domestic_price:J:005930
    - fetch_etf_domestic_price:J:294400
    - fetch_stock_info:005930:KR
    - fetch_kospi_symbols
    """
    key_parts = [method_name]
    key_parts.extend(str(arg) for arg in args)
    key_parts.extend(f"{k}={v}" for k, v in sorted(kwargs.items()))
    return ":".join(key_parts)


def cacheable(ttl: Optional[int] = None, 
              cache_condition: Optional[Callable] = None,
              key_generator: Optional[Callable] = None,
              use_dynamic_ttl: bool = True):
    """
    메서드에 캐싱 기능을 추가하는 데코레이터
    
    Args:
        ttl: 이 메서드의 TTL (None이면 기본값 사용)
        cache_condition: 캐시 여부를 결정하는 함수
        key_generator: 커스텀 캐시 키 생성 함수
        use_dynamic_ttl: 동적 TTL 사용 여부
    
    사용 예:
    @cacheable(ttl=300)  # 5분
    def fetch_domestic_price(self, market_code: str, symbol: str) -> dict:
        ...
    
    @cacheable(ttl=259200, cache_condition=lambda result: result.get('rt_cd') == '0')  # 3일
    def fetch_kospi_symbols(self) -> pd.DataFrame:
        ...
    """
    def decorator(func):
        @wraps(func)
        def wrapper(self, *args, **kwargs):
            # use_cache 파라미터 확인
            use_cache = kwargs.pop('use_cache', True)
            
            # 캐시가 비활성화되어 있거나 use_cache=False인 경우
            if not hasattr(self, '_cache') or not self._cache or not use_cache:
                return func(self, *args, **kwargs)
            
            # 캐시가 활성화되지 않은 경우
            if hasattr(self, '_cache_enabled') and not self._cache_enabled:
                return func(self, *args, **kwargs)
            
            # 캐시 키 생성
            if key_generator:
                cache_key = key_generator(func.__name__, *args, **kwargs)
            else:
                cache_key = generate_cache_key(func.__name__, *args, **kwargs)
            
            # 캐시에서 조회
            cached_value = self._cache.get(cache_key)
            if cached_value is not None:
                return cached_value
            
            # API 호출
            result = func(self, *args, **kwargs)
            
            # 캐시 조건 확인
            should_cache = True
            if cache_condition:
                should_cache = cache_condition(result)
            
            # 캐시 저장
            if should_cache:
                # TTL 결정
                method_ttl = ttl
                if method_ttl is None:
                    method_ttl = CACHE_TTL_CONFIG.get(func.__name__, 
                                                      self._cache._default_ttl)
                
                # 동적 TTL 적용
                if use_dynamic_ttl:
                    # market 파라미터 추출
                    market = 'KR'  # 기본값
                    if len(args) > 1 and isinstance(args[1], str):
                        # market이 두 번째 인자로 전달되는 경우
                        market = args[1]
                    elif 'market' in kwargs:
                        market = kwargs['market']
                    
                    # 동적 TTL 계산
                    from .market_hours import get_dynamic_ttl
                    method_ttl = get_dynamic_ttl(func.__name__, market)
                
                self._cache.set(cache_key, result, method_ttl)
            
            return result
        
        return wrapper
    return decorator 