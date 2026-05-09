package websocket

import "github.com/shopspring/decimal"

// AltMarketTradeEvent 는 NXT/통합 실시간체결가 이벤트 (H0NXCNT0 / H0UNCNT0, 46 fields).
//
// KRX H0STCNT0 (KrxTradeEvent) 와 schema 거의 동일. 차이는 22번 필드명 — KRX 는
// CCLD_DVSN, NXT/통합은 CNTG_CLS_CODE (의미 동일, KIS docs 표기 차이).
//
// NXT 와 통합은 schema 가 완전히 동일해서 base struct 1개 + type alias 2개로 처리.
type AltMarketTradeEvent struct {
	Symbol                   string          // 1  MKSC_SHRN_ISCD 단축종목코드
	Time                     string          // 2  STCK_CNTG_HOUR 체결시간 (HHMMSS)
	Price                    decimal.Decimal // 3  STCK_PRPR 현재가
	PrevDiffSign             string          // 4  PRDY_VRSS_SIGN 전일대비부호
	PrevDiff                 decimal.Decimal // 5  PRDY_VRSS 전일대비
	PrevChangeRate           float64         // 6  PRDY_CTRT 전일대비율
	WeightedAvg              decimal.Decimal // 7  WGHN_AVRG_STCK_PRC 가중평균주식가격
	Open                     decimal.Decimal // 8  STCK_OPRC 시가
	High                     decimal.Decimal // 9  STCK_HGPR 최고가
	Low                      decimal.Decimal // 10 STCK_LWPR 최저가
	Ask1                     decimal.Decimal // 11 ASKP1 매도호가1
	Bid1                     decimal.Decimal // 12 BIDP1 매수호가1
	TradeVolume              int64           // 13 CNTG_VOL 체결거래량
	AccumVolume              int64           // 14 ACML_VOL 누적거래량
	AccumValue               int64           // 15 ACML_TR_PBMN 누적거래대금
	AskCount                 int64           // 16 SELN_CNTG_CSNU 매도체결건수
	BidCount                 int64           // 17 SHNU_CNTG_CSNU 매수체결건수
	NetCount                 int64           // 18 NTBY_CNTG_CSNU 순매수체결건수
	TradeStrength            float64         // 19 CTTR 체결강도
	TotalAskVolume           int64           // 20 SELN_CNTG_SMTN 총매도수량
	TotalBidVolume           int64           // 21 SHNU_CNTG_SMTN 총매수수량
	TradeKind                string          // 22 CNTG_CLS_CODE 체결구분 (KRX 의 CCLD_DVSN 와 의미 동일)
	BidRate                  float64         // 23 SHNU_RATE 매수비율
	PrevVolRate              float64         // 24 PRDY_VOL_VRSS_ACML_VOL_RATE 전일거래량대비등락율
	OpenTime                 string          // 25 OPRC_HOUR 시가시간 (HHMMSS)
	OpenDiffSign             string          // 26 OPRC_VRSS_PRPR_SIGN 시가대비구분
	OpenDiff                 decimal.Decimal // 27 OPRC_VRSS_PRPR 시가대비
	HighTime                 string          // 28 HGPR_HOUR 최고가시간 (HHMMSS)
	HighDiffSign             string          // 29 HGPR_VRSS_PRPR_SIGN 고가대비구분
	HighDiff                 decimal.Decimal // 30 HGPR_VRSS_PRPR 고가대비
	LowTime                  string          // 31 LWPR_HOUR 최저가시간 (HHMMSS)
	LowDiffSign              string          // 32 LWPR_VRSS_PRPR_SIGN 저가대비구분
	LowDiff                  decimal.Decimal // 33 LWPR_VRSS_PRPR 저가대비
	BusinessDate             string          // 34 BSOP_DATE 영업일자 (YYYYMMDD)
	MarketOpCode             string          // 35 NEW_MKOP_CLS_CODE 신장운영구분코드
	TradeHaltYN              string          // 36 TRHT_YN 거래정지여부
	Ask1Size                 int64           // 37 ASKP_RSQN1 매도호가잔량1
	Bid1Size                 int64           // 38 BIDP_RSQN1 매수호가잔량1
	TotalAskSize             int64           // 39 TOTAL_ASKP_RSQN 총매도호가잔량
	TotalBidSize             int64           // 40 TOTAL_BIDP_RSQN 총매수호가잔량
	VolumeTurnover           float64         // 41 VOL_TNRT 거래량회전율
	PrevSameTimeAccumVol     int64           // 42 PRDY_SMNS_HOUR_ACML_VOL 전일동시간누적거래량
	PrevSameTimeAccumVolRate float64         // 43 PRDY_SMNS_HOUR_ACML_VOL_RATE 전일동시간누적거래량비율
	HourCode                 string          // 44 HOUR_CLS_CODE 시간구분코드
	MarketTermCode           string          // 45 MRKT_TRTM_CLS_CODE 임의종료구분코드
	ViStandardPrice          decimal.Decimal // 46 VI_STND_PRC 정적VI발동기준가

	Raw []string // caret 분리 원본 (escape hatch)
}

