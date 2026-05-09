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

## Available Methods (Phase 1.2 ~ 5)

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
| `Overseas.InquireMarketCap` | `ranking/market-cap` | HHDFS76350100 |
| `Overseas.InquireTradeVol` | `ranking/trade-vol` | HHDFS76310010 |
| `Overseas.InquireTradePbmn` | `ranking/trade-pbmn` | HHDFS76320010 |
| `Overseas.InquireVolumeSurge` | `ranking/volume-surge` | HHDFS76270000 |
| `Overseas.InquireVolumePower` | `ranking/volume-power` | HHDFS76280000 |
| `Overseas.InquireNewHighlow` | `ranking/new-highlow` | HHDFS76300000 |
| `Domestic.InquireKsdDividend` | `ksdinfo/dividend` | HHKDB669102C0 |
| `Domestic.InquireKsdBonusIssue` | `ksdinfo/bonus-issue` | HHKDB669101C0 |
| `Domestic.InquireKsdPaidinCapin` | `ksdinfo/paidin-capin` | HHKDB669100C0 |
| `Domestic.InquireKsdSharehldMeet` | `ksdinfo/sharehld-meet` | HHKDB669111C0 |
| `Domestic.InquireKsdMergerSplit` | `ksdinfo/merger-split` | HHKDB669104C0 |
| `Domestic.InquireKsdRevSplit` | `ksdinfo/rev-split` | HHKDB669105C0 |
| `Domestic.InquireKsdForfeit` | `ksdinfo/forfeit` | HHKDB669109C0 |
| `Domestic.InquireKsdMandDeposit` | `ksdinfo/mand-deposit` | HHKDB669110C0 |
| `Domestic.InquireKsdCapDcrs` | `ksdinfo/cap-dcrs` | HHKDB669106C0 |
| `Domestic.InquireKsdPurreq` | `ksdinfo/purreq` | HHKDB669103C0 |
| `Domestic.InquireKsdListInfo` | `ksdinfo/list-info` | HHKDB669107C0 |
| `Domestic.InquireInvestorTrendEstimate` | `quotations/investor-trend-estimate` | HHPTJ04160200 |
| `Domestic.InquireForeignInstitutionTotal` | `quotations/foreign-institution-total` | FHPTJ04400000 |
| `Domestic.InquireProgramTradeByStockDaily` | `quotations/program-trade-by-stock-daily` | FHPPG04650201 |
| `Domestic.InquireProgramTradeByStock` | `quotations/program-trade-by-stock` | FHPPG04650101 |
| `Domestic.InquireCompProgramTradeToday` | `quotations/comp-program-trade-today` | FHPPG04600101 |
| `Domestic.InquireCompProgramTradeDaily` | `quotations/comp-program-trade-daily` | FHPPG04600001 |
| `Domestic.InquireInvestorProgramTradeToday` | `quotations/investor-program-trade-today` | HHPPG046600C1 |
| `Overseas.InquireNewsTitle` | `overseas-price/v1/quotations/news-title` | HHPSTH60100C1 |
| `Overseas.InquireBrknewsTitle` | `overseas-price/v1/quotations/brknews-title` | FHKST01011801 |
| `Overseas.InquireRightsByIce` | `overseas-price/v1/quotations/rights-by-ice` | HHDFS78330900 |
| `Overseas.InquirePeriodRights` | `overseas-price/v1/quotations/period-rights` | CTRGT011R |
| `Domestic.InquireIndexDailyPrice` | `quotations/inquire-index-daily-price` | FHPUP02120000 |
| `Domestic.InquireIndexTimeprice` | `quotations/inquire-index-timeprice` | FHPUP02110200 |
| `Domestic.InquireIndexTickprice` | `quotations/inquire-index-tickprice` | FHPUP02110100 |
| `Domestic.InquireDailyIndexchartprice` | `quotations/inquire-daily-indexchartprice` | FHKUP03500100 |
| `Domestic.InquireTimeIndexchartprice` | `quotations/inquire-time-indexchartprice` | FHKUP03500200 |
| `Domestic.ExpTotalIndex` | `quotations/exp-total-index` | FHKUP11750000 |
| `Domestic.ExpIndexTrend` | `quotations/exp-index-trend` | FHPST01840000 |
| `Domestic.InquireInvestOpinion` | `quotations/invest-opinion` | FHKST663300C0 |
| `Domestic.InquireInvestOpbysec` | `quotations/invest-opbysec` | FHKST663400C0 |
| `Domestic.InquireEstimatePerform` | `quotations/estimate-perform` | HHKST668300C0 |
| `Domestic.InquireVolumePower` | `ranking/volume-power` | FHPST01680000 |
| `Domestic.InquireBulkTransNum` | `ranking/bulk-trans-num` | FHKST190900C0 |
| `Domestic.InquireTradprtByamt` | `quotations/tradprt-byamt` | FHKST111900C0 |
| `Domestic.InquireHtsTopView` | `ranking/hts-top-view` | HHMCM000100C0 |
| `Domestic.InquirePbarTraRatio` | `quotations/pbar-tratio` | FHPST01130000 |
| `Domestic.InquireExpPriceTrend` | `quotations/exp-price-trend` | FHPST01810000 |
| `Domestic.InquireExpTransUpdown` | `ranking/exp-trans-updown` | FHPST01820000 |

