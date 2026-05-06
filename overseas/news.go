package overseas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// NewsTitle 은 해외뉴스종합(제목) (HHPSTH60100C1) 응답.
//
// 한투 docs: docs/api/해외주식/해외뉴스종합(제목).md
// path: /uapi/overseas-price/v1/quotations/news-title
//
// ANOMALY: 응답 key 가 outblock1 (output/output1 아님).
type NewsTitle struct {
	Outblock1 []NewsTitleItem `json:"outblock1"`
}

// NewsTitleItem 은 해외뉴스종합(제목) 한 행. 모든 필드 string (KIS docs).
type NewsTitleItem struct {
	InfoGb     string `json:"info_gb"`     // 정보구분
	NewsKey    string `json:"news_key"`    // 뉴스키
	DataDt     string `json:"data_dt"`     // 데이터일자
	DataTm     string `json:"data_tm"`     // 데이터시간
	ClassCd    string `json:"class_cd"`    // 분류코드
	ClassName  string `json:"class_name"`  // 분류명
	Source     string `json:"source"`      // 출처
	NationCd   string `json:"nation_cd"`   // 국가코드
	ExchangeCd string `json:"exchange_cd"` // 거래소코드
	Symb       string `json:"symb"`        // 종목코드
	SymbName   string `json:"symb_name"`   // 종목명
	Title      string `json:"title"`       // 뉴스제목
}

// InquireNewsTitleParams 는 해외뉴스종합(제목) 조회 파라미터.
//
// Cts 는 페이지네이션 cursor. 첫 조회 시 "" (공백).
type InquireNewsTitleParams struct {
	InfoGb     string // INFO_GB — 정보구분. 빈 값 default
	ClassCd    string // CLASS_CD — 분류코드. 빈 값 default
	NationCd   string // NATION_CD — 국가코드. 빈 값 default
	ExchangeCd string // EXCHANGE_CD — 거래소코드. 빈 값 default
	Symb       string // SYMB — 종목코드. 빈 값 default (전체)
	DataDt     string // DATA_DT — 데이터일자 YYYYMMDD. 빈 값 default
	DataTm     string // DATA_TM — 데이터시간 HHMMSS. 빈 값 default
	Cts        string // CTS — 페이지네이션 cursor. 빈 값=첫 페이지
}

// InquireNewsTitle 은 해외뉴스종합(제목) 호출.
//
// 한투 docs: docs/api/해외주식/해외뉴스종합(제목).md
// path: /uapi/overseas-price/v1/quotations/news-title (HHPSTH60100C1)
func (c *Client) InquireNewsTitle(ctx context.Context, params InquireNewsTitleParams) (*NewsTitle, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-price/v1/quotations/news-title",
		TrID:   "HHPSTH60100C1",
		Query: map[string]string{
			"INFO_GB":     params.InfoGb,
			"CLASS_CD":    params.ClassCd,
			"NATION_CD":   params.NationCd,
			"EXCHANGE_CD": params.ExchangeCd,
			"SYMB":        params.Symb,
			"DATA_DT":     params.DataDt,
			"DATA_TM":     params.DataTm,
			"CTS":         params.Cts,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res NewsTitle
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse NewsTitle: %w", err)
	}
	return &res, nil
}

// BrknewsTitle 은 해외속보(제목) (FHKST01011801) 응답.
//
// 한투 docs: docs/api/해외주식/해외속보(제목).md
// path: /uapi/overseas-price/v1/quotations/brknews-title
//
// ANOMALY 1: 파라미터명에 FID_ prefix 사용 (일반 파라미터와 다름).
// ANOMALY 2: iscd1-10 + kor_isnm1-10 flat 20 fields (nested 배열 아님).
// ANOMALY 3: FID_COND_SCR_DIV_CODE="11801" hardcoded (Params 미노출).
type BrknewsTitle struct {
	Output []BrknewsTitleItem `json:"output"`
}

