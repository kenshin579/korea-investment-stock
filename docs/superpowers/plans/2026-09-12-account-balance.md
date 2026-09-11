# 계좌 잔고 조회 (국내·해외) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** `client.Domestic.InquireBalanceAll` / `client.Overseas.InquireBalanceAll` 로 설정된 계좌의 보유 종목 전체를 받는다(조회 전용, 주문 없음). v1.32.0 릴리스.

**Architecture:** 기존 sub-client 패턴(`domestic/`, `overseas/` 에 `balance.go` 추가, fixture 기반 httpmock 테스트) 그대로. 단, 잔고 API 는 (a) 계좌번호 쿼리(`CANO`/`ACNT_PRDT_CD`)와 (b) 연속조회 헤더(`tr_cont`)가 필요한데 `internal/httpclient` 가 둘 다 지원하지 않으므로 Task 1 에서 먼저 넓힌다. 숫자 필드는 빈 문자열·부호에 강한 `kistypes.Int`/`kistypes.Float` 만 쓴다(decimal 은 `""` 에서 실패).

**Tech Stack:** Go 1.25 · resty(내부) · `kistypes` · testify · jarcoal/httpmock

**배경:** moneyflow 포트폴리오 설계 `moneyflow.advenoh.pe.kr/docs/superpowers/specs/2026-09-12-portfolio-sync-design.md` §3·§8 의 PR 1.

**주의 (PUBLIC 저장소):** fixture·테스트·예제·커밋 메시지에 실제 계좌번호·보유 종목·금액을 넣지 않는다. 계좌번호는 `00000000-00`/`12345678-01` 같은 합성 값만.

---

## 파일 구조

| 파일 | 역할 |
|---|---|
| `internal/httpclient/client.go` (수정) | `Request.TrCont` → `tr_cont` 요청 헤더, `Response.TrCont` ← 응답 헤더, `Client.Account()` 계좌번호 분리, `HasNext()` |
| `internal/httpclient/client_test.go` (수정) | 위 3가지 테스트 |
| `domestic/balance.go` (신규) | `Balance`/`BalanceItem`/`BalanceSummary`, `InquireBalance`(1페이지), `InquireBalanceAll`(연속조회) |
| `domestic/balance_test.go` (신규) | 정상·빈 계좌·연속조회 테스트 |
| `domestic/testdata/inquire_balance_{success,empty,page1}.json` (신규) | 합성 fixture |
| `overseas/balance.go` (신규) | `OverseasBalance`/…, `InquireBalance`, `InquireBalanceAll` |
| `overseas/balance_test.go` (신규) | 정상·필수 파라미터·연속조회 테스트 |
| `overseas/testdata/inquire_balance_{success,page1}.json` (신규) | 합성 fixture |
| `examples/account_balance/main.go` (신규) | 국내+해외 잔고 출력 예제 |
| `balance_integration_test.go` (신규, 루트) | `-tags integration` 실호출 |
| `domestic/doc.go`, `overseas/doc.go`, `README.md`, `CHANGELOG.md`, `CLAUDE.md`, `*/testdata/README.md` (수정) | 문서 |

---

## 사전 준비

- [ ] **브랜치 생성**

```bash
cd /Users/user/src/workspace_moneyflow/korea-investment-stock
git checkout main && git pull origin main
git checkout -b feature/account-balance
go build ./... && go test ./... 2>&1 | tail -3
```
Expected: 모두 `ok`.

---

### Task 1: httpclient — tr_cont 헤더·계좌번호 분리

**Files:**
- Modify: `internal/httpclient/client.go`
- Test: `internal/httpclient/client_test.go`

- [ ] **Step 1: 실패하는 테스트 작성** — `internal/httpclient/client_test.go` 끝에 추가

```go
func TestClient_Do_TrCont(t *testing.T) {
	tm := &stubTokenMgr{bearer: "Bearer t"}
	c := newTestClient(t, tm)
	httpmock.RegisterResponder(http.MethodGet, "=~/inquire-balance",
		func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "N", req.Header.Get("tr_cont"), "연속조회 헤더 전달")
			resp := httpmock.NewStringResponse(200, `{"rt_cd":"0","msg_cd":"OK","msg1":"ok","output1":[]}`)
			resp.Header.Set("tr_cont", "M")
			return resp, nil
		})

	resp, err := c.Do(context.Background(), &Request{
		Method: http.MethodGet, Path: "/inquire-balance", TrID: "TTTC8434R", TrCont: "N",
	})
	require.NoError(t, err)
	assert.Equal(t, "M", resp.TrCont, "응답 헤더 tr_cont 캡처")
}

func TestClient_Do_TrContEmptyNotSent(t *testing.T) {
	tm := &stubTokenMgr{bearer: "Bearer t"}
	c := newTestClient(t, tm)
	httpmock.RegisterResponder(http.MethodGet, "=~/inquire-balance",
		func(req *http.Request) (*http.Response, error) {
			_, has := req.Header["Tr_cont"]
			assert.False(t, has, "초기 조회는 tr_cont 헤더를 보내지 않는다")
			return httpmock.NewStringResponse(200, `{"rt_cd":"0","msg_cd":"OK","msg1":"ok"}`), nil
		})
	resp, err := c.Do(context.Background(), &Request{Method: http.MethodGet, Path: "/inquire-balance", TrID: "T"})
	require.NoError(t, err)
	assert.Equal(t, "", resp.TrCont)
}

func TestSplitAccountNo(t *testing.T) {
	cases := []struct {
		in, cano, prdt string
		wantErr        bool
	}{
		{"12345678-01", "12345678", "01", false},
		{"1234567801", "12345678", "01", false},
		{" 12345678-01 ", "12345678", "01", false},
		{"1234-01", "", "", true},
		{"12345678-1", "", "", true},
		{"", "", "", true},
	}
	for _, tc := range cases {
		cano, prdt, err := SplitAccountNo(tc.in)
		if tc.wantErr {
			assert.Error(t, err, tc.in)
			continue
		}
		require.NoError(t, err, tc.in)
		assert.Equal(t, tc.cano, cano)
		assert.Equal(t, tc.prdt, prdt)
	}
}

func TestClient_Account(t *testing.T) {
	c := newTestClient(t, &stubTokenMgr{bearer: "b"})
	cano, prdt, err := c.Account()
	require.NoError(t, err)
	assert.Equal(t, "12345678", cano)
	assert.Equal(t, "01", prdt)
}

func TestHasNext(t *testing.T) {
	assert.True(t, HasNext("F"))
	assert.True(t, HasNext("M"))
	assert.False(t, HasNext("D"))
	assert.False(t, HasNext("E"))
	assert.False(t, HasNext(""))
}
```

- [ ] **Step 2: 실패 확인**

Run: `go test ./internal/httpclient/ -run 'TrCont|SplitAccountNo|Account|HasNext' 2>&1 | head -5`
Expected: `undefined: SplitAccountNo`, `HasNext`, `resp.TrCont` 컴파일 에러.

- [ ] **Step 3: 구현** — `internal/httpclient/client.go`

`Request` 에 필드 추가(`CustType` 아래):
```go
	// TrCont: 연속조회 헤더 tr_cont. 다음 페이지 조회 시 "N". 빈 문자열이면 미전송(초기 조회).
	TrCont string
```

`Response` 에 필드 추가(`Raw` 아래):
```go
	// TrCont 는 응답 헤더 tr_cont. F/M: 다음 데이터 있음, D/E: 마지막. HasNext 로 판정.
	TrCont string `json:"-"`
```

`send()` 안, `SetHeader("Content-Type", ...)` 다음의 `if req.CustType != "" {...}` 뒤에:
```go
		if req.TrCont != "" {
			r.SetHeader("tr_cont", req.TrCont)
		}
```

`send()` 안 두 반환 지점 모두에 헤더 캡처 추가. `raw := httpResp.Body()` 바로 아래에 `trCont := httpResp.Header().Get("tr_cont")` 를 두고:
```go
		if jsonErr == nil && resp.RtCode != "" {
			resp.Raw = raw
			resp.TrCont = trCont
			return &resp, nil
		}
		// ... (중략: 5xx 재시도 블록 그대로)
		resp.Raw = raw
		resp.TrCont = trCont
		return &resp, nil
```

파일 끝에 추가:
```go
// HasNext 는 응답 헤더 tr_cont 가 "다음 데이터 있음"(F/M) 인지 판정한다. D/E/빈 값은 마지막 페이지.
func HasNext(trCont string) bool { return trCont == "F" || trCont == "M" }

// Account 는 설정된 계좌번호를 CANO(종합계좌번호 8자리)와 ACNT_PRDT_CD(계좌상품코드 2자리)로 나눈다.
// 계좌·주문 계열 API 의 쿼리 파라미터에 쓴다.
func (c *Client) Account() (cano, prdtCd string, err error) {
	return SplitAccountNo(c.cfg.AccountNo)
}

// SplitAccountNo 는 "12345678-01" 또는 "1234567801" 형태의 계좌번호를 8-2 로 나눈다.
func SplitAccountNo(accountNo string) (cano, prdtCd string, err error) {
	s := strings.TrimSpace(accountNo)
	if i := strings.IndexByte(s, '-'); i >= 0 {
		cano, prdtCd = s[:i], s[i+1:]
	} else if len(s) == 10 {
		cano, prdtCd = s[:8], s[8:]
	}
	if len(cano) != 8 || len(prdtCd) != 2 {
		return "", "", fmt.Errorf("kis: account number must be 8-2 form like 12345678-01, got %q", accountNo)
	}
	return cano, prdtCd, nil
}
```

