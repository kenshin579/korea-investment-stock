# Phase 2.5+ — Read-Only Extension Continuation Design

**Status:** Active design (2026-05-05)
**Goal:** Phase 2 (4 sub-phase, 25 methods) 완료 후 read-only API 확장 계속. 5개 도메인 (외인기관, 해외뉴스, 해외권리, 업종/지수, 프로그램매매) 의 REST endpoint 20 메서드 추가.
**Out of Scope:** 실시간 (WebSocket) 형태의 endpoint (`국내지수_실시간*`, `국내주식_실시간프로그램매매_*` 등) — 별도 WebSocket phase 로 분리.

---

## §1. 목적과 결정 요약

### 목적

Phase 2 (`v1.4.0` ~ `v1.7.0`) 으로 누적 53 메서드 (Phase 1: 28 Python parity + Phase 2: 25 read-only 확장) 도달. 원본 Phase 2 design spec (`2026-05-05-phase2-readonly-extension-design.md`) §2.5+ 에 "미정" 으로 남겨둔 후보 도메인을 sub-phase 로 구체화.

### 핵심 결정

| 항목 | 결정 |
|---|---|
| Phase 2 패턴 적용 | Style A (path last segment PascalCase), Params struct (zero-value default), `decimal.Decimal`/`int64,string`/`float64,string`/`string` 매핑, Output1+Output2 verbatim. 모든 sub-phase 동일. |
| Sub-phase 분해 | 3 sub-phase (2.5 ~ 2.7), 각 별도 PR + minor release |
| Release tags | minor bump (v1.8.0 ~ v1.10.0). breaking change 없으니 v2 미사용. |
| 실시간 endpoint 제외 | KIS 의 `/tryitout/H*` 형태 (POST + 실시간) 는 WebSocket 도메인 — 별도 phase. Phase 2.5+ 는 REST `GET` 만. |

---

## §2. Sub-phase 분해

### Phase 2.5 — 투자자/매매 동향 (`v1.8.0`)

| 메서드 | 위치 | 한투 docs |
|---|---|---|
| `Domestic.InquireForeignInstEstimate` | `domestic/investor.go` (append) | `종목별_외인기관_추정가집계.md` |
| `Domestic.InquireDomInstForeignTradeAgg` | `domestic/investor.go` (append) | `국내기관_외국인_매매종목가집계.md` |
| `Domestic.InquireProgramTradeDaily` | `domestic/program_trade.go` (신규) | `종목별_프로그램매매추이(일별).md` |
| `Domestic.InquireProgramTradeExec` | `domestic/program_trade.go` | `종목별_프로그램매매추이(체결).md` |
| `Domestic.InquireProgramTradeSummaryHourly` | `domestic/program_trade.go` | `프로그램매매_종합현황(시간).md` |
| `Domestic.InquireProgramTradeSummaryDaily` | `domestic/program_trade.go` | `프로그램매매_종합현황(일별).md` |
| `Domestic.InquireProgramTradeInvestorFlow` | `domestic/program_trade.go` | `프로그램매매_투자자매매동향(당일).md` |

총 **7 메서드** (외인기관/매매집계 2 + 프로그램매매 5). 투자자 종합 흐름.

**파일**:
- 수정: `domestic/investor.go` (Phase 1.4 의 3 investor 메서드 옆에 2 메서드 append)
- 신규: `domestic/program_trade.go` (5 프로그램매매 메서드)

### Phase 2.6 — 해외 정보 (`v1.9.0`)

| 메서드 | 위치 | 한투 docs |
|---|---|---|
| `Overseas.InquireNewsTitle` | `overseas/news.go` (신규) | `해외뉴스종합(제목).md` |
| `Overseas.InquireBreakingNewsTitle` | `overseas/news.go` | `해외속보(제목).md` |
| `Overseas.InquireRightsSummary` | `overseas/rights.go` (신규) | `해외주식_권리종합.md` |
| `Overseas.InquireRightsByPeriod` | `overseas/rights.go` | `해외주식_기간별권리조회.md` |

총 **4 메서드** (뉴스 2 + 권리 2). 해외주식 보조 정보.

**파일**:
- 신규: `overseas/news.go` (2 뉴스 메서드)
- 신규: `overseas/rights.go` (2 권리 메서드)

### Phase 2.7 — 업종/지수 (`v1.10.0`)

