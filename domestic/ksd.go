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
	RecordDate    string `json:"record_date"`      // 기준일
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

// KsdBonusIssue 는 예탁원정보(무상증자) (HHKDB669101C0) 응답.
//
// 한투 docs: docs/api/국내주식/예탁원정보(무상증자).md
// path: /uapi/domestic-stock/v1/ksdinfo/bonus-issue
type KsdBonusIssue struct {
	Output1 []KsdBonusIssueItem `json:"output1"`
}

// KsdBonusIssueItem 은 무상증자 한 행. 모든 필드 string (KIS docs).
type KsdBonusIssueItem struct {
	RecordDate     string `json:"record_date"`       // 기준일
	ShtCd          string `json:"sht_cd"`            // 종목코드
	IsinName       string `json:"isin_name"`         // 종목명
	FixRate        string `json:"fix_rate"`          // 확정배정율
	OddRecPrice    string `json:"odd_rec_price"`     // 단주기준가
	RightDt        string `json:"right_dt"`          // 권리락일
	OddPayDt       string `json:"odd_pay_dt"`        // 단주대금지급일
	ListDate       string `json:"list_date"`         // 상장/등록일
	TotIssueStkQty string `json:"tot_issue_stk_qty"` // 발행주식
	IssueStkQty    string `json:"issue_stk_qty"`     // 발행할주식
	StkKind        string `json:"stk_kind"`          // 주식종류
}

// InquireKsdBonusIssueParams 는 무상증자 조회 파라미터.
type InquireKsdBonusIssueParams struct {
	Cts      string // CTS — 공백 입력 (default "")
	FromDate string // F_DT — 조회시작일 YYYYMMDD
	ToDate   string // T_DT — 조회종료일 YYYYMMDD
	Symbol   string // SHT_CD — 종목코드 (공백=전체)
}

// InquireKsdBonusIssue 호출.
//
// 한투 docs: docs/api/국내주식/예탁원정보(무상증자).md
// path: /uapi/domestic-stock/v1/ksdinfo/bonus-issue (HHKDB669101C0)
func (c *Client) InquireKsdBonusIssue(ctx context.Context, params InquireKsdBonusIssueParams) (*KsdBonusIssue, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ksdinfo/bonus-issue",
		TrID:   "HHKDB669101C0",
		Query: map[string]string{
			"CTS":    params.Cts,
			"F_DT":   params.FromDate,
			"T_DT":   params.ToDate,
			"SHT_CD": params.Symbol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res KsdBonusIssue
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse KsdBonusIssue: %w", err)
	}
	return &res, nil
}

// KsdPaidinCapin 는 예탁원정보(유상증자) (HHKDB669100C0) 응답.
//
// 한투 docs: docs/api/국내주식/예탁원정보(유상증자).md
// path: /uapi/domestic-stock/v1/ksdinfo/paidin-capin
// ANOMALY: output key 는 "output" (다른 KSD 메서드와 달리 "output1" 아님).
type KsdPaidinCapin struct {
	Output []KsdPaidinCapinItem `json:"output"`
}

// KsdPaidinCapinItem 은 유상증자 한 행. 모든 필드 string (KIS docs).
type KsdPaidinCapinItem struct {
	RecordDate     string `json:"record_date"`       // 기준일
	ShtCd          string `json:"sht_cd"`            // 종목코드
	IsinName       string `json:"isin_name"`         // 종목명
	TotIssueStkQty string `json:"tot_issue_stk_qty"` // 발행주식
	IssueStkQty    string `json:"issue_stk_qty"`     // 발행할주식
	FixRate        string `json:"fix_rate"`          // 확정배정율
	DiscRate       string `json:"disc_rate"`         // 할인율
	FixPrice       string `json:"fix_price"`         // 발행예정가
	RightDt        string `json:"right_dt"`          // 권리락일
	SubTermFt      string `json:"sub_term_ft"`       // 청약기간(시작)
	SubTerm        string `json:"sub_term"`          // 청약기간(종료)
	ListDate       string `json:"list_date"`         // 상장/등록일
	StkKind        string `json:"stk_kind"`          // 주식종류
}

// InquireKsdPaidinCapinParams 는 유상증자 조회 파라미터.
type InquireKsdPaidinCapinParams struct {
	Cts      string // CTS — 공백 입력 (default "")
	Gb1      string // GB1 — 1:청약일별, 2:기준일별. 빈 값=>"1"
	FromDate string // F_DT — 조회시작일 YYYYMMDD
	ToDate   string // T_DT — 조회종료일 YYYYMMDD
	Symbol   string // SHT_CD — 종목코드 (공백=전체)
}

// InquireKsdPaidinCapin 호출.
//
// 한투 docs: docs/api/국내주식/예탁원정보(유상증자).md
// path: /uapi/domestic-stock/v1/ksdinfo/paidin-capin (HHKDB669100C0)
func (c *Client) InquireKsdPaidinCapin(ctx context.Context, params InquireKsdPaidinCapinParams) (*KsdPaidinCapin, error) {
	gb1 := params.Gb1
	if gb1 == "" {
		gb1 = "1"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ksdinfo/paidin-capin",
		TrID:   "HHKDB669100C0",
		Query: map[string]string{
			"CTS":    params.Cts,
			"GB1":    gb1,
			"F_DT":   params.FromDate,
			"T_DT":   params.ToDate,
			"SHT_CD": params.Symbol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res KsdPaidinCapin
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse KsdPaidinCapin: %w", err)
	}
	return &res, nil
}
