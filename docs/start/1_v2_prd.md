# Phase 1: korea-investment-stock API 확장 PRD

## 1. 개요

### 1.1 목적

korea-investment-stock 라이브러리에 UI 개발에 필요한 핵심 GET API 16개를 추가하여, MoneyFlow 사이트 개발 시 라이브러리 수정 없이 UI 작업에 집중할 수 있도록 한다.

### 1.2 배경

현재 MoneyFlow 개발 프로세스의 병목:

```
현재: 4단계 파이프라인 (UI에서 새 데이터 필요할 때마다 반복)
korea-investment-stock (API 추가)
  → stock-data-batch (수집 로직 추가)
    → backend Go API (엔드포인트 추가)
      → frontend (UI 개발)
```

**문제점**: UI 개발 중 새로운 데이터가 필요하면 korea-investment-stock까지 돌아가서 API를 추가해야 하므로 개발 속도가 느려진다.

**해결**: Phase 1에서 핵심 GET API를 모두 추가 → 이후 UI 개발은 backend API + frontend만 수정하면 된다.

```
개선 후: 2단계로 축소
DB (이미 데이터 존재) → backend API (엔드포인트 추가) → frontend (UI 개발)
```

### 1.3 현재 API 커버리지

| 구분 | 전체 GET API | 현재 구현 | 커버리지 |
|------|-------------|----------|---------|
| 국내주식 | ~120개 | 5개 | 4% |
| 해외주식 | ~25개 | 2개 | 8% |

### 1.4 Phase 1 목표

- 추가 API: **16개** (국내 14개 + 해외 2개)
- 완료 후 커버리지: 핵심 조회 API 80% 이상

---

## 2. 구현 대상 API 목록

### 2.1 전체 목록 요약

| Step | 카테고리 | API 수 | 메서드명 | 우선순위 |
|------|---------|--------|---------|---------|
| Step 1 | 차트 데이터 | 3개 | `fetch_domestic_chart`, `fetch_domestic_minute_chart`, `fetch_overseas_chart` | 최우선 |
| Step 2 | 시세 순위 | 4개 | `fetch_volume_ranking`, `fetch_change_rate_ranking`, `fetch_market_cap_ranking`, `fetch_overseas_change_rate_ranking` | 높음 |
| Step 3 | 재무제표 | 5개 | `fetch_financial_ratio`, `fetch_income_statement`, `fetch_balance_sheet`, `fetch_profitability_ratio`, `fetch_growth_ratio` | 높음 |
| Step 4 | 배당 + 업종 | 4개 | `fetch_dividend_ranking`, `fetch_industry_index`, `fetch_industry_category_price`, `fetch_dividend_schedule` | 보통 |

---

## 3. Step 1: 차트 데이터 (3개 API)

### 3.1 국내주식기간별시세 - `fetch_domestic_chart`

일봉/주봉/월봉/년봉 차트 데이터 조회.

| 항목 | 값 |
|------|-----|
| API명 | 국내주식기간별시세(일/주/월/년) |
| API 경로 | `/uapi/domestic-stock/v1/quotations/inquire-daily-itemchartprice` |
| TR ID | FHKST03010100 (실전/모의 동일) |
| Method | GET |
| 문서 위치 | [docs/api/국내주식/국내주식기간별시세(일_주_월_년).md](../api/국내주식/국내주식기간별시세(일_주_월_년).md) |

**메서드 시그니처**:

```python
def fetch_domestic_chart(
    self,
    symbol: str,
    period: str = "D",
    start_date: str = "",
    end_date: str = "",
    adjusted: bool = True,
    market_code: str = "J"
) -> dict:
```

**파라미터**:

| 파라미터 | API 키 | 필수 | 설명 |
|---------|--------|------|------|
| `market_code` | `FID_COND_MRKT_DIV_CODE` | Y | J:KRX, NX:NXT, UN:통합 |
| `symbol` | `FID_INPUT_ISCD` | Y | 종목코드 (예: 005930) |
| `start_date` | `FID_INPUT_DATE_1` | Y | 조회시작일 (YYYYMMDD) |
| `end_date` | `FID_INPUT_DATE_2` | Y | 조회종료일 (YYYYMMDD), 최대 100건 |
| `period` | `FID_PERIOD_DIV_CODE` | Y | D:일, W:주, M:월, Y:년 |
| `adjusted` | `FID_ORG_ADJ_PRC` | Y | True→0(수정주가), False→1(원주가) |

**주요 응답 필드**:

| 필드 | 설명 |
|------|------|
| output1 | 요약 (현재가, 시가총액 등) |
| output2[] | 기간별 OHLCV 배열 |
| `stck_bsop_date` | 영업일자 |
| `stck_oprc` | 시가 |
| `stck_hgpr` | 고가 |
| `stck_lwpr` | 저가 |
| `stck_clpr` | 종가 |
| `acml_vol` | 거래량 |
| `acml_tr_pbmn` | 거래대금 |

**UI 활용**: TradingView 의존도 감소, 자체 차트 렌더링, 기술적 분석

---

### 3.2 주식당일분봉조회 - `fetch_domestic_minute_chart`

당일 분봉(1분/5분/10분/30분/60분) 데이터 조회.

| 항목 | 값 |
|------|-----|
| API명 | 주식당일분봉조회 |
| API 경로 | `/uapi/domestic-stock/v1/quotations/inquire-time-itemchartprice` |
| TR ID | FHKST03010200 (실전/모의 동일) |
| Method | GET |
| 문서 위치 | [docs/api/국내주식/주식당일분봉조회.md](../api/국내주식/주식당일분봉조회.md) |

**메서드 시그니처**:

```python
def fetch_domestic_minute_chart(
    self,
    symbol: str,
    time_from: str = "",
    market_code: str = "J"
) -> dict:
```