// BrknewsTitleItem 은 해외속보(제목) 한 행. 모든 필드 string (KIS docs).
// iscd1-10 / kor_isnm1-10 은 flat field (배열 아님 — KIS 원본 설계).
type BrknewsTitleItem struct {
	CnttUsiqSrno     string `json:"cntt_usiq_srno"`      // 콘텐츠고유일련번호
	NewsOferEntpCode string `json:"news_ofer_entp_code"` // 뉴스제공업체코드
	DataDt           string `json:"data_dt"`             // 데이터일자
	DataTm           string `json:"data_tm"`             // 데이터시간
	HtsPbntTitlCntt  string `json:"hts_pbnt_titl_cntt"`  // HTS게시제목내용
	NewsLrdvCode     string `json:"news_lrdv_code"`      // 뉴스대분류코드
	Dorg             string `json:"dorg"`                // 출처기관
	Iscd1            string `json:"iscd1"`               // 종목코드1
	Iscd2            string `json:"iscd2"`               // 종목코드2
	Iscd3            string `json:"iscd3"`               // 종목코드3
	Iscd4            string `json:"iscd4"`               // 종목코드4
	Iscd5            string `json:"iscd5"`               // 종목코드5
	Iscd6            string `json:"iscd6"`               // 종목코드6
	Iscd7            string `json:"iscd7"`               // 종목코드7
	Iscd8            string `json:"iscd8"`               // 종목코드8
	Iscd9            string `json:"iscd9"`               // 종목코드9
	Iscd10           string `json:"iscd10"`              // 종목코드10
	KorIsnm1         string `json:"kor_isnm1"`           // 한글종목명1
	KorIsnm2         string `json:"kor_isnm2"`           // 한글종목명2
	KorIsnm3         string `json:"kor_isnm3"`           // 한글종목명3
	KorIsnm4         string `json:"kor_isnm4"`           // 한글종목명4
	KorIsnm5         string `json:"kor_isnm5"`           // 한글종목명5
	KorIsnm6         string `json:"kor_isnm6"`           // 한글종목명6
	KorIsnm7         string `json:"kor_isnm7"`           // 한글종목명7
	KorIsnm8         string `json:"kor_isnm8"`           // 한글종목명8
	KorIsnm9         string `json:"kor_isnm9"`           // 한글종목명9
	KorIsnm10        string `json:"kor_isnm10"`          // 한글종목명10
}

// InquireBrknewsTitleParams 는 해외속보(제목) 조회 파라미터.
//
// 파라미터 wire name 에 FID_ prefix 사용 (한투 docs 원본).
// FID_COND_SCR_DIV_CODE="11801" 은 내부 hardcode (Params 미노출).
type InquireBrknewsTitleParams struct {
	NewsOferEntpCode string // FID_NEWS_OFER_ENTP_CODE — 뉴스제공업체코드. 빈 값 default
	MarketClsCode    string // FID_COND_MRKT_CLS_CODE — 시장구분코드. 빈 값 default
	Symbol           string // FID_INPUT_ISCD — 종목코드. 빈 값 default (전체)
	TitleContent     string // FID_TITL_CNTT — 제목내용 키워드. 빈 값 default
	InputDate1       string // FID_INPUT_DATE_1 — 입력일자1 YYYYMMDD. 빈 값 default
	InputHour1       string // FID_INPUT_HOUR_1 — 입력시간1 HHMMSS. 빈 값 default
	RankSortClsCode  string // FID_RANK_SORT_CLS_CODE — 순위정렬구분코드. 빈 값 default
	InputSrno        string // FID_INPUT_SRNO — 입력일련번호. 빈 값 default
}

// InquireBrknewsTitle 은 해외속보(제목) 호출.
//
// 한투 docs: docs/api/해외주식/해외속보(제목).md
// path: /uapi/overseas-price/v1/quotations/brknews-title (FHKST01011801)
func (c *Client) InquireBrknewsTitle(ctx context.Context, params InquireBrknewsTitleParams) (*BrknewsTitle, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-price/v1/quotations/brknews-title",
		TrID:   "FHKST01011801",
		Query: map[string]string{
			"FID_NEWS_OFER_ENTP_CODE": params.NewsOferEntpCode,
			"FID_COND_MRKT_CLS_CODE":  params.MarketClsCode,
			"FID_INPUT_ISCD":          params.Symbol,
			"FID_TITL_CNTT":           params.TitleContent,
			"FID_INPUT_DATE_1":        params.InputDate1,
			"FID_INPUT_HOUR_1":        params.InputHour1,
			"FID_RANK_SORT_CLS_CODE":  params.RankSortClsCode,
			"FID_INPUT_SRNO":          params.InputSrno,
			"FID_COND_SCR_DIV_CODE":   "11801", // hardcoded — Params 미노출
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res BrknewsTitle
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse BrknewsTitle: %w", err)
	}
	return &res, nil
}