- [ ] **Step 4: 통과 확인**

Run: `go test ./internal/httpclient/ 2>&1 | tail -2`
Expected: `ok  	github.com/kenshin579/korea-investment-stock/internal/httpclient`

- [ ] **Step 5: 커밋**

```bash
git add internal/httpclient/client.go internal/httpclient/client_test.go
git commit -m "$(cat <<'EOF'
feat(httpclient): tr_cont 연속조회 헤더 + 계좌번호 CANO/ACNT_PRDT_CD 분리

잔고 조회 계열 API 의 전제. Request.TrCont 로 다음 페이지를 요청하고 Response.TrCont
(응답 헤더)로 더 있는지 판정한다. Client.Account() 는 설정된 8-2 계좌번호를 나눈다.

Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>
EOF
)"
```

---

### Task 2: domestic.InquireBalance — 1페이지 조회

**Files:**
- Create: `domestic/balance.go`
- Create: `domestic/testdata/inquire_balance_success.json`, `domestic/testdata/inquire_balance_empty.json`
- Test: `domestic/balance_test.go`

- [ ] **Step 1: fixture 작성** — `domestic/testdata/inquire_balance_success.json` (합성 값)

```json
{
  "rt_cd": "0",
  "msg_cd": "KIOK0510",
  "msg1": "조회가 완료되었습니다",
  "ctx_area_fk100": "",
  "ctx_area_nk100": "",
  "output1": [
    {"pdno": "005930", "prdt_name": "삼성전자", "trad_dvsn_name": "현금", "bfdy_buy_qty": "0", "bfdy_sll_qty": "0", "thdt_buyqty": "0", "thdt_sll_qty": "0", "hldg_qty": "10", "ord_psbl_qty": "10", "pchs_avg_pric": "70000.0000", "pchs_amt": "700000", "prpr": "75000", "evlu_amt": "750000", "evlu_pfls_amt": "50000", "evlu_pfls_rt": "7.14", "evlu_erng_rt": "0.00000000", "loan_dt": "", "loan_amt": "0", "stln_slng_chgs": "0", "expd_dt": "", "fltt_rt": "1.35000000", "bfdy_cprs_icdc": "1000", "item_mgna_rt_name": "40%", "grta_rt_name": "", "sbst_pric": "52500", "stck_loan_unpr": "0.0000"},
    {"pdno": "379780", "prdt_name": "RISE 미국S&P500", "trad_dvsn_name": "현금", "bfdy_buy_qty": "0", "bfdy_sll_qty": "0", "thdt_buyqty": "0", "thdt_sll_qty": "0", "hldg_qty": "3", "ord_psbl_qty": "3", "pchs_avg_pric": "20000.0000", "pchs_amt": "60000", "prpr": "22000", "evlu_amt": "66000", "evlu_pfls_amt": "6000", "evlu_pfls_rt": "10.00", "evlu_erng_rt": "0.00000000", "loan_dt": "", "loan_amt": "0", "stln_slng_chgs": "0", "expd_dt": "", "fltt_rt": "-0.45000000", "bfdy_cprs_icdc": "-100", "item_mgna_rt_name": "20%", "grta_rt_name": "", "sbst_pric": "15400", "stck_loan_unpr": "0.0000"}
  ],
  "output2": [
    {"dnca_tot_amt": "100000", "nxdy_excc_amt": "100000", "prvs_rcdl_excc_amt": "100000", "cma_evlu_amt": "0", "bfdy_buy_amt": "0", "thdt_buy_amt": "0", "nxdy_auto_rdpt_amt": "0", "bfdy_sll_amt": "0", "thdt_sll_amt": "0", "d2_auto_rdpt_amt": "0", "bfdy_tlex_amt": "0", "thdt_tlex_amt": "0", "tot_loan_amt": "0", "scts_evlu_amt": "816000", "tot_evlu_amt": "916000", "nass_amt": "916000", "fncg_gld_auto_rdpt_yn": "", "pchs_amt_smtl_amt": "760000", "evlu_amt_smtl_amt": "816000", "evlu_pfls_smtl_amt": "56000", "tot_stln_slng_chgs": "0", "bfdy_tot_asst_evlu_amt": "910000", "asst_icdc_amt": "6000", "asst_icdc_erng_rt": "0.00000000"}
  ]
}
```

`domestic/testdata/inquire_balance_empty.json` — 보유 없는 계좌. 숫자 필드가 `""` 로 오는 경우를 흉내 낸다:
```json
{
  "rt_cd": "0",
  "msg_cd": "KIOK0510",
  "msg1": "조회가 완료되었습니다",
  "ctx_area_fk100": "",
  "ctx_area_nk100": "",
  "output1": [],
  "output2": [
    {"dnca_tot_amt": "100000", "nxdy_excc_amt": "100000", "prvs_rcdl_excc_amt": "100000", "cma_evlu_amt": "", "bfdy_buy_amt": "0", "thdt_buy_amt": "0", "nxdy_auto_rdpt_amt": "0", "bfdy_sll_amt": "0", "thdt_sll_amt": "0", "d2_auto_rdpt_amt": "0", "bfdy_tlex_amt": "0", "thdt_tlex_amt": "0", "tot_loan_amt": "0", "scts_evlu_amt": "0", "tot_evlu_amt": "100000", "nass_amt": "100000", "fncg_gld_auto_rdpt_yn": "", "pchs_amt_smtl_amt": "0", "evlu_amt_smtl_amt": "0", "evlu_pfls_smtl_amt": "0", "tot_stln_slng_chgs": "0", "bfdy_tot_asst_evlu_amt": "100000", "asst_icdc_amt": "0", "asst_icdc_erng_rt": ""}
  ]
}
```

- [ ] **Step 2: 실패하는 테스트 작성** — `domestic/balance_test.go`

```go
package domestic_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/domestic"
)

func TestClient_InquireBalance(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, `=~/trading/inquire-balance`,
		func(req *http.Request) (*http.Response, error) {
			q := req.URL.Query()
			assert.Equal(t, "00000000", q.Get("CANO"), "계좌번호 앞 8자리")
			assert.Equal(t, "00", q.Get("ACNT_PRDT_CD"), "계좌상품코드")
			assert.Equal(t, "N", q.Get("AFHR_FLPR_YN"), "기본값")
			assert.Equal(t, "02", q.Get("INQR_DVSN"), "기본값 종목별")
			assert.Equal(t, "01", q.Get("UNPR_DVSN"))
			assert.Equal(t, "N", q.Get("FUND_STTL_ICLD_YN"))
			assert.Equal(t, "N", q.Get("FNCG_AMT_AUTO_RDPT_YN"))
			assert.Equal(t, "00", q.Get("PRCS_DVSN"))
			assert.Equal(t, "", q.Get("CTX_AREA_FK100"))
			assert.Equal(t, "TTTC8434R", req.Header.Get("tr_id"))
			resp := httpmock.NewStringResponse(200, loadFixtureString(t, "inquire_balance_success.json"))
			resp.Header.Set("tr_cont", "D")
			return resp, nil
		})

	c := newTestClient(t)
	res, err := c.InquireBalance(context.Background(), domestic.InquireBalanceParams{})
	require.NoError(t, err)
	require.NotNil(t, res)

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "005930", res.Output1[0].Pdno)
	assert.Equal(t, "삼성전자", res.Output1[0].PrdtName)
	assert.Equal(t, int64(10), int64(res.Output1[0].HldgQty))
	assert.InDelta(t, 70000.0, float64(res.Output1[0].PchsAvgPric), 0.0001)
	assert.Equal(t, int64(700000), int64(res.Output1[0].PchsAmt))
	assert.Equal(t, int64(750000), int64(res.Output1[0].EvluAmt))
	assert.Equal(t, int64(50000), int64(res.Output1[0].EvluPflsAmt))
	assert.InDelta(t, 7.14, float64(res.Output1[0].EvluPflsRt), 0.001)
	assert.Equal(t, int64(-100), int64(res.Output1[1].BfdyCprsIcdc), "음수 정수")
	assert.InDelta(t, -0.45, float64(res.Output1[1].FlttRt), 0.001)

	require.Len(t, res.Output2, 1)
	assert.Equal(t, int64(100000), int64(res.Output2[0].DncaTotAmt))
	assert.Equal(t, int64(816000), int64(res.Output2[0].EvluAmtSmtlAmt))
	assert.Equal(t, int64(916000), int64(res.Output2[0].TotEvluAmt))
	assert.Equal(t, "D", res.TrCont)
}

func TestClient_InquireBalance_Empty(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodGet, `=~/trading/inquire-balance`,
		httpmock.NewStringResponder(200, loadFixtureString(t, "inquire_balance_empty.json")))

	c := newTestClient(t)
	res, err := c.InquireBalance(context.Background(), domestic.InquireBalanceParams{})
	require.NoError(t, err)
	assert.Empty(t, res.Output1)
	require.Len(t, res.Output2, 1)
	assert.Equal(t, int64(0), int64(res.Output2[0].CmaEvluAmt), `"" → 0`)
	assert.Equal(t, 0.0, float64(res.Output2[0].AsstIcdcErngRt), `"" → 0`)
	assert.Equal(t, "", res.TrCont, "헤더 없으면 빈 값 = 마지막")
}

func TestClient_InquireBalance_ParamsOverride(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodGet, `=~/trading/inquire-balance`,
		func(req *http.Request) (*http.Response, error) {
			q := req.URL.Query()
			assert.Equal(t, "01", q.Get("INQR_DVSN"))
			assert.Equal(t, "X", q.Get("AFHR_FLPR_YN"))
			assert.Equal(t, "FK", q.Get("CTX_AREA_FK100"))
			assert.Equal(t, "NK", q.Get("CTX_AREA_NK100"))
			assert.Equal(t, "N", req.Header.Get("tr_cont"))
			return httpmock.NewStringResponse(200, loadFixtureString(t, "inquire_balance_empty.json")), nil
		})
	c := newTestClient(t)
	_, err := c.InquireBalance(context.Background(), domestic.InquireBalanceParams{
		InqrDvsn: "01", AfhrFlprYn: "X", CtxAreaFk100: "FK", CtxAreaNk100: "NK", TrCont: "N",
	})
	require.NoError(t, err)
}
```

