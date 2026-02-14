"""Phase 1 API 확장 통합 테스트 (15개 API)

실제 API 자격 증명이 필요합니다.
환경 변수 설정:
    - KOREA_INVESTMENT_API_KEY
    - KOREA_INVESTMENT_API_SECRET
    - KOREA_INVESTMENT_ACCOUNT_NO

실행:
    pytest korea_investment_stock/tests/test_integration_apis.py -v
"""
import pytest
import time
from datetime import datetime, timedelta


@pytest.fixture(scope="module")
def broker():
    """실제 API 자격 증명으로 broker 생성 (모듈 단위 재사용)"""
    import os
    from korea_investment_stock import KoreaInvestment

    # Redis가 없는 환경에서도 동작하도록 file 토큰 저장소 사용
    orig = os.environ.get("KOREA_INVESTMENT_TOKEN_STORAGE")
    os.environ["KOREA_INVESTMENT_TOKEN_STORAGE"] = "file"
    try:
        b = KoreaInvestment()
    finally:
        if orig is not None:
            os.environ["KOREA_INVESTMENT_TOKEN_STORAGE"] = orig
        else:
            os.environ.pop("KOREA_INVESTMENT_TOKEN_STORAGE", None)
    return b


def _wait():
    """API 호출 간 rate limit 대기"""
    time.sleep(0.1)


# === Step 1: 차트 데이터 API ===

@pytest.mark.integration
class TestDomesticChartIntegration:

    def test_samsung_daily_chart(self, broker):
        """삼성전자 일봉 조회"""
        end_date = datetime.now().strftime("%Y%m%d")
        start_date = (datetime.now() - timedelta(days=30)).strftime("%Y%m%d")

        result = broker.fetch_domestic_chart("005930", "D", start_date, end_date)
        _wait()

        assert result['rt_cd'] == '0', f"API Error: {result.get('msg1')}"
        assert 'output2' in result
        assert len(result['output2']) > 0

    def test_daily_chart_response_fields(self, broker):
        """일봉 응답 필드 검증"""
        end_date = datetime.now().strftime("%Y%m%d")
        start_date = (datetime.now() - timedelta(days=7)).strftime("%Y%m%d")

        result = broker.fetch_domestic_chart("005930", "D", start_date, end_date)
        _wait()

        if result['rt_cd'] == '0' and result.get('output2'):
            candle = result['output2'][0]
            expected_fields = ['stck_bsop_date', 'stck_clpr', 'stck_oprc', 'stck_hgpr', 'stck_lwpr', 'acml_vol']
            for field in expected_fields:
                assert field in candle, f"Missing field: {field}"

    def test_weekly_chart(self, broker):
        """주봉 조회"""
        end_date = datetime.now().strftime("%Y%m%d")
        start_date = (datetime.now() - timedelta(days=90)).strftime("%Y%m%d")

        result = broker.fetch_domestic_chart("005930", "W", start_date, end_date)
        _wait()

        assert result['rt_cd'] == '0', f"API Error: {result.get('msg1')}"


@pytest.mark.integration
class TestDomesticMinuteChartIntegration:

    def test_samsung_minute_chart(self, broker):
        """삼성전자 분봉 조회"""
        result = broker.fetch_domestic_minute_chart("005930")
        _wait()

        assert result['rt_cd'] == '0', f"API Error: {result.get('msg1')}"
        assert 'output2' in result

    def test_minute_chart_response_fields(self, broker):
        """분봉 응답 필드 검증"""
        result = broker.fetch_domestic_minute_chart("005930")
        _wait()

        if result['rt_cd'] == '0' and result.get('output2'):
            candle = result['output2'][0]
            expected_fields = ['stck_cntg_hour', 'stck_prpr', 'cntg_vol']
            for field in expected_fields:
                assert field in candle, f"Missing field: {field}"


@pytest.mark.integration
class TestOverseasChartIntegration:

    def test_apple_daily_chart(self, broker):
        """AAPL 일봉 조회"""
        result = broker.fetch_overseas_chart("AAPL", "US", "D")
        _wait()

        assert result['rt_cd'] == '0', f"API Error: {result.get('msg1')}"
        assert 'output2' in result
        # 해외 시장은 거래 시간 외에 빈 배열 반환 가능
        assert isinstance(result['output2'], list)

    def test_overseas_chart_response_fields(self, broker):
        """해외 차트 응답 필드 검증"""
        result = broker.fetch_overseas_chart("TSLA", "US", "D")
        _wait()

        if result['rt_cd'] == '0' and result.get('output2'):
            candle = result['output2'][0]
            expected_fields = ['xymd', 'clos', 'open', 'high', 'low', 'tvol']
            for field in expected_fields:
                assert field in candle, f"Missing field: {field}"


# === Step 2: 시세 순위 API ===

