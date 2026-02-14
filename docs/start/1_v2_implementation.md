# Phase 1: korea-investment-stock API 확장 구현 문서

## 1. 코드 패턴

### 1.1 메인 클래스 메서드 패턴 (`korea_investment_stock.py`)

모든 API 메서드는 동일한 패턴을 따른다:

```python
def fetch_xxx(self, param1: str, param2: str = "default") -> dict:
    """메서드 설명"""
    path = "uapi/domestic-stock/v1/quotations/xxx"
    url = f"{self.base_url}/{path}"
    headers = {
        "content-type": "application/json",
        "authorization": self.access_token,
        "appKey": self.api_key,
        "appSecret": self.api_secret,
        "tr_id": "TR_ID_HERE"
    }
    params = {
        "PARAM_KEY": param1,
    }
    return self._request_with_token_refresh("GET", url, headers, params)
```

### 1.2 캐시 래퍼 패턴 (`cached_korea_investment.py`)

```python
def fetch_xxx(self, param1: str, param2: str = "default") -> dict:
    if not self.enable_cache:
        return self.broker.fetch_xxx(param1, param2)

    cache_key = self._make_cache_key("fetch_xxx", param1, param2)
    cached_data = self.cache.get(cache_key)
    if cached_data is not None:
        return cached_data

    result = self.broker.fetch_xxx(param1, param2)
    if result.get('rt_cd') == '0':
        self.cache.set(cache_key, result, self.ttl['price'])  # TTL 카테고리 선택
    return result
```

**TTL 카테고리 매핑**:

| API 카테고리 | TTL 키 | 기본값 |
|------------|--------|-------|
| 차트/시세/순위 | `price` | 5초 |
| 재무제표/종목정보 | `stock_info` | 300초 |
| 업종/배당 | `stock_info` | 300초 |

### 1.3 Rate Limit 래퍼 패턴 (`rate_limited_korea_investment.py`)

```python
def fetch_xxx(self, param1: str, param2: str = "default") -> Dict[str, Any]:
    self._rate_limiter.wait()
    return self._broker.fetch_xxx(param1, param2)
```

### 1.4 단위 테스트 패턴 (`tests/test_*.py`)

```python
class TestFetchXxx:
    def _create_broker_mock(self):
        from korea_investment_stock import KoreaInvestment
        with patch.object(KoreaInvestment, '__init__', lambda x: None):
            broker = KoreaInvestment.__new__(KoreaInvestment)
            broker.base_url = "https://openapi.koreainvestment.com:9443"
            broker.access_token = "Bearer test_token"
            broker.api_key = "test_api_key"
            broker.api_secret = "test_api_secret"
            broker._token_manager = Mock()
            return broker

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_success_response(self, mock_get):
        broker = self._create_broker_mock()
        mock_get.return_value = MockResponse({'rt_cd': '0', 'output': {...}})
        result = broker.fetch_xxx("param")
        assert result['rt_cd'] == '0'

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_url_and_headers(self, mock_get):
        # URL 경로, TR ID 검증

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_request_params(self, mock_get):
        # 파라미터 매핑 검증

    @patch('korea_investment_stock.korea_investment_stock.requests.get')
    def test_token_refresh_on_expiry(self, mock_get):
        # 토큰 만료 자동 재발급 검증
```

---

## 2. Step 1: 차트 데이터 API (3개)

### 2.1 `fetch_domestic_chart`

| 항목 | 값 |
|------|-----|
| 경로 | `/uapi/domestic-stock/v1/quotations/inquire-daily-itemchartprice` |
| TR ID | `FHKST03010100` |
| 문서 | `docs/api/국내주식/국내주식기간별시세(일_주_월_년).md` |

**파라미터 매핑**:

```python
params = {
    "FID_COND_MRKT_DIV_CODE": market_code,  # J, NX, UN
    "FID_INPUT_ISCD": symbol,
    "FID_INPUT_DATE_1": start_date,
    "FID_INPUT_DATE_2": end_date,
    "FID_PERIOD_DIV_CODE": period,          # D, W, M, Y
    "FID_ORG_ADJ_PRC": "0" if adjusted else "1",
}
```

### 2.2 `fetch_domestic_minute_chart`

