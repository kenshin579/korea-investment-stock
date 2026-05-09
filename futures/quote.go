package futures

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// ─── EP1: InquirePrice ───────────────────────────────────────────────────────

// InquirePriceOutput1 는 선물옵션 시세 상세 정보 (34 필드).
type InquirePriceOutput1 struct {
	HtsKorIsnm       string          `json:"hts_kor_isnm"`              // HTS 한글 종목명
	FutsPrpr         decimal.Decimal `json:"futs_prpr"`                 // 선물 현재가
	FutsPrdyVrss     decimal.Decimal `json:"futs_prdy_vrss"`            // 선물 전일 대비
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`            // 전일 대비 부호
	FutsPrdyClpr     decimal.Decimal `json:"futs_prdy_clpr"`            // 선물 전일 종가
	FutsPrdyCtrt     float64         `json:"futs_prdy_ctrt,string"`     // 선물 전일 대비율
	AcmlVol          int64           `json:"acml_vol,string"`           // 누적 거래량
	AcmlTrPbmn       int64           `json:"acml_tr_pbmn,string"`       // 누적 거래 대금
	HtsOtstStplQty   int64           `json:"hts_otst_stpl_qty,string"`  // HTS 미결제 약정 수량
	OtstStplQtyIcdc  int64           `json:"otst_stpl_qty_icdc,string"` // 미결제 약정 수량 증감
	FutsOprc         decimal.Decimal `json:"futs_oprc"`                 // 선물 시가
	FutsHgpr         decimal.Decimal `json:"futs_hgpr"`                 // 선물 최고가
	FutsLwpr         decimal.Decimal `json:"futs_lwpr"`                 // 선물 최저가
	FutsMxpr         decimal.Decimal `json:"futs_mxpr"`                 // 선물 상한가
	FutsLlam         decimal.Decimal `json:"futs_llam"`                 // 선물 하한가
	Basis            decimal.Decimal `json:"basis"`                     // 베이시스
	FutsSdpr         decimal.Decimal `json:"futs_sdpr"`                 // 선물 기준가
	HtsThpr          decimal.Decimal `json:"hts_thpr"`                  // HTS 이론가
	Dprt             float64         `json:"dprt,string"`               // 괴리율
	CrbrAplyMxpr     decimal.Decimal `json:"crbr_aply_mxpr"`            // 서킷브레이커 적용 상한가
	CrbrAplyLlam     decimal.Decimal `json:"crbr_aply_llam"`            // 서킷브레이커 적용 하한가
	FutsLastTrDate   string          `json:"futs_last_tr_date"`         // 선물 최종 거래 일자
	HtsRmnnDynu      string          `json:"hts_rmnn_dynu"`             // HTS 잔존 일수
	FutsLstnMedmHgpr decimal.Decimal `json:"futs_lstn_medm_hgpr"`       // 선물 상장 중 최고가
	FutsLstnMedmLwpr decimal.Decimal `json:"futs_lstn_medm_lwpr"`       // 선물 상장 중 최저가
	DeltaVal         float64         `json:"delta_val,string"`          // 델타 값 (옵션 지표)
	Gama             float64         `json:"gama,string"`               // 감마 (옵션 지표)
	Theta            float64         `json:"theta,string"`              // 세타 (옵션 지표)
	Vega             float64         `json:"vega,string"`               // 베가 (옵션 지표)
	Rho              float64         `json:"rho,string"`                // 로우 (옵션 지표)
	HistVltl         float64         `json:"hist_vltl,string"`          // 역사적 변동성
	HtsIntsVltl      float64         `json:"hts_ints_vltl,string"`      // HTS 내재 변동성
	MrktBasis        decimal.Decimal `json:"mrkt_basis"`                // 시장 베이시스
	Acpr             decimal.Decimal `json:"acpr"`                      // 행사가
}

// InquirePriceOutput2 는 선물 기초자산 현물 정보 (6 필드).
type InquirePriceOutput2 struct {
	BstpClsCode      string          `json:"bstp_cls_code"`              // 업종 구분 코드
	HtsKorIsnm       string          `json:"hts_kor_isnm"`               // HTS 한글 종목명
	BstpNmixPrpr     decimal.Decimal `json:"bstp_nmix_prpr"`             // 업종 지수 현재가
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`             // 전일 대비 부호
	BstpNmixPrdyVrss decimal.Decimal `json:"bstp_nmix_prdy_vrss"`        // 업종 지수 전일 대비
	BstpNmixPrdyCtrt float64         `json:"bstp_nmix_prdy_ctrt,string"` // 업종 지수 전일 대비율
}

