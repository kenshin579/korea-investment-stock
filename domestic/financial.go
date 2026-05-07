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

// IncomeStatement 는 손익계산서 (FHKST66430200) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_손익계산서.md
// path: /uapi/domestic-stock/v1/finance/income-statement
//
// 분기 데이터는 연단위 누적 합산.
type IncomeStatement struct {
	Output []IncomeStatementItem `json:"output"`
}

// IncomeStatementItem 은 손익계산서 응답의 한 행.
type IncomeStatementItem struct {
	StacYymm     string `json:"stac_yymm"`             // 결산 년월
	SaleAccount  int64  `json:"sale_account,string"`   // 매출액
	SaleCost     int64  `json:"sale_cost,string"`      // 매출 원가
	SaleTotlPrfi int64  `json:"sale_totl_prfi,string"` // 매출 총 이익
	DeprCost     string `json:"depr_cost"`             // 감가상각비 (출력 안 되면 "99.99" — string 그대로)
	SellMang     string `json:"sell_mang"`             // 판매 및 관리비 (출력 안 되면 "99.99")
	BsopPrti     int64  `json:"bsop_prti,string"`      // 영업 이익
	BsopNonErnn  string `json:"bsop_non_ernn"`         // 영업 외 수익 (출력 안 되면 "99.99")
	BsopNonExpn  string `json:"bsop_non_expn"`         // 영업 외 비용 (출력 안 되면 "99.99")
	OpPrfi       int64  `json:"op_prfi,string"`        // 경상 이익
	SpecPrfi     int64  `json:"spec_prfi,string"`      // 특별 이익
	SpecLoss     int64  `json:"spec_loss,string"`      // 특별 손실
	ThtrNtin     int64  `json:"thtr_ntin,string"`      // 당기순이익
}

// InquireIncomeStatementParams 는 손익계산서 조회 파라미터.
type InquireIncomeStatementParams struct {
	Symbol  string // fid_input_iscd (필수)
	Quarter bool   // FID_DIV_CLS_CODE — false=>년, true=>분기 (분기는 누적 합산)
}

// InquireIncomeStatement 는 손익계산서 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_손익계산서.md
// path: /uapi/domestic-stock/v1/finance/income-statement (FHKST66430200)
//
// ※ 분기 데이터는 연단위 누적 합산. depr_cost / sell_mang / bsop_non_ernn / bsop_non_expn
// 은 출력 안 되면 "99.99" — caller 는 string 으로 받아 "99.99" 검사 후 처리.
func (c *Client) InquireIncomeStatement(ctx context.Context, params InquireIncomeStatementParams) (*IncomeStatement, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-stock/v1/finance/income-statement",
		TrID:     "FHKST66430200",
		Query:    inquireFinanceQuery(params.Symbol, params.Quarter),
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res IncomeStatement
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse IncomeStatement: %w", err)
	}
	return &res, nil
}

// BalanceSheet 는 대차대조표 (FHKST66430100) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_대차대조표.md
// path: /uapi/domestic-stock/v1/finance/balance-sheet
type BalanceSheet struct {
	Output []BalanceSheetItem `json:"output"`
}

// BalanceSheetItem 은 대차대조표 응답의 한 행.
type BalanceSheetItem struct {
	StacYymm  string `json:"stac_yymm"`         // 결산 년월
	Cras      int64  `json:"cras,string"`       // 유동자산
	Fxas      int64  `json:"fxas,string"`       // 고정자산
	TotalAset int64  `json:"total_aset,string"` // 자산총계
	FlowLblt  int64  `json:"flow_lblt,string"`  // 유동부채
	FixLblt   int64  `json:"fix_lblt,string"`   // 고정부채
	TotalLblt int64  `json:"total_lblt,string"` // 부채총계
	Cpfn      int64  `json:"cpfn,string"`       // 자본금
	CfpSurp   string `json:"cfp_surp"`          // 자본 잉여금 (출력 안 되면 "99.99")
	PrfiSurp  string `json:"prfi_surp"`         // 이익 잉여금 (출력 안 되면 "99.99")
	TotalCptl int64  `json:"total_cptl,string"` // 자본총계
}

// InquireBalanceSheetParams 는 대차대조표 조회 파라미터.
type InquireBalanceSheetParams struct {
	Symbol  string // fid_input_iscd (필수)
	Quarter bool   // FID_DIV_CLS_CODE
}

// InquireBalanceSheet 는 대차대조표 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_대차대조표.md
// path: /uapi/domestic-stock/v1/finance/balance-sheet (FHKST66430100)
func (c *Client) InquireBalanceSheet(ctx context.Context, params InquireBalanceSheetParams) (*BalanceSheet, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-stock/v1/finance/balance-sheet",
		TrID:     "FHKST66430100",
		Query:    inquireFinanceQuery(params.Symbol, params.Quarter),
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res BalanceSheet
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse BalanceSheet: %w", err)
	}
	return &res, nil
}

// ProfitRatio 는 수익성비율 (FHKST66430400) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_수익성비율.md
// path: /uapi/domestic-stock/v1/finance/profit-ratio
type ProfitRatio struct {
	Output []ProfitRatioItem `json:"output"`
}

