# 시장별 투자자매매동향(시세) API 구현 PRD

## 개요

시장별 투자자매매동향(시세) API를 korea_investment_stock 라이브러리에 추가하고, stock-data-batch에서 해당 데이터를 DB에 저장하는 기능을 구현합니다.

## API 정보

| 항목 | 값 |
|------|-----|
| API명 | 시장별 투자자매매동향(시세) |
| API ID | v1_국내주식-074 |
| API 경로 | `/uapi/domestic-stock/v1/quotations/inquire-investor-time-by-market` |
| TR ID | FHPTJ04030000 |
| HTS 화면 | [0403] 시장별 시간동향 |
| 모의투자 | 미지원 |

---

## 요청 파라미터

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| fid_input_iscd | Y | 시장구분 (KSP, KSQ 등) |
| fid_input_iscd_2 | Y | 업종구분 (0001, 1001 등) |

### 시장 코드 (market_code)

| 코드 | 설명 |
|------|------|
| KSP | 코스피 |
| KSQ | 코스닥 |
| ETF | ETF |
| ELW | ELW |
| ETN | ETN |
| K2I | 선물/콜옵션/풋옵션 |
| 999 | 주식선물 |
| MKI | 미니 |
| WKM | 위클리(월) |
| WKI | 위클리(목) |
| KQI | 코스닥150 |

### 업종 코드 (sector_code)

| 코드 | 설명 |
|------|------|
| 0001 | 코스피 종합 |
| 1001 | 코스닥 종합 |
| F001 | 선물 |
| OC01 | 콜옵션 |
| OP01 | 풋옵션 |
| T000 | ETF 전체 |
| W000 | ELW 전체 |
| E199 | ETN 전체 |

---

## 응답 데이터 필드

### 투자자 유형별 접두사

| 접두사 | 투자자 유형 |
|--------|------------|
| `frgn_*` | 외국인 |
| `prsn_*` | 개인 |
| `orgn_*` | 기관계 |
| `scrt_*` | 증권 |
| `ivtr_*` | 투자신탁 |
| `pe_fund_*` | 사모펀드 |
| `bank_*` | 은행 |
| `insu_*` | 보험 |
| `mrbn_*` | 종금 |
| `fund_*` | 기금 |
| `etc_orgt_*` | 기타단체 |
| `etc_corp_*` | 기타법인 |

### 세부 필드 (접미사)

| 필드 접미사 | 설명 |
|------------|------|
| `_seln_vol` | 매도 거래량 |
| `_shnu_vol` | 매수 거래량 |
| `_ntby_qty` | 순매수 수량 |
| `_seln_tr_pbmn` | 매도 거래 대금 |
| `_shnu_tr_pbmn` | 매수 거래 대금 |
| `_ntby_tr_pbmn` | 순매수 거래 대금 |

---

## 구현 범위

### Phase 1: korea_investment_stock 라이브러리

- `fetch_investor_trend_by_market(market_code, sector_code)` 메서드 추가
- 시장/업종 코드 상수 추가 (`constants.py`)
- 테스트 코드 작성
- 버전 업데이트 (0.16.0)

### Phase 2: DB 스키마 (moneyflow.advenoh.pe.kr)

- Liquibase changelog 파일 작성
- DB 마이그레이션 실행

### Phase 3: stock-data-batch

- `MarketInvestorTrend` DB 모델 추가
- API 응답 매핑 함수
- 데이터 수집/저장 함수
- 배치 처리 함수
- CLI 옵션 (`--market-investor-trend`)

### Phase 4: 배포 및 통합 테스트

- korea_investment_stock: PR merge 후 GitHub Actions `release.yml` 실행
- moneyflow.advenoh.pe.kr: DB 마이그레이션 (Liquibase)
- stock-data-batch: 의존성 업데이트 및 배치 테스트

---

## 수집 대상

| 시장 | 업종 | 설명 |
|------|------|------|
| KSP | 0001 | 코스피 종합 |
| KSQ | 1001 | 코스닥 종합 |
| ETF | T000 | ETF 전체 |

---

## 관련 문서

- [구현 문서](4_investor_time_implementation.md)
- [TODO 체크리스트](4_investor_time_todo.md)
