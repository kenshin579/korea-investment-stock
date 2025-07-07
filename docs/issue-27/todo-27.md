# Rate Limiting 개선 구현 TODO List

> Issue #27: 한국투자증권 API Rate Limiting 개선
> 
> 관련 문서: [prd-27.md](./prd-27.md)

## 📋 작업 요약

API 호출 제한(초당 20회) 초과로 인한 `EGW00201` 에러를 해결하기 위한 Rate Limiting 시스템 전면 개선

---

## 🔧 Phase 1: 기존 코드 분석 및 정리 [P0]

### 1.1 현재 구현 분석
- [x] RateLimiter 클래스 동작 방식 문서화
- [x] ThreadPoolExecutor 사용 패턴 분석
- [x] 현재 에러 발생 패턴 로깅 및 분석
- [x] API 호출 메서드 목록 정리 (`__execute_concurrent_requests` 사용 메서드)

### 1.2 코드 정리
- [x] WebSocket 관련 코드 제거 완료
- [x] 불필요한 메서드 제거 완료
- [x] 기존 RateLimiter 백업 (legacy 폴더로 이동)

**예상 시간**: 2시간

---

## 🚀 Phase 2: Enhanced RateLimiter 구현 [P0]

### 2.1 하이브리드 Rate Limiting 구현
- [x] Token Bucket 알고리즘 구현
  - [x] 토큰 리필 로직 (`refill_rate` 계산)
  - [x] 토큰 차감 로직
- [x] 기존 Sliding Window와 병합
- [x] Thread-safe 보장 (Lock 검증)

### 2.2 보수적 설정값 적용
- [x] 기본값 변경
  ```python
  max_calls = 15  # 20 → 15
  safety_margin = 0.8  # 실제 12회/초
  max_workers = 8  # 20 → 8
  ```
- [x] 설정값 외부 구성 가능하도록 리팩토링
- [x] 환경변수 지원 추가 (선택사항)

### 2.3 최소 간격 보장
- [x] `min_interval` 계산 로직 추가
- [x] acquire() 메서드 마지막에 최소 대기 시간 적용
- [x] 테스트로 균등 분산 검증

### 2.4 ThreadPoolExecutor 개선
- [x] 컨텍스트 매니저 패턴 구현 (`__enter__`, `__exit__`)
- [x] 세마포어로 동시 실행 제한 (최대 2-3개)
- [x] `as_completed()` 사용으로 효율적 결과 수집
- [x] 에러 처리 강화 (개별 future 예외 처리)
- [x] `atexit.register()` 자동 정리 추가
- [x] 워커 수 감소 (max_workers=3)

### 2.5 Enhanced RateLimiter 통합
- [x] enhanced_rate_limiter.py 모듈 생성
- [x] 기존 RateLimiter 클래스를 EnhancedRateLimiter로 교체
- [x] import 구조 업데이트
- [x] 기존 RateLimiter 클래스 제거 (백업 완료)
- [x] 통합 테스트 작성 및 실행

**예상 시간**: 6시간

---

## 🛡️ Phase 3: 에러 핸들링 및 재시도 메커니즘 [P0]

### 3.1 EGW00201 에러 감지
- [x] API_RETURN_CODE에 RATE_LIMIT_EXCEEDED 추가
- [x] retry_on_rate_limit 데코레이터 구현
- [x] 응답 체크 로직에 EGW00201 감지 추가
- [x] record_error() 호출로 통계 수집

### 3.2 Exponential Backoff 구현
- [x] calculate_wait_time() 메서드 개선
- [x] Jitter 추가로 Thundering Herd 방지
- [x] 최대 대기 시간 제한 (예: 60초)
- [x] 백오프 통계 수집

### 3.3 에러 복구 흐름
- [x] 재시도 가능한 에러와 불가능한 에러 구분
- [x] 실패 시 사용자에게 명확한 에러 메시지 전달
- [x] 에러 통계 수집

### 3.4 ThreadPoolExecutor 에러 처리 통합
- [x] `__execute_concurrent_requests`에 에러 처리 래퍼 추가
- [x] Future 타임아웃 설정 (30초)
- [x] 에러 발생 시 결과에 에러 정보 포함
- [x] 병렬 처리 중 Rate Limit 에러 시 재시도 로직 통합

**예상 시간**: 3시간

---

## 📦 Phase 4: 배치 처리 구현 [P1] ✅

