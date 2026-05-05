package overseas

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/kenshin579/korea-investment-stock/internal/overseasmaster"
)

// OverseasSymbol 은 internal/overseasmaster 의 type alias (외부 사용자 노출).
type OverseasSymbol = overseasmaster.Symbol

// FetchOverseasSymbols 는 KIS 공개 마스터 (`<market>mst.cod.zip`) 를 다운로드/캐시 후 파싱.
//
// 한투 REST API 가 아니라 KIS 가 공개 다운로드로 제공. 토큰 인증 불필요.
// 마스터 파일은 mastercache 에 디스크 캐시 (default TTL 7일).
//
// market: "nas"(NASDAQ)/"nys"(NYSE)/"ams"(AMEX)/"shs"(상해)/"shi"(상해지수)/
// "szs"(심천)/"szi"(심천지수)/"tse"(도쿄)/"hks"(홍콩)/"hnx"(하노이)/"hsx"(호치민)
func (c *Client) FetchOverseasSymbols(ctx context.Context, market string) ([]OverseasSymbol, error) {
	url, ok := overseasmaster.MarketURLs[market]
	if !ok {
		return nil, fmt.Errorf("overseas: unknown market %q (valid: nas/nys/ams/shs/shi/szs/szi/tse/hks/hnx/hsx)", market)
	}
	cacheName := market + "mst.cod.zip"
	raw, err := c.master.Get(ctx, cacheName, func(ctx context.Context) ([]byte, error) {
		return downloadURL(ctx, url)
	})
	if err != nil {
		return nil, err
	}
	return overseasmaster.Parse(market, raw)
}

// downloadURL 은 KIS 공개 마스터 파일 단순 GET. 한투 API transport 와 분리 의도로
// http.DefaultClient 사용 (KIS 마스터 도메인은 한투 API 가 아니므로 토큰/proxy 정책 다름).
func downloadURL(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("overseas master new request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("overseas master %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("overseas master %s: HTTP %d", url, resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
