// File: domestic/program_trade.go
package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// ProgramTradeByStockDaily 는 종목별 프로그램매매추이(일별) (FHPPG04650201) 응답.
//
// 한투 docs: docs/api/국내주식/종목별_프로그램매매추이(일별).md
// path: /uapi/domestic-stock/v1/quotations/program-trade-by-stock-daily
type ProgramTradeByStockDaily struct {
	Output []ProgramTradeByStockDailyItem `json:"output"`
}

// ProgramTradeByStockDailyItem 은 응답 output 한 행 (일자별 프로그램 매매).
//
// WholNtbyTrPbmnIcdc2: 마지막 필드. 필드명 trailing "2" — EP4 의 WholNtbyTrPbmnIcdc 와 구분.
type ProgramTradeByStockDailyItem struct {
	StckBsopDate        string          `json:"stck_bsop_date"`                 // 주식 영업 일자
	StckClpr            decimal.Decimal `json:"stck_clpr"`                      // 주식 종가
	PrdyVrss            decimal.Decimal `json:"prdy_vrss"`                      // 전일 대비
	PrdyVrssSign        string          `json:"prdy_vrss_sign"`                 // 전일 대비 부호
	PrdyCtrt            float64         `json:"prdy_ctrt,string"`               // 전일 대비율
	AcmlVol             int64           `json:"acml_vol,string"`                // 누적 거래량
	AcmlTrPbmn          int64           `json:"acml_tr_pbmn,string"`            // 누적 거래 대금
	WholSmtnSelnVol     int64           `json:"whol_smtn_seln_vol,string"`      // 전체 합산 매도 수량
	WholSmtnShnuVol     int64           `json:"whol_smtn_shnu_vol,string"`      // 전체 합산 매수 수량
	WholSmtnNtbyQty     int64           `json:"whol_smtn_ntby_qty,string"`      // 전체 합산 순매수 수량
	WholSmtnSelnTrPbmn  int64           `json:"whol_smtn_seln_tr_pbmn,string"`  // 전체 합산 매도 거래대금
	WholSmtnShnuTrPbmn  int64           `json:"whol_smtn_shnu_tr_pbmn,string"`  // 전체 합산 매수 거래대금
	WholSmtnNtbyTrPbmn  int64           `json:"whol_smtn_ntby_tr_pbmn,string"`  // 전체 합산 순매수 거래대금
	WholNtbyVolIcdc     int64           `json:"whol_ntby_vol_icdc,string"`      // 전체 순매수 거래량 증감
	WholNtbyTrPbmnIcdc2 int64           `json:"whol_ntby_tr_pbmn_icdc2,string"` // 전체 순매수 거래대금 증감 (trailing "2")
}

// InquireProgramTradeByStockDailyParams 는 종목별 프로그램매매추이(일별) 조회 파라미터.
//
// BaseDate: KIS docs 예시가 "002" prefix 포함 ("0020240308") — 호출자가 raw string 그대로 전달.
// MarketCode: "J"(KRX), "NX"(NXT), "UN"(통합). 빈 값=>"J".
type InquireProgramTradeByStockDailyParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE — 빈 값=>"J"
	Symbol     string // FID_INPUT_ISCD — 필수
	BaseDate   string // FID_INPUT_DATE_1 — 필수, KIS docs 예시: "0020240308" ("002" prefix)
}

// InquireProgramTradeByStockDaily 는 종목별 프로그램매매추이(일별) 호출.
//
// 한투 docs: docs/api/국내주식/종목별_프로그램매매추이(일별).md
// path: /uapi/domestic-stock/v1/quotations/program-trade-by-stock-daily (FHPPG04650201)
func (c *Client) InquireProgramTradeByStockDaily(ctx context.Context, params InquireProgramTradeByStockDailyParams) (*ProgramTradeByStockDaily, error) {
	mkt := params.MarketCode
	if mkt == "" {
		mkt = "J"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/program-trade-by-stock-daily",
		TrID:   "FHPPG04650201",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": mkt,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_INPUT_DATE_1":       params.BaseDate,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res ProgramTradeByStockDaily
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse ProgramTradeByStockDaily: %w", err)
	}
	return &res, nil
}

// ProgramTradeByStock 는 종목별 프로그램매매추이(체결) (FHPPG04650101) 응답.
//
// 한투 docs: docs/api/국내주식/종목별_프로그램매매추이(체결).md
// path: /uapi/domestic-stock/v1/quotations/program-trade-by-stock
type ProgramTradeByStock struct {
	Output []ProgramTradeByStockItem `json:"output"`
}

