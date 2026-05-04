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
