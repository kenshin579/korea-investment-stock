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

// ─── EP10: InquirePrice ──────────────────────────────────────────────────────

// InquirePriceOutput1 는 해외선물 현재가 상세 정보 (31 필드).
type InquirePriceOutput1 struct {
	ProcDate      string          `json:"proc_date"`           // 최종처리일자
	HighPrice     decimal.Decimal `json:"high_price"`          // 고가
	ProcTime      string          `json:"proc_time"`           // 최종처리시각
	OpenPrice     decimal.Decimal `json:"open_price"`          // 시가
	TrstMgn       decimal.Decimal `json:"trst_mgn"`            // 증거금
	LowPrice      decimal.Decimal `json:"low_price"`           // 저가
	LastPrice     decimal.Decimal `json:"last_price"`          // 현재가
	Vol           int64           `json:"vol,string"`          // 누적거래수량
	PrevDiffFlag  string          `json:"prev_diff_flag"`      // 전일대비구분
	PrevDiffPrice decimal.Decimal `json:"prev_diff_price"`     // 전일대비가격
	PrevDiffRate  kistypes.Float  `json:"prev_diff_rate"`      // 전일대비율
	BidQntt       int64           `json:"bid_qntt,string"`     // 매수1수량
	BidPrice      decimal.Decimal `json:"bid_price"`           // 매수1호가
	AskQntt       int64           `json:"ask_qntt,string"`     // 매도1수량
	AskPrice      decimal.Decimal `json:"ask_price"`           // 매도1호가
	PrevPrice     decimal.Decimal `json:"prev_price"`          // 전일종가
	ExchCd        string          `json:"exch_cd"`             // 거래소코드
	CrcCd         string          `json:"crc_cd"`              // 거래통화
	TrdFrDate     string          `json:"trd_fr_date"`         // 상장일
	ExprDate      string          `json:"expr_date"`           // 만기일
	TrdToDate     string          `json:"trd_to_date"`         // 최종거래일
	RemnCnt       string          `json:"remn_cnt"`            // 잔존일수
	LastQntt      int64           `json:"last_qntt,string"`    // 체결량
	TotAskQntt    int64           `json:"tot_ask_qntt,string"` // 총매도잔량
	TotBidQntt    int64           `json:"tot_bid_qntt,string"` // 총매수잔량
	TickSize      decimal.Decimal `json:"tick_size"`           // 틱사이즈
	OpenDate      string          `json:"open_date"`           // 장개시일자
	OpenTime      string          `json:"open_time"`           // 장개시시각
	CloseDate     string          `json:"close_date"`          // 장종료일자
	CloseTime     string          `json:"close_time"`          // 장종료시각
	Sbsnsdate     string          `json:"sbsnsdate"`           // 영업일자
	SttlPrice     decimal.Decimal `json:"sttl_price"`          // 정산가
}

// InquirePriceData 는 해외선물 현재가 응답 (output1 단일).
type InquirePriceData struct {
	Output1 InquirePriceOutput1 `json:"output1"`
}

type inquirePriceResponse struct {
	RtCd    string              `json:"rt_cd"`
	MsgCd   string              `json:"msg_cd"`
	Msg1    string              `json:"msg1"`
	Output1 InquirePriceOutput1 `json:"output1"`
}

// InquirePrice 는 해외선물 현재가 조회 (HHDFC55010000).
//
// code: 종목코드 (예: CNHU24)
//
// KIS API: GET /uapi/overseas-futureoption/v1/quotations/inquire-price
func (c *Client) InquirePrice(ctx context.Context, code string) (*InquirePriceData, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/overseas-futureoption/v1/quotations/inquire-price",
		TrID:     "HHDFC55010000",
		CustType: "P",
		Query: map[string]string{
			"SRS_CD": code,
		},
	})
	if err != nil {
		return nil, err
	}
	var res inquirePriceResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InquirePrice: %w", err)
	}
	return &InquirePriceData{Output1: res.Output1}, nil
}

// ─── EP9: StockDetail ────────────────────────────────────────────────────────

// StockDetailOutput1 는 해외선물 종목 상세 정보 (23 필드).
type StockDetailOutput1 struct {
	ExchCd        string          `json:"exch_cd"`         // 거래소코드
	TickSz        decimal.Decimal `json:"tick_sz"`         // 틱사이즈
	DispDigit     string          `json:"disp_digit"`      // 가격표시진법
	TrstMgn       decimal.Decimal `json:"trst_mgn"`        // 증거금
	SttlDate      string          `json:"sttl_date"`       // 정산일
	PrevPrice     decimal.Decimal `json:"prev_price"`      // 전일종가
	CrcCd         string          `json:"crc_cd"`          // 거래통화
	ClasCd        string          `json:"clas_cd"`         // 품목종류
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
	SprdSrsCd1    string          `json:"sprd_srs_cd1"`    // 스프레드 종목 #1
	SprdSrsCd2    string          `json:"sprd_srs_cd2"`    // 스프레드 종목 #2
}

// StockDetailData 는 해외선물 종목 상세 응답 (output1 단일).
type StockDetailData struct {
	Output1 StockDetailOutput1 `json:"output1"`
}

type stockDetailResponse struct {
	RtCd    string             `json:"rt_cd"`
	MsgCd   string             `json:"msg_cd"`
	Msg1    string             `json:"msg1"`
	Output1 StockDetailOutput1 `json:"output1"`
}

