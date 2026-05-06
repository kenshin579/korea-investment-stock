package bonds_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/bonds"
)

func TestClient_SearchBondInfo(t *testing.T) {
	client, transport := newTestClient(t)
	transport.RegisterResponder(
		http.MethodGet,
		"=~/uapi/domestic-bond/v1/quotations/search-bond-info",
		httpmock.NewStringResponder(http.StatusOK, loadFixtureString(t, "search_bond_info_success.json")),
	)

	got, err := client.SearchBondInfo(context.Background(), bonds.SearchBondInfoParams{
		Pdno:       "KR103501GCC7",
		PrdtTypeCd: "300",
	})
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, "KR103501GCC7", got.Pdno)
	assert.Equal(t, "국고채권 03500-5103(21-3)", got.KsdBondItemName)
	assert.Equal(t, "3.50", got.KsdRcvgBondSrfcInrt)
	assert.Equal(t, "국채", got.BondClsfKorName)
	assert.Equal(t, "Y", got.ElecSctyYn)
}
