// Package bonds 는 장내채권 (Korean bond) API 클라이언트.
//
// Phase 3.1 메서드 (8):
//
//   - SearchBondInfo             — 채권 기본조회 (CTPF1114R, 70 fields)
//   - InquireIssueInfo           — 발행정보 (CTPF1101R, 69 fields)
//   - InquirePrice               — 현재가 시세 (FHKBJ773400C0)
//   - InquireCcnl                — 현재가 체결 (FHKBJ773403C0)
//   - InquireAskingPrice         — 현재가 호가 (FHKBJ773401C0, 5단계)
//   - InquireDailyPrice          — 현재가 일별 (FHKBJ773404C0)
//   - InquireDailyItemchartprice — 기간별 시세 (FHKBJ773701C0)
//   - InquireAvgUnit             — 평균단가조회 (CTPF2005R)
//
// 사용자는 root kis.Client.Bonds 로 접근.
package bonds
