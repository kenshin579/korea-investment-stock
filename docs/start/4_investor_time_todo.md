# 시장별 투자자매매동향(시세) API 구현 TODO

## Phase 1: korea_investment_stock 라이브러리

### API 구현

- [x] `fetch_investor_trend_by_market()` 메서드 추가
  - [x] API 경로 및 TR ID 설정
  - [x] 요청 파라미터 (market_code, sector_code) 처리
  - [x] `_request_with_token_refresh()` 사용

### 상수 정의

- [x] `constants.py`에 시장 코드 상수 추가 (`MARKET_INVESTOR_TREND_CODE`)
- [x] `constants.py`에 업종 코드 상수 추가 (`SECTOR_CODE`)

### 테스트

- [x] `tests/test_investor_trend_by_market.py` 파일 생성
- [x] 코스피 종합 조회 테스트 작성
- [x] 코스닥 종합 조회 테스트 작성

### 배포

- [x] `pyproject.toml` 버전 업데이트 (setuptools_scm으로 Git 태그 사용)
- [x] CHANGELOG.md 업데이트
- [x] PR 생성 및 merge → PR #121: https://github.com/kenshin579/korea-investment-stock/pull/121
- [x] PyPI 배포 (v0.16.1 릴리스 완료)

---

## Phase 2: DB 스키마 (moneyflow.advenoh.pe.kr)

- [x] Liquibase changelog 파일 작성 (`8_create_market_investor_trend.sql`)
  - PR #61: https://github.com/kenshin579/moneyflow.advenoh.pe.kr/pull/61
- [x] DB 마이그레이션 실행 (`./liquibase.sh update-one`)

---

## Phase 3: stock-data-batch 수정

### DB 모델

- [x] `MarketInvestorTrend` 모델 클래스 추가

### 데이터 수집

- [x] `_map_market_investor_trend()` 매핑 함수 추가
- [x] `fetch_market_investor_trend_data()` 수집 함수 추가
- [x] `upsert_market_investor_trend()` 저장 함수 추가

### 배치 처리

- [x] `process_market_investor_trend()` 배치 함수 추가
- [x] CLI 옵션 `--market-investor-trend` 추가

### 의존성

- [x] `pyproject.toml`에서 `korea-investment-stock>=0.16.1` 업데이트
- PR #70: https://github.com/kenshin579/stock-data-batch/pull/70

---

## Phase 4: 배포 및 통합 테스트

- [x] korea_investment_stock PR merge 및 `release.yml` 실행 (v0.16.1)
- [x] moneyflow.advenoh.pe.kr DB 마이그레이션 실행
- [ ] stock-data-batch 배치 테스트 (`python main.py --market-investor-trend`)
- [ ] DB 저장 결과 확인
