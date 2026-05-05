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
