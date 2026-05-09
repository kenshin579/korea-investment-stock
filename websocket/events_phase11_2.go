package websocket

import "github.com/shopspring/decimal"

// ---------------------------------------------------------------------------
// Phase 11.2 — 국내선물옵션 실시간 11 EP
//
// 11 EP 모두 Distinct (alias 없음). WebSocket 도메인: ws://ops.koreainvestment.com:21000
// 전체 모의투자 미지원.
// ---------------------------------------------------------------------------

// KrxNightFuturesTradeEvent 는 H0MFCNT0 KRX야간선물 실시간종목체결 이벤트 (49 fields).
// tr_key: 야간선물 종목코드 12자리 (예: 101W09000000)
type KrxNightFuturesTradeEvent struct {
	Symbol                 string          // 0  FUTS_SHRN_ISCD 선물단축종목코드 (9자리)
	Time                   string          // 1  BSOP_HOUR 영업시간 (HHMMSS)
	PrevDiff               decimal.Decimal // 2  FUTS_PRDY_VRSS 선물전일대비
	PrevDiffSign           string          // 3  PRDY_VRSS_SIGN 전일대비부호
	PrevChangeRate         float64         // 4  FUTS_PRDY_CTRT 선물전일대비율
	Price                  decimal.Decimal // 5  FUTS_PRPR 선물현재가
	Open                   decimal.Decimal // 6  FUTS_OPRC 선물시가2
	High                   decimal.Decimal // 7  FUTS_HGPR 선물최고가
	Low                    decimal.Decimal // 8  FUTS_LWPR 선물최저가
	LastTradeVolume        int64           // 9  LAST_CNQN 최종거래량
	AccumVolume            int64           // 10 ACML_VOL 누적거래량
	AccumValue             int64           // 11 ACML_TR_PBMN 누적거래대금
	TheoreticalPrice       decimal.Decimal // 12 HTS_THPR HTS이론가
	MarketBasis            decimal.Decimal // 13 MRKT_BASIS 시장베이시스
	DeviationRate          float64         // 14 DPRT 괴리율
	NearMonthPrice         decimal.Decimal // 15 NMSC_FCTN_STPL_PRC 근월물약정가
	FarMonthPrice          decimal.Decimal // 16 FMSC_FCTN_STPL_PRC 원월물약정가
	SpreadPrice            decimal.Decimal // 17 SPEAD_PRC 스프레드1
	OpenInterestQty        int64           // 18 HTS_OTST_STPL_QTY HTS미결제약정수량
	OpenInterestChange     int64           // 19 OTST_STPL_QTY_ICDC 미결제약정수량증감
	OpenTime               string          // 20 OPRC_HOUR 시가시간 (HHMMSS)
	OpenDiffSign           string          // 21 OPRC_VRSS_PRPR_SIGN 시가2대비현재가부호
	OpenDiff               decimal.Decimal // 22 OPRC_VRSS_NMIX_PRPR 시가대비지수현재가
	HighTime               string          // 23 HGPR_HOUR 최고가시간 (HHMMSS)
	HighDiffSign           string          // 24 HGPR_VRSS_PRPR_SIGN 최고가대비현재가부호
	HighDiff               decimal.Decimal // 25 HGPR_VRSS_NMIX_PRPR 최고가대비지수현재가
	LowTime                string          // 26 LWPR_HOUR 최저가시간 (HHMMSS)
	LowDiffSign            string          // 27 LWPR_VRSS_PRPR_SIGN 최저가대비현재가부호
	LowDiff                decimal.Decimal // 28 LWPR_VRSS_NMIX_PRPR 최저가대비지수현재가
	BidRate                float64         // 29 SHNU_RATE 매수2비율
	TradeStrength          float64         // 30 CTTR 체결강도
	DeviationDegree        decimal.Decimal // 31 ESDG 괴리도
	OpenInterestPrevChange int64           // 32 OTST_STPL_RGBF_QTY_ICDC 미결제약정직전수량증감
	TheoreticalBasis       decimal.Decimal // 33 THPR_BASIS 이론베이시스
	Ask1                   decimal.Decimal // 34 FUTS_ASKP1 선물매도호가1
	Bid1                   decimal.Decimal // 35 FUTS_BIDP1 선물매수호가1
	Ask1Size               int64           // 36 ASKP_RSQN1 매도호가잔량1
	Bid1Size               int64           // 37 BIDP_RSQN1 매수호가잔량1
	AskCount               int64           // 38 SELN_CNTG_CSNU 매도체결건수
	BidCount               int64           // 39 SHNU_CNTG_CSNU 매수체결건수
	NetCount               int64           // 40 NTBY_CNTG_CSNU 순매수체결건수
	TotalAskVolume         int64           // 41 SELN_CNTG_SMTN 총매도수량
	TotalBidVolume         int64           // 42 SHNU_CNTG_SMTN 총매수수량
	TotalAskSize           int64           // 43 TOTAL_ASKP_RSQN 총매도호가잔량
	TotalBidSize           int64           // 44 TOTAL_BIDP_RSQN 총매수호가잔량
	PrevVolRate            float64         // 45 PRDY_VOL_VRSS_ACML_VOL_RATE 전일거래량대비등락율
	DynamicUpperLimit      decimal.Decimal // 46 DYNM_MXPR 실시간상한가 (Length=8)
	DynamicLowerLimit      decimal.Decimal // 47 DYNM_LLAM 실시간하한가 (Length=8)
	DynamicPriceLimitYN    string          // 48 DYNM_PRC_LIMT_YN 실시간가격제한구분 (1자리)

	Raw []string // caret 분리 원본 (escape hatch)
}

