package bonds

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// ─── EP1: SearchBondInfo ──────────────────────────────────────────────────────

// SearchBondInfoParams 는 채권 기본조회 요청 파라미터.
type SearchBondInfoParams struct {
	Pdno       string // PDNO: 채권 종목 번호 (필수)
	PrdtTypeCd string // PRDT_TYPE_CD: 상품유형코드 (필수, 예: "300")
}

// SearchBondInfoData 는 채권 기본조회 결과. CTPF1114R — all-string 70 fields.
type SearchBondInfoData struct {
	Pdno                   string `json:"pdno"`
	PrdtTypeCd             string `json:"prdt_type_cd"`
	KsdBondItemName        string `json:"ksd_bond_item_name"`
	KsdBondItemEngName     string `json:"ksd_bond_item_eng_name"`
	KsdBondLstgTypeCd      string `json:"ksd_bond_lstg_type_cd"`
	KsdOfrgDvsnCd          string `json:"ksd_ofrg_dvsn_cd"`
	KsdBondIntDfrmDvsnCd   string `json:"ksd_bond_int_dfrm_dvsn_cd"`
	IssuDt                 string `json:"issu_dt"`
	RdptDt                 string `json:"rdpt_dt"`
	RvnuDt                 string `json:"rvnu_dt"`
	IsoCrcyCd              string `json:"iso_crcy_cd"`
	MdwyRdptDt             string `json:"mdwy_rdpt_dt"`
	KsdRcvgBondDsctRt      string `json:"ksd_rcvg_bond_dsct_rt"`
	KsdRcvgBondSrfcInrt    string `json:"ksd_rcvg_bond_srfc_inrt"`
	BondExpdRdptRt         string `json:"bond_expd_rdpt_rt"`
	KsdPrcaRdptMthdCd      string `json:"ksd_prca_rdpt_mthd_cd"`
	IntCaltmMcnt           string `json:"int_caltm_mcnt"`
	KsdIntCalcUnitCd       string `json:"ksd_int_calc_unit_cd"`
	UvalCutDvsnCd          string `json:"uval_cut_dvsn_cd"`
	UvalCutDcptDgit        string `json:"uval_cut_dcpt_dgit"`
	KsdDydvCaltmAplyDvsnCd string `json:"ksd_dydv_caltm_aply_dvsn_cd"`
	DydvCalcDcnt           string `json:"dydv_calc_dcnt"`
	BondExpdAsrcErngRt     string `json:"bond_expd_asrc_erng_rt"`
	PadfPlacHdofName       string `json:"padf_plac_hdof_name"`
	LstgDt                 string `json:"lstg_dt"`
	LstgAbolDt             string `json:"lstg_abol_dt"`
	KsdBondIssuMthdCd      string `json:"ksd_bond_issu_mthd_cd"`
	LapsIndfYn             string `json:"laps_indf_yn"`
	KsdLhdyPniaDfrmMthdCd  string `json:"ksd_lhdy_pnia_dfrm_mthd_cd"`
	FrstIntDfrmDt          string `json:"frst_int_dfrm_dt"`
	KsdPrcmLnkgGvbdYn      string `json:"ksd_prcm_lnkg_gvbd_yn"`
	DpsiEndDt              string `json:"dpsi_end_dt"`
	DpsiStrtDt             string `json:"dpsi_strt_dt"`
	DpsiPsblYn             string `json:"dpsi_psbl_yn"`
	AtypRdptBondErlmYn     string `json:"atyp_rdpt_bond_erlm_yn"`
	DshnOccrYn             string `json:"dshn_occr_yn"`
	ExpdExtsYn             string `json:"expd_exts_yn"`
	PclrPtcrText           string `json:"pclr_ptcr_text"`
	DpsiPsblExcpStatCd     string `json:"dpsi_psbl_excp_stat_cd"`
	ExpdExtsSrdpRcnt       string `json:"expd_exts_srdp_rcnt"`
	ExpdExtsSrdpRt         string `json:"expd_exts_srdp_rt"`
	ExpdRdptRt             string `json:"expd_rdpt_rt"`
	ExpdAsrcErngRt         string `json:"expd_asrc_erng_rt"`
	BondIntDfrmMthdCd      string `json:"bond_int_dfrm_mthd_cd"`
	IntDfrmDayTypeCd       string `json:"int_dfrm_day_type_cd"`
	PrcaDfmtTermMcnt       string `json:"prca_dfmt_term_mcnt"`
	SpltRdptRcnt           string `json:"splt_rdpt_rcnt"`
	RgbfIntDfrmDt          string `json:"rgbf_int_dfrm_dt"`
	NxtmIntDfrmDt          string `json:"nxtm_int_dfrm_dt"`
	SprxPsblYn             string `json:"sprx_psbl_yn"`
	IctxRtDvsnCd           string `json:"ictx_rt_dvsn_cd"`
	BondClsfCd             string `json:"bond_clsf_cd"`
	BondClsfKorName        string `json:"bond_clsf_kor_name"`
	IntMnedDvsnCd          string `json:"int_mned_dvsn_cd"`
	PniaIntCalcUnpr        string `json:"pnia_int_calc_unpr"`
	FrnIntr                string `json:"frn_intr"`
	AplyDayPrcmIdxLnkgCefc string `json:"aply_day_prcm_idx_lnkg_cefc"`
	KsdExpdDydvCalcBassCd  string `json:"ksd_expd_dydv_calc_bass_cd"`
	ExpdDydvCalcDcnt       string `json:"expd_dydv_calc_dcnt"`
	KsdCbbwDvsnCd          string `json:"ksd_cbbw_dvsn_cd"`
	CrfdItemYn             string `json:"crfd_item_yn"`
	PniaBankOfdyDfrmMthdCd string `json:"pnia_bank_ofdy_dfrm_mthd_cd"`
	QibYn                  string `json:"qib_yn"`
	QibCclcDt              string `json:"qib_cclc_dt"`
	CsbdYn                 string `json:"csbd_yn"`
	CsbdCclcDt             string `json:"csbd_cclc_dt"`
	KsdOpcbYn              string `json:"ksd_opcb_yn"`
	KsdSodnYn              string `json:"ksd_sodn_yn"`
	KsdRqdiSctyYn          string `json:"ksd_rqdi_scty_yn"`
	ElecSctyYn             string `json:"elec_scty_yn"`
	RghtEcisMbdyDvsnCd     string `json:"rght_ecis_mbdy_dvsn_cd"`
	IntRkngMthdDvsnCd      string `json:"int_rkng_mthd_dvsn_cd"`
	OfrgDvsnCd             string `json:"ofrg_dvsn_cd"`
	KsdTotIssuAmt          string `json:"ksd_tot_issu_amt"`
}

