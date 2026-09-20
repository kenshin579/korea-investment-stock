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
// KIS 는 숫자를 공백(좌측 패딩) 또는 0 패딩 문자열로 준다. 2026-09-20 실측:
//
//	fix_subscr_pri "       19500" · face_value "000000500" · assign_stk_qty "           0"
//
// 기존 픽스처는 "30000" 처럼 깨끗해서 이 사고를 잡지 못했다.
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
	require.Len(t, res.Output1, 2)

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
