package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// 5 financial 메서드 모두 동일 query (Symbol + 분기/연도) — 공통 helper.

// inquireFinanceQuery 는 finance category 의 공통 query 생성.
//
// 모든 finance API 가 fid_cond_mrkt_div_code="J" 고정, FID_DIV_CLS_CODE 는 0(년)/1(분기).
func inquireFinanceQuery(symbol string, quarter bool) map[string]string {
	div := "0"
	if quarter {
		div = "1"
	}
	return map[string]string{
		"FID_DIV_CLS_CODE":       div,
		"fid_cond_mrkt_div_code": "J",
		"fid_input_iscd":         symbol,
	}
}

// FinancialRatio 는 재무비율 (FHKST66430300) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_재무비율.md
// path: /uapi/domestic-stock/v1/finance/financial-ratio
type FinancialRatio struct {
	Output []FinancialRatioItem `json:"output"`
}

// FinancialRatioItem 은 재무비율 응답의 한 행 (분기/년).
type FinancialRatioItem struct {
	StacYymm     string          `json:"stac_yymm"`             // 결산 년월 (YYYYMM)
	Grs          float64         `json:"grs,string"`            // 매출액 증가율
	BsopPrfiInrt float64         `json:"bsop_prfi_inrt,string"` // 영업 이익 증가율 (적자지속/흑자전환/적자전환은 0)
	NtinInrt     float64         `json:"ntin_inrt,string"`      // 순이익 증가율
	RoeVal       float64         `json:"roe_val,string"`        // ROE 값
	Eps          decimal.Decimal `json:"eps"`                   // EPS
	Sps          decimal.Decimal `json:"sps"`                   // 주당매출액
	Bps          decimal.Decimal `json:"bps"`                   // BPS
	RsrvRate     float64         `json:"rsrv_rate,string"`      // 유보 비율
	LbltRate     float64         `json:"lblt_rate,string"`      // 부채 비율
}

// InquireFinancialRatioParams 는 재무비율 조회 파라미터.
//
// 5 financial 메서드 공통: Symbol (필수) + Quarter (false=년 default, true=분기).
type InquireFinancialRatioParams struct {
	Symbol  string // fid_input_iscd (필수, 종목코드)
	Quarter bool   // FID_DIV_CLS_CODE — false=>"0"(년 default), true=>"1"(분기)
}

// InquireFinancialRatio 는 재무비율 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_재무비율.md
// path: /uapi/domestic-stock/v1/finance/financial-ratio (FHKST66430300)
func (c *Client) InquireFinancialRatio(ctx context.Context, params InquireFinancialRatioParams) (*FinancialRatio, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-stock/v1/finance/financial-ratio",
		TrID:     "FHKST66430300",
		Query:    inquireFinanceQuery(params.Symbol, params.Quarter),
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res FinancialRatio
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse FinancialRatio: %w", err)
	}
	return &res, nil
}