// KrxNightFuturesAskEvent 는 H0MFASP0 KRX야간선물 실시간호가 이벤트 (38 fields).
// 5단계 호가. tr_key: 야간선물 종목코드 12자리
type KrxNightFuturesAskEvent struct {
	Symbol string // 0  FUTS_SHRN_ISCD 선물단축종목코드 (9자리)
	Time   string // 1  BSOP_HOUR 영업시간 (HHMMSS)

	Ask     [5]decimal.Decimal // 2..6  FUTS_ASKP1..5 선물매도호가1~5
	Bid     [5]decimal.Decimal // 7..11 FUTS_BIDP1..5 선물매수호가1~5
	AskCsnu [5]int64           // 12..16 ASKP_CSNU1..5 매도호가건수1~5
	BidCsnu [5]int64           // 17..21 BIDP_CSNU1..5 매수호가건수1~5
	AskSize [5]int64           // 22..26 ASKP_RSQN1..5 매도호가잔량1~5
	BidSize [5]int64           // 27..31 BIDP_RSQN1..5 매수호가잔량1~5

	TotalAskCsnu    int64 // 32 TOTAL_ASKP_CSNU 총매도호가건수
	TotalBidCsnu    int64 // 33 TOTAL_BIDP_CSNU 총매수호가건수
	TotalAskSize    int64 // 34 TOTAL_ASKP_RSQN 총매도호가잔량
	TotalBidSize    int64 // 35 TOTAL_BIDP_RSQN 총매수호가잔량
	TotalAskSizeChg int64 // 36 TOTAL_ASKP_RSQN_ICDC 총매도호가잔량증감
	TotalBidSizeChg int64 // 37 TOTAL_BIDP_RSQN_ICDC 총매수호가잔량증감

	Raw []string // caret 분리 원본 (escape hatch)
}

