// Package domestic 은 한국투자증권 OpenAPI 의 국내주식 카테고리 메서드.
//
// Phase 1.2 메서드 (7):
//
//   - InquirePrice                 — 주식현재가 시세 (FHKST01010100)
//   - SearchInfo                   — 상품기본조회 (CTPF1604R)
//   - SearchStockInfo              — 주식기본조회 (CTPF1002R)
//   - InquireDailyItemChartPrice   — 국내주식기간별시세 일/주/월/년 (FHKST03010100)
//   - InquireTimeItemChartPrice    — 주식당일분봉조회 (FHKST03010200)
//   - FetchKospiSymbols            — KRX KOSPI 마스터 (한투 API 가 아닌 KRX 공개 다운로드)
//   - FetchKosdaqSymbols           — KRX KOSDAQ 마스터
//
// Phase 1.3 메서드 (9):
//
//   - InquireVolumeRank            — 거래량순위 (FHPST01710000)
//   - InquireFluctuation           — 등락률 순위 (FHPST01700000)
//   - InquireMarketCap             — 시가총액 상위 (FHPST01740000)
//   - InquireDividendRate          — 배당률 상위 (HHKDB13470100)
//   - InquireFinancialRatio        — 재무비율 (FHKST66430300)
//   - InquireIncomeStatement       — 손익계산서 (FHKST66430200)
//   - InquireBalanceSheet          — 대차대조표 (FHKST66430100)
//   - InquireProfitRatio           — 수익성비율 (FHKST66430400)
//   - InquireGrowthRatio           — 성장성비율 (FHKST66430800)
//
// 사용자는 root kis.Client 의 Domestic 필드로 접근.
package domestic
