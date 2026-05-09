# Phase 11.7 해외선물옵션 실시간 2 EP — Schema Reference

> docs analyzer 결과 (2026-05-09). Phase 11.7 Event/Decoder/Client 의 source of truth.
> docs 위치: `docs/api/해외선물옵션/<파일>.md`

---

## EP Matrix

| # | TR_ID | 한글명 | Fields (docs count) | 모의 | tr_key 형식 | 종목코드 길이 |
|---|---|---|---|---|---|---|
| 1 | HDFFF020 | 해외선물옵션 실시간체결가 | 25 | 미지원 | 해외선물옵션 종목코드 | 6자리 |
| 2 | HDFFF010 | 해외선물옵션 실시간호가 | 35 | 미지원 | 해외선물옵션 종목코드 | 6자리 |

**전체 모의투자 미지원.** WebSocket domain: `ws://ops.koreainvestment.com:21000`

> **중요**: 해외선물옵션 시세는 CME/SGX 거래소의 경우 유료시세 신청 필수.
> 가격 출력값 해석 시 ffcode.mst(해외선물종목마스터 파일)의 sCalcDesz(계산 소수점) 값 활용 필요.

---

## Alias 분석 결과

### Phase 11.2/11.3 국내선물옵션 실시간과 비교

| 항목 | HDFFF020 (해외선물체결) | HDFFF010 (해외선물호가) | H0ZFCNT0 (주식선물 Phase11.2) | H0IFASP0 (지수선물 Phase11.3) |
|---|---|---|---|---|
| 종목코드 필드명 | SERIES_CD | SERIES_CD | FUTS_SHRN_ISCD | FUTS_SHRN_ISCD |
| 시간 필드 | RECV_DATE/TIME | RECV_DATE/TIME | BSOP_HOUR | BSOP_HOUR |
| 가격 필드 | LAST_PRICE/PREV_PRICE | BID_PRICE/ASK_PRICE | STCK_PRPR 등 | FUTS_ASKP/BIDP |
| 호가 depth | - | 5단계 (BID/ASK) | - | 5단계 (FUTS_) |
| 그릭스 포함 | **없음** | **없음** | 없음 | 없음 |
| 장운영 필드 | MRKT_OPEN/CLOSE_DATE/TIME | 없음 | 없음 | 없음 |
| 정산가 필드 | PSTTL_PRICE/SIGN/DIFF | STTL_PRICE | 없음 | 없음 |
| Field count | 25 | 35 | 49 | 38 |

**결론: 해외선물옵션 실시간 2 EP 는 국내선물옵션과 schema 가 완전히 다름. Alias 불가. Distinct.**

- 필드명 prefix/구조 모두 상이 (SERIES_CD vs FUTS_SHRN_ISCD/OPTN_SHRN_ISCD)
- 해외 특유 필드: MRKT_OPEN/CLOSE_DATE/TIME, ACTIVE_FLAG, QUOTSIGN, RECV_TIME2, PSTTL_*, BID/ASK_NUM_*
- **그릭스/IV/HV 없음**: 실시간 응답에도 DELTA/GAMA/VEGA/THETA/RHO/HTS_INTS_VLTL/UNAS_HIST_VLTL 전혀 없음
- 선물/옵션 통합 EP: HDFFF020/HDFFF010 단일 TR_ID 로 선물+옵션 모두 처리. 별도 선물 전용/옵션 전용 EP 없음.

### 선물 vs 옵션 통합 EP 확인

docs 명칭: **"해외선물옵션"** — 선물과 옵션을 구분하지 않고 동일 EP 사용.
tr_key length=6 (docs 명시): 품목코드 기반. 예: ES(S&P500), GC(금), 6A(호주달러) 등.

**Phase 11.7 scope = 정확히 2 EP** (선물 전용/옵션 전용 별도 EP 없음).

---

## TR_ID 명명 규칙

```
HD FF F020  → 해외선물옵션(HDF) 체결가(F020)
HD FF F010  → 해외선물옵션(HDF) 호가(F010)
```

---

## EP1 — HDFFF020 해외선물옵션 실시간체결가 (25 fields)

tr_key: 해외선물옵션 종목코드 6자리 (예: GCM24, 6AM24)
도메인: ws://ops.koreainvestment.com:21000
선물/옵션 통합 EP. 그릭스 없음.

