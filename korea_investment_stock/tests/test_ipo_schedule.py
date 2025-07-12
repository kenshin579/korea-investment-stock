"""
공모주 청약 일정 API 단위 테스트
"""
import pytest
from datetime import datetime, timedelta
from unittest.mock import Mock, patch
import sys
import os

# 상위 디렉토리를 sys.path에 추가
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from korea_investment_stock import KoreaInvestment


class TestIPOHelpers:
    """IPO 관련 헬퍼 함수 테스트"""
    
    def test_validate_date_format(self):
        """날짜 형식 검증 테스트"""
        # Mock broker 생성
        broker = Mock(spec=KoreaInvestment)
        broker._validate_date_format = KoreaInvestment._validate_date_format.__get__(broker)
        
        # 유효한 날짜
        assert broker._validate_date_format("20240115") == True
        assert broker._validate_date_format("20231231") == True
        
        # 유효하지 않은 날짜
        assert broker._validate_date_format("2024-01-15") == False
        assert broker._validate_date_format("240115") == False
        assert broker._validate_date_format("20241315") == False  # 13월
        assert broker._validate_date_format("20240132") == False  # 32일
        assert broker._validate_date_format("abcdefgh") == False
        
    def test_validate_date_range(self):
        """날짜 범위 검증 테스트"""
        # Mock broker 생성
        broker = Mock(spec=KoreaInvestment)
        broker._validate_date_range = KoreaInvestment._validate_date_range.__get__(broker)
        
        # 유효한 범위
        assert broker._validate_date_range("20240101", "20240131") == True
        assert broker._validate_date_range("20240115", "20240115") == True  # 같은 날짜
        
        # 유효하지 않은 범위
        assert broker._validate_date_range("20240131", "20240101") == False  # 역순
        assert broker._validate_date_range("20240101", "20239999") == False  # 유효하지 않은 날짜
        
    def test_parse_ipo_date_range(self):
        """청약기간 파싱 테스트"""
        # 정상적인 청약기간
        start, end = KoreaInvestment.parse_ipo_date_range("2024.01.15~2024.01.16")
        assert start == datetime(2024, 1, 15)
        assert end == datetime(2024, 1, 16)
        
        # 빈 문자열
        start, end = KoreaInvestment.parse_ipo_date_range("")
        assert start is None
        assert end is None
        
        # 잘못된 형식
        start, end = KoreaInvestment.parse_ipo_date_range("2024-01-15~2024-01-16")
        assert start is None
        assert end is None
        
        # 부분적으로 잘못된 형식
        start, end = KoreaInvestment.parse_ipo_date_range("2024.01.15~16")
        assert start is None
        assert end is None
        
    def test_format_ipo_date(self):
        """날짜 형식 변환 테스트"""
        # YYYYMMDD -> YYYY-MM-DD
        assert KoreaInvestment.format_ipo_date("20240115") == "2024-01-15"
        
        # YYYY.MM.DD -> YYYY-MM-DD
        assert KoreaInvestment.format_ipo_date("2024.01.15") == "2024-01-15"
        
        # 기타 형식은 그대로 반환
        assert KoreaInvestment.format_ipo_date("2024-01-15") == "2024-01-15"
        assert KoreaInvestment.format_ipo_date("240115") == "240115"
        
    def test_calculate_ipo_d_day(self):
        """D-Day 계산 테스트"""
        # 미래 날짜 테스트
        future_date = (datetime.now() + timedelta(days=7)).strftime("%Y.%m.%d")
        date_range = f"{future_date}~{future_date}"
        d_day = KoreaInvestment.calculate_ipo_d_day(date_range)
        assert 6 <= d_day <= 7  # 시간 차이로 인한 오차 허용
        
        # 과거 날짜 테스트
        past_date = (datetime.now() - timedelta(days=7)).strftime("%Y.%m.%d")
        date_range = f"{past_date}~{past_date}"
        d_day = KoreaInvestment.calculate_ipo_d_day(date_range)
        assert d_day < 0
        
        # 잘못된 형식
        assert KoreaInvestment.calculate_ipo_d_day("2024.01.15") == -999
        assert KoreaInvestment.calculate_ipo_d_day("") == -999
        
    def test_get_ipo_status(self):
        """청약 상태 판단 테스트"""
        today = datetime.now()
        
        # 청약 예정
        future_start = (today + timedelta(days=7)).strftime("%Y.%m.%d")
        future_end = (today + timedelta(days=8)).strftime("%Y.%m.%d")
        assert KoreaInvestment.get_ipo_status(f"{future_start}~{future_end}") == "예정"
        
        # 청약 진행중
        past_start = (today - timedelta(days=1)).strftime("%Y.%m.%d")
        future_end = (today + timedelta(days=1)).strftime("%Y.%m.%d")
        assert KoreaInvestment.get_ipo_status(f"{past_start}~{future_end}") == "진행중"
        
        # 청약 마감
        past_start = (today - timedelta(days=8)).strftime("%Y.%m.%d")
        past_end = (today - timedelta(days=7)).strftime("%Y.%m.%d")
        assert KoreaInvestment.get_ipo_status(f"{past_start}~{past_end}") == "마감"
        
        # 잘못된 형식
        assert KoreaInvestment.get_ipo_status("2024.01.15") == "알수없음"
        assert KoreaInvestment.get_ipo_status("") == "알수없음"
        
    def test_format_number(self):
        """숫자 포맷팅 테스트"""
        assert KoreaInvestment.format_number("1000000") == "1,000,000"
        assert KoreaInvestment.format_number("15000") == "15,000"
        assert KoreaInvestment.format_number("100") == "100"
        assert KoreaInvestment.format_number("0") == "0"
        
        # 숫자가 아닌 경우 원본 반환
        assert KoreaInvestment.format_number("abc") == "abc"
        assert KoreaInvestment.format_number("") == ""
        assert KoreaInvestment.format_number(None) == None


