# Phase 2.3 — 해외주식 추가 Ranking (6 메서드) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** `korea-investment-stock` Go 라이브러리에 해외주식 ranking 6 메서드 추가 (`v1.6.0` release). 기존 `overseas/ranking.go` 에 APPEND (신규 파일 생성 금지).

**Architecture:** Phase 1.5 / Phase 2.2 의 인프라 + 패턴 재사용. `overseas/ranking.go` (기존) 에 6 메서드 APPEND. 한투 API path 1:1 매핑 (Style A). 새 internal package 불필요. TDD: testdata fixture (KIS docs 응답 필드 기반 합성 JSON) → 실패 테스트 → struct + 메서드 구현 → 통과 → commit.

**Tech Stack:** Go 1.25+, `github.com/jarcoal/httpmock`, `github.com/shopspring/decimal`, `github.com/stretchr/testify`. 새 dependency 없음.

**참고 spec:**
- Phase 2 design spec: `docs/superpowers/specs/2026-05-05-phase2-readonly-extension-design.md`
- Phase 2.2 plan (참조 패턴): `docs/superpowers/specs/2026-05-05-phase2-2-extended-implementation-plan.md`
- 한투 API docs: `docs/api/해외주식/{해외주식_시가총액순위.md, 해외주식_거래량순위.md, 해외주식_거래대금순위.md, 해외주식_거래량급증.md, 해외주식_매수체결강도상위.md, 해외주식_신고_신저가.md}`

---

## 사전 정보

| 항목 | 값 |
|------|---|
| 작업 브랜치 | `feat/phase2-3-overseas-ranking` |
| 시작 HEAD | Phase 2.2 구현 완료 commit (v1.5.0) |
| Release 목표 | `v1.6.0` |
| PR 베이스 | `main` |
| 현재 main 상태 | v1.5.0 publish 완료 (Phase 2.2 통합, 36 메서드) |

---

## 메서드 → 한투 API 매핑

| Go 메서드 | 한투 path | TR_ID | 응답 구조 |
|---|---|---|---|
| `Overseas.InquireMarketCap(ctx, params)` | `/uapi/overseas-stock/v1/ranking/market-cap` | HHDFS76350100 | `output1: {}` (5 fields) + `output2: []` (15 fields/item) |
| `Overseas.InquireTradeVol(ctx, params)` | `/uapi/overseas-stock/v1/ranking/trade-vol` | HHDFS76310010 | `output1: {}` (5 fields) + `output2: []` (16 fields/item) |
| `Overseas.InquireTradePbmn(ctx, params)` | `/uapi/overseas-stock/v1/ranking/trade-pbmn` | HHDFS76320010 | `output1: {}` (5 fields) + `output2: []` (16 fields/item) |
| `Overseas.InquireVolumeSurge(ctx, params)` | `/uapi/overseas-stock/v1/ranking/volume-surge` | HHDFS76270000 | `output1: {}` (3 fields) + `output2: []` (16 fields/item) |
| `Overseas.InquireVolumePower(ctx, params)` | `/uapi/overseas-stock/v1/ranking/volume-power` | HHDFS76280000 | `output1: {}` (3 fields) + `output2: []` (15 fields/item) |
| `Overseas.InquireNewHighlow(ctx, params)` | `/uapi/overseas-stock/v1/ranking/new-highlow` | HHDFS76300000 | `output1: {}` (3 fields) + `output2: []` (16 fields/item) |

> **명명 노트:** Phase 2 design spec §2.3 은 메서드 명을 `InquireMarketCapRank`, `InquireVolumeRank` 등 의미 기반으로 기술했으나, §3 "Style A — path last segment PascalCase 1:1 매핑" 원칙을 따름. 실제 path last segment (`market-cap`, `trade-vol`, `trade-pbmn`, `volume-surge`, `volume-power`, `new-highlow`) PascalCase 변환 적용. 결과: `InquireMarketCap`, `InquireTradeVol`, `InquireTradePbmn`, `InquireVolumeSurge`, `InquireVolumePower`, `InquireNewHighlow`.

---

## 파일 구조

### 수정 (overseas)
- `overseas/ranking.go` — 6 메서드 + structs + Params **APPEND** (기존 `InquireUpdownRate` 유지)
- `overseas/ranking_test.go` — 6 테스트 케이스 **APPEND** (기존 `TestClient_InquireUpdownRate` 유지)

### 신규 (testdata)
- `overseas/testdata/market_cap_success.json`
- `overseas/testdata/trade_vol_success.json`
- `overseas/testdata/trade_pbmn_success.json`
- `overseas/testdata/volume_surge_success.json`
- `overseas/testdata/volume_power_success.json`
- `overseas/testdata/new_highlow_success.json`

### 신규 (examples)
- `examples/overseas_ranking/main.go` — 6 메서드 통합 예시

### 수정 (root)
- `CLAUDE.md` — banner 갱신 (v1.5.0 → v1.6.0, "Phase 2.2" → "Phase 2.3"), Phase 2.3 plan link bullet 추가
- `README.md` — Available Methods 표 갱신 (36 → 42 메서드, heading 갱신)
- `CHANGELOG.md` — `[1.6.0]` entry (above `[1.5.0]`)
- `overseas/doc.go` — Phase 2.3 section (after existing Phase 1.5 section)

---

## 타입 매핑

Phase 1.5 / Phase 2.2 와 동일 원칙:

- **가격류 → `decimal.Decimal` (bare tag)**: `last`, `diff`, `pask`, `pbid`, `tomv`, `n_base`, `n_diff` (신고신저), `n_diff` (거래량급증)
- **수량/금액 → `int64,string`**: `tvol`, `tamt`, `a_tvol`, `a_tamt`, `shar`, `n_tvol`, `crec`, `trec`, `nrec`
- **비율/강도 → `float64,string`**: `rate`, `grav`, `n_rate`, `tpow`, `powx`
- **순위 → `int64,string`**: `rank` (안전을 위해 int64 사용 — KIS는 String 6자리)
- **코드/이름/기호/Y-N → 평문 `string`**: `rsym`, `excd`, `symb`, `name`, `ename`, `knam`, `enam`, `sign`, `e_ordyn`, `zdiv`, `stat`

### output1 구조체 2 Tier

| Tier | 해당 메서드 | 필드 |
|---|---|---|
| `OverseasRankingFullSummary` (5-field) | InquireMarketCap, InquireTradeVol, InquireTradePbmn | `zdiv`, `stat`, `crec`, `trec`, `nrec` |
| `OverseasRankingMinSummary` (3-field) | InquireVolumeSurge, InquireVolumePower, InquireNewHighlow | `zdiv`, `stat`, `nrec` |

> **이유:** VolumeSurge/VolumePower/NewHighlow 의 KIS output1 에는 `crec`/`trec` 필드가 없음 (3 fields only). 두 struct 를 별도로 정의해야 JSON 역직렬화 오류 없음.

### output2 name 필드 변형

| 그룹 | 종목명 필드 | 해당 메서드 |
|---|---|---|
| `name` + `ename` | 한글 종목명 + 영문 종목명 | InquireMarketCap, InquireTradeVol, InquireTradePbmn, InquireNewHighlow |
| `knam` + `enam` | 한글 종목명 + 영문 종목명 | InquireVolumeSurge, InquireVolumePower |

> **이유:** KIS가 endpoint 별로 다른 JSON 키를 사용. 단일 output2 struct 공유 불가. 메서드별 전용 struct 정의 필수.

### output2 공통 핵심 필드 (6 메서드 전체)

```go
Rsym   string          `json:"rsym"`         // 실시간조회심볼
Excd   string          `json:"excd"`         // 거래소코드
Symb   string          `json:"symb"`         // 종목코드
Last   decimal.Decimal `json:"last"`         // 현재가
Sign   string          `json:"sign"`         // 기호
Diff   decimal.Decimal `json:"diff"`         // 대비
Rate   float64         `json:"rate,string"`  // 등락율
Tvol   int64           `json:"tvol,string"`  // 거래량
EOrdyn string          `json:"e_ordyn"`      // 매매가능
// pask/pbid: #2-#6 에만 존재 (InquireMarketCap output2 에는 없음)
```

### output2 endpoint-unique 필드

| 메서드 | 고유 필드 | 타입 |
|---|---|---|
| InquireMarketCap | `shar` (상장주식수), `tomv` (시가총액), `grav` (비중), `rank` (순위) | `int64,string`, `decimal.Decimal`, `float64,string`, `int64,string` |
| InquireTradeVol | `pask`, `pbid`, `tamt` (거래대금), `a_tvol` (평균거래량), `rank` | `decimal.Decimal`, `decimal.Decimal`, `int64,string`, `int64,string`, `int64,string` |
| InquireTradePbmn | `pask`, `pbid`, `tamt` (거래대금), `a_tamt` (평균거래대금 — `a_tvol` 아님), `rank` | `decimal.Decimal`, `decimal.Decimal`, `int64,string`, `int64,string`, `int64,string` |
| InquireVolumeSurge | `pask`, `pbid`, `n_tvol` (기준거래량), `n_diff` (증가량), `n_rate` (증가율) | `decimal.Decimal`, `decimal.Decimal`, `int64,string`, `decimal.Decimal`, `float64,string` |
| InquireVolumePower | `pask`, `pbid`, `tpow` (당일체결강도), `powx` (체결강도) | `decimal.Decimal`, `decimal.Decimal`, `float64,string`, `float64,string` |
| InquireNewHighlow | `pask`, `pbid`, `n_base` (기준가), `n_diff` (기준가대비), `n_rate` (기준가대비율) | `decimal.Decimal`, `decimal.Decimal`, `decimal.Decimal`, `decimal.Decimal`, `float64,string` |

