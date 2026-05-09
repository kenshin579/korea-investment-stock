# Phase 11.2 국내선물옵션 실시간 11 EP — Schema Reference

> docs analyzer 결과 (2026-05-09). Phase 11.2 Event/Decoder/Client 의 source of truth.
> docs 위치: `docs/api/국내선물옵션/<파일>.md`

---

## EP Matrix

| # | TR_ID | 한글명 | Fields | 모의 | tr_key 형식 | 종목코드 길이 |
|---|---|---|---|---|---|---|
| 1 | H0MFCNT0 | KRX야간선물 실시간종목체결 | 46 | 미지원 | 야간선물 종목코드 | 12자리 |
| 2 | H0MFASP0 | KRX야간선물 실시간호가 | 36 | 미지원 | 야간선물 종목코드 | 12자리 |
| 3 | H0EUCNT0 | KRX야간옵션 실시간체결가 | 51 | 미지원 | 야간옵션 종목코드 | 12자리 |
| 4 | H0EUASP0 | KRX야간옵션 실시간호가 | 36 | 미지원 | 야간옵션 종목코드 | 12자리 |
| 5 | H0EUANC0 | KRX야간옵션 실시간예상체결 | 8 | 미지원 | 야간옵션 종목코드 | 12자리 |
| 6 | H0ZFCNT0 | 주식선물 실시간체결가 | 47 | 미지원 | 주식선물 종목코드 | 6자리 |
| 7 | H0ZFASP0 | 주식선물 실시간호가 | 46 | 미지원 | 주식선물 종목코드 | 6자리 |
| 8 | H0ZFANC0 | 주식선물 실시간예상체결 | 8 | 미지원 | 주식선물 종목코드 | 12자리 |
| 9 | H0ZOCNT0 | 주식옵션 실시간체결가 | 47 | 미지원 | 주식옵션 종목코드 | 6자리 |
| 10 | H0ZOASP0 | 주식옵션 실시간호가 | 56 | 미지원 | 주식옵션 종목코드 | 6자리 |
| 11 | H0ZOANC0 | 주식옵션 실시간예상체결 | 7 | 미지원 | 주식옵션 종목코드 | 12자리 |

**전체 모의투자 미지원.** WebSocket domain: `ws://ops.koreainvestment.com:21000`

---

## Alias 분석 결과

### 카테고리별 schema 비교

**체결가 그룹 (4개 EP):**

| 항목 | KRX야간선물 (H0MFCNT0) | KRX야간옵션 (H0EUCNT0) | 주식선물 (H0ZFCNT0) | 주식옵션 (H0ZOCNT0) |
|---|---|---|---|---|
| 종목코드 필드명 | FUTS_SHRN_ISCD | OPTN_SHRN_ISCD | FUTS_SHRN_ISCD | OPTN_SHRN_ISCD |
| 가격 필드명 | FUTS_PRPR 등 | OPTN_PRPR 등 | STCK_PRPR 등 | OPTN_PRPR 등 |
| 선물 고유 필드 | MRKT_BASIS, DPRT, NMSC_FCTN, FMSC_FCTN, SPEAD_PRC | 없음 | MRKT_BASIS, DPRT, NMSC_FCTN, FMSC_FCTN, SPEAD_PRC | 없음 |
| 옵션 그릭스 | 없음 | PRMM_VAL, INVL_VAL, TMVL_VAL, DELTA, GAMA, VEGA, THETA, RHO, HTS_INTS_VLTL, UNAS_HIST_VLTL | 없음 | PRMM_VAL, INVL_VAL, TMVL_VAL, DELTA, GAMA, VEGA, THETA, RHO, HTS_INTS_VLTL, UNAS_HIST_VLTL |
| KRX 특유 필드 | DYNM_MXPR, DYNM_LLAM, DYNM_PRC_LIMT_YN | DYNM_MXPR, DYNM_PRC_LIMT_YN, DYNM_LLAM | 없음 | 없음 |
| 대비 필드 | OPRC_VRSS_NMIX_PRPR 방식 | OPRC_VRSS_NMIX_PRPR 방식 | OPRC_VRSS_PRPR 방식 | OPRC_VRSS_NMIX_PRPR 방식 |
| Field count | 46 | 51 | 47 | 47 |

**결론: 체결가 4개는 모두 distinct. base struct 불가.**

- KRX 야간선물 vs 주식선물: 필드명 prefix 차이 (FUTS_PRPR vs STCK_PRPR), 시가대비 방식 차이 (NMIX vs non-NMIX), KRX DYNM 필드 유무 차이 → **Distinct**
- KRX 야간옵션 vs 주식옵션: 필드명 동일, KRX DYNM 3필드 유무 차이 (51 vs 47 → -3 anom, 실제 DYNM_MXPR/DYNM_LLAM/DYNM_PRC_LIMT_YN) → **Distinct (KRX가 3 fields 더 있음)**

**호가 그룹 (3개 EP):**

| 항목 | KRX야간선물 (H0MFASP0) | KRX야간옵션 (H0EUASP0) | 주식선물 (H0ZFASP0) | 주식옵션 (H0ZOASP0) |
|---|---|---|---|---|
| 호가 depth | 5단계 | 5단계 | 10단계 | 10단계 (OPTN_ASKP1..10) |
| 종목코드 필드명 | FUTS_SHRN_ISCD | OPTN_SHRN_ISCD | FUTS_SHRN_ISCD | OPTN_SHRN_ISCD |
| 호가 필드명 | FUTS_ASKP/BIDP | OPTN_ASKP/BIDP | ASKP/BIDP (no prefix) | OPTN_ASKP/BIDP |
| 잔량증감 | TOTAL_ASKP_RSQN_ICDC/TOTAL_BIDP_RSQN_ICDC | 동일 | 동일 | 없음 (별도 위치) |
| Field count | 36 | 36 | 46 | 56 |