// KrxNightOptionTradeEvent 는 H0EUCNT0 KRX야간옵션 실시간체결가 이벤트 (56 fields).
// 옵션 그릭스 포함, DYNM 3 fields 포함. tr_key: 야간옵션 종목코드 12자리
type KrxNightOptionTradeEvent struct {
	Symbol                 string          // 0  OPTN_SHRN_ISCD 옵션단축종목코드 (9자리)
	Time                   string          // 1  BSOP_HOUR 영업시간 (HHMMSS)
	Price                  decimal.Decimal // 2  OPTN_PRPR 옵션현재가
	PrevDiffSign           string          // 3  PRDY_VRSS_SIGN 전일대비부호
	PrevDiff               decimal.Decimal // 4  OPTN_PRDY_VRSS 옵션전일대비
	PrevChangeRate         float64         // 5  PRDY_CTRT 전일대비율
	Open                   decimal.Decimal // 6  OPTN_OPRC 옵션시가2
	High                   decimal.Decimal // 7  OPTN_HGPR 옵션최고가
	Low                    decimal.Decimal // 8  OPTN_LWPR 옵션최저가
	LastTradeVolume        int64           // 9  LAST_CNQN 최종거래량
	AccumVolume            int64           // 10 ACML_VOL 누적거래량
	AccumValue             int64           // 11 ACML_TR_PBMN 누적거래대금
	TheoreticalPrice       decimal.Decimal // 12 HTS_THPR HTS이론가
	OpenInterestQty        int64           // 13 HTS_OTST_STPL_QTY HTS미결제약정수량
	OpenInterestChange     int64           // 14 OTST_STPL_QTY_ICDC 미결제약정수량증감
	OpenTime               string          // 15 OPRC_HOUR 시가시간 (HHMMSS)
	OpenDiffSign           string          // 16 OPRC_VRSS_PRPR_SIGN 시가2대비현재가부호
	OpenDiff               decimal.Decimal // 17 OPRC_VRSS_NMIX_PRPR 시가대비지수현재가
	HighTime               string          // 18 HGPR_HOUR 최고가시간 (HHMMSS)
	HighDiffSign           string          // 19 HGPR_VRSS_PRPR_SIGN 최고가대비현재가부호
	HighDiff               decimal.Decimal // 20 HGPR_VRSS_NMIX_PRPR 최고가대비지수현재가
	LowTime                string          // 21 LWPR_HOUR 최저가시간 (HHMMSS)
	LowDiffSign            string          // 22 LWPR_VRSS_PRPR_SIGN 최저가대비현재가부호
	LowDiff                decimal.Decimal // 23 LWPR_VRSS_NMIX_PRPR 최저가대비지수현재가
	BidRate                float64         // 24 SHNU_RATE 매수2비율
	PremiumValue           decimal.Decimal // 25 PRMM_VAL 프리미엄값
	IntrinsicValue         decimal.Decimal // 26 INVL_VAL 내재가치값
	TimeValue              decimal.Decimal // 27 TMVL_VAL 시간가치값
	Delta                  float64         // 28 DELTA 델타
	Gamma                  float64         // 29 GAMA 감마
	Vega                   float64         // 30 VEGA 베가
	Theta                  float64         // 31 THETA 세타
	Rho                    float64         // 32 RHO 로우
	ImpliedVolatility      float64         // 33 HTS_INTS_VLTL HTS내재변동성
	DeviationDegree        decimal.Decimal // 34 ESDG 괴리도
	OpenInterestPrevChange int64           // 35 OTST_STPL_RGBF_QTY_ICDC 미결제약정직전수량증감
	TheoreticalBasis       decimal.Decimal // 36 THPR_BASIS 이론베이시스
	HistoricalVolatility   float64         // 37 UNAS_HIST_VLTL 역사적변동성
	TradeStrength          float64         // 38 CTTR 체결강도
	DeviationRate          float64         // 39 DPRT 괴리율
	MarketBasis            decimal.Decimal // 40 MRKT_BASIS 시장베이시스
	Ask1                   decimal.Decimal // 41 OPTN_ASKP1 옵션매도호가1
	Bid1                   decimal.Decimal // 42 OPTN_BIDP1 옵션매수호가1
	Ask1Size               int64           // 43 ASKP_RSQN1 매도호가잔량1
	Bid1Size               int64           // 44 BIDP_RSQN1 매수호가잔량1
	AskCount               int64           // 45 SELN_CNTG_CSNU 매도체결건수
	BidCount               int64           // 46 SHNU_CNTG_CSNU 매수체결건수
	NetCount               int64           // 47 NTBY_CNTG_CSNU 순매수체결건수
	TotalAskVolume         int64           // 48 SELN_CNTG_SMTN 총매도수량
	TotalBidVolume         int64           // 49 SHNU_CNTG_SMTN 총매수수량
	TotalAskSize           int64           // 50 TOTAL_ASKP_RSQN 총매도호가잔량
	TotalBidSize           int64           // 51 TOTAL_BIDP_RSQN 총매수호가잔량
	PrevVolRate            float64         // 52 PRDY_VOL_VRSS_ACML_VOL_RATE 전일거래량대비등락율
	DynamicUpperLimit      decimal.Decimal // 53 DYNM_MXPR 실시간상한가 (Length=8)
	DynamicPriceLimitYN    string          // 54 DYNM_PRC_LIMT_YN 실시간가격제한구분 (1자리) — docs: MXPR→PRC_LIMT_YN→LLAM 순서
	DynamicLowerLimit      decimal.Decimal // 55 DYNM_LLAM 실시간하한가 (Length=8)

	Raw []string // caret 분리 원본 (escape hatch)
}

