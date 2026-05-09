package websocket

// Phase 11.7 — 해외선물옵션 실시간 2 EP decoder.
// asDecimal/asInt64/asFloat 헬퍼는 decode_krx.go 정의.

const (
	// 해외선물옵션 실시간체결가
	overseasFuturesTradeFieldCount = 25 // HDFFF020

	// 해외선물옵션 실시간호가
	overseasFuturesAskFieldCount = 35 // HDFFF010
)

// --------------------------------------------------------------------------
// HDFFF020 — 해외선물옵션 실시간체결가 (25 fields)
// --------------------------------------------------------------------------

// decodeOverseasFuturesTrade 는 HDFFF020 frame 을 OverseasFuturesTradeEvent 슬라이스로 디코딩한다.
func decodeOverseasFuturesTrade(f frame) ([]OverseasFuturesTradeEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, overseasFuturesTradeFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]OverseasFuturesTradeEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseOverseasFuturesTradeChunk(c))
	}
	return out, nil
}

func parseOverseasFuturesTradeChunk(c []string) OverseasFuturesTradeEvent {
	return OverseasFuturesTradeEvent{
		Symbol:         c[0],
		BsnsDate:       c[1],
		MrktOpenDate:   c[2],
		MrktOpenTime:   c[3],
		MrktCloseDate:  c[4],
		MrktCloseTime:  c[5],
		PrevPrice:      asDecimal(c[6]),
		RecvDate:       c[7],
		RecvTime:       c[8],
		ActiveFlag:     c[9],
		LastPrice:      asDecimal(c[10]),
		LastQntt:       asInt64(c[11]),
		PrevDiffPrice:  asDecimal(c[12]),
		PrevDiffRate:   asFloat(c[13]),
		OpenPrice:      asDecimal(c[14]),
		HighPrice:      asDecimal(c[15]),
		LowPrice:       asDecimal(c[16]),
		Vol:            asInt64(c[17]),
		PrevSign:       c[18],
		QuotSign:       c[19],
		RecvTime2:      c[20],
		PsttlPrice:     asDecimal(c[21]),
		PsttlSign:      c[22],
		PsttlDiffPrice: asDecimal(c[23]),
		PsttlDiffRate:  asFloat(c[24]),
		Raw:            c,
	}
}

// --------------------------------------------------------------------------
// HDFFF010 — 해외선물옵션 실시간호가 (35 fields)
// BID/ASK 교차 배열: 각 단계(i=0..4) → offset = 4 + i*6
//   BID_QNTT_i+1 = c[4+i*6], BID_NUM_i+1 = c[5+i*6], BID_PRICE_i+1 = c[6+i*6]
//   ASK_QNTT_i+1 = c[7+i*6], ASK_NUM_i+1 = c[8+i*6], ASK_PRICE_i+1 = c[9+i*6]
// --------------------------------------------------------------------------

// decodeOverseasFuturesAsk 는 HDFFF010 frame 을 OverseasFuturesAskEvent 슬라이스로 디코딩한다.
func decodeOverseasFuturesAsk(f frame) ([]OverseasFuturesAskEvent, error) {
	chunks, err := chunkFieldsErr(f.Fields, f.Count, overseasFuturesAskFieldCount)
	if err != nil {
		return nil, err
	}
	out := make([]OverseasFuturesAskEvent, 0, len(chunks))
	for _, c := range chunks {
		out = append(out, parseOverseasFuturesAskChunk(c))
	}
	return out, nil
}

func parseOverseasFuturesAskChunk(c []string) OverseasFuturesAskEvent {
	ev := OverseasFuturesAskEvent{
		Symbol:    c[0],
		RecvDate:  c[1],
		RecvTime:  c[2],
		PrevPrice: asDecimal(c[3]),
		SttlPrice: asDecimal(c[34]),
		Raw:       c,
	}
	for i := 0; i < 5; i++ {
		base := 4 + i*6
		ev.BidQntt[i] = asInt64(c[base])      // BID_QNTT_(i+1)
		ev.BidNum[i] = c[base+1]              // BID_NUM_(i+1)
		ev.BidPrice[i] = asDecimal(c[base+2]) // BID_PRICE_(i+1)
		ev.AskQntt[i] = asInt64(c[base+3])    // ASK_QNTT_(i+1)
		ev.AskNum[i] = c[base+4]              // ASK_NUM_(i+1)
		ev.AskPrice[i] = asDecimal(c[base+5]) // ASK_PRICE_(i+1)
	}
	return ev
}
