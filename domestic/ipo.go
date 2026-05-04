package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// PubOffer 는 예탁원정보(공모주청약일정) (HHKDB669108C0) 응답.
//
// 한투 docs: docs/api/국내주식/예탁원정보(공모주청약일정).md
// path: /uapi/domestic-stock/v1/ksdinfo/pub-offer
//
// IPO 청약일정 list. output1 (Array). 다른 ranking/financial 과 query 키 형식이 다름 (대문자+한글식).
type PubOffer struct {
	Output1 []PubOfferItem `json:"output1"`
}

// PubOfferItem 은 한 IPO 청약 일정 항목.
type PubOfferItem struct {
	RecordDate   string          `json:"record_date"`           // 기준일 (YYYYMMDD)
	ShtCd        string          `json:"sht_cd"`                // 종목코드
	IsinName     string          `json:"isin_name"`             // 종목명
	FixSubscrPri decimal.Decimal `json:"fix_subscr_pri"`        // 공모가
	FaceValue    decimal.Decimal `json:"face_value"`            // 액면가
	SubscrDt     string          `json:"subscr_dt"`             // 청약기간 (예 "20260505 ~ 20260506")
	PayDt        string          `json:"pay_dt"`                // 납입일 (YYYYMMDD)
	RefundDt     string          `json:"refund_dt"`             // 환불일
	ListDt       string          `json:"list_dt"`               // 상장/등록일
	LeadMgr      string          `json:"lead_mgr"`              // 주간사
	PubBfCap     int64           `json:"pub_bf_cap,string"`     // 공모전 자본금
	PubAfCap     int64           `json:"pub_af_cap,string"`     // 공모후 자본금
	AssignStkQty int64           `json:"assign_stk_qty,string"` // 당사 배정물량
}

// InquirePubOfferParams 는 공모주청약일정 조회 파라미터.
//
// 다른 ranking 과 query 키 형식이 다름 — KIS docs 그대로 노출 (SHT_CD, CTS, F_DT, T_DT).
type InquirePubOfferParams struct {
	Symbol   string // SHT_CD — 종목코드. 빈 값(공백) = 전체
	Cts      string // CTS — 빈 값(공백) default
	FromDate string // F_DT — 조회일자 From (YYYYMMDD)
	ToDate   string // T_DT — 조회일자 To (YYYYMMDD)
}

// InquirePubOffer 는 예탁원정보(공모주청약일정) 호출.
//
// 한투 docs: docs/api/국내주식/예탁원정보(공모주청약일정).md
// path: /uapi/domestic-stock/v1/ksdinfo/pub-offer (HHKDB669108C0)
//
// 공모주(IPO) 청약일정 list 조회. Symbol 빈 값 시 전체.
func (c *Client) InquirePubOffer(ctx context.Context, params InquirePubOfferParams) (*PubOffer, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ksdinfo/pub-offer",
		TrID:   "HHKDB669108C0",
		Query: map[string]string{
			"SHT_CD": params.Symbol,
			"CTS":    params.Cts,
			"F_DT":   params.FromDate,
			"T_DT":   params.ToDate,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res PubOffer
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse PubOffer: %w", err)
	}
	return &res, nil
}
