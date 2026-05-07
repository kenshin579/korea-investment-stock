package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// NearNewHighlow 는 국내주식 신고/신저근접종목 상위 (FHPST01870000) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_신고_신저근접종목_상위.md
// path: /uapi/domestic-stock/v1/ranking/near-new-highlow
//
// 최대 30건 확인 가능. 신고 근접 (PrcClsCode="0") 또는 신저 근접 (PrcClsCode="1").
type NearNewHighlow struct {
	Output []NearNewHighlowItem `json:"output"`
}

// NearNewHighlowItem 은 신고/신저근접종목 상위 응답의 한 행.
type NearNewHighlowItem struct {
	HtsKorIsnm   string          `json:"hts_kor_isnm"`          // HTS 한글 종목명
	MkscShrnIscd string          `json:"mksc_shrn_iscd"`        // 유가증권 단축 종목코드
	StckPrpr     decimal.Decimal `json:"stck_prpr"`             // 주식 현재가
	PrdyVrssSign string          `json:"prdy_vrss_sign"`        // 전일 대비 부호
	PrdyVrss     decimal.Decimal `json:"prdy_vrss"`             // 전일 대비
	PrdyCtrt     float64         `json:"prdy_ctrt,string"`      // 전일 대비율
	Askp         decimal.Decimal `json:"askp"`                  // 매도호가
	AskpRsqn1    int64           `json:"askp_rsqn1,string"`     // 매도호가 잔량1
	Bidp         decimal.Decimal `json:"bidp"`                  // 매수호가
	BidpRsqn1    int64           `json:"bidp_rsqn1,string"`     // 매수호가 잔량1
	AcmlVol      int64           `json:"acml_vol,string"`       // 누적 거래량
	NewHgpr      decimal.Decimal `json:"new_hgpr"`              // 신 최고가
	HprcNearRate float64         `json:"hprc_near_rate,string"` // 고가 근접 비율
	NewLwpr      decimal.Decimal `json:"new_lwpr"`              // 신 최저가
	LwprNearRate float64         `json:"lwpr_near_rate,string"` // 저가 근접 비율
	StckSdpr     decimal.Decimal `json:"stck_sdpr"`             // 주식 기준가
}

// InquireNearNewHighlowParams 는 신고/신저근접종목 상위 조회 파라미터.
//
// FID_COND_SCR_DIV_CODE = "20187" 고정 (사용자 변경 불가).
type InquireNearNewHighlowParams struct {
	MarketCode   string // fid_cond_mrkt_div_code — "J":KRX. 빈 값=>"J"
	DivClsCode   string // fid_div_cls_code — 0:전체, 1:관리종목, 2:투자주의, 3:투자경고. 빈 값=>"0"
	InputCnt1    string // fid_input_cnt_1 — 괴리율 최소. 빈 값=>"0"
	InputCnt2    string // fid_input_cnt_2 — 괴리율 최대. 빈 값=>"100"
	PrcClsCode   string // fid_prc_cls_code — 0:신고근접, 1:신저근접. 빈 값=>"0"
	InputISCD    string // fid_input_iscd — 0000:전체, 0001:거래소, 1001:코스닥, 2001:코스피200, 4001:KRX100
	TrgtClsCode  string // fid_trgt_cls_code — 0:전체. 빈 값=>"0"
	TrgtExlsCode string // fid_trgt_exls_cls_code — 0:전체. 빈 값=>"0"
	AplyRangVol  string // fid_aply_rang_vol — 0:전체, 100:100주 이상. 빈 값=>"0"
	AplyRangPrc1 string // fid_aply_rang_prc_1 — 가격 ~. 빈 값 OK
	AplyRangPrc2 string // fid_aply_rang_prc_2 — ~ 가격. 빈 값 OK
}

// InquireNearNewHighlow 는 국내주식 신고/신저근접종목 상위 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_신고_신저근접종목_상위.md
// path: /uapi/domestic-stock/v1/ranking/near-new-highlow (FHPST01870000)
//
// PrcClsCode="0" 신고 근접 / "1" 신저 근접. 최대 30건.
func (c *Client) InquireNearNewHighlow(ctx context.Context, params InquireNearNewHighlowParams) (*NearNewHighlow, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	div := params.DivClsCode
	if div == "" {
		div = "0"
	}
	cnt1 := params.InputCnt1
	if cnt1 == "" {
		cnt1 = "0"
	}
	cnt2 := params.InputCnt2
	if cnt2 == "" {
		cnt2 = "100"
	}
	prc := params.PrcClsCode
	if prc == "" {
		prc = "0"
	}
	tgt := params.TrgtClsCode
	if tgt == "" {
		tgt = "0"
	}
	tgtExcl := params.TrgtExlsCode
	if tgtExcl == "" {
		tgtExcl = "0"
	}
	vol := params.AplyRangVol
	if vol == "" {
		vol = "0"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ranking/near-new-highlow",
		TrID:   "FHPST01870000",
		Query: map[string]string{
			"fid_cond_mrkt_div_code": market,
			"fid_cond_scr_div_code":  "20187",
			"fid_div_cls_code":       div,
			"fid_input_cnt_1":        cnt1,
			"fid_input_cnt_2":        cnt2,
			"fid_prc_cls_code":       prc,
			"fid_input_iscd":         params.InputISCD,
			"fid_trgt_cls_code":      tgt,
			"fid_trgt_exls_cls_code": tgtExcl,
			"fid_aply_rang_vol":      vol,
			"fid_aply_rang_prc_1":    params.AplyRangPrc1,
			"fid_aply_rang_prc_2":    params.AplyRangPrc2,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res NearNewHighlow
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse NearNewHighlow: %w", err)
	}
	return &res, nil
}

// OvertimePrice 는 국내주식 시간외현재가 (FHPST02300000) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_시간외현재가.md
// path: /uapi/domestic-stock/v1/quotations/inquire-overtime-price
//
// 시간외 단일가 현재가 + 예상체결 + 상하한가 + 증거금비율 + 관리구분 등.
type OvertimePrice struct {
	Output OvertimePriceOutput `json:"output"`
}

// OvertimePriceOutput 은 시간외현재가 응답의 output object.
type OvertimePriceOutput struct {
	BstpKorIsnm              string          `json:"bstp_kor_isnm"`                   // 업종 한글 종목명
	MangIssuClsName          string          `json:"mang_issu_cls_name"`              // 관리 종목 구분 명
	OvtmUntpPrpr             decimal.Decimal `json:"ovtm_untp_prpr"`                  // 시간외 단일가 현재가
	OvtmUntpPrdyVrss         decimal.Decimal `json:"ovtm_untp_prdy_vrss"`             // 시간외 단일가 전일 대비
	OvtmUntpPrdyVrssSign     string          `json:"ovtm_untp_prdy_vrss_sign"`        // 시간외 단일가 전일 대비 부호
	OvtmUntpPrdyCtrt         float64         `json:"ovtm_untp_prdy_ctrt,string"`      // 시간외 단일가 전일 대비율
	OvtmUntpVol              int64           `json:"ovtm_untp_vol,string"`            // 시간외 단일가 거래량
	OvtmUntpTrPbmn           int64           `json:"ovtm_untp_tr_pbmn,string"`        // 시간외 단일가 거래 대금
	OvtmUntpMxpr             decimal.Decimal `json:"ovtm_untp_mxpr"`                  // 시간외 단일가 상한가
	OvtmUntpLlam             decimal.Decimal `json:"ovtm_untp_llam"`                  // 시간외 단일가 하한가
	OvtmUntpOprc             decimal.Decimal `json:"ovtm_untp_oprc"`                  // 시간외 단일가 시가2
	OvtmUntpHgpr             decimal.Decimal `json:"ovtm_untp_hgpr"`                  // 시간외 단일가 최고가
	OvtmUntpLwpr             decimal.Decimal `json:"ovtm_untp_lwpr"`                  // 시간외 단일가 최저가
	MargRate                 float64         `json:"marg_rate,string"`                // 증거금 비율
	OvtmUntpAntcCnpr         decimal.Decimal `json:"ovtm_untp_antc_cnpr"`             // 시간외 단일가 예상 체결가
	OvtmUntpAntcCntgVrss     decimal.Decimal `json:"ovtm_untp_antc_cntg_vrss"`        // 시간외 단일가 예상 체결 대비
	OvtmUntpAntcCntgVrssSign string          `json:"ovtm_untp_antc_cntg_vrss_sign"`   // 시간외 단일가 예상 체결 대비 부호
	OvtmUntpAntcCntgCtrt     float64         `json:"ovtm_untp_antc_cntg_ctrt,string"` // 시간외 단일가 예상 체결 대비율
	OvtmUntpAntcCnqn         int64           `json:"ovtm_untp_antc_cnqn,string"`      // 시간외 단일가 예상 체결량
	CrdtAbleYn               string          `json:"crdt_able_yn"`                    // 신용 가능 여부
	NewLstnClsName           string          `json:"new_lstn_cls_name"`               // 신규 상장 구분 명
	SltrYn                   string          `json:"sltr_yn"`                         // 정리매매 여부
	MangIssuYn               string          `json:"mang_issu_yn"`                    // 관리 종목 여부
	MrktWarnClsCode          string          `json:"mrkt_warn_cls_code"`              // 시장 경고 구분 코드
	TrhtYn                   string          `json:"trht_yn"`                         // 거래정지 여부
	VlntDealClsName          string          `json:"vlnt_deal_cls_name"`              // 임의 매매 구분 명
	OvtmUntpSdpr             decimal.Decimal `json:"ovtm_untp_sdpr"`                  // 시간외 단일가 기준가
	MrktWarnClsName          string          `json:"mrkt_warn_cls_name"`              // 시장 경고 구분 명
	RevlIssuReasName         string          `json:"revl_issu_reas_name"`             // 재평가 종목 사유 명
	InsnPbntYn               string          `json:"insn_pbnt_yn"`                    // 불성실 공시 여부
	FlngClsName              string          `json:"flng_cls_name"`                   // 락 구분 이름
	RprsMrktKorName          string          `json:"rprs_mrkt_kor_name"`              // 대표 시장 한글 명
	OvtmViClsCode            string          `json:"ovtm_vi_cls_code"`                // 시간외단일가VI적용구분코드
	Bidp                     decimal.Decimal `json:"bidp"`                            // 매수호가
	Askp                     decimal.Decimal `json:"askp"`                            // 매도호가
}

// InquireOvertimePriceParams 는 시간외현재가 조회 파라미터.
type InquireOvertimePriceParams struct {
	MarketCode string // FID_COND_MRKT_DIV_CODE — "J":KRX. 빈 값=>"J"
	Symbol     string // FID_INPUT_ISCD — 종목코드 (예 "005930")
}

// InquireOvertimePrice 는 국내주식 시간외현재가 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_시간외현재가.md
// path: /uapi/domestic-stock/v1/quotations/inquire-overtime-price (FHPST02300000)
func (c *Client) InquireOvertimePrice(ctx context.Context, params InquireOvertimePriceParams) (*OvertimePrice, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-overtime-price",
		TrID:   "FHPST02300000",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_INPUT_ISCD":         params.Symbol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res OvertimePrice
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse OvertimePrice: %w", err)
	}
	return &res, nil
}

// OvertimeAskingPrice 는 국내주식 시간외호가 (FHPST02300400) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_시간외호가.md
// path: /uapi/domestic-stock/v1/quotations/inquire-overtime-asking-price
//
// 시간외 단일가 10단계 호가/증감/잔량 + 정규장 총잔량. output1 만 존재.
type OvertimeAskingPrice struct {
	Output1 OvertimeAskingPriceOrderbook `json:"output1"`
}

