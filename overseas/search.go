package overseas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// OverseasProductInfo 는 해외주식_상품기본정보 (CTPF1702R) 응답.
//
// 한투 docs: docs/api/해외주식/해외주식_상품기본정보.md
// path: /uapi/overseas-price/v1/quotations/search-info
//
// domestic.ProductInfo 와 다른 패키지/타입 — 해외 거래소 메타정보 풍부.
type OverseasProductInfo struct {
	Output OverseasProductInfoOutput `json:"output"`
}

// OverseasProductInfoOutput 은 응답의 output (단일 객체, ~52 fields).
type OverseasProductInfoOutput struct {
	StdPdno                string `json:"std_pdno"`
	PrdtEngName            string `json:"prdt_eng_name"`
	NatnCd                 string `json:"natn_cd"`
	NatnName               string `json:"natn_name"`
	TrMketCd               string `json:"tr_mket_cd"`
	TrMketName             string `json:"tr_mket_name"`
	OvrsExcgCd             string `json:"ovrs_excg_cd"`
	OvrsExcgName           string `json:"ovrs_excg_name"`
	TrCrcyCd               string `json:"tr_crcy_cd"`
	OvrsPapr               string `json:"ovrs_papr"`
	CrcyName               string `json:"crcy_name"`
	OvrsStckDvsnCd         string `json:"ovrs_stck_dvsn_cd"`
	PrdtClsfCd             string `json:"prdt_clsf_cd"`
	PrdtClsfName           string `json:"prdt_clsf_name"`
	SllUnitQty             string `json:"sll_unit_qty"`
	BuyUnitQty             string `json:"buy_unit_qty"`
	TrUnitAmt              string `json:"tr_unit_amt"`
	LstgStckNum            int64  `json:"lstg_stck_num,string"`
	LstgDt                 string `json:"lstg_dt"`
	OvrsStckTrStopDvsnCd   string `json:"ovrs_stck_tr_stop_dvsn_cd"`
	LstgAbolItemYn         string `json:"lstg_abol_item_yn"`
	OvrsStckPrdtGrpNo      string `json:"ovrs_stck_prdt_grp_no"`
	LstgYn                 string `json:"lstg_yn"`
	TaxLevyYn              string `json:"tax_levy_yn"`
	OvrsStckErlmRosnCd     string `json:"ovrs_stck_erlm_rosn_cd"`
	OvrsStckHistRghtDvsnCd string `json:"ovrs_stck_hist_rght_dvsn_cd"`
	ChngBfPdno             string `json:"chng_bf_pdno"`
	PrdtTypeCd2            string `json:"prdt_type_cd_2"`
	OvrsItemName           string `json:"ovrs_item_name"`
	SedolNo                string `json:"sedol_no"`
	BlbgTckrText           string `json:"blbg_tckr_text"`
	OvrsStckEtfRiskDrtpCd  string `json:"ovrs_stck_etf_risk_drtp_cd"`
	EtpChasErngRtDbnb      string `json:"etp_chas_erng_rt_dbnb"`
	IsttUsgeIsinCd         string `json:"istt_usge_isin_cd"`
	MintSvcYn              string `json:"mint_svc_yn"`
	MintSvcYnChngDt        string `json:"mint_svc_yn_chng_dt"`
	PrdtName               string `json:"prdt_name"`
	LeiCd                  string `json:"lei_cd"`
	OvrsStckStopRsonCd     string `json:"ovrs_stck_stop_rson_cd"`
	LstgAbolDt             string `json:"lstg_abol_dt"`
	MiniStkTrStatDvsnCd    string `json:"mini_stk_tr_stat_dvsn_cd"`
	MintFrstSvcErlmDt      string `json:"mint_frst_svc_erlm_dt"`
	MintDcptTradPsblYn     string `json:"mint_dcpt_trad_psbl_yn"`
	MintFnumTradPsblYn     string `json:"mint_fnum_trad_psbl_yn"`
	MintCblcCvsnIpsbYn     string `json:"mint_cblc_cvsn_ipsb_yn"`
	PtpItemYn              string `json:"ptp_item_yn"`
	PtpItemTrfxExmtYn      string `json:"ptp_item_trfx_exmt_yn"`
	PtpItemTrfxExmtStrtDt  string `json:"ptp_item_trfx_exmt_strt_dt"`
	PtpItemTrfxExmtEndDt   string `json:"ptp_item_trfx_exmt_end_dt"`
	DtmTrPsblYn            string `json:"dtm_tr_psbl_yn"`
	SdrfStopEclsYn         string `json:"sdrf_stop_ecls_yn"`
	SdrfStopEclsErlmDt     string `json:"sdrf_stop_ecls_erlm_dt"`
	MemoText1              string `json:"memo_text1"`
	OvrsNowPric1           string `json:"ovrs_now_pric1"`
	LastRcvgDtime          string `json:"last_rcvg_dtime"`
}

// SearchInfoParams 는 해외주식_상품기본정보 조회 파라미터.
type SearchInfoParams struct {
	PrdtTypeCD string // PRDT_TYPE_CD — 512(나스닥)/513(뉴욕)/529(아멕스)/515(일본)/501(홍콩)/543(홍콩CNY)/558(홍콩USD)/507(베트남 하노이)/508(베트남 호치민)/551(중국 상해A)/552(중국 심천A)
	Pdno       string // PDNO — 상품번호 (예 "AAPL")
}

// SearchInfo 는 해외주식_상품기본정보 호출.
//
// 한투 docs: docs/api/해외주식/해외주식_상품기본정보.md
// path: /uapi/overseas-price/v1/quotations/search-info (CTPF1702R)
func (c *Client) SearchInfo(ctx context.Context, params SearchInfoParams) (*OverseasProductInfo, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-price/v1/quotations/search-info",
		TrID:   "CTPF1702R",
		Query: map[string]string{
			"PRDT_TYPE_CD": params.PrdtTypeCD,
			"PDNO":         params.Pdno,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res OverseasProductInfo
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse OverseasProductInfo: %w", err)
	}
	return &res, nil
}