**파라미터**:

| 파라미터 | API 키 | 필수 | 설명 |
|---------|--------|------|------|
| `market_code` | `FID_COND_MRKT_DIV_CODE` | Y | J:KRX |
| `symbol` | `FID_INPUT_ISCD` | Y | 종목코드 |
| `time_from` | `FID_INPUT_HOUR_1` | Y | 조회시작시간 (HHMMSS) |
| - | `FID_PW_DATA_INCU_YN` | Y | 과거데이터 포함여부 |
| - | `FID_ETC_CLS_CODE` | Y | 기타분류코드 |

**주요 응답 필드**:

| 필드 | 설명 |
|------|------|
| `stck_cntg_hour` | 체결시간 (HHMMSS) |
| `stck_prpr` | 현재가 |
| `stck_oprc` | 시가 |
| `stck_hgpr` | 고가 |
| `stck_lwpr` | 저가 |
| `cntg_vol` | 체결거래량 |
| `acml_tr_pbmn` | 누적거래대금 |

**제약사항**: 1회 최대 30건, 당일 데이터만 조회 가능

**UI 활용**: 인트라데이 차트, 실시간 모니터링

---

### 3.3 해외주식 기간별시세 - `fetch_overseas_chart`

해외주식 일봉/주봉/월봉 차트 데이터 조회.

| 항목 | 값 |
|------|-----|
| API명 | 해외주식 기간별시세 |
| API 경로 | `/uapi/overseas-price/v1/quotations/dailyprice` |
| TR ID | HHDFS76240000 (실전/모의 동일) |
| Method | GET |
| 문서 위치 | [docs/api/해외주식/해외주식_기간별시세.md](../api/해외주식/해외주식_기간별시세.md) |

**메서드 시그니처**:

```python
def fetch_overseas_chart(
    self,
    symbol: str,
    country_code: str = "US",
    period: str = "D",
    end_date: str = "",
    adjusted: bool = True
) -> dict:
```

**파라미터**:

| 파라미터 | API 키 | 필수 | 설명 |
|---------|--------|------|------|
| - | `AUTH` | Y | 사용자권한 (빈값) |
| `country_code` → EXCD 변환 | `EXCD` | Y | NYS, NAS, AMS, HKS, TSE, SHS, SZS, HSX, HNX |
| `symbol` | `SYMB` | Y | 종목코드 (예: TSLA) |
| `period` | `GUBN` | Y | 0:일, 1:주, 2:월 |
| `end_date` | `BYMD` | Y | 기준일 (YYYYMMDD, 빈값=오늘) |
| `adjusted` | `MODP` | Y | 0:미수정, 1:수정 |

**주요 응답 필드** (output2[]):

| 필드 | 설명 |
|------|------|
| `xymd` | 일자 (YYYYMMDD) |
| `clos` | 종가 |
| `open` | 시가 |
| `high` | 고가 |
| `low` | 저가 |
| `tvol` | 거래량 |
| `tamt` | 거래대금 |

**제약사항**: 1회 최대 100건, 무료(지연) 데이터

**UI 활용**: US/해외 주식 차트, 비교 차트

---

## 4. Step 2: 시세 순위 (4개 API)

### 4.1 거래량순위 - `fetch_volume_ranking`

| 항목 | 값 |
|------|-----|
| API명 | 거래량순위 |
| API 경로 | `/uapi/domestic-stock/v1/quotations/volume-rank` |
| TR ID | FHPST01710000 (실전 전용) |
| Method | GET |
| 문서 위치 | [docs/api/국내주식/거래량순위.md](../api/국내주식/거래량순위.md) |

**메서드 시그니처**:

```python
def fetch_volume_ranking(
    self,
    market_code: str = "J",
    sort_by: str = "0"
) -> dict:
```

**파라미터**:

| 파라미터 | API 키 | 필수 | 설명 |
|---------|--------|------|------|
| `market_code` | `FID_COND_MRKT_DIV_CODE` | Y | J:KRX, NX:NXT |
| - | `FID_COND_SCR_DIV_CODE` | Y | "20171" (고정) |
| - | `FID_INPUT_ISCD` | Y | 종목/업종코드 |
| - | `FID_DIV_CLS_CODE` | Y | 0:전체, 1:보통주, 2:우선주 |
| `sort_by` | `FID_BLNG_CLS_CODE` | Y | 0:평균거래량, 1:거래증가, 2:회전율, 3:거래대금, 4:대금회전율 |

**주요 응답 필드**:

| 필드 | 설명 |
|------|------|
| `data_rank` | 순위 |
| `hts_kor_isnm` | 종목명 |
| `mksc_shrn_iscd` | 종목코드 |
| `stck_prpr` | 현재가 |
| `acml_vol` | 거래량 |
| `vol_inrt` | 거래량 증가율 |

**제약사항**: 최대 30건

**UI 활용**: 대시보드 거래량 순위 위젯

---

### 4.2 국내주식 등락률 순위 - `fetch_change_rate_ranking`

| 항목 | 값 |
|------|-----|
| API명 | 국내주식 등락률 순위 |
| API 경로 | `/uapi/domestic-stock/v1/ranking/fluctuation` |
| TR ID | FHPST01700000 (실전 전용) |
| Method | GET |
| 문서 위치 | [docs/api/국내주식/국내주식_등락률_순위.md](../api/국내주식/국내주식_등락률_순위.md) |

**메서드 시그니처**:

```python
def fetch_change_rate_ranking(
    self,
    market_code: str = "J",
    sort_order: str = "0",
    period_days: str = "0"
) -> dict:
```

**파라미터**:

