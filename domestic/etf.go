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
