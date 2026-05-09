package websocket

import (
	"strconv"

	"github.com/shopspring/decimal"
)

// 각 EP 의 필드 수 상수
const (
	krxTradeFieldCount           = 46 // H0STCNT0
	krxAskFieldCount             = 59 // H0STASP0
	krxExpectTradeFieldCount     = 45 // H0STANC0
	krxOvernightTradeFieldCount  = 43 // H0STOUP0
	krxOvernightExpectFieldCount = 43 // H0STOAC0
)

// --------------------------------------------------------------------------
// 헬퍼 함수 (파싱 실패 시 zero 반환, 에러 안 던짐)
// --------------------------------------------------------------------------

// asDecimal 은 문자열을 decimal.Decimal 로 변환한다. 파싱 실패 시 zero 반환.
func asDecimal(s string) decimal.Decimal {
	if s == "" {
		return decimal.Zero
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero
	}
	return d
}

// asInt64 는 문자열을 int64 로 변환한다. 파싱 실패 시 0 반환.
func asInt64(s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return v
}

// asFloat 는 문자열을 float64 로 변환한다. 파싱 실패 시 0 반환.
func asFloat(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

// --------------------------------------------------------------------------
// H0STCNT0 — 실시간체결가 (KRX 본장, 46 fields)
// --------------------------------------------------------------------------

// decodeKrxTrade 는 H0STCNT0 frame 을 KrxTradeEvent 슬라이스로 디코딩한다.
func decodeKrxTrade(f frame) ([]KrxTradeEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, krxTradeFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]KrxTradeEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseKrxTradeChunk(c))
	}
	return out, nil
}

func parseKrxTradeChunk(c []string) KrxTradeEvent {
	return KrxTradeEvent{
		Symbol:                   c[0],
		Time:                     c[1],
		Price:                    asDecimal(c[2]),
		PrevDiffSign:             c[3],
		PrevDiff:                 asDecimal(c[4]),
		PrevChangeRate:           asFloat(c[5]),
		WeightedAvg:              asDecimal(c[6]),
		Open:                     asDecimal(c[7]),
		High:                     asDecimal(c[8]),
		Low:                      asDecimal(c[9]),
		Ask1:                     asDecimal(c[10]),
		Bid1:                     asDecimal(c[11]),
		TradeVolume:              asInt64(c[12]),
		AccumVolume:              asInt64(c[13]),
		AccumValue:               asInt64(c[14]),
		AskCount:                 asInt64(c[15]),
		BidCount:                 asInt64(c[16]),
		NetCount:                 asInt64(c[17]),
		TradeStrength:            asFloat(c[18]),
		TotalAskVolume:           asInt64(c[19]),
		TotalBidVolume:           asInt64(c[20]),
		TradeKind:                c[21], // CCLD_DVSN
		BidRate:                  asFloat(c[22]),
		PrevVolRate:              asFloat(c[23]),
		OpenTime:                 c[24],
		OpenDiffSign:             c[25],
		OpenDiff:                 asDecimal(c[26]),
		HighTime:                 c[27],
		HighDiffSign:             c[28],
		HighDiff:                 asDecimal(c[29]),
		LowTime:                  c[30],
		LowDiffSign:              c[31],
		LowDiff:                  asDecimal(c[32]),
		BusinessDate:             c[33],
		MarketOpCode:             c[34],
		TradeHaltYN:              c[35],
		Ask1Size:                 asInt64(c[36]),
		Bid1Size:                 asInt64(c[37]),
		TotalAskSize:             asInt64(c[38]),
		TotalBidSize:             asInt64(c[39]),
		VolumeTurnover:           asFloat(c[40]),
		PrevSameTimeAccumVol:     asInt64(c[41]),
		PrevSameTimeAccumVolRate: asFloat(c[42]),
		HourCode:                 c[43],
		MarketTermCode:           c[44],
		ViStandardPrice:          asDecimal(c[45]),
		Raw:                      c,
	}
}

// --------------------------------------------------------------------------
// H0STASP0 — 실시간호가 (KRX, 59 fields)
// --------------------------------------------------------------------------

// decodeKrxAsk 는 H0STASP0 frame 을 KrxAskEvent 슬라이스로 디코딩한다.
func decodeKrxAsk(f frame) ([]KrxAskEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, krxAskFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]KrxAskEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseKrxAskChunk(c))
	}
	return out, nil
}

