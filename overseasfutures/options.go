package overseasfutures

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// ─── EP1: InquireTimeOptchartprice ───────────────────────────────────────────

// InquireTimeOptchartpriceParams 는 해외옵션 분봉조회 파라미터.
type InquireTimeOptchartpriceParams struct {
	SrsCd         string // 종목코드 (예: OESU24 C5500)
	ExchCd        string // 거래소코드 (예: CME)
	StartDateTime string // 조회시작일시 (공백)
	CloseDateTime string // 조회종료일시
	QryTp         string // 조회구분 (Q: 최초조회, P: 다음조회)
	QryCnt        string // 요청개수 (최대 120)
	QryGap        string // 분간격 (1:1분봉, 5:5분봉 ...)
	IndexKey      string // 이전조회KEY
}

// InquireTimeOptchartpriceOutput2 는 해외옵션 분봉조회 메타 정보 (3 필드).
// 주의: docs 에서 output2 가 단일 Object (메타), output1 이 배열 (분봉). 키 역전 패턴.
// Phase 11.5 EP2(InquireTimeFuturechartprice) 와 완전 동일 구조.
type InquireTimeOptchartpriceOutput2 = InquireTimeFuturechartpriceOutput2

// InquireTimeOptchartpriceData 는 해외옵션 분봉조회 응답.
// Output2 (메타 단일) + Output1 (분봉 배열) — docs 네이밍 그대로.
type InquireTimeOptchartpriceData = InquireTimeFuturechartpriceData

type inquireTimeOptchartpriceResponse struct {
	RtCd    string                             `json:"rt_cd"`
	MsgCd   string                             `json:"msg_cd"`
	Msg1    string                             `json:"msg1"`
	Output2 InquireTimeFuturechartpriceOutput2 `json:"output2"`
	Output1 []CcnlOutput2Item                  `json:"output1"`
}

// InquireTimeOptchartprice 는 해외옵션 분봉조회 (HHDFO55020400).
//
// output2 (메타 단일) + output1 (분봉 배열) — docs 키 역전 패턴 그대로 반영.
// 최대 120건/회. 다음조회: QryTp="P", IndexKey=output2.IndexKey.
// Phase 11.5 EP2(InquireTimeFuturechartprice) 와 동일 구조.
//
// KIS API: GET /uapi/overseas-futureoption/v1/quotations/inquire-time-optchartprice
func (c *Client) InquireTimeOptchartprice(ctx context.Context, params InquireTimeOptchartpriceParams) (*InquireTimeOptchartpriceData, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/overseas-futureoption/v1/quotations/inquire-time-optchartprice",
		TrID:     "HHDFO55020400",
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
	var res inquireTimeOptchartpriceResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InquireTimeOptchartprice: %w", err)
	}
	return &InquireTimeOptchartpriceData{
		Output2: res.Output2,
		Output1: res.Output1,
	}, nil
}

// ─── EP2: SearchOptDetail ────────────────────────────────────────────────────

// SearchOptDetailOutput2Item 는 해외옵션 상품기본정보 항목 (21 필드).
// DISTINCT: 선물 SearchContractDetail 대비 sub_exch_nm 없음 (22→21 필드). clas_cd length=1 (선물=3).
type SearchOptDetailOutput2Item struct {
	ExchCd        string          `json:"exch_cd"`         // 거래소코드
	ClasCd        string          `json:"clas_cd"`         // 품목종류 (옵션: length=1)
	CrcCd         string          `json:"crc_cd"`          // 거래통화
	SttlPrice     decimal.Decimal `json:"sttl_price"`      // 정산가
	SttlDate      string          `json:"sttl_date"`       // 정산일
	TrstMgn       decimal.Decimal `json:"trst_mgn"`        // 증거금
	DispDigit     string          `json:"disp_digit"`      // 가격표시진법
	TickSz        decimal.Decimal `json:"tick_sz"`         // 틱사이즈
	TickVal       decimal.Decimal `json:"tick_val"`        // 틱가치
	MrktOpenDate  string          `json:"mrkt_open_date"`  // 장개시일자
	MrktOpenTime  string          `json:"mrkt_open_time"`  // 장개시시각
	MrktCloseDate string          `json:"mrkt_close_date"` // 장마감일자
	MrktCloseTime string          `json:"mrkt_close_time"` // 장마감시각
	TrdFrDate     string          `json:"trd_fr_date"`     // 상장일
	ExprDate      string          `json:"expr_date"`       // 만기일
	TrdToDate     string          `json:"trd_to_date"`     // 최종거래일
	RemnCnt       string          `json:"remn_cnt"`        // 잔존일수
	StatTp        string          `json:"stat_tp"`         // 매매여부
	CtrtSize      decimal.Decimal `json:"ctrt_size"`       // 계약크기
	StlTp         string          `json:"stl_tp"`          // 최종결제구분
	FrstNotiDate  string          `json:"frst_noti_date"`  // 최초식별일
}

