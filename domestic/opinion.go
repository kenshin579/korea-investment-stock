package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
	"github.com/shopspring/decimal"
)

// InvestOpinion 은 종목투자의견 (FHKST663300C0) 응답.
//
// 한투 docs: docs/api/국내주식/종목투자의견.md
// path: /uapi/domestic-stock/v1/quotations/invest-opinion
type InvestOpinion struct {
	Output []InvestOpinionItem `json:"output"`
}

// InvestOpinionItem 은 응답의 output 한 행 (12 fields).
type InvestOpinionItem struct {
	StckBsopDate        string          `json:"stck_bsop_date"`          // 영업 일자
	InvtOpnn            string          `json:"invt_opnn"`               // 투자의견
	InvtOpnnClsCode     string          `json:"invt_opnn_cls_code"`      // 투자의견 구분 코드
	RgbfInvtOpnn        string          `json:"rgbf_invt_opnn"`          // 직전 투자의견
	RgbfInvtOpnnClsCode string          `json:"rgbf_invt_opnn_cls_code"` // 직전 투자의견 구분 코드
	MbcrName            string          `json:"mbcr_name"`               // 회원사명
	HtsGoalPrc          decimal.Decimal `json:"hts_goal_prc"`            // HTS 목표 가격
	StckPrdyClpr        decimal.Decimal `json:"stck_prdy_clpr"`          // 주식 전일 종가
	StckNdayEsdg        decimal.Decimal `json:"stck_nday_esdg"`          // 주식 N일 추정 단가
	NdayDprt            float64         `json:"nday_dprt,string"`        // N일 이격도
	StftEsdg            decimal.Decimal `json:"stft_esdg"`               // 직전 추정 단가
	Dprt                float64         `json:"dprt,string"`             // 이격도
}

// InquireInvestOpinionParams 는 종목투자의견 조회 파라미터.
type InquireInvestOpinionParams struct {
	MarketCode     string // FID_COND_MRKT_DIV_CODE — 빈 값=>"J"
	CondScrDivCode string // FID_COND_SCR_DIV_CODE — 빈 값=>"16633"
	Symbol         string // FID_INPUT_ISCD — 필수, 단축 종목코드 (예 "005930")
	StartDate      string // FID_INPUT_DATE_1 — 조회 시작일 YYYYMMDD
	EndDate        string // FID_INPUT_DATE_2 — 조회 종료일 YYYYMMDD
}

// InquireInvestOpinion 은 종목투자의견 호출.
//
// 한투 docs: docs/api/국내주식/종목투자의견.md
// path: /uapi/domestic-stock/v1/quotations/invest-opinion (FHKST663300C0)
func (c *Client) InquireInvestOpinion(ctx context.Context, params InquireInvestOpinionParams) (*InvestOpinion, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "16633"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/invest-opinion",
		TrID:   "FHKST663300C0",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_COND_SCR_DIV_CODE":  scrDiv,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_INPUT_DATE_1":       params.StartDate,
			"FID_INPUT_DATE_2":       params.EndDate,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res InvestOpinion
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InvestOpinion: %w", err)
	}
	return &res, nil
}

// InvestOpbysec 은 증권사별 투자의견 (FHKST663400C0) 응답.
//
// 한투 docs: docs/api/국내주식/증권사별투자의견.md
// path: /uapi/domestic-stock/v1/quotations/invest-opbysec
type InvestOpbysec struct {
	Output []InvestOpbysecItem `json:"output"`
}

// InvestOpbysecItem 은 응답의 output 한 행 (16 fields).
type InvestOpbysecItem struct {
	StckBsopDate        string          `json:"stck_bsop_date"`          // 영업 일자
	StckShrnIscd        string          `json:"stck_shrn_iscd"`          // 주식 단축 종목코드
	HtsKorIsnm          string          `json:"hts_kor_isnm"`            // HTS 한글 종목명
	InvtOpnn            string          `json:"invt_opnn"`               // 투자의견
	InvtOpnnClsCode     string          `json:"invt_opnn_cls_code"`      // 투자의견 구분 코드
	RgbfInvtOpnn        string          `json:"rgbf_invt_opnn"`          // 직전 투자의견
	RgbfInvtOpnnClsCode string          `json:"rgbf_invt_opnn_cls_code"` // 직전 투자의견 구분 코드
	MbcrName            string          `json:"mbcr_name"`               // 회원사명
	StckPrpr            decimal.Decimal `json:"stck_prpr"`               // 주식 현재가
	PrdyVrss            decimal.Decimal `json:"prdy_vrss"`               // 전일 대비
	PrdyVrssSign        string          `json:"prdy_vrss_sign"`          // 전일 대비 부호
	PrdyCtrt            float64         `json:"prdy_ctrt,string"`        // 전일 대비율
	HtsGoalPrc          decimal.Decimal `json:"hts_goal_prc"`            // HTS 목표 가격
	StckPrdyClpr        decimal.Decimal `json:"stck_prdy_clpr"`          // 주식 전일 종가
	StftEsdg            decimal.Decimal `json:"stft_esdg"`               // 직전 추정 단가
	Dprt                float64         `json:"dprt,string"`             // 이격도
}

