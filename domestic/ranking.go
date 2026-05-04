package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// VolumeRank 는 거래량순위 (FHPST01710000) 응답.
//
// 한투 docs: docs/api/국내주식/거래량순위.md
// path: /uapi/domestic-stock/v1/quotations/volume-rank
//
// 최대 30건 확인 가능, 다음 조회 불가.
type VolumeRank struct {
	Output []VolumeRankItem `json:"Output"` // KIS docs 가 대문자 'O' 표기
}

// VolumeRankItem 은 거래량순위 응답의 한 행.
type VolumeRankItem struct {
	HtsKorIsnm            string          `json:"hts_kor_isnm"`                      // HTS 한글 종목명
	MkscShrnIscd          string          `json:"mksc_shrn_iscd"`                    // 유가증권 단축 종목코드
	DataRank              string          `json:"data_rank"`                         // 데이터 순위
	StckPrpr              decimal.Decimal `json:"stck_prpr"`                         // 주식 현재가
	PrdyVrssSign          string          `json:"prdy_vrss_sign"`                    // 전일 대비 부호
	PrdyVrss              decimal.Decimal `json:"prdy_vrss"`                         // 전일 대비
	PrdyCtrt              float64         `json:"prdy_ctrt,string"`                  // 전일 대비율
	AcmlVol               int64           `json:"acml_vol,string"`                   // 누적 거래량
	PrdyVol               int64           `json:"prdy_vol,string"`                   // 전일 거래량
	LstnStcn              int64           `json:"lstn_stcn,string"`                  // 상장 주수
	AvrgVol               int64           `json:"avrg_vol,string"`                   // 평균 거래량
	NBefRClprVrssPrprRate float64         `json:"n_befr_clpr_vrss_prpr_rate,string"` // N일전종가대비현재가대비율
	VolInrt               float64         `json:"vol_inrt,string"`                   // 거래량 증가율
	VolTnrt               float64         `json:"vol_tnrt,string"`                   // 거래량 회전율
	NdayVolTnrt           float64         `json:"nday_vol_tnrt,string"`              // N일 거래량 회전율
	AvrgTrPbmn            int64           `json:"avrg_tr_pbmn,string"`               // 평균 거래 대금
	TrPbmnTnrt            float64         `json:"tr_pbmn_tnrt,string"`               // 거래대금 회전율
	NdayTrPbmnTnrt        float64         `json:"nday_tr_pbmn_tnrt,string"`          // N일 거래대금 회전율
	AcmlTrPbmn            int64           `json:"acml_tr_pbmn,string"`               // 누적 거래 대금
}

// InquireVolumeRankParams 는 거래량순위 조회 파라미터.
//
// 필수: InputISCD (종목코드 또는 "0000" 전체).
// 나머지는 zero-value 시 sensible default 사용.
type InquireVolumeRankParams struct {
	MarketCode    string // FID_COND_MRKT_DIV_CODE — "J":KRX, "NX":NXT. 빈 값=>"J"
	ScreenCode    string // FID_COND_SCR_DIV_CODE — Unique key. 빈 값=>"20171"
	InputISCD     string // FID_INPUT_ISCD — 필수, "0000"(전체) 또는 업종코드
	DivCode       string // FID_DIV_CLS_CODE — "0":전체, "1":보통주, "2":우선주. 빈 값=>"0"
	BelongCode    string // FID_BLNG_CLS_CODE — "0":평균거래량, "1":거래증가율, "2":평균거래회전율, "3":거래금액순, "4":평균거래금액회전율. 빈 값=>"0"
	TargetCode    string // FID_TRGT_CLS_CODE — 9자리 (증거금/신용보증금 비율). 빈 값=>"111111111"
	TargetExclude string // FID_TRGT_EXLS_CLS_CODE — 10자리 (제외 항목). 빈 값=>"0000000000"
	InputPrice1   string // FID_INPUT_PRICE_1 — 가격 ~. 빈 값 OK
	InputPrice2   string // FID_INPUT_PRICE_2 — ~ 가격. 빈 값 OK
	VolCount      string // FID_VOL_CNT — 거래량 ~. 빈 값 OK
	InputDate1    string // FID_INPUT_DATE_1 — 빈 값 OK
}

// InquireVolumeRank 는 거래량순위 호출.
//
// 한투 docs: docs/api/국내주식/거래량순위.md
// path: /uapi/domestic-stock/v1/quotations/volume-rank (FHPST01710000)
func (c *Client) InquireVolumeRank(ctx context.Context, params InquireVolumeRankParams) (*VolumeRank, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scr := params.ScreenCode
	if scr == "" {
		scr = "20171"
	}
	div := params.DivCode
	if div == "" {
		div = "0"
	}
	belong := params.BelongCode
	if belong == "" {
		belong = "0"
	}
	tgt := params.TargetCode
	if tgt == "" {
		tgt = "111111111"
	}
	tgtExcl := params.TargetExclude
	if tgtExcl == "" {
		tgtExcl = "0000000000"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/volume-rank",
		TrID:   "FHPST01710000",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_COND_SCR_DIV_CODE":  scr,
			"FID_INPUT_ISCD":         params.InputISCD,
			"FID_DIV_CLS_CODE":       div,
			"FID_BLNG_CLS_CODE":      belong,
			"FID_TRGT_CLS_CODE":      tgt,
			"FID_TRGT_EXLS_CLS_CODE": tgtExcl,
			"FID_INPUT_PRICE_1":      params.InputPrice1,
			"FID_INPUT_PRICE_2":      params.InputPrice2,
			"FID_VOL_CNT":            params.VolCount,
			"FID_INPUT_DATE_1":       params.InputDate1,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	// 'Output' 키가 대문자 (KIS docs 명시) — resp.Raw 로 unmarshal.
	var rank VolumeRank
	if err := json.Unmarshal(resp.Raw, &rank); err != nil {
		return nil, fmt.Errorf("kis: parse VolumeRank: %w", err)
	}
	return &rank, nil
}
