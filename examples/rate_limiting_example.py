#!/usr/bin/env python3
"""
Rate Limiting 예제
한국투자증권 API의 속도 제한을 자동으로 관리하는 방법을 보여줍니다.
"""

import korea_investment_stock
import time
import datetime


def basic_rate_limiting_example(broker):
    """기본 Rate Limiting 예제"""
    print("=== 기본 Rate Limiting 예제 ===\n")
    
    # 여러 종목 조회 - Rate Limiting이 자동으로 적용됨
    symbols = ["005930", "000660", "035720", "051910", "005380"]  # 삼성전자, SK하이닉스, 카카오, LG화학, 현대차
    stock_list = [(symbol, "KR") for symbol in symbols]  # (symbol, market) 튜플 리스트로 변환
    
    print(f"조회할 종목 수: {len(symbols)}")
    print("Rate Limiting이 자동으로 적용되어 안전하게 처리됩니다.\n")
    
    start_time = time.time()
    
    try:
        results = broker.fetch_price_list(stock_list)
        for i, result in enumerate(results):
            if result.get('rt_cd') == '0':
                if 'output' in result:
                    symbol = stock_list[i][0]
                    price = result['output'].get('stck_prpr', 'N/A')
                    print(f"{symbol}: {price}원")
                else:
                    print(f"{stock_list[i][0]}: 데이터 없음")
            else:
                print(f"{stock_list[i][0]}: 에러 - {result.get('msg1', '알 수 없는 에러')}")
    except Exception as e:
        print(f"전체 조회 중 에러 발생: {e}")
    
    elapsed = time.time() - start_time
    print(f"\n처리 시간: {elapsed:.2f}초")
    print(f"평균 TPS: {len(symbols) / elapsed:.2f}")


def statistics_example(broker):
    """통계 확인 및 저장 예제"""
    print("\n=== 통계 확인 예제 ===\n")
    
    # 통계 출력
    broker.rate_limiter.print_stats()
    
    # 통계 데이터 가져오기
    stats = broker.rate_limiter.get_stats()
    print(f"\n상세 통계:")
    print(f"- 총 호출 수: {stats['total_calls']}")
    print(f"- 에러 수: {stats['error_count']}")
    
    # error_rate 계산
    error_rate = stats['error_count'] / stats['total_calls'] if stats['total_calls'] > 0 else 0
    print(f"- 에러율: {error_rate:.1%}")
    
    print(f"- 최대 초당 호출: {stats['max_calls_per_second']}")
    print(f"- 평균 대기 시간: {stats['avg_wait_time']:.3f}초")
    
    # 통계 파일로 저장
    filepath = broker.rate_limiter.save_stats()
    print(f"\n통계가 저장되었습니다: {filepath}")


def batch_processing_example(broker):
    """배치 처리 예제"""
    print("\n=== 배치 처리 예제 ===\n")
    
    # 많은 종목을 조회할 때 배치 처리 사용
    # KOSPI 상위 30개 종목 (예시)
    stock_list = [
        ("005930", "KR"), ("000660", "KR"), ("005490", "KR"), ("005380", "KR"),
        ("012330", "KR"), ("051910", "KR"), ("035420", "KR"), ("000270", "KR"),
        ("068270", "KR"), ("028260", "KR"), ("105560", "KR"), ("055550", "KR"),
        ("035720", "KR"), ("032830", "KR"), ("003670", "KR"), ("015760", "KR"),
        ("017670", "KR"), ("090430", "KR"), ("009150", "KR"), ("000810", "KR"),
        ("011200", "KR"), ("086790", "KR"), ("033780", "KR"), ("006400", "KR"),
        ("021240", "KR"), ("051900", "KR"), ("034730", "KR"), ("003550", "KR"),
        ("018260", "KR"), ("010130", "KR")
    ]
    
    print(f"조회할 종목 수: {len(stock_list)}")
    
    # 방법 1: 기본 방식 (한 번에 모두 처리)
    print("\n1) 기본 방식으로 처리:")
    start_time = time.time()
    results1 = broker.fetch_price_list(stock_list[:10])  # 처음 10개만
    elapsed1 = time.time() - start_time
    print(f"   처리 시간: {elapsed1:.2f}초")
    print(f"   성공: {sum(1 for r in results1 if r.get('rt_cd') == '0')}/{len(results1)}")
    
    # 방법 2: 배치 처리 방식 (권장) - fetch_price_list_with_batch 사용
    print("\n2) 배치 처리 방식 (10개씩 나누어 처리):")
    start_time = time.time()
    results2 = broker.fetch_price_list_with_batch(
        stock_list,
        batch_size=10,      # 10개씩 처리
        batch_delay=1.0,    # 배치 간 1초 대기
        progress_interval=5  # 5개마다 진행상황 출력
    )
    elapsed2 = time.time() - start_time
    print(f"\n   전체 처리 시간: {elapsed2:.2f}초")
    print(f"   성공: {sum(1 for r in results2 if r.get('rt_cd') == '0')}/{len(results2)}")


