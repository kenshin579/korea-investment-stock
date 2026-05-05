# overseasmaster testdata

`nas_code_sample.cod.zip` 는 KIS 공개 다운로드의 NASDAQ 마스터 파일 첫 3 행 sample.

## 출처

- NASDAQ 마스터: https://new.real.download.dws.co.kr/common/master/nasmst.cod.zip
- 다른 거래소: https://new.real.download.dws.co.kr/common/master/{nys,ams,shs,shi,szs,szi,tse,hks,hnx,hsx}mst.cod.zip

KIS 가 공개 다운로드로 제공. `internal/overseasmaster` 의 파서가 실제 KIS byte 와 호환되는지 검증하기 위한 단위 테스트 sample.

## 형식

- 인코딩: cp949 (EUC-KR Microsoft 확장)
- 구분자: TSV (탭, `\t`)
- 컬럼 수: 24 개
- 헤더 행: 없음 (첫 행부터 데이터)

| 인덱스 | 컬럼명 | 예시 |
|--------|--------|------|
| 0 | 국가코드 | US |
| 1 | 거래소번호 | 22 |
| 2 | 거래소코드 | NAS |
| 3 | 거래소한글명 | 나스닥 |
| 4 | **종목코드** | AACB |
| 5 | 전체코드 | NASAACB |
| 6 | **한글종목명** | 아티우스 애퀴지션 2 |
| 7 | **영문종목명** | ARTIUS II ACQUISITION INC |
| 8 | 종목구분 | 2 (보통주) / 3 (ETF) |
| 9 | **통화** | USD |
| 10 | 소수점자리수 | 4 |
| 11 | (빈 컬럼) | |
| 12 | **기준가** | 10.3800 |
| 13 | 매매수량단위 | 1 |
| 14 | 최소거래단위 | 1 |
| 15 | 시장개시시각 | 930 |
| 16 | 시장마감시각 | 1600 |
| 17 | 거래정지 | N |
| 18 | (빈 컬럼) | |
| 19 | ISIN코드 | 000 |
| 20 | 플래그1 | 0 |
| 21 | 플래그2 | 0 |
| 22 | 기타1 | |
| 23 | 기타2 | |

## 재생성 방법

```bash
cd /tmp && rm -rf overseas_master && mkdir overseas_master && cd overseas_master
curl -sSL -o nasmst.cod.zip https://new.real.download.dws.co.kr/common/master/nasmst.cod.zip
unzip -o nasmst.cod.zip
LC_ALL=C head -n 3 NASMST.COD > nas_code_sample.cod
zip nas_code_sample.cod.zip nas_code_sample.cod
```

## 라이선스

KIS 가 무료 공개 다운로드로 배포. 본 sample 은 학습/테스트 용도.