// SearchOptDetailData 는 해외옵션 상품기본정보 응답 (output2 배열만, output1 없음).
type SearchOptDetailData struct {
	Output2 []SearchOptDetailOutput2Item `json:"output2"`
}

type searchOptDetailResponse struct {
	RtCd    string                       `json:"rt_cd"`
	MsgCd   string                       `json:"msg_cd"`
	Msg1    string                       `json:"msg1"`
	Output2 []SearchOptDetailOutput2Item `json:"output2"`
}

// SearchOptDetailParams 는 해외옵션 상품기본정보 조회 파라미터.
type SearchOptDetailParams struct {
	Codes []string // 종목코드 목록 (최대 30개). QRY_CNT = len(Codes).
}

// SearchOptDetail 는 해외옵션 상품기본정보 조회 (HHDFO55200000).
//
// 한번에 최대 30개 종목 조회. output2 배열만 반환 (output1 없음).
// 선물 SearchContractDetail 대비: sub_exch_nm 없음, 최대 30개 (선물 32개).
//
// KIS API: GET /uapi/overseas-futureoption/v1/quotations/search-opt-detail
func (c *Client) SearchOptDetail(ctx context.Context, params SearchOptDetailParams) (*SearchOptDetailData, error) {
	query := map[string]string{
		"QRY_CNT": fmt.Sprintf("%d", len(params.Codes)),
	}
	for i, code := range params.Codes {
		key := fmt.Sprintf("SRS_CD_%02d", i+1)
		query[key] = code
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/overseas-futureoption/v1/quotations/search-opt-detail",
		TrID:     "HHDFO55200000",
		CustType: "P",
		Query:    query,
	})
	if err != nil {
		return nil, err
	}
	var res searchOptDetailResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse SearchOptDetail: %w", err)
	}
	return &SearchOptDetailData{Output2: res.Output2}, nil
}

// ─── 공통: OptCcnlOutput1 (옵션 체결추이 메타 — ret_cnt 통일) ─────────────────

// OptCcnlOutput1 는 해외옵션 체결추이 메타 정보 (3 필드).
// 옵션 EP3~EP6 전체: ret_cnt 통일 (선물 대응 EP 의 tret_cnt/ret_cnt 혼재와 다름).
type OptCcnlOutput1 struct {
	RetCnt   string `json:"ret_cnt"`    // 자료개수
	LastNCnt string `json:"last_n_cnt"` // N틱최종개수
	IndexKey string `json:"index_key"`  // 이전조회KEY
}

// ─── EP3: OptMonthlyCcnl ─────────────────────────────────────────────────────

// OptMonthlyCcnlData 는 해외옵션 월간 체결추이 응답 (output1 메타 + output2 배열).
// output2[] 구조는 선물 MonthlyCcnl 과 동일 (CcnlOutput2Item).
type OptMonthlyCcnlData struct {
	Output1 OptCcnlOutput1   `json:"output1"`
	Output2 []CcnlOutput2Item `json:"output2"`
}

type optMonthlyCcnlResponse struct {
	RtCd    string           `json:"rt_cd"`
	MsgCd   string           `json:"msg_cd"`
	Msg1    string           `json:"msg1"`
	Output1 OptCcnlOutput1   `json:"output1"`
	Output2 []CcnlOutput2Item `json:"output2"`
}