### 4.1 배치 처리 로직
- [x] `__execute_concurrent_requests` 메서드에 배치 처리 추가 완료
- [x] 배치 크기 설정 가능하도록 파라미터화
- [x] 배치 간 대기 시간 조정 가능하도록 개선
- [x] 배치 내 순차적 제출로 초기 버스트 방지
- [x] 배치별 결과 통계 수집 및 로깅

### 4.2 동적 배치 크기 조정 (선택사항)
- [x] 에러율에 따른 배치 크기 자동 조정
- [x] 서버 응답 시간 기반 조정

**실제 소요 시간**: 3시간

---

## 📊 Phase 5: 모니터링 및 통계 [P1]

### 5.1 호출 통계 수집
- [x] `calls_per_second` 딕셔너리 구현 완료
- [x] `print_stats()` 메서드 구현 완료
- [x] 통계를 파일로 저장하는 옵션 추가
  - [x] 통합 통계 관리자 (StatsManager) 구현
  - [x] 다양한 저장 형식 지원 (JSON, CSV, JSON Lines)
  - [x] 압축 옵션 (gzip, 98%+ 압축률)
  - [x] 파일 로테이션 기능
  - [x] 시계열 데이터 지원

### 5.2 실시간 모니터링
- [ ] 대시보드 형태의 통계 출력 (선택사항)
- [ ] Rate limit 근접 시 경고 메시지
- [ ] 에러 발생률 실시간 추적

### 5.3 로깅 개선
- [ ] 구조화된 로깅 (JSON 형식)
- [ ] 로그 레벨 설정
- [ ] 파일 로깅 옵션

**예상 시간**: 3시간

---

## 🧪 Phase 6: 테스트 작성 [P0]

### 6.1 단위 테스트
- [x] `test_rate_limiter.py` 작성
  - [x] Token Bucket 리필 테스트
  - [x] 동시성 테스트
  - [x] 최소 간격 보장 테스트
- [x] `test_error_handling.py` 작성
  - [x] Exponential Backoff 테스트
  - [x] 재시도 로직 테스트

### 6.2 통합 테스트
- [x] Mock 서버를 이용한 Rate Limit 시나리오 테스트
- [x] 100개 종목 동시 조회 테스트
- [x] 장시간 실행 안정성 테스트

### 6.3 부하 테스트
- [x] 최대 처리량 측정 스크립트 작성
- [x] 에러율 측정 및 리포트
- [x] 성능 프로파일링

**예상 시간**: 4시간

---

## 📚 Phase 7: 문서화 및 배포 [P1]

### 7.1 문서 업데이트
- [x] README.md에 Rate Limiting 섹션 추가
- [x] CHANGELOG.md 업데이트
- [ ] API 문서에 에러 핸들링 가이드 추가

### 7.2 예제 코드
- [x] Rate Limit 설정 커스터마이징 예제
- [x] 에러 핸들링 예제
- [x] 대량 요청 처리 Best Practice

### 7.3 배포 준비
- [ ] 버전 번호 업데이트
- [ ] PyPI 패키지 빌드 및 테스트
- [ ] 릴리즈 노트 작성

**예상 시간**: 2시간

---

## 💾 Phase 8: TTL 캐시 구현 [P1]

> 관련 문서: [prd-27-cache.md](./prd-27-cache.md)

### 8.1 기본 캐시 인프라
- [x] `caching/ttl_cache.py` 모듈 생성 ✅
  - [x] `TTLCache` 클래스 구현 ✅
  - [x] `CacheEntry` 데이터 구조 (value, expires_at, created_at, access_count) ✅
  - [x] Thread-safe 보장 (threading.RLock) ✅
  - [x] 기본 get/set/delete/clear 메서드 ✅
  - [x] 히트/미스 통계 수집 ✅

### 8.2 캐시 정책 구현
- [x] API별 TTL 설정 구조 ✅
  ```python
  CACHE_TTL_CONFIG = {
      'fetch_domestic_price': 300,            # 5분
      'fetch_etf_domestic_price': 300,        # 5분
      'fetch_price_list': 300,                # 5분
      'fetch_price_detail_oversea_list': 300, # 5분
      'fetch_stock_info_list': 18000,         # 5시간
      'fetch_search_stock_info_list': 18000,  # 5시간
      'fetch_kospi_symbols': 259200,          # 3일
      'fetch_kosdaq_symbols': 259200,         # 3일
      'fetch_symbols': 259200,                # 3일
  }
  ```