func parseKrxAskChunk(c []string) KrxAskEvent {
	ev := KrxAskEvent{
		Symbol:   c[0],
		Time:     c[1],
		HourCode: c[2],
		// c[3..12]  Ask[0..9]  매도호가1~10
		// c[13..22] Bid[0..9]  매수호가1~10
		// c[23..32] AskSize[0..9]
		// c[33..42] BidSize[0..9]
		TotalAskSize:            asInt64(c[43]),
		TotalBidSize:            asInt64(c[44]),
		OvernightTotalAskSize:   asInt64(c[45]),
		OvernightTotalBidSize:   asInt64(c[46]),
		ExpectPrice:             asDecimal(c[47]),
		ExpectQuantity:          asInt64(c[48]),
		ExpectVolume:            asInt64(c[49]),
		ExpectDiff:              asDecimal(c[50]),
		ExpectDiffSign:          c[51],
		ExpectChangeRate:        asFloat(c[52]),
		AccumVolume:             asInt64(c[53]),
		TotalAskSizeChange:      asInt64(c[54]),
		TotalBidSizeChange:      asInt64(c[55]),
		OvernightTotalAskChange: asInt64(c[56]),
		OvernightTotalBidChange: asInt64(c[57]),
		DealCode:                c[58],
		Raw:                     c,
	}
	for i := 0; i < 10; i++ {
		ev.Ask[i] = asDecimal(c[3+i])
		ev.Bid[i] = asDecimal(c[13+i])
		ev.AskSize[i] = asInt64(c[23+i])
		ev.BidSize[i] = asInt64(c[33+i])
	}
	return ev
}

// --------------------------------------------------------------------------
// H0STANC0 — 실시간예상체결 (KRX 본장, 45 fields)
// --------------------------------------------------------------------------

// decodeKrxExpectTrade 는 H0STANC0 frame 을 KrxExpectTradeEvent 슬라이스로 디코딩한다.
func decodeKrxExpectTrade(f frame) ([]KrxExpectTradeEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, krxExpectTradeFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]KrxExpectTradeEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseKrxExpectTradeChunk(c))
	}
	return out, nil
}

func parseKrxExpectTradeChunk(c []string) KrxExpectTradeEvent {
	return KrxExpectTradeEvent{
		Symbol:                   c[0],
		Time:                     c[1],
		Price:                    asDecimal(c[2]),
		PrevDiffSign:             c[3],
		PrevDiff:                 asDecimal(c[4]),
		PrevChangeRate:           asFloat(c[5]),
		WeightedAvg:              asDecimal(c[6]),
		Open:                     asDecimal(c[7]),
		High:                     asDecimal(c[8]),
		Low:                      asDecimal(c[9]),
		Ask1:                     asDecimal(c[10]),
		Bid1:                     asDecimal(c[11]),
		TradeVolume:              asInt64(c[12]),
		AccumVolume:              asInt64(c[13]),
		AccumValue:               asInt64(c[14]),
		AskCount:                 asInt64(c[15]),
		BidCount:                 asInt64(c[16]),
		NetCount:                 asInt64(c[17]),
		TradeStrength:            asFloat(c[18]),
		TotalAskVolume:           asInt64(c[19]),
		TotalBidVolume:           asInt64(c[20]),
		TradeKind:                c[21], // CNTG_CLS_CODE
		BidRate:                  asFloat(c[22]),
		PrevVolRate:              asFloat(c[23]),
		OpenTime:                 c[24],
		OpenDiffSign:             c[25],
		OpenDiff:                 asDecimal(c[26]),
		HighTime:                 c[27],
		HighDiffSign:             c[28],
		HighDiff:                 asDecimal(c[29]),
		LowTime:                  c[30],
		LowDiffSign:              c[31],
		LowDiff:                  asDecimal(c[32]),
		BusinessDate:             c[33],
		MarketOpCode:             c[34],
		TradeHaltYN:              c[35],
		Ask1Size:                 asInt64(c[36]),
		Bid1Size:                 asInt64(c[37]),
		TotalAskSize:             asInt64(c[38]),
		TotalBidSize:             asInt64(c[39]),
		VolumeTurnover:           asFloat(c[40]),
		PrevSameTimeAccumVol:     asInt64(c[41]),
		PrevSameTimeAccumVolRate: asFloat(c[42]),
		HourCode:                 c[43],
		MarketTermCode:           c[44],
		// VI_STND_PRC 없음 (H0STCNT0 #46 과 달리)
		Raw: c,
	}
}

// --------------------------------------------------------------------------
// H0STOUP0 — 시간외 실시간체결가 (KRX, 43 fields)
// --------------------------------------------------------------------------

// decodeKrxOvernightTrade 는 H0STOUP0 frame 을 KrxOvernightTradeEvent 슬라이스로 디코딩한다.
func decodeKrxOvernightTrade(f frame) ([]KrxOvernightTradeEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, krxOvernightTradeFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]KrxOvernightTradeEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseKrxOvernightTradeChunk(c))
	}
	return out, nil
}