// OvertimeAskingPriceOrderbook 은 시간외호가 응답 output1 — 10단계 호가+증감+잔량.
type OvertimeAskingPriceOrderbook struct {
	OvtmUntpLastHour string `json:"ovtm_untp_last_hour"` // 시간외 단일가 최종 시간 (HHMMSS)

	OvtmUntpAskp1  decimal.Decimal `json:"ovtm_untp_askp1"` // 시간외 단일가 매도호가1
	OvtmUntpAskp2  decimal.Decimal `json:"ovtm_untp_askp2"`
	OvtmUntpAskp3  decimal.Decimal `json:"ovtm_untp_askp3"`
	OvtmUntpAskp4  decimal.Decimal `json:"ovtm_untp_askp4"`
	OvtmUntpAskp5  decimal.Decimal `json:"ovtm_untp_askp5"`
	OvtmUntpAskp6  decimal.Decimal `json:"ovtm_untp_askp6"`
	OvtmUntpAskp7  decimal.Decimal `json:"ovtm_untp_askp7"`
	OvtmUntpAskp8  decimal.Decimal `json:"ovtm_untp_askp8"`
	OvtmUntpAskp9  decimal.Decimal `json:"ovtm_untp_askp9"`
	OvtmUntpAskp10 decimal.Decimal `json:"ovtm_untp_askp10"`

	OvtmUntpBidp1  decimal.Decimal `json:"ovtm_untp_bidp1"` // 시간외 단일가 매수호가1
	OvtmUntpBidp2  decimal.Decimal `json:"ovtm_untp_bidp2"`
	OvtmUntpBidp3  decimal.Decimal `json:"ovtm_untp_bidp3"`
	OvtmUntpBidp4  decimal.Decimal `json:"ovtm_untp_bidp4"`
	OvtmUntpBidp5  decimal.Decimal `json:"ovtm_untp_bidp5"`
	OvtmUntpBidp6  decimal.Decimal `json:"ovtm_untp_bidp6"`
	OvtmUntpBidp7  decimal.Decimal `json:"ovtm_untp_bidp7"`
	OvtmUntpBidp8  decimal.Decimal `json:"ovtm_untp_bidp8"`
	OvtmUntpBidp9  decimal.Decimal `json:"ovtm_untp_bidp9"`
	OvtmUntpBidp10 decimal.Decimal `json:"ovtm_untp_bidp10"`

	OvtmUntpAskpIcdc1  int64 `json:"ovtm_untp_askp_icdc1,string"` // 시간외 단일가 매도호가 증감1
	OvtmUntpAskpIcdc2  int64 `json:"ovtm_untp_askp_icdc2,string"`
	OvtmUntpAskpIcdc3  int64 `json:"ovtm_untp_askp_icdc3,string"`
	OvtmUntpAskpIcdc4  int64 `json:"ovtm_untp_askp_icdc4,string"`
	OvtmUntpAskpIcdc5  int64 `json:"ovtm_untp_askp_icdc5,string"`
	OvtmUntpAskpIcdc6  int64 `json:"ovtm_untp_askp_icdc6,string"`
	OvtmUntpAskpIcdc7  int64 `json:"ovtm_untp_askp_icdc7,string"`
	OvtmUntpAskpIcdc8  int64 `json:"ovtm_untp_askp_icdc8,string"`
	OvtmUntpAskpIcdc9  int64 `json:"ovtm_untp_askp_icdc9,string"`
	OvtmUntpAskpIcdc10 int64 `json:"ovtm_untp_askp_icdc10,string"`

	OvtmUntpBidpIcdc1  int64 `json:"ovtm_untp_bidp_icdc1,string"` // 시간외 단일가 매수호가 증감1
	OvtmUntpBidpIcdc2  int64 `json:"ovtm_untp_bidp_icdc2,string"`
	OvtmUntpBidpIcdc3  int64 `json:"ovtm_untp_bidp_icdc3,string"`
	OvtmUntpBidpIcdc4  int64 `json:"ovtm_untp_bidp_icdc4,string"`
	OvtmUntpBidpIcdc5  int64 `json:"ovtm_untp_bidp_icdc5,string"`
	OvtmUntpBidpIcdc6  int64 `json:"ovtm_untp_bidp_icdc6,string"`
	OvtmUntpBidpIcdc7  int64 `json:"ovtm_untp_bidp_icdc7,string"`
	OvtmUntpBidpIcdc8  int64 `json:"ovtm_untp_bidp_icdc8,string"`
	OvtmUntpBidpIcdc9  int64 `json:"ovtm_untp_bidp_icdc9,string"`
	OvtmUntpBidpIcdc10 int64 `json:"ovtm_untp_bidp_icdc10,string"`

	OvtmUntpAskpRsqn1  int64 `json:"ovtm_untp_askp_rsqn1,string"` // 시간외 단일가 매도호가 잔량1
	OvtmUntpAskpRsqn2  int64 `json:"ovtm_untp_askp_rsqn2,string"`
	OvtmUntpAskpRsqn3  int64 `json:"ovtm_untp_askp_rsqn3,string"`
	OvtmUntpAskpRsqn4  int64 `json:"ovtm_untp_askp_rsqn4,string"`
	OvtmUntpAskpRsqn5  int64 `json:"ovtm_untp_askp_rsqn5,string"`
	OvtmUntpAskpRsqn6  int64 `json:"ovtm_untp_askp_rsqn6,string"`
	OvtmUntpAskpRsqn7  int64 `json:"ovtm_untp_askp_rsqn7,string"`
	OvtmUntpAskpRsqn8  int64 `json:"ovtm_untp_askp_rsqn8,string"`
	OvtmUntpAskpRsqn9  int64 `json:"ovtm_untp_askp_rsqn9,string"`
	OvtmUntpAskpRsqn10 int64 `json:"ovtm_untp_askp_rsqn10,string"`

	OvtmUntpBidpRsqn1  int64 `json:"ovtm_untp_bidp_rsqn1,string"` // 시간외 단일가 매수호가 잔량1
	OvtmUntpBidpRsqn2  int64 `json:"ovtm_untp_bidp_rsqn2,string"`
	OvtmUntpBidpRsqn3  int64 `json:"ovtm_untp_bidp_rsqn3,string"`
	OvtmUntpBidpRsqn4  int64 `json:"ovtm_untp_bidp_rsqn4,string"`
	OvtmUntpBidpRsqn5  int64 `json:"ovtm_untp_bidp_rsqn5,string"`
	OvtmUntpBidpRsqn6  int64 `json:"ovtm_untp_bidp_rsqn6,string"`
	OvtmUntpBidpRsqn7  int64 `json:"ovtm_untp_bidp_rsqn7,string"`
	OvtmUntpBidpRsqn8  int64 `json:"ovtm_untp_bidp_rsqn8,string"`
	OvtmUntpBidpRsqn9  int64 `json:"ovtm_untp_bidp_rsqn9,string"`
	OvtmUntpBidpRsqn10 int64 `json:"ovtm_untp_bidp_rsqn10,string"`

	OvtmUntpTotalAskpRsqn     int64 `json:"ovtm_untp_total_askp_rsqn,string"`      // 시간외 단일가 총 매도호가 잔량
	OvtmUntpTotalBidpRsqn     int64 `json:"ovtm_untp_total_bidp_rsqn,string"`      // 시간외 단일가 총 매수호가 잔량
	OvtmUntpTotalAskpRsqnIcdc int64 `json:"ovtm_untp_total_askp_rsqn_icdc,string"` // 시간외 단일가 총 매도호가 잔량 증감
	OvtmUntpTotalBidpRsqnIcdc int64 `json:"ovtm_untp_total_bidp_rsqn_icdc,string"` // 시간외 단일가 총 매수호가 잔량 증감
	OvtmUntpNtbyBidpRsqn      int64 `json:"ovtm_untp_ntby_bidp_rsqn,string"`       // 시간외 단일가 순매수 호가 잔량
	TotalAskpRsqn             int64 `json:"total_askp_rsqn,string"`                // 총 매도호가 잔량 (정규장)
	TotalBidpRsqn             int64 `json:"total_bidp_rsqn,string"`                // 총 매수호가 잔량 (정규장)
	TotalAskpRsqnIcdc         int64 `json:"total_askp_rsqn_icdc,string"`           // 총 매도호가 잔량 증감
	TotalBidpRsqnIcdc         int64 `json:"total_bidp_rsqn_icdc,string"`           // 총 매수호가 잔량 증감
	OvtmTotalAskpRsqn         int64 `json:"ovtm_total_askp_rsqn,string"`           // 시간외 총 매도호가 잔량
	OvtmTotalBidpRsqn         int64 `json:"ovtm_total_bidp_rsqn,string"`           // 시간외 총 매수호가 잔량
	OvtmTotalAskpIcdc         int64 `json:"ovtm_total_askp_icdc,string"`           // 시간외 총 매도호가 증감
	OvtmTotalBidpIcdc         int64 `json:"ovtm_total_bidp_icdc,string"`           // 시간외 총 매수호가 증감
}

// InquireOvertimeAskingPriceParams 는 시간외호가 조회 파라미터.
type InquireOvertimeAskingPriceParams struct {
	Symbol     string // FID_INPUT_ISCD — 종목코드 (예 "005930")
	MarketCode string // FID_COND_MRKT_DIV_CODE — "J":KRX. 빈 값=>"J"
}

// InquireOvertimeAskingPrice 는 국내주식 시간외호가 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_시간외호가.md
// path: /uapi/domestic-stock/v1/quotations/inquire-overtime-asking-price (FHPST02300400)
//
// 시간외 단일가 10단계 호가/증감/잔량 (총 60 fields) + 시간외/정규장 총잔량.
func (c *Client) InquireOvertimeAskingPrice(ctx context.Context, params InquireOvertimeAskingPriceParams) (*OvertimeAskingPrice, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-overtime-asking-price",
		TrID:   "FHPST02300400",
		Query: map[string]string{
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_COND_MRKT_DIV_CODE": market,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res OvertimeAskingPrice
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse OvertimeAskingPrice: %w", err)
	}
	return &res, nil
}

// OvertimeVolume 은 국내주식 시간외거래량순위 (FHPST02350000) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_시간외거래량순위.md
// path: /uapi/domestic-stock/v1/ranking/overtime-volume
//
// output1: 거래소/코스닥 합계 (4 fields). output2: 종목별 시간외 거래량 순위 array.
type OvertimeVolume struct {
	Output1 OvertimeVolumeSummary `json:"output1"`
	Output2 []OvertimeVolumeItem  `json:"output2"`
}

// OvertimeVolumeSummary 는 시간외거래량순위 output1 — 시장 전체 합계.
type OvertimeVolumeSummary struct {
	OvtmUntpExchVol      int64 `json:"ovtm_untp_exch_vol,string"`       // 시간외 단일가 거래소 거래량
	OvtmUntpExchTrPbmn   int64 `json:"ovtm_untp_exch_tr_pbmn,string"`   // 시간외 단일가 거래소 거래대금
	OvtmUntpKosdaqVol    int64 `json:"ovtm_untp_kosdaq_vol,string"`     // 시간외 단일가 KOSDAQ 거래량
	OvtmUntpKosdaqTrPbmn int64 `json:"ovtm_untp_kosdaq_tr_pbmn,string"` // 시간외 단일가 KOSDAQ 거래대금
}

// OvertimeVolumeItem 은 시간외거래량순위 output2 의 한 행.
type OvertimeVolumeItem struct {
	StckShrnIscd         string          `json:"stck_shrn_iscd"`                 // 주식 단축 종목코드
	HtsKorIsnm           string          `json:"hts_kor_isnm"`                   // HTS 한글 종목명
	OvtmUntpPrpr         decimal.Decimal `json:"ovtm_untp_prpr"`                 // 시간외 단일가 현재가
	OvtmUntpPrdyVrss     decimal.Decimal `json:"ovtm_untp_prdy_vrss"`            // 시간외 단일가 전일 대비
	OvtmUntpPrdyVrssSign string          `json:"ovtm_untp_prdy_vrss_sign"`       // 시간외 단일가 전일 대비 부호
	OvtmUntpPrdyCtrt     float64         `json:"ovtm_untp_prdy_ctrt,string"`     // 시간외 단일가 전일 대비율
	OvtmUntpSelnRsqn     int64           `json:"ovtm_untp_seln_rsqn,string"`     // 시간외 단일가 매도 잔량
	OvtmUntpShnuRsqn     int64           `json:"ovtm_untp_shnu_rsqn,string"`     // 시간외 단일가 매수 잔량
	OvtmUntpVol          int64           `json:"ovtm_untp_vol,string"`           // 시간외 단일가 거래량
	OvtmVrssAcmlVolRlim  float64         `json:"ovtm_vrss_acml_vol_rlim,string"` // 시간외 대비 누적 거래량 비중
	StckPrpr             decimal.Decimal `json:"stck_prpr"`                      // 주식 현재가 (정규장)
	AcmlVol              int64           `json:"acml_vol,string"`                // 누적 거래량 (정규장)
	Bidp                 decimal.Decimal `json:"bidp"`                           // 매수호가
	Askp                 decimal.Decimal `json:"askp"`                           // 매도호가
}

// InquireOvertimeVolumeParams 는 시간외거래량순위 조회 파라미터.
//
// FID_COND_SCR_DIV_CODE = "20235" 고정 (사용자 변경 불가).
type InquireOvertimeVolumeParams struct {
	MarketCode   string // FID_COND_MRKT_DIV_CODE — "J":KRX. 빈 값=>"J"
	InputISCD    string // FID_INPUT_ISCD — 0000:전체, 0001:코스피, 1001:코스닥
	RankSortCode string // FID_RANK_SORT_CLS_CODE — 0:매수잔량, 1:매도잔량, 2:거래량. 빈 값=>"2"
	InputPrice1  string // FID_INPUT_PRICE_1 — 가격 ~. 빈 값 OK
	InputPrice2  string // FID_INPUT_PRICE_2 — ~ 가격. 빈 값 OK
	VolCount     string // FID_VOL_CNT — 거래량 ~. 빈 값 OK
	TrgtClsCode  string // FID_TRGT_CLS_CODE — 공백 입력
	TrgtExlsCode string // FID_TRGT_EXLS_CLS_CODE — 공백 입력
}

