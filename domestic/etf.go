package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// EtfPrice 는 ETF/ETN 현재가 (FHPST02400000) 응답.
//
// 한투 docs: docs/api/국내주식/ETF_ETN현재가.md
// path: /uapi/etfetn/v1/quotations/inquire-price
type EtfPrice struct {
	Output EtfPriceData `json:"output"`
}

// EtfPriceData 는 ETF/ETN 현재가 응답의 output object (54 fields).
type EtfPriceData struct {
	StckPrpr               decimal.Decimal `json:"stck_prpr"`                       // 주식 현재가
	PrdyVrssSign           string          `json:"prdy_vrss_sign"`                  // 전일 대비 부호
	PrdyVrss               decimal.Decimal `json:"prdy_vrss"`                       // 전일 대비
	PrdyCtrt               float64         `json:"prdy_ctrt,string"`                // 전일 대비율
	AcmlVol                int64           `json:"acml_vol,string"`                 // 누적 거래량
	PrdyVol                int64           `json:"prdy_vol,string"`                 // 전일 거래량
	StckMxpr               decimal.Decimal `json:"stck_mxpr"`                       // 주식 상한가
	StckLlam               decimal.Decimal `json:"stck_llam"`                       // 주식 하한가
	StckPrdyClpr           decimal.Decimal `json:"stck_prdy_clpr"`                  // 주식 전일 종가
	StckOprc               decimal.Decimal `json:"stck_oprc"`                       // 주식 시가
	PrdyClprVrssOprcRate   float64         `json:"prdy_clpr_vrss_oprc_rate,string"` // 전일 종가 대비 시가 비율
	StckHgpr               decimal.Decimal `json:"stck_hgpr"`                       // 주식 최고가
	PrdyClprVrssHgprRate   float64         `json:"prdy_clpr_vrss_hgpr_rate,string"` // 전일 종가 대비 최고가 비율
	StckLwpr               decimal.Decimal `json:"stck_lwpr"`                       // 주식 최저가
	PrdyClprVrssLwprRate   float64         `json:"prdy_clpr_vrss_lwpr_rate,string"` // 전일 종가 대비 최저가 비율
	PrdyLastNav            decimal.Decimal `json:"prdy_last_nav"`                   // 전일 최종 NAV
	Nav                    decimal.Decimal `json:"nav"`                             // NAV
	NavPrdyVrss            decimal.Decimal `json:"nav_prdy_vrss"`                   // NAV 전일 대비
	NavPrdyVrssSign        string          `json:"nav_prdy_vrss_sign"`              // NAV 전일 대비 부호
	NavPrdyCtrt            float64         `json:"nav_prdy_ctrt,string"`            // NAV 전일 대비율
	TrcErrt                float64         `json:"trc_errt,string"`                 // 추적 오차율
	StckSdpr               decimal.Decimal `json:"stck_sdpr"`                       // 주식 기준가
	StckSspr               decimal.Decimal `json:"stck_sspr"`                       // 주식 예상 체결가
	NmixCtrt               float64         `json:"nmix_ctrt,string"`                // 지수 대비율
	EtfCrclStcn            int64           `json:"etf_crcl_stcn,string"`            // ETF 유통 주수
	EtfNtasTtam            int64           `json:"etf_ntas_ttam,string"`            // ETF 순자산 총액
	EtfFrcrNtasTtam        int64           `json:"etf_frcr_ntas_ttam,string"`       // ETF 외화 순자산 총액
	FrgnLimtRate           float64         `json:"frgn_limt_rate,string"`           // 외국인 한도율
	FrgnOderAbleQty        int64           `json:"frgn_oder_able_qty,string"`       // 외국인 주문 가능 수량
	EtfCuUnitScrtCnt       int64           `json:"etf_cu_unit_scrt_cnt,string"`     // ETF CU 단위 증권 수
	EtfCnfgIssuCnt         int64           `json:"etf_cnfg_issu_cnt,string"`        // ETF 구성 발행 수
	EtfDvdnCycl            string          `json:"etf_dvdn_cycl"`                   // ETF 배당 주기
	Crcd                   string          `json:"crcd"`                            // 통화 코드
	EtfCrclNtasTtam        int64           `json:"etf_crcl_ntas_ttam,string"`       // ETF 유통 순자산 총액
	EtfFrcrCrclNtasTtam    int64           `json:"etf_frcr_crcl_ntas_ttam,string"`  // ETF 외화 유통 순자산 총액
	EtfFrcrLastNtasWrthVal decimal.Decimal `json:"etf_frcr_last_ntas_wrth_val"`     // ETF 외화 최종 순자산 가치
	LpOderAbleClsCode      string          `json:"lp_oder_able_cls_code"`           // LP 주문 가능 구분 코드
	StckDryyHgpr           decimal.Decimal `json:"stck_dryy_hgpr"`                  // 주식 연중 최고가
	DryyHgprVrssProrate    float64         `json:"dryy_hgpr_vrss_prpr_rate,string"` // 연중 최고가 대비율
	DryyHgprDate           string          `json:"dryy_hgpr_date"`                  // 연중 최고가 일자
	StckDryyLwpr           decimal.Decimal `json:"stck_dryy_lwpr"`                  // 주식 연중 최저가
	DryyLwprVrssProrate    float64         `json:"dryy_lwpr_vrss_prpr_rate,string"` // 연중 최저가 대비율
	DryyLwprDate           string          `json:"dryy_lwpr_date"`                  // 연중 최저가 일자
	BstpKorIsnm            string          `json:"bstp_kor_isnm"`                   // 업종 한글 종목명
	ViClsCode              string          `json:"vi_cls_code"`                     // VI 구분 코드
	LstnStcn               int64           `json:"lstn_stcn,string"`                // 상장 주수
	FrgnHldnQty            int64           `json:"frgn_hldn_qty,string"`            // 외국인 보유 수량
	FrgnHldnQtyRate        float64         `json:"frgn_hldn_qty_rate,string"`       // 외국인 보유 수량 비율
	EtfTrcErtMltp          float64         `json:"etf_trc_ert_mltp,string"`         // ETF 추적 오차율 배수
	Dprt                   float64         `json:"dprt,string"`                     // 괴리율
	MbcrName               string          `json:"mbcr_name"`                       // 운용사 명
	StckLstnDate           string          `json:"stck_lstn_date"`                  // 주식 상장 일자
	MtrtDate               string          `json:"mtrt_date"`                       // 만기 일자
	ShrgTypeCode           string          `json:"shrg_type_code"`                  // 공유 유형 코드
	LpHldnRate             float64         `json:"lp_hldn_rate,string"`             // LP 보유 비율
	EtfTrgtNmixBstpCode    string          `json:"etf_trgt_nmix_bstp_code"`         // ETF 목표 지수 업종 코드
	EtfDivName             string          `json:"etf_div_name"`                    // ETF 분류 명
	EtfRprsBstpKorIsnm     string          `json:"etf_rprs_bstp_kor_isnm"`          // ETF 대표 업종 한글 종목명
	LpHldnVol              int64           `json:"lp_hldn_vol,string"`              // LP 보유 거래량
}

