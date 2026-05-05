# Phase 2 — Read-Only API Extension Design

**Status:** Active design (2026-05-05)
**Goal:** Phase 1 의 28 fetch 메서드 (Python parity) 를 기반으로, Python wrapper 가 cover 하지 않은 KIS 의 추가 read-only API 를 도메인별 sub-phase 로 확장.
**Out of Scope (Phase 2 한정):** 주문/잔고/예약 (Trading), 실시간 WebSocket, 선물옵션, 장내채권. 이들은 별도 Phase 로 분리.

---

## §1. 목적과 결정 요약

### 목적

Phase 1 (`v1.0.0` ~ `v1.3.0`) 으로 라이브러리는 28 메서드 도메인 커버리지를 달성했지만, KIS API 의 read-only 영역에는 아직 다수의 endpoint 가 남아있음 (예: 호가/체결 디테일, 신고저가, 해외 추가 ranking, 11 종 예탁원 일정 etc.). moneyflow 백엔드 / stock-data-batch 가 활용 가능한 데이터를 늘리고, 라이브러리의 KIS API 커버리지를 폭넓게 가져가는 것이 목표.

### 핵심 결정

| 항목 | 결정 |
|---|---|
| Phase 1 패턴 적용 | Style A (path last segment PascalCase), Params struct (zero-value default), `decimal.Decimal`/`int64,string`/`float64,string`/`string` 매핑, Output1+Output2 verbatim. 모든 sub-phase 동일 |
| Sub-phase 분해 | 4 sub-phase (2.1 ~ 2.4), 각 별도 PR + minor release |
| Release tags | minor bump (v1.4.0 ~ v1.7.0). breaking change 없으니 v2 미사용 |
| 명명 충돌 회피 | Phase 1.5 의 `OverseasProductInfo` 처럼 패키지 + struct 명으로 구분 |
| Python parity | 본 phase 는 Python wrapper 가 cover **하지 않은** 영역 — parity 개념 없음. KIS docs 가 source of truth |

---

## §2. Sub-phase 분해

### Phase 2.1 — 국내 호가/체결 (`v1.4.0`)

| 메서드 | 위치 | 한투 path | TR_ID |
|---|---|---|---|
| `Domestic.InquireAskingPriceExpCcn` | `domestic/quote.go` | `quotations/inquire-asking-price-exp-ccn` | FHKST01010200 |
| `Domestic.InquireCcnl` | `domestic/quote.go` | `quotations/inquire-ccnl` | FHKST01010300 |
| `Domestic.InquireDailyPrice` | `domestic/quote.go` | `quotations/inquire-daily-price` | FHKST01010400 |

총 **3 메서드**. 시세의 깊은 분석 (호가 잔량 + 체결 streak + 일자별 주가).

**파일**: `domestic/quote.go` (price/info/chart 와 별도 — 시세 디테일).

### Phase 2.2 — 국내 신고저가 / 시간외 (`v1.5.0`)

| 메서드 | 위치 | 한투 path |
|---|---|---|
| `Domestic.InquireNewHighLowProximity` | `domestic/extended.go` | `ranking/near-new-highlow` (실 path 확인 후) |
| `Domestic.InquireAfterHourPrice` | `domestic/extended.go` | `quotations/...시간외현재가` |
| `Domestic.InquireAfterHourAskingPrice` | `domestic/extended.go` | `quotations/...시간외호가` |
| `Domestic.InquireAfterHourVolumeRank` | `domestic/extended.go` | `ranking/...시간외거래량` |
| `Domestic.InquireAfterHourFluctuationRank` | `domestic/extended.go` | `ranking/...시간외등락` |

총 **5 메서드** (정확한 path/TR_ID 는 implementation 시 docs 에서 확인 — 한글 파일명 매핑 필요).

**docs**: `docs/api/국내주식/{국내주식_시간외*.md, 국내주식_신고_신저근접종목_상위.md, 주식현재가_시간외*.md, 국내주식_시간외예상체결등락률.md}`

**파일**: `domestic/extended.go`

### Phase 2.3 — 해외 추가 Ranking (`v1.6.0`)

| 메서드 | 위치 | 한투 path |
|---|---|---|
| `Overseas.InquireMarketCapRank` | `overseas/ranking.go` (확장) | `ranking/...시가총액` |
| `Overseas.InquireVolumeRank` | `overseas/ranking.go` | `ranking/...거래량` |
| `Overseas.InquireTradeAmountRank` | `overseas/ranking.go` | `ranking/...거래대금` |
| `Overseas.InquireVolumeSurge` | `overseas/ranking.go` | `ranking/...거래량급증` |
| `Overseas.InquireBuyExecStrength` | `overseas/ranking.go` | `ranking/...매수체결강도` |
| `Overseas.InquireNewHighLow` | `overseas/ranking.go` | `ranking/...신고_신저` |

총 **6 메서드**. Phase 1.5 의 `InquireUpdownRate` 패턴 그대로 적용 (Output1 summary + Output2 array). `overseas/ranking.go` 에 누적.

**docs**: `docs/api/해외주식/{해외주식_시가총액순위.md, 해외주식_거래량순위.md, 해외주식_거래대금순위.md, 해외주식_거래량급증.md, 해외주식_매수체결강도상위.md, 해외주식_신고_신저가.md}`

### Phase 2.4 — 예탁원 정보 확장 (`v1.7.0`)

Phase 1.4 의 `InquirePubOffer` (공모주청약일정) 외 11 개 KSD 일정 endpoint. 모두 `ksdinfo/<endpoint>` path.

