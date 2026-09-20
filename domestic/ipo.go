package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// PubOffer 는 예탁원정보(공모주청약일정) (HHKDB669108C0) 응답.
//
// 한투 docs: docs/api/국내주식/예탁원정보(공모주청약일정).md
// path: /uapi/domestic-stock/v1/ksdinfo/pub-offer
//
// IPO 청약일정 list. output1 (Array). 다른 ranking/financial 과 query 키 형식이 다름 (대문자+한글식).
type PubOffer struct {
	Output1 []PubOfferItem `json:"output1"`
}

// PubOfferItem 은 한 IPO 청약 일정 항목.
//
// 날짜 필드 포맷이 섞여 있다(2026-09-20 실측):
//
//	record_date            → "20260917"   (YYYYMMDD)
//	pay_dt/refund_dt/list_dt → "2026/09/22" (YYYY/MM/DD), 없으면 ""
//	subscr_dt              → "2026/09/17 ~ 2026/09/18" (범위 문자열)
//
// list_dt 는 과거 건에서도 자주 비어 있다. 원문 그대로 주고 해석은 소비자에게 맡긴다.
type PubOfferItem struct {
	RecordDate   string          `json:"record_date"`           // 기준일 (YYYYMMDD) = 청약 첫날
	ShtCd        string          `json:"sht_cd"`                // 종목코드. 상장 전에는 임시코드("0035S0")
	IsinName     string          `json:"isin_name"`             // 종목명
	FixSubscrPri decimal.Decimal `json:"fix_subscr_pri"`        // 공모가
	FaceValue    decimal.Decimal `json:"face_value"`            // 액면가
	SubscrDt     string          `json:"subscr_dt"`             // 청약기간 "2026/09/17 ~ 2026/09/18"
	PayDt        string          `json:"pay_dt"`                // 납입일 (YYYY/MM/DD)
	RefundDt     string          `json:"refund_dt"`             // 환불일 (YYYY/MM/DD)
	ListDt       string          `json:"list_dt"`               // 상장일 (YYYY/MM/DD). 자주 빈다
	LeadMgr      string          `json:"lead_mgr"`              // 주간사. 구분자·표기가 제각각
	PubBfCap     int64           `json:"pub_bf_cap,string"`     // 공모전 자본금
	PubAfCap     int64           `json:"pub_af_cap,string"`     // 공모후 자본금
	AssignStkQty int64           `json:"assign_stk_qty,string"` // 당사(한투) 배정물량
}

// padded 는 KIS 가 주는 값을 원문 그대로 받아둔다. 문자열이든 맨 숫자든 상관없다.
//
// raw 구조체가 이 필드들을 그냥 string 으로 선언하면 맨 숫자가 왔을 때 항목
// **전체**가 깨진다 — 이 UnmarshalJSON 이 없애려던 바로 그 실패 방식이다.
type padded string

// UnmarshalJSON 은 따옴표를 벗기기만 한다. 값 해석(trim·콤마 제거·숫자 변환)은
// parsePaddedDecimal/parsePaddedInt64 가 한다.
func (p *padded) UnmarshalJSON(b []byte) error {
	*p = padded(strings.Trim(strings.TrimSpace(string(b)), `"`))
	return nil
}

