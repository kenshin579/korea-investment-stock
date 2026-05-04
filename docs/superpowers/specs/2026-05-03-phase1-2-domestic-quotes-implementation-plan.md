# Phase 1.2 — 국내 시세 + 심볼 + 차트 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** `korea-investment-stock` Go 라이브러리에 국내주식 시세/상품정보/차트/심볼 7 메서드 추가 (`v0.2.0` release).

**Architecture:** Phase 1.1 의 `internal/{httpclient,ratelimit,token,mastercache}` 인프라 위에 `domestic/` 패키지의 7 메서드 + `internal/krxmaster/` 패키지 (KOSPI/KOSDAQ 마스터 파일 cp949+fwf 파서) 추가. 한투 API 문서 1:1 매핑 (Style A — endpoint path 의 마지막 segment 를 PascalCase). TDD 흐름: testdata fixture (한투 docs 응답 필드 정의 → 합성 JSON, KRX 마스터는 실제 첫 3행) → 실패 테스트 → struct + 메서드 구현 → 통과 → commit.

**Tech Stack:** Go 1.25+, `golang.org/x/text/encoding/korean` (cp949 디코딩, 신규), `archive/zip`, `github.com/jarcoal/httpmock` (Phase 1.1 에 이미 추가), `github.com/shopspring/decimal`, `github.com/stretchr/testify`.

**참고 spec:**
- Phase 1 design spec (Phase 1.2 amendment 적용): `docs/superpowers/specs/2026-05-03-phase1-api-coverage-design.md` (commit `4dbaf51`)
- Phase 1.1 plan (참조 패턴): `docs/superpowers/specs/2026-05-03-phase1-1-infra-config-implementation-plan.md`
- 한투 API docs: `docs/api/국내주식/{주식현재가_시세.md, 상품기본조회.md, 주식기본조회.md, 국내주식기간별시세(일_주_월_년).md, 주식당일분봉조회.md}`

---

## 사전 정보

| 항목 | 값 |
|------|---|
| 작업 브랜치 | `docs/phase1-2-spec` (이미 생성됨) |
| 시작 HEAD | `4dbaf51` (Phase 1 design spec amendment commit) |
| Release 목표 | `v0.2.0` (PR merge 후 태그) |
| PR 베이스 | `main` |
| 현재 main 상태 | Phase 1.1 머지 완료 = 인프라 + Config 진입점 + `domestic/`, `overseas/` 의 placeholder Client struct |

---

## 메서드 → 한투 API 매핑

| Go 메서드 | 한투 path | TR_ID | docs |
|-----------|----------|-------|------|
| `Domestic.InquirePrice(ctx, symbol)` | `/uapi/domestic-stock/v1/quotations/inquire-price` | FHKST01010100 | 주식현재가_시세.md |
| `Domestic.SearchInfo(ctx, pdno, prdtTypeCD)` | `/uapi/domestic-stock/v1/quotations/search-info` | CTPF1604R | 상품기본조회.md |
| `Domestic.SearchStockInfo(ctx, pdno, prdtTypeCD)` | `/uapi/domestic-stock/v1/quotations/search-stock-info` | CTPF1002R | 주식기본조회.md |
| `Domestic.InquireDailyItemChartPrice(ctx, params)` | `/uapi/domestic-stock/v1/quotations/inquire-daily-itemchartprice` | FHKST03010100 | 국내주식기간별시세(일_주_월_년).md |
| `Domestic.InquireTimeItemChartPrice(ctx, params)` | `/uapi/domestic-stock/v1/quotations/inquire-time-itemchartprice` | FHKST03010200 | 주식당일분봉조회.md |
| `Domestic.FetchKospiSymbols(ctx)` | `https://new.real.download.dws.co.kr/common/master/kospi_code.mst.zip` (KRX 공개) | — | (한투 API 아님) |
| `Domestic.FetchKosdaqSymbols(ctx)` | `https://new.real.download.dws.co.kr/common/master/kosdaq_code.mst.zip` (KRX 공개) | — | (한투 API 아님) |

---

## 파일 구조

### 신규 (internal)

- `internal/krxmaster/doc.go` — 패키지 doc
- `internal/krxmaster/krxmaster.go` — `KospiSymbol`, `KosdaqSymbol` struct + `ParseKospi`, `ParseKosdaq` + `KospiURL`, `KosdaqURL` 상수 + cp949 + fwf helpers
- `internal/krxmaster/krxmaster_test.go`
- `internal/krxmaster/testdata/kospi_code_sample.mst.zip` (실제 KRX 첫 3행)
- `internal/krxmaster/testdata/kosdaq_code_sample.mst.zip` (실제 KRX 첫 3행)
- `internal/krxmaster/testdata/README.md` (출처 + 라이선스)

### 신규 (domestic)

- `domestic/price.go` — `InquirePrice` + `Price` struct
- `domestic/price_test.go`
- `domestic/info.go` — `SearchInfo` + `SearchStockInfo` + `ProductInfo` + `StockInfo` struct
- `domestic/info_test.go`
- `domestic/chart.go` — `InquireDailyItemChartPrice` + `InquireTimeItemChartPrice` + `DailyChart`, `MinuteChart`, params struct
- `domestic/chart_test.go`
- `domestic/symbols.go` — `FetchKospiSymbols`, `FetchKosdaqSymbols`, type alias re-export, `downloadURL` helper
- `domestic/symbols_test.go`
- `domestic/testhelper_test.go` — `newTestClient`, `loadFixture`, `stubTokenManager`
- `domestic/testdata/price_success.json`
- `domestic/testdata/product_info_success.json`
- `domestic/testdata/stock_info_success.json`
- `domestic/testdata/daily_chart_success.json`
- `domestic/testdata/minute_chart_success.json`
- `domestic/testdata/kospi_code_sample.mst.zip` (krxmaster 와 동일 파일 복사)
- `domestic/testdata/kosdaq_code_sample.mst.zip` (동일 복사)
- `domestic/testdata/README.md` (testdata 출처 명시)

### 신규 (examples)

- `examples/domestic_price/main.go`
- `examples/domestic_chart/main.go`
- `examples/kospi_symbols/main.go`

### 수정 (root)

- `go.mod` / `go.sum` — `golang.org/x/text` 추가
- `client.go` — `wireInfra` 안의 `c.Domestic = domestic.New(c.httpClient)` 를 `c.Domestic = domestic.New(c.httpClient, c.masterC)` 로 변경
- `CLAUDE.md` — Phase 1.2 메서드 안내
- `README.md` — Phase 1.2 메서드 사용 예시

### 수정 (sub-packages)

- `domestic/client.go` — `Client` struct 에 `master *mastercache.Cache` 추가, `New(http, master)` 시그니처 확장
- `domestic/doc.go` — Phase 1.2 메서드 안내 갱신

---

## Task 1: 의존성 추가 (`golang.org/x/text`)

**Files:**
- Modify: `go.mod`, `go.sum`

- [ ] **Step 1: 의존성 추가**

Run:
```bash
go get golang.org/x/text/encoding/korean
go mod tidy
```

- [ ] **Step 2: 검증**

Run: `go mod verify && grep 'golang.org/x/text' go.mod`
Expected: `all modules verified`. `go.mod` 에 `golang.org/x/text vX.Y.Z` 라인.

- [ ] **Step 3: Commit**

```bash
git add go.mod go.sum
git commit -m "$(cat <<'EOF'
[chore] Phase 1.2 의존성 추가 — golang.org/x/text

KOSPI/KOSDAQ 마스터 파일의 cp949 디코딩에 사용 (encoding/korean).

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 2: testdata — 실제 KRX 마스터 첫 3행 추출

**Files:**
- Create: `internal/krxmaster/testdata/kospi_code_sample.mst.zip`
- Create: `internal/krxmaster/testdata/kosdaq_code_sample.mst.zip`
- Create: `internal/krxmaster/testdata/README.md`
- Create: `domestic/testdata/kospi_code_sample.mst.zip` (위 파일 복사)
- Create: `domestic/testdata/kosdaq_code_sample.mst.zip` (위 파일 복사)
- Create: `domestic/testdata/README.md`

> **출처/라이선스**: KRX 마스터는 한국투자증권이 공개 다운로드로 제공하는 KRX 종목 마스터 파일 (`https://new.real.download.dws.co.kr/common/master/`). 본 testdata 는 첫 3 행만 추출한 sample 이며, 코드 변경 없이 라이브러리의 cp949+fwf 파서가 실제 KRX byte 와 호환되는지 검증하기 위함. README 에 출처 명시.

- [ ] **Step 1: KOSPI 마스터 다운로드 + 첫 3 행 추출**

Run:
```bash
mkdir -p /tmp/krxmaster && cd /tmp/krxmaster
curl -sSL -o kospi_code.mst.zip https://new.real.download.dws.co.kr/common/master/kospi_code.mst.zip
unzip -o kospi_code.mst.zip
head -n 3 kospi_code.mst > kospi_code_sample.mst
zip kospi_code_sample.mst.zip kospi_code_sample.mst
ls -la kospi_code_sample.mst.zip
```

Expected: `kospi_code_sample.mst.zip` 파일 생성 (수백 byte ~ 수 KB). `head -n 3 kospi_code.mst | wc -l` 출력이 3.

- [ ] **Step 2: KOSDAQ 마스터 다운로드 + 첫 3 행 추출**

Run:
```bash
cd /tmp/krxmaster
curl -sSL -o kosdaq_code.mst.zip https://new.real.download.dws.co.kr/common/master/kosdaq_code.mst.zip
unzip -o kosdaq_code.mst.zip
head -n 3 kosdaq_code.mst > kosdaq_code_sample.mst
zip kosdaq_code_sample.mst.zip kosdaq_code_sample.mst
ls -la kosdaq_code_sample.mst.zip
```

Expected: `kosdaq_code_sample.mst.zip` 파일 생성.

- [ ] **Step 3: testdata 디렉터리 생성 + 파일 복사**

Run:
```bash
mkdir -p internal/krxmaster/testdata domestic/testdata
cp /tmp/krxmaster/kospi_code_sample.mst.zip internal/krxmaster/testdata/
cp /tmp/krxmaster/kosdaq_code_sample.mst.zip internal/krxmaster/testdata/
cp /tmp/krxmaster/kospi_code_sample.mst.zip domestic/testdata/
cp /tmp/krxmaster/kosdaq_code_sample.mst.zip domestic/testdata/
```

- [ ] **Step 4: 검증 — sample 파일이 첫 3 행 포함하는지 확인**

Run:
```bash
unzip -p internal/krxmaster/testdata/kospi_code_sample.mst.zip | wc -l
unzip -p internal/krxmaster/testdata/kosdaq_code_sample.mst.zip | wc -l
```

Expected: 각각 `3` 출력. (`wc -l` 은 마지막 줄 newline 없으면 2 출력 가능. 그 경우 3 행 맞음 — `\n` 카운트 vs 줄 수 차이 인지)

- [ ] **Step 5: testdata README 작성** — `internal/krxmaster/testdata/README.md`

```markdown
# krxmaster testdata

`kospi_code_sample.mst.zip` 와 `kosdaq_code_sample.mst.zip` 는 KRX 종목 마스터 파일의 첫 3 행만 추출한 sample.

## 출처

- KOSPI 마스터: https://new.real.download.dws.co.kr/common/master/kospi_code.mst.zip
- KOSDAQ 마스터: https://new.real.download.dws.co.kr/common/master/kosdaq_code.mst.zip

한국투자증권이 공개 다운로드로 제공. `internal/krxmaster` 의 cp949+fwf 파서가 실제 KRX byte 와 호환되는지 검증하기 위한 단위 테스트 sample.

## 재생성 방법

```bash
cd /tmp && rm -rf krxmaster && mkdir krxmaster && cd krxmaster
curl -sSL -o kospi_code.mst.zip https://new.real.download.dws.co.kr/common/master/kospi_code.mst.zip
unzip -o kospi_code.mst.zip
head -n 3 kospi_code.mst > kospi_code_sample.mst
zip kospi_code_sample.mst.zip kospi_code_sample.mst
# kosdaq 동일
```

## 라이선스

KRX 종목 마스터는 한국투자증권이 무료 공개 다운로드로 배포. 본 sample 은 학습/테스트 용도이며, 라이브러리의 단위 테스트 외 사용 권장하지 않음.
```

- [ ] **Step 6: domestic/testdata README 작성** — `domestic/testdata/README.md`

```markdown
# domestic testdata

각 한투 API 메서드의 단위 테스트 fixture.

## REST API 응답 (합성 JSON)

- `price_success.json` — 주식현재가_시세 (FHKST01010100) 정상 응답
- `product_info_success.json` — 상품기본조회 (CTPF1604R) 정상 응답
- `stock_info_success.json` — 주식기본조회 (CTPF1002R) 정상 응답
- `daily_chart_success.json` — 국내주식기간별시세 (FHKST03010100) 정상 응답
- `minute_chart_success.json` — 주식당일분봉조회 (FHKST03010200) 정상 응답

각 JSON 의 필드는 `docs/api/국내주식/<API>.md` 의 응답 필드 정의에 1:1 매핑. 값은 합성 (실제 시세 아님).

## KRX 마스터 sample

- `kospi_code_sample.mst.zip`, `kosdaq_code_sample.mst.zip` — 출처 + 재생성 방법은 `internal/krxmaster/testdata/README.md` 참조
```

- [ ] **Step 7: Commit**

```bash
git add internal/krxmaster/testdata domestic/testdata/kospi_code_sample.mst.zip domestic/testdata/kosdaq_code_sample.mst.zip domestic/testdata/README.md
git commit -m "$(cat <<'EOF'
[chore] Phase 1.2 testdata — KRX 마스터 sample (첫 3행)

internal/krxmaster/testdata/ 와 domestic/testdata/ 양쪽에 KOSPI/KOSDAQ
마스터 ZIP 의 첫 3 행만 잘라낸 sample 추가. 단위 테스트의 cp949+fwf 파싱
검증용. 출처/라이선스/재생성 방법은 README.md.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 3: internal/krxmaster — KospiSymbol struct + ParseKospi

**Files:**
- Create: `internal/krxmaster/doc.go`
- Create: `internal/krxmaster/krxmaster.go`
- Create: `internal/krxmaster/krxmaster_test.go`

- [ ] **Step 1: doc.go**

```go
// Package krxmaster 는 KRX KOSPI/KOSDAQ 종목 마스터 파일의 디코딩/파싱 로직.
//
// 한투 API 가 아니라 KRX 가 공개 다운로드로 제공하는 .mst.zip 파일을 처리.
// cp949 인코딩 + fixed-width 컬럼 포맷이라 별도 파서 필요.
//
// 사용자에게 노출되지 않는 internal 패키지. domestic 패키지의 FetchKospiSymbols
// 와 FetchKosdaqSymbols 가 호출.
package krxmaster
```

- [ ] **Step 2: 테스트 작성** — `internal/krxmaster/krxmaster_test.go`

```go
package krxmaster

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseKospi(t *testing.T) {
	zipBytes, err := os.ReadFile("testdata/kospi_code_sample.mst.zip")
	require.NoError(t, err)

	syms, err := ParseKospi(zipBytes)
	require.NoError(t, err)
	require.Len(t, syms, 3, "sample 은 첫 3 행만 포함")

	// KRX 데이터는 시간에 따라 변동. 종목코드/한글명을 strict 비교 대신 패턴 검증.
	shortCodeRe := regexp.MustCompile(`^[0-9A-Z]{6}$`)
	hangulRe := regexp.MustCompile(`[\x{AC00}-\x{D7A3}]`)

	for i, s := range syms {
		assert.True(t, shortCodeRe.MatchString(s.ShortCode),
			"row %d: ShortCode %q 는 6자리 영숫자", i, s.ShortCode)
		assert.NotEmpty(t, s.StandardCode, "row %d: StandardCode 비어있음", i)
		assert.True(t, hangulRe.MatchString(s.KoreanName),
			"row %d: KoreanName %q 에 한글 포함", i, s.KoreanName)
		assert.NotEmpty(t, s.GroupCode, "row %d: GroupCode 비어있음", i)
		assert.NotNil(t, s.Raw, "row %d: Raw map 비어있지 않음", i)
		assert.GreaterOrEqual(t, len(s.Raw), 60,
			"row %d: Raw 에 ~70 컬럼 (최소 60+)", i)
	}
}

