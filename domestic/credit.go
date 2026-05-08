package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// CreditByCompany 는 당사 신용가능종목 (FHPST04770000) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_당사_신용가능종목.md
// path: /uapi/domestic-stock/v1/quotations/credit-by-company
//
// HTS [0477] 당사 신용가능 종목. 최대 100 건. tr_cont 미지원.
type CreditByCompany struct {
	Output []CreditByCompanyItem `json:"output"`
}

// CreditByCompanyItem 은 당사 신용가능종목 응답의 한 행.
type CreditByCompanyItem struct {
	StckShrnIscd string  `json:"stck_shrn_iscd"`   // 단축 종목코드
	HtsKorIsnm   string  `json:"hts_kor_isnm"`     // HTS 한글 종목명
	CrdtRate     float64 `json:"crdt_rate,string"` // 신용 비율 (%)
}

// InquireCreditByCompanyParams 는 당사 신용가능종목 조회 파라미터.
//
// 2개 query (fid_cond_scr_div_code/fid_cond_mrkt_div_code) 는 hardcoded → struct 미노출.
type InquireCreditByCompanyParams struct {
	SortCode  string // fid_rank_sort_cls_code (0=코드순, 1=이름순). 빈 값=>"0"
	SelectYN  string // fid_slct_yn (0=신용주문가능, 1=신용주문불가). 빈 값=>"0"
	InputISCD string // fid_input_iscd (0000=전체, 0001=거래소, 1001=코스닥 등). 빈 값=>"0000"
}

// InquireCreditByCompany 는 당사 신용가능종목 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_당사_신용가능종목.md
// path: /uapi/domestic-stock/v1/quotations/credit-by-company (FHPST04770000)
//
// 2 hardcoded: fid_cond_scr_div_code="20477", fid_cond_mrkt_div_code="J".
func (c *Client) InquireCreditByCompany(ctx context.Context, params InquireCreditByCompanyParams) (*CreditByCompany, error) {
	sort := params.SortCode
	if sort == "" {
		sort = "0"
	}
	sel := params.SelectYN
	if sel == "" {
		sel = "0"
	}
	iscd := params.InputISCD
	if iscd == "" {
		iscd = "0000"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/credit-by-company",
		TrID:   "FHPST04770000",
		Query: map[string]string{
			"fid_rank_sort_cls_code": sort,
			"fid_slct_yn":            sel,
			"fid_input_iscd":         iscd,
			"fid_cond_scr_div_code":  "20477",
			"fid_cond_mrkt_div_code": "J",
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res CreditByCompany
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse CreditByCompany: %w", err)
	}
	return &res, nil
}
