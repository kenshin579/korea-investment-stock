package overseas

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
	"github.com/kenshin579/korea-investment-stock/kistypes"
)

// 해외주식 잔고 TR ID. 실전/모의가 다르다.
const (
	trIDBalanceReal  = "TTTS3012R"
	trIDBalancePaper = "VTTS3012R"
)

// Balance 는 해외주식 잔고 (TTTS3012R) 응답.
//
// 한투 docs: docs/api/해외주식/해외주식_잔고.md
// path: /uapi/overseas-stock/v1/trading/inquire-balance
//
// 한 번의 호출에 실전 최대 100건. 더 있으면 TrCont 가 "F"/"M" 이고 CtxAreaFk200/CtxAreaNk200 을
// 다음 호출 파라미터로 넘긴다. 전체를 한 번에 받으려면 InquireBalanceAll 을 쓴다.
// 금액·수량은 거래 통화(TrCrcyCd) 기준 외화이며 원화 환산은 없다. 미니스탁 잔고는 포함되지 않는다.
// 모의투자 도메인(WithPaperEnv)이면 VTTS3012R 로 자동 분기한다.
type Balance struct {
	Output1      []BalanceItem  `json:"output1"`        // 보유 종목
	Output2      BalanceSummary `json:"output2"`        // 계좌 요약 (단일 객체 — 국내와 다름)
	CtxAreaFk200 string         `json:"ctx_area_fk200"` // 연속조회검색조건200
	CtxAreaNk200 string         `json:"ctx_area_nk200"` // 연속조회키200
	TrCont       string         `json:"-"`              // 응답 헤더 tr_cont. F/M: 다음 있음, D/E/"": 마지막
}

// BalanceItem 은 보유 종목 1건 (output1). 숫자는 부호(+/-)·빈 문자열을 허용하는 kistypes.Float.
type BalanceItem struct {
	Cano            string         `json:"cano"`               // 종합계좌번호
	AcntPrdtCd      string         `json:"acnt_prdt_cd"`       // 계좌상품코드
	PrdtTypeCd      string         `json:"prdt_type_cd"`       // 상품유형코드
	OvrsPdno        string         `json:"ovrs_pdno"`          // 해외상품번호 (티커)
	OvrsItemName    string         `json:"ovrs_item_name"`     // 해외종목명
	FrcrEvluPflsAmt kistypes.Float `json:"frcr_evlu_pfls_amt"` // 외화평가손익금액
	EvluPflsRt      kistypes.Float `json:"evlu_pfls_rt"`       // 평가손익율
	PchsAvgPric     kistypes.Float `json:"pchs_avg_pric"`      // 매입평균가격 (외화)
	OvrsCblcQty     kistypes.Float `json:"ovrs_cblc_qty"`      // 해외잔고수량
	OrdPsblQty      kistypes.Float `json:"ord_psbl_qty"`       // 주문가능수량
	FrcrPchsAmt1    kistypes.Float `json:"frcr_pchs_amt1"`     // 외화매입금액
	OvrsStckEvluAmt kistypes.Float `json:"ovrs_stck_evlu_amt"` // 해외주식평가금액 (외화)
	NowPric2        kistypes.Float `json:"now_pric2"`          // 현재가
	TrCrcyCd        string         `json:"tr_crcy_cd"`         // 거래통화코드 USD/HKD/CNY/JPY/VND
	OvrsExcgCd      string         `json:"ovrs_excg_cd"`       // 해외거래소코드 NASD/NYSE/AMEX/SEHK/SHAA/SZAA/TKSE/HASE/VNSE
	LoanTypeCd      string         `json:"loan_type_cd"`       // 대출유형코드 (00: 해당없음)
	LoanDt          string         `json:"loan_dt"`            // 대출일자
	ExpdDt          string         `json:"expd_dt"`            // 만기일자
}

// BalanceSummary 는 계좌 요약 (output2, 단일 객체).
type BalanceSummary struct {
	FrcrPchsAmt1     kistypes.Float `json:"frcr_pchs_amt1"`      // 외화매입금액1
	OvrsRlztPflsAmt  kistypes.Float `json:"ovrs_rlzt_pfls_amt"`  // 해외실현손익금액
	OvrsTotPfls      kistypes.Float `json:"ovrs_tot_pfls"`       // 해외총손익
	RlztErngRt       kistypes.Float `json:"rlzt_erng_rt"`        // 실현수익율
	TotEvluPflsAmt   kistypes.Float `json:"tot_evlu_pfls_amt"`   // 총평가손익금액
	TotPftrt         kistypes.Float `json:"tot_pftrt"`           // 총수익률
	FrcrBuyAmtSmtl1  kistypes.Float `json:"frcr_buy_amt_smtl1"`  // 외화매수금액합계1
	OvrsRlztPflsAmt2 kistypes.Float `json:"ovrs_rlzt_pfls_amt2"` // 해외실현손익금액2
	FrcrBuyAmtSmtl2  kistypes.Float `json:"frcr_buy_amt_smtl2"`  // 외화매수금액합계2
}

