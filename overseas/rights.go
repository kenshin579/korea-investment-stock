package overseas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// RightsByIce 는 해외주식_권리종합 (HHDFS78330900) 응답.
//
// 한투 docs: docs/api/해외주식/해외주식_권리종합.md
// path: /uapi/overseas-price/v1/quotations/rights-by-ice
//
// ANOMALY: output1 만 존재 (output2 없음).
type RightsByIce struct {
	Output1 []RightsByIceItem `json:"output1"`
}

// RightsByIceItem 은 해외주식_권리종합 한 행. 모든 필드 string (KIS docs).
type RightsByIceItem struct {
	AnnoDt         string `json:"anno_dt"`          // 공시일자
	CaTitle        string `json:"ca_title"`         // 권리종류명
	DivLockDt      string `json:"div_lock_dt"`      // 배당락일
	PayDt          string `json:"pay_dt"`           // 지급일
	RecordDt       string `json:"record_dt"`        // 기준일
	ValidityDt     string `json:"validity_dt"`      // 유효일
	LocalEndDt     string `json:"local_end_dt"`     // 현지종료일
	LockDt         string `json:"lock_dt"`          // 권리확정일
	DelistDt       string `json:"delist_dt"`        // 상장폐지일
	RedemptDt      string `json:"redempt_dt"`       // 상환일
	EarlyRedemptDt string `json:"early_redempt_dt"` // 조기상환일
	EffectiveDt    string `json:"effective_dt"`     // 효력발생일
}

// InquireRightsByIceParams 는 해외주식_권리종합 조회 파라미터.
type InquireRightsByIceParams struct {
	NCod  string // NCOD — 국가코드. 빈 값 default
	Symb  string // SYMB — 종목코드. 빈 값 default (전체)
	StYmd string // ST_YMD — 조회시작일 YYYYMMDD. 빈 값 default
	EdYmd string // ED_YMD — 조회종료일 YYYYMMDD. 빈 값 default
}

// InquireRightsByIce 는 해외주식_권리종합 호출.
//
// 한투 docs: docs/api/해외주식/해외주식_권리종합.md
// path: /uapi/overseas-price/v1/quotations/rights-by-ice (HHDFS78330900)
func (c *Client) InquireRightsByIce(ctx context.Context, params InquireRightsByIceParams) (*RightsByIce, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-price/v1/quotations/rights-by-ice",
		TrID:   "HHDFS78330900",
		Query: map[string]string{
			"NCOD":   params.NCod,
			"SYMB":   params.Symb,
			"ST_YMD": params.StYmd,
			"ED_YMD": params.EdYmd,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res RightsByIce
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse RightsByIce: %w", err)
	}
	return &res, nil
}

// PeriodRights 는 해외주식_기간별권리조회 (CTRGT011R) 응답.
//
// 한투 docs: docs/api/해외주식/해외주식_기간별권리조회.md
// path: /uapi/overseas-price/v1/quotations/period-rights
//
// ANOMALY 1: TR_ID = CTRGT011R (C prefix — 다른 해외주식 TR_ID 와 다름).
// ANOMALY 2: CTX_AREA_NK50/FK50 cursor pagination params 사용.
// ANOMALY 3: 수치 콘텐츠 필드 (비율/금액) 모두 String 타입 (KIS docs).
type PeriodRights struct {
	Output []PeriodRightsItem `json:"output"`
}

// PeriodRightsItem 은 해외주식_기간별권리조회 한 행. 모든 필드 string (KIS docs).
type PeriodRightsItem struct {
	BassDt           string `json:"bass_dt"`             // 기준일자
	RghtTypeCd       string `json:"rght_type_cd"`        // 권리유형코드
	Pdno             string `json:"pdno"`                // 상품번호
	PrdtName         string `json:"prdt_name"`           // 상품명
	PrdtTypeCd       string `json:"prdt_type_cd"`        // 상품유형코드
	StdPdno          string `json:"std_pdno"`            // 표준상품번호
	AcplBassDt       string `json:"acpl_bass_dt"`        // 발생기준일
	SbscStrtDt       string `json:"sbsc_strt_dt"`        // 청약시작일
	SbscEndDt        string `json:"sbsc_end_dt"`         // 청약종료일
	CashAlctRt       string `json:"cash_alct_rt"`        // 현금배분율
	StckAlctRt       string `json:"stck_alct_rt"`        // 주식배분율
	CrcyCd           string `json:"crcy_cd"`             // 통화코드1
	CrcyCd2          string `json:"crcy_cd2"`            // 통화코드2
	CrcyCd3          string `json:"crcy_cd3"`            // 통화코드3
	CrcyCd4          string `json:"crcy_cd4"`            // 통화코드4
	AlctFrcrUnpr     string `json:"alct_frcr_unpr"`      // 배분외화단가
	StkpDvdnFrcrAmt2 string `json:"stkp_dvdn_frcr_amt2"` // 주식배당외화금액2
	StkpDvdnFrcrAmt3 string `json:"stkp_dvdn_frcr_amt3"` // 주식배당외화금액3
	StkpDvdnFrcrAmt4 string `json:"stkp_dvdn_frcr_amt4"` // 주식배당외화금액4
	DfntYn           string `json:"dfnt_yn"`             // 확정여부
}

// InquirePeriodRightsParams 는 해외주식_기간별권리조회 파라미터.
//
// CtxAreaNk50/CtxAreaFk50 는 cursor pagination. 첫 조회 시 "" (공백).
type InquirePeriodRightsParams struct {
	RghtTypeCd  string // RGHT_TYPE_CD — 권리유형코드. 빈 값 default
	InqrDvsnCd  string // INQR_DVSN_CD — 조회구분코드. 빈 값 default
	InqrStrtDt  string // INQR_STRT_DT — 조회시작일 YYYYMMDD. 빈 값 default
	InqrEndDt   string // INQR_END_DT — 조회종료일 YYYYMMDD. 빈 값 default
	Pdno        string // PDNO — 상품번호. 빈 값 default (전체)
	PrdtTypeCd  string // PRDT_TYPE_CD — 상품유형코드. 빈 값 default
	CtxAreaNk50 string // CTX_AREA_NK50 — cursor pagination. 빈 값=첫 페이지
	CtxAreaFk50 string // CTX_AREA_FK50 — cursor pagination. 빈 값=첫 페이지
}

// InquirePeriodRights 는 해외주식_기간별권리조회 호출.
//
// 한투 docs: docs/api/해외주식/해외주식_기간별권리조회.md
// path: /uapi/overseas-price/v1/quotations/period-rights (CTRGT011R)
func (c *Client) InquirePeriodRights(ctx context.Context, params InquirePeriodRightsParams) (*PeriodRights, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-price/v1/quotations/period-rights",
		TrID:   "CTRGT011R",
		Query: map[string]string{
			"RGHT_TYPE_CD":  params.RghtTypeCd,
			"INQR_DVSN_CD":  params.InqrDvsnCd,
			"INQR_STRT_DT":  params.InqrStrtDt,
			"INQR_END_DT":   params.InqrEndDt,
			"PDNO":          params.Pdno,
			"PRDT_TYPE_CD":  params.PrdtTypeCd,
			"CTX_AREA_NK50": params.CtxAreaNk50,
			"CTX_AREA_FK50": params.CtxAreaFk50,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res PeriodRights
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse PeriodRights: %w", err)
	}
	return &res, nil
}
