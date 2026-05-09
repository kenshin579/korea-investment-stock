# Phase 11.3 국내선물옵션 실시간 6 EP — Schema Reference

> docs analyzer 결과 (2026-05-09). Phase 11.3 Event/Decoder/Client 의 source of truth.
> docs 위치: `docs/api/국내선물옵션/<파일>.md`

---

## EP Matrix

| # | TR_ID | 한글명 | Fields (docs count) | 모의 | tr_key 형식 | 종목코드 길이 |
|---|---|---|---|---|---|---|
| 1 | H0IFCNT0 | 지수선물 실시간체결가 | 50 | 미지원 | 지수선물 종목코드 | 6자리 |
| 2 | H0IFASP0 | 지수선물 실시간호가 | 38 | 미지원 | 지수선물 종목코드 | 6자리 |
| 3 | H0IOCNT0 | 지수옵션 실시간체결가 | 58 | 미지원 | 지수옵션 종목코드 | 6자리 |
| 4 | H0IOASP0 | 지수옵션 실시간호가 | 38 | 미지원 | 지수옵션 종목코드 | 6자리 |
| 5 | H0CFCNT0 | 상품선물 실시간체결가 | 50 | 미지원 | 상품선물 종목코드 | 6자리 |
| 6 | H0CFASP0 | 상품선물 실시간호가 | 38 | 미지원 | 상품선물 종목코드 | 6자리 |

**전체 모의투자 미지원.** WebSocket domain: `ws://ops.koreainvestment.com:21000`

---

## Alias 분석 결과

### 체결가 그룹

| 항목 | H0IFCNT0 (지수선물) | H0CFCNT0 (상품선물) | H0MFCNT0 (KRX야간선물 Phase11.2) | H0ZFCNT0 (주식선물 Phase11.2) |
|---|---|---|---|---|
| 종목코드 필드명 | FUTS_SHRN_ISCD | FUTS_SHRN_ISCD | FUTS_SHRN_ISCD | FUTS_SHRN_ISCD |
| 가격 prefix | FUTS_PRPR 등 | FUTS_PRPR 등 | FUTS_PRPR 등 | STCK_PRPR 등 |
| 시가대비 필드 | OPRC_VRSS_NMIX_PRPR | OPRC_VRSS_NMIX_PRPR | OPRC_VRSS_NMIX_PRPR | OPRC_VRSS_PRPR (non-NMIX) |
| DSCS_BLTR_ACML_QTY | O (협의대량거래량) | O (협의대량거래량) | X | X |
| DYNM 필드 | O (MXPR/LLAM/PRC_LIMT_YN) | O (MXPR/LLAM/PRC_LIMT_YN) | O | X |
| Field count | 50 | 50 | 49 | 49 |

**결론:**
- **H0IFCNT0 == H0CFCNT0**: field 명/순서/Type 완전 동일 → **alias 가능** (지수선물과 상품선물 동일 schema)
- **H0IFCNT0 != H0MFCNT0**: DSCS_BLTR_ACML_QTY 1 field 추가 → **Distinct**
- **H0IFCNT0 != H0ZFCNT0**: STCK_PRPR vs FUTS_PRPR, OPRC_VRSS_NMIX vs non-NMIX, DSCS_BLTR_ACML_QTY → **Distinct**

### 호가 그룹

| 항목 | H0IFASP0 (지수선물) | H0CFASP0 (상품선물) | H0IOASP0 (지수옵션) | H0MFASP0 (KRX야간선물 Phase11.2) |
|---|---|---|---|---|
| 종목코드 필드명 | FUTS_SHRN_ISCD | FUTS_SHRN_ISCD | OPTN_SHRN_ISCD | FUTS_SHRN_ISCD |
| 호가 prefix | FUTS_ASKP/BIDP | FUTS_ASKP/BIDP | OPTN_ASKP/BIDP | FUTS_ASKP/BIDP |
| 호가 depth | 5단계 | 5단계 | 5단계 | 5단계 |
| 잔량증감 | O | O | O | O |
| Field count | 38 | 38 | 38 | 38 |

