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

// IndexCategoryPrice 는 국내업종 구분별 전체시세 (FHPUP02140000) 응답.
//
// 한투 docs: docs/api/국내주식/국내업종_구분별전체시세.md
// path: /uapi/domestic-stock/v1/quotations/inquire-index-category-price
type IndexCategoryPrice struct {
	Output1 IndexCategoryPriceSummary `json:"output1"`
	Output2 []IndexCategoryPriceItem  `json:"output2"`
}

// IndexCategoryPriceSummary 는 응답의 output1 (대표 업종 지수).
type IndexCategoryPriceSummary struct {
	BstpNmixPrpr         decimal.Decimal `json:"bstp_nmix_prpr"`
	BstpNmixPrdyVrss     decimal.Decimal `json:"bstp_nmix_prdy_vrss"`
	PrdyVrssSign         string          `json:"prdy_vrss_sign"`
	BstpNmixPrdyCtrt     float64         `json:"bstp_nmix_prdy_ctrt,string"`
	AcmlVol              int64           `json:"acml_vol,string"`
	AcmlTrPbmn           int64           `json:"acml_tr_pbmn,string"`
	BstpNmixOprc         decimal.Decimal `json:"bstp_nmix_oprc"`
	BstpNmixHgpr         decimal.Decimal `json:"bstp_nmix_hgpr"`
	BstpNmixLwpr         decimal.Decimal `json:"bstp_nmix_lwpr"`
	PrdyVol              int64           `json:"prdy_vol,string"`
	AscnIssuCnt          string          `json:"ascn_issu_cnt"`
	DownIssuCnt          string          `json:"down_issu_cnt"`
	StnrIssuCnt          string          `json:"stnr_issu_cnt"`
	UplmIssuCnt          string          `json:"uplm_issu_cnt"`
	LslmIssuCnt          string          `json:"lslm_issu_cnt"`
	PrdyTrPbmn           int64           `json:"prdy_tr_pbmn,string"`
	DryyBstpNmixHgprDate string          `json:"dryy_bstp_nmix_hgpr_date"`
	DryyBstpNmixHgpr     decimal.Decimal `json:"dryy_bstp_nmix_hgpr"`
	DryyBstpNmixLwpr     decimal.Decimal `json:"dryy_bstp_nmix_lwpr"`
	DryyBstpNmixLwprDate string          `json:"dryy_bstp_nmix_lwpr_date"`
}

