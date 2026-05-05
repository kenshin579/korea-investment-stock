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

// OverseasRankingFullSummary 는 output1 5-field tier (시가총액/거래량/거래대금순위).
//
// 해당 메서드: InquireMarketCap, InquireTradeVol, InquireTradePbmn.
type OverseasRankingFullSummary struct {
	Zdiv string `json:"zdiv"`        // 소수점자리수
	Stat string `json:"stat"`        // 거래상태정보
	Crec int64  `json:"crec,string"` // 현재조회종목수
	Trec int64  `json:"trec,string"` // 전체조회종목수
	Nrec int64  `json:"nrec,string"` // RecordCount
}

// MarketCap 은 해외주식_시가총액순위 (HHDFS76350100) 응답.
//
// 한투 docs: docs/api/해외주식/해외주식_시가총액순위.md
// path: /uapi/overseas-stock/v1/ranking/market-cap
type MarketCap struct {
	Output1 OverseasRankingFullSummary `json:"output1"`
	Output2 []MarketCapItem            `json:"output2"`
}

// MarketCapItem 은 시가총액순위 output2 의 한 행 (15 fields).
type MarketCapItem struct {
	Rsym   string          `json:"rsym"`        // 실시간조회심볼
	Excd   string          `json:"excd"`        // 거래소코드
	Symb   string          `json:"symb"`        // 종목코드
	Name   string          `json:"name"`        // 종목명 (한글)
	Last   decimal.Decimal `json:"last"`        // 현재가
	Sign   string          `json:"sign"`        // 기호
	Diff   decimal.Decimal `json:"diff"`        // 대비
	Rate   float64         `json:"rate,string"` // 등락율
	Tvol   int64           `json:"tvol,string"` // 거래량
	Shar   int64           `json:"shar,string"` // 상장주식수
	Tomv   decimal.Decimal `json:"tomv"`        // 시가총액
	Grav   float64         `json:"grav,string"` // 비중
	Rank   int64           `json:"rank,string"` // 순위
	Ename  string          `json:"ename"`       // 영문종목명
	EOrdyn string          `json:"e_ordyn"`     // 매매가능
}

// InquireMarketCapParams 는 해외주식_시가총액순위 조회 파라미터.
type InquireMarketCapParams struct {
	KeyB     string // KEYB — NEXT KEY BUFF. 빈 값 default
	Auth     string // AUTH — 사용자권한정보. 빈 값 default
	ExcdCode string // EXCD — 거래소코드 (NYS/NAS/AMS/HKS/SHS/SZS/HSX/HNX/TSE). 필수
	VolRang  string // VOL_RANG — 거래량조건. 빈 값=>"0" (전체)
}