// ─── EP2: InquireIssueInfo ────────────────────────────────────────────────────

// InquireIssueInfoParams 는 발행정보 조회 요청 파라미터.
type InquireIssueInfoParams struct {
	Pdno       string // PDNO: 채권 종목 번호 (필수)
	PrdtTypeCd string // PRDT_TYPE_CD: 상품유형코드 (필수)
}

// IssueInfo 는 채권 발행정보. CTPF1101R — all-string 69 fields.
type IssueInfo struct {
	Pdno                  string `json:"pdno"`
	PrdtTypeCd            string `json:"prdt_type_cd"`
	PrdtName              string `json:"prdt_name"`
	PrdtEngName           string `json:"prdt_eng_name"`
	IvstHeedPrdtYn        string `json:"ivst_heed_prdt_yn"`
	ExtsYn                string `json:"exts_yn"`
	BondClsfCd            string `json:"bond_clsf_cd"`
	BondClsfKorName       string `json:"bond_clsf_kor_name"`
	Papr                  string `json:"papr"`
	IntMnedDvsnCd         string `json:"int_mned_dvsn_cd"`
	RvnuShapCd            string `json:"rvnu_shap_cd"`
	IssuAmt               string `json:"issu_amt"`
	LstgRmnd              string `json:"lstg_rmnd"`
	IntDfrmMcnt           string `json:"int_dfrm_mcnt"`
	BondIntDfrmMthdCd     string `json:"bond_int_dfrm_mthd_cd"`
	SpltRdptRcnt          string `json:"splt_rdpt_rcnt"`
	PrcaDfmtTermMcnt      string `json:"prca_dfmt_term_mcnt"`
	IntAnapDvsnCd         string `json:"int_anap_dvsn_cd"`
	BondRghtDvsnCd        string `json:"bond_rght_dvsn_cd"`
	PrdtPclcText          string `json:"prdt_pclc_text"`
	PrdtAbrvName          string `json:"prdt_abrv_name"`
	PrdtEngAbrvName       string `json:"prdt_eng_abrv_name"`
	SprxPsblYn            string `json:"sprx_psbl_yn"`
	PbffPplcOfrgMthdCd    string `json:"pbff_pplc_ofrg_mthd_cd"`
	CmcoCd                string `json:"cmco_cd"`
	IssuIsttCd            string `json:"issu_istt_cd"`
	IssuIsttName          string `json:"issu_istt_name"`
	PniaDfrmAgcyIsttCd    string `json:"pnia_dfrm_agcy_istt_cd"`
	DsctEcRt              string `json:"dsct_ec_rt"`
	SrfcInrt              string `json:"srfc_inrt"`
	ExpdRdptRt            string `json:"expd_rdpt_rt"`
	ExpdAsrcErngRt        string `json:"expd_asrc_erng_rt"`
	BondGrteIsttName      string `json:"bond_grte_istt_name"`
	IntDfrmDayTypeCd      string `json:"int_dfrm_day_type_cd"`
	KsdIntCalcUnitCd      string `json:"ksd_int_calc_unit_cd"`
	IntWuntUderPrcsDvsnCd string `json:"int_wunt_uder_prcs_dvsn_cd"`
	RvnuDt                string `json:"rvnu_dt"`
	IssuDt                string `json:"issu_dt"`
	LstgDt                string `json:"lstg_dt"`
	ExpdDt                string `json:"expd_dt"`
	RdptDt                string `json:"rdpt_dt"`
	SbstPric              string `json:"sbst_pric"`
	RgbfIntDfrmDt         string `json:"rgbf_int_dfrm_dt"`
	NxtmIntDfrmDt         string `json:"nxtm_int_dfrm_dt"`
	FrstIntDfrmDt         string `json:"frst_int_dfrm_dt"`
	EcisPric              string `json:"ecis_pric"`
	RghtStckStdPdno       string `json:"rght_stck_std_pdno"`
	EcisOpngDt            string `json:"ecis_opng_dt"`
	EcisEndDt             string `json:"ecis_end_dt"`
	BondRvnuMthdCd        string `json:"bond_rvnu_mthd_cd"`
	OprtStfno             string `json:"oprt_stfno"`
	OprtStffName          string `json:"oprt_stff_name"`
	RgbfIntDfrmWday       string `json:"rgbf_int_dfrm_wday"`
	NxtmIntDfrmWday       string `json:"nxtm_int_dfrm_wday"`
	KisCrdtGradText       string `json:"kis_crdt_grad_text"`
	KbpCrdtGradText       string `json:"kbp_crdt_grad_text"`
	NiceCrdtGradText      string `json:"nice_crdt_grad_text"`
	FnpCrdtGradText       string `json:"fnp_crdt_grad_text"`
	DpsiPsblYn            string `json:"dpsi_psbl_yn"`
	PniaIntCalcUnpr       string `json:"pnia_int_calc_unpr"`
	PrcmIdxBondYn         string `json:"prcm_idx_bond_yn"`
	ExpdExtsSrdpRcnt      string `json:"expd_exts_srdp_rcnt"`
	ExpdExtsSrdpRt        string `json:"expd_exts_srdp_rt"`
	LoanPsblYn            string `json:"loan_psbl_yn"`
	GrteDvsnCd            string `json:"grte_dvsn_cd"`
	FnrrRankDvsnCd        string `json:"fnrr_rank_dvsn_cd"`
	KrxLstgAbolDvsnCd     string `json:"krx_lstg_abol_dvsn_cd"`
	AsstRqdiDvsnCd        string `json:"asst_rqdi_dvsn_cd"`
	OpcbDvsnCd            string `json:"opcb_dvsn_cd"`
	CrfdItemYn            string `json:"crfd_item_yn"`
	CrfdItemRstcCclcDt    string `json:"crfd_item_rstc_cclc_dt"`
	BondNmprUnitPric      string `json:"bond_nmpr_unit_pric"`
	IvstHeedBondDvsnName  string `json:"ivst_heed_bond_dvsn_name"`
	AddErngRt             string `json:"add_erng_rt"`
	AddErngRtAplyDt       string `json:"add_erng_rt_aply_dt"`
	BondTrStopDvsnCd      string `json:"bond_tr_stop_dvsn_cd"`
	IvstHeedBondDvsnCd    string `json:"ivst_heed_bond_dvsn_cd"`
	PclrCndtText          string `json:"pclr_cndt_text"`
	HbbdYn                string `json:"hbbd_yn"`
	CdtlCptlSctyTypeCd    string `json:"cdtl_cptl_scty_type_cd"`
	ElecSctyYn            string `json:"elec_scty_yn"`
	Sq1ClopEcisOpngDt     string `json:"sq1_clop_ecis_opng_dt"`
	FrstErlmStfno         string `json:"frst_erlm_stfno"`
	FrstErlmDt            string `json:"frst_erlm_dt"`
	FrstErlmTmd           string `json:"frst_erlm_tmd"`
	TlgRcvgDtlDtime       string `json:"tlg_rcvg_dtl_dtime"`
}