func TestParseKospi_InvalidZip(t *testing.T) {
	_, err := ParseKospi([]byte("not a zip"))
	assert.Error(t, err)
}

func TestParseKospi_EmptyZip(t *testing.T) {
	_, err := ParseKospi(nil)
	assert.Error(t, err)
}
```

- [ ] **Step 3: 테스트 실행 → FAIL**

Run: `go test ./internal/krxmaster/... -run TestParseKospi -v`
Expected: 컴파일 실패 (`ParseKospi`, `KospiSymbol` 미정의).

- [ ] **Step 4: 구현** — `internal/krxmaster/krxmaster.go`

```go
package krxmaster

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

// 한국투자증권 KRX 마스터 파일 다운로드 URL (var 로 노출 — 테스트에서 override 가능).
var (
	KospiURL  = "https://new.real.download.dws.co.kr/common/master/kospi_code.mst.zip"
	KosdaqURL = "https://new.real.download.dws.co.kr/common/master/kosdaq_code.mst.zip"
)

// KospiSymbol 은 KRX KOSPI 종목 마스터 (kospi_code.mst) 한 행.
//
// 핵심 필드 typed + Raw map[한글컬럼명]값 으로 미typed 70 컬럼 fallback.
// docs: 한국투자 GitHub open-trading-api/stocks_info 의 Python 정제 코드 참조.
type KospiSymbol struct {
	ShortCode       string          // 단축코드 (예 "005930")
	StandardCode    string          // 표준코드 (ISIN, 예 "KR7005930003")
	KoreanName      string          // 한글명
	GroupCode       string          // 그룹코드 (ST=주권, EF=ETF, RT=REITs, IF=인프라, ...)
	MarketCapSize   string          // 시가총액 규모
	KOSPI200        string          // KOSPI200 섹터업종 편입 (Y/N)
	KOSPI100        string          // KOSPI100 편입
	KOSPI50         string          // KOSPI50 편입
	BasePrice       int64           // 기준가
	FaceValue       decimal.Decimal // 액면가
	ListedShares    int64           // 상장주수
	Capital         int64           // 자본금
	SettlementMonth string          // 결산월 (예 "12")
	PreferredStock  string          // 우선주 여부 (Y/N)
	SuspendedYn     string          // 거래정지 여부 (Y/N)

	Raw map[string]string // 모든 70 컬럼 (한글 키)
}

// kospiFieldSpecs 와 kospiColumns 는 한투 GitHub Python parsers/master_parser.py 의
// field_specs / part2_columns 와 1:1 매핑. 라인 마지막 228 byte 의 fwf 영역.
var kospiFieldSpecs = []int{
	2, 1, 4, 4, 4,
	1, 1, 1, 1, 1,
	1, 1, 1, 1, 1,
	1, 1, 1, 1, 1,
	1, 1, 1, 1, 1,
	1, 1, 1, 1, 1,
	1, 9, 5, 5, 1,
	1, 1, 2, 1, 1,
	1, 2, 2, 2, 3,
	1, 3, 12, 12, 8,
	15, 21, 2, 7, 1,
	1, 1, 1, 1, 9,
	9, 9, 5, 9, 8,
	9, 3, 1, 1, 1,
}

var kospiColumns = []string{
	"그룹코드", "시가총액규모", "지수업종대분류", "지수업종중분류", "지수업종소분류",
	"제조업", "저유동성", "지배구조지수종목", "KOSPI200섹터업종", "KOSPI100",
	"KOSPI50", "KRX", "ETP", "ELW발행", "KRX100",
	"KRX자동차", "KRX반도체", "KRX바이오", "KRX은행", "SPAC",
	"KRX에너지화학", "KRX철강", "단기과열", "KRX미디어통신", "KRX건설",
	"Non1", "KRX증권", "KRX선박", "KRX섹터_보험", "KRX섹터_운송",
	"SRI", "기준가", "매매수량단위", "시간외수량단위", "거래정지",
	"정리매매", "관리종목", "시장경고", "경고예고", "불성실공시",
	"우회상장", "락구분", "액면변경", "증자구분", "증거금비율",
	"신용가능", "신용기간", "전일거래량", "액면가", "상장일자",
	"상장주수", "자본금", "결산월", "공모가", "우선주",
	"공매도과열", "이상급등", "KRX300", "KOSPI", "매출액",
	"영업이익", "경상이익", "당기순이익", "ROE", "기준년월",
	"시가총액", "그룹사코드", "회사신용한도초과", "담보대출가능", "대주가능",
}

// ParseKospi 는 KOSPI 마스터 ZIP 의 byte 를 받아 종목 슬라이스로 디코딩.
func ParseKospi(zipBytes []byte) ([]KospiSymbol, error) {
	const fwfLen = 228
	mst, err := openMstFromZip(zipBytes)
	if err != nil {
		return nil, fmt.Errorf("krxmaster: kospi: %w", err)
	}
	decoded, err := decodeCP949(mst)
	if err != nil {
		return nil, fmt.Errorf("krxmaster: kospi: cp949: %w", err)
	}

	var out []KospiSymbol
	for _, line := range strings.Split(decoded, "\n") {
		line = strings.TrimRight(line, "\r")
		if len(line) < fwfLen+21 {
			continue
		}
		prefix := line[:len(line)-fwfLen]
		fwf := line[len(line)-fwfLen:]

		shortCode := strings.TrimSpace(prefix[0:9])
		standardCode := strings.TrimSpace(prefix[9:21])
		koreanName := strings.TrimSpace(prefix[21:])
		raw := parseFwf(fwf, kospiFieldSpecs, kospiColumns)

		out = append(out, KospiSymbol{
			ShortCode:       shortCode,
			StandardCode:    standardCode,
			KoreanName:      koreanName,
			GroupCode:       raw["그룹코드"],
			MarketCapSize:   raw["시가총액규모"],
			KOSPI200:        raw["KOSPI200섹터업종"],
			KOSPI100:        raw["KOSPI100"],
			KOSPI50:         raw["KOSPI50"],
			BasePrice:       atoi64(raw["기준가"]),
			FaceValue:       toDecimal(raw["액면가"]),
			ListedShares:    atoi64(raw["상장주수"]),
			Capital:         atoi64(raw["자본금"]),
			SettlementMonth: raw["결산월"],
			PreferredStock:  raw["우선주"],
			SuspendedYn:     raw["거래정지"],
			Raw:             raw,
		})
	}
	return out, nil
}

// openMstFromZip 은 ZIP byte 에서 .mst 파일 byte 를 추출.
func openMstFromZip(zipBytes []byte) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return nil, fmt.Errorf("zip open: %w", err)
	}
	for _, f := range zr.File {
		if strings.HasSuffix(f.Name, ".mst") {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("zip read %s: %w", f.Name, err)
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
	}
	return nil, fmt.Errorf(".mst file not found in zip")
}

// decodeCP949 는 cp949 byte 를 UTF-8 string 으로 변환.
// golang.org/x/text/encoding/korean.EUCKR 는 cp949 호환 (Microsoft 확장 포함).
func decodeCP949(b []byte) (string, error) {
	decoded, _, err := transform.Bytes(korean.EUCKR.NewDecoder(), b)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// parseFwf 는 fixed-width line 을 widths 별로 잘라 한글 컬럼명 → 값 map 반환.
func parseFwf(line string, widths []int, names []string) map[string]string {
	out := make(map[string]string, len(names))
	pos := 0
	for i, w := range widths {
		if pos+w > len(line) {
			break
		}
		out[names[i]] = strings.TrimSpace(line[pos : pos+w])
		pos += w
	}
	return out
}

// atoi64 는 빈 문자열/공백 → 0 fallback. 한투 마스터 데이터에 빈 값 흔함.
func atoi64(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return v
}

// toDecimal 은 빈 문자열 → decimal.Zero fallback.
func toDecimal(s string) decimal.Decimal {
	s = strings.TrimSpace(s)
	if s == "" {
		return decimal.Zero
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero
	}
	return d
}
```

- [ ] **Step 5: 테스트 실행 → PASS**

Run: `go test ./internal/krxmaster/... -run TestParseKospi -v`
Expected: 3 test cases PASS (정상 + InvalidZip + EmptyZip).

- [ ] **Step 6: Commit**

```bash
git add internal/krxmaster/doc.go internal/krxmaster/krxmaster.go internal/krxmaster/krxmaster_test.go
git commit -m "$(cat <<'EOF'
[feat] internal/krxmaster — KOSPI 마스터 파서 (cp949 + fwf)

KospiSymbol struct + ParseKospi. 핵심 15 필드 typed + Raw map[한글컬럼명]값.
한투 GitHub open-trading-api 의 Python parsers/master_parser.py 의 field_specs
와 1:1 매핑. golang.org/x/text/encoding/korean.EUCKR 로 cp949 디코딩.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 4: internal/krxmaster — ParseKosdaq

**Files:**
- Modify: `internal/krxmaster/krxmaster.go`
- Modify: `internal/krxmaster/krxmaster_test.go`

- [ ] **Step 1: 테스트 추가** — `internal/krxmaster/krxmaster_test.go` 에 함수 추가

```go
func TestParseKosdaq(t *testing.T) {
	zipBytes, err := os.ReadFile("testdata/kosdaq_code_sample.mst.zip")
	require.NoError(t, err)

	syms, err := ParseKosdaq(zipBytes)
	require.NoError(t, err)
	require.Len(t, syms, 3)

	shortCodeRe := regexp.MustCompile(`^[0-9A-Z]{6}$`)
	hangulRe := regexp.MustCompile(`[\x{AC00}-\x{D7A3}]`)

	for i, s := range syms {
		assert.True(t, shortCodeRe.MatchString(s.ShortCode),
			"row %d: ShortCode %q 는 6자리 영숫자", i, s.ShortCode)
		assert.NotEmpty(t, s.StandardCode, "row %d: StandardCode", i)
		assert.True(t, hangulRe.MatchString(s.KoreanName),
			"row %d: KoreanName 한글", i)
		assert.NotEmpty(t, s.GroupCode, "row %d: GroupCode", i)
		assert.NotNil(t, s.Raw, "row %d: Raw map", i)
		assert.GreaterOrEqual(t, len(s.Raw), 50, "row %d: Raw 60+ 컬럼", i)
	}
}
```

- [ ] **Step 2: 테스트 실행 → FAIL**

Run: `go test ./internal/krxmaster/... -run TestParseKosdaq -v`
Expected: 컴파일 실패 (`ParseKosdaq`, `KosdaqSymbol` 미정의).

- [ ] **Step 3: 구현 추가** — `internal/krxmaster/krxmaster.go` 에 추가

```go
// KosdaqSymbol 은 KRX KOSDAQ 종목 마스터 한 행.
type KosdaqSymbol struct {
	ShortCode       string
	StandardCode    string
	KoreanName      string
	GroupCode       string
	MarketCapSize   string
	KOSDAQ150       string // KOSDAQ150 편입 여부
	BasePrice       int64
	FaceValue       decimal.Decimal
	ListedShares    int64
	Capital         int64
	SettlementMonth string
	PreferredStock  string
	VentureCompany  string // 벤처기업 여부
	SuspendedYn     string

	Raw map[string]string // 모든 64 컬럼
}

// kosdaqFieldSpecs / kosdaqColumns 는 Python parsers/master_parser.py 의
// KOSDAQ field_specs / part2_columns 와 1:1 매핑. 마지막 222 byte fwf.
var kosdaqFieldSpecs = []int{
	2, 1, 4, 4, 4,
	1, 1, 1, 1, 1,
	1, 1, 1, 1, 1,
	1, 1, 1, 1, 1,
	1, 1, 1, 1, 1,
	1, 9, 5, 5, 1,
	1, 1, 2, 1, 1,
	1, 2, 2, 2, 3,
	1, 3, 12, 12, 8,
	15, 21, 2, 7, 1,
	1, 1, 1, 9, 9,
	9, 5, 9, 8, 9,
	3, 1, 1, 1,
}

var kosdaqColumns = []string{
	"그룹코드", "시가총액규모", "지수업종대분류", "지수업종중분류", "지수업종소분류",
	"벤처기업", "저유동성", "KRX", "ETP", "KRX100",
	"KRX자동차", "KRX반도체", "KRX바이오", "KRX은행", "SPAC",
	"KRX에너지화학", "KRX철강", "단기과열", "KRX미디어통신", "KRX건설",
	"투자주의", "KRX증권", "KRX선박", "KRX섹터_보험", "KRX섹터_운송",
	"KOSDAQ150", "기준가", "매매수량단위", "시간외수량단위", "거래정지",
	"정리매매", "관리종목", "시장경고", "경고예고", "불성실공시",
	"우회상장", "락구분", "액면변경", "증자구분", "증거금비율",
	"신용가능", "신용기간", "전일거래량", "액면가", "상장일자",
	"상장주수", "자본금", "결산월", "공모가", "우선주",
	"공매도과열", "이상급등", "KRX300", "매출액", "영업이익",
	"경상이익", "당기순이익", "ROE", "기준년월", "시가총액",
	"그룹사코드", "회사신용한도초과", "담보대출가능", "대주가능",
}

// ParseKosdaq 는 KOSDAQ 마스터 ZIP byte 를 종목 슬라이스로 디코딩.
func ParseKosdaq(zipBytes []byte) ([]KosdaqSymbol, error) {
	const fwfLen = 222
	mst, err := openMstFromZip(zipBytes)
	if err != nil {
		return nil, fmt.Errorf("krxmaster: kosdaq: %w", err)
	}
	decoded, err := decodeCP949(mst)
	if err != nil {
		return nil, fmt.Errorf("krxmaster: kosdaq: cp949: %w", err)
	}

	var out []KosdaqSymbol
	for _, line := range strings.Split(decoded, "\n") {
		line = strings.TrimRight(line, "\r")
		if len(line) < fwfLen+21 {
			continue
		}
		prefix := line[:len(line)-fwfLen]
		fwf := line[len(line)-fwfLen:]

		shortCode := strings.TrimSpace(prefix[0:9])
		standardCode := strings.TrimSpace(prefix[9:21])
		koreanName := strings.TrimSpace(prefix[21:])
		raw := parseFwf(fwf, kosdaqFieldSpecs, kosdaqColumns)

		out = append(out, KosdaqSymbol{
			ShortCode:       shortCode,
			StandardCode:    standardCode,
			KoreanName:      koreanName,
			GroupCode:       raw["그룹코드"],
			MarketCapSize:   raw["시가총액규모"],
			KOSDAQ150:       raw["KOSDAQ150"],
			BasePrice:       atoi64(raw["기준가"]),
			FaceValue:       toDecimal(raw["액면가"]),
			ListedShares:    atoi64(raw["상장주수"]),
			Capital:         atoi64(raw["자본금"]),
			SettlementMonth: raw["결산월"],
			PreferredStock:  raw["우선주"],
			VentureCompany:  raw["벤처기업"],
			SuspendedYn:     raw["거래정지"],
			Raw:             raw,
		})
	}
	return out, nil
}
```

- [ ] **Step 4: 테스트 실행 → PASS**

Run: `go test ./internal/krxmaster/... -v`
Expected: `TestParseKospi`, `TestParseKospi_InvalidZip`, `TestParseKospi_EmptyZip`, `TestParseKosdaq` 모두 PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/krxmaster/krxmaster.go internal/krxmaster/krxmaster_test.go
git commit -m "$(cat <<'EOF'
[feat] internal/krxmaster — KOSDAQ 파서 추가

KosdaqSymbol struct + ParseKosdaq. KOSPI 와 같은 패턴, 마지막 222 byte fwf
+ 64 컬럼. 벤처기업/KOSDAQ150 등 KOSDAQ 전용 컬럼 typed 노출.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 5: domestic.Client 시그니처 확장 + root wireInfra 갱신

**Files:**
- Modify: `domestic/client.go`
- Modify: `client.go` (root)

- [ ] **Step 1: domestic/client.go — Client 에 master 필드 추가**

```go
package domestic

import (
	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
	"github.com/kenshin579/korea-investment-stock/internal/mastercache"
)

// Client 는 국내주식 API sub-client.
//
// 사용자는 직접 생성하지 않고 kis.Client.Domestic 으로 접근.
type Client struct {
	http   *httpclient.Client
	master *mastercache.Cache // KRX 마스터 파일 디스크 캐시 (FetchKospi/Kosdaq Symbols 가 사용)
}

// New 는 internal 용도. root kis.NewClient 가 호출.
func New(http *httpclient.Client, master *mastercache.Cache) *Client {
	return &Client{http: http, master: master}
}
```

- [ ] **Step 2: client.go (root) — wireInfra 의 Domestic 생성 호출 갱신**

`client.go` 안의 `c.Domestic = domestic.New(c.httpClient)` 라인을 다음으로 교체:

```go
	c.Domestic = domestic.New(c.httpClient, c.masterC)
```

(`c.Overseas = overseas.New(c.httpClient)` 는 Phase 1.5 에서 master 주입.)

- [ ] **Step 3: 컴파일 검증**

Run: `go build ./...`
Expected: 빌드 성공.

- [ ] **Step 4: 기존 테스트 회귀 확인**

Run: `go test ./... -count=1`
Expected: 모든 패키지 PASS (Phase 1.1 의 client_test.go 등 회귀 없음).

- [ ] **Step 5: Commit**

```bash
git add domestic/client.go client.go
git commit -m "$(cat <<'EOF'
[refactor] domestic.Client 에 master *mastercache.Cache 필드 추가

Phase 1.2 의 FetchKospi/Kosdaq Symbols 가 마스터 캐시에 접근하기 위한 인프라.
domestic.New(http, master) 시그니처 변경 — internal use 만이라 BC 영향 없음.
root client.go 의 wireInfra 가 c.masterC 주입.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 6: domestic/testhelper_test.go — 공통 테스트 helper

**Files:**
- Create: `domestic/testhelper_test.go`

- [ ] **Step 1: testhelper 작성** — `domestic/testhelper_test.go`

```go
package domestic_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/domestic"
	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
	"github.com/kenshin579/korea-investment-stock/internal/mastercache"
	"github.com/kenshin579/korea-investment-stock/internal/ratelimit"
)

