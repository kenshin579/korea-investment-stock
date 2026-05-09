package websocket

import "github.com/shopspring/decimal"

// KrxTradeEvent 는 H0STCNT0 실시간체결가 (KRX 본장) 이벤트 (46 fields).
type KrxTradeEvent struct {
	Symbol                   string          // 1  MKSC_SHRN_ISCD 단축종목코드
	Time                     string          // 2  STCK_CNTG_HOUR 체결시간 (HHMMSS)
	Price                    decimal.Decimal // 3  STCK_PRPR 현재가
	PrevDiffSign             string          // 4  PRDY_VRSS_SIGN 전일대비부호 (1=상한,2=상승,3=보합,4=하한,5=하락)
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
	TradeKind                string          // 22 CCLD_DVSN 체결구분 (1=매수,3=장전,5=매도)
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

// KrxAskEvent 는 H0STASP0 실시간호가 (KRX) 이벤트 (59 fields).
type KrxAskEvent struct {
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

	DealCode string // 59 STCK_DEAL_CLS_CODE 주식매매구분코드 (사용X)

	Raw []string // caret 분리 원본 (escape hatch)
}

// KrxExpectTradeEvent 는 H0STANC0 실시간예상체결 (KRX 본장) 이벤트 (45 fields).
// H0STCNT0 와 거의 동일하나 22번이 CNTG_CLS_CODE (CCLD_DVSN 아님), VI_STND_PRC 없음.
type KrxExpectTradeEvent struct {
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
	TradeKind                string          // 22 CNTG_CLS_CODE 체결구분 (H0STCNT0 의 CCLD_DVSN 에 해당, KIS docs 불일치)
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
	// VI_STND_PRC 없음 (H0STCNT0 #46 와 달리)

	Raw []string // caret 분리 원본 (escape hatch)
}

// KrxOvernightTradeEvent 는 H0STOUP0 시간외 실시간체결가 (KRX) 이벤트 (43 fields).
// H0STCNT0 의 subset. 22번 = CNTG_CLS_CODE. #44/#45/#46 없음.
type KrxOvernightTradeEvent struct {
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
	TradeKind                string          // 22 CNTG_CLS_CODE 체결구분 (시간외, H0STOUP0 표기)
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
	// HOUR_CLS_CODE / MRKT_TRTM_CLS_CODE / VI_STND_PRC 없음 (H0STCNT0 #44~#46)

	Raw []string // caret 분리 원본 (escape hatch)
}

// KrxOvernightExpectEvent 는 H0STOAC0 시간외 실시간예상체결 (KRX) 이벤트 (43 fields).
// H0STANC0 (45) 의 subset. #44 HOUR_CLS_CODE / #45 MRKT_TRTM_CLS_CODE 없음.
// 22번 = CNTG_CLS_CODE.
type KrxOvernightExpectEvent struct {
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
	TradeKind                string          // 22 CNTG_CLS_CODE 체결구분 (시간외 예상체결)
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
	// HOUR_CLS_CODE / MRKT_TRTM_CLS_CODE 없음 (H0STANC0 #44~#45)

	Raw []string // caret 분리 원본 (escape hatch)
}
