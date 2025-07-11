#!/usr/bin/env python3
"""
Integration Tests for Rate Limiting System
Date: 2024-12-28

Phase 6.2: 통합 테스트 작성
- Mock 서버를 이용한 Rate Limit 시나리오 테스트
- 100개 종목 동시 조회 테스트
- 장시간 실행 안정성 테스트
"""

import unittest
import time
import json
import threading
from unittest.mock import Mock, patch, MagicMock
from concurrent.futures import ThreadPoolExecutor, as_completed
from datetime import datetime, timedelta
import random

from .. import KoreaInvestment
from ..rate_limiting import EnhancedRateLimiter, get_backoff_strategy
from ..error_handling import get_error_recovery_system


class MockAPIServer:
    """Mock API 서버 - Rate Limit 시뮬레이션"""
    
    def __init__(self, rate_limit=20, window=1.0):
        self.rate_limit = rate_limit
        self.window = window
        self.call_times = []
        self.lock = threading.Lock()
        self.error_count = 0
        self.total_calls = 0
        
    def check_rate_limit(self):
        """Rate limit 체크"""
        with self.lock:
            now = time.time()
            # 윈도우 밖의 호출 제거
            self.call_times = [t for t in self.call_times if now - t < self.window]
            
            if len(self.call_times) >= self.rate_limit:
                self.error_count += 1
                return False
            
            self.call_times.append(now)
            self.total_calls += 1
            return True
    
    def api_call(self, delay=0.01):
        """모의 API 호출"""
        if not self.check_rate_limit():
            return {
                "rt_cd": "EGW00201",
                "msg1": "초당 거래건수를 초과하였습니다."
            }
        
        # 정상 응답 시뮬레이션
        time.sleep(delay)  # API 처리 시간
        return {
            "rt_cd": "0",
            "msg1": "정상처리 되었습니다.",
            "output": {"price": random.randint(10000, 100000)}
        }
    
    def get_stats(self):
        """통계 반환"""
        return {
            "total_calls": self.total_calls,
            "error_count": self.error_count,
            "error_rate": self.error_count / self.total_calls if self.total_calls > 0 else 0
        }


