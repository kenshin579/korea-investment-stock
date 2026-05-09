// Package websocket 은 KIS 실시간 (WebSocket) 도메인.
//
// 한투 docs: docs/api/기타/실시간_(웹소켓)_접속키_발급.md (인증)
// docs/api/국내주식/국내주식_실시간*.md (KRX/NXT/통합 EP)
//
// 디자인:
//   - Phase 8: docs/superpowers/specs/2026-05-09-phase8-websocket-design.md
//   - Phase 9: docs/superpowers/specs/2026-05-09-phase9-websocket-nxt-unified-design.md
//   - Phase 10: docs/superpowers/specs/2026-05-09-phase10-websocket-overseas-design.md
//
// Phase 8 — 인프라 + 국내주식 KRX 시세 5 endpoint:
//
//	H0STCNT0  실시간체결가 (KRX)            SubscribeKrxTrade / OnKrxTrade
//	H0STASP0  실시간호가 (KRX)              SubscribeKrxAsk / OnKrxAsk
//	H0STANC0  실시간예상체결 (KRX)          SubscribeKrxExpectTrade / OnKrxExpectTrade
//	H0STOUP0  시간외 실시간체결가 (KRX)     SubscribeKrxOvernightTrade / OnKrxOvernightTrade
//	H0STOAC0  시간외 실시간예상체결 (KRX)   SubscribeKrxOvernightExpect / OnKrxOvernightExpect
//
// Phase 9 — NXT/통합 변형 10 endpoint (NXT 와 통합은 schema 동일, 5 base struct + 10 type alias):
//
//	H0NXCNT0  실시간체결가 (NXT)            SubscribeNxtTrade / OnNxtTrade
//	H0UNCNT0  실시간체결가 (통합)           SubscribeUnifiedTrade / OnUnifiedTrade
//	H0NXASP0  실시간호가 (NXT, +KMID/NMID)   SubscribeNxtAsk / OnNxtAsk
//	H0UNASP0  실시간호가 (통합, +KMID/NMID)  SubscribeUnifiedAsk / OnUnifiedAsk
//	H0NXANC0  실시간예상체결 (NXT, +VI_STND_PRC) SubscribeNxtExpectTrade / OnNxtExpectTrade
//	H0UNANC0  실시간예상체결 (통합, +VI_STND_PRC) SubscribeUnifiedExpectTrade / OnUnifiedExpectTrade
//	H0NXPGM0  실시간프로그램매매 (NXT)      SubscribeNxtProgramTrade / OnNxtProgramTrade
//	H0UNPGM0  실시간프로그램매매 (통합)     SubscribeUnifiedProgramTrade / OnUnifiedProgramTrade
//	H0NXMBC0  실시간회원사 (NXT)            SubscribeNxtMember / OnNxtMember
//	H0UNMBC0  실시간회원사 (통합)           SubscribeUnifiedMember / OnUnifiedMember
//
// Phase 10 — 해외주식 실시간 시세 2 endpoint:
//
//	HDFSCNT0  해외주식 실시간지연체결가      SubscribeOverseasTrade / OnOverseasTrade
//	HDFSASP0  해외주식 실시간호가 (1호가)    SubscribeOverseasAsk / OnOverseasAsk
//
// 사용자는 root kis.Client 의 WS 필드로 접근.
package websocket

import _ "github.com/coder/websocket"
