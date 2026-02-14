# Phase 1: korea-investment-stock API 확장 TODO

## Step 1: 차트 데이터 API (3개) → PR #1

### 메인 클래스 (`korea_investment_stock.py`)

- [x] `fetch_domestic_chart()` 구현 (TR ID: FHKST03010100)
- [x] `fetch_domestic_minute_chart()` 구현 (TR ID: FHKST03010200)
- [x] `fetch_overseas_chart()` 구현 (TR ID: HHDFS76240000)

### 상수 (`constants.py`)

- [x] `PERIOD_CODE` 추가 (D, W, M, Y) → 메서드 내부에서 직접 처리 (별도 상수 불필요)
- [x] `OVERSEAS_PERIOD_CODE` 추가 (0, 1, 2) → 메서드 내부 period_map으로 처리

### 캐시 래퍼 (`cached_korea_investment.py`)

- [x] `fetch_domestic_chart()` 캐시 래퍼 추가 (TTL: price)
- [x] `fetch_domestic_minute_chart()` 캐시 래퍼 추가 (TTL: price)
- [x] `fetch_overseas_chart()` 캐시 래퍼 추가 (TTL: price)

### Rate Limit 래퍼 (`rate_limited_korea_investment.py`)

- [x] `fetch_domestic_chart()` rate limit 래퍼 추가
- [x] `fetch_domestic_minute_chart()` rate limit 래퍼 추가
- [x] `fetch_overseas_chart()` rate limit 래퍼 추가

### 테스트 (`tests/test_chart_apis.py`)

- [x] `TestFetchDomesticChart` - 성공 응답, URL/헤더, 파라미터, 토큰 재발급
- [x] `TestFetchDomesticMinuteChart` - 성공 응답, URL/헤더, 파라미터, 토큰 재발급
- [x] `TestFetchOverseasChart` - 성공 응답, URL/헤더, 파라미터, 토큰 재발급

### 문서 및 버전

- [ ] `CHANGELOG.md` 업데이트
- [ ] `CLAUDE.md` 메서드 목록 업데이트
- [ ] PR 생성 및 머지

---

## Step 2: 시세 순위 API (4개) → PR #2

### 메인 클래스 (`korea_investment_stock.py`)

- [x] `fetch_volume_ranking()` 구현 (TR ID: FHPST01710000)
- [x] `fetch_change_rate_ranking()` 구현 (TR ID: FHPST01700000)
- [x] `fetch_market_cap_ranking()` 구현 (TR ID: FHPST01740000)
- [x] `fetch_overseas_change_rate_ranking()` 구현 (TR ID: HHDFS76290000)

### 상수 (`constants.py`)

- [x] `VOLUME_RANKING_SORT` 추가 → 메서드 내부에서 직접 처리 (별도 상수 불필요)
- [x] `CHANGE_RATE_SORT` 추가 → 메서드 내부에서 직접 처리 (별도 상수 불필요)
- [x] `MARKET_CAP_TARGET` 추가 → 메서드 내부에서 직접 처리 (별도 상수 불필요)

### 캐시 래퍼 (`cached_korea_investment.py`)

- [x] `fetch_volume_ranking()` 캐시 래퍼 추가 (TTL: price)
- [x] `fetch_change_rate_ranking()` 캐시 래퍼 추가 (TTL: price)
- [x] `fetch_market_cap_ranking()` 캐시 래퍼 추가 (TTL: price)
- [x] `fetch_overseas_change_rate_ranking()` 캐시 래퍼 추가 (TTL: price)

### Rate Limit 래퍼 (`rate_limited_korea_investment.py`)

- [x] `fetch_volume_ranking()` rate limit 래퍼 추가
- [x] `fetch_change_rate_ranking()` rate limit 래퍼 추가
- [x] `fetch_market_cap_ranking()` rate limit 래퍼 추가
- [x] `fetch_overseas_change_rate_ranking()` rate limit 래퍼 추가

### 테스트 (`tests/test_ranking_apis.py`)

