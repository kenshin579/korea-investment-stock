package futures

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// ─── EP3: InquireTimeFuopchartprice ──────────────────────────────────────────

// InquireTimeFuopchartpriceOutput1 는 선물옵션 분봉조회 현재 상태 요약 (31 필드).
type InquireTimeFuopchartpriceOutput1 struct {
	FutsPrdyVrss         decimal.Decimal `json:"futs_prdy_vrss"`            // 선물 전일 대비
	PrdyVrssSign         string          `json:"prdy_vrss_sign"`            // 전일 대비 부호
	FutsPrdyCtrt         float64         `json:"futs_prdy_ctrt,string"`     // 선물 전일 대비율
	FutsPrdyClpr         decimal.Decimal `json:"futs_prdy_clpr"`            // 선물 전일 종가
	PrdyNmix             decimal.Decimal `json:"prdy_nmix"`                 // 전일 지수
	AcmlVol              int64           `json:"acml_vol,string"`           // 누적 거래량
	AcmlTrPbmn           int64           `json:"acml_tr_pbmn,string"`       // 누적 거래 대금
	HtsKorIsnm           string          `json:"hts_kor_isnm"`              // HTS 한글 종목명
	FutsPrpr             decimal.Decimal `json:"futs_prpr"`                 // 선물 현재가
	FutsShrnIscd         string          `json:"futs_shrn_iscd"`            // 선물 단축 종목코드
	PrdyVol              int64           `json:"prdy_vol,string"`           // 전일 거래량
	FutsMxpr             decimal.Decimal `json:"futs_mxpr"`                 // 선물 상한가
	FutsLlam             decimal.Decimal `json:"futs_llam"`                 // 선물 하한가
	FutsOprc             decimal.Decimal `json:"futs_oprc"`                 // 선물 시가
	FutsHgpr             decimal.Decimal `json:"futs_hgpr"`                 // 선물 최고가
	FutsLwpr             decimal.Decimal `json:"futs_lwpr"`                 // 선물 최저가
	FutsPrdyOprc         decimal.Decimal `json:"futs_prdy_oprc"`            // 선물 전일 시가
	FutsPrdyHgpr         decimal.Decimal `json:"futs_prdy_hgpr"`            // 선물 전일 최고가
	FutsPrdyLwpr         decimal.Decimal `json:"futs_prdy_lwpr"`            // 선물 전일 최저가
	FutsAskp             decimal.Decimal `json:"futs_askp"`                 // 선물 매도호가
	FutsBidp             decimal.Decimal `json:"futs_bidp"`                 // 선물 매수호가
	Basis                decimal.Decimal `json:"basis"`                     // 베이시스
	Kospi200Nmix         decimal.Decimal `json:"kospi200_nmix"`             // KOSPI200 지수
	Kospi200PrdyVrss     decimal.Decimal `json:"kospi200_prdy_vrss"`        // KOSPI200 전일 대비
	Kospi200PrdyCtrt     float64         `json:"kospi200_prdy_ctrt,string"` // KOSPI200 전일 대비율
	Kospi200PrdyVrssSign string          `json:"kospi200_prdy_vrss_sign"`   // KOSPI200 전일 대비 부호
	HtsOtstStplQty       int64           `json:"hts_otst_stpl_qty,string"`  // HTS 미결제 약정 수량
	OtstStplQtyIcdc      int64           `json:"otst_stpl_qty_icdc,string"` // 미결제 약정 수량 증감
	TdayRltv             float64         `json:"tday_rltv,string"`          // 당일 체결강도
	HtsThpr              decimal.Decimal `json:"hts_thpr"`                  // HTS 이론가
	Dprt                 float64         `json:"dprt,string"`               // 괴리율
}

// InquireTimeFuopchartpriceOutput2Item 는 선물옵션 분봉 캔들 (8 필드).
type InquireTimeFuopchartpriceOutput2Item struct {
	StckBsopDate string          `json:"stck_bsop_date"`      // 주식 영업 일자
	StckCntgHour string          `json:"stck_cntg_hour"`      // 주식 체결 시간
	FutsPrpr     decimal.Decimal `json:"futs_prpr"`           // 선물 현재가
	FutsOprc     decimal.Decimal `json:"futs_oprc"`           // 선물 시가
	FutsHgpr     decimal.Decimal `json:"futs_hgpr"`           // 선물 최고가
	FutsLwpr     decimal.Decimal `json:"futs_lwpr"`           // 선물 최저가
	CntgVol      int64           `json:"cntg_vol,string"`     // 체결 거래량
	AcmlTrPbmn   int64           `json:"acml_tr_pbmn,string"` // 누적 거래 대금
}

// InquireTimeFuopchartpriceData 는 선물옵션 분봉조회 응답 (output1 + output2[]).
type InquireTimeFuopchartpriceData struct {
	Output1 InquireTimeFuopchartpriceOutput1       `json:"output1"`
	Output2 []InquireTimeFuopchartpriceOutput2Item `json:"output2"`
}

