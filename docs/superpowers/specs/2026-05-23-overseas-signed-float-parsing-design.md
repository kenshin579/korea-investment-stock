# 해외 응답 부호(+) float 파싱 버그 수정 설계

- 작성일: 2026-05-23
- 대상 저장소: `korea-investment-stock` (근본 원인) / 소비자: `moneyflow.advenoh.pe.kr`
- 릴리스: SDK `v1.27.0` (minor)

## 1. 배경 / 증상

`https://moneyflow.advenoh.pe.kr/stock/AAPL` 의 해외 주식 차트가 로딩되지 않음. 국내 주식
(예: `005930`) 차트는 정상.

프로덕션 API 재현:

| 요청 | 결과 |
|------|------|
| `GET /api/v1/stock/AAPL/chart?range=1y` | HTTP **502** `KIS_UPSTREAM_ERROR` |
| `GET /api/v1/stock/AAPL/quote` | HTTP **500** |
| `GET /api/v1/stock/005930/chart?range=1y` | HTTP **200** (정상) |
| `GET /api/v1/stock/005930/quote` | HTTP **200** (정상) |

차트뿐 아니라 **해외 KIS 조회 전반**(quote 포함)이 깨져 있었음. MSFT 등 다른 해외 종목도 동일.

## 2. 근본 원인

워크스페이스 SDK 예제로 직접 재현해 정확한 에러를 포착:

```
InquirePriceDetail: kis: parse PriceDetail: json: invalid use of ,string struct tag,
  trying to unmarshal "+1.26" into float64
InquireDailyPrice:  kis: parse DailyPrice:  json: invalid use of ,string struct tag,
  trying to unmarshal "+1.26" into float64
```

**메커니즘:** KIS 해외 API 는 등락 관련 필드를 `"+1.26"` / `"-1.26"` 처럼 **부호 prefix 가
붙은 문자열**로 반환한다. SDK 응답 구조체는 이 필드를 `float64` + `` `json:"...,string"` ``
태그로 파싱하는데, Go `encoding/json` 의 `,string` 옵션은 따옴표 안 값이 **JSON number 문법**
이어야 한다. JSON number 는 leading `-` 만 허용하고 **leading `+` 는 금지**한다. 따라서
상승(`+`) 값에서만 파싱이 실패하고, 구조체 전체 언마샬이 깨진다(해당 필드를 읽지 않아도 실패).

검증 (Go 표준 라이브러리 동작):

| 입력 (`json:",string"` + `float64`) | 결과 |
|------|------|
| `"+1.26"` (상승) | ❌ `invalid use of ,string struct tag` |
| `"-1.26"` (하락) | ✅ `-1.26` |
| `"1.26"` / `"0.00"` | ✅ |

추가 확인: `strconv.ParseFloat("+1.26")` 는 **정상**(=1.26). 즉 문제는 오직 `encoding/json`
의 `,string` 태그이며, 커스텀 `UnmarshalJSON` 에서 `strconv.ParseFloat` 를 쓰면 부호 strip
없이도 해결된다. 단 `strconv.ParseFloat("")` 는 에러이므로 빈 문자열은 별도 처리 필요.

### 왜 국내(domestic)는 멀쩡한가

KIS 국내 응답은 부호를 별도 sign 코드 필드(`prdy_vrss_sign` 등)로 주고 크기는 부호 없이
반환한다. 그래서 동일한 `float64,string` 필드라도 `+` 가 오지 않아 깨지지 않았다.

### 간헐성 주의

시장이 하락(`-`)이면 통과하므로 상승장에서만, 또는 종목별로 들쭉날쭉하게 보일 수 있다. 작성
시점에 AAPL·MSFT 등이 상승세라 전부 실패 상태였다.

## 3. 해결 방안 (Approach 1 — tolerant float 타입)

### 3.1 새 패키지 / 타입

신규 leaf 패키지 `github.com/kenshin579/korea-investment-stock/kistypes` 에 `float64` 기반
명명 타입을 정의한다. 루트 `kis` 패키지가 `overseas`/`futures`/`overseasfutures` 를 import
하므로, 새 타입은 이들이 모두 의존 가능한 leaf 패키지에 두어야 import cycle 이 없고, exported
라야 소비자(moneyflow)가 이름을 사용할 수 있다.

