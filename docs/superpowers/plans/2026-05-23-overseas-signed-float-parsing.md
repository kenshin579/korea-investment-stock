# 해외 응답 부호(+) float 파싱 버그 수정 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** KIS 해외 응답의 부호(`+`) 붙은 숫자 문자열이 `float64,string` 태그로 파싱되지 못해 해외 quote/chart 가 전부 깨지는 버그를, tolerant `kistypes.Float` 타입 도입으로 수정한다.

**Architecture:** 신규 leaf 패키지 `kistypes` 에 `float64` 기반 `Float` 타입(+ tolerant `UnmarshalJSON`)을 정의하고, `overseas`/`overseasfutures`/`futures` 3개 패키지의 모든 `float64,string` 필드(약 62개)를 이 타입으로 균일 교체한다. 소비자(moneyflow)는 SDK `v1.27.0` 으로 올린 뒤 컴파일 에러로 노출되는 변환 지점만 수정한다.

**Tech Stack:** Go 1.25+, `encoding/json`, `strconv`, `stretchr/testify`, `jarcoal/httpmock`. module `github.com/kenshin579/korea-investment-stock`.

**관련 스펙:** [`docs/superpowers/specs/2026-05-23-overseas-signed-float-parsing-design.md`](../specs/2026-05-23-overseas-signed-float-parsing-design.md)

**브랜치:** `fix/overseas-signed-float-parsing` (이미 생성됨)

---

## File Structure

| 파일 | 책임 | 작업 |
|------|------|------|
| `kistypes/float.go` | tolerant `Float` 타입 + `UnmarshalJSON` | Create |
| `kistypes/float_test.go` | `Float` 단위 테스트 | Create |
| `overseas/price.go`, `overseas/chart.go`, `overseas/ranking.go` | 응답 구조체 (21개 float64,string) | Modify |
| `overseas/testdata/price_detail_success.json`, `daily_price_success.json` | 회귀 픽스처(부호 추가) | Modify |
| `overseasfutures/quote.go`, `chart.go`, `options.go` | 응답 구조체 (5개) | Modify |
| `futures/board.go`, `chart.go`, `quote.go`, `conclusion.go` | 응답 구조체 (36개) | Modify |
| `overseas/signed_parsing_test.go`, `overseasfutures/signed_parsing_test.go`, `futures/signed_parsing_test.go` | 부호/빈문자열 회귀 테스트 | Create |
| `CHANGELOG.md` | v1.27.0 릴리스 노트 | Modify |
| (moneyflow) `backend/go.mod`, `backend/pkg/kis/client.go` | SDK 버전 bump + `float64()` 변환 | Modify |

---

## Task 1: `kistypes.Float` 타입

**Files:**
- Create: `kistypes/float.go`
- Test: `kistypes/float_test.go`

- [ ] **Step 1: 실패 테스트 작성** — `kistypes/float_test.go`

```go
package kistypes_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/kistypes"
)

func TestFloat_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    float64
		wantErr bool
	}{
		{"plus sign", `"+1.26"`, 1.26, false},   // 현재 깨지는 케이스
		{"minus sign", `"-1.26"`, -1.26, false},
		{"quoted plain", `"1.26"`, 1.26, false},
		{"unquoted number", `1.26`, 1.26, false},
		{"empty string", `""`, 0, false},        // KIS 가 종종 빈 값 반환
		{"null", `null`, 0, false},
		{"zero", `"0"`, 0, false},
		{"invalid", `"abc"`, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var f kistypes.Float
			err := json.Unmarshal([]byte(tt.input), &f)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.InDelta(t, tt.want, float64(f), 1e-9)
		})
	}
}

func TestFloat_StructField(t *testing.T) {
	var s struct {
		R kistypes.Float `json:"r"`
	}
	err := json.Unmarshal([]byte(`{"r":"+1.26"}`), &s)
	require.NoError(t, err)
	assert.InDelta(t, 1.26, float64(s.R), 1e-9)
}
```

- [ ] **Step 2: 테스트 실패 확인**

Run: `go test ./kistypes/...`
Expected: FAIL — `no required module provides package .../kistypes` 또는 `undefined: kistypes.Float` (구현 없음)

- [ ] **Step 3: 최소 구현 작성** — `kistypes/float.go`