**결론:**
- **H0IFASP0 == H0CFASP0**: field 명/순서/Type 완전 동일 → **alias 가능**
- **H0IFASP0 == H0MFASP0 (Phase 11.2)**: 완전 동일 → **alias 가능** (3개 모두 FUTS_ prefix 5단계)
- **H0IOASP0**: 종목코드 OPTN_, 호가 OPTN_ASKP/BIDP prefix → 구조 동일하나 필드명 다름 → **Distinct**

### 지수옵션 체결가

| 항목 | H0IOCNT0 (지수옵션) | H0ZOCNT0 (주식옵션 Phase11.2) | H0EUCNT0 (KRX야간옵션 Phase11.2) |
|---|---|---|---|
| 종목코드 필드명 | OPTN_SHRN_ISCD | OPTN_SHRN_ISCD | OPTN_SHRN_ISCD |
| 옵션 그릭스 | O | O | O |
| AVRG_VLTL | O (평균변동성) | X | X |
| DSCS_LRQN_VOL | O (협의대량누적거래량) | X | X |
| DYNM 필드 | O (MXPR/LLAM/PRC_LIMT_YN) | X | O |
| Field count | 58 | 53 | 56 |

**결론: H0IOCNT0 는 ZOCNT0/EUCNT0 모두와 Distinct. AVRG_VLTL + DSCS_LRQN_VOL + DYNM 3필드 추가.**

### 최종 base struct 수

**Phase 11.3 6 EP:**
- **2 distinct base types**: IndexFuturesTrade (IF/CF 공통), IndexFuturesAsk (IF/CF 공통)
- **2 distinct types**: IndexOptionTrade (IO 단독), IndexOptionAsk (IO 단독)
- **alias 결정**: H0CFCNT0 = alias of H0IFCNT0, H0CFASP0 = alias of H0IFASP0

구체적:
- `H0IFCNT0` → `IndexFuturesTradeEvent` (base)
- `H0CFCNT0` → `CommodityFuturesTradeEvent = IndexFuturesTradeEvent` (alias)
- `H0IFASP0` → `IndexFuturesAskEvent` (base, MF와 동일 schema 이지만 별도 type)
- `H0CFASP0` → `CommodityFuturesAskEvent = IndexFuturesAskEvent` (alias)
- `H0IOCNT0` → `IndexOptionTradeEvent` (distinct, 58 fields)
- `H0IOASP0` → `IndexOptionAskEvent` (distinct, OPTN_ prefix)

> **MF와 IF 호가 schema 동일**: 구현 시 공통 base struct 검토 가능하나, 도메인 명확성을 위해 별도 type 권장.

---

## TR_ID 명명 규칙

```
H0 IF CNT 0  → 지수선물(IF) 체결가(CNT)
H0 IF ASP 0  → 지수선물(IF) 호가(ASP)
H0 IO CNT 0  → 지수옵션(IO) 체결가(CNT)
H0 IO ASP 0  → 지수옵션(IO) 호가(ASP)
H0 CF CNT 0  → 상품선물(CF) 체결가(CNT)
H0 CF ASP 0  → 상품선물(CF) 호가(ASP)
```

---

## EP1 — H0IFCNT0 지수선물 실시간체결가 (50 fields)

tr_key: 지수선물 종목코드 6자리 (예: 101S12)
도메인: ws://ops.koreainvestment.com:21000
FUTS_ prefix 가격 필드, OPRC_VRSS_NMIX_PRPR (NMIX 방식), DSCS_BLTR_ACML_QTY 포함.

