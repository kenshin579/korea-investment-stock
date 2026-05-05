package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// AskingPriceExpCcn 은 주식현재가 호가/예상체결 (FHKST01010200) 응답.
//
// 한투 docs: docs/api/국내주식/주식현재가_호가_예상체결.md
// path: /uapi/domestic-stock/v1/quotations/inquire-asking-price-exp-ccn
type AskingPriceExpCcn struct {
	Output1 AskingPriceExpCcnOrderbook `json:"output1"`
	Output2 AskingPriceExpCcnExpected  `json:"output2"`
}

// AskingPriceExpCcnOrderbook 은 응답의 output1 — 10 단계 호가, 잔량, 증감.
type AskingPriceExpCcnOrderbook struct {
	AsprAcptHour string `json:"aspr_acpt_hour"` // 호가 접수 시간 (HHMMSS)

	Askp1  decimal.Decimal `json:"askp1"` // 매도호가 1
	Askp2  decimal.Decimal `json:"askp2"`
	Askp3  decimal.Decimal `json:"askp3"`
	Askp4  decimal.Decimal `json:"askp4"`
	Askp5  decimal.Decimal `json:"askp5"`
	Askp6  decimal.Decimal `json:"askp6"`
	Askp7  decimal.Decimal `json:"askp7"`
	Askp8  decimal.Decimal `json:"askp8"`
	Askp9  decimal.Decimal `json:"askp9"`
	Askp10 decimal.Decimal `json:"askp10"`

	Bidp1  decimal.Decimal `json:"bidp1"` // 매수호가 1
	Bidp2  decimal.Decimal `json:"bidp2"`
	Bidp3  decimal.Decimal `json:"bidp3"`
	Bidp4  decimal.Decimal `json:"bidp4"`
	Bidp5  decimal.Decimal `json:"bidp5"`
	Bidp6  decimal.Decimal `json:"bidp6"`
	Bidp7  decimal.Decimal `json:"bidp7"`
	Bidp8  decimal.Decimal `json:"bidp8"`
	Bidp9  decimal.Decimal `json:"bidp9"`
	Bidp10 decimal.Decimal `json:"bidp10"`

	AskpRsqn1  int64 `json:"askp_rsqn1,string"` // 매도호가 잔량 1
	AskpRsqn2  int64 `json:"askp_rsqn2,string"`
	AskpRsqn3  int64 `json:"askp_rsqn3,string"`
	AskpRsqn4  int64 `json:"askp_rsqn4,string"`
	AskpRsqn5  int64 `json:"askp_rsqn5,string"`
	AskpRsqn6  int64 `json:"askp_rsqn6,string"`
	AskpRsqn7  int64 `json:"askp_rsqn7,string"`
	AskpRsqn8  int64 `json:"askp_rsqn8,string"`
	AskpRsqn9  int64 `json:"askp_rsqn9,string"`
	AskpRsqn10 int64 `json:"askp_rsqn10,string"`

	BidpRsqn1  int64 `json:"bidp_rsqn1,string"` // 매수호가 잔량 1
	BidpRsqn2  int64 `json:"bidp_rsqn2,string"`
	BidpRsqn3  int64 `json:"bidp_rsqn3,string"`
	BidpRsqn4  int64 `json:"bidp_rsqn4,string"`
	BidpRsqn5  int64 `json:"bidp_rsqn5,string"`
	BidpRsqn6  int64 `json:"bidp_rsqn6,string"`
	BidpRsqn7  int64 `json:"bidp_rsqn7,string"`
	BidpRsqn8  int64 `json:"bidp_rsqn8,string"`
	BidpRsqn9  int64 `json:"bidp_rsqn9,string"`
	BidpRsqn10 int64 `json:"bidp_rsqn10,string"`

	AskpRsqnIcdc1  int64 `json:"askp_rsqn_icdc1,string"` // 매도호가 잔량 증감 1
	AskpRsqnIcdc2  int64 `json:"askp_rsqn_icdc2,string"`
	AskpRsqnIcdc3  int64 `json:"askp_rsqn_icdc3,string"`
	AskpRsqnIcdc4  int64 `json:"askp_rsqn_icdc4,string"`
	AskpRsqnIcdc5  int64 `json:"askp_rsqn_icdc5,string"`
	AskpRsqnIcdc6  int64 `json:"askp_rsqn_icdc6,string"`
	AskpRsqnIcdc7  int64 `json:"askp_rsqn_icdc7,string"`
	AskpRsqnIcdc8  int64 `json:"askp_rsqn_icdc8,string"`
	AskpRsqnIcdc9  int64 `json:"askp_rsqn_icdc9,string"`
	AskpRsqnIcdc10 int64 `json:"askp_rsqn_icdc10,string"`

	BidpRsqnIcdc1  int64 `json:"bidp_rsqn_icdc1,string"`
	BidpRsqnIcdc2  int64 `json:"bidp_rsqn_icdc2,string"`
	BidpRsqnIcdc3  int64 `json:"bidp_rsqn_icdc3,string"`
	BidpRsqnIcdc4  int64 `json:"bidp_rsqn_icdc4,string"`
	BidpRsqnIcdc5  int64 `json:"bidp_rsqn_icdc5,string"`
	BidpRsqnIcdc6  int64 `json:"bidp_rsqn_icdc6,string"`
	BidpRsqnIcdc7  int64 `json:"bidp_rsqn_icdc7,string"`
	BidpRsqnIcdc8  int64 `json:"bidp_rsqn_icdc8,string"`
	BidpRsqnIcdc9  int64 `json:"bidp_rsqn_icdc9,string"`
	BidpRsqnIcdc10 int64 `json:"bidp_rsqn_icdc10,string"`

	TotalAskpRsqn     int64  `json:"total_askp_rsqn,string"`      // 총 매도호가 잔량
	TotalBidpRsqn     int64  `json:"total_bidp_rsqn,string"`      // 총 매수호가 잔량
	TotalAskpRsqnIcdc int64  `json:"total_askp_rsqn_icdc,string"` // 총 매도호가 잔량 증감
	TotalBidpRsqnIcdc int64  `json:"total_bidp_rsqn_icdc,string"` // 총 매수호가 잔량 증감
	OvtmTotalAskpIcdc int64  `json:"ovtm_total_askp_icdc,string"` // 시간외 총 매도호가 증감
	OvtmTotalBidpIcdc int64  `json:"ovtm_total_bidp_icdc,string"` // 시간외 총 매수호가 증감
	OvtmTotalAskpRsqn int64  `json:"ovtm_total_askp_rsqn,string"` // 시간외 총 매도호가 잔량
	OvtmTotalBidpRsqn int64  `json:"ovtm_total_bidp_rsqn,string"` // 시간외 총 매수호가 잔량
	NtbyAsprRsqn      int64  `json:"ntby_aspr_rsqn,string"`       // 순매수 호가 잔량
	NewMkopClsCode    string `json:"new_mkop_cls_code"`           // 신 장운영 구분 코드
}