// InquireMarketCap 은 해외주식_시가총액순위 호출.
//
// 한투 docs: docs/api/해외주식/해외주식_시가총액순위.md
// path: /uapi/overseas-stock/v1/ranking/market-cap (HHDFS76350100)
func (c *Client) InquireMarketCap(ctx context.Context, params InquireMarketCapParams) (*MarketCap, error) {
	excd := params.ExcdCode
	if excd == "" {
		return nil, fmt.Errorf("kis: ExcdCode required for InquireMarketCap")
	}
	vol := params.VolRang
	if vol == "" {
		vol = "0"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-stock/v1/ranking/market-cap",
		TrID:   "HHDFS76350100",
		Query: map[string]string{
			"KEYB":     params.KeyB,
			"AUTH":     params.Auth,
			"EXCD":     excd,
			"VOL_RANG": vol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res MarketCap
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse MarketCap: %w", err)
	}
	return &res, nil
}

// TradeVol 은 해외주식_거래량순위 (HHDFS76310010) 응답.
//
// 한투 docs: docs/api/해외주식/해외주식_거래량순위.md
// path: /uapi/overseas-stock/v1/ranking/trade-vol
type TradeVol struct {
	Output1 OverseasRankingFullSummary `json:"output1"`
	Output2 []TradeVolItem             `json:"output2"`
}

// TradeVolItem 은 거래량순위 output2 의 한 행 (16 fields).
type TradeVolItem struct {
	Rsym   string          `json:"rsym"`          // 실시간조회심볼
	Excd   string          `json:"excd"`          // 거래소코드
	Symb   string          `json:"symb"`          // 종목코드
	Name   string          `json:"name"`          // 종목명 (한글)
	Last   decimal.Decimal `json:"last"`          // 현재가
	Sign   string          `json:"sign"`          // 기호
	Diff   decimal.Decimal `json:"diff"`          // 대비
	Rate   float64         `json:"rate,string"`   // 등락율
	Pask   decimal.Decimal `json:"pask"`          // 매도호가
	Pbid   decimal.Decimal `json:"pbid"`          // 매수호가
	Tvol   int64           `json:"tvol,string"`   // 거래량
	Tamt   int64           `json:"tamt,string"`   // 거래대금
	ATvol  int64           `json:"a_tvol,string"` // 평균거래량
	Rank   int64           `json:"rank,string"`   // 순위
	Ename  string          `json:"ename"`         // 영문종목명
	EOrdyn string          `json:"e_ordyn"`       // 매매가능
}

// InquireTradeVolParams 는 해외주식_거래량순위 조회 파라미터.
type InquireTradeVolParams struct {
	KeyB     string // KEYB — NEXT KEY BUFF. 빈 값 default
	Auth     string // AUTH — 사용자권한정보. 빈 값 default
	ExcdCode string // EXCD — 거래소코드. 필수
	NDay     string // NDAY — N일전: 0(당일),1(2일),2(3일),3(5일),4(10일),5(20일),6(30일),7(60일),8(120일),9(1년). 빈 값=>"0"
	Prc1     string // PRC1 — 현재가 필터범위 1 (가격 ~). 빈 값 OK
	Prc2     string // PRC2 — 현재가 필터범위 2 (~ 가격). 빈 값 OK
	VolRang  string // VOL_RANG — 거래량조건. 빈 값=>"0" (전체)
}

// InquireTradeVol 은 해외주식_거래량순위 호출.
//
// 한투 docs: docs/api/해외주식/해외주식_거래량순위.md
// path: /uapi/overseas-stock/v1/ranking/trade-vol (HHDFS76310010)
func (c *Client) InquireTradeVol(ctx context.Context, params InquireTradeVolParams) (*TradeVol, error) {
	excd := params.ExcdCode
	if excd == "" {
		return nil, fmt.Errorf("kis: ExcdCode required for InquireTradeVol")
	}
	nday := params.NDay
	if nday == "" {
		nday = "0"
	}
	vol := params.VolRang
	if vol == "" {
		vol = "0"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-stock/v1/ranking/trade-vol",
		TrID:   "HHDFS76310010",
		Query: map[string]string{
			"KEYB":     params.KeyB,
			"AUTH":     params.Auth,
			"EXCD":     excd,
			"NDAY":     nday,
			"PRC1":     params.Prc1,
			"PRC2":     params.Prc2,
			"VOL_RANG": vol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res TradeVol
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse TradeVol: %w", err)
	}
	return &res, nil
}

// TradePbmn 은 해외주식_거래대금순위 (HHDFS76320010) 응답.
//
// 한투 docs: docs/api/해외주식/해외주식_거래대금순위.md
// path: /uapi/overseas-stock/v1/ranking/trade-pbmn
type TradePbmn struct {
	Output1 OverseasRankingFullSummary `json:"output1"`
	Output2 []TradePbmnItem            `json:"output2"`
}

// TradePbmnItem 은 거래대금순위 output2 의 한 행 (16 fields).
//
// 주의: 평균 필드명이 a_tamt (평균거래대금) — TradeVolItem 의 a_tvol 과 다름.
type TradePbmnItem struct {
	Rsym   string          `json:"rsym"`          // 실시간조회심볼
	Excd   string          `json:"excd"`          // 거래소코드
	Symb   string          `json:"symb"`          // 종목코드
	Name   string          `json:"name"`          // 종목명 (한글)
	Last   decimal.Decimal `json:"last"`          // 현재가
	Sign   string          `json:"sign"`          // 기호
	Diff   decimal.Decimal `json:"diff"`          // 대비
	Rate   float64         `json:"rate,string"`   // 등락율
	Pask   decimal.Decimal `json:"pask"`          // 매도호가
	Pbid   decimal.Decimal `json:"pbid"`          // 매수호가
	Tvol   int64           `json:"tvol,string"`   // 거래량
	Tamt   int64           `json:"tamt,string"`   // 거래대금
	ATamt  int64           `json:"a_tamt,string"` // 평균거래대금 (a_tvol 아님)
	Rank   int64           `json:"rank,string"`   // 순위
	Ename  string          `json:"ename"`         // 영문종목명
	EOrdyn string          `json:"e_ordyn"`       // 매매가능
}

// InquireTradePbmnParams 는 해외주식_거래대금순위 조회 파라미터.
type InquireTradePbmnParams struct {
	KeyB     string // KEYB — NEXT KEY BUFF. 빈 값 default
	Auth     string // AUTH — 사용자권한정보. 빈 값 default
	ExcdCode string // EXCD — 거래소코드. 필수
	NDay     string // NDAY — N일전: 0(당일),1(2일),2(3일),3(5일),4(10일),5(20일),6(30일),7(60일),8(120일),9(1년). 빈 값=>"0"
	VolRang  string // VOL_RANG — 거래량조건. 빈 값=>"0" (전체)
	Prc1     string // PRC1 — 현재가 필터범위 1 (가격 ~). 빈 값 OK
	Prc2     string // PRC2 — 현재가 필터범위 2 (~ 가격). 빈 값 OK
}

// InquireTradePbmn 은 해외주식_거래대금순위 호출.
//
// 한투 docs: docs/api/해외주식/해외주식_거래대금순위.md
// path: /uapi/overseas-stock/v1/ranking/trade-pbmn (HHDFS76320010)
func (c *Client) InquireTradePbmn(ctx context.Context, params InquireTradePbmnParams) (*TradePbmn, error) {
	excd := params.ExcdCode
	if excd == "" {
		return nil, fmt.Errorf("kis: ExcdCode required for InquireTradePbmn")
	}
	nday := params.NDay
	if nday == "" {
		nday = "0"
	}
	vol := params.VolRang
	if vol == "" {
		vol = "0"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-stock/v1/ranking/trade-pbmn",
		TrID:   "HHDFS76320010",
		Query: map[string]string{
			"KEYB":     params.KeyB,
			"AUTH":     params.Auth,
			"EXCD":     excd,
			"NDAY":     nday,
			"VOL_RANG": vol,
			"PRC1":     params.Prc1,
			"PRC2":     params.Prc2,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res TradePbmn
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse TradePbmn: %w", err)
	}
	return &res, nil
}

// OverseasRankingMinSummary 는 output1 3-field tier (거래량급증/매수체결강도/신고신저).
//
// 해당 메서드: InquireVolumeSurge, InquireVolumePower, InquireNewHighlow.
// crec/trec 없음 — OverseasRankingFullSummary 와 혼용 금지.
type OverseasRankingMinSummary struct {
	Zdiv string `json:"zdiv"`        // 소수점자리수
	Stat string `json:"stat"`        // 거래상태정보
	Nrec int64  `json:"nrec,string"` // RecordCount
}

// VolumeSurge 는 해외주식_거래량급증 (HHDFS76270000) 응답.
//
// 한투 docs: docs/api/해외주식/해외주식_거래량급증.md
// path: /uapi/overseas-stock/v1/ranking/volume-surge
type VolumeSurge struct {
	Output1 OverseasRankingMinSummary `json:"output1"`
	Output2 []VolumeSurgeItem         `json:"output2"`
}

// VolumeSurgeItem 은 거래량급증 output2 의 한 행 (16 fields).
//
// 주의: 종목명 필드가 knam/enam (name/ename 아님).
type VolumeSurgeItem struct {
	Rsym   string          `json:"rsym"`          // 실시간조회심볼
	Excd   string          `json:"excd"`          // 거래소코드
	Symb   string          `json:"symb"`          // 종목코드
	Knam   string          `json:"knam"`          // 종목명 (한글) — name 아님
	Last   decimal.Decimal `json:"last"`          // 현재가
	Sign   string          `json:"sign"`          // 기호
	Diff   decimal.Decimal `json:"diff"`          // 대비
	Rate   float64         `json:"rate,string"`   // 등락율
	Tvol   int64           `json:"tvol,string"`   // 거래량
	Pask   decimal.Decimal `json:"pask"`          // 매도호가
	Pbid   decimal.Decimal `json:"pbid"`          // 매수호가
	NTvol  int64           `json:"n_tvol,string"` // 기준거래량
	NDiff  decimal.Decimal `json:"n_diff"`        // 증가량
	NRate  float64         `json:"n_rate,string"` // 증가율
	Enam   string          `json:"enam"`          // 영문종목명 — ename 아님
	EOrdyn string          `json:"e_ordyn"`       // 매매가능
}

// InquireVolumeSurgeParams 는 해외주식_거래량급증 조회 파라미터.
type InquireVolumeSurgeParams struct {
	KeyB     string // KEYB — NEXT KEY BUFF. 빈 값 default
	Auth     string // AUTH — 사용자권한정보. 빈 값 default
	ExcdCode string // EXCD — 거래소코드. 필수
	MixN     string // MIXN — N분전: 0(1분전),1(2분전),2(3분전),3(5분전),4(10분전),5(15분전),6(20분전),7(30분전),8(60분전),9(120분전). 빈 값=>"0"
	VolRang  string // VOL_RANG — 거래량조건. 빈 값=>"0" (전체)
}

// InquireVolumeSurge 는 해외주식_거래량급증 호출.
//
// 한투 docs: docs/api/해외주식/해외주식_거래량급증.md
// path: /uapi/overseas-stock/v1/ranking/volume-surge (HHDFS76270000)
func (c *Client) InquireVolumeSurge(ctx context.Context, params InquireVolumeSurgeParams) (*VolumeSurge, error) {
	excd := params.ExcdCode
	if excd == "" {
		return nil, fmt.Errorf("kis: ExcdCode required for InquireVolumeSurge")
	}
	mixn := params.MixN
	if mixn == "" {
		mixn = "0"
	}
	vol := params.VolRang
	if vol == "" {
		vol = "0"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-stock/v1/ranking/volume-surge",
		TrID:   "HHDFS76270000",
		Query: map[string]string{
			"KEYB":     params.KeyB,
			"AUTH":     params.Auth,
			"EXCD":     excd,
			"MIXN":     mixn,
			"VOL_RANG": vol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res VolumeSurge
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse VolumeSurge: %w", err)
	}
	return &res, nil
}

// VolumePower 는 해외주식_매수체결강도상위 (HHDFS76280000) 응답.
//
// 한투 docs: docs/api/해외주식/해외주식_매수체결강도상위.md
// path: /uapi/overseas-stock/v1/ranking/volume-power
type VolumePower struct {
	Output1 OverseasRankingMinSummary `json:"output1"`
	Output2 []VolumePowerItem         `json:"output2"`
}

// VolumePowerItem 은 매수체결강도상위 output2 의 한 행 (15 fields).
//
// 주의: 종목명 필드가 knam/enam (name/ename 아님). rank 필드 없음.
type VolumePowerItem struct {
	Rsym   string          `json:"rsym"`        // 실시간조회심볼
	Excd   string          `json:"excd"`        // 거래소코드
	Symb   string          `json:"symb"`        // 종목코드
	Knam   string          `json:"knam"`        // 종목명 (한글) — name 아님
	Last   decimal.Decimal `json:"last"`        // 현재가
	Sign   string          `json:"sign"`        // 기호
	Diff   decimal.Decimal `json:"diff"`        // 대비
	Rate   float64         `json:"rate,string"` // 등락율
	Tvol   int64           `json:"tvol,string"` // 거래량
	Pask   decimal.Decimal `json:"pask"`        // 매도호가
	Pbid   decimal.Decimal `json:"pbid"`        // 매수호가
	Tpow   float64         `json:"tpow,string"` // 당일체결강도
	Powx   float64         `json:"powx,string"` // 체결강도
	Enam   string          `json:"enam"`        // 영문종목명 — ename 아님
	EOrdyn string          `json:"e_ordyn"`     // 매매가능
}

// InquireVolumePowerParams 는 해외주식_매수체결강도상위 조회 파라미터.
//
// 주의: NDAY 파라미터의 설명값이 분(分) 단위 (0=1분전, 1=2분전 … 9=120분전)로
// InquireVolumeSurge 의 MIXN 과 동일한 척도. KIS docs 파라미터명 오류로 보이나
// wire name 은 NDAY 를 그대로 사용.
type InquireVolumePowerParams struct {
	KeyB     string // KEYB — NEXT KEY BUFF. 빈 값 default
	Auth     string // AUTH — 사용자권한정보. 빈 값 default
	ExcdCode string // EXCD — 거래소코드. 필수
	NDay     string // NDAY — N분전 (wire name): 0(1분전),1(2분전),2(3분전),3(5분전),4(10분전),5(15분전),6(20분전),7(30분전),8(60분전),9(120분전). 빈 값=>"0"
	VolRang  string // VOL_RANG — 거래량조건. 빈 값=>"0" (전체)
}

// InquireVolumePower 는 해외주식_매수체결강도상위 호출.
//
// 한투 docs: docs/api/해외주식/해외주식_매수체결강도상위.md
// path: /uapi/overseas-stock/v1/ranking/volume-power (HHDFS76280000)
func (c *Client) InquireVolumePower(ctx context.Context, params InquireVolumePowerParams) (*VolumePower, error) {
	excd := params.ExcdCode
	if excd == "" {
		return nil, fmt.Errorf("kis: ExcdCode required for InquireVolumePower")
	}
	nday := params.NDay
	if nday == "" {
		nday = "0"
	}
	vol := params.VolRang
	if vol == "" {
		vol = "0"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-stock/v1/ranking/volume-power",
		TrID:   "HHDFS76280000",
		Query: map[string]string{
			"KEYB":     params.KeyB,
			"AUTH":     params.Auth,
			"EXCD":     excd,
			"NDAY":     nday,
			"VOL_RANG": vol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res VolumePower
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse VolumePower: %w", err)
	}
	return &res, nil
}
