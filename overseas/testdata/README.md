# overseas testdata

각 한투 API 메서드의 단위 테스트 fixture.

## REST API 응답 (합성 JSON)

- `price_detail_success.json` — 해외주식_현재가상세 (HHDFS76200200) 정상 응답
- `search_info_success.json` — 해외주식_상품기본정보 (CTPF1702R) 정상 응답
- `daily_price_success.json` — 해외주식_기간별시세 (HHDFS76240000) 정상 응답
- `daily_chart_price_success.json` — 해외주식 종목/지수/환율 기간별시세 (FHKST03030100) 정상 응답
- `updown_rate_success.json` — 해외주식_상승율_하락율 (HHDFS76290000) 정상 응답

각 JSON의 필드는 `docs/api/해외주식/<API>.md`의 응답 필드 정의에 1:1 매핑. 값은 합성 (실제 시세 아님).

## 해외 마스터 sample

- `<market>_code_sample.cod.zip` — 출처 + 재생성 방법은 `internal/overseasmaster/testdata/README.md` 참조