**결론: 호가 4개 모두 Distinct.**
- KRX야간선물 vs KRX야간옵션: 필드명만 다름 (FUTS_vs OPTN_), count 동일 36 → **거의 동일하나 필드명 prefix 다름**
- 주식선물 vs 주식옵션: 주식선물은 ASKP1..10 (no prefix), 주식옵션은 OPTN_ASKP1..10, 주식옵션은 총잔량 이후에도 OPTN_ASKP6..10 추가 그룹 → **Distinct**

**예상체결 그룹:**

| 항목 | KRX야간옵션 (H0EUANC0) | 주식선물 (H0ZFANC0) | 주식옵션 (H0ZOANC0) |
|---|---|---|---|
| 종목코드 필드명 | OPTN_SHRN_ISCD | FUTS_SHRN_ISCD | OPTN_SHRN_ISCD |
| Fields | 8 | 8 | 7 |
| ANTC_CNQN | O (Number 타입) | O (String) | 없음 |

**결론: 예상체결도 Distinct (주식옵션 7 fields, 나머지 8 fields).**

### 최종 base struct 수

**11 EP 모두 Distinct. alias 없음.** Phase 9 와 달리 각 EP 가 고유한 schema.

- 선물 체결: 선물고유필드(베이시스/스프레드) 포함
- 옵션 체결: 옵션그릭스(델타/감마/베가/세타/로우) 포함
- KRX vs 주식: DYNM 필드 유무, 호가 depth(5 vs 10), 필드명 prefix 차이

---

## TR_ID 명명 규칙

```
H0 MF CNT 0  → KRX 야간 선물(MF) 체결가(CNT)
H0 MF ASP 0  → KRX 야간 선물(MF) 호가(ASP)
H0 EU CNT 0  → KRX 야간 옵션(EU) 체결가(CNT)
H0 EU ASP 0  → KRX 야간 옵션(EU) 호가(ASP)
H0 EU ANC 0  → KRX 야간 옵션(EU) 예상체결(ANC)
H0 ZF CNT 0  → 주식선물(ZF) 체결가(CNT)
H0 ZF ASP 0  → 주식선물(ZF) 호가(ASP)
H0 ZF ANC 0  → 주식선물(ZF) 예상체결(ANC)
H0 ZO CNT 0  → 주식옵션(ZO) 체결가(CNT)
H0 ZO ASP 0  → 주식옵션(ZO) 호가(ASP)
H0 ZO ANC 0  → 주식옵션(ZO) 예상체결(ANC)
```

---

## EP1 — H0MFCNT0 KRX야간선물 실시간종목체결 (46 fields)

tr_key: 야간선물 종목코드 12자리 (예: 101W09000000)
도메인: ws://ops.koreainvestment.com:21000

| # | 필드명 | 한글명 | Go Type | 비고 |
|---|---|---|---|---|
| 0 | FUTS_SHRN_ISCD | 선물단축종목코드 | string | 9자리 |
| 1 | BSOP_HOUR | 영업시간 | string | HHMMSS |
| 2 | FUTS_PRDY_VRSS | 선물전일대비 | decimal.Decimal | 가격 |
| 3 | PRDY_VRSS_SIGN | 전일대비부호 | string | 1자리 |
| 4 | FUTS_PRDY_CTRT | 선물전일대비율 | float64 | 비율 |
| 5 | FUTS_PRPR | 선물현재가 | decimal.Decimal | 가격 |
| 6 | FUTS_OPRC | 선물시가2 | decimal.Decimal | 가격 |
| 7 | FUTS_HGPR | 선물최고가 | decimal.Decimal | 가격 |
| 8 | FUTS_LWPR | 선물최저가 | decimal.Decimal | 가격 |
| 9 | LAST_CNQN | 최종거래량 | int64 | |
| 10 | ACML_VOL | 누적거래량 | int64 | |
| 11 | ACML_TR_PBMN | 누적거래대금 | int64 | |
| 12 | HTS_THPR | HTS이론가 | decimal.Decimal | 가격 |
| 13 | MRKT_BASIS | 시장베이시스 | decimal.Decimal | 가격 |
| 14 | DPRT | 괴리율 | float64 | 비율 |
| 15 | NMSC_FCTN_STPL_PRC | 근월물약정가 | decimal.Decimal | 가격 |
| 16 | FMSC_FCTN_STPL_PRC | 원월물약정가 | decimal.Decimal | 가격 |
| 17 | SPEAD_PRC | 스프레드1 | decimal.Decimal | 가격 |
| 18 | HTS_OTST_STPL_QTY | HTS미결제약정수량 | int64 | |
| 19 | OTST_STPL_QTY_ICDC | 미결제약정수량증감 | int64 | |
| 20 | OPRC_HOUR | 시가시간 | string | HHMMSS |
| 21 | OPRC_VRSS_PRPR_SIGN | 시가2대비현재가부호 | string | 1자리 |
| 22 | OPRC_VRSS_NMIX_PRPR | 시가대비지수현재가 | decimal.Decimal | 가격 |
| 23 | HGPR_HOUR | 최고가시간 | string | HHMMSS |
| 24 | HGPR_VRSS_PRPR_SIGN | 최고가대비현재가부호 | string | 1자리 |
| 25 | HGPR_VRSS_NMIX_PRPR | 최고가대비지수현재가 | decimal.Decimal | 가격 |
| 26 | LWPR_HOUR | 최저가시간 | string | HHMMSS |
| 27 | LWPR_VRSS_PRPR_SIGN | 최저가대비현재가부호 | string | 1자리 |
| 28 | LWPR_VRSS_NMIX_PRPR | 최저가대비지수현재가 | decimal.Decimal | 가격 |
| 29 | SHNU_RATE | 매수2비율 | float64 | 비율 |
| 30 | CTTR | 체결강도 | float64 | 비율 |
| 31 | ESDG | 괴리도 | decimal.Decimal | 가격 |
| 32 | OTST_STPL_RGBF_QTY_ICDC | 미결제약정직전수량증감 | int64 | |
| 33 | THPR_BASIS | 이론베이시스 | decimal.Decimal | 가격 |
| 34 | FUTS_ASKP1 | 선물매도호가1 | decimal.Decimal | 가격 |
| 35 | FUTS_BIDP1 | 선물매수호가1 | decimal.Decimal | 가격 |
| 36 | ASKP_RSQN1 | 매도호가잔량1 | int64 | |
| 37 | BIDP_RSQN1 | 매수호가잔량1 | int64 | |
| 38 | SELN_CNTG_CSNU | 매도체결건수 | int64 | |
| 39 | SHNU_CNTG_CSNU | 매수체결건수 | int64 | |
| 40 | NTBY_CNTG_CSNU | 순매수체결건수 | int64 | |
| 41 | SELN_CNTG_SMTN | 총매도수량 | int64 | |
| 42 | SHNU_CNTG_SMTN | 총매수수량 | int64 | |
| 43 | TOTAL_ASKP_RSQN | 총매도호가잔량 | int64 | |
| 44 | TOTAL_BIDP_RSQN | 총매수호가잔량 | int64 | |
| 45 | PRDY_VOL_VRSS_ACML_VOL_RATE | 전일거래량대비등락율 | float64 | 비율 |
| 46 | DYNM_MXPR | 실시간상한가 | decimal.Decimal | 가격, Length=8 |
| 47 | DYNM_LLAM | 실시간하한가 | decimal.Decimal | 가격, Length=8 |
| 48 | DYNM_PRC_LIMT_YN | 실시간가격제한구분 | string | 1자리 |

