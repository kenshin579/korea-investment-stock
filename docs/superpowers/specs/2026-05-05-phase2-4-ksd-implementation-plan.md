# Phase 2.4 — 예탁원 정보 확장 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** korea-investment-stock Go 라이브러리에 예탁원 정보 11 메서드 추가 (`v1.7.0` release).

**Architecture:** Phase 1 인프라 + 패턴 재사용. `domestic/ksd.go` 신규 (11 메서드 한 file). 한투 path 1:1 매핑 (Style A) + `Ksd` prefix (Phase 2 design spec §3 — 가독성 우선). 새 internal package 불필요. TDD: testdata fixture → 실패 테스트 → struct + 메서드 구현 → 통과 → commit.

**Tech Stack:** Go 1.25+, `github.com/jarcoal/httpmock`, `github.com/stretchr/testify`. 새 dependency 없음.

**참고 spec:**
- Phase 2 design spec: `docs/superpowers/specs/2026-05-05-phase2-readonly-extension-design.md` (§Phase 2.4)
- Phase 1.4 plan (참조 패턴 — InquirePubOffer): `docs/superpowers/specs/2026-05-05-phase1-4-investor-industry-ipo-implementation-plan.md`
- 한투 API docs: `docs/api/국내주식/예탁원정보(*).md` (11개)

---

## 사전 정보

| 항목 | 값 |
|------|---|
| 작업 브랜치 | `feat/phase2-4-ksd` |
| 시작 HEAD | Phase 2.3 구현 완료 commit (v1.6.0) |
| Release 목표 | `v1.7.0` |
| PR 베이스 | `main` |
| 현재 main 상태 | v1.6.0 publish 완료 (Phase 2.3 통합, 42 메서드) |

---

## 메서드 매핑

| Go 메서드 | path (last segment) | TR_ID | output key | fields | anomalies |
|---|---|---|---|---|---|
| `InquireKsdDividend` | `dividend` | HHKDB669102C0 | `output1 []` | 13 | — |
| `InquireKsdBonusIssue` | `bonus-issue` | HHKDB669101C0 | `output1 []` | 11 | — |
| `InquireKsdPaidinCapin` | `paidin-capin` | HHKDB669100C0 | `output []` | 13 | **output (not output1)** |
| `InquireKsdSharehldMeet` | `sharehld-meet` | HHKDB669111C0 | `output1 []` | 7 | — |
| `InquireKsdMergerSplit` | `merger-split` | HHKDB669104C0 | `output1 []` | 14 | **no isin_name; opp_cust_*/cust_* pairs** |
| `InquireKsdRevSplit` | `rev-split` | HHKDB669105C0 | `output1 []` | 7 | **extra MARKET_GB param** |
| `InquireKsdForfeit` | `forfeit` | HHKDB669109C0 | `output1 []` | 9 | — |
| `InquireKsdMandDeposit` | `mand-deposit` | HHKDB669110C0 | `output1 []` | 6 | **no record_date; depo_date as date key** |
| `InquireKsdCapDcrs` | `cap-dcrs` | HHKDB669106C0 | `output1 []` | 9 | — |
| `InquireKsdPurreq` | `purreq` | HHKDB669103C0 | `output1 []` | 9 | — |
| `InquireKsdListInfo` | `list-info` | HHKDB669107C0 | `output1 []` | 8 | **leading field list_dt (not record_date)** |

---

## 파일 구조

### 신규
- `domestic/ksd.go` — 11 메서드 + structs + Params
- `domestic/ksd_test.go` — 11 테스트 함수
- `domestic/testdata/ksd_dividend_success.json`
- `domestic/testdata/ksd_bonus_issue_success.json`
- `domestic/testdata/ksd_paidin_capin_success.json`
- `domestic/testdata/ksd_sharehld_meet_success.json`
- `domestic/testdata/ksd_merger_split_success.json`
- `domestic/testdata/ksd_rev_split_success.json`
- `domestic/testdata/ksd_forfeit_success.json`
- `domestic/testdata/ksd_mand_deposit_success.json`
- `domestic/testdata/ksd_cap_dcrs_success.json`
- `domestic/testdata/ksd_purreq_success.json`
- `domestic/testdata/ksd_list_info_success.json`
- `examples/domestic_ksd/main.go`

### 수정
- `CLAUDE.md` — banner Phase 2.3 → Phase 2.4, plan link 추가
- `README.md` — Available Methods 표 갱신 (42 → 53 메서드)
- `CHANGELOG.md` — `[1.7.0]` entry
- `domestic/doc.go` — Phase 2.4 section 추가

---

## 타입 매핑

KSD 의 모든 응답 필드는 KIS docs 에서 String 타입으로 명시 — Go 에서도 plain `string`. 가격/수량/비율/날짜/코드 모두 동일 (Phase 2.x 의 decimal/int64/float64 매핑 미적용). `InquirePubOffer` (`ipo.go`) 와 달리 KSD 나머지 11 메서드는 숫자 타입 변환 없음.

---

## Tasks (16 total)

---

## Task 1: testdata fixtures (11 합성 JSON)

- [ ] Step 1: `domestic/testdata/ksd_dividend_success.json`
- [ ] Step 2: `domestic/testdata/ksd_bonus_issue_success.json`
- [ ] Step 3: `domestic/testdata/ksd_paidin_capin_success.json`
- [ ] Step 4: `domestic/testdata/ksd_sharehld_meet_success.json`
- [ ] Step 5: `domestic/testdata/ksd_merger_split_success.json`
- [ ] Step 6: `domestic/testdata/ksd_rev_split_success.json`
- [ ] Step 7: `domestic/testdata/ksd_forfeit_success.json`
- [ ] Step 8: `domestic/testdata/ksd_mand_deposit_success.json`
- [ ] Step 9: `domestic/testdata/ksd_cap_dcrs_success.json`
- [ ] Step 10: `domestic/testdata/ksd_purreq_success.json`
- [ ] Step 11: `domestic/testdata/ksd_list_info_success.json`
- [ ] Step 12: validation

```bash
for f in \
  domestic/testdata/ksd_dividend_success.json \
  domestic/testdata/ksd_bonus_issue_success.json \
  domestic/testdata/ksd_paidin_capin_success.json \
  domestic/testdata/ksd_sharehld_meet_success.json \
  domestic/testdata/ksd_merger_split_success.json \
  domestic/testdata/ksd_rev_split_success.json \
  domestic/testdata/ksd_forfeit_success.json \
  domestic/testdata/ksd_mand_deposit_success.json \
  domestic/testdata/ksd_cap_dcrs_success.json \
  domestic/testdata/ksd_purreq_success.json \
  domestic/testdata/ksd_list_info_success.json; do
  python3 -c "import json; json.load(open('$f'))" && echo "$f OK"
done
# Expected: 11 OK lines
```

- [ ] Step 13: commit

