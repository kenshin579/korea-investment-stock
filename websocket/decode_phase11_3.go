package websocket

// Phase 11.3 — 국내선물옵션 실시간 6 EP decoder.
// 4 base decoders. 상품선물(CF) 은 alias 이므로 IF decoder 를 재사용.
// asDecimal/asInt64/asFloat 헬퍼는 decode_krx.go 정의.

const (
	// 지수선물
	indexFuturesTradeFieldCount = 50 // H0IFCNT0 / H0CFCNT0 (alias)
	indexFuturesAskFieldCount   = 38 // H0IFASP0 / H0CFASP0 (alias)

	// 지수옵션
	indexOptionTradeFieldCount = 58 // H0IOCNT0
	indexOptionAskFieldCount   = 38 // H0IOASP0
)

// --------------------------------------------------------------------------
// H0IFCNT0 / H0CFCNT0 — 지수선물/상품선물 실시간체결가 (50 fields)
// --------------------------------------------------------------------------

// decodeIndexFuturesTrade 는 H0IFCNT0 frame 을 IndexFuturesTradeEvent 슬라이스로 디코딩한다.
// H0CFCNT0 (상품선물) 도 동일 decoder 사용 (alias schema).
func decodeIndexFuturesTrade(f frame) ([]IndexFuturesTradeEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, indexFuturesTradeFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]IndexFuturesTradeEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseIndexFuturesTradeChunk(c))
	}
	return out, nil
}

func parseIndexFuturesTradeChunk(c []string) IndexFuturesTradeEvent {
	return IndexFuturesTradeEvent{
		Symbol:                 c[0],
		Time:                   c[1],
		PrevDiff:               asDecimal(c[2]),
		PrevDiffSign:           c[3],
		PrevChangeRate:         asFloat(c[4]),
		Price:                  asDecimal(c[5]),
		Open:                   asDecimal(c[6]),
		High:                   asDecimal(c[7]),
		Low:                    asDecimal(c[8]),
		LastTradeVolume:        asInt64(c[9]),
		AccumVolume:            asInt64(c[10]),
		AccumValue:             asInt64(c[11]),
		TheoreticalPrice:       asDecimal(c[12]),
		MarketBasis:            asDecimal(c[13]),
		DeviationRate:          asFloat(c[14]),
		NearMonthPrice:         asDecimal(c[15]),
		FarMonthPrice:          asDecimal(c[16]),
		SpreadPrice:            asDecimal(c[17]),
		OpenInterestQty:        asInt64(c[18]),
		OpenInterestChange:     asInt64(c[19]),
		OpenTime:               c[20],
		OpenDiffSign:           c[21],
		OpenDiff:               asDecimal(c[22]),
		HighTime:               c[23],
		HighDiffSign:           c[24],
		HighDiff:               asDecimal(c[25]),
		LowTime:                c[26],
		LowDiffSign:            c[27],
		LowDiff:                asDecimal(c[28]),
		BidRate:                asFloat(c[29]),
		TradeStrength:          asFloat(c[30]),
		DeviationDegree:        asDecimal(c[31]),
		OpenInterestPrevChange: asInt64(c[32]),
		TheoreticalBasis:       asDecimal(c[33]),
		Ask1:                   asDecimal(c[34]),
		Bid1:                   asDecimal(c[35]),
		Ask1Size:               asInt64(c[36]),
		Bid1Size:               asInt64(c[37]),
		AskCount:               asInt64(c[38]),
		BidCount:               asInt64(c[39]),
		NetCount:               asInt64(c[40]),
		TotalAskVolume:         asInt64(c[41]),
		TotalBidVolume:         asInt64(c[42]),
		TotalAskSize:           asInt64(c[43]),
		TotalBidSize:           asInt64(c[44]),
		PrevVolRate:            asFloat(c[45]),
		BlockTradeVolume:       asInt64(c[46]),
		DynamicUpperLimit:      asDecimal(c[47]),
		DynamicLowerLimit:      asDecimal(c[48]),
		DynamicPriceLimitYN:    c[49],
		Raw:                    c,
	}
}

// --------------------------------------------------------------------------
// H0IFASP0 / H0CFASP0 — 지수선물/상품선물 실시간호가 (38 fields)
// --------------------------------------------------------------------------

// decodeIndexFuturesAsk 는 H0IFASP0 frame 을 IndexFuturesAskEvent 슬라이스로 디코딩한다.
// H0CFASP0 (상품선물) 도 동일 decoder 사용 (alias schema).
func decodeIndexFuturesAsk(f frame) ([]IndexFuturesAskEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, indexFuturesAskFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]IndexFuturesAskEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseIndexFuturesAskChunk(c))
	}
	return out, nil
}

func parseIndexFuturesAskChunk(c []string) IndexFuturesAskEvent {
	ev := IndexFuturesAskEvent{
		Symbol:          c[0],
		Time:            c[1],
		TotalAskCsnu:    asInt64(c[32]),
		TotalBidCsnu:    asInt64(c[33]),
		TotalAskSize:    asInt64(c[34]),
		TotalBidSize:    asInt64(c[35]),
		TotalAskSizeChg: asInt64(c[36]),
		TotalBidSizeChg: asInt64(c[37]),
		Raw:             c,
	}
	for i := 0; i < 5; i++ {
		ev.Ask[i] = asDecimal(c[2+i])    // FUTS_ASKP1..5
		ev.Bid[i] = asDecimal(c[7+i])    // FUTS_BIDP1..5
		ev.AskCsnu[i] = asInt64(c[12+i]) // ASKP_CSNU1..5
		ev.BidCsnu[i] = asInt64(c[17+i]) // BIDP_CSNU1..5
		ev.AskSize[i] = asInt64(c[22+i]) // ASKP_RSQN1..5
		ev.BidSize[i] = asInt64(c[27+i]) // BIDP_RSQN1..5
	}
	return ev
}