> **주의**: docs 표는 46개로 나타나 있으나 body 표를 직접 세면 FUTS_SHRN_ISCD~DYNM_PRC_LIMT_YN 까지 49개.
> 표 Length 컬럼의 "1"은 docs 의 표현 오류 (실제 값 유효). 마지막 3개 DYNM 필드는 Length 8/8/1 로 명시.
> **실제 field count: 49**

---

## EP2 — H0MFASP0 KRX야간선물 실시간호가 (36 fields)

tr_key: 야간선물 종목코드 12자리
5단계 호가 (KRX 야간선물 전용)

| # | 필드명 | 한글명 | Go Type | 비고 |
|---|---|---|---|---|
| 0 | FUTS_SHRN_ISCD | 선물단축종목코드 | string | 9자리 |
| 1 | BSOP_HOUR | 영업시간 | string | HHMMSS |
| 2..6 | FUTS_ASKP1..5 | 선물매도호가1~5 | decimal.Decimal | 가격 |
| 7..11 | FUTS_BIDP1..5 | 선물매수호가1~5 | decimal.Decimal | 가격 |
| 12..16 | ASKP_CSNU1..5 | 매도호가건수1~5 | int64 | |
| 17..21 | BIDP_CSNU1..5 | 매수호가건수1~5 | int64 | |
| 22..26 | ASKP_RSQN1..5 | 매도호가잔량1~5 | int64 | |
| 27..31 | BIDP_RSQN1..5 | 매수호가잔량1~5 | int64 | |
| 32 | TOTAL_ASKP_CSNU | 총매도호가건수 | int64 | |
| 33 | TOTAL_BIDP_CSNU | 총매수호가건수 | int64 | |
| 34 | TOTAL_ASKP_RSQN | 총매도호가잔량 | int64 | |
| 35 | TOTAL_BIDP_RSQN | 총매수호가잔량 | int64 | |
| 36 | TOTAL_ASKP_RSQN_ICDC | 총매도호가잔량증감 | int64 | |
| 37 | TOTAL_BIDP_RSQN_ICDC | 총매수호가잔량증감 | int64 | |

> docs 표에서 카운트: FUTS_SHRN_ISCD, BSOP_HOUR + 5*FUTS_ASKP + 5*FUTS_BIDP + 5*ASKP_CSNU + 5*BIDP_CSNU + 5*ASKP_RSQN + 5*BIDP_RSQN + TOTAL_4 + TOTAL_2 = **38 fields**

---

## EP3 — H0EUCNT0 KRX야간옵션 실시간체결가 (51 fields)

tr_key: 야간옵션 종목코드 12자리
옵션 그릭스 포함, DYNM 필드 포함