// InquireBalanceParams 는 해외주식 잔고 파라미터. 계좌번호는 Client 설정값을 쓴다.
// OvrsExcgCd 와 TrCrcyCd 는 필수 — 실전 미국 전체는 ("NASD", "USD").
// InquireBalanceAll 은 CtxAreaFk200/CtxAreaNk200/TrCont 를 무시하고 첫 페이지부터 읽는다.
type InquireBalanceParams struct {
	OvrsExcgCd   string // OVRS_EXCG_CD (필수) — 실전: NASD 미국전체 / NAS 나스닥 / NYSE / AMEX · 공통: SEHK / SHAA / SZAA / TKSE / HASE / VNSE
	TrCrcyCd     string // TR_CRCY_CD (필수) — USD / HKD / CNY / JPY / VND
	CtxAreaFk200 string // CTX_AREA_FK200 — 연속조회. 첫 조회 빈 값
	CtxAreaNk200 string // CTX_AREA_NK200 — 연속조회. 첫 조회 빈 값
	TrCont       string // tr_cont 헤더 — 연속조회 "N". 첫 조회 빈 값
}

// InquireBalance 는 해외주식 잔고 1페이지 호출 (실전 최대 100건).
//
// 한투 docs: docs/api/해외주식/해외주식_잔고.md
// path: /uapi/overseas-stock/v1/trading/inquire-balance (TTTS3012R, 모의 VTTS3012R)
func (c *Client) InquireBalance(ctx context.Context, params InquireBalanceParams) (*Balance, error) {
	if params.OvrsExcgCd == "" {
		return nil, errors.New("kis: InquireBalance: OvrsExcgCd is required (e.g. NASD)")
	}
	if params.TrCrcyCd == "" {
		return nil, errors.New("kis: InquireBalance: TrCrcyCd is required (e.g. USD)")
	}
	cano, prdtCd, err := c.http.Account()
	if err != nil {
		return nil, err
	}
	trID := trIDBalanceReal
	if c.http.IsPaper() {
		trID = trIDBalancePaper
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-stock/v1/trading/inquire-balance",
		TrID:   trID,
		TrCont: params.TrCont,
		Query: map[string]string{
			"CANO":           cano,
			"ACNT_PRDT_CD":   prdtCd,
			"OVRS_EXCG_CD":   params.OvrsExcgCd,
			"TR_CRCY_CD":     params.TrCrcyCd,
			"CTX_AREA_FK200": params.CtxAreaFk200,
			"CTX_AREA_NK200": params.CtxAreaNk200,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res Balance
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse Balance: %w", err)
	}
	res.TrCont = resp.TrCont
	return &res, nil
}

// maxBalancePages 는 연속조회 상한 (100건 × 100 페이지). 잘못된 tr_cont 로 인한 무한 루프 방지.
const maxBalancePages = 100

// InquireBalanceAll 은 연속조회(tr_cont)를 따라가며 보유 종목 전체를 모은다.
// params 의 TrCont/CtxArea* 는 무시하고 첫 페이지부터 읽는다.
// Output1 은 모든 페이지를 이어 붙이고, Output2·CtxArea*·TrCont 는 마지막 페이지 값이다.
// 중간 페이지 실패 시 부분 결과 없이 error 만 반환한다.
func (c *Client) InquireBalanceAll(ctx context.Context, params InquireBalanceParams) (*Balance, error) {
	params.TrCont, params.CtxAreaFk200, params.CtxAreaNk200 = "", "", ""
	var all *Balance
	for page := 0; page < maxBalancePages; page++ {
		res, err := c.InquireBalance(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("kis: InquireBalanceAll: page %d: %w", page+1, err)
		}
		if all == nil {
			all = res
		} else {
			all.Output1 = append(all.Output1, res.Output1...)
			all.Output2 = res.Output2
			all.CtxAreaFk200, all.CtxAreaNk200, all.TrCont = res.CtxAreaFk200, res.CtxAreaNk200, res.TrCont
		}
		if !httpclient.HasNext(res.TrCont) {
			return all, nil
		}
		if res.CtxAreaFk200 == params.CtxAreaFk200 && res.CtxAreaNk200 == params.CtxAreaNk200 {
			return nil, fmt.Errorf("kis: InquireBalanceAll: tr_cont=%q but cursor did not advance (page %d)", res.TrCont, page+1)
		}
		params.TrCont, params.CtxAreaFk200, params.CtxAreaNk200 = "N", res.CtxAreaFk200, res.CtxAreaNk200
	}
	return nil, fmt.Errorf("kis: InquireBalanceAll: exceeded %d pages", maxBalancePages)
}