// UnmarshalJSON 은 KIS 가 주는 패딩된 숫자 문자열을 받아낸다.
//
// 공백 좌측 패딩("       19500")이 문제다 — 표준 JSON 숫자 문법이 아니라서
// decimal.Decimal / int64,string 디코더가 필드 하나에서 실패하면 구조체
// **전체**를 포기한다. 그래서 공모가 하나 때문에 종목명·청약기간·주간사까지
// 전부 날아간다. 0 패딩("000000500")도 실 응답에 섞여 있지만 무해하다 —
// decimal.NewFromString 과 strconv.ParseInt 둘 다 앞자리 0 을 그냥 읽는다.
//
// 그래서 숫자 필드만 원문 그대로(padded) 받아 trim 후 변환한다.
func (p *PubOfferItem) UnmarshalJSON(data []byte) error {
	// type alias 는 raw 가 이 UnmarshalJSON 을 다시 부르지 않게 하는 repo 관용구
	// (domestic/investor.go, domestic/program_trade.go 와 동일 패턴).
	//
	// 반드시 값으로 embed 한다 — *alias 로 바꾸면 "avoid the copy" 처럼 자연스러워
	// 보이지만, encoding/json 이 "cannot set embedded pointer to unexported struct
	// type" 으로 런타임에 깨진다. compile 타임 신호가 없다.
	//
	// 아래 다섯 필드가 alias 의 depth-1 필드와 이름이 겹치는데, 이건 "중복"이
	// 아니라 의도다 — encoding/json 은 이름 충돌을 depth 로 풀어서 이 depth-0
	// padded 필드가 이긴다. PubOfferItem 자체의 json 태그(fix_subscr_pri 등)를
	// 여기 맞춰 "-" 로 지우고 싶어질 수 있는데 그러면 안 된다: 그 태그가 없으면
	// json.Marshal(PubOfferItem{...}) 이 이 다섯 필드를 조용히 빠뜨리게 된다.
	// 태그는 유지한다 — TestPubOfferItem_MarshalKeepsNumbers 참고.
	type alias PubOfferItem
	var raw struct {
		alias
		FixSubscrPri padded `json:"fix_subscr_pri"`
		FaceValue    padded `json:"face_value"`
		PubBfCap     padded `json:"pub_bf_cap"`
		PubAfCap     padded `json:"pub_af_cap"`
		AssignStkQty padded `json:"assign_stk_qty"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*p = PubOfferItem(raw.alias)
	p.FixSubscrPri = parsePaddedDecimal(string(raw.FixSubscrPri))
	p.FaceValue = parsePaddedDecimal(string(raw.FaceValue))
	p.PubBfCap = parsePaddedInt64(string(raw.PubBfCap))
	p.PubAfCap = parsePaddedInt64(string(raw.PubAfCap))
	p.AssignStkQty = parsePaddedInt64(string(raw.AssignStkQty))
	return nil
}

// parsePaddedDecimal 은 공백 좌측 패딩(그리고 흔히 섞여 오는 0 패딩·천단위 콤마)
// 숫자 문자열을 decimal 로 만든다. 빈 값과 숫자가 아닌 값은 0 이다 — KIS 가 "-"
// 나 공백만 주는 칸이 있다. 0 은 이미 "공모가 미확정" 등 자연스러운 값이라
// 에러를 반환하지 않는다 — 여기서 에러를 내면 InquirePubOffer 의
// json.Unmarshal(resp.Raw, &res) 를 타고 올라가 행 하나 때문에 응답 전체가
// 깨진다. 그게 바로 이 파일이 없애려는 실패 방식이다.
//
// domestic/lenient.go 의 decodeLenientNumbers 와 헷갈리지 말 것 — 그건 빈
// 문자열("") 숫자를 0 으로 바꾸는 다른 문제(거래정지 종목)를 다룬다. 공백
// 패딩은 다루지 않는다. 여기 두 helper 가 이 파일만의 별도 대응이다.
func parsePaddedDecimal(s string) decimal.Decimal {
	t := cleanPaddedNumber(s)
	if t == "" {
		return decimal.Zero
	}
	d, err := decimal.NewFromString(t)
	if err != nil {
		return decimal.Zero
	}
	return d
}

// parsePaddedInt64 는 공백 좌측 패딩(그리고 0 패딩·천단위 콤마) 정수 문자열을
// int64 로 만든다. 실패 시 0 처리 이유는 parsePaddedDecimal 과 같다.
func parsePaddedInt64(s string) int64 {
	t := cleanPaddedNumber(s)
	if t == "" {
		return 0
	}
	n, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		return 0
	}
	return n
}

// cleanPaddedNumber 는 좌우 공백을 trim 하고 천단위 콤마를 제거한다.
//
// 콤마를 지우지 않으면 "19,500" 이 파싱 실패로 조용히 0 이 된다 — 나머지 필드는
// 멀쩡히 채워지는 행에서 공모가만 틀려 보이는, 소비자가 알아챌 수 없는 손상이다.
// 콤마 제거는 그 손상을 정상 파싱으로 바꾼다. 진짜 빈 값·"-" 같은 sentinel 만
// 0 fallback 을 타야 한다.
func cleanPaddedNumber(s string) string {
	return strings.ReplaceAll(strings.TrimSpace(s), ",", "")
}

// InquirePubOfferParams 는 공모주청약일정 조회 파라미터.
//
// 다른 ranking 과 query 키 형식이 다름 — KIS docs 그대로 노출 (SHT_CD, CTS, F_DT, T_DT).
type InquirePubOfferParams struct {
	Symbol   string // SHT_CD — 종목코드. 빈 값(공백) = 전체
	Cts      string // CTS — 빈 값(공백) default
	FromDate string // F_DT — 조회일자 From (YYYYMMDD)
	ToDate   string // T_DT — 조회일자 To (YYYYMMDD)
}

// InquirePubOffer 는 예탁원정보(공모주청약일정) 호출.
//
// 한투 docs: docs/api/국내주식/예탁원정보(공모주청약일정).md
// path: /uapi/domestic-stock/v1/ksdinfo/pub-offer (HHKDB669108C0)
//
// 공모주(IPO) 청약일정 list 조회. Symbol 빈 값 시 전체.
func (c *Client) InquirePubOffer(ctx context.Context, params InquirePubOfferParams) (*PubOffer, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ksdinfo/pub-offer",
		TrID:   "HHKDB669108C0",
		Query: map[string]string{
			"SHT_CD": params.Symbol,
			"CTS":    params.Cts,
			"F_DT":   params.FromDate,
			"T_DT":   params.ToDate,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res PubOffer
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse PubOffer: %w", err)
	}
	return &res, nil
}