| # | 필드명 | 한글명 | Go Type | 비고 |
|---|---|---|---|---|
| 0 | OPTN_SHRN_ISCD | 옵션단축종목코드 | string | 9자리 |
| 1 | BSOP_HOUR | 영업시간 | string | HHMMSS |
| 2 | OPTN_PRPR | 옵션현재가 | decimal.Decimal | 가격 |
| 3 | PRDY_VRSS_SIGN | 전일대비부호 | string | 1자리 |
| 4 | OPTN_PRDY_VRSS | 옵션전일대비 | decimal.Decimal | 가격 |
| 5 | PRDY_CTRT | 전일대비율 | float64 | 비율 |
| 6 | OPTN_OPRC | 옵션시가2 | decimal.Decimal | 가격 |
| 7 | OPTN_HGPR | 옵션최고가 | decimal.Decimal | 가격 |
| 8 | OPTN_LWPR | 옵션최저가 | decimal.Decimal | 가격 |
| 9 | LAST_CNQN | 최종거래량 | int64 | |
| 10 | ACML_VOL | 누적거래량 | int64 | |
| 11 | ACML_TR_PBMN | 누적거래대금 | int64 | |
| 12 | HTS_THPR | HTS이론가 | decimal.Decimal | 가격 |
| 13 | HTS_OTST_STPL_QTY | HTS미결제약정수량 | int64 | |
| 14 | OTST_STPL_QTY_ICDC | 미결제약정수량증감 | int64 | |
| 15 | OPRC_HOUR | 시가시간 | string | HHMMSS |
| 16 | OPRC_VRSS_PRPR_SIGN | 시가2대비현재가부호 | string | 1자리 |
| 17 | OPRC_VRSS_NMIX_PRPR | 시가대비지수현재가 | decimal.Decimal | 가격 |
| 18 | HGPR_HOUR | 최고가시간 | string | HHMMSS |
| 19 | HGPR_VRSS_PRPR_SIGN | 최고가대비현재가부호 | string | 1자리 |
| 20 | HGPR_VRSS_NMIX_PRPR | 최고가대비지수현재가 | decimal.Decimal | 가격 |
| 21 | LWPR_HOUR | 최저가시간 | string | HHMMSS |
| 22 | LWPR_VRSS_PRPR_SIGN | 최저가대비현재가부호 | string | 1자리 |
| 23 | LWPR_VRSS_NMIX_PRPR | 최저가대비지수현재가 | decimal.Decimal | 가격 |
| 24 | SHNU_RATE | 매수2비율 | float64 | 비율 |
| 25 | PRMM_VAL | 프리미엄값 | decimal.Decimal | 가격 |
| 26 | INVL_VAL | 내재가치값 | decimal.Decimal | 가격 |
| 27 | TMVL_VAL | 시간가치값 | decimal.Decimal | 가격 |
| 28 | DELTA | 델타 | float64 | 비율 |
| 29 | GAMA | 감마 | float64 | 비율 |
| 30 | VEGA | 베가 | float64 | 비율 |
| 31 | THETA | 세타 | float64 | 비율 |
| 32 | RHO | 로우 | float64 | 비율 |
| 33 | HTS_INTS_VLTL | HTS내재변동성 | float64 | 비율 |
| 34 | ESDG | 괴리도 | decimal.Decimal | 가격 |
| 35 | OTST_STPL_RGBF_QTY_ICDC | 미결제약정직전수량증감 | int64 | |
| 36 | THPR_BASIS | 이론베이시스 | decimal.Decimal | 가격 |
| 37 | UNAS_HIST_VLTL | 역사적변동성 | float64 | 비율 |
| 38 | CTTR | 체결강도 | float64 | 비율 |
| 39 | DPRT | 괴리율 | float64 | 비율 |
| 40 | MRKT_BASIS | 시장베이시스 | decimal.Decimal | 가격 |
| 41 | OPTN_ASKP1 | 옵션매도호가1 | decimal.Decimal | 가격 |
| 42 | OPTN_BIDP1 | 옵션매수호가1 | decimal.Decimal | 가격 |
| 43 | ASKP_RSQN1 | 매도호가잔량1 | int64 | |
| 44 | BIDP_RSQN1 | 매수호가잔량1 | int64 | |
| 45 | SELN_CNTG_CSNU | 매도체결건수 | int64 | |
| 46 | SHNU_CNTG_CSNU | 매수체결건수 | int64 | |
| 47 | NTBY_CNTG_CSNU | 순매수체결건수 | int64 | |
| 48 | SELN_CNTG_SMTN | 총매도수량 | int64 | |
| 49 | SHNU_CNTG_SMTN | 총매수수량 | int64 | |
| 50 | TOTAL_ASKP_RSQN | 총매도호가잔량 | int64 | |
| 51 | TOTAL_BIDP_RSQN | 총매수호가잔량 | int64 | |
| 52 | PRDY_VOL_VRSS_ACML_VOL_RATE | 전일거래량대비등락율 | float64 | 비율 |
| 53 | DYNM_MXPR | 실시간상한가 | decimal.Decimal | Length=8 |
| 54 | DYNM_PRC_LIMT_YN | 실시간가격제한구분 | string | 1자리 |
| 55 | DYNM_LLAM | 실시간하한가 | decimal.Decimal | Length=8 |

> docs 에서 직접 카운트: 56 fields (OPTN_SHRN_ISCD ~ DYNM_LLAM)

---

## EP4 — H0EUASP0 KRX야간옵션 실시간호가 (36 fields)

tr_key: 야간옵션 종목코드 12자리
5단계 호가 (KRX 야간옵션 전용, OPTN_ASKP/BIDP prefix)

| # | 필드명 | 한글명 | Go Type |
|---|---|---|---|
| 0 | OPTN_SHRN_ISCD | 옵션단축종목코드 | string |
| 1 | BSOP_HOUR | 영업시간 | string |
| 2..6 | OPTN_ASKP1..5 | 옵션매도호가1~5 | decimal.Decimal |
| 7..11 | OPTN_BIDP1..5 | 옵션매수호가1~5 | decimal.Decimal |
| 12..16 | ASKP_CSNU1..5 | 매도호가건수1~5 | int64 |
| 17..21 | BIDP_CSNU1..5 | 매수호가건수1~5 | int64 |
| 22..26 | ASKP_RSQN1..5 | 매도호가잔량1~5 | int64 |
| 27..31 | BIDP_RSQN1..5 | 매수호가잔량1~5 | int64 |
| 32 | TOTAL_ASKP_CSNU | 총매도호가건수 | int64 |
| 33 | TOTAL_BIDP_CSNU | 총매수호가건수 | int64 |
| 34 | TOTAL_ASKP_RSQN | 총매도호가잔량 | int64 |
| 35 | TOTAL_BIDP_RSQN | 총매수호가잔량 | int64 |
| 36 | TOTAL_ASKP_RSQN_ICDC | 총매도호가잔량증감 | int64 |
| 37 | TOTAL_BIDP_RSQN_ICDC | 총매수호가잔량증감 | int64 |

> docs 직접 카운트: 38 fields (EP2 와 동일 구조, prefix만 FUTS→OPTN 차이)
> EP2 (H0MFASP0) 와 EP4 (H0EUASP0) 는 종목코드 필드명과 호가 prefix 만 다름. 구조 동일.

