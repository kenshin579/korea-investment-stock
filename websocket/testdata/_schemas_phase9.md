# Phase 9 NXT/통합 변형 10 EP — Schema Reference

> docs analyzer 결과 (2026-05-09). Phase 9 Event/Decoder 의 source of truth.

## EP Matrix

| EP | NXT TR_ID | 통합 TR_ID | Body fields | KRX 차이 | 모의 |
|---|---|---|---|---|---|
| 체결가 | H0NXCNT0 | H0UNCNT0 | 46 | 동일 (22번 CNTG_CLS_CODE vs KRX CCLD_DVSN) | 미지원 |
| 호가 | H0NXASP0 | H0UNASP0 | 65 | KRX 59 + KMID/NMID 6 (중간가) | 미지원 |
| 예상체결 | H0NXANC0 | H0UNANC0 | 46 | KRX 45 + VI_STND_PRC | 미지원 |
| 프로그램매매 | H0NXPGM0 | H0UNPGM0 | 11 | 신규 EP (Phase 8 OoS) | 미지원 |
| 회원사 | H0NXMBC0 | H0UNMBC0 | 78 | 신규 EP (Phase 8 OoS) | 미지원 |

**중요한 패턴**: NXT 와 통합은 schema **완전 동일**. 차이는 TR_ID 와 모의 지원만. → 5 base struct + 10 type alias 권장.

## Plan Deviation (필요시 정정)

이전 추정과 docs 실제값:

| 항목 | 추정 | 실제 |
|---|---|---|
| NXT 호가 fields | 59 (KRX 동일) | **65** (+KMID/NMID) |
| NXT 예상체결 fields | 45 (KRX 동일) | **46** (+VI_STND_PRC) |
| 회원사 fields | 미상 | **78** (docs 응답 표 직접 검증, 2026-05-09) |
| 프로그램매매 | KRX schema | 신규 schema (11 fields, KRX 와 별도) |

---

## EP1/2 — 체결가 NXT/통합 (46 fields)

KRX H0STCNT0 (46) 와 동일 schema, 22번 필드명만 차이:
- KRX: `CCLD_DVSN` (체결구분)
- NXT/통합: `CNTG_CLS_CODE` (체결구분, 의미 동일)

전체 필드 (KRX 46 fields 와 동일 순서):
```
[0]  MKSC_SHRN_ISCD              (Symbol, string)
[1]  STCK_CNTG_HOUR              (Time, string)
[2]  STCK_PRPR                   (Price, decimal)
[3]  PRDY_VRSS_SIGN              (PrevDiffSign, string)
[4]  PRDY_VRSS                   (PrevDiff, decimal)
[5]  PRDY_CTRT                   (PrevChangeRate, float64)
[6]  WGHN_AVRG_STCK_PRC          (WeightedAvg, decimal)
[7]  STCK_OPRC                   (Open, decimal)
[8]  STCK_HGPR                   (High, decimal)
[9]  STCK_LWPR                   (Low, decimal)
[10] ASKP1                       (Ask1, decimal)
[11] BIDP1                       (Bid1, decimal)
[12] CNTG_VOL                    (TradeVolume, int64)
[13] ACML_VOL                    (AccumVolume, int64)
[14] ACML_TR_PBMN                (AccumValue, int64)
[15] SELN_CNTG_CSNU              (AskCount, int64)
[16] SHNU_CNTG_CSNU              (BidCount, int64)
[17] NTBY_CNTG_CSNU              (NetCount, int64)
[18] CTTR                        (TradeStrength, float64)
[19] SELN_CNTG_SMTN              (TotalAskVolume, int64)
[20] SHNU_CNTG_SMTN              (TotalBidVolume, int64)
[21] CNTG_CLS_CODE               (TradeKind, string)        # KRX 는 CCLD_DVSN
[22] SHNU_RATE                   (BidRate, float64)
[23] PRDY_VOL_VRSS_ACML_VOL_RATE (PrevVolRate, float64)
[24] OPRC_HOUR                   (OpenTime, string)
[25] OPRC_VRSS_PRPR_SIGN         (OpenDiffSign, string)
[26] OPRC_VRSS_PRPR              (OpenDiff, decimal)
[27] HGPR_HOUR                   (HighTime, string)
[28] HGPR_VRSS_PRPR_SIGN         (HighDiffSign, string)
[29] HGPR_VRSS_PRPR              (HighDiff, decimal)
[30] LWPR_HOUR                   (LowTime, string)
[31] LWPR_VRSS_PRPR_SIGN         (LowDiffSign, string)
[32] LWPR_VRSS_PRPR              (LowDiff, decimal)
[33] BSOP_DATE                   (BusinessDate, string)
[34] NEW_MKOP_CLS_CODE           (MarketOpCode, string)
[35] TRHT_YN                     (TradeHaltYN, string)
[36] ASKP_RSQN1                  (Ask1Size, int64)
[37] BIDP_RSQN1                  (Bid1Size, int64)
[38] TOTAL_ASKP_RSQN             (TotalAskSize, int64)
[39] TOTAL_BIDP_RSQN             (TotalBidSize, int64)
[40] VOL_TNRT                    (VolumeTurnover, float64)
[41] PRDY_SMNS_HOUR_ACML_VOL     (PrevSameTimeAccumVol, int64)
[42] PRDY_SMNS_HOUR_ACML_VOL_RATE (PrevSameTimeAccumVolRate, float64)
[43] HOUR_CLS_CODE               (HourCode, string)
[44] MRKT_TRTM_CLS_CODE          (MarketTermCode, string)
[45] VI_STND_PRC                 (ViStandardPrice, decimal)
```

