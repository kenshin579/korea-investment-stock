// Package overseas 는 한국투자증권 OpenAPI 의 해외주식 카테고리 메서드.
//
// Phase 1.5 메서드 (6):
//
//   - InquirePriceDetail        — 해외주식 현재가상세 (HHDFS76200200)
//   - SearchInfo                — 해외주식 상품기본정보 (CTPF1702R)
//   - InquireDailyPrice         — 해외주식 기간별시세 (HHDFS76240000) — 단일 종목 11 거래소
//   - InquireDailyChartPrice    — 해외주식 종목/지수/환율 기간별시세 (FHKST03030100) — 미국은 지수 한정
//   - InquireUpdownRate         — 해외주식 상승율/하락율 (HHDFS76290000)
//   - FetchOverseasSymbols      — 11 거래소 통합 마스터 (KIS 공개 다운로드)
//
// Phase 2.3 메서드 (6):
//
//   - InquireMarketCap   — 해외주식 시가총액순위 (HHDFS76350100)
//   - InquireTradeVol    — 해외주식 거래량순위 (HHDFS76310010)
//   - InquireTradePbmn   — 해외주식 거래대금순위 (HHDFS76320010)
//   - InquireVolumeSurge — 해외주식 거래량급증 (HHDFS76270000)
//   - InquireVolumePower — 해외주식 매수체결강도상위 (HHDFS76280000)
//   - InquireNewHighlow  — 해외주식 신고/신저가 (HHDFS76300000)
//
// 사용자는 root kis.Client 의 Overseas 필드로 접근.
package overseas
