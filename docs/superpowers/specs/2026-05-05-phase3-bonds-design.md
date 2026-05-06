# Phase 3 — 장내채권 (Bonds) Domain Design

**Status:** Active design (2026-05-05)
**Goal:** 장내채권 도메인 신규 sub-package 도입. Phase 3.1 (시세 8 메서드) 우선 진행. 잔고/주문/실시간은 나중 결정.
**Out of Scope (Phase 3.1 한정):** 잔고/주문조회 (계좌 인증 필요), 주문/거래 (HIGH RISK), 실시간 (WebSocket).

---

## §1. 목적과 결정 요약

### 목적

KIS API 의 18 채권 docs 중 **시세/조회 8 메서드** 만 우선 cover. 주식 (Phase 1+2+2.5+) 의 시세 패턴을 그대로 적용. 새 `bonds/` sub-package 도입.

채권 시세는 moneyflow / stock-data-batch 의 잠재 use case (채권 평가단가, 만기일 알림 등). 위험도 LOW — 시세 조회만 하면 계좌 인증 불필요.

### 핵심 결정

| 항목 | 결정 |
|---|---|
| 새 sub-package | `bonds/` (Go module path: `github.com/kenshin579/korea-investment-stock/bonds`). KIS API 의 `/uapi/domestic-bond/v1/` 와 별개로 Go 명명 단순화. |
| Client 통합 | root `client.go` 의 `Client` 에 `Bonds *bonds.Client` 필드 추가. `wireInfra` 에서 주입. |
| 호출 스타일 | `client.Bonds.InquirePrice(ctx, ...)` (도메인 1-level grouping — Phase 1+2 일관) |
| Sub-phase 분해 | Phase 3.1 (시세 8) 만 우선. 나머지는 사용자 검토 후 결정 |
| Style A 명명 | KIS path 의 last segment PascalCase (예외: `search-bond-info` → `SearchBondInfo` — `Inquire` prefix 강제 안 함) |
| 타입 매핑 | Phase 2 standard. 단 EP1+EP2 (search-bond-info, issue-info) 는 KIS docs 의 all-string 명시에 따라 plain `string` (KSD-like) |
| 인증 | OAuth access token only (계좌 CANO 불필요) |

---

## §2. Sub-phase 분해

### Phase 3.1 — 채권 시세 (`v1.11.0`)

| 메서드 | path (last seg) | TR_ID | output |
|---|---|---|---|
| `Bonds.SearchBondInfo` | `search-bond-info` | CTPF1114R | `output{}` 70 fields all-string |
| `Bonds.InquireIssueInfo` | `issue-info` | CTPF1101R | `output{}` 69 fields all-string |
| `Bonds.InquirePrice` | `inquire-price` | FHKBJ773400C0 | `output{}` 17 fields typed |
| `Bonds.InquireCcnl` | `inquire-ccnl` | FHKBJ773403C0 | `output{}` 7 fields typed (single snapshot) |
| `Bonds.InquireAskingPrice` | `inquire-asking-price` | FHKBJ773401C0 | `output{}` 34 fields typed (5 단계 호가) |
| `Bonds.InquireDailyPrice` | `inquire-daily-price` | FHKBJ773404C0 | `output{}` 9 fields typed |
| `Bonds.InquireDailyItemchartprice` | `inquire-daily-itemchartprice` | FHKBJ773701C0 | `output[]` 6 fields/item |
| `Bonds.InquireAvgUnit` | `avg-unit` | CTPF2005R | `output1+output2+output3` arrays (23+10+16) |

총 **8 메서드**. 누적 71 → 79.

**파일**: `bonds/client.go` (Client + New), `bonds/quote.go` (8 메서드 + structs), `bonds/quote_test.go`. Root `client.go` 에 `Bonds *bonds.Client` 필드 + `wireInfra` 통합.

### Phase 3.2 — 채권 잔고/주문조회 (미정)

4 메서드 (잔고/매수가능/주문체결/정정취소가능). **계좌 인증 (CANO/ACNT_PRDT_CD) 인프라 도입 필요** — Trading 도메인과 공유. Phase 3.1 publish 후 우선순위 재검토.