// InquireOvertimeVolume 은 국내주식 시간외거래량순위 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_시간외거래량순위.md
// path: /uapi/domestic-stock/v1/ranking/overtime-volume (FHPST02350000)
//
// output1: 거래소/코스닥 합계. output2: 최대 30건 종목별 시간외 거래량 순위.
func (c *Client) InquireOvertimeVolume(ctx context.Context, params InquireOvertimeVolumeParams) (*OvertimeVolume, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	sort := params.RankSortCode
	if sort == "" {
		sort = "2"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ranking/overtime-volume",
		TrID:   "FHPST02350000",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_COND_SCR_DIV_CODE":  "20235",
			"FID_INPUT_ISCD":         params.InputISCD,
			"FID_RANK_SORT_CLS_CODE": sort,
			"FID_INPUT_PRICE_1":      params.InputPrice1,
			"FID_INPUT_PRICE_2":      params.InputPrice2,
			"FID_VOL_CNT":            params.VolCount,
			"FID_TRGT_CLS_CODE":      params.TrgtClsCode,
			"FID_TRGT_EXLS_CLS_CODE": params.TrgtExlsCode,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res OvertimeVolume
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse OvertimeVolume: %w", err)
	}
	return &res, nil
}

// OvertimeFluctuation 은 국내주식 시간외등락율순위 (FHPST02340000) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_시간외등락율순위.md
// path: /uapi/domestic-stock/v1/ranking/overtime-fluctuation
//
// output1: 상한/상승/보합/하한/하락 종목 수 + 거래량/대금 합계 (11 fields).
// output2: 종목별 시간외 등락율 순위 array (16 fields/item).
type OvertimeFluctuation struct {
	Output1 OvertimeFluctuationSummary `json:"output1"`
	Output2 []OvertimeFluctuationItem  `json:"output2"`
}

// OvertimeFluctuationSummary 는 시간외등락율순위 output1 — 시장 전체 통계.
type OvertimeFluctuationSummary struct {
	OvtmUntpUplmIssuCnt  int64 `json:"ovtm_untp_uplm_issu_cnt,string"`  // 시간외 단일가 상한 종목 수
	OvtmUntpAscnIssuCnt  int64 `json:"ovtm_untp_ascn_issu_cnt,string"`  // 시간외 단일가 상승 종목 수
	OvtmUntpStnrIssuCnt  int64 `json:"ovtm_untp_stnr_issu_cnt,string"`  // 시간외 단일가 보합 종목 수
	OvtmUntpLslmIssuCnt  int64 `json:"ovtm_untp_lslm_issu_cnt,string"`  // 시간외 단일가 하한 종목 수
	OvtmUntpDownIssuCnt  int64 `json:"ovtm_untp_down_issu_cnt,string"`  // 시간외 단일가 하락 종목 수
	OvtmUntpAcmlVol      int64 `json:"ovtm_untp_acml_vol,string"`       // 시간외 단일가 누적 거래량
	OvtmUntpAcmlTrPbmn   int64 `json:"ovtm_untp_acml_tr_pbmn,string"`   // 시간외 단일가 누적 거래대금
	OvtmUntpExchVol      int64 `json:"ovtm_untp_exch_vol,string"`       // 시간외 단일가 거래소 거래량
	OvtmUntpExchTrPbmn   int64 `json:"ovtm_untp_exch_tr_pbmn,string"`   // 시간외 단일가 거래소 거래대금
	OvtmUntpKosdaqVol    int64 `json:"ovtm_untp_kosdaq_vol,string"`     // 시간외 단일가 KOSDAQ 거래량
	OvtmUntpKosdaqTrPbmn int64 `json:"ovtm_untp_kosdaq_tr_pbmn,string"` // 시간외 단일가 KOSDAQ 거래대금
}

// OvertimeFluctuationItem 은 시간외등락율순위 output2 의 한 행.
type OvertimeFluctuationItem struct {
	MkscShrnIscd         string          `json:"mksc_shrn_iscd"`                 // 유가증권 단축 종목코드
	HtsKorIsnm           string          `json:"hts_kor_isnm"`                   // HTS 한글 종목명
	OvtmUntpPrpr         decimal.Decimal `json:"ovtm_untp_prpr"`                 // 시간외 단일가 현재가
	OvtmUntpPrdyVrss     decimal.Decimal `json:"ovtm_untp_prdy_vrss"`            // 시간외 단일가 전일 대비
	OvtmUntpPrdyVrssSign string          `json:"ovtm_untp_prdy_vrss_sign"`       // 시간외 단일가 전일 대비 부호
	OvtmUntpPrdyCtrt     float64         `json:"ovtm_untp_prdy_ctrt,string"`     // 시간외 단일가 전일 대비율
	OvtmUntpAskp1        decimal.Decimal `json:"ovtm_untp_askp1"`                // 시간외 단일가 매도호가1
	OvtmUntpSelnRsqn     int64           `json:"ovtm_untp_seln_rsqn,string"`     // 시간외 단일가 매도 잔량
	OvtmUntpBidp1        decimal.Decimal `json:"ovtm_untp_bidp1"`                // 시간외 단일가 매수호가1
	OvtmUntpShnuRsqn     int64           `json:"ovtm_untp_shnu_rsqn,string"`     // 시간외 단일가 매수 잔량
	OvtmUntpVol          int64           `json:"ovtm_untp_vol,string"`           // 시간외 단일가 거래량
	OvtmVrssAcmlVolRlim  float64         `json:"ovtm_vrss_acml_vol_rlim,string"` // 시간외 대비 누적 거래량 비중
	StckPrpr             decimal.Decimal `json:"stck_prpr"`                      // 주식 현재가 (정규장)
	AcmlVol              int64           `json:"acml_vol,string"`                // 누적 거래량 (정규장)
	Bidp                 decimal.Decimal `json:"bidp"`                           // 매수호가
	Askp                 decimal.Decimal `json:"askp"`                           // 매도호가
}

// InquireOvertimeFluctuationParams 는 시간외등락율순위 조회 파라미터.
//
// FID_COND_SCR_DIV_CODE = "20234" 고정 (사용자 변경 불가).
type InquireOvertimeFluctuationParams struct {
	MarketCode   string // FID_COND_MRKT_DIV_CODE — "J":KRX. 빈 값=>"J"
	MrktClsCode  string // FID_MRKT_CLS_CODE — 공백 입력
	InputISCD    string // FID_INPUT_ISCD — 0000:전체, 0001:코스피, 1001:코스닥
	DivClsCode   string // FID_DIV_CLS_CODE — 1:상한가, 2:상승률, 3:보합, 4:하한가, 5:하락률. 빈 값=>"2"
	InputPrice1  string // FID_INPUT_PRICE_1 — 가격 ~. 빈 값 OK
	InputPrice2  string // FID_INPUT_PRICE_2 — ~ 가격. 빈 값 OK
	VolCount     string // FID_VOL_CNT — 거래량 ~. 빈 값 OK
	TrgtClsCode  string // FID_TRGT_CLS_CODE — 공백 입력
	TrgtExlsCode string // FID_TRGT_EXLS_CLS_CODE — 공백 입력
}

// InquireOvertimeFluctuation 은 국내주식 시간외등락율순위 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_시간외등락율순위.md
// path: /uapi/domestic-stock/v1/ranking/overtime-fluctuation (FHPST02340000)
//
// output1: 상한/상승/보합/하한/하락 종목 수 + 거래량/대금 합계.
// output2: 최대 30건 종목별 시간외 등락율 순위.
func (c *Client) InquireOvertimeFluctuation(ctx context.Context, params InquireOvertimeFluctuationParams) (*OvertimeFluctuation, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	div := params.DivClsCode
	if div == "" {
		div = "2"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ranking/overtime-fluctuation",
		TrID:   "FHPST02340000",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_MRKT_CLS_CODE":      params.MrktClsCode,
			"FID_COND_SCR_DIV_CODE":  "20234",
			"FID_INPUT_ISCD":         params.InputISCD,
			"FID_DIV_CLS_CODE":       div,
			"FID_INPUT_PRICE_1":      params.InputPrice1,
			"FID_INPUT_PRICE_2":      params.InputPrice2,
			"FID_VOL_CNT":            params.VolCount,
			"FID_TRGT_CLS_CODE":      params.TrgtClsCode,
			"FID_TRGT_EXLS_CLS_CODE": params.TrgtExlsCode,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res OvertimeFluctuation
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse OvertimeFluctuation: %w", err)
	}
	return &res, nil
}

// VolumePower 는 체결강도 상위 (FHPST01680000) 응답.
//
// 한투 docs: docs/api/국내주식/체결강도상위.md
// path: /uapi/domestic-stock/v1/ranking/volume-power
//
// 주의: 모든 query 파라미터가 lowercase fid_* (대문자 FID_ 아님).
type VolumePower struct {
	Output []VolumePowerItem `json:"output"`
}

// VolumePowerItem 은 응답의 output 한 행 (11 fields).
type VolumePowerItem struct {
	StckShrnIscd string          `json:"stck_shrn_iscd"`        // 주식 단축 종목코드
	DataRank     string          `json:"data_rank"`             // 데이터 순위
	HtsKorIsnm   string          `json:"hts_kor_isnm"`          // HTS 한글 종목명
	StckPrpr     decimal.Decimal `json:"stck_prpr"`             // 주식 현재가
	PrdyVrss     decimal.Decimal `json:"prdy_vrss"`             // 전일 대비
	PrdyVrssSign string          `json:"prdy_vrss_sign"`        // 전일 대비 부호
	PrdyCtrt     float64         `json:"prdy_ctrt,string"`      // 전일 대비율
	AcmlVol      int64           `json:"acml_vol,string"`       // 누적 거래량
	TdayRltv     float64         `json:"tday_rltv,string"`      // 당일 체결강도
	SelnCnqnSmtn int64           `json:"seln_cnqn_smtn,string"` // 매도 체결량 합계
	ShnuCnqnSmtn int64           `json:"shnu_cnqn_smtn,string"` // 매수 체결량 합계
}

// InquireVolumePowerParams 는 체결강도 상위 조회 파라미터.
type InquireVolumePowerParams struct {
	MarketCode     string // fid_cond_mrkt_div_code — 빈 값=>"J" (lowercase wire key)
	CondScrDivCode string // fid_cond_scr_div_code — 빈 값=>"20168"
	Symbol         string // fid_input_iscd — 0000:전체/0001:코스피/1001:코스닥
	DivClsCode     string // fid_div_cls_code
	Price1         string // fid_input_price_1
	Price2         string // fid_input_price_2
	VolCnt         string // fid_vol_cnt
	TrgtClsCode    string // fid_trgt_cls_code
	TrgtExlsCode   string // fid_trgt_exls_cls_code
}

// InquireVolumePower 는 체결강도 상위 호출.
//
// 한투 docs: docs/api/국내주식/체결강도상위.md
// path: /uapi/domestic-stock/v1/ranking/volume-power (FHPST01680000)
func (c *Client) InquireVolumePower(ctx context.Context, params InquireVolumePowerParams) (*VolumePower, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "20168"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ranking/volume-power",
		TrID:   "FHPST01680000",
		Query: map[string]string{
			"fid_cond_mrkt_div_code": market,
			"fid_cond_scr_div_code":  scrDiv,
			"fid_input_iscd":         params.Symbol,
			"fid_div_cls_code":       params.DivClsCode,
			"fid_input_price_1":      params.Price1,
			"fid_input_price_2":      params.Price2,
			"fid_vol_cnt":            params.VolCnt,
			"fid_trgt_cls_code":      params.TrgtClsCode,
			"fid_trgt_exls_cls_code": params.TrgtExlsCode,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res VolumePower
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse VolumePower: %w", err)
	}
	return &res, nil
}

// BulkTransNum 은 대량체결건수 상위 (FHKST190900C0) 응답.
//
// 한투 docs: docs/api/국내주식/대량체결건수상위.md
// path: /uapi/domestic-stock/v1/ranking/bulk-trans-num
//
// 주의 1: 모든 query 파라미터가 lowercase fid_* (대문자 FID_ 아님).
// 주의 2: 종목코드 필드는 mksc_shrn_iscd (시장구분 포함, stck_shrn_iscd 아님).
type BulkTransNum struct {
	Output []BulkTransNumItem `json:"output"`
}

// BulkTransNumItem 은 응답의 output 한 행 (11 fields).
type BulkTransNumItem struct {
	MkscShrnIscd string          `json:"mksc_shrn_iscd"`        // 시장구분+단축 종목코드
	DataRank     string          `json:"data_rank"`             // 데이터 순위
	HtsKorIsnm   string          `json:"hts_kor_isnm"`          // HTS 한글 종목명
	StckPrpr     decimal.Decimal `json:"stck_prpr"`             // 주식 현재가
	PrdyVrssSign string          `json:"prdy_vrss_sign"`        // 전일 대비 부호
	PrdyVrss     decimal.Decimal `json:"prdy_vrss"`             // 전일 대비
	PrdyCtrt     float64         `json:"prdy_ctrt,string"`      // 전일 대비율
	AcmlVol      int64           `json:"acml_vol,string"`       // 누적 거래량
	ShnuCntgCsnu int64           `json:"shnu_cntg_csnu,string"` // 매수 체결 건수
	SelnCntgCsnu int64           `json:"seln_cntg_csnu,string"` // 매도 체결 건수
	NtbyCnqn     int64           `json:"ntby_cnqn,string"`      // 순매수 체결량
}

