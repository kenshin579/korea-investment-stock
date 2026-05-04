# krxmaster testdata

`kospi_code_sample.mst.zip` 와 `kosdaq_code_sample.mst.zip` 는 KRX 종목 마스터 파일에서 6자리 단축코드 일반 주권 3 행만 추출한 sample (펀드/ETF 9자리 코드는 제외).

## 출처

- KOSPI 마스터: https://new.real.download.dws.co.kr/common/master/kospi_code.mst.zip
- KOSDAQ 마스터: https://new.real.download.dws.co.kr/common/master/kosdaq_code.mst.zip

한국투자증권이 공개 다운로드로 제공. `internal/krxmaster` 의 cp949+fwf 파서가 실제 KRX byte 와 호환되는지 검증하기 위한 단위 테스트 sample.

## 재생성 방법

```bash
cd /tmp && rm -rf krxmaster && mkdir krxmaster && cd krxmaster
curl -sSL -o kospi_code.mst.zip https://new.real.download.dws.co.kr/common/master/kospi_code.mst.zip
unzip -o kospi_code.mst.zip
LC_ALL=C grep -E '^[0-9]{6}' kospi_code.mst | head -n 3 > kospi_code_sample.mst
zip kospi_code_sample.mst.zip kospi_code_sample.mst
# KOSDAQ — 외국 기업 9자리 코드(예: 900110)가 앞에 있어 KR 표준코드 필터 필요
curl -sSL -o kosdaq_code.mst.zip https://new.real.download.dws.co.kr/common/master/kosdaq_code.mst.zip
unzip -o kosdaq_code.mst.zip
iconv -f cp949 -t utf-8 kosdaq_code.mst | grep '^[0-9]\{6\}   KR' | head -3 | iconv -f utf-8 -t cp949 > kosdaq_code_sample.mst
zip kosdaq_code_sample.mst.zip kosdaq_code_sample.mst
```

## 라이선스

KRX 종목 마스터는 한국투자증권이 무료 공개 다운로드로 배포. 본 sample 은 학습/테스트 용도이며, 라이브러리의 단위 테스트 외 사용 권장하지 않음.