type inquireIssueInfoResponse struct {
	RtCd   string    `json:"rt_cd"`
	MsgCd  string    `json:"msg_cd"`
	Msg1   string    `json:"msg1"`
	Output IssueInfo `json:"output"`
}

// InquireIssueInfo 는 채권 발행정보 조회 (CTPF1101R).
//
// KIS API: GET /uapi/domestic-bond/v1/quotations/issue-info
func (c *Client) InquireIssueInfo(ctx context.Context, params InquireIssueInfoParams) (*IssueInfo, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-bond/v1/quotations/issue-info",
		TrID:     "CTPF1101R",
		CustType: "P",
		Query: map[string]string{
			"PDNO":         params.Pdno,
			"PRDT_TYPE_CD": params.PrdtTypeCd,
		},
	})
	if err != nil {
		return nil, err
	}
	var res inquireIssueInfoResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InquireIssueInfo: %w", err)
	}
	return &res.Output, nil
}

type searchBondInfoResponse struct {
	RtCd   string             `json:"rt_cd"`
	MsgCd  string             `json:"msg_cd"`
	Msg1   string             `json:"msg1"`
	Output SearchBondInfoData `json:"output"`
}

// SearchBondInfo 는 채권 기본조회 (CTPF1114R).
//
// KIS API: GET /uapi/domestic-bond/v1/quotations/search-bond-info
func (c *Client) SearchBondInfo(ctx context.Context, params SearchBondInfoParams) (*SearchBondInfoData, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-bond/v1/quotations/search-bond-info",
		TrID:     "CTPF1114R",
		CustType: "P",
		Query: map[string]string{
			"PDNO":         params.Pdno,
			"PRDT_TYPE_CD": params.PrdtTypeCd,
		},
	})
	if err != nil {
		return nil, err
	}
	var res searchBondInfoResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse SearchBondInfo: %w", err)
	}
	return &res.Output, nil
}