// InquireBulkTransNumParams 는 대량체결건수 상위 조회 파라미터.
type InquireBulkTransNumParams struct {
	MarketCode     string // fid_cond_mrkt_div_code — 빈 값=>"J"
	CondScrDivCode string // fid_cond_scr_div_code — 빈 값=>"11909"
	Symbol         string // fid_input_iscd
	DivClsCode     string // fid_div_cls_code
	RankSortCode   string // fid_rank_sort_cls_code
	BlngClsCode    string // fid_blng_cls_code
	TrgtClsCode    string // fid_trgt_cls_code
	TrgtExlsCode   string // fid_trgt_exls_cls_code
	InputPrice1    string // fid_input_price_1
	InputPrice2    string // fid_input_price_2
	VolCnt         string // fid_vol_cnt
	InputDate1     string // fid_input_date_1
}

// InquireBulkTransNum 은 대량체결건수 상위 호출.
//
// 한투 docs: docs/api/국내주식/대량체결건수상위.md
// path: /uapi/domestic-stock/v1/ranking/bulk-trans-num (FHKST190900C0)
func (c *Client) InquireBulkTransNum(ctx context.Context, params InquireBulkTransNumParams) (*BulkTransNum, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "11909"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ranking/bulk-trans-num",
		TrID:   "FHKST190900C0",
		Query: map[string]string{
			"fid_cond_mrkt_div_code": market,
			"fid_cond_scr_div_code":  scrDiv,
			"fid_input_iscd":         params.Symbol,
			"fid_div_cls_code":       params.DivClsCode,
			"fid_rank_sort_cls_code": params.RankSortCode,
			"fid_blng_cls_code":      params.BlngClsCode,
			"fid_trgt_cls_code":      params.TrgtClsCode,
			"fid_trgt_exls_cls_code": params.TrgtExlsCode,
			"fid_input_price_1":      params.InputPrice1,
			"fid_input_price_2":      params.InputPrice2,
			"fid_vol_cnt":            params.VolCnt,
			"fid_input_date_1":       params.InputDate1,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res BulkTransNum
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse BulkTransNum: %w", err)
	}
	return &res, nil
}

// TradprtByamt 는 체결금액별 매매비중 (FHKST111900C0) 응답.
//
// 한투 docs: docs/api/국내주식/체결금액별매매비중.md
// path: /uapi/domestic-stock/v1/quotations/tradprt-byamt
//
// 주의: whol_shun_vol_rate 필드명은 KIS wire format typo (shun ≠ shnu).
// 실제 의미는 전체 매수 거래량 비율이나 KIS wire key 그대로 보존.
type TradprtByamt struct {
	Output []TradprtByamtItem `json:"output"`
}

// TradprtByamtItem 은 응답의 output 한 행 (11 fields).
type TradprtByamtItem struct {
	PrprName        string          `json:"prpr_name"`                 // 체결금액 구간명
	SmtnAvrgPrpr    decimal.Decimal `json:"smtn_avrg_prpr"`            // 합산 평균 가격
	AcmlVol         int64           `json:"acml_vol,string"`           // 누적 거래량
	WholNtbyQtyRate float64         `json:"whol_ntby_qty_rate,string"` // 전체 순매수 수량 비율
	NtbyCntgCsnu    int64           `json:"ntby_cntg_csnu,string"`     // 순매수 체결 건수
	SelnCnqnSmtn    int64           `json:"seln_cnqn_smtn,string"`     // 매도 체결량 합계
	WholSelnVolRate float64         `json:"whol_seln_vol_rate,string"` // 전체 매도 거래량 비율
	SelnCntgCsnu    int64           `json:"seln_cntg_csnu,string"`     // 매도 체결 건수
	ShnuCnqnSmtn    int64           `json:"shnu_cnqn_smtn,string"`     // 매수 체결량 합계
	WholShunVolRate float64         `json:"whol_shun_vol_rate,string"` // 전체 매수 거래량 비율 (KIS typo 보존)
	ShnuCntgCsnu    int64           `json:"shnu_cntg_csnu,string"`     // 매수 체결 건수
}

// InquireTradprtByamtParams 는 체결금액별 매매비중 조회 파라미터.
type InquireTradprtByamtParams struct {
	MarketCode     string // FID_COND_MRKT_DIV_CODE — 빈 값=>"J"
	CondScrDivCode string // FID_COND_SCR_DIV_CODE — 빈 값=>"11119"
	Symbol         string // FID_INPUT_ISCD — 필수, 단축 종목코드
}

// InquireTradprtByamt 는 체결금액별 매매비중 호출.
//
// 한투 docs: docs/api/국내주식/체결금액별매매비중.md
// path: /uapi/domestic-stock/v1/quotations/tradprt-byamt (FHKST111900C0)
func (c *Client) InquireTradprtByamt(ctx context.Context, params InquireTradprtByamtParams) (*TradprtByamt, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "11119"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/tradprt-byamt",
		TrID:   "FHKST111900C0",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_COND_SCR_DIV_CODE":  scrDiv,
			"FID_INPUT_ISCD":         params.Symbol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res TradprtByamt
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse TradprtByamt: %w", err)
	}
	return &res, nil
}

// HtsTopView 는 HTS조회상위20종목 (HHMCM000100C0) 응답.
//
// 한투 docs: docs/api/국내주식/HTS조회상위20종목.md
// path: /uapi/domestic-stock/v1/ranking/hts-top-view
//
// 특이사항: 쿼리 파라미터 없음 (zero params), 응답 키가 output1 (not output).
type HtsTopView struct {
	Output1 HtsTopViewItem `json:"output1"`
}

// HtsTopViewItem 은 HTS조회상위20종목 응답의 단일 객체.
type HtsTopViewItem struct {
	MrktDivClsCode string `json:"mrkt_div_cls_code"` // 시장구분 (J:코스피, Q:코스닥)
	MkscShrnIscd   string `json:"mksc_shrn_iscd"`    // 종목코드
}

// InquireHtsTopViewParams 는 HTS조회상위20종목 조회 파라미터.
//
// 이 endpoint 는 쿼리 파라미터가 없음 (zero params).
type InquireHtsTopViewParams struct {
	// 의도적으로 비어있음 — KIS API 가 query params 받지 않음
}

// InquireHtsTopView 는 HTS조회상위20종목 호출.
//
// 한투 docs: docs/api/국내주식/HTS조회상위20종목.md
// path: /uapi/domestic-stock/v1/ranking/hts-top-view (HHMCM000100C0)
func (c *Client) InquireHtsTopView(ctx context.Context, _ InquireHtsTopViewParams) (*HtsTopView, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-stock/v1/ranking/hts-top-view",
		TrID:     "HHMCM000100C0",
		Query:    map[string]string{},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res HtsTopView
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse HtsTopView: %w", err)
	}
	return &res, nil
}

// PbarTraRatio 는 국내주식 매물대/거래비중 (FHPST01130000) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_매물대_거래비중.md
// path: /uapi/domestic-stock/v1/quotations/pbar-tratio
//
// dual output: output1 (단건 종합 11 필드) + output2 (가격대별 array, 4 필드/item).
type PbarTraRatio struct {
	Output1 PbarTraRatioSummary `json:"output1"`
	Output2 []PbarTraRatioItem  `json:"output2"`
}

// PbarTraRatioSummary 는 매물대 거래비중 응답 종합 (output1) — 11 필드.
type PbarTraRatioSummary struct {
	RprsMrktKorName string          `json:"rprs_mrkt_kor_name"` // 대표시장한글명
	StckShrnIscd    string          `json:"stck_shrn_iscd"`     // 주식단축종목코드
	HtsKorIsnm      string          `json:"hts_kor_isnm"`       // HTS한글종목명
	StckPrpr        decimal.Decimal `json:"stck_prpr"`          // 주식현재가
	PrdyVrssSign    string          `json:"prdy_vrss_sign"`     // 전일대비부호
	PrdyVrss        decimal.Decimal `json:"prdy_vrss"`          // 전일대비
	PrdyCtrt        float64         `json:"prdy_ctrt,string"`   // 전일대비율
	AcmlVol         int64           `json:"acml_vol,string"`    // 누적거래량
	PrdyVol         int64           `json:"prdy_vol,string"`    // 전일거래량
	WghnAvrgStckPrc decimal.Decimal `json:"wghn_avrg_stck_prc"` // 가중평균주식가격
	LstnStcn        int64           `json:"lstn_stcn,string"`   // 상장주수
}

// PbarTraRatioItem 은 매물대 거래비중 output2 의 한 행 (가격대별, 4 필드).
type PbarTraRatioItem struct {
	DataRank    string          `json:"data_rank"`            // 데이터순위
	StckPrpr    decimal.Decimal `json:"stck_prpr"`            // 주식현재가 (가격대)
	CntgVol     int64           `json:"cntg_vol,string"`      // 체결거래량
	AcmlVolRlim float64         `json:"acml_vol_rlim,string"` // 누적거래량비중
}

// InquirePbarTraRatioParams 는 매물대 거래비중 조회 파라미터.
type InquirePbarTraRatioParams struct {
	MarketCode     string // FID_COND_MRKT_DIV_CODE — default "J"
	CondScrDivCode string // FID_COND_SCR_DIV_CODE — default "11130"
	Symbol         string // FID_INPUT_ISCD — 종목코드
	InputHour1     string // FID_INPUT_HOUR_1 — 입력시간1
}

// InquirePbarTraRatio 는 국내주식 매물대 거래비중 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_매물대_거래비중.md
// path: /uapi/domestic-stock/v1/quotations/pbar-tratio (FHPST01130000)
func (c *Client) InquirePbarTraRatio(ctx context.Context, params InquirePbarTraRatioParams) (*PbarTraRatio, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "11130"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/pbar-tratio",
		TrID:   "FHPST01130000",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_COND_SCR_DIV_CODE":  scrDiv,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_INPUT_HOUR_1":       params.InputHour1,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res PbarTraRatio
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse PbarTraRatio: %w", err)
	}
	return &res, nil
}

// ExpPriceTrend 는 국내주식 예상체결가 추이 (FHPST01810000) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_예상체결가_추이.md
// path: /uapi/domestic-stock/v1/quotations/exp-price-trend
//
// dual output: output1 (예상체결가 종합) + output2 (시간별 체결 추이 array).
// 특이사항: query param wire keys 가 lowercase (fid_*).
type ExpPriceTrend struct {
	Output1 ExpPriceTrendSummary `json:"output1"`
	Output2 []ExpPriceTrendItem  `json:"output2"`
}

// ExpPriceTrendSummary 는 예상체결가 추이 응답 종합 (output1) — 7 필드.
type ExpPriceTrendSummary struct {
	RprsMrktKorName  string          `json:"rprs_mrkt_kor_name"`         // 대표시장한글명
	AntcCnpr         decimal.Decimal `json:"antc_cnpr"`                  // 예상 체결가
	AntcCntgVrssSign string          `json:"antc_cntg_vrss_sign"`        // 예상 체결 대비 부호
	AntcCntgVrss     decimal.Decimal `json:"antc_cntg_vrss"`             // 예상 체결 대비
	AntcCntgPrdyCtrt float64         `json:"antc_cntg_prdy_ctrt,string"` // 예상 체결 전일 대비율
	AntcVol          int64           `json:"antc_vol,string"`            // 예상 거래량
	AntcTrPbmn       int64           `json:"antc_tr_pbmn,string"`        // 예상 거래대금
}

// ExpPriceTrendItem 은 예상체결가 추이 output2 의 한 행 — 7 필드.
type ExpPriceTrendItem struct {
	StckBsopDate string          `json:"stck_bsop_date"`   // 주식 영업 일자
	StckCntgHour string          `json:"stck_cntg_hour"`   // 주식 체결 시간
	StckPrpr     decimal.Decimal `json:"stck_prpr"`        // 주식 현재가
	PrdyVrssSign string          `json:"prdy_vrss_sign"`   // 전일 대비 부호
	PrdyVrss     decimal.Decimal `json:"prdy_vrss"`        // 전일 대비
	PrdyCtrt     float64         `json:"prdy_ctrt,string"` // 전일 대비율
	AcmlVol      int64           `json:"acml_vol,string"`  // 누적 거래량
}

// InquireExpPriceTrendParams 는 예상체결가 추이 조회 파라미터.
//
// 특이사항: KIS API 가 query param wire keys 를 lowercase 로 받음 (fid_*).
type InquireExpPriceTrendParams struct {
	MarketCode     string // fid_cond_mrkt_div_code — default "J"
	CondScrDivCode string // fid_cond_scr_div_code — default "11810"
	Symbol         string // fid_input_iscd — 종목코드
}

// InquireExpPriceTrend 는 국내주식 예상체결가 추이 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_예상체결가_추이.md
// path: /uapi/domestic-stock/v1/quotations/exp-price-trend (FHPST01810000)
// 특이사항: query param wire keys 가 lowercase (fid_*).
func (c *Client) InquireExpPriceTrend(ctx context.Context, params InquireExpPriceTrendParams) (*ExpPriceTrend, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "11810"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/exp-price-trend",
		TrID:   "FHPST01810000",
		Query: map[string]string{
			"fid_cond_mrkt_div_code": market,
			"fid_cond_scr_div_code":  scrDiv,
			"fid_input_iscd":         params.Symbol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res ExpPriceTrend
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse ExpPriceTrend: %w", err)
	}
	return &res, nil
}