class TestFetchIPOSchedule:
    """fetch_ipo_schedule 메서드 테스트"""
    
    @pytest.fixture
    def mock_broker(self):
        """테스트용 Mock broker 생성"""
        broker = Mock(spec=KoreaInvestment)
        broker.mock = False
        broker.base_url = "https://openapi.koreainvestment.com:9443"
        broker.access_token = "test_token"
        broker.api_key = "test_key"
        broker.api_secret = "test_secret"
        broker.rate_limiter = Mock()
        broker._validate_date_format = KoreaInvestment._validate_date_format.__get__(broker)
        broker._validate_date_range = KoreaInvestment._validate_date_range.__get__(broker)
        broker.fetch_ipo_schedule = KoreaInvestment.fetch_ipo_schedule.__get__(broker)
        return broker
    
    def test_mock_trading_error(self, mock_broker):
        """모의투자 에러 테스트"""
        mock_broker.mock = True
        
        with pytest.raises(ValueError) as exc_info:
            mock_broker.fetch_ipo_schedule()
        
        assert "모의투자를 지원하지 않습니다" in str(exc_info.value)
        
    def test_date_format_error(self, mock_broker):
        """날짜 형식 에러 테스트"""
        with pytest.raises(ValueError) as exc_info:
            mock_broker.fetch_ipo_schedule(from_date="2024-01-15", to_date="2024-01-31")
        
        assert "날짜 형식은 YYYYMMDD" in str(exc_info.value)
        
    def test_date_range_error(self, mock_broker):
        """날짜 범위 에러 테스트"""
        with pytest.raises(ValueError) as exc_info:
            mock_broker.fetch_ipo_schedule(from_date="20240131", to_date="20240101")
        
        assert "시작일은 종료일보다 이전" in str(exc_info.value)
        
    @patch('requests.get')
    def test_successful_request(self, mock_get, mock_broker):
        """정상적인 요청 테스트"""
        # Mock response
        mock_response = Mock()
        mock_response.json.return_value = {
            "rt_cd": "0",
            "msg_cd": "KIOK0000",
            "msg1": "정상처리되었습니다",
            "output1": [
                {
                    "record_date": "20240115",
                    "sht_cd": "123456",
                    "isin_name": "테스트회사",
                    "fix_subscr_pri": "15000",
                    "face_value": "5000",
                    "subscr_dt": "2024.01.15~2024.01.16",
                    "pay_dt": "2024.01.18",
                    "refund_dt": "2024.01.19",
                    "list_dt": "2024.01.25",
                    "lead_mgr": "한국투자증권",
                    "pub_bf_cap": "1000000000",
                    "pub_af_cap": "1500000000",
                    "assign_stk_qty": "100000"
                }
            ]
        }
        mock_get.return_value = mock_response
        
        # API 호출
        result = mock_broker.fetch_ipo_schedule(
            from_date="20240101",
            to_date="20240131"
        )
        
        # 결과 검증
        assert result["rt_cd"] == "0"
        assert len(result["output1"]) == 1
        assert result["output1"][0]["sht_cd"] == "123456"
        assert result["output1"][0]["isin_name"] == "테스트회사"
        
        # 요청 파라미터 검증
        mock_get.assert_called_once()
        call_args = mock_get.call_args
        assert call_args[1]["params"]["F_DT"] == "20240101"
        assert call_args[1]["params"]["T_DT"] == "20240131"
        assert call_args[1]["params"]["SHT_CD"] == ""
        assert call_args[1]["headers"]["tr_id"] == "HHKDB669108C0"
        assert call_args[1]["headers"]["custtype"] == "P"
        
    @patch('requests.get')
    def test_default_date_range(self, mock_get, mock_broker):
        """기본 날짜 범위 테스트"""
        # Mock response
        mock_response = Mock()
        mock_response.json.return_value = {
            "rt_cd": "0",
            "msg_cd": "KIOK0000",
            "msg1": "정상처리되었습니다",
            "output1": []
        }
        mock_get.return_value = mock_response
        
        # 날짜 파라미터 없이 호출
        result = mock_broker.fetch_ipo_schedule()
        
        # 기본값 검증
        call_args = mock_get.call_args
        from_date = call_args[1]["params"]["F_DT"]
        to_date = call_args[1]["params"]["T_DT"]
        
        # 오늘 날짜 확인
        assert from_date == datetime.now().strftime("%Y%m%d")
        
        # 30일 후 날짜 확인
        expected_to_date = (datetime.now() + timedelta(days=30)).strftime("%Y%m%d")
        assert to_date == expected_to_date
        
    @patch('requests.get')
    def test_specific_symbol_request(self, mock_get, mock_broker):
        """특정 종목 조회 테스트"""
        # Mock response
        mock_response = Mock()
        mock_response.json.return_value = {
            "rt_cd": "0",
            "msg_cd": "KIOK0000",
            "msg1": "정상처리되었습니다",
            "output1": []
        }
        mock_get.return_value = mock_response
        
        # 특정 종목 조회
        result = mock_broker.fetch_ipo_schedule(symbol="123456")
        
        # 종목코드 파라미터 검증
        call_args = mock_get.call_args
        assert call_args[1]["params"]["SHT_CD"] == "123456"
        
    @patch('requests.get')
    def test_error_response(self, mock_get, mock_broker):
        """에러 응답 테스트"""
        # Mock error response
        mock_response = Mock()
        mock_response.json.return_value = {
            "rt_cd": "7",
            "msg_cd": "KIOK0007",
            "msg1": "조회할 자료가 없습니다",
            "output1": []
        }
        mock_get.return_value = mock_response
        
        # API 호출
        result = mock_broker.fetch_ipo_schedule()
        
        # 에러 응답 확인
        assert result["rt_cd"] == "7"
        assert "조회할 자료가 없습니다" in result["msg1"]


if __name__ == "__main__":
    pytest.main([__file__, "-v"]) 