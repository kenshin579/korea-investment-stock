"""
공모주 청약 일정 조회 예제

한국투자증권 API를 사용하여 공모주 청약 일정을 조회하는 예제입니다.
"""
import os
import sys
from datetime import datetime, timedelta
from pathlib import Path
import json
import atexit

# 상위 디렉토리의 모듈을 import하기 위한 경로 추가
sys.path.append(str(Path(__file__).parent.parent))

from korea_investment_stock import KoreaInvestment


def load_credentials():
    """API 자격 증명 로드"""
    # 환경 변수에서 먼저 확인
    api_key = os.getenv('KOREA_INVESTMENT_API_KEY')
    api_secret = os.getenv('KOREA_INVESTMENT_API_SECRET')
    acc_no = os.getenv('KOREA_INVESTMENT_ACCOUNT_NO')
    
    # 키 파일에서 읽기
    if not all([api_key, api_secret, acc_no]):
        key_file = Path(__file__).parent.parent / "koreainvestment.key"
        if key_file.exists():
            with open(key_file, encoding='utf-8') as f:
                lines = f.readlines()
                if len(lines) >= 3:
                    api_key = lines[0].strip()
                    api_secret = lines[1].strip()
                    acc_no = lines[2].strip()
    
    if not all([api_key, api_secret, acc_no]):
        print("❌ API 자격 증명을 찾을 수 없습니다.")
        print("환경 변수를 설정하거나 koreainvestment.key 파일을 생성하세요.")
        sys.exit(1)
    
    return api_key, api_secret, acc_no


def example_basic_ipo_query(broker):
    """기본 공모주 조회 예제"""
    print("\n" + "="*60)
    print("📌 1. 기본 공모주 일정 조회 (이번 달)")
    print("="*60)
    
    # 이번 달 전체 공모주 조회
    today = datetime.now()
    # 이번 달 1일
    from_date = today.replace(day=1).strftime("%Y%m%d")
    # 이번 달 마지막 날
    if today.month == 12:
        next_month = today.replace(year=today.year + 1, month=1, day=1)
    else:
        next_month = today.replace(month=today.month + 1, day=1)
    to_date = (next_month - timedelta(days=1)).strftime("%Y%m%d")
    
    print(f"조회 기간: {from_date} ~ {to_date} (이번 달)")
    result = broker.fetch_ipo_schedule(from_date=from_date, to_date=to_date)
    
    if result['rt_cd'] == '0':
        ipos = result.get('output1', [])
        print(f"\n✅ 조회 성공: {len(ipos)}개의 공모주 정보를 찾았습니다.\n")
        
        # 처음 5개만 출력
        for i, ipo in enumerate(ipos[:5], 1):
            print(f"{i}. {ipo['isin_name']} ({ipo['sht_cd']})")
            print(f"   청약기간: {ipo['subscr_dt']}")
            print(f"   공모가: {broker.format_number(ipo['fix_subscr_pri'])}원")
            print(f"   주간사: {ipo['lead_mgr']}")
            print(f"   상장예정일: {ipo['list_dt']}")
            print()
        
        if len(ipos) > 5:
            print(f"   ... 외 {len(ipos) - 5}개")
    else:
        print(f"❌ 조회 실패: {result.get('msg1', 'Unknown error')}")


def example_period_query(broker):
    """특정 기간 공모주 조회 예제"""
    print("\n" + "="*60)
    print("📌 2. 특정 기간 공모주 조회 (지난달 ~ 다음달)")
    print("="*60)
    
    today = datetime.now()
    
    # 지난달 1일
    if today.month == 1:
        last_month = today.replace(year=today.year - 1, month=12, day=1)
    else:
        last_month = today.replace(month=today.month - 1, day=1)
    from_date = last_month.strftime("%Y%m%d")
    
    # 다음달 마지막 날
    if today.month >= 11:
        # 11월이나 12월인 경우
        if today.month == 11:
            next_next_month = today.replace(year=today.year + 1, month=1, day=1)
        else:  # 12월
            next_next_month = today.replace(year=today.year + 1, month=2, day=1)
    else:
        next_next_month = today.replace(month=today.month + 2, day=1)
    to_date = (next_next_month - timedelta(days=1)).strftime("%Y%m%d")
    
    print(f"조회 기간: {from_date} ~ {to_date} (지난달 ~ 다음달)")
    
    result = broker.fetch_ipo_schedule(
        from_date=from_date,
        to_date=to_date
    )
    
    if result['rt_cd'] == '0':
        ipos = result.get('output1', [])
        print(f"\n✅ 조회 성공: {len(ipos)}개의 공모주 정보를 찾았습니다.")
        
        # 상태별로 분류
        upcoming = []
        active = []
        closed = []
        
        for ipo in ipos:
            status = broker.get_ipo_status(ipo['subscr_dt'])
            if status == "예정":
                upcoming.append(ipo)
            elif status == "진행중":
                active.append(ipo)
            elif status == "마감":
                closed.append(ipo)
        
        print(f"\n📊 상태별 분류:")
        print(f"   - 청약 예정: {len(upcoming)}개")
        print(f"   - 청약 진행중: {len(active)}개")
        print(f"   - 청약 마감: {len(closed)}개")
        
        # 청약 진행중인 공모주 출력
        if active:
            print(f"\n🔥 현재 청약 진행중인 공모주:")
            for ipo in active:
                print(f"   - {ipo['isin_name']}: {ipo['subscr_dt']}")
    else:
        print(f"❌ 조회 실패: {result.get('msg1', 'Unknown error')}")