class TestRateLimitIntegration(unittest.TestCase):
    """Rate Limit 통합 테스트"""
    
    def setUp(self):
        """테스트 환경 설정"""
        # Mock KoreaInvestment 인스턴스
        self.kis = Mock(spec=KoreaInvestment)
        self.kis.APP_KEY = "test_app_key"
        self.kis.APP_SECRET = "test_app_secret"
        self.kis.TOKEN = "test_token"
        
        # EnhancedRateLimiter 설정
        self.kis.rate_limiter = EnhancedRateLimiter(
            max_calls=15,
            per_seconds=1.0,
            safety_margin=0.8
        )
        
        # Mock 서버
        self.mock_server = MockAPIServer(rate_limit=20)
        
    def test_mock_server_rate_limit(self):
        """Mock 서버 Rate Limit 시나리오 테스트"""
        print("\n=== Mock 서버 Rate Limit 시나리오 테스트 ===")
        
        success_count = 0
        error_count = 0
        
        # 30개 요청을 빠르게 전송
        for i in range(30):
            if self.kis.rate_limiter.acquire():
                response = self.mock_server.api_call()
                if response["rt_cd"] == "0":
                    success_count += 1
                else:
                    error_count += 1
                    print(f"Rate limit 에러 발생: {response['msg1']}")
            else:
                print(f"클라이언트 Rate Limiter가 요청 {i+1}을 차단")
        
        stats = self.mock_server.get_stats()
        print(f"\n서버 통계:")
        print(f"- 총 요청: {stats['total_calls']}")
        print(f"- 에러 수: {stats['error_count']}")
        print(f"- 에러율: {stats['error_rate']:.1%}")
        
        # 클라이언트 Rate Limiter로 인해 서버 에러가 없어야 함
        self.assertEqual(stats['error_count'], 0)
        print("✅ 클라이언트 Rate Limiter가 서버 Rate Limit 초과를 완벽히 방지")
    
    def test_100_symbols_concurrent(self):
        """100개 종목 동시 조회 테스트"""
        print("\n=== 100개 종목 동시 조회 테스트 ===")
        
        # 100개 종목 코드 생성
        symbols = [f"00{i:04d}" for i in range(100)]
        
        def fetch_price(symbol):
            """종목 가격 조회 시뮬레이션"""
            if self.kis.rate_limiter.acquire():
                response = self.mock_server.api_call(delay=0.02)
                return symbol, response
            else:
                return symbol, {"rt_cd": "BLOCKED", "msg1": "Rate limiter blocked"}
        
        # ThreadPoolExecutor로 병렬 처리
        start_time = time.time()
        results = {}
        
        with ThreadPoolExecutor(max_workers=3) as executor:
            futures = {executor.submit(fetch_price, symbol): symbol 
                      for symbol in symbols}
            
            for future in as_completed(futures):
                symbol = futures[future]
                try:
                    symbol, response = future.result()
                    results[symbol] = response
                except Exception as e:
                    results[symbol] = {"error": str(e)}
        
        elapsed_time = time.time() - start_time
        
        # 결과 분석
        success_count = sum(1 for r in results.values() if r.get("rt_cd") == "0")
        error_count = sum(1 for r in results.values() if r.get("rt_cd") == "EGW00201")
        blocked_count = sum(1 for r in results.values() if r.get("rt_cd") == "BLOCKED")
        
        print(f"\n실행 시간: {elapsed_time:.2f}초")
        print(f"성공: {success_count}, 서버 에러: {error_count}, 클라이언트 차단: {blocked_count}")
        print(f"초당 처리량: {success_count / elapsed_time:.1f} TPS")
        
        # 서버 에러가 없어야 함
        self.assertEqual(error_count, 0)
        # 최소 80개 이상 성공
        self.assertGreaterEqual(success_count, 80)
        print("✅ 100개 종목 조회 중 Rate Limit 에러 없이 안정적 처리")
    
    def test_long_running_stability(self):
        """장시간 실행 안정성 테스트"""
        print("\n=== 장시간 실행 안정성 테스트 ===")
        
        # 30초 동안 연속 호출
        test_duration = 30  # 초
        start_time = time.time()
        call_count = 0
        error_count = 0
        
        print("30초 동안 연속 호출 테스트 시작...")
        
        while time.time() - start_time < test_duration:
            if self.kis.rate_limiter.acquire():
                response = self.mock_server.api_call()
                call_count += 1
                
                if response["rt_cd"] == "EGW00201":
                    error_count += 1
                
                # 진행 상황 출력 (5초마다)
                if call_count % 50 == 0:
                    elapsed = time.time() - start_time
                    print(f"  {elapsed:.1f}초 경과: {call_count}개 호출, {error_count}개 에러")
            
            time.sleep(0.01)  # CPU 과부하 방지
        
        total_time = time.time() - start_time
        server_stats = self.mock_server.get_stats()
        
        print(f"\n최종 결과:")
        print(f"- 실행 시간: {total_time:.1f}초")
        print(f"- 총 호출 수: {call_count}")
        print(f"- 평균 TPS: {call_count / total_time:.1f}")
        print(f"- 서버 에러: {server_stats['error_count']}")
        
        # 안정성 검증
        self.assertEqual(server_stats['error_count'], 0)
        self.assertGreater(call_count, 300)  # 최소 300개 이상 호출
        avg_tps = call_count / total_time
        self.assertLessEqual(avg_tps, 13)  # 안전 마진 내
        self.assertGreaterEqual(avg_tps, 10)  # 적절한 처리량
        
        print("✅ 30초 동안 안정적으로 실행, Rate Limit 에러 없음")


class TestErrorRecoveryIntegration(unittest.TestCase):
    """에러 복구 통합 테스트"""
    
    def setUp(self):
        """테스트 환경 설정"""
        self.recovery_system = get_error_recovery_system()
        self.backoff_strategy = get_backoff_strategy()
        
        # Mock API 함수
        self.call_count = 0
        self.error_sequence = []
        
    def mock_api_call(self):
        """에러 시퀀스를 따르는 Mock API"""
        self.call_count += 1
        
        if self.call_count <= len(self.error_sequence):
            error_type = self.error_sequence[self.call_count - 1]
            if error_type == "rate_limit":
                return {"rt_cd": "EGW00201", "msg1": "Rate limit exceeded"}
            elif error_type == "server_error":
                raise Exception("Server error")
            elif error_type == "network":
                raise ConnectionError("Network error")
        
        return {"rt_cd": "0", "msg1": "Success", "data": "OK"}
    
    def test_error_recovery_flow(self):
        """에러 복구 플로우 통합 테스트"""
        print("\n=== 에러 복구 플로우 통합 테스트 ===")
        
        # 다양한 에러 시나리오
        test_scenarios = [
            {
                "name": "Rate Limit 후 복구",
                "errors": ["rate_limit", "rate_limit"],
                "expected_calls": 3
            },
            {
                "name": "서버 에러 후 복구",
                "errors": ["server_error"],
                "expected_calls": 2
            },
            {
                "name": "복합 에러 시나리오",
                "errors": ["network", "rate_limit", "server_error"],
                "expected_calls": 4
            }
        ]
        
        for scenario in test_scenarios:
            print(f"\n테스트: {scenario['name']}")
            
            # 초기화
            self.call_count = 0
            self.error_sequence = scenario["errors"]
            self.backoff_strategy.reset()
            
            # retry_on_rate_limit 데코레이터 시뮬레이션
            from .enhanced_retry_decorator import enhanced_retry
            
            @enhanced_retry(max_retries=5)
            def api_call_with_retry():
                return self.mock_api_call()
            
            # 실행
            try:
                result = api_call_with_retry()
                self.assertEqual(result["rt_cd"], "0")
                self.assertEqual(self.call_count, scenario["expected_calls"])
                print(f"✅ {self.call_count}번 시도 후 성공")
            except Exception as e:
                self.fail(f"예상치 못한 실패: {e}")
    
    def test_circuit_breaker_integration(self):
        """Circuit Breaker 통합 테스트"""
        print("\n=== Circuit Breaker 통합 테스트 ===")
        
        # 연속 실패로 Circuit Breaker 열기
        self.error_sequence = ["server_error"] * 10
        
        from .enhanced_retry_decorator import enhanced_retry
        
        @enhanced_retry(max_retries=3, enable_circuit_breaker=True)
        def failing_api_call():
            return self.mock_api_call()
        
        # 여러 번 호출하여 Circuit Breaker 상태 확인
        failures = 0
        for i in range(5):
            try:
                self.call_count = 0
                result = failing_api_call()
            except Exception as e:
                failures += 1
                if i >= 2:  # Circuit Breaker가 열려야 함
                    self.assertIn("Circuit", str(e))
                    print(f"✅ Circuit Breaker 활성화 확인: {e}")
        
        self.assertEqual(failures, 5)
        
        # Circuit Breaker 상태 확인
        state = self.backoff_strategy.state
        print(f"✅ Circuit Breaker 상태: {state}")