// OptMonthlyCcnl 는 해외옵션 월간 체결추이 조회 (HHDFO55020300).
//
// output1 카운트 필드: ret_cnt (선물 EP4 MonthlyCcnl 의 tret_cnt 와 다름).
// 최대 120건/회.
//
// KIS API: GET /uapi/overseas-futureoption/v1/quotations/opt-monthly-ccnl
func (c *Client) OptMonthlyCcnl(ctx context.Context, params CcnlParams) (*OptMonthlyCcnlData, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/overseas-futureoption/v1/quotations/opt-monthly-ccnl",
		TrID:     "HHDFO55020300",
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
	var res optMonthlyCcnlResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse OptMonthlyCcnl: %w", err)
	}
	return &OptMonthlyCcnlData{
		Output1: res.Output1,
		Output2: res.Output2,
	}, nil
}

// ─── EP4: OptDailyCcnl ───────────────────────────────────────────────────────

// OptDailyCcnlData 는 해외옵션 일간 체결추이 응답 (output1 메타 + output2 배열).
type OptDailyCcnlData struct {
	Output1 OptCcnlOutput1   `json:"output1"`
	Output2 []CcnlOutput2Item `json:"output2"`
}

type optDailyCcnlResponse struct {
	RtCd    string           `json:"rt_cd"`
	MsgCd   string           `json:"msg_cd"`
	Msg1    string           `json:"msg1"`
	Output1 OptCcnlOutput1   `json:"output1"`
	Output2 []CcnlOutput2Item `json:"output2"`
}

// OptDailyCcnl 는 해외옵션 일간 체결추이 조회 (HHDFO55020100).
//
// output1 카운트 필드: ret_cnt (선물 EP5 DailyCcnl 의 tret_cnt 와 다름).
// QRY_CNT 최대 119 (실제 조회 = 입력+1 개).
//
// KIS API: GET /uapi/overseas-futureoption/v1/quotations/opt-daily-ccnl
func (c *Client) OptDailyCcnl(ctx context.Context, params CcnlParams) (*OptDailyCcnlData, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/overseas-futureoption/v1/quotations/opt-daily-ccnl",
		TrID:     "HHDFO55020100",
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
	var res optDailyCcnlResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse OptDailyCcnl: %w", err)
	}
	return &OptDailyCcnlData{
		Output1: res.Output1,
		Output2: res.Output2,
	}, nil
}

// ─── EP5: OptWeeklyCcnl ──────────────────────────────────────────────────────

// OptWeeklyCcnlData 는 해외옵션 주간 체결추이 응답 (output1 메타 + output2 배열).
type OptWeeklyCcnlData struct {
	Output1 OptCcnlOutput1   `json:"output1"`
	Output2 []CcnlOutput2Item `json:"output2"`
}

type optWeeklyCcnlResponse struct {
	RtCd    string           `json:"rt_cd"`
	MsgCd   string           `json:"msg_cd"`
	Msg1    string           `json:"msg1"`
	Output1 OptCcnlOutput1   `json:"output1"`
	Output2 []CcnlOutput2Item `json:"output2"`
}

// OptWeeklyCcnl 는 해외옵션 주간 체결추이 조회 (HHDFO55020000).
//
// output1 카운트 필드: ret_cnt (선물 EP6 WeeklyCcnl 도 ret_cnt — 동일).
//
// KIS API: GET /uapi/overseas-futureoption/v1/quotations/opt-weekly-ccnl
func (c *Client) OptWeeklyCcnl(ctx context.Context, params CcnlParams) (*OptWeeklyCcnlData, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/overseas-futureoption/v1/quotations/opt-weekly-ccnl",
		TrID:     "HHDFO55020000",
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
	var res optWeeklyCcnlResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse OptWeeklyCcnl: %w", err)
	}
	return &OptWeeklyCcnlData{
		Output1: res.Output1,
		Output2: res.Output2,
	}, nil
}

// ─── EP6: OptTickCcnl ────────────────────────────────────────────────────────

// OptTickCcnlData 는 해외옵션 틱 체결추이 응답 (output1 메타 + output2 배열).
type OptTickCcnlData struct {
	Output1 OptCcnlOutput1   `json:"output1"`
	Output2 []CcnlOutput2Item `json:"output2"`
}