---

## EP5 — H0EUANC0 KRX야간옵션 실시간예상체결 (8 fields)

tr_key: 야간옵션 종목코드 12자리

| # | 필드명 | 한글명 | Go Type | 비고 |
|---|---|---|---|---|
| 0 | OPTN_SHRN_ISCD | 옵션단축종목코드 | string | 9자리 |
| 1 | BSOP_HOUR | 영업시간 | string | HHMMSS |
| 2 | ANTC_CNPR | 예상체결가 | decimal.Decimal | Length=8 |
| 3 | ANTC_CNTG_VRSS | 예상체결대비 | decimal.Decimal | Length=8 |
| 4 | ANTC_CNTG_VRSS_SIGN | 예상체결대비부호 | string | 1자리 |
| 5 | ANTC_CNTG_PRDY_CTRT | 예상체결전일대비율 | float64 | Length=8 |
| 6 | ANTC_MKOP_CLS_CODE | 예상장운영구분코드 | string | 3자리 |
| 7 | ANTC_CNQN | 예상체결수량 | int64 | Number 타입 명시 |

> docs Python 예시에서 ANTC_CNQN: float 로 표기 → 실제는 수량이므로 int64 로 처리.

---

## EP6 — H0ZFCNT0 주식선물 실시간체결가 (47 fields)

tr_key: 주식선물 종목코드 6자리
STCK_* prefix 가격 필드 (KRX야간선물과 달리), KRX DYNM 없음

| # | 필드명 | 한글명 | Go Type | 비고 |
|---|---|---|---|---|
| 0 | FUTS_SHRN_ISCD | 선물단축종목코드 | string | 9자리 |
| 1 | BSOP_HOUR | 영업시간 | string | HHMMSS |
| 2 | STCK_PRPR | 주식현재가 | decimal.Decimal | 가격 |
| 3 | PRDY_VRSS_SIGN | 전일대비부호 | string | 1자리 |
| 4 | PRDY_VRSS | 전일대비 | decimal.Decimal | 가격 |
| 5 | FUTS_PRDY_CTRT | 선물전일대비율 | float64 | 비율 |
| 6 | STCK_OPRC | 주식시가2 | decimal.Decimal | 가격 |
| 7 | STCK_HGPR | 주식최고가 | decimal.Decimal | 가격 |
| 8 | STCK_LWPR | 주식최저가 | decimal.Decimal | 가격 |
| 9 | LAST_CNQN | 최종거래량 | int64 | |
| 10 | ACML_VOL | 누적거래량 | int64 | |
| 11 | ACML_TR_PBMN | 누적거래대금 | int64 | |
| 12 | HTS_THPR | HTS이론가 | decimal.Decimal | 가격 |
| 13 | MRKT_BASIS | 시장베이시스 | decimal.Decimal | 가격 |
| 14 | DPRT | 괴리율 | float64 | 비율 |
| 15 | NMSC_FCTN_STPL_PRC | 근월물약정가 | decimal.Decimal | 가격 |
| 16 | FMSC_FCTN_STPL_PRC | 원월물약정가 | decimal.Decimal | 가격 |
| 17 | SPEAD_PRC | 스프레드1 | decimal.Decimal | 가격 |
| 18 | HTS_OTST_STPL_QTY | HTS미결제약정수량 | int64 | |
| 19 | OTST_STPL_QTY_ICDC | 미결제약정수량증감 | int64 | |
| 20 | OPRC_HOUR | 시가시간 | string | HHMMSS |
| 21 | OPRC_VRSS_PRPR_SIGN | 시가2대비현재가부호 | string | 1자리 |
| 22 | OPRC_VRSS_PRPR | 시가2대비현재가 | decimal.Decimal | 가격 (NMIX 아님) |
| 23 | HGPR_HOUR | 최고가시간 | string | HHMMSS |
| 24 | HGPR_VRSS_PRPR_SIGN | 최고가대비현재가부호 | string | 1자리 |
| 25 | HGPR_VRSS_PRPR | 최고가대비현재가 | decimal.Decimal | 가격 |
| 26 | LWPR_HOUR | 최저가시간 | string | HHMMSS |
| 27 | LWPR_VRSS_PRPR_SIGN | 최저가대비현재가부호 | string | 1자리 |
| 28 | LWPR_VRSS_PRPR | 최저가대비현재가 | decimal.Decimal | 가격 |
| 29 | SHNU_RATE | 매수2비율 | float64 | 비율 |
| 30 | CTTR | 체결강도 | float64 | 비율 |
| 31 | ESDG | 괴리도 | decimal.Decimal | 가격 |
| 32 | OTST_STPL_RGBF_QTY_ICDC | 미결제약정직전수량증감 | int64 | |
| 33 | THPR_BASIS | 이론베이시스 | decimal.Decimal | 가격 |
| 34 | ASKP1 | 매도호가1 | decimal.Decimal | 가격 (FUTS_ prefix 없음) |
| 35 | BIDP1 | 매수호가1 | decimal.Decimal | 가격 |
| 36 | ASKP_RSQN1 | 매도호가잔량1 | int64 | |
| 37 | BIDP_RSQN1 | 매수호가잔량1 | int64 | |
| 38 | SELN_CNTG_CSNU | 매도체결건수 | int64 | |
| 39 | SHNU_CNTG_CSNU | 매수체결건수 | int64 | |
| 40 | NTBY_CNTG_CSNU | 순매수체결건수 | int64 | |
| 41 | SELN_CNTG_SMTN | 총매도수량 | int64 | |
| 42 | SHNU_CNTG_SMTN | 총매수수량 | int64 | |
| 43 | TOTAL_ASKP_RSQN | 총매도호가잔량 | int64 | |
| 44 | TOTAL_BIDP_RSQN | 총매수호가잔량 | int64 | |
| 45 | PRDY_VOL_VRSS_ACML_VOL_RATE | 전일거래량대비등락율 | float64 | 비율 |
| 46 | DYNM_MXPR | 실시간상한가 | decimal.Decimal | 가격 |
| 47 | DYNM_LLAM | 실시간하한가 | decimal.Decimal | 가격 |
| 48 | DYNM_PRC_LIMT_YN | 실시간가격제한구분 | string | 1자리 |