| # | 필드명 | 한글명 | Go Type | 비고 |
|---|---|---|---|---|
| 0 | SERIES_CD | 종목코드 | string | 32자리 |
| 1 | BSNS_DATE | 영업일자 | string | YYYYMMDD |
| 2 | MRKT_OPEN_DATE | 장개시일자 | string | YYYYMMDD |
| 3 | MRKT_OPEN_TIME | 장개시시각 | string | HHMMSS |
| 4 | MRKT_CLOSE_DATE | 장종료일자 | string | YYYYMMDD |
| 5 | MRKT_CLOSE_TIME | 장종료시각 | string | HHMMSS |
| 6 | PREV_PRICE | 전일종가 | decimal.Decimal | 가격 (sCalcDesz 참고) |
| 7 | RECV_DATE | 수신일자 | string | YYYYMMDD |
| 8 | RECV_TIME | 수신시각 | string | HHMMSS (실제 체결시각) |
| 9 | ACTIVE_FLAG | 본장_전산장구분 | string | 1자리 |
| 10 | LAST_PRICE | 체결가격 | decimal.Decimal | 가격 |
| 11 | LAST_QNTT | 체결수량 | int64 | 거래량 |
| 12 | PREV_DIFF_PRICE | 전일대비가 | decimal.Decimal | 가격 |
| 13 | PREV_DIFF_RATE | 등락률 | float64 | 비율 |
| 14 | OPEN_PRICE | 시가 | decimal.Decimal | 가격 |
| 15 | HIGH_PRICE | 고가 | decimal.Decimal | 가격 |
| 16 | LOW_PRICE | 저가 | decimal.Decimal | 가격 |
| 17 | VOL | 누적거래량 | int64 | 거래량 |
| 18 | PREV_SIGN | 전일대비부호 | string | 1자리 |
| 19 | QUOTSIGN | 체결구분 | string | 1자리 (2:매수체결, 5:매도체결) |
| 20 | RECV_TIME2 | 수신시각2 만분의일초 | string | 4자리 |
| 21 | PSTTL_PRICE | 전일정산가 | decimal.Decimal | 가격 |
| 22 | PSTTL_SIGN | 전일정산가대비 | string | 1자리 |
| 23 | PSTTL_DIFF_PRICE | 전일정산가대비가격 | decimal.Decimal | 가격 |
| 24 | PSTTL_DIFF_RATE | 전일정산가대비율 | float64 | 비율 |

> docs 직접 카운트: 25 fields.
> **그릭스 없음 확인**: DELTA/GAMA/VEGA/THETA/RHO/HTS_INTS_VLTL/UNAS_HIST_VLTL 전혀 없음.
> 선물/옵션 통합 단일 schema. REST Phase11.6 해외옵션 응답과 일치 (그릭스 미포함).

---

## EP2 — HDFFF010 해외선물옵션 실시간호가 (35 fields)

tr_key: 해외선물옵션 종목코드 6자리
도메인: ws://ops.koreainvestment.com:21000
5단계 호가. BID(매수)/ASK(매도) 교차 구조 (BID_QNTT/BID_NUM/BID_PRICE + ASK_QNTT/ASK_NUM/ASK_PRICE 그룹).
국내선물옵션과 달리 건수(CSNU) 대신 번호(NUM) 사용. 총잔량 합계 필드 없음.

| # | 필드명 | 한글명 | Go Type | 비고 |
|---|---|---|---|---|
| 0 | SERIES_CD | 종목코드 | string | 32자리 |
| 1 | RECV_DATE | 수신일자 | string | YYYYMMDD |
| 2 | RECV_TIME | 수신시각 | string | 12자리 (나노초 포함) |
| 3 | PREV_PRICE | 전일종가 | decimal.Decimal | 가격 |
| 4 | BID_QNTT_1 | 매수1수량 | int64 | 잔량 |
| 5 | BID_NUM_1 | 매수1번호 | string | 10자리 |
| 6 | BID_PRICE_1 | 매수1호가 | decimal.Decimal | 가격 |
| 7 | ASK_QNTT_1 | 매도1수량 | int64 | 잔량 |
| 8 | ASK_NUM_1 | 매도1번호 | string | 10자리 |
| 9 | ASK_PRICE_1 | 매도1호가 | decimal.Decimal | 가격 |
| 10 | BID_QNTT_2 | 매수2수량 | int64 | 잔량 |
| 11 | BID_NUM_2 | 매수2번호 | string | 10자리 |
| 12 | BID_PRICE_2 | 매수2호가 | decimal.Decimal | 가격 |
| 13 | ASK_QNTT_2 | 매도2수량 | int64 | 잔량 |
| 14 | ASK_NUM_2 | 매도2번호 | string | 10자리 |
| 15 | ASK_PRICE_2 | 매도2호가 | decimal.Decimal | 가격 |
| 16 | BID_QNTT_3 | 매수3수량 | int64 | 잔량 |
| 17 | BID_NUM_3 | 매수3번호 | string | 10자리 |
| 18 | BID_PRICE_3 | 매수3호가 | decimal.Decimal | 가격 |
| 19 | ASK_QNTT_3 | 매도3수량 | int64 | 잔량 |
| 20 | ASK_NUM_3 | 매도3번호 | string | 10자리 |
| 21 | ASK_PRICE_3 | 매도3호가 | decimal.Decimal | 가격 |
| 22 | BID_QNTT_4 | 매수4수량 | int64 | 잔량 |
| 23 | BID_NUM_4 | 매수4번호 | string | 10자리 |
| 24 | BID_PRICE_4 | 매수4호가 | decimal.Decimal | 가격 |
| 25 | ASK_QNTT_4 | 매도4수량 | int64 | 잔량 |
| 26 | ASK_NUM_4 | 매도4번호 | string | 10자리 |
| 27 | ASK_PRICE_4 | 매도4호가 | decimal.Decimal | 가격 |
| 28 | BID_QNTT_5 | 매수5수량 | int64 | 잔량 |
| 29 | BID_NUM_5 | 매수5번호 | string | 10자리 |
| 30 | BID_PRICE_5 | 매수5호가 | decimal.Decimal | 가격 |
| 31 | ASK_QNTT_5 | 매도5수량 | int64 | 잔량 |
| 32 | ASK_NUM_5 | 매도5번호 | string | 10자리 |
| 33 | ASK_PRICE_5 | 매도5호가 | decimal.Decimal | 가격 |
| 34 | STTL_PRICE | 전일정산가 | decimal.Decimal | 가격 |

