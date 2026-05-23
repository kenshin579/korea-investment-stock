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

// ─── EP8: DisplayBoardTop ─────────────────────────────────────────────────────

// DisplayBoardTopOutput1 는 국내선물 기초자산 + 선물 현재 시세 (10 필드).
type DisplayBoardTopOutput1 struct {
	UnasPrpr         decimal.Decimal `json:"unas_prpr"`            // 기초자산 현재가
	UnasPrdyVrss     decimal.Decimal `json:"unas_prdy_vrss"`       // 기초자산 전일 대비
	UnasPrdyVrssSign string          `json:"unas_prdy_vrss_sign"`  // 기초자산 전일 대비 부호
	UnasPrdyCtrt     kistypes.Float  `json:"unas_prdy_ctrt"`       // 기초자산 전일 대비율
	UnasAcmlVol      int64           `json:"unas_acml_vol,string"` // 기초자산 누적 거래량
	HtsKorIsnm       string          `json:"hts_kor_isnm"`         // HTS 한글 종목명
	FutsPrpr         decimal.Decimal `json:"futs_prpr"`            // 선물 현재가
	FutsPrdyVrss     decimal.Decimal `json:"futs_prdy_vrss"`       // 선물 전일 대비
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`       // 전일 대비 부호
	FutsPrdyCtrt     kistypes.Float  `json:"futs_prdy_ctrt"`       // 선물 전일 대비율
}

// DisplayBoardTopOutput2Item 는 월물별 잔존일수 항목 (1 필드).
type DisplayBoardTopOutput2Item struct {
	HtsRmnnDynu string `json:"hts_rmnn_dynu"` // HTS 잔존 일수
}

// DisplayBoardTopData 는 국내선물 기초자산 시세 응답 (output1 + output2[]).
type DisplayBoardTopData struct {
	Output1 DisplayBoardTopOutput1       `json:"output1"`
	Output2 []DisplayBoardTopOutput2Item `json:"output2"`
}

type displayBoardTopResponse struct {
	RtCd    string                       `json:"rt_cd"`
	MsgCd   string                       `json:"msg_cd"`
	Msg1    string                       `json:"msg1"`
	Output1 DisplayBoardTopOutput1       `json:"output1"`
	Output2 []DisplayBoardTopOutput2Item `json:"output2"`
}

// DisplayBoardTopParams 는 국내선물 기초자산 시세 파라미터.
type DisplayBoardTopParams struct {
	MarketCode  string // 조건 시장 분류 코드 (F:선물), 기본값 "F"
	Code        string // 입력 종목코드 (선물최근월물, 예: 101V06)
	MarketCode1 string // 조건 시장 분류 코드1 (공백)
	ScrDivCode  string // 조건 화면 분류 코드 (공백)
	MtrtCnt     string // 만기 수 (공백)
	MrktClsCode string // 조건 시장 구분 코드 (공백)
}

// DisplayBoardTop 는 국내선물 기초자산 시세 조회 (FHPIF05030000).
//
// 모의: 미지원 (실전 only)
// output2 는 월물별 잔존일수 목록 (hts_rmnn_dynu 1 필드만 존재).
//
// KIS API: GET /uapi/domestic-futureoption/v1/quotations/display-board-top
func (c *Client) DisplayBoardTop(ctx context.Context, params DisplayBoardTopParams) (*DisplayBoardTopData, error) {
	marketCode := params.MarketCode
	if marketCode == "" {
		marketCode = "F"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-futureoption/v1/quotations/display-board-top",
		TrID:     "FHPIF05030000",
		CustType: "P",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE":  marketCode,
			"FID_INPUT_ISCD":          params.Code,
			"FID_COND_MRKT_DIV_CODE1": params.MarketCode1,
			"FID_COND_SCR_DIV_CODE":   params.ScrDivCode,
			"FID_MTRT_CNT":            params.MtrtCnt,
			"FID_COND_MRKT_CLS_CODE":  params.MrktClsCode,
		},
	})
	if err != nil {
		return nil, err
	}
	var res displayBoardTopResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse DisplayBoardTop: %w", err)
	}
	return &DisplayBoardTopData{
		Output1: res.Output1,
		Output2: res.Output2,
	}, nil
}

// ─── EP9: DisplayBoardFutures ─────────────────────────────────────────────────

// DisplayBoardFuturesItem 는 선물 월물별 시세 항목 (20 필드).
type DisplayBoardFuturesItem struct {
	FutsShrnIscd     string          `json:"futs_shrn_iscd"`           // 선물 단축 종목코드
	HtsKorIsnm       string          `json:"hts_kor_isnm"`             // HTS 한글 종목명
	FutsPrpr         decimal.Decimal `json:"futs_prpr"`                // 선물 현재가
	FutsPrdyVrss     decimal.Decimal `json:"futs_prdy_vrss"`           // 선물 전일 대비
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`           // 전일 대비 부호
	FutsPrdyCtrt     kistypes.Float  `json:"futs_prdy_ctrt"`           // 선물 전일 대비율
	HtsThpr          decimal.Decimal `json:"hts_thpr"`                 // HTS 이론가
	AcmlVol          int64           `json:"acml_vol,string"`          // 누적 거래량
	FutsAskp         decimal.Decimal `json:"futs_askp"`                // 선물 매도호가
	FutsBidp         decimal.Decimal `json:"futs_bidp"`                // 선물 매수호가
	HtsOtstStplQty   int64           `json:"hts_otst_stpl_qty,string"` // HTS 미결제 약정 수량
	FutsHgpr         decimal.Decimal `json:"futs_hgpr"`                // 선물 최고가
	FutsLwpr         decimal.Decimal `json:"futs_lwpr"`                // 선물 최저가
	HtsRmnnDynu      string          `json:"hts_rmnn_dynu"`            // HTS 잔존 일수
	TotalAskpRsqn    int64           `json:"total_askp_rsqn,string"`   // 총 매도호가 잔량
	TotalBidpRsqn    int64           `json:"total_bidp_rsqn,string"`   // 총 매수호가 잔량
	FutsAntcCnpr     decimal.Decimal `json:"futs_antc_cnpr"`           // 선물 예상체결가
	FutsAntcCntgVrss decimal.Decimal `json:"futs_antc_cntg_vrss"`      // 선물 예상체결대비
	AntcCntgVrssSign string          `json:"antc_cntg_vrss_sign"`      // 예상 체결 대비 부호
	AntcCntgPrdyCtrt kistypes.Float  `json:"antc_cntg_prdy_ctrt"`      // 예상 체결 전일 대비율
}

