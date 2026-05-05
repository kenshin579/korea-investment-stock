package overseas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// PriceDetail 은 해외주식_현재가상세 (HHDFS76200200) 응답.
//
// 한투 docs: docs/api/해외주식/해외주식_현재가상세.md
// path: /uapi/overseas-price/v1/quotations/price-detail
type PriceDetail struct {
	Output PriceDetailSnapshot `json:"output"`
}

// PriceDetailSnapshot 은 응답의 output (단일 객체, ~40 fields).
type PriceDetailSnapshot struct {
	Rsym    string          `json:"rsym"`          // 실시간조회종목코드
	Pvol    int64           `json:"pvol,string"`   // 전일거래량
	Open    decimal.Decimal `json:"open"`          // 시가
	High    decimal.Decimal `json:"high"`          // 고가
	Low     decimal.Decimal `json:"low"`           // 저가
	Last    decimal.Decimal `json:"last"`          // 현재가
	Base    decimal.Decimal `json:"base"`          // 전일종가
	Tomv    int64           `json:"tomv,string"`   // 시가총액
	Pamt    int64           `json:"pamt,string"`   // 전일거래대금
	Uplp    decimal.Decimal `json:"uplp"`          // 상한가
	Dnlp    decimal.Decimal `json:"dnlp"`          // 하한가
	H52p    decimal.Decimal `json:"h52p"`          // 52주최고가
	H52d    string          `json:"h52d"`          // 52주최고일자
	L52p    decimal.Decimal `json:"l52p"`          // 52주최저가
	L52d    string          `json:"l52d"`          // 52주최저일자
	Perx    float64         `json:"perx,string"`   // PER
	Pbrx    float64         `json:"pbrx,string"`   // PBR
	Epsx    decimal.Decimal `json:"epsx"`          // EPS
	Bpsx    decimal.Decimal `json:"bpsx"`          // BPS
	Shar    int64           `json:"shar,string"`   // 상장주수
	Mcap    int64           `json:"mcap,string"`   // 자본금
	Curr    string          `json:"curr"`          // 통화
	Zdiv    string          `json:"zdiv"`          // 소수점자리수
	Vnit    string          `json:"vnit"`          // 매매단위
	TXprc   decimal.Decimal `json:"t_xprc"`        // 원환산당일가격
	TXdif   decimal.Decimal `json:"t_xdif"`        // 원환산당일대비
	TXrat   float64         `json:"t_xrat,string"` // 원환산당일등락
	PXprc   decimal.Decimal `json:"p_xprc"`        // 원환산전일가격
	PXdif   decimal.Decimal `json:"p_xdif"`        // 원환산전일대비
	PXrat   float64         `json:"p_xrat,string"` // 원환산전일등락
	TRate   float64         `json:"t_rate,string"` // 당일환율
	PRate   float64         `json:"p_rate,string"` // 전일환율
	TXsgn   string          `json:"t_xsgn"`        // 원환산당일기호
	PXsng   string          `json:"p_xsng"`        // 원환산전일기호
	EOrdyn  string          `json:"e_ordyn"`        // 거래가능여부
	EHogau  string          `json:"e_hogau"`        // 호가단위
	EIcod   string          `json:"e_icod"`         // 업종(섹터)
	EParp   decimal.Decimal `json:"e_parp"`         // 액면가
	Tvol    int64           `json:"tvol,string"`   // 거래량
	Tamt    int64           `json:"tamt,string"`   // 거래대금
	EtypNm  string          `json:"etyp_nm"`        // ETP 분류명
}

// InquirePriceDetailParams 는 해외주식_현재가상세 조회 파라미터.
type InquirePriceDetailParams struct {
	Auth string // AUTH — 사용자권한정보. 빈 값 default
	Excd string // EXCD — 거래소명 (HKS/NYS/NAS/AMS/TSE/SHS/SZS/SHI/SZI/HSX/HNX/BAY/BAQ/BAA)
	Symb string // SYMB — 종목코드
}

// InquirePriceDetail 은 해외주식_현재가상세 호출.
//
// 한투 docs: docs/api/해외주식/해외주식_현재가상세.md
// path: /uapi/overseas-price/v1/quotations/price-detail (HHDFS76200200)
func (c *Client) InquirePriceDetail(ctx context.Context, params InquirePriceDetailParams) (*PriceDetail, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-price/v1/quotations/price-detail",
		TrID:   "HHDFS76200200",
		Query: map[string]string{
			"AUTH": params.Auth,
			"EXCD": params.Excd,
			"SYMB": params.Symb,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res PriceDetail
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse PriceDetail: %w", err)
	}
	return &res, nil
}
