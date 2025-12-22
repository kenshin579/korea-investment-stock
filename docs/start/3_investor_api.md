# 투자자별 매매현황 조회 API

개인, 외국인, 기관 등 투자자별 매매현황을 조회할 수 있는 API 목록입니다.

## API 목록

| API명 | API 경로 | TR ID | 설명 |
|-------|----------|-------|------|
| **주식현재가 투자자** | `/uapi/domestic-stock/v1/quotations/inquire-investor` | FHKST01010900 | 종목별 개인/외국인/기관 일자별 순매수 수량·금액 조회 |
| **국내기관_외국인 매매종목가집계** | `/uapi/domestic-stock/v1/quotations/foreign-institution-total` | FHPTJ04400000 | 장중 외국인/기관 매매종목 순매수 상위 조회 (가집계) |
| **시장별 투자자매매동향(시세)** | `/uapi/domestic-stock/v1/quotations/inquire-investor-time-by-market` | FHPTJ04030000 | 코스피/코스닥 등 시장별 투자자 실시간 매매 현황 |
| **시장별 투자자매매동향(일별)** | `/uapi/domestic-stock/v1/quotations/inquire-investor-daily-by-market` | FHPTJ04040000 | 시장별 투자자 일별 매매 동향 |
| **종목별 투자자매매동향(일별)** | `/uapi/domestic-stock/v1/quotations/investor-trade-by-stock-daily` | FHPTJ04160001 | 특정 종목의 투자자별 일별 매매 내역 |
| **종목별 외인기관 추정가집계** | `/uapi/domestic-stock/v1/quotations/investor-trend-estimate` | HHPTJ04160200 | 장중 종목별 외국인/기관 추정 순매수 |
| **외국계 매매종목 가집계** | `/uapi/domestic-stock/v1/quotations/frgnmem-trade-estimate` | FHKST644100C0 | 외국계 증권사 매매종목 가집계 |
| **종목별 외국계 순매수추이** | `/uapi/domestic-stock/v1/quotations/frgnmem-pchs-trend` | FHKST644400C0 | 특정 종목의 외국계 시간별 순매수 추이 |

---

## 상세 설명

### 1. 주식현재가 투자자 [v1_국내주식-012]

- **API 경로**: `/uapi/domestic-stock/v1/quotations/inquire-investor`
- **TR ID**: FHKST01010900
- **HTS 화면**: -

#### 개요
주식현재가 투자자 API입니다. 개인, 외국인, 기관 등 투자 정보를 확인할 수 있습니다.

#### 주요 파라미터
| 파라미터 | 설명 |
|----------|------|
| FID_COND_MRKT_DIV_CODE | J: KRX, NX: NXT, UN: 통합 |
| FID_INPUT_ISCD | 종목코드 (ex: 005930) |

#### 응답 데이터
- `prsn_ntby_qty`: 개인 순매수 수량
- `frgn_ntby_qty`: 외국인 순매수 수량
- `orgn_ntby_qty`: 기관계 순매수 수량
- `prsn_ntby_tr_pbmn`: 개인 순매수 거래 대금
- `frgn_ntby_tr_pbmn`: 외국인 순매수 거래 대금
- `orgn_ntby_tr_pbmn`: 기관계 순매수 거래 대금

#### 유의사항
- 외국인은 외국인(외국인투자등록 고유번호가 있는 경우) + 기타 외국인을 지칭
- 당일 데이터는 장 종료 후 제공

---

### 2. 국내기관_외국인 매매종목가집계 [국내주식-037]

- **API 경로**: `/uapi/domestic-stock/v1/quotations/foreign-institution-total`
- **TR ID**: FHPTJ04400000
- **HTS 화면**: [0440] 외국인/기관 매매종목 가집계

#### 개요
장중 외국인/기관 매매종목 가집계 정보를 조회합니다.

#### 주요 파라미터
| 파라미터 | 설명 |
|----------|------|
| FID_INPUT_ISCD | 0000: 전체, 0001: 코스피, 1001: 코스닥 |
| FID_DIV_CLS_CODE | 0: 수량정렬, 1: 금액정렬 |
| FID_RANK_SORT_CLS_CODE | 0: 순매수상위, 1: 순매도상위 |
| FID_ETC_CLS_CODE | 0: 전체, 1: 외국인, 2: 기관계, 3: 기타 |

