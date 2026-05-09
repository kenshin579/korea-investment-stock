package websocket

// Phase 10 — 해외주식 실시간 시세 2 EP decoder.

const (
	overseasTradeFieldCount = 26 // HDFSCNT0
	overseasAskFieldCount   = 17 // HDFSASP0
)

// --------------------------------------------------------------------------
// HDFSCNT0 — 해외주식 실시간지연체결가 (26 fields)
// --------------------------------------------------------------------------

func decodeOverseasTrade(f frame) ([]OverseasTradeEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, overseasTradeFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]OverseasTradeEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseOverseasTradeChunk(c))
	}
	return out, nil
}

func parseOverseasTradeChunk(c []string) OverseasTradeEvent {
	return OverseasTradeEvent{
		Symbol:        c[0],
		SymbolCode:    c[1],
		Decimals:      c[2],
		LocalDate:     c[3],
		LocalDayDate:  c[4],
		LocalTime:     c[5],
		KrDate:        c[6],
		KrTime:        c[7],
		Open:          asDecimal(c[8]),
		High:          asDecimal(c[9]),
		Low:           asDecimal(c[10]),
		Last:          asDecimal(c[11]),
		PrevDiffSign:  c[12],
		PrevDiff:      asDecimal(c[13]),
		ChangeRate:    asFloat(c[14]),
		Bid:           asDecimal(c[15]),
		Ask:           asDecimal(c[16]),
		BidSize:       asInt64(c[17]),
		AskSize:       asInt64(c[18]),
		TradeVolume:   asInt64(c[19]),
		AccumVolume:   asInt64(c[20]),
		AccumValue:    asInt64(c[21]),
		AskTradeVol:   asInt64(c[22]),
		BidTradeVol:   asInt64(c[23]),
		TradeStrength: asFloat(c[24]),
		MarketKind:    c[25],
		Raw:           c,
	}
}

// --------------------------------------------------------------------------
// HDFSASP0 — 해외주식 실시간호가 (17 fields, 1호가만)
// --------------------------------------------------------------------------

func decodeOverseasAsk(f frame) ([]OverseasAskEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, overseasAskFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]OverseasAskEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseOverseasAskChunk(c))
	}
	return out, nil
}

func parseOverseasAskChunk(c []string) OverseasAskEvent {
	return OverseasAskEvent{
		Symbol:             c[0],
		SymbolCode:         c[1],
		Decimals:           c[2],
		LocalDayDate:       c[3],
		LocalTime:          c[4],
		KrDate:             c[5],
		KrTime:             c[6],
		TotalBidSize:       asInt64(c[7]),
		TotalAskSize:       asInt64(c[8]),
		TotalBidSizeChange: asInt64(c[9]),
		TotalAskSizeChange: asInt64(c[10]),
		Bid1:               asDecimal(c[11]),
		Ask1:               asDecimal(c[12]),
		Bid1Size:           asInt64(c[13]),
		Ask1Size:           asInt64(c[14]),
		Bid1SizeChange:     asInt64(c[15]),
		Ask1SizeChange:     asInt64(c[16]),
		Raw:                c,
	}
}