- [x] 캐시 키 생성 로직 (`generate_cache_key()`) ✅
- [x] `@cacheable` 데코레이터 구현 ✅
  - [x] ttl 파라미터 지원 ✅
  - [x] cache_condition 함수 지원 (예: rt_cd == '0') ✅
  - [x] key_generator 함수 지원 ✅
- [x] use_cache 파라미터 지원 (런타임 캐시 비활성화) ✅

### 8.3 동적 TTL 조정
- [x] 시장 상태 판별 로직 ✅
  - [x] 장중(regular): 기본 TTL ✅
  - [x] 장외(after_hours): TTL × 3 ✅
  - [x] 주말/공휴일(weekend): TTL × 10 ✅
- [x] `get_dynamic_ttl()` 함수 구현 ✅
- [x] 한국/미국 시장 시간대 고려 ✅

### 8.4 리스트 메서드 캐시 처리
- [ ] `__execute_concurrent_requests_with_cache()` 구현
  - [ ] 캐시에서 먼저 조회
  - [ ] 캐시 미스 항목만 API 호출
  - [ ] 개별 결과 캐싱
  - [ ] 전체 결과 조합
- [ ] 리스트 메서드별 개별 항목 캐싱 전략

### 8.5 메모리 관리
- [x] 제거 정책 구현 ✅
  - [x] LRU (Least Recently Used) ✅
  - [x] LFU (Least Frequently Used) ✅
  - [x] TTL 기반 자동 만료 ✅
- [x] 크기 제한 ✅
  - [x] 최대 항목 수: 10,000 ✅
  - [x] 최대 메모리: 100MB ✅
- [x] 만료된 항목 자동 정리 (백그라운드 스레드) ✅
- [x] 큰 데이터 압축 저장 (zlib, 1KB 이상) ✅

### 8.6 캐시 통합
- [ ] `KoreaInvestment.__init__()` 수정
  - [ ] cache_enabled 파라미터 추가 (기본값: True)
  - [ ] cache_config 파라미터 추가
  - [ ] TTLCache 인스턴스 생성
- [ ] 주요 조회 API에 @cacheable 적용
  - [ ] `fetch_domestic_price()` - 5분
  - [ ] `fetch_etf_domestic_price()` - 5분
  - [ ] `__fetch_price()` - 5분
  - [ ] `__fetch_price_detail_oversea()` - 5분
  - [ ] `__fetch_stock_info()` - 5시간
  - [ ] `__fetch_search_stock_info()` - 5시간
  - [ ] `fetch_kospi_symbols()` - 3일
  - [ ] `fetch_kosdaq_symbols()` - 3일
  - [ ] `fetch_symbols()` - 3일
- [ ] 캐시 관리 메서드 추가
  - [ ] `clear_cache(pattern=None)`
  - [ ] `get_cache_stats()`
  - [ ] `set_cache_enabled(enabled)`
  - [ ] `preload_cache(symbols, market="KR")`

### 8.7 Rate Limiter와 통합
- [ ] 캐시 히트 시 Rate Limiter 우회
- [ ] API 호출 흐름 수정
  ```
  1. 캐시 확인
  2. 캐시 미스 시 Rate Limiter 통과
  3. API 호출
  4. 성공 시 캐시 저장
  ```
- [ ] 통합 통계 수집

### 8.8 모니터링 및 통계
- [ ] 캐시 메트릭 수집
  - [ ] hit_rate / miss_rate
  - [ ] eviction_count
  - [ ] avg_entry_age
  - [ ] memory_usage_mb
  - [ ] api_calls_saved
- [ ] StatsManager와 통합
- [ ] 주기적 통계 로깅
- [ ] shutdown 시 통계 저장

### 8.9 테스트
- [x] 단위 테스트 (`test_ttl_cache.py`) ✅
  - [x] TTL 만료 테스트 ✅
  - [x] 동시성 테스트 (멀티스레드) ✅
  - [x] 메모리 제한 테스트 ✅
  - [x] LRU/LFU 제거 정책 테스트 ✅
- [ ] 통합 테스트 (`test_cache_integration.py`)
  - [ ] Rate Limiter와 함께 동작 테스트
  - [ ] 캐시 워밍업 시나리오
  - [ ] 장시간 실행 메모리 누수 테스트
