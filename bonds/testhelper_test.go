package bonds_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/kenshin579/korea-investment-stock/bonds"
	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

func newTestClient(t *testing.T) (*bonds.Client, *httpmock.MockTransport) {
	t.Helper()
	transport := httpmock.NewMockTransport()
	hc := httpclient.NewForTest(transport)
	return bonds.New(hc), transport
}

func loadFixtureString(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("loadFixtureString: %v", err)
	}
	return string(b)
}