// AskingPriceExpCcnExpected 은 응답의 output2 — 예상체결 + 시세.
type AskingPriceExpCcnExpected struct {
	AntcMkopClsCode  string          `json:"antc_mkop_cls_code"`         // 예상 장운영 구분 코드
	StckPrpr         decimal.Decimal `json:"stck_prpr"`                  // 주식 현재가
	StckOprc         decimal.Decimal `json:"stck_oprc"`                  // 주식 시가
	StckHgpr         decimal.Decimal `json:"stck_hgpr"`                  // 주식 최고가
	StckLwpr         decimal.Decimal `json:"stck_lwpr"`                  // 주식 최저가
	StckSdpr         decimal.Decimal `json:"stck_sdpr"`                  // 주식 기준가
	AntcCnpr         decimal.Decimal `json:"antc_cnpr"`                  // 예상 체결가
	AntcCntgVrssSign string          `json:"antc_cntg_vrss_sign"`        // 예상 체결 대비 부호
	AntcCntgVrss     decimal.Decimal `json:"antc_cntg_vrss"`             // 예상 체결 대비
	AntcCntgPrdyCtrt float64         `json:"antc_cntg_prdy_ctrt,string"` // 예상 체결 전일 대비율
	AntcVol          int64           `json:"antc_vol,string"`            // 예상 거래량
	StckShrnIscd     string          `json:"stck_shrn_iscd"`             // 주식 단축 종목코드
	ViClsCode        string          `json:"vi_cls_code"`                // VI 적용 구분 코드
}

// InquireAskingPriceExpCcnParams 는 호가/예상체결 조회 파라미터.
type InquireAskingPriceExpCcnParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE — "J":KRX/"NX":NXT/"UN":통합. 빈 값=>"J"
	Symbol     string // FID_INPUT_ISCD — 종목코드 (예 "005930")
}

// Ccnl 은 주식현재가 체결 (FHKST01010300) 응답.
//
// 한투 docs: docs/api/국내주식/주식현재가_체결.md
// path: /uapi/domestic-stock/v1/quotations/inquire-ccnl
//
// 최근 체결 list (~30건). 체결강도 (tday_rltv) 포함.
type Ccnl struct {
	Output []CcnlItem `json:"output"`
}

// CcnlItem 은 체결 한 건.
type CcnlItem struct {
	StckCntgHour string          `json:"stck_cntg_hour"`   // 체결 시간 (HHMMSS)
	StckPrpr     decimal.Decimal `json:"stck_prpr"`        // 현재가
	PrdyVrss     decimal.Decimal `json:"prdy_vrss"`        // 전일 대비
	PrdyVrssSign string          `json:"prdy_vrss_sign"`   // 전일 대비 부호
	CntgVol      int64           `json:"cntg_vol,string"`  // 체결 거래량
	TdayRltv     float64         `json:"tday_rltv,string"` // 당일 체결강도
	PrdyCtrt     float64         `json:"prdy_ctrt,string"` // 전일 대비율
}

// InquireCcnlParams 는 체결 조회 파라미터.
type InquireCcnlParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE — 빈 값=>"J"
	Symbol     string // FID_INPUT_ISCD
}

// InquireCcnl 은 주식현재가 체결 호출.
//
// 한투 docs: docs/api/국내주식/주식현재가_체결.md
// path: /uapi/domestic-stock/v1/quotations/inquire-ccnl (FHKST01010300)
func (c *Client) InquireCcnl(ctx context.Context, params InquireCcnlParams) (*Ccnl, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-ccnl",
		TrID:   "FHKST01010300",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_INPUT_ISCD":         params.Symbol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res Ccnl
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse Ccnl: %w", err)
	}
	return &res, nil
}

// InquireAskingPriceExpCcn 은 주식현재가 호가/예상체결 호출.
//
// 한투 docs: docs/api/국내주식/주식현재가_호가_예상체결.md
// path: /uapi/domestic-stock/v1/quotations/inquire-asking-price-exp-ccn (FHKST01010200)
//
// output1: 10 단계 호가/잔량/증감 + 시간외 + VI 등.
// output2: 예상체결가 + 시세.
func (c *Client) InquireAskingPriceExpCcn(ctx context.Context, params InquireAskingPriceExpCcnParams) (*AskingPriceExpCcn, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-asking-price-exp-ccn",
		TrID:   "FHKST01010200",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_INPUT_ISCD":         params.Symbol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res AskingPriceExpCcn
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse AskingPriceExpCcn: %w", err)
	}
	return &res, nil
}