@pytest.mark.integration
class TestVolumeRankingIntegration:

    def test_volume_ranking(self, broker):
        """거래량 순위 조회"""
        result = broker.fetch_volume_ranking()
        _wait()

        assert result['rt_cd'] == '0', f"API Error: {result.get('msg1')}"
        assert 'output' in result
        assert len(result['output']) > 0

    def test_volume_ranking_response_fields(self, broker):
        """거래량 순위 응답 필드 검증"""
        result = broker.fetch_volume_ranking()
        _wait()

        if result['rt_cd'] == '0' and result.get('output'):
            stock = result['output'][0]
            expected_fields = ['data_rank', 'hts_kor_isnm', 'mksc_shrn_iscd', 'stck_prpr', 'acml_vol']
            for field in expected_fields:
                assert field in stock, f"Missing field: {field}"


@pytest.mark.integration
class TestChangeRateRankingIntegration:

    def test_rise_ranking(self, broker):
        """상승률 순위 조회"""
        result = broker.fetch_change_rate_ranking("J", "0")
        _wait()

        assert result['rt_cd'] == '0', f"API Error: {result.get('msg1')}"
        assert 'output' in result

    def test_fall_ranking(self, broker):
        """하락률 순위 조회"""
        result = broker.fetch_change_rate_ranking("J", "1")
        _wait()

        assert result['rt_cd'] == '0', f"API Error: {result.get('msg1')}"


@pytest.mark.integration
class TestMarketCapRankingIntegration:

    def test_market_cap_ranking(self, broker):
        """시가총액 상위 조회"""
        result = broker.fetch_market_cap_ranking()
        _wait()

        assert result['rt_cd'] == '0', f"API Error: {result.get('msg1')}"
        assert 'output' in result
        assert len(result['output']) > 0

    def test_market_cap_response_fields(self, broker):
        """시가총액 응답 필드 검증"""
        result = broker.fetch_market_cap_ranking()
        _wait()

        if result['rt_cd'] == '0' and result.get('output'):
            stock = result['output'][0]
            expected_fields = ['data_rank', 'hts_kor_isnm', 'stck_prpr', 'stck_avls']
            for field in expected_fields:
                assert field in stock, f"Missing field: {field}"


@pytest.mark.integration
class TestOverseasChangeRateRankingIntegration:

    def test_us_rise_ranking(self, broker):
        """미국 상승률 순위 조회"""
        result = broker.fetch_overseas_change_rate_ranking("US", "1")
        _wait()

        assert result['rt_cd'] == '0', f"API Error: {result.get('msg1')}"
        assert 'output2' in result

    def test_overseas_ranking_response_fields(self, broker):
        """해외 순위 응답 필드 검증"""
        result = broker.fetch_overseas_change_rate_ranking("US", "1")
        _wait()

        if result['rt_cd'] == '0' and result.get('output2'):
            stock = result['output2'][0]
            expected_fields = ['rank', 'name', 'symb', 'last', 'rate']
            for field in expected_fields:
                assert field in stock, f"Missing field: {field}"


# === Step 3: 재무제표 API ===

@pytest.mark.integration
class TestFinancialRatioIntegration:

    def test_samsung_financial_ratio_annual(self, broker):
        """삼성전자 재무비율 (연간)"""
        result = broker.fetch_financial_ratio("005930", "0")
        _wait()

        assert result['rt_cd'] == '0', f"API Error: {result.get('msg1')}"
        assert 'output' in result
        assert len(result['output']) > 0

    def test_financial_ratio_response_fields(self, broker):
        """재무비율 응답 필드 검증"""
        result = broker.fetch_financial_ratio("005930", "0")
        _wait()

        if result['rt_cd'] == '0' and result.get('output'):
            row = result['output'][0]
            expected_fields = ['stac_yymm', 'roe_val', 'eps', 'bps', 'lblt_rate']
            for field in expected_fields:
                assert field in row, f"Missing field: {field}"

    def test_financial_ratio_quarterly(self, broker):
        """재무비율 (분기)"""
        result = broker.fetch_financial_ratio("005930", "1")
        _wait()

        assert result['rt_cd'] == '0', f"API Error: {result.get('msg1')}"


@pytest.mark.integration
class TestIncomeStatementIntegration:

    def test_samsung_income_statement(self, broker):
        """삼성전자 손익계산서"""
        result = broker.fetch_income_statement("005930", "0")
        _wait()

        assert result['rt_cd'] == '0', f"API Error: {result.get('msg1')}"
        assert 'output' in result
        assert len(result['output']) > 0

    def test_income_statement_response_fields(self, broker):
        """손익계산서 응답 필드 검증"""
        result = broker.fetch_income_statement("005930", "0")
        _wait()

        if result['rt_cd'] == '0' and result.get('output'):
            row = result['output'][0]
            expected_fields = ['stac_yymm', 'sale_account', 'bsop_prti', 'thtr_ntin']
            for field in expected_fields:
                assert field in row, f"Missing field: {field}"


@pytest.mark.integration
class TestBalanceSheetIntegration:

    def test_samsung_balance_sheet(self, broker):
        """삼성전자 대차대조표"""
        result = broker.fetch_balance_sheet("005930", "0")
        _wait()

        assert result['rt_cd'] == '0', f"API Error: {result.get('msg1')}"
        assert 'output' in result
        assert len(result['output']) > 0

    def test_balance_sheet_response_fields(self, broker):
        """대차대조표 응답 필드 검증"""
        result = broker.fetch_balance_sheet("005930", "0")
        _wait()

        if result['rt_cd'] == '0' and result.get('output'):
            row = result['output'][0]
            expected_fields = ['stac_yymm', 'total_aset', 'total_lblt', 'total_cptl']
            for field in expected_fields:
                assert field in row, f"Missing field: {field}"


