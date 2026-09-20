package domestic_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/domestic"
)

func TestClient_InquirePubOffer(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ksdinfo/pub-offer`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "pub_offer_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquirePubOffer(context.Background(), domestic.InquirePubOfferParams{
		FromDate: "20260501",
		ToDate:   "20260531",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// query 키가 대문자 + 한글식
	assert.Equal(t, "", capturedQuery.Get("SHT_CD"))
	assert.Equal(t, "", capturedQuery.Get("CTS"))
	assert.Equal(t, "20260501", capturedQuery.Get("F_DT"))
	assert.Equal(t, "20260531", capturedQuery.Get("T_DT"))

	require.Len(t, res.Output1, 2)
	assert.Equal(t, "999998", res.Output1[0].ShtCd)
	assert.Equal(t, "샘플바이오", res.Output1[0].IsinName)
	assert.Equal(t, decimal.NewFromInt(30000), res.Output1[0].FixSubscrPri)
	assert.Equal(t, decimal.NewFromInt(100), res.Output1[0].FaceValue)
	assert.Equal(t, "20260505 ~ 20260506", res.Output1[0].SubscrDt)
	assert.Equal(t, "한국투자증권", res.Output1[0].LeadMgr)
}

// TestClient_InquirePubOffer_PaddedNumbers 는 실 API 응답 모양을 그대로 넣는다.
// KIS 는 숫자를 공백 좌측 패딩 문자열로 준다 — 이게 파싱을 깨뜨리는 원인이다.
// 0 패딩("000000500")도 실 응답에 섞여 있지만 무해해서(leading zero 는
// decimal/ParseInt 둘 다 그냥 읽는다) 사고의 원인이 아니다. 2026-09-20 실측:
//
//	fix_subscr_pri "       19500" · pub_bf_cap "     5063824"
//
// 기존 픽스처(pub_offer_success.json)는 "30000" 처럼 패딩 없이 깨끗해서 이
// 사고를 잡지 못했다.
//
// 픽스처 세 행 모두 2026-09-20 실 API 응답을 그대로 옮긴 것이며 숫자를 손대지
// 않았다. 1·2번 행은 pub_bf_cap(공모전 자본금) > pub_af_cap(공모후 자본금) 인데
// "공모전 < 공모후" 라는 필드명 상 직관과 반대다 — 실제로 그렇게 내려온 값이라
// 고치지 않는다. pub_af_cap 이 "공모후 총자본금"이 아니라 "이번 공모분 자본금"
// 을 뜻할 가능성이 있어 보이지만 확인되지 않았으니 필드 주석을 추측으로
// 바꾸지는 않는다.
func TestClient_InquirePubOffer_PaddedNumbers(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/ksdinfo/pub-offer`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, loadFixtureString(t, "pub_offer_padded.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquirePubOffer(context.Background(), domestic.InquirePubOfferParams{
		FromDate: "20260901", ToDate: "20260930",
	})
	require.NoError(t, err)
	require.Len(t, res.Output1, 3)

	first := res.Output1[0]
	assert.Equal(t, "468670", first.ShtCd)
	assert.Equal(t, decimal.NewFromInt(19500), first.FixSubscrPri)
	assert.Equal(t, decimal.NewFromInt(500), first.FaceValue)
	assert.Equal(t, int64(5063824), first.PubBfCap)
	assert.Equal(t, int64(150000), first.PubAfCap)
	assert.Equal(t, int64(0), first.AssignStkQty)
	// 상장일이 빈 종목이 실제로 많다 — 빈 문자열 그대로 통과해야 한다.
	assert.Equal(t, "", first.ListDt)
	// 범위 문자열은 손대지 않고 원문 그대로 준다. 쪼개는 것은 소비자 몫이다.
	assert.Equal(t, "2026/09/17 ~ 2026/09/18", first.SubscrDt)

	// 상장 전 임시 종목코드는 6자리 숫자가 아니다.
	assert.Equal(t, "0035S0", res.Output1[1].ShtCd)
	assert.Equal(t, decimal.NewFromInt(18000), res.Output1[1].FixSubscrPri)

	// 세 번째 행은 배정물량이 0 이 아니다 — 앞의 두 행은 둘 다
	// assign_stk_qty 가 "           0" 이라 helper 가 무조건 0 을 반환해도
	// 통과했다. 이 행이 그 맹점을 막는다.
	assert.Equal(t, int64(1000000), res.Output1[2].AssignStkQty)
}

// TestPubOfferItem_MarshalKeepsNumbers 는 UnmarshalJSON 을 붙이면서 직렬화를
// 깨뜨리지 않았는지 본다. 숫자 필드에 json:"-" 를 달면 역직렬화는 멀쩡한데
// json.Marshal 이 그 필드들을 조용히 빠뜨린다 — 라이브러리를 쓰는 쪽에서만
// 드러나는 종류의 회귀라 여기서 붙잡는다.
func TestPubOfferItem_MarshalKeepsNumbers(t *testing.T) {
	item := domestic.PubOfferItem{
		ShtCd:        "468670",
		FixSubscrPri: decimal.NewFromInt(19500),
		FaceValue:    decimal.NewFromInt(500),
		PubBfCap:     5063824,
		PubAfCap:     150000,
		AssignStkQty: 0,
	}
	b, err := json.Marshal(item)
	require.NoError(t, err)

	var back domestic.PubOfferItem
	require.NoError(t, json.Unmarshal(b, &back))
	assert.Equal(t, item.FixSubscrPri, back.FixSubscrPri)
	assert.Equal(t, item.FaceValue, back.FaceValue)
	assert.Equal(t, item.PubBfCap, back.PubBfCap)
	assert.Equal(t, item.PubAfCap, back.PubAfCap)
	assert.Equal(t, item.AssignStkQty, back.AssignStkQty)
}

// TestPubOfferItem_UnmarshalAcceptsBareNumbers 는 fix_subscr_pri/pub_bf_cap 등이
// 패딩된 문자열이 아니라 맨 JSON 숫자로 오는 경우도 항목 전체를 깨뜨리지 않아야
// 함을 본다. UnmarshalJSON 의 raw 구조체가 이 필드들을 그냥 string 으로 선언하면
// 맨 숫자에서 "cannot unmarshal number into ... of type string" 으로 항목 전체가
// 깨진다 — 이 파일이 원래 없애려던 바로 그 실패 방식이라 회귀를 막는다.
func TestPubOfferItem_UnmarshalAcceptsBareNumbers(t *testing.T) {
	var item domestic.PubOfferItem
	raw := []byte(`{
		"sht_cd": "468670",
		"fix_subscr_pri": 19500,
		"pub_bf_cap": 5063824
	}`)
	require.NoError(t, json.Unmarshal(raw, &item))
	assert.Equal(t, "468670", item.ShtCd)
	assert.Equal(t, decimal.NewFromInt(19500), item.FixSubscrPri)
	assert.Equal(t, int64(5063824), item.PubBfCap)
}
