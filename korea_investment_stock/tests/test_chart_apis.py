"""차트 데이터 API 단위 테스트 (Step 1)"""
import pytest
from unittest.mock import Mock, patch


class MockResponse:
    """Mock HTTP 응답 객체"""

    def __init__(self, json_data):
        self._json = json_data

    def json(self):
        return self._json


def _create_broker_mock():
    """테스트용 KoreaInvestment 인스턴스 생성"""
    from korea_investment_stock import KoreaInvestment

    with patch.object(KoreaInvestment, '__init__', lambda x: None):
        broker = KoreaInvestment.__new__(KoreaInvestment)
        broker.base_url = "https://openapi.koreainvestment.com:9443"
        broker.access_token = "Bearer test_token"
        broker.api_key = "test_api_key"
        broker.api_secret = "test_api_secret"
        broker._token_manager = Mock()
        return broker


class TestFetchDomesticChart:
    """fetch_domestic_chart 단위 테스트"""

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_success_response(self, mock_get):
        """성공 응답 테스트"""
        broker = _create_broker_mock()

        mock_response = {
            'rt_cd': '0',
            'msg_cd': 'MCA00000',
            'msg1': '정상처리되었습니다',
            'output1': {
                'stck_prpr': '70000',
                'hts_kor_isnm': '삼성전자',
            },
            'output2': [
                {
                    'stck_bsop_date': '20241231',
                    'stck_clpr': '70000',
                    'stck_oprc': '69500',
                    'stck_hgpr': '70500',
                    'stck_lwpr': '69000',
                    'acml_vol': '10000000',
                    'acml_tr_pbmn': '700000000000',
                }
            ]
        }
        mock_get.return_value = MockResponse(mock_response)

        result = broker.fetch_domestic_chart("005930", "D", "20240101", "20241231")

        assert result['rt_cd'] == '0'
        assert 'output1' in result
        assert 'output2' in result
        assert len(result['output2']) > 0
        assert result['output2'][0]['stck_clpr'] == '70000'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_url_and_headers(self, mock_get):
        """요청 URL 및 헤더 검증"""
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})

        broker.fetch_domestic_chart("005930", "D", "20240101", "20241231")

        call_args = mock_get.call_args
        url = call_args[0][0]
        headers = call_args[1]['headers']

        assert "inquire-daily-itemchartprice" in url
        assert headers['tr_id'] == 'FHKST03010100'
        assert headers['appKey'] == 'test_api_key'
        assert headers['appSecret'] == 'test_api_secret'
        assert headers['authorization'] == 'Bearer test_token'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_params_default(self, mock_get):
        """기본 파라미터 검증"""
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})

        broker.fetch_domestic_chart("005930", "D", "20240101", "20241231")

        params = mock_get.call_args[1]['params']

        assert params['FID_COND_MRKT_DIV_CODE'] == 'J'
        assert params['FID_INPUT_ISCD'] == '005930'
        assert params['FID_INPUT_DATE_1'] == '20240101'
        assert params['FID_INPUT_DATE_2'] == '20241231'
        assert params['FID_PERIOD_DIV_CODE'] == 'D'
        assert params['FID_ORG_ADJ_PRC'] == '0'  # adjusted=True → "0"

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_period_options(self, mock_get):
        """기간 옵션 테스트 (D, W, M, Y)"""
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})

        for period in ["D", "W", "M", "Y"]:
            broker.fetch_domestic_chart("005930", period, "20240101", "20241231")
            params = mock_get.call_args[1]['params']
            assert params['FID_PERIOD_DIV_CODE'] == period

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_adjusted_false(self, mock_get):
        """원주가 옵션 테스트"""
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})

        broker.fetch_domestic_chart("005930", adjusted=False)

        params = mock_get.call_args[1]['params']
        assert params['FID_ORG_ADJ_PRC'] == '1'  # adjusted=False → "1"

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_token_refresh_on_expiry(self, mock_get):
        """토큰 만료 시 자동 재발급 테스트"""
        broker = _create_broker_mock()

        mock_get.side_effect = [
            MockResponse({"rt_cd": "1", "msg1": "기간이 만료된 token 입니다"}),
            MockResponse({"rt_cd": "0", "output1": {}, "output2": []})
        ]

        def refresh_token(force=False):
            if force:
                broker.access_token = "Bearer new_token"

        broker.issue_access_token = Mock(side_effect=refresh_token)

        result = broker.fetch_domestic_chart("005930", "D", "20240101", "20241231")

        broker.issue_access_token.assert_called_once_with(force=True)
        assert result["rt_cd"] == "0"