| 파라미터 | API 키 | 필수 | 설명 |
|---------|--------|------|------|
| `market_code` | `fid_cond_mrkt_div_code` | Y | J:KRX, NX:NXT |
| `sort_order` | `fid_rank_sort_cls_code` | Y | 0:상승률, 1:하락률, 2:보합상승, 3:보합하락, 4:변동률 |
| `period_days` | `fid_input_cnt_1` | Y | 0:당일, 1:2일, 2:3일... |

**주요 응답 필드**:

| 필드 | 설명 |
|------|------|
| `data_rank` | 순위 |
| `hts_kor_isnm` | 종목명 |
| `stck_shrn_iscd` | 종목코드 |
| `stck_prpr` | 현재가 |
| `prdy_ctrt` | 등락률 (%) |
| `prd_rsfl_rate` | 기간 등락률 |

**제약사항**: 최대 30건

**UI 활용**: 상승/하락 종목 랭킹, 대시보드 위젯

---

### 4.3 국내주식 시가총액 상위 - `fetch_market_cap_ranking`

| 항목 | 값 |
|------|-----|
| API명 | 국내주식 시가총액 상위 |
| API 경로 | `/uapi/domestic-stock/v1/ranking/market-cap` |
| TR ID | FHPST01740000 (실전 전용) |
| Method | GET |
| 문서 위치 | [docs/api/국내주식/국내주식_시가총액_상위.md](../api/국내주식/국내주식_시가총액_상위.md) |

**메서드 시그니처**:

```python
def fetch_market_cap_ranking(
    self,
    market_code: str = "J",
    target_market: str = "0000"
) -> dict:
```

**파라미터**:

| 파라미터 | API 키 | 필수 | 설명 |
|---------|--------|------|------|
| `market_code` | `fid_cond_mrkt_div_code` | Y | J:KRX |
| `target_market` | `fid_input_iscd` | Y | 0000:전체, 0001:KRX, 1001:KOSDAQ, 2001:KOSPI200 |

**주요 응답 필드**:

| 필드 | 설명 |
|------|------|
| `data_rank` | 순위 |
| `hts_kor_isnm` | 종목명 |
| `mksc_shrn_iscd` | 종목코드 |
| `stck_prpr` | 현재가 |
| `stck_avls` | 시가총액 |
| `mrkt_whol_avls_rlim` | 시가총액 비중 (%) |

**제약사항**: 최대 30건

**UI 활용**: 시가총액 순위 테이블

---

### 4.4 해외주식 상승율/하락율 - `fetch_overseas_change_rate_ranking`

| 항목 | 값 |
|------|-----|
| API명 | 해외주식 상승율/하락율 |
| API 경로 | `/uapi/overseas-stock/v1/ranking/updown-rate` |
| TR ID | HHDFS76290000 (실전 전용) |
| Method | GET |
| 문서 위치 | [docs/api/해외주식/해외주식_상승율_하락율.md](../api/해외주식/해외주식_상승율_하락율.md) |

**메서드 시그니처**:

```python
def fetch_overseas_change_rate_ranking(
    self,
    country_code: str = "US",
    sort_order: str = "1",
    period: str = "0",
    volume_filter: str = "0"
) -> dict:
```

**파라미터**:

| 파라미터 | API 키 | 필수 | 설명 |
|---------|--------|------|------|
| `country_code` → EXCD 변환 | `EXCD` | Y | NYS, NAS, AMS, HKS 등 |
| `sort_order` | `GUBN` | Y | 0:하락률, 1:상승률 |
| `period` | `NDAY` | Y | 0:당일, 1:2일... 9:1년 |
| `volume_filter` | `VOL_RANG` | Y | 0:전체, 1:100주이상... |

**주요 응답 필드** (output2[]):

| 필드 | 설명 |
|------|------|
| `rank` | 순위 |
| `name` | 종목명 |
| `symb` | 종목코드 |
| `last` | 현재가 |
| `rate` | 등락률 (%) |
| `tvol` | 거래량 |

**UI 활용**: 해외주식 상승/하락 랭킹

---

## 5. Step 3: 재무제표 (5개 API)

### 5.1 국내주식 재무비율 - `fetch_financial_ratio`

| 항목 | 값 |
|------|-----|
| API명 | 국내주식 재무비율 |
| API 경로 | `/uapi/domestic-stock/v1/finance/financial-ratio` |
| TR ID | FHKST66430300 (실전 전용) |
| Method | GET |
| 문서 위치 | [docs/api/국내주식/국내주식_재무비율.md](../api/국내주식/국내주식_재무비율.md) |

**메서드 시그니처**:

```python
def fetch_financial_ratio(
    self,
    symbol: str,
    period_type: str = "0"
) -> dict:
```

**파라미터**:

| 파라미터 | API 키 | 필수 | 설명 |
|---------|--------|------|------|
| `period_type` | `FID_DIV_CLS_CODE` | Y | 0:연간, 1:분기 |
| - | `fid_cond_mrkt_div_code` | Y | "J" (고정) |
| `symbol` | `fid_input_iscd` | Y | 종목코드 |

**주요 응답 필드**:

| 필드 | 설명 |
|------|------|
| `stac_yymm` | 결산년월 |
| `grs` | 매출액 증가율 |
| `bsop_prfi_inrt` | 영업이익 증가율 |
| `ntin_inrt` | 순이익 증가율 |
| `roe_val` | ROE |
| `eps` | EPS |
| `bps` | BPS |
| `rsrv_rate` | 유보율 |
| `lblt_rate` | 부채비율 |

**UI 활용**: 종목 상세 페이지 재무지표 섹션

---

### 5.2 국내주식 손익계산서 - `fetch_income_statement`