#### 응답 데이터
- `frgn_ntby_qty`: 외국인 순매수 수량
- `orgn_ntby_qty`: 기관계 순매수 수량
- `ivtr_ntby_qty`: 투자신탁 순매수 수량
- `bank_ntby_qty`: 은행 순매수 수량
- `insu_ntby_qty`: 보험 순매수 수량
- `fund_ntby_qty`: 기금 순매수 수량

#### 유의사항
- 증권사 직원이 장중에 집계/입력한 자료의 단순 누계
- 입력시간: 외국인 09:30, 11:20, 13:20, 14:30 / 기관종합 10:00, 11:20, 13:20, 14:30
- 입력 시간은 ±10분 차이 발생 가능

---

### 3. 시장별 투자자매매동향(시세) [v1_국내주식-074]

- **API 경로**: `/uapi/domestic-stock/v1/quotations/inquire-investor-time-by-market`
- **TR ID**: FHPTJ04030000
- **HTS 화면**: [0403] 시장별 시간동향

#### 개요
시장별 투자자 실시간 매매동향을 조회합니다.

#### 주요 파라미터
| 파라미터 | 설명 |
|----------|------|
| fid_input_iscd | 시장구분 - KSP: 코스피, KSQ: 코스닥, ETF, ELW, ETN 등 |
| fid_input_iscd_2 | 업종구분 - 0001: 코스피 종합, 1001: 코스닥 종합 등 |

#### 응답 데이터
투자자별 매도/매수/순매수 거래량 및 거래대금:
- 외국인 (`frgn_*`)
- 개인 (`prsn_*`)
- 기관계 (`orgn_*`)
- 증권 (`scrt_*`)
- 투자신탁 (`ivtr_*`)
- 사모펀드 (`pe_fund_*`)
- 은행 (`bank_*`)
- 보험 (`insu_*`)
- 종금 (`mrbn_*`)
- 기금 (`fund_*`)
- 기타단체 (`etc_orgt_*`)
- 기타법인 (`etc_corp_*`)

---

### 4. 시장별 투자자매매동향(일별) [국내주식-075]

- **API 경로**: `/uapi/domestic-stock/v1/quotations/inquire-investor-daily-by-market`
- **TR ID**: FHPTJ04040000
- **HTS 화면**: [0404] 시장별 일별동향

#### 개요
시장별 투자자 일별 매매동향을 조회합니다.

#### 주요 파라미터
| 파라미터 | 설명 |
|----------|------|
| FID_COND_MRKT_DIV_CODE | 시장구분코드 (업종 U) |
| FID_INPUT_ISCD | 업종분류코드 |
| FID_INPUT_DATE_1 | 조회 시작일 (ex: 20240517) |
| FID_INPUT_ISCD_1 | KSP: 코스피, KSQ: 코스닥 |

#### 응답 데이터
- 업종 지수 정보 (현재가, 전일대비 등)
- 투자자별 순매수 수량/금액 (외국인, 개인, 기관 등)
- 외국인 등록/비등록 구분 데이터 포함

---

### 5. 종목별 투자자매매동향(일별)

- **API 경로**: `/uapi/domestic-stock/v1/quotations/investor-trade-by-stock-daily`
- **TR ID**: FHPTJ04160001
- **HTS 화면**: [0416] 종목별 일별동향

#### 개요
특정 종목의 투자자별 일별 매매 내역을 조회합니다.

#### 주요 파라미터
| 파라미터 | 설명 |
|----------|------|
| FID_COND_MRKT_DIV_CODE | J: KRX, NX: NXT, UN: 통합 |
| FID_INPUT_ISCD | 종목코드 (6자리) |
| FID_INPUT_DATE_1 | 조회일 (ex: 20250812) |

#### 응답 데이터
- 주가 정보 (종가, 시가, 고가, 저가)
- 투자자별 순매수 수량/금액
- 투자자별 매도/매수 거래량/대금
- 외국인 등록/비등록 구분 데이터

#### 유의사항
- 단위: 금액(백만원), 수량(주)
- 해당일 조회는 장 종료 후 정상 조회 가능

---

### 6. 종목별 외인기관 추정가집계 [v1_국내주식-046]