class TestFetchDomesticMinuteChart:
    """fetch_domestic_minute_chart 단위 테스트"""

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_success_response(self, mock_get):
        """성공 응답 테스트"""
        broker = _create_broker_mock()

        mock_response = {
            'rt_cd': '0',
            'msg_cd': 'MCA00000',
            'msg1': '정상처리되었습니다',
            'output1': {
                'stck_prpr': '70000',
                'hts_kor_isnm': '삼성전자',
            },
            'output2': [
                {
                    'stck_bsop_date': '20241231',
                    'stck_cntg_hour': '100000',
                    'stck_prpr': '70000',
                    'stck_oprc': '69500',
                    'stck_hgpr': '70500',
                    'stck_lwpr': '69000',
                    'cntg_vol': '50000',
                    'acml_tr_pbmn': '700000000',
                }
            ]
        }
        mock_get.return_value = MockResponse(mock_response)

        result = broker.fetch_domestic_minute_chart("005930", "090000")

        assert result['rt_cd'] == '0'
        assert 'output2' in result
        assert result['output2'][0]['stck_cntg_hour'] == '100000'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_url_and_headers(self, mock_get):
        """요청 URL 및 헤더 검증"""
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})

        broker.fetch_domestic_minute_chart("005930", "090000")

        call_args = mock_get.call_args
        url = call_args[0][0]
        headers = call_args[1]['headers']

        assert "inquire-time-itemchartprice" in url
        assert headers['tr_id'] == 'FHKST03010200'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_params(self, mock_get):
        """파라미터 검증"""
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})

        broker.fetch_domestic_minute_chart("005930", "093000")

        params = mock_get.call_args[1]['params']

        assert params['FID_COND_MRKT_DIV_CODE'] == 'J'
        assert params['FID_INPUT_ISCD'] == '005930'
        assert params['FID_INPUT_HOUR_1'] == '093000'
        assert params['FID_PW_DATA_INCU_YN'] == 'N'
        assert params['FID_ETC_CLS_CODE'] == ''

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_token_refresh_on_expiry(self, mock_get):
        """토큰 만료 시 자동 재발급 테스트"""
        broker = _create_broker_mock()

        mock_get.side_effect = [
            MockResponse({"rt_cd": "1", "msg1": "기간이 만료된 token 입니다"}),
            MockResponse({"rt_cd": "0", "output1": {}, "output2": []})
        ]

        def refresh_token(force=False):
            if force:
                broker.access_token = "Bearer new_token"

        broker.issue_access_token = Mock(side_effect=refresh_token)

        result = broker.fetch_domestic_minute_chart("005930", "090000")

        broker.issue_access_token.assert_called_once_with(force=True)
        assert result["rt_cd"] == "0"


class TestFetchOverseasChart:
    """fetch_overseas_chart 단위 테스트"""

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_success_response(self, mock_get):
        """성공 응답 테스트"""
        broker = _create_broker_mock()

        mock_response = {
            'rt_cd': '0',
            'msg_cd': 'MCA00000',
            'msg1': '정상처리되었습니다',
            'output1': {
                'rsym': 'DNASAAPL',
                'zdiv': '2',
                'nrec': '150.00',
            },
            'output2': [
                {
                    'xymd': '20241231',
                    'clos': '250.00',
                    'open': '248.00',
                    'high': '252.00',
                    'low': '247.00',
                    'tvol': '5000000',
                    'tamt': '1250000000',
                    'sign': '2',
                    'diff': '2.00',
                    'rate': '0.81',
                }
            ]
        }
        mock_get.return_value = MockResponse(mock_response)

        result = broker.fetch_overseas_chart("AAPL", "US", "D")

        assert result['rt_cd'] == '0'
        assert 'output2' in result
        assert result['output2'][0]['clos'] == '250.00'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_url_and_headers(self, mock_get):
        """요청 URL 및 헤더 검증"""
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})

        broker.fetch_overseas_chart("AAPL", "US", "D")

        call_args = mock_get.call_args
        url = call_args[0][0]
        headers = call_args[1]['headers']

        assert "dailyprice" in url
        assert headers['tr_id'] == 'HHDFS76240000'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_params_us(self, mock_get):
        """미국 주식 파라미터 검증"""
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})

        broker.fetch_overseas_chart("AAPL", "US", "D", "20241231", True)

        params = mock_get.call_args[1]['params']

        assert params['AUTH'] == ''
        assert params['EXCD'] == 'NYS'
        assert params['SYMB'] == 'AAPL'
        assert params['GUBN'] == '0'  # D → 0
        assert params['BYMD'] == '20241231'
        assert params['MODP'] == '1'  # adjusted=True → "1"

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_period_mapping(self, mock_get):
        """기간 코드 변환 테스트"""
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})

        expected = {"D": "0", "W": "1", "M": "2"}
        for period, gubn in expected.items():
            broker.fetch_overseas_chart("AAPL", "US", period)
            params = mock_get.call_args[1]['params']
            assert params['GUBN'] == gubn, f"period={period} should map to GUBN={gubn}"

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_adjusted_false(self, mock_get):
        """수정주가 미반영 테스트"""
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})

        broker.fetch_overseas_chart("AAPL", adjusted=False)

        params = mock_get.call_args[1]['params']
        assert params['MODP'] == '0'  # adjusted=False → "0"

    def test_invalid_country_code(self):
        """지원하지 않는 국가 코드 에러 테스트"""
        broker = _create_broker_mock()

        with pytest.raises(ValueError, match="지원하지 않는 country_code"):
            broker.fetch_overseas_chart("AAPL", "XX")

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_token_refresh_on_expiry(self, mock_get):
        """토큰 만료 시 자동 재발급 테스트"""
        broker = _create_broker_mock()

        mock_get.side_effect = [
            MockResponse({"rt_cd": "1", "msg1": "기간이 만료된 token 입니다"}),
            MockResponse({"rt_cd": "0", "output1": {}, "output2": []})
        ]

        def refresh_token(force=False):
            if force:
                broker.access_token = "Bearer new_token"

        broker.issue_access_token = Mock(side_effect=refresh_token)

        result = broker.fetch_overseas_chart("AAPL", "US", "D")

        broker.issue_access_token.assert_called_once_with(force=True)
        assert result["rt_cd"] == "0"


