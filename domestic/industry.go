package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// IndexPrice 는 국내업종 현재지수 (FHPUP02100000) 응답.
//
// 한투 docs: docs/api/국내주식/국내업종_현재지수.md
// path: /uapi/domestic-stock/v1/quotations/inquire-index-price
type IndexPrice struct {
	Output IndexPriceSnapshot `json:"output"`
}

// IndexPriceSnapshot 은 응답의 output (단일 객체).
//
// KIS docs 의 line 73~108 모든 필드 1:1 매핑 (~36 fields).
type IndexPriceSnapshot struct {
	// 지수 + 변동
	BstpNmixPrpr     decimal.Decimal `json:"bstp_nmix_prpr"`             // 업종 지수 현재가
	BstpNmixPrdyVrss decimal.Decimal `json:"bstp_nmix_prdy_vrss"`        // 업종 지수 전일 대비
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`             // 전일 대비 부호
	BstpNmixPrdyCtrt float64         `json:"bstp_nmix_prdy_ctrt,string"` // 업종 지수 전일 대비율

	// 거래량/거래대금
	AcmlVol    int64 `json:"acml_vol,string"`     // 누적 거래량
	PrdyVol    int64 `json:"prdy_vol,string"`     // 전일 거래량
	AcmlTrPbmn int64 `json:"acml_tr_pbmn,string"` // 누적 거래 대금
	PrdyTrPbmn int64 `json:"prdy_tr_pbmn,string"` // 전일 거래 대금

	// 시가 + 시가대비
	BstpNmixOprc         decimal.Decimal `json:"bstp_nmix_oprc"`                  // 업종 지수 시가
	PrdyNmixVrssNmixOprc decimal.Decimal `json:"prdy_nmix_vrss_nmix_oprc"`        // 전일 지수 대비 지수 시가
	OprcVrssPrprSign     string          `json:"oprc_vrss_prpr_sign"`             // 시가 대비 현재가 부호
	BstpNmixOprcPrdyCtrt float64         `json:"bstp_nmix_oprc_prdy_ctrt,string"` // 업종 지수 시가 전일 대비율

	// 최고가
	BstpNmixHgpr         decimal.Decimal `json:"bstp_nmix_hgpr"`                  // 업종 지수 최고가
	PrdyNmixVrssNmixHgpr decimal.Decimal `json:"prdy_nmix_vrss_nmix_hgpr"`        // 전일 지수 대비 지수 최고가
	HgprVrssPrprSign     string          `json:"hgpr_vrss_prpr_sign"`             // 최고가 대비 현재가 부호
	BstpNmixHgprPrdyCtrt float64         `json:"bstp_nmix_hgpr_prdy_ctrt,string"` // 업종 지수 최고가 전일 대비율

	// 최저가
	BstpNmixLwpr         decimal.Decimal `json:"bstp_nmix_lwpr"`                  // 업종 지수 최저가
	PrdyClprVrssLwpr     decimal.Decimal `json:"prdy_clpr_vrss_lwpr"`             // 전일 종가 대비 최저가
	LwprVrssPrprSign     string          `json:"lwpr_vrss_prpr_sign"`             // 최저가 대비 현재가 부호
	PrdyClprVrssLwprRate float64         `json:"prdy_clpr_vrss_lwpr_rate,string"` // 전일 종가 대비 최저가 비율

	// 종목수 (5 fields, KIS 가 string 으로 줌 — 작은 정수)
	AscnIssuCnt string `json:"ascn_issu_cnt"` // 상승 종목 수
	UplmIssuCnt string `json:"uplm_issu_cnt"` // 상한 종목 수
	StnrIssuCnt string `json:"stnr_issu_cnt"` // 보합 종목 수
	DownIssuCnt string `json:"down_issu_cnt"` // 하락 종목 수
	LslmIssuCnt string `json:"lslm_issu_cnt"` // 하한 종목 수

	// 연중 (6 fields)
	DryyBstpNmixHgpr     decimal.Decimal `json:"dryy_bstp_nmix_hgpr"`             // 연중 업종 지수 최고가
	DryyHgprVrssPrprRate float64         `json:"dryy_hgpr_vrss_prpr_rate,string"` // 연중 최고가 대비 현재가 비율
	DryyBstpNmixHgprDate string          `json:"dryy_bstp_nmix_hgpr_date"`        // 연중 업종 지수 최고가 일자
	DryyBstpNmixLwpr     decimal.Decimal `json:"dryy_bstp_nmix_lwpr"`             // 연중 업종 지수 최저가
	DryyLwprVrssPrprRate float64         `json:"dryy_lwpr_vrss_prpr_rate,string"` // 연중 최저가 대비 현재가 비율
	DryyBstpNmixLwprDate string          `json:"dryy_bstp_nmix_lwpr_date"`        // 연중 업종 지수 최저가 일자

	// 호가 잔량 (5 fields)
	TotalAskpRsqn int64   `json:"total_askp_rsqn,string"` // 총 매도호가 잔량
	TotalBidpRsqn int64   `json:"total_bidp_rsqn,string"` // 총 매수호가 잔량
	SelnRsqnRate  float64 `json:"seln_rsqn_rate,string"`  // 매도 잔량 비율
	ShnuRsqnRate  float64 `json:"shnu_rsqn_rate,string"`  // 매수 잔량 비율
	NtbyRsqn      int64   `json:"ntby_rsqn,string"`       // 순매수 잔량
}

// InquireIndexPriceParams 는 국내업종 현재지수 조회 파라미터.
type InquireIndexPriceParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE — 빈 값=>"U" (업종)
	Symbol     string // FID_INPUT_ISCD — 필수, 업종 코드 (예 "0001":코스피, "1001":코스닥, "2001":코스피200)
}

// InquireIndexPrice 는 국내업종 현재지수 호출.
//
// 한투 docs: docs/api/국내주식/국내업종_현재지수.md
// path: /uapi/domestic-stock/v1/quotations/inquire-index-price (FHPUP02100000)
func (c *Client) InquireIndexPrice(ctx context.Context, params InquireIndexPriceParams) (*IndexPrice, error) {
	market := params.MarketCode
	if market == "" {
		market = "U"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-index-price",
		TrID:   "FHPUP02100000",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_INPUT_ISCD":         params.Symbol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res IndexPrice
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse IndexPrice: %w", err)
	}
	return &res, nil
}