```go
package kistypes

// Float 는 KIS 응답의 부호(+/-) 붙은 숫자 문자열과 빈 문자열을 안전하게 파싱하는 float64.
// 표준 encoding/json 의 `,string` 태그는 leading '+' 를 거부하므로 그 대체.
type Float float64

func (f *Float) UnmarshalJSON(b []byte) error {
    // 1) "null" 또는 빈 입력 → 0
    // 2) 양끝 따옴표 제거 (KIS 는 문자열로 주지만 숫자로 올 가능성도 허용)
    // 3) 트림 후 "" → 0  (KIS 가 종종 빈 값 반환)
    // 4) strconv.ParseFloat (부호 +/- 자동 처리) → 실패 시 error
}
```

**동작 규약:**

| 입력 | 결과 |
|------|------|
| `"+1.26"` | `1.26` (현재 깨지는 케이스) |
| `"-1.26"` | `-1.26` |
| `"1.26"` / `1.26` (따옴표 없음) | `1.26` |
| `""` / `null` | `0` (에러 아님) |
| `"abc"` | error |

소비자는 `float64(x.Field)` 로 변환해 사용한다(필드 타입 변경 → 컴파일 타임에 모든 변경 지점 노출).

### 3.2 적용 범위 — 균일 교체

대상 3개 패키지(`overseas`, `overseasfutures`, `futures`) 내의 **모든 `float64,string`
필드(약 62개: overseas 21 / overseasfutures 5 / futures 36)를 `kistypes.Float` 로 균일
교체**하고 `,string` 태그를 제거한다.

"부호 가능 필드만 선별"하지 않는 이유: (1) 빈 문자열 `""` 도 `,string` 에서 동일하게 깨지므로
부호 없는 필드(PER/PBR/환율)도 잠복 버그다. (2) `kistypes.Float` 는 부호 없는 값도 동일하게
처리하므로 단점이 없고, 균일 규칙이 리뷰를 단순화한다.

| 패키지 | 파일 | 대표 필드 |
|--------|------|-----------|
| `overseas` (21) | `price.go`(6), `chart.go`(2), `ranking.go`(13) | `t_xrat`/`p_xrat`(등락·**확인된 버그**), `rate`/`prdy_ctrt`, `perx`/`pbrx`, 환율, 체결강도 |
| `overseasfutures` (5) | `chart.go`, `quote.go`, `options.go` | `prev_diff_rate`(전일대비율) |
| `futures` (36) | `board.go`, `chart.go`, `quote.go`, `conclusion.go` | `*_prdy_ctrt`(대비율), 그릭스(`delta`/`gama`/`vega`/`theta`/`rho`), `dprt`/`esdg`(괴리), 변동성 |

### 3.3 범위 밖 (의도적 제외 / 알려진 후속 과제)

- **`domestic` (271개 `float64,string`)** — 현재 sign-코드 방식이라 `+` 가 오지 않음. 단 빈
  문자열 `""` 잠복 위험은 남음. 변경 폭이 과대하여 이번 범위에서 제외. follow-up 후보.
- **`int64,string` 필드** (거래량 등) — `+` 없음, `""` 잠복 위험만. 필요 시 추후 `kistypes.Int`
  도입. 현재는 YAGNI.

## 4. 소비자(moneyflow) 영향 & 버전

### 4.1 영향 (최소)

`moneyflow.advenoh.pe.kr/backend` 에서 마이그레이션 대상 구조체의 float 필드를 읽는 곳:

- `pkg/kis/client.go:313` — `decimal.NewFromFloat(s.PrdyCtrt)` (overseas `InquireDailyChartPrice`
  Output1) → `decimal.NewFromFloat(float64(s.PrdyCtrt))` 로 변환 필요.
- 그 외 사용처 없음: 해외 quote 는 가격(`decimal`)·거래량(`int`)만 읽고, chart 어댑터는 OHLC 만
  읽는다. `client.go:238`(domestic `InquirePrice`), `client.go:512`(domestic `InquireEtfPrice`)
  는 domestic 이라 미대상.
