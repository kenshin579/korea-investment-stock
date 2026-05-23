package overseasfutures

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
	"github.com/kenshin579/korea-investment-stock/kistypes"
)

// ─── 공통: CcnlParams ────────────────────────────────────────────────────────

// CcnlParams 는 해외선물 체결추이 조회 공통 파라미터 (EP4~EP7).
type CcnlParams struct {
	SrsCd         string // 종목코드 (예: 6AM24)
	ExchCd        string // 거래소코드 (예: CME)
	StartDateTime string // 조회시작일시 (공백 가능)
	CloseDateTime string // 조회종료일시 (YYYYMMDD)
	QryTp         string // 조회구분 (Q: 최초조회, P: 다음조회)
	QryCnt        string // 요청개수 (최대 40)
	QryGap        string // 묶음개수 (공백, 분만 사용)
	IndexKey      string // 이전조회KEY
}

// ─── 공통: CcnlOutput2Item ───────────────────────────────────────────────────

// CcnlOutput2Item 는 해외선물 체결추이 배열 항목 공통 구조 (11 필드).
type CcnlOutput2Item struct {
	DataDate      string          `json:"data_date"`        // 일자
	DataTime      string          `json:"data_time"`        // 시각
	OpenPrice     decimal.Decimal `json:"open_price"`       // 시가
	HighPrice     decimal.Decimal `json:"high_price"`       // 고가
	LowPrice      decimal.Decimal `json:"low_price"`        // 저가
	LastPrice     decimal.Decimal `json:"last_price"`       // 체결가격
	LastQntt      int64           `json:"last_qntt,string"` // 체결수량
	Vol           int64           `json:"vol,string"`       // 누적거래수량
	PrevDiffFlag  string          `json:"prev_diff_flag"`   // 전일대비구분
	PrevDiffPrice decimal.Decimal `json:"prev_diff_price"`  // 전일대비가격
	PrevDiffRate  kistypes.Float  `json:"prev_diff_rate"`   // 전일대비율
}

// ─── EP2: InquireTimeFuturechartprice ────────────────────────────────────────

// InquireTimeFuturechartpriceOutput2 는 분봉조회 메타 정보 (3 필드).
// 주의: docs 에서 output2 가 단일 Object (메타), output1 이 배열 (분봉). 키 역전 패턴.
type InquireTimeFuturechartpriceOutput2 struct {
	RetCnt   string `json:"ret_cnt"`    // 자료개수
	LastNCnt string `json:"last_n_cnt"` // N틱최종개수
	IndexKey string `json:"index_key"`  // 이전조회KEY
}

// InquireTimeFuturechartpriceData 는 해외선물 분봉조회 응답.
// Output2 (메타 단일) + Output1 (분봉 배열) — docs 네이밍 그대로.
type InquireTimeFuturechartpriceData struct {
	Output2 InquireTimeFuturechartpriceOutput2 `json:"output2"` // 메타 (단일)
	Output1 []CcnlOutput2Item                  `json:"output1"` // 분봉 배열
}

type inquireTimeFuturechartpriceResponse struct {
	RtCd    string                             `json:"rt_cd"`
	MsgCd   string                             `json:"msg_cd"`
	Msg1    string                             `json:"msg1"`
	Output2 InquireTimeFuturechartpriceOutput2 `json:"output2"`
	Output1 []CcnlOutput2Item                  `json:"output1"`
}

// InquireTimeFuturechartpriceParams 는 해외선물 분봉조회 파라미터.
type InquireTimeFuturechartpriceParams struct {
	SrsCd         string // 종목코드 (예: CNHU24)
	ExchCd        string // 거래소코드 (예: CME)
	StartDateTime string // 조회시작일시 (공백 가능)
	CloseDateTime string // 조회종료일시 (예: 20231214)
	QryTp         string // 조회구분 (Q: 최초조회, P: 다음조회)
	QryCnt        string // 요청개수 (최대 120)
	QryGap        string // 묶음개수/분간격 (예: 5)
	IndexKey      string // 이전조회KEY
}

// InquireTimeFuturechartprice 는 해외선물 분봉조회 (HHDFC55020400).
//
// output2 (메타 단일) + output1 (분봉 배열) — docs 키 역전 패턴 그대로 반영.
// 최대 120건/회. 다음조회: QryTp="P", IndexKey=output2.IndexKey.
//
// KIS API: GET /uapi/overseas-futureoption/v1/quotations/inquire-time-futurechartprice
func (c *Client) InquireTimeFuturechartprice(ctx context.Context, params InquireTimeFuturechartpriceParams) (*InquireTimeFuturechartpriceData, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/overseas-futureoption/v1/quotations/inquire-time-futurechartprice",
		TrID:     "HHDFC55020400",
		CustType: "P",
		Query: map[string]string{
			"SRS_CD":          params.SrsCd,
			"EXCH_CD":         params.ExchCd,
			"START_DATE_TIME": params.StartDateTime,
			"CLOSE_DATE_TIME": params.CloseDateTime,
			"QRY_TP":          params.QryTp,
			"QRY_CNT":         params.QryCnt,
			"QRY_GAP":         params.QryGap,
			"INDEX_KEY":       params.IndexKey,
		},
	})
	if err != nil {
		return nil, err
	}
	var res inquireTimeFuturechartpriceResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InquireTimeFuturechartprice: %w", err)
	}
	return &InquireTimeFuturechartpriceData{
		Output2: res.Output2,
		Output1: res.Output1,
	}, nil
}

