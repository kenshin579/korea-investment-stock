# Phase 10 해외주식 실시간 2 EP — Schema Reference

> docs/api/해외주식/해외주식_실시간*.md (2026-05-09 직접 검증). Phase 10 Event/Decoder 의 source of truth.

## EP Matrix

| EP | TR_ID | Body fields | 모의 | 비고 |
|---|---|---|---|---|
| 체결가 (지연) | HDFSCNT0 | 26 | 미지원 | 미국 0분지연 무료 / 아시아 15분지연 |
| 호가 | HDFSASP0 | 17 | 미지원 | 미국 1호가 무료, 아시아 유료 신청 시 |

**중요**: KIS docs 가 모든 필드를 String 으로 표기하지만, type 매핑은 KRX/NXT 패턴 따라 가격→`decimal.Decimal`, 수량/거래량→`int64`, 비율/체결강도→`float64`, 그 외 (코드/시간/일자/sign)→`string`.

**tr_key 형식** (체결가/호가 공통):
- `D`+시장구분(3자리)+종목코드 — 무료 시세 (예: `DNASAAPL`)
- `R`+시장구분(3자리)+종목코드 — 유료 시세 + 미국주간 (예: `RBAQAAPL`)
- 시장구분: `NYS` 뉴욕, `NAS` 나스닥, `AMS` 아멕스, `TSE` 도쿄, `HKS` 홍콩, `SHS` 상해, `SZS` 심천, `HSX` 호치민, `HNX` 하노이, `BAY/BAQ/BAA` 미국 주간

---

## EP1 — HDFSCNT0 체결가 (지연) (26 fields)

```
[0]  RSYM  실시간종목코드 (16자, e.g. "DNASAAPL")    string
[1]  SYMB  종목코드                                  string
[2]  ZDIV  소수점자리수                              string
[3]  TYMD  현지영업일자 (YYYYMMDD)                   string
[4]  XYMD  현지일자                                  string
[5]  XHMS  현지시간 (HHMMSS)                         string
[6]  KYMD  한국일자                                  string
[7]  KHMS  한국시간 (HHMMSS)                         string
[8]  OPEN  시가                                      decimal
[9]  HIGH  고가                                      decimal
[10] LOW   저가                                      decimal
[11] LAST  현재가                                    decimal
[12] SIGN  대비구분                                  string
[13] DIFF  전일대비                                  decimal
[14] RATE  등락율                                    float64
[15] PBID  매수호가                                  decimal
[16] PASK  매도호가                                  decimal
[17] VBID  매수잔량                                  int64
[18] VASK  매도잔량                                  int64
[19] EVOL  체결량                                    int64
[20] TVOL  거래량                                    int64
[21] TAMT  거래대금                                  int64
[22] BIVL  매도체결량                                int64
[23] ASVL  매수체결량                                int64
[24] STRN  체결강도                                  float64
[25] MTYP  시장구분 (1:장중, 2:장전, 3:장후)         string
```

```go
type OverseasTradeEvent struct {
    Symbol       string  // RSYM 실시간종목코드 (D/R + 시장 + 종목)
    SymbolCode   string  // SYMB 종목코드
    Decimals     string  // ZDIV 소수점자리수
    LocalDate    string  // TYMD 현지영업일자 (YYYYMMDD)
    LocalDayDate string  // XYMD 현지일자
    LocalTime    string  // XHMS 현지시간 (HHMMSS)
    KrDate       string  // KYMD 한국일자
    KrTime       string  // KHMS 한국시간 (HHMMSS)
    Open, High, Low, Last decimal.Decimal
    PrevDiffSign string          // SIGN
    PrevDiff     decimal.Decimal // DIFF
    ChangeRate   float64         // RATE
    Bid, Ask     decimal.Decimal // PBID, PASK
    BidSize, AskSize int64       // VBID, VASK
    TradeVolume  int64           // EVOL 체결량
    AccumVolume  int64           // TVOL 거래량
    AccumValue   int64           // TAMT 거래대금
    AskTradeVol  int64           // BIVL 매도체결량 (매수가 매도주문 따라가서 체결)
    BidTradeVol  int64           // ASVL 매수체결량 (매도가 매수주문 따라가서 체결)
    TradeStrength float64        // STRN
    MarketKind   string          // MTYP 1:장중, 2:장전, 3:장후
    Raw []string
}
```

---

## EP2 — HDFSASP0 호가 (17 fields)

```
[0]  RSYM  실시간종목코드                            string
[1]  SYMB  종목코드                                  string
[2]  ZDIV  소수점자리수                              string
[3]  XYMD  현지일자                                  string
[4]  XHMS  현지시간                                  string
[5]  KYMD  한국일자                                  string
[6]  KHMS  한국시간                                  string
[7]  BVOL  매수총잔량                                int64
[8]  AVOL  매도총잔량                                int64
[9]  BDVL  매수총잔량대비                            int64
[10] ADVL  매도총잔량대비                            int64
[11] PBID1 매수호가1                                 decimal
[12] PASK1 매도호가1                                 decimal
[13] VBID1 매수잔량1                                 int64
[14] VASK1 매도잔량1                                 int64
[15] DBID1 매수잔량대비1                             int64
[16] DASK1 매도잔량대비1                             int64
```

```go
type OverseasAskEvent struct {
    Symbol       string
    SymbolCode   string
    Decimals     string
    LocalDayDate string
    LocalTime    string
    KrDate       string
    KrTime       string
    TotalBidSize int64           // BVOL
    TotalAskSize int64           // AVOL
    TotalBidSizeChange int64     // BDVL
    TotalAskSizeChange int64     // ADVL
    Bid1   decimal.Decimal       // PBID1
    Ask1   decimal.Decimal       // PASK1
    Bid1Size int64               // VBID1
    Ask1Size int64               // VASK1
    Bid1SizeChange int64         // DBID1
    Ask1SizeChange int64         // DASK1
    Raw []string
}
```

해외는 KRX 의 10단계 호가와 다르게 **1호가만** 제공.

---

## TR_ID 상수 (client.go 추가)

```go
const (
    trIDOverseasTrade = "HDFSCNT0" // 해외주식 실시간지연체결가
    trIDOverseasAsk   = "HDFSASP0" // 해외주식 실시간호가
)
```
