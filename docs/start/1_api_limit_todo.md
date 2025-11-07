# API 호출 속도 제한 구현 TODO

## 1단계: 핵심 구현

### 1.1 디렉토리 구조 생성
- [ ] `korea_investment_stock/rate_limit/` 디렉토리 생성
- [ ] `korea_investment_stock/rate_limit/__init__.py` 생성

### 1.2 RateLimiter 구현
- [ ] `rate_limiter.py` 파일 생성
- [ ] `RateLimiter` 클래스 구현
  - [ ] `__init__()`: 초기화 및 검증
  - [ ] `wait()`: 속도 제한 대기 로직
  - [ ] `get_stats()`: 통계 조회
  - [ ] `adjust_rate_limit()`: 동적 속도 조정
  - [ ] `threading.Lock` 스레드 안전성 구현

### 1.3 RateLimitedKoreaInvestment 구현
- [ ] `rate_limited_korea_investment.py` 파일 생성
- [ ] `RateLimitedKoreaInvestment` 클래스 구현
  - [ ] `__init__()`: 브로커 래핑 초기화
  - [ ] Context Manager 지원 (`__enter__`, `__exit__`)
  - [ ] 18개 API 메서드 래핑
    - [ ] `fetch_price()`
    - [ ] `fetch_domestic_price()`
    - [ ] `fetch_etf_domestic_price()`
    - [ ] `fetch_price_detail_oversea()`
    - [ ] `fetch_stock_info()`
    - [ ] `fetch_search_stock_info()`
    - [ ] `fetch_kospi_symbols()`
    - [ ] `fetch_kosdaq_symbols()`
    - [ ] `fetch_ipo_schedule()`
    - [ ] 9개 IPO 헬퍼 메서드
  - [ ] `get_rate_limit_stats()` 구현
  - [ ] `adjust_rate_limit()` 구현

### 1.4 패키지 통합
- [ ] `rate_limit/__init__.py`에 exports 추가
- [ ] `korea_investment_stock/__init__.py`에 exports 추가

## 2단계: 테스트 구현

### 2.1 RateLimiter 단위 테스트
- [ ] `test_rate_limiter.py` 파일 생성
- [ ] 테스트 케이스 작성
  - [ ] `test_rate_limiter_basic()`: 기본 속도 제한
  - [ ] `test_rate_limiter_thread_safe()`: 스레드 안전성
  - [ ] `test_rate_limiter_stats()`: 통계 조회
  - [ ] `test_rate_limiter_adjust()`: 동적 속도 조정
  - [ ] `test_rate_limiter_invalid_input()`: 입력 검증

### 2.2 통합 테스트
- [ ] `test_rate_limited_integration.py` 파일 생성
- [ ] 테스트 케이스 작성
  - [ ] `test_rate_limited_basic()`: 기본 API 호출
  - [ ] `test_rate_limited_context_manager()`: Context Manager
  - [ ] `test_rate_limited_preserves_functionality()`: 기능 보존
  - [ ] `test_rate_limited_stats()`: 통계 조회

### 2.3 Stress Test 업데이트
- [ ] `examples/stress_test.py` 수정
  - [ ] `RateLimitedKoreaInvestment` import 추가
  - [ ] `rate_limited_broker` 생성 코드 추가
  - [ ] 500회 호출 테스트

## 3단계: 검증 및 문서화

### 3.1 단위 테스트 실행
- [ ] `pytest korea_investment_stock/rate_limit/test_rate_limiter.py -v`
- [ ] 모든 테스트 통과 확인

### 3.2 통합 테스트 실행
- [ ] 환경 변수 설정 확인
- [ ] `pytest korea_investment_stock/rate_limit/test_rate_limited_integration.py -v`
- [ ] 모든 테스트 통과 확인

### 3.3 Stress Test 실행
- [ ] `python examples/stress_test.py` 실행
- [ ] 성공 기준 확인
  - [ ] 500회 API 호출 완료
  - [ ] 성공률 100%
  - [ ] 실행 시간 33-40초
  - [ ] API 속도 제한 에러 0건

### 3.4 CLAUDE.md 업데이트
- [ ] Rate Limiting 섹션 추가
- [ ] 사용 예제 추가
- [ ] Cache와 결합 사용 예제 추가

### 3.5 CHANGELOG.md 업데이트
- [ ] v0.8.0 섹션 생성
- [ ] 새 기능 설명
- [ ] Breaking Changes (없음) 명시

## 성공 기준 체크리스트

### 필수 (P0)
- [ ] `examples/stress_test.py` 500회 호출 100% 성공
- [ ] API 호출 속도 제한 에러 0건
- [ ] 스레드 안전 구현 검증
- [ ] `KoreaInvestment` 클래스 변경 없음 (래퍼만 추가)

### 권장 (P1)
- [ ] `CLAUDE.md` 문서화 완료
- [ ] 사용 예제 3개 이상 작성
- [ ] 단위 테스트 90% 이상 커버리지
- [ ] 통합 테스트 모두 통과

## 예상 소요 시간

- **1단계 (핵심 구현)**: 3-4시간
- **2단계 (테스트 구현)**: 2-3시간
- **3단계 (검증 및 문서화)**: 1-2시간
- **총 예상 시간**: 6-9시간

## 주의사항

1. **스레드 안전성**: `threading.Lock` 사용 필수
2. **성능**: 호출당 오버헤드 5ms 미만 유지
3. **호환성**: 기존 코드 변경 없이 opt-in 방식
4. **테스트**: 실제 API 호출 필요 (환경 변수 설정)