### VolumePower NDAY 주의사항

KIS docs 의 `InquireVolumePower` query 파라미터 이름은 `NDAY` 이지만 설명값이 `0(1분전), 1(2분전), ..., 9(120분전)` — 분(分) 단위. `InquireVolumeSurge` 의 `MIXN` 과 동일한 척도. **wire name 은 `NDAY` 를 그대로 사용**. Params struct 주석에 "실제로는 분(分) 단위 (1분~120분), 파라미터명 오류로 보임 — wire name `NDAY` 그대로 사용" 를 명기.

---

## Task 1: testdata fixtures (6 합성 JSON)

**Files (Create):**
- `overseas/testdata/market_cap_success.json`
- `overseas/testdata/trade_vol_success.json`
- `overseas/testdata/trade_pbmn_success.json`
- `overseas/testdata/volume_surge_success.json`
- `overseas/testdata/volume_power_success.json`
- `overseas/testdata/new_highlow_success.json`

> KIS docs 응답 필드 기반 합성. 각 파일은 2 샘플 종목 포함 (NAS 계열은 AAPL/MSFT, HKS 계열은 0700/9988). 필드명은 KIS docs 의 JSON key 를 verbatim 사용.

- [ ] **Step 1: market_cap_success.json**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": {
    "zdiv": "2",
    "stat": "정상",
    "crec": "2",
    "trec": "500",
    "nrec": "30"
  },
  "output2": [
    {
      "rsym": "DNASNAAPL",
      "excd": "NAS",
      "symb": "AAPL",
      "name": "애플",
      "last": "189.30",
      "sign": "2",
      "diff": "2.50",
      "rate": "1.34",
      "tvol": "55000000",
      "shar": "15634232000",
      "tomv": "2958652560000",
      "grav": "6.85",
      "rank": "1",
      "ename": "APPLE INC",
      "e_ordyn": "Y"
    },
    {
      "rsym": "DNASNMSFT",
      "excd": "NAS",
      "symb": "MSFT",
      "name": "마이크로소프트",
      "last": "415.20",
      "sign": "2",
      "diff": "4.10",
      "rate": "1.00",
      "tvol": "22000000",
      "shar": "7435000000",
      "tomv": "3086100000000",
      "grav": "7.14",
      "rank": "2",
      "ename": "MICROSOFT CORP",
      "e_ordyn": "Y"
    }
  ]
}
```

- [ ] **Step 2: trade_vol_success.json**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": {
    "zdiv": "2",
    "stat": "정상",
    "crec": "2",
    "trec": "500",
    "nrec": "30"
  },
  "output2": [
    {
      "rsym": "DNASNAAPL",
      "excd": "NAS",
      "symb": "AAPL",
      "name": "애플",
      "last": "189.30",
      "sign": "2",
      "diff": "2.50",
      "rate": "1.34",
      "pask": "189.35",
      "pbid": "189.25",
      "tvol": "55000000",
      "tamt": "10411500000",
      "a_tvol": "48000000",
      "rank": "1",
      "ename": "APPLE INC",
      "e_ordyn": "Y"
    },
    {
      "rsym": "DNASNMSFT",
      "excd": "NAS",
      "symb": "MSFT",
      "name": "마이크로소프트",
      "last": "415.20",
      "sign": "2",
      "diff": "4.10",
      "rate": "1.00",
      "pask": "415.25",
      "pbid": "415.15",
      "tvol": "22000000",
      "tamt": "9134400000",
      "a_tvol": "20000000",
      "rank": "2",
      "ename": "MICROSOFT CORP",
      "e_ordyn": "Y"
    }
  ]
}
```

- [ ] **Step 3: trade_pbmn_success.json**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": {
    "zdiv": "2",
    "stat": "정상",
    "crec": "2",
    "trec": "500",
    "nrec": "30"
  },
  "output2": [
    {
      "rsym": "DNASNMSFT",
      "excd": "NAS",
      "symb": "MSFT",
      "name": "마이크로소프트",
      "last": "415.20",
      "sign": "2",
      "diff": "4.10",
      "rate": "1.00",
      "pask": "415.25",
      "pbid": "415.15",
      "tvol": "22000000",
      "tamt": "9134400000",
      "a_tamt": "8500000000",
      "rank": "1",
      "ename": "MICROSOFT CORP",
      "e_ordyn": "Y"
    },
    {
      "rsym": "DNASNAAPL",
      "excd": "NAS",
      "symb": "AAPL",
      "name": "애플",
      "last": "189.30",
      "sign": "2",
      "diff": "2.50",
      "rate": "1.34",
      "pask": "189.35",
      "pbid": "189.25",
      "tvol": "55000000",
      "tamt": "10411500000",
      "a_tamt": "9800000000",
      "rank": "2",
      "ename": "APPLE INC",
      "e_ordyn": "Y"
    }
  ]
}
```

- [ ] **Step 4: volume_surge_success.json**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": {
    "zdiv": "2",
    "stat": "정상",
    "nrec": "30"
  },
  "output2": [
    {
      "rsym": "DNASNAAPL",
      "excd": "NAS",
      "symb": "AAPL",
      "knam": "애플",
      "last": "189.30",
      "sign": "2",
      "diff": "2.50",
      "rate": "1.34",
      "tvol": "55000000",
      "pask": "189.35",
      "pbid": "189.25",
      "n_tvol": "8000000",
      "n_diff": "47000000",
      "n_rate": "587.50",
      "enam": "APPLE INC",
      "e_ordyn": "Y"
    },
    {
      "rsym": "DNASNMSFT",
      "excd": "NAS",
      "symb": "MSFT",
      "knam": "마이크로소프트",
      "last": "415.20",
      "sign": "2",
      "diff": "4.10",
      "rate": "1.00",
      "tvol": "22000000",
      "pask": "415.25",
      "pbid": "415.15",
      "n_tvol": "4000000",
      "n_diff": "18000000",
      "n_rate": "450.00",
      "enam": "MICROSOFT CORP",
      "e_ordyn": "Y"
    }
  ]
}
```

- [ ] **Step 5: volume_power_success.json**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": {
    "zdiv": "2",
    "stat": "정상",
    "nrec": "30"
  },
  "output2": [
    {
      "rsym": "DNASNAAPL",
      "excd": "NAS",
      "symb": "AAPL",
      "knam": "애플",
      "last": "189.30",
      "sign": "2",
      "diff": "2.50",
      "rate": "1.34",
      "tvol": "55000000",
      "pask": "189.35",
      "pbid": "189.25",
      "tpow": "143.25",
      "powx": "138.90",
      "enam": "APPLE INC",
      "e_ordyn": "Y"
    },
    {
      "rsym": "DNASNMSFT",
      "excd": "NAS",
      "symb": "MSFT",
      "knam": "마이크로소프트",
      "last": "415.20",
      "sign": "2",
      "diff": "4.10",
      "rate": "1.00",
      "tvol": "22000000",
      "pask": "415.25",
      "pbid": "415.15",
      "tpow": "128.70",
      "powx": "125.40",
      "enam": "MICROSOFT CORP",
      "e_ordyn": "Y"
    }
  ]
}
```

- [ ] **Step 6: new_highlow_success.json**

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": {
    "zdiv": "2",
    "stat": "정상",
    "nrec": "30"
  },
  "output2": [
    {
      "rsym": "DNASNAAPL",
      "excd": "NAS",
      "symb": "AAPL",
      "name": "애플",
      "last": "189.30",
      "sign": "2",
      "diff": "2.50",
      "rate": "1.34",
      "tvol": "55000000",
      "pask": "189.35",
      "pbid": "189.25",
      "n_base": "160.00",
      "n_diff": "29.30",
      "n_rate": "18.31",
      "ename": "APPLE INC",
      "e_ordyn": "Y"
    },
    {
      "rsym": "DNASNMSFT",
      "excd": "NAS",
      "symb": "MSFT",
      "name": "마이크로소프트",
      "last": "415.20",
      "sign": "2",
      "diff": "4.10",
      "rate": "1.00",
      "tvol": "22000000",
      "pask": "415.25",
      "pbid": "415.15",
      "n_base": "350.00",
      "n_diff": "65.20",
      "n_rate": "18.63",
      "ename": "MICROSOFT CORP",
      "e_ordyn": "Y"
    }
  ]
}
```

- [ ] **Step 7: 검증**

```bash
for f in overseas/testdata/{market_cap,trade_vol,trade_pbmn,volume_surge,volume_power,new_highlow}_success.json; do
  python3 -c "import json; json.load(open('$f'))" && echo "$f OK" || echo "$f BROKEN"
done
```
Expected: 6 lines `OK`.