**Event 구조 (5 base + alias 패턴)**:
```go
// base struct — NXT/통합 공유
type AltMarketTradeEvent struct { ... 46 fields ... Raw []string }

// type alias (시장 구분, 컴파일 타임에 해소)
type NxtTradeEvent = AltMarketTradeEvent
type UnifiedTradeEvent = AltMarketTradeEvent
```

KRX 는 별도 (`KrxTradeEvent`, 22번=CCLD_DVSN 의미 동일이지만 type 분리 — Phase 8 그대로).

---

## EP3/4 — 호가 NXT/통합 (65 fields)

KRX H0STASP0 (59) 의 superset. 끝부분에 6 fields 추가:

```
[0..58] KRX 59 fields 와 동일 순서:
[0]  MKSC_SHRN_ISCD
[1]  BSOP_HOUR
[2]  HOUR_CLS_CODE
[3..12]  ASKP1..10
[13..22] BIDP1..10
[23..32] ASKP_RSQN1..10
[33..42] BIDP_RSQN1..10
[43] TOTAL_ASKP_RSQN
[44] TOTAL_BIDP_RSQN
[45] OVTM_TOTAL_ASKP_RSQN
[46] OVTM_TOTAL_BIDP_RSQN
[47] ANTC_CNPR
[48] ANTC_CNQN
[49] ANTC_VOL
[50] ANTC_CNTG_VRSS
[51] ANTC_CNTG_VRSS_SIGN
[52] ANTC_CNTG_PRDY_CTRT
[53] ACML_VOL
[54] TOTAL_ASKP_RSQN_ICDC
[55] TOTAL_BIDP_RSQN_ICDC
[56] OVTM_TOTAL_ASKP_ICDC
[57] OVTM_TOTAL_BIDP_ICDC
[58] STCK_DEAL_CLS_CODE

# NXT/통합 추가 6 fields:
[59] KMID_PRC          (KrxMidPrice, decimal)
[60] KMID_TOTAL_RSQN   (KrxMidTotalSize, int64)
[61] KMID_CLS_CODE     (KrxMidCode, string)
[62] NMID_PRC          (NxtMidPrice, decimal)
[63] NMID_TOTAL_RSQN   (NxtMidTotalSize, int64)
[64] NMID_CLS_CODE     (NxtMidCode, string)
```

**Event 구조**:
```go
type AltMarketAskEvent struct {
    Symbol  string
    Time    string
    HourCode string
    Ask     [10]decimal.Decimal
    Bid     [10]decimal.Decimal
    AskSize [10]int64
    BidSize [10]int64
    // ... 16 추가 (TotalAskSize 부터 STCK_DEAL_CLS_CODE 까지)
    // KRX/NXT 중간가 6 fields:
    KrxMidPrice      decimal.Decimal
    KrxMidTotalSize  int64
    KrxMidCode       string
    NxtMidPrice      decimal.Decimal
    NxtMidTotalSize  int64
    NxtMidCode       string
    Raw              []string
}
type NxtAskEvent = AltMarketAskEvent
type UnifiedAskEvent = AltMarketAskEvent
```