// ProgramTradeByStockItem 은 응답 output 한 행 (시간대별 프로그램 체결).
//
// WholNtbyTrPbmnIcdc: 마지막 필드. trailing "2" 없음 — EP3 의 WholNtbyTrPbmnIcdc2 와 구분.
type ProgramTradeByStockItem struct {
	BsopHour           string          `json:"bsop_hour"`                     // 영업 시간
	StckPrpr           decimal.Decimal `json:"stck_prpr"`                     // 주식 현재가
	PrdyVrss           decimal.Decimal `json:"prdy_vrss"`                     // 전일 대비
	PrdyVrssSign       string          `json:"prdy_vrss_sign"`                // 전일 대비 부호
	PrdyCtrt           float64         `json:"prdy_ctrt,string"`              // 전일 대비율
	AcmlVol            int64           `json:"acml_vol,string"`               // 누적 거래량
	WholSmtnSelnVol    int64           `json:"whol_smtn_seln_vol,string"`     // 전체 합산 매도 수량
	WholSmtnShnuVol    int64           `json:"whol_smtn_shnu_vol,string"`     // 전체 합산 매수 수량
	WholSmtnNtbyQty    int64           `json:"whol_smtn_ntby_qty,string"`     // 전체 합산 순매수 수량
	WholSmtnSelnTrPbmn int64           `json:"whol_smtn_seln_tr_pbmn,string"` // 전체 합산 매도 거래대금
	WholSmtnShnuTrPbmn int64           `json:"whol_smtn_shnu_tr_pbmn,string"` // 전체 합산 매수 거래대금
	WholSmtnNtbyTrPbmn int64           `json:"whol_smtn_ntby_tr_pbmn,string"` // 전체 합산 순매수 거래대금
	WholNtbyVolIcdc    int64           `json:"whol_ntby_vol_icdc,string"`     // 전체 순매수 거래량 증감
	WholNtbyTrPbmnIcdc int64           `json:"whol_ntby_tr_pbmn_icdc,string"` // 전체 순매수 거래대금 증감 (no trailing "2")
}

// InquireProgramTradeByStockParams 는 종목별 프로그램매매추이(체결) 조회 파라미터.
type InquireProgramTradeByStockParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE — 빈 값=>"J"
	Symbol     string // FID_INPUT_ISCD — 필수
}

// InquireProgramTradeByStock 는 종목별 프로그램매매추이(체결) 호출.
//
// 한투 docs: docs/api/국내주식/종목별_프로그램매매추이(체결).md
// path: /uapi/domestic-stock/v1/quotations/program-trade-by-stock (FHPPG04650101)
func (c *Client) InquireProgramTradeByStock(ctx context.Context, params InquireProgramTradeByStockParams) (*ProgramTradeByStock, error) {
	mkt := params.MarketCode
	if mkt == "" {
		mkt = "J"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/program-trade-by-stock",
		TrID:   "FHPPG04650101",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": mkt,
			"FID_INPUT_ISCD":         params.Symbol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res ProgramTradeByStock
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse ProgramTradeByStock: %w", err)
	}
	return &res, nil
}

// CompProgramTradeToday 는 프로그램매매 종합현황(시간) (FHPPG04600101) 응답.
//
// 한투 docs: docs/api/국내주식/프로그램매매_종합현황(시간).md
// path: /uapi/domestic-stock/v1/quotations/comp-program-trade-today
//
// output1 (시간별 Array). 18 fields.
// 비고: `smtm`(rate 계열)과 `smtn`(금액/수량 계열) 혼용 — KIS API 원문 그대로.
// `shun` typo 2개 (arbt_smtm_shun_tr_pbmn_rate, nabt_smtm_shun_tr_pbmn_rate) 원문 보존.
type CompProgramTradeToday struct {
	Output1 []CompProgramTradeTodayItem `json:"output1"`
}