- **API 경로**: `/uapi/domestic-stock/v1/quotations/investor-trend-estimate`
- **TR ID**: HHPTJ04160200
- **MTS 화면**: 국내 현재가 > 투자자 > 투자자동향 > 추정(주)

#### 개요
장중 종목별 외국인/기관 추정 순매수를 조회합니다.

#### 주요 파라미터
| 파라미터 | 설명 |
|----------|------|
| MKSC_SHRN_ISCD | 종목코드 |

#### 응답 데이터
| 필드 | 설명 |
|------|------|
| bsop_hour_gb | 입력구분 (1: 09:30, 2: 10:00, 3: 11:20, 4: 13:20, 5: 14:30) |
| frgn_fake_ntby_qty | 외국인 수량 (가집계) |
| orgn_fake_ntby_qty | 기관 수량 (가집계) |
| sum_fake_ntby_qty | 합산 수량 (가집계) |

#### 유의사항
- 증권사 직원이 장중에 집계/입력한 자료의 단순 누계
- 사정에 따라 입력 시간 변동 가능

---

### 7. 외국계 매매종목 가집계 [국내주식-161]

- **API 경로**: `/uapi/domestic-stock/v1/quotations/frgnmem-trade-estimate`
- **TR ID**: FHKST644100C0
- **HTS 화면**: [0430] 외국계 매매종목 가집계

#### 개요
외국계 증권사의 매매종목 가집계를 조회합니다.

#### 주요 파라미터
| 파라미터 | 설명 |
|----------|------|
| FID_INPUT_ISCD | 0000: 전체, 1001: 코스피, 2001: 코스닥 |
| FID_RANK_SORT_CLS_CODE | 0: 금액순, 1: 수량순 |
| FID_RANK_SORT_CLS_CODE_2 | 0: 매수순, 1: 매도순 |

#### 응답 데이터
- `glob_ntsl_qty`: 외국계 순매도 수량
- `glob_total_seln_qty`: 외국계 총매도 수량
- `glob_total_shnu_qty`: 외국계 총매수 수량

---

### 8. 종목별 외국계 순매수추이 [국내주식-164]

- **API 경로**: `/uapi/domestic-stock/v1/quotations/frgnmem-pchs-trend`
- **TR ID**: FHKST644400C0
- **HTS 화면**: [0433] 종목별 외국계 순매수추이

#### 개요
특정 종목의 외국계 시간별 순매수 추이를 조회합니다.

#### 주요 파라미터
| 파라미터 | 설명 |
|----------|------|
| FID_INPUT_ISCD | 종목코드 (ex: 005930) |
| FID_INPUT_ISCD_2 | 외국계 전체 (99999) |
| FID_COND_MRKT_DIV_CODE | J (KRX만 지원) |

#### 응답 데이터
- `bsop_hour`: 영업시간
- `frgn_seln_vol`: 외국인 매도 거래량
- `frgn_shnu_vol`: 외국인 매수 거래량
- `glob_ntby_qty`: 외국계 순매수 수량
- `frgn_ntby_qty_icdc`: 외국인 순매수 수량 증감

---

## 용도별 API 추천

### 종목별 조회

| 용도 | 추천 API |
|------|----------|
| 일별 투자자 동향 | 주식현재가 투자자, 종목별 투자자매매동향(일별) |
| 장중 외인/기관 추정 | 종목별 외인기관 추정가집계 |
| 외국계 시간별 추이 | 종목별 외국계 순매수추이 |

### 시장 전체 조회

| 용도 | 추천 API |
|------|----------|
| 실시간 동향 | 시장별 투자자매매동향(시세) |
| 일별 동향 | 시장별 투자자매매동향(일별) |
| 외인/기관 순매수 상위 | 국내기관_외국인 매매종목가집계, 외국계 매매종목 가집계 |

---

## 공통 유의사항

1. **가집계 데이터**: 증권사 직원이 장중 입력한 데이터로 확정 수치가 아닙니다.
2. **당일 확정 데이터**: 장 종료 후 제공됩니다.
3. **모의투자**: 대부분의 투자자 동향 API는 모의투자를 지원하지 않습니다.
4. **외국인 구분**: 외국인 = 외국인(등록) + 기타 외국인(비등록)
