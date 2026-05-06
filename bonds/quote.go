package bonds

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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