// CompProgramTradeTodayItem 은 응답 output1 한 행 (시간대별 종합현황).
//
// 필드명 내 typo 주의:
// - ArbtSmtmShunTrPbmnRate: "shun" (KIS docs typo) — 매수 율이지만 필드명 오기
// - NabtSmtmShunTrPbmnRate: 동일 패턴
type CompProgramTradeTodayItem struct {
	BsopHour               string          `json:"bsop_hour"`                          // 영업 시간
	ArbtSmtnSelnTrPbmn     int64           `json:"arbt_smtn_seln_tr_pbmn,string"`      // 차익 합산 매도 거래대금
	ArbtSmtmSelnTrPbmnRate float64         `json:"arbt_smtm_seln_tr_pbmn_rate,string"` // 차익 합산 매도 거래대금 비율
	ArbtSmtnShnuTrPbmn     int64           `json:"arbt_smtn_shnu_tr_pbmn,string"`      // 차익 합산 매수 거래대금
	ArbtSmtmShunTrPbmnRate float64         `json:"arbt_smtm_shun_tr_pbmn_rate,string"` // 차익 합산 매수 거래대금 비율 ("shun" typo)
	NabtSmtnSelnTrPbmn     int64           `json:"nabt_smtn_seln_tr_pbmn,string"`      // 비차익 합산 매도 거래대금
	NabtSmtmSelnTrPbmnRate float64         `json:"nabt_smtm_seln_tr_pbmn_rate,string"` // 비차익 합산 매도 거래대금 비율
	NabtSmtnShnuTrPbmn     int64           `json:"nabt_smtn_shnu_tr_pbmn,string"`      // 비차익 합산 매수 거래대금
	NabtSmtmShunTrPbmnRate float64         `json:"nabt_smtm_shun_tr_pbmn_rate,string"` // 비차익 합산 매수 거래대금 비율 ("shun" typo)
	ArbtSmtnNtbyTrPbmn     int64           `json:"arbt_smtn_ntby_tr_pbmn,string"`      // 차익 합산 순매수 거래대금
	ArbtSmtmNtbyTrPbmnRate float64         `json:"arbt_smtm_ntby_tr_pbmn_rate,string"` // 차익 합산 순매수 거래대금 비율
	NabtSmtnNtbyTrPbmn     int64           `json:"nabt_smtn_ntby_tr_pbmn,string"`      // 비차익 합산 순매수 거래대금
	NabtSmtmNtbyTrPbmnRate float64         `json:"nabt_smtm_ntby_tr_pbmn_rate,string"` // 비차익 합산 순매수 거래대금 비율
	WholSmtnNtbyTrPbmn     int64           `json:"whol_smtn_ntby_tr_pbmn,string"`      // 전체 합산 순매수 거래대금
	WholNtbyTrPbmnRate     float64         `json:"whol_ntby_tr_pbmn_rate,string"`      // 전체 순매수 거래대금 비율
	BstpNmixPrpr           decimal.Decimal `json:"bstp_nmix_prpr"`                     // 업종 지수 현재가
	BstpNmixPrdyVrss       decimal.Decimal `json:"bstp_nmix_prdy_vrss"`                // 업종 지수 전일 대비
	PrdyVrssSign           string          `json:"prdy_vrss_sign"`                     // 전일 대비 부호
}

// InquireCompProgramTradeTodayParams 는 프로그램매매 종합현황(시간) 조회 파라미터.
//
// 6개 query 파라미터 중 첫 2개만 의미있음. 나머지 4개는 빈 문자열 전송.
// MrktClsCode: "K"=코스피, "Q"=코스닥.
type InquireCompProgramTradeTodayParams struct {
	MarketCode  string // FID_COND_MRKT_DIV_CODE — 필수
	MrktClsCode string // FID_MRKT_CLS_CODE — K:코스피, Q:코스닥
	// 아래 4개는 항상 "" 전송 (KIS docs 명시)
	// SctnClsCode FID_SCTN_CLS_CODE
	// Symbol      FID_INPUT_ISCD
	// MarketCode1 FID_COND_MRKT_DIV_CODE1
	// InputHour1  FID_INPUT_HOUR_1
}

// InquireCompProgramTradeToday 는 프로그램매매 종합현황(시간) 호출.
//
// 한투 docs: docs/api/국내주식/프로그램매매_종합현황(시간).md
// path: /uapi/domestic-stock/v1/quotations/comp-program-trade-today (FHPPG04600101)
func (c *Client) InquireCompProgramTradeToday(ctx context.Context, params InquireCompProgramTradeTodayParams) (*CompProgramTradeToday, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/comp-program-trade-today",
		TrID:   "FHPPG04600101",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE":  params.MarketCode,
			"FID_MRKT_CLS_CODE":       params.MrktClsCode,
			"FID_SCTN_CLS_CODE":       "",
			"FID_INPUT_ISCD":          "",
			"FID_COND_MRKT_DIV_CODE1": "",
			"FID_INPUT_HOUR_1":        "",
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res CompProgramTradeToday
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse CompProgramTradeToday: %w", err)
	}
	return &res, nil
}

