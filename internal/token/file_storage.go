package token

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// FileStorage 는 토큰을 로컬 파일에 JSON 으로 저장.
type FileStorage struct {
	path string
}

// NewFileStorage 는 지정된 경로에 토큰을 저장하는 FileStorage 생성.
// 파일은 0600 권한으로 작성.
func NewFileStorage(path string) *FileStorage {
	return &FileStorage{path: path}
}

type fileToken struct {
	Value     string    `json:"value"`
	TokenType string    `json:"token_type"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Save 는 토큰을 파일에 저장.
func (s *FileStorage) Save(_ context.Context, token *AccessToken) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0700); err != nil {
		return err
	}
	data, err := json.Marshal(fileToken{
		Value: token.Value, TokenType: token.TokenType, ExpiresAt: token.ExpiresAt,
	})
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0600)
}

// Load 는 파일에서 토큰을 읽음. 파일이 없으면 nil, nil.
func (s *FileStorage) Load(_ context.Context) (*AccessToken, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var ft fileToken
	if err := json.Unmarshal(data, &ft); err != nil {
		return nil, err
	}
	return &AccessToken{
		Value: ft.Value, TokenType: ft.TokenType, ExpiresAt: ft.ExpiresAt,
	}, nil
}

// Clear 는 파일을 삭제. 파일이 없으면 에러 없음.
func (s *FileStorage) Clear(_ context.Context) error {
	err := os.Remove(s.path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	return err
}