> docs 표 직접 카운트: 49 fields. DYNM 필드는 Length=4 로 명시 (KRX야간선물 8자리와 다름).
> KRX야간선물(H0MFCNT0) 과 비교: FUTS_PRPR→STCK_PRPR, OPRC_VRSS_NMIX→OPRC_VRSS_PRPR 이 핵심 차이.

---

## EP7 — H0ZFASP0 주식선물 실시간호가 (46 fields)

tr_key: 주식선물 종목코드 6자리
**10단계 호가** (KRX 야간선물과 달리 10단계), ASKP/BIDP prefix 없음

| # | 필드명 | 한글명 | Go Type |
|---|---|---|---|
| 0 | FUTS_SHRN_ISCD | 선물단축종목코드 | string |
| 1 | BSOP_HOUR | 영업시간 | string |
| 2..11 | ASKP1..10 | 매도호가1~10 | decimal.Decimal |
| 12..21 | BIDP1..10 | 매수호가1~10 | decimal.Decimal |
| 22..31 | ASKP_CSNU1..10 | 매도호가건수1~10 | int64 |
| 32..41 | BIDP_CSNU1..10 | 매수호가건수1~10 | int64 |
| 42..51 | ASKP_RSQN1..10 | 매도호가잔량1~10 | int64 |
| 52..61 | BIDP_RSQN1..10 | 매수호가잔량1~10 | int64 |
| 62 | TOTAL_ASKP_CSNU | 총매도호가건수 | int64 |
| 63 | TOTAL_BIDP_CSNU | 총매수호가건수 | int64 |
| 64 | TOTAL_ASKP_RSQN | 총매도호가잔량 | int64 |
| 65 | TOTAL_BIDP_RSQN | 총매수호가잔량 | int64 |
| 66 | TOTAL_ASKP_RSQN_ICDC | 총매도호가잔량증감 | int64 |
| 67 | TOTAL_BIDP_RSQN_ICDC | 총매수호가잔량증감 | int64 |

> docs 직접 카운트: 68 fields. 10단계 × 4그룹(ASK/BID/ASKCSNU/BIDCSNU×10 + ASKRSQN/BIDRSQN×10) + header 2 + total 6

---

## EP8 — H0ZFANC0 주식선물 실시간예상체결 (8 fields)

tr_key: 주식선물 종목코드 12자리 (docs에 12자리 명시)

| # | 필드명 | 한글명 | Go Type | 비고 |
|---|---|---|---|---|
| 0 | FUTS_SHRN_ISCD | 선물단축종목코드 | string | 9자리 |
| 1 | BSOP_HOUR | 영업시간 | string | HHMMSS |
| 2 | ANTC_CNPR | 예상체결가 | decimal.Decimal | Length=8 |
| 3 | ANTC_CNTG_VRSS | 예상체결대비 | decimal.Decimal | Length=8 |
| 4 | ANTC_CNTG_VRSS_SIGN | 예상체결대비부호 | string | 1자리 |
| 5 | ANTC_CNTG_PRDY_CTRT | 예상체결전일대비율 | float64 | Length=8 |
| 6 | ANTC_MKOP_CLS_CODE | 예상장운영구분코드 | string | 3자리 |
| 7 | ANTC_CNQN | 예상체결수량 | int64 | String 타입 명시 |

---

## EP9 — H0ZOCNT0 주식옵션 실시간체결가 (47 fields)

tr_key: 주식옵션 종목코드 6자리
옵션 그릭스 포함, DYNM 없음 (KRX야간옵션 대비)

