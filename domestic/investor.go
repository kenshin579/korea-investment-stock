package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// InvestorTradeByStockDaily 는 종목별 투자자매매동향(일별) (FHPTJ04160001) 응답.
//
// 한투 docs: docs/api/국내주식/종목별_투자자매매동향(일별).md
// path: /uapi/domestic-stock/v1/quotations/investor-trade-by-stock-daily
//
// output1 (요약) + output2 (일별 Array). 각 일별 행에 외국인/개인/기관 등
// 13개 투자자 type 의 매수/매도/순매수 수량 + 거래대금 (~95 필드).
type InvestorTradeByStockDaily struct {
	Output1 InvestorTradeByStockDailySummary `json:"output1"`
	Output2 []InvestorTradeByStockDailyItem  `json:"output2"`
}

// InvestorTradeByStockDailySummary 는 응답의 output1 (단일 객체, 요약).
type InvestorTradeByStockDailySummary struct {
	StckPrpr        decimal.Decimal `json:"stck_prpr"`          // 주식 현재가
	PrdyVrss        decimal.Decimal `json:"prdy_vrss"`          // 전일 대비
	PrdyVrssSign    string          `json:"prdy_vrss_sign"`     // 전일 대비 부호
	PrdyCtrt        float64         `json:"prdy_ctrt,string"`   // 전일 대비율
	AcmlVol         int64           `json:"acml_vol,string"`    // 누적 거래량
	PrdyVol         int64           `json:"prdy_vol,string"`    // 전일 거래량
	RprsMrktKorName string          `json:"rprs_mrkt_kor_name"` // 대표 시장 한글명
}