// IndexCategoryPriceItem 은 응답의 output2 한 행 (구분별 업종).
type IndexCategoryPriceItem struct {
	BstpClsCode      string          `json:"bstp_cls_code"`              // 업종 구분 코드
	HtsKorIsnm       string          `json:"hts_kor_isnm"`               // HTS 한글 종목명
	BstpNmixPrpr     decimal.Decimal `json:"bstp_nmix_prpr"`             // 업종 지수 현재가
	BstpNmixPrdyVrss decimal.Decimal `json:"bstp_nmix_prdy_vrss"`        // 업종 지수 전일 대비
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`             // 전일 대비 부호
	BstpNmixPrdyCtrt float64         `json:"bstp_nmix_prdy_ctrt,string"` // 업종 지수 전일 대비율
	AcmlVol          int64           `json:"acml_vol,string"`            // 누적 거래량
	AcmlTrPbmn       int64           `json:"acml_tr_pbmn,string"`        // 누적 거래 대금
	AcmlVolRlim      float64         `json:"acml_vol_rlim,string"`       // 누적 거래량 비중
	AcmlTrPbmnRlim   float64         `json:"acml_tr_pbmn_rlim,string"`   // 누적 거래 대금 비중
}

// InquireIndexCategoryPriceParams 는 국내업종 구분별 전체시세 조회 파라미터.
type InquireIndexCategoryPriceParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE — 빈 값=>"U" (업종)
	Symbol     string // FID_INPUT_ISCD — 필수, 업종 코드 (코스피 0001 등)
	ScreenCode string // FID_COND_SCR_DIV_CODE — 빈 값=>"20214"
	MarketCls  string // FID_MRKT_CLS_CODE — "K":거래소, "Q":코스닥, "K2":코스피200
	BelongCls  string // FID_BLNG_CLS_CODE — 빈 값=>"0" (전업종)
}

// InquireIndexCategoryPrice 는 국내업종 구분별 전체시세 호출.
//
// 한투 docs: docs/api/국내주식/국내업종_구분별전체시세.md
// path: /uapi/domestic-stock/v1/quotations/inquire-index-category-price (FHPUP02140000)
func (c *Client) InquireIndexCategoryPrice(ctx context.Context, params InquireIndexCategoryPriceParams) (*IndexCategoryPrice, error) {
	market := params.MarketCode
	if market == "" {
		market = "U"
	}
	scr := params.ScreenCode
	if scr == "" {
		scr = "20214"
	}
	belong := params.BelongCls
	if belong == "" {
		belong = "0"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-index-category-price",
		TrID:   "FHPUP02140000",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_COND_SCR_DIV_CODE":  scr,
			"FID_MRKT_CLS_CODE":      params.MarketCls,
			"FID_BLNG_CLS_CODE":      belong,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res IndexCategoryPrice
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse IndexCategoryPrice: %w", err)
	}
	return &res, nil
}

// IndexDailyPrice 는 국내업종 일별지수 (FHPUP02120000) 응답.
//
// 한투 docs: docs/api/국내주식/국내업종_일별지수.md
// path: /uapi/domestic-stock/v1/quotations/inquire-index-daily-price
type IndexDailyPrice struct {
	Output1 IndexDailyPriceSummary `json:"output1"`
	Output2 []IndexDailyPriceItem  `json:"output2"`
}

// IndexDailyPriceSummary 는 응답의 output1 (대표 스냅샷, 20 fields).
type IndexDailyPriceSummary struct {
	BstpNmixPrpr         decimal.Decimal `json:"bstp_nmix_prpr"`             // 업종 지수 현재가
	BstpNmixPrdyVrss     decimal.Decimal `json:"bstp_nmix_prdy_vrss"`        // 업종 지수 전일 대비
	PrdyVrssSign         string          `json:"prdy_vrss_sign"`             // 전일 대비 부호
	BstpNmixPrdyCtrt     float64         `json:"bstp_nmix_prdy_ctrt,string"` // 업종 지수 전일 대비율
	AcmlVol              int64           `json:"acml_vol,string"`            // 누적 거래량
	AcmlTrPbmn           int64           `json:"acml_tr_pbmn,string"`        // 누적 거래 대금
	BstpNmixOprc         decimal.Decimal `json:"bstp_nmix_oprc"`             // 업종 지수 시가
	BstpNmixHgpr         decimal.Decimal `json:"bstp_nmix_hgpr"`             // 업종 지수 최고가
	BstpNmixLwpr         decimal.Decimal `json:"bstp_nmix_lwpr"`             // 업종 지수 최저가
	PrdyVol              int64           `json:"prdy_vol,string"`            // 전일 거래량
	AscnIssuCnt          string          `json:"ascn_issu_cnt"`              // 상승 종목 수
	DownIssuCnt          string          `json:"down_issu_cnt"`              // 하락 종목 수
	StnrIssuCnt          string          `json:"stnr_issu_cnt"`              // 보합 종목 수
	UplmIssuCnt          string          `json:"uplm_issu_cnt"`              // 상한 종목 수
	LslmIssuCnt          string          `json:"lslm_issu_cnt"`              // 하한 종목 수
	PrdyTrPbmn           int64           `json:"prdy_tr_pbmn,string"`        // 전일 거래 대금
	DryyBstpNmixHgprDate string          `json:"dryy_bstp_nmix_hgpr_date"`   // 연중 최고가 일자
	DryyBstpNmixHgpr     decimal.Decimal `json:"dryy_bstp_nmix_hgpr"`        // 연중 업종 지수 최고가
	DryyBstpNmixLwpr     decimal.Decimal `json:"dryy_bstp_nmix_lwpr"`        // 연중 업종 지수 최저가
	DryyBstpNmixLwprDate string          `json:"dryy_bstp_nmix_lwpr_date"`   // 연중 최저가 일자
}

// IndexDailyPriceItem 은 응답의 output2 한 행 (일별, 13 fields).
type IndexDailyPriceItem struct {
	StckBsopDate     string          `json:"stck_bsop_date"`             // 영업 일자
	BstpNmixPrpr     decimal.Decimal `json:"bstp_nmix_prpr"`             // 업종 지수 현재가
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`             // 전일 대비 부호
	BstpNmixPrdyVrss decimal.Decimal `json:"bstp_nmix_prdy_vrss"`        // 업종 지수 전일 대비
	BstpNmixPrdyCtrt float64         `json:"bstp_nmix_prdy_ctrt,string"` // 업종 지수 전일 대비율
	BstpNmixOprc     decimal.Decimal `json:"bstp_nmix_oprc"`             // 업종 지수 시가
	BstpNmixHgpr     decimal.Decimal `json:"bstp_nmix_hgpr"`             // 업종 지수 최고가
	BstpNmixLwpr     decimal.Decimal `json:"bstp_nmix_lwpr"`             // 업종 지수 최저가
	AcmlVolRlim      float64         `json:"acml_vol_rlim,string"`       // 누적 거래량 비중
	AcmlVol          int64           `json:"acml_vol,string"`            // 누적 거래량
	AcmlTrPbmn       int64           `json:"acml_tr_pbmn,string"`        // 누적 거래 대금
	InvtNewPsdg      decimal.Decimal `json:"invt_new_psdg"`              // 투자자 순매수 주도
	D20Dsrt          decimal.Decimal `json:"d20_dsrt"`                   // 20일 이격도
}