const testBaseURL = "https://openapi.koreainvestment.com:9443"

// loadFixture 는 testdata/<name> 파일 byte 를 로드.
func loadFixture(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	require.NoError(t, err)
	return b
}

// loadFixtureString 은 string 으로 로드 (httpmock.NewStringResponder 용).
func loadFixtureString(t *testing.T, name string) string {
	return string(loadFixture(t, name))
}

// stubTokenManager 는 httpclient.TokenManager 의 stub. 항상 "Bearer test" 반환.
type stubTokenManager struct{}

func (stubTokenManager) Get(ctx context.Context) (string, error)     { return "Bearer test", nil }
func (stubTokenManager) Refresh(ctx context.Context) (string, error) { return "Bearer test", nil }

// newTestClient 는 httpmock 활성 상태에서 사용할 domestic.Client 생성.
// 호출자는 httpmock.Activate() / httpmock.DeactivateAndReset() 직접 관리.
func newTestClient(t *testing.T) *domestic.Client {
	t.Helper()
	httpClient := &http.Client{Transport: httpmock.DefaultTransport}
	httpcli := httpclient.New(httpclient.Config{
		BaseURL:    testBaseURL,
		AppKey:     "test-key",
		AppSecret:  "test-secret",
		AccountNo:  "00000000-00",
		Limiter:    ratelimit.New(1000),
		TokenMgr:   stubTokenManager{},
		Retries:    0,
		Timeout:    5 * time.Second,
		HTTPClient: httpClient,
	})
	master := mastercache.New(t.TempDir(), time.Hour)
	return domestic.New(httpcli, master)
}
```

- [ ] **Step 2: 컴파일 검증**

Run: `go test ./domestic/... -run NoSuchTest -v`
Expected: 빌드 성공 (테스트 0개 실행되지만 컴파일 에러 없음).

- [ ] **Step 3: Commit**

```bash
git add domestic/testhelper_test.go
git commit -m "$(cat <<'EOF'
[test] domestic — 공통 테스트 helper

newTestClient (httpmock + stub token manager + temp master cache),
loadFixture, stubTokenManager. Phase 1.2 의 메서드 테스트가 공통 사용.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 7: domestic/symbols.go — FetchKospiSymbols + FetchKosdaqSymbols

**Files:**
- Create: `domestic/symbols.go`
- Create: `domestic/symbols_test.go`

- [ ] **Step 1: 테스트 작성** — `domestic/symbols_test.go`

```go
package domestic_test

import (
	"context"
	"net/http"
	"regexp"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/internal/krxmaster"
)

func TestClient_FetchKospiSymbols(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	zipBytes := loadFixture(t, "kospi_code_sample.mst.zip")
	httpmock.RegisterResponder(http.MethodGet, krxmaster.KospiURL,
		httpmock.NewBytesResponder(200, zipBytes))

	c := newTestClient(t)
	syms, err := c.FetchKospiSymbols(context.Background())
	require.NoError(t, err)
	require.Len(t, syms, 3)

	hangulRe := regexp.MustCompile(`[\x{AC00}-\x{D7A3}]`)
	for i, s := range syms {
		assert.Regexp(t, `^[0-9A-Z]{6}$`, s.ShortCode, "row %d ShortCode", i)
		assert.True(t, hangulRe.MatchString(s.KoreanName), "row %d KoreanName 한글", i)
	}
}

func TestClient_FetchKosdaqSymbols(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	zipBytes := loadFixture(t, "kosdaq_code_sample.mst.zip")
	httpmock.RegisterResponder(http.MethodGet, krxmaster.KosdaqURL,
		httpmock.NewBytesResponder(200, zipBytes))

	c := newTestClient(t)
	syms, err := c.FetchKosdaqSymbols(context.Background())
	require.NoError(t, err)
	require.Len(t, syms, 3)
	for _, s := range syms {
		assert.NotEmpty(t, s.ShortCode)
		assert.NotEmpty(t, s.KoreanName)
	}
}

func TestClient_FetchKospiSymbols_DownloadError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, krxmaster.KospiURL,
		httpmock.NewStringResponder(500, "internal error"))

	c := newTestClient(t)
	_, err := c.FetchKospiSymbols(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP 500")
}
```

- [ ] **Step 2: 테스트 실행 → FAIL**

Run: `go test ./domestic/... -run FetchKospi -v`
Expected: 컴파일 실패 (`FetchKospiSymbols` 미정의).

- [ ] **Step 3: 구현** — `domestic/symbols.go`

```go
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
```

- [ ] **Step 4: 테스트 실행 → PASS**

Run: `go test ./domestic/... -run FetchKospi -run FetchKosdaq -v`
Expected: 3 test cases PASS.

- [ ] **Step 5: Commit**

```bash
git add domestic/symbols.go domestic/symbols_test.go
git commit -m "$(cat <<'EOF'
[feat] domestic — FetchKospiSymbols + FetchKosdaqSymbols

KRX 공개 마스터 파일 다운로드 + mastercache 디스크 캐시 + krxmaster 파싱.
http.DefaultClient 로 KRX 도메인 직접 호출 (한투 API transport 와 분리).
KospiSymbol/KosdaqSymbol 은 krxmaster 의 type alias 로 외부 노출.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 8: domestic/price.go — InquirePrice + Price struct

**Files:**
- Create: `domestic/price.go`
- Create: `domestic/price_test.go`
- Create: `domestic/testdata/price_success.json`

- [ ] **Step 1: testdata fixture 작성** — `domestic/testdata/price_success.json`

> 한투 docs (`docs/api/국내주식/주식현재가_시세.md`) 의 응답 필드 정의 기반 합성. 핵심 필드 + 단위 테스트가 검증할 값.

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": {
    "iscd_stat_cls_code": "55",
    "marg_rate": "20.00",
    "rprs_mrkt_kor_name": "KOSPI200",
    "new_hgpr_lwpr_cls_code": "",
    "bstp_kor_isnm": "전기.전자",
    "temp_stop_yn": "N",
    "oprc_rang_cont_yn": "N",
    "clpr_rang_cont_yn": "N",
    "crdt_able_yn": "Y",
    "grmn_rate_cls_code": "40",
    "elw_pblc_yn": "Y",
    "stck_prpr": "75800",
    "prdy_vrss": "-200",
    "prdy_vrss_sign": "5",
    "prdy_ctrt": "-0.26",
    "acml_tr_pbmn": "938223456000",
    "acml_vol": "12345678",
    "prdy_vrss_vol_rate": "85.42",
    "stck_oprc": "76000",
    "stck_hgpr": "76200",
    "stck_lwpr": "75500",
    "stck_mxpr": "98500",
    "stck_llam": "53100",
    "stck_sdpr": "76000",
    "wghn_avrg_stck_prc": "75865.21",
    "hts_frgn_ehrt": "53.42",
    "frgn_ntby_qty": "-123456",
    "pgtr_ntby_qty": "-89000",
    "pvt_scnd_dmrs_prc": "0",
    "pvt_frst_dmrs_prc": "0",
    "pvt_pont_val": "0",
    "pvt_frst_dmsp_prc": "0",
    "pvt_scnd_dmsp_prc": "0",
    "dmrs_val": "0",
    "dmsp_val": "0",
    "cpfn": "7780383",
    "rstc_wdth_prc": "22700",
    "stck_fcam": "100",
    "stck_sspr": "60640",
    "aspr_unit": "100",
    "hts_deal_qty_unit_val": "1",
    "lstn_stcn": "5969782550",
    "hts_avls": "452312",
    "per": "11.42",
    "pbr": "1.32",
    "stac_month": "12",
    "vol_tnrt": "0.21",
    "eps": "6638",
    "bps": "57420",
    "d250_hgpr": "88000",
    "d250_hgpr_date": "20251015",
    "d250_hgpr_vrss_prpr_rate": "-13.86",
    "d250_lwpr": "65000",
    "d250_lwpr_date": "20260120",
    "d250_lwpr_vrss_prpr_rate": "16.62",
    "stck_dryy_hgpr": "88000",
    "dryy_hgpr_vrss_prpr_rate": "-13.86",
    "dryy_hgpr_date": "20251015",
    "stck_dryy_lwpr": "65000",
    "dryy_lwpr_vrss_prpr_rate": "16.62",
    "dryy_lwpr_date": "20260120",
    "w52_hgpr": "88000",
    "w52_hgpr_vrss_prpr_ctrt": "-13.86",
    "w52_hgpr_date": "20251015",
    "w52_lwpr": "65000",
    "w52_lwpr_vrss_prpr_ctrt": "16.62",
    "w52_lwpr_date": "20260120",
    "whol_loan_rmnd_rate": "0.45",
    "ssts_yn": "Y",
    "stck_shrn_iscd": "005930",
    "fcam_cnnm": "원",
    "cpfn_cnnm": "원",
    "apprch_rate": "70.50",
    "frgn_hldn_qty": "3187654321",
    "vi_cls_code": "N",
    "ovtm_vi_cls_code": "N",
    "last_ssts_cntg_qty": "0",
    "invt_caful_yn": "N",
    "mrkt_warn_cls_code": "00",
    "short_over_yn": "N",
    "sltr_yn": "N",
    "mang_issu_cls_code": "N"
  }
}
```

- [ ] **Step 2: 테스트 작성** — `domestic/price_test.go`

```go
package domestic_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_InquirePrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-price`,
		httpmock.NewStringResponder(200, loadFixtureString(t, "price_success.json")),
	)

	c := newTestClient(t)
	price, err := c.InquirePrice(context.Background(), "005930")
	require.NoError(t, err)
	require.NotNil(t, price)

	assert.Equal(t, decimal.NewFromInt(75800), price.StckPrpr)
	assert.Equal(t, decimal.NewFromInt(-200), price.PrdyVrss)
	assert.Equal(t, "5", price.PrdyVrssSign)
	assert.InDelta(t, -0.26, price.PrdyCtrt, 0.001)
	assert.Equal(t, int64(12345678), price.AcmlVol)
	assert.Equal(t, int64(938223456000), price.AcmlTrPbmn)
	assert.Equal(t, decimal.NewFromInt(76000), price.StckOprc)
	assert.Equal(t, decimal.NewFromInt(76200), price.StckHgpr)
	assert.Equal(t, decimal.NewFromInt(75500), price.StckLwpr)
	assert.Equal(t, "005930", price.StckShrnIscd)
	assert.Equal(t, "Y", price.SstsYn)
	assert.Equal(t, "N", price.MangIssuClsCode)
	assert.InDelta(t, 11.42, price.Per, 0.001)
	assert.InDelta(t, 1.32, price.Pbr, 0.001)
}

func TestClient_InquirePrice_APIError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-price`,
		httpmock.NewStringResponder(200, `{"rt_cd":"1","msg_cd":"MCA00001","msg1":"잘못된 요청","output":null}`),
	)

	c := newTestClient(t)
	_, err := c.InquirePrice(context.Background(), "005930")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "MCA00001")
}
```

- [ ] **Step 3: 테스트 실행 → FAIL**

Run: `go test ./domestic/... -run InquirePrice -v`
Expected: 컴파일 실패 (`InquirePrice`, `Price` 미정의).

- [ ] **Step 4: 구현** — `domestic/price.go`

