package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// Price 는 주식현재가_시세 (FHKST01010100) 의 output 응답.
//
// 한투 docs: docs/api/국내주식/주식현재가_시세.md
// path: /uapi/domestic-stock/v1/quotations/inquire-price
type Price struct {
	IscdStatClsCode      string          `json:"iscd_stat_cls_code"`              // 종목 상태 구분 코드
	MargRate             float64         `json:"marg_rate,string"`                // 증거금 비율
	RprsMrktKorName      string          `json:"rprs_mrkt_kor_name"`              // 대표 시장 한글명
	NewHgprLwprClsCode   string          `json:"new_hgpr_lwpr_cls_code"`          // 신 고가/저가 구분 코드
	BstpKorIsnm          string          `json:"bstp_kor_isnm"`                   // 업종 한글 종목명
	TempStopYn           string          `json:"temp_stop_yn"`                    // 임시 정지 여부
	OprcRangContYn       string          `json:"oprc_rang_cont_yn"`               // 시가 범위 연장 여부
	ClprRangContYn       string          `json:"clpr_rang_cont_yn"`               // 종가 범위 연장 여부
	CrdtAbleYn           string          `json:"crdt_able_yn"`                    // 신용 가능 여부
	GrmnRateClsCode      string          `json:"grmn_rate_cls_code"`              // 보증금 비율 구분 코드
	ElwPblcYn            string          `json:"elw_pblc_yn"`                     // ELW 발행 여부
	StckPrpr             decimal.Decimal `json:"stck_prpr"`                       // 주식 현재가
	PrdyVrss             decimal.Decimal `json:"prdy_vrss"`                       // 전일 대비
	PrdyVrssSign         string          `json:"prdy_vrss_sign"`                  // 전일 대비 부호 (1:상한 2:상승 3:보합 4:하락 5:하한)
	PrdyCtrt             float64         `json:"prdy_ctrt,string"`                // 전일 대비율
	AcmlTrPbmn           int64           `json:"acml_tr_pbmn,string"`             // 누적 거래 대금
	AcmlVol              int64           `json:"acml_vol,string"`                 // 누적 거래량
	PrdyVrssVolRate      float64         `json:"prdy_vrss_vol_rate,string"`       // 전일 대비 거래량 비율
	StckOprc             decimal.Decimal `json:"stck_oprc"`                       // 주식 시가
	StckHgpr             decimal.Decimal `json:"stck_hgpr"`                       // 주식 최고가
	StckLwpr             decimal.Decimal `json:"stck_lwpr"`                       // 주식 최저가
	StckMxpr             decimal.Decimal `json:"stck_mxpr"`                       // 주식 상한가
	StckLlam             decimal.Decimal `json:"stck_llam"`                       // 주식 하한가
	StckSdpr             decimal.Decimal `json:"stck_sdpr"`                       // 주식 기준가
	WghnAvrgStckPrc      decimal.Decimal `json:"wghn_avrg_stck_prc"`              // 가중 평균 주식 가격
	HtsFrgnEhrt          float64         `json:"hts_frgn_ehrt,string"`            // HTS 외국인 소진율
	FrgnNtbyQty          int64           `json:"frgn_ntby_qty,string"`            // 외국인 순매수 수량
	PgtrNtbyQty          int64           `json:"pgtr_ntby_qty,string"`            // 프로그램매매 순매수 수량
	PvtScndDmrsPrc       decimal.Decimal `json:"pvt_scnd_dmrs_prc"`               // 피벗 2차 디저항 가격
	PvtFrstDmrsPrc       decimal.Decimal `json:"pvt_frst_dmrs_prc"`               // 피벗 1차 디저항 가격
	PvtPontVal           decimal.Decimal `json:"pvt_pont_val"`                    // 피벗 포인트 값
	PvtFrstDmspPrc       decimal.Decimal `json:"pvt_frst_dmsp_prc"`               // 피벗 1차 디지지 가격
	PvtScndDmspPrc       decimal.Decimal `json:"pvt_scnd_dmsp_prc"`               // 피벗 2차 디지지 가격
	DmrsVal              decimal.Decimal `json:"dmrs_val"`                        // 디저항 값
	DmspVal              decimal.Decimal `json:"dmsp_val"`                        // 디지지 값
	Cpfn                 int64           `json:"cpfn,string"`                     // 자본금
	RstcWdthPrc          decimal.Decimal `json:"rstc_wdth_prc"`                   // 제한 폭 가격
	StckFcam             decimal.Decimal `json:"stck_fcam"`                       // 주식 액면가
	StckSspr             decimal.Decimal `json:"stck_sspr"`                       // 주식 대용가
	AsprUnit             decimal.Decimal `json:"aspr_unit"`                       // 호가 단위
	HtsDealQtyUnitVal    int64           `json:"hts_deal_qty_unit_val,string"`    // HTS 매매 수량 단위 값
	LstnStcn             int64           `json:"lstn_stcn,string"`                // 상장 주수
	HtsAvls              int64           `json:"hts_avls,string"`                 // HTS 시가총액 (억원)
	Per                  float64         `json:"per,string"`                      // PER
	Pbr                  float64         `json:"pbr,string"`                      // PBR
	StacMonth            string          `json:"stac_month"`                      // 결산 월
	VolTnrt              float64         `json:"vol_tnrt,string"`                 // 거래량 회전율
	Eps                  decimal.Decimal `json:"eps"`                             // EPS
	Bps                  decimal.Decimal `json:"bps"`                             // BPS
	D250Hgpr             decimal.Decimal `json:"d250_hgpr"`                       // 250일 최고가
	D250HgprDate         string          `json:"d250_hgpr_date"`                  // 250일 최고가 일자
	D250HgprVrssPrprRate float64         `json:"d250_hgpr_vrss_prpr_rate,string"` // 250일 최고가 대비 현재가 비율
	D250Lwpr             decimal.Decimal `json:"d250_lwpr"`                       // 250일 최저가
	D250LwprDate         string          `json:"d250_lwpr_date"`                  // 250일 최저가 일자
	D250LwprVrssPrprRate float64         `json:"d250_lwpr_vrss_prpr_rate,string"` // 250일 최저가 대비 현재가 비율
	StckDryyHgpr         decimal.Decimal `json:"stck_dryy_hgpr"`                  // 주식 연중 최고가
	DryyHgprVrssPrprRate float64         `json:"dryy_hgpr_vrss_prpr_rate,string"` // 연중 최고가 대비 현재가 비율
	DryyHgprDate         string          `json:"dryy_hgpr_date"`                  // 연중 최고가 일자
	StckDryyLwpr         decimal.Decimal `json:"stck_dryy_lwpr"`                  // 주식 연중 최저가
	DryyLwprVrssPrprRate float64         `json:"dryy_lwpr_vrss_prpr_rate,string"` // 연중 최저가 대비 현재가 비율
	DryyLwprDate         string          `json:"dryy_lwpr_date"`                  // 연중 최저가 일자
	W52Hgpr              decimal.Decimal `json:"w52_hgpr"`                        // 52주 최고가
	W52HgprVrssPrprCtrt  float64         `json:"w52_hgpr_vrss_prpr_ctrt,string"`  // 52주 최고가 대비 현재가 대비
	W52HgprDate          string          `json:"w52_hgpr_date"`                   // 52주 최고가 일자
	W52Lwpr              decimal.Decimal `json:"w52_lwpr"`                        // 52주 최저가
	W52LwprVrssPrprCtrt  float64         `json:"w52_lwpr_vrss_prpr_ctrt,string"`  // 52주 최저가 대비 현재가 대비
	W52LwprDate          string          `json:"w52_lwpr_date"`                   // 52주 최저가 일자
	WholLoanRmndRate     float64         `json:"whol_loan_rmnd_rate,string"`      // 전체 융자 잔고 비율
	SstsYn               string          `json:"ssts_yn"`                         // 공매도 가능 여부
	StckShrnIscd         string          `json:"stck_shrn_iscd"`                  // 주식 단축 종목코드
	FcamCnnm             string          `json:"fcam_cnnm"`                       // 액면가 통화명
	CpfnCnnm             string          `json:"cpfn_cnnm"`                       // 자본금 통화명
	ApprchRate           float64         `json:"apprch_rate,string"`              // 접근도
	FrgnHldnQty          int64           `json:"frgn_hldn_qty,string"`            // 외국인 보유 수량
	ViClsCode            string          `json:"vi_cls_code"`                     // VI 구분 코드
	OvtmViClsCode        string          `json:"ovtm_vi_cls_code"`                // 시간외 VI 구분 코드
	LastSstsCntgQty      int64           `json:"last_ssts_cntg_qty,string"`       // 최종 공매도 체결 수량
	InvtCafulYn          string          `json:"invt_caful_yn"`                   // 투자 유의 여부
	MrktWarnClsCode      string          `json:"mrkt_warn_cls_code"`              // 시장 경고 구분 코드
	ShortOverYn          string          `json:"short_over_yn"`                   // 단기 과열 여부
	SltrYn               string          `json:"sltr_yn"`                         // 정리매매 여부
	MangIssuClsCode      string          `json:"mang_issu_cls_code"`              // 관리 종목 구분 코드
}

// InquirePrice 는 주식현재가 시세 호출.
//
// 한투 docs: docs/api/국내주식/주식현재가_시세.md
// path: /uapi/domestic-stock/v1/quotations/inquire-price (FHKST01010100)
//
// FID_COND_MRKT_DIV_CODE 는 "J" (KRX) 고정. NXT/통합 시장은 별도 메서드 후일 추가 검토.
func (c *Client) InquirePrice(ctx context.Context, symbol string) (*Price, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-price",
		TrID:   "FHKST01010100",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": "J",
			"FID_INPUT_ISCD":         symbol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var p Price
	if err := json.Unmarshal(resp.Output, &p); err != nil {
		return nil, fmt.Errorf("kis: parse Price: %w", err)
	}
	return &p, nil
}