```go
// Package kistypes 는 KIS API 응답 파싱을 위한 공용 타입을 제공한다.
package kistypes

import (
	"strconv"
	"strings"
)

// Float 는 KIS 응답의 부호(+/-) 붙은 숫자 문자열과 빈 문자열을 안전하게 파싱하는 float64.
//
// 표준 encoding/json 의 `,string` 태그는 따옴표 안 값이 JSON number 문법이길 요구하는데,
// JSON number 는 leading '+' 를 금지한다. KIS 해외 API 는 등락 필드를 "+1.26" 형태로 주므로
// 그 태그로는 상승 값에서 파싱이 실패한다. 이 타입이 그 대체다.
type Float float64

// UnmarshalJSON 은 다음을 허용한다:
//   - "null" / 빈 입력 → 0
//   - 따옴표로 감싼 숫자 문자열(KIS 기본): "+1.26", "-1.26", "1.26", ""(→0)
//   - 따옴표 없는 JSON number: 1.26
func (f *Float) UnmarshalJSON(b []byte) error {
	s := strings.TrimSpace(string(b))
	if s == "" || s == "null" {
		*f = 0
		return nil
	}
	s = strings.TrimSpace(strings.Trim(s, `"`))
	if s == "" {
		*f = 0
		return nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*f = Float(v)
	return nil
}
```

- [ ] **Step 4: 테스트 통과 확인**

Run: `go test ./kistypes/...`
Expected: PASS (8 subtests + struct field test)

- [ ] **Step 5: 커밋**

```bash
git add kistypes/float.go kistypes/float_test.go
git commit -m "feat(kistypes): + 부호/빈문자열 허용하는 Float 타입 추가"
```

---

## Task 2: `overseas` 패키지 마이그레이션 + 회귀

**Files:**
- Modify: `overseas/testdata/price_detail_success.json`, `overseas/testdata/daily_price_success.json`
- Modify: `overseas/price.go`, `overseas/chart.go`, `overseas/ranking.go`
- Create: `overseas/signed_parsing_test.go`

- [ ] **Step 1: 픽스처에 부호/빈문자열 반영(테스트 갭 메움)**

`overseas/testdata/price_detail_success.json` 의 `output` 안 값 변경:
- `"t_xrat": "0.69"` → `"t_xrat": "+0.69"`
- `"p_xrat": "0"` → `"p_xrat": ""`  (빈 문자열 케이스)

`overseas/testdata/daily_price_success.json` 의 `output2` 첫 행:
- `"rate": "0.69"` → `"rate": "+0.69"`

- [ ] **Step 2: 회귀 테스트 작성** — `overseas/signed_parsing_test.go`

```go
package overseas_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/overseas"
)

// KIS 해외 응답은 등락 필드를 "+0.69" 처럼 부호 붙은 문자열로, 빈 값은 "" 로 준다.
// 회귀: float64,string 태그로는 '+' 에서 파싱이 깨졌다 (kistypes.Float 로 수정).
func TestPriceDetailSnapshot_SignedRate(t *testing.T) {
	const body = `{"rsym":"DNASAAPL","t_xrat":"+0.69","p_xrat":"","perx":"29.45","tvol":"100","pvol":"200","tomv":"1","pamt":"1","shar":"1","mcap":"1","tamt":"1"}`
	var got overseas.PriceDetailSnapshot
	err := json.Unmarshal([]byte(body), &got)
	require.NoError(t, err)
	assert.InDelta(t, 0.69, float64(got.TXrat), 1e-9)
	assert.InDelta(t, 0, float64(got.PXrat), 1e-9)
	assert.InDelta(t, 29.45, float64(got.Perx), 1e-9)
}

func TestDailyPriceCandle_SignedRate(t *testing.T) {
	const body = `{"xymd":"20260505","rate":"+0.69","tvol":"1","tamt":"1","vbid":"1","vask":"1"}`
	var got overseas.DailyPriceCandle
	err := json.Unmarshal([]byte(body), &got)
	require.NoError(t, err)
	assert.InDelta(t, 0.69, float64(got.Rate), 1e-9)
}
```

- [ ] **Step 3: 테스트 실패 확인 (RED)**

Run: `go test ./overseas/... -run 'SignedRate|InquirePriceDetail|InquireDailyPrice'`
Expected: FAIL — `invalid use of ,string struct tag, trying to unmarshal "+0.69" into float64`

- [ ] **Step 4: float64,string → kistypes.Float 균일 교체**

```bash
perl -0777 -pi -e 's/float64(\s+)`json:"([^"]+),string"`/kistypes.Float$1`json:"$2"`/g' \
  overseas/price.go overseas/chart.go overseas/ranking.go
