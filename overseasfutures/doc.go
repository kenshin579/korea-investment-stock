// Package overseasfutures 는 KIS 해외선물옵션 도메인.
//
// 한투 docs: docs/api/해외선물옵션/*.md
// base path: /uapi/overseas-futureoption/v1/
//
// Phase 11.5 — 해외선물 시세/조회 10 endpoint (모의 미지원):
//
//	HHDDB95030000  해외선물 미결제추이              InvestorUnpdTrend
//	HHDFC55020400  해외선물 분봉                    InquireTimeFuturechartprice
//	HHDFC55200000  해외선물 상품기본정보 (32 bulk)  SearchContractDetail
//	HHDFC55020300  해외선물 체결추이 (월간)         MonthlyCcnl
//	HHDFC55020100  해외선물 체결추이 (일간)         DailyCcnl
//	HHDFC55020000  해외선물 체결추이 (주간)         WeeklyCcnl
//	HHDFC55020200  해외선물 체결추이 (틱)           TickCcnl
//	HHDFC86000000  해외선물 호가                    InquireAskingPrice
//	HHDFC55010100  해외선물 종목상세                StockDetail
//	HHDFC55010000  해외선물 종목현재가              InquirePrice
//
// 사용자는 root kis.Client 의 OverseasFutures 필드로 접근.
package overseasfutures
