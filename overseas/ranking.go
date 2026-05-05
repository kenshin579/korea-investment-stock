package overseas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// UpdownRate 은 해외주식_상승율_하락율 (HHDFS76290000) 응답.
//
// 한투 docs: docs/api/해외주식/해외주식_상승율_하락율.md
// path: /uapi/overseas-stock/v1/ranking/updown-rate
type UpdownRate struct {
	Output1 UpdownRateSummary `json:"output1"`
	Output2 []UpdownRateItem  `json:"output2"`
}

// UpdownRateSummary 는 응답의 output1 (단일 객체).
type UpdownRateSummary struct {
	Zdiv string `json:"zdiv"`
	Stat string `json:"stat"`
	Crec int64  `json:"crec,string"`
	Trec int64  `json:"trec,string"`
	Nrec int64  `json:"nrec,string"`
}

// UpdownRateItem 은 응답의 output2 한 행.
type UpdownRateItem struct {
	Rsym   string          `json:"rsym"`
	Excd   string          `json:"excd"`
	Symb   string          `json:"symb"`
	Name   string          `json:"name"`
	Last   decimal.Decimal `json:"last"`
	Sign   string          `json:"sign"`
	Diff   decimal.Decimal `json:"diff"`
	Rate   float64         `json:"rate,string"`
	Tvol   int64           `json:"tvol,string"`
	Pask   decimal.Decimal `json:"pask"`
	Pbid   decimal.Decimal `json:"pbid"`
	NBase  decimal.Decimal `json:"n_base"`
	NDiff  decimal.Decimal `json:"n_diff"`
	NRate  float64         `json:"n_rate,string"`
	Rank   int64           `json:"rank,string"`
	Ename  string          `json:"ename"`
	EOrdyn string          `json:"e_ordyn"`
}

// InquireUpdownRateParams 는 해외주식_상승율_하락율 조회 파라미터.
type InquireUpdownRateParams struct {
	Keyb    string // KEYB — NEXT KEY BUFF. 빈 값 default
	Auth    string // AUTH — 사용자권한정보. 빈 값 default
	Excd    string // EXCD — 거래소코드 (NYS/NAS/AMS/HKS/SHS/SZS/HSX/HNX/TSE)
	Gubn    string // GUBN — "0":하락율/"1":상승율
	Nday    string // NDAY — N일전 ("0"~"9"). 빈 값=>"0" (당일)
	VolRang string // VOL_RANG — 거래량조건. 빈 값=>"0" (전체)
}

// InquireUpdownRate 는 해외주식_상승율_하락율 호출.
//
// 한투 docs: docs/api/해외주식/해외주식_상승율_하락율.md
// path: /uapi/overseas-stock/v1/ranking/updown-rate (HHDFS76290000)
func (c *Client) InquireUpdownRate(ctx context.Context, params InquireUpdownRateParams) (*UpdownRate, error) {
	nday := params.Nday
	if nday == "" {
		nday = "0"
	}
	vol := params.VolRang
	if vol == "" {
		vol = "0"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-stock/v1/ranking/updown-rate",
		TrID:   "HHDFS76290000",
		Query: map[string]string{
			"KEYB":     params.Keyb,
			"AUTH":     params.Auth,
			"EXCD":     params.Excd,
			"GUBN":     params.Gubn,
			"NDAY":     nday,
			"VOL_RANG": vol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res UpdownRate
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse UpdownRate: %w", err)
	}
	return &res, nil
}