// KrxNightOptionAskEvent 는 H0EUASP0 KRX야간옵션 실시간호가 이벤트 (38 fields).
// 5단계 호가, OPTN_ prefix. tr_key: 야간옵션 종목코드 12자리
type KrxNightOptionAskEvent struct {
	Symbol string // 0  OPTN_SHRN_ISCD 옵션단축종목코드 (9자리)
	Time   string // 1  BSOP_HOUR 영업시간 (HHMMSS)

	Ask     [5]decimal.Decimal // 2..6  OPTN_ASKP1..5 옵션매도호가1~5
	Bid     [5]decimal.Decimal // 7..11 OPTN_BIDP1..5 옵션매수호가1~5
	AskCsnu [5]int64           // 12..16 ASKP_CSNU1..5 매도호가건수1~5
	BidCsnu [5]int64           // 17..21 BIDP_CSNU1..5 매수호가건수1~5
	AskSize [5]int64           // 22..26 ASKP_RSQN1..5 매도호가잔량1~5
	BidSize [5]int64           // 27..31 BIDP_RSQN1..5 매수호가잔량1~5

	TotalAskCsnu    int64 // 32 TOTAL_ASKP_CSNU 총매도호가건수
	TotalBidCsnu    int64 // 33 TOTAL_BIDP_CSNU 총매수호가건수
	TotalAskSize    int64 // 34 TOTAL_ASKP_RSQN 총매도호가잔량
	TotalBidSize    int64 // 35 TOTAL_BIDP_RSQN 총매수호가잔량
	TotalAskSizeChg int64 // 36 TOTAL_ASKP_RSQN_ICDC 총매도호가잔량증감
	TotalBidSizeChg int64 // 37 TOTAL_BIDP_RSQN_ICDC 총매수호가잔량증감

	Raw []string // caret 분리 원본 (escape hatch)
}

// KrxNightOptionExpectTradeEvent 는 H0EUANC0 KRX야간옵션 실시간예상체결 이벤트 (8 fields).
// tr_key: 야간옵션 종목코드 12자리
type KrxNightOptionExpectTradeEvent struct {
	Symbol           string          // 0  OPTN_SHRN_ISCD 옵션단축종목코드 (9자리)
	Time             string          // 1  BSOP_HOUR 영업시간 (HHMMSS)
	ExpectPrice      decimal.Decimal // 2  ANTC_CNPR 예상체결가 (Length=8)
	ExpectDiff       decimal.Decimal // 3  ANTC_CNTG_VRSS 예상체결대비 (Length=8)
	ExpectDiffSign   string          // 4  ANTC_CNTG_VRSS_SIGN 예상체결대비부호 (1자리)
	ExpectChangeRate float64         // 5  ANTC_CNTG_PRDY_CTRT 예상체결전일대비율 (Length=8)
	ExpectMarketCode string          // 6  ANTC_MKOP_CLS_CODE 예상장운영구분코드 (3자리)
	ExpectQuantity   int64           // 7  ANTC_CNQN 예상체결수량 (Number 타입)

	Raw []string // caret 분리 원본 (escape hatch)
}

