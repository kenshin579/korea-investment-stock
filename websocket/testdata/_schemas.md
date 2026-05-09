# Phase 8 KRX 5 EP — Schema Reference

> docs analyzer (Task 2) 결과. 다음 task 들의 Event struct / decoder / fixture 의 source of truth.

## Plan Deviation

Plan 의 TR_ID 추정이 시간외 EP 2개에서 틀림. 정정:

| EP | Plan 추정 | **실제** |
|---|---|---|
| 시간외 체결가 KRX | H0STOAC0 | **H0STOUP0** |
| 시간외 예상체결 KRX | H0STOAA0 | **H0STOAC0** |

## EP 정리

| EP | TR_ID | 한글명 | Field count | 모의 |
|---|---|---|---|---|
| 1 | H0STCNT0 | 실시간체결가 (KRX) | 46 | 지원 |
| 2 | H0STASP0 | 실시간호가 (KRX) | 59 | 지원 |
| 3 | H0STANC0 | 실시간예상체결 (KRX) | 45 | 미지원 |
| 4 | H0STOUP0 | 시간외 실시간체결가 (KRX) | 43 | 미지원 |
| 5 | H0STOAC0 | 시간외 실시간예상체결 (KRX) | 43 | 미지원 |

공통:
- tr_key: 종목번호 6자리 (ETN 은 Q-prefix, 예: Q500001)
- 응답 포맷: `0|TR_ID|count|f1^f2^...^fN`
- count > 1 페이징 가능 (모든 EP)
- encrypted=1 가능 (Phase 8 미지원 — ErrWSEncryptedNotSupported)

---

## EP1 H0STCNT0 — 체결가 KRX (46 fields)

| # | 필드명 | 한글 | Type | 길이 |
|---|---|---|---|---|
| 1 | MKSC_SHRN_ISCD | 단축종목코드 | String | 9 |
| 2 | STCK_CNTG_HOUR | 체결시간 | String | 6 |
| 3 | STCK_PRPR | 현재가 | Number | 4 |
| 4 | PRDY_VRSS_SIGN | 전일대비부호 | String | 1 |
| 5 | PRDY_VRSS | 전일대비 | Number | 4 |
| 6 | PRDY_CTRT | 전일대비율 | Number | 8 |
| 7 | WGHN_AVRG_STCK_PRC | 가중평균주식가격 | Number | 8 |
| 8 | STCK_OPRC | 시가 | Number | 4 |
| 9 | STCK_HGPR | 최고가 | Number | 4 |
| 10 | STCK_LWPR | 최저가 | Number | 4 |
| 11 | ASKP1 | 매도호가1 | Number | 4 |
| 12 | BIDP1 | 매수호가1 | Number | 4 |
| 13 | CNTG_VOL | 체결거래량 | Number | 8 |
| 14 | ACML_VOL | 누적거래량 | Number | 8 |
| 15 | ACML_TR_PBMN | 누적거래대금 | Number | 8 |
| 16 | SELN_CNTG_CSNU | 매도체결건수 | Number | 4 |
| 17 | SHNU_CNTG_CSNU | 매수체결건수 | Number | 4 |
| 18 | NTBY_CNTG_CSNU | 순매수체결건수 | Number | 4 |
| 19 | CTTR | 체결강도 | Number | 8 |
| 20 | SELN_CNTG_SMTN | 총매도수량 | Number | 8 |
| 21 | SHNU_CNTG_SMTN | 총매수수량 | Number | 8 |
| 22 | **CCLD_DVSN** | 체결구분 | String | 1 |
| 23 | SHNU_RATE | 매수비율 | Number | 8 |
| 24 | PRDY_VOL_VRSS_ACML_VOL_RATE | 전일거래량대비등락율 | Number | 8 |
| 25 | OPRC_HOUR | 시가시간 | String | 6 |
| 26 | OPRC_VRSS_PRPR_SIGN | 시가대비구분 | String | 1 |
| 27 | OPRC_VRSS_PRPR | 시가대비 | Number | 4 |
| 28 | HGPR_HOUR | 최고가시간 | String | 6 |
| 29 | HGPR_VRSS_PRPR_SIGN | 고가대비구분 | String | 1 |
| 30 | HGPR_VRSS_PRPR | 고가대비 | Number | 4 |
| 31 | LWPR_HOUR | 최저가시간 | String | 6 |
| 32 | LWPR_VRSS_PRPR_SIGN | 저가대비구분 | String | 1 |
| 33 | LWPR_VRSS_PRPR | 저가대비 | Number | 4 |
| 34 | BSOP_DATE | 영업일자 | String | 8 |
| 35 | NEW_MKOP_CLS_CODE | 신장운영구분코드 | String | 2 |
| 36 | TRHT_YN | 거래정지여부 | String | 1 |
| 37 | ASKP_RSQN1 | 매도호가잔량1 | Number | 8 |
| 38 | BIDP_RSQN1 | 매수호가잔량1 | Number | 8 |
| 39 | TOTAL_ASKP_RSQN | 총매도호가잔량 | Number | 8 |
| 40 | TOTAL_BIDP_RSQN | 총매수호가잔량 | Number | 8 |
| 41 | VOL_TNRT | 거래량회전율 | Number | 8 |
| 42 | PRDY_SMNS_HOUR_ACML_VOL | 전일동시간누적거래량 | Number | 8 |
| 43 | PRDY_SMNS_HOUR_ACML_VOL_RATE | 전일동시간누적거래량비율 | Number | 8 |
| 44 | HOUR_CLS_CODE | 시간구분코드 | String | 1 |
| 45 | MRKT_TRTM_CLS_CODE | 임의종료구분코드 | String | 1 |
| 46 | VI_STND_PRC | 정적VI발동기준가 | Number | 4 |