```bash
git commit -m "$(cat <<'EOF'
[chore] testdata — 11 KSD fixture JSON (Phase 2.4)

합성 JSON fixtures (2 records each, 005930/000660):
- ksd_dividend_success.json (13 fields)
- ksd_bonus_issue_success.json (11 fields)
- ksd_paidin_capin_success.json (13 fields, output key)
- ksd_sharehld_meet_success.json (7 fields)
- ksd_merger_split_success.json (14 fields, no isin_name)
- ksd_rev_split_success.json (7 fields)
- ksd_forfeit_success.json (9 fields)
- ksd_mand_deposit_success.json (6 fields, depo_date)
- ksd_cap_dcrs_success.json (9 fields)
- ksd_purreq_success.json (9 fields)
- ksd_list_info_success.json (8 fields, list_dt)

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

### Fixture content

**Step 1 — `ksd_dividend_success.json`**
```json
{
  "output1": [
    {
      "record_date": "20260331",
      "sht_cd": "005930",
      "isin_name": "삼성전자",
      "divi_kind": "결산",
      "face_val": "100",
      "per_sto_divi_amt": "361",
      "divi_rate": "3.61",
      "stk_divi_rate": "0.00",
      "divi_pay_dt": "20260430",
      "stk_div_pay_dt": "",
      "odd_pay_dt": "",
      "stk_kind": "보통주",
      "high_divi_gb": "N"
    },
    {
      "record_date": "20260331",
      "sht_cd": "000660",
      "isin_name": "SK하이닉스",
      "divi_kind": "결산",
      "face_val": "5000",
      "per_sto_divi_amt": "1200",
      "divi_rate": "2.40",
      "stk_divi_rate": "0.00",
      "divi_pay_dt": "20260430",
      "stk_div_pay_dt": "",
      "odd_pay_dt": "",
      "stk_kind": "보통주",
      "high_divi_gb": "N"
    }
  ]
}
```

**Step 2 — `ksd_bonus_issue_success.json`**
```json
{
  "output1": [
    {
      "record_date": "20260315",
      "sht_cd": "005930",
      "isin_name": "삼성전자",
      "fix_rate": "0.05",
      "odd_rec_price": "70000",
      "right_dt": "20260314",
      "odd_pay_dt": "20260416",
      "list_date": "20260417",
      "tot_issue_stk_qty": "5969782550",
      "issue_stk_qty": "298489127",
      "stk_kind": "보통주"
    },
    {
      "record_date": "20260315",
      "sht_cd": "000660",
      "isin_name": "SK하이닉스",
      "fix_rate": "0.10",
      "odd_rec_price": "130000",
      "right_dt": "20260314",
      "odd_pay_dt": "20260416",
      "list_date": "20260417",
      "tot_issue_stk_qty": "728002365",
      "issue_stk_qty": "72800236",
      "stk_kind": "보통주"
    }
  ]
}
```

**Step 3 — `ksd_paidin_capin_success.json`** (output key, not output1)
```json
{
  "output": [
    {
      "record_date": "20260310",
      "sht_cd": "005930",
      "isin_name": "삼성전자",
      "tot_issue_stk_qty": "5969782550",
      "issue_stk_qty": "100000000",
      "fix_rate": "0.10",
      "disc_rate": "0.00",
      "fix_price": "68000",
      "right_dt": "20260309",
      "sub_term_ft": "20260312",
      "sub_term": "20260313",
      "list_date": "20260325",
      "stk_kind": "보통주"
    },
    {
      "record_date": "20260310",
      "sht_cd": "000660",
      "isin_name": "SK하이닉스",
      "tot_issue_stk_qty": "728002365",
      "issue_stk_qty": "30000000",
      "fix_rate": "0.10",
      "disc_rate": "10.00",
      "fix_price": "115000",
      "right_dt": "20260309",
      "sub_term_ft": "20260312",
      "sub_term": "20260313",
      "list_date": "20260325",
      "stk_kind": "보통주"
    }
  ]
}
```

**Step 4 — `ksd_sharehld_meet_success.json`**
```json
{
  "output1": [
    {
      "record_date": "20260331",
      "sht_cd": "005930",
      "isin_name": "삼성전자",
      "gen_meet_dt": "20260326",
      "gen_meet_type": "정기주총",
      "agenda": "이익배당 승인의 건",
      "vote_tot_qty": "5969782550"
    },
    {
      "record_date": "20260331",
      "sht_cd": "000660",
      "isin_name": "SK하이닉스",
      "gen_meet_dt": "20260327",
      "gen_meet_type": "정기주총",
      "agenda": "재무제표 승인의 건",
      "vote_tot_qty": "728002365"
    }
  ]
}
```

**Step 5 — `ksd_merger_split_success.json`** (no isin_name)
```json
{
  "output1": [
    {
      "record_date": "20260401",
      "sht_cd": "005930",
      "opp_cust_cd": "999999",
      "opp_cust_nm": "흡수대상회사A",
      "cust_cd": "005930",
      "cust_nm": "삼성전자",
      "merge_type": "흡수합병",
      "merge_rate": "1.00",
      "td_stop_dt": "20260330 ~ 20260401",
      "list_dt": "20260402",
      "odd_amt_pay_dt": "20260410",
      "tot_issue_stk_qty": "5969782550",
      "issue_stk_qty": "1000000",
      "seq": "1"
    },
    {
      "record_date": "20260401",
      "sht_cd": "000660",
      "opp_cust_cd": "888888",
      "opp_cust_nm": "흡수대상회사B",
      "cust_cd": "000660",
      "cust_nm": "SK하이닉스",
      "merge_type": "흡수합병",
      "merge_rate": "0.50",
      "td_stop_dt": "20260330 ~ 20260401",
      "list_dt": "20260402",
      "odd_amt_pay_dt": "20260410",
      "tot_issue_stk_qty": "728002365",
      "issue_stk_qty": "500000",
      "seq": "2"
    }
  ]
}
```

**Step 6 — `ksd_rev_split_success.json`**
```json
{
  "output1": [
    {
      "record_date": "20260415",
      "sht_cd": "005930",
      "isin_name": "삼성전자",
      "inter_bf_face_amt": "100",
      "inter_af_face_amt": "500",
      "td_stop_dt": "20260414 ~ 20260416",
      "list_dt": "20260417"
    },
    {
      "record_date": "20260415",
      "sht_cd": "000660",
      "isin_name": "SK하이닉스",
      "inter_bf_face_amt": "5000",
      "inter_af_face_amt": "10000",
      "td_stop_dt": "20260414 ~ 20260416",
      "list_dt": "20260417"
    }
  ]
}
```

**Step 7 — `ksd_forfeit_success.json`**
```json
{
  "output1": [
    {
      "record_date": "20260310",
      "sht_cd": "005930",
      "isin_name": "삼성전자",
      "subscr_dt": "20260312 ~ 20260313",
      "subscr_price": "68000",
      "subscr_stk_qty": "100000000",
      "refund_dt": "20260318",
      "list_dt": "20260325",
      "lead_mgr": "한국투자증권"
    },
    {
      "record_date": "20260310",
      "sht_cd": "000660",
      "isin_name": "SK하이닉스",
      "subscr_dt": "20260312 ~ 20260313",
      "subscr_price": "115000",
      "subscr_stk_qty": "30000000",
      "refund_dt": "20260318",
      "list_dt": "20260325",
      "lead_mgr": "미래에셋증권"
    }
  ]
}
```

**Step 8 — `ksd_mand_deposit_success.json`** (no record_date, depo_date as date key)
```json
{
  "output1": [
    {
      "sht_cd": "005930",
      "isin_name": "삼성전자",
      "stk_qty": "50000000",
      "depo_date": "20260101",
      "depo_reason": "의무보호예수",
      "tot_issue_qty_per_rate": "0.84"
    },
    {
      "sht_cd": "000660",
      "isin_name": "SK하이닉스",
      "stk_qty": "10000000",
      "depo_date": "20260101",
      "depo_reason": "의무보호예수",
      "tot_issue_qty_per_rate": "1.37"
    }
  ]
}
```

**Step 9 — `ksd_cap_dcrs_success.json`**
```json
{
  "output1": [
    {
      "record_date": "20260320",
      "sht_cd": "005930",
      "isin_name": "삼성전자",
      "stk_kind": "보통주",
      "reduce_cap_type": "유상감자",
      "reduce_cap_rate": "0.20",
      "comp_way": "주식병합",
      "td_stop_dt": "20260319 ~ 20260321",
      "list_dt": "20260322"
    },
    {
      "record_date": "20260320",
      "sht_cd": "000660",
      "isin_name": "SK하이닉스",
      "stk_kind": "보통주",
      "reduce_cap_type": "무상감자",
      "reduce_cap_rate": "0.50",
      "comp_way": "주식소각",
      "td_stop_dt": "20260319 ~ 20260321",
      "list_dt": "20260322"
    }
  ]
}
```

**Step 10 — `ksd_purreq_success.json`**
```json
{
  "output1": [
    {
      "record_date": "20260331",
      "sht_cd": "005930",
      "isin_name": "삼성전자",
      "stk_kind": "보통주",
      "opp_opi_rcpt_term": "20260310 ~ 20260322",
      "buy_req_rcpt_term": "20260324 ~ 20260326",
      "buy_req_price": "69000",
      "buy_amt_pay_dt": "20260410",
      "get_meet_dt": "20260326"
    },
    {
      "record_date": "20260331",
      "sht_cd": "000660",
      "isin_name": "SK하이닉스",
      "stk_kind": "보통주",
      "opp_opi_rcpt_term": "20260310 ~ 20260322",
      "buy_req_rcpt_term": "20260324 ~ 20260326",
      "buy_req_price": "120000",
      "buy_amt_pay_dt": "20260410",
      "get_meet_dt": "20260327"
    }
  ]
}
```

**Step 11 — `ksd_list_info_success.json`** (leading field is list_dt, not record_date)
```json
{
  "output1": [
    {
      "list_dt": "20260102",
      "sht_cd": "005930",
      "isin_name": "삼성전자",
      "stk_kind": "보통주",
      "issue_type": "유상증자",
      "issue_stk_qty": "50000000",
      "tot_issue_stk_qty": "5969782550",
      "issue_price": "68000"
    },
    {
      "list_dt": "20260102",
      "sht_cd": "000660",
      "isin_name": "SK하이닉스",
      "stk_kind": "보통주",
      "issue_type": "유상증자",
      "issue_stk_qty": "30000000",
      "tot_issue_stk_qty": "728002365",
      "issue_price": "115000"
    }
  ]
}
```

---

## Task 2: `domestic/ksd.go` base + InquireKsdDividend

**Files:**
- Create: `domestic/ksd.go`
- Create: `domestic/ksd_test.go`

- [ ] Step 1: create `domestic/ksd.go` with package header and `InquireKsdDividend`
- [ ] Step 2: create `domestic/ksd_test.go` with `TestClient_InquireKsdDividend`
- [ ] Step 3: `go test ./domestic/... -run TestClient_InquireKsdDividend -v` — PASS
- [ ] Step 4: commit

```go
// File: domestic/ksd.go
package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// KsdDividend 는 예탁원정보(배당일정) (HHKDB669102C0) 응답.
//
// 한투 docs: docs/api/국내주식/예탁원정보(배당일정).md
// path: /uapi/domestic-stock/v1/ksdinfo/dividend
type KsdDividend struct {
	Output1 []KsdDividendItem `json:"output1"`
}

