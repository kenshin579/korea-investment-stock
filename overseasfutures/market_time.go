package overseasfutures

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// ─── EP10: MarketTime ────────────────────────────────────────────────────────

// MarketTimeItem 는 해외선물옵션 장운영시간 항목 (15 필드).
// 주의: 응답 키가 output1/output2 가 아닌 `output` 직접 사용 — 공통 TR 패턴.
type MarketTimeItem struct {
	FmPdgrCd        string `json:"fm_pdgr_cd"`         // FM상품군코드
	FmPdgrName      string `json:"fm_pdgr_name"`       // FM상품군명
	FmExcgCd        string `json:"fm_excg_cd"`         // FM거래소코드
	FmExcgName      string `json:"fm_excg_name"`       // FM거래소명
	FuopDvsnName    string `json:"fuop_dvsn_name"`     // 선물옵션구분명
	FmClasCd        string `json:"fm_clas_cd"`         // FM클래스코드
	FmClasName      string `json:"fm_clas_name"`       // FM클래스명
	AmMkmnStrtTmd   string `json:"am_mkmn_strt_tmd"`   // 오전장운영시작시각
	AmMkmnEndTmd    string `json:"am_mkmn_end_tmd"`    // 오전장운영종료시각
	PmMkmnStrtTmd   string `json:"pm_mkmn_strt_tmd"`   // 오후장운영시작시각
	PmMkmnEndTmd    string `json:"pm_mkmn_end_tmd"`    // 오후장운영종료시각
	MkmnNxdyStrtTmd string `json:"mkmn_nxdy_strt_tmd"` // 장운영익일시작시각
	MkmnNxdyEndTmd  string `json:"mkmn_nxdy_end_tmd"`  // 장운영익일종료시각
	BaseMketStrtTmd string `json:"base_mket_strt_tmd"` // 기본시장시작시각
	BaseMketEndTmd  string `json:"base_mket_end_tmd"`  // 기본시장종료시각
}

// MarketTimeData 는 해외선물옵션 장운영시간 응답.
// 주의: 응답 키가 `output` (output1/output2 아님) — 공통 TR 특이 패턴.
type MarketTimeData struct {
	Output []MarketTimeItem `json:"output"`
}

type marketTimeResponse struct {
	RtCd   string           `json:"rt_cd"`
	MsgCd  string           `json:"msg_cd"`
	Msg1   string           `json:"msg1"`
	Output []MarketTimeItem `json:"output"`
}

// MarketTimeParams 는 해외선물옵션 장운영시간 조회 파라미터.
type MarketTimeParams struct {
	FmPdgrCd     string // FM상품군코드 (공백)
	FmClasCd     string // FM클래스코드 ('공백':전체, '001':통화, '002':금리, '003':지수, '004':농산물, '005':축산물, '006':금속, '007':에너지)
	FmExcgCd     string // FM거래소코드 (CME/EUREX/HKEx/ICE/SGX/OSE/ASX/CBOE/MDEX/NYSE/BMF/FTX/HNX/ETC)
	OptYn        string // 옵션여부 (%:전체, N:선물, Y:옵션)
	CtxAreaNK200 string // 연속조회키200
	CtxAreaFK200 string // 연속조회검색조건200
}

// MarketTime 는 해외선물옵션 장운영시간 조회 (OTFM2229R).
//
// 선물옵션 공통 TR. output 키 직접 사용 (output1/output2 아님).
// 연속조회 지원: CtxAreaNK200/CtxAreaFK200 이용.
// OptYn=Y 로 옵션 전용 조회, OptYn=% 로 선물+옵션 전체 조회.
//
// KIS API: GET /uapi/overseas-futureoption/v1/quotations/market-time
func (c *Client) MarketTime(ctx context.Context, params MarketTimeParams) (*MarketTimeData, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/overseas-futureoption/v1/quotations/market-time",
		TrID:     "OTFM2229R",
		CustType: "P",
		Query: map[string]string{
			"FM_PDGR_CD":     params.FmPdgrCd,
			"FM_CLAS_CD":     params.FmClasCd,
			"FM_EXCG_CD":     params.FmExcgCd,
			"OPT_YN":         params.OptYn,
			"CTX_AREA_NK200": params.CtxAreaNK200,
			"CTX_AREA_FK200": params.CtxAreaFK200,
		},
	})
	if err != nil {
		return nil, err
	}
	var res marketTimeResponse
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse MarketTime: %w", err)
	}
	return &MarketTimeData{Output: res.Output}, nil
}