| 항목 | 값 |
|------|-----|
| API명 | 국내주식 손익계산서 |
| API 경로 | `/uapi/domestic-stock/v1/finance/income-statement` |
| TR ID | FHKST66430200 (실전 전용) |
| Method | GET |
| 문서 위치 | [docs/api/국내주식/국내주식_손익계산서.md](../api/국내주식/국내주식_손익계산서.md) |

**메서드 시그니처**:

```python
def fetch_income_statement(
    self,
    symbol: str,
    period_type: str = "0"
) -> dict:
```

**파라미터**:

| 파라미터 | API 키 | 필수 | 설명 |
|---------|--------|------|------|
| `period_type` | `FID_DIV_CLS_CODE` | Y | 0:연간, 1:분기(누적) |
| - | `fid_cond_mrkt_div_code` | Y | "J" (고정) |
| `symbol` | `fid_input_iscd` | Y | 종목코드 |

**주요 응답 필드**:

| 필드 | 설명 |
|------|------|
| `stac_yymm` | 결산년월 |
| `sale_account` | 매출액 |
| `sale_cost` | 매출원가 |
| `sale_totl_prfi` | 매출총이익 |
| `bsop_prti` | 영업이익 |
| `op_prfi` | 경상이익 |
| `thtr_ntin` | 당기순이익 |

**UI 활용**: 실적 추이 차트 (매출/영업이익/순이익)

---

### 5.3 국내주식 대차대조표 - `fetch_balance_sheet`

| 항목 | 값 |
|------|-----|
| API명 | 국내주식 대차대조표 |
| API 경로 | `/uapi/domestic-stock/v1/finance/balance-sheet` |
| TR ID | FHKST66430100 (실전 전용) |
| Method | GET |
| 문서 위치 | [docs/api/국내주식/국내주식_대차대조표.md](../api/국내주식/국내주식_대차대조표.md) |

**메서드 시그니처**:

```python
def fetch_balance_sheet(
    self,
    symbol: str,
    period_type: str = "0"
) -> dict:
```

**파라미터**:

| 파라미터 | API 키 | 필수 | 설명 |
|---------|--------|------|------|
| `period_type` | `FID_DIV_CLS_CODE` | Y | 0:연간, 1:분기 |
| - | `fid_cond_mrkt_div_code` | Y | "J" (고정) |
| `symbol` | `fid_input_iscd` | Y | 종목코드 |

**주요 응답 필드**:

| 필드 | 설명 |
|------|------|
| `stac_yymm` | 결산년월 |
| `cras` | 유동자산 |
| `fxas` | 고정자산 |
| `total_aset` | 자산총계 |
| `flow_lblt` | 유동부채 |
| `fix_lblt` | 고정부채 |
| `total_lblt` | 부채총계 |
| `cpfn` | 자본금 |
| `total_cptl` | 자본총계 |

**UI 활용**: 재무 구조 분석 (자산/부채/자본 비율)

---

### 5.4 국내주식 수익성비율 - `fetch_profitability_ratio`

| 항목 | 값 |
|------|-----|
| API명 | 국내주식 수익성비율 |
| API 경로 | `/uapi/domestic-stock/v1/finance/profit-ratio` |
| TR ID | FHKST66430400 (실전 전용) |
| Method | GET |
| 문서 위치 | [docs/api/국내주식/국내주식_수익성비율.md](../api/국내주식/국내주식_수익성비율.md) |

**메서드 시그니처**:

```python
def fetch_profitability_ratio(
    self,
    symbol: str,
    period_type: str = "0"
) -> dict:
```

**파라미터**:

| 파라미터 | API 키 | 필수 | 설명 |
|---------|--------|------|------|
| `symbol` | `fid_input_iscd` | Y | 종목코드 |
| `period_type` | `FID_DIV_CLS_CODE` | Y | 0:연간, 1:분기 |
| - | `fid_cond_mrkt_div_code` | Y | "J" (고정) |

**주요 응답 필드**:

| 필드 | 설명 |
|------|------|
| `stac_yymm` | 결산년월 |
| `cptl_ntin_rate` | 총자본 순이익률 |
| `self_cptl_ntin_inrt` | 자기자본 순이익률 (ROE) |
| `sale_ntin_rate` | 매출액 순이익률 |
| `sale_totl_rate` | 매출액 총이익률 |

**UI 활용**: 수익성 지표 (ROE, 매출이익률 등)

---

### 5.5 국내주식 성장성비율 - `fetch_growth_ratio`

| 항목 | 값 |
|------|-----|
| API명 | 국내주식 성장성비율 |
| API 경로 | `/uapi/domestic-stock/v1/finance/growth-ratio` |
| TR ID | FHKST66430800 (실전 전용) |
| Method | GET |
| 문서 위치 | [docs/api/국내주식/국내주식_성장성비율.md](../api/국내주식/국내주식_성장성비율.md) |

**메서드 시그니처**:

```python
def fetch_growth_ratio(
    self,
    symbol: str,
    period_type: str = "0"
) -> dict:
```

**파라미터**:

| 파라미터 | API 키 | 필수 | 설명 |
|---------|--------|------|------|
| `symbol` | `fid_input_iscd` | Y | 종목코드 |
| `period_type` | `fid_div_cls_code` | Y | 0:연간, 1:분기 |
| - | `fid_cond_mrkt_div_code` | Y | "J" (고정) |

**주요 응답 필드**:

| 필드 | 설명 |
|------|------|
| `stac_yymm` | 결산년월 |
| `grs` | 매출액 증가율 |
| `bsop_prfi_inrt` | 영업이익 증가율 |
| `equt_inrt` | 자기자본 증가율 |
| `totl_aset_inrt` | 총자산 증가율 |

**UI 활용**: 성장성 분석 (매출/이익 성장률 추이)

---

## 6. Step 4: 배당 + 업종 (4개 API)

