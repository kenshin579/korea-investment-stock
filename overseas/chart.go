package overseas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
	"github.com/kenshin579/korea-investment-stock/kistypes"
)

// DailyPrice 는 해외주식_기간별시세 (HHDFS76240000) 응답.
type DailyPrice struct {
	Output1 DailyPriceSummary  `json:"output1"`
	Output2 []DailyPriceCandle `json:"output2"`
}

// DailyPriceSummary 는 응답의 output1 (단일 객체).
type DailyPriceSummary struct {
	Rsym string `json:"rsym"`
	Zdiv string `json:"zdiv"`
	Nrec string `json:"nrec"`
}

// DailyPriceCandle 은 응답의 output2 한 행 (한 일자).
type DailyPriceCandle struct {
	Xymd string          `json:"xymd"`
	Clos decimal.Decimal `json:"clos"`
	Sign string          `json:"sign"`
	Diff decimal.Decimal `json:"diff"`
	Rate kistypes.Float  `json:"rate"`
	Open decimal.Decimal `json:"open"`
	High decimal.Decimal `json:"high"`
	Low  decimal.Decimal `json:"low"`
	Tvol int64           `json:"tvol,string"`
	Tamt int64           `json:"tamt,string"`
	Pbid decimal.Decimal `json:"pbid"`
	Vbid int64           `json:"vbid,string"`
	Pask decimal.Decimal `json:"pask"`
	Vask int64           `json:"vask,string"`
}

// InquireDailyPriceParams 는 해외주식_기간별시세 조회 파라미터.
type InquireDailyPriceParams struct {
	Auth string // AUTH — 빈 값 default
	Excd string // EXCD — 거래소코드 (HKS/NYS/NAS/AMS/TSE/SHS/SZS/SHI/SZI/HSX/HNX)
	Symb string // SYMB — 종목코드
	Gubn string // GUBN — "0":일/"1":주/"2":월. 빈 값=>"0"
	Bymd string // BYMD — 조회기준일자 (YYYYMMDD). 빈 값=>오늘
	Modp string // MODP — "0":미반영/"1":반영(수정주가). 빈 값=>"0"
	Keyb string // KEYB — NEXT KEY BUFF
}

// InquireDailyPrice 는 해외주식_기간별시세 호출.
//
// 한투 docs: docs/api/해외주식/해외주식_기간별시세.md
// path: /uapi/overseas-price/v1/quotations/dailyprice (HHDFS76240000)
//
// 한 번의 호출에 최대 100건. 미국은 0분지연, 홍콩/베트남/중국/일본은 15분지연시세.
func (c *Client) InquireDailyPrice(ctx context.Context, params InquireDailyPriceParams) (*DailyPrice, error) {
	gubn := params.Gubn
	if gubn == "" {
		gubn = "0"
	}
	modp := params.Modp
	if modp == "" {
		modp = "0"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-price/v1/quotations/dailyprice",
		TrID:   "HHDFS76240000",
		Query: map[string]string{
			"AUTH": params.Auth,
			"EXCD": params.Excd,
			"SYMB": params.Symb,
			"GUBN": gubn,
			"BYMD": params.Bymd,
			"MODP": modp,
			"KEYB": params.Keyb,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res DailyPrice
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse DailyPrice: %w", err)
	}
	return &res, nil
}

// DailyChartPrice 는 해외주식 종목/지수/환율 기간별시세 (FHKST03030100) 응답.
//
// 한투 docs: docs/api/해외주식/해외주식_종목_지수_환율기간별시세(일_주_월_년).md
// path: /uapi/overseas-price/v1/quotations/inquire-daily-chartprice
//
// 미국 주식은 다우30/나스닥100/S&P500 만 조회 가능 (다른 미국 종목은 InquireDailyPrice).
type DailyChartPrice struct {
	Output1 DailyChartPriceSummary  `json:"output1"`
	Output2 []DailyChartPriceCandle `json:"output2"`
}

// DailyChartPriceSummary 는 응답의 output1 (단일 객체, 기본정보).
type DailyChartPriceSummary struct {
	OvrsNmixPrdyVrss decimal.Decimal `json:"ovrs_nmix_prdy_vrss"`
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`
	PrdyCtrt         kistypes.Float  `json:"prdy_ctrt"`
	OvrsNmixPrdyClpr decimal.Decimal `json:"ovrs_nmix_prdy_clpr"`
	AcmlVol          int64           `json:"acml_vol,string"`
	HtsKorIsnm       string          `json:"hts_kor_isnm"`
	OvrsNmixPrpr     decimal.Decimal `json:"ovrs_nmix_prpr"`
	StckShrnIscd     string          `json:"stck_shrn_iscd"`
	PrdyVol          int64           `json:"prdy_vol,string"`
	OvrsProdOprc     decimal.Decimal `json:"ovrs_prod_oprc"`
	OvrsProdHgpr     decimal.Decimal `json:"ovrs_prod_hgpr"`
	OvrsProdLwpr     decimal.Decimal `json:"ovrs_prod_lwpr"`
}

// DailyChartPriceCandle 은 응답의 output2 한 행 (한 일자/주/월/년 봉).
type DailyChartPriceCandle struct {
	StckBsopDate string          `json:"stck_bsop_date"`
	OvrsNmixPrpr decimal.Decimal `json:"ovrs_nmix_prpr"`
	OvrsNmixOprc decimal.Decimal `json:"ovrs_nmix_oprc"`
	OvrsNmixHgpr decimal.Decimal `json:"ovrs_nmix_hgpr"`
	OvrsNmixLwpr decimal.Decimal `json:"ovrs_nmix_lwpr"`
	AcmlVol      int64           `json:"acml_vol,string"`
	ModYn        string          `json:"mod_yn"`
}

// InquireDailyChartPriceParams 는 해외주식 종목/지수/환율 기간별시세 조회 파라미터.
type InquireDailyChartPriceParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE — "N":해외지수/"X":환율/"I":국채/"S":금선물
	Symbol     string // FID_INPUT_ISCD
	FromDate   string // FID_INPUT_DATE_1 (YYYYMMDD)
	ToDate     string // FID_INPUT_DATE_2 (YYYYMMDD)
	Period     string // FID_PERIOD_DIV_CODE — "D"/"W"/"M"/"Y"
}

// InquireDailyChartPrice 는 해외주식 종목/지수/환율 기간별시세 호출.
//
// 한투 docs: docs/api/해외주식/해외주식_종목_지수_환율기간별시세(일_주_월_년).md
// path: /uapi/overseas-price/v1/quotations/inquire-daily-chartprice (FHKST03030100)
//
// ※ 미국 주식 조회 시 다우30/나스닥100/S&P500 종목만 가능. 다른 미국 종목은 InquireDailyPrice 사용.
func (c *Client) InquireDailyChartPrice(ctx context.Context, params InquireDailyChartPriceParams) (*DailyChartPrice, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-price/v1/quotations/inquire-daily-chartprice",
		TrID:   "FHKST03030100",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": params.MarketCode,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_INPUT_DATE_1":       params.FromDate,
			"FID_INPUT_DATE_2":       params.ToDate,
			"FID_PERIOD_DIV_CODE":    params.Period,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res DailyChartPrice
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse DailyChartPrice: %w", err)
	}
	return &res, nil
}
