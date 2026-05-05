// File: domestic/ksd.go
package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// KsdDividend 는 예탁원정보(배당일정) (HHKDB669102C0) 응답.
//
// 한투 docs: docs/api/국내주식/예탁원정보(배당일정).md
// path: /uapi/domestic-stock/v1/ksdinfo/dividend
type KsdDividend struct {
	Output1 []KsdDividendItem `json:"output1"`
}

// KsdDividendItem 은 배당일정 한 행. 모든 필드 string (KIS docs).
type KsdDividendItem struct {
	RecordDate    string `json:"record_date"`      // 기준일
	ShtCd         string `json:"sht_cd"`           // 종목코드
	IsinName      string `json:"isin_name"`        // 종목명
	DiviKind      string `json:"divi_kind"`        // 배당종류
	FaceVal       string `json:"face_val"`         // 액면가
	PerStoDiviAmt string `json:"per_sto_divi_amt"` // 현금배당금
	DiviRate      string `json:"divi_rate"`        // 현금배당률(%)
	StkDiviRate   string `json:"stk_divi_rate"`    // 주식배당률(%)
	DiviPayDt     string `json:"divi_pay_dt"`      // 배당금지급일
	StkDivPayDt   string `json:"stk_div_pay_dt"`   // 주식배당지급일
	OddPayDt      string `json:"odd_pay_dt"`       // 단주대금지급일
	StkKind       string `json:"stk_kind"`         // 주식종류
	HighDiviGb    string `json:"high_divi_gb"`     // 고배당종목여부
}

// InquireKsdDividendParams 는 배당일정 조회 파라미터.
type InquireKsdDividendParams struct {
	Cts      string // CTS — 공백 입력 (default "")
	Gb1      string // GB1 — 0:전체, 1:결산배당, 2:중간배당. 빈 값=>"0"
	FromDate string // F_DT — 조회시작일 YYYYMMDD
	ToDate   string // T_DT — 조회종료일 YYYYMMDD
	Symbol   string // SHT_CD — 종목코드 (공백=전체)
	HighGb   string // HIGH_GB — 공백 입력
}

