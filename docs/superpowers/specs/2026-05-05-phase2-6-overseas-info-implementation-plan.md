# Phase 2.6 — 해외 정보 확장 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** korea-investment-stock Go 라이브러리에 해외 정보 4 메서드 추가 (`v1.9.0` release).

**Architecture:** Phase 1 인프라 + 패턴 재사용. `overseas/news.go` 신규 (2 뉴스 메서드), `overseas/rights.go` 신규 (2 권리 메서드). 새 internal package 불필요. TDD: testdata fixture → 실패 테스트 → struct + 메서드 구현 → 통과 → commit.

**Tech Stack:** Go 1.25+, `github.com/jarcoal/httpmock`, `github.com/stretchr/testify`. 새 dependency 없음.

**참고 spec:**
- Phase 2.5+ design spec: `docs/superpowers/specs/2026-05-05-phase2-5plus-extension-design.md` (§Phase 2.6)
- Phase 2.4 plan (참조 패턴 — KSD all-string): `docs/superpowers/specs/2026-05-05-phase2-4-ksd-implementation-plan.md`
- Phase 2.3 구현 참조 (same module, similar param style): `overseas/ranking.go`

---

## 사전 정보

| 항목 | 값 |
|------|---|
| 작업 브랜치 | `feat/phase2-6-overseas-info` |
| 시작 HEAD | Phase 2.5 구현 완료 commit (v1.8.0) |
| Release 목표 | `v1.9.0` |
| PR 베이스 | `main` |
| 현재 main 상태 | v1.8.0 publish 완료 (Phase 2.5 통합, 60 메서드) |

---

## 메서드 매핑

| Go 메서드 | path (last segment) | TR_ID | output key | fields | anomalies |
|---|---|---|---|---|---|
| `InquireNewsTitle` | `news-title` | HHPSTH60100C1 | `outblock1 []` | 12 | **unusual key `outblock1`**, CTS pagination cursor |
| `InquireBrknewsTitle` | `brknews-title` | FHKST01011801 | `output []` | 27 | **FID_ prefix params**, iscd1-10/kor_isnm1-10 flat arrays, FID_COND_SCR_DIV_CODE="11801" hardcoded |
| `InquireRightsByIce` | `rights-by-ice` | HHDFS78330900 | `output1 []` | 12 | **output1 only** (no output2) |
| `InquirePeriodRights` | `period-rights` | CTRGT011R | `output []` | 20 | **CTX_AREA_NK50/FK50 cursor pagination**, TR_ID `C` prefix unique, numeric-content-as-String fields |

---

## 파일 구조

### 신규
- `overseas/news.go` — 2 뉴스 메서드 + structs + Params
- `overseas/news_test.go` — 2 테스트 함수
- `overseas/rights.go` — 2 권리 메서드 + structs + Params
- `overseas/rights_test.go` — 2 테스트 함수
- `overseas/testdata/news_title_success.json`
- `overseas/testdata/brknews_title_success.json`
- `overseas/testdata/rights_by_ice_success.json`
- `overseas/testdata/period_rights_success.json`
- `examples/overseas_info/main.go`

### 수정
- `CLAUDE.md` — banner Phase 2.5 → Phase 2.6, plan link 추가
- `README.md` — Available Methods 표 갱신 (60 → 64 메서드)
- `CHANGELOG.md` — `[1.9.0]` entry ABOVE `[1.8.0]`
- `overseas/doc.go` — Phase 2.6 section 추가

---

## 타입 매핑

KIS docs 는 이 4 endpoint 의 모든 응답 필드를 String 타입으로 명시 — Phase 2.4 KSD 와 동일 패턴. Go 에서도 plain `string`. 뉴스 제목/날짜/코드, 권리 비율/날짜 등 모두 string. `decimal.Decimal`/`int64`/`float64` 매핑 없음.

---

## Tasks (9 total)

---

## Task 1: testdata fixtures (4 합성 JSON)

- [ ] Step 1: `overseas/testdata/news_title_success.json` (EP1 — `outblock1`, 12 fields)
- [ ] Step 2: `overseas/testdata/brknews_title_success.json` (EP2 — `output`, 27 fields)
- [ ] Step 3: `overseas/testdata/rights_by_ice_success.json` (EP3 — `output1`, 12 fields)
- [ ] Step 4: `overseas/testdata/period_rights_success.json` (EP4 — `output`, 20 fields)
- [ ] Step 5: validation

```bash
for f in \
  overseas/testdata/news_title_success.json \
  overseas/testdata/brknews_title_success.json \
  overseas/testdata/rights_by_ice_success.json \
  overseas/testdata/period_rights_success.json; do
  python3 -c "import json; json.load(open('$f'))" && echo "$f OK"
done
# Expected: 4 OK lines
```

- [ ] Step 6: commit

