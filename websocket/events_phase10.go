package websocket

import "github.com/shopspring/decimal"

// OverseasTradeEvent 는 해외주식 실시간 (지연) 체결가 이벤트 (HDFSCNT0, 26 fields).
//
// 해외 시세는 미국이 0분 지연 (무료), 아시아는 15분 지연 (유료 신청 시 실시간).
// tr_key 형식: D|R + 시장구분(3자리, NAS/NYS/AMS/HKS/TSE 등) + 종목코드 (예: "DNASAAPL").
//
// KIS docs 는 모든 응답 필드를 String 으로 표기하지만, 본 struct 는 KRX 패턴 따라
// 가격→decimal, 수량→int64, 비율→float64 매핑 (decoder 가 string→타입 변환).
type OverseasTradeEvent struct {
	Symbol       string // 1  RSYM 실시간종목코드 (D/R+시장+종목)
	SymbolCode   string // 2  SYMB 종목코드
	Decimals     string // 3  ZDIV 소수점자리수
	LocalDate    string // 4  TYMD 현지영업일자 (YYYYMMDD)
	LocalDayDate string // 5  XYMD 현지일자
	LocalTime    string // 6  XHMS 현지시간 (HHMMSS)
	KrDate       string // 7  KYMD 한국일자
	KrTime       string // 8  KHMS 한국시간 (HHMMSS)

	Open decimal.Decimal // 9  OPEN 시가
	High decimal.Decimal // 10 HIGH 고가
	Low  decimal.Decimal // 11 LOW 저가
	Last decimal.Decimal // 12 LAST 현재가

	PrevDiffSign string          // 13 SIGN 대비구분
	PrevDiff     decimal.Decimal // 14 DIFF 전일대비
	ChangeRate   float64         // 15 RATE 등락율

	Bid     decimal.Decimal // 16 PBID 매수호가
	Ask     decimal.Decimal // 17 PASK 매도호가
	BidSize int64           // 18 VBID 매수잔량
	AskSize int64           // 19 VASK 매도잔량

	TradeVolume int64 // 20 EVOL 체결량
	AccumVolume int64 // 21 TVOL 거래량
	AccumValue  int64 // 22 TAMT 거래대금

	// BIVL: 매수호가가 매도주문 수량을 따라가서 체결 (KIS docs 명시).
	// ASVL: 매도호가가 매수주문 수량을 따라가서 체결.
	AskTradeVol int64 // 23 BIVL 매도체결량
	BidTradeVol int64 // 24 ASVL 매수체결량

	TradeStrength float64 // 25 STRN 체결강도

	MarketKind string // 26 MTYP 시장구분 (1:장중, 2:장전, 3:장후)

	Raw []string // caret 분리 원본 (escape hatch)
}

// OverseasAskEvent 는 해외주식 실시간호가 이벤트 (HDFSASP0, 17 fields).
//
// 해외는 KRX 의 10단계 호가와 다르게 1단계 호가만 제공.
// 미국은 무료 (0분 지연), 아시아는 유료 신청 시.
type OverseasAskEvent struct {
	Symbol       string // 1  RSYM 실시간종목코드
	SymbolCode   string // 2  SYMB 종목코드
	Decimals     string // 3  ZDIV 소수점자리수
	LocalDayDate string // 4  XYMD 현지일자
	LocalTime    string // 5  XHMS 현지시간
	KrDate       string // 6  KYMD 한국일자
	KrTime       string // 7  KHMS 한국시간

	TotalBidSize       int64 // 8  BVOL 매수총잔량
	TotalAskSize       int64 // 9  AVOL 매도총잔량
	TotalBidSizeChange int64 // 10 BDVL 매수총잔량대비
	TotalAskSizeChange int64 // 11 ADVL 매도총잔량대비

	Bid1           decimal.Decimal // 12 PBID1 매수호가1
	Ask1           decimal.Decimal // 13 PASK1 매도호가1
	Bid1Size       int64           // 14 VBID1 매수잔량1
	Ask1Size       int64           // 15 VASK1 매도잔량1
	Bid1SizeChange int64           // 16 DBID1 매수잔량대비1
	Ask1SizeChange int64           // 17 DASK1 매도잔량대비1

	Raw []string // caret 분리 원본 (escape hatch)
}