// ─── EP4: MonthlyCcnl ────────────────────────────────────────────────────────

// MonthlyCcnlOutput1 는 월간 체결추이 메타 정보 (3 필드).
type MonthlyCcnlOutput1 struct {
	TretCnt  string `json:"tret_cnt"`   // 자료개수
	LastNCnt string `json:"last_n_cnt"` // N틱최종개수
	IndexKey string `json:"index_key"`  // 이전조회KEY
}

// MonthlyCcnlData 는 해외선물 월간 체결추이 응답 (output1 메타 + output2 배열).
type MonthlyCcnlData struct {
	Output1 MonthlyCcnlOutput1 `json:"output1"`
	Output2 []CcnlOutput2Item  `json:"output2"`
}

type monthlyCcnlResponse struct {
	RtCd    string             `json:"rt_cd"`
	MsgCd   string             `json:"msg_cd"`
	Msg1    string             `json:"msg1"`
	Output1 MonthlyCcnlOutput1 `json:"output1"`
	Output2 []CcnlOutput2Item  `json:"output2"`
}

// MonthlyCcnl 는 해외선물 월간 체결추이 조회 (HHDFC55020300).
//
// 최대 40건/회. 다음조회: QryTp="P", IndexKey=output1.IndexKey.
//
// KIS API: GET /uapi/overseas-futureoption/v1/quotations/monthly-ccnl
func (c *Client) MonthlyCcnl(ctx context.Context, params CcnlParams) (*MonthlyCcnlData, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/overseas-futureoption/v1/quotations/monthly-ccnl",
		TrID:     "HHDFC55020300",
		CustType: "P",
		Query: map[string]string{
			"SRS_CD":          params.SrsCd,
			"EXCH_CD":         params.ExchCd,
			"START_DATE_TIME": params.StartDateTime,
			"CLOSE_DATE_TIME": params.CloseDateTime,
			"QRY_TP":          params.QryTp,
			"QRY_CNT":         params.QryCnt,
			"QRY_GAP":         params.QryGap,
			"INDEX_KEY":       params.IndexKey,
		},
	})
	if err != nil {
		return nil, err
	}
	var res monthlyCcnlResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse MonthlyCcnl: %w", err)
	}
	return &MonthlyCcnlData{
		Output1: res.Output1,
		Output2: res.Output2,
	}, nil
}

// ─── EP5: DailyCcnl ──────────────────────────────────────────────────────────

// DailyCcnlOutput1 는 일간 체결추이 메타 정보 (3 필드).
type DailyCcnlOutput1 struct {
	TretCnt  string `json:"tret_cnt"`   // 자료개수
	LastNCnt string `json:"last_n_cnt"` // N틱최종개수
	IndexKey string `json:"index_key"`  // 이전조회KEY
}

// DailyCcnlData 는 해외선물 일간 체결추이 응답 (output1 메타 + output2 배열).
type DailyCcnlData struct {
	Output1 DailyCcnlOutput1  `json:"output1"`
	Output2 []CcnlOutput2Item `json:"output2"`
}

type dailyCcnlResponse struct {
	RtCd    string            `json:"rt_cd"`
	MsgCd   string            `json:"msg_cd"`
	Msg1    string            `json:"msg1"`
	Output1 DailyCcnlOutput1  `json:"output1"`
	Output2 []CcnlOutput2Item `json:"output2"`
}

// DailyCcnl 는 해외선물 일간 체결추이 조회 (HHDFC55020100).
//
// 최대 40건/회. 다음조회: QryTp="P", IndexKey=output1.IndexKey.
//
// KIS API: GET /uapi/overseas-futureoption/v1/quotations/daily-ccnl
func (c *Client) DailyCcnl(ctx context.Context, params CcnlParams) (*DailyCcnlData, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/overseas-futureoption/v1/quotations/daily-ccnl",
		TrID:     "HHDFC55020100",
		CustType: "P",
		Query: map[string]string{
			"SRS_CD":          params.SrsCd,
			"EXCH_CD":         params.ExchCd,
			"START_DATE_TIME": params.StartDateTime,
			"CLOSE_DATE_TIME": params.CloseDateTime,
			"QRY_TP":          params.QryTp,
			"QRY_CNT":         params.QryCnt,
			"QRY_GAP":         params.QryGap,
			"INDEX_KEY":       params.IndexKey,
		},
	})
	if err != nil {
		return nil, err
	}
	var res dailyCcnlResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse DailyCcnl: %w", err)
	}
	return &DailyCcnlData{
		Output1: res.Output1,
		Output2: res.Output2,
	}, nil
}