| 항목 | 값 |
|------|-----|
| 경로 | `/uapi/domestic-stock/v1/quotations/inquire-time-itemchartprice` |
| TR ID | `FHKST03010200` |
| 문서 | `docs/api/국내주식/주식당일분봉조회.md` |

**파라미터 매핑**:

```python
params = {
    "FID_COND_MRKT_DIV_CODE": market_code,  # J
    "FID_INPUT_ISCD": symbol,
    "FID_INPUT_HOUR_1": time_from,          # HHMMSS
    "FID_PW_DATA_INCU_YN": "N",            # 고정
    "FID_ETC_CLS_CODE": "",                 # 고정
}
```

### 2.3 `fetch_overseas_chart`

| 항목 | 값 |
|------|-----|
| 경로 | `/uapi/overseas-price/v1/quotations/dailyprice` |
| TR ID | `HHDFS76240000` |
| 문서 | `docs/api/해외주식/해외주식_기간별시세.md` |

**파라미터 매핑**:

```python
# period 변환: D→0, W→1, M→2
period_map = {"D": "0", "W": "1", "M": "2"}

# country_code → EXCD 변환 (EXCD_BY_COUNTRY 사용, 첫 번째 값)
params = {
    "AUTH": "",
    "EXCD": EXCD_BY_COUNTRY[country_code][0],  # NYS, HKS 등
    "SYMB": symbol,
    "GUBN": period_map.get(period, "0"),
    "BYMD": end_date,
    "MODP": "1" if adjusted else "0",
}
```

---

## 3. Step 2: 시세 순위 API (4개)

### 3.1 `fetch_volume_ranking`

| 항목 | 값 |
|------|-----|
| 경로 | `/uapi/domestic-stock/v1/quotations/volume-rank` |
| TR ID | `FHPST01710000` |
| 문서 | `docs/api/국내주식/거래량순위.md` |

**파라미터 매핑**:

```python
params = {
    "FID_COND_MRKT_DIV_CODE": market_code,    # J, NX
    "FID_COND_SCR_DIV_CODE": "20171",         # 고정
    "FID_INPUT_ISCD": "0000",                 # 전체
    "FID_DIV_CLS_CODE": "0",                  # 전체
    "FID_BLNG_CLS_CODE": sort_by,             # 0~4
    "FID_TRGT_CLS_CODE": "",
    "FID_TRGT_EXLS_CLS_CODE": "",
    "FID_INPUT_PRICE_1": "",
    "FID_INPUT_PRICE_2": "",
    "FID_VOL_CNT": "",
    "FID_INPUT_DATE_1": "",
}
```

- 고정값 파라미터가 많으므로 API 문서를 참조하여 정확한 기본값 설정 필요

### 3.2 `fetch_change_rate_ranking`

| 항목 | 값 |
|------|-----|
| 경로 | `/uapi/domestic-stock/v1/ranking/fluctuation` |
| TR ID | `FHPST01700000` |
| 문서 | `docs/api/국내주식/국내주식_등락률_순위.md` |

**파라미터 매핑**:

```python
params = {
    "fid_cond_mrkt_div_code": market_code,
    "fid_cond_scr_div_code": "20170",         # 고정
    "fid_input_iscd": "0000",                 # 전체
    "fid_rank_sort_cls_code": sort_order,     # 0~4
    "fid_input_cnt_1": period_days,           # 0~N
    "fid_prc_cls_code": "0",                  # 고정
    "fid_input_price_1": "",
    "fid_input_price_2": "",
    "fid_vol_cnt": "",
    "fid_trgt_cls_code": "0",
    "fid_trgt_exls_cls_code": "0",
    "fid_div_cls_code": "0",
    "fid_rsfl_rate1": "",
    "fid_rsfl_rate2": "",
}
```

### 3.3 `fetch_market_cap_ranking`

| 항목 | 값 |
|------|-----|
| 경로 | `/uapi/domestic-stock/v1/ranking/market-cap` |
| TR ID | `FHPST01740000` |
| 문서 | `docs/api/국내주식/국내주식_시가총액_상위.md` |

**파라미터 매핑**:

```python
params = {
    "fid_cond_mrkt_div_code": market_code,    # J
    "fid_cond_scr_div_code": "20174",         # 고정
    "fid_input_iscd": target_market,          # 0000, 0001, 1001, 2001
    "fid_div_cls_code": "0",                  # 고정
    "fid_trgt_cls_code": "0",
    "fid_trgt_exls_cls_code": "0",
    "fid_input_price_1": "",
    "fid_input_price_2": "",
    "fid_vol_cnt": "",
}
```

### 3.4 `fetch_overseas_change_rate_ranking`

| 항목 | 값 |
|------|-----|
| 경로 | `/uapi/overseas-stock/v1/ranking/updown-rate` |
| TR ID | `HHDFS76290000` |
| 문서 | `docs/api/해외주식/해외주식_상승율_하락율.md` |

**파라미터 매핑**:

```python
# country_code → EXCD 변환
params = {
    "AUTH": "",
    "EXCD": EXCD_BY_COUNTRY[country_code][0],
    "GUBN": sort_order,      # 0:하락률, 1:상승률
    "NDAY": period,           # 0~9
    "VOL_RANG": volume_filter, # 0~9
}
```

---

## 4. Step 3: 재무제표 API (5개)

5개 모두 동일한 파라미터 구조:

```python
params = {
    "fid_div_cls_code" 또는 "FID_DIV_CLS_CODE": period_type,  # 0:연간, 1:분기
    "fid_cond_mrkt_div_code": "J",    # 고정
    "fid_input_iscd": symbol,
}
```

### 4.1 `fetch_financial_ratio`

| 항목 | 값 |
|------|-----|
| 경로 | `/uapi/domestic-stock/v1/finance/financial-ratio` |
| TR ID | `FHKST66430300` |
| 문서 | `docs/api/국내주식/국내주식_재무비율.md` |

### 4.2 `fetch_income_statement`

| 항목 | 값 |
|------|-----|
| 경로 | `/uapi/domestic-stock/v1/finance/income-statement` |
| TR ID | `FHKST66430200` |
| 문서 | `docs/api/국내주식/국내주식_손익계산서.md` |

### 4.3 `fetch_balance_sheet`

| 항목 | 값 |
|------|-----|
| 경로 | `/uapi/domestic-stock/v1/finance/balance-sheet` |
| TR ID | `FHKST66430100` |
| 문서 | `docs/api/국내주식/국내주식_대차대조표.md` |

### 4.4 `fetch_profitability_ratio`

| 항목 | 값 |
|------|-----|
| 경로 | `/uapi/domestic-stock/v1/finance/profit-ratio` |
| TR ID | `FHKST66430400` |
| 문서 | `docs/api/국내주식/국내주식_수익성비율.md` |

### 4.5 `fetch_growth_ratio`

| 항목 | 값 |
|------|-----|
| 경로 | `/uapi/domestic-stock/v1/finance/growth-ratio` |
| TR ID | `FHKST66430800` |
| 문서 | `docs/api/국내주식/국내주식_성장성비율.md` |

---

## 5. Step 4: 배당 + 업종 API (3~4개)

### 5.1 `fetch_dividend_ranking`

| 항목 | 값 |
|------|-----|
| 경로 | `/uapi/domestic-stock/v1/ranking/dividend-rate` |
| TR ID | `HHKDB13470100` |
| 문서 | `docs/api/국내주식/국내주식_배당률_상위.md` |

**파라미터 매핑**:

```python
params = {
    "CTS": "",
    "GB1": market_type,       # 0:전체, 1:KOSPI, 2:KOSPI200, 3:KOSDAQ
    "UPJONG": "",             # 업종코드 (빈값=전체)
    "GB2": "0",               # 전체
    "GB3": dividend_type,     # 1:주식, 2:현금
    "F_DT": start_date,
    "T_DT": end_date,
}
```

### 5.2 `fetch_industry_index`

| 항목 | 값 |
|------|-----|
| 경로 | `/uapi/domestic-stock/v1/quotations/inquire-index-price` |
| TR ID | `FHPUP02100000` |
| 문서 | `docs/api/국내주식/국내업종_현재지수.md` |

**파라미터 매핑**:

```python
params = {
    "FID_COND_MRKT_DIV_CODE": "U",  # 고정 (업종)
    "FID_INPUT_ISCD": industry_code, # 0001, 1001, 2001 등
}
```

### 5.3 `fetch_industry_category_price`