| 메서드 | docs |
|---|---|
| `Domestic.InquireKsdRights` (배당일정) | 예탁원정보(배당일정).md |
| `Domestic.InquireKsdBonus` (무상증자일정) | 예탁원정보(무상증자일정).md |
| `Domestic.InquireKsdRightsOffer` (유상증자일정) | 예탁원정보(유상증자일정).md |
| `Domestic.InquireKsdMeeting` (주주총회일정) | 예탁원정보(주주총회일정).md |
| `Domestic.InquireKsdMerge` (합병_분할일정) | 예탁원정보(합병_분할일정).md |
| `Domestic.InquireKsdParChange` (액면교체일정) | 예탁원정보(액면교체일정).md |
| `Domestic.InquireKsdRevoked` (실권주일정) | 예탁원정보(실권주일정).md |
| `Domestic.InquireKsdEscrow` (의무예치일정) | 예탁원정보(의무예치일정).md |
| `Domestic.InquireKsdCapReduce` (자본감소일정) | 예탁원정보(자본감소일정).md |
| `Domestic.InquireKsdAppraisal` (주식매수청구일정) | 예탁원정보(주식매수청구일정).md |
| `Domestic.InquireKsdListingNotice` (상장정보일정) | 예탁원정보(상장정보일정).md |

총 **11 메서드**. 메서드 명명 prefix `InquireKsd*` (Phase 1.4 의 `InquirePubOffer` 와 일관 — `ksdinfo/` path 의 last segment 가 다양해서 prefix 통합 가독성 우선). 정확한 path/TR_ID 는 implementation 시 확인.

**파일**: `domestic/ksd.go` (Phase 1.4 의 `domestic/ipo.go` 와 분리 — 11 메서드는 별도 file 가독성 우선).

### Phase 2.5+ (미정)

| 후보 | 설명 |
|---|---|
| 외인기관 추정가 | `종목별_외인기관_추정가집계`, `국내기관_외국인_매매종목가집계` |
| 해외뉴스/속보 | `해외뉴스종합(제목)`, `해외속보(제목)` |
| 권리/잔고 디테일 | `해외주식_권리종합`, `해외주식_기간별권리조회` |
| 업종 chart | `국내업종_시간별지수(분/초)`, `국내업종_일자별지수`, `업종_분봉조회` |
| 프로그램 매매 | `프로그램매매_투자자매매동향(당일)` |

Phase 2.4 완료 후 우선순위 재검토.

### 합계 (Phase 2.1 ~ 2.4)

- **메서드**: 3 + 5 + 6 + 11 = **25 메서드** (누적 28 → 53)
- **Release tags**: `v1.4.0` (2.1), `v1.5.0` (2.2), `v1.6.0` (2.3), `v1.7.0` (2.4)
- **PR**: 4 sub-phase = 4 PR

---

## §3. 명명 / 패턴 (Phase 1 그대로)

### Style A 메서드 명명

한투 endpoint path 의 last segment 를 PascalCase 로 1:1 매핑. 단, Phase 2.4 처럼 path 가 복잡하거나 의미 약하면 prefix (`InquireKsd*`) 통합 가독성 우선.

### 응답 typed struct

KIS docs 1:1 매핑 (PascalCase + 한투 약어 보존). Output1+Output2 verbatim. 충돌 시 패키지 + 명시 prefix (e.g., `OverseasMarketCap` vs `MarketCap`).

### Params struct

`Inquire<X>Params` (zero-value default — 빈 값 시 KIS docs 의 default 적용).

### 타입 매핑

- **가격/지수 → `decimal.Decimal` (bare tag)**
- **수량/금액 → `int64,string`**
- **비율 → `float64,string`**
- **코드/이름/날짜/Y-N → 평문 `string`**

---

## §4. 인프라 변경 (없음)

Phase 1 의 인프라 (`internal/{httpclient,ratelimit,token,mastercache,krxmaster,overseasmaster}`) 그대로 사용. 새 internal package 불필요 (모든 endpoint 가 KIS REST API).

`domestic.New(http, master)` / `overseas.New(http, master)` 시그니처 변경 없음.

---

## §5. 테스트 / 문서 / Release 흐름 (Phase 1 동일)

각 sub-phase 별 implementation plan 에서 다음을 따름:

1. testdata fixtures (KIS docs 의 응답 필드 정의 기반 합성 JSON)
2. TDD: failing test → struct + method 구현 → PASS → commit
3. examples (`examples/<sub-phase>/main.go`)
4. CLAUDE.md / README.md / CHANGELOG.md / `<package>/doc.go` 갱신
5. 최종 점검 (build/vet/fmt/race/coverage ≥ 80%)
6. PR 생성 (사용자 승인 후), merge, tag, GitHub Release

---

## §6. 진입/종료 조건

### 진입 조건

- main HEAD = v1.3.0 (Phase 1.5 publish 완료)
- Phase 2 design spec (이 문서) 사용자 승인

### 종료 조건 (각 sub-phase)

- PR merge 완료, CI clean
- minor version tag push (v1.4.0 ~ v1.7.0)
- GitHub Release publish
- memory 갱신 (다음 sub-phase 시작 절차)

### Phase 2 전체 종료 조건

- Phase 2.1 ~ 2.4 모두 publish (v1.7.0)
- 누적 53 메서드 커버리지
- 사용자가 Phase 2.5+ 진행 또는 다른 도메인 (Trading/WebSocket) 으로 전환 결정