def example_upcoming_ipos(broker):
    """청약 예정 공모주 D-Day 표시 예제"""
    print("\n" + "="*60)
    print("📌 3. 청약 예정 공모주 (D-Day 표시)")
    print("="*60)
    
    today = datetime.now()
    
    # 오늘부터 다음달 말까지 조회
    from_date = today.strftime("%Y%m%d")
    
    # 다음달 마지막 날
    if today.month == 12:
        next_next_month = today.replace(year=today.year + 1, month=2, day=1)
    else:
        if today.month == 11:
            next_next_month = today.replace(year=today.year + 1, month=1, day=1)
        else:
            next_next_month = today.replace(month=today.month + 2, day=1)
    to_date = (next_next_month - timedelta(days=1)).strftime("%Y%m%d")
    
    print(f"조회 기간: {from_date} ~ {to_date} (오늘 ~ 다음달 말)")
    result = broker.fetch_ipo_schedule(from_date=from_date, to_date=to_date)
    
    if result['rt_cd'] == '0':
        upcoming_ipos = []
        
        for ipo in result.get('output1', []):
            status = broker.get_ipo_status(ipo['subscr_dt'])
            if status == "예정":
                d_day = broker.calculate_ipo_d_day(ipo['subscr_dt'])
                if 0 <= d_day <= 30:  # 30일 이내
                    upcoming_ipos.append({
                        'name': ipo['isin_name'],
                        'code': ipo['sht_cd'],
                        'subscr_dt': ipo['subscr_dt'],
                        'd_day': d_day,
                        'price': ipo['fix_subscr_pri'],
                        'lead_mgr': ipo['lead_mgr']
                    })
        
        # D-Day 기준 정렬
        upcoming_ipos.sort(key=lambda x: x['d_day'])
        
        if upcoming_ipos:
            print(f"\n✅ 향후 30일 이내 청약 예정: {len(upcoming_ipos)}개\n")
            for ipo in upcoming_ipos[:10]:  # 최대 10개만 표시
                print(f"D-{ipo['d_day']:2d} | {ipo['name']} ({ipo['code']})")
                print(f"      | 청약: {ipo['subscr_dt']}")
                print(f"      | 공모가: {broker.format_number(ipo['price'])}원")
                print(f"      | 주간사: {ipo['lead_mgr']}")
                print()
        else:
            print("향후 30일 이내에 청약 예정인 공모주가 없습니다.")
    else:
        print(f"❌ 조회 실패: {result.get('msg1', 'Unknown error')}")


def example_ipo_details(broker):
    """공모주 상세 정보 출력 예제"""
    print("\n" + "="*60)
    print("📌 4. 공모주 상세 정보")
    print("="*60)
    
    # 이번 달 공모주 조회
    today = datetime.now()
    from_date = today.replace(day=1).strftime("%Y%m%d")
    if today.month == 12:
        next_month = today.replace(year=today.year + 1, month=1, day=1)
    else:
        next_month = today.replace(month=today.month + 1, day=1)
    to_date = (next_month - timedelta(days=1)).strftime("%Y%m%d")
    
    result = broker.fetch_ipo_schedule(from_date=from_date, to_date=to_date)
    
    if result['rt_cd'] == '0' and result.get('output1'):
        # 첫 번째 공모주의 상세 정보 출력
        ipo = result['output1'][0]
        
        print(f"\n📋 공모주 상세 정보:")
        print(f"{'종목명':　<15}: {ipo['isin_name']} ({ipo['sht_cd']})")
        print(f"{'공모가':　<15}: {broker.format_number(ipo['fix_subscr_pri'])}원")
        print(f"{'액면가':　<15}: {broker.format_number(ipo['face_value'])}원")
        print(f"{'청약기간':　<15}: {ipo['subscr_dt']}")
        print(f"{'납입일':　<15}: {broker.format_ipo_date(ipo['pay_dt'])}")
        print(f"{'환불일':　<15}: {broker.format_ipo_date(ipo['refund_dt'])}")
        print(f"{'상장예정일':　<15}: {broker.format_ipo_date(ipo['list_dt'])}")
        print(f"{'주간사':　<15}: {ipo['lead_mgr']}")
        print(f"{'공모전 자본금':　<15}: {broker.format_number(ipo['pub_bf_cap'])}원")
        print(f"{'공모후 자본금':　<15}: {broker.format_number(ipo['pub_af_cap'])}원")
        print(f"{'당사배정물량':　<15}: {broker.format_number(ipo['assign_stk_qty'])}주")
        
        # 청약 상태 및 D-Day
        status = broker.get_ipo_status(ipo['subscr_dt'])
        d_day = broker.calculate_ipo_d_day(ipo['subscr_dt'])
        print(f"\n📊 청약 상태: {status}")
        if status == "예정" and d_day >= 0:
            print(f"📅 D-{d_day}")
        elif status == "진행중":
            print(f"🔥 현재 청약 진행중!")
    else:
        print("조회된 공모주가 없습니다.")


