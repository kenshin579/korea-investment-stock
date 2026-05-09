package futures_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jarcoal/httpmock"
)

// ─── EP1: InquirePrice ───────────────────────────────────────────────────────

func TestClient_InquirePrice(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-futureoption/v1/quotations/inquire-price",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "inquire_price_success.json")),
	)

	got, err := client.InquirePrice(context.Background(), "F", "101S06")
	require.NoError(t, err)
	require.NotNil(t, got)

	// output1 assertions
	assert.Equal(t, "KOSPI200선물 2506", got.Output1.HtsKorIsnm)
	prpr, _ := decimal.NewFromString("362.50")
	assert.True(t, prpr.Equal(got.Output1.FutsPrpr))
	assert.Equal(t, int64(125430), got.Output1.AcmlVol)
	assert.Equal(t, int64(310250), got.Output1.HtsOtstStplQty)
	assert.Equal(t, "20250612", got.Output1.FutsLastTrDate)
	assert.InDelta(t, 0.35, got.Output1.FutsPrdyCtrt, 0.001)

	// output2 assertions
	assert.Equal(t, "코스피200", got.Output2.HtsKorIsnm)
	nmixPrpr, _ := decimal.NewFromString("361.65")
	assert.True(t, nmixPrpr.Equal(got.Output2.BstpNmixPrpr))

	// output3 assertions
	assert.Equal(t, "코스피200", got.Output3.HtsKorIsnm)
}

func TestClient_InquirePrice_InvalidJSON(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-futureoption/v1/quotations/inquire-price",
		httpmock.NewStringResponder(http.StatusOK, `{"rt_cd":"0","msg_cd":"X","msg1":"x","output1":"not-an-object"}`),
	)
	_, err := client.InquirePrice(context.Background(), "F", "101S06")
	require.Error(t, err)
}

// ─── EP2: InquireAskingPrice ──────────────────────────────────────────────────

func TestClient_InquireAskingPrice(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-futureoption/v1/quotations/inquire-asking-price",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "inquire_asking_price_success.json")),
	)

	got, err := client.InquireAskingPrice(context.Background(), "F", "101S06")
	require.NoError(t, err)
	require.NotNil(t, got)

	// output1 assertions
	assert.Equal(t, "KOSPI200선물 2506", got.Output1.HtsKorIsnm)
	prpr, _ := decimal.NewFromString("362.50")
	assert.True(t, prpr.Equal(got.Output1.FutsPrpr))
	assert.Equal(t, int64(125430), got.Output1.AcmlVol)
	assert.Equal(t, "101S06", got.Output1.FutsShrnIscd)

	// output2 array assertions
	require.GreaterOrEqual(t, len(got.Output2), 1)
	q := got.Output2[0]
	askp1, _ := decimal.NewFromString("362.55")
	bidp1, _ := decimal.NewFromString("362.45")
	assert.True(t, askp1.Equal(q.FutsAskp1))
	assert.True(t, bidp1.Equal(q.FutsBidp1))
	assert.Equal(t, int64(250), q.AskpRsqn1)
	assert.Equal(t, int64(980), q.TotalAskpRsqn)
	assert.Equal(t, "141523", q.AsprAcptHour)
}

func TestClient_InquireAskingPrice_InvalidJSON(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-futureoption/v1/quotations/inquire-asking-price",
		httpmock.NewStringResponder(http.StatusOK, `{"rt_cd":"0","msg_cd":"X","msg1":"x","output1":"not-an-object"}`),
	)
	_, err := client.InquireAskingPrice(context.Background(), "F", "101S06")
	require.Error(t, err)
}
