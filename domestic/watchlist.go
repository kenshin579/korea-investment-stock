package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// IntstockMultprice 는 관심종목 멀티 시세 (FHKST11300006) 응답.
//
// 한투 docs: docs/api/국내주식/관심종목멀티시세.md
// path: /uapi/domestic-stock/v1/quotations/intstock-multprice
//
// 30종목 batch 입력 (FID_COND_MRKT_DIV_CODE_1..30 + FID_INPUT_ISCD_1..30),
// 단일 output 객체 반환.
type IntstockMultprice struct {
	Output IntstockMultpriceData `json:"output"`
}

// IntstockMultpriceData 는 관심종목 멀티 시세 응답의 output object (29 fields).
type IntstockMultpriceData struct {
	KospiKosdaqClsName   string          `json:"kospi_kosdaq_cls_name"`           // 코스피 코스닥 구분명
	MrktTrtmClsName      string          `json:"mrkt_trtm_cls_name"`              // 시장 처리 구분명
	HourClsCode          string          `json:"hour_cls_code"`                   // 시간 구분 코드
	InterShrnIscd        string          `json:"inter_shrn_iscd"`                 // 관심 단축 종목코드
	InterKorIsnm         string          `json:"inter_kor_isnm"`                  // 관심 한글 종목명
	Inter2Prpr           decimal.Decimal `json:"inter2_prpr"`                     // 관심2 현재가
	Inter2PrdyVrss       decimal.Decimal `json:"inter2_prdy_vrss"`                // 관심2 전일 대비
	PrdyVrssSign         string          `json:"prdy_vrss_sign"`                  // 전일 대비 부호
	PrdyCtrt             float64         `json:"prdy_ctrt,string"`                // 전일 대비율
	AcmlVol              int64           `json:"acml_vol,string"`                 // 누적 거래량
	Inter2Oprc           decimal.Decimal `json:"inter2_oprc"`                     // 관심2 시가
	Inter2Hgpr           decimal.Decimal `json:"inter2_hgpr"`                     // 관심2 최고가
	Inter2Lwpr           decimal.Decimal `json:"inter2_lwpr"`                     // 관심2 최저가
	Inter2Llam           decimal.Decimal `json:"inter2_llam"`                     // 관심2 하한가
	Inter2Mxpr           decimal.Decimal `json:"inter2_mxpr"`                     // 관심2 상한가
	Inter2Askp           decimal.Decimal `json:"inter2_askp"`                     // 관심2 매도호가
	Inter2Bidp           decimal.Decimal `json:"inter2_bidp"`                     // 관심2 매수호가
	SelnRsqn             int64           `json:"seln_rsqn,string"`                // 매도 잔량
	ShnuRsqn             int64           `json:"shnu_rsqn,string"`                // 매수 잔량
	TotalAskpRsqn        int64           `json:"total_askp_rsqn,string"`          // 총 매도호가 잔량
	TotalBidpRsqn        int64           `json:"total_bidp_rsqn,string"`          // 총 매수호가 잔량
	AcmlTrPbmn           int64           `json:"acml_tr_pbmn,string"`             // 누적 거래 대금
	Inter2PrdyClpr       decimal.Decimal `json:"inter2_prdy_clpr"`                // 관심2 전일 종가
	OprcVrssHgprRate     float64         `json:"oprc_vrss_hgpr_rate,string"`      // 시가 대비 최고가 비율
	IntrAntcCntgVrss     decimal.Decimal `json:"intr_antc_cntg_vrss"`             // 관심 예상 체결 대비
	IntrAntcCntgVrssSign string          `json:"intr_antc_cntg_vrss_sign"`        // 관심 예상 체결 대비 부호
	IntrAntcCntgPrdyCtrt float64         `json:"intr_antc_cntg_prdy_ctrt,string"` // 관심 예상 체결 전일 대비율
	IntrAntcVol          int64           `json:"intr_antc_vol,string"`            // 관심 예상 거래량
	Inter2Sdpr           decimal.Decimal `json:"inter2_sdpr"`                     // 관심2 기준가
}

// InquireIntstockMultpriceParams 는 관심종목 멀티 시세 조회 파라미터.
//
// MarketCodes/Symbols 는 최대 30 항목 슬라이스.
// 각 항목이 FID_COND_MRKT_DIV_CODE_1..30 / FID_INPUT_ISCD_1..30 로 전송됨.
type InquireIntstockMultpriceParams struct {
	MarketCodes []string // FID_COND_MRKT_DIV_CODE_1..30 — 시장 구분 코드 목록 (예 "J")
	Symbols     []string // FID_INPUT_ISCD_1..30 — 종목코드 목록 (예 "005930")
}

// InquireIntstockMultprice 는 관심종목 멀티 시세 호출.
//
// 한투 docs: docs/api/국내주식/관심종목멀티시세.md
// path: /uapi/domestic-stock/v1/quotations/intstock-multprice (FHKST11300006)
//
// 최대 30종목 batch 입력: MarketCodes[i] + Symbols[i] 쌍이 FID_COND_MRKT_DIV_CODE_{i+1},
// FID_INPUT_ISCD_{i+1} 로 전송. min(len(MarketCodes), len(Symbols), 30) 만큼 처리.
func (c *Client) InquireIntstockMultprice(ctx context.Context, params InquireIntstockMultpriceParams) (*IntstockMultprice, error) {
	n := len(params.MarketCodes)
	if len(params.Symbols) < n {
		n = len(params.Symbols)
	}
	if n > 30 {
		n = 30
	}

	query := make(map[string]string, n*2)
	for i := 0; i < n; i++ {
		idx := fmt.Sprintf("%d", i+1)
		query["FID_COND_MRKT_DIV_CODE_"+idx] = params.MarketCodes[i]
		query["FID_INPUT_ISCD_"+idx] = params.Symbols[i]
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-stock/v1/quotations/intstock-multprice",
		TrID:     "FHKST11300006",
		Query:    query,
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res IntstockMultprice
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse IntstockMultprice: %w", err)
	}
	return &res, nil
}
