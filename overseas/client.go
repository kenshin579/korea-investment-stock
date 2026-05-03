package overseas

import "github.com/kenshin579/korea-investment-stock/internal/httpclient"

// Client 는 해외주식 API sub-client. Phase 1.5 부터 메서드 추가.
//
// 사용자는 직접 생성하지 않고 kis.Client.Overseas 로 접근.
type Client struct {
	http *httpclient.Client
}

// New 는 internal 용도. root kis.NewClient 가 호출.
func New(http *httpclient.Client) *Client {
	return &Client{http: http}
}