@pytest.mark.integration
class TestProfitabilityRatioIntegration:

    def test_samsung_profitability_ratio(self, broker):
        """삼성전자 수익성비율"""
        result = broker.fetch_profitability_ratio("005930", "0")
        _wait()

        assert result['rt_cd'] == '0', f"API Error: {result.get('msg1')}"
        assert 'output' in result

    def test_profitability_ratio_response_fields(self, broker):
        """수익성비율 응답 필드 검증"""
        result = broker.fetch_profitability_ratio("005930", "0")
        _wait()

        if result['rt_cd'] == '0' and result.get('output'):
            row = result['output'][0]
            expected_fields = ['stac_yymm', 'cptl_ntin_rate', 'sale_ntin_rate', 'sale_totl_rate']
            for field in expected_fields:
                assert field in row, f"Missing field: {field}"


@pytest.mark.integration
class TestGrowthRatioIntegration:

    def test_samsung_growth_ratio(self, broker):
        """삼성전자 성장성비율"""
        result = broker.fetch_growth_ratio("005930", "0")
        _wait()

        assert result['rt_cd'] == '0', f"API Error: {result.get('msg1')}"
        assert 'output' in result

    def test_growth_ratio_response_fields(self, broker):
        """성장성비율 응답 필드 검증"""
        result = broker.fetch_growth_ratio("005930", "0")
        _wait()

        if result['rt_cd'] == '0' and result.get('output'):
            row = result['output'][0]
            expected_fields = ['stac_yymm', 'grs', 'bsop_prfi_inrt']
            for field in expected_fields:
                assert field in row, f"Missing field: {field}"


# === Step 4: 배당 + 업종 API ===

@pytest.mark.integration
class TestDividendRankingIntegration:

    def test_dividend_ranking(self, broker):
        """배당률 상위 조회"""
        end_date = datetime.now().strftime("%Y%m%d")
        start_date = (datetime.now() - timedelta(days=365)).strftime("%Y%m%d")

        result = broker.fetch_dividend_ranking("0", "2", start_date, end_date)
        _wait()

        assert result['rt_cd'] == '0', f"API Error: {result.get('msg1')}"

    def test_dividend_ranking_response_fields(self, broker):
        """배당률 응답 필드 검증"""
        end_date = datetime.now().strftime("%Y%m%d")
        start_date = (datetime.now() - timedelta(days=365)).strftime("%Y%m%d")

        result = broker.fetch_dividend_ranking("0", "2", start_date, end_date)
        _wait()

        if result['rt_cd'] == '0' and result.get('output1'):
            row = result['output1'][0]
            expected_fields = ['rank', 'isin_name', 'divi_rate']
            for field in expected_fields:
                assert field in row, f"Missing field: {field}"


@pytest.mark.integration
class TestIndustryIndexIntegration:

    def test_kospi_index(self, broker):
        """코스피 종합 지수 조회"""
        result = broker.fetch_industry_index("0001")
        _wait()

        assert result['rt_cd'] == '0', f"API Error: {result.get('msg1')}"
        assert 'output' in result

    def test_industry_index_response_fields(self, broker):
        """업종 지수 응답 필드 검증"""
        result = broker.fetch_industry_index("0001")
        _wait()

        if result['rt_cd'] == '0' and result.get('output'):
            output = result['output']
            expected_fields = ['bstp_nmix_prpr', 'bstp_nmix_prdy_vrss', 'bstp_nmix_prdy_ctrt']
            for field in expected_fields:
                assert field in output, f"Missing field: {field}"

    def test_kosdaq_index(self, broker):
        """코스닥 종합 지수 조회"""
        result = broker.fetch_industry_index("1001")
        _wait()

        assert result['rt_cd'] == '0', f"API Error: {result.get('msg1')}"


@pytest.mark.integration
class TestIndustryCategoryPriceIntegration:

    def test_kospi_category_price(self, broker):
        """거래소 업종별 전체 시세"""
        result = broker.fetch_industry_category_price("K", "0")
        _wait()

        assert result['rt_cd'] == '0', f"API Error: {result.get('msg1')}"
        assert 'output2' in result

    def test_category_price_response_fields(self, broker):
        """업종별 시세 응답 필드 검증"""
        result = broker.fetch_industry_category_price("K", "0")
        _wait()

        if result['rt_cd'] == '0' and result.get('output2'):
            sector = result['output2'][0]
            expected_fields = ['bstp_cls_code', 'hts_kor_isnm', 'bstp_nmix_prpr']
            for field in expected_fields:
                assert field in sector, f"Missing field: {field}"

    def test_kosdaq_category_price(self, broker):
        """코스닥 업종별 전체 시세"""
        result = broker.fetch_industry_category_price("Q", "0")
        _wait()

        assert result['rt_cd'] == '0', f"API Error: {result.get('msg1')}"