| # | 필드명 | 한글명 | Go Type | 비고 |
|---|---|---|---|---|
| 0 | FUTS_SHRN_ISCD | 선물 단축 종목코드 | string | 종목코드 |
| 1 | BSOP_HOUR | 영업 시간 | string | HHMMSS |
| 2 | FUTS_PRDY_VRSS | 선물 전일 대비 | decimal.Decimal | 가격 |
| 3 | PRDY_VRSS_SIGN | 전일 대비 부호 | string | 1자리 |
| 4 | FUTS_PRDY_CTRT | 선물 전일 대비율 | float64 | 비율 |
| 5 | FUTS_PRPR | 선물 현재가 | decimal.Decimal | 가격 |
| 6 | FUTS_OPRC | 선물 시가2 | decimal.Decimal | 가격 |
| 7 | FUTS_HGPR | 선물 최고가 | decimal.Decimal | 가격 |
| 8 | FUTS_LWPR | 선물 최저가 | decimal.Decimal | 가격 |
| 9 | LAST_CNQN | 최종 거래량 | int64 | 체결량 |
| 10 | ACML_VOL | 누적 거래량 | int64 | |
| 11 | ACML_TR_PBMN | 누적 거래 대금 | int64 | |
| 12 | HTS_THPR | HTS 이론가 | decimal.Decimal | 가격 |
| 13 | MRKT_BASIS | 시장 베이시스 | decimal.Decimal | 가격 |
| 14 | DPRT | 괴리율 | float64 | 비율 |
| 15 | NMSC_FCTN_STPL_PRC | 근월물 약정가 | decimal.Decimal | 가격 |
| 16 | FMSC_FCTN_STPL_PRC | 원월물 약정가 | decimal.Decimal | 가격 |
| 17 | SPEAD_PRC | 스프레드1 | decimal.Decimal | 가격 |
| 18 | HTS_OTST_STPL_QTY | HTS 미결제 약정 수량 | int64 | |
| 19 | OTST_STPL_QTY_ICDC | 미결제 약정 수량 증감 | int64 | |
| 20 | OPRC_HOUR | 시가 시간 | string | HHMMSS |
| 21 | OPRC_VRSS_PRPR_SIGN | 시가2 대비 현재가 부호 | string | 1자리 |
| 22 | OPRC_VRSS_NMIX_PRPR | 시가 대비 지수 현재가 | decimal.Decimal | 가격 (NMIX) |
| 23 | HGPR_HOUR | 최고가 시간 | string | HHMMSS |
| 24 | HGPR_VRSS_PRPR_SIGN | 최고가 대비 현재가 부호 | string | 1자리 |
| 25 | HGPR_VRSS_NMIX_PRPR | 최고가 대비 지수 현재가 | decimal.Decimal | 가격 (NMIX) |
| 26 | LWPR_HOUR | 최저가 시간 | string | HHMMSS |
| 27 | LWPR_VRSS_PRPR_SIGN | 최저가 대비 현재가 부호 | string | 1자리 |
| 28 | LWPR_VRSS_NMIX_PRPR | 최저가 대비 지수 현재가 | decimal.Decimal | 가격 (NMIX) |
| 29 | SHNU_RATE | 매수2 비율 | float64 | 비율 |
| 30 | CTTR | 체결강도 | float64 | 비율 |
| 31 | ESDG | 괴리도 | decimal.Decimal | 가격 |
| 32 | OTST_STPL_RGBF_QTY_ICDC | 미결제 약정 직전 수량 증감 | int64 | |
| 33 | THPR_BASIS | 이론 베이시스 | decimal.Decimal | 가격 |
| 34 | FUTS_ASKP1 | 선물 매도호가1 | decimal.Decimal | 가격 |
| 35 | FUTS_BIDP1 | 선물 매수호가1 | decimal.Decimal | 가격 |
| 36 | ASKP_RSQN1 | 매도호가 잔량1 | int64 | |
| 37 | BIDP_RSQN1 | 매수호가 잔량1 | int64 | |
| 38 | SELN_CNTG_CSNU | 매도 체결 건수 | int64 | |
| 39 | SHNU_CNTG_CSNU | 매수 체결 건수 | int64 | |
| 40 | NTBY_CNTG_CSNU | 순매수 체결 건수 | int64 | |
| 41 | SELN_CNTG_SMTN | 총 매도 수량 | int64 | |
| 42 | SHNU_CNTG_SMTN | 총 매수 수량 | int64 | |
| 43 | TOTAL_ASKP_RSQN | 총 매도호가 잔량 | int64 | |
| 44 | TOTAL_BIDP_RSQN | 총 매수호가 잔량 | int64 | |
| 45 | PRDY_VOL_VRSS_ACML_VOL_RATE | 전일 거래량 대비 등락율 | float64 | 비율 |
| 46 | DSCS_BLTR_ACML_QTY | 협의 대량 거래량 | int64 | H0MFCNT0 없음 |
| 47 | DYNM_MXPR | 실시간상한가 | decimal.Decimal | 가격 |
| 48 | DYNM_LLAM | 실시간하한가 | decimal.Decimal | 가격 |
| 49 | DYNM_PRC_LIMT_YN | 실시간가격제한구분 | string | 1자리 |