### 6.1 국내주식 배당률 상위 - `fetch_dividend_ranking`

| 항목 | 값 |
|------|-----|
| API명 | 국내주식 배당률 상위 |
| API 경로 | `/uapi/domestic-stock/v1/ranking/dividend-rate` |
| TR ID | HHKDB13470100 (실전 전용) |
| Method | GET |
| 문서 위치 | [docs/api/국내주식/국내주식_배당률_상위.md](../api/국내주식/국내주식_배당률_상위.md) |

**메서드 시그니처**:

```python
def fetch_dividend_ranking(
    self,
    market_type: str = "0",
    start_date: str = "",
    end_date: str = "",
    dividend_type: str = "2"
) -> dict:
```

**파라미터**:

| 파라미터 | API 키 | 필수 | 설명 |
|---------|--------|------|------|
| `market_type` | `GB1` | Y | 0:전체, 1:KOSPI, 2:KOSPI200, 3:KOSDAQ |
| - | `UPJONG` | Y | 업종코드 |
| - | `GB2` | Y | 0:전체, 6:보통주, 7:우선주 |
| `dividend_type` | `GB3` | Y | 1:주식, 2:현금 |
| `start_date` | `F_DT` | Y | 기준일 from (YYYYMMDD) |
| `end_date` | `T_DT` | Y | 기준일 to (YYYYMMDD) |

**주요 응답 필드**:

| 필드 | 설명 |
|------|------|
| `rank` | 순위 |
| `sht_cd` | 종목코드 |
| `isin_name` | 종목명 |
| `record_date` | 기준일 |
| `per_sto_divi_amt` | 주당배당금 |
| `divi_rate` | 배당률 (%) |

**제약사항**: 최대 30건

**UI 활용**: 고배당 종목 랭킹

---

### 6.2 국내업종 현재지수 - `fetch_industry_index`

| 항목 | 값 |
|------|-----|
| API명 | 국내업종 현재지수 |
| API 경로 | `/uapi/domestic-stock/v1/quotations/inquire-index-price` |
| TR ID | FHPUP02100000 (실전 전용) |
| Method | GET |
| 문서 위치 | [docs/api/국내주식/국내업종_현재지수.md](../api/국내주식/국내업종_현재지수.md) |

**메서드 시그니처**:

```python
def fetch_industry_index(
    self,
    industry_code: str = "0001"
) -> dict:
```

**파라미터**:

| 파라미터 | API 키 | 필수 | 설명 |
|---------|--------|------|------|
| - | `FID_COND_MRKT_DIV_CODE` | Y | "U" (고정, 업종) |
| `industry_code` | `FID_INPUT_ISCD` | Y | 0001:KOSPI, 1001:KOSDAQ, 2001:KOSPI200 |

**주요 응답 필드**:

| 필드 | 설명 |
|------|------|
| `bstp_nmix_prpr` | 업종지수 현재가 |
| `bstp_nmix_prdy_vrss` | 전일대비 |
| `bstp_nmix_prdy_ctrt` | 등락률 (%) |
| `bstp_nmix_oprc` | 시가 |
| `bstp_nmix_hgpr` | 고가 |
| `bstp_nmix_lwpr` | 저가 |
| `ascn_issu_cnt` | 상승 종목수 |
| `down_issu_cnt` | 하락 종목수 |

**UI 활용**: 시장 지수 현황, 섹터별 등락

---

### 6.3 국내업종 구분별전체시세 - `fetch_industry_category_price`

| 항목 | 값 |
|------|-----|
| API명 | 국내업종 구분별전체시세 |
| API 경로 | `/uapi/domestic-stock/v1/quotations/inquire-index-category-price` |
| TR ID | FHPUP02140000 (실전 전용) |
| Method | GET |
| 문서 위치 | [docs/api/국내주식/국내업종_구분별전체시세.md](../api/국내주식/국내업종_구분별전체시세.md) |

**메서드 시그니처**:

```python
def fetch_industry_category_price(
    self,
    market_type: str = "K"
) -> dict:
```

**파라미터**:

| 파라미터 | API 키 | 필수 | 설명 |
|---------|--------|------|------|
| - | `FID_COND_MRKT_DIV_CODE` | Y | "U" (고정) |
| - | `FID_COND_SCR_DIV_CODE` | Y | "20214" (고정) |
| `market_type` | `FID_MRKT_CLS_CODE` | Y | K:KRX, Q:KOSDAQ, K2:KOSPI200 |
| - | `FID_BLNG_CLS_CODE` | Y | 카테고리 코드 (시장별 상이) |
| - | `FID_INPUT_ISCD` | Y | 업종코드 |

**주요 응답 필드**:

| 필드 | 설명 |
|------|------|
| output1 | 전체 지수 요약 |
| output2[] | 업종별 상세 배열 |
| `bstp_cls_code` | 업종 분류 코드 |
| `hts_kor_isnm` | 업종명 |
| `bstp_nmix_prpr` | 업종 지수 |
| `acml_vol` | 거래량 |
| `acml_tr_pbmn` | 거래대금 |

**UI 활용**: 섹터 히트맵, 업종별 등락 현황

---

### 6.4 예탁원정보(배당일정) - `fetch_dividend_schedule`

| 항목 | 값 |
|------|-----|
| API명 | 예탁원정보(배당일정) |
| API 경로 | `/uapi/domestic-stock/v1/ksdinfo/dividend-info` (확인 필요) |
| TR ID | 확인 필요 |
| Method | GET |
| 문서 위치 | 문서 미작성 (API 문서 확인 후 업데이트 필요) |

**⚠️ 참고**: 이 API는 docs에 문서가 아직 없습니다. 한국투자증권 OpenAPI 공식 문서에서 확인 후 구현이 필요합니다. 나머지 15개 API를 먼저 구현하고, 이 API는 후순위로 진행합니다.

