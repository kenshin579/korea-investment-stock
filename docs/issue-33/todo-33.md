# TODO-33: 미국 주식 현재가 조회 통합 인터페이스 구현

> PRD-33 기반 구현 작업 목록  
> 예상 소요 시간: 5-7 영업일  
> 우선순위: P0 (Critical)

## 📋 작업 현황
- [ ] Phase 1: 핵심 구현 (2-3일)
- [ ] Phase 2: 안정화 및 최적화 (3-4일)  
- [ ] Phase 3: 문서화 및 배포 (1-2일)

---

## Phase 1: 핵심 구현 (2-3일)

### 1.1 `__fetch_price()` 메서드 수정
- [ ] `korea_investment_stock.py`에서 `__fetch_price()` 메서드 찾기
- [ ] US market 처리 로직 수정
  - [ ] 기존 `self.fetch_oversea_price(symbol)` 호출 제거
  - [ ] `self.__fetch_price_detail_oversea(symbol, market)` 호출로 변경
  - [ ] 주석으로 변경 이유 명시 (기존 메서드 미구현, 상세 시세 API 활용)
- [ ] 에러 처리 확인 및 개선

### 1.2 국내 가격 조회 메서드 캡슐화
- [ ] `fetch_etf_domestic_price()` → `__fetch_etf_domestic_price()` 변경
  - [ ] 메서드명 변경
  - [ ] 관련 호출 부분 모두 수정
  - [ ] docstring 업데이트 (internal method 명시)
- [ ] `fetch_domestic_price()` → `__fetch_domestic_price()` 변경
  - [ ] 메서드명 변경
  - [ ] 관련 호출 부분 모두 수정
  - [ ] docstring 업데이트 (internal method 명시)
- [ ] 변경 후 기존 코드 호환성 확인

### 1.3 기본 동작 테스트
- [ ] 국내 주식 조회 테스트 (기존 동작 확인)
  - [ ] 삼성전자(005930) 조회
  - [ ] ETF 조회 테스트
- [ ] 미국 주식 조회 테스트 (새 기능)
  - [ ] AAPL 조회
  - [ ] TSLA 조회
  - [ ] NVDA 조회
- [ ] 혼합 조회 테스트
  - [ ] 국내 + 미국 주식 동시 조회

---

## Phase 2: 안정화 및 최적화 (3-4일)

### 2.1 응답 형식 분석 및 통일
- [ ] 국내 주식 응답 형식 분석
  - [ ] 필수 필드 목록 작성
  - [ ] 데이터 타입 확인
- [ ] 미국 주식 응답 형식 분석 (`__fetch_price_detail_oversea()` 응답)
  - [ ] 필수 필드 목록 작성
  - [ ] 추가 필드 확인 (PER, PBR, EPS, BPS 등)
- [ ] 응답 변환 로직 필요성 검토
  - [ ] 필드명 매핑 테이블 작성
  - [ ] `_normalize_response()` 메서드 구현 여부 결정

### 2.2 통합 테스트 케이스 작성
- [ ] `tests/test_integration_us_stocks.py` 파일 생성
- [ ] 테스트 시나리오 구현:
  - [ ] `test_unified_price_interface()` - 통합 인터페이스 테스트
  - [ ] `test_fetch_price_internal_routing()` - 내부 라우팅 테스트
  - [ ] `test_us_stock_response_format()` - 응답 형식 검증
  - [ ] `test_mixed_market_batch()` - 국내/미국 혼합 배치 조회
  - [ ] `test_invalid_market_type()` - 잘못된 market 타입 처리
- [ ] 기존 테스트와의 통합성 확인

### 2.3 에러 처리 및 복구
- [ ] Rate Limit 처리 확인
  - [ ] 미국 주식 API의 Rate Limit 확인
  - [ ] 기존 Rate Limiter와의 통합