> docs 직접 카운트: 50 fields. H0MFCNT0 (49) 대비 DSCS_BLTR_ACML_QTY 1 field 추가.
> H0CFCNT0 와 field 명/순서/Type 완전 동일 → alias 패턴 적용.

---

## EP2 — H0IFASP0 지수선물 실시간호가 (38 fields)

tr_key: 지수선물 종목코드 6자리
**5단계 호가**, FUTS_ASKP/BIDP prefix
H0MFASP0 (Phase 11.2 KRX야간선물) 와 field 명/순서/Type 완전 동일.
H0CFASP0 와도 완전 동일 → alias 가능.

| # | 필드명 | 한글명 | Go Type |
|---|---|---|---|
| 0 | FUTS_SHRN_ISCD | 선물 단축 종목코드 | string |
| 1 | BSOP_HOUR | 영업 시간 | string |
| 2..6 | FUTS_ASKP1..5 | 선물 매도호가1~5 | decimal.Decimal |
| 7..11 | FUTS_BIDP1..5 | 선물 매수호가1~5 | decimal.Decimal |
| 12..16 | ASKP_CSNU1..5 | 매도호가 건수1~5 | int64 |
| 17..21 | BIDP_CSNU1..5 | 매수호가 건수1~5 | int64 |
| 22..26 | ASKP_RSQN1..5 | 매도호가 잔량1~5 | int64 |
| 27..31 | BIDP_RSQN1..5 | 매수호가 잔량1~5 | int64 |
| 32 | TOTAL_ASKP_CSNU | 총 매도호가 건수 | int64 |
| 33 | TOTAL_BIDP_CSNU | 총 매수호가 건수 | int64 |
| 34 | TOTAL_ASKP_RSQN | 총 매도호가 잔량 | int64 |
| 35 | TOTAL_BIDP_RSQN | 총 매수호가 잔량 | int64 |
| 36 | TOTAL_ASKP_RSQN_ICDC | 총 매도호가 잔량 증감 | int64 |
| 37 | TOTAL_BIDP_RSQN_ICDC | 총 매수호가 잔량 증감 | int64 |

> docs 직접 카운트: 38 fields. Phase 11.2 H0MFASP0 와 완전 동일 schema.

---

## EP3 — H0IOCNT0 지수옵션 실시간체결가 (58 fields)

tr_key: 지수옵션 종목코드 6자리 (예: 201S11305)
옵션 그릭스 포함, DYNM 필드 포함, AVRG_VLTL + DSCS_LRQN_VOL 추가.
Phase 11.2 주식옵션(53) / KRX야간옵션(56) 과 모두 Distinct.