### Domestic ranking/흐름 — Phase 4.3 (v1.14.0)

| Method | 한투 path | TR_ID |
|--------|----------|-------|
| `Domestic.InquireShortSale` | `ranking/short-sale` | FHPST04820000 |
| `Domestic.InquireDailyShortSale` | `quotations/daily-short-sale` | FHPST04830000 |
| `Domestic.InquireCreditBalance` | `ranking/credit-balance` | FHKST17010000 |
| `Domestic.InquireDailyCreditBalance` | `quotations/daily-credit-balance` | FHPST04760000 |
| `Domestic.InquireLendableByCompany` | `quotations/lendable-by-company` | CTSC2702R |
| `Domestic.InquireQuoteBalance` | `ranking/quote-balance` | FHPST01720000 |
| `Domestic.InquireAfterHourBalance` | `ranking/after-hour-balance` | FHPST01760000 |
| `Domestic.InquireOvertimeExpTransFluct` | `ranking/overtime-exp-trans-fluct` | FHKST11860000 |
| `Domestic.InquireMarketValue` | `ranking/market-value` | FHPST01790000 |
| `Domestic.InquireDisparity` | `ranking/disparity` | FHPST01780000 |
| `Domestic.InquirePreferDisparateRatio` | `ranking/prefer-disparate-ratio` | FHPST01770000 |
| `Domestic.InquireProfitAssetIndex` | `ranking/profit-asset-index` | FHPST01730000 |
| `Domestic.InquireMktfunds` | `quotations/mktfunds` | FHKST649100C0 |

### Domestic 시장운영/특수상태 — Phase 4.2 (v1.13.0)

| Method | 한투 path | TR_ID |
|--------|----------|-------|
| `Domestic.InquireExpClosingPrice` | `quotations/exp-closing-price` | FHKST117300C0 |
| `Domestic.InquireChkHoliday` | `quotations/chk-holiday` | CTCA0903R |
| `Domestic.InquireViStatus` | `quotations/inquire-vi-status` | FHPST01390000 |
| `Domestic.InquireCaptureUplowprice` | `quotations/capture-uplowprice` | FHKST130000C0 |

### ETF/NAV/관심종목 — Phase 5 (v1.15.0)

| Method | 한투 path | TR_ID |
|--------|----------|-------|
| `Domestic.InquireEtfPrice` | `etfetn/inquire-price` | FHPST02400000 |
| `Domestic.InquireComponentStockPrice` | `etfetn/inquire-component-stock-price` | FHKST121600C0 |
| `Domestic.InquireNavComparisonTimeTrend` | `etfetn/nav-comparison-time-trend` | FHPST02440100 |
| `Domestic.InquireNavComparisonDailyTrend` | `etfetn/nav-comparison-daily-trend` | FHPST02440200 |
| `Domestic.InquireNavComparisonTrend` | `etfetn/nav-comparison-trend` | FHPST02440000 |
| `Domestic.InquireIntstockMultprice` | `quotations/intstock-multprice` | FHKST11300006 |
| `Domestic.InquireIntstockStocklistByGroup` | `quotations/intstock-stocklist-by-group` | HHKCM113004C6 |
| `Domestic.InquireIntstockGrouplist` | `quotations/intstock-grouplist` | HHKCM113004C7 |
| `Domestic.InquireTopInterestStock` | `ranking/top-interest-stock` | FHPST01800000 |