type inquireTimeFuopchartpriceResponse struct {
	RtCd    string                                 `json:"rt_cd"`
	MsgCd   string                                 `json:"msg_cd"`
	Msg1    string                                 `json:"msg1"`
	Output1 InquireTimeFuopchartpriceOutput1       `json:"output1"`
	Output2 []InquireTimeFuopchartpriceOutput2Item `json:"output2"`
}

// InquireTimeFuopchartpriceParams 는 선물옵션 분봉조회 파라미터.
type InquireTimeFuopchartpriceParams struct {
	MarketCode     string // 조건 시장 분류 코드 (F/O/JF/JO/CF/CM/EU)
	Code           string // 입력 종목코드
	HourClsCode    string // 시간 구분 코드 (30:30초, 60:1분, 3600:1시간)
	PwDataIncuYn   string // 과거 데이터 포함 여부 (Y:과거, N:당일)
	FakeTickIncuYn string // 허봉 포함 여부 (N 입력)
	InputDate      string // 입력 날짜 (YYYYMMDD)
	InputHour      string // 입력 시간 (HHMMSS)
}

// InquireTimeFuopchartprice 는 선물옵션 분봉조회 (FHKIF03020200).
//
// 모의: 미지원 (실전 only)
// 한 번 호출에 최대 102건. 연속 조회: FID_INPUT_DATE_1, FID_INPUT_HOUR_1 이용.
//
// KIS API: GET /uapi/domestic-futureoption/v1/quotations/inquire-time-fuopchartprice
func (c *Client) InquireTimeFuopchartprice(ctx context.Context, params InquireTimeFuopchartpriceParams) (*InquireTimeFuopchartpriceData, error) {
	fakeTickYn := params.FakeTickIncuYn
	if fakeTickYn == "" {
		fakeTickYn = "N"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-futureoption/v1/quotations/inquire-time-fuopchartprice",
		TrID:     "FHKIF03020200",
		CustType: "P",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": params.MarketCode,
			"FID_INPUT_ISCD":         params.Code,
			"FID_HOUR_CLS_CODE":      params.HourClsCode,
			"FID_PW_DATA_INCU_YN":    params.PwDataIncuYn,
			"FID_FAKE_TICK_INCU_YN":  fakeTickYn,
			"FID_INPUT_DATE_1":       params.InputDate,
			"FID_INPUT_HOUR_1":       params.InputHour,
		},
	})
	if err != nil {
		return nil, err
	}
	var res inquireTimeFuopchartpriceResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InquireTimeFuopchartprice: %w", err)
	}
	return &InquireTimeFuopchartpriceData{
		Output1: res.Output1,
		Output2: res.Output2,
	}, nil
}

// ─── EP6: InquireDailyFuopchartprice ─────────────────────────────────────────

// InquireDailyFuopchartpriceOutput1 는 선물옵션 기간별 시세 현재 상태 요약 (30 필드).
type InquireDailyFuopchartpriceOutput1 struct {
	FutsPrdyVrss         decimal.Decimal `json:"futs_prdy_vrss"`            // 전일 대비
	PrdyVrssSign         string          `json:"prdy_vrss_sign"`            // 전일 대비 부호
	FutsPrdyCtrt         float64         `json:"futs_prdy_ctrt,string"`     // 선물 전일 대비율
	FutsPrdyClpr         decimal.Decimal `json:"futs_prdy_clpr"`            // 선물 전일 종가
	AcmlVol              int64           `json:"acml_vol,string"`           // 누적 거래량
	AcmlTrPbmn           int64           `json:"acml_tr_pbmn,string"`       // 누적 거래 대금
	HtsKorIsnm           string          `json:"hts_kor_isnm"`              // HTS 한글 종목명
	FutsPrpr             decimal.Decimal `json:"futs_prpr"`                 // 현재가
	FutsShrnIscd         string          `json:"futs_shrn_iscd"`            // 단축 종목코드
	PrdyVol              int64           `json:"prdy_vol,string"`           // 전일 거래량
	FutsMxpr             decimal.Decimal `json:"futs_mxpr"`                 // 상한가
	FutsLlam             decimal.Decimal `json:"futs_llam"`                 // 하한가
	FutsOprc             decimal.Decimal `json:"futs_oprc"`                 // 시가
	FutsHgpr             decimal.Decimal `json:"futs_hgpr"`                 // 최고가
	FutsLwpr             decimal.Decimal `json:"futs_lwpr"`                 // 최저가
	FutsPrdyOprc         decimal.Decimal `json:"futs_prdy_oprc"`            // 전일 시가
	FutsPrdyHgpr         decimal.Decimal `json:"futs_prdy_hgpr"`            // 전일 최고가
	FutsPrdyLwpr         decimal.Decimal `json:"futs_prdy_lwpr"`            // 전일 최저가
	FutsAskp             decimal.Decimal `json:"futs_askp"`                 // 매도호가
	FutsBidp             decimal.Decimal `json:"futs_bidp"`                 // 매수호가
	Basis                decimal.Decimal `json:"basis"`                     // 베이시스
	Kospi200Nmix         decimal.Decimal `json:"kospi200_nmix"`             // KOSPI200 지수
	Kospi200PrdyVrss     decimal.Decimal `json:"kospi200_prdy_vrss"`        // KOSPI200 전일 대비
	Kospi200PrdyCtrt     float64         `json:"kospi200_prdy_ctrt,string"` // KOSPI200 전일 대비율
	Kospi200PrdyVrssSign string          `json:"kospi200_prdy_vrss_sign"`   // 전일 대비 부호
	HtsOtstStplQty       int64           `json:"hts_otst_stpl_qty,string"`  // HTS 미결제 약정 수량
	OtstStplQtyIcdc      int64           `json:"otst_stpl_qty_icdc,string"` // 미결제 약정 수량 증감
	TdayRltv             float64         `json:"tday_rltv,string"`          // 당일 체결강도
	HtsThpr              decimal.Decimal `json:"hts_thpr"`                  // HTS 이론가
	Dprt                 float64         `json:"dprt,string"`               // 괴리율
}