type optTickCcnlResponse struct {
	RtCd    string           `json:"rt_cd"`
	MsgCd   string           `json:"msg_cd"`
	Msg1    string           `json:"msg1"`
	Output1 OptCcnlOutput1   `json:"output1"`
	Output2 []CcnlOutput2Item `json:"output2"`
}

// OptTickCcnl 는 해외옵션 틱 체결추이 조회 (HHDFO55020200).
//
// output1 카운트 필드: ret_cnt (선물 EP7 TickCcnl 의 tret_cnt 와 다름).
// 최대 40건/회. 다음조회: QryTp="P", IndexKey=output1.IndexKey.
//
// KIS API: GET /uapi/overseas-futureoption/v1/quotations/opt-tick-ccnl
func (c *Client) OptTickCcnl(ctx context.Context, params CcnlParams) (*OptTickCcnlData, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/overseas-futureoption/v1/quotations/opt-tick-ccnl",
		TrID:     "HHDFO55020200",
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
	var res optTickCcnlResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse OptTickCcnl: %w", err)
	}
	return &OptTickCcnlData{
		Output1: res.Output1,
		Output2: res.Output2,
	}, nil
}

// ─── EP7: OptAskingPrice ─────────────────────────────────────────────────────

// OptAskingPriceOutput1 는 해외옵션 호가 현재 시세 요약 (10 필드).
// DISTINCT: 선물 InquireAskingPrice 대비 output1[5] 가 sttl_price(정산가) vs prev_price(전일종가).
type OptAskingPriceOutput1 struct {
	OpenPrice     decimal.Decimal `json:"open_price"`            // 시가
	HighPrice     decimal.Decimal `json:"high_price"`            // 고가
	LowpRice      decimal.Decimal `json:"lowp_rice"`             // 저가 (docs 오타 그대로 보존)
	LastPrice     decimal.Decimal `json:"last_price"`            // 현재가
	SttlPrice     decimal.Decimal `json:"sttl_price"`            // 정산가 (옵션 추가 필드 — 선물 prev_price 와 상이)
	Vol           int64           `json:"vol,string"`            // 거래량
	PrevDiffPrice decimal.Decimal `json:"prev_diff_price"`       // 전일대비가
	PrevDiffRate  float64         `json:"prev_diff_rate,string"` // 전일대비율
	QuotDate      string          `json:"quot_date"`             // 호가수신일자
	QuotTime      string          `json:"quot_time"`             // 호가수신시각
}

// OptAskingPriceOutput2Item 는 해외옵션 호가 레벨 (6 필드).
type OptAskingPriceOutput2Item struct {
	BidQntt  int64           `json:"bid_qntt,string"` // 매수수량
	BidNum   string          `json:"bid_num"`         // 매수번호
	BidPrice decimal.Decimal `json:"bid_price"`       // 매수호가
	AskQntt  int64           `json:"ask_qntt,string"` // 매도수량
	AskNum   string          `json:"ask_num"`         // 매도번호
	AskPrice decimal.Decimal `json:"ask_price"`       // 매도호가
}

// OptAskingPriceData 는 해외옵션 호가 응답 (output1 + output2[]).
type OptAskingPriceData struct {
	Output1 OptAskingPriceOutput1       `json:"output1"`
	Output2 []OptAskingPriceOutput2Item `json:"output2"`
}

type optAskingPriceResponse struct {
	RtCd    string                     `json:"rt_cd"`
	MsgCd   string                     `json:"msg_cd"`
	Msg1    string                     `json:"msg1"`
	Output1 OptAskingPriceOutput1      `json:"output1"`
	Output2 []OptAskingPriceOutput2Item `json:"output2"`
}

