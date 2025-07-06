#!/usr/bin/env python3
"""
Dynamic Batch Controller 테스트
Date: 2024-12-28
Issue: #27 - Phase 4.2 동적 배치 조정 테스트

동적 배치 조정 기능을 다양한 시나리오에서 테스트합니다.
"""

import os
import sys
import time
import logging
from unittest.mock import MagicMock, patch
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from korea_investment_stock.koreainvestmentstock import KoreaInvestment
from korea_investment_stock.dynamic_batch_controller import DynamicBatchController

# 로깅 설정
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)

def create_test_stock_list(count):
    """테스트용 주식 리스트 생성"""
    # 실제 한국 주식 종목 코드 예시
    sample_stocks = [
        ("005930", "KR"),  # 삼성전자
        ("000660", "KR"),  # SK하이닉스
        ("035420", "KR"),  # NAVER
        ("035720", "KR"),  # 카카오
        ("051910", "KR"),  # LG화학
        ("005380", "KR"),  # 현대차
        ("000270", "KR"),  # 기아
        ("006400", "KR"),  # 삼성SDI
        ("068270", "KR"),  # 셀트리온
        ("105560", "KR"),  # KB금융
    ]
    
    # count만큼 반복해서 리스트 생성
    result = []
    for i in range(count):
        result.append(sample_stocks[i % len(sample_stocks)])
    return result


def test_scenario_1_stable_requests():
    """시나리오 1: 안정적인 요청 (에러율 낮음)"""
    print("\n" + "="*60)
    print("시나리오 1: 안정적인 요청 테스트")
    print("="*60)
    
    # Mock KoreaInvestment 생성
    mock_ki = MagicMock(spec=KoreaInvestment)
    mock_ki.rate_limiter = MagicMock()
    mock_ki.executor = MagicMock()
    
    # 실제 메서드 바인딩
    from korea_investment_stock.koreainvestmentstock import KoreaInvestment
    mock_ki._KoreaInvestment__execute_concurrent_requests = KoreaInvestment.__execute_concurrent_requests.__get__(mock_ki)
    mock_ki.concurrent_limit = MagicMock()
    
    # Dynamic Batch Controller 생성
    controller = DynamicBatchController(
        initial_batch_size=20,
        initial_batch_delay=0.5,
        target_error_rate=0.02  # 2% 목표
    )
    
    # 테스트 데이터
    stock_list = create_test_stock_list(100)
    
    # Mock method - 대부분 성공
    success_count = 0
    def mock_method(symbol, market):
        nonlocal success_count
        time.sleep(0.01)  # 짧은 지연
        success_count += 1
        if success_count % 50 == 0:  # 2% 에러율
            raise Exception("Simulated error")
        return {"symbol": symbol, "market": market, "price": 50000}
    
    # 테스트 실행
    with patch('concurrent.futures.ThreadPoolExecutor') as mock_executor_class:
        mock_executor = MagicMock()
        mock_executor_class.return_value = mock_executor
        
        # Future 객체들 생성
        from concurrent.futures import Future
        futures = []
        for symbol, market in stock_list:
            future = Future()
            try:
                result = mock_method(symbol, market)
                future.set_result(result)
            except Exception as e:
                future.set_exception(e)
            futures.append(future)
        
        mock_executor.submit.side_effect = lambda fn, *args: futures.pop(0)
        
        # 실행
        results = mock_ki._KoreaInvestment__execute_concurrent_requests(
            mock_method,
            stock_list,
            dynamic_batch_controller=controller
        )
    
    # 결과 확인
    print(f"\n처리 결과:")
    print(f"- 총 항목: {len(stock_list)}")
    print(f"- 성공: {sum(1 for r in results if not isinstance(r, dict) or not r.get('error'))}")
    print(f"- 실패: {sum(1 for r in results if isinstance(r, dict) and r.get('error'))}")
    
    # 컨트롤러 통계
    stats = controller.get_stats()
    print(f"\n동적 배치 조정 결과:")
    print(f"- 초기 배치 크기: 20 → 최종: {stats['current_batch_size']}")
    print(f"- 초기 대기 시간: 0.5s → 최종: {stats['current_batch_delay']:.1f}s")
    print(f"- 조정 횟수: {stats['adjustment_count']}")