### 재무 추가 — Phase 6 (v1.16.0)

| Method | 한투 path | TR_ID |
|--------|----------|-------|
| `Domestic.InquireOtherMajorRatios` | `finance/other-major-ratios` | FHKST66430500 |
| `Domestic.InquireFinanceRatioRanking` | `ranking/finance-ratio` | FHPST01750000 |

### 헬퍼 — Phase 7 (v1.17.0)

| Method | 한투 path | TR_ID |
|--------|----------|-------|
| `Domestic.InquireMarketTime` | `quotations/market-time` | HHMCM000002C0 |
| `Domestic.InquireCompInterest` | `quotations/comp-interest` | FHPST07020000 |
| `Domestic.InquireTradedByCompany` | `ranking/traded-by-company` | FHPST01860000 |
| `Domestic.InquireCreditByCompany` | `quotations/credit-by-company` | FHPST04770000 |

### WebSocket — Phase 8 (v1.18.0)

| Method | TR_ID | 설명 |
|--------|-------|------|
| `WS.SubscribeKrxTrade` / `OnKrxTrade` | H0STCNT0 | 국내주식 실시간체결가 (KRX) |
| `WS.SubscribeKrxAsk` / `OnKrxAsk` | H0STASP0 | 국내주식 실시간호가 (KRX) |
| `WS.SubscribeKrxExpectTrade` / `OnKrxExpectTrade` | H0STANC0 | 국내주식 실시간예상체결 (KRX) |
| `WS.SubscribeKrxOvernightTrade` / `OnKrxOvernightTrade` | H0STOUP0 | 국내주식 시간외 실시간체결가 (KRX) |
| `WS.SubscribeKrxOvernightExpect` / `OnKrxOvernightExpect` | H0STOAC0 | 국내주식 시간외 실시간예상체결 (KRX) |

자동 재연결 + 구독 자동 복원 (exp backoff, max 10 attempts). 사용 예: `examples/ws_krx_basic/`.

### WebSocket — Phase 9 (v1.19.0) — NXT/통합 변형

NXT (대체거래소) 와 통합 (KRX+NXT) 시장의 실시간 EP. NXT 와 통합은 schema 동일 (5 base struct + 10 type alias 패턴). 모든 EP 모의 미지원.

| Method | TR_ID | 설명 |
|--------|-------|------|
| `WS.SubscribeNxtTrade` / `OnNxtTrade` | H0NXCNT0 | 실시간체결가 NXT |
| `WS.SubscribeUnifiedTrade` / `OnUnifiedTrade` | H0UNCNT0 | 실시간체결가 통합 |
| `WS.SubscribeNxtAsk` / `OnNxtAsk` | H0NXASP0 | 실시간호가 NXT (KMID/NMID 중간가 6 fields 포함) |
| `WS.SubscribeUnifiedAsk` / `OnUnifiedAsk` | H0UNASP0 | 실시간호가 통합 (KMID/NMID 중간가 6 fields 포함) |
| `WS.SubscribeNxtExpectTrade` / `OnNxtExpectTrade` | H0NXANC0 | 실시간예상체결 NXT (VI_STND_PRC 추가) |
| `WS.SubscribeUnifiedExpectTrade` / `OnUnifiedExpectTrade` | H0UNANC0 | 실시간예상체결 통합 (VI_STND_PRC 추가) |
| `WS.SubscribeNxtProgramTrade` / `OnNxtProgramTrade` | H0NXPGM0 | 실시간프로그램매매 NXT (신규 EP) |
| `WS.SubscribeUnifiedProgramTrade` / `OnUnifiedProgramTrade` | H0UNPGM0 | 실시간프로그램매매 통합 (신규 EP) |
| `WS.SubscribeNxtMember` / `OnNxtMember` | H0NXMBC0 | 실시간회원사 NXT (5단계 매도/매수 + 외국계 + 영문) |
| `WS.SubscribeUnifiedMember` / `OnUnifiedMember` | H0UNMBC0 | 실시간회원사 통합 (5단계 매도/매수 + 외국계 + 영문) |

