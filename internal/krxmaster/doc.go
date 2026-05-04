// Package krxmaster 는 KRX KOSPI/KOSDAQ 종목 마스터 파일의 디코딩/파싱 로직.
//
// 한투 API 가 아니라 KRX 가 공개 다운로드로 제공하는 .mst.zip 파일을 처리.
// cp949 인코딩 + fixed-width 컬럼 포맷이라 별도 파서 필요.
//
// 사용자에게 노출되지 않는 internal 패키지. domestic 패키지의 FetchKospiSymbols
// 와 FetchKosdaqSymbols 가 호출.
package krxmaster