// KsdDividendItem 은 배당일정 한 행. 모든 필드 string (KIS docs).
type KsdDividendItem struct {
	RecordDate    string `json:"record_date"`     // 기준일
	ShtCd         string `json:"sht_cd"`           // 종목코드
	IsinName      string `json:"isin_name"`        // 종목명
	DiviKind      string `json:"divi_kind"`        // 배당종류
	FaceVal       string `json:"face_val"`         // 액면가
	PerStoDiviAmt string `json:"per_sto_divi_amt"` // 현금배당금
	DiviRate      string `json:"divi_rate"`        // 현금배당률(%)
	StkDiviRate   string `json:"stk_divi_rate"`    // 주식배당률(%)
	DiviPayDt     string `json:"divi_pay_dt"`      // 배당금지급일
	StkDivPayDt   string `json:"stk_div_pay_dt"`   // 주식배당지급일
	OddPayDt      string `json:"odd_pay_dt"`       // 단주대금지급일
	StkKind       string `json:"stk_kind"`         // 주식종류
	HighDiviGb    string `json:"high_divi_gb"`     // 고배당종목여부
}

// InquireKsdDividendParams 는 배당일정 조회 파라미터.
type InquireKsdDividendParams struct {
	Cts      string // CTS — 공백 입력 (default "")
	Gb1      string // GB1 — 0:전체, 1:결산배당, 2:중간배당. 빈 값=>"0"
	FromDate string // F_DT — 조회시작일 YYYYMMDD
	ToDate   string // T_DT — 조회종료일 YYYYMMDD
	Symbol   string // SHT_CD — 종목코드 (공백=전체)
	HighGb   string // HIGH_GB — 공백 입력
}

// InquireKsdDividend 호출.
//
// 한투 docs: docs/api/국내주식/예탁원정보(배당일정).md
// path: /uapi/domestic-stock/v1/ksdinfo/dividend (HHKDB669102C0)
func (c *Client) InquireKsdDividend(ctx context.Context, params InquireKsdDividendParams) (*KsdDividend, error) {
	gb1 := params.Gb1
	if gb1 == "" {
		gb1 = "0"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ksdinfo/dividend",
		TrID:   "HHKDB669102C0",
		Query: map[string]string{
			"CTS":     params.Cts,
			"GB1":     gb1,
			"F_DT":    params.FromDate,
			"T_DT":    params.ToDate,
			"SHT_CD":  params.Symbol,
			"HIGH_GB": params.HighGb,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res KsdDividend
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse KsdDividend: %w", err)
	}
	return &res, nil
}
```

```go
// File: domestic/ksd_test.go
package domestic_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/domestic"
)

func TestClient_InquireKsdDividend(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ksdinfo/dividend`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "ksd_dividend_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireKsdDividend(context.Background(), domestic.InquireKsdDividendParams{
		FromDate: "20260101",
		ToDate:   "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "0", capturedQuery.Get("GB1"))
	assert.Equal(t, "20260101", capturedQuery.Get("F_DT"))
	assert.Equal(t, "20260505", capturedQuery.Get("T_DT"))

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "005930", res.Output1[0].ShtCd)
	assert.Equal(t, "삼성전자", res.Output1[0].IsinName)
	assert.Equal(t, "20260331", res.Output1[0].RecordDate)
}
```

Commit:
```bash
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireKsdDividend (예탁원정보 배당일정, HHKDB669102C0)

- KsdDividend / KsdDividendItem (13 fields, all string)
- InquireKsdDividendParams (6 query params: CTS/GB1/F_DT/T_DT/SHT_CD/HIGH_GB)
- GB1 default "0" (전체) when empty
- TestClient_InquireKsdDividend — fixture ksd_dividend_success.json

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 3: InquireKsdBonusIssue

**Files:** Modify `domestic/ksd.go` (append), `domestic/ksd_test.go` (append)

무상증자 일정 조회. `output1 []KsdBonusIssueItem` 11 fields.

- [ ] Step 1: append struct + method to `domestic/ksd.go`
- [ ] Step 2: append test to `domestic/ksd_test.go`
- [ ] Step 3: `go test ./domestic/... -run TestClient_InquireKsdBonusIssue -v` — PASS
- [ ] Step 4: commit

### 메서드 매핑
- Path: `/uapi/domestic-stock/v1/ksdinfo/bonus-issue`
- TR_ID: `HHKDB669101C0`
- Output key: `output1 []KsdBonusIssueItem`

### 응답 struct 필드 (모두 string)

| Go field | json tag | 설명 |
|---|---|---|
| `RecordDate` | `record_date` | 기준일 |
| `ShtCd` | `sht_cd` | 종목코드 |
| `IsinName` | `isin_name` | 종목명 |
| `FixRate` | `fix_rate` | 확정배정율 |
| `OddRecPrice` | `odd_rec_price` | 단주기준가 |
| `RightDt` | `right_dt` | 권리락일 |
| `OddPayDt` | `odd_pay_dt` | 단주대금지급일 |
| `ListDate` | `list_date` | 상장/등록일 |
| `TotIssueStkQty` | `tot_issue_stk_qty` | 발행주식 |
| `IssueStkQty` | `issue_stk_qty` | 발행할주식 |
| `StkKind` | `stk_kind` | 주식종류 |

### Params struct fields

| Go field | wire name | default |
|---|---|---|
| `Cts` | `CTS` | `""` |
| `FromDate` | `F_DT` | required |
| `ToDate` | `T_DT` | required |
| `Symbol` | `SHT_CD` | `""` (전체) |

### Test code

```go
func TestClient_InquireKsdBonusIssue(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ksdinfo/bonus-issue`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "ksd_bonus_issue_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireKsdBonusIssue(context.Background(), domestic.InquireKsdBonusIssueParams{
		FromDate: "20260101",
		ToDate:   "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "20260101", capturedQuery.Get("F_DT"))
	assert.Equal(t, "20260505", capturedQuery.Get("T_DT"))

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "005930", res.Output1[0].ShtCd)
	assert.Equal(t, "삼성전자", res.Output1[0].IsinName)
	assert.Equal(t, "20260315", res.Output1[0].RecordDate)
}
```

### Commit
```bash
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireKsdBonusIssue (예탁원정보 무상증자, HHKDB669101C0)

- KsdBonusIssue / KsdBonusIssueItem (11 fields, all string)
- InquireKsdBonusIssueParams (4 query params: CTS/F_DT/T_DT/SHT_CD)
- TestClient_InquireKsdBonusIssue — fixture ksd_bonus_issue_success.json

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 4: InquireKsdPaidinCapin

**Files:** Modify `domestic/ksd.go` (append), `domestic/ksd_test.go` (append)

유상증자 일정 조회. **ANOMALY: output key는 `output` (not `output1`).** `output []KsdPaidinCapinItem` 13 fields.

- [ ] Step 1: append struct + method to `domestic/ksd.go`
- [ ] Step 2: append test to `domestic/ksd_test.go`
- [ ] Step 3: `go test ./domestic/... -run TestClient_InquireKsdPaidinCapin -v` — PASS
- [ ] Step 4: commit

### 메서드 매핑
- Path: `/uapi/domestic-stock/v1/ksdinfo/paidin-capin`
- TR_ID: `HHKDB669100C0`
- Output key: **`output []KsdPaidinCapinItem`** (NOT `output1`)

### 응답 struct 필드 (모두 string)

| Go field | json tag | 설명 |
|---|---|---|
| `RecordDate` | `record_date` | 기준일 |
| `ShtCd` | `sht_cd` | 종목코드 |
| `IsinName` | `isin_name` | 종목명 |
| `TotIssueStkQty` | `tot_issue_stk_qty` | 발행주식 |
| `IssueStkQty` | `issue_stk_qty` | 발행할주식 |
| `FixRate` | `fix_rate` | 확정배정율 |
| `DiscRate` | `disc_rate` | 할인율 |
| `FixPrice` | `fix_price` | 발행예정가 |
| `RightDt` | `right_dt` | 권리락일 |
| `SubTermFt` | `sub_term_ft` | 청약기간(시작) |
| `SubTerm` | `sub_term` | 청약기간(종료) |
| `ListDate` | `list_date` | 상장/등록일 |
| `StkKind` | `stk_kind` | 주식종류 |

### Params struct fields

| Go field | wire name | default |
|---|---|---|
| `Cts` | `CTS` | `""` |
| `Gb1` | `GB1` | `"1"` (1:청약일별, 2:기준일별) |
| `FromDate` | `F_DT` | required |
| `ToDate` | `T_DT` | required |
| `Symbol` | `SHT_CD` | `""` (전체) |

### Anomaly callout
**`KsdPaidinCapin` struct 는 `Output []KsdPaidinCapinItem \`json:"output"\`` (NOT `output1`).** 다른 KSD 메서드와 다름 — KIS API 응답 구조 그대로.

### Test code

```go
func TestClient_InquireKsdPaidinCapin(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ksdinfo/paidin-capin`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "ksd_paidin_capin_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireKsdPaidinCapin(context.Background(), domestic.InquireKsdPaidinCapinParams{
		FromDate: "20260101",
		ToDate:   "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "1", capturedQuery.Get("GB1")) // default
	assert.Equal(t, "20260101", capturedQuery.Get("F_DT"))

	require.Len(t, res.Output, 2) // output (not Output1)
	assert.Equal(t, "005930", res.Output[0].ShtCd)
	assert.Equal(t, "삼성전자", res.Output[0].IsinName)
	assert.Equal(t, "68000", res.Output[0].FixPrice)
}
```

### Commit
```bash
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireKsdPaidinCapin (예탁원정보 유상증자, HHKDB669100C0)

- KsdPaidinCapin / KsdPaidinCapinItem (13 fields, all string)
- Output key is `output` (not output1) — KIS API 응답 구조
- InquireKsdPaidinCapinParams (5 query params: CTS/GB1/F_DT/T_DT/SHT_CD)
- GB1 default "1" (청약일별) when empty
- TestClient_InquireKsdPaidinCapin — fixture ksd_paidin_capin_success.json

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 5: InquireKsdSharehldMeet

**Files:** Modify `domestic/ksd.go` (append), `domestic/ksd_test.go` (append)

주주총회 일정 조회. `output1 []KsdSharehldMeetItem` 7 fields.

- [ ] Step 1: append struct + method to `domestic/ksd.go`
- [ ] Step 2: append test to `domestic/ksd_test.go`
- [ ] Step 3: `go test ./domestic/... -run TestClient_InquireKsdSharehldMeet -v` — PASS
- [ ] Step 4: commit

### 메서드 매핑
- Path: `/uapi/domestic-stock/v1/ksdinfo/sharehld-meet`
- TR_ID: `HHKDB669111C0`
- Output key: `output1 []KsdSharehldMeetItem`

### 응답 struct 필드 (모두 string)

| Go field | json tag | 설명 |
|---|---|---|
| `RecordDate` | `record_date` | 기준일 |
| `ShtCd` | `sht_cd` | 종목코드 |
| `IsinName` | `isin_name` | 종목명 |
| `GenMeetDt` | `gen_meet_dt` | 주총일자 |
| `GenMeetType` | `gen_meet_type` | 주총사유 |
| `Agenda` | `agenda` | 주총의안 |
| `VoteTotQty` | `vote_tot_qty` | 의결권주식총수 |

### Params struct fields

| Go field | wire name | default |
|---|---|---|
| `Cts` | `CTS` | `""` |
| `FromDate` | `F_DT` | required |
| `ToDate` | `T_DT` | required |
| `Symbol` | `SHT_CD` | `""` (전체) |

### Test code

```go
func TestClient_InquireKsdSharehldMeet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ksdinfo/sharehld-meet`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "ksd_sharehld_meet_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireKsdSharehldMeet(context.Background(), domestic.InquireKsdSharehldMeetParams{
		FromDate: "20260101",
		ToDate:   "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "20260101", capturedQuery.Get("F_DT"))
	assert.Equal(t, "20260505", capturedQuery.Get("T_DT"))

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "005930", res.Output1[0].ShtCd)
	assert.Equal(t, "정기주총", res.Output1[0].GenMeetType)
	assert.Equal(t, "20260326", res.Output1[0].GenMeetDt)
}
```

### Commit
```bash
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireKsdSharehldMeet (예탁원정보 주주총회, HHKDB669111C0)

