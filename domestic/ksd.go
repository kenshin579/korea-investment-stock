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