- [ ] **Step 3: 실패 확인**

Run: `go test ./domestic/ -run InquireBalance 2>&1 | head -3`
Expected: `undefined: domestic.InquireBalanceParams` 컴파일 에러. (`InquireBalanceSheet` 와 이름이 다르므로 충돌 없음.)

- [ ] **Step 4: 구현** — `domestic/balance.go`

```go
package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
	"github.com/kenshin579/korea-investment-stock/kistypes"
)

// Balance 는 주식잔고조회 (TTTC8434R) 응답.
//
// 한투 docs: docs/api/국내주식/주식잔고조회.md
// path: /uapi/domestic-stock/v1/trading/inquire-balance
//
// 한 번의 호출에 최대 50건. 더 있으면 TrCont 가 "F"/"M" 이고 CtxAreaFk100/CtxAreaNk100 을
// 다음 호출 파라미터로 넘긴다. 전체를 한 번에 받으려면 InquireBalanceAll 을 쓴다.
// 당일 전량 매도한 종목은 HldgQty 0 으로 남아 있을 수 있다(D-2 이후 사라짐).
// 모의투자(VTTC8434R)는 지원하지 않는다 — 실전 TR 고정.
type Balance struct {
	Output1      []BalanceItem    `json:"output1"`        // 보유 종목
	Output2      []BalanceSummary `json:"output2"`        // 계좌 요약 (1행)
	CtxAreaFk100 string           `json:"ctx_area_fk100"` // 연속조회검색조건100
	CtxAreaNk100 string           `json:"ctx_area_nk100"` // 연속조회키100
	TrCont       string           `json:"-"`              // 응답 헤더 tr_cont. F/M: 다음 있음, D/E/"": 마지막
}

// BalanceItem 은 보유 종목 1건 (output1). 금액은 원 단위 정수, 단가·비율은 실수.
type BalanceItem struct {
	Pdno           string         `json:"pdno"`              // 종목번호 (6자리)
	PrdtName       string         `json:"prdt_name"`         // 종목명
	TradDvsnName   string         `json:"trad_dvsn_name"`    // 매매구분명
	BfdyBuyQty     kistypes.Int   `json:"bfdy_buy_qty"`      // 전일매수수량
	BfdySllQty     kistypes.Int   `json:"bfdy_sll_qty"`      // 전일매도수량
	ThdtBuyqty     kistypes.Int   `json:"thdt_buyqty"`       // 금일매수수량
	ThdtSllQty     kistypes.Int   `json:"thdt_sll_qty"`      // 금일매도수량
	HldgQty        kistypes.Int   `json:"hldg_qty"`          // 보유수량
	OrdPsblQty     kistypes.Int   `json:"ord_psbl_qty"`      // 주문가능수량
	PchsAvgPric    kistypes.Float `json:"pchs_avg_pric"`     // 매입평균가격 (매입금액/보유수량)
	PchsAmt        kistypes.Int   `json:"pchs_amt"`          // 매입금액
	Prpr           kistypes.Int   `json:"prpr"`              // 현재가
	EvluAmt        kistypes.Int   `json:"evlu_amt"`          // 평가금액
	EvluPflsAmt    kistypes.Int   `json:"evlu_pfls_amt"`     // 평가손익금액 (평가금액 - 매입금액)
	EvluPflsRt     kistypes.Float `json:"evlu_pfls_rt"`      // 평가손익율
	EvluErngRt     kistypes.Float `json:"evlu_erng_rt"`      // 평가수익율 (미사용, 0)
	LoanDt         string         `json:"loan_dt"`           // 대출일자 (INQR_DVSN=01 일 때만)
	LoanAmt        kistypes.Int   `json:"loan_amt"`          // 대출금액
	StlnSlngChgs   kistypes.Int   `json:"stln_slng_chgs"`    // 대주매각대금
	ExpdDt         string         `json:"expd_dt"`           // 만기일자
	FlttRt         kistypes.Float `json:"fltt_rt"`           // 등락율
	BfdyCprsIcdc   kistypes.Int   `json:"bfdy_cprs_icdc"`    // 전일대비증감
	ItemMgnaRtName string         `json:"item_mgna_rt_name"` // 종목증거금율명
	GrtaRtName     string         `json:"grta_rt_name"`      // 보증금율명
	SbstPric       kistypes.Int   `json:"sbst_pric"`         // 대용가격
	StckLoanUnpr   kistypes.Float `json:"stck_loan_unpr"`    // 주식대출단가
}

// BalanceSummary 는 계좌 요약 (output2, 1행).
type BalanceSummary struct {
	DncaTotAmt         kistypes.Int   `json:"dnca_tot_amt"`           // 예수금총금액
	NxdyExccAmt        kistypes.Int   `json:"nxdy_excc_amt"`          // 익일정산금액 (D+1 예수금)
	PrvsRcdlExccAmt    kistypes.Int   `json:"prvs_rcdl_excc_amt"`     // 가수도정산금액 (D+2 예수금)
	CmaEvluAmt         kistypes.Int   `json:"cma_evlu_amt"`           // CMA평가금액
	BfdyBuyAmt         kistypes.Int   `json:"bfdy_buy_amt"`           // 전일매수금액
	ThdtBuyAmt         kistypes.Int   `json:"thdt_buy_amt"`           // 금일매수금액
	NxdyAutoRdptAmt    kistypes.Int   `json:"nxdy_auto_rdpt_amt"`     // 익일자동상환금액
	BfdySllAmt         kistypes.Int   `json:"bfdy_sll_amt"`           // 전일매도금액
	ThdtSllAmt         kistypes.Int   `json:"thdt_sll_amt"`           // 금일매도금액
	D2AutoRdptAmt      kistypes.Int   `json:"d2_auto_rdpt_amt"`       // D+2자동상환금액
	BfdyTlexAmt        kistypes.Int   `json:"bfdy_tlex_amt"`          // 전일제비용금액
	ThdtTlexAmt        kistypes.Int   `json:"thdt_tlex_amt"`          // 금일제비용금액
	TotLoanAmt         kistypes.Int   `json:"tot_loan_amt"`           // 총대출금액
	SctsEvluAmt        kistypes.Int   `json:"scts_evlu_amt"`          // 유가평가금액
	TotEvluAmt         kistypes.Int   `json:"tot_evlu_amt"`           // 총평가금액 (유가증권 평가 + D+2 예수금)
	NassAmt            kistypes.Int   `json:"nass_amt"`               // 순자산금액
	FncgGldAutoRdptYn  string         `json:"fncg_gld_auto_rdpt_yn"`  // 융자금자동상환여부
	PchsAmtSmtlAmt     kistypes.Int   `json:"pchs_amt_smtl_amt"`      // 매입금액합계금액
	EvluAmtSmtlAmt     kistypes.Int   `json:"evlu_amt_smtl_amt"`      // 평가금액합계금액 (유가증권)
	EvluPflsSmtlAmt    kistypes.Int   `json:"evlu_pfls_smtl_amt"`     // 평가손익합계금액
	TotStlnSlngChgs    kistypes.Int   `json:"tot_stln_slng_chgs"`     // 총대주매각대금
	BfdyTotAsstEvluAmt kistypes.Int   `json:"bfdy_tot_asst_evlu_amt"` // 전일총자산평가금액
	AsstIcdcAmt        kistypes.Int   `json:"asst_icdc_amt"`          // 자산증감액
	AsstIcdcErngRt     kistypes.Float `json:"asst_icdc_erng_rt"`      // 자산증감수익율 (데이터 미제공)
}

// InquireBalanceParams 는 주식잔고조회 파라미터. 계좌번호(CANO/ACNT_PRDT_CD)는 Client 설정값을 쓴다.
//
// 빈 값은 한투 기본값으로 채운다: AfhrFlprYn "N", InqrDvsn "02"(종목별), UnprDvsn "01",
// FundSttlIcldYn "N", FncgAmtAutoRdptYn "N", PrcsDvsn "00"(전일매매포함).
// 연속조회는 이전 응답의 CtxAreaFk100/CtxAreaNk100 과 TrCont "N" 을 넘긴다.
type InquireBalanceParams struct {
	AfhrFlprYn        string // AFHR_FLPR_YN — N: 기본, Y: 시간외단일가, X: NXT 정규장
	InqrDvsn          string // INQR_DVSN — 01: 대출일별, 02: 종목별
	UnprDvsn          string // UNPR_DVSN — 01: 기본
	FundSttlIcldYn    string // FUND_STTL_ICLD_YN — 펀드결제분 포함 N/Y
	FncgAmtAutoRdptYn string // FNCG_AMT_AUTO_RDPT_YN — 융자금액자동상환 N
	PrcsDvsn          string // PRCS_DVSN — 00: 전일매매포함, 01: 전일매매미포함
	CtxAreaFk100      string // CTX_AREA_FK100 — 연속조회. 첫 조회 빈 값
	CtxAreaNk100      string // CTX_AREA_NK100 — 연속조회. 첫 조회 빈 값
	TrCont            string // tr_cont 헤더 — 연속조회 "N". 첫 조회 빈 값
}

// balanceParamOr 는 빈 파라미터를 한투 기본값으로 채운다.
func balanceParamOr(v, def string) string {
	if v == "" {
		return def
	}
	return v
}

// InquireBalance 는 주식잔고조회 1페이지 호출 (최대 50건).
//
// 한투 docs: docs/api/국내주식/주식잔고조회.md
// path: /uapi/domestic-stock/v1/trading/inquire-balance (TTTC8434R)
func (c *Client) InquireBalance(ctx context.Context, params InquireBalanceParams) (*Balance, error) {
	cano, prdtCd, err := c.http.Account()
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/trading/inquire-balance",
		TrID:   "TTTC8434R",
		TrCont: params.TrCont,
		Query: map[string]string{
			"CANO":                  cano,
			"ACNT_PRDT_CD":          prdtCd,
			"AFHR_FLPR_YN":          balanceParamOr(params.AfhrFlprYn, "N"),
			"OFL_YN":                "",
			"INQR_DVSN":             balanceParamOr(params.InqrDvsn, "02"),
			"UNPR_DVSN":             balanceParamOr(params.UnprDvsn, "01"),
			"FUND_STTL_ICLD_YN":     balanceParamOr(params.FundSttlIcldYn, "N"),
			"FNCG_AMT_AUTO_RDPT_YN": balanceParamOr(params.FncgAmtAutoRdptYn, "N"),
			"PRCS_DVSN":             balanceParamOr(params.PrcsDvsn, "00"),
			"CTX_AREA_FK100":        params.CtxAreaFk100,
			"CTX_AREA_NK100":        params.CtxAreaNk100,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res Balance
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse Balance: %w", err)
	}
	res.TrCont = resp.TrCont
	return &res, nil
}
```