// StockFuturesTradeEvent 는 H0ZFCNT0 주식선물 실시간체결가 이벤트 (49 fields).
// STCK_* prefix 가격 필드. tr_key: 주식선물 종목코드 6자리
type StockFuturesTradeEvent struct {
	Symbol                 string          // 0  FUTS_SHRN_ISCD 선물단축종목코드 (9자리)
	Time                   string          // 1  BSOP_HOUR 영업시간 (HHMMSS)
	Price                  decimal.Decimal // 2  STCK_PRPR 주식현재가
	PrevDiffSign           string          // 3  PRDY_VRSS_SIGN 전일대비부호
	PrevDiff               decimal.Decimal // 4  PRDY_VRSS 전일대비
	PrevChangeRate         float64         // 5  FUTS_PRDY_CTRT 선물전일대비율
	Open                   decimal.Decimal // 6  STCK_OPRC 주식시가2
	High                   decimal.Decimal // 7  STCK_HGPR 주식최고가
	Low                    decimal.Decimal // 8  STCK_LWPR 주식최저가
	LastTradeVolume        int64           // 9  LAST_CNQN 최종거래량
	AccumVolume            int64           // 10 ACML_VOL 누적거래량
	AccumValue             int64           // 11 ACML_TR_PBMN 누적거래대금
	TheoreticalPrice       decimal.Decimal // 12 HTS_THPR HTS이론가
	MarketBasis            decimal.Decimal // 13 MRKT_BASIS 시장베이시스
	DeviationRate          float64         // 14 DPRT 괴리율
	NearMonthPrice         decimal.Decimal // 15 NMSC_FCTN_STPL_PRC 근월물약정가
	FarMonthPrice          decimal.Decimal // 16 FMSC_FCTN_STPL_PRC 원월물약정가
	SpreadPrice            decimal.Decimal // 17 SPEAD_PRC 스프레드1
	OpenInterestQty        int64           // 18 HTS_OTST_STPL_QTY HTS미결제약정수량
	OpenInterestChange     int64           // 19 OTST_STPL_QTY_ICDC 미결제약정수량증감
	OpenTime               string          // 20 OPRC_HOUR 시가시간 (HHMMSS)
	OpenDiffSign           string          // 21 OPRC_VRSS_PRPR_SIGN 시가2대비현재가부호
	OpenDiff               decimal.Decimal // 22 OPRC_VRSS_PRPR 시가2대비현재가 (NMIX 아님)
	HighTime               string          // 23 HGPR_HOUR 최고가시간 (HHMMSS)
	HighDiffSign           string          // 24 HGPR_VRSS_PRPR_SIGN 최고가대비현재가부호
	HighDiff               decimal.Decimal // 25 HGPR_VRSS_PRPR 최고가대비현재가
	LowTime                string          // 26 LWPR_HOUR 최저가시간 (HHMMSS)
	LowDiffSign            string          // 27 LWPR_VRSS_PRPR_SIGN 최저가대비현재가부호
	LowDiff                decimal.Decimal // 28 LWPR_VRSS_PRPR 최저가대비현재가
	BidRate                float64         // 29 SHNU_RATE 매수2비율
	TradeStrength          float64         // 30 CTTR 체결강도
	DeviationDegree        decimal.Decimal // 31 ESDG 괴리도
	OpenInterestPrevChange int64           // 32 OTST_STPL_RGBF_QTY_ICDC 미결제약정직전수량증감
	TheoreticalBasis       decimal.Decimal // 33 THPR_BASIS 이론베이시스
	Ask1                   decimal.Decimal // 34 ASKP1 매도호가1 (FUTS_ prefix 없음)
	Bid1                   decimal.Decimal // 35 BIDP1 매수호가1
	Ask1Size               int64           // 36 ASKP_RSQN1 매도호가잔량1
	Bid1Size               int64           // 37 BIDP_RSQN1 매수호가잔량1
	AskCount               int64           // 38 SELN_CNTG_CSNU 매도체결건수
	BidCount               int64           // 39 SHNU_CNTG_CSNU 매수체결건수
	NetCount               int64           // 40 NTBY_CNTG_CSNU 순매수체결건수
	TotalAskVolume         int64           // 41 SELN_CNTG_SMTN 총매도수량
	TotalBidVolume         int64           // 42 SHNU_CNTG_SMTN 총매수수량
	TotalAskSize           int64           // 43 TOTAL_ASKP_RSQN 총매도호가잔량
	TotalBidSize           int64           // 44 TOTAL_BIDP_RSQN 총매수호가잔량
	PrevVolRate            float64         // 45 PRDY_VOL_VRSS_ACML_VOL_RATE 전일거래량대비등락율
	DynamicUpperLimit      decimal.Decimal // 46 DYNM_MXPR 실시간상한가 (Length=4)
	DynamicLowerLimit      decimal.Decimal // 47 DYNM_LLAM 실시간하한가 (Length=4)
	DynamicPriceLimitYN    string          // 48 DYNM_PRC_LIMT_YN 실시간가격제한구분 (1자리)

	Raw []string // caret 분리 원본 (escape hatch)
}