// DisplayBoardFuturesData 는 국내옵션전광판 선물 응답 (output1[]).
type DisplayBoardFuturesData struct {
	Output1 []DisplayBoardFuturesItem `json:"output1"`
}

type displayBoardFuturesResponse struct {
	RtCd    string                    `json:"rt_cd"`
	MsgCd   string                    `json:"msg_cd"`
	Msg1    string                    `json:"msg1"`
	Output1 []DisplayBoardFuturesItem `json:"output1"`
}

// DisplayBoardFuturesParams 는 국내옵션전광판 선물 파라미터.
type DisplayBoardFuturesParams struct {
	MarketCode  string // 조건 시장 분류 코드 (F:선물), 기본값 "F"
	ScrDivCode  string // 조건 화면 분류 코드, 기본값 "20503"
	MrktClsCode string // 조건 시장 구분 코드 (공백:KOSPI200, MKI:미니, WKM:위클리(월), WKI:위클리(목), KQI:KOSDAQ150)
}

// DisplayBoardFutures 는 국내옵션전광판 선물 조회 (FHPIF05030200).
//
// 모의: 미지원 (실전 only)
// 복수 월물의 선물 목록 반환.
//
// KIS API: GET /uapi/domestic-futureoption/v1/quotations/display-board-futures
func (c *Client) DisplayBoardFutures(ctx context.Context, params DisplayBoardFuturesParams) (*DisplayBoardFuturesData, error) {
	marketCode := params.MarketCode
	if marketCode == "" {
		marketCode = "F"
	}
	scrDivCode := params.ScrDivCode
	if scrDivCode == "" {
		scrDivCode = "20503"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-futureoption/v1/quotations/display-board-futures",
		TrID:     "FHPIF05030200",
		CustType: "P",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": marketCode,
			"FID_COND_SCR_DIV_CODE":  scrDivCode,
			"FID_COND_MRKT_CLS_CODE": params.MrktClsCode,
		},
	})
	if err != nil {
		return nil, err
	}
	var res displayBoardFuturesResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse DisplayBoardFutures: %w", err)
	}
	return &DisplayBoardFuturesData{
		Output1: res.Output1,
	}, nil
}

// ─── EP10: DisplayBoardOptionList ────────────────────────────────────────────

// DisplayBoardOptionListItem 는 옵션 월물 항목 (2 필드).
type DisplayBoardOptionListItem struct {
	MtrtYymmCode string `json:"mtrt_yymm_code"` // 만기 년월 코드
	MtrtYymm     string `json:"mtrt_yymm"`      // 만기 년월
}

// DisplayBoardOptionListData 는 국내옵션전광판 옵션월물리스트 응답 (output1[]).
type DisplayBoardOptionListData struct {
	Output1 []DisplayBoardOptionListItem `json:"output1"`
}

type displayBoardOptionListResponse struct {
	RtCd    string                       `json:"rt_cd"`
	MsgCd   string                       `json:"msg_cd"`
	Msg1    string                       `json:"msg1"`
	Output1 []DisplayBoardOptionListItem `json:"output1"`
}