- [ ] **Step 5: 통과 확인**

Run: `go test ./domestic/ -run InquireBalance -v 2>&1 | grep -E '^(---|ok|FAIL)'`
Expected: `--- PASS: TestClient_InquireBalance`, `_Empty`, `_ParamsOverride` 3개 PASS. (`TestClient_InquireBalanceSheet` 도 같이 잡히지만 기존 PASS.)

- [ ] **Step 6: 커밋**

```bash
git add domestic/balance.go domestic/balance_test.go domestic/testdata/inquire_balance_success.json domestic/testdata/inquire_balance_empty.json
git commit -m "$(cat <<'EOF'
feat(domestic): 주식잔고조회 InquireBalance (TTTC8434R)

Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>
EOF
)"
```

---

### Task 3: domestic.InquireBalanceAll — 연속조회

**Files:**
- Modify: `domestic/balance.go`
- Create: `domestic/testdata/inquire_balance_page1.json`
- Test: `domestic/balance_test.go`

- [ ] **Step 1: fixture** — `domestic/testdata/inquire_balance_page1.json` (첫 페이지: 1건 + 커서)

```json
{
  "rt_cd": "0",
  "msg_cd": "KIOK0510",
  "msg1": "조회가 완료되었습니다",
  "ctx_area_fk100": "FK-PAGE1",
  "ctx_area_nk100": "NK-PAGE1",
  "output1": [
    {"pdno": "000660", "prdt_name": "SK하이닉스", "trad_dvsn_name": "현금", "bfdy_buy_qty": "0", "bfdy_sll_qty": "0", "thdt_buyqty": "0", "thdt_sll_qty": "0", "hldg_qty": "2", "ord_psbl_qty": "2", "pchs_avg_pric": "150000.0000", "pchs_amt": "300000", "prpr": "160000", "evlu_amt": "320000", "evlu_pfls_amt": "20000", "evlu_pfls_rt": "6.66", "evlu_erng_rt": "0.00000000", "loan_dt": "", "loan_amt": "0", "stln_slng_chgs": "0", "expd_dt": "", "fltt_rt": "0.50000000", "bfdy_cprs_icdc": "800", "item_mgna_rt_name": "40%", "grta_rt_name": "", "sbst_pric": "112000", "stck_loan_unpr": "0.0000"}
  ],
  "output2": [
    {"dnca_tot_amt": "100000", "nxdy_excc_amt": "100000", "prvs_rcdl_excc_amt": "100000", "cma_evlu_amt": "0", "bfdy_buy_amt": "0", "thdt_buy_amt": "0", "nxdy_auto_rdpt_amt": "0", "bfdy_sll_amt": "0", "thdt_sll_amt": "0", "d2_auto_rdpt_amt": "0", "bfdy_tlex_amt": "0", "thdt_tlex_amt": "0", "tot_loan_amt": "0", "scts_evlu_amt": "1136000", "tot_evlu_amt": "1236000", "nass_amt": "1236000", "fncg_gld_auto_rdpt_yn": "", "pchs_amt_smtl_amt": "1060000", "evlu_amt_smtl_amt": "1136000", "evlu_pfls_smtl_amt": "76000", "tot_stln_slng_chgs": "0", "bfdy_tot_asst_evlu_amt": "1230000", "asst_icdc_amt": "6000", "asst_icdc_erng_rt": "0.00000000"}
  ]
}
```

- [ ] **Step 2: 실패하는 테스트** — `domestic/balance_test.go` 끝에 추가

```go
func TestClient_InquireBalanceAll_FollowsTrCont(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	calls := 0
	httpmock.RegisterResponder(http.MethodGet, `=~/trading/inquire-balance`,
		func(req *http.Request) (*http.Response, error) {
			calls++
			q := req.URL.Query()
			switch calls {
			case 1:
				assert.Equal(t, "", req.Header.Get("tr_cont"))
				assert.Equal(t, "", q.Get("CTX_AREA_NK100"))
				resp := httpmock.NewStringResponse(200, loadFixtureString(t, "inquire_balance_page1.json"))
				resp.Header.Set("tr_cont", "M")
				return resp, nil
			case 2:
				assert.Equal(t, "N", req.Header.Get("tr_cont"), "2페이지는 tr_cont=N")
				assert.Equal(t, "FK-PAGE1", q.Get("CTX_AREA_FK100"), "이전 응답 커서 전달")
				assert.Equal(t, "NK-PAGE1", q.Get("CTX_AREA_NK100"))
				resp := httpmock.NewStringResponse(200, loadFixtureString(t, "inquire_balance_success.json"))
				resp.Header.Set("tr_cont", "D")
				return resp, nil
			}
			t.Fatalf("unexpected call #%d", calls)
			return nil, nil
		})

	c := newTestClient(t)
	// 호출자가 커서를 넣어도 첫 페이지부터 다시 읽는다.
	res, err := c.InquireBalanceAll(context.Background(), domestic.InquireBalanceParams{TrCont: "N", CtxAreaNk100: "junk"})
	require.NoError(t, err)
	assert.Equal(t, 2, calls)
	require.Len(t, res.Output1, 3, "1 + 2 건 이어 붙임")
	assert.Equal(t, "000660", res.Output1[0].Pdno)
	assert.Equal(t, "005930", res.Output1[1].Pdno)
	assert.Equal(t, "379780", res.Output1[2].Pdno)
	require.Len(t, res.Output2, 1, "요약은 마지막 페이지 것")
	assert.Equal(t, int64(916000), int64(res.Output2[0].TotEvluAmt))
	assert.Equal(t, "D", res.TrCont)
}

func TestClient_InquireBalanceAll_SinglePage(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodGet, `=~/trading/inquire-balance`,
		httpmock.NewStringResponder(200, loadFixtureString(t, "inquire_balance_success.json")))

	c := newTestClient(t)
	res, err := c.InquireBalanceAll(context.Background(), domestic.InquireBalanceParams{})
	require.NoError(t, err)
	assert.Len(t, res.Output1, 2)
	assert.Equal(t, 1, httpmock.GetTotalCallCount(), "tr_cont 헤더 없음 = 1페이지로 끝")
}

func TestClient_InquireBalanceAll_PageCap(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodGet, `=~/trading/inquire-balance`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, loadFixtureString(t, "inquire_balance_page1.json"))
			resp.Header.Set("tr_cont", "M") // 영원히 다음 있음
			return resp, nil
		})

	c := newTestClient(t)
	_, err := c.InquireBalanceAll(context.Background(), domestic.InquireBalanceParams{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceeded")
	assert.Equal(t, 100, httpmock.GetTotalCallCount())
}
```

- [ ] **Step 3: 실패 확인**

Run: `go test ./domestic/ -run InquireBalanceAll 2>&1 | head -3`
Expected: `c.InquireBalanceAll undefined` 컴파일 에러.

- [ ] **Step 4: 구현** — `domestic/balance.go` 끝에 추가