---

## EP5/6 — 예상체결 NXT/통합 (46 fields)

KRX H0STANC0 (45) + 1 field. 끝에 `VI_STND_PRC` 추가.

```
[0..44] KRX H0STANC0 와 동일 순서 (22번=CNTG_CLS_CODE — KRX 와 동일)
[45] VI_STND_PRC  (ViStandardPrice, decimal)
```

```go
type AltMarketExpectTradeEvent struct {
    // KRX H0STANC0 의 45 fields + VI_STND_PRC
    Symbol, Time string
    Price        decimal.Decimal
    // ... (체결가와 거의 동일하지만 ASKP_RSQN/TOTAL_*_RSQN 일부 누락 + VI_STND_PRC)
    ViStandardPrice decimal.Decimal  // 새 필드
    Raw             []string
}
type NxtExpectTradeEvent = AltMarketExpectTradeEvent
type UnifiedExpectTradeEvent = AltMarketExpectTradeEvent
```

---

## EP7/8 — 프로그램매매 NXT/통합 (11 fields, 신규)

```
[0]  MKSC_SHRN_ISCD     (Symbol, string)
[1]  STCK_CNTG_HOUR     (Time, string)
[2]  SELN_CNQN          (AskQuantity, int64)
[3]  SELN_TR_PBMN       (AskValue, int64)
[4]  SHNU_CNQN          (BidQuantity, int64)
[5]  SHNU_TR_PBMN       (BidValue, int64)
[6]  NTBY_CNQN          (NetQuantity, int64)
[7]  NTBY_TR_PBMN       (NetValue, int64)
[8]  SELN_RSQN          (AskRemainingSize, int64)
[9]  SHNU_RSQN          (BidRemainingSize, int64)
[10] WHOL_NTBY_QTY      (TotalNetQuantity, int64)
```

```go
type ProgramTradeEvent struct {
    Symbol, Time string
    AskQuantity, AskValue, BidQuantity, BidValue, NetQuantity, NetValue int64
    AskRemainingSize, BidRemainingSize, TotalNetQuantity int64
    Raw []string
}
type NxtProgramTradeEvent = ProgramTradeEvent
type UnifiedProgramTradeEvent = ProgramTradeEvent
```

---

## EP9/10 — 회원사 NXT/통합 (78 fields, 신규)

5단계 매도/매수 회원사 + 외국계 통계.

구조 (인덱스 그룹 단위):
```
[0]                          MKSC_SHRN_ISCD          (Symbol, string)
[1..5]                       SELN2_MBCR_NAME1..5     (SellBrokerNames, [5]string)
[6..10]                      BYOV_MBCR_NAME1..5      (BuyBrokerNames, [5]string)
[11..15]                     TOTAL_SELN_QTY1..5      (TotalSellQty, [5]int64)
[16..20]                     TOTAL_SHNU_QTY1..5      (TotalBuyQty, [5]int64)
[21..25]                     SELN_MBCR_GLOB_YN_1..5  (SellGlobalYN, [5]string)
[26..30]                     SHNU_MBCR_GLOB_YN_1..5  (BuyGlobalYN, [5]string)
[31..35]                     SELN_MBCR_NO1..5        (SellBrokerCodes, [5]string)
[36..40]                     SHNU_MBCR_NO1..5        (BuyBrokerCodes, [5]string)
[41..45]                     SELN_MBCR_RLIM1..5      (SellRatio, [5]float64)
[46..50]                     SHNU_MBCR_RLIM1..5      (BuyRatio, [5]float64)
[51..55]                     SELN_QTY_ICDC1..5       (SellQtyChange, [5]int64)
[56..60]                     SHNU_QTY_ICDC1..5       (BuyQtyChange, [5]int64)
[61]                         GLOB_TOTAL_SELN_QTY     (GlobalTotalSellQty, int64)
[62]                         GLOB_TOTAL_SHNU_QTY     (GlobalTotalBuyQty, int64)
[63]                         GLOB_TOTAL_SELN_QTY_ICDC (GlobalSellQtyChange, int64)
[64]                         GLOB_TOTAL_SHNU_QTY_ICDC (GlobalBuyQtyChange, int64)
[65]                         GLOB_NTBY_QTY            (GlobalNetBuyQty, int64)
[66]                         GLOB_SELN_RLIM           (GlobalSellRatio, float64)
[67]                         GLOB_SHNU_RLIM           (GlobalBuyRatio, float64)
[68..72]                     SELN2_MBCR_ENG_NAME1..5  (SellBrokerEngNames, [5]string)
[73..77]                     BYOV_MBCR_ENG_NAME1..5   (BuyBrokerEngNames, [5]string)
```

