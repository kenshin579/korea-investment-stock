package overseasfutures_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
	"github.com/kenshin579/korea-investment-stock/overseasfutures"
)

func newTestClient(t *testing.T) (*overseasfutures.Client, *httpmock.MockTransport) {
	t.Helper()
	transport := httpmock.NewMockTransport()
	hc := httpclient.NewForTest(transport)
	return overseasfutures.New(hc), transport
}

func loadFixtureString(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("loadFixtureString: %v", err)
	}
	return string(b)
}