// DisplayBoardOptionListParams 는 국내옵션전광판 옵션월물리스트 파라미터.
type DisplayBoardOptionListParams struct {
	ScrDivCode  string // 조건 화면 분류 코드, 기본값 "509"
	MarketCode  string // 조건 시장 분류 코드 (공백)
	MrktClsCode string // 조건 시장 구분 코드 (공백)
}

// DisplayBoardOptionList 는 국내옵션전광판 옵션월물리스트 조회 (FHPIO056104C0).
//
// 모의: 미지원 (실전 only)
// TR_ID 가 O prefix — 옵션(Option) 전용 EP.
// 옵션 월물 목록만 반환하는 단순 EP.
//
// KIS API: GET /uapi/domestic-futureoption/v1/quotations/display-board-option-list
func (c *Client) DisplayBoardOptionList(ctx context.Context, params DisplayBoardOptionListParams) (*DisplayBoardOptionListData, error) {
	scrDivCode := params.ScrDivCode
	if scrDivCode == "" {
		scrDivCode = "509"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-futureoption/v1/quotations/display-board-option-list",
		TrID:     "FHPIO056104C0",
		CustType: "P",
		Query: map[string]string{
			"FID_COND_SCR_DIV_CODE":  scrDivCode,
			"FID_COND_MRKT_DIV_CODE": params.MarketCode,
			"FID_COND_MRKT_CLS_CODE": params.MrktClsCode,
		},
	})
	if err != nil {
		return nil, err
	}
	var res displayBoardOptionListResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse DisplayBoardOptionList: %w", err)
	}
	return &DisplayBoardOptionListData{
		Output1: res.Output1,
	}, nil
}

// ─── EP11: DisplayBoardCallput ────────────────────────────────────────────────