**메서드 시그니처** (예상):

```python
def fetch_dividend_schedule(
    self,
    symbol: str = "",
    start_date: str = "",
    end_date: str = ""
) -> dict:
```

**UI 활용**: 배당 캘린더, 배당 이벤트 타임라인

---

## 7. 수정 대상 파일

### 7.1 korea-investment-stock (라이브러리)

| 파일 | 변경 내용 |
|------|----------|
| `korea_investment_stock/korea_investment_stock.py` | 15~16개 `fetch_*` 메서드 추가 |
| `korea_investment_stock/constants.py` | 업종 코드, 순위 정렬 코드 등 상수 추가 |
| `korea_investment_stock/cache/cached_korea_investment.py` | 새 메서드에 대한 캐시 래퍼 추가 |
| `korea_investment_stock/rate_limit/rate_limited_korea_investment.py` | 새 메서드에 대한 rate limit 래퍼 추가 |
| `korea_investment_stock/tests/` | 각 API별 단위 테스트 추가 |
| `CHANGELOG.md` | 버전 업데이트 |
| `pyproject.toml` | 버전 번호 업데이트 |
| `CLAUDE.md` | 새 메서드 문서 추가 |

### 7.2 stock-data-batch (데이터 수집) - 라이브러리 구현 후

| 파일 | 변경 내용 |
|------|----------|
| `database/models.py` | 새 DB 모델 추가 (재무제표, 차트 히스토리 등) |
| `services/stock_collector.py` | 새 API 호출 로직 추가 |
| `services/data_saver.py` | 새 UPSERT 함수 추가 |
| `main.py` | CLI 옵션 추가 |
| `config/config.yaml` | 수집 대상 설정 추가 |

### 7.3 moneyflow.advenoh.pe.kr (DB 스키마) - 라이브러리 구현 후

| 파일 | 변경 내용 |
|------|----------|
| `backend/db/changelog/YYYY-MM/` | Liquibase changelog 추가 |

---

## 8. 구현 순서 및 작업 단위

### Phase 1-1: korea-investment-stock 라이브러리 확장

각 Step은 독립적인 PR로 진행:

```
Step 1: 차트 데이터 API 3개 → PR #1
  - fetch_domestic_chart
  - fetch_domestic_minute_chart
  - fetch_overseas_chart
  - 단위 테스트
  - 버전 업데이트

Step 2: 시세 순위 API 4개 → PR #2
  - fetch_volume_ranking
  - fetch_change_rate_ranking
  - fetch_market_cap_ranking
  - fetch_overseas_change_rate_ranking
  - 단위 테스트
  - 버전 업데이트

Step 3: 재무제표 API 5개 → PR #3
  - fetch_financial_ratio
  - fetch_income_statement
  - fetch_balance_sheet
  - fetch_profitability_ratio
  - fetch_growth_ratio
  - 단위 테스트
  - 버전 업데이트

Step 4: 배당 + 업종 API 3~4개 → PR #4
  - fetch_dividend_ranking
  - fetch_industry_index
  - fetch_industry_category_price
  - (fetch_dividend_schedule - 문서 확인 후)
  - 단위 테스트
  - 버전 업데이트
```

### Phase 1-2: DB 스키마 + stock-data-batch (라이브러리 완료 후)

```
Step 5: DB 스키마 추가 → PR #5 (moneyflow.advenoh.pe.kr)
  - 재무제표 테이블
  - 차트 히스토리 테이블 (필요 시)
  - 순위 데이터 테이블 (필요 시)
  - Liquibase changelog

Step 6: stock-data-batch 확장 → PR #6
  - 새 DB 모델 추가
  - 수집 로직 추가
  - CLI 옵션 추가
  - 통합 테스트
```

---

## 9. 설계 원칙

### 9.1 메서드 네이밍 컨벤션

기존 패턴을 따름:

```python
# 국내주식: fetch_domestic_*
fetch_domestic_chart()
fetch_domestic_minute_chart()

# 해외주식: fetch_overseas_*
fetch_overseas_chart()
fetch_overseas_change_rate_ranking()

# 재무제표: fetch_* (종목 단위)
fetch_financial_ratio()
fetch_income_statement()
fetch_balance_sheet()

# 순위/랭킹: fetch_*_ranking
fetch_volume_ranking()
fetch_change_rate_ranking()
fetch_market_cap_ranking()
fetch_dividend_ranking()

# 업종: fetch_industry_*
fetch_industry_index()
fetch_industry_category_price()
```

### 9.2 파라미터 설계

- `country_code` 사용 (기존 패턴과 일관성)
- 내부에서 `EXCD` 등 API 파라미터로 변환
- 기본값은 가장 일반적인 값으로 설정
- 필수값이 아닌 필터링 파라미터는 Optional로

### 9.3 응답 형식

- API 원본 응답을 그대로 반환 (기존 패턴 유지)
- 사용자가 필요한 필드만 선택하여 사용
- 라이브러리 철학: "Simple, transparent, and flexible"

### 9.4 토큰 자동 갱신

- 모든 새 메서드에 토큰 만료 시 자동 재발급 적용
- 기존 `_request_with_token_refresh` 패턴 활용

### 9.5 DB 필드 네이밍 규칙: API 원본 필드명 사용

**결정사항**: stock-data-batch의 DB 컬럼명을 한국투자 API 응답 필드명과 동일하게 사용한다.

**이유**:
- `docs/api/` 디렉토리에 모든 API 문서가 md로 작성되어 있어, Claude Code가 필드 의미를 즉시 참조 가능
- API 응답 → DB 저장 시 필드 매핑 변환 코드 불필요 (코드 단순화)
- API 문서와 DB 스키마가 1:1 대응되어 디버깅 용이

