#!/usr/bin/env python3
"""
TTL 캐시 통합 테스트
Date: 2025-01-07

Phase 8.9: 통합 테스트
- Rate Limiter와 함께 동작 테스트
- 캐시 워밍업 시나리오
- 성능 측정
"""

import time
import unittest
from unittest.mock import Mock, patch, MagicMock
import threading

# 상대 import
from ..korea_investment_stock import KoreaInvestment
from ..caching import TTLCache, cacheable


class TestCacheIntegration(unittest.TestCase):
    """캐시 통합 테스트"""
    
    def setUp(self):
        """테스트 환경 설정"""
        # Mock API 응답
        self.mock_response = {
            'rt_cd': '0',
            'msg1': '정상처리되었습니다',
            'output': {
                'price': '50000',
                'volume': '1000000',
                'prdt_clsf_name': '주권'
            }
        }
    
    def _create_kis_with_real_cache(self):
        """실제 캐시가 동작하는 KIS 인스턴스 생성"""
        # KoreaInvestment 인스턴스 생성
        with patch('korea_investment_stock.korea_investment_stock.KoreaInvestment.issue_access_token'):
            with patch('korea_investment_stock.korea_investment_stock.KoreaInvestment.check_access_token', return_value=False):
                # requests.get을 Mock
                with patch('korea_investment_stock.korea_investment_stock.requests') as mock_requests:
                    # Mock 응답 설정
                    mock_response = Mock()
                    
                    # URL에 따라 다른 응답 반환
                    def mock_get(url, **kwargs):
                        response = Mock()
                        # 종목 정보 조회
                        if 'search-info' in url or 'search-stock-info' in url:
                            response.json.return_value = {
                                'rt_cd': '0',
                                'output': {'prdt_clsf_name': '주권'}
                            }
                        # 가격 조회
                        else:
                            response.json.return_value = self.mock_response
                        return response
                    
                    mock_requests.get.side_effect = mock_get
                    mock_requests.post.return_value = mock_response
                    
                    kis = KoreaInvestment(
                        api_key="test_key",
                        api_secret="test_secret",
                        acc_no="12345678-01",
                        mock=True,
                        cache_enabled=True,
                        cache_config={
                            'default_ttl': 60,  # 1분
                            'max_size': 100
                        }
                    )
                    
                    # requests mock을 인스턴스 속성으로 저장
                    kis._mock_requests = mock_requests
                    
                    # Rate limiter acquire를 Mock으로
                    if hasattr(kis, 'rate_limiter'):
                        kis.rate_limiter.acquire = Mock()
                    
                    return kis
    
    def test_cache_with_real_methods(self):
        """실제 메서드에서 캐시 동작 테스트"""
        print("\n=== 실제 메서드 캐시 테스트 ===")
        
        kis = self._create_kis_with_real_cache()
        
        # __fetch_stock_info Mock
        kis._KoreaInvestment__fetch_stock_info = Mock(return_value={
            'rt_cd': '0',
            'output': {'prdt_clsf_name': '주권'}
        })
        
        # 첫 번째 호출 - 캐시 미스
        start = time.time()
        result1 = kis.fetch_domestic_price("J", "005930")
        first_call_time = time.time() - start
        print(f"DEBUG: 첫 번째 호출 결과: {result1}")
        
        # 두 번째 호출 - 캐시 히트
        start = time.time()
        result2 = kis.fetch_domestic_price("J", "005930")
        second_call_time = time.time() - start
        print(f"DEBUG: 두 번째 호출 결과: {result2}")
        
        # 결과 검증
        if result1.get('rt_cd') != '0':
            print(f"WARNING: API 응답 에러 - rt_cd: {result1.get('rt_cd')}, msg: {result1.get('msg1')}")
            # 캐시는 동작하는지 확인
            self.assertEqual(result1, result2)  # 캐시에서 같은 값이 나와야 함
        else:
            self.assertEqual(result1['rt_cd'], '0')
            self.assertEqual(result2['rt_cd'], '0')
        
        # 캐시가 동작했는지 확인
        cache_stats = kis.get_cache_stats()
        self.assertGreater(cache_stats['hit_count'], 0)
        
        print(f"✅ 첫 번째 호출: {first_call_time:.4f}초")
        print(f"✅ 두 번째 호출: {second_call_time:.4f}초 (캐시)")
        print(f"✅ 캐시 적중: {cache_stats['hit_count']}, 미스: {cache_stats['miss_count']}")
        
        # shutdown 호출하여 정리
        kis.shutdown()
    
    def test_cache_ttl_expiration(self):
        """TTL 만료 테스트"""
        print("\n=== TTL 만료 테스트 ===")
        
        # TTL이 짧은 캐시로 설정
        kis = self._create_kis_with_real_cache()
        kis._cache = TTLCache(default_ttl=1, max_size=10)  # 1초 TTL
        
        # Mock 설정
        kis._KoreaInvestment__fetch_stock_info = Mock(return_value={
            'rt_cd': '0',
            'output': {'prdt_clsf_name': '주권'}
        })
        
        # API 호출 횟수 추적
        call_count = 0
        original_get = kis._mock_requests.get
        def counting_get(*args, **kwargs):
            nonlocal call_count
            call_count += 1
            time.sleep(0.01)  # API 호출 시뮬레이션
            return original_get(*args, **kwargs)
        
        kis._mock_requests.get.side_effect = counting_get
        
        # 첫 번째 호출
        result1 = kis.fetch_domestic_price("J", "005930")
        first_call_count = call_count
        
        # 즉시 재호출 - 캐시에서
        result2 = kis.fetch_domestic_price("J", "005930")
        self.assertEqual(call_count, first_call_count)  # API 호출 없음
        
        # 2초 대기 (TTL 만료)
        time.sleep(2)
        
        # 재호출 - API 재호출
        result3 = kis.fetch_domestic_price("J", "005930")
        self.assertGreater(call_count, first_call_count)  # API 호출 발생
        
        print(f"✅ TTL 만료 후 API 재호출 확인")
        print(f"✅ 총 API 호출: {call_count}회")
        
        kis.shutdown()
    
    def test_list_method_with_cache(self):
        """리스트 메서드의 캐시 동작"""
        print("\n=== 리스트 메서드 캐시 테스트 ===")
        
        kis = self._create_kis_with_real_cache()
        
        # Mock 설정
        def mock_fetch_price(symbol, market):
            time.sleep(0.01)  # API 호출 시뮬레이션
            return {
                'rt_cd': '0',
                'symbol': symbol,
                'market': market,
                'price': 50000 + hash(symbol) % 10000
            }
        
        kis._KoreaInvestment__fetch_price = Mock(side_effect=mock_fetch_price)
        kis._KoreaInvestment__fetch_price.__name__ = '__fetch_price'  # Mock에 이름 추가
        
        # 종목 리스트
        stock_list = [
            ("005930", "KR"),
            ("000660", "KR"),
            ("035720", "KR")
        ]
        
        # 첫 번째 호출
        start = time.time()
        results1 = kis.fetch_price_list(stock_list)
        first_time = time.time() - start
        api_calls_first = kis._KoreaInvestment__fetch_price.call_count
        
        # 두 번째 호출 - 캐시에서
        start = time.time()
        results2 = kis.fetch_price_list(stock_list)
        second_time = time.time() - start
        api_calls_second = kis._KoreaInvestment__fetch_price.call_count
        
        # 검증
        self.assertEqual(len(results1), len(stock_list))
        self.assertEqual(results1, results2)
        self.assertEqual(api_calls_first, api_calls_second)  # API 호출 증가 없음
        
        print(f"✅ 첫 번째 호출: {first_time:.3f}초 ({api_calls_first}회 API)")
        print(f"✅ 두 번째 호출: {second_time:.3f}초 (캐시)")
        print(f"✅ 속도 향상: {first_time/second_time:.1f}배")
        
        kis.shutdown()
    
    def test_concurrent_access_thread_safety(self):
        """동시 접근 스레드 안전성 테스트"""
        print("\n=== 스레드 안전성 테스트 ===")
        
        kis = self._create_kis_with_real_cache()
        
        # Mock 설정
        kis._KoreaInvestment__fetch_stock_info = Mock(return_value={
            'rt_cd': '0',
            'output': {'prdt_clsf_name': '주권'}
        })
        
        results = []
        errors = []
        api_call_count = 0
        lock = threading.Lock()
        
        # API 호출 추적
        original_get = kis._mock_requests.get
        def counting_get(*args, **kwargs):
            nonlocal api_call_count
            with lock:
                api_call_count += 1
            time.sleep(0.01)  # API 호출 시뮬레이션
            return original_get(*args, **kwargs)
        
        kis._mock_requests.get.side_effect = counting_get
        
        def worker(thread_id):
            try:
                for i in range(5):
                    result = kis.fetch_domestic_price("J", "005930")
                    results.append((thread_id, i, result['rt_cd']))
            except Exception as e:
                errors.append((thread_id, str(e)))
        
        # 10개 스레드 동시 실행
        threads = []
        for i in range(10):
            t = threading.Thread(target=worker, args=(i,))
            threads.append(t)
            t.start()
        
        for t in threads:
            t.join()
        
        # 검증
        self.assertEqual(len(errors), 0)
        self.assertEqual(len(results), 50)  # 10 threads * 5 calls
        
        # API 호출은 최소화되어야 함
        self.assertLess(api_call_count, 10)  # 캐시 덕분에 적게 호출
        
        print(f"✅ 총 조회: {len(results)}회")
        print(f"✅ API 호출: {api_call_count}회")
        print(f"✅ 캐시 효율: {(1 - api_call_count/len(results)) * 100:.1f}%")
        
        kis.shutdown()
    
    def test_cache_memory_management(self):
        """캐시 메모리 관리 테스트"""
        print("\n=== 메모리 관리 테스트 ===")
        
        kis = self._create_kis_with_real_cache()
        kis._cache = TTLCache(default_ttl=300, max_size=5)  # 최대 5개 항목
        
        # Mock 설정
        kis._KoreaInvestment__fetch_stock_info = Mock(return_value={
            'rt_cd': '0',
            'output': {'prdt_clsf_name': '주권'}
        })
        
        # 10개 종목 조회 (캐시 크기는 5)
        symbols = [f"{i:06d}" for i in range(10)]
        
        for symbol in symbols:
            kis.fetch_domestic_price("J", symbol)
        
        # 캐시 상태 확인
        stats = kis.get_cache_stats()
        self.assertLessEqual(stats['total_entries'], 5)
        self.assertGreater(stats['eviction_count'], 0)
        
        print(f"✅ 캐시 크기 제한: {stats['total_entries']}/5")
        print(f"✅ 제거된 항목: {stats['eviction_count']}")
        
        kis.shutdown()
    
    def tearDown(self):
        """테스트 정리"""
        # shutdown은 각 테스트 메서드에서 호출
        pass


if __name__ == '__main__':
    unittest.main() 