```go
// maxBalancePages 는 연속조회 상한. 50건 × 100 = 5,000 종목 — 실계좌에서 도달할 수 없는 값이며
// 한투가 tr_cont 를 잘못 주는 경우의 무한 루프 방지용.
const maxBalancePages = 100

// InquireBalanceAll 은 연속조회(tr_cont)를 따라가며 보유 종목 전체를 모은다.
// params 의 TrCont/CtxArea* 는 무시하고 첫 페이지부터 읽는다.
// Output1 은 모든 페이지를 이어 붙이고, Output2·CtxArea*·TrCont 는 마지막 페이지 값이다.
func (c *Client) InquireBalanceAll(ctx context.Context, params InquireBalanceParams) (*Balance, error) {
	params.TrCont, params.CtxAreaFk100, params.CtxAreaNk100 = "", "", ""
	var all *Balance
	for page := 0; page < maxBalancePages; page++ {
		res, err := c.InquireBalance(ctx, params)
		if err != nil {
			return nil, err
		}
		if all == nil {
			all = res
		} else {
			all.Output1 = append(all.Output1, res.Output1...)
			all.Output2 = res.Output2
			all.CtxAreaFk100, all.CtxAreaNk100, all.TrCont = res.CtxAreaFk100, res.CtxAreaNk100, res.TrCont
		}
		if !httpclient.HasNext(res.TrCont) {
			return all, nil
		}
		params.TrCont, params.CtxAreaFk100, params.CtxAreaNk100 = "N", res.CtxAreaFk100, res.CtxAreaNk100
	}
	return nil, fmt.Errorf("kis: InquireBalanceAll: exceeded %d pages", maxBalancePages)
}
```

- [ ] **Step 5: 통과 확인**

Run: `go test ./domestic/ -run 'InquireBalance' -v 2>&1 | grep -E '^(---|ok|FAIL)'`
Expected: 6개 `--- PASS` (Sheet 포함 7개), `ok`.

- [ ] **Step 6: 커밋**

```bash
git add domestic/balance.go domestic/balance_test.go domestic/testdata/inquire_balance_page1.json
git commit -m "$(cat <<'EOF'
feat(domestic): InquireBalanceAll — tr_cont 연속조회로 보유 종목 전체 수집

Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>
EOF
)"
```

---

### Task 4: overseas.InquireBalance / InquireBalanceAll

**Files:**
- Create: `overseas/balance.go`
- Create: `overseas/testdata/inquire_balance_success.json`, `overseas/testdata/inquire_balance_page1.json`
- Test: `overseas/balance_test.go`

해외 잔고는 output2 가 **배열이 아니라 단일 객체**, 커서는 `CTX_AREA_*200`, 거래소·통화 코드가 필수다. 실전에서 `OVRS_EXCG_CD=NASD` 가 "미국전체"(NAS/NYSE/AMEX 모두)다. 미니스탁 잔고는 이 API 에 안 나온다.

- [ ] **Step 1: fixture** — `overseas/testdata/inquire_balance_success.json`

```json
{
  "rt_cd": "0",
  "msg_cd": "KIOK0510",
  "msg1": "조회가 완료되었습니다",
  "ctx_area_fk200": "",
  "ctx_area_nk200": "",
  "output1": [
    {"cano": "00000000", "acnt_prdt_cd": "00", "prdt_type_cd": "512", "ovrs_pdno": "AAPL", "ovrs_item_name": "APPLE INC", "frcr_evlu_pfls_amt": "+120.50000000", "evlu_pfls_rt": "+6.02", "pchs_avg_pric": "200.00000000", "ovrs_cblc_qty": "10", "ord_psbl_qty": "10", "frcr_pchs_amt1": "2000.00000000", "ovrs_stck_evlu_amt": "2120.50000000", "now_pric2": "212.05000000", "tr_crcy_cd": "USD", "ovrs_excg_cd": "NASD", "loan_type_cd": "00", "loan_dt": "", "expd_dt": ""},
    {"cano": "00000000", "acnt_prdt_cd": "00", "prdt_type_cd": "512", "ovrs_pdno": "O", "ovrs_item_name": "REALTY INCOME CORP", "frcr_evlu_pfls_amt": "-15.00000000", "evlu_pfls_rt": "-2.50", "pchs_avg_pric": "60.00000000", "ovrs_cblc_qty": "10", "ord_psbl_qty": "10", "frcr_pchs_amt1": "600.00000000", "ovrs_stck_evlu_amt": "585.00000000", "now_pric2": "58.50000000", "tr_crcy_cd": "USD", "ovrs_excg_cd": "NYSE", "loan_type_cd": "00", "loan_dt": "", "expd_dt": ""}
  ],
  "output2": {"frcr_pchs_amt1": "2600.00000000", "ovrs_rlzt_pfls_amt": "0.00000000", "ovrs_tot_pfls": "+105.50000000", "rlzt_erng_rt": "0.00000000", "tot_evlu_pfls_amt": "+105.50000000", "tot_pftrt": "+4.05", "frcr_buy_amt_smtl1": "2600.00000000", "ovrs_rlzt_pfls_amt2": "0.00000000", "frcr_buy_amt_smtl2": "2600.00000000"}
}
```

`overseas/testdata/inquire_balance_page1.json`:
```json
{
  "rt_cd": "0",
  "msg_cd": "KIOK0510",
  "msg1": "조회가 완료되었습니다",
  "ctx_area_fk200": "FK-PAGE1",
  "ctx_area_nk200": "NK-PAGE1",
  "output1": [
    {"cano": "00000000", "acnt_prdt_cd": "00", "prdt_type_cd": "512", "ovrs_pdno": "QQQM", "ovrs_item_name": "INVESCO NASDAQ 100 ETF", "frcr_evlu_pfls_amt": "+30.00000000", "evlu_pfls_rt": "+1.50", "pchs_avg_pric": "200.00000000", "ovrs_cblc_qty": "10", "ord_psbl_qty": "10", "frcr_pchs_amt1": "2000.00000000", "ovrs_stck_evlu_amt": "2030.00000000", "now_pric2": "203.00000000", "tr_crcy_cd": "USD", "ovrs_excg_cd": "NASD", "loan_type_cd": "00", "loan_dt": "", "expd_dt": ""}
  ],
  "output2": {"frcr_pchs_amt1": "4600.00000000", "ovrs_rlzt_pfls_amt": "0.00000000", "ovrs_tot_pfls": "+135.50000000", "rlzt_erng_rt": "0.00000000", "tot_evlu_pfls_amt": "+135.50000000", "tot_pftrt": "+2.94", "frcr_buy_amt_smtl1": "4600.00000000", "ovrs_rlzt_pfls_amt2": "0.00000000", "frcr_buy_amt_smtl2": "4600.00000000"}
}
```

- [ ] **Step 2: 실패하는 테스트** — `overseas/balance_test.go`

```go
package overseas_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/overseas"
)

func TestClient_InquireBalance(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, `=~/overseas-stock/v1/trading/inquire-balance`,
		func(req *http.Request) (*http.Response, error) {
			q := req.URL.Query()
			assert.Equal(t, "00000000", q.Get("CANO"))
			assert.Equal(t, "00", q.Get("ACNT_PRDT_CD"))
			assert.Equal(t, "NASD", q.Get("OVRS_EXCG_CD"))
			assert.Equal(t, "USD", q.Get("TR_CRCY_CD"))
			assert.Equal(t, "", q.Get("CTX_AREA_FK200"))
			assert.Equal(t, "TTTS3012R", req.Header.Get("tr_id"))
			resp := httpmock.NewStringResponse(200, loadFixtureString(t, "inquire_balance_success.json"))
			resp.Header.Set("tr_cont", "D")
			return resp, nil
		})

	c := newTestClient(t)
	res, err := c.InquireBalance(context.Background(), overseas.InquireBalanceParams{OvrsExcgCd: "NASD", TrCrcyCd: "USD"})
	require.NoError(t, err)

	require.Len(t, res.Output1, 2)
	a := res.Output1[0]
	assert.Equal(t, "AAPL", a.OvrsPdno)
	assert.Equal(t, "APPLE INC", a.OvrsItemName)
	assert.Equal(t, "NASD", a.OvrsExcgCd)
	assert.Equal(t, "USD", a.TrCrcyCd)
	assert.InDelta(t, 10, float64(a.OvrsCblcQty), 0.0001)
	assert.InDelta(t, 200.0, float64(a.PchsAvgPric), 0.0001)
	assert.InDelta(t, 2120.5, float64(a.OvrsStckEvluAmt), 0.0001)
	assert.InDelta(t, 120.5, float64(a.FrcrEvluPflsAmt), 0.0001, "부호 + 문자열")
	assert.InDelta(t, 6.02, float64(a.EvluPflsRt), 0.001)
	assert.InDelta(t, -15.0, float64(res.Output1[1].FrcrEvluPflsAmt), 0.0001, "부호 - 문자열")

	assert.InDelta(t, 2600.0, float64(res.Output2.FrcrPchsAmt1), 0.0001)
	assert.InDelta(t, 105.5, float64(res.Output2.TotEvluPflsAmt), 0.0001)
	assert.InDelta(t, 4.05, float64(res.Output2.TotPftrt), 0.001)
	assert.Equal(t, "D", res.TrCont)
}

func TestClient_InquireBalance_RequiredParams(t *testing.T) {
	c := newTestClient(t)
	_, err := c.InquireBalance(context.Background(), overseas.InquireBalanceParams{TrCrcyCd: "USD"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "OvrsExcgCd")
	_, err = c.InquireBalance(context.Background(), overseas.InquireBalanceParams{OvrsExcgCd: "NASD"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "TrCrcyCd")
}

func TestClient_InquireBalanceAll_FollowsTrCont(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	calls := 0
	httpmock.RegisterResponder(http.MethodGet, `=~/overseas-stock/v1/trading/inquire-balance`,
		func(req *http.Request) (*http.Response, error) {
			calls++
			q := req.URL.Query()
			switch calls {
			case 1:
				assert.Equal(t, "", req.Header.Get("tr_cont"))
				resp := httpmock.NewStringResponse(200, loadFixtureString(t, "inquire_balance_page1.json"))
				resp.Header.Set("tr_cont", "F")
				return resp, nil
			case 2:
				assert.Equal(t, "N", req.Header.Get("tr_cont"))
				assert.Equal(t, "FK-PAGE1", q.Get("CTX_AREA_FK200"))
				assert.Equal(t, "NK-PAGE1", q.Get("CTX_AREA_NK200"))
				resp := httpmock.NewStringResponse(200, loadFixtureString(t, "inquire_balance_success.json"))
				resp.Header.Set("tr_cont", "E")
				return resp, nil
			}
			t.Fatalf("unexpected call #%d", calls)
			return nil, nil
		})

	c := newTestClient(t)
	res, err := c.InquireBalanceAll(context.Background(), overseas.InquireBalanceParams{OvrsExcgCd: "NASD", TrCrcyCd: "USD"})
	require.NoError(t, err)
	assert.Equal(t, 2, calls)
	require.Len(t, res.Output1, 3)
	assert.Equal(t, "QQQM", res.Output1[0].OvrsPdno)
	assert.Equal(t, "AAPL", res.Output1[1].OvrsPdno)
	assert.InDelta(t, 2600.0, float64(res.Output2.FrcrPchsAmt1), 0.0001, "요약은 마지막 페이지")
	assert.Equal(t, "E", res.TrCont)
}

func TestClient_InquireBalanceAll_PageCap(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodGet, `=~/overseas-stock/v1/trading/inquire-balance`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, loadFixtureString(t, "inquire_balance_page1.json"))
			resp.Header.Set("tr_cont", "M")
			return resp, nil
		})
	c := newTestClient(t)
	_, err := c.InquireBalanceAll(context.Background(), overseas.InquireBalanceParams{OvrsExcgCd: "NASD", TrCrcyCd: "USD"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceeded")
	assert.Equal(t, 100, httpmock.GetTotalCallCount())
}
```

