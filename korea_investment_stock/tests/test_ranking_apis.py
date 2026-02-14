"""시세 순위 API 단위 테스트 (Step 2)"""
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


class TestFetchVolumeRanking:

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_success_response(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({
            'rt_cd': '0',
            'output': [{'data_rank': '1', 'hts_kor_isnm': '삼성전자', 'mksc_shrn_iscd': '005930', 'stck_prpr': '70000', 'acml_vol': '50000000'}]
        })
        result = broker.fetch_volume_ranking()
        assert result['rt_cd'] == '0'
        assert result['output'][0]['data_rank'] == '1'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_url_and_headers(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        broker.fetch_volume_ranking()
        url = mock_get.call_args[0][0]
        headers = mock_get.call_args[1]['headers']
        assert "volume-rank" in url
        assert headers['tr_id'] == 'FHPST01710000'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_params_default(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        broker.fetch_volume_ranking()
        params = mock_get.call_args[1]['params']
        assert params['FID_COND_MRKT_DIV_CODE'] == 'J'
        assert params['FID_COND_SCR_DIV_CODE'] == '20171'
        assert params['FID_INPUT_ISCD'] == '0000'
        assert params['FID_BLNG_CLS_CODE'] == '0'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_sort_by_options(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        for sort_by in ["0", "1", "2", "3", "4"]:
            broker.fetch_volume_ranking(sort_by=sort_by)
            params = mock_get.call_args[1]['params']
            assert params['FID_BLNG_CLS_CODE'] == sort_by

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
        result = broker.fetch_volume_ranking()
        broker.issue_access_token.assert_called_once_with(force=True)
        assert result["rt_cd"] == "0"


class TestFetchChangeRateRanking:

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_success_response(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({
            'rt_cd': '0',
            'output': [{'data_rank': '1', 'hts_kor_isnm': '테스트종목', 'stck_prpr': '10000', 'prdy_ctrt': '29.90'}]
        })
        result = broker.fetch_change_rate_ranking()
        assert result['rt_cd'] == '0'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_url_and_headers(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        broker.fetch_change_rate_ranking()
        url = mock_get.call_args[0][0]
        headers = mock_get.call_args[1]['headers']
        assert "fluctuation" in url
        assert headers['tr_id'] == 'FHPST01700000'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_params(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        broker.fetch_change_rate_ranking("J", "1", "3")
        params = mock_get.call_args[1]['params']
        assert params['fid_cond_mrkt_div_code'] == 'J'
        assert params['fid_cond_scr_div_code'] == '20170'
        assert params['fid_rank_sort_cls_code'] == '1'
        assert params['fid_input_cnt_1'] == '3'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_sort_order_options(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        for sort_order in ["0", "1", "2", "3", "4"]:
            broker.fetch_change_rate_ranking(sort_order=sort_order)
            params = mock_get.call_args[1]['params']
            assert params['fid_rank_sort_cls_code'] == sort_order


class TestFetchMarketCapRanking:

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_success_response(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({
            'rt_cd': '0',
            'output': [{'data_rank': '1', 'hts_kor_isnm': '삼성전자', 'stck_avls': '5000000', 'mrkt_whol_avls_rlim': '20.5'}]
        })
        result = broker.fetch_market_cap_ranking()
        assert result['rt_cd'] == '0'
        assert result['output'][0]['stck_avls'] == '5000000'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_url_and_headers(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        broker.fetch_market_cap_ranking()
        url = mock_get.call_args[0][0]
        headers = mock_get.call_args[1]['headers']
        assert "market-cap" in url
        assert headers['tr_id'] == 'FHPST01740000'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_params(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        broker.fetch_market_cap_ranking("J", "1001")
        params = mock_get.call_args[1]['params']
        assert params['fid_cond_mrkt_div_code'] == 'J'
        assert params['fid_cond_scr_div_code'] == '20174'
        assert params['fid_input_iscd'] == '1001'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_target_market_options(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        for target in ["0000", "0001", "1001", "2001"]:
            broker.fetch_market_cap_ranking(target_market=target)
            params = mock_get.call_args[1]['params']
            assert params['fid_input_iscd'] == target


class TestFetchOverseasChangeRateRanking:

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_success_response(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({
            'rt_cd': '0',
            'output1': {'nrec': '30'},
            'output2': [{'rank': '1', 'name': 'TESLA', 'symb': 'TSLA', 'last': '250.00', 'rate': '5.20', 'tvol': '100000000'}]
        })
        result = broker.fetch_overseas_change_rate_ranking()
        assert result['rt_cd'] == '0'
        assert result['output2'][0]['rank'] == '1'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_url_and_headers(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        broker.fetch_overseas_change_rate_ranking()
        url = mock_get.call_args[0][0]
        headers = mock_get.call_args[1]['headers']
        assert "updown-rate" in url
        assert headers['tr_id'] == 'HHDFS76290000'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_params_us(self, mock_get):
        broker = _create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0'})
        broker.fetch_overseas_change_rate_ranking("US", "1", "0", "0")
        params = mock_get.call_args[1]['params']
        assert params['AUTH'] == ''
        assert params['EXCD'] == 'NYS'
        assert params['GUBN'] == '1'
        assert params['NDAY'] == '0'
        assert params['VOL_RANG'] == '0'

    def test_invalid_country_code(self):
        broker = _create_broker_mock()
        with pytest.raises(ValueError, match="지원하지 않는 country_code"):
            broker.fetch_overseas_change_rate_ranking("XX")

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
        result = broker.fetch_overseas_change_rate_ranking()
        broker.issue_access_token.assert_called_once_with(force=True)
        assert result["rt_cd"] == "0"


class TestCachedRankingApis:

    def test_volume_ranking_cache_hit(self):
        from korea_investment_stock import KoreaInvestment
        from korea_investment_stock.cache import CachedKoreaInvestment
        mock_broker = Mock(spec=KoreaInvestment)
        mock_broker.fetch_volume_ranking.return_value = {'rt_cd': '0', 'output': []}
        cached = CachedKoreaInvestment(mock_broker, enable_cache=True)
        cached.fetch_volume_ranking()
        cached.fetch_volume_ranking()
        assert mock_broker.fetch_volume_ranking.call_count == 1

    def test_market_cap_ranking_cache_disabled(self):
        from korea_investment_stock import KoreaInvestment
        from korea_investment_stock.cache import CachedKoreaInvestment
        mock_broker = Mock(spec=KoreaInvestment)
        mock_broker.fetch_market_cap_ranking.return_value = {'rt_cd': '0', 'output': []}
        cached = CachedKoreaInvestment(mock_broker, enable_cache=False)
        cached.fetch_market_cap_ranking()
        cached.fetch_market_cap_ranking()
        assert mock_broker.fetch_market_cap_ranking.call_count == 2


class TestRateLimitedRankingApis:

    def test_volume_ranking_rate_limit(self):
        from korea_investment_stock import KoreaInvestment
        from korea_investment_stock.rate_limit import RateLimitedKoreaInvestment
        mock_broker = Mock(spec=KoreaInvestment)
        mock_broker.fetch_volume_ranking.return_value = {'rt_cd': '0', 'output': []}
        rate_limited = RateLimitedKoreaInvestment(mock_broker, calls_per_second=10)
        result = rate_limited.fetch_volume_ranking()
        assert result['rt_cd'] == '0'
        mock_broker.fetch_volume_ranking.assert_called_once()

    def test_overseas_ranking_rate_limit(self):
        from korea_investment_stock import KoreaInvestment
        from korea_investment_stock.rate_limit import RateLimitedKoreaInvestment
        mock_broker = Mock(spec=KoreaInvestment)
        mock_broker.fetch_overseas_change_rate_ranking.return_value = {'rt_cd': '0', 'output2': []}
        rate_limited = RateLimitedKoreaInvestment(mock_broker, calls_per_second=10)
        result = rate_limited.fetch_overseas_change_rate_ranking("US", "1")
        assert result['rt_cd'] == '0'