```bash
git commit -m "$(cat <<'EOF'
[chore] testdata — 4 overseas info fixture JSON (Phase 2.6)

합성 JSON fixtures (2 records each):
- news_title_success.json (outblock1 key, 12 fields)
- brknews_title_success.json (output key, 27 fields incl. iscd1-10/kor_isnm1-10)
- rights_by_ice_success.json (output1 key, 12 fields)
- period_rights_success.json (output key, 20 fields, CTX cursor)

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

### Fixture content

**Step 1 — `news_title_success.json`** (ANOMALY: `outblock1` key, not `output`/`output1`)
```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "outblock1": [
    {
      "info_gb": "1",
      "news_key": "20260505120001",
      "data_dt": "20260505",
      "data_tm": "120001",
      "class_cd": "001",
      "class_name": "증시",
      "source": "연합뉴스",
      "nation_cd": "US",
      "exchange_cd": "NAS",
      "symb": "AAPL",
      "symb_name": "Apple Inc.",
      "title": "애플 2분기 실적 서프라이즈 발표"
    },
    {
      "info_gb": "1",
      "news_key": "20260505120045",
      "data_dt": "20260505",
      "data_tm": "120045",
      "class_cd": "002",
      "class_name": "경제",
      "source": "블룸버그",
      "nation_cd": "US",
      "exchange_cd": "NAS",
      "symb": "MSFT",
      "symb_name": "Microsoft Corp.",
      "title": "마이크로소프트 AI 투자 확대 계획 발표"
    }
  ]
}
```

**Step 2 — `brknews_title_success.json`** (ANOMALY: FID_ prefix params, 27 fields incl. flat iscd1-10/kor_isnm1-10)
```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": [
    {
      "cntt_usiq_srno": "00001",
      "news_ofer_entp_code": "001",
      "data_dt": "20260505",
      "data_tm": "120100",
      "hts_pbnt_titl_cntt": "미 연준 금리 동결 결정",
      "news_lrdv_code": "A",
      "dorg": "로이터",
      "iscd1": "AAPL",
      "iscd2": "MSFT",
      "iscd3": "",
      "iscd4": "",
      "iscd5": "",
      "iscd6": "",
      "iscd7": "",
      "iscd8": "",
      "iscd9": "",
      "iscd10": "",
      "kor_isnm1": "애플",
      "kor_isnm2": "마이크로소프트",
      "kor_isnm3": "",
      "kor_isnm4": "",
      "kor_isnm5": "",
      "kor_isnm6": "",
      "kor_isnm7": "",
      "kor_isnm8": "",
      "kor_isnm9": "",
      "kor_isnm10": ""
    },
    {
      "cntt_usiq_srno": "00002",
      "news_ofer_entp_code": "002",
      "data_dt": "20260505",
      "data_tm": "121500",
      "hts_pbnt_titl_cntt": "나스닥 강세 지속 전망",
      "news_lrdv_code": "B",
      "dorg": "블룸버그",
      "iscd1": "NVDA",
      "iscd2": "",
      "iscd3": "",
      "iscd4": "",
      "iscd5": "",
      "iscd6": "",
      "iscd7": "",
      "iscd8": "",
      "iscd9": "",
      "iscd10": "",
      "kor_isnm1": "엔비디아",
      "kor_isnm2": "",
      "kor_isnm3": "",
      "kor_isnm4": "",
      "kor_isnm5": "",
      "kor_isnm6": "",
      "kor_isnm7": "",
      "kor_isnm8": "",
      "kor_isnm9": "",
      "kor_isnm10": ""
    }
  ]
}
```

**Step 3 — `rights_by_ice_success.json`** (ANOMALY: `output1` only, no output2)
```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": [
    {
      "anno_dt": "20260401",
      "ca_title": "주식배당",
      "div_lock_dt": "20260315",
      "pay_dt": "20260430",
      "record_dt": "20260331",
      "validity_dt": "20260401",
      "local_end_dt": "20260330",
      "lock_dt": "20260328",
      "delist_dt": "",
      "redempt_dt": "",
      "early_redempt_dt": "",
      "effective_dt": "20260401"
    },
    {
      "anno_dt": "20260402",
      "ca_title": "현금배당",
      "div_lock_dt": "20260316",
      "pay_dt": "20260501",
      "record_dt": "20260401",
      "validity_dt": "20260402",
      "local_end_dt": "20260331",
      "lock_dt": "20260329",
      "delist_dt": "",
      "redempt_dt": "",
      "early_redempt_dt": "",
      "effective_dt": "20260402"
    }
  ]
}
```

**Step 4 — `period_rights_success.json`** (ANOMALY: CTX_AREA_NK50/FK50 cursor pagination, TR_ID `C` prefix)
```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": [
    {
      "bass_dt": "20260401",
      "rght_type_cd": "D",
      "pdno": "AAPL",
      "prdt_name": "Apple Inc.",
      "prdt_type_cd": "512",
      "std_pdno": "US0378331005",
      "acpl_bass_dt": "20260315",
      "sbsc_strt_dt": "20260316",
      "sbsc_end_dt": "20260322",
      "cash_alct_rt": "0.05",
      "stck_alct_rt": "0.00",
      "crcy_cd": "USD",
      "crcy_cd2": "",
      "crcy_cd3": "",
      "crcy_cd4": "",
      "alct_frcr_unpr": "1.50",
      "stkp_dvdn_frcr_amt2": "0",
      "stkp_dvdn_frcr_amt3": "0",
      "stkp_dvdn_frcr_amt4": "0",
      "dfnt_yn": "Y"
    },
    {
      "bass_dt": "20260402",
      "rght_type_cd": "R",
      "pdno": "MSFT",
      "prdt_name": "Microsoft Corp.",
      "prdt_type_cd": "512",
      "std_pdno": "US5949181045",
      "acpl_bass_dt": "20260316",
      "sbsc_strt_dt": "20260317",
      "sbsc_end_dt": "20260323",
      "cash_alct_rt": "0.00",
      "stck_alct_rt": "0.10",
      "crcy_cd": "USD",
      "crcy_cd2": "",
      "crcy_cd3": "",
      "crcy_cd4": "",
      "alct_frcr_unpr": "0",
      "stkp_dvdn_frcr_amt2": "0",
      "stkp_dvdn_frcr_amt3": "0",
      "stkp_dvdn_frcr_amt4": "0",
      "dfnt_yn": "N"
    }
  ]
}
```

---

## Task 2: `overseas/news.go` base + InquireNewsTitle (EP1)

**Files:**
- Create: `overseas/news.go`
- Create: `overseas/news_test.go`

- [ ] Step 1: APPEND test code to `overseas/news_test.go`
- [ ] Step 2: Verify FAIL — `go test ./overseas/... -run TestClient_InquireNewsTitle -v` (compile error expected)
- [ ] Step 3: Create `overseas/news.go` with package header, imports, NewsTitle struct + Params + method
- [ ] Step 4: Verify PASS — `go test ./overseas/... -run TestClient_InquireNewsTitle -v`
- [ ] Step 5: `gofmt -w overseas/news.go overseas/news_test.go && go vet ./overseas/...`
- [ ] Step 6: commit

### 메서드 매핑
- Path: `/uapi/overseas-price/v1/quotations/news-title`
- TR_ID: `HHPSTH60100C1`
- Output key: **`outblock1 []NewsTitleItem`** (NOT `output`/`output1` — anomaly)

### 응답 struct 필드 (모두 string)

| Go field | json tag | 설명 |
|---|---|---|
| `InfoGb` | `info_gb` | 정보구분 |
| `NewsKey` | `news_key` | 뉴스키 |
| `DataDt` | `data_dt` | 데이터일자 |
| `DataTm` | `data_tm` | 데이터시간 |
| `ClassCd` | `class_cd` | 분류코드 |
| `ClassName` | `class_name` | 분류명 |
| `Source` | `source` | 출처 |
| `NationCd` | `nation_cd` | 국가코드 |
| `ExchangeCd` | `exchange_cd` | 거래소코드 |
| `Symb` | `symb` | 종목코드 |
| `SymbName` | `symb_name` | 종목명 |
| `Title` | `title` | 뉴스제목 |

### Params struct fields

| Go field | wire name | default |
|---|---|---|
| `InfoGb` | `INFO_GB` | `""` (공백) |
| `ClassCd` | `CLASS_CD` | `""` (공백) |
| `NationCd` | `NATION_CD` | `""` (공백) |
| `ExchangeCd` | `EXCHANGE_CD` | `""` (공백) |
| `Symb` | `SYMB` | `""` (공백) |
| `DataDt` | `DATA_DT` | `""` (공백) |
| `DataTm` | `DATA_TM` | `""` (공백) |
| `Cts` | `CTS` | `""` (페이지네이션 cursor) |

### 구현 코드

```go
// File: overseas/news.go
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
```

### 테스트 코드

```go
// File: overseas/news_test.go
package overseas_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/overseas"
)

