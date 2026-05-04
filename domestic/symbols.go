package domestic

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/kenshin579/korea-investment-stock/internal/krxmaster"
)

// KospiSymbol 은 internal/krxmaster 의 type alias (외부 사용자 노출).
type KospiSymbol = krxmaster.KospiSymbol

// KosdaqSymbol 은 internal/krxmaster 의 type alias (외부 사용자 노출).
type KosdaqSymbol = krxmaster.KosdaqSymbol

const (
	kospiCacheName  = "kospi_code.mst.zip"
	kosdaqCacheName = "kosdaq_code.mst.zip"
)

// FetchKospiSymbols 는 KRX KOSPI 종목 마스터 (kospi_code.mst.zip) 를 다운로드/캐시 후 파싱.
//
// 한투 REST API 가 아니라 KRX 가 공개 다운로드로 제공. 토큰 인증 불필요.
// 마스터 파일은 mastercache 에 디스크 캐시 (default TTL 7일). cp949 + fwf 포맷.
func (c *Client) FetchKospiSymbols(ctx context.Context) ([]KospiSymbol, error) {
	raw, err := c.master.Get(ctx, kospiCacheName, func(ctx context.Context) ([]byte, error) {
		return downloadURL(ctx, krxmaster.KospiURL)
	})
	if err != nil {
		return nil, err
	}
	return krxmaster.ParseKospi(raw)
}

// FetchKosdaqSymbols 는 KRX KOSDAQ 종목 마스터 다운로드/캐시 후 파싱.
func (c *Client) FetchKosdaqSymbols(ctx context.Context) ([]KosdaqSymbol, error) {
	raw, err := c.master.Get(ctx, kosdaqCacheName, func(ctx context.Context) ([]byte, error) {
		return downloadURL(ctx, krxmaster.KosdaqURL)
	})
	if err != nil {
		return nil, err
	}
	return krxmaster.ParseKosdaq(raw)
}

// downloadURL 은 KRX 공개 마스터 파일 단순 GET. 한투 API transport 와 분리 의도로
// http.DefaultClient 사용 (KRX 도메인은 한투 API 가 아니므로 토큰/proxy 정책 다름).
func downloadURL(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("krx download new request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("krx download %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("krx download %s: HTTP %d", url, resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
