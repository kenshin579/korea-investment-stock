package overseasfutures

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// ─── EP1: InvestorUnpdTrend ──────────────────────────────────────────────────

// InvestorUnpdTrendOutput1 는 미결제추이 레코드 카운트 (1 필드).
type InvestorUnpdTrendOutput1 struct {
	RowCnt string `json:"row_cnt"` // 응답레코드카운트
}

// InvestorUnpdTrendOutput2Item 는 투자자별 미결제 추이 항목 (17 필드).
type InvestorUnpdTrendOutput2Item struct {
	ProdIscd       string `json:"prod_iscd"`        // 상품 코드
	CftcIscd       string `json:"cftc_iscd"`        // CFTC 코드
	BsopDate       string `json:"bsop_date"`        // 일자
	BidpSpec       string `json:"bidp_spec"`        // 매수투기
	AskpSpec       string `json:"askp_spec"`        // 매도투기
	SpreadSpec     string `json:"spread_spec"`      // 스프레드투기
	BidpHedge      string `json:"bidp_hedge"`       // 매수헤지
	AskpHedge      string `json:"askp_hedge"`       // 매도헤지
	HtsOtstSmtn    string `json:"hts_otst_smtn"`    // 미결제합계
	BidpMissing    string `json:"bidp_missing"`     // 매수누락
	AskpMissing    string `json:"askp_missing"`     // 매도누락
	BidpSpecCust   string `json:"bidp_spec_cust"`   // 매수투기고객
	AskpSpecCust   string `json:"askp_spec_cust"`   // 매도투기고객
	SpreadSpecCust string `json:"spread_spec_cust"` // 스프레드투기고객
	BidpHedgeCust  string `json:"bidp_hedge_cust"`  // 매수헤지고객
	AskpHedgeCust  string `json:"askp_hedge_cust"`  // 매도헤지고객
	CustSmtn       string `json:"cust_smtn"`        // 고객합계
}

// InvestorUnpdTrendData 는 해외선물 미결제추이 응답 (output1 + output2[]).
type InvestorUnpdTrendData struct {
	Output1 InvestorUnpdTrendOutput1       `json:"output1"`
	Output2 []InvestorUnpdTrendOutput2Item `json:"output2"`
}

type investorUnpdTrendResponse struct {
	RtCd    string                         `json:"rt_cd"`
	MsgCd   string                         `json:"msg_cd"`
	Msg1    string                         `json:"msg1"`
	Output1 InvestorUnpdTrendOutput1       `json:"output1"`
	Output2 []InvestorUnpdTrendOutput2Item `json:"output2"`
}

// InvestorUnpdTrendParams 는 해외선물 미결제추이 조회 파라미터.
type InvestorUnpdTrendParams struct {
	ProdIscd  string // 상품 코드 (예: GE/ZB/ZF/NQ/ES 등)
	BsopDate  string // 기준일 (YYYYMMDD)
	UpmuGubun string // 구분 (0:수량, 1:증감)
	CtsKey    string // 연속조회키 (공백 기본)
}

// InvestorUnpdTrend 는 해외선물 미결제추이 조회 (HHDDB95030000).
//
// CFTC 자료 — 매주 토요일 업데이트 (기준일: 화요일, 발표: 금요일).
// 연속조회 미지원 (tr_cont 이용 불가).
//
// KIS API: GET /uapi/overseas-futureoption/v1/quotations/investor-unpd-trend
func (c *Client) InvestorUnpdTrend(ctx context.Context, params InvestorUnpdTrendParams) (*InvestorUnpdTrendData, error) {
	ctsKey := params.CtsKey
	if ctsKey == "" {
		ctsKey = " "
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/overseas-futureoption/v1/quotations/investor-unpd-trend",
		TrID:     "HHDDB95030000",
		CustType: "P",
		Query: map[string]string{
			"PROD_ISCD":  params.ProdIscd,
			"BSOP_DATE":  params.BsopDate,
			"UPMU_GUBUN": params.UpmuGubun,
			"CTS_KEY":    ctsKey,
		},
	})
	if err != nil {
		return nil, err
	}
	var res investorUnpdTrendResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InvestorUnpdTrend: %w", err)
	}
	return &InvestorUnpdTrendData{
		Output1: res.Output1,
		Output2: res.Output2,
	}, nil
}
