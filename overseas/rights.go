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
// ANOMALY 4: 문서는 "tr_cont 를 이용한 다음조회 불가 API" 라고 하지만 실 응답(2026-11-01 실측)은 한 번에 최대
// 100행이고, 더 있으면 응답 헤더 tr_cont "M" + 본문 ctx_area_nk50/fk50(뒤 공백 포함 50자) 를 준다. 다음 페이지는
// 그 커서를 그대로 돌려보내면서 **요청 헤더 tr_cont: "N"** 을 붙여야 나온다 — 커서만 보내면 첫 페이지가 되풀이된다.
// 종목(PDNO)을 비우면 그 기간의 시장 전체(미국·홍콩·일본 등 섞임)가 온다. 자료가 없으면 rt_cd "0" + msg_cd KIER2520 +
// 빈 output 이다(에러가 아니다).
type PeriodRights struct {
	Output      []PeriodRightsItem `json:"output"`
	CtxAreaNk50 string             `json:"ctx_area_nk50"` // 다음 페이지 커서 — 공백 포함 그대로 다음 호출에 넘긴다
	CtxAreaFk50 string             `json:"ctx_area_fk50"`
	TrCont      string             `json:"-"` // 응답 헤더 tr_cont. M/F: 다음 있음, D/E/"": 마지막
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
// 다음 페이지는 직전 응답의 커서 두 개와 TrCont "N" 을 함께 넘긴다(ANOMALY 4). 전부 받으려면 InquirePeriodRightsAll.
type InquirePeriodRightsParams struct {
	RghtTypeCd  string // RGHT_TYPE_CD — 권리유형코드. %%(전체) 03(배당) 75(특별배당) 등. 빈 값 default
	InqrDvsnCd  string // INQR_DVSN_CD — 조회구분코드. 02(현지기준일) 03(청약시작일) 04(청약종료일). 빈 값 default
	InqrStrtDt  string // INQR_STRT_DT — 조회시작일 YYYYMMDD. 빈 값 default
	InqrEndDt   string // INQR_END_DT — 조회종료일 YYYYMMDD. 빈 값 default
	Pdno        string // PDNO — 상품번호. 빈 값 default (전체)
	PrdtTypeCd  string // PRDT_TYPE_CD — 상품유형코드. 빈 값 default
	CtxAreaNk50 string // CTX_AREA_NK50 — cursor pagination. 빈 값=첫 페이지
	CtxAreaFk50 string // CTX_AREA_FK50 — cursor pagination. 빈 값=첫 페이지
	TrCont      string // tr_cont 요청 헤더 — 연속조회 "N". 첫 조회 빈 값
}

// InquirePeriodRights 는 해외주식_기간별권리조회 1페이지 호출 (최대 100행).
//
// 한투 docs: docs/api/해외주식/해외주식_기간별권리조회.md
// path: /uapi/overseas-price/v1/quotations/period-rights (CTRGT011R)
func (c *Client) InquirePeriodRights(ctx context.Context, params InquirePeriodRightsParams) (*PeriodRights, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-price/v1/quotations/period-rights",
		TrID:   "CTRGT011R",
		TrCont: params.TrCont,
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
	res.TrCont = resp.TrCont
	return &res, nil
}

// maxPeriodRightsPages 는 연속조회 상한 (100행 × 1000 페이지). 시장 전체 배당 한 달이 약 30페이지(2026-09 실측
// 2,923행)라 여유 있게 잡았다. 잘못된 tr_cont 로 인한 무한 루프 방지.
const maxPeriodRightsPages = 1000

// InquirePeriodRightsAll 은 연속조회(tr_cont)를 따라가며 기간 전체 행을 모은다.
// params 의 TrCont/CtxArea* 는 무시하고 첫 페이지부터 읽는다.
// Output 은 모든 페이지를 이어 붙이고, CtxArea*·TrCont 는 마지막 페이지 값이다.
// 중간 페이지 실패 시 부분 결과 없이 error 만 반환한다. 긴 기간은 호출자가 달 단위 등으로 나눠 부르는 편이 낫다.
func (c *Client) InquirePeriodRightsAll(ctx context.Context, params InquirePeriodRightsParams) (*PeriodRights, error) {
	params.TrCont, params.CtxAreaNk50, params.CtxAreaFk50 = "", "", ""
	var all *PeriodRights
	for page := 0; page < maxPeriodRightsPages; page++ {
		res, err := c.InquirePeriodRights(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("kis: InquirePeriodRightsAll: page %d: %w", page+1, err)
		}
		if all == nil {
			all = res
		} else {
			all.Output = append(all.Output, res.Output...)
			all.CtxAreaNk50, all.CtxAreaFk50, all.TrCont = res.CtxAreaNk50, res.CtxAreaFk50, res.TrCont
		}
		if !httpclient.HasNext(res.TrCont) {
			return all, nil
		}
		if res.CtxAreaNk50 == params.CtxAreaNk50 && res.CtxAreaFk50 == params.CtxAreaFk50 {
			return nil, fmt.Errorf("kis: InquirePeriodRightsAll: tr_cont=%q but cursor did not advance (page %d)", res.TrCont, page+1)
		}
		params.TrCont, params.CtxAreaNk50, params.CtxAreaFk50 = "N", res.CtxAreaNk50, res.CtxAreaFk50
	}
	return nil, fmt.Errorf("kis: InquirePeriodRightsAll: exceeded %d pages", maxPeriodRightsPages)
}