// StockDetail 는 해외선물 종목 상세 조회 (HHDFC55010100).
//
// code: 종목코드 (예: CNHU24)
//
// KIS API: GET /uapi/overseas-futureoption/v1/quotations/stock-detail
func (c *Client) StockDetail(ctx context.Context, code string) (*StockDetailData, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/overseas-futureoption/v1/quotations/stock-detail",
		TrID:     "HHDFC55010100",
		CustType: "P",
		Query: map[string]string{
			"SRS_CD": code,
		},
	})
	if err != nil {
		return nil, err
	}
	var res stockDetailResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse StockDetail: %w", err)
	}
	return &StockDetailData{Output1: res.Output1}, nil
}

// ─── EP3: SearchContractDetail ───────────────────────────────────────────────

// SearchContractDetailOutput2Item 는 해외선물 상품기본정보 항목 (22 필드).
type SearchContractDetailOutput2Item struct {
	ExchCd        string          `json:"exch_cd"`         // 거래소코드
	ClasCd        string          `json:"clas_cd"`         // 품목종류
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
	SubExchNm     string          `json:"sub_exch_nm"`     // 서브거래소코드
}

// SearchContractDetailData 는 해외선물 상품기본정보 응답 (output2 배열).
type SearchContractDetailData struct {
	Output2 []SearchContractDetailOutput2Item `json:"output2"`
}

type searchContractDetailResponse struct {
	RtCd    string                            `json:"rt_cd"`
	MsgCd   string                            `json:"msg_cd"`
	Msg1    string                            `json:"msg1"`
	Output2 []SearchContractDetailOutput2Item `json:"output2"`
}

// SearchContractDetailParams 는 해외선물 상품기본정보 조회 파라미터.
type SearchContractDetailParams struct {
	Codes []string // 종목코드 목록 (최대 32개). QRY_CNT = len(Codes).
}

// SearchContractDetail 는 해외선물 상품기본정보 조회 (HHDFC55200000).
//
// 한번에 최대 32개 종목 조회. output2 배열만 반환 (output1 없음).
//
// KIS API: GET /uapi/overseas-futureoption/v1/quotations/search-contract-detail
func (c *Client) SearchContractDetail(ctx context.Context, params SearchContractDetailParams) (*SearchContractDetailData, error) {
	query := map[string]string{
		"QRY_CNT": fmt.Sprintf("%d", len(params.Codes)),
	}
	for i, code := range params.Codes {
		key := fmt.Sprintf("SRS_CD_%02d", i+1)
		query[key] = code
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/overseas-futureoption/v1/quotations/search-contract-detail",
		TrID:     "HHDFC55200000",
		CustType: "P",
		Query:    query,
	})
	if err != nil {
		return nil, err
	}
	var res searchContractDetailResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse SearchContractDetail: %w", err)
	}
	return &SearchContractDetailData{Output2: res.Output2}, nil
}

// ─── EP8: InquireAskingPrice ─────────────────────────────────────────────────

// InquireAskingPriceOutput1 는 해외선물 호가 현재 시세 요약 (10 필드).
type InquireAskingPriceOutput1 struct {
	OpenPrice     decimal.Decimal `json:"open_price"`      // 시가
	HighPrice     decimal.Decimal `json:"high_price"`      // 고가
	LowpRice      decimal.Decimal `json:"lowp_rice"`       // 저가 (docs 오타 그대로 보존)
	LastPrice     decimal.Decimal `json:"last_price"`      // 현재가
	PrevPrice     decimal.Decimal `json:"prev_price"`      // 전일종가
	Vol           int64           `json:"vol,string"`      // 거래량
	PrevDiffPrice decimal.Decimal `json:"prev_diff_price"` // 전일대비가
	PrevDiffRate  kistypes.Float  `json:"prev_diff_rate"`  // 전일대비율
	QuotDate      string          `json:"quot_date"`       // 호가수신일자
	QuotTime      string          `json:"quot_time"`       // 호가수신시각
}

// InquireAskingPriceOutput2Item 는 해외선물 호가 레벨 (6 필드).
type InquireAskingPriceOutput2Item struct {
	BidQntt  int64           `json:"bid_qntt,string"` // 매수수량
	BidNum   string          `json:"bid_num"`         // 매수번호
	BidPrice decimal.Decimal `json:"bid_price"`       // 매수호가
	AskQntt  int64           `json:"ask_qntt,string"` // 매도수량
	AskNum   string          `json:"ask_num"`         // 매도번호
	AskPrice decimal.Decimal `json:"ask_price"`       // 매도호가
}

// InquireAskingPriceData 는 해외선물 호가 응답 (output1 + output2[]).
type InquireAskingPriceData struct {
	Output1 InquireAskingPriceOutput1       `json:"output1"`
	Output2 []InquireAskingPriceOutput2Item `json:"output2"`
}

type inquireAskingPriceResponse struct {
	RtCd    string                          `json:"rt_cd"`
	MsgCd   string                          `json:"msg_cd"`
	Msg1    string                          `json:"msg1"`
	Output1 InquireAskingPriceOutput1       `json:"output1"`
	Output2 []InquireAskingPriceOutput2Item `json:"output2"`
}

// InquireAskingPrice 는 해외선물 호가 조회 (HHDFC86000000).
//
// code: 종목코드
//
// KIS API: GET /uapi/overseas-futureoption/v1/quotations/inquire-asking-price
func (c *Client) InquireAskingPrice(ctx context.Context, code string) (*InquireAskingPriceData, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/overseas-futureoption/v1/quotations/inquire-asking-price",
		TrID:     "HHDFC86000000",
		CustType: "P",
		Query: map[string]string{
			"SRS_CD": code,
		},
	})
	if err != nil {
		return nil, err
	}
	var res inquireAskingPriceResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InquireAskingPrice: %w", err)
	}
	return &InquireAskingPriceData{
		Output1: res.Output1,
		Output2: res.Output2,
	}, nil
}