// StockFuturesAskEvent 는 H0ZFASP0 주식선물 실시간호가 이벤트 (68 fields).
// 10단계 호가, ASKP/BIDP prefix 없음. tr_key: 주식선물 종목코드 6자리
type StockFuturesAskEvent struct {
	Symbol string // 0  FUTS_SHRN_ISCD 선물단축종목코드
	Time   string // 1  BSOP_HOUR 영업시간 (HHMMSS)

	Ask     [10]decimal.Decimal // 2..11  ASKP1..10 매도호가1~10
	Bid     [10]decimal.Decimal // 12..21 BIDP1..10 매수호가1~10
	AskCsnu [10]int64           // 22..31 ASKP_CSNU1..10 매도호가건수1~10
	BidCsnu [10]int64           // 32..41 BIDP_CSNU1..10 매수호가건수1~10
	AskSize [10]int64           // 42..51 ASKP_RSQN1..10 매도호가잔량1~10
	BidSize [10]int64           // 52..61 BIDP_RSQN1..10 매수호가잔량1~10

	TotalAskCsnu    int64 // 62 TOTAL_ASKP_CSNU 총매도호가건수
	TotalBidCsnu    int64 // 63 TOTAL_BIDP_CSNU 총매수호가건수
	TotalAskSize    int64 // 64 TOTAL_ASKP_RSQN 총매도호가잔량
	TotalBidSize    int64 // 65 TOTAL_BIDP_RSQN 총매수호가잔량
	TotalAskSizeChg int64 // 66 TOTAL_ASKP_RSQN_ICDC 총매도호가잔량증감
	TotalBidSizeChg int64 // 67 TOTAL_BIDP_RSQN_ICDC 총매수호가잔량증감

	Raw []string // caret 분리 원본 (escape hatch)
}

// StockFuturesExpectTradeEvent 는 H0ZFANC0 주식선물 실시간예상체결 이벤트 (8 fields).
// tr_key: 주식선물 종목코드 12자리 (docs 에 12자리 명시)
type StockFuturesExpectTradeEvent struct {
	Symbol           string          // 0  FUTS_SHRN_ISCD 선물단축종목코드 (9자리)
	Time             string          // 1  BSOP_HOUR 영업시간 (HHMMSS)
	ExpectPrice      decimal.Decimal // 2  ANTC_CNPR 예상체결가 (Length=8)
	ExpectDiff       decimal.Decimal // 3  ANTC_CNTG_VRSS 예상체결대비 (Length=8)
	ExpectDiffSign   string          // 4  ANTC_CNTG_VRSS_SIGN 예상체결대비부호 (1자리)
	ExpectChangeRate float64         // 5  ANTC_CNTG_PRDY_CTRT 예상체결전일대비율 (Length=8)
	ExpectMarketCode string          // 6  ANTC_MKOP_CLS_CODE 예상장운영구분코드 (3자리)
	ExpectQuantity   int64           // 7  ANTC_CNQN 예상체결수량 (String 타입)

	Raw []string // caret 분리 원본 (escape hatch)
}