// InquireKsdDividend 호출.
//
// 한투 docs: docs/api/국내주식/예탁원정보(배당일정).md
// path: /uapi/domestic-stock/v1/ksdinfo/dividend (HHKDB669102C0)
func (c *Client) InquireKsdDividend(ctx context.Context, params InquireKsdDividendParams) (*KsdDividend, error) {
	gb1 := params.Gb1
	if gb1 == "" {
		gb1 = "0"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ksdinfo/dividend",
		TrID:   "HHKDB669102C0",
		Query: map[string]string{
			"CTS":     params.Cts,
			"GB1":     gb1,
			"F_DT":    params.FromDate,
			"T_DT":    params.ToDate,
			"SHT_CD":  params.Symbol,
			"HIGH_GB": params.HighGb,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res KsdDividend
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse KsdDividend: %w", err)
	}
	return &res, nil
}

// KsdBonusIssue 는 예탁원정보(무상증자) (HHKDB669101C0) 응답.
//
// 한투 docs: docs/api/국내주식/예탁원정보(무상증자).md
// path: /uapi/domestic-stock/v1/ksdinfo/bonus-issue
type KsdBonusIssue struct {
	Output1 []KsdBonusIssueItem `json:"output1"`
}

// KsdBonusIssueItem 은 무상증자 한 행. 모든 필드 string (KIS docs).
type KsdBonusIssueItem struct {
	RecordDate     string `json:"record_date"`       // 기준일
	ShtCd          string `json:"sht_cd"`            // 종목코드
	IsinName       string `json:"isin_name"`         // 종목명
	FixRate        string `json:"fix_rate"`          // 확정배정율
	OddRecPrice    string `json:"odd_rec_price"`     // 단주기준가
	RightDt        string `json:"right_dt"`          // 권리락일
	OddPayDt       string `json:"odd_pay_dt"`        // 단주대금지급일
	ListDate       string `json:"list_date"`         // 상장/등록일
	TotIssueStkQty string `json:"tot_issue_stk_qty"` // 발행주식
	IssueStkQty    string `json:"issue_stk_qty"`     // 발행할주식
	StkKind        string `json:"stk_kind"`          // 주식종류
}

// InquireKsdBonusIssueParams 는 무상증자 조회 파라미터.
type InquireKsdBonusIssueParams struct {
	Cts      string // CTS — 공백 입력 (default "")
	FromDate string // F_DT — 조회시작일 YYYYMMDD
	ToDate   string // T_DT — 조회종료일 YYYYMMDD
	Symbol   string // SHT_CD — 종목코드 (공백=전체)
}

// InquireKsdBonusIssue 호출.
//
// 한투 docs: docs/api/국내주식/예탁원정보(무상증자).md
// path: /uapi/domestic-stock/v1/ksdinfo/bonus-issue (HHKDB669101C0)
func (c *Client) InquireKsdBonusIssue(ctx context.Context, params InquireKsdBonusIssueParams) (*KsdBonusIssue, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ksdinfo/bonus-issue",
		TrID:   "HHKDB669101C0",
		Query: map[string]string{
			"CTS":    params.Cts,
			"F_DT":   params.FromDate,
			"T_DT":   params.ToDate,
			"SHT_CD": params.Symbol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res KsdBonusIssue
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse KsdBonusIssue: %w", err)
	}
	return &res, nil
}

// KsdPaidinCapin 는 예탁원정보(유상증자) (HHKDB669100C0) 응답.
//
// 한투 docs: docs/api/국내주식/예탁원정보(유상증자).md
// path: /uapi/domestic-stock/v1/ksdinfo/paidin-capin
// ANOMALY: output key 는 "output" (다른 KSD 메서드와 달리 "output1" 아님).
type KsdPaidinCapin struct {
	Output []KsdPaidinCapinItem `json:"output"`
}

// KsdPaidinCapinItem 은 유상증자 한 행. 모든 필드 string (KIS docs).
type KsdPaidinCapinItem struct {
	RecordDate     string `json:"record_date"`       // 기준일
	ShtCd          string `json:"sht_cd"`            // 종목코드
	IsinName       string `json:"isin_name"`         // 종목명
	TotIssueStkQty string `json:"tot_issue_stk_qty"` // 발행주식
	IssueStkQty    string `json:"issue_stk_qty"`     // 발행할주식
	FixRate        string `json:"fix_rate"`          // 확정배정율
	DiscRate       string `json:"disc_rate"`         // 할인율
	FixPrice       string `json:"fix_price"`         // 발행예정가
	RightDt        string `json:"right_dt"`          // 권리락일
	SubTermFt      string `json:"sub_term_ft"`       // 청약기간(시작)
	SubTerm        string `json:"sub_term"`          // 청약기간(종료)
	ListDate       string `json:"list_date"`         // 상장/등록일
	StkKind        string `json:"stk_kind"`          // 주식종류
}

// InquireKsdPaidinCapinParams 는 유상증자 조회 파라미터.
type InquireKsdPaidinCapinParams struct {
	Cts      string // CTS — 공백 입력 (default "")
	Gb1      string // GB1 — 1:청약일별, 2:기준일별. 빈 값=>"1"
	FromDate string // F_DT — 조회시작일 YYYYMMDD
	ToDate   string // T_DT — 조회종료일 YYYYMMDD
	Symbol   string // SHT_CD — 종목코드 (공백=전체)
}

// InquireKsdPaidinCapin 호출.
//
// 한투 docs: docs/api/국내주식/예탁원정보(유상증자).md
// path: /uapi/domestic-stock/v1/ksdinfo/paidin-capin (HHKDB669100C0)
func (c *Client) InquireKsdPaidinCapin(ctx context.Context, params InquireKsdPaidinCapinParams) (*KsdPaidinCapin, error) {
	gb1 := params.Gb1
	if gb1 == "" {
		gb1 = "1"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ksdinfo/paidin-capin",
		TrID:   "HHKDB669100C0",
		Query: map[string]string{
			"CTS":    params.Cts,
			"GB1":    gb1,
			"F_DT":   params.FromDate,
			"T_DT":   params.ToDate,
			"SHT_CD": params.Symbol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res KsdPaidinCapin
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse KsdPaidinCapin: %w", err)
	}
	return &res, nil
}

// KsdSharehldMeet 는 예탁원정보(주주총회) (HHKDB669111C0) 응답.
//
// 한투 docs: docs/api/국내주식/예탁원정보(주주총회).md
// path: /uapi/domestic-stock/v1/ksdinfo/sharehld-meet
type KsdSharehldMeet struct {
	Output1 []KsdSharehldMeetItem `json:"output1"`
}

// KsdSharehldMeetItem 은 주주총회 한 행. 모든 필드 string (KIS docs).
type KsdSharehldMeetItem struct {
	RecordDate  string `json:"record_date"`   // 기준일
	ShtCd       string `json:"sht_cd"`        // 종목코드
	IsinName    string `json:"isin_name"`     // 종목명
	GenMeetDt   string `json:"gen_meet_dt"`   // 주총일자
	GenMeetType string `json:"gen_meet_type"` // 주총사유
	Agenda      string `json:"agenda"`        // 주총의안
	VoteTotQty  string `json:"vote_tot_qty"`  // 의결권주식총수
}

// InquireKsdSharehldMeetParams 는 주주총회 조회 파라미터.
type InquireKsdSharehldMeetParams struct {
	Cts      string // CTS — 공백 입력 (default "")
	FromDate string // F_DT — 조회시작일 YYYYMMDD
	ToDate   string // T_DT — 조회종료일 YYYYMMDD
	Symbol   string // SHT_CD — 종목코드 (공백=전체)
}

// InquireKsdSharehldMeet 호출.
//
// 한투 docs: docs/api/국내주식/예탁원정보(주주총회).md
// path: /uapi/domestic-stock/v1/ksdinfo/sharehld-meet (HHKDB669111C0)
func (c *Client) InquireKsdSharehldMeet(ctx context.Context, params InquireKsdSharehldMeetParams) (*KsdSharehldMeet, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ksdinfo/sharehld-meet",
		TrID:   "HHKDB669111C0",
		Query: map[string]string{
			"CTS":    params.Cts,
			"F_DT":   params.FromDate,
			"T_DT":   params.ToDate,
			"SHT_CD": params.Symbol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res KsdSharehldMeet
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse KsdSharehldMeet: %w", err)
	}
	return &res, nil
}

// KsdMergerSplit 는 예탁원정보(합병분할) (HHKDB669104C0) 응답.
//
// 한투 docs: docs/api/국내주식/예탁원정보(합병분할).md
// path: /uapi/domestic-stock/v1/ksdinfo/merger-split
// ANOMALY: isin_name 없음 — opp_cust_cd/opp_cust_nm (피합병) + cust_cd/cust_nm (합병) 사용.
type KsdMergerSplit struct {
	Output1 []KsdMergerSplitItem `json:"output1"`
}

// KsdMergerSplitItem 은 합병분할 한 행. 모든 필드 string (KIS docs).
// isin_name 없음 — opp_cust_*/cust_* 쌍으로 양사 정보 표현.
type KsdMergerSplitItem struct {
	RecordDate     string `json:"record_date"`       // 기준일
	ShtCd          string `json:"sht_cd"`            // 종목코드
	OppCustCd      string `json:"opp_cust_cd"`       // 피합병(피분할)회사코드
	OppCustNm      string `json:"opp_cust_nm"`       // 피합병(피분할)회사명
	CustCd         string `json:"cust_cd"`           // 합병(분할)회사코드
	CustNm         string `json:"cust_nm"`           // 합병(분할)회사명
	MergeType      string `json:"merge_type"`        // 합병사유
	MergeRate      string `json:"merge_rate"`        // 비율
	TdStopDt       string `json:"td_stop_dt"`        // 매매거래정지기간
	ListDt         string `json:"list_dt"`           // 상장/등록일
	OddAmtPayDt    string `json:"odd_amt_pay_dt"`    // 단주대금지급일
	TotIssueStkQty string `json:"tot_issue_stk_qty"` // 발행주식
	IssueStkQty    string `json:"issue_stk_qty"`     // 발행할주식
	Seq            string `json:"seq"`               // 연번
}

// InquireKsdMergerSplitParams 는 합병분할 조회 파라미터.
type InquireKsdMergerSplitParams struct {
	Cts      string // CTS — 공백 입력 (default "")
	FromDate string // F_DT — 조회시작일 YYYYMMDD
	ToDate   string // T_DT — 조회종료일 YYYYMMDD
	Symbol   string // SHT_CD — 종목코드 (공백=전체)
}

// InquireKsdMergerSplit 호출.
//
// 한투 docs: docs/api/국내주식/예탁원정보(합병분할).md
// path: /uapi/domestic-stock/v1/ksdinfo/merger-split (HHKDB669104C0)
func (c *Client) InquireKsdMergerSplit(ctx context.Context, params InquireKsdMergerSplitParams) (*KsdMergerSplit, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ksdinfo/merger-split",
		TrID:   "HHKDB669104C0",
		Query: map[string]string{
			"CTS":    params.Cts,
			"F_DT":   params.FromDate,
			"T_DT":   params.ToDate,
			"SHT_CD": params.Symbol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res KsdMergerSplit
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse KsdMergerSplit: %w", err)
	}
	return &res, nil
}

// KsdRevSplit 는 예탁원정보(액면변경) (HHKDB669105C0) 응답.
//
// 한투 docs: docs/api/국내주식/예탁원정보(액면변경).md
// path: /uapi/domestic-stock/v1/ksdinfo/rev-split
// ANOMALY: extra MARKET_GB query param (default "0" = 전체).
type KsdRevSplit struct {
	Output1 []KsdRevSplitItem `json:"output1"`
}

// KsdRevSplitItem 은 액면변경 한 행. 모든 필드 string (KIS docs).
type KsdRevSplitItem struct {
	RecordDate     string `json:"record_date"`       // 기준일
	ShtCd          string `json:"sht_cd"`            // 종목코드
	IsinName       string `json:"isin_name"`         // 종목명
	InterBfFaceAmt string `json:"inter_bf_face_amt"` // 변경전액면가
	InterAfFaceAmt string `json:"inter_af_face_amt"` // 변경후액면가
	TdStopDt       string `json:"td_stop_dt"`        // 매매거래정지기간
	ListDt         string `json:"list_dt"`           // 상장/등록일
}

// InquireKsdRevSplitParams 는 액면변경 조회 파라미터.
type InquireKsdRevSplitParams struct {
	Symbol   string // SHT_CD — 종목코드 (공백=전체)
	Cts      string // CTS — 공백 입력 (default "")
	FromDate string // F_DT — 조회시작일 YYYYMMDD
	ToDate   string // T_DT — 조회종료일 YYYYMMDD
	MarketGb string // MARKET_GB — 0:전체. 빈 값=>"0"
}

// InquireKsdRevSplit 호출.
//
// 한투 docs: docs/api/국내주식/예탁원정보(액면변경).md
// path: /uapi/domestic-stock/v1/ksdinfo/rev-split (HHKDB669105C0)
func (c *Client) InquireKsdRevSplit(ctx context.Context, params InquireKsdRevSplitParams) (*KsdRevSplit, error) {
	marketGb := params.MarketGb
	if marketGb == "" {
		marketGb = "0"
	}
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ksdinfo/rev-split",
		TrID:   "HHKDB669105C0",
		Query: map[string]string{
			"SHT_CD":    params.Symbol,
			"CTS":       params.Cts,
			"F_DT":      params.FromDate,
			"T_DT":      params.ToDate,
			"MARKET_GB": marketGb,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res KsdRevSplit
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse KsdRevSplit: %w", err)
	}
	return &res, nil
}

// KsdForfeit 는 예탁원정보(실권주청약) (HHKDB669109C0) 응답.
//
// 한투 docs: docs/api/국내주식/예탁원정보(실권주청약).md
// path: /uapi/domestic-stock/v1/ksdinfo/forfeit
type KsdForfeit struct {
	Output1 []KsdForfeitItem `json:"output1"`
}

// KsdForfeitItem 은 실권주청약 한 행. 모든 필드 string (KIS docs).
type KsdForfeitItem struct {
	RecordDate   string `json:"record_date"`    // 기준일
	ShtCd        string `json:"sht_cd"`         // 종목코드
	IsinName     string `json:"isin_name"`      // 종목명
	SubscrDt     string `json:"subscr_dt"`      // 청약일
	SubscrPrice  string `json:"subscr_price"`   // 공모가
	SubscrStkQty string `json:"subscr_stk_qty"` // 공모주식수
	RefundDt     string `json:"refund_dt"`      // 환불일
	ListDt       string `json:"list_dt"`        // 상장/등록일
	LeadMgr      string `json:"lead_mgr"`       // 주간사
}

// InquireKsdForfeitParams 는 실권주청약 조회 파라미터.
type InquireKsdForfeitParams struct {
	Symbol   string // SHT_CD — 종목코드 (공백=전체)
	ToDate   string // T_DT — 조회종료일 YYYYMMDD
	FromDate string // F_DT — 조회시작일 YYYYMMDD
	Cts      string // CTS — 공백 입력 (default "")
}

// InquireKsdForfeit 호출.
//
// 한투 docs: docs/api/국내주식/예탁원정보(실권주청약).md
// path: /uapi/domestic-stock/v1/ksdinfo/forfeit (HHKDB669109C0)
func (c *Client) InquireKsdForfeit(ctx context.Context, params InquireKsdForfeitParams) (*KsdForfeit, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ksdinfo/forfeit",
		TrID:   "HHKDB669109C0",
		Query: map[string]string{
			"SHT_CD": params.Symbol,
			"T_DT":   params.ToDate,
			"F_DT":   params.FromDate,
			"CTS":    params.Cts,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res KsdForfeit
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse KsdForfeit: %w", err)
	}
	return &res, nil
}

// KsdMandDeposit 는 예탁원정보(의무보호예수) (HHKDB669110C0) 응답.
//
// 한투 docs: docs/api/국내주식/예탁원정보(의무보호예수).md
// path: /uapi/domestic-stock/v1/ksdinfo/mand-deposit
// ANOMALY: record_date 없음 — depo_date (예치일) 가 날짜 key.
type KsdMandDeposit struct {
	Output1 []KsdMandDepositItem `json:"output1"`
}

// KsdMandDepositItem 은 의무보호예수 한 행. 모든 필드 string (KIS docs).
// record_date 없음 — depo_date 가 날짜 기준 필드.
type KsdMandDepositItem struct {
	ShtCd              string `json:"sht_cd"`                 // 종목코드
	IsinName           string `json:"isin_name"`              // 종목명
	StkQty             string `json:"stk_qty"`                // 주식수
	DepoDate           string `json:"depo_date"`              // 예치일 (날짜 key)
	DepoReason         string `json:"depo_reason"`            // 사유
	TotIssueQtyPerRate string `json:"tot_issue_qty_per_rate"` // 총발행주식수대비비율(%)
}

// InquireKsdMandDepositParams 는 의무보호예수 조회 파라미터.
type InquireKsdMandDepositParams struct {
	ToDate   string // T_DT — 조회종료일 YYYYMMDD
	Symbol   string // SHT_CD — 종목코드 (공백=전체)
	FromDate string // F_DT — 조회시작일 YYYYMMDD
	Cts      string // CTS — 공백 입력 (default "")
}

// InquireKsdMandDeposit 호출.
//
// 한투 docs: docs/api/국내주식/예탁원정보(의무보호예수).md
// path: /uapi/domestic-stock/v1/ksdinfo/mand-deposit (HHKDB669110C0)
func (c *Client) InquireKsdMandDeposit(ctx context.Context, params InquireKsdMandDepositParams) (*KsdMandDeposit, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ksdinfo/mand-deposit",
		TrID:   "HHKDB669110C0",
		Query: map[string]string{
			"T_DT":   params.ToDate,
			"SHT_CD": params.Symbol,
			"F_DT":   params.FromDate,
			"CTS":    params.Cts,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res KsdMandDeposit
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse KsdMandDeposit: %w", err)
	}
	return &res, nil
}