- [ ] **Step 8: Commit**

```bash
git add overseas/testdata/{market_cap,trade_vol,trade_pbmn,volume_surge,volume_power,new_highlow}_success.json
git commit -m "$(cat <<'EOF'
[chore] Phase 2.3 testdata — 6 합성 JSON fixtures

시가총액순위 (market_cap, HHDFS76350100) + 거래량순위 (trade_vol, HHDFS76310010)
+ 거래대금순위 (trade_pbmn, HHDFS76320010) + 거래량급증 (volume_surge,
HHDFS76270000) + 매수체결강도상위 (volume_power, HHDFS76280000) + 신고/신저가
(new_highlow, HHDFS76300000). KIS docs 응답 필드 기반 합성.
AAPL/MSFT (NAS) 샘플. output1 2 tier (5-field: #1-#3, 3-field: #4-#6) 검증 가능.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 2: overseas/ranking.go — InquireMarketCap (OverseasRankingFullSummary 정의)

**Files:**
- Modify: `overseas/ranking.go` (APPEND)
- Modify: `overseas/ranking_test.go` (APPEND)

> 5-field output1 Tier 첫 메서드. `OverseasRankingFullSummary` struct (5 fields) 신규 정의. `MarketCapItem` output2 struct 신규 정의. InquireMarketCap 에는 `pask`/`pbid` 가 없음 (output2 15 fields). 4 query params (KEYB, AUTH, EXCD, VOL_RANG).

- [ ] **Step 1: 테스트 추가 — APPEND to `overseas/ranking_test.go`**

```go
func TestClient_InquireMarketCap(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/market-cap`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "market_cap_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireMarketCap(context.Background(), overseas.InquireMarketCapParams{
		ExcdCode: "NAS",
		VolRang:  "0",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "NAS", capturedQuery.Get("EXCD"))
	assert.Equal(t, "0", capturedQuery.Get("VOL_RANG"))

	// output1 검증
	assert.Equal(t, "2", res.Output1.Zdiv)
	assert.Equal(t, int64(2), res.Output1.Crec)
	assert.Equal(t, int64(500), res.Output1.Trec)
	assert.Equal(t, int64(30), res.Output1.Nrec)

	// output2[0] 검증
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "AAPL", res.Output2[0].Symb)
	assert.Equal(t, "애플", res.Output2[0].Name)
	d, _ := decimal.NewFromString("189.30")
	assert.True(t, d.Equal(res.Output2[0].Last))
	assert.InDelta(t, 1.34, res.Output2[0].Rate, 0.001)
	assert.Equal(t, int64(55000000), res.Output2[0].Tvol)
	assert.Equal(t, int64(15634232000), res.Output2[0].Shar)
	tomv, _ := decimal.NewFromString("2958652560000")
	assert.True(t, tomv.Equal(res.Output2[0].Tomv))
	assert.InDelta(t, 6.85, res.Output2[0].Grav, 0.001)
	assert.Equal(t, int64(1), res.Output2[0].Rank)
	assert.Equal(t, "APPLE INC", res.Output2[0].Ename)
}
```

- [ ] **Step 2: FAIL**

`go test ./overseas/... -run InquireMarketCap -v` — 컴파일 실패.

- [ ] **Step 3: 구현 추가 — APPEND to `overseas/ranking.go`**

```go
// OverseasRankingFullSummary 는 output1 5-field tier (시가총액/거래량/거래대금순위).
//
// 해당 메서드: InquireMarketCap, InquireTradeVol, InquireTradePbmn.
type OverseasRankingFullSummary struct {
	Zdiv string `json:"zdiv"`        // 소수점자리수
	Stat string `json:"stat"`        // 거래상태정보
	Crec int64  `json:"crec,string"` // 현재조회종목수
	Trec int64  `json:"trec,string"` // 전체조회종목수
	Nrec int64  `json:"nrec,string"` // RecordCount
}

// MarketCap 은 해외주식_시가총액순위 (HHDFS76350100) 응답.
//
// 한투 docs: docs/api/해외주식/해외주식_시가총액순위.md
// path: /uapi/overseas-stock/v1/ranking/market-cap
type MarketCap struct {
	Output1 OverseasRankingFullSummary `json:"output1"`
	Output2 []MarketCapItem            `json:"output2"`
}

// MarketCapItem 은 시가총액순위 output2 의 한 행 (15 fields).
type MarketCapItem struct {
	Rsym   string          `json:"rsym"`         // 실시간조회심볼
	Excd   string          `json:"excd"`         // 거래소코드
	Symb   string          `json:"symb"`         // 종목코드
	Name   string          `json:"name"`         // 종목명 (한글)
	Last   decimal.Decimal `json:"last"`         // 현재가
	Sign   string          `json:"sign"`         // 기호
	Diff   decimal.Decimal `json:"diff"`         // 대비
	Rate   float64         `json:"rate,string"`  // 등락율
	Tvol   int64           `json:"tvol,string"`  // 거래량
	Shar   int64           `json:"shar,string"`  // 상장주식수
	Tomv   decimal.Decimal `json:"tomv"`         // 시가총액
	Grav   float64         `json:"grav,string"`  // 비중
	Rank   int64           `json:"rank,string"`  // 순위
	Ename  string          `json:"ename"`        // 영문종목명
	EOrdyn string          `json:"e_ordyn"`      // 매매가능
}

// InquireMarketCapParams 는 해외주식_시가총액순위 조회 파라미터.
type InquireMarketCapParams struct {
	KeyB     string // KEYB — NEXT KEY BUFF. 빈 값 default
	Auth     string // AUTH — 사용자권한정보. 빈 값 default
	ExcdCode string // EXCD — 거래소코드 (NYS/NAS/AMS/HKS/SHS/SZS/HSX/HNX/TSE). 필수
	VolRang  string // VOL_RANG — 거래량조건. 빈 값=>"0" (전체)
}

// InquireMarketCap 은 해외주식_시가총액순위 호출.
//
// 한투 docs: docs/api/해외주식/해외주식_시가총액순위.md
// path: /uapi/overseas-stock/v1/ranking/market-cap (HHDFS76350100)
func (c *Client) InquireMarketCap(ctx context.Context, params InquireMarketCapParams) (*MarketCap, error) {
	excd := params.ExcdCode
	if excd == "" {
		return nil, fmt.Errorf("kis: ExcdCode required for InquireMarketCap")
	}
	vol := params.VolRang
	if vol == "" {
		vol = "0"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-stock/v1/ranking/market-cap",
		TrID:   "HHDFS76350100",
		Query: map[string]string{
			"KEYB":     params.KeyB,
			"AUTH":     params.Auth,
			"EXCD":     excd,
			"VOL_RANG": vol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res MarketCap
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse MarketCap: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: PASS**

`go test ./overseas/... -run InquireMarketCap -v`

- [ ] **Step 5: Commit**

```bash
git add overseas/ranking.go overseas/ranking_test.go
git commit -m "$(cat <<'EOF'
[feat] overseas — InquireMarketCap (시가총액순위, HHDFS76350100)

OverseasRankingFullSummary (output1 5-field tier 공유 struct) + MarketCap +
MarketCapItem (15 필드: rsym/excd/symb/name/last/sign/diff/rate/tvol/shar/
tomv/grav/rank/ename/e_ordyn) + InquireMarketCapParams (4 query params).
pask/pbid 없음 (시가총액순위 output2 특성).

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 3: InquireTradeVol

**Files:**
- Modify: `overseas/ranking.go` (APPEND)
- Modify: `overseas/ranking_test.go` (APPEND)

> `OverseasRankingFullSummary` 재사용. `TradeVolItem` 신규 정의 (16 fields, pask/pbid/tamt/a_tvol/rank 포함). 6 query params (KEYB, AUTH, EXCD, NDAY, PRC1, PRC2, VOL_RANG — 실제 7개이나 KIS 문서에 6개로 기술, 전체 포함).

- [ ] **Step 1: 테스트 추가 — APPEND to `overseas/ranking_test.go`**

```go
func TestClient_InquireTradeVol(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/trade-vol`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "trade_vol_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireTradeVol(context.Background(), overseas.InquireTradeVolParams{
		ExcdCode: "NAS",
		NDay:     "0",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "NAS", capturedQuery.Get("EXCD"))
	assert.Equal(t, "0", capturedQuery.Get("NDAY"))
	assert.Equal(t, "0", capturedQuery.Get("VOL_RANG"))

	// output1 검증
	assert.Equal(t, int64(2), res.Output1.Crec)
	assert.Equal(t, int64(30), res.Output1.Nrec)

	// output2[0] 검증
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "AAPL", res.Output2[0].Symb)
	d, _ := decimal.NewFromString("189.30")
	assert.True(t, d.Equal(res.Output2[0].Last))
	pask, _ := decimal.NewFromString("189.35")
	assert.True(t, pask.Equal(res.Output2[0].Pask))
	pbid, _ := decimal.NewFromString("189.25")
	assert.True(t, pbid.Equal(res.Output2[0].Pbid))
	assert.Equal(t, int64(55000000), res.Output2[0].Tvol)
	assert.Equal(t, int64(10411500000), res.Output2[0].Tamt)
	assert.Equal(t, int64(48000000), res.Output2[0].ATvol)
	assert.Equal(t, int64(1), res.Output2[0].Rank)
}
```

- [ ] **Step 2: FAIL**

`go test ./overseas/... -run InquireTradeVol -v` — 컴파일 실패.

- [ ] **Step 3: 구현 추가 — APPEND to `overseas/ranking.go`**

```go
// TradeVol 은 해외주식_거래량순위 (HHDFS76310010) 응답.
//
// 한투 docs: docs/api/해외주식/해외주식_거래량순위.md
// path: /uapi/overseas-stock/v1/ranking/trade-vol
type TradeVol struct {
	Output1 OverseasRankingFullSummary `json:"output1"`
	Output2 []TradeVolItem             `json:"output2"`
}