// StockOptionTradeEvent 는 H0ZOCNT0 주식옵션 실시간체결가 이벤트 (53 fields).
// 옵션 그릭스 포함, DYNM 없음. tr_key: 주식옵션 종목코드 6자리
type StockOptionTradeEvent struct {
	Symbol                 string          // 0  OPTN_SHRN_ISCD 옵션단축종목코드 (9자리)
	Time                   string          // 1  BSOP_HOUR 영업시간 (HHMMSS)
	Price                  decimal.Decimal // 2  OPTN_PRPR 옵션현재가
	PrevDiffSign           string          // 3  PRDY_VRSS_SIGN 전일대비부호
	PrevDiff               decimal.Decimal // 4  OPTN_PRDY_VRSS 옵션전일대비
	PrevChangeRate         float64         // 5  PRDY_CTRT 전일대비율
	Open                   decimal.Decimal // 6  OPTN_OPRC 옵션시가2
	High                   decimal.Decimal // 7  OPTN_HGPR 옵션최고가
	Low                    decimal.Decimal // 8  OPTN_LWPR 옵션최저가
	LastTradeVolume        int64           // 9  LAST_CNQN 최종거래량
	AccumVolume            int64           // 10 ACML_VOL 누적거래량
	AccumValue             int64           // 11 ACML_TR_PBMN 누적거래대금
	TheoreticalPrice       decimal.Decimal // 12 HTS_THPR HTS이론가
	OpenInterestQty        int64           // 13 HTS_OTST_STPL_QTY HTS미결제약정수량
	OpenInterestChange     int64           // 14 OTST_STPL_QTY_ICDC 미결제약정수량증감
	OpenTime               string          // 15 OPRC_HOUR 시가시간 (HHMMSS)
	OpenDiffSign           string          // 16 OPRC_VRSS_PRPR_SIGN 시가2대비현재가부호
	OpenDiff               decimal.Decimal // 17 OPRC_VRSS_NMIX_PRPR 시가대비지수현재가
	HighTime               string          // 18 HGPR_HOUR 최고가시간 (HHMMSS)
	HighDiffSign           string          // 19 HGPR_VRSS_PRPR_SIGN 최고가대비현재가부호
	HighDiff               decimal.Decimal // 20 HGPR_VRSS_NMIX_PRPR 최고가대비지수현재가
	LowTime                string          // 21 LWPR_HOUR 최저가시간 (HHMMSS)
	LowDiffSign            string          // 22 LWPR_VRSS_PRPR_SIGN 최저가대비현재가부호
	LowDiff                decimal.Decimal // 23 LWPR_VRSS_NMIX_PRPR 최저가대비지수현재가
	BidRate                float64         // 24 SHNU_RATE 매수2비율
	PremiumValue           decimal.Decimal // 25 PRMM_VAL 프리미엄값
	IntrinsicValue         decimal.Decimal // 26 INVL_VAL 내재가치값
	TimeValue              decimal.Decimal // 27 TMVL_VAL 시간가치값
	Delta                  float64         // 28 DELTA 델타
	Gamma                  float64         // 29 GAMA 감마
	Vega                   float64         // 30 VEGA 베가
	Theta                  float64         // 31 THETA 세타
	Rho                    float64         // 32 RHO 로우
	ImpliedVolatility      float64         // 33 HTS_INTS_VLTL HTS내재변동성
	DeviationDegree        decimal.Decimal // 34 ESDG 괴리도
	OpenInterestPrevChange int64           // 35 OTST_STPL_RGBF_QTY_ICDC 미결제약정직전수량증감
	TheoreticalBasis       decimal.Decimal // 36 THPR_BASIS 이론베이시스
	HistoricalVolatility   float64         // 37 UNAS_HIST_VLTL 역사적변동성
	TradeStrength          float64         // 38 CTTR 체결강도
	DeviationRate          float64         // 39 DPRT 괴리율
	MarketBasis            decimal.Decimal // 40 MRKT_BASIS 시장베이시스
	Ask1                   decimal.Decimal // 41 OPTN_ASKP1 옵션매도호가1
	Bid1                   decimal.Decimal // 42 OPTN_BIDP1 옵션매수호가1
	Ask1Size               int64           // 43 ASKP_RSQN1 매도호가잔량1
	Bid1Size               int64           // 44 BIDP_RSQN1 매수호가잔량1
	AskCount               int64           // 45 SELN_CNTG_CSNU 매도체결건수
	BidCount               int64           // 46 SHNU_CNTG_CSNU 매수체결건수
	NetCount               int64           // 47 NTBY_CNTG_CSNU 순매수체결건수
	TotalAskVolume         int64           // 48 SELN_CNTG_SMTN 총매도수량
	TotalBidVolume         int64           // 49 SHNU_CNTG_SMTN 총매수수량
	TotalAskSize           int64           // 50 TOTAL_ASKP_RSQN 총매도호가잔량
	TotalBidSize           int64           // 51 TOTAL_BIDP_RSQN 총매수호가잔량
	PrevVolRate            float64         // 52 PRDY_VOL_VRSS_ACML_VOL_RATE 전일거래량대비등락율

	Raw []string // caret 분리 원본 (escape hatch)
}

