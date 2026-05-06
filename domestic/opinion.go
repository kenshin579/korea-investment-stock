package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
	"github.com/shopspring/decimal"
)

// InvestOpinion 은 종목투자의견 (FHKST663300C0) 응답.
//
// 한투 docs: docs/api/국내주식/종목투자의견.md
// path: /uapi/domestic-stock/v1/quotations/invest-opinion
type InvestOpinion struct {
	Output []InvestOpinionItem `json:"output"`
}

// InvestOpinionItem 은 응답의 output 한 행 (12 fields).
type InvestOpinionItem struct {
	StckBsopDate        string          `json:"stck_bsop_date"`          // 영업 일자
	InvtOpnn            string          `json:"invt_opnn"`               // 투자의견
	InvtOpnnClsCode     string          `json:"invt_opnn_cls_code"`      // 투자의견 구분 코드
	RgbfInvtOpnn        string          `json:"rgbf_invt_opnn"`          // 직전 투자의견
	RgbfInvtOpnnClsCode string          `json:"rgbf_invt_opnn_cls_code"` // 직전 투자의견 구분 코드
	MbcrName            string          `json:"mbcr_name"`               // 회원사명
	HtsGoalPrc          decimal.Decimal `json:"hts_goal_prc"`            // HTS 목표 가격
	StckPrdyClpr        decimal.Decimal `json:"stck_prdy_clpr"`          // 주식 전일 종가
	StckNdayEsdg        decimal.Decimal `json:"stck_nday_esdg"`          // 주식 N일 추정 단가
	NdayDprt            float64         `json:"nday_dprt,string"`        // N일 이격도
	StftEsdg            decimal.Decimal `json:"stft_esdg"`               // 직전 추정 단가
	Dprt                float64         `json:"dprt,string"`             // 이격도
}

// InquireInvestOpinionParams 는 종목투자의견 조회 파라미터.
type InquireInvestOpinionParams struct {
	MarketCode     string // FID_COND_MRKT_DIV_CODE — 빈 값=>"J"
	CondScrDivCode string // FID_COND_SCR_DIV_CODE — 빈 값=>"16633"
	Symbol         string // FID_INPUT_ISCD — 필수, 단축 종목코드 (예 "005930")
	StartDate      string // FID_INPUT_DATE_1 — 조회 시작일 YYYYMMDD
	EndDate        string // FID_INPUT_DATE_2 — 조회 종료일 YYYYMMDD
}

// InquireInvestOpinion 은 종목투자의견 호출.
//
// 한투 docs: docs/api/국내주식/종목투자의견.md
// path: /uapi/domestic-stock/v1/quotations/invest-opinion (FHKST663300C0)
func (c *Client) InquireInvestOpinion(ctx context.Context, params InquireInvestOpinionParams) (*InvestOpinion, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "16633"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/invest-opinion",
		TrID:   "FHKST663300C0",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_COND_SCR_DIV_CODE":  scrDiv,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_INPUT_DATE_1":       params.StartDate,
			"FID_INPUT_DATE_2":       params.EndDate,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res InvestOpinion
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InvestOpinion: %w", err)
	}
	return &res, nil
}