```go
package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// Price 는 주식현재가_시세 (FHKST01010100) 의 output 응답.
//
// 한투 docs: docs/api/국내주식/주식현재가_시세.md
// path: /uapi/domestic-stock/v1/quotations/inquire-price
type Price struct {
	IscdStatClsCode      string          `json:"iscd_stat_cls_code"`             // 종목 상태 (51=관리, 52=투자위험, 58=거래정지, ...)
	MargRate             float64         `json:"marg_rate,string"`               // 증거금 비율
	RprsMrktKorName      string          `json:"rprs_mrkt_kor_name"`             // 대표 시장 한글명
	NewHgprLwprClsCode   string          `json:"new_hgpr_lwpr_cls_code"`         // 신 고가/저가 구분
	BstpKorIsnm          string          `json:"bstp_kor_isnm"`                  // 업종 한글명
	TempStopYn           string          `json:"temp_stop_yn"`                   // 임시 정지 (Y/N)
	OprcRangContYn       string          `json:"oprc_rang_cont_yn"`              // 시가 범위 연장
	ClprRangContYn       string          `json:"clpr_rang_cont_yn"`              // 종가 범위 연장
	CrdtAbleYn           string          `json:"crdt_able_yn"`                   // 신용 가능 (Y/N)
	GrmnRateClsCode      string          `json:"grmn_rate_cls_code"`             // 보증금 비율 구분
	ElwPblcYn            string          `json:"elw_pblc_yn"`                    // ELW 발행 여부 (Y/N)
	StckPrpr             decimal.Decimal `json:"stck_prpr"`                      // 주식 현재가
	PrdyVrss             decimal.Decimal `json:"prdy_vrss"`                      // 전일 대비
	PrdyVrssSign         string          `json:"prdy_vrss_sign"`                 // 전일 대비 부호 (1=상한, 2=상승, 3=보합, 4=하한, 5=하락)
	PrdyCtrt             float64         `json:"prdy_ctrt,string"`               // 전일 대비율
	AcmlTrPbmn           int64           `json:"acml_tr_pbmn,string"`            // 누적 거래 대금
	AcmlVol              int64           `json:"acml_vol,string"`                // 누적 거래량
	PrdyVrssVolRate      float64         `json:"prdy_vrss_vol_rate,string"`      // 전일 대비 거래량 비율
	StckOprc             decimal.Decimal `json:"stck_oprc"`                      // 시가
	StckHgpr             decimal.Decimal `json:"stck_hgpr"`                      // 고가
	StckLwpr             decimal.Decimal `json:"stck_lwpr"`                      // 저가
	StckMxpr             decimal.Decimal `json:"stck_mxpr"`                      // 상한가
	StckLlam             decimal.Decimal `json:"stck_llam"`                      // 하한가
	StckSdpr             decimal.Decimal `json:"stck_sdpr"`                      // 기준가
	WghnAvrgStckPrc      decimal.Decimal `json:"wghn_avrg_stck_prc"`             // 가중 평균 주식 가격
	HtsFrgnEhrt          float64         `json:"hts_frgn_ehrt,string"`           // HTS 외국인 소진율
	FrgnNtbyQty          int64           `json:"frgn_ntby_qty,string"`           // 외국인 순매수 수량
	PgtrNtbyQty          int64           `json:"pgtr_ntby_qty,string"`           // 프로그램매매 순매수 수량
	PvtScndDmrsPrc       decimal.Decimal `json:"pvt_scnd_dmrs_prc"`              // 피벗 2차 디저항 가격
	PvtFrstDmrsPrc       decimal.Decimal `json:"pvt_frst_dmrs_prc"`              // 피벗 1차 디저항 가격
	PvtPontVal           decimal.Decimal `json:"pvt_pont_val"`                   // 피벗 포인트 값
	PvtFrstDmspPrc       decimal.Decimal `json:"pvt_frst_dmsp_prc"`              // 피벗 1차 디지지
	PvtScndDmspPrc       decimal.Decimal `json:"pvt_scnd_dmsp_prc"`              // 피벗 2차 디지지
	DmrsVal              decimal.Decimal `json:"dmrs_val"`                       // 디저항 값
	DmspVal              decimal.Decimal `json:"dmsp_val"`                       // 디지지 값
	Cpfn                 int64           `json:"cpfn,string"`                    // 자본금 (백만원)
	RstcWdthPrc          decimal.Decimal `json:"rstc_wdth_prc"`                  // 제한 폭 가격
	StckFcam             decimal.Decimal `json:"stck_fcam"`                      // 주식 액면가
	StckSspr             decimal.Decimal `json:"stck_sspr"`                      // 주식 대용가
	AsprUnit             decimal.Decimal `json:"aspr_unit"`                      // 호가 단위
	HtsDealQtyUnitVal    int64           `json:"hts_deal_qty_unit_val,string"`   // HTS 매매 수량 단위
	LstnStcn             int64           `json:"lstn_stcn,string"`               // 상장 주수
	HtsAvls              int64           `json:"hts_avls,string"`                // HTS 시가총액 (억원)
	Per                  float64         `json:"per,string"`                     // PER
	Pbr                  float64         `json:"pbr,string"`                     // PBR
	StacMonth            string          `json:"stac_month"`                     // 결산월
	VolTnrt              float64         `json:"vol_tnrt,string"`                // 거래량 회전율
	Eps                  decimal.Decimal `json:"eps"`                            // EPS
	Bps                  decimal.Decimal `json:"bps"`                            // BPS
	D250Hgpr             decimal.Decimal `json:"d250_hgpr"`                      // 250일 최고가
	D250HgprDate         string          `json:"d250_hgpr_date"`                 // 250일 최고가 일자 (YYYYMMDD)
	D250HgprVrssPrprRate float64         `json:"d250_hgpr_vrss_prpr_rate,string"` // 250일 최고가 대비 현재가 비율
	D250Lwpr             decimal.Decimal `json:"d250_lwpr"`                      // 250일 최저가
	D250LwprDate         string          `json:"d250_lwpr_date"`                 // 250일 최저가 일자
	D250LwprVrssPrprRate float64         `json:"d250_lwpr_vrss_prpr_rate,string"` // 250일 최저가 대비 현재가 비율
	StckDryyHgpr         decimal.Decimal `json:"stck_dryy_hgpr"`                 // 연중 최고가
	DryyHgprVrssPrprRate float64         `json:"dryy_hgpr_vrss_prpr_rate,string"` // 연중 최고가 대비 현재가 비율
	DryyHgprDate         string          `json:"dryy_hgpr_date"`                 // 연중 최고가 일자
	StckDryyLwpr         decimal.Decimal `json:"stck_dryy_lwpr"`                 // 연중 최저가
	DryyLwprVrssPrprRate float64         `json:"dryy_lwpr_vrss_prpr_rate,string"` // 연중 최저가 대비 현재가 비율
	DryyLwprDate         string          `json:"dryy_lwpr_date"`                 // 연중 최저가 일자
	W52Hgpr              decimal.Decimal `json:"w52_hgpr"`                       // 52주 최고가
	W52HgprVrssPrprCtrt  float64         `json:"w52_hgpr_vrss_prpr_ctrt,string"` // 52주 최고가 대비
	W52HgprDate          string          `json:"w52_hgpr_date"`                  // 52주 최고가 일자
	W52Lwpr              decimal.Decimal `json:"w52_lwpr"`                       // 52주 최저가
	W52LwprVrssPrprCtrt  float64         `json:"w52_lwpr_vrss_prpr_ctrt,string"` // 52주 최저가 대비
	W52LwprDate          string          `json:"w52_lwpr_date"`                  // 52주 최저가 일자
	WholLoanRmndRate     float64         `json:"whol_loan_rmnd_rate,string"`     // 전체 융자 잔고 비율
	SstsYn               string          `json:"ssts_yn"`                        // 공매도 가능 (Y/N)
	StckShrnIscd         string          `json:"stck_shrn_iscd"`                 // 단축 종목코드
	FcamCnnm             string          `json:"fcam_cnnm"`                      // 액면가 통화명
	CpfnCnnm             string          `json:"cpfn_cnnm"`                      // 자본금 통화명
	ApprchRate           float64         `json:"apprch_rate,string"`             // 접근도
	FrgnHldnQty          int64           `json:"frgn_hldn_qty,string"`           // 외국인 보유 수량
	ViClsCode            string          `json:"vi_cls_code"`                    // VI 적용 구분
	OvtmViClsCode        string          `json:"ovtm_vi_cls_code"`               // 시간외단일가 VI 구분
	LastSstsCntgQty      int64           `json:"last_ssts_cntg_qty,string"`      // 최종 공매도 체결 수량
	InvtCafulYn          string          `json:"invt_caful_yn"`                  // 투자유의 (Y/N)
	MrktWarnClsCode      string          `json:"mrkt_warn_cls_code"`             // 시장경고 코드
	ShortOverYn          string          `json:"short_over_yn"`                  // 단기과열 (Y/N)
	SltrYn               string          `json:"sltr_yn"`                        // 정리매매 (Y/N)
	MangIssuClsCode      string          `json:"mang_issu_cls_code"`             // 관리종목 (Y/N)
}

// InquirePrice 는 주식현재가 시세 호출.
//
// 한투 docs: docs/api/국내주식/주식현재가_시세.md
// path: /uapi/domestic-stock/v1/quotations/inquire-price (FHKST01010100)
//
// FID_COND_MRKT_DIV_CODE 는 "J" (KRX) 고정. NXT/통합 시장은 별도 메서드 후일 추가 검토.
func (c *Client) InquirePrice(ctx context.Context, symbol string) (*Price, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-price",
		TrID:   "FHKST01010100",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": "J",
			"FID_INPUT_ISCD":         symbol,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var p Price
	if err := json.Unmarshal(resp.Output, &p); err != nil {
		return nil, fmt.Errorf("kis: parse Price: %w", err)
	}
	return &p, nil
}
```

- [ ] **Step 5: 테스트 실행 → PASS**

Run: `go test ./domestic/... -run InquirePrice -v`
Expected: 2 test cases PASS.

- [ ] **Step 6: Commit**

```bash
git add domestic/price.go domestic/price_test.go domestic/testdata/price_success.json
git commit -m "$(cat <<'EOF'
[feat] domestic — InquirePrice (주식현재가 시세, FHKST01010100)

Price struct 의 80+ 필드 한투 docs (docs/api/국내주식/주식현재가_시세.md) 와
1:1 매핑, 한투 약어를 PascalCase 로 변환 + 인라인 한글 코멘트. FID_COND_MRKT_DIV_CODE
는 "J" (KRX) 고정.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 9: domestic/info.go — SearchInfo + ProductInfo struct

**Files:**
- Create: `domestic/info.go`
- Create: `domestic/info_test.go`
- Create: `domestic/testdata/product_info_success.json`

- [ ] **Step 1: testdata fixture** — `domestic/testdata/product_info_success.json`

> 한투 docs (`docs/api/국내주식/상품기본조회.md`) 의 응답 필드 정의 기반.

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": {
    "pdno": "005930",
    "prdt_type_cd": "300",
    "prdt_name": "삼성전자",
    "prdt_name120": "삼성전자보통주",
    "prdt_abrv_name": "삼성전자",
    "prdt_eng_name": "SAMSUNG ELECTRONICS",
    "prdt_eng_name120": "SAMSUNG ELECTRONICS CO.,LTD.",
    "prdt_eng_abrv_name": "SAMSUNG",
    "std_pdno": "KR7005930003",
    "shtn_pdno": "005930",
    "prdt_sale_stat_cd": "01",
    "prdt_risk_grade_cd": "5",
    "prdt_clsf_cd": "STK",
    "prdt_clsf_name": "주권",
    "sale_strt_dt": "19750611",
    "sale_end_dt": "99991231",
    "wrap_asst_type_cd": "00",
    "ivst_prdt_type_cd": "100",
    "ivst_prdt_type_cd_name": "국내주식",
    "frst_erlm_dt": "19750611"
  }
}
```

- [ ] **Step 2: 테스트 작성** — `domestic/info_test.go` (SearchInfo 부분)

```go
package domestic_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_SearchInfo(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/search-info`,
		httpmock.NewStringResponder(200, loadFixtureString(t, "product_info_success.json")),
	)

	c := newTestClient(t)
	info, err := c.SearchInfo(context.Background(), "005930", "300")
	require.NoError(t, err)
	require.NotNil(t, info)

	assert.Equal(t, "005930", info.Pdno)
	assert.Equal(t, "300", info.PrdtTypeCd)
	assert.Equal(t, "삼성전자", info.PrdtName)
	assert.Equal(t, "주권", info.PrdtClsfName)
	assert.Equal(t, "KR7005930003", info.StdPdno)
}
```

- [ ] **Step 3: 테스트 실행 → FAIL**

Run: `go test ./domestic/... -run SearchInfo -v`
Expected: 컴파일 실패 (`SearchInfo`, `ProductInfo` 미정의).

- [ ] **Step 4: 구현** — `domestic/info.go` (SearchInfo + ProductInfo 부분)

```go
package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// ProductInfo 는 상품기본조회 (CTPF1604R) 의 output.
//
// 한투 docs: docs/api/국내주식/상품기본조회.md
// path: /uapi/domestic-stock/v1/quotations/search-info
//
// 다국가 종목 (KR/US/JP/HK/CN/VN) 의 가벼운 기본 정보. PDNO + PRDT_TYPE_CD 입력.
type ProductInfo struct {
	Pdno               string `json:"pdno"`                  // 상품번호
	PrdtTypeCd         string `json:"prdt_type_cd"`          // 상품유형 (300=주식, 512=US 나스닥, ...)
	PrdtName           string `json:"prdt_name"`             // 상품명
	PrdtName120        string `json:"prdt_name120"`          // 상품명120
	PrdtAbrvName       string `json:"prdt_abrv_name"`        // 상품 약어명
	PrdtEngName        string `json:"prdt_eng_name"`         // 상품 영문명
	PrdtEngName120     string `json:"prdt_eng_name120"`      // 상품 영문명120
	PrdtEngAbrvName    string `json:"prdt_eng_abrv_name"`    // 상품 영문 약어명
	StdPdno            string `json:"std_pdno"`              // 표준상품번호 (ISIN)
	ShtnPdno           string `json:"shtn_pdno"`             // 단축 상품번호
	PrdtSaleStatCd     string `json:"prdt_sale_stat_cd"`     // 판매상태
	PrdtRiskGradeCd    string `json:"prdt_risk_grade_cd"`    // 위험등급
	PrdtClsfCd         string `json:"prdt_clsf_cd"`          // 상품분류 (STK 등)
	PrdtClsfName       string `json:"prdt_clsf_name"`        // 상품분류명 (주권/ETF/REITs/...)
	SaleStrtDt         string `json:"sale_strt_dt"`          // 판매 시작일
	SaleEndDt          string `json:"sale_end_dt"`           // 판매 종료일
	WrapAsstTypeCd     string `json:"wrap_asst_type_cd"`     // wrap 자산유형
	IvstPrdtTypeCd     string `json:"ivst_prdt_type_cd"`     // 투자상품 유형
	IvstPrdtTypeCdName string `json:"ivst_prdt_type_cd_name"` // 투자상품 유형명
	FrstErlmDt         string `json:"frst_erlm_dt"`          // 최초 등록일
}