def example_save_to_file(broker):
    """공모주 정보를 파일로 저장하는 예제"""
    print("\n" + "="*60)
    print("📌 5. 공모주 정보 파일 저장 (이번 달)")
    print("="*60)
    
    # 이번 달 공모주 조회
    today = datetime.now()
    from_date = today.replace(day=1).strftime("%Y%m%d")
    if today.month == 12:
        next_month = today.replace(year=today.year + 1, month=1, day=1)
    else:
        next_month = today.replace(month=today.month + 1, day=1)
    to_date = (next_month - timedelta(days=1)).strftime("%Y%m%d")
    
    result = broker.fetch_ipo_schedule(from_date=from_date, to_date=to_date)
    
    if result['rt_cd'] == '0' and result.get('output1'):
        # JSON 파일로 저장
        filename = f"ipo_schedule_{datetime.now().strftime('%Y%m%d_%H%M%S')}.json"
        
        # 저장할 데이터 정리
        save_data = {
            'query_time': datetime.now().isoformat(),
            'total_count': len(result['output1']),
            'ipo_list': []
        }
        
        for ipo in result['output1']:
            status = broker.get_ipo_status(ipo['subscr_dt'])
            d_day = broker.calculate_ipo_d_day(ipo['subscr_dt'])
            
            save_data['ipo_list'].append({
                'name': ipo['isin_name'],
                'code': ipo['sht_cd'],
                'status': status,
                'd_day': d_day if d_day != -999 else None,
                'subscription_period': ipo['subscr_dt'],
                'ipo_price': ipo['fix_subscr_pri'],
                'listing_date': ipo['list_dt'],
                'lead_manager': ipo['lead_mgr'],
                'allocation_qty': ipo['assign_stk_qty']
            })
        
        # 파일 저장
        with open(filename, 'w', encoding='utf-8') as f:
            json.dump(save_data, f, ensure_ascii=False, indent=2)
        
        print(f"\n✅ 공모주 정보가 '{filename}' 파일로 저장되었습니다.")
        print(f"   - 총 {len(result['output1'])}개 공모주 정보 저장")
    else:
        print("저장할 공모주 정보가 없습니다.")


def main():
    """메인 함수"""
    print("\n" + "="*60)
    print("🎯 한국투자증권 공모주 청약 일정 조회 예제")
    print("="*60)
    
    # API 자격 증명 로드
    api_key, api_secret, acc_no = load_credentials()
    
    # broker 인스턴스 생성
    try:
        broker = KoreaInvestment(
            api_key=api_key,
            api_secret=api_secret,
            acc_no=acc_no,
            mock=False  # 실전투자 (공모주 조회는 모의투자 미지원)
        )
        
        # 통계 저장 비활성화
        if hasattr(broker.rate_limiter, 'enable_stats'):
            broker.rate_limiter.enable_stats = False
        if hasattr(broker.rate_limiter, 'disable_auto_save'):
            broker.rate_limiter.disable_auto_save()
        
        # atexit 핸들러 제거 (통계 저장 방지)
        # atexit에 등록된 모든 핸들러 중 broker.shutdown 관련 제거
        try:
            # atexit 내부 리스트 접근 (Python 버전에 따라 다를 수 있음)
            if hasattr(atexit, '_exithandlers'):
                # shutdown 관련 핸들러 제거
                atexit._exithandlers = [
                    (func, args, kwargs) 
                    for func, args, kwargs in atexit._exithandlers 
                    if not (hasattr(func, '__self__') and func.__self__ == broker)
                ]
        except:
            pass
        
        # shutdown 메서드를 빈 함수로 오버라이드 (통계 저장 방지)
        def empty_shutdown():
            if hasattr(broker, 'executor') and broker.executor:
                print("ThreadPoolExecutor 종료 중...")
                broker.executor.shutdown(wait=True)
                broker.executor = None
                print("ThreadPoolExecutor 종료 완료")
        
        broker.shutdown = empty_shutdown
            
        print("\n✅ API 연결 성공!")
        
        # 예제 실행
        example_basic_ipo_query(broker)
        example_period_query(broker)
        example_upcoming_ipos(broker)
        example_ipo_details(broker)
        example_save_to_file(broker)
        
        print("\n" + "="*60)
        print("✅ 모든 예제 실행 완료!")
        print("="*60)
        
        # 강제 종료 (통계 저장 방지)
        import os
        os._exit(0)
            
    except ValueError as e:
        if "모의투자를 지원하지 않습니다" in str(e):
            print("\n❌ 공모주 청약 일정 조회는 실전투자만 지원합니다.")
            print("mock=False로 설정하세요.")
        else:
            print(f"\n❌ 오류 발생: {e}")
    except Exception as e:
        print(f"\n❌ 예상치 못한 오류 발생: {e}")


if __name__ == "__main__":
    main() 