// TradeVolItem 은 거래량순위 output2 의 한 행 (16 fields).
type TradeVolItem struct {
	Rsym   string          `json:"rsym"`          // 실시간조회심볼
	Excd   string          `json:"excd"`          // 거래소코드
	Symb   string          `json:"symb"`          // 종목코드
	Name   string          `json:"name"`          // 종목명 (한글)
	Last   decimal.Decimal `json:"last"`          // 현재가
	Sign   string          `json:"sign"`          // 기호
	Diff   decimal.Decimal `json:"diff"`          // 대비
	Rate   float64         `json:"rate,string"`   // 등락율
	Pask   decimal.Decimal `json:"pask"`          // 매도호가
	Pbid   decimal.Decimal `json:"pbid"`          // 매수호가
	Tvol   int64           `json:"tvol,string"`   // 거래량
	Tamt   int64           `json:"tamt,string"`   // 거래대금
	ATvol  int64           `json:"a_tvol,string"` // 평균거래량
	Rank   int64           `json:"rank,string"`   // 순위
	Ename  string          `json:"ename"`         // 영문종목명
	EOrdyn string          `json:"e_ordyn"`       // 매매가능
}

// InquireTradeVolParams 는 해외주식_거래량순위 조회 파라미터.
type InquireTradeVolParams struct {
	KeyB     string // KEYB — NEXT KEY BUFF. 빈 값 default
	Auth     string // AUTH — 사용자권한정보. 빈 값 default
	ExcdCode string // EXCD — 거래소코드. 필수
	NDay     string // NDAY — N일전: 0(당일),1(2일),2(3일),3(5일),4(10일),5(20일),6(30일),7(60일),8(120일),9(1년). 빈 값=>"0"
	Prc1     string // PRC1 — 현재가 필터범위 1 (가격 ~). 빈 값 OK
	Prc2     string // PRC2 — 현재가 필터범위 2 (~ 가격). 빈 값 OK
	VolRang  string // VOL_RANG — 거래량조건. 빈 값=>"0" (전체)
}