// DisplayBoardCallputItem 는 콜옵션 또는 풋옵션 행사가별 시세 항목 (41 필드).
//
// output1 (콜옵션) / output2 (풋옵션) 가 동일 구조이므로 하나의 struct 재사용.
type DisplayBoardCallputItem struct {
	Acpr             decimal.Decimal `json:"acpr"`                      // 행사가
	UnchPrpr         decimal.Decimal `json:"unch_prpr"`                 // 환산 현재가
	OptnShrnIscd     string          `json:"optn_shrn_iscd"`            // 옵션 단축 종목코드
	OptnPrpr         decimal.Decimal `json:"optn_prpr"`                 // 옵션 현재가
	OptnPrdyVrss     decimal.Decimal `json:"optn_prdy_vrss"`            // 옵션 전일 대비
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`            // 전일 대비 부호
	OptnPrdyCtrt     kistypes.Float  `json:"optn_prdy_ctrt"`            // 옵션 전일 대비율
	OptnBidp         decimal.Decimal `json:"optn_bidp"`                 // 옵션 매수호가
	OptnAskp         decimal.Decimal `json:"optn_askp"`                 // 옵션 매도호가
	TmvlVal          decimal.Decimal `json:"tmvl_val"`                  // 시간가치 값
	NmixSdpr         decimal.Decimal `json:"nmix_sdpr"`                 // 지수 기준가
	AcmlVol          int64           `json:"acml_vol,string"`           // 누적 거래량
	SelnRsqn         int64           `json:"seln_rsqn,string"`          // 매도 잔량
	ShnuRsqn         int64           `json:"shnu_rsqn,string"`          // 매수 잔량
	AcmlTrPbmn       int64           `json:"acml_tr_pbmn,string"`       // 누적 거래 대금
	HtsOtstStplQty   int64           `json:"hts_otst_stpl_qty,string"`  // HTS 미결제 약정 수량
	OtstStplQtyIcdc  int64           `json:"otst_stpl_qty_icdc,string"` // 미결제 약정 수량 증감
	DeltaVal         kistypes.Float  `json:"delta_val"`                 // 델타 값
	Gama             kistypes.Float  `json:"gama"`                      // 감마
	Vega             kistypes.Float  `json:"vega"`                      // 베가
	Theta            kistypes.Float  `json:"theta"`                     // 세타
	Rho              kistypes.Float  `json:"rho"`                       // 로우
	HtsIntsVltl      kistypes.Float  `json:"hts_ints_vltl"`             // HTS 내재 변동성
	InvlVal          decimal.Decimal `json:"invl_val"`                  // 내재가치 값
	Esdg             kistypes.Float  `json:"esdg"`                      // 괴리도
	Dprt             kistypes.Float  `json:"dprt"`                      // 괴리율
	HistVltl         kistypes.Float  `json:"hist_vltl"`                 // 역사적 변동성
	HtsThpr          decimal.Decimal `json:"hts_thpr"`                  // HTS 이론가
	OptnOprc         decimal.Decimal `json:"optn_oprc"`                 // 옵션 시가
	OptnHgpr         decimal.Decimal `json:"optn_hgpr"`                 // 옵션 최고가
	OptnLwpr         decimal.Decimal `json:"optn_lwpr"`                 // 옵션 최저가
	OptnMxpr         decimal.Decimal `json:"optn_mxpr"`                 // 옵션 상한가
	OptnLlam         decimal.Decimal `json:"optn_llam"`                 // 옵션 하한가
	AtmClsName       string          `json:"atm_cls_name"`              // ATM 구분 명
	RgbfVrssIcdc     string          `json:"rgbf_vrss_icdc"`            // 직전 대비 증감
	TotalAskpRsqn    int64           `json:"total_askp_rsqn,string"`    // 총 매도호가 잔량
	TotalBidpRsqn    int64           `json:"total_bidp_rsqn,string"`    // 총 매수호가 잔량
	FutsAntcCnpr     decimal.Decimal `json:"futs_antc_cnpr"`            // 선물 예상체결가
	FutsAntcCntgVrss decimal.Decimal `json:"futs_antc_cntg_vrss"`       // 선물 예상체결대비
	AntcCntgVrssSign string          `json:"antc_cntg_vrss_sign"`       // 예상 체결 대비 부호
	AntcCntgPrdyCtrt kistypes.Float  `json:"antc_cntg_prdy_ctrt"`       // 예상 체결 전일 대비율
}

// DisplayBoardCallputData 는 국내옵션전광판 콜풋 응답 (output1[] 콜 + output2[] 풋).
type DisplayBoardCallputData struct {
	Output1 []DisplayBoardCallputItem `json:"output1"` // 콜옵션 목록 (최대 100건)
	Output2 []DisplayBoardCallputItem `json:"output2"` // 풋옵션 목록 (최대 100건)
}

type displayBoardCallputResponse struct {
	RtCd    string                    `json:"rt_cd"`
	MsgCd   string                    `json:"msg_cd"`
	Msg1    string                    `json:"msg1"`
	Output1 []DisplayBoardCallputItem `json:"output1"`
	Output2 []DisplayBoardCallputItem `json:"output2"`
}

// DisplayBoardCallputParams 는 국내옵션전광판 콜풋 파라미터.
type DisplayBoardCallputParams struct {
	MarketCode   string // 조건 시장 분류 코드 (O:옵션), 기본값 "O"
	ScrDivCode   string // 조건 화면 분류 코드, 기본값 "20503"
	MrktClsCode  string // 시장 구분 코드 (CO:콜옵션), 기본값 "CO"
	MtrtCnt      string // 만기년월(YYYYMM) 또는 만기년월주차(YYMMWW)
	MrktClsCode1 string // 조건 시장 구분 코드 (공백:KOSPI200, MKI, WKM, WKI, KQI)
	MrktClsCode2 string // 시장 구분 코드 (PO:풋옵션), 기본값 "PO"
}

// DisplayBoardCallput 는 국내옵션전광판 콜풋 조회 (FHPIF05030100).
//
// 모의: 미지원 (실전 only)
// output1(콜옵션) + output2(풋옵션) 각각 행사가별 목록 (각 최대 100건).
// 조회 속도가 느린 API — 1초당 최대 1건 권장.
//
// KIS API: GET /uapi/domestic-futureoption/v1/quotations/display-board-callput
func (c *Client) DisplayBoardCallput(ctx context.Context, params DisplayBoardCallputParams) (*DisplayBoardCallputData, error) {
	marketCode := params.MarketCode
	if marketCode == "" {
		marketCode = "O"
	}
	scrDivCode := params.ScrDivCode
	if scrDivCode == "" {
		scrDivCode = "20503"
	}
	mrktClsCode := params.MrktClsCode
	if mrktClsCode == "" {
		mrktClsCode = "CO"
	}
	mrktClsCode2 := params.MrktClsCode2
	if mrktClsCode2 == "" {
		mrktClsCode2 = "PO"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-futureoption/v1/quotations/display-board-callput",
		TrID:     "FHPIF05030100",
		CustType: "P",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": marketCode,
			"FID_COND_SCR_DIV_CODE":  scrDivCode,
			"FID_MRKT_CLS_CODE":      mrktClsCode,
			"FID_MTRT_CNT":           params.MtrtCnt,
			"FID_COND_MRKT_CLS_CODE": params.MrktClsCode1,
			"FID_MRKT_CLS_CODE1":     mrktClsCode2,
		},
	})
	if err != nil {
		return nil, err
	}
	var res displayBoardCallputResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse DisplayBoardCallput: %w", err)
	}
	return &DisplayBoardCallputData{
		Output1: res.Output1,
		Output2: res.Output2,
	}, nil
}