// OptAskingPrice 는 해외옵션 호가 조회 (HHDFO86000000).
//
// DISTINCT: 선물 InquireAskingPrice 대비 output1 에 sttl_price(정산가) 포함 (선물 prev_price 와 상이).
// lowp_rice 필드명 — docs 원본 오타 보존. output2 는 1~5호가 배열.
//
// code: 종목코드 (예: OESM24 C5340)
//
// KIS API: GET /uapi/overseas-futureoption/v1/quotations/opt-asking-price
func (c *Client) OptAskingPrice(ctx context.Context, code string) (*OptAskingPriceData, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/overseas-futureoption/v1/quotations/opt-asking-price",
		TrID:     "HHDFO86000000",
		CustType: "P",
		Query: map[string]string{
			"SRS_CD": code,
		},
	})
	if err != nil {
		return nil, err
	}
	var res optAskingPriceResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse OptAskingPrice: %w", err)
	}
	return &OptAskingPriceData{
		Output1: res.Output1,
		Output2: res.Output2,
	}, nil
}

// ─── EP8: OptDetail ──────────────────────────────────────────────────────────

// OptDetailOutput1 는 해외옵션 종목 상세 정보 (21 필드).
// DISTINCT: 선물 StockDetail 대비 sprd_srs_cd1/sprd_srs_cd2 없음 (23→21 필드).
// 주의: sttl_price 필드는 docs 에서 "정산가 X 전일종가 O" 명시 — 필드명과 실제 값 불일치.
type OptDetailOutput1 struct {
	ExchCd        string          `json:"exch_cd"`         // 거래소코드
	ClasCd        string          `json:"clas_cd"`         // 품목종류
	CrcCd         string          `json:"crc_cd"`          // 거래통화
	SttlPrice     decimal.Decimal `json:"sttl_price"`      // 전일종가 (★주의: 필드명 sttl_price 이나 실제 전일종가 수신)
	SttlDate      string          `json:"sttl_date"`       // 정산일
	TrstMgn       decimal.Decimal `json:"trst_mgn"`        // 증거금
	DispDigit     string          `json:"disp_digit"`      // 가격표시진법
	TickSz        decimal.Decimal `json:"tick_sz"`         // 틱사이즈
	TickVal       decimal.Decimal `json:"tick_val"`        // 틱가치
	MrktOpenDate  string          `json:"mrkt_open_date"`  // 장개시일자
	MrktOpenTime  string          `json:"mrkt_open_time"`  // 장개시시각
	MrktCloseDate string          `json:"mrkt_close_date"` // 장마감일자
	MrktCloseTime string          `json:"mrkt_close_time"` // 장마감시각
	TrdFrDate     string          `json:"trd_fr_date"`     // 상장일
	ExprDate      string          `json:"expr_date"`       // 만기일
	TrdToDate     string          `json:"trd_to_date"`     // 최종거래일
	RemnCnt       string          `json:"remn_cnt"`        // 잔존일수
	StatTp        string          `json:"stat_tp"`         // 매매여부
	CtrtSize      decimal.Decimal `json:"ctrt_size"`       // 계약크기
	StlTp         string          `json:"stl_tp"`          // 최종결제구분
	FrstNotiDate  string          `json:"frst_noti_date"`  // 최초식별일
}

// OptDetailData 는 해외옵션 종목 상세 응답 (output1 단일).
type OptDetailData struct {
	Output1 OptDetailOutput1 `json:"output1"`
}

type optDetailResponse struct {
	RtCd    string           `json:"rt_cd"`
	MsgCd   string           `json:"msg_cd"`
	Msg1    string           `json:"msg1"`
	Output1 OptDetailOutput1 `json:"output1"`
}

// OptDetail 는 해외옵션 종목 상세 조회 (HHDFO55010100).
//
// DISTINCT: 선물 StockDetail 대비 sprd_srs_cd1/2 없음 (21 필드).
// sttl_price 필드는 docs 주의: 전일종가 수신 (정산가 아님).
//
// code: 종목코드 (예: OESU24 C5500)
//
// KIS API: GET /uapi/overseas-futureoption/v1/quotations/opt-detail
func (c *Client) OptDetail(ctx context.Context, code string) (*OptDetailData, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/overseas-futureoption/v1/quotations/opt-detail",
		TrID:     "HHDFO55010100",
		CustType: "P",
		Query: map[string]string{
			"SRS_CD": code,
		},
	})
	if err != nil {
		return nil, err
	}
	var res optDetailResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse OptDetail: %w", err)
	}
	return &OptDetailData{Output1: res.Output1}, nil
}

