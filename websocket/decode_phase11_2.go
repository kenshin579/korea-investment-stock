package websocket

// Phase 11.2 — 국내선물옵션 실시간 11 EP decoder.
// 11 EP 모두 Distinct (alias 없음). asDecimal/asInt64/asFloat 헬퍼는 decode_krx.go 정의.

const (
	// KRX 야간 선물
	krxNightFuturesTradeFieldCount = 49 // H0MFCNT0
	krxNightFuturesAskFieldCount   = 38 // H0MFASP0

	// KRX 야간 옵션
	krxNightOptionTradeFieldCount       = 56 // H0EUCNT0
	krxNightOptionAskFieldCount         = 38 // H0EUASP0
	krxNightOptionExpectTradeFieldCount = 8  // H0EUANC0

	// 주식 선물
	stockFuturesTradeFieldCount       = 49 // H0ZFCNT0
	stockFuturesAskFieldCount         = 68 // H0ZFASP0
	stockFuturesExpectTradeFieldCount = 8  // H0ZFANC0

	// 주식 옵션
	stockOptionTradeFieldCount       = 53 // H0ZOCNT0
	stockOptionAskFieldCount         = 68 // H0ZOASP0
	stockOptionExpectTradeFieldCount = 7  // H0ZOANC0
)

// --------------------------------------------------------------------------
// H0MFCNT0 — KRX야간선물 실시간종목체결 (49 fields)
// --------------------------------------------------------------------------

// decodeKrxNightFuturesTrade 는 H0MFCNT0 frame 을 KrxNightFuturesTradeEvent 슬라이스로 디코딩한다.
func decodeKrxNightFuturesTrade(f frame) ([]KrxNightFuturesTradeEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, krxNightFuturesTradeFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]KrxNightFuturesTradeEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseKrxNightFuturesTradeChunk(c))
	}
	return out, nil
}

func parseKrxNightFuturesTradeChunk(c []string) KrxNightFuturesTradeEvent {
	return KrxNightFuturesTradeEvent{
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
		DynamicUpperLimit:      asDecimal(c[46]),
		DynamicLowerLimit:      asDecimal(c[47]),
		DynamicPriceLimitYN:    c[48],
		Raw:                    c,
	}
}

// --------------------------------------------------------------------------
// H0MFASP0 — KRX야간선물 실시간호가 (38 fields)
// --------------------------------------------------------------------------

// decodeKrxNightFuturesAsk 는 H0MFASP0 frame 을 KrxNightFuturesAskEvent 슬라이스로 디코딩한다.
func decodeKrxNightFuturesAsk(f frame) ([]KrxNightFuturesAskEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, krxNightFuturesAskFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]KrxNightFuturesAskEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseKrxNightFuturesAskChunk(c))
	}
	return out, nil
}

func parseKrxNightFuturesAskChunk(c []string) KrxNightFuturesAskEvent {
	ev := KrxNightFuturesAskEvent{
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
// H0EUCNT0 — KRX야간옵션 실시간체결가 (56 fields)
// --------------------------------------------------------------------------

// decodeKrxNightOptionTrade 는 H0EUCNT0 frame 을 KrxNightOptionTradeEvent 슬라이스로 디코딩한다.
func decodeKrxNightOptionTrade(f frame) ([]KrxNightOptionTradeEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, krxNightOptionTradeFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]KrxNightOptionTradeEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseKrxNightOptionTradeChunk(c))
	}
	return out, nil
}

func parseKrxNightOptionTradeChunk(c []string) KrxNightOptionTradeEvent {
	return KrxNightOptionTradeEvent{
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
		DynamicUpperLimit:      asDecimal(c[53]),
		DynamicPriceLimitYN:    c[54], // docs: MXPR→PRC_LIMT_YN→LLAM 순서 (anomaly)
		DynamicLowerLimit:      asDecimal(c[55]),
		Raw:                    c,
	}
}

// --------------------------------------------------------------------------
// H0EUASP0 — KRX야간옵션 실시간호가 (38 fields)
// --------------------------------------------------------------------------

// decodeKrxNightOptionAsk 는 H0EUASP0 frame 을 KrxNightOptionAskEvent 슬라이스로 디코딩한다.
func decodeKrxNightOptionAsk(f frame) ([]KrxNightOptionAskEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, krxNightOptionAskFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]KrxNightOptionAskEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseKrxNightOptionAskChunk(c))
	}
	return out, nil
}