- KsdSharehldMeet / KsdSharehldMeetItem (7 fields, all string)
- InquireKsdSharehldMeetParams (4 query params: CTS/F_DT/T_DT/SHT_CD)
- TestClient_InquireKsdSharehldMeet — fixture ksd_sharehld_meet_success.json

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 6: InquireKsdMergerSplit

**Files:** Modify `domestic/ksd.go` (append), `domestic/ksd_test.go` (append)

합병/분할 일정 조회. `output1 []KsdMergerSplitItem` 14 fields. **ANOMALY: `isin_name` 없음 — `opp_cust_*`/`cust_*` pair 사용.**

- [ ] Step 1: append struct + method to `domestic/ksd.go`
- [ ] Step 2: append test to `domestic/ksd_test.go`
- [ ] Step 3: `go test ./domestic/... -run TestClient_InquireKsdMergerSplit -v` — PASS
- [ ] Step 4: commit

### 메서드 매핑
- Path: `/uapi/domestic-stock/v1/ksdinfo/merger-split`
- TR_ID: `HHKDB669104C0`
- Output key: `output1 []KsdMergerSplitItem`

### 응답 struct 필드 (모두 string)

| Go field | json tag | 설명 |
|---|---|---|
| `RecordDate` | `record_date` | 기준일 |
| `ShtCd` | `sht_cd` | 종목코드 |
| `OppCustCd` | `opp_cust_cd` | 피합병(피분할)회사코드 |
| `OppCustNm` | `opp_cust_nm` | 피합병(피분할)회사명 |
| `CustCd` | `cust_cd` | 합병(분할)회사코드 |
| `CustNm` | `cust_nm` | 합병(분할)회사명 |
| `MergeType` | `merge_type` | 합병사유 |
| `MergeRate` | `merge_rate` | 비율 |
| `TdStopDt` | `td_stop_dt` | 매매거래정지기간 |
| `ListDt` | `list_dt` | 상장/등록일 |
| `OddAmtPayDt` | `odd_amt_pay_dt` | 단주대금지급일 |
| `TotIssueStkQty` | `tot_issue_stk_qty` | 발행주식 |
| `IssueStkQty` | `issue_stk_qty` | 발행할주식 |
| `Seq` | `seq` | 연번 |

### Params struct fields

| Go field | wire name | default |
|---|---|---|
| `Cts` | `CTS` | `""` |
| `FromDate` | `F_DT` | required |
| `ToDate` | `T_DT` | required |
| `Symbol` | `SHT_CD` | `""` (전체) |

### Anomaly callout
**`isin_name` 없음** — 합병/피합병 양사 정보를 `opp_cust_cd`/`opp_cust_nm` (피합병측) + `cust_cd`/`cust_nm` (합병측) 으로 표현. 다른 KSD 메서드의 `IsinName` 패턴과 다름.

### Test code

```go
func TestClient_InquireKsdMergerSplit(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ksdinfo/merger-split`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "ksd_merger_split_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireKsdMergerSplit(context.Background(), domestic.InquireKsdMergerSplitParams{
		FromDate: "20260101",
		ToDate:   "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "20260101", capturedQuery.Get("F_DT"))

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "005930", res.Output1[0].ShtCd)
	assert.Equal(t, "삼성전자", res.Output1[0].CustNm)       // cust_nm (합병측)
	assert.Equal(t, "흡수대상회사A", res.Output1[0].OppCustNm) // opp_cust_nm (피합병측)
	assert.Equal(t, "흡수합병", res.Output1[0].MergeType)
}
```

### Commit
```bash
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireKsdMergerSplit (예탁원정보 합병분할, HHKDB669104C0)

- KsdMergerSplit / KsdMergerSplitItem (14 fields, all string)
- No isin_name — uses opp_cust_cd/opp_cust_nm (피합병) + cust_cd/cust_nm (합병)
- InquireKsdMergerSplitParams (4 query params: CTS/F_DT/T_DT/SHT_CD)
- TestClient_InquireKsdMergerSplit — fixture ksd_merger_split_success.json

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 7: InquireKsdRevSplit

**Files:** Modify `domestic/ksd.go` (append), `domestic/ksd_test.go` (append)

액면분할/병합 일정 조회. `output1 []KsdRevSplitItem` 7 fields. **ANOMALY: extra `MARKET_GB` query param (default "0").**

- [ ] Step 1: append struct + method to `domestic/ksd.go`
- [ ] Step 2: append test to `domestic/ksd_test.go`
- [ ] Step 3: `go test ./domestic/... -run TestClient_InquireKsdRevSplit -v` — PASS
- [ ] Step 4: commit