- [x] `TestFetchVolumeRanking` - 성공 응답, URL/헤더, 파라미터, 정렬 옵션
- [x] `TestFetchChangeRateRanking` - 성공 응답, URL/헤더, 파라미터, 정렬 옵션
- [x] `TestFetchMarketCapRanking` - 성공 응답, URL/헤더, 파라미터, 시장 구분
- [x] `TestFetchOverseasChangeRateRanking` - 성공 응답, URL/헤더, 파라미터

### 문서 및 버전

- [ ] `CHANGELOG.md` 업데이트
- [ ] `CLAUDE.md` 메서드 목록 업데이트
- [ ] PR 생성 및 머지

---

## Step 3: 재무제표 API (5개) → PR #3

### 메인 클래스 (`korea_investment_stock.py`)

- [ ] `fetch_financial_ratio()` 구현 (TR ID: FHKST66430300)
- [ ] `fetch_income_statement()` 구현 (TR ID: FHKST66430200)
- [ ] `fetch_balance_sheet()` 구현 (TR ID: FHKST66430100)
- [ ] `fetch_profitability_ratio()` 구현 (TR ID: FHKST66430400)
- [ ] `fetch_growth_ratio()` 구현 (TR ID: FHKST66430800)

### 캐시 래퍼 (`cached_korea_investment.py`)

- [ ] `fetch_financial_ratio()` 캐시 래퍼 추가 (TTL: stock_info)
- [ ] `fetch_income_statement()` 캐시 래퍼 추가 (TTL: stock_info)
- [ ] `fetch_balance_sheet()` 캐시 래퍼 추가 (TTL: stock_info)
- [ ] `fetch_profitability_ratio()` 캐시 래퍼 추가 (TTL: stock_info)
- [ ] `fetch_growth_ratio()` 캐시 래퍼 추가 (TTL: stock_info)

### Rate Limit 래퍼 (`rate_limited_korea_investment.py`)

- [ ] `fetch_financial_ratio()` rate limit 래퍼 추가
- [ ] `fetch_income_statement()` rate limit 래퍼 추가
- [ ] `fetch_balance_sheet()` rate limit 래퍼 추가
- [ ] `fetch_profitability_ratio()` rate limit 래퍼 추가
- [ ] `fetch_growth_ratio()` rate limit 래퍼 추가

### 테스트 (`tests/test_financial_apis.py`)

- [ ] `TestFetchFinancialRatio` - 성공 응답, URL/헤더, 파라미터 (연간/분기)
- [ ] `TestFetchIncomeStatement` - 성공 응답, URL/헤더, 파라미터 (연간/분기)
- [ ] `TestFetchBalanceSheet` - 성공 응답, URL/헤더, 파라미터 (연간/분기)
- [ ] `TestFetchProfitabilityRatio` - 성공 응답, URL/헤더, 파라미터 (연간/분기)
- [ ] `TestFetchGrowthRatio` - 성공 응답, URL/헤더, 파라미터 (연간/분기)

### 문서 및 버전

- [ ] `CHANGELOG.md` 업데이트
- [ ] `CLAUDE.md` 메서드 목록 업데이트
- [ ] PR 생성 및 머지

---

## Step 4: 배당 + 업종 API (3~4개) → PR #4

### 메인 클래스 (`korea_investment_stock.py`)

- [ ] `fetch_dividend_ranking()` 구현 (TR ID: HHKDB13470100)
- [ ] `fetch_industry_index()` 구현 (TR ID: FHPUP02100000)
- [ ] `fetch_industry_category_price()` 구현 (TR ID: FHPUP02140000)
- [ ] `fetch_dividend_schedule()` 구현 (API 문서 확인 후, 후순위)

### 상수 (`constants.py`)

- [ ] `INDUSTRY_INDEX_CODE` 추가

### 캐시 래퍼 (`cached_korea_investment.py`)

- [ ] `fetch_dividend_ranking()` 캐시 래퍼 추가 (TTL: stock_info)
- [ ] `fetch_industry_index()` 캐시 래퍼 추가 (TTL: stock_info)
- [ ] `fetch_industry_category_price()` 캐시 래퍼 추가 (TTL: stock_info)
- [ ] `fetch_dividend_schedule()` 캐시 래퍼 추가 (문서 확인 후)