// SearchInfo 는 상품기본조회 호출.
//
// 한투 docs: docs/api/국내주식/상품기본조회.md
// path: /uapi/domestic-stock/v1/quotations/search-info (CTPF1604R)
//
// PRDT_TYPE_CD 예: "300"(국내 주식), "512"(US 나스닥), "513"(US 뉴욕), "529"(US 아멕스).
// 호출자가 직접 명시 — Python 의 country_code fallback 루프는 한투 spec 에 없음.
func (c *Client) SearchInfo(ctx context.Context, pdno, prdtTypeCD string) (*ProductInfo, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/search-info",
		TrID:   "CTPF1604R",
		Query: map[string]string{
			"PDNO":         pdno,
			"PRDT_TYPE_CD": prdtTypeCD,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var info ProductInfo
	if err := json.Unmarshal(resp.Output, &info); err != nil {
		return nil, fmt.Errorf("kis: parse ProductInfo: %w", err)
	}
	return &info, nil
}
```

- [ ] **Step 5: 테스트 실행 → PASS**

Run: `go test ./domestic/... -run SearchInfo -v`
Expected: 1 test case PASS.

- [ ] **Step 6: Commit**

```bash
git add domestic/info.go domestic/info_test.go domestic/testdata/product_info_success.json
git commit -m "$(cat <<'EOF'
[feat] domestic — SearchInfo (상품기본조회, CTPF1604R)

ProductInfo struct + SearchInfo. PDNO + PRDT_TYPE_CD 인자, Python 의
country_code fallback 루프 없음 (한투 spec 충실).

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 10: domestic/info.go — SearchStockInfo + StockInfo struct

**Files:**
- Modify: `domestic/info.go`
- Modify: `domestic/info_test.go`
- Create: `domestic/testdata/stock_info_success.json`

- [ ] **Step 1: testdata fixture** — `domestic/testdata/stock_info_success.json`

> 한투 docs (`docs/api/국내주식/주식기본조회.md`) 의 응답 필드 정의 기반.

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output": {
    "pdno": "005930",
    "prdt_type_cd": "300",
    "mket_id_cd": "STK",
    "scrt_grp_id_cd": "ST",
    "excg_dvsn_cd": "02",
    "setl_mmdd": "1231",
    "lstg_stqt": "5969782550",
    "lstg_cptl_amt": "778038308000",
    "cpta": "778038308000",
    "papr": "100",
    "issu_pric": "1000",
    "kospi200_item_yn": "Y",
    "scts_mket_lstg_dt": "19750611",
    "scts_mket_lstg_abol_dt": "",
    "kosdaq_mket_lstg_dt": "",
    "kosdaq_mket_lstg_abol_dt": "",
    "frbd_mket_lstg_dt": "",
    "frbd_mket_lstg_abol_dt": "",
    "reits_kind_cd": "",
    "etf_dvsn_cd": "",
    "oilf_fund_yn": "N",
    "idx_bztp_lcls_cd": "001",
    "idx_bztp_mcls_cd": "013",
    "idx_bztp_scls_cd": "01301",
    "idx_bztp_lcls_cd_name": "전기.전자",
    "idx_bztp_mcls_cd_name": "전기.전자",
    "idx_bztp_scls_cd_name": "전기.전자",
    "stck_kind_cd": "101",
    "mfnd_opng_dt": "",
    "mfnd_end_dt": "",
    "dpsi_erlm_cncl_dt": "",
    "etf_cu_qty": "",
    "prdt_name": "삼성전자",
    "prdt_name120": "삼성전자보통주",
    "prdt_abrv_name": "삼성전자",
    "std_pdno": "KR7005930003",
    "prdt_eng_name": "SAMSUNG ELECTRONICS",
    "prdt_eng_name120": "SAMSUNG ELECTRONICS CO.,LTD.",
    "prdt_eng_abrv_name": "SAMSUNG",
    "dpsi_aptm_erlm_yn": "N",
    "etf_txtn_type_cd": "",
    "etf_type_cd": "",
    "lstg_abol_dt": "",
    "nwst_odst_dvsn_cd": "",
    "sbst_pric": "60640",
    "thco_sbst_pric": "60640",
    "thco_sbst_pric_chng_dt": "20260101",
    "tr_stop_yn": "N",
    "admn_item_yn": "N",
    "thdt_clpr": "75800",
    "bfdy_clpr": "76000",
    "clpr_chng_dt": "20260502",
    "std_idst_clsf_cd": "C26",
    "std_idst_clsf_cd_name": "전자부품 제조업",
    "idx_bztp_lcls_cd_eng_name": "Electronics",
    "idx_bztp_mcls_cd_eng_name": "Electronics",
    "idx_bztp_scls_cd_eng_name": "Electronics"
  }
}
```

- [ ] **Step 2: 테스트 추가** — `domestic/info_test.go` 에 함수 추가

```go
func TestClient_SearchStockInfo(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/search-stock-info`,
		httpmock.NewStringResponder(200, loadFixtureString(t, "stock_info_success.json")),
	)

	c := newTestClient(t)
	info, err := c.SearchStockInfo(context.Background(), "005930", "300")
	require.NoError(t, err)
	require.NotNil(t, info)

	assert.Equal(t, "005930", info.Pdno)
	assert.Equal(t, "300", info.PrdtTypeCd)
	assert.Equal(t, "STK", info.MketIdCd)
	assert.Equal(t, "ST", info.ScrtGrpIdCd)
	assert.Equal(t, "Y", info.Kospi200ItemYn)
	assert.Equal(t, "삼성전자", info.PrdtName)
	assert.Equal(t, "1231", info.SetlMmdd)
	assert.Equal(t, "N", info.AdmnItemYn)
	assert.Equal(t, "N", info.TrStopYn)
}
```

- [ ] **Step 3: 테스트 실행 → FAIL**

Run: `go test ./domestic/... -run SearchStockInfo -v`
Expected: 컴파일 실패 (`SearchStockInfo`, `StockInfo` 미정의).

- [ ] **Step 4: 구현 추가** — `domestic/info.go` 에 추가

```go
// StockInfo 는 주식기본조회 (CTPF1002R) 의 output.
//
// 한투 docs: docs/api/국내주식/주식기본조회.md
// path: /uapi/domestic-stock/v1/quotations/search-stock-info
//
// 국내주식 종목의 상세 정보 (시장ID, 증권그룹, 상장정보, 업종분류, KOSPI200 편입 등).
// 한투 spec 자체는 다국가 PRDT_TYPE_CD 받지만 endpoint 는 국내 전용 데이터.
type StockInfo struct {
	Pdno                  string `json:"pdno"`                       // 상품번호
	PrdtTypeCd            string `json:"prdt_type_cd"`               // 상품유형코드
	MketIdCd              string `json:"mket_id_cd"`                 // 시장ID 코드 (STK=KOSPI, KSQ=KOSDAQ, ...)
	ScrtGrpIdCd           string `json:"scrt_grp_id_cd"`             // 증권그룹 ID (ST=주권, EF=ETF, ...)
	ExcgDvsnCd            string `json:"excg_dvsn_cd"`               // 거래소 구분 코드
	SetlMmdd              string `json:"setl_mmdd"`                  // 결산 월일 (MMDD)
	LstgStqt              int64  `json:"lstg_stqt,string"`           // 상장주수
	LstgCptlAmt           int64  `json:"lstg_cptl_amt,string"`       // 상장 자본금
	Cpta                  int64  `json:"cpta,string"`                // 자본금
	Papr                  string `json:"papr"`                       // 액면가
	IssuPric              string `json:"issu_pric"`                  // 발행가
	Kospi200ItemYn        string `json:"kospi200_item_yn"`           // KOSPI200 편입 (Y/N)
	SctsMketLstgDt        string `json:"scts_mket_lstg_dt"`          // 유가증권 시장 상장일
	SctsMketLstgAbolDt    string `json:"scts_mket_lstg_abol_dt"`     // 유가증권 시장 상장폐지일
	KosdaqMketLstgDt      string `json:"kosdaq_mket_lstg_dt"`        // KOSDAQ 시장 상장일
	KosdaqMketLstgAbolDt  string `json:"kosdaq_mket_lstg_abol_dt"`   // KOSDAQ 시장 상장폐지일
	FrbdMketLstgDt        string `json:"frbd_mket_lstg_dt"`          // 외국 시장 상장일
	FrbdMketLstgAbolDt    string `json:"frbd_mket_lstg_abol_dt"`     // 외국 시장 상장폐지일
	ReitsKindCd           string `json:"reits_kind_cd"`              // 리츠 종류
	EtfDvsnCd             string `json:"etf_dvsn_cd"`                // ETF 구분
	OilfFundYn            string `json:"oilf_fund_yn"`               // 유전 펀드 (Y/N)
	IdxBztpLclsCd         string `json:"idx_bztp_lcls_cd"`           // 업종 대분류 코드
	IdxBztpMclsCd         string `json:"idx_bztp_mcls_cd"`           // 업종 중분류 코드
	IdxBztpSclsCd         string `json:"idx_bztp_scls_cd"`           // 업종 소분류 코드
	IdxBztpLclsCdName     string `json:"idx_bztp_lcls_cd_name"`      // 업종 대분류명
	IdxBztpMclsCdName     string `json:"idx_bztp_mcls_cd_name"`      // 업종 중분류명
	IdxBztpSclsCdName     string `json:"idx_bztp_scls_cd_name"`      // 업종 소분류명
	StckKindCd            string `json:"stck_kind_cd"`               // 주식 종류 코드
	MfndOpngDt            string `json:"mfnd_opng_dt"`               // 펀드 개시일
	MfndEndDt             string `json:"mfnd_end_dt"`                // 펀드 종료일
	DpsiErlmCnclDt        string `json:"dpsi_erlm_cncl_dt"`          // 예수금 등록취소일
	EtfCuQty              string `json:"etf_cu_qty"`                 // ETF CU 수량
	PrdtName              string `json:"prdt_name"`                  // 상품명
	PrdtName120           string `json:"prdt_name120"`               // 상품명120
	PrdtAbrvName          string `json:"prdt_abrv_name"`             // 상품 약어명
	StdPdno               string `json:"std_pdno"`                   // 표준상품번호
	PrdtEngName           string `json:"prdt_eng_name"`              // 영문명
	PrdtEngName120        string `json:"prdt_eng_name120"`           // 영문명120
	PrdtEngAbrvName       string `json:"prdt_eng_abrv_name"`         // 영문 약어명
	DpsiAptmErlmYn        string `json:"dpsi_aptm_erlm_yn"`          // 예수금 적용 등록 (Y/N)
	EtfTxtnTypeCd         string `json:"etf_txtn_type_cd"`           // ETF 과세 유형
	EtfTypeCd             string `json:"etf_type_cd"`                // ETF 유형
	LstgAbolDt            string `json:"lstg_abol_dt"`               // 상장 폐지일
	NwstOdstDvsnCd        string `json:"nwst_odst_dvsn_cd"`          // 신규/구주 구분
	SbstPric              string `json:"sbst_pric"`                  // 대용가
	ThcoSbstPric          string `json:"thco_sbst_pric"`             // 당사 대용가
	ThcoSbstPricChngDt    string `json:"thco_sbst_pric_chng_dt"`     // 당사 대용가 변경일
	TrStopYn              string `json:"tr_stop_yn"`                 // 거래 정지 (Y/N)
	AdmnItemYn            string `json:"admn_item_yn"`               // 관리종목 (Y/N)
	ThdtClpr              string `json:"thdt_clpr"`                  // 당일 종가
	BfdyClpr              string `json:"bfdy_clpr"`                  // 전일 종가
	ClprChngDt            string `json:"clpr_chng_dt"`               // 종가 변경일
	StdIdstClsfCd         string `json:"std_idst_clsf_cd"`           // 표준 산업분류 코드
	StdIdstClsfCdName     string `json:"std_idst_clsf_cd_name"`      // 표준 산업분류명
	IdxBztpLclsCdEngName  string `json:"idx_bztp_lcls_cd_eng_name"`  // 업종 대분류 영문명
	IdxBztpMclsCdEngName  string `json:"idx_bztp_mcls_cd_eng_name"`  // 업종 중분류 영문명
	IdxBztpSclsCdEngName  string `json:"idx_bztp_scls_cd_eng_name"`  // 업종 소분류 영문명
}

// SearchStockInfo 는 주식기본조회 호출.
//
// 한투 docs: docs/api/국내주식/주식기본조회.md
// path: /uapi/domestic-stock/v1/quotations/search-stock-info (CTPF1002R)
//
// 한투 spec 충실: PDNO + PRDT_TYPE_CD 명시. Python 의 "KR" country_code 검사 없음.
func (c *Client) SearchStockInfo(ctx context.Context, pdno, prdtTypeCD string) (*StockInfo, error) {
	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/search-stock-info",
		TrID:   "CTPF1002R",
		Query: map[string]string{
			"PDNO":         pdno,
			"PRDT_TYPE_CD": prdtTypeCD,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}
	var info StockInfo
	if err := json.Unmarshal(resp.Output, &info); err != nil {
		return nil, fmt.Errorf("kis: parse StockInfo: %w", err)
	}
	return &info, nil
}
```

- [ ] **Step 5: 테스트 실행 → PASS**

Run: `go test ./domestic/... -run 'Search' -v`
Expected: SearchInfo + SearchStockInfo 모두 PASS.

- [ ] **Step 6: Commit**

```bash
git add domestic/info.go domestic/info_test.go domestic/testdata/stock_info_success.json
git commit -m "$(cat <<'EOF'
[feat] domestic — SearchStockInfo (주식기본조회, CTPF1002R)

StockInfo struct + SearchStockInfo. 시장ID/증권그룹/업종분류/KOSPI200 편입 등
국내주식 디테일. 한투 spec 충실: PDNO + PRDT_TYPE_CD 명시 받음, Python 의
country_code "KR" 검사 없음.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 11: domestic/chart.go — InquireDailyItemChartPrice + DailyChart

**Files:**
- Create: `domestic/chart.go`
- Create: `domestic/chart_test.go`
- Create: `domestic/testdata/daily_chart_success.json`

- [ ] **Step 1: testdata fixture** — `domestic/testdata/daily_chart_success.json`

> 한투 docs (`docs/api/국내주식/국내주식기간별시세(일_주_월_년).md`). output1 (요약) + output2 (3 캔들).

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": {
    "prdy_vrss": "-200",
    "prdy_vrss_sign": "5",
    "prdy_ctrt": "-0.26",
    "stck_prdy_clpr": "76000",
    "acml_vol": "12345678",
    "acml_tr_pbmn": "938223456000",
    "hts_kor_isnm": "삼성전자",
    "stck_prpr": "75800",
    "stck_shrn_iscd": "005930",
    "prdy_vol": "11000000",
    "stck_mxpr": "98500",
    "stck_llam": "53100",
    "stck_oprc": "76000",
    "stck_hgpr": "76200",
    "stck_lwpr": "75500",
    "stck_prdy_oprc": "76500",
    "stck_prdy_hgpr": "76800",
    "stck_prdy_lwpr": "75900",
    "askp": "75900",
    "bidp": "75800",
    "prdy_vrss_vol": "1345678",
    "vol_tnrt": "0.21",
    "stck_fcam": "100",
    "lstn_stcn": "5969782550",
    "cpfn": "778038308",
    "hts_avls": "452312",
    "per": "11.42",
    "eps": "6638",
    "pbr": "1.32",
    "itewhol_loan_rmnd_ratem name": "0.45"
  },
  "output2": [
    {
      "stck_bsop_date": "20260502",
      "stck_clpr": "75800",
      "stck_oprc": "76000",
      "stck_hgpr": "76200",
      "stck_lwpr": "75500",
      "acml_vol": "12345678",
      "acml_tr_pbmn": "938223456000",
      "flng_cls_code": "00",
      "prtt_rate": "0.00",
      "mod_yn": "N",
      "prdy_vrss_sign": "5",
      "prdy_vrss": "-200",
      "revl_issu_reas": ""
    },
    {
      "stck_bsop_date": "20260501",
      "stck_clpr": "76000",
      "stck_oprc": "76500",
      "stck_hgpr": "76800",
      "stck_lwpr": "75900",
      "acml_vol": "11000000",
      "acml_tr_pbmn": "836000000000",
      "flng_cls_code": "00",
      "prtt_rate": "0.00",
      "mod_yn": "N",
      "prdy_vrss_sign": "2",
      "prdy_vrss": "100",
      "revl_issu_reas": ""
    },
    {
      "stck_bsop_date": "20260430",
      "stck_clpr": "75900",
      "stck_oprc": "76200",
      "stck_hgpr": "76400",
      "stck_lwpr": "75600",
      "acml_vol": "10500000",
      "acml_tr_pbmn": "796950000000",
      "flng_cls_code": "00",
      "prtt_rate": "0.00",
      "mod_yn": "N",
      "prdy_vrss_sign": "5",
      "prdy_vrss": "-300",
      "revl_issu_reas": ""
    }
  ]
}
```

- [ ] **Step 2: 테스트 작성** — `domestic/chart_test.go` (Daily 부분)

```go
package domestic_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/domestic"
)

func TestClient_InquireDailyItemChartPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery map[string][]string
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-daily-itemchartprice`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "daily_chart_success.json")), nil
		},
	)

	c := newTestClient(t)
	chart, err := c.InquireDailyItemChartPrice(context.Background(), domestic.InquireDailyItemChartPriceParams{
		Symbol:   "005930",
		FromDate: "20260430",
		ToDate:   "20260502",
	})
	require.NoError(t, err)
	require.NotNil(t, chart)

	// query default 검증 (zero-value → "D" / "J" / 수정주가)
	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "005930", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "20260430", capturedQuery.Get("FID_INPUT_DATE_1"))
	assert.Equal(t, "20260502", capturedQuery.Get("FID_INPUT_DATE_2"))
	assert.Equal(t, "D", capturedQuery.Get("FID_PERIOD_DIV_CODE"))
	assert.Equal(t, "0", capturedQuery.Get("FID_ORG_ADJ_PRC")) // 0 = 수정주가

	// output1 검증
	assert.Equal(t, "삼성전자", chart.Output1.HtsKorIsnm)
	assert.Equal(t, decimal.NewFromInt(75800), chart.Output1.StckPrpr)
	assert.Equal(t, "005930", chart.Output1.StckShrnIscd)

	// output2 검증
	require.Len(t, chart.Output2, 3)
	assert.Equal(t, "20260502", chart.Output2[0].StckBsopDate)
	assert.Equal(t, decimal.NewFromInt(75800), chart.Output2[0].StckClpr)
	assert.Equal(t, decimal.NewFromInt(76000), chart.Output2[0].StckOprc)
	assert.Equal(t, int64(12345678), chart.Output2[0].AcmlVol)
}