### 메서드 매핑
- Path: `/uapi/domestic-stock/v1/ksdinfo/rev-split`
- TR_ID: `HHKDB669105C0`
- Output key: `output1 []KsdRevSplitItem`

### 응답 struct 필드 (모두 string)

| Go field | json tag | 설명 |
|---|---|---|
| `RecordDate` | `record_date` | 기준일 |
| `ShtCd` | `sht_cd` | 종목코드 |
| `IsinName` | `isin_name` | 종목명 |
| `InterBfFaceAmt` | `inter_bf_face_amt` | 변경전액면가 |
| `InterAfFaceAmt` | `inter_af_face_amt` | 변경후액면가 |
| `TdStopDt` | `td_stop_dt` | 매매거래정지기간 |
| `ListDt` | `list_dt` | 상장/등록일 |

### Params struct fields

| Go field | wire name | default |
|---|---|---|
| `Symbol` | `SHT_CD` | `""` (전체) |
| `Cts` | `CTS` | `""` |
| `FromDate` | `F_DT` | required |
| `ToDate` | `T_DT` | required |
| `MarketGb` | `MARKET_GB` | `"0"` |

### Anomaly callout
**`MARKET_GB` 추가 query param** — 다른 KSD 메서드에 없는 파라미터. KIS docs 명시. 빈 값 시 `"0"` 기본값 적용.

### Test code

```go
func TestClient_InquireKsdRevSplit(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ksdinfo/rev-split`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "ksd_rev_split_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireKsdRevSplit(context.Background(), domestic.InquireKsdRevSplitParams{
		FromDate: "20260101",
		ToDate:   "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "0", capturedQuery.Get("MARKET_GB")) // default
	assert.Equal(t, "20260101", capturedQuery.Get("F_DT"))

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "005930", res.Output1[0].ShtCd)
	assert.Equal(t, "100", res.Output1[0].InterBfFaceAmt)
	assert.Equal(t, "500", res.Output1[0].InterAfFaceAmt)
}
```

### Commit
```bash
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireKsdRevSplit (예탁원정보 액면변경, HHKDB669105C0)

- KsdRevSplit / KsdRevSplitItem (7 fields, all string)
- Extra MARKET_GB param (default "0") vs other KSD methods
- InquireKsdRevSplitParams (5 query params: SHT_CD/CTS/F_DT/T_DT/MARKET_GB)
- TestClient_InquireKsdRevSplit — fixture ksd_rev_split_success.json

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 8: InquireKsdForfeit

**Files:** Modify `domestic/ksd.go` (append), `domestic/ksd_test.go` (append)

실권주 청약 일정 조회. `output1 []KsdForfeitItem` 9 fields.

- [ ] Step 1: append struct + method to `domestic/ksd.go`
- [ ] Step 2: append test to `domestic/ksd_test.go`
- [ ] Step 3: `go test ./domestic/... -run TestClient_InquireKsdForfeit -v` — PASS
- [ ] Step 4: commit

### 메서드 매핑
- Path: `/uapi/domestic-stock/v1/ksdinfo/forfeit`
- TR_ID: `HHKDB669109C0`
- Output key: `output1 []KsdForfeitItem`

### 응답 struct 필드 (모두 string)

| Go field | json tag | 설명 |
|---|---|---|
| `RecordDate` | `record_date` | 기준일 |
| `ShtCd` | `sht_cd` | 종목코드 |
| `IsinName` | `isin_name` | 종목명 |
| `SubscrDt` | `subscr_dt` | 청약일 |
| `SubscrPrice` | `subscr_price` | 공모가 |
| `SubscrStkQty` | `subscr_stk_qty` | 공모주식수 |
| `RefundDt` | `refund_dt` | 환불일 |
| `ListDt` | `list_dt` | 상장/등록일 |
| `LeadMgr` | `lead_mgr` | 주간사 |

### Params struct fields

| Go field | wire name | default |
|---|---|---|
| `Symbol` | `SHT_CD` | `""` (전체) |
| `ToDate` | `T_DT` | required |
| `FromDate` | `F_DT` | required |
| `Cts` | `CTS` | `""` |

### Test code

```go
func TestClient_InquireKsdForfeit(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ksdinfo/forfeit`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "ksd_forfeit_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireKsdForfeit(context.Background(), domestic.InquireKsdForfeitParams{
		FromDate: "20260101",
		ToDate:   "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "20260505", capturedQuery.Get("T_DT"))

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "005930", res.Output1[0].ShtCd)
	assert.Equal(t, "68000", res.Output1[0].SubscrPrice)
	assert.Equal(t, "한국투자증권", res.Output1[0].LeadMgr)
}
```

### Commit
```bash
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireKsdForfeit (예탁원정보 실권주청약, HHKDB669109C0)

- KsdForfeit / KsdForfeitItem (9 fields, all string)
- InquireKsdForfeitParams (4 query params: SHT_CD/T_DT/F_DT/CTS)
- TestClient_InquireKsdForfeit — fixture ksd_forfeit_success.json

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 9: InquireKsdMandDeposit

**Files:** Modify `domestic/ksd.go` (append), `domestic/ksd_test.go` (append)

의무보호예수 일정 조회. `output1 []KsdMandDepositItem` 6 fields. **ANOMALY: `record_date` 없음 — `depo_date` 가 날짜 key.**

- [ ] Step 1: append struct + method to `domestic/ksd.go`
- [ ] Step 2: append test to `domestic/ksd_test.go`
- [ ] Step 3: `go test ./domestic/... -run TestClient_InquireKsdMandDeposit -v` — PASS
- [ ] Step 4: commit

### 메서드 매핑
- Path: `/uapi/domestic-stock/v1/ksdinfo/mand-deposit`
- TR_ID: `HHKDB669110C0`
- Output key: `output1 []KsdMandDepositItem`

### 응답 struct 필드 (모두 string)

| Go field | json tag | 설명 |
|---|---|---|
| `ShtCd` | `sht_cd` | 종목코드 |
| `IsinName` | `isin_name` | 종목명 |
| `StkQty` | `stk_qty` | 주식수 |
| `DepoDate` | `depo_date` | 예치일 (날짜 key) |
| `DepoReason` | `depo_reason` | 사유 |
| `TotIssueQtyPerRate` | `tot_issue_qty_per_rate` | 총발행주식수대비비율(%) |

### Params struct fields

| Go field | wire name | default |
|---|---|---|
| `ToDate` | `T_DT` | required |
| `Symbol` | `SHT_CD` | `""` (전체) |
| `FromDate` | `F_DT` | required |
| `Cts` | `CTS` | `""` |

### Anomaly callout
**`record_date` 없음** — 보호예수는 기준일 개념이 아닌 `depo_date` (예치일) 기반. struct 에 `RecordDate` 필드 없음. 테스트도 `DepoDate` 로 검증.

### Test code

```go
func TestClient_InquireKsdMandDeposit(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ksdinfo/mand-deposit`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "ksd_mand_deposit_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireKsdMandDeposit(context.Background(), domestic.InquireKsdMandDepositParams{
		FromDate: "20260101",
		ToDate:   "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "20260505", capturedQuery.Get("T_DT"))

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "005930", res.Output1[0].ShtCd)
	assert.Equal(t, "20260101", res.Output1[0].DepoDate) // depo_date (not record_date)
	assert.Equal(t, "의무보호예수", res.Output1[0].DepoReason)
	assert.Equal(t, "0.84", res.Output1[0].TotIssueQtyPerRate)
}
```

### Commit
```bash
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireKsdMandDeposit (예탁원정보 의무보호예수, HHKDB669110C0)

- KsdMandDeposit / KsdMandDepositItem (6 fields, all string)
- No record_date — depo_date used as date key
- InquireKsdMandDepositParams (4 query params: T_DT/SHT_CD/F_DT/CTS)
- TestClient_InquireKsdMandDeposit — fixture ksd_mand_deposit_success.json

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 10: InquireKsdCapDcrs

**Files:** Modify `domestic/ksd.go` (append), `domestic/ksd_test.go` (append)

감자 일정 조회. `output1 []KsdCapDcrsItem` 9 fields.

- [ ] Step 1: append struct + method to `domestic/ksd.go`
- [ ] Step 2: append test to `domestic/ksd_test.go`
- [ ] Step 3: `go test ./domestic/... -run TestClient_InquireKsdCapDcrs -v` — PASS
- [ ] Step 4: commit

### 메서드 매핑
- Path: `/uapi/domestic-stock/v1/ksdinfo/cap-dcrs`
- TR_ID: `HHKDB669106C0`
- Output key: `output1 []KsdCapDcrsItem`

### 응답 struct 필드 (모두 string)