// InquireInvestOpbysecParams 는 증권사별 투자의견 조회 파라미터.
type InquireInvestOpbysecParams struct {
	MarketCode     string // FID_COND_MRKT_DIV_CODE — 빈 값=>"J"
	CondScrDivCode string // FID_COND_SCR_DIV_CODE — 빈 값=>"16634"
	SecBrokerCode  string // FID_INPUT_ISCD — 필수, 증권사코드 (종목코드 아님!)
	DivClsCode     string // FID_DIV_CLS_CODE — 0=전체
	StartDate      string // FID_INPUT_DATE_1 — 조회 시작일 YYYYMMDD
	EndDate        string // FID_INPUT_DATE_2 — 조회 종료일 YYYYMMDD
}

// InquireInvestOpbysec 은 증권사별 투자의견 호출.
//
// 한투 docs: docs/api/국내주식/증권사별투자의견.md
// path: /uapi/domestic-stock/v1/quotations/invest-opbysec (FHKST663400C0)
//
// 주의: FID_INPUT_ISCD 는 종목코드가 아닌 증권사코드를 입력한다.
func (c *Client) InquireInvestOpbysec(ctx context.Context, params InquireInvestOpbysecParams) (*InvestOpbysec, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "16634"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/invest-opbysec",
		TrID:   "FHKST663400C0",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_COND_SCR_DIV_CODE":  scrDiv,
			"FID_INPUT_ISCD":         params.SecBrokerCode,
			"FID_DIV_CLS_CODE":       params.DivClsCode,
			"FID_INPUT_DATE_1":       params.StartDate,
			"FID_INPUT_DATE_2":       params.EndDate,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res InvestOpbysec
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InvestOpbysec: %w", err)
	}
	return &res, nil
}

// EstimatePerform 은 종목 추정실적 (HHKST668300C0) 응답.
//
// 한투 docs: docs/api/국내주식/종목추정실적.md
// path: /uapi/domestic-stock/v1/quotations/estimate-perform
//
// QUAD OUTPUT: output1(요약) + output2(추정손익계산서 6행) + output3(투자지표 8행) + output4(결산년월 5행).
// output1 KIS docs body labels 오표기 — Python dataclass field names 사용.
type EstimatePerform struct {
	Output1 EstimatePerformSummary           `json:"output1"`
	Output2 []EstimatePerformIncomeStatement `json:"output2"`
	Output3 []EstimatePerformInvestIndicator `json:"output3"`
	Output4 []EstimatePerformPeriod          `json:"output4"`
}

// EstimatePerformSummary 는 응답의 output1 (종목 요약, 8 fields).
// KIS docs body table 오표기 — Python dataclass 기준 field names 사용.
type EstimatePerformSummary struct {
	ShtCd         string `json:"sht_cd"`          // 단축 종목코드
	ItemKorNm     string `json:"item_kor_nm"`     // 한글 종목명
	Name1         string `json:"name1"`           // 현재가 (opaque)
	Name2         string `json:"name2"`           // 전일 대비 (opaque)
	Estdate       string `json:"estdate"`         // 결산 일자
	RcmdName      string `json:"rcmd_name"`       // 투자의견명
	Capital       string `json:"capital"`         // 자본금/거래량 (opaque)
	FornItemLmtrt string `json:"forn_item_lmtrt"` // 외국인 한도 비율 (opaque)
}

// EstimatePerformIncomeStatement 는 output2(추정손익계산서) 한 행 (5 fields).
type EstimatePerformIncomeStatement struct {
	Data1 string `json:"data1"` // 항목명
	Data2 string `json:"data2"` // 기간1 값
	Data3 string `json:"data3"` // 기간2 값
	Data4 string `json:"data4"` // 기간3 값
	Data5 string `json:"data5"` // 기간4 값
}

// EstimatePerformInvestIndicator 는 output3(투자지표) 한 행 (5 fields).
type EstimatePerformInvestIndicator struct {
	Data1 string `json:"data1"` // 항목명
	Data2 string `json:"data2"` // 기간1 값
	Data3 string `json:"data3"` // 기간2 값
	Data4 string `json:"data4"` // 기간3 값
	Data5 string `json:"data5"` // 기간4 값
}

// EstimatePerformPeriod 는 output4 한 행 (결산년월, 1 field).
type EstimatePerformPeriod struct {
	Dt string `json:"dt"` // 결산 년월 (예 "202412", "202612E")
}

// InquireEstimatePerformParams 는 종목 추정실적 조회 파라미터.
type InquireEstimatePerformParams struct {
	Symbol string // SHT_CD — 필수, 6자리 단축 종목코드 (비표준 param명, FID_ 접두어 없음)
}

// InquireEstimatePerform 은 종목 추정실적 호출.
//
// 한투 docs: docs/api/국내주식/종목추정실적.md
// path: /uapi/domestic-stock/v1/quotations/estimate-perform (HHKST668300C0)
//
// 주의: query param 이름이 SHT_CD (FID_ 접두어 없음).
func (c *Client) InquireEstimatePerform(ctx context.Context, params InquireEstimatePerformParams) (*EstimatePerform, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/estimate-perform",
		TrID:   "HHKST668300C0",
		Query: map[string]string{
			"SHT_CD": params.Symbol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res EstimatePerform
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse EstimatePerform: %w", err)
	}
	return &res, nil
}
