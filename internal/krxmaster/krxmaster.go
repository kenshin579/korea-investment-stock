package krxmaster

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

// 한국투자증권 KRX 마스터 파일 다운로드 URL (var 로 노출 — 테스트에서 override 가능).
var (
	KospiURL  = "https://new.real.download.dws.co.kr/common/master/kospi_code.mst.zip"
	KosdaqURL = "https://new.real.download.dws.co.kr/common/master/kosdaq_code.mst.zip"
)

// KospiSymbol 은 KRX KOSPI 종목 마스터 (kospi_code.mst) 한 행.
//
// 핵심 필드 typed + Raw map[한글컬럼명]값 으로 미typed 70 컬럼 fallback.
// docs: 한국투자 GitHub open-trading-api/stocks_info 의 Python 정제 코드 참조.
type KospiSymbol struct {
	ShortCode       string          // 단축코드 (예 "005930")
	StandardCode    string          // 표준코드 (ISIN, 예 "KR7005930003")
	KoreanName      string          // 한글명
	GroupCode       string          // 그룹코드 (ST=주권, EF=ETF, RT=REITs, IF=인프라, ...)
	MarketCapSize   string          // 시가총액 규모
	KOSPI200        string          // KOSPI200 섹터업종 편입 (Y/N)
	KOSPI100        string          // KOSPI100 편입
	KOSPI50         string          // KOSPI50 편입
	BasePrice       int64           // 기준가
	FaceValue       decimal.Decimal // 액면가
	ListedShares    int64           // 상장주수
	Capital         int64           // 자본금
	SettlementMonth string          // 결산월 (예 "12")
	PreferredStock  string          // 우선주 여부 (Y/N)
	SuspendedYn     string          // 거래정지 여부 (Y/N)

	Raw map[string]string // 모든 70 컬럼 (한글 키)
}

// kospiFieldSpecs 와 kospiColumns 는 한투 GitHub Python parsers/master_parser.py 의
// field_specs / part2_columns 와 1:1 매핑. 라인 마지막 227 byte 의 fwf 영역.
var kospiFieldSpecs = []int{
	2, 1, 4, 4, 4,
	1, 1, 1, 1, 1,
	1, 1, 1, 1, 1,
	1, 1, 1, 1, 1,
	1, 1, 1, 1, 1,
	1, 1, 1, 1, 1,
	1, 9, 5, 5, 1,
	1, 1, 2, 1, 1,
	1, 2, 2, 2, 3,
	1, 3, 12, 12, 8,
	15, 21, 2, 7, 1,
	1, 1, 1, 1, 9,
	9, 9, 5, 9, 8,
	9, 3, 1, 1, 1,
}

var kospiColumns = []string{
	"그룹코드", "시가총액규모", "지수업종대분류", "지수업종중분류", "지수업종소분류",
	"제조업", "저유동성", "지배구조지수종목", "KOSPI200섹터업종", "KOSPI100",
	"KOSPI50", "KRX", "ETP", "ELW발행", "KRX100",
	"KRX자동차", "KRX반도체", "KRX바이오", "KRX은행", "SPAC",
	"KRX에너지화학", "KRX철강", "단기과열", "KRX미디어통신", "KRX건설",
	"Non1", "KRX증권", "KRX선박", "KRX섹터_보험", "KRX섹터_운송",
	"SRI", "기준가", "매매수량단위", "시간외수량단위", "거래정지",
	"정리매매", "관리종목", "시장경고", "경고예고", "불성실공시",
	"우회상장", "락구분", "액면변경", "증자구분", "증거금비율",
	"신용가능", "신용기간", "전일거래량", "액면가", "상장일자",
	"상장주수", "자본금", "결산월", "공모가", "우선주",
	"공매도과열", "이상급등", "KRX300", "KOSPI", "매출액",
	"영업이익", "경상이익", "당기순이익", "ROE", "기준년월",
	"시가총액", "그룹사코드", "회사신용한도초과", "담보대출가능", "대주가능",
}

// ParseKospi 는 KOSPI 마스터 ZIP 의 byte 를 받아 종목 슬라이스로 디코딩.
func ParseKospi(zipBytes []byte) ([]KospiSymbol, error) {
	const fwfLen = 227
	mst, err := openMstFromZip(zipBytes)
	if err != nil {
		return nil, fmt.Errorf("krxmaster: kospi: %w", err)
	}
	decoded, err := decodeCP949(mst)
	if err != nil {
		return nil, fmt.Errorf("krxmaster: kospi: cp949: %w", err)
	}

	var out []KospiSymbol
	for _, line := range strings.Split(decoded, "\n") {
		line = strings.TrimRight(line, "\r")
		if len(line) < fwfLen+21 {
			continue
		}
		prefix := line[:len(line)-fwfLen]
		fwf := line[len(line)-fwfLen:]

		shortCode := strings.TrimSpace(prefix[0:9])
		standardCode := strings.TrimSpace(prefix[9:21])
		koreanName := strings.TrimSpace(prefix[21:])
		raw := parseFwf(fwf, kospiFieldSpecs, kospiColumns)

		out = append(out, KospiSymbol{
			ShortCode:       shortCode,
			StandardCode:    standardCode,
			KoreanName:      koreanName,
			GroupCode:       raw["그룹코드"],
			MarketCapSize:   raw["시가총액규모"],
			KOSPI200:        raw["KOSPI200섹터업종"],
			KOSPI100:        raw["KOSPI100"],
			KOSPI50:         raw["KOSPI50"],
			BasePrice:       atoi64(raw["기준가"]),
			FaceValue:       toDecimal(raw["액면가"]),
			ListedShares:    atoi64(raw["상장주수"]),
			Capital:         atoi64(raw["자본금"]),
			SettlementMonth: raw["결산월"],
			PreferredStock:  raw["우선주"],
			SuspendedYn:     raw["거래정지"],
			Raw:             raw,
		})
	}
	return out, nil
}

// openMstFromZip 은 ZIP byte 에서 .mst 파일 byte 를 추출.
func openMstFromZip(zipBytes []byte) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return nil, fmt.Errorf("zip open: %w", err)
	}
	for _, f := range zr.File {
		if strings.HasSuffix(f.Name, ".mst") {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("zip read %s: %w", f.Name, err)
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
	}
	return nil, fmt.Errorf(".mst file not found in zip")
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

// parseFwf 는 fixed-width line 을 widths 별로 잘라 한글 컬럼명 → 값 map 반환.
func parseFwf(line string, widths []int, names []string) map[string]string {
	out := make(map[string]string, len(names))
	pos := 0
	for i, w := range widths {
		if pos+w > len(line) {
			break
		}
		out[names[i]] = strings.TrimSpace(line[pos : pos+w])
		pos += w
	}
	return out
}

// atoi64 는 빈 문자열/공백 → 0 fallback. 한투 마스터 데이터에 빈 값 흔함.
func atoi64(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return v
}

// toDecimal 은 빈 문자열 → decimal.Zero fallback.
func toDecimal(s string) decimal.Decimal {
	s = strings.TrimSpace(s)
	if s == "" {
		return decimal.Zero
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero
	}
	return d
}
