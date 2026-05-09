package websocket

import "github.com/shopspring/decimal"

// ---------------------------------------------------------------------------
// Phase 11.3 — 국내선물옵션 실시간 6 EP
//
// 4 base types + 2 alias.
// WebSocket 도메인: ws://ops.koreainvestment.com:21000
// 전체 모의투자 미지원.
//
// Base types:
//   IndexFuturesTradeEvent  (H0IFCNT0, 50 fields)
//   IndexFuturesAskEvent    (H0IFASP0, 38 fields)
//   IndexOptionTradeEvent   (H0IOCNT0, 58 fields)
//   IndexOptionAskEvent     (H0IOASP0, 38 fields)
//
// Alias:
//   CommodityFuturesTradeEvent = IndexFuturesTradeEvent (H0CFCNT0 — 완전 동일 schema)
//   CommodityFuturesAskEvent   = IndexFuturesAskEvent   (H0CFASP0 — 완전 동일 schema)
// ---------------------------------------------------------------------------

// IndexFuturesTradeEvent 는 H0IFCNT0 지수선물 실시간체결가 이벤트 (50 fields).
// tr_key: 지수선물 종목코드 6자리 (예: 101S12)
// H0CFCNT0 상품선물과 field 명/순서/Type 완전 동일 → CommodityFuturesTradeEvent alias.
type IndexFuturesTradeEvent struct {
	Symbol                 string          // 0  FUTS_SHRN_ISCD 선물단축종목코드
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
	OpenDiff               decimal.Decimal // 22 OPRC_VRSS_NMIX_PRPR 시가대비지수현재가 (NMIX)
	HighTime               string          // 23 HGPR_HOUR 최고가시간 (HHMMSS)
	HighDiffSign           string          // 24 HGPR_VRSS_PRPR_SIGN 최고가대비현재가부호
	HighDiff               decimal.Decimal // 25 HGPR_VRSS_NMIX_PRPR 최고가대비지수현재가 (NMIX)
	LowTime                string          // 26 LWPR_HOUR 최저가시간 (HHMMSS)
	LowDiffSign            string          // 27 LWPR_VRSS_PRPR_SIGN 최저가대비현재가부호
	LowDiff                decimal.Decimal // 28 LWPR_VRSS_NMIX_PRPR 최저가대비지수현재가 (NMIX)
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
	BlockTradeVolume       int64           // 46 DSCS_BLTR_ACML_QTY 협의대량거래량 (H0MFCNT0 없음)
	DynamicUpperLimit      decimal.Decimal // 47 DYNM_MXPR 실시간상한가
	DynamicLowerLimit      decimal.Decimal // 48 DYNM_LLAM 실시간하한가
	DynamicPriceLimitYN    string          // 49 DYNM_PRC_LIMT_YN 실시간가격제한구분 (1자리)

	Raw []string // caret 분리 원본 (escape hatch)
}

// IndexFuturesAskEvent 는 H0IFASP0 지수선물 실시간호가 이벤트 (38 fields).
// 5단계 호가, FUTS_ASKP/BIDP prefix.
// tr_key: 지수선물 종목코드 6자리
// H0CFASP0 상품선물과 완전 동일 schema → CommodityFuturesAskEvent alias.
type IndexFuturesAskEvent struct {
	Symbol string // 0  FUTS_SHRN_ISCD 선물단축종목코드
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

// IndexOptionTradeEvent 는 H0IOCNT0 지수옵션 실시간체결가 이벤트 (58 fields).
// 옵션 그릭스 + AVRG_VLTL + DSCS_LRQN_VOL + DYNM 3 fields 포함.
// tr_key: 지수옵션 종목코드 6자리 (예: 201S11305)
type IndexOptionTradeEvent struct {
	Symbol                 string          // 0  OPTN_SHRN_ISCD 옵션단축종목코드
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
	OpenDiff               decimal.Decimal // 17 OPRC_VRSS_NMIX_PRPR 시가대비지수현재가 (NMIX)
	HighTime               string          // 18 HGPR_HOUR 최고가시간 (HHMMSS)
	HighDiffSign           string          // 19 HGPR_VRSS_PRPR_SIGN 최고가대비현재가부호
	HighDiff               decimal.Decimal // 20 HGPR_VRSS_NMIX_PRPR 최고가대비지수현재가 (NMIX)
	LowTime                string          // 21 LWPR_HOUR 최저가시간 (HHMMSS)
	LowDiffSign            string          // 22 LWPR_VRSS_PRPR_SIGN 최저가대비현재가부호
	LowDiff                decimal.Decimal // 23 LWPR_VRSS_NMIX_PRPR 최저가대비지수현재가 (NMIX)
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
	AvgVolatility          float64         // 53 AVRG_VLTL 평균변동성 (ZOCNT0/EUCNT0 없음)
	BlockTradeVolume       int64           // 54 DSCS_LRQN_VOL 협의대량누적거래량 (ZOCNT0/EUCNT0 없음)
	DynamicUpperLimit      decimal.Decimal // 55 DYNM_MXPR 실시간상한가
	DynamicLowerLimit      decimal.Decimal // 56 DYNM_LLAM 실시간하한가
	DynamicPriceLimitYN    string          // 57 DYNM_PRC_LIMT_YN 실시간가격제한구분 (1자리)

	Raw []string // caret 분리 원본 (escape hatch)
}

// IndexOptionAskEvent 는 H0IOASP0 지수옵션 실시간호가 이벤트 (38 fields).
// 5단계 호가, OPTN_ASKP/BIDP prefix (IndexFuturesAskEvent 의 FUTS_ 와 다름).
// tr_key: 지수옵션 종목코드 6자리
type IndexOptionAskEvent struct {
	Symbol string // 0  OPTN_SHRN_ISCD 옵션단축종목코드
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

// ---------------------------------------------------------------------------
// Alias (2)
// ---------------------------------------------------------------------------

// CommodityFuturesTradeEvent 는 H0CFCNT0 상품선물 실시간체결가 이벤트.
// H0IFCNT0 와 field 명/순서/Type 완전 동일 → alias.
type CommodityFuturesTradeEvent = IndexFuturesTradeEvent

// CommodityFuturesAskEvent 는 H0CFASP0 상품선물 실시간호가 이벤트.
// H0IFASP0 와 field 명/순서/Type 완전 동일 → alias.
type CommodityFuturesAskEvent = IndexFuturesAskEvent