// InquireIndexDailyPriceParams 는 국내업종 일별지수 조회 파라미터.
type InquireIndexDailyPriceParams struct {
	MarketCode    string // FID_COND_MRKT_DIV_CODE — 빈 값=>"U" (업종)
	Symbol        string // FID_INPUT_ISCD — 필수, 업종 코드 (예 "0001":코스피)
	PeriodDivCode string // FID_PERIOD_DIV_CODE — D:일 W:주 M:월 Y:년
	InputDate1    string // FID_INPUT_DATE_1 — 조회 시작일 YYYYMMDD
}

// InquireIndexDailyPrice 는 국내업종 일별지수 호출.
//
// 한투 docs: docs/api/국내주식/국내업종_일별지수.md
// path: /uapi/domestic-stock/v1/quotations/inquire-index-daily-price (FHPUP02120000)
func (c *Client) InquireIndexDailyPrice(ctx context.Context, params InquireIndexDailyPriceParams) (*IndexDailyPrice, error) {
	market := params.MarketCode
	if market == "" {
		market = "U"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-index-daily-price",
		TrID:   "FHPUP02120000",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_PERIOD_DIV_CODE":    params.PeriodDivCode,
			"FID_INPUT_DATE_1":       params.InputDate1,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res IndexDailyPrice
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse IndexDailyPrice: %w", err)
	}
	return &res, nil
}

// IndexTimeprice 는 국내업종 시간별 지수 (FHPUP02110200) 응답.
//
// 한투 docs: docs/api/국내주식/국내업종_시간별지수.md
// path: /uapi/domestic-stock/v1/quotations/inquire-index-timeprice
//
// BsopHour 는 HHMMSS 형식 타임스탬프. FID_INPUT_HOUR_1 파라미터로 집계 단위 설정 (60/300/600초).
type IndexTimeprice struct {
	Output []IndexTimepriceItem `json:"output"`
}