// InquireTradeVol 은 해외주식_거래량순위 호출.
//
// 한투 docs: docs/api/해외주식/해외주식_거래량순위.md
// path: /uapi/overseas-stock/v1/ranking/trade-vol (HHDFS76310010)
func (c *Client) InquireTradeVol(ctx context.Context, params InquireTradeVolParams) (*TradeVol, error) {
	excd := params.ExcdCode
	if excd == "" {
		return nil, fmt.Errorf("kis: ExcdCode required for InquireTradeVol")
	}
	nday := params.NDay
	if nday == "" {
		nday = "0"
	}
	vol := params.VolRang
	if vol == "" {
		vol = "0"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-stock/v1/ranking/trade-vol",
		TrID:   "HHDFS76310010",
		Query: map[string]string{
			"KEYB":     params.KeyB,
			"AUTH":     params.Auth,
			"EXCD":     excd,
			"NDAY":     nday,
			"PRC1":     params.Prc1,
			"PRC2":     params.Prc2,
			"VOL_RANG": vol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res TradeVol
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse TradeVol: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: PASS**

`go test ./overseas/... -run InquireTradeVol -v`

- [ ] **Step 5: Commit**

```bash
git add overseas/ranking.go overseas/ranking_test.go
git commit -m "$(cat <<'EOF'
[feat] overseas — InquireTradeVol (거래량순위, HHDFS76310010)

TradeVol + TradeVolItem (16 필드: 공통 8 + pask/pbid/tamt/a_tvol/rank/ename/
e_ordyn) + InquireTradeVolParams (7 query params: KEYB/AUTH/EXCD/NDAY/PRC1/
PRC2/VOL_RANG). OverseasRankingFullSummary 재사용.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 4: InquireTradePbmn

**Files:**
- Modify: `overseas/ranking.go` (APPEND)
- Modify: `overseas/ranking_test.go` (APPEND)

> `OverseasRankingFullSummary` 재사용. `TradePbmnItem` 신규 정의 (16 fields). 핵심 차이: `a_tvol` 대신 `a_tamt` (평균거래대금) 필드. 동일한 7 query params (KEYB, AUTH, EXCD, NDAY, VOL_RANG, PRC1, PRC2).

- [ ] **Step 1: 테스트 추가 — APPEND to `overseas/ranking_test.go`**

```go
func TestClient_InquireTradePbmn(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/trade-pbmn`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "trade_pbmn_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireTradePbmn(context.Background(), overseas.InquireTradePbmnParams{
		ExcdCode: "NAS",
		NDay:     "0",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "NAS", capturedQuery.Get("EXCD"))
	assert.Equal(t, "0", capturedQuery.Get("NDAY"))
	assert.Equal(t, "0", capturedQuery.Get("VOL_RANG"))

	// output1 검증
	assert.Equal(t, int64(30), res.Output1.Nrec)

	// output2[0] 검증 — MSFT 가 순위 1 (거래대금 기준)
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "MSFT", res.Output2[0].Symb)
	d, _ := decimal.NewFromString("415.20")
	assert.True(t, d.Equal(res.Output2[0].Last))
	assert.Equal(t, int64(9134400000), res.Output2[0].Tamt)
	assert.Equal(t, int64(8500000000), res.Output2[0].ATamt) // a_tamt, not a_tvol
	assert.Equal(t, int64(1), res.Output2[0].Rank)
}
```

- [ ] **Step 2: FAIL**

`go test ./overseas/... -run InquireTradePbmn -v` — 컴파일 실패.

- [ ] **Step 3: 구현 추가 — APPEND to `overseas/ranking.go`**

```go
// TradePbmn 은 해외주식_거래대금순위 (HHDFS76320010) 응답.
//
// 한투 docs: docs/api/해외주식/해외주식_거래대금순위.md
// path: /uapi/overseas-stock/v1/ranking/trade-pbmn
type TradePbmn struct {
	Output1 OverseasRankingFullSummary `json:"output1"`
	Output2 []TradePbmnItem            `json:"output2"`
}

// TradePbmnItem 은 거래대금순위 output2 의 한 행 (16 fields).
//
// 주의: 평균 필드명이 a_tamt (평균거래대금) — TradeVolItem 의 a_tvol 과 다름.
type TradePbmnItem struct {
	Rsym   string          `json:"rsym"`          // 실시간조회심볼
	Excd   string          `json:"excd"`          // 거래소코드
	Symb   string          `json:"symb"`          // 종목코드
	Name   string          `json:"name"`          // 종목명 (한글)
	Last   decimal.Decimal `json:"last"`          // 현재가
	Sign   string          `json:"sign"`          // 기호
	Diff   decimal.Decimal `json:"diff"`          // 대비
	Rate   float64         `json:"rate,string"`   // 등락율
	Pask   decimal.Decimal `json:"pask"`          // 매도호가
	Pbid   decimal.Decimal `json:"pbid"`          // 매수호가
	Tvol   int64           `json:"tvol,string"`   // 거래량
	Tamt   int64           `json:"tamt,string"`   // 거래대금
	ATamt  int64           `json:"a_tamt,string"` // 평균거래대금 (a_tvol 아님)
	Rank   int64           `json:"rank,string"`   // 순위
	Ename  string          `json:"ename"`         // 영문종목명
	EOrdyn string          `json:"e_ordyn"`       // 매매가능
}

// InquireTradePbmnParams 는 해외주식_거래대금순위 조회 파라미터.
type InquireTradePbmnParams struct {
	KeyB     string // KEYB — NEXT KEY BUFF. 빈 값 default
	Auth     string // AUTH — 사용자권한정보. 빈 값 default
	ExcdCode string // EXCD — 거래소코드. 필수
	NDay     string // NDAY — N일전: 0(당일),1(2일),2(3일),3(5일),4(10일),5(20일),6(30일),7(60일),8(120일),9(1년). 빈 값=>"0"
	VolRang  string // VOL_RANG — 거래량조건. 빈 값=>"0" (전체)
	Prc1     string // PRC1 — 현재가 필터범위 1 (가격 ~). 빈 값 OK
	Prc2     string // PRC2 — 현재가 필터범위 2 (~ 가격). 빈 값 OK
}

// InquireTradePbmn 은 해외주식_거래대금순위 호출.
//
// 한투 docs: docs/api/해외주식/해외주식_거래대금순위.md
// path: /uapi/overseas-stock/v1/ranking/trade-pbmn (HHDFS76320010)
func (c *Client) InquireTradePbmn(ctx context.Context, params InquireTradePbmnParams) (*TradePbmn, error) {
	excd := params.ExcdCode
	if excd == "" {
		return nil, fmt.Errorf("kis: ExcdCode required for InquireTradePbmn")
	}
	nday := params.NDay
	if nday == "" {
		nday = "0"
	}
	vol := params.VolRang
	if vol == "" {
		vol = "0"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-stock/v1/ranking/trade-pbmn",
		TrID:   "HHDFS76320010",
		Query: map[string]string{
			"KEYB":     params.KeyB,
			"AUTH":     params.Auth,
			"EXCD":     excd,
			"NDAY":     nday,
			"VOL_RANG": vol,
			"PRC1":     params.Prc1,
			"PRC2":     params.Prc2,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res TradePbmn
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse TradePbmn: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: PASS**

`go test ./overseas/... -run InquireTradePbmn -v`

- [ ] **Step 5: Commit**

```bash
git add overseas/ranking.go overseas/ranking_test.go
git commit -m "$(cat <<'EOF'
[feat] overseas — InquireTradePbmn (거래대금순위, HHDFS76320010)

TradePbmn + TradePbmnItem (16 필드: a_tamt=평균거래대금, TradeVolItem 의
a_tvol 과 다름 주의) + InquireTradePbmnParams (7 query params).
OverseasRankingFullSummary 재사용.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 5: InquireVolumeSurge (OverseasRankingMinSummary 정의)

**Files:**
- Modify: `overseas/ranking.go` (APPEND)
- Modify: `overseas/ranking_test.go` (APPEND)

> 3-field output1 Tier 첫 메서드. `OverseasRankingMinSummary` struct (3 fields) 신규 정의. `VolumeSurgeItem` output2 struct 신규 정의 (16 fields, `knam`/`enam` 사용). 5 query params (KEYB, AUTH, EXCD, MIXN, VOL_RANG).

- [ ] **Step 1: 테스트 추가 — APPEND to `overseas/ranking_test.go`**

```go
func TestClient_InquireVolumeSurge(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/volume-surge`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "volume_surge_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireVolumeSurge(context.Background(), overseas.InquireVolumeSurgeParams{
		ExcdCode: "NAS",
		MixN:     "0",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "NAS", capturedQuery.Get("EXCD"))
	assert.Equal(t, "0", capturedQuery.Get("MIXN"))
	assert.Equal(t, "0", capturedQuery.Get("VOL_RANG"))

	// output1 검증 — 3-field MinSummary (crec/trec 없음)
	assert.Equal(t, "2", res.Output1.Zdiv)
	assert.Equal(t, int64(30), res.Output1.Nrec)

	// output2[0] 검증 — knam/enam (name/ename 아님)
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "AAPL", res.Output2[0].Symb)
	assert.Equal(t, "애플", res.Output2[0].Knam)
	assert.Equal(t, "APPLE INC", res.Output2[0].Enam)
	d, _ := decimal.NewFromString("189.30")
	assert.True(t, d.Equal(res.Output2[0].Last))
	assert.Equal(t, int64(55000000), res.Output2[0].Tvol)
	assert.Equal(t, int64(8000000), res.Output2[0].NTvol)
	ndiff, _ := decimal.NewFromString("47000000")
	assert.True(t, ndiff.Equal(res.Output2[0].NDiff))
	assert.InDelta(t, 587.50, res.Output2[0].NRate, 0.01)
}
```

- [ ] **Step 2: FAIL**

`go test ./overseas/... -run InquireVolumeSurge -v` — 컴파일 실패.

- [ ] **Step 3: 구현 추가 — APPEND to `overseas/ranking.go`**

```go
// OverseasRankingMinSummary 는 output1 3-field tier (거래량급증/매수체결강도/신고신저).
//
// 해당 메서드: InquireVolumeSurge, InquireVolumePower, InquireNewHighlow.
// crec/trec 없음 — OverseasRankingFullSummary 와 혼용 금지.
type OverseasRankingMinSummary struct {
	Zdiv string `json:"zdiv"`        // 소수점자리수
	Stat string `json:"stat"`        // 거래상태정보
	Nrec int64  `json:"nrec,string"` // RecordCount
}

// VolumeSurge 는 해외주식_거래량급증 (HHDFS76270000) 응답.
//
// 한투 docs: docs/api/해외주식/해외주식_거래량급증.md
// path: /uapi/overseas-stock/v1/ranking/volume-surge
type VolumeSurge struct {
	Output1 OverseasRankingMinSummary `json:"output1"`
	Output2 []VolumeSurgeItem         `json:"output2"`
}

// VolumeSurgeItem 은 거래량급증 output2 의 한 행 (16 fields).
//
// 주의: 종목명 필드가 knam/enam (name/ename 아님).
type VolumeSurgeItem struct {
	Rsym   string          `json:"rsym"`          // 실시간조회심볼
	Excd   string          `json:"excd"`          // 거래소코드
	Symb   string          `json:"symb"`          // 종목코드
	Knam   string          `json:"knam"`          // 종목명 (한글) — name 아님
	Last   decimal.Decimal `json:"last"`          // 현재가
	Sign   string          `json:"sign"`          // 기호
	Diff   decimal.Decimal `json:"diff"`          // 대비
	Rate   float64         `json:"rate,string"`   // 등락율
	Tvol   int64           `json:"tvol,string"`   // 거래량
	Pask   decimal.Decimal `json:"pask"`          // 매도호가
	Pbid   decimal.Decimal `json:"pbid"`          // 매수호가
	NTvol  int64           `json:"n_tvol,string"` // 기준거래량
	NDiff  decimal.Decimal `json:"n_diff"`        // 증가량
	NRate  float64         `json:"n_rate,string"` // 증가율
	Enam   string          `json:"enam"`          // 영문종목명 — ename 아님
	EOrdyn string          `json:"e_ordyn"`       // 매매가능
}

// InquireVolumeSurgeParams 는 해외주식_거래량급증 조회 파라미터.
type InquireVolumeSurgeParams struct {
	KeyB     string // KEYB — NEXT KEY BUFF. 빈 값 default
	Auth     string // AUTH — 사용자권한정보. 빈 값 default
	ExcdCode string // EXCD — 거래소코드. 필수
	MixN     string // MIXN — N분전: 0(1분전),1(2분전),2(3분전),3(5분전),4(10분전),5(15분전),6(20분전),7(30분전),8(60분전),9(120분전). 빈 값=>"0"
	VolRang  string // VOL_RANG — 거래량조건. 빈 값=>"0" (전체)
}

// InquireVolumeSurge 는 해외주식_거래량급증 호출.
//
// 한투 docs: docs/api/해외주식/해외주식_거래량급증.md
// path: /uapi/overseas-stock/v1/ranking/volume-surge (HHDFS76270000)
func (c *Client) InquireVolumeSurge(ctx context.Context, params InquireVolumeSurgeParams) (*VolumeSurge, error) {
	excd := params.ExcdCode
	if excd == "" {
		return nil, fmt.Errorf("kis: ExcdCode required for InquireVolumeSurge")
	}
	mixn := params.MixN
	if mixn == "" {
		mixn = "0"
	}
	vol := params.VolRang
	if vol == "" {
		vol = "0"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-stock/v1/ranking/volume-surge",
		TrID:   "HHDFS76270000",
		Query: map[string]string{
			"KEYB":     params.KeyB,
			"AUTH":     params.Auth,
			"EXCD":     excd,
			"MIXN":     mixn,
			"VOL_RANG": vol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res VolumeSurge
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse VolumeSurge: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: PASS**

`go test ./overseas/... -run InquireVolumeSurge -v`

- [ ] **Step 5: Commit**

```bash
git add overseas/ranking.go overseas/ranking_test.go
git commit -m "$(cat <<'EOF'
[feat] overseas — InquireVolumeSurge (거래량급증, HHDFS76270000)

OverseasRankingMinSummary (output1 3-field tier 공유 struct: zdiv/stat/nrec,
crec/trec 없음) + VolumeSurge + VolumeSurgeItem (16 필드: knam/enam 사용,
name/ename 아님) + InquireVolumeSurgeParams (5 query params: MIXN=분 단위).

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 6: InquireVolumePower

**Files:**
- Modify: `overseas/ranking.go` (APPEND)
- Modify: `overseas/ranking_test.go` (APPEND)

> `OverseasRankingMinSummary` 재사용. `VolumePowerItem` 신규 정의 (15 fields, `knam`/`enam` 사용, `tpow`/`powx` 고유). 5 query params (KEYB, AUTH, EXCD, NDAY, VOL_RANG). **주의:** NDAY 파라미터명 사용 — 실제로는 분(分) 단위.

- [ ] **Step 1: 테스트 추가 — APPEND to `overseas/ranking_test.go`**

```go
func TestClient_InquireVolumePower(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/volume-power`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "volume_power_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireVolumePower(context.Background(), overseas.InquireVolumePowerParams{
		ExcdCode: "NAS",
		NDay:     "0",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "NAS", capturedQuery.Get("EXCD"))
	assert.Equal(t, "0", capturedQuery.Get("NDAY")) // wire name NDAY (분 단위)
	assert.Equal(t, "0", capturedQuery.Get("VOL_RANG"))

	// output1 검증 — 3-field MinSummary
	assert.Equal(t, int64(30), res.Output1.Nrec)

	// output2[0] 검증 — knam/enam + tpow/powx
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "AAPL", res.Output2[0].Symb)
	assert.Equal(t, "애플", res.Output2[0].Knam)
	assert.Equal(t, "APPLE INC", res.Output2[0].Enam)
	d, _ := decimal.NewFromString("189.30")
	assert.True(t, d.Equal(res.Output2[0].Last))
	assert.InDelta(t, 143.25, res.Output2[0].Tpow, 0.01)
	assert.InDelta(t, 138.90, res.Output2[0].Powx, 0.01)
}
```

- [ ] **Step 2: FAIL**

`go test ./overseas/... -run InquireVolumePower -v` — 컴파일 실패.

- [ ] **Step 3: 구현 추가 — APPEND to `overseas/ranking.go`**

```go
// VolumePower 는 해외주식_매수체결강도상위 (HHDFS76280000) 응답.
//
// 한투 docs: docs/api/해외주식/해외주식_매수체결강도상위.md
// path: /uapi/overseas-stock/v1/ranking/volume-power
type VolumePower struct {
	Output1 OverseasRankingMinSummary `json:"output1"`
	Output2 []VolumePowerItem         `json:"output2"`
}

// VolumePowerItem 은 매수체결강도상위 output2 의 한 행 (15 fields).
//
// 주의: 종목명 필드가 knam/enam (name/ename 아님). rank 필드 없음.
type VolumePowerItem struct {
	Rsym   string          `json:"rsym"`         // 실시간조회심볼
	Excd   string          `json:"excd"`         // 거래소코드
	Symb   string          `json:"symb"`         // 종목코드
	Knam   string          `json:"knam"`         // 종목명 (한글) — name 아님
	Last   decimal.Decimal `json:"last"`         // 현재가
	Sign   string          `json:"sign"`         // 기호
	Diff   decimal.Decimal `json:"diff"`         // 대비
	Rate   float64         `json:"rate,string"`  // 등락율
	Tvol   int64           `json:"tvol,string"`  // 거래량
	Pask   decimal.Decimal `json:"pask"`         // 매도호가
	Pbid   decimal.Decimal `json:"pbid"`         // 매수호가
	Tpow   float64         `json:"tpow,string"`  // 당일체결강도
	Powx   float64         `json:"powx,string"`  // 체결강도
	Enam   string          `json:"enam"`         // 영문종목명 — ename 아님
	EOrdyn string          `json:"e_ordyn"`      // 매매가능
}

// InquireVolumePowerParams 는 해외주식_매수체결강도상위 조회 파라미터.
//
// 주의: NDAY 파라미터의 설명값이 분(分) 단위 (0=1분전, 1=2분전 … 9=120분전)로
// InquireVolumeSurge 의 MIXN 과 동일한 척도. KIS docs 파라미터명 오류로 보이나
// wire name 은 NDAY 를 그대로 사용.
type InquireVolumePowerParams struct {
	KeyB     string // KEYB — NEXT KEY BUFF. 빈 값 default
	Auth     string // AUTH — 사용자권한정보. 빈 값 default
	ExcdCode string // EXCD — 거래소코드. 필수
	NDay     string // NDAY — N분전 (wire name): 0(1분전),1(2분전),2(3분전),3(5분전),4(10분전),5(15분전),6(20분전),7(30분전),8(60분전),9(120분전). 빈 값=>"0"
	VolRang  string // VOL_RANG — 거래량조건. 빈 값=>"0" (전체)
}

// InquireVolumePower 는 해외주식_매수체결강도상위 호출.
//
// 한투 docs: docs/api/해외주식/해외주식_매수체결강도상위.md
// path: /uapi/overseas-stock/v1/ranking/volume-power (HHDFS76280000)
func (c *Client) InquireVolumePower(ctx context.Context, params InquireVolumePowerParams) (*VolumePower, error) {
	excd := params.ExcdCode
	if excd == "" {
		return nil, fmt.Errorf("kis: ExcdCode required for InquireVolumePower")
	}
	nday := params.NDay
	if nday == "" {
		nday = "0"
	}
	vol := params.VolRang
	if vol == "" {
		vol = "0"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-stock/v1/ranking/volume-power",
		TrID:   "HHDFS76280000",
		Query: map[string]string{
			"KEYB":     params.KeyB,
			"AUTH":     params.Auth,
			"EXCD":     excd,
			"NDAY":     nday,
			"VOL_RANG": vol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res VolumePower
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse VolumePower: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: PASS**

`go test ./overseas/... -run InquireVolumePower -v`

- [ ] **Step 5: Commit**

```bash
git add overseas/ranking.go overseas/ranking_test.go
git commit -m "$(cat <<'EOF'
[feat] overseas — InquireVolumePower (매수체결강도상위, HHDFS76280000)

VolumePower + VolumePowerItem (15 필드: knam/enam 사용, tpow/powx 체결강도)
+ InquireVolumePowerParams (5 query params: NDAY wire name 사용 — 실제 분 단위,
KIS docs 파라미터명 오류 추정). OverseasRankingMinSummary 재사용.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 7: InquireNewHighlow

**Files:**
- Modify: `overseas/ranking.go` (APPEND)
- Modify: `overseas/ranking_test.go` (APPEND)

> `OverseasRankingMinSummary` 재사용. `NewHighlowItem` 신규 정의 (16 fields, `name`/`ename` 사용 — VolumeSurge/VolumePower 와 다름). 7 query params (KEYB, AUTH, EXCD, GUBN, GUBN2, NDAY, VOL_RANG).

- [ ] **Step 1: 테스트 추가 — APPEND to `overseas/ranking_test.go`**

```go
func TestClient_InquireNewHighlow(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ranking/new-highlow`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "new_highlow_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireNewHighlow(context.Background(), overseas.InquireNewHighlowParams{
		ExcdCode: "NAS",
		Gubn:     "1", // 신고(1)
		Gubn2:    "1", // 돌파유지(1)
		NDay:     "6", // 52주
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "NAS", capturedQuery.Get("EXCD"))
	assert.Equal(t, "1", capturedQuery.Get("GUBN"))
	assert.Equal(t, "1", capturedQuery.Get("GUBN2"))
	assert.Equal(t, "6", capturedQuery.Get("NDAY"))
	assert.Equal(t, "0", capturedQuery.Get("VOL_RANG"))

	// output1 검증 — 3-field MinSummary
	assert.Equal(t, int64(30), res.Output1.Nrec)

	// output2[0] 검증 — name/ename (knam/enam 아님) + n_base/n_diff/n_rate
	require.Len(t, res.Output2, 2)
	assert.Equal(t, "AAPL", res.Output2[0].Symb)
	assert.Equal(t, "애플", res.Output2[0].Name)
	assert.Equal(t, "APPLE INC", res.Output2[0].Ename)
	d, _ := decimal.NewFromString("189.30")
	assert.True(t, d.Equal(res.Output2[0].Last))
	nbase, _ := decimal.NewFromString("160.00")
	assert.True(t, nbase.Equal(res.Output2[0].NBase))
	ndiff, _ := decimal.NewFromString("29.30")
	assert.True(t, ndiff.Equal(res.Output2[0].NDiff))
	assert.InDelta(t, 18.31, res.Output2[0].NRate, 0.01)
}
```

- [ ] **Step 2: FAIL**

`go test ./overseas/... -run InquireNewHighlow -v` — 컴파일 실패.

- [ ] **Step 3: 구현 추가 — APPEND to `overseas/ranking.go`**

```go
// NewHighlow 는 해외주식_신고/신저가 (HHDFS76300000) 응답.
//
// 한투 docs: docs/api/해외주식/해외주식_신고_신저가.md
// path: /uapi/overseas-stock/v1/ranking/new-highlow
type NewHighlow struct {
	Output1 OverseasRankingMinSummary `json:"output1"`
	Output2 []NewHighlowItem          `json:"output2"`
}

// NewHighlowItem 은 신고/신저가 output2 의 한 행 (16 fields).
//
// 주의: 종목명 필드가 name/ename — VolumeSurge/VolumePower 의 knam/enam 와 다름.
type NewHighlowItem struct {
	Rsym   string          `json:"rsym"`          // 실시간조회심볼
	Excd   string          `json:"excd"`          // 거래소코드
	Symb   string          `json:"symb"`          // 종목코드
	Name   string          `json:"name"`          // 종목명 (한글) — knam 아님
	Last   decimal.Decimal `json:"last"`          // 현재가
	Sign   string          `json:"sign"`          // 기호
	Diff   decimal.Decimal `json:"diff"`          // 대비
	Rate   float64         `json:"rate,string"`   // 등락율
	Tvol   int64           `json:"tvol,string"`   // 거래량
	Pask   decimal.Decimal `json:"pask"`          // 매도호가
	Pbid   decimal.Decimal `json:"pbid"`          // 매수호가
	NBase  decimal.Decimal `json:"n_base"`        // 기준가
	NDiff  decimal.Decimal `json:"n_diff"`        // 기준가대비
	NRate  float64         `json:"n_rate,string"` // 기준가대비율
	Ename  string          `json:"ename"`         // 영문종목명 — enam 아님
	EOrdyn string          `json:"e_ordyn"`       // 매매가능
}

// InquireNewHighlowParams 는 해외주식_신고/신저가 조회 파라미터.
type InquireNewHighlowParams struct {
	KeyB     string // KEYB — NEXT KEY BUFF. 빈 값 default
	Auth     string // AUTH — 사용자권한정보. 빈 값 default
	ExcdCode string // EXCD — 거래소코드. 필수
	Gubn     string // GUBN — 신고(1)/신저(0) 구분. 빈 값=>"1"
	Gubn2    string // GUBN2 — 일시돌파(0)/돌파유지(1) 구분. 빈 값=>"1"
	NDay     string // NDAY — N일자값: 0(5일),1(10일),2(20일),3(30일),4(60일),5(120일),6(52주),7(1년). 빈 값=>"6" (52주)
	VolRang  string // VOL_RANG — 거래량조건. 빈 값=>"0" (전체)
}

// InquireNewHighlow 는 해외주식_신고/신저가 호출.
//
// 한투 docs: docs/api/해외주식/해외주식_신고_신저가.md
// path: /uapi/overseas-stock/v1/ranking/new-highlow (HHDFS76300000)
func (c *Client) InquireNewHighlow(ctx context.Context, params InquireNewHighlowParams) (*NewHighlow, error) {
	excd := params.ExcdCode
	if excd == "" {
		return nil, fmt.Errorf("kis: ExcdCode required for InquireNewHighlow")
	}
	gubn := params.Gubn
	if gubn == "" {
		gubn = "1"
	}
	gubn2 := params.Gubn2
	if gubn2 == "" {
		gubn2 = "1"
	}
	nday := params.NDay
	if nday == "" {
		nday = "6"
	}
	vol := params.VolRang
	if vol == "" {
		vol = "0"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-stock/v1/ranking/new-highlow",
		TrID:   "HHDFS76300000",
		Query: map[string]string{
			"KEYB":     params.KeyB,
			"AUTH":     params.Auth,
			"EXCD":     excd,
			"GUBN":     gubn,
			"GUBN2":    gubn2,
			"NDAY":     nday,
			"VOL_RANG": vol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res NewHighlow
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse NewHighlow: %w", err)
	}
	return &res, nil
}
```

- [ ] **Step 4: PASS**

`go test ./overseas/... -run InquireNewHighlow -v`

- [ ] **Step 5: 전체 회귀 테스트**

`go test ./... -count=1` — all PASS.

- [ ] **Step 6: Commit**

```bash
git add overseas/ranking.go overseas/ranking_test.go
git commit -m "$(cat <<'EOF'
[feat] overseas — InquireNewHighlow (신고/신저가, HHDFS76300000)

NewHighlow + NewHighlowItem (16 필드: name/ename 사용 — knam/enam 아님;
n_base/n_diff/n_rate 기준가 관련) + InquireNewHighlowParams (7 query params:
GUBN 신고/신저, GUBN2 일시/유지, NDAY 기간). OverseasRankingMinSummary 재사용.
Phase 2.3 해외주식 ranking 6 메서드 구현 완료.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 8: examples/overseas_ranking/main.go

**Files:**
- Create: `examples/overseas_ranking/main.go`

- [ ] **Step 1: example 작성**

```go
// overseas_ranking example: InquireMarketCap + InquireTradeVol + InquireTradePbmn
// + InquireVolumeSurge + InquireVolumePower + InquireNewHighlow.
//
// Run: KIS credentials env vars 후 go run ./examples/overseas_ranking
package main

import (
	"context"
	"fmt"
	"log"

	kis "github.com/kenshin579/korea-investment-stock"
	"github.com/kenshin579/korea-investment-stock/overseas"
)

func main() {
	client, err := kis.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	excd := "NAS"

	// 1. 시가총액순위
	mc, err := client.Overseas.InquireMarketCap(ctx, overseas.InquireMarketCapParams{
		ExcdCode: excd,
		VolRang:  "0",
	})
	if err != nil {
		log.Fatalf("InquireMarketCap: %v", err)
	}
	fmt.Printf("[%s 시가총액 상위 %d 건]\n", excd, len(mc.Output2))
	for i, item := range mc.Output2 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  #%d %s (%s): %s USD (시총 %s, 비중 %v%%)\n",
			item.Rank, item.Name, item.Symb, item.Last, item.Tomv, item.Grav)
	}

	// 2. 거래량순위
	tv, err := client.Overseas.InquireTradeVol(ctx, overseas.InquireTradeVolParams{
		ExcdCode: excd,
		NDay:     "0",
	})
	if err != nil {
		log.Fatalf("InquireTradeVol: %v", err)
	}
	fmt.Printf("\n[%s 거래량 상위 %d 건]\n", excd, len(tv.Output2))
	for i, item := range tv.Output2 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  #%d %s (%s): %d주 (평균 %d주)\n",
			item.Rank, item.Name, item.Symb, item.Tvol, item.ATvol)
	}

	// 3. 거래대금순위
	tp, err := client.Overseas.InquireTradePbmn(ctx, overseas.InquireTradePbmnParams{
		ExcdCode: excd,
		NDay:     "0",
	})
	if err != nil {
		log.Fatalf("InquireTradePbmn: %v", err)
	}
	fmt.Printf("\n[%s 거래대금 상위 %d 건]\n", excd, len(tp.Output2))
	for i, item := range tp.Output2 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  #%d %s (%s): %d USD (평균 %d USD)\n",
			item.Rank, item.Name, item.Symb, item.Tamt, item.ATamt)
	}

	// 4. 거래량급증
	vs, err := client.Overseas.InquireVolumeSurge(ctx, overseas.InquireVolumeSurgeParams{
		ExcdCode: excd,
		MixN:     "0", // 1분전
	})
	if err != nil {
		log.Fatalf("InquireVolumeSurge: %v", err)
	}
	fmt.Printf("\n[%s 거래량급증 %d 건]\n", excd, len(vs.Output2))
	for i, item := range vs.Output2 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s (%s): 현재 %d주 / 기준 %d주 (+%v, %v%%)\n",
			item.Knam, item.Symb, item.Tvol, item.NTvol, item.NDiff, item.NRate)
	}

	// 5. 매수체결강도상위
	vp, err := client.Overseas.InquireVolumePower(ctx, overseas.InquireVolumePowerParams{
		ExcdCode: excd,
		NDay:     "0", // 1분전 (wire name NDAY, 실제 분 단위)
	})
	if err != nil {
		log.Fatalf("InquireVolumePower: %v", err)
	}
	fmt.Printf("\n[%s 매수체결강도 상위 %d 건]\n", excd, len(vp.Output2))
	for i, item := range vp.Output2 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s (%s): 당일체결강도 %v / 체결강도 %v\n",
			item.Knam, item.Symb, item.Tpow, item.Powx)
	}

	// 6. 신고/신저가
	nh, err := client.Overseas.InquireNewHighlow(ctx, overseas.InquireNewHighlowParams{
		ExcdCode: excd,
		Gubn:     "1", // 신고(1)
		Gubn2:    "1", // 돌파유지(1)
		NDay:     "6", // 52주
	})
	if err != nil {
		log.Fatalf("InquireNewHighlow: %v", err)
	}
	fmt.Printf("\n[%s 52주 신고가 %d 건]\n", excd, len(nh.Output2))
	for i, item := range nh.Output2 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s (%s): %s USD (기준가 %s, 대비 %s, %v%%)\n",
			item.Name, item.Symb, item.Last, item.NBase, item.NDiff, item.NRate)
	}
}
```

- [ ] **Step 2: 컴파일 검증**

`go build ./examples/overseas_ranking && echo OK`

- [ ] **Step 3: Commit**

```bash
git add examples/overseas_ranking
git commit -m "$(cat <<'EOF'
[feat] examples/overseas_ranking — 해외주식 ranking 6 메서드 통합 예시

시가총액순위 + 거래량순위 + 거래대금순위 + 거래량급증 + 매수체결강도상위 +
신고/신저가 (52주) 출력. NAS 거래소 + AAPL 샘플.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 9: 문서 갱신

**Files:**
- Modify: `CLAUDE.md`
- Modify: `README.md`
- Modify: `CHANGELOG.md`
- Modify: `overseas/doc.go`

- [ ] **Step 1: CLAUDE.md — banner 갱신**

Replace:
```
> **Phase 2.2 — 국내 신고저/시간외 (v1.5.0).** Phase 2.3+ 는 추후 sub-plan 으로.
```
With:
```
> **Phase 2.3 — 해외주식 추가 Ranking (v1.6.0).** Phase 2.4+ 는 추후 sub-plan 으로.
```

ADD spec link bullet after Phase 2.2 plan link:
```markdown
- Phase 2.3 implementation plan: [`docs/superpowers/specs/2026-05-05-phase2-3-overseas-ranking-implementation-plan.md`](docs/superpowers/specs/2026-05-05-phase2-3-overseas-ranking-implementation-plan.md)
```

- [ ] **Step 2: README.md — Available Methods 표 갱신**

Find existing `## Available Methods (Phase 1.2 ~ 2.2)` heading. Update heading to `## Available Methods (Phase 1.2 ~ 2.3)` and APPEND 6 rows AT THE END:

```markdown
| `Overseas.InquireMarketCap` | `ranking/market-cap` | HHDFS76350100 |
| `Overseas.InquireTradeVol` | `ranking/trade-vol` | HHDFS76310010 |
| `Overseas.InquireTradePbmn` | `ranking/trade-pbmn` | HHDFS76320010 |
| `Overseas.InquireVolumeSurge` | `ranking/volume-surge` | HHDFS76270000 |
| `Overseas.InquireVolumePower` | `ranking/volume-power` | HHDFS76280000 |
| `Overseas.InquireNewHighlow` | `ranking/new-highlow` | HHDFS76300000 |
```

Also update the method count in the heading if present (36 → 42).

- [ ] **Step 3: CHANGELOG.md — `[1.6.0]` entry**

ADD AT THE TOP (above `## [1.5.0]`):

```markdown
## [1.6.0] - 2026-05-05

### Added — Phase 2.3 (해외주식 추가 Ranking)

- `Overseas.InquireMarketCap` — 해외주식 시가총액순위 (HHDFS76350100) — output1 5-field + output2 15 fields/item
- `Overseas.InquireTradeVol` — 해외주식 거래량순위 (HHDFS76310010) — output1 5-field + output2 16 fields/item
- `Overseas.InquireTradePbmn` — 해외주식 거래대금순위 (HHDFS76320010) — output1 5-field + output2 16 fields/item (a_tamt)
- `Overseas.InquireVolumeSurge` — 해외주식 거래량급증 (HHDFS76270000) — output1 3-field + output2 16 fields/item (knam/enam)
- `Overseas.InquireVolumePower` — 해외주식 매수체결강도상위 (HHDFS76280000) — output1 3-field + output2 15 fields/item (knam/enam, tpow/powx)
- `Overseas.InquireNewHighlow` — 해외주식 신고/신저가 (HHDFS76300000) — output1 3-field + output2 16 fields/item (n_base/n_diff/n_rate)
- examples: `overseas_ranking`

### Notes

- output1 2 tier: `OverseasRankingFullSummary` (5-field: #1-#3) / `OverseasRankingMinSummary` (3-field: #4-#6, crec/trec 없음)
- output2 종목명 키 분기: InquireMarketCap/InquireTradeVol/InquireTradePbmn/InquireNewHighlow 는 `name`/`ename`, InquireVolumeSurge/InquireVolumePower 는 `knam`/`enam`
- `InquireVolumePower` 의 query 파라미터 `NDAY` 는 실제로 분(分) 단위 — KIS docs 명명 이슈. wire name 그대로 사용
- Phase 2.3 완료 — 누적 42 메서드 (Phase 2.2: 36 → Phase 2.3: 42)
```

- [ ] **Step 4: overseas/doc.go 갱신**

ADD Phase 2.3 section after existing section (before closing `package overseas` comment):

```go
// Phase 2.3 메서드 (6):
//
//   - InquireMarketCap   — 해외주식 시가총액순위 (HHDFS76350100)
//   - InquireTradeVol    — 해외주식 거래량순위 (HHDFS76310010)
//   - InquireTradePbmn   — 해외주식 거래대금순위 (HHDFS76320010)
//   - InquireVolumeSurge — 해외주식 거래량급증 (HHDFS76270000)
//   - InquireVolumePower — 해외주식 매수체결강도상위 (HHDFS76280000)
//   - InquireNewHighlow  — 해외주식 신고/신저가 (HHDFS76300000)
```

- [ ] **Step 5: 검증**

```bash
go build ./... && go vet ./... && gofmt -l .
```
Expected: silent.

- [ ] **Step 6: Commit**

```bash
git add CLAUDE.md README.md CHANGELOG.md overseas/doc.go
git commit -m "$(cat <<'EOF'
[doc] Phase 2.3 메서드 문서 갱신 — CLAUDE/README/CHANGELOG/overseas/doc.go

Phase 2.3 의 6 메서드 (시가총액/거래량/거래대금순위 + 거래량급증 + 매수체결강도
+ 신고신저가) 목록 + CHANGELOG [1.6.0] entry. CLAUDE.md banner 갱신
(Phase 2.2 → 2.3, v1.5.0 → v1.6.0). overseas/doc.go 패키지 doc 에 Phase 2.3
section 추가. README Available Methods 36 → 42.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 10: 최종 점검

- [ ] **Step 1: gofmt cleanup (필요 시)**

`gofmt -w overseas/*.go && gofmt -l .` — empty output.

- [ ] **Step 2: 빌드/vet/test**

```bash
go build ./... && go vet ./... && go test ./... -race -count=1
```
Expected: silent + all PASS.

- [ ] **Step 3: Coverage**

```bash
go test ./... -coverprofile=/tmp/cov.out -covermode=atomic
go tool cover -func=/tmp/cov.out | tail -10
```
Expected: overseas/ ≥ 80%, root kis ≥ 80%.

- [ ] **Step 4: 디렉터리 구조 확인**

```bash
ls -la \
  overseas/ranking.go \
  overseas/ranking_test.go \
  overseas/testdata/market_cap_success.json \
  overseas/testdata/trade_vol_success.json \
  overseas/testdata/trade_pbmn_success.json \
  overseas/testdata/volume_surge_success.json \
  overseas/testdata/volume_power_success.json \
  overseas/testdata/new_highlow_success.json \
  examples/overseas_ranking/main.go \
  2>&1 | wc -l
```
Expected: 9 lines.

- [ ] **Step 5: Commit history**

`git log main..HEAD --oneline | wc -l` — should be ~12-15.

---

## Task 11: PR 생성 (사용자 승인 후)

> Claude 는 push / PR 생성을 사용자 명시적 승인 후에만 실행 (글로벌 정책).

- [ ] **Step 1: 사용자 승인 요청**

작업 진행 보고 + PR 생성 가능 여부 confirm.

- [ ] **Step 2: Push branch**

`git push -u origin feat/phase2-3-overseas-ranking`

- [ ] **Step 3: PR 생성**

```bash
gh pr create --title "Phase 2.3 — 해외주식 추가 Ranking (v1.6.0)" --reviewer kenshin579 --base main --head feat/phase2-3-overseas-ranking --body "$(cat <<'EOF'
## Summary

- 해외주식 ranking 6 메서드 추가 (Phase 2 세 번째 sub-phase)
- Phase 1.5 / 2.2 패턴 그대로 재사용 (Style A, Params struct, KIS docs 1:1)
- v1.6.0 release 대상 (누적 36 → 42 메서드)

## 메서드 → 한투 API 매핑

| Go 메서드 | path | TR_ID |
|---|---|---|
| InquireMarketCap | ranking/market-cap | HHDFS76350100 |
| InquireTradeVol | ranking/trade-vol | HHDFS76310010 |
| InquireTradePbmn | ranking/trade-pbmn | HHDFS76320010 |
| InquireVolumeSurge | ranking/volume-surge | HHDFS76270000 |
| InquireVolumePower | ranking/volume-power | HHDFS76280000 |
| InquireNewHighlow | ranking/new-highlow | HHDFS76300000 |

## 구현 이슈 및 대응

| 이슈 | 대응 |
|---|---|
| output1 2 tier (5-field vs 3-field) | `OverseasRankingFullSummary` (#1-#3) / `OverseasRankingMinSummary` (#4-#6) 분리 정의 |
| output2 종목명 키 분기 (`name`/`ename` vs `knam`/`enam`) | 메서드별 전용 output2 struct 정의 |
| VolumePower NDAY 파라미터 분(分) 단위 불일치 | wire name NDAY 유지 + Params struct 주석에 명기 |

## Test Plan

- [x] go build/vet/fmt clean
- [x] go test ./... -race -count=1 모든 패키지 PASS
- [x] Coverage overseas/ >= 80%
- [x] httpmock 단위 테스트 (6 메서드)
- [x] output1 2-tier 타입 분리 검증 (Crec/Trec 유무)
- [x] output2 knam/enam vs name/ename 분기 검증

## Breaking Changes

없음 — 신규 메서드 추가만.

## 참고 문서

- Phase 2 design spec: `docs/superpowers/specs/2026-05-05-phase2-readonly-extension-design.md`
- Phase 2.3 implementation plan: `docs/superpowers/specs/2026-05-05-phase2-3-overseas-ranking-implementation-plan.md`

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

- [ ] **Step 4: Merge (사용자 승인 후)** — `gh pr merge <PR#> --merge`

- [ ] **Step 5: 후속 작업**

```bash
git tag -a v1.6.0 -m "Phase 2.3 — 해외주식 추가 Ranking (6 메서드, 누적 42)"
git push origin v1.6.0
gh release create v1.6.0 --title "v1.6.0 — Phase 2.3 해외주식 추가 Ranking" \
  --notes-file <(awk '/^## \[1\.6\.0\]/{p=1} p && /^## \[1\.5\.0\]/{exit} p' CHANGELOG.md)
```