| Go field | json tag | 설명 |
|---|---|---|
| `RecordDate` | `record_date` | 기준일 |
| `ShtCd` | `sht_cd` | 종목코드 |
| `IsinName` | `isin_name` | 종목명 |
| `StkKind` | `stk_kind` | 주식종류 |
| `ReduceCapType` | `reduce_cap_type` | 감자구분 |
| `ReduceCapRate` | `reduce_cap_rate` | 감자배정율 |
| `CompWay` | `comp_way` | 계산방법 |
| `TdStopDt` | `td_stop_dt` | 매매거래정지기간 |
| `ListDt` | `list_dt` | 상장/등록일 |

### Params struct fields

| Go field | wire name | default |
|---|---|---|
| `Cts` | `CTS` | `""` |
| `FromDate` | `F_DT` | required |
| `ToDate` | `T_DT` | required |
| `Symbol` | `SHT_CD` | `""` (전체) |

### Test code

```go
func TestClient_InquireKsdCapDcrs(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ksdinfo/cap-dcrs`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "ksd_cap_dcrs_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireKsdCapDcrs(context.Background(), domestic.InquireKsdCapDcrsParams{
		FromDate: "20260101",
		ToDate:   "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "20260101", capturedQuery.Get("F_DT"))

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "005930", res.Output1[0].ShtCd)
	assert.Equal(t, "유상감자", res.Output1[0].ReduceCapType)
	assert.Equal(t, "주식병합", res.Output1[0].CompWay)
}
```

### Commit
```bash
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireKsdCapDcrs (예탁원정보 감자, HHKDB669106C0)

- KsdCapDcrs / KsdCapDcrsItem (9 fields, all string)
- InquireKsdCapDcrsParams (4 query params: CTS/F_DT/T_DT/SHT_CD)
- TestClient_InquireKsdCapDcrs — fixture ksd_cap_dcrs_success.json

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 11: InquireKsdPurreq

**Files:** Modify `domestic/ksd.go` (append), `domestic/ksd_test.go` (append)

주식매수청구 일정 조회. `output1 []KsdPurreqItem` 9 fields.

- [ ] Step 1: append struct + method to `domestic/ksd.go`
- [ ] Step 2: append test to `domestic/ksd_test.go`
- [ ] Step 3: `go test ./domestic/... -run TestClient_InquireKsdPurreq -v` — PASS
- [ ] Step 4: commit

### 메서드 매핑
- Path: `/uapi/domestic-stock/v1/ksdinfo/purreq`
- TR_ID: `HHKDB669103C0`
- Output key: `output1 []KsdPurreqItem`

### 응답 struct 필드 (모두 string)

| Go field | json tag | 설명 |
|---|---|---|
| `RecordDate` | `record_date` | 기준일 |
| `ShtCd` | `sht_cd` | 종목코드 |
| `IsinName` | `isin_name` | 종목명 |
| `StkKind` | `stk_kind` | 주식종류 |
| `OppOpiRcptTerm` | `opp_opi_rcpt_term` | 반대의사접수시한 |
| `BuyReqRcptTerm` | `buy_req_rcpt_term` | 매수청구접수시한 |
| `BuyReqPrice` | `buy_req_price` | 매수청구가격 |
| `BuyAmtPayDt` | `buy_amt_pay_dt` | 매수대금지급일 |
| `GetMeetDt` | `get_meet_dt` | 주총일 |

### Params struct fields

| Go field | wire name | default |
|---|---|---|
| `Symbol` | `SHT_CD` | `""` (전체) |
| `ToDate` | `T_DT` | required |
| `FromDate` | `F_DT` | required |
| `Cts` | `CTS` | `""` |

### Test code

```go
func TestClient_InquireKsdPurreq(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ksdinfo/purreq`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "ksd_purreq_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireKsdPurreq(context.Background(), domestic.InquireKsdPurreqParams{
		FromDate: "20260101",
		ToDate:   "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "20260505", capturedQuery.Get("T_DT"))

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "005930", res.Output1[0].ShtCd)
	assert.Equal(t, "69000", res.Output1[0].BuyReqPrice)
	assert.Equal(t, "20260326", res.Output1[0].GetMeetDt)
}
```

### Commit
```bash
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireKsdPurreq (예탁원정보 주식매수청구, HHKDB669103C0)

- KsdPurreq / KsdPurreqItem (9 fields, all string)
- InquireKsdPurreqParams (4 query params: SHT_CD/T_DT/F_DT/CTS)
- TestClient_InquireKsdPurreq — fixture ksd_purreq_success.json

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 12: InquireKsdListInfo

**Files:** Modify `domestic/ksd.go` (append), `domestic/ksd_test.go` (append)

주식상장정보 조회. `output1 []KsdListInfoItem` 8 fields. **ANOMALY: leading field `list_dt` (not `record_date`).**

- [ ] Step 1: append struct + method to `domestic/ksd.go`
- [ ] Step 2: append test to `domestic/ksd_test.go`
- [ ] Step 3: `go test ./domestic/... -run TestClient_InquireKsdListInfo -v` — PASS
- [ ] Step 4: commit

### 메서드 매핑
- Path: `/uapi/domestic-stock/v1/ksdinfo/list-info`
- TR_ID: `HHKDB669107C0`
- Output key: `output1 []KsdListInfoItem`

### 응답 struct 필드 (모두 string)

| Go field | json tag | 설명 |
|---|---|---|
| `ListDt` | `list_dt` | 상장/등록일 (날짜 key — not record_date) |
| `ShtCd` | `sht_cd` | 종목코드 |
| `IsinName` | `isin_name` | 종목명 |
| `StkKind` | `stk_kind` | 주식종류 |
| `IssueType` | `issue_type` | 사유 |
| `IssueStkQty` | `issue_stk_qty` | 상장주식수 |
| `TotIssueStkQty` | `tot_issue_stk_qty` | 총발행주식수 |
| `IssuePrice` | `issue_price` | 발행가 |

### Params struct fields

| Go field | wire name | default |
|---|---|---|
| `Symbol` | `SHT_CD` | `""` (전체) |
| `ToDate` | `T_DT` | required |
| `FromDate` | `F_DT` | required |
| `Cts` | `CTS` | `""` |

### Anomaly callout
**Leading field `list_dt` (not `record_date`)** — `KsdListInfoItem` struct 에 `RecordDate` 없음. `list_dt` 가 날짜 기준 필드. 테스트도 `ListDt` 로 검증.

### Test code

```go
func TestClient_InquireKsdListInfo(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ksdinfo/list-info`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "ksd_list_info_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireKsdListInfo(context.Background(), domestic.InquireKsdListInfoParams{
		FromDate: "20260101",
		ToDate:   "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "20260101", capturedQuery.Get("F_DT"))

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "005930", res.Output1[0].ShtCd)
	assert.Equal(t, "20260102", res.Output1[0].ListDt) // list_dt (not record_date)
	assert.Equal(t, "유상증자", res.Output1[0].IssueType)
	assert.Equal(t, "68000", res.Output1[0].IssuePrice)
}
```

### Commit
```bash
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireKsdListInfo (예탁원정보 주식상장정보, HHKDB669107C0)

- KsdListInfo / KsdListInfoItem (8 fields, all string)
- Leading date field is list_dt (not record_date)
- InquireKsdListInfoParams (4 query params: SHT_CD/T_DT/F_DT/CTS)
- TestClient_InquireKsdListInfo — fixture ksd_list_info_success.json

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 13: `examples/domestic_ksd/main.go`

11-call example demonstrating all KSD methods. F_DT="20260101", T_DT="20260505". SHT_CD="005930".

- [ ] Step 1: create `examples/domestic_ksd/main.go`
- [ ] Step 2: `go build ./examples/domestic_ksd/` — success
- [ ] Step 3: commit