- **안전망:** SDK 교체 후 moneyflow `go build ./...` 시 누락된 변환이 컴파일 에러로 전부 노출.
  보이지 않는 런타임 깨짐 없음.

### 4.2 버전

- exported 응답 구조체 필드 타입이 `float64` → `kistypes.Float` 로 바뀌어 엄밀히는 breaking.
- 이 SDK 는 workspace 내부(moneyflow 단일 소비자) 사용이고 변경 본질은 버그 수정.
- **`v1.27.0` (minor bump)** + `CHANGELOG.md` 에 `BREAKING (response struct field types):
  float64 → kistypes.Float` 명시.

## 5. 테스트 계획 (TDD)

기존 테스트가 통과하면서 운영이 깨진 이유는 **픽스처(`testdata/*.json`)에 `+` 부호 값이 없었기**
때문이다(테스트 갭). 이 갭을 메우는 것이 회귀 방지의 핵심이다.

1. **`kistypes.Float` 단위 테스트** (`kistypes/float_test.go`) — 실패 케이스부터(RED→GREEN):
   테이블 테스트로 §3.1 동작 규약 전부(`"+1.26"`, `"-1.26"`, `"1.26"`, 따옴표 없는 `1.26`,
   `""`, `null`, `"abc"`) + 구조체 필드 임베드 언마샬 1개.
2. **회귀 테스트 — 픽스처 보강:** `overseas/testdata/price_detail_success.json`,
   `daily_price_success.json` 의 등락 필드를 `"+1.26"` 형태로 갱신(또는 `_signed.json` 추가).
   `overseas`/`overseasfutures`/`futures` 각 패키지에 `+` 부호 + `""` 포함 픽스처로 언마샬
   성공 + 값 검증 테스트 추가. 교체 전 RED → 후 GREEN.
3. **실 API 재현 확인(수동/integration):** `go run ./examples/overseas_price`,
   `./examples/overseas_chart` 가 현재 파싱 에러 → 수정 후 정상 출력.
4. **moneyflow 검증:** `go build ./...` 로 안전망 동작 확인 후 `float64()` 변환 수정;
   `/api/v1/stock/AAPL/chart` 200 + candles 존재(integration 또는 E2E).
5. **전 패키지 회귀:** `go test ./...` 전체 GREEN.

## 6. 작업 순서 & 롤아웃

**SDK (`korea-investment-stock`, 브랜치 `fix/overseas-signed-float-parsing`):**

1. `kistypes` 패키지 + `Float`/`UnmarshalJSON` (단위 테스트 RED→GREEN).
2. 픽스처에 `+`/`""` 케이스 반영(회귀 테스트 RED 확인).
3. 3개 패키지의 `float64,string` 약 62개 필드 → `kistypes.Float` 균일 교체(`,string` 제거).
4. `go test ./...` 전체 GREEN + 예제 2종 수동 확인.
5. `CHANGELOG.md` `## v1.27.0` 항목(버그 수정 + BREAKING 명시).
6. PR(`gh pr create` + HEREDOC, 리뷰어 미지정) → main merge 후 `v1.27.0` 태그.

**moneyflow (`moneyflow.advenoh.pe.kr`, 별도 브랜치):**

7. `go.mod` `v1.26.0 → v1.27.0` (`go get` / `go mod tidy`).
8. `go build ./...` → 컴파일 에러 노출 지점 수정(최소 `client.go:313`).
9. AAPL chart/quote integration·E2E → PR → merge → 재배포.
10. 운영에서 `https://moneyflow.advenoh.pe.kr/stock/AAPL` 차트 정상 확인(상승장 포함).

**위험/롤백:**

- SDK 변경은 타입 교체뿐이라 위험 낮음. moneyflow 는 SDK 버전 pin 이므로 문제 시 `go.mod` 를
  `v1.26.0` 으로 되돌리면 즉시 롤백(단, 해외 버그는 재현).
- §3.3 의 잠복 위험(domestic 271개, `int64,string` 의 `""`)은 이번 범위 밖, 후속 과제로 기록.

**의존:** 8~10 은 6(`v1.27.0` 태그) 완료에 의존. SDK PR 이 먼저.
