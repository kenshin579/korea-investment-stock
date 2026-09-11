package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
	"github.com/kenshin579/korea-investment-stock/kistypes"
)

// 주식잔고조회 TR ID. 실전/모의가 다르다.
const (
	trIDBalanceReal  = "TTTC8434R"
	trIDBalancePaper = "VTTC8434R"
)

// Balance 는 주식잔고조회 (TTTC8434R/VTTC8434R) 응답.
//
// 한투 docs: docs/api/국내주식/주식잔고조회.md
// path: /uapi/domestic-stock/v1/trading/inquire-balance
//
// 한 번의 호출에 실전 최대 50건(모의 20건). 더 있으면 TrCont 가 "F"/"M" 이고 CtxAreaFk100/CtxAreaNk100 을
// 다음 호출 파라미터로 넘긴다. 전체를 한 번에 받으려면 InquireBalanceAll 을 쓴다.
// 당일 전량 매도한 종목은 HldgQty 0 으로 남아 있을 수 있다(D-2 이후 사라짐).
// 모의투자 도메인(WithPaperEnv)이면 VTTC8434R 로 자동 분기한다(모의는 한 번에 최대 20건).
type Balance struct {
	Output1      []BalanceItem    `json:"output1"`        // 보유 종목
	Output2      []BalanceSummary `json:"output2"`        // 계좌 요약 (1행)
	CtxAreaFk100 string           `json:"ctx_area_fk100"` // 연속조회검색조건100
	CtxAreaNk100 string           `json:"ctx_area_nk100"` // 연속조회키100
	TrCont       string           `json:"-"`              // 응답 헤더 tr_cont. F/M: 다음 있음, D/E/"": 마지막
}

// BalanceItem 은 보유 종목 1건 (output1). 금액은 원 단위 정수, 단가·비율은 실수.
type BalanceItem struct {
	Pdno           string         `json:"pdno"`              // 종목번호 (6자리)
	PrdtName       string         `json:"prdt_name"`         // 종목명
	TradDvsnName   string         `json:"trad_dvsn_name"`    // 매매구분명
	BfdyBuyQty     kistypes.Int   `json:"bfdy_buy_qty"`      // 전일매수수량
	BfdySllQty     kistypes.Int   `json:"bfdy_sll_qty"`      // 전일매도수량
	ThdtBuyqty     kistypes.Int   `json:"thdt_buyqty"`       // 금일매수수량
	ThdtSllQty     kistypes.Int   `json:"thdt_sll_qty"`      // 금일매도수량
	HldgQty        kistypes.Int   `json:"hldg_qty"`          // 보유수량
	OrdPsblQty     kistypes.Int   `json:"ord_psbl_qty"`      // 주문가능수량
	PchsAvgPric    kistypes.Float `json:"pchs_avg_pric"`     // 매입평균가격 (매입금액/보유수량)
	PchsAmt        kistypes.Int   `json:"pchs_amt"`          // 매입금액
	Prpr           kistypes.Int   `json:"prpr"`              // 현재가
	EvluAmt        kistypes.Int   `json:"evlu_amt"`          // 평가금액
	EvluPflsAmt    kistypes.Int   `json:"evlu_pfls_amt"`     // 평가손익금액 (평가금액 - 매입금액)
	EvluPflsRt     kistypes.Float `json:"evlu_pfls_rt"`      // 평가손익율
	EvluErngRt     kistypes.Float `json:"evlu_erng_rt"`      // 평가수익율 (미사용, 0)
	LoanDt         string         `json:"loan_dt"`           // 대출일자 (INQR_DVSN=01 일 때만)
	LoanAmt        kistypes.Int   `json:"loan_amt"`          // 대출금액
	StlnSlngChgs   kistypes.Int   `json:"stln_slng_chgs"`    // 대주매각대금
	ExpdDt         string         `json:"expd_dt"`           // 만기일자
	FlttRt         kistypes.Float `json:"fltt_rt"`           // 등락율
	BfdyCprsIcdc   kistypes.Int   `json:"bfdy_cprs_icdc"`    // 전일대비증감
	ItemMgnaRtName string         `json:"item_mgna_rt_name"` // 종목증거금율명
	GrtaRtName     string         `json:"grta_rt_name"`      // 보증금율명
	SbstPric       kistypes.Int   `json:"sbst_pric"`         // 대용가격
	StckLoanUnpr   kistypes.Float `json:"stck_loan_unpr"`    // 주식대출단가
}

