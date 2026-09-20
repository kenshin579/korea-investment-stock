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

// pubOfferFields 는 PubOfferItem 과 같은 필드를 갖되 UnmarshalJSON 을 물려받지 않는
// 별칭이다. 아래 UnmarshalJSON 이 자기 자신을 무한히 부르지 않게 한다.
//
// 위 다섯 숫자 필드의 json 태그(fix_subscr_pri 등)를 "-" 로 지우고 싶어질 수 있는데
// (아래 raw 구조체가 어차피 같은 이름의 string 필드를 depth 0 에서 다시 선언하니
// "중복"으로 보인다) 그러면 안 된다. encoding/json 은 이름 충돌을 depth 로 푼다 —
// raw 의 depth-0 string 필드가 이겨서 역직렬화는 원래대로 그쪽을 타지만, 원래 태그가
// 없으면 json.Marshal(PubOfferItem{...}) 이 이 다섯 필드를 조용히 빠뜨리게 된다.
// 이 라이브러리는 태그 버전으로 배포되는 공개 API 라 그 회귀는 여기서가 아니라
// 쓰는 쪽에서만 드러난다. 태그는 유지한다 — TestPubOfferItem_MarshalKeepsNumbers 참고.
//
// int64 세 필드(pub_bf_cap 등)는 ",string" 옵션도 그대로 유지한다 — 이게 빠지면
// Marshal 이 맨 숫자(bare number)를 내보내는데, 위 raw 구조체는 그 자리를 string
// 으로 받으므로 라운드트립이 "cannot unmarshal number into ... of type string" 으로
// 깨진다. decimal.Decimal 두 필드는 기본이 이미 따옴표 붙은 문자열이라 옵션이 없다.
type pubOfferFields PubOfferItem

// UnmarshalJSON 은 KIS 가 주는 패딩된 숫자 문자열을 받아낸다.
//
// KIS 는 숫자를 고정폭 문자열로 준다 — 공백 좌측 패딩("       19500") 또는
// 0 패딩("000000500"). 기본 디코더는 여기서 필드 하나가 아니라 **구조체 전체**를
// 포기하므로, 공모가 하나 때문에 종목명·청약기간·주간사까지 전부 날아간다.
// 그래서 숫자 필드만 문자열로 받아 trim 후 변환한다.
func (p *PubOfferItem) UnmarshalJSON(data []byte) error {
	var raw struct {
		pubOfferFields
		FixSubscrPri string `json:"fix_subscr_pri"`
		FaceValue    string `json:"face_value"`
		PubBfCap     string `json:"pub_bf_cap"`
		PubAfCap     string `json:"pub_af_cap"`
		AssignStkQty string `json:"assign_stk_qty"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*p = PubOfferItem(raw.pubOfferFields)
	p.FixSubscrPri = paddedDecimal(raw.FixSubscrPri)
	p.FaceValue = paddedDecimal(raw.FaceValue)
	p.PubBfCap = paddedInt64(raw.PubBfCap)
	p.PubAfCap = paddedInt64(raw.PubAfCap)
	p.AssignStkQty = paddedInt64(raw.AssignStkQty)
	return nil
}

// paddedDecimal 은 공백·0 패딩 숫자 문자열을 decimal 로 만든다.
// 빈 값과 숫자가 아닌 값은 0 이다 — KIS 가 "-" 나 공백만 주는 칸이 있다.
func paddedDecimal(s string) decimal.Decimal {
	t := strings.TrimSpace(s)
	if t == "" {
		return decimal.Zero
	}
	d, err := decimal.NewFromString(t)
	if err != nil {
		return decimal.Zero
	}
	return d
}

// paddedInt64 는 공백·0 패딩 정수 문자열을 int64 로 만든다.
func paddedInt64(s string) int64 {
	t := strings.TrimSpace(s)
	if t == "" {
		return 0
	}
	n, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		return 0
	}
	return n
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