func parseKrxNightOptionAskChunk(c []string) KrxNightOptionAskEvent {
	ev := KrxNightOptionAskEvent{
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

// --------------------------------------------------------------------------
// H0EUANC0 — KRX야간옵션 실시간예상체결 (8 fields)
// --------------------------------------------------------------------------

// decodeKrxNightOptionExpectTrade 는 H0EUANC0 frame 을 KrxNightOptionExpectTradeEvent 슬라이스로 디코딩한다.
func decodeKrxNightOptionExpectTrade(f frame) ([]KrxNightOptionExpectTradeEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, krxNightOptionExpectTradeFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]KrxNightOptionExpectTradeEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseKrxNightOptionExpectTradeChunk(c))
	}
	return out, nil
}

func parseKrxNightOptionExpectTradeChunk(c []string) KrxNightOptionExpectTradeEvent {
	return KrxNightOptionExpectTradeEvent{
		Symbol:           c[0],
		Time:             c[1],
		ExpectPrice:      asDecimal(c[2]),
		ExpectDiff:       asDecimal(c[3]),
		ExpectDiffSign:   c[4],
		ExpectChangeRate: asFloat(c[5]),
		ExpectMarketCode: c[6],
		ExpectQuantity:   asInt64(c[7]),
		Raw:              c,
	}
}

// --------------------------------------------------------------------------
// H0ZFCNT0 — 주식선물 실시간체결가 (49 fields)
// --------------------------------------------------------------------------

// decodeStockFuturesTrade 는 H0ZFCNT0 frame 을 StockFuturesTradeEvent 슬라이스로 디코딩한다.
func decodeStockFuturesTrade(f frame) ([]StockFuturesTradeEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, stockFuturesTradeFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]StockFuturesTradeEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseStockFuturesTradeChunk(c))
	}
	return out, nil
}

func parseStockFuturesTradeChunk(c []string) StockFuturesTradeEvent {
	return StockFuturesTradeEvent{
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
		DynamicUpperLimit:      asDecimal(c[46]),
		DynamicLowerLimit:      asDecimal(c[47]),
		DynamicPriceLimitYN:    c[48],
		Raw:                    c,
	}
}

// --------------------------------------------------------------------------
// H0ZFASP0 — 주식선물 실시간호가 (68 fields)
// --------------------------------------------------------------------------

// decodeStockFuturesAsk 는 H0ZFASP0 frame 을 StockFuturesAskEvent 슬라이스로 디코딩한다.
func decodeStockFuturesAsk(f frame) ([]StockFuturesAskEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, stockFuturesAskFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]StockFuturesAskEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseStockFuturesAskChunk(c))
	}
	return out, nil
}

func parseStockFuturesAskChunk(c []string) StockFuturesAskEvent {
	ev := StockFuturesAskEvent{
		Symbol:          c[0],
		Time:            c[1],
		TotalAskCsnu:    asInt64(c[62]),
		TotalBidCsnu:    asInt64(c[63]),
		TotalAskSize:    asInt64(c[64]),
		TotalBidSize:    asInt64(c[65]),
		TotalAskSizeChg: asInt64(c[66]),
		TotalBidSizeChg: asInt64(c[67]),
		Raw:             c,
	}
	for i := 0; i < 10; i++ {
		ev.Ask[i] = asDecimal(c[2+i])    // ASKP1..10
		ev.Bid[i] = asDecimal(c[12+i])   // BIDP1..10
		ev.AskCsnu[i] = asInt64(c[22+i]) // ASKP_CSNU1..10
		ev.BidCsnu[i] = asInt64(c[32+i]) // BIDP_CSNU1..10
		ev.AskSize[i] = asInt64(c[42+i]) // ASKP_RSQN1..10
		ev.BidSize[i] = asInt64(c[52+i]) // BIDP_RSQN1..10
	}
	return ev
}

// --------------------------------------------------------------------------
// H0ZFANC0 — 주식선물 실시간예상체결 (8 fields)
// --------------------------------------------------------------------------

// decodeStockFuturesExpectTrade 는 H0ZFANC0 frame 을 StockFuturesExpectTradeEvent 슬라이스로 디코딩한다.
func decodeStockFuturesExpectTrade(f frame) ([]StockFuturesExpectTradeEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, stockFuturesExpectTradeFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]StockFuturesExpectTradeEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseStockFuturesExpectTradeChunk(c))
	}
	return out, nil
}