// ExpTransUpdown 는 국내주식 예상체결 상승/하락 상위 (FHPST01820000) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_예상체결_상승_하락상위.md
// path: /uapi/domestic-stock/v1/ranking/exp-trans-updown
//
// 특이사항: query param wire keys 가 lowercase (fid_*). 10개 query params.
type ExpTransUpdown struct {
	Output []ExpTransUpdownItem `json:"output"`
}

// ExpTransUpdownItem 은 예상체결 상승/하락 상위 한 행 — 15 필드.
type ExpTransUpdownItem struct {
	StckShrnIscd  string          `json:"stck_shrn_iscd"`         // 주식 단축 종목코드
	HtsKorIsnm    string          `json:"hts_kor_isnm"`           // HTS 한글 종목명
	StckPrpr      decimal.Decimal `json:"stck_prpr"`              // 주식 현재가
	PrdyVrss      decimal.Decimal `json:"prdy_vrss"`              // 전일 대비
	PrdyVrssSign  string          `json:"prdy_vrss_sign"`         // 전일 대비 부호
	PrdyCtrt      float64         `json:"prdy_ctrt,string"`       // 전일 대비율
	StckSdpr      decimal.Decimal `json:"stck_sdpr"`              // 주식 기준가
	SelnRsqn      int64           `json:"seln_rsqn,string"`       // 매도 잔량
	Askp          decimal.Decimal `json:"askp"`                   // 매도호가
	Bidp          decimal.Decimal `json:"bidp"`                   // 매수호가
	ShnuRsqn      int64           `json:"shnu_rsqn,string"`       // 매수 잔량
	CntgVol       int64           `json:"cntg_vol,string"`        // 체결 거래량
	AntcTrPbmn    int64           `json:"antc_tr_pbmn,string"`    // 체결 거래대금 (예상)
	TotalAskpRsqn int64           `json:"total_askp_rsqn,string"` // 총 매도호가 잔량
	TotalBidpRsqn int64           `json:"total_bidp_rsqn,string"` // 총 매수호가 잔량
}

// InquireExpTransUpdownParams 는 예상체결 상승/하락 상위 조회 파라미터.
//
// 특이사항: KIS API 가 query param wire keys 를 lowercase 로 받음 (fid_*).
type InquireExpTransUpdownParams struct {
	MarketCode     string // fid_cond_mrkt_div_code — default "J"
	CondScrDivCode string // fid_cond_scr_div_code — default "11820"
	Symbol         string // fid_input_iscd — 종목코드 / 0000:전체, 0001:코스피, 1001:코스닥
	DivClsCode     string // fid_div_cls_code — 0:전체, 1:관리종목 등
	RankSortCode   string // fid_rank_sort_cls_code — 0:상승률, 1:하락률 등
	InputPrice1    string // fid_input_price_1 — 가격 ~
	InputPrice2    string // fid_input_price_2 — ~ 가격
	VolCnt         string // fid_vol_cnt — 거래량 ~
	TrgtClsCode    string // fid_trgt_cls_code — 대상 구분
	TrgtExlsCode   string // fid_trgt_exls_cls_code — 대상 제외 구분
}

// InquireExpTransUpdown 는 국내주식 예상체결 상승/하락 상위 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_예상체결_상승_하락상위.md
// path: /uapi/domestic-stock/v1/ranking/exp-trans-updown (FHPST01820000)
// 특이사항: query param wire keys 가 lowercase (fid_*).
func (c *Client) InquireExpTransUpdown(ctx context.Context, params InquireExpTransUpdownParams) (*ExpTransUpdown, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "11820"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ranking/exp-trans-updown",
		TrID:   "FHPST01820000",
		Query: map[string]string{
			"fid_cond_mrkt_div_code": market,
			"fid_cond_scr_div_code":  scrDiv,
			"fid_input_iscd":         params.Symbol,
			"fid_div_cls_code":       params.DivClsCode,
			"fid_rank_sort_cls_code": params.RankSortCode,
			"fid_input_price_1":      params.InputPrice1,
			"fid_input_price_2":      params.InputPrice2,
			"fid_vol_cnt":            params.VolCnt,
			"fid_trgt_cls_code":      params.TrgtClsCode,
			"fid_trgt_exls_cls_code": params.TrgtExlsCode,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res ExpTransUpdown
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse ExpTransUpdown: %w", err)
	}
	return &res, nil
}

// ShortSale 은 국내주식 공매도 상위 (FHPST04820000) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_공매도_상위.md
// path: /uapi/domestic-stock/v1/ranking/short-sale
//
// 공매도 체결량/비중 상위 종목 목록.
type ShortSale struct {
	Output []ShortSaleItem `json:"output"`
}

// ShortSaleItem 은 공매도 상위 응답의 한 행.
type ShortSaleItem struct {
	MkscShrnIscd   string          `json:"mksc_shrn_iscd"`           // 유가증권 단축 종목코드
	HtsKorIsnm     string          `json:"hts_kor_isnm"`             // HTS 한글 종목명
	StckPrpr       decimal.Decimal `json:"stck_prpr"`                // 주식 현재가
	PrdyVrss       decimal.Decimal `json:"prdy_vrss"`                // 전일 대비
	PrdyVrssSign   string          `json:"prdy_vrss_sign"`           // 전일 대비 부호
	PrdyCtrt       float64         `json:"prdy_ctrt,string"`         // 전일 대비율
	AcmlVol        int64           `json:"acml_vol,string"`          // 누적 거래량
	AcmlTrPbmn     int64           `json:"acml_tr_pbmn,string"`      // 누적 거래 대금
	SstsCntgQty    int64           `json:"ssts_cntg_qty,string"`     // 공매도 체결 수량
	SstsVolRlim    float64         `json:"ssts_vol_rlim,string"`     // 공매도 거래량 비중
	SstsTrPbmn     int64           `json:"ssts_tr_pbmn,string"`      // 공매도 거래 대금
	SstsTrPbmnRlim float64         `json:"ssts_tr_pbmn_rlim,string"` // 공매도 거래 대금 비중
	StndDate1      string          `json:"stnd_date1"`               // 기준 일자1
	StndDate2      string          `json:"stnd_date2"`               // 기준 일자2
	AvrgPrc        decimal.Decimal `json:"avrg_prc"`                 // 평균가
}

// InquireShortSaleParams 는 공매도 상위 조회 파라미터.
type InquireShortSaleParams struct {
	AplyRangVol    string // FID_APLY_RANG_VOL — 적용범위 거래량. 빈 값 OK
	MarketCode     string // FID_COND_MRKT_DIV_CODE — "J":KRX. 빈 값=>"J"
	CondScrDivCode string // FID_COND_SCR_DIV_CODE — 고정 "20482". 빈 값=>"20482"
	Symbol         string // FID_INPUT_ISCD — 종목코드 (예 "005930")
	PeriodDivCode  string // FID_PERIOD_DIV_CODE — D:일, W:주, M:월
	InputCnt1      string // FID_INPUT_CNT_1 — 입력 개수1
	TrgtExlsCode   string // FID_TRGT_EXLS_CLS_CODE — 대상 제외 구분 코드
	TrgtClsCode    string // FID_TRGT_CLS_CODE — 대상 구분 코드
	AplyRangPrc1   string // FID_APLY_RANG_PRC_1 — 가격 ~. 빈 값 OK
	AplyRangPrc2   string // FID_APLY_RANG_PRC_2 — ~ 가격. 빈 값 OK
}

// InquireShortSale 은 국내주식 공매도 상위 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_공매도_상위.md
// path: /uapi/domestic-stock/v1/ranking/short-sale (FHPST04820000)
func (c *Client) InquireShortSale(ctx context.Context, params InquireShortSaleParams) (*ShortSale, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "20482"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ranking/short-sale",
		TrID:   "FHPST04820000",
		Query: map[string]string{
			"FID_APLY_RANG_VOL":      params.AplyRangVol,
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_COND_SCR_DIV_CODE":  scrDiv,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_PERIOD_DIV_CODE":    params.PeriodDivCode,
			"FID_INPUT_CNT_1":        params.InputCnt1,
			"FID_TRGT_EXLS_CLS_CODE": params.TrgtExlsCode,
			"FID_TRGT_CLS_CODE":      params.TrgtClsCode,
			"FID_APLY_RANG_PRC_1":    params.AplyRangPrc1,
			"FID_APLY_RANG_PRC_2":    params.AplyRangPrc2,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res ShortSale
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse ShortSale: %w", err)
	}
	return &res, nil
}

// DailyShortSale 은 국내주식 공매도 일별추이 (FHPST04830000) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_공매도_일별추이.md
// path: /uapi/domestic-stock/v1/quotations/daily-short-sale
//
// output1: 현재가 요약 (단일 객체), output2: 일자별 공매도 추이 목록.
type DailyShortSale struct {
	Output1 DailyShortSaleSummary `json:"output1"`
	Output2 []DailyShortSaleItem  `json:"output2"`
}

// DailyShortSaleSummary 는 공매도 일별추이 현재가 요약 (output1).
type DailyShortSaleSummary struct {
	StckPrpr     decimal.Decimal `json:"stck_prpr"`        // 주식 현재가
	PrdyVrss     decimal.Decimal `json:"prdy_vrss"`        // 전일 대비
	PrdyVrssSign string          `json:"prdy_vrss_sign"`   // 전일 대비 부호
	PrdyCtrt     float64         `json:"prdy_ctrt,string"` // 전일 대비율
	AcmlVol      int64           `json:"acml_vol,string"`  // 누적 거래량
	PrdyVol      int64           `json:"prdy_vol,string"`  // 전일 거래량
}

// DailyShortSaleItem 은 공매도 일별추이 응답의 한 행 (output2).
type DailyShortSaleItem struct {
	StckBsopDate        string          `json:"stck_bsop_date"`                 // 주식 영업 일자
	StckClpr            decimal.Decimal `json:"stck_clpr"`                      // 주식 종가
	PrdyVrss            decimal.Decimal `json:"prdy_vrss"`                      // 전일 대비
	PrdyVrssSign        string          `json:"prdy_vrss_sign"`                 // 전일 대비 부호
	PrdyCtrt            float64         `json:"prdy_ctrt,string"`               // 전일 대비율
	AcmlVol             int64           `json:"acml_vol,string"`                // 누적 거래량
	StndVolSmtn         int64           `json:"stnd_vol_smtn,string"`           // 기준 거래량 합계
	SstsCntgQty         int64           `json:"ssts_cntg_qty,string"`           // 공매도 체결 수량
	SstsVolRlim         float64         `json:"ssts_vol_rlim,string"`           // 공매도 거래량 비중
	AcmlSstsCntgQty     int64           `json:"acml_ssts_cntg_qty,string"`      // 누적 공매도 체결 수량
	AcmlSstsCntgQtyRlim float64         `json:"acml_ssts_cntg_qty_rlim,string"` // 누적 공매도 수량 비중
	AcmlTrPbmn          int64           `json:"acml_tr_pbmn,string"`            // 누적 거래 대금
	StndTrPbmnSmtn      int64           `json:"stnd_tr_pbmn_smtn,string"`       // 기준 거래 대금 합계
	SstsTrPbmn          int64           `json:"ssts_tr_pbmn,string"`            // 공매도 거래 대금
	SstsTrPbmnRlim      float64         `json:"ssts_tr_pbmn_rlim,string"`       // 공매도 거래 대금 비중
	AcmlSstsTrPbmn      int64           `json:"acml_ssts_tr_pbmn,string"`       // 누적 공매도 거래 대금
	AcmlSstsTrPbmnRlim  float64         `json:"acml_ssts_tr_pbmn_rlim,string"`  // 누적 공매도 대금 비중
	StckOprc            decimal.Decimal `json:"stck_oprc"`                      // 주식 시가
	StckHgpr            decimal.Decimal `json:"stck_hgpr"`                      // 주식 최고가
	StckLwpr            decimal.Decimal `json:"stck_lwpr"`                      // 주식 최저가
	AvrgPrc             decimal.Decimal `json:"avrg_prc"`                       // 평균가
}

// InquireDailyShortSaleParams 는 공매도 일별추이 조회 파라미터.
type InquireDailyShortSaleParams struct {
	InputDate2 string // FID_INPUT_DATE_2 — 조회 종료 일자 (YYYYMMDD)
	MarketCode string // FID_COND_MRKT_DIV_CODE — "J":KRX. 빈 값=>"J"
	Symbol     string // FID_INPUT_ISCD — 종목코드 (예 "005930")
	InputDate1 string // FID_INPUT_DATE_1 — 조회 시작 일자 (YYYYMMDD)
}

