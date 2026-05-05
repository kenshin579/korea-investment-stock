// Package overseasmaster 는 KIS 공개 다운로드의 해외 거래소 마스터 파일
// (NASDAQ/NYSE/AMEX/홍콩/일본 등 11 거래소) 의 디코딩/파싱 로직.
//
// 한투 API 가 아니라 KIS 가 공개 다운로드로 제공하는 .cod.zip 파일을 처리.
// 형식: cp949 인코딩 + TSV (탭 구분) 24 컬럼.
//
// 사용자에게 노출되지 않는 internal 패키지. overseas 패키지의 FetchOverseasSymbols 가 호출.
package overseasmaster