### WebSocket — Phase 10 (v1.20.0) — 해외주식 시세

해외주식 실시간 시세 2 EP. tr_key 형식: `D`/`R` + 시장구분(NAS/NYS/AMS/TSE/HKS/SHS/SZS/HSX/HNX 등) + 종목코드 (예: `DNASAAPL`). 모든 EP 모의 미지원.

| Method | TR_ID | 설명 |
|--------|-------|------|
| `WS.SubscribeOverseasTrade` / `OnOverseasTrade` | HDFSCNT0 | 해외주식 실시간지연체결가 (미국 0분지연 / 아시아 15분지연) |
| `WS.SubscribeOverseasAsk` / `OnOverseasAsk` | HDFSASP0 | 해외주식 실시간호가 (1호가만, 미국 무료 / 아시아 유료) |

### WebSocket — Phase 11.2 (v1.22.0) — 국내선물옵션 실시간

국내선물옵션 실시간 11 EP — KRX 야간 5 + 주식 선물옵션 6. 모든 EP 모의 미지원. 11 EP 모두 distinct schema (선물 vs 옵션 + KRX 야간 vs 주식 모두 다름). 옵션 EP 들은 그릭스 (DELTA / GAMA / VEGA / THETA / RHO) + IV / HV 포함.

| Method | TR_ID | 설명 |
|--------|-------|------|
| `WS.SubscribeKrxNightFuturesTrade` / `OnKrxNightFuturesTrade` | H0MFCNT0 | KRX 야간 선물 실시간 체결 |
| `WS.SubscribeKrxNightFuturesAsk` / `OnKrxNightFuturesAsk` | H0MFASP0 | KRX 야간 선물 실시간 호가 (5단계) |
| `WS.SubscribeKrxNightOptionTrade` / `OnKrxNightOptionTrade` | H0EUCNT0 | KRX 야간 옵션 실시간 체결가 (그릭스) |
| `WS.SubscribeKrxNightOptionAsk` / `OnKrxNightOptionAsk` | H0EUASP0 | KRX 야간 옵션 실시간 호가 (5단계) |
| `WS.SubscribeKrxNightOptionExpectTrade` / `OnKrxNightOptionExpectTrade` | H0EUANC0 | KRX 야간 옵션 실시간 예상체결 |
| `WS.SubscribeStockFuturesTrade` / `OnStockFuturesTrade` | H0ZFCNT0 | 주식 선물 실시간 체결가 |
| `WS.SubscribeStockFuturesAsk` / `OnStockFuturesAsk` | H0ZFASP0 | 주식 선물 실시간 호가 (10단계) |
| `WS.SubscribeStockFuturesExpectTrade` / `OnStockFuturesExpectTrade` | H0ZFANC0 | 주식 선물 실시간 예상체결 |
| `WS.SubscribeStockOptionTrade` / `OnStockOptionTrade` | H0ZOCNT0 | 주식 옵션 실시간 체결가 (그릭스) |
| `WS.SubscribeStockOptionAsk` / `OnStockOptionAsk` | H0ZOASP0 | 주식 옵션 실시간 호가 (10단계) |
| `WS.SubscribeStockOptionExpectTrade` / `OnStockOptionExpectTrade` | H0ZOANC0 | 주식 옵션 실시간 예상체결 |

### Bonds (장내채권) — Phase 3.1

