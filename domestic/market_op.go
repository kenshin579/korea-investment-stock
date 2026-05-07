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

// ─── EP5: InquireChkHoliday ──────────────────────────────────────────────────

// InquireChkHolidayParams 는 휴장일 조회 파라미터.
// 주의: 파라미터명이 FID_ 접두어 없는 비표준 UPPERCASE 형식 (BASS_DT / CTX_AREA_NK / CTX_AREA_FK).
// 주의: 단시간 다수 호출 자제 (KIS docs 권장 1일 1회).
type InquireChkHolidayParams struct {
	BassDt    string // BASS_DT (Y): 조회기준일 YYYYMMDD
	CtxAreaNk string // CTX_AREA_NK (Y): 연속조회검색조건 (공란 가능)
	CtxAreaFk string // CTX_AREA_FK (Y): 연속조회키 (공란 가능)
}

// ChkHolidayItem 은 휴장일 조회 단일 응답 객체.
// wire key "bass_dt" → Go 필드 Bassdt.
type ChkHolidayItem struct {
	Bassdt     string `json:"bass_dt"`      // 기준일자 YYYYMMDD
	WdayDvsnCd string `json:"wday_dvsn_cd"` // 요일구분코드 01(일)~07(토)
	BzdyYn     string `json:"bzdy_yn"`      // 영업일여부 Y/N
	TrDayYn    string `json:"tr_day_yn"`    // 거래일여부 Y/N
	OpndYn     string `json:"opnd_yn"`      // 개장일여부 Y/N
	SttlDayYn  string `json:"sttl_day_yn"`  // 결제일여부 Y/N
}

// InquireChkHolidayResponse 는 휴장일 조회 응답.
type InquireChkHolidayResponse struct {
	RtCd   string          `json:"rt_cd"`
	MsgCd  string          `json:"msg_cd"`
	Msg1   string          `json:"msg1"`
	Output *ChkHolidayItem `json:"output"`
}

// InquireChkHoliday 는 휴장일을 조회한다 (CTCA0903R).
//
// 주의: 단시간 다수 호출 자제 (KIS docs 권장 1일 1회).
// 파라미터명이 FID_ 접두어 없는 비표준 UPPERCASE 형식임에 유의 (BASS_DT/CTX_AREA_NK/CTX_AREA_FK).
func (c *Client) InquireChkHoliday(ctx context.Context, params InquireChkHolidayParams) (*InquireChkHolidayResponse, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/chk-holiday",
		TrID:   "CTCA0903R",
		Query: map[string]string{
			"BASS_DT":     params.BassDt,
			"CTX_AREA_NK": params.CtxAreaNk,
			"CTX_AREA_FK": params.CtxAreaFk,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res InquireChkHolidayResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InquireChkHolidayResponse: %w", err)
	}
	return &res, nil
}

// ─── EP6: InquireViStatus ────────────────────────────────────────────────────

// InquireViStatusParams 는 변동성완화장치(VI) 현황 조회 파라미터.
type InquireViStatusParams struct {
	DivClsCode      string // FID_DIV_CLS_CODE (Y): 0=전체/1=상승/2=하락
	CondScrDivCode  string // FID_COND_SCR_DIV_CODE: 기본 "20139" (hardcoded)
	MrktClsCode     string // FID_MRKT_CLS_CODE (Y): 0=전체/K=거래소/Q=코스닥
	Symbol          string // FID_INPUT_ISCD (Y, 공란 가능)
	RankSortClsCode string // FID_RANK_SORT_CLS_CODE (Y): 0=전체/1=정적/2=동적/3=정적&동적
	InputDate1      string // FID_INPUT_DATE_1 (Y): YYYYMMDD
	TrgtClsCode     string // FID_TRGT_CLS_CODE (Y, 공란 가능)
	TrgtExlsCode    string // FID_TRGT_EXLS_CLS_CODE (Y, 공란 가능)
}

// ViStatusOutput 는 변동성완화장치(VI) 현황 응답 단일 객체.
//
// KIS 공식 문서는 output 을 단일 Object({})로 선언하나, 실제 응답에서 배열([])을 반환할 수 있음.
// ("30건" 등 복수 건 문구 포함). 실 API 호출 시 배열 반환 확인 시 []ViStatusOutput 로 전환 필요.
type ViStatusOutput struct {
	HtsKorIsnm   string          `json:"hts_kor_isnm"`
	MkscShrnIscd string          `json:"mksc_shrn_iscd"`
	ViClsCode    string          `json:"vi_cls_code"`  // Y=발동/N=해제
	BsopDate     string          `json:"bsop_date"`    // YYYYMMDD
	CntgViHour   string          `json:"cntg_vi_hour"` // HHMMSS
	ViCnclHour   string          `json:"vi_cncl_hour"` // HHMMSS
	ViKindCode   string          `json:"vi_kind_code"` // 1=정적/2=동적/3=정적&동적
	ViPrc        decimal.Decimal `json:"vi_prc"`
	ViStndPrc    decimal.Decimal `json:"vi_stnd_prc"`
	ViDprt       float64         `json:"vi_dprt,string"`
	ViDmcStndPrc decimal.Decimal `json:"vi_dmc_stnd_prc"`
	ViDmcDprt    float64         `json:"vi_dmc_dprt,string"`
	ViCount      int64           `json:"vi_count,string"`
}

// InquireViStatusResponse 는 변동성완화장치(VI) 현황 응답.
type InquireViStatusResponse struct {
	RtCd   string          `json:"rt_cd"`
	MsgCd  string          `json:"msg_cd"`
	Msg1   string          `json:"msg1"`
	Output *ViStatusOutput `json:"output"`
}

// InquireViStatus 는 변동성완화장치(VI) 현황을 조회한다 (FHPST01390000).
//
// KIS 문서가 output 을 단일 Object로 선언. 실 API 에서 배열 반환 시 struct 변경 필요.
// FID_COND_SCR_DIV_CODE 는 "20139" 로 hardcoded.
func (c *Client) InquireViStatus(ctx context.Context, params InquireViStatusParams) (*InquireViStatusResponse, error) {
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "20139"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-vi-status",
		TrID:   "FHPST01390000",
		Query: map[string]string{
			"FID_DIV_CLS_CODE":       params.DivClsCode,
			"FID_COND_SCR_DIV_CODE":  scrDiv,
			"FID_MRKT_CLS_CODE":      params.MrktClsCode,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_RANK_SORT_CLS_CODE": params.RankSortClsCode,
			"FID_INPUT_DATE_1":       params.InputDate1,
			"FID_TRGT_CLS_CODE":      params.TrgtClsCode,
			"FID_TRGT_EXLS_CLS_CODE": params.TrgtExlsCode,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res InquireViStatusResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InquireViStatusResponse: %w", err)
	}
	return &res, nil
}