```go
// examples/domestic_ksd/main.go — 예탁원 정보 11 메서드 사용 예시
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kenshin579/korea-investment-stock"
	"github.com/kenshin579/korea-investment-stock/domestic"
)

func main() {
	c, err := kis.NewClient(
		kis.WithAppKey(os.Getenv("KIS_APP_KEY")),
		kis.WithAppSecret(os.Getenv("KIS_APP_SECRET")),
	)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	fromDate := "20260101"
	toDate := "20260505"
	symbol := "005930"

	// 1. 배당일정
	div, err := c.Domestic.InquireKsdDividend(ctx, domestic.InquireKsdDividendParams{
		FromDate: fromDate, ToDate: toDate,
	})
	if err != nil {
		log.Printf("InquireKsdDividend: %v", err)
	} else {
		fmt.Printf("InquireKsdDividend: %d rows\n", len(div.Output1))
		for i, item := range div.Output1 {
			if i >= 3 { break }
			fmt.Printf("  [%d] %s %s 배당금=%s\n", i, item.RecordDate, item.IsinName, item.PerStoDiviAmt)
		}
	}

	// 2. 무상증자
	bonus, err := c.Domestic.InquireKsdBonusIssue(ctx, domestic.InquireKsdBonusIssueParams{
		FromDate: fromDate, ToDate: toDate,
	})
	if err != nil {
		log.Printf("InquireKsdBonusIssue: %v", err)
	} else {
		fmt.Printf("InquireKsdBonusIssue: %d rows\n", len(bonus.Output1))
	}

	// 3. 유상증자
	paid, err := c.Domestic.InquireKsdPaidinCapin(ctx, domestic.InquireKsdPaidinCapinParams{
		FromDate: fromDate, ToDate: toDate,
	})
	if err != nil {
		log.Printf("InquireKsdPaidinCapin: %v", err)
	} else {
		fmt.Printf("InquireKsdPaidinCapin: %d rows\n", len(paid.Output)) // output (not output1)
	}

	// 4. 주주총회
	meet, err := c.Domestic.InquireKsdSharehldMeet(ctx, domestic.InquireKsdSharehldMeetParams{
		FromDate: fromDate, ToDate: toDate,
	})
	if err != nil {
		log.Printf("InquireKsdSharehldMeet: %v", err)
	} else {
		fmt.Printf("InquireKsdSharehldMeet: %d rows\n", len(meet.Output1))
	}

	// 5. 합병/분할
	merge, err := c.Domestic.InquireKsdMergerSplit(ctx, domestic.InquireKsdMergerSplitParams{
		FromDate: fromDate, ToDate: toDate,
	})
	if err != nil {
		log.Printf("InquireKsdMergerSplit: %v", err)
	} else {
		fmt.Printf("InquireKsdMergerSplit: %d rows\n", len(merge.Output1))
		for i, item := range merge.Output1 {
			if i >= 3 { break }
			fmt.Printf("  [%d] %s → %s (%s)\n", i, item.OppCustNm, item.CustNm, item.MergeType)
		}
	}

	// 6. 액면변경
	rev, err := c.Domestic.InquireKsdRevSplit(ctx, domestic.InquireKsdRevSplitParams{
		FromDate: fromDate, ToDate: toDate,
	})
	if err != nil {
		log.Printf("InquireKsdRevSplit: %v", err)
	} else {
		fmt.Printf("InquireKsdRevSplit: %d rows\n", len(rev.Output1))
	}

	// 7. 실권주청약
	forf, err := c.Domestic.InquireKsdForfeit(ctx, domestic.InquireKsdForfeitParams{
		Symbol: symbol, FromDate: fromDate, ToDate: toDate,
	})
	if err != nil {
		log.Printf("InquireKsdForfeit: %v", err)
	} else {
		fmt.Printf("InquireKsdForfeit: %d rows\n", len(forf.Output1))
	}

	// 8. 의무보호예수
	dep, err := c.Domestic.InquireKsdMandDeposit(ctx, domestic.InquireKsdMandDepositParams{
		Symbol: symbol, FromDate: fromDate, ToDate: toDate,
	})
	if err != nil {
		log.Printf("InquireKsdMandDeposit: %v", err)
	} else {
		fmt.Printf("InquireKsdMandDeposit: %d rows\n", len(dep.Output1))
		for i, item := range dep.Output1 {
			if i >= 3 { break }
			fmt.Printf("  [%d] depo_date=%s qty=%s reason=%s\n", i, item.DepoDate, item.StkQty, item.DepoReason)
		}
	}

	// 9. 감자
	cap, err := c.Domestic.InquireKsdCapDcrs(ctx, domestic.InquireKsdCapDcrsParams{
		FromDate: fromDate, ToDate: toDate,
	})
	if err != nil {
		log.Printf("InquireKsdCapDcrs: %v", err)
	} else {
		fmt.Printf("InquireKsdCapDcrs: %d rows\n", len(cap.Output1))
	}

	// 10. 주식매수청구
	pur, err := c.Domestic.InquireKsdPurreq(ctx, domestic.InquireKsdPurreqParams{
		FromDate: fromDate, ToDate: toDate,
	})
	if err != nil {
		log.Printf("InquireKsdPurreq: %v", err)
	} else {
		fmt.Printf("InquireKsdPurreq: %d rows\n", len(pur.Output1))
	}

	// 11. 주식상장정보
	lst, err := c.Domestic.InquireKsdListInfo(ctx, domestic.InquireKsdListInfoParams{
		Symbol: symbol, FromDate: fromDate, ToDate: toDate,
	})
	if err != nil {
		log.Printf("InquireKsdListInfo: %v", err)
	} else {
		fmt.Printf("InquireKsdListInfo: %d rows\n", len(lst.Output1))
		for i, item := range lst.Output1 {
			if i >= 3 { break }
			fmt.Printf("  [%d] list_dt=%s %s 발행가=%s\n", i, item.ListDt, item.IsinName, item.IssuePrice)
		}
	}
}
```

Commit:
```bash
git commit -m "$(cat <<'EOF'
[feat] examples — domestic_ksd/main.go (Phase 2.4 KSD 11 메서드 사용 예시)

examples/domestic_ksd/main.go:
- 11 KSD 메서드 순차 호출 예시
- 주요 필드 출력 (배당금, 합병사명, depo_date, list_dt 등 anomaly 포함)
- F_DT/T_DT "20260101"/"20260505", SHT_CD="005930"

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 14: 문서 갱신

- [ ] Step 1: `CLAUDE.md` — banner Phase 2.3 → Phase 2.4, plan link 추가
- [ ] Step 2: `README.md` — Available Methods 표에 11행 추가, heading `Phase 1.2 ~ 2.4`, count 42 → 53
- [ ] Step 3: `CHANGELOG.md` — `[1.7.0]` entry above `[1.6.0]`
- [ ] Step 4: `domestic/doc.go` — Phase 2.4 section 추가 (Phase 2.2 section 뒤에)
- [ ] Step 5: commit

### CLAUDE.md 변경사항

Banner 행 교체:
```
# 변경 전
> **Phase 2.3 — 해외주식 추가 Ranking (v1.6.0).** Phase 2.4+ 는 추후 sub-plan 으로.

# 변경 후
> **Phase 2.4 — 예탁원 정보 확장 11 메서드 (v1.7.0).** Phase 2.5+ 는 추후 sub-plan 으로.
```

Plan link bullet 추가 (Phase 2.3 link 뒤에):
```
- Phase 2.4 implementation plan: [`docs/superpowers/specs/2026-05-05-phase2-4-ksd-implementation-plan.md`](docs/superpowers/specs/2026-05-05-phase2-4-ksd-implementation-plan.md)
```

### CHANGELOG.md `[1.7.0]` entry

```markdown
## [1.7.0] - 2026-05-05

### Added
- `domestic.InquireKsdDividend` — 예탁원정보 배당일정 (HHKDB669102C0)
- `domestic.InquireKsdBonusIssue` — 예탁원정보 무상증자 (HHKDB669101C0)
- `domestic.InquireKsdPaidinCapin` — 예탁원정보 유상증자 (HHKDB669100C0)
- `domestic.InquireKsdSharehldMeet` — 예탁원정보 주주총회 (HHKDB669111C0)
- `domestic.InquireKsdMergerSplit` — 예탁원정보 합병/분할 (HHKDB669104C0)
- `domestic.InquireKsdRevSplit` — 예탁원정보 액면변경 (HHKDB669105C0)
- `domestic.InquireKsdForfeit` — 예탁원정보 실권주청약 (HHKDB669109C0)
- `domestic.InquireKsdMandDeposit` — 예탁원정보 의무보호예수 (HHKDB669110C0)
- `domestic.InquireKsdCapDcrs` — 예탁원정보 감자 (HHKDB669106C0)
- `domestic.InquireKsdPurreq` — 예탁원정보 주식매수청구 (HHKDB669103C0)
- `domestic.InquireKsdListInfo` — 예탁원정보 주식상장정보 (HHKDB669107C0)
- `examples/domestic_ksd/main.go` — KSD 11 메서드 통합 예시

