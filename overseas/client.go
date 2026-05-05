package overseas

import (
	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
	"github.com/kenshin579/korea-investment-stock/internal/mastercache"
)

// Client 는 해외주식 API sub-client.
//
// 사용자는 직접 생성하지 않고 kis.Client.Overseas 으로 접근.
type Client struct {
	http   *httpclient.Client
	master *mastercache.Cache // 해외 마스터 파일 디스크 캐시 (FetchOverseasSymbols 가 사용)
}

// New 는 internal 용도. root kis.NewClient 가 호출.
func New(http *httpclient.Client, master *mastercache.Cache) *Client {
	return &Client{http: http, master: master}
}
