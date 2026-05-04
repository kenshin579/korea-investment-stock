package domestic_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_SearchInfo(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/search-info`,
		httpmock.NewStringResponder(200, loadFixtureString(t, "product_info_success.json")),
	)

	c := newTestClient(t)
	info, err := c.SearchInfo(context.Background(), "005930", "300")
	require.NoError(t, err)
	require.NotNil(t, info)

	assert.Equal(t, "005930", info.Pdno)
	assert.Equal(t, "300", info.PrdtTypeCd)
	assert.Equal(t, "삼성전자", info.PrdtName)
	assert.Equal(t, "주권", info.PrdtClsfName)
	assert.Equal(t, "KR7005930003", info.StdPdno)
}

func TestClient_SearchStockInfo(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/search-stock-info`,
		httpmock.NewStringResponder(200, loadFixtureString(t, "stock_info_success.json")),
	)

	c := newTestClient(t)
	info, err := c.SearchStockInfo(context.Background(), "005930", "300")
	require.NoError(t, err)
	require.NotNil(t, info)

	assert.Equal(t, "005930", info.Pdno)
	assert.Equal(t, "300", info.PrdtTypeCd)
	assert.Equal(t, "STK", info.MketIdCd)
	assert.Equal(t, "ST", info.ScrtGrpIdCd)
	assert.Equal(t, "Y", info.Kospi200ItemYn)
	assert.Equal(t, "삼성전자", info.PrdtName)
	assert.Equal(t, "1231", info.SetlMmdd)
	assert.Equal(t, "N", info.AdmnItemYn)
	assert.Equal(t, "N", info.TrStopYn)
}