// --------------------------------------------------------------------------
// H0IOCNT0 — 지수옵션 실시간체결가 (58 fields)
// --------------------------------------------------------------------------

// decodeIndexOptionTrade 는 H0IOCNT0 frame 을 IndexOptionTradeEvent 슬라이스로 디코딩한다.
func decodeIndexOptionTrade(f frame) ([]IndexOptionTradeEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, indexOptionTradeFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]IndexOptionTradeEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseIndexOptionTradeChunk(c))
	}
	return out, nil
}

func parseIndexOptionTradeChunk(c []string) IndexOptionTradeEvent {
	return IndexOptionTradeEvent{
		Symbol:                 c[0],
		Time:                   c[1],
		Price:                  asDecimal(c[2]),
		PrevDiffSign:           c[3],
		PrevDiff:               asDecimal(c[4]),
		PrevChangeRate:         asFloat(c[5]),
		Open:                   asDecimal(c[6]),
		High:                   asDecimal(c[7]),
		Low:                    asDecimal(c[8]),
		LastTradeVolume:        asInt64(c[9]),
		AccumVolume:            asInt64(c[10]),
		AccumValue:             asInt64(c[11]),
		TheoreticalPrice:       asDecimal(c[12]),
		OpenInterestQty:        asInt64(c[13]),
		OpenInterestChange:     asInt64(c[14]),
		OpenTime:               c[15],
		OpenDiffSign:           c[16],
		OpenDiff:               asDecimal(c[17]),
		HighTime:               c[18],
		HighDiffSign:           c[19],
		HighDiff:               asDecimal(c[20]),
		LowTime:                c[21],
		LowDiffSign:            c[22],
		LowDiff:                asDecimal(c[23]),
		BidRate:                asFloat(c[24]),
		PremiumValue:           asDecimal(c[25]),
		IntrinsicValue:         asDecimal(c[26]),
		TimeValue:              asDecimal(c[27]),
		Delta:                  asFloat(c[28]),
		Gamma:                  asFloat(c[29]),
		Vega:                   asFloat(c[30]),
		Theta:                  asFloat(c[31]),
		Rho:                    asFloat(c[32]),
		ImpliedVolatility:      asFloat(c[33]),
		DeviationDegree:        asDecimal(c[34]),
		OpenInterestPrevChange: asInt64(c[35]),
		TheoreticalBasis:       asDecimal(c[36]),
		HistoricalVolatility:   asFloat(c[37]),
		TradeStrength:          asFloat(c[38]),
		DeviationRate:          asFloat(c[39]),
		MarketBasis:            asDecimal(c[40]),
		Ask1:                   asDecimal(c[41]),
		Bid1:                   asDecimal(c[42]),
		Ask1Size:               asInt64(c[43]),
		Bid1Size:               asInt64(c[44]),
		AskCount:               asInt64(c[45]),
		BidCount:               asInt64(c[46]),
		NetCount:               asInt64(c[47]),
		TotalAskVolume:         asInt64(c[48]),
		TotalBidVolume:         asInt64(c[49]),
		TotalAskSize:           asInt64(c[50]),
		TotalBidSize:           asInt64(c[51]),
		PrevVolRate:            asFloat(c[52]),
		AvgVolatility:          asFloat(c[53]),
		BlockTradeVolume:       asInt64(c[54]),
		DynamicUpperLimit:      asDecimal(c[55]),
		DynamicLowerLimit:      asDecimal(c[56]),
		DynamicPriceLimitYN:    c[57],
		Raw:                    c,
	}
}

// --------------------------------------------------------------------------
// H0IOASP0 — 지수옵션 실시간호가 (38 fields)
// --------------------------------------------------------------------------

// decodeIndexOptionAsk 는 H0IOASP0 frame 을 IndexOptionAskEvent 슬라이스로 디코딩한다.
func decodeIndexOptionAsk(f frame) ([]IndexOptionAskEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, indexOptionAskFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]IndexOptionAskEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseIndexOptionAskChunk(c))
	}
	return out, nil
}

func parseIndexOptionAskChunk(c []string) IndexOptionAskEvent {
	ev := IndexOptionAskEvent{
		Symbol:          c[0],
		Time:            c[1],
		TotalAskCsnu:    asInt64(c[32]),
		TotalBidCsnu:    asInt64(c[33]),
		TotalAskSize:    asInt64(c[34]),
		TotalBidSize:    asInt64(c[35]),
		TotalAskSizeChg: asInt64(c[36]),
		TotalBidSizeChg: asInt64(c[37]),
		Raw:             c,
	}
	for i := 0; i < 5; i++ {
		ev.Ask[i] = asDecimal(c[2+i])    // OPTN_ASKP1..5
		ev.Bid[i] = asDecimal(c[7+i])    // OPTN_BIDP1..5
		ev.AskCsnu[i] = asInt64(c[12+i]) // ASKP_CSNU1..5
		ev.BidCsnu[i] = asInt64(c[17+i]) // BIDP_CSNU1..5
		ev.AskSize[i] = asInt64(c[22+i]) // ASKP_RSQN1..5
		ev.BidSize[i] = asInt64(c[27+i]) // BIDP_RSQN1..5
	}
	return ev
}