def auto_save_example(broker):
    """자동 저장 예제"""
    print("\n=== 자동 저장 예제 ===\n")
    
    # 자동 저장 활성화 (30초마다)
    broker.rate_limiter.enable_auto_save(interval_seconds=30)
    print("자동 저장이 활성화되었습니다 (30초마다)")
    
    # 일정 시간 동안 작업 수행
    print("\n30초 동안 주기적으로 API 호출을 수행합니다...")
    symbols = ["005930", "000660", "035720"]
    
    start_time = time.time()
    while time.time() - start_time < 35:  # 35초 동안 실행
        for symbol in symbols:
            try:
                broker.fetch_price(symbol)
            except Exception:
                pass
        time.sleep(5)  # 5초마다 반복
        
        # 현재 상태 출력
        stats = broker.rate_limiter.get_stats()
        print(f"  {int(time.time() - start_time)}초 경과 - 총 호출: {stats['total_calls']}")
    
    print("\n자동 저장이 실행되었는지 확인:")
    print("logs/rate_limiter_stats/rate_limiter_stats_latest.json 파일을 확인하세요.")
    
    # 자동 저장 비활성화
    broker.rate_limiter.disable_auto_save()
    print("자동 저장이 비활성화되었습니다.")


def peak_hour_example(broker):
    """피크 시간대 처리 예제"""
    print("\n=== 피크 시간대 처리 예제 ===\n")
    
    current_hour = datetime.datetime.now().hour
    
    # 시간대별 설정
    if 9 <= current_hour <= 10 or 15 <= current_hour <= 16:
        # 장 시작/종료 시간대
        batch_size = 20
        batch_delay = 2.0
        print(f"피크 시간대 ({current_hour}시): 보수적 설정 사용")
    else:
        # 일반 시간대
        batch_size = 50
        batch_delay = 0.5
        print(f"일반 시간대 ({current_hour}시): 표준 설정 사용")
    
    print(f"- 배치 크기: {batch_size}")
    print(f"- 배치 간 대기: {batch_delay}초")
    
    # 예제 실행
    stock_list = [(f"{i:06d}", "KR") for i in range(5000, 5020)]
    
    results = broker.fetch_price_list_with_batch(
        stock_list,
        batch_size=batch_size,
        batch_delay=batch_delay
    )
    
    print(f"\n처리 완료: {len(results)}개 항목")


def error_monitoring_example(broker):
    """에러 모니터링 예제"""
    print("\n=== 에러 모니터링 예제 ===\n")
    
    # 일부러 에러를 발생시킬 잘못된 종목 코드 포함
    symbols = ["005930", "INVALID", "000660", "999999", "035720"]
    stock_list = [(symbol, "KR") for symbol in symbols]
    
    try:
        results = broker.fetch_price_list(stock_list)
        for i, result in enumerate(results):
            symbol = symbols[i]
            if result.get('rt_cd') == '0':
                print(f"✓ {symbol}: 성공")
            else:
                print(f"✗ {symbol}: 실패 - {result.get('msg1', '알 수 없는 에러')[:50]}...")
    except Exception as e:
        print(f"에러 모니터링 중 예외 발생: {e}")
    
    # 에러율 확인
    stats = broker.rate_limiter.get_stats()
    error_rate = stats['error_count'] / stats['total_calls'] if stats['total_calls'] > 0 else 0
    
    print(f"\n현재 에러율: {error_rate:.1%}")
    
    if error_rate > 0.01:  # 1% 초과
        print("⚠️  경고: 에러율이 1%를 초과했습니다!")
        print("조치 필요: API 호출 패턴을 검토하세요.")
    else:
        print("✅ 에러율이 정상 범위 내에 있습니다.")


def main():
    """메인 함수"""
    print("한국투자증권 Rate Limiting 예제")
    print("=" * 50)
    
    # API 키 설정 (실제 사용 시 본인의 키로 교체)
    key = "YOUR_API_KEY"
    secret = "YOUR_API_SECRET"
    acc_no = "12345678-01"
    
    # 환경 변수에서 API 키 가져오기 (선택적)
    import os
    key = os.getenv("KOREA_INVESTMENT_API_KEY", key)
    secret = os.getenv("KOREA_INVESTMENT_API_SECRET", secret)
    acc_no = os.getenv("KOREA_INVESTMENT_ACC_NO", acc_no)
    
    if key == "YOUR_API_KEY" or secret == "YOUR_API_SECRET":
        print("\n⚠️  주의: API 키가 설정되지 않았습니다!")
        print("다음 중 하나의 방법으로 API 키를 설정하세요:")
        print("1. 이 파일의 key, secret, acc_no 변수를 직접 수정")
        print("2. 환경 변수 설정:")
        print("   export KOREA_INVESTMENT_API_KEY='your_api_key'")
        print("   export KOREA_INVESTMENT_API_SECRET='your_api_secret'")
        print("   export KOREA_INVESTMENT_ACC_NO='your_account_number'")
        print("\nAPI 키는 한국투자증권 KIS Developers에서 발급받을 수 있습니다.")
        print("https://apiportal.koreainvestment.com/")
        return
    
    # 브로커 객체 생성
    print(f"\nAPI 키 확인: {key[:10]}...")
    print(f"계좌 번호: {acc_no}")
    
    broker = korea_investment_stock.KoreaInvestment(
        api_key=key, 
        api_secret=secret, 
        acc_no=acc_no,
        mock=False  # 실제 서버 사용
    )
    
    try:
        # 각 예제 실행
        basic_rate_limiting_example(broker)
        statistics_example(broker)
        batch_processing_example(broker)
        # auto_save_example(broker)  # 시간이 오래 걸림
        peak_hour_example(broker)
        error_monitoring_example(broker)
        
    finally:
        # 리소스 정리
        print("\n리소스 정리 중...")
        broker.shutdown()
        print("완료!")


if __name__ == "__main__":
    main() 