// ─── EP3: InquirePrice ────────────────────────────────────────────────────────

// InquirePriceParams 는 채권 현재가 시세 요청 파라미터.
type InquirePriceParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE: 기본 "B"
	Symbol     string // FID_INPUT_ISCD: 채권 단축 종목코드 (필수)
}

// BondPrice 는 채권 현재가 시세. FHKBJ773400C0 — 17 fields typed.
type BondPrice struct {
	StndIscd     string          `json:"stnd_iscd"`
	HtsKorIsnm   string          `json:"hts_kor_isnm"`
	BondPrpr     decimal.Decimal `json:"bond_prpr"`
	PrdyVrssSign string          `json:"prdy_vrss_sign"`
	BondPrdyVrss decimal.Decimal `json:"bond_prdy_vrss"`
	PrdyCtrt     float64         `json:"prdy_ctrt,string"`
	AcmlVol      int64           `json:"acml_vol,string"`
	BondPrdyClpr decimal.Decimal `json:"bond_prdy_clpr"`
	BondOprc     decimal.Decimal `json:"bond_oprc"`
	BondHgpr     decimal.Decimal `json:"bond_hgpr"`
	BondLwpr     decimal.Decimal `json:"bond_lwpr"`
	ErnnRate     float64         `json:"ernn_rate,string"`
	OprcErt      float64         `json:"oprc_ert,string"`
	HgprErt      float64         `json:"hgpr_ert,string"`
	LwprErt      float64         `json:"lwpr_ert,string"`
	BondMxpr     decimal.Decimal `json:"bond_mxpr"`
	BondLlam     decimal.Decimal `json:"bond_llam"`
}

type inquirePriceResponse struct {
	RtCd   string    `json:"rt_cd"`
	MsgCd  string    `json:"msg_cd"`
	Msg1   string    `json:"msg1"`
	Output BondPrice `json:"output"`
}

