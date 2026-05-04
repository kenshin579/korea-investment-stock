package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// ProductInfo 는 상품기본조회 (CTPF1604R) 의 output.
//
// 한투 docs: docs/api/국내주식/상품기본조회.md
// path: /uapi/domestic-stock/v1/quotations/search-info
//
// 다국가 종목 (KR/US/JP/HK/CN/VN) 의 가벼운 기본 정보. PDNO + PRDT_TYPE_CD 입력.
type ProductInfo struct {
	Pdno               string `json:"pdno"`
	PrdtTypeCd         string `json:"prdt_type_cd"`
	PrdtName           string `json:"prdt_name"`
	PrdtName120        string `json:"prdt_name120"`
	PrdtAbrvName       string `json:"prdt_abrv_name"`
	PrdtEngName        string `json:"prdt_eng_name"`
	PrdtEngName120     string `json:"prdt_eng_name120"`
	PrdtEngAbrvName    string `json:"prdt_eng_abrv_name"`
	StdPdno            string `json:"std_pdno"`
	ShtnPdno           string `json:"shtn_pdno"`
	PrdtSaleStatCd     string `json:"prdt_sale_stat_cd"`
	PrdtRiskGradeCd    string `json:"prdt_risk_grade_cd"`
	PrdtClsfCd         string `json:"prdt_clsf_cd"`
	PrdtClsfName       string `json:"prdt_clsf_name"`
	SaleStrtDt         string `json:"sale_strt_dt"`
	SaleEndDt          string `json:"sale_end_dt"`
	WrapAsstTypeCd     string `json:"wrap_asst_type_cd"`
	IvstPrdtTypeCd     string `json:"ivst_prdt_type_cd"`
	IvstPrdtTypeCdName string `json:"ivst_prdt_type_cd_name"`
	FrstErlmDt         string `json:"frst_erlm_dt"`
}

