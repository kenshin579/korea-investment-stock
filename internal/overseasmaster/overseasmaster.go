package overseasmaster

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"strings"

	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

// MarketURLs 는 KIS 공개 다운로드의 11 거래소 마스터 파일 URL.
var MarketURLs = map[string]string{
	"nas": "https://new.real.download.dws.co.kr/common/master/nasmst.cod.zip",
	"nys": "https://new.real.download.dws.co.kr/common/master/nysmst.cod.zip",
	"ams": "https://new.real.download.dws.co.kr/common/master/amsmst.cod.zip",
	"shs": "https://new.real.download.dws.co.kr/common/master/shsmst.cod.zip",
	"shi": "https://new.real.download.dws.co.kr/common/master/shimst.cod.zip",
	"szs": "https://new.real.download.dws.co.kr/common/master/szsmst.cod.zip",
	"szi": "https://new.real.download.dws.co.kr/common/master/szimst.cod.zip",
	"tse": "https://new.real.download.dws.co.kr/common/master/tsemst.cod.zip",
	"hks": "https://new.real.download.dws.co.kr/common/master/hksmst.cod.zip",
	"hnx": "https://new.real.download.dws.co.kr/common/master/hnxmst.cod.zip",
	"hsx": "https://new.real.download.dws.co.kr/common/master/hsxmst.cod.zip",
}

// TSV 컬럼 인덱스 (0-based). 형식: cp949 인코딩 + TSV 24 컬럼.
// KIS 공개 다운로드 NASMST.COD / NYSMST.COD 등 실측 결과 기반.
const (
	colCountryCode  = 0  // 국가코드 (US/CN/JP/HK/VN)
	colMarketNum    = 1  // 거래소 숫자코드
	colMarketCode   = 2  // 거래소 코드 (NAS/NYS/AMS/SHS/TSE/HKS 등)
	colMarketKorean = 3  // 거래소 한글명 (나스닥/뉴욕/상해 등)
	colSymbol       = 4  // 종목코드 (AAPL/A/1301 등)
	colFullCode     = 5  // 거래소+종목코드 (NASAAPL/NYSA 등)
	colKoreanName   = 6  // 한글 종목명
	colEnglishName  = 7  // 영문 종목명
	colStockType    = 8  // 종목구분 (2=보통주, 3=ETF/ETC 등)
	colCurrency     = 9  // 거래통화 (USD/JPY/CNY/HKD/VND)
	colDecimals     = 10 // 소수점 자리수
	// col 11: 빈 컬럼
	colBasePrice   = 12 // 기준가/종가
	colTradeUnit   = 13 // 매매수량단위
	colMinUnit     = 14 // 최소거래단위
	colOpenTime    = 15 // 시장개시시각 (930/900 등)
	colCloseTime   = 16 // 시장마감시각 (1600/1530 등)
	colSuspended   = 17 // 거래정지 여부 (N/Y)
	// col 18: 빈 컬럼
	colISINCode = 19 // ISIN 유사코드 (000 = 없음)
)

// columnNames 는 Raw map 의 key 로 사용하는 컬럼명 (인덱스 순).
var columnNames = []string{
	"국가코드", "거래소번호", "거래소코드", "거래소한글명",
	"종목코드", "전체코드", "한글종목명", "영문종목명",
	"종목구분", "통화", "소수점자리수", "예비",
	"기준가", "매매수량단위", "최소거래단위", "시장개시시각", "시장마감시각", "거래정지",
	"예비2", "ISIN코드", "플래그1", "플래그2", "기타1", "기타2",
}

// Symbol 은 해외 거래소 마스터 한 행 (typed 핵심 필드 + Raw map fallback).
//
// 형식: cp949+TSV 24 컬럼. 모든 컬럼은 Raw 에 한글 컬럼명 → 값 으로 저장.
type Symbol struct {
	Symbol      string // 종목코드 (예 "AAPL", "1301")
	EnglishName string // 영문 종목명 (예 "APPLE INC")
	KoreanName  string // 한글 종목명 (예 "애플")
	Currency    string // 거래 통화 (예 "USD", "JPY", "CNY")
	MarketCode  string // 거래소 코드 (예 "NAS", "NYS", "TSE")
	StockType   string // 종목구분 (2=보통주, 3=ETF 등)
	BasePrice   string // 기준가/종가 문자열 (소수점 포함 가능)
	SuspendedYn string // 거래정지 여부 (N=정상, Y=정지)

	Raw map[string]string // 모든 24 컬럼 (한글 컬럼명 → 값)
}

// Parse 는 market ("nas"/"nys"/"ams" 등) 별 마스터 ZIP byte 를 파싱해 Symbol 슬라이스 반환.
//
// 형식: cp949 인코딩 + TSV (탭 구분) 24 컬럼.
// 헤더 행 없음 — 첫 행부터 바로 데이터.
func Parse(market string, zipBytes []byte) ([]Symbol, error) {
	if _, ok := MarketURLs[market]; !ok {
		return nil, fmt.Errorf("overseasmaster: unknown market %q", market)
	}
	mst, err := openCodFromZip(zipBytes)
	if err != nil {
		return nil, fmt.Errorf("overseasmaster: %s: %w", market, err)
	}

	decoded, err := decodeCP949(mst)
	if err != nil {
		return nil, fmt.Errorf("overseasmaster: %s: decode: %w", market, err)
	}

	var out []Symbol
	for _, line := range strings.Split(decoded, "\n") {
		line = strings.TrimRight(line, "\r")
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) < len(columnNames) {
			// 24 컬럼 미만이면 파싱 불가 — skip
			continue
		}

		raw := make(map[string]string, len(columnNames))
		for i, name := range columnNames {
			if i < len(parts) {
				raw[name] = strings.TrimSpace(parts[i])
			}
		}

		sym := strings.TrimSpace(parts[colSymbol])
		if sym == "" {
			continue
		}

		out = append(out, Symbol{
			Symbol:      sym,
			EnglishName: strings.TrimSpace(parts[colEnglishName]),
			KoreanName:  strings.TrimSpace(parts[colKoreanName]),
			Currency:    strings.TrimSpace(parts[colCurrency]),
			MarketCode:  strings.TrimSpace(parts[colMarketCode]),
			StockType:   strings.TrimSpace(parts[colStockType]),
			BasePrice:   strings.TrimSpace(parts[colBasePrice]),
			SuspendedYn: strings.TrimSpace(parts[colSuspended]),
			Raw:         raw,
		})
	}
	return out, nil
}

// openCodFromZip 은 ZIP byte 에서 .cod 또는 .COD 파일 byte 를 추출.
// KIS 의 ZIP 내부 파일명은 대문자 (NASMST.COD) 임에 유의.
func openCodFromZip(zipBytes []byte) ([]byte, error) {
	if len(zipBytes) == 0 {
		return nil, fmt.Errorf("zip open: empty input")
	}
	zr, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return nil, fmt.Errorf("zip open: %w", err)
	}
	for _, f := range zr.File {
		lower := strings.ToLower(f.Name)
		if strings.HasSuffix(lower, ".cod") || strings.HasSuffix(lower, ".mst") {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("zip read %s: %w", f.Name, err)
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
	}
	return nil, fmt.Errorf(".cod or .mst file not found in zip")
}

// decodeCP949 는 cp949 byte 를 UTF-8 string 으로 변환.
// golang.org/x/text/encoding/korean.EUCKR 는 cp949 호환 (Microsoft 확장 포함).
func decodeCP949(b []byte) (string, error) {
	decoded, _, err := transform.Bytes(korean.EUCKR.NewDecoder(), b)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