func TestClient_InquireDailyItemChartPrice_OriginalPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery map[string][]string
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-daily-itemchartprice`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "daily_chart_success.json")), nil
		},
	)

	c := newTestClient(t)
	_, err := c.InquireDailyItemChartPrice(context.Background(), domestic.InquireDailyItemChartPriceParams{
		Symbol:        "005930",
		Period:        "W",
		FromDate:      "20260101",
		ToDate:        "20260502",
		OriginalPrice: true,
		MarketCode:    "NX",
	})
	require.NoError(t, err)
	assert.Equal(t, "NX", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "W", capturedQuery.Get("FID_PERIOD_DIV_CODE"))
	assert.Equal(t, "1", capturedQuery.Get("FID_ORG_ADJ_PRC")) // 1 = 원주가
}
```

- [ ] **Step 3: 테스트 실행 → FAIL**

Run: `go test ./domestic/... -run InquireDaily -v`
Expected: 컴파일 실패 (`InquireDailyItemChartPrice`, `DailyChart`, `InquireDailyItemChartPriceParams` 미정의).

- [ ] **Step 4: 구현** — `domestic/chart.go`

```go
package domestic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kenshin579/korea-investment-stock/internal/httpclient"
)

// DailyChart 는 국내주식기간별시세(일/주/월/년) (FHKST03010100) 응답.
//
// 한투 docs: docs/api/국내주식/국내주식기간별시세(일_주_월_년).md
type DailyChart struct {
	Output1 DailyChartSummary `json:"output1"` // 요약 (현재가, 전일대비 등)
	Output2 []DailyChartCandle `json:"output2"` // 일/주/월/년 봉 배열 (최대 100건)
}

// DailyChartSummary 는 차트 응답의 output1 (단일 객체, 요약 정보).
type DailyChartSummary struct {
	PrdyVrss               decimal.Decimal `json:"prdy_vrss"`                 // 전일 대비
	PrdyVrssSign           string          `json:"prdy_vrss_sign"`            // 전일 대비 부호
	PrdyCtrt               float64         `json:"prdy_ctrt,string"`          // 전일 대비율
	StckPrdyClpr           decimal.Decimal `json:"stck_prdy_clpr"`            // 전일 종가
	AcmlVol                int64           `json:"acml_vol,string"`           // 누적 거래량
	AcmlTrPbmn             int64           `json:"acml_tr_pbmn,string"`       // 누적 거래 대금
	HtsKorIsnm             string          `json:"hts_kor_isnm"`              // 종목 한글명
	StckPrpr               decimal.Decimal `json:"stck_prpr"`                 // 현재가
	StckShrnIscd           string          `json:"stck_shrn_iscd"`            // 단축 종목코드
	PrdyVol                int64           `json:"prdy_vol,string"`           // 전일 거래량
	StckMxpr               decimal.Decimal `json:"stck_mxpr"`                 // 상한가
	StckLlam               decimal.Decimal `json:"stck_llam"`                 // 하한가
	StckOprc               decimal.Decimal `json:"stck_oprc"`                 // 시가
	StckHgpr               decimal.Decimal `json:"stck_hgpr"`                 // 고가
	StckLwpr               decimal.Decimal `json:"stck_lwpr"`                 // 저가
	StckPrdyOprc           decimal.Decimal `json:"stck_prdy_oprc"`            // 전일 시가
	StckPrdyHgpr           decimal.Decimal `json:"stck_prdy_hgpr"`            // 전일 고가
	StckPrdyLwpr           decimal.Decimal `json:"stck_prdy_lwpr"`            // 전일 저가
	Askp                   decimal.Decimal `json:"askp"`                      // 매도호가
	Bidp                   decimal.Decimal `json:"bidp"`                      // 매수호가
	PrdyVrssVol            int64           `json:"prdy_vrss_vol,string"`      // 전일 대비 거래량
	VolTnrt                float64         `json:"vol_tnrt,string"`           // 거래량 회전율
	StckFcam               decimal.Decimal `json:"stck_fcam"`                 // 액면가
	LstnStcn               int64           `json:"lstn_stcn,string"`          // 상장주수
	Cpfn                   int64           `json:"cpfn,string"`               // 자본금 (백만원)
	HtsAvls                int64           `json:"hts_avls,string"`           // HTS 시가총액
	Per                    float64         `json:"per,string"`                // PER
	Eps                    decimal.Decimal `json:"eps"`                       // EPS
	Pbr                    float64         `json:"pbr,string"`                // PBR
}

// DailyChartCandle 은 차트 응답의 output2 한 행 (한 캔들).
type DailyChartCandle struct {
	StckBsopDate string          `json:"stck_bsop_date"`        // 영업일자 (YYYYMMDD)
	StckClpr     decimal.Decimal `json:"stck_clpr"`             // 종가
	StckOprc     decimal.Decimal `json:"stck_oprc"`             // 시가
	StckHgpr     decimal.Decimal `json:"stck_hgpr"`             // 고가
	StckLwpr     decimal.Decimal `json:"stck_lwpr"`             // 저가
	AcmlVol      int64           `json:"acml_vol,string"`       // 거래량
	AcmlTrPbmn   int64           `json:"acml_tr_pbmn,string"`   // 거래 대금
	FlngClsCode  string          `json:"flng_cls_code"`         // 락 구분 (00=일반)
	PrttRate     float64         `json:"prtt_rate,string"`      // 분할 비율
	ModYn        string          `json:"mod_yn"`                // 분할 변경 (Y/N)
	PrdyVrssSign string          `json:"prdy_vrss_sign"`        // 전일 대비 부호
	PrdyVrss     decimal.Decimal `json:"prdy_vrss"`             // 전일 대비
	RevlIssuReas string          `json:"revl_issu_reas"`        // 재평가 사유
}

// InquireDailyItemChartPriceParams 는 일/주/월/년 봉 조회 파라미터.
type InquireDailyItemChartPriceParams struct {
	Symbol        string // 필수, 종목코드 (예 "005930")
	Period        string // "D"/"W"/"M"/"Y", 빈 값이면 "D"
	FromDate      string // YYYYMMDD, 필수
	ToDate        string // YYYYMMDD, 필수, 1회 최대 100건
	OriginalPrice bool   // false=수정주가(default), true=원주가
	MarketCode    string // "J"/"NX"/"UN", 빈 값이면 "J"
}

// InquireDailyItemChartPrice 는 국내주식기간별시세(일/주/월/년) 호출.
//
// 한투 docs: docs/api/국내주식/국내주식기간별시세(일_주_월_년).md
// path: /uapi/domestic-stock/v1/quotations/inquire-daily-itemchartprice (FHKST03010100)
func (c *Client) InquireDailyItemChartPrice(ctx context.Context, params InquireDailyItemChartPriceParams) (*DailyChart, error) {
	period := params.Period
	if period == "" {
		period = "D"
	}
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	adjPrc := "0" // 수정주가
	if params.OriginalPrice {
		adjPrc = "1"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-daily-itemchartprice",
		TrID:   "FHKST03010100",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_INPUT_DATE_1":       params.FromDate,
			"FID_INPUT_DATE_2":       params.ToDate,
			"FID_PERIOD_DIV_CODE":    period,
			"FID_ORG_ADJ_PRC":        adjPrc,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	// output1 + output2 동시 unmarshal: resp.Raw 에서 파싱
	var chart DailyChart
	if err := json.Unmarshal(resp.Raw, &chart); err != nil {
		return nil, fmt.Errorf("kis: parse DailyChart: %w", err)
	}
	return &chart, nil
}
```

- [ ] **Step 5: 테스트 실행 → PASS**

Run: `go test ./domestic/... -run InquireDaily -v`
Expected: 2 test cases PASS.

- [ ] **Step 6: Commit**

```bash
git add domestic/chart.go domestic/chart_test.go domestic/testdata/daily_chart_success.json
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireDailyItemChartPrice (국내주식기간별시세, FHKST03010100)

DailyChart (Output1 요약 + Output2 캔들 배열) + InquireDailyItemChartPriceParams.
zero-value default: Period="D", MarketCode="J", OriginalPrice=false(수정주가).
output1+output2 동시 파싱은 resp.Raw 에서 직접 unmarshal.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 12: domestic/chart.go — InquireTimeItemChartPrice + MinuteChart

**Files:**
- Modify: `domestic/chart.go`
- Modify: `domestic/chart_test.go`
- Create: `domestic/testdata/minute_chart_success.json`

- [ ] **Step 1: testdata fixture** — `domestic/testdata/minute_chart_success.json`

```json
{
  "rt_cd": "0",
  "msg_cd": "MCA00000",
  "msg1": "정상처리 되었습니다.",
  "output1": {
    "prdy_vrss": "-200",
    "prdy_vrss_sign": "5",
    "prdy_ctrt": "-0.26",
    "stck_prdy_clpr": "76000",
    "acml_vol": "12345678",
    "acml_tr_pbmn": "938223456000",
    "hts_kor_isnm": "삼성전자",
    "stck_prpr": "75800"
  },
  "output2": [
    {
      "stck_bsop_date": "20260503",
      "stck_cntg_hour": "150000",
      "stck_prpr": "75800",
      "stck_oprc": "75900",
      "stck_hgpr": "75900",
      "stck_lwpr": "75700",
      "cntg_vol": "12345",
      "acml_tr_pbmn": "936000000"
    },
    {
      "stck_bsop_date": "20260503",
      "stck_cntg_hour": "145900",
      "stck_prpr": "75900",
      "stck_oprc": "75850",
      "stck_hgpr": "75950",
      "stck_lwpr": "75800",
      "cntg_vol": "8765",
      "acml_tr_pbmn": "664500000"
    },
    {
      "stck_bsop_date": "20260503",
      "stck_cntg_hour": "145800",
      "stck_prpr": "75850",
      "stck_oprc": "75800",
      "stck_hgpr": "75900",
      "stck_lwpr": "75750",
      "cntg_vol": "10234",
      "acml_tr_pbmn": "775944000"
    }
  ]
}
```

- [ ] **Step 2: 테스트 추가** — `domestic/chart_test.go` 에 함수 추가

```go
func TestClient_InquireTimeItemChartPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery map[string][]string
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-time-itemchartprice`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "minute_chart_success.json")), nil
		},
	)

	c := newTestClient(t)
	chart, err := c.InquireTimeItemChartPrice(context.Background(), domestic.InquireTimeItemChartPriceParams{
		Symbol:   "005930",
		TimeFrom: "150000",
	})
	require.NoError(t, err)
	require.NotNil(t, chart)

	assert.Equal(t, "J", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
	assert.Equal(t, "005930", capturedQuery.Get("FID_INPUT_ISCD"))
	assert.Equal(t, "150000", capturedQuery.Get("FID_INPUT_HOUR_1"))
	assert.Equal(t, "N", capturedQuery.Get("FID_PW_DATA_INCU_YN"))

	require.Len(t, chart.Output2, 3)
	assert.Equal(t, "150000", chart.Output2[0].StckCntgHour)
	assert.Equal(t, decimal.NewFromInt(75800), chart.Output2[0].StckPrpr)
	assert.Equal(t, int64(12345), chart.Output2[0].CntgVol)
}

func TestClient_InquireTimeItemChartPrice_PastDataInclude(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var capturedQuery map[string][]string
	httpmock.RegisterResponder(
		http.MethodGet,
		`=~/quotations/inquire-time-itemchartprice`,
		func(req *http.Request) (*http.Response, error) {
			capturedQuery = req.URL.Query()
			return httpmock.NewStringResponse(200, loadFixtureString(t, "minute_chart_success.json")), nil
		},
	)

	c := newTestClient(t)
	_, err := c.InquireTimeItemChartPrice(context.Background(), domestic.InquireTimeItemChartPriceParams{
		Symbol:          "005930",
		TimeFrom:        "150000",
		PastDataInclude: true,
		MarketCode:      "UN",
	})
	require.NoError(t, err)
	assert.Equal(t, "Y", capturedQuery.Get("FID_PW_DATA_INCU_YN"))
	assert.Equal(t, "UN", capturedQuery.Get("FID_COND_MRKT_DIV_CODE"))
}
```

- [ ] **Step 3: 테스트 실행 → FAIL**

Run: `go test ./domestic/... -run InquireTime -v`
Expected: 컴파일 실패 (`InquireTimeItemChartPrice`, `MinuteChart`, `InquireTimeItemChartPriceParams` 미정의).

- [ ] **Step 4: 구현 추가** — `domestic/chart.go` 에 추가