| # | 필드명 | 한글명 | Go Type | 비고 |
|---|---|---|---|---|
| 0 | OPTN_SHRN_ISCD | 옵션 단축 종목코드 | string | 종목코드 |
| 1 | BSOP_HOUR | 영업 시간 | string | HHMMSS |
| 2 | OPTN_PRPR | 옵션 현재가 | decimal.Decimal | 가격 |
| 3 | PRDY_VRSS_SIGN | 전일 대비 부호 | string | 1자리 |
| 4 | OPTN_PRDY_VRSS | 옵션 전일 대비 | decimal.Decimal | 가격 |
| 5 | PRDY_CTRT | 전일 대비율 | float64 | 비율 |
| 6 | OPTN_OPRC | 옵션 시가2 | decimal.Decimal | 가격 |
| 7 | OPTN_HGPR | 옵션 최고가 | decimal.Decimal | 가격 |
| 8 | OPTN_LWPR | 옵션 최저가 | decimal.Decimal | 가격 |
| 9 | LAST_CNQN | 최종 거래량 | int64 | |
| 10 | ACML_VOL | 누적 거래량 | int64 | |
| 11 | ACML_TR_PBMN | 누적 거래 대금 | int64 | |
| 12 | HTS_THPR | HTS 이론가 | decimal.Decimal | 가격 |
| 13 | HTS_OTST_STPL_QTY | HTS 미결제 약정 수량 | int64 | |
| 14 | OTST_STPL_QTY_ICDC | 미결제 약정 수량 증감 | int64 | |
| 15 | OPRC_HOUR | 시가 시간 | string | HHMMSS |
| 16 | OPRC_VRSS_PRPR_SIGN | 시가2 대비 현재가 부호 | string | 1자리 |
| 17 | OPRC_VRSS_NMIX_PRPR | 시가 대비 지수 현재가 | decimal.Decimal | 가격 (NMIX) |
| 18 | HGPR_HOUR | 최고가 시간 | string | HHMMSS |
| 19 | HGPR_VRSS_PRPR_SIGN | 최고가 대비 현재가 부호 | string | 1자리 |
| 20 | HGPR_VRSS_NMIX_PRPR | 최고가 대비 지수 현재가 | decimal.Decimal | 가격 (NMIX) |
| 21 | LWPR_HOUR | 최저가 시간 | string | HHMMSS |
| 22 | LWPR_VRSS_PRPR_SIGN | 최저가 대비 현재가 부호 | string | 1자리 |
| 23 | LWPR_VRSS_NMIX_PRPR | 최저가 대비 지수 현재가 | decimal.Decimal | 가격 (NMIX) |
| 24 | SHNU_RATE | 매수2 비율 | float64 | 비율 |
| 25 | PRMM_VAL | 프리미엄 값 | decimal.Decimal | 가격 |
| 26 | INVL_VAL | 내재가치 값 | decimal.Decimal | 가격 |
| 27 | TMVL_VAL | 시간가치 값 | decimal.Decimal | 가격 |
| 28 | DELTA | 델타 | float64 | 그릭스 |
| 29 | GAMA | 감마 | float64 | 그릭스 |
| 30 | VEGA | 베가 | float64 | 그릭스 |
| 31 | THETA | 세타 | float64 | 그릭스 |
| 32 | RHO | 로우 | float64 | 그릭스 |
| 33 | HTS_INTS_VLTL | HTS 내재 변동성 | float64 | 비율 |
| 34 | ESDG | 괴리도 | decimal.Decimal | 가격 |
| 35 | OTST_STPL_RGBF_QTY_ICDC | 미결제 약정 직전 수량 증감 | int64 | |
| 36 | THPR_BASIS | 이론 베이시스 | decimal.Decimal | 가격 |
| 37 | UNAS_HIST_VLTL | 역사적변동성 | float64 | 비율 |
| 38 | CTTR | 체결강도 | float64 | 비율 |
| 39 | DPRT | 괴리율 | float64 | 비율 |
| 40 | MRKT_BASIS | 시장 베이시스 | decimal.Decimal | 가격 |
| 41 | OPTN_ASKP1 | 옵션 매도호가1 | decimal.Decimal | 가격 |
| 42 | OPTN_BIDP1 | 옵션 매수호가1 | decimal.Decimal | 가격 |
| 43 | ASKP_RSQN1 | 매도호가 잔량1 | int64 | |
| 44 | BIDP_RSQN1 | 매수호가 잔량1 | int64 | |
| 45 | SELN_CNTG_CSNU | 매도 체결 건수 | int64 | |
| 46 | SHNU_CNTG_CSNU | 매수 체결 건수 | int64 | |
| 47 | NTBY_CNTG_CSNU | 순매수 체결 건수 | int64 | |
| 48 | SELN_CNTG_SMTN | 총 매도 수량 | int64 | |
| 49 | SHNU_CNTG_SMTN | 총 매수 수량 | int64 | |
| 50 | TOTAL_ASKP_RSQN | 총 매도호가 잔량 | int64 | |
| 51 | TOTAL_BIDP_RSQN | 총 매수호가 잔량 | int64 | |
| 52 | PRDY_VOL_VRSS_ACML_VOL_RATE | 전일 거래량 대비 등락율 | float64 | 비율 |
| 53 | AVRG_VLTL | 평균 변동성 | float64 | 비율, ZOCNT0/EUCNT0 없음 |
| 54 | DSCS_LRQN_VOL | 협의대량누적 거래량 | int64 | ZOCNT0/EUCNT0 없음 |
| 55 | DYNM_MXPR | 실시간상한가 | decimal.Decimal | 가격 |
| 56 | DYNM_LLAM | 실시간하한가 | decimal.Decimal | 가격 |
| 57 | DYNM_PRC_LIMT_YN | 실시간가격제한구분 | string | 1자리 |