| # | 필드명 | 한글명 | Go Type | 비고 |
|---|---|---|---|---|
| 0 | OPTN_SHRN_ISCD | 옵션단축종목코드 | string | 9자리 |
| 1 | BSOP_HOUR | 영업시간 | string | HHMMSS |
| 2 | OPTN_PRPR | 옵션현재가 | decimal.Decimal | 가격 |
| 3 | PRDY_VRSS_SIGN | 전일대비부호 | string | 1자리 |
| 4 | OPTN_PRDY_VRSS | 옵션전일대비 | decimal.Decimal | 가격 |
| 5 | PRDY_CTRT | 전일대비율 | float64 | 비율 |
| 6 | OPTN_OPRC | 옵션시가2 | decimal.Decimal | 가격 |
| 7 | OPTN_HGPR | 옵션최고가 | decimal.Decimal | 가격 |
| 8 | OPTN_LWPR | 옵션최저가 | decimal.Decimal | 가격 |
| 9 | LAST_CNQN | 최종거래량 | int64 | |
| 10 | ACML_VOL | 누적거래량 | int64 | |
| 11 | ACML_TR_PBMN | 누적거래대금 | int64 | |
| 12 | HTS_THPR | HTS이론가 | decimal.Decimal | 가격 |
| 13 | HTS_OTST_STPL_QTY | HTS미결제약정수량 | int64 | |
| 14 | OTST_STPL_QTY_ICDC | 미결제약정수량증감 | int64 | |
| 15 | OPRC_HOUR | 시가시간 | string | HHMMSS |
| 16 | OPRC_VRSS_PRPR_SIGN | 시가2대비현재가부호 | string | 1자리 |
| 17 | OPRC_VRSS_NMIX_PRPR | 시가대비지수현재가 | decimal.Decimal | 가격 |
| 18 | HGPR_HOUR | 최고가시간 | string | HHMMSS |
| 19 | HGPR_VRSS_PRPR_SIGN | 최고가대비현재가부호 | string | 1자리 |
| 20 | HGPR_VRSS_NMIX_PRPR | 최고가대비지수현재가 | decimal.Decimal | 가격 |
| 21 | LWPR_HOUR | 최저가시간 | string | HHMMSS |
| 22 | LWPR_VRSS_PRPR_SIGN | 최저가대비현재가부호 | string | 1자리 |
| 23 | LWPR_VRSS_NMIX_PRPR | 최저가대비지수현재가 | decimal.Decimal | 가격 |
| 24 | SHNU_RATE | 매수2비율 | float64 | 비율 |
| 25 | PRMM_VAL | 프리미엄값 | decimal.Decimal | 가격 |
| 26 | INVL_VAL | 내재가치값 | decimal.Decimal | 가격 |
| 27 | TMVL_VAL | 시간가치값 | decimal.Decimal | 가격 |
| 28 | DELTA | 델타 | float64 | 비율 |
| 29 | GAMA | 감마 | float64 | 비율 |
| 30 | VEGA | 베가 | float64 | 비율 |
| 31 | THETA | 세타 | float64 | 비율 |
| 32 | RHO | 로우 | float64 | 비율 |
| 33 | HTS_INTS_VLTL | HTS내재변동성 | float64 | 비율 |
| 34 | ESDG | 괴리도 | decimal.Decimal | 가격 |
| 35 | OTST_STPL_RGBF_QTY_ICDC | 미결제약정직전수량증감 | int64 | |
| 36 | THPR_BASIS | 이론베이시스 | decimal.Decimal | 가격 |
| 37 | UNAS_HIST_VLTL | 역사적변동성 | float64 | 비율 |
| 38 | CTTR | 체결강도 | float64 | 비율 |
| 39 | DPRT | 괴리율 | float64 | 비율 |
| 40 | MRKT_BASIS | 시장베이시스 | decimal.Decimal | 가격 |
| 41 | OPTN_ASKP1 | 옵션매도호가1 | decimal.Decimal | 가격 |
| 42 | OPTN_BIDP1 | 옵션매수호가1 | decimal.Decimal | 가격 |
| 43 | ASKP_RSQN1 | 매도호가잔량1 | int64 | |
| 44 | BIDP_RSQN1 | 매수호가잔량1 | int64 | |
| 45 | SELN_CNTG_CSNU | 매도체결건수 | int64 | |
| 46 | SHNU_CNTG_CSNU | 매수체결건수 | int64 | |
| 47 | NTBY_CNTG_CSNU | 순매수체결건수 | int64 | |
| 48 | SELN_CNTG_SMTN | 총매도수량 | int64 | |
| 49 | SHNU_CNTG_SMTN | 총매수수량 | int64 | |
| 50 | TOTAL_ASKP_RSQN | 총매도호가잔량 | int64 | |
| 51 | TOTAL_BIDP_RSQN | 총매수호가잔량 | int64 | |
| 52 | PRDY_VOL_VRSS_ACML_VOL_RATE | 전일거래량대비등락율 | float64 | 비율 |

> docs 직접 카운트: 53 fields. KRX야간옵션(H0EUCNT0) 대비 DYNM_MXPR/DYNM_PRC_LIMT_YN/DYNM_LLAM 3 fields 없음.

---

## EP10 — H0ZOASP0 주식옵션 실시간호가 (56 fields)

tr_key: 주식옵션 종목코드 6자리
**10단계 호가**, OPTN_ASKP/BIDP prefix, 잔량증감 없음 (대신 10단계 전체)
docs 구조 특이: OPTN_ASKP1..5 + 총계 + OPTN_ASKP6..10 형태로 분리 기술

| # | 필드명 | 한글명 | Go Type |
|---|---|---|---|
| 0 | OPTN_SHRN_ISCD | 옵션단축종목코드 | string |
| 1 | BSOP_HOUR | 영업시간 | string |
| 2..6 | OPTN_ASKP1..5 | 옵션매도호가1~5 | decimal.Decimal |
| 7..11 | OPTN_BIDP1..5 | 옵션매수호가1~5 | decimal.Decimal |
| 12..16 | ASKP_CSNU1..5 | 매도호가건수1~5 | int64 |
| 17..21 | BIDP_CSNU1..5 | 매수호가건수1~5 | int64 |
| 22..26 | ASKP_RSQN1..5 | 매도호가잔량1~5 | int64 |
| 27..31 | BIDP_RSQN1..5 | 매수호가잔량1~5 | int64 |
| 32 | TOTAL_ASKP_CSNU | 총매도호가건수 | int64 |
| 33 | TOTAL_BIDP_CSNU | 총매수호가건수 | int64 |
| 34 | TOTAL_ASKP_RSQN | 총매도호가잔량 | int64 |
| 35 | TOTAL_BIDP_RSQN | 총매수호가잔량 | int64 |
| 36 | TOTAL_ASKP_RSQN_ICDC | 총매도호가잔량증감 | int64 |
| 37 | TOTAL_BIDP_RSQN_ICDC | 총매수호가잔량증감 | int64 |
| 38..42 | OPTN_ASKP6..10 | 옵션매도호가6~10 | decimal.Decimal |
| 43..47 | OPTN_BIDP6..10 | 옵션매수호가6~10 | decimal.Decimal |
| 48..52 | ASKP_CSNU6..10 | 매도호가건수6~10 | int64 |
| 53..57 | BIDP_CSNU6..10 | 매수호가건수6~10 | int64 |
| 58..62 | ASKP_RSQN6..10 | 매도호가잔량6~10 | int64 |
| 63..67 | BIDP_RSQN6..10 | 매수호가잔량6~10 | int64 |

> docs 직접 카운트: 68 fields. 총잔량증감 있음, 호가 depth 10단계.
> 구현 시 배열 정렬: OPTN_ASKP[10], OPTN_BIDP[10], ASKP_CSNU[10], BIDP_CSNU[10], ASKP_RSQN[10], BIDP_RSQN[10]