class TestRealWorldScenarios(unittest.TestCase):
    """실제 사용 시나리오 테스트"""
    
    @patch('korea_investment_stock.koreainvestmentstock.KoreaInvestment.__execute_concurrent_requests')
    def test_batch_processing_scenario(self, mock_execute):
        """배치 처리 시나리오 테스트"""
        print("\n=== 배치 처리 시나리오 테스트 ===")
        
        # Mock 설정
        mock_execute.return_value = [
            {"rt_cd": "0", "output": {"price": i * 1000}} 
            for i in range(10)
        ]
        
        # KoreaInvestment 인스턴스 생성
        kis = KoreaInvestment(
            api_key="test_key",
            api_secret="test_secret",
            is_mock=True
        )
        
        # 50개 종목을 5개씩 배치로 처리
        symbols = [f"00{i:04d}" for i in range(50)]
        batch_size = 5
        
        start_time = time.time()
        all_results = []
        
        for i in range(0, len(symbols), batch_size):
            batch = symbols[i:i+batch_size]
            print(f"배치 {i//batch_size + 1}: {len(batch)}개 처리")
            
            # 배치 처리 시뮬레이션
            results = []
            for symbol in batch:
                if kis.rate_limiter.acquire():
                    results.append({"symbol": symbol, "price": random.randint(10000, 100000)})
            
            all_results.extend(results)
            
            # 배치 간 대기
            if i + batch_size < len(symbols):
                time.sleep(0.5)
        
        elapsed = time.time() - start_time
        
        print(f"\n결과:")
        print(f"- 처리 시간: {elapsed:.2f}초")
        print(f"- 처리된 종목: {len(all_results)}개")
        print(f"- 평균 처리 속도: {len(all_results) / elapsed:.1f} 종목/초")
        
        self.assertEqual(len(all_results), 50)
        self.assertLess(elapsed, 10)  # 10초 이내 완료
        print("✅ 배치 처리로 효율적인 대량 조회 완료")


def run_integration_tests():
    """통합 테스트 실행"""
    print("\n" + "="*60)
    print("Rate Limiting 시스템 통합 테스트 시작")
    print("="*60)
    
    # 테스트 스위트 생성
    suite = unittest.TestSuite()
    
    # Rate Limit 통합 테스트
    suite.addTests(unittest.TestLoader().loadTestsFromTestCase(TestRateLimitIntegration))
    
    # 에러 복구 통합 테스트
    suite.addTests(unittest.TestLoader().loadTestsFromTestCase(TestErrorRecoveryIntegration))
    
    # 실제 시나리오 테스트
    suite.addTests(unittest.TestLoader().loadTestsFromTestCase(TestRealWorldScenarios))
    
    # 테스트 실행
    runner = unittest.TextTestRunner(verbosity=2)
    result = runner.run(suite)
    
    print("\n" + "="*60)
    print("통합 테스트 완료")
    print(f"실행: {result.testsRun}, 성공: {result.testsRun - len(result.failures) - len(result.errors)}")
    print("="*60)
    
    return result


if __name__ == "__main__":
    run_integration_tests() 