PRDY_VRSS_SIGN 의미: 1=상한, 2=상승, 3=보합, 4=하한, 5=하락. CCLD_DVSN: 1=매수(+), 3=장전, 5=매도(-).

---

## EP2 H0STASP0 — 호가 KRX (59 fields)

| # | 필드명 | 한글 | Type | 길이 |
|---|---|---|---|---|
| 1 | MKSC_SHRN_ISCD | 단축종목코드 | String | 9 |
| 2 | BSOP_HOUR | 영업시간 | String | 6 |
| 3 | HOUR_CLS_CODE | 시간구분코드 | String | 1 |
| 4-13 | ASKP1..10 | 매도호가1~10 | Number | 4 |
| 14-23 | BIDP1..10 | 매수호가1~10 | Number | 4 |
| 24-33 | ASKP_RSQN1..10 | 매도호가잔량1~10 | Number | 8 |
| 34-43 | BIDP_RSQN1..10 | 매수호가잔량1~10 | Number | 8 |
| 44 | TOTAL_ASKP_RSQN | 총매도호가잔량 | Number | 8 |
| 45 | TOTAL_BIDP_RSQN | 총매수호가잔량 | Number | 8 |
| 46 | OVTM_TOTAL_ASKP_RSQN | 시간외총매도호가잔량 | Number | 8 |
| 47 | OVTM_TOTAL_BIDP_RSQN | 시간외총매수호가잔량 | Number | 8 |
| 48 | ANTC_CNPR | 예상체결가 | Number | 4 |
| 49 | ANTC_CNQN | 예상체결량 | Number | 8 |
| 50 | ANTC_VOL | 예상거래량 | Number | 8 |
| 51 | ANTC_CNTG_VRSS | 예상체결대비 | Number | 4 |
| 52 | ANTC_CNTG_VRSS_SIGN | 예상체결대비부호 | String | 1 |
| 53 | ANTC_CNTG_PRDY_CTRT | 예상체결전일대비율 | Number | 8 |
| 54 | ACML_VOL | 누적거래량 | Number | 8 |
| 55 | TOTAL_ASKP_RSQN_ICDC | 총매도호가잔량증감 | Number | 4 |
| 56 | TOTAL_BIDP_RSQN_ICDC | 총매수호가잔량증감 | Number | 4 |
| 57 | OVTM_TOTAL_ASKP_ICDC | 시간외총매도호가증감 | Number | 4 |
| 58 | OVTM_TOTAL_BIDP_ICDC | 시간외총매수호가증감 | Number | 4 |
| 59 | STCK_DEAL_CLS_CODE | 주식매매구분코드 (사용X) | String | 2 |

---

## EP3 H0STANC0 — 예상체결 KRX (45 fields)

H0STCNT0 와 거의 동일하나:
- 22번: **CNTG_CLS_CODE** (CCLD_DVSN 아님 — KIS docs inconsistency)
- 끝부분에 VI_STND_PRC 없음 (45 fields, 46 아님)
- 모든 필드 명시적으로 String 타입 (실제로는 Number 도 wire format 으로는 string 이지만 docs 표기 차이)

H0STCNT0 의 #44 HOUR_CLS_CODE / #45 MRKT_TRTM_CLS_CODE 는 H0STANC0 에서도 존재. #46 VI_STND_PRC 만 없음.

---

## EP4 H0STOUP0 — 시간외 체결가 KRX (43 fields)

H0STCNT0 의 subset. 22번 = **CNTG_CLS_CODE** (CCLD_DVSN 아님).
H0STCNT0 의 #44/#45/#46 (HOUR_CLS_CODE / MRKT_TRTM_CLS_CODE / VI_STND_PRC) 없음.

---

## EP5 H0STOAC0 — 시간외 예상체결 KRX (43 fields)

H0STANC0 (45) 의 subset. H0STANC0 의 #44 HOUR_CLS_CODE / #45 MRKT_TRTM_CLS_CODE 없음.
22번 = CNTG_CLS_CODE.

---

## Event Struct 매핑 결정

5 separate Event struct (재사용 불가):

```go
KrxTradeEvent           // H0STCNT0 (46 fields)
KrxAskEvent             // H0STASP0 (42 fields)
KrxExpectTradeEvent     // H0STANC0 (45 fields)
KrxOvernightTradeEvent  // H0STOUP0 (43 fields)
KrxOvernightExpectEvent // H0STOAC0 (43 fields)
```

이전 plan 의 "OvernightTrade/OvernightExpect 가 KrxTradeEvent/KrxExpectTradeEvent 재사용" 가정은 **틀림**. schema 다름.

## Type 매핑

| KIS Type | Go type |
|---|---|
| String (HHMMSS, code, sign) | `string` |
| Number 가격 (4자리) | `decimal.Decimal` (소수 가능) |
| Number 비율 (8자리) | `float64,string` |
| Number 거래량/금액 (8자리, 정수) | `int64,string` |
| Number 건수 (4자리) | `int64,string` |

CCLD_DVSN/CNTG_CLS_CODE → string. PRDY_VRSS_SIGN → string.
