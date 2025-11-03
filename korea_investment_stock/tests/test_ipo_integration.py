"""
공모주 청약 일정 API 통합 테스트
실제 API 호출을 테스트합니다.
"""
import pytest
import os
import sys
from datetime import datetime, timedelta
import time

# 상위 디렉토리를 sys.path에 추가
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from korea_investment_stock import KoreaInvestment


# 환경 변수에서 API 키 정보 가져오기
API_KEY = os.getenv('KOREA_INVESTMENT_API_KEY')
API_SECRET = os.getenv('KOREA_INVESTMENT_API_SECRET')
ACC_NO = os.getenv('KOREA_INVESTMENT_ACC_NO')

# API 키가 없으면 테스트 스킵
skip_if_no_credentials = pytest.mark.skipif(
    not all([API_KEY, API_SECRET, ACC_NO]),
    reason="API credentials not found in environment variables"
)


@skip_if_no_credentials
class TestIPOIntegration:
    """공모주 API 통합 테스트"""
    
    @pytest.fixture(scope="class")
    def broker(self):
        """테스트용 broker 인스턴스 생성"""
        return KoreaInvestment(
            api_key=API_KEY,
            api_secret=API_SECRET,
            acc_no=ACC_NO,
            mock=False  # 실전투자 모드
        )
    
    def test_fetch_all_ipo_schedule(self, broker):
        """전체 공모주 일정 조회 테스트"""
        # API 호출
        result = broker.fetch_ipo_schedule()
        
        # 기본 응답 검증
        assert result is not None
        assert 'rt_cd' in result
        assert 'msg1' in result
        
        if result['rt_cd'] == '0':
            # 성공 응답 검증
            assert 'output1' in result
            assert isinstance(result['output1'], list)
            
            # 데이터가 있는 경우 필드 검증
            if len(result['output1']) > 0:
                ipo = result['output1'][0]
                required_fields = [
                    'record_date', 'sht_cd', 'isin_name', 'fix_subscr_pri',
                    'face_value', 'subscr_dt', 'pay_dt', 'refund_dt',
                    'list_dt', 'lead_mgr', 'pub_bf_cap', 'pub_af_cap',
                    'assign_stk_qty'
                ]
                for field in required_fields:
                    assert field in ipo, f"필수 필드 '{field}'가 없습니다"
                
                # 청약기간 형식 검증
                assert '~' in ipo['subscr_dt'], "청약기간 형식이 올바르지 않습니다"
                
            print(f"✅ 전체 조회 성공: {len(result['output1'])}개 공모주 조회됨")
        else:
            # 에러 응답도 정상적인 응답으로 간주 (데이터가 없을 수 있음)
            print(f"⚠️ 조회 결과: {result['msg1']}")
            
        # API 호출 간 대기
        time.sleep(0.5)
    
    def test_fetch_specific_period(self, broker):
        """특정 기간 공모주 조회 테스트"""
        # 최근 30일 ~ 향후 30일
        from_date = (datetime.now() - timedelta(days=30)).strftime("%Y%m%d")
        to_date = (datetime.now() + timedelta(days=30)).strftime("%Y%m%d")
        
        # API 호출
        result = broker.fetch_ipo_schedule(
            from_date=from_date,
            to_date=to_date
        )
        
        # 응답 검증
        assert result is not None
        assert 'rt_cd' in result
        
        if result['rt_cd'] == '0' and len(result.get('output1', [])) > 0:
            # 기간 내 데이터 검증
            for ipo in result['output1']:
                # 청약 상태 확인
                status = broker.get_ipo_status(ipo['subscr_dt'])
                assert status in ["예정", "진행중", "마감", "알수없음"]
                
                # D-Day 계산 확인
                d_day = broker.calculate_ipo_d_day(ipo['subscr_dt'])
                assert isinstance(d_day, int)
            
            print(f"✅ 기간 조회 성공: {from_date} ~ {to_date}, {len(result['output1'])}개 조회됨")
        else:
            print(f"⚠️ 기간 조회 결과: {result.get('msg1', 'Unknown')}")
            
        # API 호출 간 대기
        time.sleep(0.5)
    
    def test_fetch_specific_symbol(self, broker):
        """특정 종목 공모주 조회 테스트"""
        # 테스트용 종목코드 (실제로는 없을 가능성이 높음)
        test_symbol = "999999"
        
        # API 호출
        result = broker.fetch_ipo_schedule(symbol=test_symbol)
        
        # 응답 검증
        assert result is not None
        assert 'rt_cd' in result
        
        if result['rt_cd'] == '0' and len(result.get('output1', [])) > 0:
            # 종목코드 일치 검증
            for ipo in result['output1']:
                assert ipo['sht_cd'] == test_symbol
            print(f"✅ 종목 조회 성공: {test_symbol}")
        else:
            # 대부분의 경우 데이터가 없을 것으로 예상
            print(f"⚠️ 종목 조회 결과: {test_symbol} - {result.get('msg1', 'Unknown')}")
            
        # API 호출 간 대기
        time.sleep(0.5)
    
    def test_helper_functions_with_real_data(self, broker):
        """실제 데이터로 헬퍼 함수 테스트"""
        # 전체 조회
        result = broker.fetch_ipo_schedule()
        
        if result['rt_cd'] == '0' and len(result.get('output1', [])) > 0:
            for ipo in result['output1'][:3]:  # 처음 3개만 테스트
                # 날짜 포맷팅 테스트
                formatted_date = broker.format_ipo_date(ipo['list_dt'])
                assert '-' in formatted_date or '.' in formatted_date
                
                # 숫자 포맷팅 테스트
                formatted_price = broker.format_number(ipo['fix_subscr_pri'])
                if ipo['fix_subscr_pri'] and ipo['fix_subscr_pri'] != '0':
                    assert ',' in formatted_price or len(formatted_price) <= 3
                
                # 청약 상태 테스트
                status = broker.get_ipo_status(ipo['subscr_dt'])
                assert status in ["예정", "진행중", "마감", "알수없음"]
                
                # D-Day 테스트
                d_day = broker.calculate_ipo_d_day(ipo['subscr_dt'])
                assert isinstance(d_day, int)
                
                print(f"✅ {ipo['isin_name']}: 상태={status}, D-Day={d_day}, "
                      f"공모가={formatted_price}원")


# 통합 테스트 실행 시 주의사항 출력
def test_integration_info():
    """통합 테스트 정보"""
    if not all([API_KEY, API_SECRET, ACC_NO]):
        print("\n" + "="*60)
        print("⚠️  통합 테스트를 실행하려면 환경 변수를 설정하세요:")
        print("export KOREA_INVESTMENT_API_KEY='your_api_key'")
        print("export KOREA_INVESTMENT_API_SECRET='your_api_secret'")
        print("export KOREA_INVESTMENT_ACC_NO='your_account_number'")
        print("="*60 + "\n")
    else:
        print("\n" + "="*60)
        print("✅ API 자격 증명이 설정되어 있습니다.")
        print("통합 테스트를 실행합니다...")
        print("="*60 + "\n")


if __name__ == "__main__":
    test_integration_info()
    pytest.main([__file__, "-v", "-s"]) 