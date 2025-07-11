# 개선된 ThreadPoolExecutor 패턴 예제

import threading
from concurrent.futures import ThreadPoolExecutor, as_completed
import atexit
from typing import List, Tuple, Dict, Any, Callable

class ImprovedKoreaInvestment:
    """개선된 병렬 처리 패턴"""
    
    def __init__(self, api_key: str, api_secret: str, acc_no: str, mock: bool = False):
        # ... 기존 초기화 코드 ...
        
        # 개선점 1: 더 보수적인 워커 수
        self.executor = ThreadPoolExecutor(max_workers=3)
        
        # 개선점 2: 동시 실행 제한 세마포어
        self.api_semaphore = threading.Semaphore(2)  # 동시에 2개만 API 호출
        
        # 개선점 3: 자동 정리 등록
        atexit.register(self.shutdown)
        
        # Rate Limiter
        self.rate_limiter = RateLimiter(12, 1, safety_margin=0.9)
    
    # 개선점 4: 컨텍스트 매니저 지원
    def __enter__(self):
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        self.shutdown()
        return False
    
    def shutdown(self):
        """안전한 종료"""
        if hasattr(self, 'executor') and self.executor:
            self.executor.shutdown(wait=True)
            self.executor = None
    
    def __execute_concurrent_requests(self, method: Callable, stock_list: List[Tuple[str, str]]) -> List[Dict[str, Any]]:
        """개선된 병렬 요청 실행"""
        
        # 개선점 5: 래퍼 함수로 세마포어 적용
        def rate_limited_method(symbol: str, market: str) -> Dict[str, Any]:
            # 동시 실행 제한
            with self.api_semaphore:
                # Rate limiting
                self.rate_limiter.acquire()
                
                # 재시도 로직 포함
                for retry in range(3):
                    try:
                        result = method(symbol, market)
                        
                        # EGW00201 에러 체크
                        if result.get('msg_cd') == 'EGW00201':
                            self.__handle_rate_limit_error(retry)
                            continue
                        
                        return result
                    
                    except Exception as e:
                        if retry == 2:  # 마지막 시도
                            raise
                        print(f"Retry {retry + 1} for {symbol}: {e}")
                        time.sleep(0.5 * (retry + 1))
                
                return result
        
        # 개선점 6: 효율적인 작업 제출 및 결과 수집
        futures_to_stock = {}
        results = []
        
        # 모든 작업을 먼저 제출
        for symbol, market in stock_list:
            future = self.executor.submit(rate_limited_method, symbol, market)
            futures_to_stock[future] = (symbol, market)
        
        # 개선점 7: as_completed로 완료되는 대로 처리
        for future in as_completed(futures_to_stock):
            symbol, market = futures_to_stock[future]
            
            try:
                result = future.result(timeout=30)  # 타임아웃 설정
                results.append(result)
                
            except Exception as e:
                # 개선점 8: 에러도 결과에 포함
                error_result = {
                    'error': str(e),
                    'error_type': type(e).__name__,
                    'symbol': symbol,
                    'market': market,
                    'rt_cd': 'ERROR'
                }
                results.append(error_result)
                print(f"Error processing {symbol}: {e}")
        
        # 통계 출력
        self.rate_limiter.print_stats()
        
        # 결과 정렬 (원래 순서대로)
        # 순서가 중요한 경우 인덱스를 추가로 관리
        return results
    
    # 개선점 9: 배치 처리 개선 (선택사항)
    def __execute_concurrent_requests_batched(self, method: Callable, stock_list: List[Tuple[str, str]], 
                                            batch_size: int = 3) -> List[Dict[str, Any]]:
        """배치 기반 병렬 처리 (더 안정적)"""
        all_results = []
        
        for i in range(0, len(stock_list), batch_size):
            batch = stock_list[i:i + batch_size]
            
            # 배치 내 순차적 제출로 초기 버스트 방지
            futures = []
            for idx, (symbol, market) in enumerate(batch):
                if idx > 0:
                    time.sleep(0.05)  # 50ms 간격으로 제출
                
                future = self.executor.submit(
                    self._rate_limited_api_call, 
                    method, symbol, market
                )
                futures.append((future, symbol, market))
            
            # 배치 결과 수집
            batch_results = []
            for future, symbol, market in futures:
                try:
                    result = future.result(timeout=30)
                    batch_results.append(result)
                except Exception as e:
                    batch_results.append({
                        'error': str(e),
                        'symbol': symbol,
                        'market': market
                    })
            
            all_results.extend(batch_results)
            
            # 다음 배치 전 대기
            if i + batch_size < len(stock_list):
                time.sleep(0.2)
        
        return all_results
    
    def _rate_limited_api_call(self, method: Callable, symbol: str, market: str) -> Dict[str, Any]:
        """Rate limiting이 적용된 API 호출"""
        with self.api_semaphore:
            self.rate_limiter.acquire()
            return method(symbol, market)


# 사용 예시
if __name__ == "__main__":
    # 컨텍스트 매니저로 자동 정리
    with ImprovedKoreaInvestment(api_key, api_secret, acc_no) as broker:
        stocks = [
            ("005930", "KR"),  # 삼성전자
            ("000660", "KR"),  # SK하이닉스
            ("035720", "KR"),  # 카카오
            # ... 더 많은 종목
        ]
        
        # 병렬 처리
        results = broker.fetch_price_list(stocks)
        
        # 에러와 성공 분리
        successes = [r for r in results if 'error' not in r]
        errors = [r for r in results if 'error' in r]
        
        print(f"성공: {len(successes)}, 실패: {len(errors)}")
        
        # 에러 상세 확인
        for error in errors:
            print(f"Error for {error['symbol']}: {error['error']}")
    
    # 여기서 자동으로 executor.shutdown() 호출됨 