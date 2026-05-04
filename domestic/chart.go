package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// DailyChart 는 국내주식기간별시세(일/주/월/년) (FHKST03010100) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식기간별시세(일_주_월_년).md
type DailyChart struct {
	Output1 DailyChartSummary  `json:"output1"`
	Output2 []DailyChartCandle `json:"output2"`
}

// DailyChartSummary 는 차트 응답의 output1 (단일 객체, 요약 정보).
type DailyChartSummary struct {
	PrdyVrss     decimal.Decimal `json:"prdy_vrss"`
	PrdyVrssSign string          `json:"prdy_vrss_sign"`
	PrdyCtrt     float64         `json:"prdy_ctrt,string"`
	StckPrdyClpr decimal.Decimal `json:"stck_prdy_clpr"`
	AcmlVol      int64           `json:"acml_vol,string"`
	AcmlTrPbmn   int64           `json:"acml_tr_pbmn,string"`
	HtsKorIsnm   string          `json:"hts_kor_isnm"`
	StckPrpr     decimal.Decimal `json:"stck_prpr"`
	StckShrnIscd string          `json:"stck_shrn_iscd"`
	PrdyVol      int64           `json:"prdy_vol,string"`
	StckMxpr     decimal.Decimal `json:"stck_mxpr"`
	StckLlam     decimal.Decimal `json:"stck_llam"`
	StckOprc     decimal.Decimal `json:"stck_oprc"`
	StckHgpr     decimal.Decimal `json:"stck_hgpr"`
	StckLwpr     decimal.Decimal `json:"stck_lwpr"`
	StckPrdyOprc decimal.Decimal `json:"stck_prdy_oprc"`
	StckPrdyHgpr decimal.Decimal `json:"stck_prdy_hgpr"`
	StckPrdyLwpr decimal.Decimal `json:"stck_prdy_lwpr"`
	Askp         decimal.Decimal `json:"askp"`
	Bidp         decimal.Decimal `json:"bidp"`
	PrdyVrssVol  int64           `json:"prdy_vrss_vol,string"`
	VolTnrt      float64         `json:"vol_tnrt,string"`
	StckFcam     decimal.Decimal `json:"stck_fcam"`
	LstnStcn     int64           `json:"lstn_stcn,string"`
	Cpfn         int64           `json:"cpfn,string"`
	HtsAvls      int64           `json:"hts_avls,string"`
	Per          float64         `json:"per,string"`
	Eps          decimal.Decimal `json:"eps"`
	Pbr          float64         `json:"pbr,string"`
}

// DailyChartCandle 은 차트 응답의 output2 한 행 (한 캔들).
type DailyChartCandle struct {
	StckBsopDate string          `json:"stck_bsop_date"`
	StckClpr     decimal.Decimal `json:"stck_clpr"`
	StckOprc     decimal.Decimal `json:"stck_oprc"`
	StckHgpr     decimal.Decimal `json:"stck_hgpr"`
	StckLwpr     decimal.Decimal `json:"stck_lwpr"`
	AcmlVol      int64           `json:"acml_vol,string"`
	AcmlTrPbmn   int64           `json:"acml_tr_pbmn,string"`
	FlngClsCode  string          `json:"flng_cls_code"`
	PrttRate     float64         `json:"prtt_rate,string"`
	ModYn        string          `json:"mod_yn"`
	PrdyVrssSign string          `json:"prdy_vrss_sign"`
	PrdyVrss     decimal.Decimal `json:"prdy_vrss"`
	RevlIssuReas string          `json:"revl_issu_reas"`
}

// InquireDailyItemChartPriceParams 는 일/주/월/년 봉 조회 파라미터.
type InquireDailyItemChartPriceParams struct {
	Symbol        string // 필수, 종목코드 (예 "005930")
	Period        string // "D"/"W"/"M"/"Y", 빈 값이면 "D"
	FromDate      string // YYYYMMDD, 필수
	ToDate        string // YYYYMMDD, 필수, 1회 최대 100건
	OriginalPrice bool   // false=수정주가(default), true=원주가
	MarketCode    string // "J"/"NX"/"UN", 빈 값이면 "J"
}

// InquireDailyItemChartPrice 는 국내주식기간별시세(일/주/월/년) 호출.
//
// 한투 docs: docs/api/국내주식/국내주식기간별시세(일_주_월_년).md
// path: /uapi/domestic-stock/v1/quotations/inquire-daily-itemchartprice (FHKST03010100)
func (c *Client) InquireDailyItemChartPrice(ctx context.Context, params InquireDailyItemChartPriceParams) (*DailyChart, error) {
	period := params.Period
	if period == "" {
		period = "D"
	}
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	adjPrc := "0" // 수정주가
	if params.OriginalPrice {
		adjPrc = "1"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-daily-itemchartprice",
		TrID:   "FHKST03010100",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_INPUT_DATE_1":       params.FromDate,
			"FID_INPUT_DATE_2":       params.ToDate,
			"FID_PERIOD_DIV_CODE":    period,
			"FID_ORG_ADJ_PRC":        adjPrc,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	// output1 + output2 동시 unmarshal: resp.Raw 에서 파싱
	var chart DailyChart
	if err := json.Unmarshal(resp.Raw, &chart); err != nil {
		return nil, fmt.Errorf("kis: parse DailyChart: %w", err)
	}
	return &chart, nil
}