// ProfitRatioItem 은 수익성비율 응답의 한 행.
type ProfitRatioItem struct {
	StacYymm         string  `json:"stac_yymm"`                  // 결산 년월
	CptlNtinRate     float64 `json:"cptl_ntin_rate,string"`      // 총자본 순이익율
	SelfCptlNtinInrt float64 `json:"self_cptl_ntin_inrt,string"` // 자기자본 순이익율
	SaleNtinRate     float64 `json:"sale_ntin_rate,string"`      // 매출액 순이익율
	SaleTotlRate     float64 `json:"sale_totl_rate,string"`      // 매출액 총이익율
}

// InquireProfitRatioParams 는 수익성비율 조회 파라미터.
type InquireProfitRatioParams struct {
	Symbol  string // fid_input_iscd (필수)
	Quarter bool   // FID_DIV_CLS_CODE
}

// InquireProfitRatio 는 수익성비율 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_수익성비율.md
// path: /uapi/domestic-stock/v1/finance/profit-ratio (FHKST66430400)
func (c *Client) InquireProfitRatio(ctx context.Context, params InquireProfitRatioParams) (*ProfitRatio, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-stock/v1/finance/profit-ratio",
		TrID:     "FHKST66430400",
		Query:    inquireFinanceQuery(params.Symbol, params.Quarter),
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res ProfitRatio
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse ProfitRatio: %w", err)
	}
	return &res, nil
}

// GrowthRatio 는 성장성비율 (FHKST66430800) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_성장성비율.md
// path: /uapi/domestic-stock/v1/finance/growth-ratio
type GrowthRatio struct {
	Output []GrowthRatioItem `json:"output"`
}

// GrowthRatioItem 은 성장성비율 응답의 한 행.
type GrowthRatioItem struct {
	StacYymm     string  `json:"stac_yymm"`             // 결산 년월
	Grs          float64 `json:"grs,string"`            // 매출액 증가율
	BsopPrfiInrt float64 `json:"bsop_prfi_inrt,string"` // 영업 이익 증가율
	EqutInrt     float64 `json:"equt_inrt,string"`      // 자기자본 증가율
	TotlAsetInrt float64 `json:"totl_aset_inrt,string"` // 총자산 증가율
}

// InquireGrowthRatioParams 는 성장성비율 조회 파라미터.
type InquireGrowthRatioParams struct {
	Symbol  string // fid_input_iscd (필수)
	Quarter bool   // fid_div_cls_code (소문자 — KIS docs 그대로). false=>"0"(년), true=>"1"(분기)
}

// InquireGrowthRatio 는 성장성비율 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_성장성비율.md
// path: /uapi/domestic-stock/v1/finance/growth-ratio (FHKST66430800)
//
// 다른 finance API 와 다르게 query 키가 fid_div_cls_code (소문자) — 다른 4 메서드는 FID_DIV_CLS_CODE (대문자).
// 그래서 inquireFinanceQuery helper 는 사용하지 않고 inline query.
func (c *Client) InquireGrowthRatio(ctx context.Context, params InquireGrowthRatioParams) (*GrowthRatio, error) {
	div := "0"
	if params.Quarter {
		div = "1"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/finance/growth-ratio",
		TrID:   "FHKST66430800",
		Query: map[string]string{
			"fid_input_iscd":         params.Symbol,
			"fid_div_cls_code":       div,
			"fid_cond_mrkt_div_code": "J",
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res GrowthRatio
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse GrowthRatio: %w", err)
	}
	return &res, nil
}

// OtherMajorRatios 는 기타주요비율 (FHKST66430500) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_기타주요비율.md
// path: /uapi/domestic-stock/v1/finance/other-major-ratios
type OtherMajorRatios struct {
	Output []OtherMajorRatiosItem `json:"output"`
}

// OtherMajorRatiosItem 은 기타주요비율 응답의 한 행.
type OtherMajorRatiosItem struct {
	StacYymm   string          `json:"stac_yymm"`        // 결산 년월
	PayoutRate string          `json:"payout_rate"`      // 배당 성향 — KIS 측 비정상 출력으로 무시 권고. string 보존.
	Eva        decimal.Decimal `json:"eva"`              // EVA (경제적 부가가치)
	Ebitda     decimal.Decimal `json:"ebitda"`           // EBITDA
	EvEbitda   float64         `json:"ev_ebitda,string"` // EV/EBITDA (배수)
}

// InquireOtherMajorRatiosParams 는 기타주요비율 조회 파라미터.
type InquireOtherMajorRatiosParams struct {
	Symbol  string // fid_input_iscd (필수)
	Quarter bool   // fid_div_cls_code (소문자 — InquireGrowthRatio 와 동일 패턴). false=>"0"(년), true=>"1"(분기)
}

// InquireOtherMajorRatios 는 기타주요비율 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_기타주요비율.md
// path: /uapi/domestic-stock/v1/finance/other-major-ratios (FHKST66430500)
//
// fid_div_cls_code 가 소문자 — InquireGrowthRatio 와 동일 패턴.
// inquireFinanceQuery helper 사용 불가 (helper 는 FID_DIV_CLS_CODE 대문자).
func (c *Client) InquireOtherMajorRatios(ctx context.Context, params InquireOtherMajorRatiosParams) (*OtherMajorRatios, error) {
	div := "0"
	if params.Quarter {
		div = "1"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/finance/other-major-ratios",
		TrID:   "FHKST66430500",
		Query: map[string]string{
			"fid_input_iscd":         params.Symbol,
			"fid_div_cls_code":       div,
			"fid_cond_mrkt_div_code": "J",
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res OtherMajorRatios
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse OtherMajorRatios: %w", err)
	}
	return &res, nil
}