- [ ] **Step 3: 실패 확인**

Run: `go test ./overseas/ -run InquireBalance 2>&1 | head -3`
Expected: `undefined: overseas.InquireBalanceParams` 컴파일 에러.

- [ ] **Step 4: 구현** — `overseas/balance.go`

```go
package overseas

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
	"github.com/kenshin579/korea-investment-stock/kistypes"
)

// OverseasBalance 는 해외주식 잔고 (TTTS3012R) 응답.
//
// 한투 docs: docs/api/해외주식/해외주식_잔고.md
// path: /uapi/overseas-stock/v1/trading/inquire-balance
//
// 한 번의 호출에 최대 100건. 더 있으면 TrCont 가 "F"/"M" 이고 CtxAreaFk200/CtxAreaNk200 을
// 다음 호출 파라미터로 넘긴다. 전체를 한 번에 받으려면 InquireBalanceAll 을 쓴다.
// 금액·수량은 거래 통화(TrCrcyCd) 기준 외화이며 원화 환산은 없다. 미니스탁 잔고는 포함되지 않는다.
// 모의투자(VTTS3012R)는 지원하지 않는다 — 실전 TR 고정.
type OverseasBalance struct {
	Output1      []OverseasBalanceItem  `json:"output1"`        // 보유 종목
	Output2      OverseasBalanceSummary `json:"output2"`        // 계좌 요약 (단일 객체 — 국내와 다름)
	CtxAreaFk200 string                 `json:"ctx_area_fk200"` // 연속조회검색조건200
	CtxAreaNk200 string                 `json:"ctx_area_nk200"` // 연속조회키200
	TrCont       string                 `json:"-"`              // 응답 헤더 tr_cont. F/M: 다음 있음, D/E/"": 마지막
}

// OverseasBalanceItem 은 보유 종목 1건 (output1). 숫자는 부호(+/-)·빈 문자열을 허용하는 kistypes.Float.
type OverseasBalanceItem struct {
	Cano            string         `json:"cano"`               // 종합계좌번호
	AcntPrdtCd      string         `json:"acnt_prdt_cd"`       // 계좌상품코드
	PrdtTypeCd      string         `json:"prdt_type_cd"`       // 상품유형코드
	OvrsPdno        string         `json:"ovrs_pdno"`          // 해외상품번호 (티커)
	OvrsItemName    string         `json:"ovrs_item_name"`     // 해외종목명
	FrcrEvluPflsAmt kistypes.Float `json:"frcr_evlu_pfls_amt"` // 외화평가손익금액
	EvluPflsRt      kistypes.Float `json:"evlu_pfls_rt"`       // 평가손익율
	PchsAvgPric     kistypes.Float `json:"pchs_avg_pric"`      // 매입평균가격 (외화)
	OvrsCblcQty     kistypes.Float `json:"ovrs_cblc_qty"`      // 해외잔고수량
	OrdPsblQty      kistypes.Float `json:"ord_psbl_qty"`       // 주문가능수량
	FrcrPchsAmt1    kistypes.Float `json:"frcr_pchs_amt1"`     // 외화매입금액
	OvrsStckEvluAmt kistypes.Float `json:"ovrs_stck_evlu_amt"` // 해외주식평가금액 (외화)
	NowPric2        kistypes.Float `json:"now_pric2"`          // 현재가
	TrCrcyCd        string         `json:"tr_crcy_cd"`         // 거래통화코드 USD/HKD/CNY/JPY/VND
	OvrsExcgCd      string         `json:"ovrs_excg_cd"`       // 해외거래소코드 NASD/NYSE/AMEX/SEHK/SHAA/SZAA/TKSE/HASE/VNSE
	LoanTypeCd      string         `json:"loan_type_cd"`       // 대출유형코드 (00: 해당없음)
	LoanDt          string         `json:"loan_dt"`            // 대출일자
	ExpdDt          string         `json:"expd_dt"`            // 만기일자
}

// OverseasBalanceSummary 는 계좌 요약 (output2, 단일 객체).
type OverseasBalanceSummary struct {
	FrcrPchsAmt1     kistypes.Float `json:"frcr_pchs_amt1"`      // 외화매입금액1
	OvrsRlztPflsAmt  kistypes.Float `json:"ovrs_rlzt_pfls_amt"`  // 해외실현손익금액
	OvrsTotPfls      kistypes.Float `json:"ovrs_tot_pfls"`       // 해외총손익
	RlztErngRt       kistypes.Float `json:"rlzt_erng_rt"`        // 실현수익율
	TotEvluPflsAmt   kistypes.Float `json:"tot_evlu_pfls_amt"`   // 총평가손익금액
	TotPftrt         kistypes.Float `json:"tot_pftrt"`           // 총수익률
	FrcrBuyAmtSmtl1  kistypes.Float `json:"frcr_buy_amt_smtl1"`  // 외화매수금액합계1
	OvrsRlztPflsAmt2 kistypes.Float `json:"ovrs_rlzt_pfls_amt2"` // 해외실현손익금액2
	FrcrBuyAmtSmtl2  kistypes.Float `json:"frcr_buy_amt_smtl2"`  // 외화매수금액합계2
}

// InquireBalanceParams 는 해외주식 잔고 파라미터. 계좌번호는 Client 설정값을 쓴다.
// OvrsExcgCd 와 TrCrcyCd 는 필수 — 실전 미국 전체는 ("NASD", "USD").
type InquireBalanceParams struct {
	OvrsExcgCd   string // OVRS_EXCG_CD (필수) — 실전: NASD 미국전체 / NAS 나스닥 / NYSE / AMEX · 공통: SEHK / SHAA / SZAA / TKSE / HASE / VNSE
	TrCrcyCd     string // TR_CRCY_CD (필수) — USD / HKD / CNY / JPY / VND
	CtxAreaFk200 string // CTX_AREA_FK200 — 연속조회. 첫 조회 빈 값
	CtxAreaNk200 string // CTX_AREA_NK200 — 연속조회. 첫 조회 빈 값
	TrCont       string // tr_cont 헤더 — 연속조회 "N". 첫 조회 빈 값
}

// InquireBalance 는 해외주식 잔고 1페이지 호출 (최대 100건).
//
// 한투 docs: docs/api/해외주식/해외주식_잔고.md
// path: /uapi/overseas-stock/v1/trading/inquire-balance (TTTS3012R)
func (c *Client) InquireBalance(ctx context.Context, params InquireBalanceParams) (*OverseasBalance, error) {
	if params.OvrsExcgCd == "" {
		return nil, errors.New("kis: InquireBalance: OvrsExcgCd is required (e.g. NASD)")
	}
	if params.TrCrcyCd == "" {
		return nil, errors.New("kis: InquireBalance: TrCrcyCd is required (e.g. USD)")
	}
	cano, prdtCd, err := c.http.Account()
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-stock/v1/trading/inquire-balance",
		TrID:   "TTTS3012R",
		TrCont: params.TrCont,
		Query: map[string]string{
			"CANO":           cano,
			"ACNT_PRDT_CD":   prdtCd,
			"OVRS_EXCG_CD":   params.OvrsExcgCd,
			"TR_CRCY_CD":     params.TrCrcyCd,
			"CTX_AREA_FK200": params.CtxAreaFk200,
			"CTX_AREA_NK200": params.CtxAreaNk200,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res OverseasBalance
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse OverseasBalance: %w", err)
	}
	res.TrCont = resp.TrCont
	return &res, nil
}

// maxBalancePages 는 연속조회 상한 (100건 × 100 페이지). 잘못된 tr_cont 로 인한 무한 루프 방지.
const maxBalancePages = 100

// InquireBalanceAll 은 연속조회(tr_cont)를 따라가며 보유 종목 전체를 모은다.
// params 의 TrCont/CtxArea* 는 무시하고 첫 페이지부터 읽는다.
// Output1 은 모든 페이지를 이어 붙이고, Output2·CtxArea*·TrCont 는 마지막 페이지 값이다.
func (c *Client) InquireBalanceAll(ctx context.Context, params InquireBalanceParams) (*OverseasBalance, error) {
	params.TrCont, params.CtxAreaFk200, params.CtxAreaNk200 = "", "", ""
	var all *OverseasBalance
	for page := 0; page < maxBalancePages; page++ {
		res, err := c.InquireBalance(ctx, params)
		if err != nil {
			return nil, err
		}
		if all == nil {
			all = res
		} else {
			all.Output1 = append(all.Output1, res.Output1...)
			all.Output2 = res.Output2
			all.CtxAreaFk200, all.CtxAreaNk200, all.TrCont = res.CtxAreaFk200, res.CtxAreaNk200, res.TrCont
		}
		if !httpclient.HasNext(res.TrCont) {
			return all, nil
		}
		params.TrCont, params.CtxAreaFk200, params.CtxAreaNk200 = "N", res.CtxAreaFk200, res.CtxAreaNk200
	}
	return nil, fmt.Errorf("kis: InquireBalanceAll: exceeded %d pages", maxBalancePages)
}
```