```go
// MinuteChart 는 주식당일분봉조회 (FHKST03010200) 응답.
//
// 한투 docs: docs/api/국내주식/주식당일분봉조회.md
// 실전/모의 1회 최대 30건. 당일 분봉만 제공 (전일자 미제공).
type MinuteChart struct {
	Output1 MinuteChartSummary  `json:"output1"` // 요약
	Output2 []MinuteChartCandle `json:"output2"` // 분봉 배열 (최대 30건)
}

type MinuteChartSummary struct {
	PrdyVrss     decimal.Decimal `json:"prdy_vrss"`           // 전일 대비
	PrdyVrssSign string          `json:"prdy_vrss_sign"`      // 전일 대비 부호
	PrdyCtrt     float64         `json:"prdy_ctrt,string"`    // 전일 대비율
	StckPrdyClpr decimal.Decimal `json:"stck_prdy_clpr"`      // 전일 종가
	AcmlVol      int64           `json:"acml_vol,string"`     // 누적 거래량
	AcmlTrPbmn   int64           `json:"acml_tr_pbmn,string"` // 누적 거래 대금
	HtsKorIsnm   string          `json:"hts_kor_isnm"`        // 종목 한글명
	StckPrpr     decimal.Decimal `json:"stck_prpr"`           // 현재가
}

type MinuteChartCandle struct {
	StckBsopDate string          `json:"stck_bsop_date"`      // 영업일자 (YYYYMMDD)
	StckCntgHour string          `json:"stck_cntg_hour"`      // 체결 시간 (HHMMSS)
	StckPrpr     decimal.Decimal `json:"stck_prpr"`           // 현재가
	StckOprc     decimal.Decimal `json:"stck_oprc"`           // 시가
	StckHgpr     decimal.Decimal `json:"stck_hgpr"`           // 고가
	StckLwpr     decimal.Decimal `json:"stck_lwpr"`           // 저가
	CntgVol      int64           `json:"cntg_vol,string"`     // 체결 거래량
	AcmlTrPbmn   int64           `json:"acml_tr_pbmn,string"` // 누적 거래 대금
}

// InquireTimeItemChartPriceParams 는 분봉 조회 파라미터.
type InquireTimeItemChartPriceParams struct {
	Symbol          string // 필수, 종목코드
	TimeFrom        string // 필수, HHMMSS 시작 시간
	PastDataInclude bool   // false="N"(default), true="Y" — FID_PW_DATA_INCU_YN
	EtcClassCode    string // FID_ETC_CLS_CODE, 빈 값 default
	MarketCode      string // "J"/"NX"/"UN", 빈 값이면 "J"
}

// InquireTimeItemChartPrice 는 주식당일분봉조회 호출.
//
// 한투 docs: docs/api/국내주식/주식당일분봉조회.md
// path: /uapi/domestic-stock/v1/quotations/inquire-time-itemchartprice (FHKST03010200)
//
// ※ 당일 분봉만 제공 (전일자 미제공). FID_INPUT_HOUR_1 에 미래시각 입력 시 현재가로 조회됨.
func (c *Client) InquireTimeItemChartPrice(ctx context.Context, params InquireTimeItemChartPriceParams) (*MinuteChart, error) {
	market := params.MarketCode
	if market == "" {
		market = "J"
	}
	pastIncl := "N"
	if params.PastDataInclude {
		pastIncl = "Y"
	}

	resp, err := c.http.Do(ctx, &httpclient.Request{
		Method: http.MethodGet,
		Path:   "/uapi/domestic-stock/v1/quotations/inquire-time-itemchartprice",
		TrID:   "FHKST03010200",
		Query: map[string]string{
			"FID_COND_MRKT_DIV_CODE": market,
			"FID_INPUT_ISCD":         params.Symbol,
			"FID_INPUT_HOUR_1":       params.TimeFrom,
			"FID_PW_DATA_INCU_YN":    pastIncl,
			"FID_ETC_CLS_CODE":       params.EtcClassCode,
		},
		CustType: "P",
	})
	if err != nil {
		return nil, err
	}

	var chart MinuteChart
	if err := json.Unmarshal(resp.Raw, &chart); err != nil {
		return nil, fmt.Errorf("kis: parse MinuteChart: %w", err)
	}
	return &chart, nil
}
```

- [ ] **Step 5: 테스트 실행 → PASS**

Run: `go test ./domestic/... -run InquireTime -v`
Expected: 2 test cases PASS.

- [ ] **Step 6: 전체 회귀 테스트**

Run: `go test ./... -count=1`
Expected: 모든 패키지 PASS.

- [ ] **Step 7: Commit**

```bash
git add domestic/chart.go domestic/chart_test.go domestic/testdata/minute_chart_success.json
git commit -m "$(cat <<'EOF'
[feat] domestic — InquireTimeItemChartPrice (주식당일분봉조회, FHKST03010200)

MinuteChart (Output1 요약 + Output2 분봉 배열, 최대 30건) +
InquireTimeItemChartPriceParams. zero-value default: MarketCode="J",
PastDataInclude=false → "N".

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 13: examples/domestic_price/main.go

**Files:**
- Create: `examples/domestic_price/main.go`

- [ ] **Step 1: example 작성**

```go
// domestic_price example: InquirePrice + SearchInfo + SearchStockInfo.
//
// Run:
//   export KOREA_INVESTMENT_API_KEY=...
//   export KOREA_INVESTMENT_API_SECRET=...
//   export KOREA_INVESTMENT_ACCOUNT_NO=...
//   go run ./examples/domestic_price
package main

import (
	"context"
	"fmt"
	"log"

	kis "github.com/kenshin579/korea-investment-stock"
)

func main() {
	client, err := kis.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	symbol := "005930" // 삼성전자

	price, err := client.Domestic.InquirePrice(ctx, symbol)
	if err != nil {
		log.Fatalf("InquirePrice: %v", err)
	}
	fmt.Printf("[%s] 현재가 %s, 전일대비 %s (%v%%)\n",
		symbol, price.StckPrpr.String(), price.PrdyVrss.String(), price.PrdyCtrt)
	fmt.Printf("  시가/고가/저가: %s / %s / %s\n",
		price.StckOprc.String(), price.StckHgpr.String(), price.StckLwpr.String())
	fmt.Printf("  거래량: %d, PER: %v, PBR: %v\n", price.AcmlVol, price.Per, price.Pbr)

	info, err := client.Domestic.SearchInfo(ctx, symbol, "300")
	if err != nil {
		log.Fatalf("SearchInfo: %v", err)
	}
	fmt.Printf("  상품명: %s (%s, %s)\n", info.PrdtName, info.PrdtClsfName, info.StdPdno)

	stockInfo, err := client.Domestic.SearchStockInfo(ctx, symbol, "300")
	if err != nil {
		log.Fatalf("SearchStockInfo: %v", err)
	}
	fmt.Printf("  시장: %s (%s), 업종: %s, KOSPI200=%s\n",
		stockInfo.MketIdCd, stockInfo.ScrtGrpIdCd,
		stockInfo.IdxBztpLclsCdName, stockInfo.Kospi200ItemYn)
}
```

- [ ] **Step 2: 컴파일 검증**

Run: `go build ./examples/domestic_price && echo OK`
Expected: `OK`.

- [ ] **Step 3: Commit**

```bash
git add examples/domestic_price
git commit -m "$(cat <<'EOF'
[feat] examples/domestic_price — InquirePrice + Search* 통합 사용 예

삼성전자 (005930) 의 현재가, 상품정보, 주식정보 호출 + 출력. 실제 KIS
credentials 필요 (KOREA_INVESTMENT_* env vars).

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 14: examples/domestic_chart/main.go

**Files:**
- Create: `examples/domestic_chart/main.go`

- [ ] **Step 1: example 작성**

```go
// domestic_chart example: InquireDailyItemChartPrice + InquireTimeItemChartPrice.
//
// Run: KIS credentials env vars 후 go run ./examples/domestic_chart
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	kis "github.com/kenshin579/korea-investment-stock"
	"github.com/kenshin579/korea-investment-stock/domestic"
)

func main() {
	client, err := kis.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	symbol := "005930"

	// 일봉 (최근 약 30일)
	to := time.Now().Format("20060102")
	from := time.Now().AddDate(0, 0, -30).Format("20060102")
	daily, err := client.Domestic.InquireDailyItemChartPrice(ctx, domestic.InquireDailyItemChartPriceParams{
		Symbol:   symbol,
		Period:   "D",
		FromDate: from,
		ToDate:   to,
	})
	if err != nil {
		log.Fatalf("InquireDailyItemChartPrice: %v", err)
	}
	fmt.Printf("[%s %s] 일봉 %d 캔들\n", symbol, daily.Output1.HtsKorIsnm, len(daily.Output2))
	for i, c := range daily.Output2 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s: O=%s H=%s L=%s C=%s V=%d\n",
			c.StckBsopDate, c.StckOprc, c.StckHgpr, c.StckLwpr, c.StckClpr, c.AcmlVol)
	}

	// 분봉 (장 마감 직전부터 30개)
	minute, err := client.Domestic.InquireTimeItemChartPrice(ctx, domestic.InquireTimeItemChartPriceParams{
		Symbol:   symbol,
		TimeFrom: "153000",
	})
	if err != nil {
		log.Fatalf("InquireTimeItemChartPrice: %v", err)
	}
	fmt.Printf("\n[%s] 당일 분봉 %d 개 (15:30 시작):\n", symbol, len(minute.Output2))
	for i, c := range minute.Output2 {
		if i >= 5 {
			fmt.Println("  ... (이하 생략)")
			break
		}
		fmt.Printf("  %s %s: O=%s H=%s L=%s C=%s V=%d\n",
			c.StckBsopDate, c.StckCntgHour,
			c.StckOprc, c.StckHgpr, c.StckLwpr, c.StckPrpr, c.CntgVol)
	}
}
```

- [ ] **Step 2: 컴파일 검증**

Run: `go build ./examples/domestic_chart && echo OK`
Expected: `OK`.

- [ ] **Step 3: Commit**

```bash
git add examples/domestic_chart
git commit -m "$(cat <<'EOF'
[feat] examples/domestic_chart — Daily + Minute 차트 출력

최근 30일 일봉 + 당일 15:30 분봉 30개 출력 예시.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 15: examples/kospi_symbols/main.go

**Files:**
- Create: `examples/kospi_symbols/main.go`

- [ ] **Step 1: example 작성**

```go
// kospi_symbols example: FetchKospiSymbols — KRX KOSPI 마스터 다운로드 + 통계.
//
// Run: KIS credentials env vars 후 go run ./examples/kospi_symbols
//
// 첫 실행 시 ~수 MB ZIP 다운로드 (디스크 캐시, default TTL 7일).
package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	kis "github.com/kenshin579/korea-investment-stock"
)