// StockOptionAskEvent 는 H0ZOASP0 주식옵션 실시간호가 이벤트 (68 fields).
// 10단계 호가, OPTN_ASKP/BIDP prefix.
// docs 구조 특이: OPTN_ASKP1..5 + 총계 + OPTN_ASKP6..10 형태이나 배열로 통합 표현.
// tr_key: 주식옵션 종목코드 6자리
type StockOptionAskEvent struct {
	Symbol string // 0  OPTN_SHRN_ISCD 옵션단축종목코드
	Time   string // 1  BSOP_HOUR 영업시간 (HHMMSS)

	// 1~5단계 (index 0..4)
	Ask1to5     [5]decimal.Decimal // 2..6   OPTN_ASKP1..5 옵션매도호가1~5
	Bid1to5     [5]decimal.Decimal // 7..11  OPTN_BIDP1..5 옵션매수호가1~5
	AskCsnu1to5 [5]int64           // 12..16 ASKP_CSNU1..5 매도호가건수1~5
	BidCsnu1to5 [5]int64           // 17..21 BIDP_CSNU1..5 매수호가건수1~5
	AskSize1to5 [5]int64           // 22..26 ASKP_RSQN1..5 매도호가잔량1~5
	BidSize1to5 [5]int64           // 27..31 BIDP_RSQN1..5 매수호가잔량1~5

	TotalAskCsnu    int64 // 32 TOTAL_ASKP_CSNU 총매도호가건수
	TotalBidCsnu    int64 // 33 TOTAL_BIDP_CSNU 총매수호가건수
	TotalAskSize    int64 // 34 TOTAL_ASKP_RSQN 총매도호가잔량
	TotalBidSize    int64 // 35 TOTAL_BIDP_RSQN 총매수호가잔량
	TotalAskSizeChg int64 // 36 TOTAL_ASKP_RSQN_ICDC 총매도호가잔량증감
	TotalBidSizeChg int64 // 37 TOTAL_BIDP_RSQN_ICDC 총매수호가잔량증감

	// 6~10단계 (index 5..9)
	Ask6to10     [5]decimal.Decimal // 38..42 OPTN_ASKP6..10 옵션매도호가6~10
	Bid6to10     [5]decimal.Decimal // 43..47 OPTN_BIDP6..10 옵션매수호가6~10
	AskCsnu6to10 [5]int64           // 48..52 ASKP_CSNU6..10 매도호가건수6~10
	BidCsnu6to10 [5]int64           // 53..57 BIDP_CSNU6..10 매수호가건수6~10
	AskSize6to10 [5]int64           // 58..62 ASKP_RSQN6..10 매도호가잔량6~10
	BidSize6to10 [5]int64           // 63..67 BIDP_RSQN6..10 매수호가잔량6~10

	Raw []string // caret 분리 원본 (escape hatch)
}

// StockOptionExpectTradeEvent 는 H0ZOANC0 주식옵션 실시간예상체결 이벤트 (7 fields).
// 주의: ANTC_CNQN 없음 (다른 예상체결 EP 는 8 fields).
// tr_key: 주식옵션 종목코드 12자리
type StockOptionExpectTradeEvent struct {
	Symbol           string          // 0  OPTN_SHRN_ISCD 옵션단축종목코드 (9자리)
	Time             string          // 1  BSOP_HOUR 영업시간 (HHMMSS)
	ExpectPrice      decimal.Decimal // 2  ANTC_CNPR 예상체결가 (Length=8)
	ExpectDiff       decimal.Decimal // 3  ANTC_CNTG_VRSS 예상체결대비 (Length=8)
	ExpectDiffSign   string          // 4  ANTC_CNTG_VRSS_SIGN 예상체결대비부호 (1자리)
	ExpectChangeRate float64         // 5  ANTC_CNTG_PRDY_CTRT 예상체결전일대비율 (Length=8)
	ExpectMarketCode string          // 6  ANTC_MKOP_CLS_CODE 예상장운영구분코드 (3자리)
	// ANTC_CNQN 없음 — docs anomaly

	Raw []string // caret 분리 원본 (escape hatch)
}