| Go 메서드 | path | TR_ID |
|---|---|---|
| `Bonds.SearchBondInfo` | `quotations/search-bond-info` | CTPF1114R |
| `Bonds.InquireIssueInfo` | `quotations/issue-info` | CTPF1101R |
| `Bonds.InquirePrice` | `quotations/inquire-price` | FHKBJ773400C0 |
| `Bonds.InquireCcnl` | `quotations/inquire-ccnl` | FHKBJ773403C0 |
| `Bonds.InquireAskingPrice` | `quotations/inquire-asking-price` | FHKBJ773401C0 |
| `Bonds.InquireDailyPrice` | `quotations/inquire-daily-price` | FHKBJ773404C0 |
| `Bonds.InquireDailyItemchartprice` | `quotations/inquire-daily-itemchartprice` | FHKBJ773701C0 |
| `Bonds.InquireAvgUnit` | `quotations/avg-unit` | CTPF2005R |

### Futures (국내선물옵션) — Phase 11.1

종목코드 9자리 alphanumeric (예: `101W3000` 선물, `201X3300` 옵션). MarketCode 인자 caller 가 입력 (`F`/`O`/`JF`/`JO`/`CF` 등).

| Go 메서드 | path | TR_ID | 모의 |
|---|---|---|---|
| `Futures.InquirePrice` | `quotations/inquire-price` | FHMIF10000000 | 지원 |
| `Futures.InquireAskingPrice` | `quotations/inquire-asking-price` | FHMIF10010000 | 지원 |
| `Futures.InquireTimeFuopchartprice` | `quotations/inquire-time-fuopchartprice` | FHKIF03020200 | 미지원 |
| `Futures.ExpPriceTrend` | `quotations/exp-price-trend` | FHPIF05110100 | 미지원 |
| `Futures.InquireDailyFuopchartprice` | `quotations/inquire-daily-fuopchartprice` | FHKIF03020100 | 지원 |
| `Futures.DisplayBoardTop` | `quotations/display-board-top` | FHPIF05030000 | 미지원 |
| `Futures.DisplayBoardFutures` | `quotations/display-board-futures` | FHPIF05030200 | 미지원 |
| `Futures.DisplayBoardOptionList` | `quotations/display-board-option-list` | FHPIO056104C0 | 미지원 |
| `Futures.DisplayBoardCallput` | `quotations/display-board-callput` | FHPIF05030100 | 미지원 |

## Design

- **호출 스타일**: `client.Domestic.<Method>(ctx, ...)` 1단계 그룹화 (go-github / stripe-go 패턴)
- **응답**: typed struct, 한투 API 의 한글 약어 필드는 JSON 태그로 매핑하고 영문 필드명으로 노출
- **에러**: `error` 반환. `error.Error()` 메시지에 `msg_cd` / `msg1` 가 포함됩니다 (예: `"kis: API error [EGW00201] 초당 거래건수를 초과하였습니다."`). typed error 분기는 추후 사용자 demand 시 재도입 예정.
- **자동 처리**: 토큰 갱신, rate limit (token bucket, 기본 15 req/sec), 429/5xx 재시도
- **HTTP**: 내부적으로 [resty](https://github.com/go-resty/resty) 사용 (사용자는 표준 `*http.Client` 만 알면 됨)
- **금융 정밀도**: 가격 필드는 [shopspring/decimal](https://github.com/shopspring/decimal)

상세 설계: [Phase 0 design spec](docs/superpowers/specs/2026-05-03-korea-investment-go-migration-design.md)

## Scope

- ✅ 국내주식 (시세, 차트, 순위, 재무, 투자자 동향, IPO/예탁원, ETF/NAV, 관심종목, 심볼)
- ✅ 해외주식 (시세, 차트, 순위)
- ✅ 장내채권 (시세, 발행정보, 호가, 기간별, 평균단가 — Phase 3.1)
- ✅ 국내선물옵션 시세/조회 9 EP (Phase 11.1; v1.21.0). 실시간/Trading 은 Phase 11.2+.
- ✅ 실시간 WebSocket — KRX 5 EP (Phase 8; v1.18.0) + NXT/통합 10 EP (Phase 9; v1.19.0) + 해외 시세 2 EP (Phase 10; v1.20.0) + 국내선물옵션 11 EP (Phase 11.2; v1.22.0)
- ❌ 주식 주문/잔고/예약주문 — 본 spec 에서 다루지 않음

## License

MIT — 기존 Python 라이브러리와 동일.
