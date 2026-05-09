package websocket

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFrame_RealtimeData(t *testing.T) {
	raw := "0|H0STCNT0|001|005930^123929^73100^2"
	f, err := parseFrame(raw)
	require.NoError(t, err)
	assert.Equal(t, frameKindRealtime, f.Kind)
	assert.False(t, f.Encrypted)
	assert.Equal(t, "H0STCNT0", f.TrID)
	assert.Equal(t, 1, f.Count)
	assert.Equal(t, []string{"005930", "123929", "73100", "2"}, f.Fields)
}

func TestParseFrame_Encrypted(t *testing.T) {
	raw := "1|H0STCNI0|001|encrypted-payload"
	f, err := parseFrame(raw)
	require.NoError(t, err)
	assert.Equal(t, frameKindRealtime, f.Kind)
	assert.True(t, f.Encrypted)
}

func TestParseFrame_RealtimePaging(t *testing.T) {
	// count=2 인 페이징
	raw := "0|H0STCNT0|002|f1^f2^f3^f4^f1b^f2b^f3b^f4b"
	f, err := parseFrame(raw)
	require.NoError(t, err)
	assert.Equal(t, 2, f.Count)
	assert.Len(t, f.Fields, 8)
}

func TestParseFrame_JSON_SubscribeSuccess(t *testing.T) {
	raw := `{"header":{"tr_id":"H0STCNT0"},"body":{"rt_cd":"0","msg_cd":"OPSP0000","msg1":"SUBSCRIBE SUCCESS","output":{"iv":"abc","key":"def"}}}`
	f, err := parseFrame(raw)
	require.NoError(t, err)
	assert.Equal(t, frameKindJSON, f.Kind)
	assert.Equal(t, "0", f.JSON.Body.RtCd)
	assert.Equal(t, "OPSP0000", f.JSON.Body.MsgCd)
	assert.Contains(t, f.JSON.Body.Msg1, "SUBSCRIBE")
}

func TestParseFrame_JSON_PingPong(t *testing.T) {
	raw := `{"header":{"tr_id":"PINGPONG"}}`
	f, err := parseFrame(raw)
	require.NoError(t, err)
	assert.Equal(t, frameKindPingPong, f.Kind)
}

func TestParseFrame_Invalid(t *testing.T) {
	_, err := parseFrame("garbage-data-no-pipe-no-brace")
	assert.True(t, errors.Is(err, ErrWSInvalidFrame))
}

func TestChunkFields(t *testing.T) {
	// 1 frame = 2 chunks (count=2), 4 fields/chunk
	chunks := chunkFields([]string{"a", "b", "c", "d", "e", "f", "g", "h"}, 2, 4)
	assert.Equal(t, [][]string{{"a", "b", "c", "d"}, {"e", "f", "g", "h"}}, chunks)
}

func TestChunkFields_MismatchedLength(t *testing.T) {
	// 7 fields, count=2, fieldsPerChunk=4 → mismatch
	_, err := chunkFieldsErr([]string{"a", "b", "c", "d", "e", "f", "g"}, 2, 4)
	assert.True(t, errors.Is(err, ErrWSInvalidFrame))
}
