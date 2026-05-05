// File: domestic/program_trade.go
package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// ProgramTradeByStockDaily 는 종목별 프로그램매매추이(일별) (FHPPG04650201) 응답.
//
// 한투 docs: docs/api/국내주식/종목별_프로그램매매추이(일별).md
// path: /uapi/domestic-stock/v1/quotations/program-trade-by-stock-daily
type ProgramTradeByStockDaily struct {
	Output []ProgramTradeByStockDailyItem `json:"output"`
}

// ProgramTradeByStockDailyItem 은 응답 output 한 행 (일자별 프로그램 매매).
//
// WholNtbyTrPbmnIcdc2: 마지막 필드. 필드명 trailing "2" — EP4 의 WholNtbyTrPbmnIcdc 와 구분.
type ProgramTradeByStockDailyItem struct {
	StckBsopDate        string          `json:"stck_bsop_date"`                 // 주식 영업 일자
	StckClpr            decimal.Decimal `json:"stck_clpr"`                      // 주식 종가
	PrdyVrss            decimal.Decimal `json:"prdy_vrss"`                      // 전일 대비
	PrdyVrssSign        string          `json:"prdy_vrss_sign"`                 // 전일 대비 부호
	PrdyCtrt            float64         `json:"prdy_ctrt,string"`               // 전일 대비율
	AcmlVol             int64           `json:"acml_vol,string"`                // 누적 거래량
	AcmlTrPbmn          int64           `json:"acml_tr_pbmn,string"`            // 누적 거래 대금
	WholSmtnSelnVol     int64           `json:"whol_smtn_seln_vol,string"`      // 전체 합산 매도 수량
	WholSmtnShnuVol     int64           `json:"whol_smtn_shnu_vol,string"`      // 전체 합산 매수 수량
	WholSmtnNtbyQty     int64           `json:"whol_smtn_ntby_qty,string"`      // 전체 합산 순매수 수량
	WholSmtnSelnTrPbmn  int64           `json:"whol_smtn_seln_tr_pbmn,string"`  // 전체 합산 매도 거래대금
	WholSmtnShnuTrPbmn  int64           `json:"whol_smtn_shnu_tr_pbmn,string"`  // 전체 합산 매수 거래대금
	WholSmtnNtbyTrPbmn  int64           `json:"whol_smtn_ntby_tr_pbmn,string"`  // 전체 합산 순매수 거래대금
	WholNtbyVolIcdc     int64           `json:"whol_ntby_vol_icdc,string"`      // 전체 순매수 거래량 증감
	WholNtbyTrPbmnIcdc2 int64           `json:"whol_ntby_tr_pbmn_icdc2,string"` // 전체 순매수 거래대금 증감 (trailing "2")
}

// InquireProgramTradeByStockDailyParams 는 종목별 프로그램매매추이(일별) 조회 파라미터.
//
// BaseDate: KIS docs 예시가 "002" prefix 포함 ("0020240308") — 호출자가 raw string 그대로 전달.
// MarketCode: "J"(KRX), "NX"(NXT), "UN"(통합). 빈 값=>"J".
type InquireProgramTradeByStockDailyParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE — 빈 값=>"J"
	Symbol     string // FID_INPUT_ISCD — 필수
	BaseDate   string // FID_INPUT_DATE_1 — 필수, KIS docs 예시: "0020240308" ("002" prefix)
}

// InquireProgramTradeByStockDaily 는 종목별 프로그램매매추이(일별) 호출.
//
// 한투 docs: docs/api/국내주식/종목별_프로그램매매추이(일별).md
// path: /uapi/domestic-stock/v1/quotations/program-trade-by-stock-daily (FHPPG04650201)
func (c *Client) InquireProgramTradeByStockDaily(ctx context.Context, params InquireProgramTradeByStockDailyParams) (*ProgramTradeByStockDaily, error) {
	mkt := params.MarketCode
	if mkt == "" {
		mkt = "J"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/program-trade-by-stock-daily",
		TrID:   "FHPPG04650201",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": mkt,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_INPUT_DATE_1":       params.BaseDate,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res ProgramTradeByStockDaily
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse ProgramTradeByStockDaily: %w", err)
	}
	return &res, nil
}