// AltMarketAskEvent 는 NXT/통합 실시간호가 이벤트 (H0NXASP0 / H0UNASP0, 65 fields).
//
// KRX H0STASP0 (59 fields) 의 superset. 끝부분에 KRX/NXT 중간가 6 fields 추가:
// KMID_PRC, KMID_TOTAL_RSQN, KMID_CLS_CODE, NMID_PRC, NMID_TOTAL_RSQN, NMID_CLS_CODE.
type AltMarketAskEvent struct {
	Symbol   string // 1  MKSC_SHRN_ISCD 단축종목코드
	Time     string // 2  BSOP_HOUR 영업시간 (HHMMSS)
	HourCode string // 3  HOUR_CLS_CODE 시간구분코드

	Ask     [10]decimal.Decimal // 4-13  ASKP1..10 매도호가1~10
	Bid     [10]decimal.Decimal // 14-23 BIDP1..10 매수호가1~10
	AskSize [10]int64           // 24-33 ASKP_RSQN1..10 매도호가잔량1~10
	BidSize [10]int64           // 34-43 BIDP_RSQN1..10 매수호가잔량1~10

	TotalAskSize int64 // 44 TOTAL_ASKP_RSQN 총매도호가잔량
	TotalBidSize int64 // 45 TOTAL_BIDP_RSQN 총매수호가잔량

	OvernightTotalAskSize int64 // 46 OVTM_TOTAL_ASKP_RSQN 시간외총매도호가잔량
	OvernightTotalBidSize int64 // 47 OVTM_TOTAL_BIDP_RSQN 시간외총매수호가잔량

	ExpectPrice      decimal.Decimal // 48 ANTC_CNPR 예상체결가
	ExpectQuantity   int64           // 49 ANTC_CNQN 예상체결량
	ExpectVolume     int64           // 50 ANTC_VOL 예상거래량
	ExpectDiff       decimal.Decimal // 51 ANTC_CNTG_VRSS 예상체결대비
	ExpectDiffSign   string          // 52 ANTC_CNTG_VRSS_SIGN 예상체결대비부호
	ExpectChangeRate float64         // 53 ANTC_CNTG_PRDY_CTRT 예상체결전일대비율

	AccumVolume int64 // 54 ACML_VOL 누적거래량

	TotalAskSizeChange      int64 // 55 TOTAL_ASKP_RSQN_ICDC 총매도호가잔량증감
	TotalBidSizeChange      int64 // 56 TOTAL_BIDP_RSQN_ICDC 총매수호가잔량증감
	OvernightTotalAskChange int64 // 57 OVTM_TOTAL_ASKP_ICDC 시간외총매도호가증감
	OvernightTotalBidChange int64 // 58 OVTM_TOTAL_BIDP_ICDC 시간외총매수호가증감

	DealCode string // 59 STCK_DEAL_CLS_CODE 주식매매구분코드

	// NXT/통합 추가 6 필드 (KRX/NXT 중간가)
	KrxMidPrice     decimal.Decimal // 60 KMID_PRC KRX 중간가
	KrxMidTotalSize int64           // 61 KMID_TOTAL_RSQN KRX 중간가 총잔량
	KrxMidCode      string          // 62 KMID_CLS_CODE KRX 중간가 구분
	NxtMidPrice     decimal.Decimal // 63 NMID_PRC NXT 중간가
	NxtMidTotalSize int64           // 64 NMID_TOTAL_RSQN NXT 중간가 총잔량
	NxtMidCode      string          // 65 NMID_CLS_CODE NXT 중간가 구분

	Raw []string // caret 분리 원본 (escape hatch)
}