```

각 파일 import 블록에 추가(이미 있으면 생략):

```go
	"github.com/kenshin579/korea-investment-stock/kistypes"
```

정렬/import 정리:

```bash
gofmt -w overseas/price.go overseas/chart.go overseas/ranking.go
goimports -w overseas/price.go overseas/chart.go overseas/ranking.go 2>/dev/null || true
```

- [ ] **Step 5: 잔여 float64,string 0건 확인**

Run: `grep -rn 'float64.*,string' overseas/*.go | grep -v _test`
Expected: (출력 없음)

- [ ] **Step 6: 테스트 통과 확인 (GREEN)**

Run: `go test ./overseas/...`
Expected: PASS (회귀 테스트 + 기존 테스트 전부)

- [ ] **Step 7: 커밋**

```bash
git add overseas/
git commit -m "fix(overseas): 부호 붙은 등락 필드를 kistypes.Float 로 파싱"
```

---

## Task 3: `overseasfutures` 패키지 마이그레이션 + 회귀

**Files:**
- Modify: `overseasfutures/quote.go`, `overseasfutures/chart.go`, `overseasfutures/options.go`
- Create: `overseasfutures/signed_parsing_test.go`

- [ ] **Step 1: 회귀 테스트 작성** — `overseasfutures/signed_parsing_test.go`

```go
package overseasfutures_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/overseasfutures"
)

// 회귀: 전일대비율(prev_diff_rate)이 "+1.50" 부호 문자열로 와도 파싱돼야 한다.
func TestInquirePriceOutput1_SignedDiffRate(t *testing.T) {
	const body = `{"prev_diff_rate":"+1.50"}`
	var got overseasfutures.InquirePriceOutput1
	err := json.Unmarshal([]byte(body), &got)
	require.NoError(t, err)
	assert.InDelta(t, 1.50, float64(got.PrevDiffRate), 1e-9)
}
```

- [ ] **Step 2: 테스트 실패 확인 (RED)**

Run: `go test ./overseasfutures/... -run SignedDiffRate`
Expected: FAIL — `invalid use of ,string struct tag, trying to unmarshal "+1.50" into float64`

- [ ] **Step 3: float64,string → kistypes.Float 균일 교체**

```bash
perl -0777 -pi -e 's/float64(\s+)`json:"([^"]+),string"`/kistypes.Float$1`json:"$2"`/g' \
  overseasfutures/quote.go overseasfutures/chart.go overseasfutures/options.go
```

각 파일 import 블록에 추가(이미 있으면 생략):

```go
	"github.com/kenshin579/korea-investment-stock/kistypes"
```

정렬/import 정리:

```bash
gofmt -w overseasfutures/quote.go overseasfutures/chart.go overseasfutures/options.go
goimports -w overseasfutures/quote.go overseasfutures/chart.go overseasfutures/options.go 2>/dev/null || true
```

- [ ] **Step 4: 잔여 float64,string 0건 확인**

Run: `grep -rn 'float64.*,string' overseasfutures/*.go | grep -v _test`
Expected: (출력 없음)

- [ ] **Step 5: 테스트 통과 확인 (GREEN)**

Run: `go test ./overseasfutures/...`
Expected: PASS

- [ ] **Step 6: 커밋**

```bash
git add overseasfutures/
git commit -m "fix(overseasfutures): prev_diff_rate 를 kistypes.Float 로 파싱"
```

---

## Task 4: `futures` 패키지 마이그레이션 + 회귀

**Files:**
- Modify: `futures/board.go`, `futures/chart.go`, `futures/quote.go`, `futures/conclusion.go`
- Create: `futures/signed_parsing_test.go`

- [ ] **Step 1: 회귀 테스트 작성** — `futures/signed_parsing_test.go`

```go
package futures_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/korea-investment-stock/futures"
)

// 회귀: 대비율(+0.85)·괴리율(빈 값) 등이 부호/빈문자열로 와도 파싱돼야 한다.
func TestInquirePriceOutput1_SignedFields(t *testing.T) {
	const body = `{"futs_prdy_ctrt":"+0.85","dprt":""}`
	var got futures.InquirePriceOutput1
	err := json.Unmarshal([]byte(body), &got)
	require.NoError(t, err)
	assert.InDelta(t, 0.85, float64(got.FutsPrdyCtrt), 1e-9)
	assert.InDelta(t, 0, float64(got.Dprt), 1e-9)
}
```

- [ ] **Step 2: 테스트 실패 확인 (RED)**

Run: `go test ./futures/... -run SignedFields`
Expected: FAIL — `invalid use of ,string struct tag, trying to unmarshal "+0.85" into float64`

- [ ] **Step 3: float64,string → kistypes.Float 균일 교체**

```bash
perl -0777 -pi -e 's/float64(\s+)`json:"([^"]+),string"`/kistypes.Float$1`json:"$2"`/g' \
  futures/board.go futures/chart.go futures/quote.go futures/conclusion.go
```

각 파일 import 블록에 추가(이미 있으면 생략):

```go
	"github.com/kenshin579/korea-investment-stock/kistypes"
```

정렬/import 정리:

```bash
gofmt -w futures/board.go futures/chart.go futures/quote.go futures/conclusion.go
goimports -w futures/board.go futures/chart.go futures/quote.go futures/conclusion.go 2>/dev/null || true
```

- [ ] **Step 4: 잔여 float64,string 0건 확인**

Run: `grep -rn 'float64.*,string' futures/*.go | grep -v _test`
Expected: (출력 없음)

- [ ] **Step 5: 테스트 통과 확인 (GREEN)**

Run: `go test ./futures/...`
Expected: PASS

- [ ] **Step 6: 커밋**

```bash
git add futures/
git commit -m "fix(futures): 대비율/그릭스/괴리 필드를 kistypes.Float 로 파싱"
```

---

## Task 5: SDK 전체 검증 + 릴리스 (v1.27.0)

**Files:**
- Modify: `CHANGELOG.md`

- [ ] **Step 1: 전 패키지 빌드/vet/테스트**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: 전부 통과 (0 fail)

- [ ] **Step 2: 잔여 대상 0건 최종 확인**

Run: `grep -rn 'float64.*,string' overseas/*.go overseasfutures/*.go futures/*.go | grep -v _test`
Expected: (출력 없음)

- [ ] **Step 3: 실 API 재현 확인 (수동, KIS 키 필요)**

Run:
```bash
go run ./examples/overseas_price
go run ./examples/overseas_chart
```
Expected: 파싱 에러 없이 AAPL 현재가/일봉 출력. (redis 연결 거부 로그는 무관 — 토큰 캐시용)
Note: KIS 토큰 발급은 분당 1회 제한 — 두 예제 사이 간격 필요. redis 미연결 시 매 호출 토큰 재발급으로 403 가능 → 시간 두고 재시도.

- [ ] **Step 4: CHANGELOG 갱신** — `CHANGELOG.md` 최상단에 추가

```markdown
## v1.27.0

### Fixed
- 해외 응답의 부호(`+`) 붙은 등락 문자열(`"+1.26"`) 파싱 실패 수정. `float64,string` 태그는
  JSON number 의 leading `+` 를 거부해 해외 quote/chart 가 상승 값에서 전부 깨졌다. 빈 문자열
  `""` 도 동일 실패였다.

### Changed (BREAKING — response struct field types)
- `overseas`/`overseasfutures`/`futures` 패키지의 `float64,string` 응답 필드(약 62개)를
  신규 `kistypes.Float` 타입으로 교체. 소비자는 `float64(x.Field)` 변환 필요(컴파일 타임 노출).

### Added
- `kistypes` 패키지 — `+`/`-`/빈 문자열을 허용하는 tolerant `Float` 타입.
```

- [ ] **Step 5: 커밋**

```bash
git add CHANGELOG.md
git commit -m "docs: CHANGELOG v1.27.0 (해외 부호 float 파싱 수정)"
```

- [ ] **Step 6: PR 생성 (리뷰어 미지정)**

```bash
git push -u origin fix/overseas-signed-float-parsing
gh pr create --title "fix: 해외 응답 부호(+) float 파싱 버그 수정 (v1.27.0)" --body "$(cat <<'EOF'
## Summary
- KIS 해외 응답의 "+1.26" 부호 붙은 등락 필드가 `float64,string` 태그로 파싱 불가(JSON number 는 leading + 금지)하여 해외 quote/chart 전반이 깨지던 버그 수정
- 신규 `kistypes.Float`(tolerant UnmarshalJSON) 도입, overseas/overseasfutures/futures 의 float64,string ~62개 필드 균일 교체
- 빈 문자열 "" 잠복 버그도 함께 해소

## Test plan
- [ ] `go test ./...` 전체 통과
- [ ] `kistypes` 단위 테스트(+/-/""/null/invalid)
- [ ] overseas/overseasfutures/futures 회귀 테스트(부호 문자열)
- [ ] `go run ./examples/overseas_price`, `./examples/overseas_chart` 정상 출력
EOF
)"
```

- [ ] **Step 7: (머지 후) 태그**

> 사람이 PR 리뷰/머지 후 진행. main 에서:
```bash
git tag v1.27.0 && git push origin v1.27.0
```

---

## Task 6: moneyflow 소비자 업데이트

> **의존:** Task 5 Step 7 (`v1.27.0` 태그) 완료 후 진행. 작업 디렉토리: `moneyflow.advenoh.pe.kr/backend`. 별도 feature 브랜치(`fix/kis-signed-float-v1.27.0`).

**Files:**
- Modify: `backend/go.mod`
- Modify: `backend/pkg/kis/client.go:313`

- [ ] **Step 1: feature 브랜치 + SDK 버전 bump**

```bash
cd moneyflow.advenoh.pe.kr && git checkout main && git pull origin main
git checkout -b fix/kis-signed-float-v1.27.0
cd backend
go get github.com/kenshin579/korea-investment-stock@v1.27.0
go mod tidy
```

- [ ] **Step 2: 빌드로 변환 지점 노출 (안전망 확인)**

Run: `go build ./...`
Expected: FAIL — `pkg/kis/client.go:313` 근처 `cannot use s.PrdyCtrt (variable of type kistypes.Float) as float64 value in argument to decimal.NewFromFloat`

- [ ] **Step 3: float64 변환 적용** — `backend/pkg/kis/client.go:313`

기존:
```go
	pct := decimal.NewFromFloat(s.PrdyCtrt)
```
변경:
```go
	pct := decimal.NewFromFloat(float64(s.PrdyCtrt))
```

(빌드가 추가 지점을 노출하면 동일하게 `float64(...)` 로 감싼다. 단, domestic 필드는 미대상이므로 건드리지 않는다.)

- [ ] **Step 4: 빌드/테스트 통과 확인**

Run: `go build ./... && go test ./...`
Expected: PASS

- [ ] **Step 5: AAPL 해외 조회 동작 확인 (로컬 또는 배포 후)**

Run (백엔드 기동 후):
```bash
curl -s -w "\n[%{http_code}]\n" "http://localhost:8080/v1/stock/AAPL/chart?range=1y" | head -c 300
curl -s -w "\n[%{http_code}]\n" "http://localhost:8080/v1/stock/AAPL/quote" | head -c 300
```
Expected: 둘 다 `[200]`, chart 는 `candles` 배열 존재.

- [ ] **Step 6: 커밋 + PR**

```bash
git add backend/go.mod backend/go.sum backend/pkg/kis/client.go
git commit -m "fix: KIS SDK v1.27.0 (해외 부호 float 파싱) 반영"
git push -u origin fix/kis-signed-float-v1.27.0
gh pr create --title "fix: 해외 주식 차트/시세 로딩 복구 (KIS SDK v1.27.0)" --body "$(cat <<'EOF'
## Summary
- KIS SDK v1.27.0 으로 업데이트하여 해외 주식 quote/chart 502/500 버그 수정
- `client.go:313` overseas 지수 등락률을 `float64()` 변환

## Test plan
- [ ] `go build ./... && go test ./...`
- [ ] `/v1/stock/AAPL/chart` 200 + candles 존재
- [ ] `/v1/stock/AAPL/quote` 200
EOF
)"
```

- [ ] **Step 7: 배포 후 운영 확인**

머지/재배포 후 `https://moneyflow.advenoh.pe.kr/stock/AAPL` 차트가 보이는지 확인(상승장 포함).

---

## Known Follow-ups (이번 범위 밖)

- `domestic` 패키지 271개 `float64,string` — 현재 sign-코드 방식이라 `+` 안 옴. 빈 문자열 `""` 잠복 위험만 존재. 필요 시 동일 `kistypes.Float` 점진 적용.
- `int64,string` 필드(거래량 등) — `+` 없음, `""` 잠복 위험만. 필요 시 `kistypes.Int` 도입.
