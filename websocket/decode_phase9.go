package websocket

// Phase 9 — NXT/통합 실시간 변형 5 base decoder.
// NXT 와 통합은 schema 가 완전히 동일해서 base decoder 1개를 공유한다.
// 시장 구분은 client.routeRealtime 의 TR_ID 분기에서 결정.

const (
	altMarketTradeFieldCount       = 46 // H0NXCNT0 / H0UNCNT0
	altMarketAskFieldCount         = 65 // H0NXASP0 / H0UNASP0
	altMarketExpectTradeFieldCount = 46 // H0NXANC0 / H0UNANC0
	programTradeFieldCount         = 11 // H0NXPGM0 / H0UNPGM0
	memberFieldCount               = 78 // H0NXMBC0 / H0UNMBC0
)

// --------------------------------------------------------------------------
// H0NXCNT0 / H0UNCNT0 — NXT/통합 실시간체결가 (46 fields)
// --------------------------------------------------------------------------

func decodeAltMarketTrade(f frame) ([]AltMarketTradeEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, altMarketTradeFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]AltMarketTradeEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseAltMarketTradeChunk(c))
	}
	return out, nil
}

func parseAltMarketTradeChunk(c []string) AltMarketTradeEvent {
	return AltMarketTradeEvent{
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
		ViStandardPrice:          asDecimal(c[45]),
		Raw:                      c,
	}
}

// --------------------------------------------------------------------------
// H0NXASP0 / H0UNASP0 — NXT/통합 실시간호가 (65 fields)
// --------------------------------------------------------------------------

func decodeAltMarketAsk(f frame) ([]AltMarketAskEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, altMarketAskFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]AltMarketAskEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseAltMarketAskChunk(c))
	}
	return out, nil
}

func parseAltMarketAskChunk(c []string) AltMarketAskEvent {
	ev := AltMarketAskEvent{
		Symbol:                  c[0],
		Time:                    c[1],
		HourCode:                c[2],
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
		KrxMidPrice:             asDecimal(c[59]),
		KrxMidTotalSize:         asInt64(c[60]),
		KrxMidCode:              c[61],
		NxtMidPrice:             asDecimal(c[62]),
		NxtMidTotalSize:         asInt64(c[63]),
		NxtMidCode:              c[64],
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
// H0NXANC0 / H0UNANC0 — NXT/통합 실시간예상체결 (46 fields)
// --------------------------------------------------------------------------

func decodeAltMarketExpectTrade(f frame) ([]AltMarketExpectTradeEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, altMarketExpectTradeFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]AltMarketExpectTradeEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseAltMarketExpectTradeChunk(c))
	}
	return out, nil
}

func parseAltMarketExpectTradeChunk(c []string) AltMarketExpectTradeEvent {
	return AltMarketExpectTradeEvent{
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
		ViStandardPrice:          asDecimal(c[45]),
		Raw:                      c,
	}
}

// --------------------------------------------------------------------------
// H0NXPGM0 / H0UNPGM0 — NXT/통합 실시간프로그램매매 (11 fields)
// --------------------------------------------------------------------------

func decodeProgramTrade(f frame) ([]ProgramTradeEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, programTradeFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]ProgramTradeEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseProgramTradeChunk(c))
	}
	return out, nil
}

func parseProgramTradeChunk(c []string) ProgramTradeEvent {
	return ProgramTradeEvent{
		Symbol:           c[0],
		Time:             c[1],
		AskQuantity:      asInt64(c[2]),
		AskValue:         asInt64(c[3]),
		BidQuantity:      asInt64(c[4]),
		BidValue:         asInt64(c[5]),
		NetQuantity:      asInt64(c[6]),
		NetValue:         asInt64(c[7]),
		AskRemainingSize: asInt64(c[8]),
		BidRemainingSize: asInt64(c[9]),
		TotalNetQuantity: asInt64(c[10]),
		Raw:              c,
	}
}

// --------------------------------------------------------------------------
// H0NXMBC0 / H0UNMBC0 — NXT/통합 실시간회원사 (78 fields)
// --------------------------------------------------------------------------

func decodeMember(f frame) ([]MemberEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, memberFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]MemberEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseMemberChunk(c))
	}
	return out, nil
}

func parseMemberChunk(c []string) MemberEvent {
	ev := MemberEvent{
		Symbol:              c[0],
		GlobalTotalSellQty:  asInt64(c[61]),
		GlobalTotalBuyQty:   asInt64(c[62]),
		GlobalSellQtyChange: asInt64(c[63]),
		GlobalBuyQtyChange:  asInt64(c[64]),
		GlobalNetBuyQty:     asInt64(c[65]),
		GlobalSellRatio:     asFloat(c[66]),
		GlobalBuyRatio:      asFloat(c[67]),
		Raw:                 c,
	}
	for i := 0; i < 5; i++ {
		ev.SellBrokerNames[i] = c[1+i]         // 1..5
		ev.BuyBrokerNames[i] = c[6+i]          // 6..10
		ev.TotalSellQty[i] = asInt64(c[11+i])  // 11..15
		ev.TotalBuyQty[i] = asInt64(c[16+i])   // 16..20
		ev.SellGlobalYN[i] = c[21+i]           // 21..25
		ev.BuyGlobalYN[i] = c[26+i]            // 26..30
		ev.SellBrokerCodes[i] = c[31+i]        // 31..35
		ev.BuyBrokerCodes[i] = c[36+i]         // 36..40
		ev.SellRatio[i] = asFloat(c[41+i])     // 41..45
		ev.BuyRatio[i] = asFloat(c[46+i])      // 46..50
		ev.SellQtyChange[i] = asInt64(c[51+i]) // 51..55
		ev.BuyQtyChange[i] = asInt64(c[56+i])  // 56..60
		ev.SellBrokerEngNames[i] = c[68+i]     // 68..72
		ev.BuyBrokerEngNames[i] = c[73+i]      // 73..77
	}
	return ev
}