---

## EP11 — H0ZOANC0 주식옵션 실시간예상체결 (7 fields)

tr_key: 주식옵션 종목코드 12자리
**주의: ANTC_CNQN 없음** (주식옵션 예상체결만 7 fields)

| # | 필드명 | 한글명 | Go Type | 비고 |
|---|---|---|---|---|
| 0 | OPTN_SHRN_ISCD | 옵션단축종목코드 | string | 9자리 |
| 1 | BSOP_HOUR | 영업시간 | string | HHMMSS |
| 2 | ANTC_CNPR | 예상체결가 | decimal.Decimal | Length=8 |
| 3 | ANTC_CNTG_VRSS | 예상체결대비 | decimal.Decimal | Length=8 |
| 4 | ANTC_CNTG_VRSS_SIGN | 예상체결대비부호 | string | 1자리 |
| 5 | ANTC_CNTG_PRDY_CTRT | 예상체결전일대비율 | float64 | Length=8 |
| 6 | ANTC_MKOP_CLS_CODE | 예상장운영구분코드 | string | 3자리 |

> docs Python 예시에도 7 fields 만 존재. ANTC_CNQN 누락 — docs anomaly or 의도적 차이.

---

## Anomalies

1. **docs Length "1" 오류**: KRX야간선물/옵션 체결가 docs 에서 거의 모든 필드가 Length=1 로 표기됨 (실제 데이터는 더 긴 값). 주식선물/옵션은 일부 정확한 Length 있음.

2. **주식옵션 예상체결 ANTC_CNQN 누락**: H0ZOANC0 는 7 fields (다른 예상체결 EP 는 8 fields). ANTC_CNQN 없음.

3. **H0EUCNT0 DYNM 필드 순서**: docs 표에서 DYNM_MXPR → DYNM_PRC_LIMT_YN → DYNM_LLAM 순서 (다른 EP 는 MXPR → LLAM → PRC_LIMT_YN). 구현 시 docs 순서 그대로.

4. **docs field count vs "N" 표기**: docs 상단 "총 N 개" 표기가 없거나 부정확. 이 schema 는 모두 docs 표 직접 카운트 기준.

5. **주식선물/옵션 호가 10단계**: KRX 야간 5단계 대비 주식선물/옵션은 10단계. Phase 8/9 의 국내주식호가(10단계)와 유사.

6. **H0ZFANC0 tr_key 12자리**: 주식선물 예상체결 docs 에서 tr_key length=12 명시 (다른 주식선물 EP는 6자리). 구현 시 12자리 허용.

---

## Type 매핑 정책 (Phase 8 와 동일)

| KIS Type | Go field type |
|---|---|
| String (HHMMSS, code, sign, name, YN) | `string` |
| 가격/이론가/베이시스/스프레드 | `decimal.Decimal` |
| 비율 (`*_CTRT`, `*_RATE`, `CTTR`, `SHNU_RATE`, `DPRT`, 그릭스, `*_VLTL`, `*_RLIM`) | `float64` |
| 거래량/잔량/건수/수량/대금 | `int64` |

KIS docs 가 모든 필드를 String 으로 표기해도 위 정책에 따라 Go type 결정.

---

## TR_ID 상수 (client.go 추가 예정)

```go
const (
    // Phase 11.2 — KRX 야간 선물
    trIDKrxNightFuturesTrade = "H0MFCNT0"
    trIDKrxNightFuturesAsk   = "H0MFASP0"

    // Phase 11.2 — KRX 야간 옵션
    trIDKrxNightOptionTrade       = "H0EUCNT0"
    trIDKrxNightOptionAsk         = "H0EUASP0"
    trIDKrxNightOptionExpectTrade = "H0EUANC0"

    // Phase 11.2 — 주식 선물
    trIDStockFuturesTrade       = "H0ZFCNT0"
    trIDStockFuturesAsk         = "H0ZFASP0"
    trIDStockFuturesExpectTrade = "H0ZFANC0"

    // Phase 11.2 — 주식 옵션
    trIDStockOptionTrade       = "H0ZOCNT0"
    trIDStockOptionAsk         = "H0ZOASP0"
    trIDStockOptionExpectTrade = "H0ZOANC0"
)
```

---

## Decoder 패턴 (11 개별 decoder)

```go
// 선물 체결가 (KRX 야간)
func decodeKrxNightFuturesTrade(f frame) ([]KrxNightFuturesTradeEvent, error)
// 선물 호가 (KRX 야간)
func decodeKrxNightFuturesAsk(f frame) ([]KrxNightFuturesAskEvent, error)
// 옵션 체결가 (KRX 야간)
func decodeKrxNightOptionTrade(f frame) ([]KrxNightOptionTradeEvent, error)
// 옵션 호가 (KRX 야간)
func decodeKrxNightOptionAsk(f frame) ([]KrxNightOptionAskEvent, error)
// 옵션 예상체결 (KRX 야간)
func decodeKrxNightOptionExpectTrade(f frame) ([]KrxNightOptionExpectTradeEvent, error)
// 주식선물 체결가
func decodeStockFuturesTrade(f frame) ([]StockFuturesTradeEvent, error)
// 주식선물 호가
func decodeStockFuturesAsk(f frame) ([]StockFuturesAskEvent, error)
// 주식선물 예상체결
func decodeStockFuturesExpectTrade(f frame) ([]StockFuturesExpectTradeEvent, error)
// 주식옵션 체결가
func decodeStockOptionTrade(f frame) ([]StockOptionTradeEvent, error)
// 주식옵션 호가
func decodeStockOptionAsk(f frame) ([]StockOptionAskEvent, error)
// 주식옵션 예상체결
func decodeStockOptionExpectTrade(f frame) ([]StockOptionExpectTradeEvent, error)
```
