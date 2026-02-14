"""재무제표 API 단위 테스트 (Step 3)"""
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


class TestFetchFinancialRatio:

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_success_response(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({
            'rt_cd': '0',
            'output': [{'stac_yymm': '202312', 'roe_val': '10.5', 'eps': '5000', 'bps': '40000', 'lblt_rate': '30.2'}]
        })
        result = broker.fetch_financial_ratio("005930")
        assert result['rt_cd'] == '0'
        assert result['output'][0]['roe_val'] == '10.5'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_url_and_headers(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        broker.fetch_financial_ratio("005930")
        url = mock_get.call_args[0][0]
        headers = mock_get.call_args[1]['headers']
        assert "financial-ratio" in url
        assert headers['tr_id'] == 'FHKST66430300'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_params_default(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        broker.fetch_financial_ratio("005930")
        params = mock_get.call_args[1]['params']
        assert params['FID_DIV_CLS_CODE'] == '0'
        assert params['fid_cond_mrkt_div_code'] == 'J'
        assert params['fid_input_iscd'] == '005930'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_quarterly_period(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        broker.fetch_financial_ratio("005930", period_type="1")
        params = mock_get.call_args[1]['params']
        assert params['FID_DIV_CLS_CODE'] == '1'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_token_refresh_on_expiry(self, mock_get):
        broker = _create_broker_mock()
        mock_get.side_effect = [
            MockResponse({"rt_cd": "1", "msg1": "기간이 만료된 token 입니다"}),
            MockResponse({"rt_cd": "0", "output": []})
        ]
        def refresh_token(force=False):
            if force:
                broker.access_token = "Bearer new_token"
        broker.issue_access_token = Mock(side_effect=refresh_token)
        result = broker.fetch_financial_ratio("005930")
        broker.issue_access_token.assert_called_once_with(force=True)
        assert result["rt_cd"] == "0"


class TestFetchIncomeStatement:

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_success_response(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({
            'rt_cd': '0',
            'output': [{'stac_yymm': '202312', 'sale_account': '300000000', 'bsop_prti': '50000000', 'thtr_ntin': '40000000'}]
        })
        result = broker.fetch_income_statement("005930")
        assert result['rt_cd'] == '0'
        assert result['output'][0]['sale_account'] == '300000000'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_url_and_headers(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        broker.fetch_income_statement("005930")
        url = mock_get.call_args[0][0]
        headers = mock_get.call_args[1]['headers']
        assert "income-statement" in url
        assert headers['tr_id'] == 'FHKST66430200'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_params(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        broker.fetch_income_statement("000660", "1")
        params = mock_get.call_args[1]['params']
        assert params['FID_DIV_CLS_CODE'] == '1'
        assert params['fid_cond_mrkt_div_code'] == 'J'
        assert params['fid_input_iscd'] == '000660'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_token_refresh_on_expiry(self, mock_get):
        broker = _create_broker_mock()
        mock_get.side_effect = [
            MockResponse({"rt_cd": "1", "msg1": "기간이 만료된 token 입니다"}),
            MockResponse({"rt_cd": "0", "output": []})
        ]
        def refresh_token(force=False):
            if force:
                broker.access_token = "Bearer new_token"
        broker.issue_access_token = Mock(side_effect=refresh_token)
        result = broker.fetch_income_statement("005930")
        broker.issue_access_token.assert_called_once_with(force=True)
        assert result["rt_cd"] == "0"


class TestFetchBalanceSheet:

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_success_response(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({
            'rt_cd': '0',
            'output': [{'stac_yymm': '202312', 'total_aset': '500000000', 'total_lblt': '100000000', 'total_cptl': '400000000'}]
        })
        result = broker.fetch_balance_sheet("005930")
        assert result['rt_cd'] == '0'
        assert result['output'][0]['total_aset'] == '500000000'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_url_and_headers(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        broker.fetch_balance_sheet("005930")
        url = mock_get.call_args[0][0]
        headers = mock_get.call_args[1]['headers']
        assert "balance-sheet" in url
        assert headers['tr_id'] == 'FHKST66430100'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_params(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        broker.fetch_balance_sheet("000660", "1")
        params = mock_get.call_args[1]['params']
        assert params['FID_DIV_CLS_CODE'] == '1'
        assert params['fid_cond_mrkt_div_code'] == 'J'
        assert params['fid_input_iscd'] == '000660'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_token_refresh_on_expiry(self, mock_get):
        broker = _create_broker_mock()
        mock_get.side_effect = [
            MockResponse({"rt_cd": "1", "msg1": "기간이 만료된 token 입니다"}),
            MockResponse({"rt_cd": "0", "output": []})
        ]
        def refresh_token(force=False):
            if force:
                broker.access_token = "Bearer new_token"
        broker.issue_access_token = Mock(side_effect=refresh_token)
        result = broker.fetch_balance_sheet("005930")
        broker.issue_access_token.assert_called_once_with(force=True)
        assert result["rt_cd"] == "0"


class TestFetchProfitabilityRatio:

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_success_response(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({
            'rt_cd': '0',
            'output': [{'stac_yymm': '202312', 'cptl_ntin_rate': '8.5', 'sale_ntin_rate': '15.2', 'sale_totl_rate': '35.0'}]
        })
        result = broker.fetch_profitability_ratio("005930")
        assert result['rt_cd'] == '0'
        assert result['output'][0]['sale_totl_rate'] == '35.0'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_url_and_headers(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        broker.fetch_profitability_ratio("005930")
        url = mock_get.call_args[0][0]
        headers = mock_get.call_args[1]['headers']
        assert "profit-ratio" in url
        assert headers['tr_id'] == 'FHKST66430400'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_params(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        broker.fetch_profitability_ratio("000660", "1")
        params = mock_get.call_args[1]['params']
        assert params['FID_DIV_CLS_CODE'] == '1'
        assert params['fid_cond_mrkt_div_code'] == 'J'
        assert params['fid_input_iscd'] == '000660'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_token_refresh_on_expiry(self, mock_get):
        broker = _create_broker_mock()
        mock_get.side_effect = [
            MockResponse({"rt_cd": "1", "msg1": "기간이 만료된 token 입니다"}),
            MockResponse({"rt_cd": "0", "output": []})
        ]
        def refresh_token(force=False):
            if force:
                broker.access_token = "Bearer new_token"
        broker.issue_access_token = Mock(side_effect=refresh_token)
        result = broker.fetch_profitability_ratio("005930")
        broker.issue_access_token.assert_called_once_with(force=True)
        assert result["rt_cd"] == "0"


class TestFetchGrowthRatio:

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_success_response(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({
            'rt_cd': '0',
            'output': [{'stac_yymm': '202312', 'grs': '12.5', 'bsop_prfi_inrt': '15.0', 'equt_inrt': '8.0', 'totl_aset_inrt': '10.0'}]
        })
        result = broker.fetch_growth_ratio("005930")
        assert result['rt_cd'] == '0'
        assert result['output'][0]['grs'] == '12.5'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_url_and_headers(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        broker.fetch_growth_ratio("005930")
        url = mock_get.call_args[0][0]
        headers = mock_get.call_args[1]['headers']
        assert "growth-ratio" in url
        assert headers['tr_id'] == 'FHKST66430800'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_params_lowercase_key(self, mock_get):
        """성장성비율 API는 fid_div_cls_code (소문자) 사용"""
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        broker.fetch_growth_ratio("000660", "1")
        params = mock_get.call_args[1]['params']
        assert params['fid_div_cls_code'] == '1'
        assert params['fid_cond_mrkt_div_code'] == 'J'
        assert params['fid_input_iscd'] == '000660'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_token_refresh_on_expiry(self, mock_get):
        broker = _create_broker_mock()
        mock_get.side_effect = [
            MockResponse({"rt_cd": "1", "msg1": "기간이 만료된 token 입니다"}),
            MockResponse({"rt_cd": "0", "output": []})
        ]
        def refresh_token(force=False):
            if force:
                broker.access_token = "Bearer new_token"
        broker.issue_access_token = Mock(side_effect=refresh_token)
        result = broker.fetch_growth_ratio("005930")
        broker.issue_access_token.assert_called_once_with(force=True)
        assert result["rt_cd"] == "0"


class TestCachedFinancialApis:

    def test_financial_ratio_cache_hit(self):
        from korea_investment_stock import KoreaInvestment
        from korea_investment_stock.cache import CachedKoreaInvestment
        mock_broker = Mock(spec=KoreaInvestment)
        mock_broker.fetch_financial_ratio.return_value = {'rt_cd': '0', 'output': []}
        cached = CachedKoreaInvestment(mock_broker, enable_cache=True)
        cached.fetch_financial_ratio("005930")
        cached.fetch_financial_ratio("005930")
        assert mock_broker.fetch_financial_ratio.call_count == 1

    def test_balance_sheet_cache_disabled(self):
        from korea_investment_stock import KoreaInvestment
        from korea_investment_stock.cache import CachedKoreaInvestment
        mock_broker = Mock(spec=KoreaInvestment)
        mock_broker.fetch_balance_sheet.return_value = {'rt_cd': '0', 'output': []}
        cached = CachedKoreaInvestment(mock_broker, enable_cache=False)
        cached.fetch_balance_sheet("005930")
        cached.fetch_balance_sheet("005930")
        assert mock_broker.fetch_balance_sheet.call_count == 2


class TestRateLimitedFinancialApis:

    def test_financial_ratio_rate_limit(self):
        from korea_investment_stock import KoreaInvestment
        from korea_investment_stock.rate_limit import RateLimitedKoreaInvestment
        mock_broker = Mock(spec=KoreaInvestment)
        mock_broker.fetch_financial_ratio.return_value = {'rt_cd': '0', 'output': []}
        rate_limited = RateLimitedKoreaInvestment(mock_broker, calls_per_second=10)
        result = rate_limited.fetch_financial_ratio("005930")
        assert result['rt_cd'] == '0'
        mock_broker.fetch_financial_ratio.assert_called_once()

    def test_growth_ratio_rate_limit(self):
        from korea_investment_stock import KoreaInvestment
        from korea_investment_stock.rate_limit import RateLimitedKoreaInvestment
        mock_broker = Mock(spec=KoreaInvestment)
        mock_broker.fetch_growth_ratio.return_value = {'rt_cd': '0', 'output': []}
        rate_limited = RateLimitedKoreaInvestment(mock_broker, calls_per_second=10)
        result = rate_limited.fetch_growth_ratio("005930", "1")
        assert result['rt_cd'] == '0'