// AltMarketExpectTradeEvent 는 NXT/통합 실시간예상체결 이벤트 (H0NXANC0 / H0UNANC0, 46 fields).
//
// KRX H0STANC0 (45 fields) 의 superset. 끝에 VI_STND_PRC 1 field 추가.
// 22번 = CNTG_CLS_CODE (KRX 도 CNTG_CLS_CODE — 동일).
type AltMarketExpectTradeEvent struct {
	Symbol                   string          // 1  MKSC_SHRN_ISCD 단축종목코드
	Time                     string          // 2  STCK_CNTG_HOUR 체결시간 (HHMMSS)
	Price                    decimal.Decimal // 3  STCK_PRPR 현재가
	PrevDiffSign             string          // 4  PRDY_VRSS_SIGN 전일대비부호
	PrevDiff                 decimal.Decimal // 5  PRDY_VRSS 전일대비
	PrevChangeRate           float64         // 6  PRDY_CTRT 전일대비율
	WeightedAvg              decimal.Decimal // 7  WGHN_AVRG_STCK_PRC 가중평균주식가격
	Open                     decimal.Decimal // 8  STCK_OPRC 시가
	High                     decimal.Decimal // 9  STCK_HGPR 최고가
	Low                      decimal.Decimal // 10 STCK_LWPR 최저가
	Ask1                     decimal.Decimal // 11 ASKP1 매도호가1
	Bid1                     decimal.Decimal // 12 BIDP1 매수호가1
	TradeVolume              int64           // 13 CNTG_VOL 체결거래량
	AccumVolume              int64           // 14 ACML_VOL 누적거래량
	AccumValue               int64           // 15 ACML_TR_PBMN 누적거래대금
	AskCount                 int64           // 16 SELN_CNTG_CSNU 매도체결건수
	BidCount                 int64           // 17 SHNU_CNTG_CSNU 매수체결건수
	NetCount                 int64           // 18 NTBY_CNTG_CSNU 순매수체결건수
	TradeStrength            float64         // 19 CTTR 체결강도
	TotalAskVolume           int64           // 20 SELN_CNTG_SMTN 총매도수량
	TotalBidVolume           int64           // 21 SHNU_CNTG_SMTN 총매수수량
	TradeKind                string          // 22 CNTG_CLS_CODE 체결구분
	BidRate                  float64         // 23 SHNU_RATE 매수비율
	PrevVolRate              float64         // 24 PRDY_VOL_VRSS_ACML_VOL_RATE 전일거래량대비등락율
	OpenTime                 string          // 25 OPRC_HOUR 시가시간 (HHMMSS)
	OpenDiffSign             string          // 26 OPRC_VRSS_PRPR_SIGN 시가대비구분
	OpenDiff                 decimal.Decimal // 27 OPRC_VRSS_PRPR 시가대비
	HighTime                 string          // 28 HGPR_HOUR 최고가시간 (HHMMSS)
	HighDiffSign             string          // 29 HGPR_VRSS_PRPR_SIGN 고가대비구분
	HighDiff                 decimal.Decimal // 30 HGPR_VRSS_PRPR 고가대비
	LowTime                  string          // 31 LWPR_HOUR 최저가시간 (HHMMSS)
	LowDiffSign              string          // 32 LWPR_VRSS_PRPR_SIGN 저가대비구분
	LowDiff                  decimal.Decimal // 33 LWPR_VRSS_PRPR 저가대비
	BusinessDate             string          // 34 BSOP_DATE 영업일자 (YYYYMMDD)
	MarketOpCode             string          // 35 NEW_MKOP_CLS_CODE 신장운영구분코드
	TradeHaltYN              string          // 36 TRHT_YN 거래정지여부
	Ask1Size                 int64           // 37 ASKP_RSQN1 매도호가잔량1
	Bid1Size                 int64           // 38 BIDP_RSQN1 매수호가잔량1
	TotalAskSize             int64           // 39 TOTAL_ASKP_RSQN 총매도호가잔량
	TotalBidSize             int64           // 40 TOTAL_BIDP_RSQN 총매수호가잔량
	VolumeTurnover           float64         // 41 VOL_TNRT 거래량회전율
	PrevSameTimeAccumVol     int64           // 42 PRDY_SMNS_HOUR_ACML_VOL 전일동시간누적거래량
	PrevSameTimeAccumVolRate float64         // 43 PRDY_SMNS_HOUR_ACML_VOL_RATE 전일동시간누적거래량비율
	HourCode                 string          // 44 HOUR_CLS_CODE 시간구분코드
	MarketTermCode           string          // 45 MRKT_TRTM_CLS_CODE 임의종료구분코드
	ViStandardPrice          decimal.Decimal // 46 VI_STND_PRC 정적VI발동기준가 (KRX H0STANC0 와 차이)

	Raw []string // caret 분리 원본 (escape hatch)
}