// InquireEtfPriceParams 는 ETF/ETN 현재가 조회 파라미터.
type InquireEtfPriceParams struct {
	Symbol     string // fid_input_iscd — 종목코드 (예 "069500") Y
	MarketCode string // fid_cond_mrkt_div_code — "J":KRX. 빈 값=>"J"
}

// InquireEtfPrice 는 ETF/ETN 현재가 호출.
//
// 한투 docs: docs/api/국내주식/ETF_ETN현재가.md
// path: /uapi/etfetn/v1/quotations/inquire-price (FHPST02400000)
//
// NOTE: 메서드명은 InquireEtfPrice (InquirePrice 와 충돌 방지).
func (c *Client) InquireEtfPrice(ctx context.Context, params InquireEtfPriceParams) (*EtfPrice, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/etfetn/v1/quotations/inquire-price",
		TrID:   "FHPST02400000",
		Query: map[string]string{
			"fid_cond_mrkt_div_code": market,
			"fid_input_iscd":         params.Symbol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res EtfPrice
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse EtfPrice: %w", err)
	}
	return &res, nil
}

// ComponentStockPrice 는 ETF 구성종목 시세 (FHKST121600C0) 응답.
//
// 한투 docs: docs/api/국내주식/ETF_구성종목시세.md
// path: /uapi/etfetn/v1/quotations/inquire-component-stock-price
type ComponentStockPrice struct {
	Output1 ComponentStockPriceSummary `json:"output1"`
	Output2 []ComponentStockPriceItem  `json:"output2"`
}

