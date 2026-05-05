# korea-investment-stock (Go)

[![Go Reference](https://pkg.go.dev/badge/github.com/kenshin579/korea-investment-stock.svg)](https://pkg.go.dev/github.com/kenshin579/korea-investment-stock)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**한국투자증권 OpenAPI Go 클라이언트** — typed struct, context-first, functional options, 자동 토큰 갱신/rate limit 내장.

> ⚠️ **Work in progress.** 이 라이브러리는 현재 초기 개발 단계입니다 (`v0.x`). 안정화 시 `v1.0.0` 으로 올릴 예정입니다.

## Python 사용자에게

이 repo 는 **2026-05-03 부로 Python → Go 로 전환**되었습니다.

- 기존 Python 코드 (`v0.18.0` 까지 + `v0.19.0` deprecation release): [`python-final`](https://github.com/kenshin579/korea-investment-stock/tree/python-final) 태그로 영구 보존
- PyPI 패키지 (`korea-investment-stock`): `v0.19.0` 까지 그대로 유지. critical security fix 외 신규 기능 없음.
- 마이그레이션 배경: [Phase 0 design spec](docs/superpowers/specs/2026-05-03-korea-investment-go-migration-design.md)

## Install

```bash
go get github.com/kenshin579/korea-investment-stock
```

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

## Available Methods (Phase 1.2 ~ 2.2)

| Method | 한투 path | TR_ID |
|--------|----------|-------|
| `Domestic.InquirePrice` | `inquire-price` | FHKST01010100 |
| `Domestic.SearchInfo` | `search-info` | CTPF1604R |
| `Domestic.SearchStockInfo` | `search-stock-info` | CTPF1002R |
| `Domestic.InquireDailyItemChartPrice` | `inquire-daily-itemchartprice` | FHKST03010100 |
| `Domestic.InquireTimeItemChartPrice` | `inquire-time-itemchartprice` | FHKST03010200 |
| `Domestic.FetchKospiSymbols` | (KRX 공개 마스터) | — |
| `Domestic.FetchKosdaqSymbols` | (KRX 공개 마스터) | — |
| `Domestic.InquireVolumeRank` | `quotations/volume-rank` | FHPST01710000 |
| `Domestic.InquireFluctuation` | `ranking/fluctuation` | FHPST01700000 |
| `Domestic.InquireMarketCap` | `ranking/market-cap` | FHPST01740000 |
| `Domestic.InquireDividendRate` | `ranking/dividend-rate` | HHKDB13470100 |
| `Domestic.InquireFinancialRatio` | `finance/financial-ratio` | FHKST66430300 |
| `Domestic.InquireIncomeStatement` | `finance/income-statement` | FHKST66430200 |
| `Domestic.InquireBalanceSheet` | `finance/balance-sheet` | FHKST66430100 |
| `Domestic.InquireProfitRatio` | `finance/profit-ratio` | FHKST66430400 |
| `Domestic.InquireGrowthRatio` | `finance/growth-ratio` | FHKST66430800 |
| `Domestic.InquireInvestorTradeByStockDaily` | `quotations/investor-trade-by-stock-daily` | FHPTJ04160001 |
| `Domestic.InquireInvestorDailyByMarket` | `quotations/inquire-investor-daily-by-market` | FHPTJ04040000 |
| `Domestic.InquireInvestorTimeByMarket` | `quotations/inquire-investor-time-by-market` | FHPTJ04030000 |
| `Domestic.InquireIndexPrice` | `quotations/inquire-index-price` | FHPUP02100000 |
| `Domestic.InquireIndexCategoryPrice` | `quotations/inquire-index-category-price` | FHPUP02140000 |
| `Domestic.InquirePubOffer` | `ksdinfo/pub-offer` | HHKDB669108C0 |
| `Overseas.InquirePriceDetail` | `overseas-price/v1/quotations/price-detail` | HHDFS76200200 |
| `Overseas.SearchInfo` | `overseas-price/v1/quotations/search-info` | CTPF1702R |
| `Overseas.InquireDailyPrice` | `overseas-price/v1/quotations/dailyprice` | HHDFS76240000 |
| `Overseas.InquireDailyChartPrice` | `overseas-price/v1/quotations/inquire-daily-chartprice` | FHKST03030100 |
| `Overseas.InquireUpdownRate` | `overseas-stock/v1/ranking/updown-rate` | HHDFS76290000 |
| `Overseas.FetchOverseasSymbols(market)` | (KIS 공개 마스터 — 11 거래소) | — |
| `Domestic.InquireAskingPriceExpCcn` | `quotations/inquire-asking-price-exp-ccn` | FHKST01010200 |
| `Domestic.InquireCcnl` | `quotations/inquire-ccnl` | FHKST01010300 |
| `Domestic.InquireDailyPrice` | `quotations/inquire-daily-price` | FHKST01010400 |
| `Domestic.InquireNearNewHighlow` | `ranking/near-new-highlow` | FHPST01870000 |
| `Domestic.InquireOvertimePrice` | `quotations/inquire-overtime-price` | FHPST02300000 |
| `Domestic.InquireOvertimeAskingPrice` | `quotations/inquire-overtime-asking-price` | FHPST02300400 |
| `Domestic.InquireOvertimeVolume` | `ranking/overtime-volume` | FHPST02350000 |
| `Domestic.InquireOvertimeFluctuation` | `ranking/overtime-fluctuation` | FHPST02340000 |

## Design

- **호출 스타일**: `client.Domestic.<Method>(ctx, ...)` 1단계 그룹화 (go-github / stripe-go 패턴)
- **응답**: typed struct, 한투 API 의 한글 약어 필드는 JSON 태그로 매핑하고 영문 필드명으로 노출
- **에러**: `error` 반환. `error.Error()` 메시지에 `msg_cd` / `msg1` 가 포함됩니다 (예: `"kis: API error [EGW00201] 초당 거래건수를 초과하였습니다."`). typed error 분기는 추후 사용자 demand 시 재도입 예정.
- **자동 처리**: 토큰 갱신, rate limit (token bucket, 기본 15 req/sec), 429/5xx 재시도
- **HTTP**: 내부적으로 [resty](https://github.com/go-resty/resty) 사용 (사용자는 표준 `*http.Client` 만 알면 됨)
- **금융 정밀도**: 가격 필드는 [shopspring/decimal](https://github.com/shopspring/decimal)

상세 설계: [Phase 0 design spec](docs/superpowers/specs/2026-05-03-korea-investment-go-migration-design.md)

## Scope

- ✅ 국내주식 (시세, 차트, 순위, 재무, 투자자 동향, IPO/예탁원, 심볼)
- ✅ 해외주식 (시세, 차트, 순위)
- ❌ 선물옵션 — 영구 제외
- ❌ 장내채권 — 영구 제외
- ❌ 실시간 WebSocket — 추후 별도 spec
- ❌ 주식 주문/잔고/예약주문 — 본 spec 에서 다루지 않음

## License

MIT — 기존 Python 라이브러리와 동일.