// CompProgramTradeDailyItem 은 프로그램매매 종합현황(일별) 한 행을 나타낸다.
// FHPPG04600001 output 배열 원소. 8개월 lookback 한도.
// shun 타이포 필드(arbt_smtm_shun_vol_rate 등)는 KIS docs 원문 보존.
type CompProgramTradeDailyItem struct {
	StckBsopDate           string  `json:"stck_bsop_date"`                     // 주식 영업 일자
	NabtEntmSelnTrPbmn     int64   `json:"nabt_entm_seln_tr_pbmn,string"`      // 비차익 위탁 매도 거래대금
	NabtOnslSelnVol        int64   `json:"nabt_onsl_seln_vol,string"`          // 비차익 자기 매도 거래량
	WholOnslSelnTrPbmn     int64   `json:"whol_onsl_seln_tr_pbmn,string"`      // 전체 자기 매도 거래대금
	ArbtSmtnShnuVol        int64   `json:"arbt_smtn_shnu_vol,string"`          // 차익 합계 매수 거래량
	NabtSmtnShnuTrPbmn     int64   `json:"nabt_smtn_shnu_tr_pbmn,string"`      // 비차익 합계 매수 거래대금
	ArbtEntmNtbyQty        int64   `json:"arbt_entm_ntby_qty,string"`          // 차익 위탁 순매수 수량
	NabtEntmNtbyTrPbmn     int64   `json:"nabt_entm_ntby_tr_pbmn,string"`      // 비차익 위탁 순매수 거래대금
	ArbtEntmSelnVol        int64   `json:"arbt_entm_seln_vol,string"`          // 차익 위탁 매도 거래량
	NabtEntmSelnVolRate    float64 `json:"nabt_entm_seln_vol_rate,string"`     // 비차익 위탁 매도 거래량 비율
	NabtOnslSelnVolRate    float64 `json:"nabt_onsl_seln_vol_rate,string"`     // 비차익 자기 매도 거래량 비율
	WholOnslSelnTrPbmnRate float64 `json:"whol_onsl_seln_tr_pbmn_rate,string"` // 전체 자기 매도 거래대금 비율
	ArbtSmtmShunVolRate    float64 `json:"arbt_smtm_shun_vol_rate,string"`     // 차익 합계 매수 거래량 비율 (shun 타이포)
	NabtSmtmShunTrPbmnRate float64 `json:"nabt_smtm_shun_tr_pbmn_rate,string"` // 비차익 합계 매수 거래대금 비율 (shun 타이포)
	ArbtEntmNtbyQtyRate    float64 `json:"arbt_entm_ntby_qty_rate,string"`     // 차익 위탁 순매수 수량 비율
	NabtEntmNtbyTrPbmnRate float64 `json:"nabt_entm_ntby_tr_pbmn_rate,string"` // 비차익 위탁 순매수 거래대금 비율
	ArbtEntmSelnVolRate    float64 `json:"arbt_entm_seln_vol_rate,string"`     // 차익 위탁 매도 거래량 비율
	NabtEntmSelnTrPbmnRate float64 `json:"nabt_entm_seln_tr_pbmn_rate,string"` // 비차익 위탁 매도 거래대금 비율
	NabtOnslSelnTrPbmn     int64   `json:"nabt_onsl_seln_tr_pbmn,string"`      // 비차익 자기 매도 거래대금
	WholSmtnSelnVol        int64   `json:"whol_smtn_seln_vol,string"`          // 전체 합계 매도 거래량
	ArbtSmtnShnuTrPbmn     int64   `json:"arbt_smtn_shnu_tr_pbmn,string"`      // 차익 합계 매수 거래대금
	WholEntmShnuVol        int64   `json:"whol_entm_shnu_vol,string"`          // 전체 위탁 매수 거래량
	ArbtEntmNtbyTrPbmn     int64   `json:"arbt_entm_ntby_tr_pbmn,string"`      // 차익 위탁 순매수 거래대금
	NabtOnslNtbyQty        int64   `json:"nabt_onsl_ntby_qty,string"`          // 비차익 자기 순매수 수량
	ArbtEntmSelnTrPbmn     int64   `json:"arbt_entm_seln_tr_pbmn,string"`      // 차익 위탁 매도 거래대금
	NabtOnslSelnTrPbmnRate float64 `json:"nabt_onsl_seln_tr_pbmn_rate,string"` // 비차익 자기 매도 거래대금 비율
	WholSelnVolRate        float64 `json:"whol_seln_vol_rate,string"`          // 전체 매도 거래량 비율
	ArbtSmtmShunTrPbmnRate float64 `json:"arbt_smtm_shun_tr_pbmn_rate,string"` // 차익 합계 매수 거래대금 비율 (shun 타이포)
	WholEntmShnuVolRate    float64 `json:"whol_entm_shnu_vol_rate,string"`     // 전체 위탁 매수 거래량 비율
	ArbtEntmNtbyTrPbmnRate float64 `json:"arbt_entm_ntby_tr_pbmn_rate,string"` // 차익 위탁 순매수 거래대금 비율
	NabtOnslNtbyQtyRate    float64 `json:"nabt_onsl_ntby_qty_rate,string"`     // 비차익 자기 순매수 수량 비율
	ArbtEntmSelnTrPbmnRate float64 `json:"arbt_entm_seln_tr_pbmn_rate,string"` // 차익 위탁 매도 거래대금 비율
	NabtSmtnSelnVol        int64   `json:"nabt_smtn_seln_vol,string"`          // 비차익 합계 매도 거래량
	WholSmtnSelnTrPbmn     int64   `json:"whol_smtn_seln_tr_pbmn,string"`      // 전체 합계 매도 거래대금
	NabtEntmShnuVol        int64   `json:"nabt_entm_shnu_vol,string"`          // 비차익 위탁 매수 거래량
	WholEntmShnuTrPbmn     int64   `json:"whol_entm_shnu_tr_pbmn,string"`      // 전체 위탁 매수 거래대금
	ArbtOnslNtbyQty        int64   `json:"arbt_onsl_ntby_qty,string"`          // 차익 자기 순매수 수량
	NabtOnslNtbyTrPbmn     int64   `json:"nabt_onsl_ntby_tr_pbmn,string"`      // 비차익 자기 순매수 거래대금
	ArbtOnslSelnTrPbmn     int64   `json:"arbt_onsl_seln_tr_pbmn,string"`      // 차익 자기 매도 거래대금
	NabtSmtmSelnVolRate    float64 `json:"nabt_smtm_seln_vol_rate,string"`     // 비차익 합계 매도 거래량 비율
	WholSelnTrPbmnRate     float64 `json:"whol_seln_tr_pbmn_rate,string"`      // 전체 매도 거래대금 비율
	NabtEntmShnuVolRate    float64 `json:"nabt_entm_shnu_vol_rate,string"`     // 비차익 위탁 매수 거래량 비율
	WholEntmShnuTrPbmnRate float64 `json:"whol_entm_shnu_tr_pbmn_rate,string"` // 전체 위탁 매수 거래대금 비율
	ArbtOnslNtbyQtyRate    float64 `json:"arbt_onsl_ntby_qty_rate,string"`     // 차익 자기 순매수 수량 비율
	NabtOnslNtbyTrPbmnRate float64 `json:"nabt_onsl_ntby_tr_pbmn_rate,string"` // 비차익 자기 순매수 거래대금 비율
	ArbtOnslSelnTrPbmnRate float64 `json:"arbt_onsl_seln_tr_pbmn_rate,string"` // 차익 자기 매도 거래대금 비율
	NabtSmtnSelnTrPbmn     int64   `json:"nabt_smtn_seln_tr_pbmn,string"`      // 비차익 합계 매도 거래대금
	ArbtEntmShnuVol        int64   `json:"arbt_entm_shnu_vol,string"`          // 차익 위탁 매수 거래량
	NabtEntmShnuTrPbmn     int64   `json:"nabt_entm_shnu_tr_pbmn,string"`      // 비차익 위탁 매수 거래대금
	WholOnslShnuVol        int64   `json:"whol_onsl_shnu_vol,string"`          // 전체 자기 매수 거래량
	ArbtOnslNtbyTrPbmn     int64   `json:"arbt_onsl_ntby_tr_pbmn,string"`      // 차익 자기 순매수 거래대금
	NabtSmtnNtbyQty        int64   `json:"nabt_smtn_ntby_qty,string"`          // 비차익 합계 순매수 수량
	ArbtOnslSelnVol        int64   `json:"arbt_onsl_seln_vol,string"`          // 차익 자기 매도 거래량
	NabtSmtmSelnTrPbmnRate float64 `json:"nabt_smtm_seln_tr_pbmn_rate,string"` // 비차익 합계 매도 거래대금 비율
	ArbtEntmShnuVolRate    float64 `json:"arbt_entm_shnu_vol_rate,string"`     // 차익 위탁 매수 거래량 비율
	NabtEntmShnuTrPbmnRate float64 `json:"nabt_entm_shnu_tr_pbmn_rate,string"` // 비차익 위탁 매수 거래대금 비율
	WholOnslShnuTrPbmn     int64   `json:"whol_onsl_shnu_tr_pbmn,string"`      // 전체 자기 매수 거래대금
	ArbtOnslNtbyTrPbmnRate float64 `json:"arbt_onsl_ntby_tr_pbmn_rate,string"` // 차익 자기 순매수 거래대금 비율
	NabtSmtmNtbyQtyRate    float64 `json:"nabt_smtm_ntby_qty_rate,string"`     // 비차익 합계 순매수 수량 비율
	ArbtOnslSelnVolRate    float64 `json:"arbt_onsl_seln_vol_rate,string"`     // 차익 자기 매도 거래량 비율
	WholEntmSelnVol        int64   `json:"whol_entm_seln_vol,string"`          // 전체 위탁 매도 거래량
	ArbtEntmShnuTrPbmn     int64   `json:"arbt_entm_shnu_tr_pbmn,string"`      // 차익 위탁 매수 거래대금
	NabtOnslShnuVol        int64   `json:"nabt_onsl_shnu_vol,string"`          // 비차익 자기 매수 거래량
	WholOnslShnuTrPbmnRate float64 `json:"whol_onsl_shnu_tr_pbmn_rate,string"` // 전체 자기 매수 거래대금 비율
	ArbtSmtnNtbyQty        int64   `json:"arbt_smtn_ntby_qty,string"`          // 차익 합계 순매수 수량
	NabtSmtnNtbyTrPbmn     int64   `json:"nabt_smtn_ntby_tr_pbmn,string"`      // 비차익 합계 순매수 거래대금
	ArbtSmtnSelnVol        int64   `json:"arbt_smtn_seln_vol,string"`          // 차익 합계 매도 거래량
	WholEntmSelnTrPbmn     int64   `json:"whol_entm_seln_tr_pbmn,string"`      // 전체 위탁 매도 거래대금
	ArbtEntmShnuTrPbmnRate float64 `json:"arbt_entm_shnu_tr_pbmn_rate,string"` // 차익 위탁 매수 거래대금 비율
	NabtOnslShnuVolRate    float64 `json:"nabt_onsl_shnu_vol_rate,string"`     // 비차익 자기 매수 거래량 비율
	WholOnslShnuVolRate    float64 `json:"whol_onsl_shnu_vol_rate,string"`     // 전체 자기 매수 거래량 비율
	ArbtSmtmNtbyQtyRate    float64 `json:"arbt_smtm_ntby_qty_rate,string"`     // 차익 합계 순매수 수량 비율
	NabtSmtmNtbyTrPbmnRate float64 `json:"nabt_smtm_ntby_tr_pbmn_rate,string"` // 비차익 합계 순매수 거래대금 비율
	ArbtSmtmSelnVolRate    float64 `json:"arbt_smtm_seln_vol_rate,string"`     // 차익 합계 매도 거래량 비율
	WholEntmSelnVolRate    float64 `json:"whol_entm_seln_vol_rate,string"`     // 전체 위탁 매도 거래량 비율
	ArbtOnslShnuVol        int64   `json:"arbt_onsl_shnu_vol,string"`          // 차익 자기 매수 거래량
	NabtOnslShnuTrPbmn     int64   `json:"nabt_onsl_shnu_tr_pbmn,string"`      // 비차익 자기 매수 거래대금
	WholSmtnShnuVol        int64   `json:"whol_smtn_shnu_vol,string"`          // 전체 합계 매수 거래량
	ArbtSmtnNtbyTrPbmn     int64   `json:"arbt_smtn_ntby_tr_pbmn,string"`      // 차익 합계 순매수 거래대금
	WholEntmNtbyQty        int64   `json:"whol_entm_ntby_qty,string"`          // 전체 위탁 순매수 수량
	ArbtSmtnSelnTrPbmn     int64   `json:"arbt_smtn_seln_tr_pbmn,string"`      // 차익 합계 매도 거래대금
	WholEntmSelnTrPbmnRate float64 `json:"whol_entm_seln_tr_pbmn_rate,string"` // 전체 위탁 매도 거래대금 비율
	ArbtOnslShnuVolRate    float64 `json:"arbt_onsl_shnu_vol_rate,string"`     // 차익 자기 매수 거래량 비율
	NabtOnslShnuTrPbmnRate float64 `json:"nabt_onsl_shnu_tr_pbmn_rate,string"` // 비차익 자기 매수 거래대금 비율
	WholShunVolRate        float64 `json:"whol_shun_vol_rate,string"`          // 전체 매수 거래량 비율 (shun 타이포)
	ArbtSmtmNtbyTrPbmnRate float64 `json:"arbt_smtm_ntby_tr_pbmn_rate,string"` // 차익 합계 순매수 거래대금 비율
	WholEntmNtbyQtyRate    float64 `json:"whol_entm_ntby_qty_rate,string"`     // 전체 위탁 순매수 수량 비율
	ArbtSmtmSelnTrPbmnRate float64 `json:"arbt_smtm_seln_tr_pbmn_rate,string"` // 차익 합계 매도 거래대금 비율
	WholOnslSelnVol        int64   `json:"whol_onsl_seln_vol,string"`          // 전체 자기 매도 거래량
	ArbtOnslShnuTrPbmn     int64   `json:"arbt_onsl_shnu_tr_pbmn,string"`      // 차익 자기 매수 거래대금
	NabtSmtnShnuVol        int64   `json:"nabt_smtn_shnu_vol,string"`          // 비차익 합계 매수 거래량
	WholSmtnShnuTrPbmn     int64   `json:"whol_smtn_shnu_tr_pbmn,string"`      // 전체 합계 매수 거래대금
	NabtEntmNtbyQty        int64   `json:"nabt_entm_ntby_qty,string"`          // 비차익 위탁 순매수 수량
	WholEntmNtbyTrPbmn     int64   `json:"whol_entm_ntby_tr_pbmn,string"`      // 전체 위탁 순매수 거래대금
	NabtEntmSelnVol        int64   `json:"nabt_entm_seln_vol,string"`          // 비차익 위탁 매도 거래량
	WholOnslSelnVolRate    float64 `json:"whol_onsl_seln_vol_rate,string"`     // 전체 자기 매도 거래량 비율
	ArbtOnslShnuTrPbmnRate float64 `json:"arbt_onsl_shnu_tr_pbmn_rate,string"` // 차익 자기 매수 거래대금 비율
}

