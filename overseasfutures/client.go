package overseasfutures

import "github.com/kenshin579/korea-investment-stock/internal/httpclient"

// Client 는 해외선물옵션 도메인 진입점.
//
// kis.NewClient(...) 가 자동으로 wireInfra 에서 New(httpClient) 호출하여
// kis.Client.OverseasFutures 필드에 주입한다. 직접 인스턴스화 불필요.
type Client struct {
	http *httpclient.Client
}

// New 는 새 OverseasFutures Client 를 생성한다.
func New(http *httpclient.Client) *Client {
	return &Client{http: http}
}
