package websocket

import "github.com/shopspring/decimal"

// ---------------------------------------------------------------------------
// Phase 11.7 — 해외선물옵션 실시간 2 EP
//
// 선물/옵션 통합 EP — 단일 TR_ID 로 선물+옵션 모두 처리. 그릭스 없음.
// WebSocket 도메인: ws://ops.koreainvestment.com:21000
// 전체 모의투자 미지원.
//
// EP:
//   OverseasFuturesTradeEvent  (HDFFF020, 25 fields)
//   OverseasFuturesAskEvent    (HDFFF010, 35 fields — BID/ASK 교차 배열)
// ---------------------------------------------------------------------------

// OverseasFuturesTradeEvent 는 HDFFF020 해외선물옵션 실시간체결가 이벤트 (25 fields).
// 선물/옵션 통합 EP. 그릭스 없음.
// tr_key: 해외선물옵션 종목코드 6자리 (예: GCM24, 6AM24)
type OverseasFuturesTradeEvent struct {
	Symbol         string          // 0  SERIES_CD 종목코드
	BsnsDate       string          // 1  BSNS_DATE 영업일자 (YYYYMMDD)
	MrktOpenDate   string          // 2  MRKT_OPEN_DATE 장개시일자 (YYYYMMDD)
	MrktOpenTime   string          // 3  MRKT_OPEN_TIME 장개시시각 (HHMMSS)
	MrktCloseDate  string          // 4  MRKT_CLOSE_DATE 장종료일자 (YYYYMMDD)
	MrktCloseTime  string          // 5  MRKT_CLOSE_TIME 장종료시각 (HHMMSS)
	PrevPrice      decimal.Decimal // 6  PREV_PRICE 전일종가
	RecvDate       string          // 7  RECV_DATE 수신일자 (YYYYMMDD)
	RecvTime       string          // 8  RECV_TIME 수신시각 (HHMMSS)
	ActiveFlag     string          // 9  ACTIVE_FLAG 본장_전산장구분 (1자리)
	LastPrice      decimal.Decimal // 10 LAST_PRICE 체결가격
	LastQntt       int64           // 11 LAST_QNTT 체결수량
	PrevDiffPrice  decimal.Decimal // 12 PREV_DIFF_PRICE 전일대비가
	PrevDiffRate   float64         // 13 PREV_DIFF_RATE 등락률
	OpenPrice      decimal.Decimal // 14 OPEN_PRICE 시가
	HighPrice      decimal.Decimal // 15 HIGH_PRICE 고가
	LowPrice       decimal.Decimal // 16 LOW_PRICE 저가
	Vol            int64           // 17 VOL 누적거래량
	PrevSign       string          // 18 PREV_SIGN 전일대비부호 (1자리)
	QuotSign       string          // 19 QUOTSIGN 체결구분 (1자리; 2:매수체결, 5:매도체결)
	RecvTime2      string          // 20 RECV_TIME2 수신시각2 만분의일초 (4자리)
	PsttlPrice     decimal.Decimal // 21 PSTTL_PRICE 전일정산가
	PsttlSign      string          // 22 PSTTL_SIGN 전일정산가대비 (1자리)
	PsttlDiffPrice decimal.Decimal // 23 PSTTL_DIFF_PRICE 전일정산가대비가격
	PsttlDiffRate  float64         // 24 PSTTL_DIFF_RATE 전일정산가대비율

	Raw []string // caret 분리 원본 (escape hatch)
}

// OverseasFuturesAskEvent 는 HDFFF010 해외선물옵션 실시간호가 이벤트 (35 fields).
// 5단계 호가. BID/ASK 교차 구조 (BID_QNTT/BID_NUM/BID_PRICE → ASK_QNTT/ASK_NUM/ASK_PRICE 그룹).
// 국내선물옵션과 달리 건수(CSNU) 대신 번호(NUM) 사용. 총잔량 합계 필드 없음.
// tr_key: 해외선물옵션 종목코드 6자리
type OverseasFuturesAskEvent struct {
	Symbol    string          // 0  SERIES_CD 종목코드
	RecvDate  string          // 1  RECV_DATE 수신일자 (YYYYMMDD)
	RecvTime  string          // 2  RECV_TIME 수신시각 (12자리, 나노초 포함)
	PrevPrice decimal.Decimal // 3  PREV_PRICE 전일종가

	// 5단계 호가 — BID/ASK 교차 배열
	// 필드 순서: BID_QNTT_i → BID_NUM_i → BID_PRICE_i → ASK_QNTT_i → ASK_NUM_i → ASK_PRICE_i (i=1..5)
	BidQntt  [5]int64           // 4,10,16,22,28 BID_QNTT_1..5 매수수량
	BidNum   [5]string          // 5,11,17,23,29 BID_NUM_1..5  매수번호 (10자리)
	BidPrice [5]decimal.Decimal // 6,12,18,24,30 BID_PRICE_1..5 매수호가
	AskQntt  [5]int64           // 7,13,19,25,31 ASK_QNTT_1..5 매도수량
	AskNum   [5]string          // 8,14,20,26,32 ASK_NUM_1..5  매도번호 (10자리)
	AskPrice [5]decimal.Decimal // 9,15,21,27,33 ASK_PRICE_1..5 매도호가

	SttlPrice decimal.Decimal // 34 STTL_PRICE 전일정산가

	Raw []string // caret 분리 원본 (escape hatch)
}
