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
	NBefrClprVrssPrprRate float64         `json:"n_befr_clpr_vrss_prpr_rate,string"` // N일전종가대비현재가대비율
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

// Fluctuation 은 등락률 순위 (FHPST01700000) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_등락률_순위.md
// path: /uapi/domestic-stock/v1/ranking/fluctuation
type Fluctuation struct {
	Output []FluctuationItem `json:"output"`
}

// FluctuationItem 은 등락률 순위 응답의 한 행.
type FluctuationItem struct {
	StckShrnIscd             string          `json:"stck_shrn_iscd"`                       // 주식 단축 종목코드
	DataRank                 string          `json:"data_rank"`                            // 데이터 순위
	HtsKorIsnm               string          `json:"hts_kor_isnm"`                         // HTS 한글 종목명
	StckPrpr                 decimal.Decimal `json:"stck_prpr"`                            // 주식 현재가
	PrdyVrss                 decimal.Decimal `json:"prdy_vrss"`                            // 전일 대비
	PrdyVrssSign             string          `json:"prdy_vrss_sign"`                       // 전일 대비 부호
	PrdyCtrt                 float64         `json:"prdy_ctrt,string"`                     // 전일 대비율
	AcmlVol                  int64           `json:"acml_vol,string"`                      // 누적 거래량
	StckHgpr                 decimal.Decimal `json:"stck_hgpr"`                            // 주식 최고가
	HgprHour                 string          `json:"hgpr_hour"`                            // 최고가 시간 (HHMMSS)
	AcmlHgprDate             string          `json:"acml_hgpr_date"`                       // 누적 최고가 일자
	StckLwpr                 decimal.Decimal `json:"stck_lwpr"`                            // 주식 최저가
	LwprHour                 string          `json:"lwpr_hour"`                            // 최저가 시간
	AcmlLwprDate             string          `json:"acml_lwpr_date"`                       // 누적 최저가 일자
	LwprVrssPrprRate         float64         `json:"lwpr_vrss_prpr_rate,string"`           // 최저가 대비 현재가 비율
	DsgtDateClprVrssPrprRate float64         `json:"dsgt_date_clpr_vrss_prpr_rate,string"` // 지정 일자 종가 대비 현재가 비율
	CnntAscnDynu             int64           `json:"cnnt_ascn_dynu,string"`                // 연속 상승 일수
	HgprVrssPrprRate         float64         `json:"hgpr_vrss_prpr_rate,string"`           // 최고가 대비 현재가 비율
	CnntDownDynu             int64           `json:"cnnt_down_dynu,string"`                // 연속 하락 일수
	OprcVrssPrprSign         string          `json:"oprc_vrss_prpr_sign"`                  // 시가 대비 현재가 부호
	OprcVrssPrpr             decimal.Decimal `json:"oprc_vrss_prpr"`                       // 시가 대비 현재가
	OprcVrssPrprRate         float64         `json:"oprc_vrss_prpr_rate,string"`           // 시가 대비 현재가 비율
	PrdRsfl                  decimal.Decimal `json:"prd_rsfl"`                             // 기간 등락
	PrdRsflRate              float64         `json:"prd_rsfl_rate,string"`                 // 기간 등락 비율
}

// InquireFluctuationParams 는 등락률 순위 조회 파라미터.
type InquireFluctuationParams struct {
	RsflRate2     string // fid_rsfl_rate2 — 등락 비율2 (~ 비율). 빈 값 OK
	MarketCode    string // fid_cond_mrkt_div_code — "J":KRX, "NX":NXT. 빈 값=>"J"
	ScreenCode    string // fid_cond_scr_div_code — Unique key. 빈 값=>"20170"
	InputISCD     string // fid_input_iscd — 필수, "0000"(전체)/"0001"(코스피)/"1001"(코스닥)/"2001"(코스피200)
	SortCode      string // fid_rank_sort_cls_code — "0":상승율순, "1":하락율순, "2":시가대비상승, "3":시가대비하락, "4":변동율. 빈 값=>"0"
	InputCnt1     string // fid_input_cnt_1 — "0":전체, 또는 누적일수. 빈 값=>"0"
	PriceCode     string // fid_prc_cls_code — 가격 구분. 빈 값=>"0"
	InputPrice1   string // fid_input_price_1
	InputPrice2   string // fid_input_price_2
	VolCount      string // fid_vol_cnt
	TargetCode    string // fid_trgt_cls_code. 빈 값=>"0"
	TargetExclude string // fid_trgt_exls_cls_code. 빈 값=>"0"
	DivCode       string // fid_div_cls_code. 빈 값=>"0"
	RsflRate1     string // fid_rsfl_rate1 — 등락 비율1 (비율 ~). 빈 값 OK
}

