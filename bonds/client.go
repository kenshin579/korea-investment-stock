package bonds

import (
	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// Client 는 장내채권 (Korean bond) API 클라이언트.
type Client struct {
	http *httpclient.Client
}

// New 는 채권 Client 생성. http 는 root client 의 internal httpclient.
func New(http *httpclient.Client) *Client {
	return &Client{http: http}
}
