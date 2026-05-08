package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// CompInterest 는 금리 종합 (FHPST07020000) 응답.
//
// 한투 docs: docs/api/국내주식/금리_종합(국내채권_금리).md
// path: /uapi/domestic-stock/v1/quotations/comp-interest
//
// output1 (단일 객체): 대표 금리 metadata. output2 (배열): 개별 금리 항목.
// 11:30 이후 신규 데이터 수신 → 그 전 조회 시 변경 없음.
type CompInterest struct {
	Output1 CompInterestSummary `json:"output1"`
	Output2 []CompInterestItem  `json:"output2"`
}

// CompInterestSummary 는 output1 (단일 대표 항목).
type CompInterestSummary struct {
	BcdtCode         string          `json:"bcdt_code"`           // 자료코드
	HtsKorIsnm       string          `json:"hts_kor_isnm"`        // HTS 한글 종목명
	BondMnrtPrpr     decimal.Decimal `json:"bond_mnrt_prpr"`      // 채권금리 현재가
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`      // 전일대비 부호
	BondMnrtPrdyVrss decimal.Decimal `json:"bond_mnrt_prdy_vrss"` // 채권금리 전일대비
	PrdyCtrt         float64         `json:"prdy_ctrt,string"`    // 전일대비율
	StckBsopDate     string          `json:"stck_bsop_date"`      // 영업일자
}

// CompInterestItem 은 output2 (개별 금리 항목).
//
// output1 과 거의 동일하지만 prdy_ctrt 가 bstp_nmix_prdy_ctrt (업종지수전일대비율) 로 키 다름.
type CompInterestItem struct {
	BcdtCode         string          `json:"bcdt_code"`                  // 자료코드
	HtsKorIsnm       string          `json:"hts_kor_isnm"`               // HTS 한글 종목명
	BondMnrtPrpr     decimal.Decimal `json:"bond_mnrt_prpr"`             // 채권금리 현재가
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`             // 전일대비 부호
	BondMnrtPrdyVrss decimal.Decimal `json:"bond_mnrt_prdy_vrss"`        // 채권금리 전일대비
	BstpNmixPrdyCtrt float64         `json:"bstp_nmix_prdy_ctrt,string"` // 업종지수 전일대비율
	StckBsopDate     string          `json:"stck_bsop_date"`             // 영업일자
}

// InquireCompInterestParams 는 금리 종합 조회 파라미터.
//
// 4개 query 모두 hardcoded → struct 비움. 향후 DivCode 변경이 필요하면 옵션 필드 추가.
type InquireCompInterestParams struct{}

// InquireCompInterest 는 금리 종합 호출.
//
// 한투 docs: docs/api/국내주식/금리_종합(국내채권_금리).md
// path: /uapi/domestic-stock/v1/quotations/comp-interest (FHPST07020000)
//
// 4 query 모두 UPPERCASE + hardcoded:
//   - FID_COND_MRKT_DIV_CODE = "I" (Unique key)
//   - FID_COND_SCR_DIV_CODE  = "20702"
//   - FID_DIV_CLS_CODE       = "1" (해외금리지표)
//   - FID_DIV_CLS_CODE1      = ""  (공백=전체)
func (c *Client) InquireCompInterest(ctx context.Context, _ InquireCompInterestParams) (*CompInterest, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/comp-interest",
		TrID:   "FHPST07020000",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": "I",
			"FID_COND_SCR_DIV_CODE":  "20702",
			"FID_DIV_CLS_CODE":       "1",
			"FID_DIV_CLS_CODE1":      "",
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res CompInterest
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse CompInterest: %w", err)
	}
	return &res, nil
}
