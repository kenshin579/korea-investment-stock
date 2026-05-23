package futures

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
	"github.com/kenshin579/korea-investment-stock/kistypes"
)

// ─── EP5: ExpPriceTrend ──────────────────────────────────────────────────────

// ExpPriceTrendOutput1 는 선물옵션 일중예상체결추이 현재 상태 (6 필드).
//
// 주의: docs 한글명에 오기 존재 (hts_kor_isnm "영업 시간", futs_antc_cnpr "업종 지수 현재가" 등).
// 필드명 기준으로 의미 해석.
type ExpPriceTrendOutput1 struct {
	HtsKorIsnm       string          `json:"hts_kor_isnm"`        // HTS 한글 종목명
	FutsAntcCnpr     decimal.Decimal `json:"futs_antc_cnpr"`      // 선물 예상 체결가
	AntcCntgVrssSign string          `json:"antc_cntg_vrss_sign"` // 예상 체결 대비 부호
	FutsAntcCntgVrss decimal.Decimal `json:"futs_antc_cntg_vrss"` // 선물 예상 체결 대비
	AntcCntgPrdyCtrt kistypes.Float  `json:"antc_cntg_prdy_ctrt"` // 예상 체결 전일 대비율
	FutsSdpr         decimal.Decimal `json:"futs_sdpr"`           // 선물 기준가
}

// ExpPriceTrendOutput2Item 는 선물옵션 일중예상체결추이 시계열 (5 필드).
type ExpPriceTrendOutput2Item struct {
	StckCntgHour     string          `json:"stck_cntg_hour"`      // 주식 체결 시간
	FutsAntcCnpr     decimal.Decimal `json:"futs_antc_cnpr"`      // 선물 예상 체결가
	AntcCntgVrssSign string          `json:"antc_cntg_vrss_sign"` // 예상 체결 대비 부호
	FutsAntcCntgVrss decimal.Decimal `json:"futs_antc_cntg_vrss"` // 선물 예상 체결 대비
	AntcCntgPrdyCtrt kistypes.Float  `json:"antc_cntg_prdy_ctrt"` // 예상 체결 전일 대비율
}

// ExpPriceTrendData 는 선물옵션 일중예상체결추이 응답 (output1 + output2[]).
type ExpPriceTrendData struct {
	Output1 ExpPriceTrendOutput1       `json:"output1"`
	Output2 []ExpPriceTrendOutput2Item `json:"output2"`
}

type expPriceTrendResponse struct {
	RtCd    string                     `json:"rt_cd"`
	MsgCd   string                     `json:"msg_cd"`
	Msg1    string                     `json:"msg1"`
	Output1 ExpPriceTrendOutput1       `json:"output1"`
	Output2 []ExpPriceTrendOutput2Item `json:"output2"`
}

// ExpPriceTrendParams 는 선물옵션 일중예상체결추이 파라미터.
type ExpPriceTrendParams struct {
	Code       string // 입력 종목코드 (지수선물:6자리, 지수옵션:9자리)
	MarketCode string // 조건 시장 분류 코드 (F:지수선물, O:지수옵션)
}

// ExpPriceTrend 는 선물옵션 일중예상체결추이 조회 (FHPIF05110100).
//
// 모의: 미지원 (실전 only)
// F/O 시장만 지원 (지수선물/지수옵션).
//
// KIS API: GET /uapi/domestic-futureoption/v1/quotations/exp-price-trend
func (c *Client) ExpPriceTrend(ctx context.Context, params ExpPriceTrendParams) (*ExpPriceTrendData, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-futureoption/v1/quotations/exp-price-trend",
		TrID:     "FHPIF05110100",
		CustType: "P",
		Query: map[string]string{
			"FID_INPUT_ISCD":         params.Code,
			"FID_COND_MRKT_DIV_CODE": params.MarketCode,
		},
	})
	if err != nil {
		return nil, err
	}
	var res expPriceTrendResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse ExpPriceTrend: %w", err)
	}
	return &ExpPriceTrendData{
		Output1: res.Output1,
		Output2: res.Output2,
	}, nil
}