// BalanceSummary 는 계좌 요약 (output2, 1행).
type BalanceSummary struct {
	DncaTotAmt         kistypes.Int   `json:"dnca_tot_amt"`           // 예수금총금액
	NxdyExccAmt        kistypes.Int   `json:"nxdy_excc_amt"`          // 익일정산금액 (D+1 예수금)
	PrvsRcdlExccAmt    kistypes.Int   `json:"prvs_rcdl_excc_amt"`     // 가수도정산금액 (D+2 예수금)
	CmaEvluAmt         kistypes.Int   `json:"cma_evlu_amt"`           // CMA평가금액
	BfdyBuyAmt         kistypes.Int   `json:"bfdy_buy_amt"`           // 전일매수금액
	ThdtBuyAmt         kistypes.Int   `json:"thdt_buy_amt"`           // 금일매수금액
	NxdyAutoRdptAmt    kistypes.Int   `json:"nxdy_auto_rdpt_amt"`     // 익일자동상환금액
	BfdySllAmt         kistypes.Int   `json:"bfdy_sll_amt"`           // 전일매도금액
	ThdtSllAmt         kistypes.Int   `json:"thdt_sll_amt"`           // 금일매도금액
	D2AutoRdptAmt      kistypes.Int   `json:"d2_auto_rdpt_amt"`       // D+2자동상환금액
	BfdyTlexAmt        kistypes.Int   `json:"bfdy_tlex_amt"`          // 전일제비용금액
	ThdtTlexAmt        kistypes.Int   `json:"thdt_tlex_amt"`          // 금일제비용금액
	TotLoanAmt         kistypes.Int   `json:"tot_loan_amt"`           // 총대출금액
	SctsEvluAmt        kistypes.Int   `json:"scts_evlu_amt"`          // 유가평가금액
	TotEvluAmt         kistypes.Int   `json:"tot_evlu_amt"`           // 총평가금액 (유가증권 평가 + D+2 예수금)
	NassAmt            kistypes.Int   `json:"nass_amt"`               // 순자산금액
	FncgGldAutoRdptYn  string         `json:"fncg_gld_auto_rdpt_yn"`  // 융자금자동상환여부
	PchsAmtSmtlAmt     kistypes.Int   `json:"pchs_amt_smtl_amt"`      // 매입금액합계금액
	EvluAmtSmtlAmt     kistypes.Int   `json:"evlu_amt_smtl_amt"`      // 평가금액합계금액 (유가증권)
	EvluPflsSmtlAmt    kistypes.Int   `json:"evlu_pfls_smtl_amt"`     // 평가손익합계금액
	TotStlnSlngChgs    kistypes.Int   `json:"tot_stln_slng_chgs"`     // 총대주매각대금
	BfdyTotAsstEvluAmt kistypes.Int   `json:"bfdy_tot_asst_evlu_amt"` // 전일총자산평가금액
	AsstIcdcAmt        kistypes.Int   `json:"asst_icdc_amt"`          // 자산증감액
	AsstIcdcErngRt     kistypes.Float `json:"asst_icdc_erng_rt"`      // 자산증감수익율 (데이터 미제공)
}

// InquireBalanceParams 는 주식잔고조회 파라미터. 계좌번호(CANO/ACNT_PRDT_CD)는 Client 설정값을 쓴다.
//
// 빈 값은 한투 기본값으로 채운다: AfhrFlprYn "N", InqrDvsn "02"(종목별), UnprDvsn "01",
// FundSttlIcldYn "N", FncgAmtAutoRdptYn "N", PrcsDvsn "00"(전일매매포함).
// 연속조회는 이전 응답의 CtxAreaFk100/CtxAreaNk100 과 TrCont "N" 을 넘긴다.
// InquireBalanceAll 은 CtxAreaFk100/CtxAreaNk100/TrCont 를 무시하고 첫 페이지부터 읽는다.
type InquireBalanceParams struct {
	AfhrFlprYn        string // AFHR_FLPR_YN — N: 기본, Y: 시간외단일가, X: NXT 정규장
	InqrDvsn          string // INQR_DVSN — 01: 대출일별, 02: 종목별
	UnprDvsn          string // UNPR_DVSN — 01: 기본
	FundSttlIcldYn    string // FUND_STTL_ICLD_YN — 펀드결제분 포함 N/Y
	FncgAmtAutoRdptYn string // FNCG_AMT_AUTO_RDPT_YN — 융자금액자동상환 N
	PrcsDvsn          string // PRCS_DVSN — 00: 전일매매포함, 01: 전일매매미포함
	CtxAreaFk100      string // CTX_AREA_FK100 — 연속조회. 첫 조회 빈 값
	CtxAreaNk100      string // CTX_AREA_NK100 — 연속조회. 첫 조회 빈 값
	TrCont            string // tr_cont 헤더 — 연속조회 "N". 첫 조회 빈 값
}

