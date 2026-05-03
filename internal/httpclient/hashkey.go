package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const hashkeyPath = "/uapi/hashkey"

type hashkeyResp struct {
	Hash string `json:"HASH"`
}

// Hashkey 는 한투의 hashkey 엔드포인트로 body 를 보내고 hash 문자열 반환.
// 주문 등 일부 API 가 요구.
func (c *Client) Hashkey(ctx context.Context, body any) (string, error) {
	httpResp, err := c.resty.R().
		SetContext(ctx).
		SetHeader("appkey", c.cfg.AppKey).
		SetHeader("appsecret", c.cfg.AppSecret).
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetBody(body).
		Execute(http.MethodPost, hashkeyPath)
	if err != nil {
		return "", fmt.Errorf("kis: hashkey: %w", err)
	}
	if httpResp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("kis: hashkey: HTTP %d", httpResp.StatusCode())
	}
	var r hashkeyResp
	if err := json.Unmarshal(httpResp.Body(), &r); err != nil {
		return "", fmt.Errorf("kis: hashkey parse: %w", err)
	}
	return r.Hash, nil
}