- [ ] **Step 5: 통과 확인**

Run: `go test ./overseas/ -run InquireBalance -v 2>&1 | grep -E '^(---|ok|FAIL)'`
Expected: 4개 `--- PASS`, `ok`.

- [ ] **Step 6: 전체 빌드·vet·테스트**

Run: `go build ./... && go vet ./... && go test ./... 2>&1 | grep -v '^ok' ; echo "exit=$?"`
Expected: `ok` 아닌 줄이 없다(`no test files` 는 무시).

- [ ] **Step 7: 커밋**

```bash
git add overseas/balance.go overseas/balance_test.go overseas/testdata/inquire_balance_success.json overseas/testdata/inquire_balance_page1.json
git commit -m "$(cat <<'EOF'
feat(overseas): 해외주식 잔고 InquireBalance / InquireBalanceAll (TTTS3012R)

output2 가 단일 객체, 커서는 CTX_AREA_*200, 거래소·통화 코드 필수(실전 미국 전체 NASD/USD).
숫자는 부호·빈 문자열을 허용하는 kistypes.Float.

Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>
EOF
)"
```

---

### Task 5: 예제 + integration 테스트

**Files:**
- Create: `examples/account_balance/main.go`
- Create: `balance_integration_test.go` (루트, `package kis_test`)

- [ ] **Step 1: 예제** — `examples/account_balance/main.go`

```go
// account_balance example: 설정된 계좌의 국내·해외 보유 종목 전체 (InquireBalanceAll).
//
// Run: KOREA_INVESTMENT_API_KEY / API_SECRET / ACCOUNT_NO 설정 후 go run ./examples/account_balance
package main

import (
	"context"
	"fmt"
	"log"

	kis "github.com/kenshin579/korea-investment-stock"
	"github.com/kenshin579/korea-investment-stock/domestic"
	"github.com/kenshin579/korea-investment-stock/overseas"
)

func main() {
	client, err := kis.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	// 1. 국내주식 잔고 — 계좌번호는 클라이언트 설정값, 나머지는 한투 기본값
	dom, err := client.Domestic.InquireBalanceAll(ctx, domestic.InquireBalanceParams{})
	if err != nil {
		log.Fatalf("Domestic.InquireBalanceAll: %v", err)
	}
	fmt.Printf("국내 보유 %d 종목\n", len(dom.Output1))
	for _, it := range dom.Output1 {
		fmt.Printf("  %s %-20s 수량=%d 평단=%.0f 평가=%d 손익=%d (%.2f%%)\n",
			it.Pdno, it.PrdtName, int64(it.HldgQty), float64(it.PchsAvgPric),
			int64(it.EvluAmt), int64(it.EvluPflsAmt), float64(it.EvluPflsRt))
	}
	if len(dom.Output2) > 0 {
		s := dom.Output2[0]
		fmt.Printf("  예수금=%d 유가평가=%d 총평가=%d\n", int64(s.DncaTotAmt), int64(s.SctsEvluAmt), int64(s.TotEvluAmt))
	}

	// 2. 해외주식 잔고 — 실전 미국 전체(NASD) + USD
	ovs, err := client.Overseas.InquireBalanceAll(ctx, overseas.InquireBalanceParams{OvrsExcgCd: "NASD", TrCrcyCd: "USD"})
	if err != nil {
		log.Fatalf("Overseas.InquireBalanceAll: %v", err)
	}
	fmt.Printf("해외(미국) 보유 %d 종목\n", len(ovs.Output1))
	for _, it := range ovs.Output1 {
		fmt.Printf("  %-6s %-30s %s 수량=%.4f 평단=%.2f 평가=%.2f 손익=%.2f (%.2f%%)\n",
			it.OvrsPdno, it.OvrsItemName, it.OvrsExcgCd, float64(it.OvrsCblcQty), float64(it.PchsAvgPric),
			float64(it.OvrsStckEvluAmt), float64(it.FrcrEvluPflsAmt), float64(it.EvluPflsRt))
	}
	fmt.Printf("  외화매입합계=%.2f 총평가손익=%.2f (%.2f%%)\n",
		float64(ovs.Output2.FrcrPchsAmt1), float64(ovs.Output2.TotEvluPflsAmt), float64(ovs.Output2.TotPftrt))
}
```

- [ ] **Step 2: 예제 빌드 확인**

Run: `go build -o /dev/null ./examples/account_balance && echo BUILD_OK`
Expected: `BUILD_OK`

- [ ] **Step 3: integration 테스트** — 루트 `balance_integration_test.go`. 이 저장소에는 아직 `integration` 태그 테스트가 없다(형제 라이브러리 opendart-go/fmp-go 관례를 따른다). 기본 `go test ./...` 에는 포함되지 않는다.

```go
//go:build integration

package kis_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	kis "github.com/kenshin579/korea-investment-stock"
	"github.com/kenshin579/korea-investment-stock/domestic"
	"github.com/kenshin579/korea-investment-stock/overseas"
)

// TestIntegration_InquireBalance 는 실계좌 잔고 실호출. KOREA_INVESTMENT_* env 필요.
// 실행: go test -tags integration -run TestIntegration_InquireBalance -v .
// 금액·종목은 로그에 남기지 않는다(개수만).
func TestIntegration_InquireBalance(t *testing.T) {
	c, err := kis.NewClientFromEnv()
	if err != nil {
		t.Skipf("skip: %v", err)
	}
	ctx := context.Background()

	dom, err := c.Domestic.InquireBalanceAll(ctx, domestic.InquireBalanceParams{})
	require.NoError(t, err)
	require.Len(t, dom.Output2, 1, "국내 요약 1행")
	t.Logf("domestic: %d holdings, last tr_cont=%q", len(dom.Output1), dom.TrCont)

	ovs, err := c.Overseas.InquireBalanceAll(ctx, overseas.InquireBalanceParams{OvrsExcgCd: "NASD", TrCrcyCd: "USD"})
	require.NoError(t, err)
	t.Logf("overseas(NASD/USD): %d holdings, last tr_cont=%q", len(ovs.Output1), ovs.TrCont)
}
```

- [ ] **Step 4: 실호출 (사용자 env 가 있을 때만)**

Run: `go vet -tags integration . && go test -tags integration -run TestIntegration_InquireBalance -v . 2>&1 | tail -6`
Expected: env 없으면 `--- SKIP`; 있으면 `--- PASS` 와 `domestic: N holdings`. 현재 한투 위탁계좌는 보유 0 이 정상이다(spec §3). 실패하면 응답 필드 타입 불일치가 원인일 가능성이 높으니 에러의 `parse ...` 메시지로 fixture 를 실제 형식에 맞춘다.

- [ ] **Step 5: 커밋**

```bash
git add examples/account_balance/main.go balance_integration_test.go
git commit -m "$(cat <<'EOF'
feat: account_balance 예제 + 잔고 integration 테스트 (-tags integration)

Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>
EOF
)"
```

---

### Task 6: 문서

**Files:**
- Modify: `domestic/doc.go`, `overseas/doc.go`, `domestic/testdata/README.md`, `overseas/testdata/README.md`, `README.md`, `CHANGELOG.md`, `CLAUDE.md`

- [ ] **Step 1: `domestic/doc.go`** — `// 사용자는 root kis.Client 의 Domestic 필드로 접근.` 바로 위에 추가

```go
// 계좌 잔고 (v1.32.0) — 조회 전용, 주문 없음
//
//	InquireBalance     — 주식잔고조회 1페이지 (최대 50건)  TTTC8434R
//	InquireBalanceAll  — tr_cont 연속조회로 전체 수집
//
// Anomalies (잔고):
//
//	계좌번호(CANO/ACNT_PRDT_CD)는 Client 설정값을 httpclient.Account() 로 분리해 자동 주입
//	연속조회는 응답 헤더 tr_cont(F/M 다음 있음, D/E 마지막) + ctx_area_fk100/nk100 커서
//	숫자 필드는 kistypes.Int/Float (빈 문자열·부호 허용) — decimal 미사용
//	모의투자(VTTC8434R) 미지원
//
```