**적용 범위**:
- 한국투자 API 데이터 테이블: API 원본 필드명 사용
- 비한국투자 데이터 테이블 (Yahoo Finance, CNN): 각 소스의 필드명 또는 기존 이름 유지

**예시**:

```python
# 변경 전 (가독성 영어)
class KRStock(Model):
    current_price = IntegerField()
    trading_volume = BigIntegerField()
    market_cap = BigIntegerField()
    price_change_rate = DecimalField()

# 변경 후 (API 원본 필드명)
class KRStock(Model):
    stck_prpr = IntegerField()      # 주식현재가
    acml_vol = BigIntegerField()    # 누적거래량
    hts_avls = BigIntegerField()    # HTS시가총액(억)
    prdy_ctrt = DecimalField()      # 전일대비율
```

**필드 참조 문서**: 각 API의 응답 필드 설명은 `docs/api/국내주식/`, `docs/api/해외주식/` 참조

---

## 10. 기존 테이블 마이그레이션

### 10.1 마이그레이션 방침

- **기존 테이블 전체 DROP 후 재생성** (데이터 재수집 가능하므로 마이그레이션 불필요)
- 새 테이블은 API 원본 필드명으로 작성
- `symbols` 테이블은 메타 정보이므로 현재 필드명 유지 (API 응답 필드가 아님)

### 10.2 DROP 대상 테이블

| 테이블 | 데이터 소스 | 재수집 | 비고 |
|--------|-----------|--------|------|
| `kr_stocks` | 한국투자 API | batch로 재수집 | 필드명 변경 |
| `kr_etf` | 한국투자 API | batch로 재수집 | 필드명 변경 |
| `us_stocks` | 한국투자 API | batch로 재수집 | 필드명 변경 |
| `us_etf` | 한국투자 API | batch로 재수집 | 필드명 변경 |
| `kr_investor_trading` | 한국투자 API | batch로 재수집 | 필드명 변경 |
| `market_investor_trend` | 한국투자 API | batch로 재수집 | 필드명 변경 |

### 10.3 유지 대상 테이블 (비한국투자 데이터)

| 테이블 | 데이터 소스 | 비고 |
|--------|-----------|------|
| `symbols` | 메타 정보 | 필드명 유지 |
| `market_index_history` | Yahoo Finance | 필드명 유지 |
| `forex_rate_history` | Yahoo Finance | 필드명 유지 |
| `commodity_history` | Yahoo Finance | 필드명 유지 |
| `futures_history` | Yahoo Finance | 필드명 유지 |
| `fear_greed_index` | CNN | 필드명 유지 |

### 10.4 필드명 변경 매핑 (기존 → API 원본)

#### kr_stocks (국내주식 현재가 시세 API: FHKST01010100)

| 기존 필드 | API 원본 필드 | 설명 |
|----------|-------------|------|
| `current_price` | `stck_prpr` | 주식현재가 |
| `high_price` | `stck_hgpr` | 주식최고가 |
| `low_price` | `stck_lwpr` | 주식최저가 |
| `w52_high_price` | `w52_hgpr` | 52주최고가 |
| `w52_low_price` | `w52_lwpr` | 52주최저가 |
| `d250_high_price` | `d250_hgpr` | 250일최고가 |
| `d250_low_price` | `d250_lwpr` | 250일최저가 |
| `per` | `per` | PER (동일) |
| `pbr` | `pbr` | PBR (동일) |
| `eps` | `eps` | EPS (동일) |
| `bps` | `bps` | BPS (동일) |
| `dividend_yield` | `divi_rate` | 배당수익률 |
| `sector` | `bstp_kor_isnm` | 업종한글종목명 |
| `industry` | `idx_bztp_scls_cd_name` | 지수업종소분류코드명 |
| `foreign_holding_exhaustion_rate` | `frgn_hldn_qty_rt` | 외국인보유수량비율 |
| `price_change_rate` | `prdy_ctrt` | 전일대비율 |
| `trading_volume` | `acml_vol` | 누적거래량 |
| `market_cap` | `hts_avls` | HTS시가총액(억) |
| `listed_date` | `ssts_cntg_dtm` | 상장일 (구현 시 확인) |
| `listed_stock_count` | `lstn_stcn` | 상장주수 |

#### kr_etf (ETF/ETN 현재가 API: FHPST02400000)

| 기존 필드 | API 원본 필드 | 설명 |
|----------|-------------|------|
| `current_price` | `stck_prpr` | 현재가 |
| `high_price` | `stck_hgpr` | 최고가 |
| `low_price` | `stck_lwpr` | 최저가 |
| `w52_high_price` | `w52_hgpr` | 52주최고가 |
| `w52_low_price` | `w52_lwpr` | 52주최저가 |
| `d250_high_price` | `d250_hgpr` | 250일최고가 |
| `d250_low_price` | `d250_lwpr` | 250일최저가 |
| `m1_yield` | `m1_divi_rate` | 1개월수익률 (구현 시 확인) |
| `m3_yield` | `m3_divi_rate` | 3개월수익률 (구현 시 확인) |
| `m6_yield` | `m6_divi_rate` | 6개월수익률 (구현 시 확인) |
| `y1_yield` | `y1_divi_rate` | 1년수익률 (구현 시 확인) |
| `nav` | `nav` | NAV (동일) |
| `operating_cost` | `opng_cost` | 운영비용 (구현 시 확인) |
| `price_change_rate` | `prdy_ctrt` | 전일대비율 |
| `trading_volume` | `acml_vol` | 누적거래량 |
| `market_cap` | `hts_avls` | HTS시가총액 |