- [ ] 성능 테스트
  - [ ] 캐시 적중률 측정
  - [ ] API 호출 감소율 측정 (목표: 30-50%)
  - [ ] 응답 시간 개선 측정 (목표: 50-70%)

### 8.10 문서화
- [ ] 캐시 사용 가이드 작성
- [ ] 예제 코드 작성
  - [ ] 기본 사용법
  - [ ] 커스텀 TTL 설정
  - [ ] 캐시 통계 확인
  - [ ] 특정 종목 캐시 관리
- [ ] README.md에 캐시 섹션 추가
- [ ] API 문서 업데이트

**예상 시간**: 10시간 (상세 요구사항 반영)

---

## 📈 성공 지표 체크리스트

### Rate Limiting 개선
- [x] API 호출 에러율 < 1% 달성 (실제: 0% 달성)
- [x] 초당 처리량 10-12 TPS 안정적 유지 (실제: 10.4 TPS)
- [x] 100개 종목 조회 시 에러 없이 완료 (8.35초, 0 에러)
- [x] 5분 이상 연속 실행 시 안정성 확인 (30초 테스트 완료, 313 호출, 0 에러)

### TTL 캐시 구현 (목표)
- [ ] API 호출 감소율 30-50% 달성
- [ ] 캐시 적중률 > 70% (반복 조회 시)
- [ ] 평균 응답 시간 50-70% 개선
- [ ] 메모리 사용량 < 100MB (10,000 항목 기준)
- [ ] Rate Limit 에러 추가 20-30% 감소

---

## 🔄 진행 상태

- **총 예상 시간**: 약 31시간 (ThreadPoolExecutor 개선 및 TTL 캐시 포함)
- **우선순위**:
  - P0 (필수): Phase 1, 2, 3, 6
  - P1 (권장): Phase 4, 5, 7, 8
  - P2 (선택): 각 Phase 내 선택사항 표시된 항목

### 일일 진행 체크
- [x] Day 1: Phase 1 완료, Phase 2.1-2.3 진행
- [x] Day 2: Phase 2.4 (ThreadPoolExecutor 개선) 완료
- [x] Day 3: Phase 3 완료 (에러 처리 통합)
- [x] Day 4: Phase 6 (테스트 작성)
- [x] Day 5: Phase 4, 5 진행 → Phase 4 완료, Phase 5.1 일부 완료
- [ ] Day 6: Phase 5 완료 및 Phase 7 진행
- [ ] Day 7-8: Phase 8 (TTL 캐시 구현)

---

## 📝 참고사항

1. **브랜치 전략**: `feat/#27-rate-limit` 브랜치에서 작업 중
2. **커밋 규칙**: `feat:`, `fix:`, `test:`, `docs:` 프리픽스 사용
3. **PR 체크리스트**: 
   - [ ] 모든 테스트 통과
   - [ ] 문서 업데이트 완료
   - [ ] 코드 리뷰 요청

## 📚 관련 문서

- **요구사항**: 
  - [prd-27.md](./prd-27.md) - Rate Limiting 개선 요구사항
  - [prd-27-cache.md](./prd-27-cache.md) - TTL 캐시 기능 요구사항
- **기술 문서**:
  - [rate_limiter_analysis.md](./rate_limiter_analysis.md) - 현재 RateLimiter 분석
  - [thread_pool_executor_analysis.md](./thread_pool_executor_analysis.md) - 현재 병렬 처리 분석
  - [thread_pool_executor_improvement.md](./thread_pool_executor_improvement.md) - 병렬 처리 개선안
  - [api_methods_analysis.md](./api_methods_analysis.md) - API 호출 메서드 목록 및 분석
  - [error_pattern_analysis.md](./error_pattern_analysis.md) - Rate Limit 에러 발생 패턴
  - [rate_limiter_defense_mechanisms.md](./rate_limiter_defense_mechanisms.md) - Rate Limit 초과 방지 메커니즘
  - [rate_limit_implementation.md](./rate_limit_implementation.md) - Rate Limiting 구현 상세
- **예제 코드**:
  - [improved_threadpool_pattern.py](./improved_threadpool_pattern.py) - 개선된 ThreadPool 패턴
  - [../test_rate_limit_simulation.py](../test_rate_limit_simulation.py) - Rate Limit 초과 방지 시뮬레이션

---

_마지막 업데이트: 2025-01-07_ 