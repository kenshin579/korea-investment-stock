"""
Stress Test Example

종목 리스트를 순회하며 종목 정보와 가격을 조회하는 간단한 stress test
각 API 호출 사이에 100ms sleep을 적용합니다.
"""

import os
import time
import yaml
from pathlib import Path
from korea_investment_stock import KoreaInvestment


def load_stock_list(yaml_path: str) -> list:
    """
    YAML 파일에서 종목 리스트 로드

    Args:
        yaml_path: YAML 파일 경로

    Returns:
        종목 리스트 [["symbol", "market"], ...]
    """
    with open(yaml_path, 'r', encoding='utf-8') as f:
        data = yaml.safe_load(f)
    return data['stock_list']


def run_stress_test():
    """
    종목 리스트를 순회하며 API 호출 stress test 실행

    각 종목에 대해:
    1. fetch_stock_info() 호출
    2. 100ms sleep
    3. fetch_price() 호출
    4. 100ms sleep
    """
    # Environment variables
    api_key = os.environ.get('KOREA_INVESTMENT_API_KEY')
    api_secret = os.environ.get('KOREA_INVESTMENT_API_SECRET')
    acc_no = os.environ.get('KOREA_INVESTMENT_ACCOUNT_NO')

    if not all([api_key, api_secret, acc_no]):
        print("❌ Error: 환경변수를 설정해주세요:")
        print("  - KOREA_INVESTMENT_API_KEY")
        print("  - KOREA_INVESTMENT_API_SECRET")
        print("  - KOREA_INVESTMENT_ACCOUNT_NO")
        return

    # Load stock list
    yaml_path = Path(__file__).parent / 'testdata' / 'stock_list.yaml'
    stock_list = load_stock_list(yaml_path)

    print(f"📋 총 {len(stock_list)}개 종목 stress test 시작")
    print("=" * 60)

    success_count = 0
    error_count = 0
    start_time = time.time()

    # Initialize broker with context manager
    with KoreaInvestment(api_key, api_secret, acc_no) as broker:
        for i, (symbol, market) in enumerate(stock_list, 1):
            print(f"\n[{i}/{len(stock_list)}] {symbol} ({market})")

            # 1. fetch_stock_info
            try:
                info_result = broker.fetch_stock_info(symbol, market)
                if info_result['rt_cd'] == '0':
                    print(f"  ✅ Stock Info: Success")
                    success_count += 1
                else:
                    print(f"  ⚠️  Stock Info: {info_result['msg1']}")
                    error_count += 1
                    print("\n🚨 실패 감지: Stress test 중단")
                    break
            except Exception as e:
                print(f"  ❌ Stock Info Error: {e}")
                error_count += 1
                print("\n🚨 예외 발생: Stress test 중단")
                break

            # time.sleep(0.1)  # 100ms sleep

            # 2. fetch_price
            try:
                price_result = broker.fetch_price(symbol, market)
                if price_result['rt_cd'] == '0':
                    print(f"  ✅ Price: Success")
                    success_count += 1
                else:
                    print(f"  ⚠️  Price: {price_result['msg1']}")
                    error_count += 1
                    print("\n🚨 실패 감지: Stress test 중단")
                    break
            except Exception as e:
                print(f"  ❌ Price Error: {e}")
                error_count += 1
                print("\n🚨 예외 발생: Stress test 중단")
                break

            # time.sleep(0.1)  # 100ms sleep

    # Summary
    elapsed_time = time.time() - start_time
    total_calls = success_count + error_count
    avg_time = elapsed_time / total_calls if total_calls > 0 else 0

    print("\n" + "=" * 60)
    print("📊 Stress Test 결과")
    print("=" * 60)
    print(f"총 API 호출: {total_calls}회")
    print(f"성공: {success_count}회")
    print(f"실패: {error_count}회")
    print(f"성공률: {success_count / total_calls * 100:.1f}%" if total_calls > 0 else "N/A")
    print(f"실행 시간: {elapsed_time:.2f}초")
    print(f"평균 응답 시간: {avg_time:.3f}초/호출")


if __name__ == "__main__":
    run_stress_test()