// ─── EP6: WeeklyCcnl ─────────────────────────────────────────────────────────

// WeeklyCcnlOutput1 는 주간 체결추이 메타 정보 (3 필드).
// 주의: 자료개수 필드명이 `ret_cnt` (EP4/EP5/EP7 은 `tret_cnt`) — docs 불일치.
type WeeklyCcnlOutput1 struct {
	RetCnt   string `json:"ret_cnt"`    // 자료개수 (anomaly: ret_cnt, not tret_cnt)
	LastNCnt string `json:"last_n_cnt"` // N틱최종개수
	IndexKey string `json:"index_key"`  // 이전조회KEY
}

// WeeklyCcnlData 는 해외선물 주간 체결추이 응답 (output1 메타 + output2 배열).
type WeeklyCcnlData struct {
	Output1 WeeklyCcnlOutput1 `json:"output1"`
	Output2 []CcnlOutput2Item `json:"output2"`
}

type weeklyCcnlResponse struct {
	RtCd    string            `json:"rt_cd"`
	MsgCd   string            `json:"msg_cd"`
	Msg1    string            `json:"msg1"`
	Output1 WeeklyCcnlOutput1 `json:"output1"`
	Output2 []CcnlOutput2Item `json:"output2"`
}

// WeeklyCcnl 는 해외선물 주간 체결추이 조회 (HHDFC55020000).
//
// output1 자료개수 필드명이 ret_cnt (EP4/EP5/EP7 은 tret_cnt) — docs 원본 반영.
// 최대 40건/회. 다음조회: QryTp="P", IndexKey=output1.IndexKey.
//
// KIS API: GET /uapi/overseas-futureoption/v1/quotations/weekly-ccnl
func (c *Client) WeeklyCcnl(ctx context.Context, params CcnlParams) (*WeeklyCcnlData, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/overseas-futureoption/v1/quotations/weekly-ccnl",
		TrID:     "HHDFC55020000",
		CustType: "P",
		Query: map[string]string{
			"SRS_CD":          params.SrsCd,
			"EXCH_CD":         params.ExchCd,
			"START_DATE_TIME": params.StartDateTime,
			"CLOSE_DATE_TIME": params.CloseDateTime,
			"QRY_TP":          params.QryTp,
			"QRY_CNT":         params.QryCnt,
			"QRY_GAP":         params.QryGap,
			"INDEX_KEY":       params.IndexKey,
		},
	})
	if err != nil {
		return nil, err
	}
	var res weeklyCcnlResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse WeeklyCcnl: %w", err)
	}
	return &WeeklyCcnlData{
		Output1: res.Output1,
		Output2: res.Output2,
	}, nil
}

// ─── EP7: TickCcnl ───────────────────────────────────────────────────────────

// TickCcnlOutput1 는 틱 체결추이 메타 정보 (3 필드).
type TickCcnlOutput1 struct {
	TretCnt  string `json:"tret_cnt"`   // 자료개수
	LastNCnt string `json:"last_n_cnt"` // N틱최종개수
	IndexKey string `json:"index_key"`  // 이전조회KEY
}

// TickCcnlData 는 해외선물 틱 체결추이 응답 (output1 메타 + output2 배열).
type TickCcnlData struct {
	Output1 TickCcnlOutput1   `json:"output1"`
	Output2 []CcnlOutput2Item `json:"output2"`
}

type tickCcnlResponse struct {
	RtCd    string            `json:"rt_cd"`
	MsgCd   string            `json:"msg_cd"`
	Msg1    string            `json:"msg1"`
	Output1 TickCcnlOutput1   `json:"output1"`
	Output2 []CcnlOutput2Item `json:"output2"`
}

// TickCcnl 는 해외선물 틱 체결추이 조회 (HHDFC55020200).
//
// 최대 40건/회. 다음조회: QryTp="P", IndexKey=output1.IndexKey.
//
// KIS API: GET /uapi/overseas-futureoption/v1/quotations/tick-ccnl
func (c *Client) TickCcnl(ctx context.Context, params CcnlParams) (*TickCcnlData, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/overseas-futureoption/v1/quotations/tick-ccnl",
		TrID:     "HHDFC55020200",
		CustType: "P",
		Query: map[string]string{
			"SRS_CD":          params.SrsCd,
			"EXCH_CD":         params.ExchCd,
			"START_DATE_TIME": params.StartDateTime,
			"CLOSE_DATE_TIME": params.CloseDateTime,
			"QRY_TP":          params.QryTp,
			"QRY_CNT":         params.QryCnt,
			"QRY_GAP":         params.QryGap,
			"INDEX_KEY":       params.IndexKey,
		},
	})
	if err != nil {
		return nil, err
	}
	var res tickCcnlResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse TickCcnl: %w", err)
	}
	return &TickCcnlData{
		Output1: res.Output1,
		Output2: res.Output2,
	}, nil
}