### Phase 3.3 — 채권 주문/거래 (HIGH RISK, 미정)

3 메서드 (매수/매도/정정취소). 별도 design spec + safety review 필요.

### Phase 3.4 — 채권 실시간 (WebSocket, 미정)

3 메서드 (실시간 체결/호가/지수). WebSocket architecture phase 와 묶음.

---

## §3. 명명 / 패턴

### Style A 메서드 명명

KIS path 의 last segment PascalCase. 단:
- `search-bond-info` → `SearchBondInfo` (path 자체에 동사 `search` 포함 → `Inquire` prefix 강제 안 함)
- `avg-unit` → `InquireAvgUnit` (기존 `Inquire*` 패턴과 일관)
- 패키지 경계로 다른 도메인 (`domestic.InquireCcnl`, `domestic.InquireDailyPrice` 등) 과 동일한 메서드명 충돌 OK

### 응답 typed struct

KIS docs 1:1 매핑. Output1+Output2+Output3 (EP8) verbatim. EP1+EP2 의 70/69 large struct 는 KSD-style all-string.

### Params struct

`Inquire<X>Params` / `Search<X>Params` (zero-value default).

### 타입 매핑

EP1+EP2 (CTPF11* TR_IDs, KSD-like reference data): **모든 필드 plain `string`** — KIS docs 가 모두 `String` 타입으로 명시.

EP3-EP8 (FHKBJ* market data + CTPF2005R avg-unit): **Phase 2 standard typed mapping**:
- 채권 가격 (`bond_prpr`, `bond_oprc`, `bond_hgpr`, `bond_lwpr`, `bond_prdy_vrss`, `bond_mxpr`, `bond_llam`, `bond_askp1-5`, `bond_bidp1-5`, `kis_unpr`, `kbp_unpr`, `nice_evlu_unpr`, `fnp_unpr`, `avg_evlu_unpr`, `kis_rf_unpr` 등) → `decimal.Decimal` (bare)
- 평가금액 → `int64,string`
- 거래량/잔량 → `int64,string`
- 비율/수익률/등락률 → `float64,string`
- 코드/날짜/시간/통화코드/Y-N → 평문 `string`

---

## §4. 인프라 변경

**최소 변경**:

1. 새 패키지 `bonds/`:
   - `bonds/client.go` — Client struct + New(http) constructor (master cache 불필요 — 채권은 KRX master 형식과 다름)
   - `bonds/quote.go` — 8 메서드 + 응답 structs + Params

2. Root `client.go`:
   - `type Client struct { ... Bonds *bonds.Client }`
   - `wireInfra` 에서 `c.Bonds = bonds.New(http)` 추가

새 internal package 불필요. Mastercache 미사용 (채권 종목 코드는 KIS docs 에서 별도 master 미제공 — 사용자가 PDNO 직접 알고 있다고 가정).

---

## §5. 테스트 / 문서 / Release 흐름 (Phase 2 동일)

각 sub-phase 별 implementation plan 에서:

1. testdata fixtures (rt_cd/msg_cd/msg1 envelope 필수)
2. TDD: failing test → struct + method 구현 → PASS → commit
3. examples (`examples/bonds_quote/main.go`)
4. CLAUDE.md / README.md / CHANGELOG.md / `bonds/doc.go` 갱신
5. 최종 점검 (build/vet/fmt/race/coverage ≥ 80%)
6. PR 생성 (사용자 승인 후), merge, tag, GitHub Release

---

## §6. 진입/종료 조건

### 진입 조건

- main HEAD = v1.10.0 (Phase 2.5+ 완료, 누적 71 메서드)
- Phase 3 design spec (이 문서) 사용자 승인

### Phase 3.1 종료 조건

- PR merge 완료, CI clean
- minor version tag push (v1.11.0)
- GitHub Release publish
- memory 갱신 (Phase 3.2 진행 또는 다른 도메인 결정)

### Phase 3 전체 향후 결정

- Phase 3.2-3.4 는 별도 사용자 결정 시점에 재검토. Phase 3.1 단독 publish 도 가능.