| 항목 | 값 |
|------|-----|
| 경로 | `/uapi/domestic-stock/v1/quotations/inquire-index-category-price` |
| TR ID | `FHPUP02140000` |
| 문서 | `docs/api/국내주식/국내업종_구분별전체시세.md` |

**파라미터 매핑**:

```python
params = {
    "FID_COND_MRKT_DIV_CODE": "U",       # 고정
    "FID_COND_SCR_DIV_CODE": "20214",    # 고정
    "FID_MRKT_CLS_CODE": market_type,    # K, Q, K2
    "FID_BLNG_CLS_CODE": "",             # 카테고리 코드 (API 문서 확인)
    "FID_INPUT_ISCD": "",                # 업종코드
}
```

### 5.4 `fetch_dividend_schedule` (후순위)

- API 문서 미작성 상태 → 한국투자 공식 문서 확인 후 구현
- 나머지 15개 API 완료 후 진행

---

## 6. 상수 추가 (`constants.py`)

### 6.1 차트 기간 코드

```python
# 국내주식 기간 코드
PERIOD_CODE = {
    "DAY": "D",
    "WEEK": "W",
    "MONTH": "M",
    "YEAR": "Y",
}

# 해외주식 기간 코드 (숫자)
OVERSEAS_PERIOD_CODE = {
    "DAY": "0",
    "WEEK": "1",
    "MONTH": "2",
}
```

### 6.2 순위 관련 코드

```python
# 거래량 순위 정렬
VOLUME_RANKING_SORT = {
    "AVG_VOLUME": "0",
    "VOLUME_INCREASE": "1",
    "TURNOVER_RATE": "2",
    "TRADING_AMOUNT": "3",
    "AMOUNT_TURNOVER": "4",
}

# 등락률 순위 정렬
CHANGE_RATE_SORT = {
    "RISE": "0",
    "FALL": "1",
    "FLAT_RISE": "2",
    "FLAT_FALL": "3",
    "CHANGE_RATE": "4",
}

# 시가총액 상위 시장 구분
MARKET_CAP_TARGET = {
    "ALL": "0000",
    "KRX": "0001",
    "KOSDAQ": "1001",
    "KOSPI200": "2001",
}
```

### 6.3 업종 지수 코드

```python
# 업종 지수 코드
INDUSTRY_INDEX_CODE = {
    "KOSPI": "0001",
    "KOSDAQ": "1001",
    "KOSPI200": "2001",
}
```

---

## 7. 테스트 파일 구성

각 Step별 테스트 파일 생성:

| Step | 테스트 파일 | 테스트 클래스 |
|------|-----------|-------------|
| Step 1 | `tests/test_chart_apis.py` | `TestFetchDomesticChart`, `TestFetchDomesticMinuteChart`, `TestFetchOverseasChart` |
| Step 2 | `tests/test_ranking_apis.py` | `TestFetchVolumeRanking`, `TestFetchChangeRateRanking`, `TestFetchMarketCapRanking`, `TestFetchOverseasChangeRateRanking` |
| Step 3 | `tests/test_financial_apis.py` | `TestFetchFinancialRatio`, `TestFetchIncomeStatement`, `TestFetchBalanceSheet`, `TestFetchProfitabilityRatio`, `TestFetchGrowthRatio` |
| Step 4 | `tests/test_dividend_industry_apis.py` | `TestFetchDividendRanking`, `TestFetchIndustryIndex`, `TestFetchIndustryCategoryPrice` |

각 테스트 클래스에 포함할 테스트:
1. `test_success_response` - 성공 응답
2. `test_request_url_and_headers` - URL, TR ID 검증
3. `test_request_params` - 파라미터 매핑 검증
4. `test_token_refresh_on_expiry` - 토큰 만료 재발급

---

## 8. `__init__.py` 내보내기

새 메서드는 `KoreaInvestment` 클래스에 추가되므로 `__init__.py` 변경 불필요.
상수가 외부에서 필요한 경우에만 `__init__.py`에 추가.

---

## 9. 버전 관리

각 Step PR 완료 시:
- `.bumpversion.cfg`의 `current_version` 업데이트
- `CHANGELOG.md`에 변경 내역 추가
- 최종 Step 4 완료 시 PyPI 릴리스