// CompProgramTradeDaily 는 FHPPG04600001 전체 응답이다.
type CompProgramTradeDaily struct {
	RTCd   string                      `json:"rt_cd"`
	MsgCd  string                      `json:"msg_cd"`
	Msg1   string                      `json:"msg1"`
	Output []CompProgramTradeDailyItem `json:"output"`
}

// InquireCompProgramTradeDailyParams 는 EP6 쿼리 파라미터다.
type InquireCompProgramTradeDailyParams struct {
	MarketCode  string // FID_COND_MRKT_DIV_CODE
	MrktClsCode string // FID_MRKT_CLS_CODE  K:코스피  Q:코스닥
	StartDate   string // FID_INPUT_DATE_1  blank or YYYYMMDD  8개월 한도
	EndDate     string // FID_INPUT_DATE_2  blank or YYYYMMDD
}

// InquireCompProgramTradeDaily 는 프로그램매매 종합현황(일별)을 조회한다 (FHPPG04600001).
func (c *Client) InquireCompProgramTradeDaily(ctx context.Context, p InquireCompProgramTradeDailyParams) (*CompProgramTradeDaily, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/comp-program-trade-daily",
		TrID:   "FHPPG04600001",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": p.MarketCode,
			"FID_MRKT_CLS_CODE":      p.MrktClsCode,
			"FID_INPUT_DATE_1":       p.StartDate,
			"FID_INPUT_DATE_2":       p.EndDate,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res CompProgramTradeDaily
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse CompProgramTradeDaily: %w", err)
	}
	return &res, nil
}