> docs 직접 카운트: 58 fields. ZOCNT0(53) 대비 +5: AVRG_VLTL, DSCS_LRQN_VOL, DYNM_MXPR, DYNM_LLAM, DYNM_PRC_LIMT_YN.
> EUCNT0(56) 대비 +2: AVRG_VLTL, DSCS_LRQN_VOL. DYNM 순서: DYNM_MXPR→DYNM_LLAM→DYNM_PRC_LIMT_YN (EUCNT0 와 다름).

---

## EP4 — H0IOASP0 지수옵션 실시간호가 (38 fields)

tr_key: 지수옵션 종목코드 6자리
**5단계 호가**, OPTN_ASKP/BIDP prefix (IFASP0 의 FUTS_ 와 다름)
구조는 IFASP0/CFASP0 와 동일하나 종목코드 필드명(OPTN_SHRN_ISCD), 호가 prefix(OPTN_) 차이.

| # | 필드명 | 한글명 | Go Type |
|---|---|---|---|
| 0 | OPTN_SHRN_ISCD | 옵션 단축 종목코드 | string |
| 1 | BSOP_HOUR | 영업 시간 | string |
| 2..6 | OPTN_ASKP1..5 | 옵션 매도호가1~5 | decimal.Decimal |
| 7..11 | OPTN_BIDP1..5 | 옵션 매수호가1~5 | decimal.Decimal |
| 12..16 | ASKP_CSNU1..5 | 매도호가 건수1~5 | int64 |
| 17..21 | BIDP_CSNU1..5 | 매수호가 건수1~5 | int64 |
| 22..26 | ASKP_RSQN1..5 | 매도호가 잔량1~5 | int64 |
| 27..31 | BIDP_RSQN1..5 | 매수호가 잔량1~5 | int64 |
| 32 | TOTAL_ASKP_CSNU | 총 매도호가 건수 | int64 |
| 33 | TOTAL_BIDP_CSNU | 총 매수호가 건수 | int64 |
| 34 | TOTAL_ASKP_RSQN | 총 매도호가 잔량 | int64 |
| 35 | TOTAL_BIDP_RSQN | 총 매수호가 잔량 | int64 |
| 36 | TOTAL_ASKP_RSQN_ICDC | 총 매도호가 잔량 증감 | int64 |
| 37 | TOTAL_BIDP_RSQN_ICDC | 총 매수호가 잔량 증감 | int64 |

> docs 직접 카운트: 38 fields. H0EUASP0 (Phase 11.2 KRX야간옵션) 와 완전 동일.
> IFASP0/CFASP0 와 prefix 만 다름 (OPTN_ vs FUTS_) → Distinct type 필요.

---

## EP5 — H0CFCNT0 상품선물 실시간체결가 (50 fields)

tr_key: 상품선물 종목코드 6자리
**H0IFCNT0 와 field 명/순서/Type 완전 동일** → alias 패턴 적용.
종목코드 필드명: FUTS_SHRN_ISCD (동일).

