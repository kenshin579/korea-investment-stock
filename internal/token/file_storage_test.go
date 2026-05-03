package token

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStorage_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "token.json")
	s := NewFileStorage(path)
	ctx := context.Background()

	tok := &AccessToken{
		Value:     "abc123",
		TokenType: "Bearer",
		ExpiresAt: time.Now().Add(1 * time.Hour).Round(time.Second),
	}
	require.NoError(t, s.Save(ctx, tok))

	loaded, err := s.Load(ctx)
	require.NoError(t, err)
	require.NotNil(t, loaded)
	assert.Equal(t, "abc123", loaded.Value)
	assert.Equal(t, "Bearer", loaded.TokenType)
	assert.True(t, loaded.ExpiresAt.Equal(tok.ExpiresAt))
}

func TestFileStorage_LoadEmpty(t *testing.T) {
	dir := t.TempDir()
	s := NewFileStorage(filepath.Join(dir, "token.json"))
	loaded, err := s.Load(context.Background())
	require.NoError(t, err)
	assert.Nil(t, loaded)
}

func TestFileStorage_Clear(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "token.json")
	s := NewFileStorage(path)
	ctx := context.Background()

	require.NoError(t, s.Save(ctx, &AccessToken{Value: "x", ExpiresAt: time.Now().Add(time.Hour)}))
	require.NoError(t, s.Clear(ctx))

	_, err := os.Stat(path)
	assert.True(t, os.IsNotExist(err))
}

func TestFileStorage_Permissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "token.json")
	s := NewFileStorage(path)
	require.NoError(t, s.Save(context.Background(), &AccessToken{
		Value: "x", ExpiresAt: time.Now().Add(time.Hour),
	}))

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
}
