#!/usr/bin/env python3
"""
공모주청약일정 API 사용 예제
한국투자증권 API의 공모주 정보 조회 기능을 보여줍니다.
"""

from korea_investment_stock import KoreaInvestment
from datetime import datetime, timedelta


def main():
    # API 인증 정보
    api_key = "YOUR_API_KEY"
    api_secret = "YOUR_API_SECRET" 
    account_no = "12345678-01"
    
    # 클라이언트 생성
    broker = KoreaInvestment(
        api_key=api_key,
        api_secret=api_secret,
        acc_no=account_no,
        mock=False  # 실거래 환경
    )
    
    print("=== 공모주청약일정 API 사용 예제 ===\n")
    
    # 1. 전체 공모주 일정 조회
    print("1. 이번 달 공모주 일정 조회")
    print("-" * 50)
    
    today = datetime.now()
    month_end = (today.replace(day=1) + timedelta(days=32)).replace(day=1) - timedelta(days=1)
    
    result = broker.fetch_ipo_schedule(
        from_date=today.strftime("%Y%m%d"),
        to_date=month_end.strftime("%Y%m%d")
    )
    
    if result.get('rt_cd') == '0':
        ipos = result.get('output2', [])
        print(f"조회된 공모주 수: {len(ipos)}개\n")
        
        for ipo in ipos:
            print(f"종목명: {ipo.get('stck_name')}")
            print(f"종목코드: {ipo.get('stck_shrn_iscd')}")
            print(f"청약일: {format_date(ipo.get('pbof_date'))} ~ {format_date(ipo.get('pbof_end_date'))}")
            print(f"상장예정일: {format_date(ipo.get('list_date'))}")
            print(f"공모가: {format_number(ipo.get('pbof_pric'))}원")
            print(f"경쟁률: {ipo.get('pbof_comp_rate')}:1")
            print(f"상태: {get_status_name(ipo.get('pbof_stat'))}")
            print(f"주간사: {ipo.get('lead_mgr')}")
            print("-" * 30)
    else:
        print(f"조회 실패: {result.get('msg1')}")
    
    # 2. 청약 진행중인 공모주 조회
    print("\n2. 현재 청약 진행중인 공모주")
    print("-" * 50)
    
    active_ipos = broker.fetch_active_ipos()
    
    if active_ipos:
        print(f"진행중인 공모주: {len(active_ipos)}개\n")
        for ipo in active_ipos:
            print(f"- {ipo.get('stck_name')} ({ipo.get('stck_shrn_iscd')})")
            print(f"  청약마감: {format_date(ipo.get('pbof_end_date'))}")
            print(f"  경쟁률: {ipo.get('pbof_comp_rate')}:1")
    else:
        print("현재 청약 진행중인 공모주가 없습니다.")
    
    # 3. 향후 청약 예정 공모주
    print("\n3. 향후 7일 이내 청약 예정 공모주")
    print("-" * 50)
    
    upcoming_ipos = broker.fetch_upcoming_ipos(days=7)
    
    if upcoming_ipos:
        print(f"예정 공모주: {len(upcoming_ipos)}개\n")
        for ipo in upcoming_ipos:
            d_day = calculate_d_day(ipo.get('pbof_date'))
            print(f"- {ipo.get('stck_name')} (D{d_day:+d})")
            print(f"  청약일: {format_date(ipo.get('pbof_date'))}")
            print(f"  공모가: {format_number(ipo.get('pbof_pric'))}원")
    else:
        print("7일 이내 청약 예정 공모주가 없습니다.")
    
    # 4. 높은 경쟁률 공모주 조회
    print("\n4. 경쟁률 100:1 이상 공모주 (최근 3개월)")
    print("-" * 50)
    
    hot_ipos = broker.fetch_ipo_by_competition(min_rate=100.0)
    
    if hot_ipos:
        print(f"조회된 공모주: {len(hot_ipos)}개\n")
        for i, ipo in enumerate(hot_ipos[:5]):  # 상위 5개만 표시
            print(f"{i+1}. {ipo.get('stck_name')}")
            print(f"   경쟁률: {ipo.get('pbof_comp_rate')}:1")
            print(f"   청약일: {format_date(ipo.get('pbof_date'))}")
    else:
        print("경쟁률 100:1 이상 공모주가 없습니다.")
    
    # 5. 특정 종목 공모 정보 조회
    print("\n5. 특정 종목 공모 정보 조회")
    print("-" * 50)
    
    symbol = input("조회할 종목코드 입력 (예: 123456): ")
    if symbol:
        ipo_info = broker.fetch_ipo_info(symbol)
        
        if ipo_info.get('rt_cd') == '0':
            ipo = ipo_info.get('output')
            print(f"\n종목명: {ipo.get('stck_name')}")
            print(f"청약일: {format_date(ipo.get('pbof_date'))} ~ {format_date(ipo.get('pbof_end_date'))}")
            print(f"공모가: {format_number(ipo.get('pbof_pric'))}원")
            print(f"공모수량: {format_number(ipo.get('pbof_qty'))}주")
            print(f"최소/최대 청약: {ipo.get('min_sbsc_qty')}주 ~ {ipo.get('max_sbsc_qty')}주")
            print(f"개인배정비율: {ipo.get('indv_allm_rate')}%")
        else:
            print(f"조회 실패: {ipo_info.get('msg1')}")


def format_date(date_str: str) -> str:
    """날짜 형식 변환 (YYYYMMDD -> YYYY-MM-DD)"""
    if date_str and len(date_str) == 8:
        return f"{date_str[:4]}-{date_str[4:6]}-{date_str[6:8]}"
    return date_str


def format_number(num_str: str) -> str:
    """숫자 천단위 콤마 추가"""
    try:
        return f"{int(num_str):,}"
    except (ValueError, TypeError):
        return num_str


def get_status_name(status_code: str) -> str:
    """공모상태 코드를 한글로 변환"""
    status_map = {
        "01": "청약예정",
        "02": "청약진행중", 
        "03": "청약마감"
    }
    return status_map.get(status_code, status_code)


def calculate_d_day(date_str: str) -> int:
    """D-Day 계산"""
    try:
        target = datetime.strptime(date_str, "%Y%m%d")
        today = datetime.now()
        return (target - today).days
    except (ValueError, TypeError):
        return 0


if __name__ == "__main__":
    main() 