// InvestorTradeByStockDailyItem 은 응답의 output2 한 행 (한 일자).
//
// KIS docs 의 line 86~180+ 모든 필드 1:1 매핑. 13 투자자 type
// (외국인/개인/기관계/증권/투자신탁/사모펀드/은행/보험/종금/기금/기타/기타법인/기타단체)
// 각각 ntby_qty + seln_vol + shnu_vol + seln_tr_pbmn + shnu_tr_pbmn + ntby_tr_pbmn = 6 fields.
type InvestorTradeByStockDailyItem struct {
	// 일자 + 시세 (10 fields)
	StckBsopDate string          `json:"stck_bsop_date"`      // 주식 영업 일자
	StckClpr     decimal.Decimal `json:"stck_clpr"`           // 주식 종가
	PrdyVrss     decimal.Decimal `json:"prdy_vrss"`           // 전일 대비
	PrdyVrssSign string          `json:"prdy_vrss_sign"`      // 전일 대비 부호
	PrdyCtrt     float64         `json:"prdy_ctrt,string"`    // 전일 대비율
	AcmlVol      int64           `json:"acml_vol,string"`     // 누적 거래량 (주)
	AcmlTrPbmn   int64           `json:"acml_tr_pbmn,string"` // 누적 거래 대금 (백만원)
	StckOprc     decimal.Decimal `json:"stck_oprc"`           // 시가
	StckHgpr     decimal.Decimal `json:"stck_hgpr"`           // 최고가
	StckLwpr     decimal.Decimal `json:"stck_lwpr"`           // 최저가

	// 외국인 (10 fields)
	FrgnNtbyQty       int64 `json:"frgn_ntby_qty,string"`        // 외국인 순매수 수량
	FrgnRegNtbyQty    int64 `json:"frgn_reg_ntby_qty,string"`    // 외국인 등록 순매수 수량
	FrgnNregNtbyQty   int64 `json:"frgn_nreg_ntby_qty,string"`   // 외국인 비등록 순매수 수량
	FrgnRegNtbyPbmn   int64 `json:"frgn_reg_ntby_pbmn,string"`   // 외국인 등록 순매수 대금
	FrgnNtbyTrPbmn    int64 `json:"frgn_ntby_tr_pbmn,string"`    // 외국인 순매수 거래 대금
	FrgnNregNtbyPbmn  int64 `json:"frgn_nreg_ntby_pbmn,string"`  // 외국인 비등록 순매수 대금
	FrgnSelnVol       int64 `json:"frgn_seln_vol,string"`        // 외국인 매도 거래량
	FrgnShnuVol       int64 `json:"frgn_shnu_vol,string"`        // 외국인 매수 거래량
	FrgnSelnTrPbmn    int64 `json:"frgn_seln_tr_pbmn,string"`    // 외국인 매도 거래 대금
	FrgnShnuTrPbmn    int64 `json:"frgn_shnu_tr_pbmn,string"`    // 외국인 매수 거래 대금

	// 외국인 등록/비등록 매도매수 (8 fields)
	FrgnRegAskpQty   int64 `json:"frgn_reg_askp_qty,string"`    // 외국인 등록 매도 수량
	FrgnRegBidpQty   int64 `json:"frgn_reg_bidp_qty,string"`    // 외국인 등록 매수 수량
	FrgnRegAskpPbmn  int64 `json:"frgn_reg_askp_pbmn,string"`   // 외국인 등록 매도 대금
	FrgnRegBidpPbmn  int64 `json:"frgn_reg_bidp_pbmn,string"`   // 외국인 등록 매수 대금
	FrgnNregAskpQty  int64 `json:"frgn_nreg_askp_qty,string"`   // 외국인 비등록 매도 수량
	FrgnNregBidpQty  int64 `json:"frgn_nreg_bidp_qty,string"`   // 외국인 비등록 매수 수량
	FrgnNregAskpPbmn int64 `json:"frgn_nreg_askp_pbmn,string"`  // 외국인 비등록 매도 대금
	FrgnNregBidpPbmn int64 `json:"frgn_nreg_bidp_pbmn,string"`  // 외국인 비등록 매수 대금

	// 개인 (6 fields)
	PrsnNtbyQty    int64 `json:"prsn_ntby_qty,string"`      // 개인 순매수 수량
	PrsnNtbyTrPbmn int64 `json:"prsn_ntby_tr_pbmn,string"`  // 개인 순매수 거래 대금
	PrsnSelnVol    int64 `json:"prsn_seln_vol,string"`      // 개인 매도 거래량
	PrsnShnuVol    int64 `json:"prsn_shnu_vol,string"`      // 개인 매수 거래량
	PrsnSelnTrPbmn int64 `json:"prsn_seln_tr_pbmn,string"`  // 개인 매도 거래 대금
	PrsnShnuTrPbmn int64 `json:"prsn_shnu_tr_pbmn,string"`  // 개인 매수 거래 대금

	// 기관계 (6 fields)
	OrgnNtbyQty    int64 `json:"orgn_ntby_qty,string"`      // 기관계 순매수 수량
	OrgnNtbyTrPbmn int64 `json:"orgn_ntby_tr_pbmn,string"`  // 기관계 순매수 거래 대금
	OrgnSelnVol    int64 `json:"orgn_seln_vol,string"`      // 기관계 매도 거래량
	OrgnShnuVol    int64 `json:"orgn_shnu_vol,string"`      // 기관계 매수 거래량
	OrgnSelnTrPbmn int64 `json:"orgn_seln_tr_pbmn,string"`  // 기관계 매도 거래 대금
	OrgnShnuTrPbmn int64 `json:"orgn_shnu_tr_pbmn,string"`  // 기관계 매수 거래 대금

	// 증권 (6 fields)
	ScrtNtbyQty    int64 `json:"scrt_ntby_qty,string"`
	ScrtNtbyTrPbmn int64 `json:"scrt_ntby_tr_pbmn,string"`
	ScrtSelnVol    int64 `json:"scrt_seln_vol,string"`
	ScrtShnuVol    int64 `json:"scrt_shnu_vol,string"`
	ScrtSelnTrPbmn int64 `json:"scrt_seln_tr_pbmn,string"`
	ScrtShnuTrPbmn int64 `json:"scrt_shnu_tr_pbmn,string"`

	// 투자신탁 (6 fields)
	IvtrNtbyQty    int64 `json:"ivtr_ntby_qty,string"`
	IvtrNtbyTrPbmn int64 `json:"ivtr_ntby_tr_pbmn,string"`
	IvtrSelnVol    int64 `json:"ivtr_seln_vol,string"`
	IvtrShnuVol    int64 `json:"ivtr_shnu_vol,string"`
	IvtrSelnTrPbmn int64 `json:"ivtr_seln_tr_pbmn,string"`
	IvtrShnuTrPbmn int64 `json:"ivtr_shnu_tr_pbmn,string"`

	// 사모펀드 (6 fields, KIS docs 가 vol/qty 혼용 — vol 사용)
	PeFundNtbyVol    int64 `json:"pe_fund_ntby_vol,string"`
	PeFundNtbyTrPbmn int64 `json:"pe_fund_ntby_tr_pbmn,string"`
	PeFundSelnVol    int64 `json:"pe_fund_seln_vol,string"`
	PeFundShnuVol    int64 `json:"pe_fund_shnu_vol,string"`
	PeFundSelnTrPbmn int64 `json:"pe_fund_seln_tr_pbmn,string"`
	PeFundShnuTrPbmn int64 `json:"pe_fund_shnu_tr_pbmn,string"`

	// 은행 (6 fields)
	BankNtbyQty    int64 `json:"bank_ntby_qty,string"`
	BankNtbyTrPbmn int64 `json:"bank_ntby_tr_pbmn,string"`
	BankSelnVol    int64 `json:"bank_seln_vol,string"`
	BankShnuVol    int64 `json:"bank_shnu_vol,string"`
	BankSelnTrPbmn int64 `json:"bank_seln_tr_pbmn,string"`
	BankShnuTrPbmn int64 `json:"bank_shnu_tr_pbmn,string"`

	// 보험 (6 fields)
	InsuNtbyQty    int64 `json:"insu_ntby_qty,string"`
	InsuNtbyTrPbmn int64 `json:"insu_ntby_tr_pbmn,string"`
	InsuSelnVol    int64 `json:"insu_seln_vol,string"`
	InsuShnuVol    int64 `json:"insu_shnu_vol,string"`
	InsuSelnTrPbmn int64 `json:"insu_seln_tr_pbmn,string"`
	InsuShnuTrPbmn int64 `json:"insu_shnu_tr_pbmn,string"`

	// 종금 (6 fields)
	MrbnNtbyQty    int64 `json:"mrbn_ntby_qty,string"`
	MrbnNtbyTrPbmn int64 `json:"mrbn_ntby_tr_pbmn,string"`
	MrbnSelnVol    int64 `json:"mrbn_seln_vol,string"`
	MrbnShnuVol    int64 `json:"mrbn_shnu_vol,string"`
	MrbnSelnTrPbmn int64 `json:"mrbn_seln_tr_pbmn,string"`
	MrbnShnuTrPbmn int64 `json:"mrbn_shnu_tr_pbmn,string"`

	// 기금 (6 fields)
	FundNtbyQty    int64 `json:"fund_ntby_qty,string"`
	FundNtbyTrPbmn int64 `json:"fund_ntby_tr_pbmn,string"`
	FundSelnVol    int64 `json:"fund_seln_vol,string"`
	FundShnuVol    int64 `json:"fund_shnu_vol,string"`
	FundSelnTrPbmn int64 `json:"fund_seln_tr_pbmn,string"`
	FundShnuTrPbmn int64 `json:"fund_shnu_tr_pbmn,string"`

	// 기타 (6 fields)
	EtcNtbyQty    int64 `json:"etc_ntby_qty,string"`
	EtcNtbyTrPbmn int64 `json:"etc_ntby_tr_pbmn,string"`
	EtcSelnVol    int64 `json:"etc_seln_vol,string"`
	EtcShnuVol    int64 `json:"etc_shnu_vol,string"`
	EtcSelnTrPbmn int64 `json:"etc_seln_tr_pbmn,string"`
	EtcShnuTrPbmn int64 `json:"etc_shnu_tr_pbmn,string"`

	// 기타 법인 (2 fields, KIS docs 가 vol 사용 — Plan comment says 3 but KIS only has 2)
	EtcCorpNtbyVol    int64 `json:"etc_corp_ntby_vol,string"`
	EtcCorpNtbyTrPbmn int64 `json:"etc_corp_ntby_tr_pbmn,string"`

	// 기타 단체 (6 fields, vol)
	EtcOrgtNtbyVol    int64 `json:"etc_orgt_ntby_vol,string"`
	EtcOrgtNtbyTrPbmn int64 `json:"etc_orgt_ntby_tr_pbmn,string"`
	EtcOrgtSelnVol    int64 `json:"etc_orgt_seln_vol,string"`
	EtcOrgtShnuVol    int64 `json:"etc_orgt_shnu_vol,string"`
	EtcOrgtSelnTrPbmn int64 `json:"etc_orgt_seln_tr_pbmn,string"`
	EtcOrgtShnuTrPbmn int64 `json:"etc_orgt_shnu_tr_pbmn,string"`
}

