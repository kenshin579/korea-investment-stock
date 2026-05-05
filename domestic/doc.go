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
// Phase 1.4 메서드 (6):
//
//   - InquireInvestorTradeByStockDaily — 종목별 투자자매매동향 일별 (FHPTJ04160001)
//   - InquireInvestorDailyByMarket    — 시장별 투자자매매동향 일별 (FHPTJ04040000)
//   - InquireInvestorTimeByMarket     — 시장별 투자자매매동향 시세 (FHPTJ04030000)
//   - InquireIndexPrice               — 국내업종 현재지수 (FHPUP02100000)
//   - InquireIndexCategoryPrice       — 국내업종 구분별 전체시세 (FHPUP02140000)
//   - InquirePubOffer                 — 예탁원정보 공모주청약일정 (HHKDB669108C0)
//
// Phase 2.1 메서드 (3):
//
//   - InquireAskingPriceExpCcn  — 주식현재가 호가/예상체결 (FHKST01010200)
//   - InquireCcnl               — 주식현재가 체결 (FHKST01010300)
//   - InquireDailyPrice         — 주식현재가 일자별 (FHKST01010400)
//
// Phase 2.2 메서드 (5):
//
//   - InquireNearNewHighlow      — 국내주식 신고/신저근접종목 상위 (FHPST01870000)
//   - InquireOvertimePrice       — 국내주식 시간외현재가 (FHPST02300000)
//   - InquireOvertimeAskingPrice — 국내주식 시간외호가 (FHPST02300400)
//   - InquireOvertimeVolume      — 국내주식 시간외거래량순위 (FHPST02350000)
//   - InquireOvertimeFluctuation — 국내주식 시간외등락율순위 (FHPST02340000)
//
// 사용자는 root kis.Client 의 Domestic 필드로 접근.
package domestic