def test_scenario_2_high_error_rate():
    """시나리오 2: 높은 에러율 (서버 부하 시뮬레이션)"""
    print("\n" + "="*60)
    print("시나리오 2: 높은 에러율 테스트")
    print("="*60)
    
    controller = DynamicBatchController(
        initial_batch_size=50,
        initial_batch_delay=0.5,
        target_error_rate=0.01  # 1% 목표
    )
    
    # 시뮬레이션: 점진적으로 에러율 증가 후 회복
    test_phases = [
        (30, 0.05),   # 30개, 5% 에러
        (30, 0.20),   # 30개, 20% 에러 (높음)
        (30, 0.40),   # 30개, 40% 에러 (매우 높음)
        (30, 0.10),   # 30개, 10% 에러 (회복 중)
        (30, 0.02),   # 30개, 2% 에러 (안정화)
    ]
    
    for phase_idx, (count, error_rate) in enumerate(test_phases):
        print(f"\n단계 {phase_idx + 1}: {count}개 항목, {error_rate:.0%} 에러율")
        
        # 현재 파라미터
        batch_size, batch_delay = controller.get_current_parameters()
        print(f"현재 설정: batch_size={batch_size}, batch_delay={batch_delay:.1f}s")
        
        # 배치 시뮬레이션
        success = int(count * (1 - error_rate))
        error = count - success
        elapsed = count * 0.05  # 항목당 0.05초
        
        # 결과 기록
        controller.record_batch_result(
            batch_size=batch_size,
            success_count=success,
            error_count=error,
            elapsed_time=elapsed
        )
        
        time.sleep(0.5)  # 단계 간 지연
    
    # 최종 통계
    stats = controller.get_stats()
    print(f"\n최종 통계:")
    print(f"- 총 처리 항목: {stats['total_items']}")
    print(f"- 전체 에러율: {stats['overall_error_rate']:.1%}")
    print(f"- 파라미터 조정 횟수: {stats['adjustment_count']}")
    print(f"- 최종 배치 크기: {stats['current_batch_size']}")
    print(f"- 최종 대기 시간: {stats['current_batch_delay']:.1f}s")


def test_scenario_3_integration_test():
    """시나리오 3: 실제 KoreaInvestment와 통합 테스트"""
    print("\n" + "="*60)
    print("시나리오 3: KoreaInvestment 통합 테스트")
    print("="*60)
    
    # 환경 변수에서 인증 정보 가져오기
    api_key = os.environ.get('KI_API_KEY', 'test_key')
    api_secret = os.environ.get('KI_API_SECRET', 'test_secret')
    acc_no = os.environ.get('KI_ACC_NO', '00000000-00')
    
    if api_key == 'test_key':
        print("⚠️ 실제 API 키가 설정되지 않았습니다. Mock 모드로 실행합니다.")
        # Mock 모드로 실행
        return test_scenario_1_stable_requests()
    
    # 실제 KoreaInvestment 인스턴스 생성
    ki = KoreaInvestment(
        api_key=api_key,
        api_secret=api_secret,
        acc_no=acc_no,
        mock=True  # 모의투자 서버 사용
    )
    
    # Dynamic Batch Controller 생성
    controller = DynamicBatchController(
        initial_batch_size=10,  # 보수적으로 시작
        initial_batch_delay=1.0,
        target_error_rate=0.01,
        min_batch_size=5,
        max_batch_size=30
    )
    
    # 테스트할 종목 리스트 (30개)
    stock_list = create_test_stock_list(30)
    
    try:
        # 가격 조회 실행
        print(f"\n{len(stock_list)}개 종목 가격 조회 시작...")
        start_time = time.time()
        
        results = ki.fetch_price_list_with_dynamic_batch(
            stock_list,
            dynamic_batch_controller=controller
        )
        
        elapsed_time = time.time() - start_time
        
        # 결과 분석
        success_count = sum(1 for r in results if r.get('rt_cd') == '0')
        error_count = len(results) - success_count
        
        print(f"\n처리 완료:")
        print(f"- 소요 시간: {elapsed_time:.2f}초")
        print(f"- 성공: {success_count}")
        print(f"- 실패: {error_count}")
        print(f"- TPS: {len(results) / elapsed_time:.1f}")
        
        # 동적 배치 조정 통계
        stats = controller.get_stats()
        print(f"\n동적 배치 조정 통계:")
        print(f"- 초기 → 최종 배치 크기: 10 → {stats['current_batch_size']}")
        print(f"- 초기 → 최종 대기 시간: 1.0s → {stats['current_batch_delay']:.1f}s")
        print(f"- 조정 횟수: {stats['adjustment_count']}")
        
    except Exception as e:
        print(f"❌ 에러 발생: {e}")
    finally:
        ki.shutdown()


def main():
    """메인 테스트 실행"""
    print("🧪 Dynamic Batch Controller 테스트 시작")
    print("Phase 4.2: 에러율 기반 동적 배치 조정")
    
    # 시나리오 1: 안정적인 요청
    test_scenario_1_stable_requests()
    
    # 시나리오 2: 높은 에러율
    test_scenario_2_high_error_rate()
    
    # 시나리오 3: 통합 테스트 (실제 API 키가 있을 때만)
    # test_scenario_3_integration_test()
    
    print("\n✅ 모든 테스트 완료!")


if __name__ == "__main__":
    main() 