> docs 직접 카운트: 35 fields.
> **국내선물옵션과 구조 차이**: BID/ASK 교차 배열 (BID_1 → ASK_1 → BID_2 → ASK_2 ...) vs 국내는 전체매도 → 전체매수 분리.
> **번호(NUM) 필드**: 국내 CSNU(건수) 대신 NUM(번호) 사용 — string 처리.
> **총잔량 합계 없음**: TOTAL_ASKP_RSQN/TOTAL_BIDP_RSQN 없음.
> RECV_TIME Length=12 (나노초 포함, 국내 HHMMSS 6자리와 다름).

---

## Anomalies

1. **RECV_TIME 길이**: HDFFF010 의 RECV_TIME Length=12 (국내 선물옵션 BSOP_HOUR=6과 다름). 나노초 포함 가능성. string 처리.

2. **BID/ASK 교차 구조**: HDFFF010 은 BID_QNTT_1→BID_NUM_1→BID_PRICE_1→ASK_QNTT_1→ASK_NUM_1→ASK_PRICE_1 순서로 반복. 국내는 전체 매도 먼저, 전체 매수 나중.

3. **BID_NUM/ASK_NUM**: 국내선물옵션 CSNU(건수, int64) 와 달리 NUM(번호, string). docs Length=10.

4. **SERIES_CD docs Type**: docs 에서 HDFFF010 의 SERIES_CD Type이 "Object" 로 표기됨. 실제는 String (다른 EP 참고). string 처리.

5. **선물/옵션 단일 EP**: 국내선물옵션은 선물/옵션 별도 TR_ID 인 반면, 해외선물옵션은 HDFFF020/HDFFF010 단일 TR_ID 로 선물+옵션 통합 처리. Phase scope = 2 EP (선물용/옵션용 분리 없음).

6. **그릭스 완전 미포함**: 실시간 응답에 DELTA/GAMA/VEGA/THETA/RHO/HTS_INTS_VLTL/UNAS_HIST_VLTL 없음. REST Phase11.6 해외옵션 현재가 응답과 동일하게 미포함.

7. **Endpoint URL 동일**: 실전 Domain = `ws://ops.koreainvestment.com:21000` (기존 인프라와 동일). 별도 endpoint 불필요.

---

## Type 매핑 정책 (Phase 8 와 동일)

| KIS Type | Go field type |
|---|---|
| String (YYYYMMDD, HHMMSS, code, sign, YN, flag) | `string` |
| 가격 (LAST_PRICE, PREV_PRICE, OPEN/HIGH/LOW, PSTTL_*, BID_PRICE, ASK_PRICE, STTL_PRICE) | `decimal.Decimal` |
| 비율 (PREV_DIFF_RATE, PSTTL_DIFF_RATE) | `float64` |
| 거래량/잔량 (LAST_QNTT, VOL, BID_QNTT, ASK_QNTT) | `int64` |
| 번호 (BID_NUM, ASK_NUM) | `string` |

---

## TR_ID 상수 (client.go 추가 예정)

```go
const (
    // Phase 11.7 — 해외선물옵션 실시간
    trIDOverseasFuturesOptionTrade = "HDFFF020"
    trIDOverseasFuturesOptionAsk   = "HDFFF010"
)
```

---

## Decoder 패턴 (2 개별 decoder)

```go
// 해외선물옵션 실시간체결가
func decodeOverseasFuturesOptionTrade(f frame) ([]OverseasFuturesOptionTradeEvent, error)
// 해외선물옵션 실시간호가
func decodeOverseasFuturesOptionAsk(f frame) ([]OverseasFuturesOptionAskEvent, error)
```