// InquirePrice 는 채권 현재가 시세 (FHKBJ773400C0).
//
// KIS API: GET /uapi/domestic-bond/v1/quotations/inquire-price
func (c *Client) InquirePrice(ctx context.Context, params InquirePriceParams) (*BondPrice, error) {
	if params.MarketCode == "" {
		params.MarketCode = "B"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-bond/v1/quotations/inquire-price",
		TrID:     "FHKBJ773400C0",
		CustType: "P",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": params.MarketCode,
			"FID_INPUT_ISCD":         params.Symbol,
		},
	})
	if err != nil {
		return nil, err
	}
	var res inquirePriceResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InquirePrice: %w", err)
	}
	return &res.Output, nil
}

// ─── EP4: InquireCcnl ─────────────────────────────────────────────────────────

// InquireCcnlParams 는 채권 현재가 체결 요청 파라미터.
type InquireCcnlParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE: 기본 "B"
	Symbol     string // FID_INPUT_ISCD: 채권 단축 종목코드 (필수)
}

// BondCcnl 는 채권 현재가 체결 (single snapshot). FHKBJ773403C0 — 7 fields typed.
type BondCcnl struct {
	StckCntgHour string          `json:"stck_cntg_hour"`
	BondPrpr     decimal.Decimal `json:"bond_prpr"`
	BondPrdyVrss decimal.Decimal `json:"bond_prdy_vrss"`
	PrdyVrssSign string          `json:"prdy_vrss_sign"`
	PrdyCtrt     float64         `json:"prdy_ctrt,string"`
	CntgVol      int64           `json:"cntg_vol,string"`
	AcmlVol      int64           `json:"acml_vol,string"`
}

type inquireCcnlResponse struct {
	RtCd   string   `json:"rt_cd"`
	MsgCd  string   `json:"msg_cd"`
	Msg1   string   `json:"msg1"`
	Output BondCcnl `json:"output"`
}

// InquireCcnl 는 채권 현재가 체결 (FHKBJ773403C0).
//
// KIS API: GET /uapi/domestic-bond/v1/quotations/inquire-ccnl
func (c *Client) InquireCcnl(ctx context.Context, params InquireCcnlParams) (*BondCcnl, error) {
	if params.MarketCode == "" {
		params.MarketCode = "B"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-bond/v1/quotations/inquire-ccnl",
		TrID:     "FHKBJ773403C0",
		CustType: "P",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": params.MarketCode,
			"FID_INPUT_ISCD":         params.Symbol,
		},
	})
	if err != nil {
		return nil, err
	}
	var res inquireCcnlResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InquireCcnl: %w", err)
	}
	return &res.Output, nil
}

// ─── EP5: InquireAskingPrice ──────────────────────────────────────────────────

// InquireAskingPriceParams 는 채권 현재가 호가 요청 파라미터.
type InquireAskingPriceParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE: 기본 "B"
	Symbol     string // FID_INPUT_ISCD: 채권 단축 종목코드 (필수)
}

// BondAskingPrice 는 채권 현재가 호가 (5단계). FHKBJ773401C0 — 34 fields typed.
type BondAskingPrice struct {
	AsprAcptHour  string          `json:"aspr_acpt_hour"`
	BondAskp1     decimal.Decimal `json:"bond_askp1"`
	BondAskp2     decimal.Decimal `json:"bond_askp2"`
	BondAskp3     decimal.Decimal `json:"bond_askp3"`
	BondAskp4     decimal.Decimal `json:"bond_askp4"`
	BondAskp5     decimal.Decimal `json:"bond_askp5"`
	BondBidp1     decimal.Decimal `json:"bond_bidp1"`
	BondBidp2     decimal.Decimal `json:"bond_bidp2"`
	BondBidp3     decimal.Decimal `json:"bond_bidp3"`
	BondBidp4     decimal.Decimal `json:"bond_bidp4"`
	BondBidp5     decimal.Decimal `json:"bond_bidp5"`
	AskpRsqn1     int64           `json:"askp_rsqn1,string"`
	AskpRsqn2     int64           `json:"askp_rsqn2,string"`
	AskpRsqn3     int64           `json:"askp_rsqn3,string"`
	AskpRsqn4     int64           `json:"askp_rsqn4,string"`
	AskpRsqn5     int64           `json:"askp_rsqn5,string"`
	BidpRsqn1     int64           `json:"bidp_rsqn1,string"`
	BidpRsqn2     int64           `json:"bidp_rsqn2,string"`
	BidpRsqn3     int64           `json:"bidp_rsqn3,string"`
	BidpRsqn4     int64           `json:"bidp_rsqn4,string"`
	BidpRsqn5     int64           `json:"bidp_rsqn5,string"`
	TotalAskpRsqn int64           `json:"total_askp_rsqn,string"`
	TotalBidpRsqn int64           `json:"total_bidp_rsqn,string"`
	NtbyAsprRsqn  int64           `json:"ntby_aspr_rsqn,string"`
	SelnErnnRate1 float64         `json:"seln_ernn_rate1,string"`
	SelnErnnRate2 float64         `json:"seln_ernn_rate2,string"`
	SelnErnnRate3 float64         `json:"seln_ernn_rate3,string"`
	SelnErnnRate4 float64         `json:"seln_ernn_rate4,string"`
	SelnErnnRate5 float64         `json:"seln_ernn_rate5,string"`
	ShnuErnnRate1 float64         `json:"shnu_ernn_rate1,string"`
	ShnuErnnRate2 float64         `json:"shnu_ernn_rate2,string"`
	ShnuErnnRate3 float64         `json:"shnu_ernn_rate3,string"`
	ShnuErnnRate4 float64         `json:"shnu_ernn_rate4,string"`
	ShnuErnnRate5 float64         `json:"shnu_ernn_rate5,string"`
}

