package kis

import "context"

// IssueAccessToken 은 OAuth 토큰을 발급하고 Bearer 문자열 반환.
// 일반적으로 사용자가 명시 호출할 필요는 없음 — 라이브러리가 자동 발급.
// 디버깅 목적이나 사전 warmup 시에 유용.
func (c *Client) IssueAccessToken(ctx context.Context) (string, error) {
	return c.tokenMgr.Get(ctx)
}

// RefreshAccessToken 은 캐시 무시하고 강제로 새 토큰 발급.
func (c *Client) RefreshAccessToken(ctx context.Context) (string, error) {
	return c.tokenMgr.Refresh(ctx)
}
