// Package domestic — market_op.go
// Phase 4.2: 시장운영/특수상태 4 메서드 (EP4~EP7)
//
// EP4  InquireExpClosingPrice    — 장마감 예상체결가   FHKST117300C0
// EP5  InquireChkHoliday         — 휴장일 조회         CTCA0903R
// EP6  InquireViStatus           — 변동성완화장치 현황 FHPST01390000
// EP7  InquireCaptureUplowprice  — 상하한가 포착       FHKST130000C0
//
// WebSocket 제외: 장운영정보 KRX(H0STMKO0) / NXT(H0NXMKO0) / 통합(H0UNMKO0) 은
// REST GET 이 아닌 WebSocket push API — Phase 4.2 범위 외. Phase 5 (WebSocket) 에서 처리 예정.
package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
	"github.com/shopspring/decimal"
)

// ─── EP4: InquireExpClosingPrice ────────────────────────────────────────────

// InquireExpClosingPriceParams 는 장마감 예상체결가 조회 파라미터.
// FID_INPUT_ISCD 는 종목코드가 아닌 시장 구분코드: 0000(전체)/0001(코스피)/1001(코스닥)/2001(코스피200)/4001(KRX100).
type InquireExpClosingPriceParams struct {
	RankSortClsCode string // FID_RANK_SORT_CLS_CODE: 0=전체/1=상한가마감/2=하한가마감/3=상승률상위/4=하락률상위
	MarketCode      string // FID_COND_MRKT_DIV_CODE: 기본 "J"
	CondScrDivCode  string // FID_COND_SCR_DIV_CODE: 기본 "11173" (hardcoded)
	Symbol          string // FID_INPUT_ISCD: 시장구분코드 0000/0001/1001/2001/4001
	BlngClsCode     string // FID_BLNG_CLS_CODE: 0=전체/1=종가범위연장
}

// ExpClosingPriceItem 은 장마감 예상체결가 종목별 데이터.
type ExpClosingPriceItem struct {
	StckShrnIscd     string          `json:"stck_shrn_iscd"`
	HtsKorIsnm       string          `json:"hts_kor_isnm"`
	StckPrpr         decimal.Decimal `json:"stck_prpr"`
	PrdyVrss         decimal.Decimal `json:"prdy_vrss"`
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`
	PrdyCtrt         float64         `json:"prdy_ctrt,string"`
	SdprVrssPrpr     decimal.Decimal `json:"sdpr_vrss_prpr"`
	SdprVrssPrprRate float64         `json:"sdpr_vrss_prpr_rate,string"`
	CntgVol          int64           `json:"cntg_vol,string"`
}

// InquireExpClosingPriceResponse 는 장마감 예상체결가 응답.
type InquireExpClosingPriceResponse struct {
	RtCd    string                `json:"rt_cd"`
	MsgCd   string                `json:"msg_cd"`
	Msg1    string                `json:"msg1"`
	Output1 []ExpClosingPriceItem `json:"output1"`
}

// InquireExpClosingPrice 는 장마감 예상체결가를 조회한다 (FHKST117300C0).
func (c *Client) InquireExpClosingPrice(ctx context.Context, params InquireExpClosingPriceParams) (*InquireExpClosingPriceResponse, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "11173"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/exp-closing-price",
		TrID:   "FHKST117300C0",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_COND_SCR_DIV_CODE":  scrDiv,
			"FID_RANK_SORT_CLS_CODE": params.RankSortClsCode,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_BLNG_CLS_CODE":      params.BlngClsCode,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res InquireExpClosingPriceResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InquireExpClosingPriceResponse: %w", err)
	}
	return &res, nil
}
