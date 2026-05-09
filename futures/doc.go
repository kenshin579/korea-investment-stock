// Package futures 는 KIS 국내선물옵션 도메인.
//
// 한투 docs: docs/api/국내선물옵션/*.md
//
// Phase 11.1 — 시세/조회 9 endpoint:
//
//	FHMIF10000000  선물옵션 시세                  InquirePrice
//	FHMIF10010000  선물옵션 시세호가              InquireAskingPrice
//	FHKIF03020200  선물옵션 분봉                  InquireTimeFuopchartprice
//	FHPIF05110100  선물옵션 일중예상체결추이      ExpPriceTrend
//	FHKIF03020100  선물옵션 기간별시세 (일/주/월/년) InquireDailyFuopchartprice
//	FHPIF05030000  국내선물 전광판 top (기초자산) DisplayBoardTop
//	FHPIF05030200  옵션 전광판 선물               DisplayBoardFutures
//	FHPIO056104C0  옵션 전광판 옵션월물리스트     DisplayBoardOptionList
//	FHPIF05030100  옵션 전광판 콜풋               DisplayBoardCallput
//
// EP4 (InquireCcnlBstime, CTFO5139R) + EP7 (InquireDailyAmountFee, CTFO6119R) 는
// 계좌 정보 (CANO/ACNT_PRDT_CD) 필요로 Phase 11.4 (Trading) 에서 구현.
//
// 사용자는 root kis.Client 의 Futures 필드로 접근.
package futures