// InquireDailyShortSale 은 국내주식 공매도 일별추이 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_공매도_일별추이.md
// path: /uapi/domestic-stock/v1/quotations/daily-short-sale (FHPST04830000)
func (c *Client) InquireDailyShortSale(ctx context.Context, params InquireDailyShortSaleParams) (*DailyShortSale, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/daily-short-sale",
		TrID:   "FHPST04830000",
		Query: map[string]string{
			"FID_INPUT_DATE_2":       params.InputDate2,
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_INPUT_DATE_1":       params.InputDate1,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res DailyShortSale
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse DailyShortSale: %w", err)
	}
	return &res, nil
}

// CreditBalance 는 국내주식 신용잔고 상위 (FHKST17010000) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_신용잔고_상위.md
// path: /uapi/domestic-stock/v1/ranking/credit-balance
//
// output1: 날짜 헤더 배열, output2: 신용잔고 상위 종목 배열.
type CreditBalance struct {
	Output1 []CreditBalanceHeader `json:"output1"`
	Output2 []CreditBalanceItem   `json:"output2"`
}

// CreditBalanceHeader 는 신용잔고 상위 날짜 헤더 (output1 한 행).
type CreditBalanceHeader struct {
	BstpClsCode string `json:"bstp_cls_code"` // 업종 구분 코드
	HtsKorIsnm  string `json:"hts_kor_isnm"`  // HTS 한글 종목명
	StndDate1   string `json:"stnd_date1"`    // 기준 일자1
	StndDate2   string `json:"stnd_date2"`    // 기준 일자2
}

// CreditBalanceItem 은 신용잔고 상위 응답의 한 행 (output2).
type CreditBalanceItem struct {
	MkscShrnIscd         string          `json:"mksc_shrn_iscd"`                  // 유가증권 단축 종목코드
	HtsKorIsnm           string          `json:"hts_kor_isnm"`                    // HTS 한글 종목명
	StckPrpr             decimal.Decimal `json:"stck_prpr"`                       // 주식 현재가
	PrdyVrss             decimal.Decimal `json:"prdy_vrss"`                       // 전일 대비
	PrdyVrssSign         string          `json:"prdy_vrss_sign"`                  // 전일 대비 부호
	PrdyCtrt             float64         `json:"prdy_ctrt,string"`                // 전일 대비율
	AcmlVol              int64           `json:"acml_vol,string"`                 // 누적 거래량
	WholLoanRmndStcn     int64           `json:"whol_loan_rmnd_stcn,string"`      // 전체 융자 잔고 수량
	WholLoanRmndAmt      int64           `json:"whol_loan_rmnd_amt,string"`       // 전체 융자 잔고 금액
	WholLoanRmndRate     float64         `json:"whol_loan_rmnd_rate,string"`      // 전체 융자 잔고 비율
	WholStlnRmndStcn     int64           `json:"whol_stln_rmnd_stcn,string"`      // 전체 대주 잔고 수량
	WholStlnRmndAmt      int64           `json:"whol_stln_rmnd_amt,string"`       // 전체 대주 잔고 금액
	WholStlnRmndRate     float64         `json:"whol_stln_rmnd_rate,string"`      // 전체 대주 잔고 비율
	NdayVrssLoanRmndInrt float64         `json:"nday_vrss_loan_rmnd_inrt,string"` // N일 대비 융자 잔고 증가율
	NdayVrssStlnRmndInrt float64         `json:"nday_vrss_stln_rmnd_inrt,string"` // N일 대비 대주 잔고 증가율
}

// InquireCreditBalanceParams 는 신용잔고 상위 조회 파라미터.
type InquireCreditBalanceParams struct {
	CondScrDivCode string // FID_COND_SCR_DIV_CODE — 고정 "11701". 빈 값=>"11701"
	Symbol         string // FID_INPUT_ISCD — 종목코드 또는 시장코드
	Option         string // FID_OPTION — N일 윈도우 (2-999)
	MarketCode     string // FID_COND_MRKT_DIV_CODE — "J":KRX. 빈 값=>"J"
	RankSortCode   string // FID_RANK_SORT_CLS_CODE — 정렬 구분 코드
}

// InquireCreditBalance 는 국내주식 신용잔고 상위 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_신용잔고_상위.md
// path: /uapi/domestic-stock/v1/ranking/credit-balance (FHKST17010000)
func (c *Client) InquireCreditBalance(ctx context.Context, params InquireCreditBalanceParams) (*CreditBalance, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "11701"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ranking/credit-balance",
		TrID:   "FHKST17010000",
		Query: map[string]string{
			"FID_COND_SCR_DIV_CODE":  scrDiv,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_OPTION":             params.Option,
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_RANK_SORT_CLS_CODE": params.RankSortCode,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res CreditBalance
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse CreditBalance: %w", err)
	}
	return &res, nil
}

// DailyCreditBalance 는 국내주식 신용잔고 일별추이 (FHPST04760000) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_신용잔고_일별추이.md
// path: /uapi/domestic-stock/v1/quotations/daily-credit-balance
//
// 특정 종목의 일별 융자/대주 신규·상환·잔고 추이.
type DailyCreditBalance struct {
	Output []DailyCreditBalanceItem `json:"output"`
}