// IndexTimepriceItem 은 응답의 output 한 행 (시간별, 8 fields).
type IndexTimepriceItem struct {
	BsopHour         string          `json:"bsop_hour"`                  // 영업 시간 HHMMSS
	BstpNmixPrpr     decimal.Decimal `json:"bstp_nmix_prpr"`             // 업종 지수 현재가
	BstpNmixPrdyVrss decimal.Decimal `json:"bstp_nmix_prdy_vrss"`        // 업종 지수 전일 대비
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`             // 전일 대비 부호
	BstpNmixPrdyCtrt float64         `json:"bstp_nmix_prdy_ctrt,string"` // 업종 지수 전일 대비율
	AcmlTrPbmn       int64           `json:"acml_tr_pbmn,string"`        // 누적 거래 대금
	AcmlVol          int64           `json:"acml_vol,string"`            // 누적 거래량
	CntgVol          int64           `json:"cntg_vol,string"`            // 체결 거래량
}

// InquireIndexTimepriceParams 는 국내업종 시간별 지수 조회 파라미터.
type InquireIndexTimepriceParams struct {
	InputHour1 string // FID_INPUT_HOUR_1 — 집계 단위: "60"(1분)/"300"(5분)/"600"(10분)
	Symbol     string // FID_INPUT_ISCD — 필수, 업종 코드
	MarketCode string // FID_COND_MRKT_DIV_CODE — 빈 값=>"U" (업종)
}

// InquireIndexTimeprice 는 국내업종 시간별 지수 호출.
//
// 한투 docs: docs/api/국내주식/국내업종_시간별지수.md
// path: /uapi/domestic-stock/v1/quotations/inquire-index-timeprice (FHPUP02110200)
func (c *Client) InquireIndexTimeprice(ctx context.Context, params InquireIndexTimepriceParams) (*IndexTimeprice, error) {
	market := params.MarketCode
	if market == "" {
		market = "U"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-index-timeprice",
		TrID:   "FHPUP02110200",
		Query: map[string]string{
			"FID_INPUT_HOUR_1":       params.InputHour1,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_COND_MRKT_DIV_CODE": market,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res IndexTimeprice
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse IndexTimeprice: %w", err)
	}
	return &res, nil
}

// IndexTickprice 는 국내업종 틱별 지수 (FHPUP02110100) 응답.
//
// 한투 docs: docs/api/국내주식/국내업종_틱별지수.md
// path: /uapi/domestic-stock/v1/quotations/inquire-index-tickprice
//
// StckCntgHour 는 HHMMSS 형식 틱 타임스탬프.
type IndexTickprice struct {
	Output []IndexTickpriceItem `json:"output"`
}

// IndexTickpriceItem 은 응답의 output 한 행 (틱별, 8 fields).
type IndexTickpriceItem struct {
	StckCntgHour     string          `json:"stck_cntg_hour"`             // 주식 체결 시간 HHMMSS
	BstpNmixPrpr     decimal.Decimal `json:"bstp_nmix_prpr"`             // 업종 지수 현재가
	BstpNmixPrdyVrss decimal.Decimal `json:"bstp_nmix_prdy_vrss"`        // 업종 지수 전일 대비
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`             // 전일 대비 부호
	BstpNmixPrdyCtrt float64         `json:"bstp_nmix_prdy_ctrt,string"` // 업종 지수 전일 대비율
	AcmlTrPbmn       int64           `json:"acml_tr_pbmn,string"`        // 누적 거래 대금
	AcmlVol          int64           `json:"acml_vol,string"`            // 누적 거래량
	CntgVol          int64           `json:"cntg_vol,string"`            // 체결 거래량
}

// InquireIndexTickpriceParams 는 국내업종 틱별 지수 조회 파라미터.
type InquireIndexTickpriceParams struct {
	Symbol     string // FID_INPUT_ISCD — 필수, 업종 코드
	MarketCode string // FID_COND_MRKT_DIV_CODE — 빈 값=>"U" (업종)
}

