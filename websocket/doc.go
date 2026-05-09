// Package websocket 은 KIS 실시간 (WebSocket) 도메인.
//
// 한투 docs: docs/api/기타/실시간_(웹소켓)_접속키_발급.md (인증)
// docs/api/국내주식/국내주식_실시간*.md (KRX 5 EP)
//
// 디자인: docs/superpowers/specs/2026-05-09-phase8-websocket-design.md
//
// Phase 8 — 인프라 + 국내주식 KRX 시세 5 endpoint:
//
//	H0STCNT0  실시간체결가 (KRX)        SubscribeKrxTrade / OnKrxTrade
//	H0STASP0  실시간호가 (KRX)          SubscribeKrxAsk / OnKrxAsk
//	H0STANC0  실시간예상체결 (KRX)      SubscribeKrxExpectTrade / OnKrxExpectTrade
//	H0STOAC0  시간외 실시간체결가 (KRX) SubscribeKrxOvernightTrade / OnKrxOvernightTrade
//	H0STOAA0  시간외 실시간예상체결 (KRX) SubscribeKrxOvernightExpect / OnKrxOvernightExpect
//
// 사용자는 root kis.Client 의 WS 필드로 접근.
package websocket

import _ "github.com/coder/websocket"