// ComponentStockPriceSummary 는 ETF 구성종목 시세 응답 output1 (단일 객체 16 fields).
type ComponentStockPriceSummary struct {
	StckPrpr         decimal.Decimal `json:"stck_prpr"`                   // 주식 현재가
	PrdyVrss         decimal.Decimal `json:"prdy_vrss"`                   // 전일 대비
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`              // 전일 대비 부호
	PrdyCtrt         float64         `json:"prdy_ctrt,string"`            // 전일 대비율
	EtfCnfgIssuAvls  int64           `json:"etf_cnfg_issu_avls,string"`   // ETF 구성 발행 시가총액
	Nav              decimal.Decimal `json:"nav"`                         // NAV
	NavPrdyVrssSign  string          `json:"nav_prdy_vrss_sign"`          // NAV 전일 대비 부호
	NavPrdyVrss      decimal.Decimal `json:"nav_prdy_vrss"`               // NAV 전일 대비
	NavPrdyCtrt      float64         `json:"nav_prdy_ctrt,string"`        // NAV 전일 대비율
	EtfNtasTtam      int64           `json:"etf_ntas_ttam,string"`        // ETF 순자산 총액
	PrdyClprNav      decimal.Decimal `json:"prdy_clpr_nav"`               // 전일 종가 NAV
	OprcNav          decimal.Decimal `json:"oprc_nav"`                    // 시가 NAV
	HprcNav          decimal.Decimal `json:"hprc_nav"`                    // 최고가 NAV
	LprcNav          decimal.Decimal `json:"lprc_nav"`                    // 최저가 NAV
	EtfCuUnitScrtCnt int64           `json:"etf_cu_unit_scrt_cnt,string"` // ETF CU 단위 증권 수
	EtfCnfgIssuCnt   int64           `json:"etf_cnfg_issu_cnt,string"`    // ETF 구성 발행 수
}

// ComponentStockPriceItem 은 ETF 구성종목 시세 응답 output2 의 한 행 (15 fields/item).
type ComponentStockPriceItem struct {
	StckShrnIscd    string          `json:"stck_shrn_iscd"`            // 주식 단축 종목코드
	HtsKorIsnm      string          `json:"hts_kor_isnm"`              // HTS 한글 종목명
	StckPrpr        decimal.Decimal `json:"stck_prpr"`                 // 주식 현재가
	PrdyVrss        decimal.Decimal `json:"prdy_vrss"`                 // 전일 대비
	PrdyVrssSign    string          `json:"prdy_vrss_sign"`            // 전일 대비 부호
	PrdyCtrt        float64         `json:"prdy_ctrt,string"`          // 전일 대비율
	AcmlVol         int64           `json:"acml_vol,string"`           // 누적 거래량
	AcmlTrPbmn      int64           `json:"acml_tr_pbmn,string"`       // 누적 거래 대금
	TdayRsflRate    float64         `json:"tday_rsfl_rate,string"`     // 당일 등락률
	PrdyVrssVol     int64           `json:"prdy_vrss_vol,string"`      // 전일 대비 거래량
	TrPbmnTnrt      float64         `json:"tr_pbmn_tnrt,string"`       // 거래 대금 회전율
	HtsAvls         int64           `json:"hts_avls,string"`           // HTS 시가총액
	EtfCnfgIssuAvls int64           `json:"etf_cnfg_issu_avls,string"` // ETF 구성 발행 시가총액
	EtfCnfgIssuRlim float64         `json:"etf_cnfg_issu_rlim,string"` // ETF 구성 발행 비율
	EtfVltnAmt      int64           `json:"etf_vltn_amt,string"`       // ETF 평가 금액
}

// InquireComponentStockPriceParams 는 ETF 구성종목 시세 조회 파라미터.
//
// FID_COND_SCR_DIV_CODE = "11216" 고정 (사용자 변경 불가).
type InquireComponentStockPriceParams struct {
	MarketCode     string // FID_COND_MRKT_DIV_CODE — "J":KRX. 빈 값=>"J"
	Symbol         string // FID_INPUT_ISCD — 종목코드 (예 "069500") Y
	CondScrDivCode string // FID_COND_SCR_DIV_CODE — "11216" 고정. 빈 값=>"11216"
}

// InquireComponentStockPrice 는 ETF 구성종목 시세 호출.
//
// 한투 docs: docs/api/국내주식/ETF_구성종목시세.md
// path: /uapi/etfetn/v1/quotations/inquire-component-stock-price (FHKST121600C0)
func (c *Client) InquireComponentStockPrice(ctx context.Context, params InquireComponentStockPriceParams) (*ComponentStockPrice, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "11216"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/etfetn/v1/quotations/inquire-component-stock-price",
		TrID:   "FHKST121600C0",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_COND_SCR_DIV_CODE":  scrDiv,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res ComponentStockPrice
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse ComponentStockPrice: %w", err)
	}
	return &res, nil
}

// NavComparisonTimeTrend 는 NAV 비교 시간 추이 (FHPST02440100) 응답.
//
// 한투 docs: docs/api/국내주식/NAV비교시간추이.md
// path: /uapi/etfetn/v1/quotations/nav-comparison-time-trend
type NavComparisonTimeTrend struct {
	Output []NavComparisonTimeTrendItem `json:"output"`
}

// NavComparisonTimeTrendItem 은 NAV 비교 시간 추이 응답의 한 행 (13 fields/item).
type NavComparisonTimeTrendItem struct {
	BsopHour        string          `json:"bsop_hour"`            // 영업 시간
	Nav             decimal.Decimal `json:"nav"`                  // NAV
	NavPrdyVrssSign string          `json:"nav_prdy_vrss_sign"`   // NAV 전일 대비 부호
	NavPrdyVrss     decimal.Decimal `json:"nav_prdy_vrss"`        // NAV 전일 대비
	NavPrdyCtrt     float64         `json:"nav_prdy_ctrt,string"` // NAV 전일 대비율
	NavVrssPrpr     decimal.Decimal `json:"nav_vrss_prpr"`        // NAV 대비 현재가
	Dprt            float64         `json:"dprt,string"`          // 괴리율
	StckPrpr        decimal.Decimal `json:"stck_prpr"`            // 주식 현재가
	PrdyVrss        decimal.Decimal `json:"prdy_vrss"`            // 전일 대비
	PrdyVrssSign    string          `json:"prdy_vrss_sign"`       // 전일 대비 부호
	PrdyCtrt        float64         `json:"prdy_ctrt,string"`     // 전일 대비율
	AcmlVol         int64           `json:"acml_vol,string"`      // 누적 거래량
	CntgVol         int64           `json:"cntg_vol,string"`      // 체결 거래량
}

// InquireNavComparisonTimeTrendParams 는 NAV 비교 시간 추이 조회 파라미터.
type InquireNavComparisonTimeTrendParams struct {
	HourClsCode string // fid_hour_cls_code — 시간 구분 코드 (예 "60") Y
	MarketCode  string // fid_cond_mrkt_div_code — "E":ETF. 빈 값=>"E"
	Symbol      string // fid_input_iscd — 종목코드 (예 "069500") Y
}

// InquireNavComparisonTimeTrend 는 NAV 비교 시간 추이 호출.
//
// 한투 docs: docs/api/국내주식/NAV비교시간추이.md
// path: /uapi/etfetn/v1/quotations/nav-comparison-time-trend (FHPST02440100)
func (c *Client) InquireNavComparisonTimeTrend(ctx context.Context, params InquireNavComparisonTimeTrendParams) (*NavComparisonTimeTrend, error) {
	market := params.MarketCode
	if market == "" {
		market = "E"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/etfetn/v1/quotations/nav-comparison-time-trend",
		TrID:   "FHPST02440100",
		Query: map[string]string{
			"fid_hour_cls_code":      params.HourClsCode,
			"fid_cond_mrkt_div_code": market,
			"fid_input_iscd":         params.Symbol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res NavComparisonTimeTrend
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse NavComparisonTimeTrend: %w", err)
	}
	return &res, nil
}