// ProgramTradeEvent 는 NXT/통합 실시간프로그램매매 이벤트 (H0NXPGM0 / H0UNPGM0, 11 fields).
//
// KRX H0STPGM0 와는 schema 별개 (Phase 8 OoS — Phase 9 NXT/통합만 우선 구현).
type ProgramTradeEvent struct {
	Symbol           string // 1  MKSC_SHRN_ISCD 단축종목코드
	Time             string // 2  STCK_CNTG_HOUR 체결시간 (HHMMSS)
	AskQuantity      int64  // 3  SELN_CNQN 매도체결수량
	AskValue         int64  // 4  SELN_TR_PBMN 매도거래대금
	BidQuantity      int64  // 5  SHNU_CNQN 매수체결수량
	BidValue         int64  // 6  SHNU_TR_PBMN 매수거래대금
	NetQuantity      int64  // 7  NTBY_CNQN 순매수수량
	NetValue         int64  // 8  NTBY_TR_PBMN 순매수거래대금
	AskRemainingSize int64  // 9  SELN_RSQN 매도호가잔량
	BidRemainingSize int64  // 10 SHNU_RSQN 매수호가잔량
	TotalNetQuantity int64  // 11 WHOL_NTBY_QTY 전체순매수수량

	Raw []string // caret 분리 원본 (escape hatch)
}

// MemberEvent 는 NXT/통합 실시간회원사 이벤트 (H0NXMBC0 / H0UNMBC0, 78 fields).
//
// 5단계 매도/매수 회원사 + 외국계 통계 + 영문 회원사명. KRX H0STMBC0 와는 schema 별개.
type MemberEvent struct {
	Symbol string // 1  MKSC_SHRN_ISCD 단축종목코드

	SellBrokerNames [5]string // 2-6   SELN2_MBCR_NAME1..5 매도2 회원사명1~5
	BuyBrokerNames  [5]string // 7-11  BYOV_MBCR_NAME1..5 매수 회원사명1~5

	TotalSellQty [5]int64 // 12-16 TOTAL_SELN_QTY1..5 총매도수량1~5
	TotalBuyQty  [5]int64 // 17-21 TOTAL_SHNU_QTY1..5 총매수수량1~5

	SellGlobalYN [5]string // 22-26 SELN_MBCR_GLOB_YN_1..5 매도거래원구분1~5
	BuyGlobalYN  [5]string // 27-31 SHNU_MBCR_GLOB_YN_1..5 매수거래원구분1~5

	SellBrokerCodes [5]string // 32-36 SELN_MBCR_NO1..5 매도거래원코드1~5
	BuyBrokerCodes  [5]string // 37-41 SHNU_MBCR_NO1..5 매수거래원코드1~5

	SellRatio [5]float64 // 42-46 SELN_MBCR_RLIM1..5 매도 회원사 비중1~5
	BuyRatio  [5]float64 // 47-51 SHNU_MBCR_RLIM1..5 매수 회원사 비중1~5

	SellQtyChange [5]int64 // 52-56 SELN_QTY_ICDC1..5 매도수량증감1~5
	BuyQtyChange  [5]int64 // 57-61 SHNU_QTY_ICDC1..5 매수수량증감1~5

	GlobalTotalSellQty  int64   // 62 GLOB_TOTAL_SELN_QTY 외국계 총매도수량
	GlobalTotalBuyQty   int64   // 63 GLOB_TOTAL_SHNU_QTY 외국계 총매수수량
	GlobalSellQtyChange int64   // 64 GLOB_TOTAL_SELN_QTY_ICDC 외국계 총매도증감
	GlobalBuyQtyChange  int64   // 65 GLOB_TOTAL_SHNU_QTY_ICDC 외국계 총매수증감
	GlobalNetBuyQty     int64   // 66 GLOB_NTBY_QTY 외국계 순매수수량
	GlobalSellRatio     float64 // 67 GLOB_SELN_RLIM 외국계 매도비중
	GlobalBuyRatio      float64 // 68 GLOB_SHNU_RLIM 외국계 매수비중

	SellBrokerEngNames [5]string // 69-73 SELN2_MBCR_ENG_NAME1..5 매도2 영문회원사명1~5
	BuyBrokerEngNames  [5]string // 74-78 BYOV_MBCR_ENG_NAME1..5 매수 영문회원사명1~5

	Raw []string // caret 분리 원본 (escape hatch)
}

// 시장별 type alias — NXT 와 통합은 schema 가 동일하므로 base struct 의 alias.
//
// alias 는 Go 컴파일 시 base type 으로 해소되므로 handler 시그니처는 base type 과 호환.
// 사용자는 시장 명시적 이름 (NxtTradeEvent, UnifiedTradeEvent) 으로 코드 가독성 향상.
type (
	NxtTradeEvent            = AltMarketTradeEvent
	UnifiedTradeEvent        = AltMarketTradeEvent
	NxtAskEvent              = AltMarketAskEvent
	UnifiedAskEvent          = AltMarketAskEvent
	NxtExpectTradeEvent      = AltMarketExpectTradeEvent
	UnifiedExpectTradeEvent  = AltMarketExpectTradeEvent
	NxtProgramTradeEvent     = ProgramTradeEvent
	UnifiedProgramTradeEvent = ProgramTradeEvent
	NxtMemberEvent           = MemberEvent
	UnifiedMemberEvent       = MemberEvent
)