#### us_stocks (해외주식 현재가상세 API: HHDFS76200200)

| 기존 필드 | API 원본 필드 | 설명 |
|----------|-------------|------|
| `current_price` | `last` | 현재가 |
| `w52_high_price` | `h52p` | 52주최고가 |
| `w52_low_price` | `l52p` | 52주최저가 |
| `per` | `perx` | PER |
| `pbr` | `pbrx` | PBR |
| `eps` | `epsx` | EPS |
| `bps` | `bpsx` | BPS |
| `dividend_yield` | `divi_rate` | 배당수익률 (구현 시 확인) |
| `price_change_rate` | `t_xrat` | 등락률 |
| `trading_volume` | `tvol` | 거래량 |
| `market_cap` | `tomv` | 시가총액 |
| `shares_outstanding` | `shar` | 상장주수 |
| `sector` | `e_icod` | 업종코드 (구현 시 확인) |
| `industry` | `e_icod` | 업종코드 (구현 시 확인) |

#### us_etf (해외주식 현재가상세 API: HHDFS76200200)

us_stocks와 동일한 매핑 적용 (필드 수 감소)

#### kr_investor_trading (종목별 투자자매매동향 API: FHPTJ04160001)

| 기존 필드 | API 원본 필드 | 설명 |
|----------|-------------|------|
| `close_price` | `stck_prpr` | 주식현재가 |
| `price_change` | `prdy_vrss` | 전일대비 |
| `price_change_rate` | `prdy_ctrt` | 전일대비율 |
| `trading_volume` | `acml_vol` | 누적거래량 |
| `trading_amount` | `acml_tr_pbmn` | 누적거래대금 |
| `foreign_net_qty` | `frgn_ntby_qty` | 외국인순매수수량 |
| `foreign_net_amount` | `frgn_ntby_tr_pbmn` | 외국인순매수거래대금 |
| `institution_net_qty` | `orgn_ntby_qty` | 기관순매수수량 |
| `institution_net_amount` | `orgn_ntby_tr_pbmn` | 기관순매수거래대금 |
| `individual_net_qty` | `prsn_ntby_qty` | 개인순매수수량 |
| `individual_net_amount` | `prsn_ntby_tr_pbmn` | 개인순매수거래대금 |

#### market_investor_trend (시장별 투자자매매동향 API: FHPTJ04030000)

| 기존 필드 | API 원본 필드 | 설명 |
|----------|-------------|------|
| `foreign_net_qty` | `frgn_ntby_qty` | 외국인순매수수량 |
| `foreign_net_amount` | `frgn_ntby_tr_pbmn` | 외국인순매수거래대금 |
| `individual_net_qty` | `prsn_ntby_qty` | 개인순매수수량 |
| `individual_net_amount` | `prsn_ntby_tr_pbmn` | 개인순매수거래대금 |
| `institution_net_qty` | `orgn_ntby_qty` | 기관순매수수량 |
| `institution_net_amount` | `orgn_ntby_tr_pbmn` | 기관순매수거래대금 |
| `securities_net_qty` | `scrt_ntby_qty` | 증권순매수수량 |
| `securities_net_amount` | `scrt_ntby_tr_pbmn` | 증권순매수거래대금 |
| `investment_trust_net_qty` | `ivtr_ntby_qty` | 투자신탁순매수수량 |
| `investment_trust_net_amount` | `ivtr_ntby_tr_pbmn` | 투자신탁순매수거래대금 |
| `pe_fund_net_qty` | `pe_fund_ntby_qty` | 사모펀드순매수수량 |
| `pe_fund_net_amount` | `pe_fund_ntby_tr_pbmn` | 사모펀드순매수거래대금 |
| `bank_net_qty` | `bank_ntby_qty` | 은행순매수수량 |
| `bank_net_amount` | `bank_ntby_tr_pbmn` | 은행순매수거래대금 |
| `insurance_net_qty` | `insu_ntby_qty` | 보험순매수수량 |
| `insurance_net_amount` | `insu_ntby_tr_pbmn` | 보험순매수거래대금 |
| `pension_fund_net_qty` | `fund_ntby_qty` | 기금순매수수량 |
| `pension_fund_net_amount` | `fund_ntby_tr_pbmn` | 기금순매수거래대금 |

### 10.5 구현 시 주의사항

- 일부 필드명은 `(구현 시 확인)` 표시 → 실제 API 응답과 대조 후 확정
- `symbols` 테이블의 FK 관계는 유지
- `trade_date`, `created_at`, `updated_at` 등 메타 필드는 DB 전용이므로 가독성 이름 유지
- 새로 추가되는 재무제표 등 테이블도 동일한 원칙 적용 (API 원본 필드명 사용)

### 10.6 마이그레이션 실행 순서

```
1. Liquibase changelog 작성 (기존 테이블 DROP + 새 테이블 CREATE)
2. DB 마이그레이션 실행
3. stock-data-batch 모델 수정 (API 원본 필드명)
4. backend Go API 수정 (새 필드명에 맞게)
5. frontend API 타입 수정
6. batch 실행하여 데이터 재수집
```

---

## 11. 참고 자료

### API 문서 위치

- 국내주식: `docs/api/국내주식/`
- 해외주식: `docs/api/해외주식/`
- 한국투자증권 공식: https://wikidocs.net/book/7845

### 필드명 참조

- API 응답 필드 설명 → `docs/api/` 하위 각 API 문서 참조
- DB 필드와 API 필드가 동일하므로, API 문서가 곧 DB 스키마 문서

### 관련 문서

- 기존 PRD 예시: [docs/done/6_stock_prd.md](../done/6_stock_prd.md)
- 시장별 투자자매매동향 PRD: [docs/start/4_investor_time_prd.md](4_investor_time_prd.md)