// ─── EP9: OptPrice ───────────────────────────────────────────────────────────

// OptPriceOutput1 는 해외옵션 현재가 상세 정보 (31 필드).
// DISTINCT: 선물 InquirePrice 대비 필드 순서/구성 상이. proc_date/proc_time 앞에 위치.
// sttl_price Required=N (Optional).
type OptPriceOutput1 struct {
	ProcDate      string          `json:"proc_date"`             // 최종처리일자
	ProcTime      string          `json:"proc_time"`             // 최종처리시각
	OpenPrice     decimal.Decimal `json:"open_price"`            // 시가
	HighPrice     decimal.Decimal `json:"high_price"`            // 고가
	LowPrice      decimal.Decimal `json:"low_price"`             // 저가
	LastPrice     decimal.Decimal `json:"last_price"`            // 현재가
	Vol           int64           `json:"vol,string"`            // 누적거래수량
	PrevDiffFlag  string          `json:"prev_diff_flag"`        // 전일대비구분
	PrevDiffPrice decimal.Decimal `json:"prev_diff_price"`       // 전일대비가격
	PrevDiffRate  float64         `json:"prev_diff_rate,string"` // 전일대비율
	BidQntt       int64           `json:"bid_qntt,string"`       // 매수1수량
	BidPrice      decimal.Decimal `json:"bid_price"`             // 매수1호가
	AskQntt       int64           `json:"ask_qntt,string"`       // 매도1수량
	AskPrice      decimal.Decimal `json:"ask_price"`             // 매도1호가
	TrstMgn       decimal.Decimal `json:"trst_mgn"`              // 증거금
	ExchCd        string          `json:"exch_cd"`               // 거래소코드
	CrcCd         string          `json:"crc_cd"`                // 거래통화
	TrdFrDate     string          `json:"trd_fr_date"`           // 상장일
	ExprDate      string          `json:"expr_date"`             // 만기일
	TrdToDate     string          `json:"trd_to_date"`           // 최종거래일
	RemnCnt       string          `json:"remn_cnt"`              // 잔존일수
	LastQntt      int64           `json:"last_qntt,string"`      // 체결량
	TotAskQntt    int64           `json:"tot_ask_qntt,string"`   // 총매도잔량
	TotBidQntt    int64           `json:"tot_bid_qntt,string"`   // 총매수잔량
	TickSize      decimal.Decimal `json:"tick_size"`             // 틱사이즈
	OpenDate      string          `json:"open_date"`             // 장개시일자
	OpenTime      string          `json:"open_time"`             // 장개시시각
	CloseDate     string          `json:"close_date"`            // 장종료일자
	CloseTime     string          `json:"close_time"`            // 장종료시각
	Sbsnsdate     string          `json:"sbsnsdate"`             // 영업일자
	SttlPrice     decimal.Decimal `json:"sttl_price"`            // 정산가 (Optional — Required=N)
}

// OptPriceData 는 해외옵션 현재가 응답 (output1 단일).
type OptPriceData struct {
	Output1 OptPriceOutput1 `json:"output1"`
}

type optPriceResponse struct {
	RtCd    string          `json:"rt_cd"`
	MsgCd   string          `json:"msg_cd"`
	Msg1    string          `json:"msg1"`
	Output1 OptPriceOutput1 `json:"output1"`
}

// OptPrice 는 해외옵션 현재가 조회 (HHDFO55010000).
//
// DISTINCT: 선물 InquirePrice 대비 필드 순서/구성 상이 (proc_date/proc_time 앞에 위치).
// sttl_price Required=N (Optional).
//
// code: 종목코드 (예: OESU24 C5500)
//
// KIS API: GET /uapi/overseas-futureoption/v1/quotations/opt-price
func (c *Client) OptPrice(ctx context.Context, code string) (*OptPriceData, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/overseas-futureoption/v1/quotations/opt-price",
		TrID:     "HHDFO55010000",
		CustType: "P",
		Query: map[string]string{
			"SRS_CD": code,
		},
	})
	if err != nil {
		return nil, err
	}
	var res optPriceResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse OptPrice: %w", err)
	}
	return &OptPriceData{Output1: res.Output1}, nil
}