> alias: `type CommodityFuturesTradeEvent = IndexFuturesTradeEvent`
> 구현 시 별도 decoder 필요 없음.

---

## EP6 — H0CFASP0 상품선물 실시간호가 (38 fields)

tr_key: 상품선물 종목코드 6자리
**H0IFASP0 와 field 명/순서/Type 완전 동일** → alias 패턴 적용.

> alias: `type CommodityFuturesAskEvent = IndexFuturesAskEvent`

---

## Anomalies

1. **docs Length 비일관성**: H0CFCNT0 docs 에서 Length 값이 1/4/6/8/9 등 혼재. H0IFCNT0 보다 Length 컬럼이 더 불규칙. 실제 데이터 기준으로 type 결정.

2. **H0IOCNT0 DYNM 순서**: docs 에서 DYNM_MXPR → DYNM_LLAM → DYNM_PRC_LIMT_YN 순서. H0EUCNT0 는 DYNM_MXPR → DYNM_PRC_LIMT_YN → DYNM_LLAM 순서로 다름. 이 schema 는 docs 순서 그대로 (MXPR/LLAM/PRC_LIMT_YN).

3. **H0IFCNT0 DSCS_BLTR_ACML_QTY vs H0IOCNT0 DSCS_LRQN_VOL**: 두 EP 모두 "협의대량거래량" 성격이나 필드명이 다름. 지수선물 = DSCS_BLTR_ACML_QTY, 지수옵션 = DSCS_LRQN_VOL. 두 필드 모두 int64 처리.

4. **H0IOCNT0 AVRG_VLTL (평균변동성)**: ZOCNT0/EUCNT0 에 없는 신규 필드. float64 처리.

5. **tr_key 예시**: H0IFCNT0 docs tr_key 예시 "101S12" (6자리), H0IOCNT0 "201S11305" (9자리) — docs 에는 6자리 명시이나 실제 옵션코드는 더 길 수 있음. tr_key 길이 제한 없이 허용.

---

## Type 매핑 정책 (Phase 8 와 동일)

| KIS Type | Go field type |
|---|---|
| String (HHMMSS, code, sign, name, YN) | `string` |
| 가격/이론가/베이시스/스프레드/호가 | `decimal.Decimal` |
| 비율 (`*_CTRT`, `*_RATE`, `CTTR`, `SHNU_RATE`, `DPRT`, 그릭스, `*_VLTL`, `AVRG_VLTL`) | `float64` |
| 거래량/잔량/건수/수량/대금 | `int64` |

---

## TR_ID 상수 (client.go 추가 예정)

```go
const (
    // Phase 11.3 — 지수선물
    trIDIndexFuturesTrade = "H0IFCNT0"
    trIDIndexFuturesAsk   = "H0IFASP0"

    // Phase 11.3 — 지수옵션
    trIDIndexOptionTrade = "H0IOCNT0"
    trIDIndexOptionAsk   = "H0IOASP0"

    // Phase 11.3 — 상품선물 (alias → IF schema 재사용)
    trIDCommodityFuturesTrade = "H0CFCNT0"
    trIDCommodityFuturesAsk   = "H0CFASP0"
)
```

---

## Decoder 패턴 (4 개별 decoder + 2 alias)

```go
// 지수선물 체결가
func decodeIndexFuturesTrade(f frame) ([]IndexFuturesTradeEvent, error)
// 지수선물 호가
func decodeIndexFuturesAsk(f frame) ([]IndexFuturesAskEvent, error)
// 지수옵션 체결가
func decodeIndexOptionTrade(f frame) ([]IndexOptionTradeEvent, error)
// 지수옵션 호가
func decodeIndexOptionAsk(f frame) ([]IndexOptionAskEvent, error)
// 상품선물 → alias, decoder 재사용
// decodeIndexFuturesTrade → CommodityFuturesTradeEvent alias
// decodeIndexFuturesAsk   → CommodityFuturesAskEvent alias
```