func main() {
	client, err := kis.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	syms, err := client.Domestic.FetchKospiSymbols(ctx)
	if err != nil {
		log.Fatalf("FetchKospiSymbols: %v", err)
	}
	fmt.Printf("KOSPI 종목 %d 개\n", len(syms))

	// 그룹코드별 분포
	byGroup := make(map[string]int)
	for _, s := range syms {
		byGroup[s.GroupCode]++
	}
	fmt.Println("그룹별 분포:")
	for k, v := range byGroup {
		fmt.Printf("  %s: %d\n", k, v)
	}

	// 우선주만 추리기
	var preferred []string
	for _, s := range syms {
		if s.PreferredStock == "Y" {
			preferred = append(preferred, fmt.Sprintf("%s:%s", s.ShortCode, s.KoreanName))
		}
	}
	fmt.Printf("\n우선주 %d 개 (앞 10):\n  %s\n",
		len(preferred), strings.Join(preferred[:min(10, len(preferred))], ", "))

	// KOSPI200 편입 종목
	var kospi200 int
	for _, s := range syms {
		if s.KOSPI200 == "Y" {
			kospi200++
		}
	}
	fmt.Printf("\nKOSPI200 편입: %d 개\n", kospi200)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
```

- [ ] **Step 2: 컴파일 검증**

Run: `go build ./examples/kospi_symbols && echo OK`
Expected: `OK`.

- [ ] **Step 3: Commit**

```bash
git add examples/kospi_symbols
git commit -m "$(cat <<'EOF'
[feat] examples/kospi_symbols — KOSPI 마스터 다운로드 + 통계

종목 수, 그룹코드 분포, 우선주, KOSPI200 편입 종목 출력. 첫 실행 시
KRX 마스터 ZIP 다운로드 + 디스크 캐시 (default TTL 7일).

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 16: 문서 갱신 (CLAUDE.md, README.md, domestic/doc.go)

**Files:**
- Modify: `CLAUDE.md`
- Modify: `README.md`
- Modify: `domestic/doc.go`

- [ ] **Step 1: CLAUDE.md 갱신** — "Phase 0 — skeleton" 줄을 갱신

`CLAUDE.md` 의 첫 line block (`> **Phase 0 — skeleton.**...`) 을 다음으로 교체:

```markdown
> **Phase 1.2 — domestic 시세/심볼/차트 (v0.2.0).** Phase 1.3+ 메서드는 추후 sub-plan 으로.
```

또한 spec 링크 line 에 Phase 1.2 plan 추가:

```markdown
- Phase 1.2 implementation plan: [`docs/superpowers/specs/2026-05-03-phase1-2-domestic-quotes-implementation-plan.md`](docs/superpowers/specs/2026-05-03-phase1-2-domestic-quotes-implementation-plan.md)
```

- [ ] **Step 2: README.md 의 Quick Start 갱신**

README 의 Quick Start 코드 (Phase 1.1 의 IssueAccessToken 호출만 있을 것) 아래에 Phase 1.2 메서드 사용 예 추가:

```markdown
## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    kis "github.com/kenshin579/korea-investment-stock"
    "github.com/kenshin579/korea-investment-stock/domestic"
)

func main() {
    client, err := kis.NewClientFromEnv()
    if err != nil {
        log.Fatal(err)
    }
    ctx := context.Background()

    // 1. 현재가
    price, _ := client.Domestic.InquirePrice(ctx, "005930")
    fmt.Printf("삼성전자 현재가: %s\n", price.StckPrpr)

    // 2. 일봉 차트
    chart, _ := client.Domestic.InquireDailyItemChartPrice(ctx, domestic.InquireDailyItemChartPriceParams{
        Symbol:   "005930",
        Period:   "D",
        FromDate: "20260401",
        ToDate:   "20260503",
    })
    fmt.Printf("일봉 %d 개\n", len(chart.Output2))

    // 3. KOSPI 종목 마스터
    syms, _ := client.Domestic.FetchKospiSymbols(ctx)
    fmt.Printf("KOSPI 종목 %d\n", len(syms))
}
```

## Available Methods (Phase 1.2)

| Method | 한투 path | TR_ID |
|--------|----------|-------|
| `Domestic.InquirePrice` | `inquire-price` | FHKST01010100 |
| `Domestic.SearchInfo` | `search-info` | CTPF1604R |
| `Domestic.SearchStockInfo` | `search-stock-info` | CTPF1002R |
| `Domestic.InquireDailyItemChartPrice` | `inquire-daily-itemchartprice` | FHKST03010100 |
| `Domestic.InquireTimeItemChartPrice` | `inquire-time-itemchartprice` | FHKST03010200 |
| `Domestic.FetchKospiSymbols` | (KRX 공개 마스터) | — |
| `Domestic.FetchKosdaqSymbols` | (KRX 공개 마스터) | — |
```

- [ ] **Step 3: domestic/doc.go 갱신**

`domestic/doc.go` 를 다음으로 교체:

```go
// Package domestic 은 한국투자증권 OpenAPI 의 국내주식 카테고리 메서드.
//
// Phase 1.2 의 7 메서드:
//
//   - InquirePrice                 — 주식현재가 시세 (FHKST01010100)
//   - SearchInfo                   — 상품기본조회 (CTPF1604R)
//   - SearchStockInfo              — 주식기본조회 (CTPF1002R)
//   - InquireDailyItemChartPrice   — 국내주식기간별시세 일/주/월/년 (FHKST03010100)
//   - InquireTimeItemChartPrice    — 주식당일분봉조회 (FHKST03010200)
//   - FetchKospiSymbols            — KRX KOSPI 마스터 (한투 API 가 아닌 KRX 공개 다운로드)
//   - FetchKosdaqSymbols           — KRX KOSDAQ 마스터
//
// 사용자는 root kis.Client 의 Domestic 필드로 접근.
package domestic
```

- [ ] **Step 4: 검증**

Run: `go build ./... && go vet ./... && gofmt -l . | tee /tmp/fmt.out`
Expected: 모두 출력 없음.

- [ ] **Step 5: Commit**

```bash
git add CLAUDE.md README.md domestic/doc.go
git commit -m "$(cat <<'EOF'
[doc] Phase 1.2 메서드 문서 갱신 — CLAUDE.md, README.md, domestic/doc.go

Phase 1.2 의 7 메서드 목록 + 한투 path/TR_ID 매핑 + Quick Start 코드 예시.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

## Task 17: 최종 점검

- [ ] **Step 1: 빌드/vet/fmt**

Run:
```bash
go build ./... && go vet ./... && gofmt -l . | tee /tmp/fmt.out
```
Expected: 빌드/vet 출력 없음. `gofmt -l` 출력 빈 파일.

- [ ] **Step 2: 모든 테스트 통과 + race**

Run: `go test ./... -race -count=1`
Expected: 모든 패키지 PASS.

- [ ] **Step 3: Coverage 측정**

Run:
```bash
go test ./... -coverprofile=/tmp/cov.out -covermode=atomic
go tool cover -func=/tmp/cov.out | tail -10
```
Expected: 마지막 줄 `total: (statements) ...` ≥ 80%. `domestic/` ≥ 80%, `internal/krxmaster/` ≥ 85%.

만약 부족하면 다음 영역 추가 테스트 후 fix commit:
- `domestic/symbols.go` 의 `downloadURL` 의 에러 경로 (network 에러, 빈 body 등)
- `internal/krxmaster/krxmaster.go` 의 `parseFwf`/`atoi64`/`toDecimal` 의 edge case (빈 문자열 → 0 등)

- [ ] **Step 4: 디렉터리 구조 확인**

Run:
```bash
find domestic internal/krxmaster examples -name '*.go' -o -name '*.json' -o -name '*.zip' | sort
```
Expected: 다음 패턴 모두 존재 —

```
domestic/chart.go
domestic/chart_test.go
domestic/client.go
domestic/doc.go
domestic/info.go
domestic/info_test.go
domestic/price.go
domestic/price_test.go
domestic/symbols.go
domestic/symbols_test.go
domestic/testdata/daily_chart_success.json
domestic/testdata/kosdaq_code_sample.mst.zip
domestic/testdata/kospi_code_sample.mst.zip
domestic/testdata/minute_chart_success.json
domestic/testdata/price_success.json
domestic/testdata/product_info_success.json
domestic/testdata/stock_info_success.json
domestic/testhelper_test.go
examples/domestic_chart/main.go
examples/domestic_price/main.go
examples/kospi_symbols/main.go
internal/krxmaster/doc.go
internal/krxmaster/krxmaster.go
internal/krxmaster/krxmaster_test.go
internal/krxmaster/testdata/kosdaq_code_sample.mst.zip
internal/krxmaster/testdata/kospi_code_sample.mst.zip
```

- [ ] **Step 5: Commit history 확인**

Run: `git log main..HEAD --oneline`
Expected: ~17 commits (각 task 의 commit + amendment commit).

(이 task 는 fix 가 필요하면 수행 후 commit; 이상 없으면 commit 없이 다음 단계로.)

---

## Task 18: PR 생성 및 push (사용자 승인 후)

> Claude 는 push / PR 생성을 사용자 명시적 승인 후에만 실행 (글로벌 정책).

- [ ] **Step 1: 사용자 승인 요청**

Claude 가 사용자에게:
- "Phase 1.2 모든 task 완료 + 최종 점검 통과. 지금 push + PR 생성하라" 또는
- "수정 사항 있다"

응답 받기.

- [ ] **Step 2: 승인 시 push**

Run:
```bash
git push -u origin docs/phase1-2-spec
```

- [ ] **Step 3: PR 생성** (gh CLI + HEREDOC, 글로벌 정책)

```bash
gh pr create --base main --title "Phase 1.2: 국내 시세+심볼+차트 (v0.2.0)" --body "$(cat <<'EOF'
## Summary

`korea-investment-stock` Go 라이브러리의 두 번째 release (v0.2.0) — 국내주식
시세/상품정보/차트/심볼 7 메서드 추가. 한투 API 문서 1:1 매핑 (Style A).

### 포함 내용

- **Phase 1 design spec amendment** — Phase 1.2 결정 반영 (root FetchPrice 통합 진입점 제거, 명명 스타일 A)
- **Phase 1.2 implementation plan** — `docs/superpowers/specs/2026-05-03-phase1-2-domestic-quotes-implementation-plan.md`
- **internal/krxmaster** — KOSPI/KOSDAQ 마스터 파서 (cp949 + fixed-width + struct 매핑). 핵심 필드 typed + Raw map fallback
- **domestic/price.go** — InquirePrice (주식현재가 시세, FHKST01010100)
- **domestic/info.go** — SearchInfo (상품기본조회, CTPF1604R), SearchStockInfo (주식기본조회, CTPF1002R)
- **domestic/chart.go** — InquireDailyItemChartPrice (FHKST03010100), InquireTimeItemChartPrice (FHKST03010200) + Params struct
- **domestic/symbols.go** — FetchKospiSymbols, FetchKosdaqSymbols (KRX 공개 마스터 다운로드 + mastercache)
- **domestic.New(http, master)** 시그니처 확장 — root client.go 의 wireInfra 가 c.masterC 주입
- **examples/** — domestic_price, domestic_chart, kospi_symbols
- **testdata/** — 한투 docs 응답 필드 정의 기반 합성 JSON 5종 + 실제 KRX 마스터 첫 3행 sample (KOSPI/KOSDAQ)

### 한투 spec 충실 원칙

Python wrapper 의 자체 편의 동작 (root FetchPrice 통합 진입점, ETF tr_id 분기,
fetch_stock_info 의 country_code fallback 루프, KR 검사 등) 모두 미구현.
한투 docs 의 query param 그대로 노출, PRDT_TYPE_CD 등은 사용자가 명시.

### 검증

- ✅ `go build ./...`
- ✅ `go vet ./...`
- ✅ `gofmt -l .`
- ✅ `go test ./... -race -count=1`
- ✅ Coverage ≥ 80% (`domestic/`, `internal/krxmaster/`)

### 다음 단계

- **Phase 1.3** (`v0.3.0`): 국내 순위 + 재무 (9 메서드)
- **Phase 1.4 ~ 1.5** 순차 진행
- 모두 끝나면 **`v1.0.0`** = Python parity 완성

## Test plan

- [x] 모든 단위 테스트 통과 (httpmock + testdata fixture)
- [x] race detector 통과
- [x] krxmaster 의 cp949 + fwf 파싱이 실제 KRX byte 첫 3행으로 검증
- [ ] examples 실행 검증 (사용자 manual, 실제 KIS credentials + KRX 다운로드)
- [ ] PR merge 후 v0.2.0 태그
EOF
)"
```

- [ ] **Step 4: PR URL 사용자에게 전달**

Run: `gh pr view --json url --jq '.url'`

PR URL 사용자에게 보여주고 review 안내.

- [ ] **Step 5: PR merge 후 v0.2.0 태그 (사용자가 승인 후)**

PR merge 후:
```bash
git checkout main && git pull
git tag -a v0.2.0 -m "v0.2.0: 국내 시세+심볼+차트 (Phase 1.2)"
git push origin v0.2.0
```

---

## Summary

### 포함 내용

- 7 메서드 (`Domestic.InquirePrice`, `SearchInfo`, `SearchStockInfo`, `InquireDailyItemChartPrice`, `InquireTimeItemChartPrice`, `FetchKospiSymbols`, `FetchKosdaqSymbols`)
- `internal/krxmaster/` 신규 패키지 (KOSPI/KOSDAQ 마스터 파서)
- `domestic.Client.New(http, master)` 시그니처 확장
- 5 typed struct (`Price`, `ProductInfo`, `StockInfo`, `DailyChart`, `MinuteChart`) + `KospiSymbol`, `KosdaqSymbol`
- 2 Params struct (`InquireDailyItemChartPriceParams`, `InquireTimeItemChartPriceParams`)
- 7 testdata fixture (5 합성 JSON + 2 KRX byte sample)
- 3 examples
- Phase 1 design spec amendment + Phase 1.2 plan

### 검증

- `go build ./...`
- `go vet ./...`
- `gofmt -l .`
- `go test ./... -race -count=1`
- Coverage ≥ 80%

### 다음 단계

- **Phase 1.3** (`v0.3.0`): 국내 순위 + 재무 (9 메서드 — `FetchVolumeRanking`, `FetchChangeRateRanking`, `FetchMarketCapRanking`, `FetchDividendRanking`, `FetchFinancialRatio`, `FetchIncomeStatement`, `FetchBalanceSheet`, `FetchProfitabilityRatio`, `FetchGrowthRatio`)
- 각 sub-plan 완료 후 release tag → 5개 모두 끝나면 **`v1.0.0`**

## Test plan

- [x] 17 task 모두 TDD 흐름 (테스트 작성 → fail → 구현 → pass → commit)
- [x] httpmock + testdata fixture 로 REST API 5종 단위 테스트
- [x] internal/krxmaster 의 cp949 + fwf 파싱을 실제 KRX byte sample 로 검증
- [x] domestic/symbols.go 의 다운로드 + 캐시 + 파싱 통합 테스트 (httpmock + 같은 sample)
- [x] race detector 통과
- [x] coverage ≥ 80%
- [ ] examples 실행 검증 (사용자 manual)
- [ ] PR merge 후 v0.2.0 태그

---

## Self-Review

### 1. Spec coverage

| spec 요구사항 (Phase 1 design + Phase 1.2 amendment) | 구현 task |
|---|---|
| `Domestic.InquirePrice` | Task 8 |
| `Domestic.SearchInfo` | Task 9 |
| `Domestic.SearchStockInfo` | Task 10 |
| `Domestic.InquireDailyItemChartPrice` | Task 11 |
| `Domestic.InquireTimeItemChartPrice` | Task 12 |
| `Domestic.FetchKospiSymbols` | Task 7 (+ Task 3 의 ParseKospi) |
| `Domestic.FetchKosdaqSymbols` | Task 7 (+ Task 4 의 ParseKosdaq) |
| `internal/krxmaster/` 신규 패키지 (cp949 + fwf) | Task 3, 4 |
| 핵심 typed + Raw map[한글컬럼명]값 | Task 3, 4 |
| `domestic.Client` 시그니처 확장 (master 추가) | Task 5 |
| Output1+Output2 한투 키 그대로 typed struct | Task 11, 12 |
| Params struct + zero-value default + OriginalPrice invert | Task 11 |
| 한투 spec 충실 원칙 (Python wrapper 동작 미구현) | Task 8 ~ 12 godoc 명시 |
| 메서드 명명 스타일 A (path 1:1) | 모든 task |
| KRX 마스터 다운로드: http.DefaultClient | Task 7 |
| testdata sample: 실제 KRX 첫 3행 | Task 2 |
| testdata 출처/라이선스 README | Task 2 |
| 응답 typed struct godoc 첫 줄에 path/TR_ID | Task 8, 9, 10, 11, 12 |
| Examples 3개 (domestic_price, domestic_chart, kospi_symbols) | Task 13, 14, 15 |
| 문서 갱신 (CLAUDE.md, README.md, domestic/doc.go) | Task 16 |
| 의존성 추가 (golang.org/x/text) | Task 1 |
| 검증 (build/vet/fmt/test/race/coverage) | Task 17 |
| PR 생성 (사용자 승인 후, gh + HEREDOC) | Task 18 |

### 2. Placeholder scan

- 모든 step 에 실제 코드 / 명령어 / expected 결과 포함
- "TBD", "TODO", "implement later" 사용 없음
- (Task 17 Step 3 의 coverage 부족 시 fix 안내가 "다음 영역 추가 테스트 후 fix commit" 으로 일반화 — 어느 영역인지 구체 명시. implementer 가 즉시 판단 가능)

### 3. Type consistency

- `Client` struct 의 `master *mastercache.Cache` 필드 (Task 5 추가) → Task 7 의 `c.master.Get(...)` 호출과 일치
- `KospiSymbol` / `KosdaqSymbol` struct 필드 (Task 3, 4 정의) → Task 7 의 `[]KospiSymbol` 반환 type alias (`type KospiSymbol = krxmaster.KospiSymbol`) 와 일치
- `Price` struct 의 80+ 필드 (Task 8) → Task 8 테스트의 `assert.Equal` 으로 검증되는 ~14 필드와 일치
- `InquireDailyItemChartPriceParams` (Task 11) 의 `Symbol`, `Period`, `FromDate`, `ToDate`, `OriginalPrice`, `MarketCode` 필드 → Task 11, 14 의 사용처와 일치
- `DailyChart{Output1, Output2}` / `MinuteChart{Output1, Output2}` 구조 일관 (output1+output2 한투 키 그대로)
- `InquireTimeItemChartPriceParams` 의 `PastDataInclude bool` → query param `FID_PW_DATA_INCU_YN` 의 "Y"/"N" 변환과 일치 (Task 12)
- `httpclient.Request` struct 의 필드명 (`Method`, `Path`, `TrID`, `Query`, `CustType`) → Phase 1.1 의 정의와 일치 (Task 8 ~ 12 모두 사용)

### 4. 위험 / 결함

- **cp949 vs EUC-KR**: `golang.org/x/text/encoding/korean.EUCKR` 가 cp949 확장 한자에 100% 대응 안 할 가능성. Task 3 의 단위 테스트가 한글 유니코드 범위 (`\x{AC00}-\x{D7A3}`) 매칭으로 검증하지만 한자 (Hanja) 종목명 들어 있는 행이 sample 첫 3행에 없으면 잠재 미발견. → Task 17 Step 3 의 coverage 단계에서 추가 sample 라인 (한자 포함) 으로 회귀 검증 가능.
- **KRX 마스터 fwf widths 오차**: Task 3, 4 의 `kospiFieldSpecs` / `kosdaqFieldSpecs` 가 Python `parsers/master_parser.py` 의 `field_specs` 와 1:1 매핑 명시했지만, byte-by-byte 검증은 단위 테스트의 `len(s.Raw) >= 60` 와 핵심 필드 비어있지 않음만 확인. 만약 widths 가 한 행 오차나면 컬럼 전체 어긋나는데, 한투 데이터의 `그룹코드` (앞부분) 와 `결산월` (중간) 둘 다 비어있지 않다면 정확도 신뢰 가능. PR review 단계에서 사용자가 examples/kospi_symbols 실행으로 실제 데이터 검증.
- **mastercache 디스크 캐시 격리**: Task 7 의 단위 테스트는 `mastercache.New(t.TempDir(), time.Hour)` 로 임시 디렉터리 사용 → 테스트 격리 OK. CI 에서 디스크 IO 권한 문제 없음.
- **httpmock + http.DefaultClient**: Task 7 의 `downloadURL` 이 `http.DefaultClient` 사용 → httpmock.Activate() 시 `http.DefaultTransport` 가 mock 으로 교체되어 자동 intercept. `Client.http` (resty 안의 customized transport) 와는 별개로 동작. 만약 사용자가 production 에서 httpmock.Activate 상태에서 KRX 다운로드 시도 시 mock 응답 반환 — 단 production 에서 httpmock 활성 안 되므로 무관.
- **Phase 1.1 의 client.go BC**: Task 5 가 `c.Domestic = domestic.New(c.httpClient, c.masterC)` 로 한 줄 변경. `c.masterC` 는 Phase 1.1 에서 이미 wireInfra 에 정의되어 있음 (`c.masterC = mastercache.New(masterDir, 7*24*time.Hour)`). → 의존성 깨짐 없음.
- **Output1+Output2 동시 unmarshal**: Task 11, 12 가 `httpclient.Response.Output` 만 unmarshal 안 하고 `resp.Raw` 전체에서 unmarshal. Phase 1.1 의 `httpclient.Response` struct 에 `Raw []byte` 필드가 있는지 확인 — 있음 (Phase 1.1 의 Task 9 에서 정의). OK.

(self-review 통과)