// InquireInvestorTradeByStockDailyParams 는 종목별 투자자매매동향(일별) 조회 파라미터.
type InquireInvestorTradeByStockDailyParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE — "J":KRX, "NX":NXT, "UN":통합. 빈 값=>"J"
	Symbol     string // FID_INPUT_ISCD — 필수, 종목코드 (6자리)
	BaseDate   string // FID_INPUT_DATE_1 — 필수, YYYYMMDD (해당일 조회는 장 종료 후 가능)
	OrgAdjPrc  string // FID_ORG_ADJ_PRC — 빈 값(공란) default
	EtcClsCode string // FID_ETC_CLS_CODE — 빈 값(공란) default
}

// InquireInvestorTradeByStockDaily 는 종목별 투자자매매동향(일별) 호출.
//
// 한투 docs: docs/api/국내주식/종목별_투자자매매동향(일별).md
// path: /uapi/domestic-stock/v1/quotations/investor-trade-by-stock-daily (FHPTJ04160001)
func (c *Client) InquireInvestorTradeByStockDaily(ctx context.Context, params InquireInvestorTradeByStockDailyParams) (*InvestorTradeByStockDaily, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/investor-trade-by-stock-daily",
		TrID:   "FHPTJ04160001",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_INPUT_DATE_1":       params.BaseDate,
			"FID_ORG_ADJ_PRC":        params.OrgAdjPrc,
			"FID_ETC_CLS_CODE":       params.EtcClsCode,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res InvestorTradeByStockDaily
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InvestorTradeByStockDaily: %w", err)
	}
	return &res, nil
}

