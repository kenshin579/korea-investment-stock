# fetch_price_detail_oversea 리팩토링 TODO

## 1단계: constants.py 수정

- [x] EXCD 상수 키 변경 (`"NYSE"` → `"NYS"`, `"NASDAQ"` → `"NAS"` 등)
- [x] EXCD_BY_COUNTRY 상수 추가
- [x] __init__.py에 EXCD_BY_COUNTRY export 추가

## 2단계: korea_investment_stock.py 수정

- [x] EXCD_BY_COUNTRY import 추가
- [x] fetch_price_detail_oversea 인자명 변경 (`market` → `country_code`)
- [x] 기본값 변경 (`"KR"` → `"US"`)
- [x] KR/KRX 체크 로직 제거
- [x] EXCD_BY_COUNTRY를 활용한 거래소 순회 로직 적용
- [x] ValueError 메시지 한글화
- [x] docstring 업데이트 (Query Parameters, Returns, Raises 포함)
- [x] return type hint 추가 (`-> dict`)

## 3단계: Wrapper 클래스 수정

- [x] cached_korea_investment.py: `market` → `country_code` 변경
- [x] cached_korea_investment.py: 기본값 `"US"` 적용
- [x] rate_limited_korea_investment.py: `market` → `country_code` 변경
- [x] rate_limited_korea_investment.py: 기본값 `"US"` 적용

## 4단계: 테스트 수정

- [ ] 기존 테스트에서 `market` → `country_code` 변경
- [ ] `fetch_price_detail_oversea("AAPL")` 기본값 테스트 추가
- [ ] `fetch_price_detail_oversea("AAPL", country_code="US")` 테스트
- [ ] 홍콩 주식 테스트 (`country_code="HK"`)
- [ ] 일본 주식 테스트 (`country_code="JP"`)
- [ ] 지원하지 않는 country_code에 대한 ValueError 테스트
- [ ] CachedKoreaInvestment 래퍼 테스트
- [ ] RateLimitedKoreaInvestment 래퍼 테스트

## 5단계: 문서 업데이트

- [ ] CLAUDE.md API Response Fields Reference 섹션 업데이트
- [ ] CHANGELOG.md에 Breaking Changes 추가
- [ ] examples/ 예제 코드 업데이트 (필요시)

## 6단계: 검증

- [ ] `pytest` 전체 테스트 통과 확인
- [ ] 기존 테스트 통과 확인
- [ ] 신규 테스트 통과 확인