// balanceParamOr 는 빈 파라미터를 한투 기본값으로 채운다.
func balanceParamOr(v, def string) string {
	if v == "" {
		return def
	}
	return v
}

// InquireBalance 는 주식잔고조회 1페이지 호출 (실전 최대 50건, 모의 20건).
//
// 한투 docs: docs/api/국내주식/주식잔고조회.md
// path: /uapi/domestic-stock/v1/trading/inquire-balance (실전 TTTC8434R, 모의 VTTC8434R)
func (c *Client) InquireBalance(ctx context.Context, params InquireBalanceParams) (*Balance, error) {
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
		Path:   "/uapi/domestic-stock/v1/trading/inquire-balance",
		TrID:   trID,
		TrCont: params.TrCont,
		Query: map[string]string{
			"CANO":                  cano,
			"ACNT_PRDT_CD":          prdtCd,
			"AFHR_FLPR_YN":          balanceParamOr(params.AfhrFlprYn, "N"),
			"OFL_YN":                "",
			"INQR_DVSN":             balanceParamOr(params.InqrDvsn, "02"),
			"UNPR_DVSN":             balanceParamOr(params.UnprDvsn, "01"),
			"FUND_STTL_ICLD_YN":     balanceParamOr(params.FundSttlIcldYn, "N"),
			"FNCG_AMT_AUTO_RDPT_YN": balanceParamOr(params.FncgAmtAutoRdptYn, "N"),
			"PRCS_DVSN":             balanceParamOr(params.PrcsDvsn, "00"),
			"CTX_AREA_FK100":        params.CtxAreaFk100,
			"CTX_AREA_NK100":        params.CtxAreaNk100,
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

// maxBalancePages 는 연속조회 상한. 50건 × 100 = 5,000 종목 — 실계좌에서 도달할 수 없는 값이며
// 한투가 tr_cont 를 잘못 주는 경우의 무한 루프 방지용.
const maxBalancePages = 100

// InquireBalanceAll 은 연속조회(tr_cont)를 따라가며 보유 종목 전체를 모은다.
// params 의 TrCont/CtxArea* 는 무시하고 첫 페이지부터 읽는다.
// Output1 은 모든 페이지를 이어 붙이고, Output2·CtxArea*·TrCont 는 마지막 페이지 값이다.
// 중간 페이지 실패 시 부분 결과 없이 error 만 반환한다.
func (c *Client) InquireBalanceAll(ctx context.Context, params InquireBalanceParams) (*Balance, error) {
	params.TrCont, params.CtxAreaFk100, params.CtxAreaNk100 = "", "", ""
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
			all.CtxAreaFk100, all.CtxAreaNk100, all.TrCont = res.CtxAreaFk100, res.CtxAreaNk100, res.TrCont
		}
		if !httpclient.HasNext(res.TrCont) {
			return all, nil
		}
		if res.CtxAreaFk100 == params.CtxAreaFk100 && res.CtxAreaNk100 == params.CtxAreaNk100 {
			return nil, fmt.Errorf("kis: InquireBalanceAll: tr_cont=%q but cursor did not advance (page %d)", res.TrCont, page+1)
		}
		params.TrCont, params.CtxAreaFk100, params.CtxAreaNk100 = "N", res.CtxAreaFk100, res.CtxAreaNk100
	}
	return nil, fmt.Errorf("kis: InquireBalanceAll: exceeded %d pages", maxBalancePages)
}