func parseStockFuturesExpectTradeChunk(c []string) StockFuturesExpectTradeEvent {
	return StockFuturesExpectTradeEvent{
		Symbol:           c[0],
		Time:             c[1],
		ExpectPrice:      asDecimal(c[2]),
		ExpectDiff:       asDecimal(c[3]),
		ExpectDiffSign:   c[4],
		ExpectChangeRate: asFloat(c[5]),
		ExpectMarketCode: c[6],
		ExpectQuantity:   asInt64(c[7]),
		Raw:              c,
	}
}

// --------------------------------------------------------------------------
// H0ZOCNT0 — 주식옵션 실시간체결가 (53 fields)
// --------------------------------------------------------------------------

// decodeStockOptionTrade 는 H0ZOCNT0 frame 을 StockOptionTradeEvent 슬라이스로 디코딩한다.
func decodeStockOptionTrade(f frame) ([]StockOptionTradeEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, stockOptionTradeFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]StockOptionTradeEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseStockOptionTradeChunk(c))
	}
	return out, nil
}

func parseStockOptionTradeChunk(c []string) StockOptionTradeEvent {
	return StockOptionTradeEvent{
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
		Raw:                    c,
	}
}

// --------------------------------------------------------------------------
// H0ZOASP0 — 주식옵션 실시간호가 (68 fields)
// --------------------------------------------------------------------------

// decodeStockOptionAsk 는 H0ZOASP0 frame 을 StockOptionAskEvent 슬라이스로 디코딩한다.
func decodeStockOptionAsk(f frame) ([]StockOptionAskEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, stockOptionAskFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]StockOptionAskEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseStockOptionAskChunk(c))
	}
	return out, nil
}

func parseStockOptionAskChunk(c []string) StockOptionAskEvent {
	ev := StockOptionAskEvent{
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
	// 1~5단계 (index 2..31)
	for i := 0; i < 5; i++ {
		ev.Ask1to5[i] = asDecimal(c[2+i])    // OPTN_ASKP1..5
		ev.Bid1to5[i] = asDecimal(c[7+i])    // OPTN_BIDP1..5
		ev.AskCsnu1to5[i] = asInt64(c[12+i]) // ASKP_CSNU1..5
		ev.BidCsnu1to5[i] = asInt64(c[17+i]) // BIDP_CSNU1..5
		ev.AskSize1to5[i] = asInt64(c[22+i]) // ASKP_RSQN1..5
		ev.BidSize1to5[i] = asInt64(c[27+i]) // BIDP_RSQN1..5
	}
	// 6~10단계 (index 38..67)
	for i := 0; i < 5; i++ {
		ev.Ask6to10[i] = asDecimal(c[38+i])   // OPTN_ASKP6..10
		ev.Bid6to10[i] = asDecimal(c[43+i])   // OPTN_BIDP6..10
		ev.AskCsnu6to10[i] = asInt64(c[48+i]) // ASKP_CSNU6..10
		ev.BidCsnu6to10[i] = asInt64(c[53+i]) // BIDP_CSNU6..10
		ev.AskSize6to10[i] = asInt64(c[58+i]) // ASKP_RSQN6..10
		ev.BidSize6to10[i] = asInt64(c[63+i]) // BIDP_RSQN6..10
	}
	return ev
}

// --------------------------------------------------------------------------
// H0ZOANC0 — 주식옵션 실시간예상체결 (7 fields)
// --------------------------------------------------------------------------

// decodeStockOptionExpectTrade 는 H0ZOANC0 frame 을 StockOptionExpectTradeEvent 슬라이스로 디코딩한다.
// 주의: ANTC_CNQN 없음 — 7 fields (다른 예상체결 EP 는 8 fields).
func decodeStockOptionExpectTrade(f frame) ([]StockOptionExpectTradeEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, stockOptionExpectTradeFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]StockOptionExpectTradeEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseStockOptionExpectTradeChunk(c))
	}
	return out, nil
}

func parseStockOptionExpectTradeChunk(c []string) StockOptionExpectTradeEvent {
	return StockOptionExpectTradeEvent{
		Symbol:           c[0],
		Time:             c[1],
		ExpectPrice:      asDecimal(c[2]),
		ExpectDiff:       asDecimal(c[3]),
		ExpectDiffSign:   c[4],
		ExpectChangeRate: asFloat(c[5]),
		ExpectMarketCode: c[6],
		Raw:              c,
	}
}
