"""배당 + 업종 API 단위 테스트 (Step 4)"""
import pytest
from unittest.mock import Mock, patch


class MockResponse:
    def __init__(self, json_data):
        self._json = json_data

    def json(self):
        return self._json


def _create_broker_mock():
    from korea_investment_stock import KoreaInvestment
    with patch.object(KoreaInvestment, '__init__', lambda x: None):
        broker = KoreaInvestment.__new__(KoreaInvestment)
        broker.base_url = "https://openapi.koreainvestment.com:9443"
        broker.access_token = "Bearer test_token"
        broker.api_key = "test_api_key"
        broker.api_secret = "test_api_secret"
        broker._token_manager = Mock()
        return broker


class TestFetchDividendRanking:

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_success_response(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({
            'rt_cd': '0',
            'output1': [{'rank': '1', 'isin_name': '테스트종목', 'sht_cd': '005930', 'divi_rate': '5.2', 'per_sto_divi_amt': '1500'}]
        })
        result = broker.fetch_dividend_ranking()
        assert result['rt_cd'] == '0'
        assert result['output1'][0]['rank'] == '1'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_url_and_headers(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        broker.fetch_dividend_ranking()
        url = mock_get.call_args[0][0]
        headers = mock_get.call_args[1]['headers']
        assert "dividend-rate" in url
        assert headers['tr_id'] == 'HHKDB13470100'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_params_default(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        broker.fetch_dividend_ranking()
        params = mock_get.call_args[1]['params']
        assert params['GB1'] == '0'
        assert params['GB3'] == '2'
        assert params['GB4'] == '0'
        assert params['CTS_AREA'] == ''

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_params_custom(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        broker.fetch_dividend_ranking("1", "1", "20240101", "20241231", "1")
        params = mock_get.call_args[1]['params']
        assert params['GB1'] == '1'
        assert params['GB3'] == '1'
        assert params['F_DT'] == '20240101'
        assert params['T_DT'] == '20241231'
        assert params['GB4'] == '1'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_token_refresh_on_expiry(self, mock_get):
        broker = _create_broker_mock()
        mock_get.side_effect = [
            MockResponse({"rt_cd": "1", "msg1": "기간이 만료된 token 입니다"}),
            MockResponse({"rt_cd": "0", "output1": []})
        ]
        def refresh_token(force=False):
            if force:
                broker.access_token = "Bearer new_token"
        broker.issue_access_token = Mock(side_effect=refresh_token)
        result = broker.fetch_dividend_ranking()
        broker.issue_access_token.assert_called_once_with(force=True)
        assert result["rt_cd"] == "0"


class TestFetchIndustryIndex:

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_success_response(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({
            'rt_cd': '0',
            'output': {'bstp_nmix_prpr': '2500.00', 'bstp_nmix_prdy_vrss': '10.50', 'bstp_nmix_prdy_ctrt': '0.42'}
        })
        result = broker.fetch_industry_index("0001")
        assert result['rt_cd'] == '0'
        assert result['output']['bstp_nmix_prpr'] == '2500.00'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_url_and_headers(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        broker.fetch_industry_index()
        url = mock_get.call_args[0][0]
        headers = mock_get.call_args[1]['headers']
        assert "inquire-index-price" in url
        assert headers['tr_id'] == 'FHPUP02100000'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_params(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        broker.fetch_industry_index("1001")
        params = mock_get.call_args[1]['params']
        assert params['FID_COND_MRKT_DIV_CODE'] == 'U'
        assert params['FID_INPUT_ISCD'] == '1001'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_industry_code_options(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        for code in ["0001", "1001", "2001"]:
            broker.fetch_industry_index(code)
            params = mock_get.call_args[1]['params']
            assert params['FID_INPUT_ISCD'] == code

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_token_refresh_on_expiry(self, mock_get):
        broker = _create_broker_mock()
        mock_get.side_effect = [
            MockResponse({"rt_cd": "1", "msg1": "기간이 만료된 token 입니다"}),
            MockResponse({"rt_cd": "0", "output": {}})
        ]
        def refresh_token(force=False):
            if force:
                broker.access_token = "Bearer new_token"
        broker.issue_access_token = Mock(side_effect=refresh_token)
        result = broker.fetch_industry_index()
        broker.issue_access_token.assert_called_once_with(force=True)
        assert result["rt_cd"] == "0"


class TestFetchIndustryCategoryPrice:

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_success_response(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({
            'rt_cd': '0',
            'output1': {'bstp_nmix_prpr': '2500.00'},
            'output2': [{'bstp_cls_code': '0001', 'hts_kor_isnm': '종합', 'bstp_nmix_prpr': '2500.00'}]
        })
        result = broker.fetch_industry_category_price()
        assert result['rt_cd'] == '0'
        assert result['output2'][0]['hts_kor_isnm'] == '종합'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_url_and_headers(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        broker.fetch_industry_category_price()
        url = mock_get.call_args[0][0]
        headers = mock_get.call_args[1]['headers']
        assert "inquire-index-category-price" in url
        assert headers['tr_id'] == 'FHPUP02140000'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_params_kospi(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        broker.fetch_industry_category_price("K", "0")
        params = mock_get.call_args[1]['params']
        assert params['FID_COND_MRKT_DIV_CODE'] == 'U'
        assert params['FID_INPUT_ISCD'] == '0001'
        assert params['FID_COND_SCR_DIV_CODE'] == '20214'
        assert params['FID_MRKT_CLS_CODE'] == 'K'
        assert params['FID_BLNG_CLS_CODE'] == '0'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_params_kosdaq(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        broker.fetch_industry_category_price("Q", "2")
        params = mock_get.call_args[1]['params']
        assert params['FID_INPUT_ISCD'] == '1001'
        assert params['FID_MRKT_CLS_CODE'] == 'Q'
        assert params['FID_BLNG_CLS_CODE'] == '2'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_market_type_options(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        expected = {"K": "0001", "Q": "1001", "K2": "2001"}
        for market_type, expected_code in expected.items():
            broker.fetch_industry_category_price(market_type)
            params = mock_get.call_args[1]['params']
            assert params['FID_INPUT_ISCD'] == expected_code
            assert params['FID_MRKT_CLS_CODE'] == market_type

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_token_refresh_on_expiry(self, mock_get):
        broker = _create_broker_mock()
        mock_get.side_effect = [
            MockResponse({"rt_cd": "1", "msg1": "기간이 만료된 token 입니다"}),
            MockResponse({"rt_cd": "0", "output1": {}, "output2": []})
        ]
        def refresh_token(force=False):
            if force:
                broker.access_token = "Bearer new_token"
        broker.issue_access_token = Mock(side_effect=refresh_token)
        result = broker.fetch_industry_category_price()
        broker.issue_access_token.assert_called_once_with(force=True)
        assert result["rt_cd"] == "0"


class TestCachedDividendIndustryApis:

    def test_dividend_ranking_cache_hit(self):
        from korea_investment_stock import KoreaInvestment
        from korea_investment_stock.cache import CachedKoreaInvestment
        mock_broker = Mock(spec=KoreaInvestment)
        mock_broker.fetch_dividend_ranking.return_value = {'rt_cd': '0', 'output1': []}
        cached = CachedKoreaInvestment(mock_broker, enable_cache=True)
        cached.fetch_dividend_ranking()
        cached.fetch_dividend_ranking()
        assert mock_broker.fetch_dividend_ranking.call_count == 1

    def test_industry_index_cache_disabled(self):
        from korea_investment_stock import KoreaInvestment
        from korea_investment_stock.cache import CachedKoreaInvestment
        mock_broker = Mock(spec=KoreaInvestment)
        mock_broker.fetch_industry_index.return_value = {'rt_cd': '0', 'output': {}}
        cached = CachedKoreaInvestment(mock_broker, enable_cache=False)
        cached.fetch_industry_index()
        cached.fetch_industry_index()
        assert mock_broker.fetch_industry_index.call_count == 2


class TestRateLimitedDividendIndustryApis:

    def test_dividend_ranking_rate_limit(self):
        from korea_investment_stock import KoreaInvestment
        from korea_investment_stock.rate_limit import RateLimitedKoreaInvestment
        mock_broker = Mock(spec=KoreaInvestment)
        mock_broker.fetch_dividend_ranking.return_value = {'rt_cd': '0', 'output1': []}
        rate_limited = RateLimitedKoreaInvestment(mock_broker, calls_per_second=10)
        result = rate_limited.fetch_dividend_ranking()
        assert result['rt_cd'] == '0'
        mock_broker.fetch_dividend_ranking.assert_called_once()

    def test_industry_category_price_rate_limit(self):
        from korea_investment_stock import KoreaInvestment
        from korea_investment_stock.rate_limit import RateLimitedKoreaInvestment
        mock_broker = Mock(spec=KoreaInvestment)
        mock_broker.fetch_industry_category_price.return_value = {'rt_cd': '0', 'output2': []}
        rate_limited = RateLimitedKoreaInvestment(mock_broker, calls_per_second=10)
        result = rate_limited.fetch_industry_category_price("K")
        assert result['rt_cd'] == '0'