- [ ] 네트워크 에러 처리
  - [ ] 재시도 로직 확인
  - [ ] Exponential backoff 적용 여부
- [ ] 심볼 에러 처리
  - [ ] 잘못된 심볼 입력 시 에러 메시지
  - [ ] 거래소별 심볼 차이 처리

### 2.4 성능 최적화
- [ ] 캐시 통합 확인
  - [ ] TTL 캐시 적용 여부
  - [ ] 미국 주식 캐시 키 형식 확인
- [ ] 배치 처리 성능 테스트
  - [ ] 100개 종목 동시 조회 테스트
  - [ ] 병렬 처리 최적화
- [ ] 모니터링 통합
  - [ ] StatsManager 통합 확인
  - [ ] 미국 주식 조회 통계 수집

---

## Phase 3: 문서화 및 배포 (1-2일)

### 3.1 API 문서 업데이트
- [ ] `README.md` 업데이트
  - [ ] 통합 인터페이스 사용법 추가
  - [ ] 미국 주식 조회 예제 추가
  - [ ] 제약사항 명시 (모의투자 미지원 등)
- [ ] docstring 업데이트
  - [ ] `fetch_price_list()` 설명에 미국 주식 지원 추가
  - [ ] `__fetch_price()` 내부 로직 설명
- [ ] API 변경사항 문서화
  - [ ] Private 메서드 변경 내역
  - [ ] Migration guide 작성 (필요시)

### 3.2 예제 코드 작성
- [ ] `examples/us_stock_price_example.py` 생성
  - [ ] 기본 사용법
  - [ ] 혼합 조회 예제
  - [ ] 에러 처리 예제
- [ ] 기존 예제 업데이트
  - [ ] 통합 인터페이스 사용 권장
  - [ ] Deprecated 메서드 표시

### 3.3 릴리즈 준비
- [ ] `CHANGELOG.md` 업데이트
  - [ ] 새로운 기능 설명
  - [ ] Breaking changes 명시 (private 메서드 변경)
  - [ ] Migration 가이드 링크
- [ ] 버전 번호 결정
  - [ ] Minor version bump 검토 (기능 추가)
  - [ ] Patch version 검토 (호환성 유지 시)
- [ ] 릴리즈 노트 작성
  - [ ] 주요 변경사항 요약
  - [ ] 알려진 이슈 (모의투자 미지원)
  - [ ] 감사 인사

---

## 🔍 검증 체크리스트

### 기능 검증
- [ ] 국내 주식 조회 정상 동작
- [ ] 미국 주식 조회 정상 동작
- [ ] 혼합 조회 정상 동작
- [ ] ETF 조회 정상 동작
- [ ] 에러 상황 처리 적절

### 성능 검증
- [ ] API 응답 시간 < 500ms
- [ ] 캐시 적중률 > 80%
- [ ] 배치 처리 100개 < 10초
- [ ] Rate Limit 준수

### 품질 검증
- [ ] 테스트 커버리지 > 90%
- [ ] 문서화 완성도 100%
- [ ] 코드 리뷰 완료
- [ ] 보안 취약점 없음

---

## 📌 참고사항

### API 정보
- 해외주식 현재가상세: `/uapi/overseas-price/v1/quotations/price-detail`
- TR ID: `HHDFS76200200`
- 지원 거래소: NAS(나스닥), NYS(뉴욕), AMS(아멕스)
- 제약사항: 모의투자 미지원, 미국은 실시간 무료시세

### 리스크
- 모의투자 계정에서는 미국 주식 조회 불가
- API 응답 형식이 국내/해외 간 상이할 수 있음
- 거래소별 심볼 형식 차이 존재

### 연락처
- 프로젝트 관리자: [담당자명]
- 기술 문의: [이메일/슬랙]
- 긴급 연락처: [전화번호]

---

**작성일**: 2025-01-13  
**최종 수정**: 2025-01-13  
**상태**: 작업 대기 