class TestCachedChartApis:
    """CachedKoreaInvestment의 차트 API 캐싱 테스트"""

    def test_domestic_chart_cache_hit(self):
        """국내 차트 캐시 히트 테스트"""
        from korea_investment_stock import KoreaInvestment
        from korea_investment_stock.cache import CachedKoreaInvestment

        mock_broker = Mock(spec=KoreaInvestment)
        mock_broker.fetch_domestic_chart.return_value = {
            'rt_cd': '0',
            'output2': [{'stck_clpr': '70000'}]
        }

        cached = CachedKoreaInvestment(mock_broker, enable_cache=True)

        result1 = cached.fetch_domestic_chart("005930", "D", "20240101", "20241231")
        result2 = cached.fetch_domestic_chart("005930", "D", "20240101", "20241231")

        assert result1['rt_cd'] == '0'
        assert result2['rt_cd'] == '0'
        assert mock_broker.fetch_domestic_chart.call_count == 1

    def test_overseas_chart_cache_disabled(self):
        """해외 차트 캐시 비활성화 테스트"""
        from korea_investment_stock import KoreaInvestment
        from korea_investment_stock.cache import CachedKoreaInvestment

        mock_broker = Mock(spec=KoreaInvestment)
        mock_broker.fetch_overseas_chart.return_value = {
            'rt_cd': '0',
            'output2': []
        }

        cached = CachedKoreaInvestment(mock_broker, enable_cache=False)

        cached.fetch_overseas_chart("AAPL", "US", "D")
        cached.fetch_overseas_chart("AAPL", "US", "D")

        assert mock_broker.fetch_overseas_chart.call_count == 2


class TestRateLimitedChartApis:
    """RateLimitedKoreaInvestment의 차트 API 테스트"""

    def test_domestic_chart_rate_limit(self):
        """국내 차트 속도 제한 테스트"""
        from korea_investment_stock import KoreaInvestment
        from korea_investment_stock.rate_limit import RateLimitedKoreaInvestment

        mock_broker = Mock(spec=KoreaInvestment)
        mock_broker.fetch_domestic_chart.return_value = {
            'rt_cd': '0',
            'output2': []
        }

        rate_limited = RateLimitedKoreaInvestment(mock_broker, calls_per_second=10)
        result = rate_limited.fetch_domestic_chart("005930", "D", "20240101", "20241231")

        assert result['rt_cd'] == '0'
        mock_broker.fetch_domestic_chart.assert_called_once()

    def test_overseas_chart_rate_limit(self):
        """해외 차트 속도 제한 테스트"""
        from korea_investment_stock import KoreaInvestment
        from korea_investment_stock.rate_limit import RateLimitedKoreaInvestment

        mock_broker = Mock(spec=KoreaInvestment)
        mock_broker.fetch_overseas_chart.return_value = {
            'rt_cd': '0',
            'output2': []
        }

        rate_limited = RateLimitedKoreaInvestment(mock_broker, calls_per_second=10)
        result = rate_limited.fetch_overseas_chart("AAPL", "US", "D")

        assert result['rt_cd'] == '0'
        mock_broker.fetch_overseas_chart.assert_called_once()