// InquirePriceData 는 선물옵션 시세 응답 (output1 + output2 + output3).
type InquirePriceData struct {
	Output1 InquirePriceOutput1 `json:"output1"`
	Output2 InquirePriceOutput2 `json:"output2"`
	Output3 InquirePriceOutput2 `json:"output3"` // output2 와 동일 구조
}

type inquirePriceResponse struct {
	RtCd    string              `json:"rt_cd"`
	MsgCd   string              `json:"msg_cd"`
	Msg1    string              `json:"msg1"`
	Output1 InquirePriceOutput1 `json:"output1"`
	Output2 InquirePriceOutput2 `json:"output2"`
	Output3 InquirePriceOutput2 `json:"output3"`
}

// InquirePrice 는 선물옵션 시세 조회 (FHMIF10000000).
//
// marketCode: 조건 시장 분류 코드 (F:지수선물, O:지수옵션, JF:주식선물, JO:주식옵션, CF:상품선물, CM:야간선물, EU:야간옵션)
// code: 입력 종목코드 (예: 101S06)
//
// KIS API: GET /uapi/domestic-futureoption/v1/quotations/inquire-price
func (c *Client) InquirePrice(ctx context.Context, marketCode, code string) (*InquirePriceData, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-futureoption/v1/quotations/inquire-price",
		TrID:     "FHMIF10000000",
		CustType: "P",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": marketCode,
			"FID_INPUT_ISCD":         code,
		},
	})
	if err != nil {
		return nil, err
	}
	var res inquirePriceResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse InquirePrice: %w", err)
	}
	return &InquirePriceData{
		Output1: res.Output1,
		Output2: res.Output2,
		Output3: res.Output3,
	}, nil
}

// ─── EP2: InquireAskingPrice ──────────────────────────────────────────────────

// InquireAskingPriceOutput1 는 선물옵션 시세호가 요약 정보 (8 필드).
type InquireAskingPriceOutput1 struct {
	HtsKorIsnm   string          `json:"hts_kor_isnm"`          // HTS 한글 종목명
	FutsPrpr     decimal.Decimal `json:"futs_prpr"`             // 선물 현재가
	PrdyVrssSign string          `json:"prdy_vrss_sign"`        // 전일 대비 부호
	FutsPrdyVrss decimal.Decimal `json:"futs_prdy_vrss"`        // 선물 전일 대비
	FutsPrdyCtrt float64         `json:"futs_prdy_ctrt,string"` // 선물 전일 대비율
	AcmlVol      int64           `json:"acml_vol,string"`       // 누적 거래량
	FutsPrdyClpr decimal.Decimal `json:"futs_prdy_clpr"`        // 선물 전일 종가
	FutsShrnIscd string          `json:"futs_shrn_iscd"`        // 선물 단축 종목코드
}