// InquireIndexTickprice 는 국내업종 틱별 지수 호출.
//
// 한투 docs: docs/api/국내주식/국내업종_틱별지수.md
// path: /uapi/domestic-stock/v1/quotations/inquire-index-tickprice (FHPUP02110100)
func (c *Client) InquireIndexTickprice(ctx context.Context, params InquireIndexTickpriceParams) (*IndexTickprice, error) {
	market := params.MarketCode
	if market == "" {
		market = "U"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-index-tickprice",
		TrID:   "FHPUP02110100",
		Query: map[string]string{
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_COND_MRKT_DIV_CODE": market,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res IndexTickprice
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse IndexTickprice: %w", err)
	}
	return &res, nil
}

// DailyIndexchartprice 는 국내업종 일봉 차트 (FHKUP03500100) 응답.
//
// 한투 docs: docs/api/국내주식/국내업종_일봉차트.md
// path: /uapi/domestic-stock/v1/quotations/inquire-daily-indexchartprice
//
// output1 에 futs_prdy_* (선물 전일 시가/고가/저가) 3 필드 포함 — 업종+선물 복합 스냅샷.
type DailyIndexchartprice struct {
	Output1 DailyIndexchartpriceSummary `json:"output1"`
	Output2 []DailyIndexchartpriceItem  `json:"output2"`
}

// DailyIndexchartpriceSummary 는 응답의 output1 (현재 스냅샷 + 선물 전일 OHLC, 15 fields).
type DailyIndexchartpriceSummary struct {
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`             // 전일 대비 부호
	BstpNmixPrdyCtrt float64         `json:"bstp_nmix_prdy_ctrt,string"` // 업종 지수 전일 대비율
	PrdyNmix         decimal.Decimal `json:"prdy_nmix"`                  // 전일 지수
	AcmlVol          int64           `json:"acml_vol,string"`            // 누적 거래량
	AcmlTrPbmn       int64           `json:"acml_tr_pbmn,string"`        // 누적 거래 대금
	HtsKorIsnm       string          `json:"hts_kor_isnm"`               // HTS 한글 종목명
	BstpNmixPrpr     decimal.Decimal `json:"bstp_nmix_prpr"`             // 업종 지수 현재가
	BstpClsCode      string          `json:"bstp_cls_code"`              // 업종 구분 코드
	PrdyVol          int64           `json:"prdy_vol,string"`            // 전일 거래량
	BstpNmixOprc     decimal.Decimal `json:"bstp_nmix_oprc"`             // 업종 지수 시가
	BstpNmixHgpr     decimal.Decimal `json:"bstp_nmix_hgpr"`             // 업종 지수 최고가
	BstpNmixLwpr     decimal.Decimal `json:"bstp_nmix_lwpr"`             // 업종 지수 최저가
	FutsPrdyOprc     decimal.Decimal `json:"futs_prdy_oprc"`             // 선물 전일 시가
	FutsPrdyHgpr     decimal.Decimal `json:"futs_prdy_hgpr"`             // 선물 전일 고가
	FutsPrdyLwpr     decimal.Decimal `json:"futs_prdy_lwpr"`             // 선물 전일 저가
}

// DailyIndexchartpriceItem 은 응답의 output2 한 행 (일봉, 8 fields).
type DailyIndexchartpriceItem struct {
	StckBsopDate string          `json:"stck_bsop_date"`      // 영업 일자
	BstpNmixPrpr decimal.Decimal `json:"bstp_nmix_prpr"`      // 업종 지수 현재가
	BstpNmixOprc decimal.Decimal `json:"bstp_nmix_oprc"`      // 업종 지수 시가
	BstpNmixHgpr decimal.Decimal `json:"bstp_nmix_hgpr"`      // 업종 지수 최고가
	BstpNmixLwpr decimal.Decimal `json:"bstp_nmix_lwpr"`      // 업종 지수 최저가
	AcmlVol      int64           `json:"acml_vol,string"`     // 누적 거래량
	AcmlTrPbmn   int64           `json:"acml_tr_pbmn,string"` // 누적 거래 대금
	ModYn        string          `json:"mod_yn"`              // 수정 여부
}

// InquireDailyIndexchartpriceParams 는 국내업종 일봉 차트 조회 파라미터.
type InquireDailyIndexchartpriceParams struct {
	MarketCode    string // FID_COND_MRKT_DIV_CODE — 빈 값=>"U" (업종)
	Symbol        string // FID_INPUT_ISCD — 필수, 업종 코드
	InputDate1    string // FID_INPUT_DATE_1 — 조회 시작일 YYYYMMDD
	InputDate2    string // FID_INPUT_DATE_2 — 조회 종료일 YYYYMMDD
	PeriodDivCode string // FID_PERIOD_DIV_CODE — D:일 W:주 M:월 Y:년
}

// InquireDailyIndexchartprice 는 국내업종 일봉 차트 호출.
//
// 한투 docs: docs/api/국내주식/국내업종_일봉차트.md
// path: /uapi/domestic-stock/v1/quotations/inquire-daily-indexchartprice (FHKUP03500100)
func (c *Client) InquireDailyIndexchartprice(ctx context.Context, params InquireDailyIndexchartpriceParams) (*DailyIndexchartprice, error) {
	market := params.MarketCode
	if market == "" {
		market = "U"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-daily-indexchartprice",
		TrID:   "FHKUP03500100",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_INPUT_DATE_1":       params.InputDate1,
			"FID_INPUT_DATE_2":       params.InputDate2,
			"FID_PERIOD_DIV_CODE":    params.PeriodDivCode,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res DailyIndexchartprice
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse DailyIndexchartprice: %w", err)
	}
	return &res, nil
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