type inquireAskingPriceResponse struct {
	RtCd   string          `json:"rt_cd"`
	MsgCd  string          `json:"msg_cd"`
	Msg1   string          `json:"msg1"`
	Output BondAskingPrice `json:"output"`
}

// InquireAskingPrice 는 채권 현재가 호가 5단계 (FHKBJ773401C0).
//
// KIS API: GET /uapi/domestic-bond/v1/quotations/inquire-asking-price
func (c *Client) InquireAskingPrice(ctx context.Context, params InquireAskingPriceParams) (*BondAskingPrice, error) {
	if params.MarketCode == "" {
		params.MarketCode = "B"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-bond/v1/quotations/inquire-asking-price",
		TrID:     "FHKBJ773401C0",
		CustType: "P",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": params.MarketCode,
			"FID_INPUT_ISCD":         params.Symbol,
		},
	})
	if err != nil {
		return nil, err
	}
	var res inquireAskingPriceResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InquireAskingPrice: %w", err)
	}
	return &res.Output, nil
}

// ─── EP6: InquireDailyPrice ───────────────────────────────────────────────────

// InquireDailyPriceParams 는 채권 현재가 일별 요청 파라미터.
type InquireDailyPriceParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE: 기본 "B"
	Symbol     string // FID_INPUT_ISCD: 채권 단축 종목코드 (필수)
}

// BondDailyPrice 는 채권 현재가 일별 시세. FHKBJ773404C0 — 9 fields typed.
//
// Note: KIS docs 는 output{} (object) 로 명시. 실제 API 가 array 를 반환하면 patch 에서 수정.
type BondDailyPrice struct {
	StckBsopDate string          `json:"stck_bsop_date"`
	BondPrpr     decimal.Decimal `json:"bond_prpr"`
	BondPrdyVrss decimal.Decimal `json:"bond_prdy_vrss"`
	PrdyVrssSign string          `json:"prdy_vrss_sign"`
	PrdyCtrt     float64         `json:"prdy_ctrt,string"`
	AcmlVol      int64           `json:"acml_vol,string"`
	BondOprc     decimal.Decimal `json:"bond_oprc"`
	BondHgpr     decimal.Decimal `json:"bond_hgpr"`
	BondLwpr     decimal.Decimal `json:"bond_lwpr"`
}

type inquireDailyPriceResponse struct {
	RtCd   string         `json:"rt_cd"`
	MsgCd  string         `json:"msg_cd"`
	Msg1   string         `json:"msg1"`
	Output BondDailyPrice `json:"output"`
}

// InquireDailyPrice 는 채권 현재가 일별 시세 (FHKBJ773404C0).
//
// KIS API: GET /uapi/domestic-bond/v1/quotations/inquire-daily-price
func (c *Client) InquireDailyPrice(ctx context.Context, params InquireDailyPriceParams) (*BondDailyPrice, error) {
	if params.MarketCode == "" {
		params.MarketCode = "B"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-bond/v1/quotations/inquire-daily-price",
		TrID:     "FHKBJ773404C0",
		CustType: "P",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": params.MarketCode,
			"FID_INPUT_ISCD":         params.Symbol,
		},
	})
	if err != nil {
		return nil, err
	}
	var res inquireDailyPriceResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InquireDailyPrice: %w", err)
	}
	return &res.Output, nil
}