// DailyCreditBalanceItem 은 신용잔고 일별추이 응답의 한 행.
type DailyCreditBalanceItem struct {
	DealDate         string          `json:"deal_date"`                  // 결제 일자
	StckPrpr         decimal.Decimal `json:"stck_prpr"`                  // 주식 현재가
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`             // 전일 대비 부호
	PrdyVrss         decimal.Decimal `json:"prdy_vrss"`                  // 전일 대비
	PrdyCtrt         float64         `json:"prdy_ctrt,string"`           // 전일 대비율
	AcmlVol          int64           `json:"acml_vol,string"`            // 누적 거래량
	StlmDate         string          `json:"stlm_date"`                  // 결산 일자
	WholLoanNewStcn  int64           `json:"whol_loan_new_stcn,string"`  // 전체 융자 신규 수량
	WholLoanRdmpStcn int64           `json:"whol_loan_rdmp_stcn,string"` // 전체 융자 상환 수량
	WholLoanRmndStcn int64           `json:"whol_loan_rmnd_stcn,string"` // 전체 융자 잔고 수량
	WholLoanNewAmt   int64           `json:"whol_loan_new_amt,string"`   // 전체 융자 신규 금액
	WholLoanRdmpAmt  int64           `json:"whol_loan_rdmp_amt,string"`  // 전체 융자 상환 금액
	WholLoanRmndAmt  int64           `json:"whol_loan_rmnd_amt,string"`  // 전체 융자 잔고 금액
	WholLoanRmndRate float64         `json:"whol_loan_rmnd_rate,string"` // 전체 융자 잔고 비율
	WholLoanGvrt     float64         `json:"whol_loan_gvrt,string"`      // 전체 융자 담보비율
	WholStlnNewStcn  int64           `json:"whol_stln_new_stcn,string"`  // 전체 대주 신규 수량
	WholStlnRdmpStcn int64           `json:"whol_stln_rdmp_stcn,string"` // 전체 대주 상환 수량
	WholStlnRmndStcn int64           `json:"whol_stln_rmnd_stcn,string"` // 전체 대주 잔고 수량
	WholStlnNewAmt   int64           `json:"whol_stln_new_amt,string"`   // 전체 대주 신규 금액
	WholStlnRdmpAmt  int64           `json:"whol_stln_rdmp_amt,string"`  // 전체 대주 상환 금액
	WholStlnRmndAmt  int64           `json:"whol_stln_rmnd_amt,string"`  // 전체 대주 잔고 금액
	WholStlnRmndRate float64         `json:"whol_stln_rmnd_rate,string"` // 전체 대주 잔고 비율
	WholStlnGvrt     float64         `json:"whol_stln_gvrt,string"`      // 전체 대주 담보비율
	StckOprc         decimal.Decimal `json:"stck_oprc"`                  // 주식 시가
	StckHgpr         decimal.Decimal `json:"stck_hgpr"`                  // 주식 최고가
	StckLwpr         decimal.Decimal `json:"stck_lwpr"`                  // 주식 최저가
}

// InquireDailyCreditBalanceParams 는 신용잔고 일별추이 조회 파라미터.
//
// 쿼리 파라미터는 lowercase fid_* 형식으로 전송.
type InquireDailyCreditBalanceParams struct {
	MarketCode     string // fid_cond_mrkt_div_code — "J":KRX. 빈 값=>"J"
	CondScrDivCode string // fid_cond_scr_div_code — 고정 "20476". 빈 값=>"20476"
	Symbol         string // fid_input_iscd — 종목코드 (예 "005930")
	InputDate1     string // fid_input_date_1 — 조회 기준일 (YYYYMMDD)
}

// InquireDailyCreditBalance 는 국내주식 신용잔고 일별추이 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_신용잔고_일별추이.md
// path: /uapi/domestic-stock/v1/quotations/daily-credit-balance (FHPST04760000)
//
// 쿼리 파라미터는 lowercase fid_* 형식 사용.
func (c *Client) InquireDailyCreditBalance(ctx context.Context, params InquireDailyCreditBalanceParams) (*DailyCreditBalance, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "20476"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/daily-credit-balance",
		TrID:   "FHPST04760000",
		Query: map[string]string{
			"fid_cond_mrkt_div_code": market,
			"fid_cond_scr_div_code":  scrDiv,
			"fid_input_iscd":         params.Symbol,
			"fid_input_date_1":       params.InputDate1,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res DailyCreditBalance
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse DailyCreditBalance: %w", err)
	}
	return &res, nil
}

// LendableByCompany 는 당사 대주가능 종목 (CTSC2702R) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_당사대주가능종목.md
// path: /uapi/domestic-stock/v1/quotations/lendable-by-company
//
// output1: 대주가능 종목 배열, output2: 한도 합계 요약.
type LendableByCompany struct {
	Output1 []LendableByCompanyItem  `json:"output1"`
	Output2 LendableByCompanySummary `json:"output2"`
}

// LendableByCompanyItem 은 당사 대주가능 종목 응답의 한 행 (output1).
type LendableByCompanyItem struct {
	Pdno           string          `json:"pdno"`                  // 상품번호 (종목코드)
	PrdtName       string          `json:"prdt_name"`             // 상품명
	Papr           decimal.Decimal `json:"papr"`                  // 액면가
	BfdyClpr       decimal.Decimal `json:"bfdy_clpr"`             // 전일 종가
	SbstPrvs       decimal.Decimal `json:"sbst_prvs"`             // 대용가
	TrStopDvsnName string          `json:"tr_stop_dvsn_name"`     // 거래정지 구분 명
	PsblYnName     string          `json:"psbl_yn_name"`          // 가능 여부 명
	LmtQty1        int64           `json:"lmt_qty1,string"`       // 한도 수량1
	UseQty1        int64           `json:"use_qty1,string"`       // 사용 수량1
	TradPsblQty2   int64           `json:"trad_psbl_qty2,string"` // 거래 가능 수량2
	RghtTypeCd     string          `json:"rght_type_cd"`          // 권리 유형 코드
	BassDt         string          `json:"bass_dt"`               // 기준 일자
	PsblYn         string          `json:"psbl_yn"`               // 가능 여부
}

// LendableByCompanySummary 는 당사 대주가능 종목 한도 합계 요약 (output2).
type LendableByCompanySummary struct {
	TotStupLmtQty int64 `json:"tot_stup_lmt_qty,string"` // 총 설정 한도 수량
	BrchLmtQty    int64 `json:"brch_lmt_qty,string"`     // 지점 한도 수량
	RqstPsblQty   int64 `json:"rqst_psbl_qty,string"`    // 신청 가능 수량
}

// InquireLendableByCompanyParams 는 당사 대주가능 종목 조회 파라미터.
//
// 쿼리 파라미터는 non-FID UPPERCASE 형식 (EXCG_DVSN_CD, PDNO 등).
type InquireLendableByCompanyParams struct {
	ExcgDvsnCd     string // EXCG_DVSN_CD — 거래소 구분 코드
	Pdno           string // PDNO — 상품번호 (종목코드)
	ThcoStlnPsblYn string // THCO_STLN_PSBL_YN — 당사 대주 가능 여부
	InqrDvsn1      string // INQR_DVSN_1 — 조회 구분1
	CtxAreaFk200   string // CTX_AREA_FK200 — 연속조회 커서 (빈 값 가능)
	CtxAreaNk100   string // CTX_AREA_NK100 — 연속조회 키 (빈 값 가능)
}

// InquireLendableByCompany 는 당사 대주가능 종목 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_당사대주가능종목.md
// path: /uapi/domestic-stock/v1/quotations/lendable-by-company (CTSC2702R)
//
// 쿼리 파라미터는 non-FID UPPERCASE 형식 사용. CANO 불필요.
func (c *Client) InquireLendableByCompany(ctx context.Context, params InquireLendableByCompanyParams) (*LendableByCompany, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/lendable-by-company",
		TrID:   "CTSC2702R",
		Query: map[string]string{
			"EXCG_DVSN_CD":      params.ExcgDvsnCd,
			"PDNO":              params.Pdno,
			"THCO_STLN_PSBL_YN": params.ThcoStlnPsblYn,
			"INQR_DVSN_1":       params.InqrDvsn1,
			"CTX_AREA_FK200":    params.CtxAreaFk200,
			"CTX_AREA_NK100":    params.CtxAreaNk100,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res LendableByCompany
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse LendableByCompany: %w", err)
	}
	return &res, nil
}

// QuoteBalance 는 국내주식 호가잔량 순위 (FHPST01720000) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_호가잔량_순위.md
// path: /uapi/domestic-stock/v1/ranking/quote-balance
//
// 매도/매수 호가잔량 상위 종목 순위.
type QuoteBalance struct {
	Output []QuoteBalanceItem `json:"output"`
}

// QuoteBalanceItem 은 호가잔량 순위 응답의 한 행.
type QuoteBalanceItem struct {
	MkscShrnIscd      string          `json:"mksc_shrn_iscd"`              // 유가증권 단축 종목코드
	DataRank          string          `json:"data_rank"`                   // 데이터 순위
	HtsKorIsnm        string          `json:"hts_kor_isnm"`                // HTS 한글 종목명
	StckPrpr          decimal.Decimal `json:"stck_prpr"`                   // 주식 현재가
	PrdyVrss          decimal.Decimal `json:"prdy_vrss"`                   // 전일 대비
	PrdyVrssSign      string          `json:"prdy_vrss_sign"`              // 전일 대비 부호
	PrdyCtrt          float64         `json:"prdy_ctrt,string"`            // 전일 대비율
	AcmlVol           int64           `json:"acml_vol,string"`             // 누적 거래량
	TotalAskpRsqn     int64           `json:"total_askp_rsqn,string"`      // 총 매도호가 잔량
	TotalBidpRsqn     int64           `json:"total_bidp_rsqn,string"`      // 총 매수호가 잔량
	TotalNtslBidpRsqn int64           `json:"total_ntsl_bidp_rsqn,string"` // 총 순매수 매수호가 잔량
	ShnuRsqnRate      float64         `json:"shnu_rsqn_rate,string"`       // 매수 잔량 비율
	SelnRsqnRate      float64         `json:"seln_rsqn_rate,string"`       // 매도 잔량 비율
}

// InquireQuoteBalanceParams 는 호가잔량 순위 조회 파라미터.
//
// 쿼리 파라미터는 lowercase fid_* 형식으로 전송.
type InquireQuoteBalanceParams struct {
	VolCnt         string // fid_vol_cnt — 조회 건수
	MarketCode     string // fid_cond_mrkt_div_code — "J":KRX. 빈 값=>"J"
	CondScrDivCode string // fid_cond_scr_div_code — 고정 "20172". 빈 값=>"20172"
	Symbol         string // fid_input_iscd — 종목코드 또는 시장코드
	RankSortCode   string // fid_rank_sort_cls_code — 정렬 구분 코드
	DivClsCode     string // fid_div_cls_code — 구분 코드
	TrgtClsCode    string // fid_trgt_cls_code — 대상 구분 코드
	TrgtExlsCode   string // fid_trgt_exls_cls_code — 대상 제외 구분 코드
	InputPrice1    string // fid_input_price_1 — 입력 가격1
	InputPrice2    string // fid_input_price_2 — 입력 가격2
}

// InquireQuoteBalance 는 국내주식 호가잔량 순위 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_호가잔량_순위.md
// path: /uapi/domestic-stock/v1/ranking/quote-balance (FHPST01720000)
//
// 쿼리 파라미터는 lowercase fid_* 형식 사용.
func (c *Client) InquireQuoteBalance(ctx context.Context, params InquireQuoteBalanceParams) (*QuoteBalance, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "20172"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ranking/quote-balance",
		TrID:   "FHPST01720000",
		Query: map[string]string{
			"fid_vol_cnt":            params.VolCnt,
			"fid_cond_mrkt_div_code": market,
			"fid_cond_scr_div_code":  scrDiv,
			"fid_input_iscd":         params.Symbol,
			"fid_rank_sort_cls_code": params.RankSortCode,
			"fid_div_cls_code":       params.DivClsCode,
			"fid_trgt_cls_code":      params.TrgtClsCode,
			"fid_trgt_exls_cls_code": params.TrgtExlsCode,
			"fid_input_price_1":      params.InputPrice1,
			"fid_input_price_2":      params.InputPrice2,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res QuoteBalance
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse QuoteBalance: %w", err)
	}
	return &res, nil
}

// AfterHourBalance 는 국내주식 시간외잔량 순위 (FHPST01760000) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_시간외잔량_순위.md
// path: /uapi/domestic-stock/v1/ranking/after-hour-balance
//
// 시간외잔량 상위 종목 순위. 쿼리 파라미터는 lowercase fid_* 형식 사용.
type AfterHourBalance struct {
	Output []AfterHourBalanceItem `json:"output"`
}

// AfterHourBalanceItem 은 시간외잔량 순위 응답의 한 행.
type AfterHourBalanceItem struct {
	StckShrnIscd      string          `json:"stck_shrn_iscd"`              // 주식 단축 종목코드 (NOT mksc_shrn_iscd)
	DataRank          string          `json:"data_rank"`                   // 데이터 순위
	HtsKorIsnm        string          `json:"hts_kor_isnm"`                // HTS 한글 종목명
	StckPrpr          decimal.Decimal `json:"stck_prpr"`                   // 주식 현재가
	PrdyVrss          decimal.Decimal `json:"prdy_vrss"`                   // 전일 대비
	PrdyVrssSign      string          `json:"prdy_vrss_sign"`              // 전일 대비 부호
	PrdyCtrt          float64         `json:"prdy_ctrt,string"`            // 전일 대비율
	OvtmTotalAskpRsqn int64           `json:"ovtm_total_askp_rsqn,string"` // 시간외 총 매도호가 잔량
	OvtmTotalBidpRsqn int64           `json:"ovtm_total_bidp_rsqn,string"` // 시간외 총 매수호가 잔량
	MkobOtcpVol       int64           `json:"mkob_otcp_vol,string"`        // 장개시전 시간외종가 거래량
	MkfaOtcpVol       int64           `json:"mkfa_otcp_vol,string"`        // 장종료후 시간외종가 거래량
}

// InquireAfterHourBalanceParams 는 시간외잔량 순위 조회 파라미터.
//
// 쿼리 파라미터는 lowercase fid_* 형식으로 전송.
type InquireAfterHourBalanceParams struct {
	InputPrice1    string // fid_input_price_1 — 입력 가격1
	MarketCode     string // fid_cond_mrkt_div_code — "J":KRX. 빈 값=>"J"
	CondScrDivCode string // fid_cond_scr_div_code — 고정 "20176". 빈 값=>"20176"
	RankSortCode   string // fid_rank_sort_cls_code — 정렬 구분 코드
	DivClsCode     string // fid_div_cls_code — 구분 코드
	Symbol         string // fid_input_iscd — 종목코드 또는 시장코드
	TrgtExlsCode   string // fid_trgt_exls_cls_code — 대상 제외 구분 코드
	TrgtClsCode    string // fid_trgt_cls_code — 대상 구분 코드
	VolCnt         string // fid_vol_cnt — 조회 건수
	InputPrice2    string // fid_input_price_2 — 입력 가격2
}

// InquireAfterHourBalance 는 국내주식 시간외잔량 순위 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_시간외잔량_순위.md
// path: /uapi/domestic-stock/v1/ranking/after-hour-balance (FHPST01760000)
//
// 쿼리 파라미터는 lowercase fid_* 형식 사용. StckShrnIscd (NOT mksc_shrn_iscd) 주의.
func (c *Client) InquireAfterHourBalance(ctx context.Context, params InquireAfterHourBalanceParams) (*AfterHourBalance, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "20176"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ranking/after-hour-balance",
		TrID:   "FHPST01760000",
		Query: map[string]string{
			"fid_input_price_1":      params.InputPrice1,
			"fid_cond_mrkt_div_code": market,
			"fid_cond_scr_div_code":  scrDiv,
			"fid_rank_sort_cls_code": params.RankSortCode,
			"fid_div_cls_code":       params.DivClsCode,
			"fid_input_iscd":         params.Symbol,
			"fid_trgt_exls_cls_code": params.TrgtExlsCode,
			"fid_trgt_cls_code":      params.TrgtClsCode,
			"fid_vol_cnt":            params.VolCnt,
			"fid_input_price_2":      params.InputPrice2,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res AfterHourBalance
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse AfterHourBalance: %w", err)
	}
	return &res, nil
}

// OvertimeExpTransFluct 는 국내주식 시간외 예상체결 등락률 (FHKST11860000) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_시간외_예상체결_등락률.md
// path: /uapi/domestic-stock/v1/ranking/overtime-exp-trans-fluct
//
// output 은 배열이 아닌 단일 객체. 쿼리 파라미터는 UPPERCASE FID_ 형식 사용.
// InquireOvertimeFluctuation (Phase 2.2) 과 다른 별도 API.
type OvertimeExpTransFluct struct {
	Output OvertimeExpTransFluctData `json:"output"`
}

// OvertimeExpTransFluctData 는 시간외 예상체결 등락률 응답의 output 단일 객체.
//
// NOTE: ovtm_untp_antc_cntg_vrsssign — KIS docs 오타 보존 (vrss+sign 연결, 밑줄 없음).
type OvertimeExpTransFluctData struct {
	DataRank                 string          `json:"data_rank"`                       // 데이터 순위
	IscdStatClsCode          string          `json:"iscd_stat_cls_code"`              // 종목 상태 구분 코드
	StckShrnIscd             string          `json:"stck_shrn_iscd"`                  // 주식 단축 종목코드
	HtsKorIsnm               string          `json:"hts_kor_isnm"`                    // HTS 한글 종목명
	OvtmUntpAntcCnpr         decimal.Decimal `json:"ovtm_untp_antc_cnpr"`             // 시간외 단일가 예상 체결가
	OvtmUntpAntcCntgVrss     decimal.Decimal `json:"ovtm_untp_antc_cntg_vrss"`        // 시간외 단일가 예상 체결 대비
	OvtmUntpAntcCntgVrsssign string          `json:"ovtm_untp_antc_cntg_vrsssign"`    // 시간외 단일가 예상 체결 대비 부호 (KIS 오타: vrsssign)
	OvtmUntpAntcCntgCtrt     float64         `json:"ovtm_untp_antc_cntg_ctrt,string"` // 시간외 단일가 예상 체결 대비율
	OvtmUntpAskpRsqn1        int64           `json:"ovtm_untp_askp_rsqn1,string"`     // 시간외 단일가 매도호가 잔량1
	OvtmUntpBidpRsqn1        int64           `json:"ovtm_untp_bidp_rsqn1,string"`     // 시간외 단일가 매수호가 잔량1
	OvtmUntpAntcCnqn         int64           `json:"ovtm_untp_antc_cnqn,string"`      // 시간외 단일가 예상 체결량
	ItmtVol                  int64           `json:"itmt_vol,string"`                 // 중간 거래량
	StckPrpr                 decimal.Decimal `json:"stck_prpr"`                       // 주식 현재가
}

// InquireOvertimeExpTransFluctParams 는 시간외 예상체결 등락률 조회 파라미터.
//
// 쿼리 파라미터는 UPPERCASE FID_ 형식으로 전송.
type InquireOvertimeExpTransFluctParams struct {
	MarketCode     string // FID_COND_MRKT_DIV_CODE — "J":KRX. 빈 값=>"J"
	CondScrDivCode string // FID_COND_SCR_DIV_CODE — 고정 "11186". 빈 값=>"11186"
	Symbol         string // FID_INPUT_ISCD — 종목코드 또는 시장코드
	RankSortCode   string // FID_RANK_SORT_CLS_CODE — 정렬 구분 코드
	DivClsCode     string // FID_DIV_CLS_CODE — 구분 코드
	InputPrice1    string // FID_INPUT_PRICE_1 — 입력 가격1
	InputPrice2    string // FID_INPUT_PRICE_2 — 입력 가격2
	InputVol1      string // FID_INPUT_VOL_1 — 입력 거래량1
}

// InquireOvertimeExpTransFluct 는 국내주식 시간외 예상체결 등락률 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_시간외_예상체결_등락률.md
// path: /uapi/domestic-stock/v1/ranking/overtime-exp-trans-fluct (FHKST11860000)
//
// output 은 단일 객체 (배열 아님). 쿼리 파라미터는 UPPERCASE FID_ 형식 사용.
func (c *Client) InquireOvertimeExpTransFluct(ctx context.Context, params InquireOvertimeExpTransFluctParams) (*OvertimeExpTransFluct, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "11186"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ranking/overtime-exp-trans-fluct",
		TrID:   "FHKST11860000",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_COND_SCR_DIV_CODE":  scrDiv,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_RANK_SORT_CLS_CODE": params.RankSortCode,
			"FID_DIV_CLS_CODE":       params.DivClsCode,
			"FID_INPUT_PRICE_1":      params.InputPrice1,
			"FID_INPUT_PRICE_2":      params.InputPrice2,
			"FID_INPUT_VOL_1":        params.InputVol1,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res OvertimeExpTransFluct
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse OvertimeExpTransFluct: %w", err)
	}
	return &res, nil
}

// MarketValue 는 국내주식 시장가치 순위 (FHPST01790000) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_시장가치_순위.md
// path: /uapi/domestic-stock/v1/ranking/market-value
//
// PER/PBR/PCR/PSR/EPS/EVA/EBITDA 등 가치지표 기반 순위.
// 쿼리 파라미터는 lowercase fid_* 형식 사용.
type MarketValue struct {
	Output []MarketValueItem `json:"output"`
}

// MarketValueItem 은 시장가치 순위 응답의 한 행.
type MarketValueItem struct {
	DataRank          string          `json:"data_rank"`                   // 데이터 순위
	HtsKorIsnm        string          `json:"hts_kor_isnm"`                // HTS 한글 종목명
	MkscShrnIscd      string          `json:"mksc_shrn_iscd"`              // 유가증권 단축 종목코드
	StckPrpr          decimal.Decimal `json:"stck_prpr"`                   // 주식 현재가
	PrdyVrss          decimal.Decimal `json:"prdy_vrss"`                   // 전일 대비
	PrdyVrssSign      string          `json:"prdy_vrss_sign"`              // 전일 대비 부호
	PrdyCtrt          float64         `json:"prdy_ctrt,string"`            // 전일 대비율
	AcmlVol           int64           `json:"acml_vol,string"`             // 누적 거래량
	Per               float64         `json:"per,string"`                  // PER
	Pbr               float64         `json:"pbr,string"`                  // PBR
	Pcr               float64         `json:"pcr,string"`                  // PCR
	Psr               float64         `json:"psr,string"`                  // PSR
	Eps               float64         `json:"eps,string"`                  // EPS
	Eva               float64         `json:"eva,string"`                  // EVA
	Ebitda            float64         `json:"ebitda,string"`               // EBITDA
	PvDivEbitda       float64         `json:"pv_div_ebitda,string"`        // EV/EBITDA
	EbitdaDivFnncExpn float64         `json:"ebitda_div_fnnc_expn,string"` // EBITDA/금융비용
	StacMonth         string          `json:"stac_month"`                  // 결산 월
	StacMonthClsCode  string          `json:"stac_month_cls_code"`         // 결산 월 구분 코드
	IqryCsnu          string          `json:"iqry_csnu"`                   // 조회 건수
}

// InquireMarketValueParams 는 시장가치 순위 조회 파라미터.
//
// 쿼리 파라미터는 lowercase fid_* 형식으로 전송.
type InquireMarketValueParams struct {
	TrgtClsCode    string // fid_trgt_cls_code — 대상 구분 코드
	MarketCode     string // fid_cond_mrkt_div_code — "J":KRX. 빈 값=>"J"
	CondScrDivCode string // fid_cond_scr_div_code — 고정 "20179". 빈 값=>"20179"
	Symbol         string // fid_input_iscd — 종목코드 또는 시장코드
	DivClsCode     string // fid_div_cls_code — 구분 코드
	InputPrice1    string // fid_input_price_1 — 입력 가격1
	InputPrice2    string // fid_input_price_2 — 입력 가격2
	VolCnt         string // fid_vol_cnt — 조회 건수
	InputOption1   string // fid_input_option_1 — 회계연도
	InputOption2   string // fid_input_option_2 — 분기구분
	RankSortCode   string // fid_rank_sort_cls_code — 정렬 구분 코드
	BlngClsCode    string // fid_blng_cls_code — 소속 구분 코드
	TrgtExlsCode   string // fid_trgt_exls_cls_code — 대상 제외 구분 코드
}

// InquireMarketValue 는 국내주식 시장가치 순위 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_시장가치_순위.md
// path: /uapi/domestic-stock/v1/ranking/market-value (FHPST01790000)
//
// 쿼리 파라미터는 lowercase fid_* 형식 사용.
func (c *Client) InquireMarketValue(ctx context.Context, params InquireMarketValueParams) (*MarketValue, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "20179"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ranking/market-value",
		TrID:   "FHPST01790000",
		Query: map[string]string{
			"fid_trgt_cls_code":      params.TrgtClsCode,
			"fid_cond_mrkt_div_code": market,
			"fid_cond_scr_div_code":  scrDiv,
			"fid_input_iscd":         params.Symbol,
			"fid_div_cls_code":       params.DivClsCode,
			"fid_input_price_1":      params.InputPrice1,
			"fid_input_price_2":      params.InputPrice2,
			"fid_vol_cnt":            params.VolCnt,
			"fid_input_option_1":     params.InputOption1,
			"fid_input_option_2":     params.InputOption2,
			"fid_rank_sort_cls_code": params.RankSortCode,
			"fid_blng_cls_code":      params.BlngClsCode,
			"fid_trgt_exls_cls_code": params.TrgtExlsCode,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res MarketValue
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse MarketValue: %w", err)
	}
	return &res, nil
}

// Disparity 는 국내주식 이격도 순위 (FHPST01780000) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_이격도_순위.md
// path: /uapi/domestic-stock/v1/ranking/disparity
//
// 이동평균선(5/10/20/60/120일)과 현재가의 이격도 순위.
type Disparity struct {
	Output []DisparityItem `json:"output"`
}

// DisparityItem 은 이격도 순위 응답의 한 행.
type DisparityItem struct {
	MkscShrnIscd string          `json:"mksc_shrn_iscd"`   // 유가증권 단축 종목코드
	DataRank     string          `json:"data_rank"`        // 데이터 순위
	HtsKorIsnm   string          `json:"hts_kor_isnm"`     // HTS 한글 종목명
	StckPrpr     decimal.Decimal `json:"stck_prpr"`        // 주식 현재가
	PrdyVrss     decimal.Decimal `json:"prdy_vrss"`        // 전일 대비
	PrdyCtrt     float64         `json:"prdy_ctrt,string"` // 전일 대비율
	PrdyVrssSign string          `json:"prdy_vrss_sign"`   // 전일 대비 부호
	AcmlVol      int64           `json:"acml_vol,string"`  // 누적 거래량
	D5Dsrt       float64         `json:"d5_dsrt,string"`   // 5일 이격도
	D10Dsrt      float64         `json:"d10_dsrt,string"`  // 10일 이격도
	D20Dsrt      float64         `json:"d20_dsrt,string"`  // 20일 이격도
	D60Dsrt      float64         `json:"d60_dsrt,string"`  // 60일 이격도
	D120Dsrt     float64         `json:"d120_dsrt,string"` // 120일 이격도
}

// InquireDisparityParams 는 이격도 순위 조회 파라미터.
//
// FID_COND_SCR_DIV_CODE = "20178" 고정.
// HourClsCode: 5/10/20/60/120 (일수).
type InquireDisparityParams struct {
	InputPrice2    string // fid_input_price_2 — 입력 가격2
	MarketCode     string // fid_cond_mrkt_div_code — "J":KRX. 빈 값=>"J"
	CondScrDivCode string // fid_cond_scr_div_code — 빈 값=>"20178"
	DivClsCode     string // fid_div_cls_code — 구분 코드
	RankSortCode   string // fid_rank_sort_cls_code — 정렬 구분 코드
	HourClsCode    string // fid_hour_cls_code — 5/10/20/60/120일
	Symbol         string // fid_input_iscd — 종목코드 (0000:전체)
	TrgtClsCode    string // fid_trgt_cls_code — 대상 구분 코드
	TrgtExlsCode   string // fid_trgt_exls_cls_code — 대상 제외 구분 코드
	InputPrice1    string // fid_input_price_1 — 입력 가격1
	VolCnt         string // fid_vol_cnt — 조회 건수
}

// InquireDisparity 는 국내주식 이격도 순위 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_이격도_순위.md
// path: /uapi/domestic-stock/v1/ranking/disparity (FHPST01780000)
//
// 쿼리 파라미터는 lowercase fid_* 형식 사용.
// HourClsCode: 5/10/20/60/120 (이동평균일수).
func (c *Client) InquireDisparity(ctx context.Context, params InquireDisparityParams) (*Disparity, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "20178"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ranking/disparity",
		TrID:   "FHPST01780000",
		Query: map[string]string{
			"fid_input_price_2":      params.InputPrice2,
			"fid_cond_mrkt_div_code": market,
			"fid_cond_scr_div_code":  scrDiv,
			"fid_div_cls_code":       params.DivClsCode,
			"fid_rank_sort_cls_code": params.RankSortCode,
			"fid_hour_cls_code":      params.HourClsCode,
			"fid_input_iscd":         params.Symbol,
			"fid_trgt_cls_code":      params.TrgtClsCode,
			"fid_trgt_exls_cls_code": params.TrgtExlsCode,
			"fid_input_price_1":      params.InputPrice1,
			"fid_vol_cnt":            params.VolCnt,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res Disparity
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse Disparity: %w", err)
	}
	return &res, nil
}

// PreferDisparateRatio 는 국내주식 우선주 괴리율 상위 (FHPST01770000) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_우선주_괴리율_상위.md
// path: /uapi/domestic-stock/v1/ranking/prefer-disparate-ratio
//
// 보통주-우선주 간 괴리율 순위.
type PreferDisparateRatio struct {
	Output []PreferDisparateRatioItem `json:"output"`
}

// PreferDisparateRatioItem 은 우선주 괴리율 상위 응답의 한 행.
type PreferDisparateRatioItem struct {
	MkscShrnIscd     string          `json:"mksc_shrn_iscd"`        // 유가증권 단축 종목코드 (보통주)
	DataRank         string          `json:"data_rank"`             // 데이터 순위
	HtsKorIsnm       string          `json:"hts_kor_isnm"`          // HTS 한글 종목명 (보통주)
	StckPrpr         decimal.Decimal `json:"stck_prpr"`             // 주식 현재가 (보통주)
	PrdyVrss         decimal.Decimal `json:"prdy_vrss"`             // 전일 대비 (보통주)
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`        // 전일 대비 부호 (보통주)
	AcmlVol          int64           `json:"acml_vol,string"`       // 누적 거래량 (보통주)
	PrstIscd         string          `json:"prst_iscd"`             // 우선주 종목코드
	PrstKorIsnm      string          `json:"prst_kor_isnm"`         // 우선주 한글 종목명
	PrstPrpr         decimal.Decimal `json:"prst_prpr"`             // 우선주 현재가
	PrstPrdyVrss     decimal.Decimal `json:"prst_prdy_vrss"`        // 우선주 전일 대비
	PrstPrdyVrssSign string          `json:"prst_prdy_vrss_sign"`   // 우선주 전일 대비 부호
	PrstAcmlVol      int64           `json:"prst_acml_vol,string"`  // 우선주 누적 거래량
	DiffPrpr         decimal.Decimal `json:"diff_prpr"`             // 차이 현재가
	Dprt             float64         `json:"dprt,string"`           // 괴리율
	PrdyCtrt         float64         `json:"prdy_ctrt,string"`      // 보통주 전일 대비율
	PrstPrdyCtrt     float64         `json:"prst_prdy_ctrt,string"` // 우선주 전일 대비율
}