// SearchInfo 는 상품기본조회 호출.
//
// 한투 docs: docs/api/국내주식/상품기본조회.md
// path: /uapi/domestic-stock/v1/quotations/search-info (CTPF1604R)
//
// PRDT_TYPE_CD 예: "300"(국내 주식), "512"(US 나스닥), "513"(US 뉴욕), "529"(US 아멕스).
func (c *Client) SearchInfo(ctx context.Context, pdno, prdtTypeCD string) (*ProductInfo, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/search-info",
		TrID:   "CTPF1604R",
		Query: map[string]string{
			"PDNO":         pdno,
			"PRDT_TYPE_CD": prdtTypeCD,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var info ProductInfo
	if err := json.Unmarshal(resp.Output, &info); err != nil {
		return nil, fmt.Errorf("kis: parse ProductInfo: %w", err)
	}
	return &info, nil
}

// StockInfo 는 주식기본조회 (CTPF1002R) 의 output.
//
// 한투 docs: docs/api/국내주식/주식기본조회.md
// path: /uapi/domestic-stock/v1/quotations/search-stock-info
//
// 국내주식 종목의 상세 정보 (시장ID, 증권그룹, 상장정보, 업종분류, KOSPI200 편입 등).
// 한투 spec 자체는 다국가 PRDT_TYPE_CD 받지만 endpoint 는 국내 전용 데이터.
type StockInfo struct {
	Pdno                 string `json:"pdno"`
	PrdtTypeCd           string `json:"prdt_type_cd"`
	MketIdCd             string `json:"mket_id_cd"`
	ScrtGrpIdCd          string `json:"scrt_grp_id_cd"`
	ExcgDvsnCd           string `json:"excg_dvsn_cd"`
	SetlMmdd             string `json:"setl_mmdd"`
	LstgStqt             int64  `json:"lstg_stqt,string"`
	LstgCptlAmt          int64  `json:"lstg_cptl_amt,string"`
	Cpta                 int64  `json:"cpta,string"`
	Papr                 string `json:"papr"`
	IssuPric             string `json:"issu_pric"`
	Kospi200ItemYn       string `json:"kospi200_item_yn"`
	SctsMketLstgDt       string `json:"scts_mket_lstg_dt"`
	SctsMketLstgAbolDt   string `json:"scts_mket_lstg_abol_dt"`
	KosdaqMketLstgDt     string `json:"kosdaq_mket_lstg_dt"`
	KosdaqMketLstgAbolDt string `json:"kosdaq_mket_lstg_abol_dt"`
	FrbdMketLstgDt       string `json:"frbd_mket_lstg_dt"`
	FrbdMketLstgAbolDt   string `json:"frbd_mket_lstg_abol_dt"`
	ReitsKindCd          string `json:"reits_kind_cd"`
	EtfDvsnCd            string `json:"etf_dvsn_cd"`
	OilfFundYn           string `json:"oilf_fund_yn"`
	IdxBztpLclsCd        string `json:"idx_bztp_lcls_cd"`
	IdxBztpMclsCd        string `json:"idx_bztp_mcls_cd"`
	IdxBztpSclsCd        string `json:"idx_bztp_scls_cd"`
	IdxBztpLclsCdName    string `json:"idx_bztp_lcls_cd_name"`
	IdxBztpMclsCdName    string `json:"idx_bztp_mcls_cd_name"`
	IdxBztpSclsCdName    string `json:"idx_bztp_scls_cd_name"`
	StckKindCd           string `json:"stck_kind_cd"`
	MfndOpngDt           string `json:"mfnd_opng_dt"`
	MfndEndDt            string `json:"mfnd_end_dt"`
	DpsiErlmCnclDt       string `json:"dpsi_erlm_cncl_dt"`
	EtfCuQty             string `json:"etf_cu_qty"`
	PrdtName             string `json:"prdt_name"`
	PrdtName120          string `json:"prdt_name120"`
	PrdtAbrvName         string `json:"prdt_abrv_name"`
	StdPdno              string `json:"std_pdno"`
	PrdtEngName          string `json:"prdt_eng_name"`
	PrdtEngName120       string `json:"prdt_eng_name120"`
	PrdtEngAbrvName      string `json:"prdt_eng_abrv_name"`
	DpsiAptmErlmYn       string `json:"dpsi_aptm_erlm_yn"`
	EtfTxtnTypeCd        string `json:"etf_txtn_type_cd"`
	EtfTypeCd            string `json:"etf_type_cd"`
	LstgAbolDt           string `json:"lstg_abol_dt"`
	NwstOdstDvsnCd       string `json:"nwst_odst_dvsn_cd"`
	SbstPric             string `json:"sbst_pric"`
	ThcoSbstPric         string `json:"thco_sbst_pric"`
	ThcoSbstPricChngDt   string `json:"thco_sbst_pric_chng_dt"`
	TrStopYn             string `json:"tr_stop_yn"`
	AdmnItemYn           string `json:"admn_item_yn"`
	ThdtClpr             string `json:"thdt_clpr"`
	BfdyClpr             string `json:"bfdy_clpr"`
	ClprChngDt           string `json:"clpr_chng_dt"`
	StdIdstClsfCd        string `json:"std_idst_clsf_cd"`
	StdIdstClsfCdName    string `json:"std_idst_clsf_cd_name"`
	IdxBztpLclsCdEngName string `json:"idx_bztp_lcls_cd_eng_name"`
	IdxBztpMclsCdEngName string `json:"idx_bztp_mcls_cd_eng_name"`
	IdxBztpSclsCdEngName string `json:"idx_bztp_scls_cd_eng_name"`
}

// SearchStockInfo 는 주식기본조회 호출.
//
// 한투 docs: docs/api/국내주식/주식기본조회.md
// path: /uapi/domestic-stock/v1/quotations/search-stock-info (CTPF1002R)
//
// 한투 spec 충실: PDNO + PRDT_TYPE_CD 명시. Python 의 "KR" country_code 검사 없음.
func (c *Client) SearchStockInfo(ctx context.Context, pdno, prdtTypeCD string) (*StockInfo, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/search-stock-info",
		TrID:   "CTPF1002R",
		Query: map[string]string{
			"PDNO":         pdno,
			"PRDT_TYPE_CD": prdtTypeCD,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var info StockInfo
	if err := json.Unmarshal(resp.Output, &info); err != nil {
		return nil, fmt.Errorf("kis: parse StockInfo: %w", err)
	}
	return &info, nil
}