- [ ] **Step 2: `overseas/doc.go`** — `// 사용자는 root kis.Client 의 Overseas 필드로 접근.` 바로 위에 추가

```go
// 계좌 잔고 (v1.32.0) — 조회 전용, 주문 없음
//
//   - InquireBalance     — 해외주식 잔고 1페이지 (최대 100건) (TTTS3012R) — OvrsExcgCd/TrCrcyCd 필수, output2 단일 객체
//   - InquireBalanceAll  — tr_cont + CTX_AREA_*200 연속조회로 전체 수집
//
```

- [ ] **Step 3: testdata README** — `domestic/testdata/README.md` 의 "REST API 응답 (합성 JSON)" 목록 끝에

```markdown
- `inquire_balance_success.json` / `inquire_balance_empty.json` / `inquire_balance_page1.json` — 주식잔고조회 (TTTC8434R) 정상 / 보유 없음(빈 문자열 숫자) / 연속조회 첫 페이지(커서 포함). 계좌·종목·금액 모두 합성.
```

`overseas/testdata/README.md` 끝에(형식은 파일의 기존 목록을 따른다):
```markdown
- `inquire_balance_success.json` / `inquire_balance_page1.json` — 해외주식 잔고 (TTTS3012R) 정상 / 연속조회 첫 페이지. 계좌·종목·금액 모두 합성.
```

- [ ] **Step 4: `README.md`** — `### Futures (국내선물옵션) — Phase 11.1` 섹션 앞에 새 섹션

```markdown
### 계좌 잔고 조회 — v1.32.0 (조회 전용)

계좌번호는 `NewClient(..., accountNo)` / `KOREA_INVESTMENT_ACCOUNT_NO` 의 값을 자동으로 `CANO`/`ACNT_PRDT_CD` 로 나눠 보낸다.

| 메서드 | TR | 설명 |
|---|---|---|
| `Domestic.InquireBalance` | TTTC8434R | 주식잔고조회 1페이지(50건). `Output2[0]` 에 예수금·총평가 |
| `Domestic.InquireBalanceAll` | TTTC8434R | `tr_cont` 연속조회로 전체 |
| `Overseas.InquireBalance` | TTTS3012R | 해외주식 잔고 1페이지(100건). `OvrsExcgCd`(실전 미국 전체 `NASD`)·`TrCrcyCd` 필수, 외화 기준 |
| `Overseas.InquireBalanceAll` | TTTS3012R | 연속조회로 전체 |

```go
dom, _ := client.Domestic.InquireBalanceAll(ctx, domestic.InquireBalanceParams{})
ovs, _ := client.Overseas.InquireBalanceAll(ctx, overseas.InquireBalanceParams{OvrsExcgCd: "NASD", TrCrcyCd: "USD"})
```

예제: `go run ./examples/account_balance`. 실호출 테스트: `go test -tags integration -run TestIntegration_InquireBalance .`
주문·예약주문은 여전히 범위 밖.
```

그리고 `## Scope` 의 `- ❌ 주식 주문/잔고/예약주문 — 본 spec 에서 다루지 않음` 을 다음으로 교체:
```markdown
- ✅ 주식 잔고 조회 (국내·해외, v1.32.0) — 조회 전용
- ❌ 주식 주문/예약주문 — 본 spec 에서 다루지 않음
```

- [ ] **Step 5: `CHANGELOG.md`** — 파일 최상단 `# CHANGELOG` 다음에 (1.28~1.31 항목이 빠져 있지만 이번에 채우지 않는다 — 별건)

```markdown
## [1.32.0] - 2026-09-12

### Added — 계좌 잔고 조회 (조회 전용, 2 REST EP)
- `Domestic.InquireBalance` / `InquireBalanceAll` — 주식잔고조회 (TTTC8434R), 50건 페이지 + `tr_cont` 연속조회
- `Overseas.InquireBalance` / `InquireBalanceAll` — 해외주식 잔고 (TTTS3012R), 100건 페이지, `OvrsExcgCd`/`TrCrcyCd` 필수
- `internal/httpclient`: `Request.TrCont`(요청 헤더) / `Response.TrCont`(응답 헤더) / `Client.Account()` 계좌번호 8-2 분리 / `HasNext()`
- examples: `account_balance`. 루트 `balance_integration_test.go` (`-tags integration`)

### Notes
- 모의투자 TR(VTTC8434R/VTTS3012R) 미지원 — 실전 고정.
- 잔고 숫자 필드는 `kistypes.Int`/`Float` — 빈 문자열·부호 허용. 해외는 외화 기준(원화 환산 없음).
- 누적 152 REST + 36 WS = 188 endpoints.

```

- [ ] **Step 6: `CLAUDE.md`** — 헤더 인용문과 Out of Scope 갱신

`> **Phase 11.7 — 해외선물옵션 실시간 2 WS (v1.26.0). 누적 150 REST + 36 WS = 186 endpoints.**` →
`> **v1.32.0 — 계좌 잔고 조회 2 REST (국내 TTTC8434R · 해외 TTTS3012R). 누적 152 REST + 36 WS = 188 endpoints.** 이전: Phase 11.7 해외선물옵션 실시간 2 WS (v1.26.0).`

`## Out of Scope (Phase 0)` 의 `선물옵션 · 장내채권 · 주문/잔고/예약주문` →
`주문/예약주문 (잔고 조회는 v1.32.0 에서 추가 — `domestic/balance.go`, `overseas/balance.go`)`

`## Common Commands` 블록에 한 줄 추가:
```bash
go test -tags integration -run TestIntegration_InquireBalance -v .   # 실계좌 잔고 실호출 (KOREA_INVESTMENT_* env)
```

- [ ] **Step 7: 검증 후 커밋**

Run: `go vet ./... && go test ./... 2>&1 | grep -c '^ok'`
Expected: 패키지 수만큼 `ok` (숫자 > 0), 실패 없음.

```bash
git add domestic/doc.go overseas/doc.go domestic/testdata/README.md overseas/testdata/README.md README.md CHANGELOG.md CLAUDE.md
git commit -m "$(cat <<'EOF'
docs: 계좌 잔고 조회 v1.32.0 — README/CHANGELOG/CLAUDE.md/doc.go

Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>
EOF
)"
```

---

### Task 7: PR → 릴리스

- [ ] **Step 1: 푸시 + PR** (리뷰어 지정 금지)

```bash
git push -u origin feature/account-balance
gh pr create --title "feat: 계좌 잔고 조회 — 국내 TTTC8434R · 해외 TTTS3012R (v1.32.0)" --body "$(cat <<'EOF'
## Summary
- `Domestic.InquireBalance/InquireBalanceAll` (주식잔고조회, 50건 페이지 + tr_cont 연속조회)
- `Overseas.InquireBalance/InquireBalanceAll` (해외주식 잔고, 100건 페이지, OvrsExcgCd/TrCrcyCd 필수)
- httpclient: `tr_cont` 요청/응답 헤더, `Account()` 계좌번호 8-2 분리, `HasNext()`
- 조회 전용 — 주문 API 없음. 모의투자 TR 미지원.
- 배경: moneyflow 포트폴리오 수집 (moneyflow `docs/superpowers/specs/2026-09-12-portfolio-sync-design.md` PR 1)

## Test plan
- [x] `go build ./... && go vet ./... && go test ./...`
- [x] httpmock fixture 단위 테스트: 정상 / 빈 계좌("" 숫자) / 파라미터 override / 연속조회 2페이지 / 페이지 상한
- [ ] `go test -tags integration -run TestIntegration_InquireBalance -v .` 실계좌 실호출 (보유 0 정상)
- [ ] `go run ./examples/account_balance`

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

- [ ] **Step 2: 머지 후 릴리스** (사용자가 머지한 뒤)

```bash
git checkout main && git pull origin main
./scripts/release.sh v1.32.0
```
Expected: build/vet/test → 모듈 zip 검증 → 태그 push → GitHub Release 생성. 이후 moneyflow `go.mod` 를 `v1.32.0` 으로 bump (계획 2).

---

## Self-review

- **Spec 커버리지 (portfolio spec §3·§8 PR 1):** 국내·해외 잔고 조회 ✅ (Task 2~4), 조회 전용 ✅, v1.32.0 릴리스 ✅ (Task 7), PUBLIC 저장소 fixture 합성 값 ✅, integration 테스트 ✅ (Task 5).
- **타입 일관성:** `httpclient.HasNext`·`Account`·`Request.TrCont`·`Response.TrCont` (Task 1) 를 Task 2~4 가 그대로 사용. `domestic.Balance.Output2` 는 `[]BalanceSummary`(배열), `overseas.OverseasBalance.Output2` 는 `OverseasBalanceSummary`(객체) — 예제·테스트가 이 차이를 따른다. `maxBalancePages` 는 두 패키지에 각각 unexported 로 선언(패키지가 달라 충돌 없음).
- **알려진 불확실성:** 실제 KIS 응답에서 국내 금액 필드가 `"123.00"` 같은 소수 형식이면 `kistypes.Int` 파싱이 실패한다. Task 5 Step 4 실호출에서 확인하고, 그 경우 해당 필드를 `kistypes.Float` 로 바꾼다(테스트 fixture 도 함께).