### Rate Limit 래퍼 (`rate_limited_korea_investment.py`)

- [ ] `fetch_dividend_ranking()` rate limit 래퍼 추가
- [ ] `fetch_industry_index()` rate limit 래퍼 추가
- [ ] `fetch_industry_category_price()` rate limit 래퍼 추가
- [ ] `fetch_dividend_schedule()` rate limit 래퍼 추가 (문서 확인 후)

### 테스트 (`tests/test_dividend_industry_apis.py`)

- [ ] `TestFetchDividendRanking` - 성공 응답, URL/헤더, 파라미터
- [ ] `TestFetchIndustryIndex` - 성공 응답, URL/헤더, 파라미터 (업종코드)
- [ ] `TestFetchIndustryCategoryPrice` - 성공 응답, URL/헤더, 파라미터 (시장구분)
- [ ] `TestFetchDividendSchedule` - (문서 확인 후)

### 문서 및 버전

- [ ] `CHANGELOG.md` 업데이트
- [ ] `CLAUDE.md` 메서드 목록 업데이트
- [ ] 버전 업데이트 (PyPI 릴리스용)
- [ ] PR 생성 및 머지

---

## Step 5: DB 스키마 마이그레이션 → PR #5 (moneyflow.advenoh.pe.kr)

### 기존 테이블 DROP

- [ ] `kr_stocks` DROP changelog 작성
- [ ] `kr_etf` DROP changelog 작성
- [ ] `us_stocks` DROP changelog 작성
- [ ] `us_etf` DROP changelog 작성
- [ ] `kr_investor_trading` DROP changelog 작성
- [ ] `market_investor_trend` DROP changelog 작성

### 새 테이블 CREATE (API 원본 필드명)

- [ ] `kr_stocks` 재생성 (API 필드명: stck_prpr, acml_vol 등)
- [ ] `kr_etf` 재생성 (API 필드명)
- [ ] `us_stocks` 재생성 (API 필드명: last, tvol 등)
- [ ] `us_etf` 재생성 (API 필드명)
- [ ] `kr_investor_trading` 재생성 (API 필드명: frgn_ntby_qty 등)
- [ ] `market_investor_trend` 재생성 (API 필드명)

### 새 테이블 추가 (재무제표 등)

- [ ] 재무비율 테이블 생성 (필요 시)
- [ ] 손익계산서 테이블 생성 (필요 시)
- [ ] 대차대조표 테이블 생성 (필요 시)

### 마이그레이션 실행

- [ ] Liquibase changelog 작성
- [ ] DB 마이그레이션 실행 (`./liquibase.sh update`)
- [ ] 마이그레이션 결과 확인

---

## Step 6: stock-data-batch 수정 → PR #6

### DB 모델 수정 (`database/models.py`)

- [ ] `KRStock` 모델 필드명 변경 (API 원본 필드명)
- [ ] `KRETF` 모델 필드명 변경
- [ ] `USStock` 모델 필드명 변경
- [ ] `USETF` 모델 필드명 변경
- [ ] `KRInvestorTrading` 모델 필드명 변경
- [ ] `MarketInvestorTrend` 모델 필드명 변경

### 수집 로직 수정

- [ ] `services/stock_collector.py` 필드 매핑 수정
- [ ] `services/data_saver.py` UPSERT 함수 수정

### korea-investment-stock 의존성

- [ ] `pyproject.toml` 또는 `requirements.txt` 버전 업데이트

### 테스트

- [ ] 배치 테스트 실행 (`python test_integration.py --quick`)
- [ ] 데이터 재수집 확인

---

## Step 7: Backend/Frontend 수정 → 별도 PR

### Backend (Go API)

- [ ] DB 모델 구조체 필드명 수정
- [ ] API 엔드포인트 응답 필드 조정

### Frontend (Next.js)

- [ ] API 타입 정의 수정
- [ ] 컴포넌트 필드 참조 수정