| 메서드 | 위치 | 한투 docs |
|---|---|---|
| `Domestic.InquireIndustryCurrent` | `domestic/industry.go` (신규) | `국내업종_현재지수.md` |
| `Domestic.InquireIndustryFullPrice` | `domestic/industry.go` | `국내업종_구분별전체시세.md` |
| `Domestic.InquireIndustryDaily` | `domestic/industry.go` | `국내업종_일자별지수.md` |
| `Domestic.InquireIndustryMinute` | `domestic/industry.go` | `국내업종_시간별지수(분).md` |
| `Domestic.InquireIndustrySecond` | `domestic/industry.go` | `국내업종_시간별지수(초).md` |
| `Domestic.InquireIndustryPeriodPrice` | `domestic/industry.go` | `국내주식업종기간별시세(일_주_월_년).md` |
| `Domestic.InquireIndustryMinuteChart` | `domestic/industry.go` | `업종_분봉조회.md` |
| `Domestic.InquireExpectedIndexAll` | `domestic/industry.go` | `국내주식_예상체결_전체지수.md` |
| `Domestic.InquireExpectedIndexTrend` | `domestic/industry.go` | `국내주식_예상체결지수_추이.md` |

총 **9 메서드**. 업종/지수 종합 (현재/일자별/분/초 + 기간별 + 분봉 + 예상체결 지수).

**파일**: `domestic/industry.go` (신규, 9 메서드)

> 메서드 명명 일부 조정: KIS docs 의 한글 파일명 자체로는 path last segment 가 곧바로 영문으로 매핑되지 않으므로, implementation 시 docs 의 `URL` (e.g., `inquire-index-price`) 을 확인해 Style A 적용. 위 표는 후보 명명이며 실 path 검증 후 확정.

### 합계 (Phase 2.5 ~ 2.7)

- **메서드**: 7 + 4 + 9 = **20 메서드** (누적 53 → 73)
- **Release tags**: `v1.8.0` (2.5), `v1.9.0` (2.6), `v1.10.0` (2.7)
- **PR**: 3 sub-phase = 3 PR

---

## §3. 명명 / 패턴 (Phase 2 그대로)

### Style A 메서드 명명

한투 endpoint path 의 last segment 를 PascalCase 로 1:1 매핑. 단, path 가 모호하거나 의미 약하면 prefix 통합 (Phase 1.4 의 `InquireKsd*`, Phase 2.4 적용 사례) 가능.

### 응답 typed struct

KIS docs 1:1 매핑 (PascalCase + 한투 약어 보존). Output1+Output2 verbatim. 충돌 시 패키지 + 명시 prefix.

### Params struct

`Inquire<X>Params` (zero-value default — 빈 값 시 KIS docs 의 default 적용).

### 타입 매핑

- **가격/지수 → `decimal.Decimal` (bare tag)**
- **수량/금액 → `int64,string`**
- **비율 → `float64,string`**
- **코드/이름/날짜/Y-N → 평문 `string`**

KSD endpoint 의 all-string convention (Phase 2.4) 은 KSD 한정 — Phase 2.5+ 는 일반 매핑 적용.

---

## §4. 인프라 변경 (없음)

Phase 1 의 인프라 (`internal/{httpclient,ratelimit,token,mastercache,krxmaster,overseasmaster}`) 그대로 사용. 새 internal package 불필요.

`domestic.New(http, master)` / `overseas.New(http, master)` 시그니처 변경 없음.

---

## §5. 테스트 / 문서 / Release 흐름 (Phase 2 동일)

각 sub-phase 별 implementation plan 에서:

1. testdata fixtures (KIS docs 응답 필드 정의 기반 합성 JSON, **rt_cd/msg_cd/msg1 envelope 필수** — Phase 2.4 hazard 반영)
2. TDD: failing test → struct + method 구현 → PASS → commit
3. examples (`examples/<sub-phase>/main.go`)
4. CLAUDE.md / README.md / CHANGELOG.md / `<package>/doc.go` 갱신
5. 최종 점검 (build/vet/fmt/race/coverage ≥ 80%)
6. PR 생성 (사용자 승인 후), merge, tag, GitHub Release

---

## §6. 진입/종료 조건

### 진입 조건

- main HEAD = v1.7.0 (Phase 2.4 publish 완료)
- Phase 2.5+ design spec (이 문서) 사용자 승인

### 종료 조건 (각 sub-phase)

- PR merge 완료, CI clean
- minor version tag push (v1.8.0 ~ v1.10.0)
- GitHub Release publish
- memory 갱신 (다음 sub-phase 시작 절차)

### Phase 2.5+ 전체 종료 조건

- Phase 2.5 ~ 2.7 모두 publish (v1.10.0)
- 누적 73 메서드 커버리지
- 사용자가 다음 도메인 (Trading/WebSocket/선물옵션/장내채권/v1.x 유지보수) 으로 전환 결정
