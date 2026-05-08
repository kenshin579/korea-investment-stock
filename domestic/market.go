package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// MarketTime 은 국내선물 영업일조회 (HHMCM000002C0) 응답.
//
// 한투 docs: docs/api/국내주식/국내선물_영업일조회.md
// path: /uapi/domestic-stock/v1/quotations/market-time
//
// 영업일 D-2 ~ D+2 (date1~date5) + today + time(현재시간) + s_time/e_time(장 시작/마감).
// date3 가 영업일 당일.
type MarketTime struct {
	Output1 []MarketTimeItem `json:"output1"`
}

// MarketTimeItem 은 영업일조회 응답의 한 행.
type MarketTimeItem struct {
	Date1 string `json:"date1"`  // 영업일1 (YYYYMMDD)
	Date2 string `json:"date2"`  // 영업일2 (YYYYMMDD)
	Date3 string `json:"date3"`  // 영업일3 — 당일 (YYYYMMDD)
	Date4 string `json:"date4"`  // 영업일4 (YYYYMMDD)
	Date5 string `json:"date5"`  // 영업일5 (YYYYMMDD)
	Today string `json:"today"`  // 오늘일자 (YYYYMMDD)
	Time  string `json:"time"`   // 현재시간 (HHmmss)
	STime string `json:"s_time"` // 장시작시간 (HHmmss)
	ETime string `json:"e_time"` // 장마감시간 (HHmmss)
}

// InquireMarketTime 은 국내선물 영업일조회 호출.
//
// 한투 docs: docs/api/국내주식/국내선물_영업일조회.md
// path: /uapi/domestic-stock/v1/quotations/market-time (HHMCM000002C0)
//
// 파라미터 없음 — path + tr_id 만으로 호출. 모의투자 미지원.
func (c *Client) InquireMarketTime(ctx context.Context) (*MarketTime, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method:   http.MethodGet,
		Path:     "/uapi/domestic-stock/v1/quotations/market-time",
		TrID:     "HHMCM000002C0",
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res MarketTime
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse MarketTime: %w", err)
	}
	return &res, nil
}