// InquireDailyFuopchartpriceOutput2Item 는 선물옵션 기간별 시세 캔들 (8 필드).
type InquireDailyFuopchartpriceOutput2Item struct {
	StckBsopDate string          `json:"stck_bsop_date"`      // 영업 일자
	FutsPrpr     decimal.Decimal `json:"futs_prpr"`           // 현재가
	FutsOprc     decimal.Decimal `json:"futs_oprc"`           // 시가
	FutsHgpr     decimal.Decimal `json:"futs_hgpr"`           // 최고가
	FutsLwpr     decimal.Decimal `json:"futs_lwpr"`           // 최저가
	AcmlVol      int64           `json:"acml_vol,string"`     // 누적 거래량
	AcmlTrPbmn   int64           `json:"acml_tr_pbmn,string"` // 누적 거래 대금
	ModYn        string          `json:"mod_yn"`              // 변경 여부
}

// InquireDailyFuopchartpriceData 는 선물옵션 기간별 시세 응답 (output1 + output2[]).
type InquireDailyFuopchartpriceData struct {
	Output1 InquireDailyFuopchartpriceOutput1       `json:"output1"`
	Output2 []InquireDailyFuopchartpriceOutput2Item `json:"output2"`
}

type inquireDailyFuopchartpriceResponse struct {
	RtCd    string                                  `json:"rt_cd"`
	MsgCd   string                                  `json:"msg_cd"`
	Msg1    string                                  `json:"msg1"`
	Output1 InquireDailyFuopchartpriceOutput1       `json:"output1"`
	Output2 []InquireDailyFuopchartpriceOutput2Item `json:"output2"`
}

// InquireDailyFuopchartpriceParams 는 선물옵션 기간별 시세 파라미터.
type InquireDailyFuopchartpriceParams struct {
	MarketCode string // 조건 시장 분류 코드 (F/O/JF/JO/CF/CM/EU)
	Code       string // 종목코드
	FromDate   string // 조회 시작일자 (YYYYMMDD)
	ToDate     string // 조회 종료일자 (YYYYMMDD)
	Period     string // 기간분류코드 (D:일봉, W:주봉, M:월봉, Y:년봉), 빈 값이면 "D"
}

// InquireDailyFuopchartprice 는 선물옵션 기간별 시세(일/주/월/년) 조회 (FHKIF03020100).
//
// 모의: 지원
// 최대 100건 / 연속 조회 지원 (tr_cont).
//
// KIS API: GET /uapi/domestic-futureoption/v1/quotations/inquire-daily-fuopchartprice
func (c *Client) InquireDailyFuopchartprice(ctx context.Context, params InquireDailyFuopchartpriceParams) (*InquireDailyFuopchartpriceData, error) {
	period := params.Period
	if period == "" {
		period = "D"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-futureoption/v1/quotations/inquire-daily-fuopchartprice",
		TrID:     "FHKIF03020100",
		CustType: "P",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": params.MarketCode,
			"FID_INPUT_ISCD":         params.Code,
			"FID_INPUT_DATE_1":       params.FromDate,
			"FID_INPUT_DATE_2":       params.ToDate,
			"FID_PERIOD_DIV_CODE":    period,
		},
	})
	if err != nil {
		return nil, err
	}
	var res inquireDailyFuopchartpriceResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InquireDailyFuopchartprice: %w", err)
	}
	return &InquireDailyFuopchartpriceData{
		Output1: res.Output1,
		Output2: res.Output2,
	}, nil
}