func TestClient_InquireNewsTitle(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/news-title`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "news_title_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireNewsTitle(context.Background(), overseas.InquireNewsTitleParams{
		NationCd:   "US",
		ExchangeCd: "NAS",
		DataDt:     "20260505",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "US", capturedQuery.Get("NATION_CD"))
	assert.Equal(t, "NAS", capturedQuery.Get("EXCHANGE_CD"))
	assert.Equal(t, "20260505", capturedQuery.Get("DATA_DT"))

	require.Len(t, res.Outblock1, 2)
	assert.Equal(t, "AAPL", res.Outblock1[0].Symb)
	assert.Equal(t, "20260505120001", res.Outblock1[0].NewsKey)
	assert.Equal(t, "애플 2분기 실적 서프라이즈 발표", res.Outblock1[0].Title)
}
```

### Commit

```bash
git commit -m "$(cat <<'EOF'
[feat] overseas — InquireNewsTitle (해외뉴스종합 제목, HHPSTH60100C1)

- NewsTitle / NewsTitleItem (12 fields, all string)
- InquireNewsTitleParams (8 query params: INFO_GB/CLASS_CD/NATION_CD/EXCHANGE_CD/SYMB/DATA_DT/DATA_TM/CTS)
- ANOMALY: 응답 key outblock1 (output/output1 아님) — json:"outblock1" 태그
- CTS pagination cursor 지원
- TestClient_InquireNewsTitle — fixture news_title_success.json

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 3: InquireBrknewsTitle (EP2)

**Files:** Modify `overseas/news.go` (append), `overseas/news_test.go` (append)

해외속보(제목) 조회. **ANOMALY 1: 파라미터명에 FID_ prefix 사용.** **ANOMALY 2: iscd1-10 + kor_isnm1-10 (flat 20 fields).** **ANOMALY 3: FID_COND_SCR_DIV_CODE="11801" hardcoded in Query map (Params 미노출).** `output []BrknewsTitleItem` 27 fields.

- [ ] Step 1: APPEND test code to `overseas/news_test.go`
- [ ] Step 2: Verify FAIL — `go test ./overseas/... -run TestClient_InquireBrknewsTitle -v`
- [ ] Step 3: APPEND struct + Params + method to `overseas/news.go`
- [ ] Step 4: Verify PASS — `go test ./overseas/... -run TestClient_InquireBrknewsTitle -v`
- [ ] Step 5: `gofmt -w overseas/news.go overseas/news_test.go && go vet ./overseas/...`
- [ ] Step 6: commit

### 메서드 매핑
- Path: `/uapi/overseas-price/v1/quotations/brknews-title`
- TR_ID: `FHKST01011801`
- Output key: `output []BrknewsTitleItem`

### 응답 struct 필드 (모두 string, 27 fields)

| Go field | json tag | 설명 |
|---|---|---|
| `CnttUsiqSrno` | `cntt_usiq_srno` | 콘텐츠고유일련번호 |
| `NewsOferEntpCode` | `news_ofer_entp_code` | 뉴스제공업체코드 |
| `DataDt` | `data_dt` | 데이터일자 |
| `DataTm` | `data_tm` | 데이터시간 |
| `HtsPbntTitlCntt` | `hts_pbnt_titl_cntt` | HTS게시제목내용 |
| `NewsLrdvCode` | `news_lrdv_code` | 뉴스대분류코드 |
| `Dorg` | `dorg` | 출처기관 |
| `Iscd1` | `iscd1` | 종목코드1 |
| `Iscd2` | `iscd2` | 종목코드2 |
| `Iscd3` | `iscd3` | 종목코드3 |
| `Iscd4` | `iscd4` | 종목코드4 |
| `Iscd5` | `iscd5` | 종목코드5 |
| `Iscd6` | `iscd6` | 종목코드6 |
| `Iscd7` | `iscd7` | 종목코드7 |
| `Iscd8` | `iscd8` | 종목코드8 |
| `Iscd9` | `iscd9` | 종목코드9 |
| `Iscd10` | `iscd10` | 종목코드10 |
| `KorIsnm1` | `kor_isnm1` | 한글종목명1 |
| `KorIsnm2` | `kor_isnm2` | 한글종목명2 |
| `KorIsnm3` | `kor_isnm3` | 한글종목명3 |
| `KorIsnm4` | `kor_isnm4` | 한글종목명4 |
| `KorIsnm5` | `kor_isnm5` | 한글종목명5 |
| `KorIsnm6` | `kor_isnm6` | 한글종목명6 |
| `KorIsnm7` | `kor_isnm7` | 한글종목명7 |
| `KorIsnm8` | `kor_isnm8` | 한글종목명8 |
| `KorIsnm9` | `kor_isnm9` | 한글종목명9 |
| `KorIsnm10` | `kor_isnm10` | 한글종목명10 |

### Params struct fields

| Go field | wire name | default | 비고 |
|---|---|---|---|
| `NewsOferEntpCode` | `FID_NEWS_OFER_ENTP_CODE` | `""` | 뉴스제공업체코드 |
| `MarketClsCode` | `FID_COND_MRKT_CLS_CODE` | `""` | 시장구분코드 |
| `Symbol` | `FID_INPUT_ISCD` | `""` | 종목코드 |
| `TitleContent` | `FID_TITL_CNTT` | `""` | 제목내용 |
| `InputDate1` | `FID_INPUT_DATE_1` | `""` | 입력일자1 |
| `InputHour1` | `FID_INPUT_HOUR_1` | `""` | 입력시간1 |
| `RankSortClsCode` | `FID_RANK_SORT_CLS_CODE` | `""` | 순위정렬구분코드 |
| `InputSrno` | `FID_INPUT_SRNO` | `""` | 입력일련번호 |
| — | `FID_COND_SCR_DIV_CODE` | `"11801"` (hardcoded) | Params 미노출 |

### 구현 코드

```go
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
	CnttUsiqSrno    string `json:"cntt_usiq_srno"`    // 콘텐츠고유일련번호
	NewsOferEntpCode string `json:"news_ofer_entp_code"` // 뉴스제공업체코드
	DataDt          string `json:"data_dt"`            // 데이터일자
	DataTm          string `json:"data_tm"`            // 데이터시간
	HtsPbntTitlCntt string `json:"hts_pbnt_titl_cntt"` // HTS게시제목내용
	NewsLrdvCode    string `json:"news_lrdv_code"`     // 뉴스대분류코드
	Dorg            string `json:"dorg"`               // 출처기관
	Iscd1           string `json:"iscd1"`              // 종목코드1
	Iscd2           string `json:"iscd2"`              // 종목코드2
	Iscd3           string `json:"iscd3"`              // 종목코드3
	Iscd4           string `json:"iscd4"`              // 종목코드4
	Iscd5           string `json:"iscd5"`              // 종목코드5
	Iscd6           string `json:"iscd6"`              // 종목코드6
	Iscd7           string `json:"iscd7"`              // 종목코드7
	Iscd8           string `json:"iscd8"`              // 종목코드8
	Iscd9           string `json:"iscd9"`              // 종목코드9
	Iscd10          string `json:"iscd10"`             // 종목코드10
	KorIsnm1        string `json:"kor_isnm1"`          // 한글종목명1
	KorIsnm2        string `json:"kor_isnm2"`          // 한글종목명2
	KorIsnm3        string `json:"kor_isnm3"`          // 한글종목명3
	KorIsnm4        string `json:"kor_isnm4"`          // 한글종목명4
	KorIsnm5        string `json:"kor_isnm5"`          // 한글종목명5
	KorIsnm6        string `json:"kor_isnm6"`          // 한글종목명6
	KorIsnm7        string `json:"kor_isnm7"`          // 한글종목명7
	KorIsnm8        string `json:"kor_isnm8"`          // 한글종목명8
	KorIsnm9        string `json:"kor_isnm9"`          // 한글종목명9
	KorIsnm10       string `json:"kor_isnm10"`         // 한글종목명10
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
```

### 테스트 코드

```go
func TestClient_InquireBrknewsTitle(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/brknews-title`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "brknews_title_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireBrknewsTitle(context.Background(), overseas.InquireBrknewsTitleParams{
		InputDate1: "20260505",
		Symbol:     "AAPL",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// FID_ prefix wire name 검증
	assert.Equal(t, "20260505", capturedQuery.Get("FID_INPUT_DATE_1"))
	assert.Equal(t, "AAPL", capturedQuery.Get("FID_INPUT_ISCD"))
	// hardcoded param 검증
	assert.Equal(t, "11801", capturedQuery.Get("FID_COND_SCR_DIV_CODE"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, "미 연준 금리 동결 결정", res.Output[0].HtsPbntTitlCntt)
	assert.Equal(t, "AAPL", res.Output[0].Iscd1)
	assert.Equal(t, "애플", res.Output[0].KorIsnm1)
	// 빈 iscd/kor_isnm 필드 확인
	assert.Equal(t, "", res.Output[0].Iscd3)
	assert.Equal(t, "", res.Output[0].KorIsnm3)
}
```

### Commit

```bash
git commit -m "$(cat <<'EOF'
[feat] overseas — InquireBrknewsTitle (해외속보 제목, FHKST01011801)

- BrknewsTitle / BrknewsTitleItem (27 fields, all string)
- InquireBrknewsTitleParams (8 query params, FID_ prefix wire name)
- ANOMALY 1: FID_ prefix 파라미터명 (일반 파라미터와 다름)
- ANOMALY 2: iscd1-10 + kor_isnm1-10 flat 20 fields (nested 배열 아님)
- ANOMALY 3: FID_COND_SCR_DIV_CODE="11801" hardcoded in Query map
- TestClient_InquireBrknewsTitle — fixture brknews_title_success.json

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 4: `overseas/rights.go` base + InquireRightsByIce (EP3)

**Files:**
- Create: `overseas/rights.go`
- Create: `overseas/rights_test.go`

해외주식_권리종합 조회. **ANOMALY: `output1 []` 만 존재 (output2 없음).** 12 fields.

- [ ] Step 1: APPEND test code to `overseas/rights_test.go`
- [ ] Step 2: Verify FAIL — `go test ./overseas/... -run TestClient_InquireRightsByIce -v`
- [ ] Step 3: Create `overseas/rights.go` with package header, imports, RightsByIce struct + Params + method
- [ ] Step 4: Verify PASS — `go test ./overseas/... -run TestClient_InquireRightsByIce -v`
- [ ] Step 5: `gofmt -w overseas/rights.go overseas/rights_test.go && go vet ./overseas/...`
- [ ] Step 6: commit

### 메서드 매핑
- Path: `/uapi/overseas-price/v1/quotations/rights-by-ice`
- TR_ID: `HHDFS78330900`
- Output key: **`output1 []RightsByIceItem`** (output2 없음 — anomaly)

### 응답 struct 필드 (모두 string, 12 fields)

| Go field | json tag | 설명 |
|---|---|---|
| `AnnoDt` | `anno_dt` | 공시일자 |
| `CaTitle` | `ca_title` | 권리종류명 |
| `DivLockDt` | `div_lock_dt` | 배당락일 |
| `PayDt` | `pay_dt` | 지급일 |
| `RecordDt` | `record_dt` | 기준일 |
| `ValidityDt` | `validity_dt` | 유효일 |
| `LocalEndDt` | `local_end_dt` | 현지종료일 |
| `LockDt` | `lock_dt` | 권리확정일 |
| `DelistDt` | `delist_dt` | 상장폐지일 |
| `RedemptDt` | `redempt_dt` | 상환일 |
| `EarlyRedemptDt` | `early_redempt_dt` | 조기상환일 |
| `EffectiveDt` | `effective_dt` | 효력발생일 |

### Params struct fields

| Go field | wire name | default |
|---|---|---|
| `NCod` | `NCOD` | `""` (국가코드) |
| `Symb` | `SYMB` | `""` (종목코드) |
| `StYmd` | `ST_YMD` | `""` (조회시작일 YYYYMMDD) |
| `EdYmd` | `ED_YMD` | `""` (조회종료일 YYYYMMDD) |

### 구현 코드

```go
// File: overseas/rights.go
package overseas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// RightsByIce 는 해외주식_권리종합 (HHDFS78330900) 응답.
//
// 한투 docs: docs/api/해외주식/해외주식_권리종합.md
// path: /uapi/overseas-price/v1/quotations/rights-by-ice
//
// ANOMALY: output1 만 존재 (output2 없음).
type RightsByIce struct {
	Output1 []RightsByIceItem `json:"output1"`
}

// RightsByIceItem 은 해외주식_권리종합 한 행. 모든 필드 string (KIS docs).
type RightsByIceItem struct {
	AnnoDt        string `json:"anno_dt"`         // 공시일자
	CaTitle       string `json:"ca_title"`        // 권리종류명
	DivLockDt     string `json:"div_lock_dt"`     // 배당락일
	PayDt         string `json:"pay_dt"`          // 지급일
	RecordDt      string `json:"record_dt"`       // 기준일
	ValidityDt    string `json:"validity_dt"`     // 유효일
	LocalEndDt    string `json:"local_end_dt"`    // 현지종료일
	LockDt        string `json:"lock_dt"`         // 권리확정일
	DelistDt      string `json:"delist_dt"`       // 상장폐지일
	RedemptDt     string `json:"redempt_dt"`      // 상환일
	EarlyRedemptDt string `json:"early_redempt_dt"` // 조기상환일
	EffectiveDt   string `json:"effective_dt"`    // 효력발생일
}

// InquireRightsByIceParams 는 해외주식_권리종합 조회 파라미터.
type InquireRightsByIceParams struct {
	NCod  string // NCOD — 국가코드. 빈 값 default
	Symb  string // SYMB — 종목코드. 빈 값 default (전체)
	StYmd string // ST_YMD — 조회시작일 YYYYMMDD. 빈 값 default
	EdYmd string // ED_YMD — 조회종료일 YYYYMMDD. 빈 값 default
}

// InquireRightsByIce 는 해외주식_권리종합 호출.
//
// 한투 docs: docs/api/해외주식/해외주식_권리종합.md
// path: /uapi/overseas-price/v1/quotations/rights-by-ice (HHDFS78330900)
func (c *Client) InquireRightsByIce(ctx context.Context, params InquireRightsByIceParams) (*RightsByIce, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-price/v1/quotations/rights-by-ice",
		TrID:   "HHDFS78330900",
		Query: map[string]string{
			"NCOD":   params.NCod,
			"SYMB":   params.Symb,
			"ST_YMD": params.StYmd,
			"ED_YMD": params.EdYmd,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res RightsByIce
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse RightsByIce: %w", err)
	}
	return &res, nil
}
```

### 테스트 코드

```go
// File: overseas/rights_test.go
package overseas_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/overseas"
)

func TestClient_InquireRightsByIce(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/rights-by-ice`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "rights_by_ice_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquireRightsByIce(context.Background(), overseas.InquireRightsByIceParams{
		NCod:  "US",
		Symb:  "AAPL",
		StYmd: "20260401",
		EdYmd: "20260430",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "US", capturedQuery.Get("NCOD"))
	assert.Equal(t, "AAPL", capturedQuery.Get("SYMB"))
	assert.Equal(t, "20260401", capturedQuery.Get("ST_YMD"))
	assert.Equal(t, "20260430", capturedQuery.Get("ED_YMD"))

	// output1 only — output2 없음 anomaly 검증
	require.Len(t, res.Output1, 2)
	assert.Equal(t, "20260401", res.Output1[0].AnnoDt)
	assert.Equal(t, "주식배당", res.Output1[0].CaTitle)
	assert.Equal(t, "20260430", res.Output1[0].PayDt)
}
```

### Commit

```bash
git commit -m "$(cat <<'EOF'
[feat] overseas — InquireRightsByIce (해외주식 권리종합, HHDFS78330900)

- RightsByIce / RightsByIceItem (12 fields, all string)
- InquireRightsByIceParams (4 query params: NCOD/SYMB/ST_YMD/ED_YMD)
- ANOMALY: output1 only (output2 없음)
- TestClient_InquireRightsByIce — fixture rights_by_ice_success.json

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 5: InquirePeriodRights (EP4)

**Files:** Modify `overseas/rights.go` (append), `overseas/rights_test.go` (append)

해외주식_기간별권리조회. **ANOMALY 1: TR_ID `CTRGT011R` — `C` prefix (다른 해외주식 endpoint 와 다름).** **ANOMALY 2: CTX_AREA_NK50/FK50 cursor pagination params.** **ANOMALY 3: numeric-content 필드 모두 String.** `output []PeriodRightsItem` 20 fields.

- [ ] Step 1: APPEND test code to `overseas/rights_test.go`
- [ ] Step 2: Verify FAIL — `go test ./overseas/... -run TestClient_InquirePeriodRights -v`
- [ ] Step 3: APPEND struct + Params + method to `overseas/rights.go`
- [ ] Step 4: Verify PASS — `go test ./overseas/... -run TestClient_InquirePeriodRights -v`
- [ ] Step 5: `gofmt -w overseas/rights.go overseas/rights_test.go && go vet ./overseas/...`
- [ ] Step 6: commit

### 메서드 매핑
- Path: `/uapi/overseas-price/v1/quotations/period-rights`
- TR_ID: **`CTRGT011R`** (`C` prefix — 다른 해외주식 TR_ID 와 다름)
- Output key: `output []PeriodRightsItem`

### 응답 struct 필드 (모두 string, 20 fields)

| Go field | json tag | 설명 |
|---|---|---|
| `BassDt` | `bass_dt` | 기준일자 |
| `RghtTypeCd` | `rght_type_cd` | 권리유형코드 |
| `Pdno` | `pdno` | 상품번호 |
| `PrdtName` | `prdt_name` | 상품명 |
| `PrdtTypeCd` | `prdt_type_cd` | 상품유형코드 |
| `StdPdno` | `std_pdno` | 표준상품번호 |
| `AcplBassDt` | `acpl_bass_dt` | 발생기준일 |
| `SbscStrtDt` | `sbsc_strt_dt` | 청약시작일 |
| `SbscEndDt` | `sbsc_end_dt` | 청약종료일 |
| `CashAlctRt` | `cash_alct_rt` | 현금배분율 |
| `StckAlctRt` | `stck_alct_rt` | 주식배분율 |
| `CrcyCd` | `crcy_cd` | 통화코드1 |
| `CrcyCd2` | `crcy_cd2` | 통화코드2 |
| `CrcyCd3` | `crcy_cd3` | 통화코드3 |
| `CrcyCd4` | `crcy_cd4` | 통화코드4 |
| `AlctFrcrUnpr` | `alct_frcr_unpr` | 배분외화단가 |
| `StkpDvdnFrcrAmt2` | `stkp_dvdn_frcr_amt2` | 주식배당외화금액2 |
| `StkpDvdnFrcrAmt3` | `stkp_dvdn_frcr_amt3` | 주식배당외화금액3 |
| `StkpDvdnFrcrAmt4` | `stkp_dvdn_frcr_amt4` | 주식배당외화금액4 |
| `DfntYn` | `dfnt_yn` | 확정여부 |

### Params struct fields

| Go field | wire name | default | 비고 |
|---|---|---|---|
| `RghtTypeCd` | `RGHT_TYPE_CD` | `""` | 권리유형코드 |
| `InqrDvsnCd` | `INQR_DVSN_CD` | `""` | 조회구분코드 |
| `InqrStrtDt` | `INQR_STRT_DT` | `""` | 조회시작일 YYYYMMDD |
| `InqrEndDt` | `INQR_END_DT` | `""` | 조회종료일 YYYYMMDD |
| `Pdno` | `PDNO` | `""` | 상품번호 |
| `PrdtTypeCd` | `PRDT_TYPE_CD` | `""` | 상품유형코드 |
| `CtxAreaNk50` | `CTX_AREA_NK50` | `""` | cursor pagination (NK50) |
| `CtxAreaFk50` | `CTX_AREA_FK50` | `""` | cursor pagination (FK50) |

### 구현 코드

```go
// PeriodRights 는 해외주식_기간별권리조회 (CTRGT011R) 응답.
//
// 한투 docs: docs/api/해외주식/해외주식_기간별권리조회.md
// path: /uapi/overseas-price/v1/quotations/period-rights
//
// ANOMALY 1: TR_ID = CTRGT011R (C prefix — 다른 해외주식 TR_ID 와 다름).
// ANOMALY 2: CTX_AREA_NK50/FK50 cursor pagination params 사용.
// ANOMALY 3: 수치 콘텐츠 필드 (비율/금액) 모두 String 타입 (KIS docs).
type PeriodRights struct {
	Output []PeriodRightsItem `json:"output"`
}

// PeriodRightsItem 은 해외주식_기간별권리조회 한 행. 모든 필드 string (KIS docs).
type PeriodRightsItem struct {
	BassDt            string `json:"bass_dt"`             // 기준일자
	RghtTypeCd        string `json:"rght_type_cd"`        // 권리유형코드
	Pdno              string `json:"pdno"`                // 상품번호
	PrdtName          string `json:"prdt_name"`           // 상품명
	PrdtTypeCd        string `json:"prdt_type_cd"`        // 상품유형코드
	StdPdno           string `json:"std_pdno"`            // 표준상품번호
	AcplBassDt        string `json:"acpl_bass_dt"`        // 발생기준일
	SbscStrtDt        string `json:"sbsc_strt_dt"`        // 청약시작일
	SbscEndDt         string `json:"sbsc_end_dt"`         // 청약종료일
	CashAlctRt        string `json:"cash_alct_rt"`        // 현금배분율
	StckAlctRt        string `json:"stck_alct_rt"`        // 주식배분율
	CrcyCd            string `json:"crcy_cd"`             // 통화코드1
	CrcyCd2           string `json:"crcy_cd2"`            // 통화코드2
	CrcyCd3           string `json:"crcy_cd3"`            // 통화코드3
	CrcyCd4           string `json:"crcy_cd4"`            // 통화코드4
	AlctFrcrUnpr      string `json:"alct_frcr_unpr"`      // 배분외화단가
	StkpDvdnFrcrAmt2  string `json:"stkp_dvdn_frcr_amt2"` // 주식배당외화금액2
	StkpDvdnFrcrAmt3  string `json:"stkp_dvdn_frcr_amt3"` // 주식배당외화금액3
	StkpDvdnFrcrAmt4  string `json:"stkp_dvdn_frcr_amt4"` // 주식배당외화금액4
	DfntYn            string `json:"dfnt_yn"`             // 확정여부
}

// InquirePeriodRightsParams 는 해외주식_기간별권리조회 파라미터.
//
// CtxAreaNk50/CtxAreaFk50 는 cursor pagination. 첫 조회 시 "" (공백).
type InquirePeriodRightsParams struct {
	RghtTypeCd  string // RGHT_TYPE_CD — 권리유형코드. 빈 값 default
	InqrDvsnCd  string // INQR_DVSN_CD — 조회구분코드. 빈 값 default
	InqrStrtDt  string // INQR_STRT_DT — 조회시작일 YYYYMMDD. 빈 값 default
	InqrEndDt   string // INQR_END_DT — 조회종료일 YYYYMMDD. 빈 값 default
	Pdno        string // PDNO — 상품번호. 빈 값 default (전체)
	PrdtTypeCd  string // PRDT_TYPE_CD — 상품유형코드. 빈 값 default
	CtxAreaNk50 string // CTX_AREA_NK50 — cursor pagination. 빈 값=첫 페이지
	CtxAreaFk50 string // CTX_AREA_FK50 — cursor pagination. 빈 값=첫 페이지
}

// InquirePeriodRights 는 해외주식_기간별권리조회 호출.
//
// 한투 docs: docs/api/해외주식/해외주식_기간별권리조회.md
// path: /uapi/overseas-price/v1/quotations/period-rights (CTRGT011R)
func (c *Client) InquirePeriodRights(ctx context.Context, params InquirePeriodRightsParams) (*PeriodRights, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/overseas-price/v1/quotations/period-rights",
		TrID:   "CTRGT011R",
		Query: map[string]string{
			"RGHT_TYPE_CD":  params.RghtTypeCd,
			"INQR_DVSN_CD":  params.InqrDvsnCd,
			"INQR_STRT_DT":  params.InqrStrtDt,
			"INQR_END_DT":   params.InqrEndDt,
			"PDNO":          params.Pdno,
			"PRDT_TYPE_CD":  params.PrdtTypeCd,
			"CTX_AREA_NK50": params.CtxAreaNk50,
			"CTX_AREA_FK50": params.CtxAreaFk50,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var res PeriodRights
	if err := json.Unmarshal(resp.Raw, &res); err != nil {
		return nil, fmt.Errorf("kis: parse PeriodRights: %w", err)
	}
	return &res, nil
}
```

### 테스트 코드

```go
func TestClient_InquirePeriodRights(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery url.Values
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/period-rights`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "period_rights_success.json")), nil
		},
	)

	c := newTestClient(t)
	res, err := c.InquirePeriodRights(context.Background(), overseas.InquirePeriodRightsParams{
		InqrStrtDt: "20260401",
		InqrEndDt:  "20260430",
		Pdno:       "AAPL",
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, "20260401", capturedQuery.Get("INQR_STRT_DT"))
	assert.Equal(t, "20260430", capturedQuery.Get("INQR_END_DT"))
	assert.Equal(t, "AAPL", capturedQuery.Get("PDNO"))
	// cursor pagination 첫 조회 시 빈 값 전송 확인
	assert.Equal(t, "", capturedQuery.Get("CTX_AREA_NK50"))
	assert.Equal(t, "", capturedQuery.Get("CTX_AREA_FK50"))

	require.Len(t, res.Output, 2)
	assert.Equal(t, "20260401", res.Output[0].BassDt)
	assert.Equal(t, "AAPL", res.Output[0].Pdno)
	assert.Equal(t, "Y", res.Output[0].DfntYn)
	assert.Equal(t, "USD", res.Output[0].CrcyCd)
}
```

### Commit

```bash
git commit -m "$(cat <<'EOF'
[feat] overseas — InquirePeriodRights (해외주식 기간별권리조회, CTRGT011R)

- PeriodRights / PeriodRightsItem (20 fields, all string)
- InquirePeriodRightsParams (8 query params incl. CTX_AREA_NK50/FK50 cursor)
- ANOMALY 1: TR_ID CTRGT011R (C prefix — 해외주식 TR_ID 패턴 예외)
- ANOMALY 2: CTX_AREA_NK50/FK50 cursor pagination 지원
- ANOMALY 3: 비율/금액 numeric-content 필드 모두 String
- TestClient_InquirePeriodRights — fixture period_rights_success.json

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 6: `examples/overseas_info/main.go`

**Files:**
- Create: `examples/overseas_info/main.go`

- [ ] Step 1: create example file with 4 메서드 통합 호출
- [ ] Step 2: `go build ./examples/overseas_info/...` — PASS
- [ ] Step 3: commit

```go
// File: examples/overseas_info/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	kis "github.com/kenshin579/korea-investment-stock"
	"github.com/kenshin579/korea-investment-stock/overseas"
)

func main() {
	client, err := kis.NewClient(kis.Config{
		AppKey:    os.Getenv("KIS_APP_KEY"),
		AppSecret: os.Getenv("KIS_APP_SECRET"),
		IsMock:    true,
	})
	if err != nil {
		log.Fatalf("client init: %v", err)
	}
	ctx := context.Background()

	// EP1: InquireNewsTitle — 해외뉴스종합(제목)
	// ANOMALY: 응답 key outblock1
	news, err := client.Overseas.InquireNewsTitle(ctx, overseas.InquireNewsTitleParams{
		NationCd:   "US",
		ExchangeCd: "NAS",
		DataDt:     "20260505",
	})
	if err != nil {
		log.Printf("InquireNewsTitle error: %v", err)
	} else {
		fmt.Printf("=== InquireNewsTitle (outblock1 key) ===\n")
		for _, item := range news.Outblock1 {
			fmt.Printf("  [%s %s] %s — %s\n", item.DataDt, item.DataTm, item.Symb, item.Title)
		}
	}

	// EP2: InquireBrknewsTitle — 해외속보(제목)
	// ANOMALY: FID_ prefix params, iscd1-10/kor_isnm1-10 flat, hardcoded FID_COND_SCR_DIV_CODE
	brknews, err := client.Overseas.InquireBrknewsTitle(ctx, overseas.InquireBrknewsTitleParams{
		InputDate1: "20260505",
	})
	if err != nil {
		log.Printf("InquireBrknewsTitle error: %v", err)
	} else {
		fmt.Printf("\n=== InquireBrknewsTitle (FID_ prefix params) ===\n")
		for _, item := range brknews.Output {
			fmt.Printf("  [%s %s] %s (iscd1=%s, kor_isnm1=%s)\n",
				item.DataDt, item.DataTm, item.HtsPbntTitlCntt, item.Iscd1, item.KorIsnm1)
		}
	}

	// EP3: InquireRightsByIce — 해외주식_권리종합
	// ANOMALY: output1 only (no output2)
	rights, err := client.Overseas.InquireRightsByIce(ctx, overseas.InquireRightsByIceParams{
		NCod:  "US",
		Symb:  "AAPL",
		StYmd: "20260401",
		EdYmd: "20260430",
	})
	if err != nil {
		log.Printf("InquireRightsByIce error: %v", err)
	} else {
		fmt.Printf("\n=== InquireRightsByIce (output1 only) ===\n")
		for _, item := range rights.Output1 {
			fmt.Printf("  [%s] %s — pay_dt=%s\n", item.AnnoDt, item.CaTitle, item.PayDt)
		}
	}

	// EP4: InquirePeriodRights — 해외주식_기간별권리조회
	// ANOMALY: TR_ID C prefix, CTX_AREA_NK50/FK50 cursor pagination
	periodRights, err := client.Overseas.InquirePeriodRights(ctx, overseas.InquirePeriodRightsParams{
		InqrStrtDt: "20260401",
		InqrEndDt:  "20260430",
	})
	if err != nil {
		log.Printf("InquirePeriodRights error: %v", err)
	} else {
		fmt.Printf("\n=== InquirePeriodRights (CTRGT011R, CTX cursor) ===\n")
		for _, item := range periodRights.Output {
			fmt.Printf("  [%s] %s (%s) — cash_alct_rt=%s dfnt_yn=%s\n",
				item.BassDt, item.Pdno, item.PrdtName, item.CashAlctRt, item.DfntYn)
		}
	}
}
```

### Commit

```bash
git commit -m "$(cat <<'EOF'
[feat] examples — overseas_info: 4 메서드 통합 예제 (Phase 2.6)

InquireNewsTitle / InquireBrknewsTitle / InquireRightsByIce / InquirePeriodRights
각 anomaly 주석 포함 (outblock1 key, FID_ prefix, output1-only, CTX cursor)

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 7: 문서 갱신

**Files:** `CLAUDE.md`, `README.md`, `CHANGELOG.md`, `overseas/doc.go`

- [ ] Step 1: `CLAUDE.md` — banner Phase 2.5 → Phase 2.6, plan link 추가
- [ ] Step 2: `README.md` — Available Methods 표 Phase 2.6 4행 추가 + 누적 60 → 64
- [ ] Step 3: `CHANGELOG.md` — `[1.9.0]` entry ABOVE `[1.8.0]`
- [ ] Step 4: `overseas/doc.go` — Phase 2.6 section 추가
- [ ] Step 5: commit

### CLAUDE.md 변경

```
# 변경 전
> **Phase 2.5 — 투자자/매매 동향 7 메서드 (v1.8.0).** Phase 2.5+ design spec 및 plan 참고.

# 변경 후
> **Phase 2.6 — 해외 정보 4 메서드 (v1.9.0).** Phase 2.5+ design spec 및 plan 참고.
```

plan link 추가 (기존 목록 맨 아래):
```
- Phase 2.6 implementation plan: [`docs/superpowers/specs/2026-05-05-phase2-6-overseas-info-implementation-plan.md`](docs/superpowers/specs/2026-05-05-phase2-6-overseas-info-implementation-plan.md)
```

### README.md 변경

Available Methods 표에 Phase 2.6 행 추가:

```markdown
| `Overseas.InquireNewsTitle`      | 해외뉴스종합(제목)       | HHPSTH60100C1  | Phase 2.6 |
| `Overseas.InquireBrknewsTitle`   | 해외속보(제목)           | FHKST01011801  | Phase 2.6 |
| `Overseas.InquireRightsByIce`    | 해외주식 권리종합        | HHDFS78330900  | Phase 2.6 |
| `Overseas.InquirePeriodRights`   | 해외주식 기간별권리조회  | CTRGT011R      | Phase 2.6 |
```

헤더 메서드 수: `60 methods` → `64 methods`

### CHANGELOG.md 변경

`[1.8.0]` ABOVE 위치에 신규 entry 추가:

```markdown
## [1.9.0] - 2026-05-06

### Added
- `Overseas.InquireNewsTitle` — 해외뉴스종합(제목) (HHPSTH60100C1)
  - ANOMALY: 응답 key `outblock1` (output/output1 아님), CTS pagination cursor
- `Overseas.InquireBrknewsTitle` — 해외속보(제목) (FHKST01011801)
  - ANOMALY: FID_ prefix 파라미터명, iscd1-10/kor_isnm1-10 flat 20 fields, FID_COND_SCR_DIV_CODE="11801" hardcoded
- `Overseas.InquireRightsByIce` — 해외주식 권리종합 (HHDFS78330900)
  - ANOMALY: output1 only (output2 없음)
- `Overseas.InquirePeriodRights` — 해외주식 기간별권리조회 (CTRGT011R)
  - ANOMALY: TR_ID C prefix, CTX_AREA_NK50/FK50 cursor pagination, numeric-content-as-String

### Notes
- 누적 메서드: 60 → 64
- 신규 파일: `overseas/news.go`, `overseas/rights.go`
```

### overseas/doc.go 변경

Phase 2.6 section 추가:

```go
// Phase 2.6 메서드 (4):
//
//   - InquireNewsTitle      — 해외뉴스종합(제목) (HHPSTH60100C1) — outblock1 key
//   - InquireBrknewsTitle   — 해외속보(제목) (FHKST01011801) — FID_ prefix params
//   - InquireRightsByIce    — 해외주식 권리종합 (HHDFS78330900) — output1 only
//   - InquirePeriodRights   — 해외주식 기간별권리조회 (CTRGT011R) — CTX cursor
```

### Commit

```bash
git commit -m "$(cat <<'EOF'
[docs] Phase 2.6 문서 갱신 (v1.9.0)

- CLAUDE.md: banner Phase 2.5 → Phase 2.6, plan link 추가
- README.md: Available Methods 4행 추가 (60 → 64 메서드)
- CHANGELOG.md: [1.9.0] entry (4 메서드 + anomaly 설명)
- overseas/doc.go: Phase 2.6 section (4 메서드)

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 8: 최종 점검

- [ ] Step 1: `go build ./...` — 에러 없음
- [ ] Step 2: `go vet ./...` — 경고 없음
- [ ] Step 3: `gofmt -l .` — 변경 파일 없음 (이미 fmt 적용)
- [ ] Step 4: `go test -race ./overseas/... -v` — 4 신규 테스트 PASS
- [ ] Step 5: coverage 확인

```bash
# build + vet
go build ./...
go vet ./...

# gofmt 검사
gofmt -l overseas/news.go overseas/news_test.go overseas/rights.go overseas/rights_test.go
# Expected: 출력 없음 (이미 fmt)

# race + coverage
go test -race -coverprofile=coverage.out ./overseas/...
go tool cover -func=coverage.out | grep total
# Expected: overseas 커버리지 ≥ 80%

go test -race -coverprofile=coverage_root.out ./...
go tool cover -func=coverage_root.out | grep total
# Expected: root 커버리지 ≥ 80%

# 신규 테스트 4개 명시적 실행
go test ./overseas/... -run "TestClient_InquireNewsTitle|TestClient_InquireBrknewsTitle|TestClient_InquireRightsByIce|TestClient_InquirePeriodRights" -v
# Expected: 4 PASS
```

- [ ] Step 6: commit (이상 발견 시 fix-up commit, 정상이면 skip)

---

## Task 9: PR 생성 (사용자 승인 후)

> **사용자 승인 후 진행.** 아래 단계는 승인 없이 실행 금지.

- [ ] Step 1: push feature branch

```bash
git push -u origin feat/phase2-6-overseas-info
```

- [ ] Step 2: PR 생성 (HEREDOC 방식 — MCP GitHub 도구 `\n` 리터럴 문제 회피)

```bash
gh pr create --title "feat: Phase 2.6 — 해외 정보 4 메서드 (v1.9.0)" --body "$(cat <<'EOF'
## Summary

Phase 2.6 — 해외 정보 4 메서드 추가 (누적 60 → 64, v1.9.0).

- `Overseas.InquireNewsTitle` — 해외뉴스종합(제목) (HHPSTH60100C1)
- `Overseas.InquireBrknewsTitle` — 해외속보(제목) (FHKST01011801)
- `Overseas.InquireRightsByIce` — 해외주식 권리종합 (HHDFS78330900)
- `Overseas.InquirePeriodRights` — 해외주식 기간별권리조회 (CTRGT011R)

## Anomalies documented

| EP | Anomaly |
|---|---|
| EP1 InquireNewsTitle | 응답 key `outblock1` (output/output1 아님), CTS pagination cursor |
| EP2 InquireBrknewsTitle | FID_ prefix 파라미터명, iscd1-10/kor_isnm1-10 flat 20 fields, FID_COND_SCR_DIV_CODE="11801" hardcoded |
| EP3 InquireRightsByIce | output1 only (output2 없음) |
| EP4 InquirePeriodRights | TR_ID C prefix (CTRGT011R), CTX_AREA_NK50/FK50 cursor pagination |

## Files changed

**New:**
- `overseas/news.go` (EP1+EP2)
- `overseas/news_test.go`
- `overseas/rights.go` (EP3+EP4)
- `overseas/rights_test.go`
- `overseas/testdata/news_title_success.json`
- `overseas/testdata/brknews_title_success.json`
- `overseas/testdata/rights_by_ice_success.json`
- `overseas/testdata/period_rights_success.json`
- `examples/overseas_info/main.go`

**Modified:**
- `CLAUDE.md`, `README.md`, `CHANGELOG.md`, `overseas/doc.go`

## Test plan

- [ ] `go build ./...` — clean
- [ ] `go vet ./...` — clean
- [ ] `go test -race ./overseas/... -v` — 4 신규 테스트 PASS
- [ ] overseas 커버리지 ≥ 80%
- [ ] root 커버리지 ≥ 80%
EOF
)" --reviewer kenshin579
```

- [ ] Step 3: PR 머지 (리뷰 완료 후)
- [ ] Step 4: tag + GitHub Release

```bash
git tag v1.9.0
git push origin v1.9.0

gh release create v1.9.0 \
  --title "v1.9.0 — Phase 2.6 해외 정보 4 메서드" \
  --notes "$(cat <<'EOF'
## Phase 2.6 — 해외 정보 (v1.9.0)

누적 메서드: 60 → **64**

### Added

- `Overseas.InquireNewsTitle` — 해외뉴스종합(제목) (HHPSTH60100C1)
- `Overseas.InquireBrknewsTitle` — 해외속보(제목) (FHKST01011801)
- `Overseas.InquireRightsByIce` — 해외주식 권리종합 (HHDFS78330900)
- `Overseas.InquirePeriodRights` — 해외주식 기간별권리조회 (CTRGT011R)

### Next

Phase 2.7 — 업종/지수 9 메서드 (v1.10.0)
EOF
)"
```

- [ ] Step 5: memory 갱신 — Phase 2.6 완료 (v1.9.0, 64 메서드), 다음 Phase 2.7 시작 준비

---

## Self-review checklist

- [x] 0 placeholders (TBD/TODO/Similar 없음)
- [x] 4 메서드 전체 커버: EP1 InquireNewsTitle, EP2 InquireBrknewsTitle, EP3 InquireRightsByIce, EP4 InquirePeriodRights
- [x] EP1 `Outblock1 []NewsTitleItem` + `json:"outblock1"` 태그 (NOT Output/Output1)
- [x] EP2 27 fields: iscd1-10 + kor_isnm1-10 flat (nested 배열 아님)
- [x] EP2 `FID_COND_SCR_DIV_CODE="11801"` hardcoded in Query map (Params 미노출)
- [x] EP3 `Output1 []RightsByIceItem` only (output2 없음)
- [x] HEREDOC commit messages 전체 적용
- [x] Task 9 "사용자 승인 후" 명시