func parseKrxOvernightTradeChunk(c []string) KrxOvernightTradeEvent {
	return KrxOvernightTradeEvent{
		Symbol:                   c[0],
		Time:                     c[1],
		Price:                    asDecimal(c[2]),
		PrevDiffSign:             c[3],
		PrevDiff:                 asDecimal(c[4]),
		PrevChangeRate:           asFloat(c[5]),
		WeightedAvg:              asDecimal(c[6]),
		Open:                     asDecimal(c[7]),
		High:                     asDecimal(c[8]),
		Low:                      asDecimal(c[9]),
		Ask1:                     asDecimal(c[10]),
		Bid1:                     asDecimal(c[11]),
		TradeVolume:              asInt64(c[12]),
		AccumVolume:              asInt64(c[13]),
		AccumValue:               asInt64(c[14]),
		AskCount:                 asInt64(c[15]),
		BidCount:                 asInt64(c[16]),
		NetCount:                 asInt64(c[17]),
		TradeStrength:            asFloat(c[18]),
		TotalAskVolume:           asInt64(c[19]),
		TotalBidVolume:           asInt64(c[20]),
		TradeKind:                c[21], // CNTG_CLS_CODE
		BidRate:                  asFloat(c[22]),
		PrevVolRate:              asFloat(c[23]),
		OpenTime:                 c[24],
		OpenDiffSign:             c[25],
		OpenDiff:                 asDecimal(c[26]),
		HighTime:                 c[27],
		HighDiffSign:             c[28],
		HighDiff:                 asDecimal(c[29]),
		LowTime:                  c[30],
		LowDiffSign:              c[31],
		LowDiff:                  asDecimal(c[32]),
		BusinessDate:             c[33],
		MarketOpCode:             c[34],
		TradeHaltYN:              c[35],
		Ask1Size:                 asInt64(c[36]),
		Bid1Size:                 asInt64(c[37]),
		TotalAskSize:             asInt64(c[38]),
		TotalBidSize:             asInt64(c[39]),
		VolumeTurnover:           asFloat(c[40]),
		PrevSameTimeAccumVol:     asInt64(c[41]),
		PrevSameTimeAccumVolRate: asFloat(c[42]),
		// HOUR_CLS_CODE / MRKT_TRTM_CLS_CODE / VI_STND_PRC 없음
		Raw: c,
	}
}

// --------------------------------------------------------------------------
// H0STOAC0 — 시간외 실시간예상체결 (KRX, 43 fields)
// --------------------------------------------------------------------------

// decodeKrxOvernightExpect 는 H0STOAC0 frame 을 KrxOvernightExpectEvent 슬라이스로 디코딩한다.
func decodeKrxOvernightExpect(f frame) ([]KrxOvernightExpectEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, krxOvernightExpectFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]KrxOvernightExpectEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseKrxOvernightExpectChunk(c))
	}
	return out, nil
}

func parseKrxOvernightExpectChunk(c []string) KrxOvernightExpectEvent {
	return KrxOvernightExpectEvent{
		Symbol:                   c[0],
		Time:                     c[1],
		Price:                    asDecimal(c[2]),
		PrevDiffSign:             c[3],
		PrevDiff:                 asDecimal(c[4]),
		PrevChangeRate:           asFloat(c[5]),
		WeightedAvg:              asDecimal(c[6]),
		Open:                     asDecimal(c[7]),
		High:                     asDecimal(c[8]),
		Low:                      asDecimal(c[9]),
		Ask1:                     asDecimal(c[10]),
		Bid1:                     asDecimal(c[11]),
		TradeVolume:              asInt64(c[12]),
		AccumVolume:              asInt64(c[13]),
		AccumValue:               asInt64(c[14]),
		AskCount:                 asInt64(c[15]),
		BidCount:                 asInt64(c[16]),
		NetCount:                 asInt64(c[17]),
		TradeStrength:            asFloat(c[18]),
		TotalAskVolume:           asInt64(c[19]),
		TotalBidVolume:           asInt64(c[20]),
		TradeKind:                c[21], // CNTG_CLS_CODE
		BidRate:                  asFloat(c[22]),
		PrevVolRate:              asFloat(c[23]),
		OpenTime:                 c[24],
		OpenDiffSign:             c[25],
		OpenDiff:                 asDecimal(c[26]),
		HighTime:                 c[27],
		HighDiffSign:             c[28],
		HighDiff:                 asDecimal(c[29]),
		LowTime:                  c[30],
		LowDiffSign:              c[31],
		LowDiff:                  asDecimal(c[32]),
		BusinessDate:             c[33],
		MarketOpCode:             c[34],
		TradeHaltYN:              c[35],
		Ask1Size:                 asInt64(c[36]),
		Bid1Size:                 asInt64(c[37]),
		TotalAskSize:             asInt64(c[38]),
		TotalBidSize:             asInt64(c[39]),
		VolumeTurnover:           asFloat(c[40]),
		PrevSameTimeAccumVol:     asInt64(c[41]),
		PrevSameTimeAccumVolRate: asFloat(c[42]),
		// HOUR_CLS_CODE / MRKT_TRTM_CLS_CODE 없음
		Raw: c,
	}
}
