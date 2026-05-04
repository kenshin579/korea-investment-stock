# domestic testdata

각 한투 API 메서드의 단위 테스트 fixture.

## REST API 응답 (합성 JSON)

- `price_success.json` — 주식현재가_시세 (FHKST01010100) 정상 응답
- `product_info_success.json` — 상품기본조회 (CTPF1604R) 정상 응답
- `stock_info_success.json` — 주식기본조회 (CTPF1002R) 정상 응답
- `daily_chart_success.json` — 국내주식기간별시세 (FHKST03010100) 정상 응답
- `minute_chart_success.json` — 주식당일분봉조회 (FHKST03010200) 정상 응답

각 JSON 의 필드는 `docs/api/국내주식/<API>.md` 의 응답 필드 정의에 1:1 매핑. 값은 합성 (실제 시세 아님).

## KRX 마스터 sample

- `kospi_code_sample.mst.zip`, `kosdaq_code_sample.mst.zip` — 출처 + 재생성 방법은 `internal/krxmaster/testdata/README.md` 참조