// InquirePreferDisparateRatioParams 는 우선주 괴리율 상위 조회 파라미터.
//
// FID_COND_SCR_DIV_CODE = "20177" 고정.
type InquirePreferDisparateRatioParams struct {
	VolCnt         string // fid_vol_cnt — 조회 건수
	MarketCode     string // fid_cond_mrkt_div_code — "J":KRX. 빈 값=>"J"
	CondScrDivCode string // fid_cond_scr_div_code — 빈 값=>"20177"
	DivClsCode     string // fid_div_cls_code — 구분 코드
	Symbol         string // fid_input_iscd — 종목코드 (0000:전체)
	TrgtClsCode    string // fid_trgt_cls_code — 대상 구분 코드
	TrgtExlsCode   string // fid_trgt_exls_cls_code — 대상 제외 구분 코드
	InputPrice1    string // fid_input_price_1 — 입력 가격1
	InputPrice2    string // fid_input_price_2 — 입력 가격2
}

// InquirePreferDisparateRatio 는 국내주식 우선주 괴리율 상위 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_우선주_괴리율_상위.md
// path: /uapi/domestic-stock/v1/ranking/prefer-disparate-ratio (FHPST01770000)
//
// 쿼리 파라미터는 lowercase fid_* 형식 사용.
func (c *Client) InquirePreferDisparateRatio(ctx context.Context, params InquirePreferDisparateRatioParams) (*PreferDisparateRatio, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scrDiv := params.CondScrDivCode
	if scrDiv == "" {
		scrDiv = "20177"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ranking/prefer-disparate-ratio",
		TrID:   "FHPST01770000",
		Query: map[string]string{
			"fid_vol_cnt":            params.VolCnt,
			"fid_cond_mrkt_div_code": market,
			"fid_cond_scr_div_code":  scrDiv,
			"fid_div_cls_code":       params.DivClsCode,
			"fid_input_iscd":         params.Symbol,
			"fid_trgt_cls_code":      params.TrgtClsCode,
			"fid_trgt_exls_cls_code": params.TrgtExlsCode,
			"fid_input_price_1":      params.InputPrice1,
			"fid_input_price_2":      params.InputPrice2,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res PreferDisparateRatio
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse PreferDisparateRatio: %w", err)
	}
	return &res, nil
}