// InvestorDailyByMarket 은 시장별 투자자매매동향(일별) (FHPTJ04040000) 응답.
//
// 한투 docs: docs/api/국내주식/시장별_투자자매매동향(일별).md
// path: /uapi/domestic-stock/v1/quotations/inquire-investor-daily-by-market
type InvestorDailyByMarket struct {
	Output []InvestorDailyByMarketItem `json:"output"`
}

// InvestorDailyByMarketItem 은 응답 한 행 (한 일자).
type InvestorDailyByMarketItem struct {
	// 일자 + 지수 (9 fields)
	StckBsopDate     string          `json:"stck_bsop_date"`             // 주식 영업 일자
	BstpNmixPrpr     decimal.Decimal `json:"bstp_nmix_prpr"`             // 업종 지수 현재가
	BstpNmixPrdyVrss decimal.Decimal `json:"bstp_nmix_prdy_vrss"`        // 업종 지수 전일 대비
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`             // 전일 대비 부호
	BstpNmixPrdyCtrt float64         `json:"bstp_nmix_prdy_ctrt,string"` // 업종 지수 전일 대비율
	BstpNmixOprc     decimal.Decimal `json:"bstp_nmix_oprc"`             // 업종 지수 시가
	BstpNmixHgpr     decimal.Decimal `json:"bstp_nmix_hgpr"`             // 업종 지수 최고가
	BstpNmixLwpr     decimal.Decimal `json:"bstp_nmix_lwpr"`             // 업종 지수 최저가
	StckPrdyClpr     decimal.Decimal `json:"stck_prdy_clpr"`             // 전일 종가

	// 13 type ntby_qty (수량)
	FrgnNtbyQty     int64 `json:"frgn_ntby_qty,string"`      // 외국인 순매수 수량
	FrgnRegNtbyQty  int64 `json:"frgn_reg_ntby_qty,string"`  // 외국인 등록 순매수 수량
	FrgnNregNtbyQty int64 `json:"frgn_nreg_ntby_qty,string"` // 외국인 비등록 순매수 수량
	PrsnNtbyQty     int64 `json:"prsn_ntby_qty,string"`      // 개인 순매수 수량
	OrgnNtbyQty     int64 `json:"orgn_ntby_qty,string"`      // 기관계 순매수 수량
	ScrtNtbyQty     int64 `json:"scrt_ntby_qty,string"`      // 증권 순매수 수량
	IvtrNtbyQty     int64 `json:"ivtr_ntby_qty,string"`      // 투자신탁 순매수 수량
	PeFundNtbyVol   int64 `json:"pe_fund_ntby_vol,string"`   // 사모 펀드 순매수 거래량 (vol)
	BankNtbyQty     int64 `json:"bank_ntby_qty,string"`      // 은행 순매수 수량
	InsuNtbyQty     int64 `json:"insu_ntby_qty,string"`      // 보험 순매수 수량
	MrbnNtbyQty     int64 `json:"mrbn_ntby_qty,string"`      // 종금 순매수 수량
	FundNtbyQty     int64 `json:"fund_ntby_qty,string"`      // 기금 순매수 수량
	EtcNtbyQty      int64 `json:"etc_ntby_qty,string"`       // 기타 순매수 수량
	EtcOrgtNtbyVol  int64 `json:"etc_orgt_ntby_vol,string"`  // 기타 단체 순매수 거래량
	EtcCorpNtbyVol  int64 `json:"etc_corp_ntby_vol,string"`  // 기타 법인 순매수 거래량

	// 14 type ntby_tr_pbmn (거래대금)
	FrgnNtbyTrPbmn    int64 `json:"frgn_ntby_tr_pbmn,string"`     // 외국인 순매수 거래 대금
	FrgnRegNtbyPbmn   int64 `json:"frgn_reg_ntby_pbmn,string"`    // 외국인 등록 순매수 대금
	FrgnNregNtbyPbmn  int64 `json:"frgn_nreg_ntby_pbmn,string"`   // 외국인 비등록 순매수 대금
	PrsnNtbyTrPbmn    int64 `json:"prsn_ntby_tr_pbmn,string"`     // 개인 순매수 거래 대금
	OrgnNtbyTrPbmn    int64 `json:"orgn_ntby_tr_pbmn,string"`     // 기관계 순매수 거래 대금
	ScrtNtbyTrPbmn    int64 `json:"scrt_ntby_tr_pbmn,string"`     // 증권
	IvtrNtbyTrPbmn    int64 `json:"ivtr_ntby_tr_pbmn,string"`     // 투자신탁
	PeFundNtbyTrPbmn  int64 `json:"pe_fund_ntby_tr_pbmn,string"`  // 사모 펀드
	BankNtbyTrPbmn    int64 `json:"bank_ntby_tr_pbmn,string"`     // 은행
	InsuNtbyTrPbmn    int64 `json:"insu_ntby_tr_pbmn,string"`     // 보험
	MrbnNtbyTrPbmn    int64 `json:"mrbn_ntby_tr_pbmn,string"`     // 종금
	FundNtbyTrPbmn    int64 `json:"fund_ntby_tr_pbmn,string"`     // 기금
	EtcNtbyTrPbmn     int64 `json:"etc_ntby_tr_pbmn,string"`      // 기타
	EtcOrgtNtbyTrPbmn int64 `json:"etc_orgt_ntby_tr_pbmn,string"` // 기타 단체
	EtcCorpNtbyTrPbmn int64 `json:"etc_corp_ntby_tr_pbmn,string"` // 기타 법인
}

// InquireInvestorDailyByMarketParams 는 시장별 투자자매매동향(일별) 조회 파라미터.
//
// KIS docs 의 query 키 그대로 노출. FID_INPUT_ISCD = 업종분류코드, FID_INPUT_ISCD_1 = 시장 (KSP/KSQ), FID_INPUT_ISCD_2 = 하위 분류.
type InquireInvestorDailyByMarketParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE — 빈 값=>"J"
	Symbol     string // FID_INPUT_ISCD — 업종분류코드 (예 "0001":코스피 종합)
	BaseDate   string // FID_INPUT_DATE_1 — YYYYMMDD
	Market     string // FID_INPUT_ISCD_1 — "KSP"(코스피) 또는 "KSQ"(코스닥)
	BaseDate2  string // FID_INPUT_DATE_2 — BaseDate 와 동일 일자
	SubCode    string // FID_INPUT_ISCD_2 — 하위 분류코드 (업종분류코드)
}

// InquireInvestorDailyByMarket 은 시장별 투자자매매동향(일별) 호출.
//
// 한투 docs: docs/api/국내주식/시장별_투자자매매동향(일별).md
// path: /uapi/domestic-stock/v1/quotations/inquire-investor-daily-by-market (FHPTJ04040000)
func (c *Client) InquireInvestorDailyByMarket(ctx context.Context, params InquireInvestorDailyByMarketParams) (*InvestorDailyByMarket, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-investor-daily-by-market",
		TrID:   "FHPTJ04040000",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_INPUT_DATE_1":       params.BaseDate,
			"FID_INPUT_ISCD_1":       params.Market,
			"FID_INPUT_DATE_2":       params.BaseDate2,
			"FID_INPUT_ISCD_2":       params.SubCode,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res InvestorDailyByMarket
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InvestorDailyByMarket: %w", err)
	}
	return &res, nil
}
