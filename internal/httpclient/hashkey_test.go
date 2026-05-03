package httpclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Hashkey(t *testing.T) {
	tm := &stubTokenMgr{bearer: "Bearer T"}
	c := newTestClient(t, tm)
	httpmock.RegisterResponder(http.MethodPost, "=~/uapi/hashkey",
		httpmock.NewStringResponder(200, `{"HASH":"abcdef","BODY":{}}`))

	hk, err := c.Hashkey(context.Background(), map[string]string{"k": "v"})
	require.NoError(t, err)
	assert.Equal(t, "abcdef", hk)
}