### Notes
- KSD 모든 응답 필드는 KIS docs 명시 String — Go plain `string` (decimal/int64 변환 미적용)
- `InquireKsdPaidinCapin`: output key `output` (not `output1`) — KIS API 응답 구조 그대로
- `InquireKsdMergerSplit`: `isin_name` 없음; `opp_cust_cd`/`opp_cust_nm` + `cust_cd`/`cust_nm` pair
- `InquireKsdRevSplit`: extra `MARKET_GB` query param (default "0")
- `InquireKsdMandDeposit`: `record_date` 없음; `depo_date` 가 날짜 key
- `InquireKsdListInfo`: leading date field `list_dt` (not `record_date`)
- Total methods: 42 → 53
```

### domestic/doc.go Phase 2.4 section

Phase 2.2 section 뒤, `사용자는 root kis.Client` 행 전에 삽입:

```go
// Phase 2.4 메서드 (11):
//
//   - InquireKsdDividend    — 예탁원정보 배당일정 (HHKDB669102C0)
//   - InquireKsdBonusIssue  — 예탁원정보 무상증자 (HHKDB669101C0)
//   - InquireKsdPaidinCapin — 예탁원정보 유상증자 (HHKDB669100C0) [output key: output]
//   - InquireKsdSharehldMeet — 예탁원정보 주주총회 (HHKDB669111C0)
//   - InquireKsdMergerSplit  — 예탁원정보 합병/분할 (HHKDB669104C0) [no isin_name]
//   - InquireKsdRevSplit     — 예탁원정보 액면변경 (HHKDB669105C0) [+MARKET_GB]
//   - InquireKsdForfeit      — 예탁원정보 실권주청약 (HHKDB669109C0)
//   - InquireKsdMandDeposit  — 예탁원정보 의무보호예수 (HHKDB669110C0) [depo_date]
//   - InquireKsdCapDcrs      — 예탁원정보 감자 (HHKDB669106C0)
//   - InquireKsdPurreq       — 예탁원정보 주식매수청구 (HHKDB669103C0)
//   - InquireKsdListInfo     — 예탁원정보 주식상장정보 (HHKDB669107C0) [list_dt]
```

### Commit
```bash
git commit -m "$(cat <<'EOF'
[docs] Phase 2.4 문서 갱신 (v1.7.0)

- CLAUDE.md: banner Phase 2.4, plan link 추가
- README.md: Available Methods 42 → 53 (heading Phase 1.2 ~ 2.4)
- CHANGELOG.md: [1.7.0] entry (11 메서드 + KSD all-string convention notes)
- domestic/doc.go: Phase 2.4 section (11 메서드, anomaly 주석 포함)

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 15: 최종 점검

- [ ] Step 1: format check
  ```bash
  gofmt -l domestic/ksd.go domestic/ksd_test.go examples/domestic_ksd/main.go
  # Expected: no output (no format issues)
  ```

- [ ] Step 2: build
  ```bash
  go build ./...
  # Expected: success (no output)
  ```

- [ ] Step 3: vet
  ```bash
  go vet ./...
  # Expected: success (no output)
  ```

- [ ] Step 4: test with race detector
  ```bash
  go test -race ./domestic/... -v -count=1 2>&1 | tail -20
  # Expected: PASS, no DATA RACE
  ```

- [ ] Step 5: coverage
  ```bash
  go test ./domestic/... -coverprofile=coverage.out
  go tool cover -func=coverage.out | grep -E "^(total|domestic)"
  # Expected: domestic >= 80%, total >= 80%
  ```

- [ ] Step 6: file count check
  ```bash
  ls domestic/ksd.go domestic/ksd_test.go \
     domestic/testdata/ksd_dividend_success.json \
     domestic/testdata/ksd_bonus_issue_success.json \
     domestic/testdata/ksd_paidin_capin_success.json \
     domestic/testdata/ksd_sharehld_meet_success.json \
     domestic/testdata/ksd_merger_split_success.json \
     domestic/testdata/ksd_rev_split_success.json \
     domestic/testdata/ksd_forfeit_success.json \
     domestic/testdata/ksd_mand_deposit_success.json \
     domestic/testdata/ksd_cap_dcrs_success.json \
     domestic/testdata/ksd_purreq_success.json \
     domestic/testdata/ksd_list_info_success.json \
     examples/domestic_ksd/main.go | wc -l
  # Expected: 14
  ```

- [ ] Step 7: commit count ahead of main
  ```bash
  git log main..HEAD --oneline | wc -l
  # Expected: 14 (Task 1 step 13 + Tasks 2-12 each + Task 13 + Task 14)
  ```

---

## Task 16: PR 생성 (사용자 승인 후)

> **사용자 승인 후** 아래 단계를 실행한다.

- [ ] Step 1: push feature branch
  ```bash
  git push -u origin feat/phase2-4-ksd
  ```

- [ ] Step 2: create PR
  ```bash
  gh pr create --title "feat: Phase 2.4 — 예탁원 정보 11 메서드 추가 (v1.7.0)" --body "$(cat <<'EOF'
  ## Summary
  
  예탁원(KSD) 정보 11 메서드를 `domestic/ksd.go` 에 추가합니다.
  
  - `InquireKsdDividend` — 배당일정 (HHKDB669102C0)
  - `InquireKsdBonusIssue` — 무상증자 (HHKDB669101C0)
  - `InquireKsdPaidinCapin` — 유상증자 (HHKDB669100C0)
  - `InquireKsdSharehldMeet` — 주주총회 (HHKDB669111C0)
  - `InquireKsdMergerSplit` — 합병/분할 (HHKDB669104C0)
  - `InquireKsdRevSplit` — 액면변경 (HHKDB669105C0)
  - `InquireKsdForfeit` — 실권주청약 (HHKDB669109C0)
  - `InquireKsdMandDeposit` — 의무보호예수 (HHKDB669110C0)
  - `InquireKsdCapDcrs` — 감자 (HHKDB669106C0)
  - `InquireKsdPurreq` — 주식매수청구 (HHKDB669103C0)
  - `InquireKsdListInfo` — 주식상장정보 (HHKDB669107C0)
  
  총 메서드 수: 42 → **53**
  
  ## KSD 특이사항 (4 anomalies)
  
  | 메서드 | 특이사항 |
  |---|---|
  | `InquireKsdPaidinCapin` | output key `output` (not `output1`) |
  | `InquireKsdMergerSplit` | `isin_name` 없음; `opp_cust_*/cust_*` pair |
  | `InquireKsdRevSplit` | extra `MARKET_GB` query param (default "0") |
  | `InquireKsdMandDeposit` | `record_date` 없음; `depo_date` 가 날짜 key |
  | `InquireKsdListInfo` | leading date field `list_dt` (not `record_date`) |
  
  ## 타입 규칙
  
  KSD 모든 응답 필드 → plain `string` (KIS docs String 명시 — decimal/int64 변환 미적용).
  
  ## Test plan
  
  - [ ] `go test ./domestic/... -race` — PASS, no DATA RACE
  - [ ] `go test ./domestic/... -cover` — domestic >= 80%
  - [ ] `go build ./...` — success
  - [ ] `go vet ./...` — no issues
  - [ ] 11 testdata JSON fixtures valid
  - [ ] `examples/domestic_ksd/main.go` builds
  - [ ] CHANGELOG.md, README.md, CLAUDE.md, domestic/doc.go 갱신 확인
  
  🤖 Generated with [Claude Code](https://claude.com/claude-code)
  EOF
  )"
  ```

- [ ] Step 3: merge PR (after review)

- [ ] Step 4: tag and release
  ```bash
  git checkout main && git pull origin main
  git tag v1.7.0
  git push origin v1.7.0
  gh release create v1.7.0 --title "v1.7.0 — KSD 예탁원 정보 11 메서드 (Phase 2.4)" --notes "$(cat <<'EOF'
  ## What's new
  
  예탁원(KSD) 정보 11 메서드 추가. 총 메서드 수 42 → **53**.
  
  자세한 내용: CHANGELOG.md `[1.7.0]` 참조.
  
  ### 주요 anomaly
  - `InquireKsdPaidinCapin`: output key `output` (not `output1`)
  - `InquireKsdMergerSplit`: `isin_name` 없음; `opp_cust_*/cust_*` pair
  - `InquireKsdRevSplit`: extra `MARKET_GB` query param
  - `InquireKsdMandDeposit`: `record_date` 없음; `depo_date` 사용
  - `InquireKsdListInfo`: leading date `list_dt` (not `record_date`)
  EOF
  )"
  ```

---

## Self-review checklist

- [x] 0 placeholders (TBD/TODO/Similar) — 없음
- [x] 11 메서드 모두 field table 포함 (Tasks 2-12)
- [x] Task 4 `output` (not `output1`) 명시 — anomaly callout + struct comment + test assertion
- [x] Task 9 `no record_date`, `depo_date` 명시 — anomaly callout + test assertion
- [x] Task 6 `no isin_name` + `opp_cust_*/cust_*` pair 명시 — anomaly callout + test assertion
- [x] Task 12 `list_dt` (not `record_date`) 명시 — anomaly callout + test assertion
- [x] Task 7 `MARKET_GB` extra param 명시 — anomaly callout + test assertion
- [x] HEREDOC commits 모두 적용 (Tasks 1-14)
- [x] Task 16 "사용자 승인 후" 명시
- [x] 모든 fixture 2 records (005930/000660), string values
- [x] Reviewer: kenshin579 (global CLAUDE.md policy)