// InquireAskingPriceOutput2Item 는 선물옵션 호가 스냅샷 (35 필드).
type InquireAskingPriceOutput2Item struct {
	FutsAskp1     decimal.Decimal `json:"futs_askp1"`             // 선물 매도호가1
	FutsAskp2     decimal.Decimal `json:"futs_askp2"`             // 선물 매도호가2
	FutsAskp3     decimal.Decimal `json:"futs_askp3"`             // 선물 매도호가3
	FutsAskp4     decimal.Decimal `json:"futs_askp4"`             // 선물 매도호가4
	FutsAskp5     decimal.Decimal `json:"futs_askp5"`             // 선물 매도호가5
	FutsBidp1     decimal.Decimal `json:"futs_bidp1"`             // 선물 매수호가1
	FutsBidp2     decimal.Decimal `json:"futs_bidp2"`             // 선물 매수호가2
	FutsBidp3     decimal.Decimal `json:"futs_bidp3"`             // 선물 매수호가3
	FutsBidp4     decimal.Decimal `json:"futs_bidp4"`             // 선물 매수호가4
	FutsBidp5     decimal.Decimal `json:"futs_bidp5"`             // 선물 매수호가5
	AskpRsqn1     int64           `json:"askp_rsqn1,string"`      // 매도호가 잔량1
	AskpRsqn2     int64           `json:"askp_rsqn2,string"`      // 매도호가 잔량2
	AskpRsqn3     int64           `json:"askp_rsqn3,string"`      // 매도호가 잔량3
	AskpRsqn4     int64           `json:"askp_rsqn4,string"`      // 매도호가 잔량4
	AskpRsqn5     int64           `json:"askp_rsqn5,string"`      // 매도호가 잔량5
	BidpRsqn1     int64           `json:"bidp_rsqn1,string"`      // 매수호가 잔량1
	BidpRsqn2     int64           `json:"bidp_rsqn2,string"`      // 매수호가 잔량2
	BidpRsqn3     int64           `json:"bidp_rsqn3,string"`      // 매수호가 잔량3
	BidpRsqn4     int64           `json:"bidp_rsqn4,string"`      // 매수호가 잔량4
	BidpRsqn5     int64           `json:"bidp_rsqn5,string"`      // 매수호가 잔량5
	AskpCsnu1     int64           `json:"askp_csnu1,string"`      // 매도호가 건수1
	AskpCsnu2     int64           `json:"askp_csnu2,string"`      // 매도호가 건수2
	AskpCsnu3     int64           `json:"askp_csnu3,string"`      // 매도호가 건수3
	AskpCsnu4     int64           `json:"askp_csnu4,string"`      // 매도호가 건수4
	AskpCsnu5     int64           `json:"askp_csnu5,string"`      // 매도호가 건수5
	BidpCsnu1     int64           `json:"bidp_csnu1,string"`      // 매수호가 건수1
	BidpCsnu2     int64           `json:"bidp_csnu2,string"`      // 매수호가 건수2
	BidpCsnu3     int64           `json:"bidp_csnu3,string"`      // 매수호가 건수3
	BidpCsnu4     int64           `json:"bidp_csnu4,string"`      // 매수호가 건수4
	BidpCsnu5     int64           `json:"bidp_csnu5,string"`      // 매수호가 건수5
	TotalAskpRsqn int64           `json:"total_askp_rsqn,string"` // 총 매도호가 잔량
	TotalBidpRsqn int64           `json:"total_bidp_rsqn,string"` // 총 매수호가 잔량
	TotalAskpCsnu int64           `json:"total_askp_csnu,string"` // 총 매도호가 건수
	TotalBidpCsnu int64           `json:"total_bidp_csnu,string"` // 총 매수호가 건수
	AsprAcptHour  string          `json:"aspr_acpt_hour"`         // 호가 접수 시간
}

// InquireAskingPriceData 는 선물옵션 시세호가 응답 (output1 + output2[]).
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

// InquireAskingPrice 는 선물옵션 시세호가 조회 (FHMIF10010000).
//
// marketCode: 조건 시장 분류 코드 (F:지수선물, O:지수옵션, JF:주식선물, JO:주식옵션, CF:상품선물, CM:야간선물, EU:야간옵션)
// code: 입력 종목코드 (예: 101S06)
//
// KIS API: GET /uapi/domestic-futureoption/v1/quotations/inquire-asking-price
func (c *Client) InquireAskingPrice(ctx context.Context, marketCode, code string) (*InquireAskingPriceData, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-futureoption/v1/quotations/inquire-asking-price",
		TrID:     "FHMIF10010000",
		CustType: "P",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": marketCode,
			"FID_INPUT_ISCD":         code,
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