총 **78 fields** (실제 docs 응답 body 표 기준, 2026-05-09 직접 검증).
Request header (approval_key/custtype/tr_type/content-type) 와 request body (tr_id/tr_key) 는 응답 body 와 별개.
NXT 와 통합 docs 모두 응답 body 78 fields 동일.

```go
type MemberEvent struct {
    Symbol             string
    SellBrokerNames    [5]string
    BuyBrokerNames     [5]string
    TotalSellQty       [5]int64
    TotalBuyQty        [5]int64
    SellGlobalYN       [5]string
    BuyGlobalYN        [5]string
    SellBrokerCodes    [5]string
    BuyBrokerCodes     [5]string
    SellRatio          [5]float64
    BuyRatio           [5]float64
    SellQtyChange      [5]int64
    BuyQtyChange       [5]int64
    GlobalTotalSellQty   int64
    GlobalTotalBuyQty    int64
    GlobalSellQtyChange  int64
    GlobalBuyQtyChange   int64
    GlobalNetBuyQty      int64
    GlobalSellRatio      float64
    GlobalBuyRatio       float64
    SellBrokerEngNames [5]string
    BuyBrokerEngNames  [5]string
    Raw                []string
}
type NxtMemberEvent = MemberEvent
type UnifiedMemberEvent = MemberEvent
```

---

## Type 매핑 정책 (Phase 8 와 동일)

| KIS Type | Go field type |
|---|---|
| String (HHMMSS, code, sign, name) | `string` |
| Number 가격/금액 (4자리, 8자리) | `decimal.Decimal` |
| Number 비율 (`*_RATE`, `CTTR`, `SHNU_RATE`, `*_RLIM`) | `float64` |
| Number 거래량/금액/수량/체결건수 | `int64` |

NXT/통합 docs 가 모든 필드를 `String` 으로 표기해도 type 매핑은 KRX 와 동일하게 (decoder 가 string→decimal/int64/float64 변환).

---

## Decoder 패턴

5 base decoder (NXT/통합 공유):
```go
func decodeAltMarketTrade(f frame) ([]AltMarketTradeEvent, error)
func decodeAltMarketAsk(f frame) ([]AltMarketAskEvent, error)
func decodeAltMarketExpectTrade(f frame) ([]AltMarketExpectTradeEvent, error)
func decodeProgramTrade(f frame) ([]ProgramTradeEvent, error)
func decodeMember(f frame) ([]MemberEvent, error)
```

routeRealtime 에서 TR_ID 별 분기 (10 cases):
```go
case trIDNxtTrade, trIDUnifiedTrade:    // H0NXCNT0 / H0UNCNT0
    evs, err := decodeAltMarketTrade(f)
    // 시장별 RouteXxx 분기:
    if f.TrID == trIDNxtTrade {
        for _, ev := range evs { dispatcher.RouteNxtTrade(ev) }
    } else {
        for _, ev := range evs { dispatcher.RouteUnifiedTrade(ev) }
    }
```

또는 라우팅 helper:
```go
case trIDNxtTrade:
    routeAltMarketTrade(f, dispatcher.RouteNxtTrade)
case trIDUnifiedTrade:
    routeAltMarketTrade(f, dispatcher.RouteUnifiedTrade)
```

---

## TR_ID 상수 (client.go 추가)

```go
const (
    // Phase 9 — NXT
    trIDNxtTrade           = "H0NXCNT0"
    trIDNxtAsk             = "H0NXASP0"
    trIDNxtExpectTrade     = "H0NXANC0"
    trIDNxtProgramTrade    = "H0NXPGM0"
    trIDNxtMember          = "H0NXMBC0"
    
    // Phase 9 — 통합
    trIDUnifiedTrade        = "H0UNCNT0"
    trIDUnifiedAsk          = "H0UNASP0"
    trIDUnifiedExpectTrade  = "H0UNANC0"
    trIDUnifiedProgramTrade = "H0UNPGM0"
    trIDUnifiedMember       = "H0UNMBC0"
)
```