// InquireFluctuation 는 등락률 순위 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_등락률_순위.md
// path: /uapi/domestic-stock/v1/ranking/fluctuation (FHPST01700000)
func (c *Client) InquireFluctuation(ctx context.Context, params InquireFluctuationParams) (*Fluctuation, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scr := params.ScreenCode
	if scr == "" {
		scr = "20170"
	}
	sort := params.SortCode
	if sort == "" {
		sort = "0"
	}
	cnt := params.InputCnt1
	if cnt == "" {
		cnt = "0"
	}
	prc := params.PriceCode
	if prc == "" {
		prc = "0"
	}
	tgt := params.TargetCode
	if tgt == "" {
		tgt = "0"
	}
	tgtExcl := params.TargetExclude
	if tgtExcl == "" {
		tgtExcl = "0"
	}
	div := params.DivCode
	if div == "" {
		div = "0"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ranking/fluctuation",
		TrID:   "FHPST01700000",
		Query: map[string]string{
			"fid_rsfl_rate2":         params.RsflRate2,
			"fid_cond_mrkt_div_code": market,
			"fid_cond_scr_div_code":  scr,
			"fid_input_iscd":         params.InputISCD,
			"fid_rank_sort_cls_code": sort,
			"fid_input_cnt_1":        cnt,
			"fid_prc_cls_code":       prc,
			"fid_input_price_1":      params.InputPrice1,
			"fid_input_price_2":      params.InputPrice2,
			"fid_vol_cnt":            params.VolCount,
			"fid_trgt_cls_code":      tgt,
			"fid_trgt_exls_cls_code": tgtExcl,
			"fid_div_cls_code":       div,
			"fid_rsfl_rate1":         params.RsflRate1,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res Fluctuation
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse Fluctuation: %w", err)
	}
	return &res, nil
}

// MarketCap 은 시가총액 상위 (FHPST01740000) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_시가총액_상위.md
// path: /uapi/domestic-stock/v1/ranking/market-cap
type MarketCap struct {
	Output []MarketCapItem `json:"output"`
}

// MarketCapItem 은 시가총액 상위 응답의 한 행.
type MarketCapItem struct {
	MkscShrnIscd     string          `json:"mksc_shrn_iscd"`             // 유가증권 단축 종목코드
	DataRank         string          `json:"data_rank"`                  // 데이터 순위
	HtsKorIsnm       string          `json:"hts_kor_isnm"`               // HTS 한글 종목명
	StckPrpr         decimal.Decimal `json:"stck_prpr"`                  // 주식 현재가
	PrdyVrss         decimal.Decimal `json:"prdy_vrss"`                  // 전일 대비
	PrdyVrssSign     string          `json:"prdy_vrss_sign"`             // 전일 대비 부호
	PrdyCtrt         float64         `json:"prdy_ctrt,string"`           // 전일 대비율
	AcmlVol          int64           `json:"acml_vol,string"`            // 누적 거래량
	LstnStcn         int64           `json:"lstn_stcn,string"`           // 상장 주수
	StckAvls         int64           `json:"stck_avls,string"`           // 시가 총액
	MrktWholAvlsRlim float64         `json:"mrkt_whol_avls_rlim,string"` // 시장 전체 시가총액 비중
}

// InquireMarketCapParams 는 시가총액 상위 조회 파라미터.
type InquireMarketCapParams struct {
	InputPrice2   string // fid_input_price_2
	MarketCode    string // fid_cond_mrkt_div_code — 빈 값=>"J"
	ScreenCode    string // fid_cond_scr_div_code — 빈 값=>"20174"
	DivCode       string // fid_div_cls_code — "0":전체, "1":보통주, "2":우선주. 빈 값=>"0"
	InputISCD     string // fid_input_iscd — 필수, "0000"(전체)/"0001"(거래소)/"1001"(코스닥)/"2001"(코스피200)
	TargetCode    string // fid_trgt_cls_code — 빈 값=>"0"
	TargetExclude string // fid_trgt_exls_cls_code — 빈 값=>"0"
	InputPrice1   string // fid_input_price_1
	VolCount      string // fid_vol_cnt
}

// InquireMarketCap 은 시가총액 상위 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_시가총액_상위.md
// path: /uapi/domestic-stock/v1/ranking/market-cap (FHPST01740000)
func (c *Client) InquireMarketCap(ctx context.Context, params InquireMarketCapParams) (*MarketCap, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	scr := params.ScreenCode
	if scr == "" {
		scr = "20174"
	}
	div := params.DivCode
	if div == "" {
		div = "0"
	}
	tgt := params.TargetCode
	if tgt == "" {
		tgt = "0"
	}
	tgtExcl := params.TargetExclude
	if tgtExcl == "" {
		tgtExcl = "0"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ranking/market-cap",
		TrID:   "FHPST01740000",
		Query: map[string]string{
			"fid_input_price_2":      params.InputPrice2,
			"fid_cond_mrkt_div_code": market,
			"fid_cond_scr_div_code":  scr,
			"fid_div_cls_code":       div,
			"fid_input_iscd":         params.InputISCD,
			"fid_trgt_cls_code":      tgt,
			"fid_trgt_exls_cls_code": tgtExcl,
			"fid_input_price_1":      params.InputPrice1,
			"fid_vol_cnt":            params.VolCount,
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

// DividendRate 는 배당률 상위 (HHKDB13470100) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식_배당률_상위.md
// path: /uapi/domestic-stock/v1/ranking/dividend-rate
type DividendRate struct {
	Output1 []DividendRateItem `json:"output1"` // 응답상세 (output1)
}

// DividendRateItem 은 배당률 상위 응답의 한 행.
type DividendRateItem struct {
	Rank          string          `json:"rank"`             // 순위
	ShtCd         string          `json:"sht_cd"`           // 종목코드
	IsinName      string          `json:"isin_name"`        // 종목명
	RecordDate    string          `json:"record_date"`      // 기준일 (YYYYMMDD)
	PerStoDiviAmt decimal.Decimal `json:"per_sto_divi_amt"` // 현금/주식배당금
	DiviRate      float64         `json:"divi_rate,string"` // 현금/주식배당률 (%)
	DiviKind      string          `json:"divi_kind"`        // 배당종류
}

// InquireDividendRateParams 는 배당률 상위 조회 파라미터.
//
// 다른 ranking 과 query 형식이 다름. KIS docs 의 query 키 (CTS_AREA, GB1~GB4, UPJONG, F_DT, T_DT) 그대로 노출.
type InquireDividendRateParams struct {
	CtsArea     string // CTS_AREA — 빈 값(공백) default
	Market      string // GB1 — KOSPI 구분: "0":전체, "1":코스피, "2":코스피200, "3":코스닥. 빈 값=>"0"
	Sector      string // UPJONG — 업종구분 (필수). 예: "0001"(코스피 종합), "1001"(코스닥 종합)
	StockType   string // GB2 — 종목선택: "0":전체, "6":보통주, "7":우선주. 빈 값=>"0"
	DividendCls string // GB3 — 배당구분: "1":주식배당, "2":현금배당. 빈 값=>"1"
	FromDate    string // F_DT — 기준일 From (필수, YYYYMMDD)
	ToDate      string // T_DT — 기준일 To (필수, YYYYMMDD)
	YearCls     string // GB4 — 결산/중간배당: "0":전체, "1":결산배당, "2":중간배당. 빈 값=>"0"
}

// InquireDividendRate 는 배당률 상위 호출.
//
// 한투 docs: docs/api/국내주식/국내주식_배당률_상위.md
// path: /uapi/domestic-stock/v1/ranking/dividend-rate (HHKDB13470100)
func (c *Client) InquireDividendRate(ctx context.Context, params InquireDividendRateParams) (*DividendRate, error) {
	gb1 := params.Market
	if gb1 == "" {
		gb1 = "0"
	}
	gb2 := params.StockType
	if gb2 == "" {
		gb2 = "0"
	}
	gb3 := params.DividendCls
	if gb3 == "" {
		gb3 = "1"
	}
	gb4 := params.YearCls
	if gb4 == "" {
		gb4 = "0"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/ranking/dividend-rate",
		TrID:   "HHKDB13470100",
		Query: map[string]string{
			"CTS_AREA": params.CtsArea,
			"GB1":      gb1,
			"UPJONG":   params.Sector,
			"GB2":      gb2,
			"GB3":      gb3,
			"F_DT":     params.FromDate,
			"T_DT":     params.ToDate,
			"GB4":      gb4,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var res DividendRate
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse DividendRate: %w", err)
	}
	return &res, nil
}
