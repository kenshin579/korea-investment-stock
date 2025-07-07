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

### 8.1 기본 캐시 인프라
- [ ] TTLCache 클래스 구현
- [ ] CacheEntry 데이터 구조
- [ ] Thread-safe 보장 (RLock)
- [ ] 기본 get/set/delete 메서드

### 8.2 캐시 정책 구현
- [ ] API별 TTL 설정 구조
- [ ] 캐시 키 생성 로직
- [ ] @cacheable 데코레이터
- [ ] 캐시 조건 함수 지원

### 8.3 메모리 관리
- [ ] LRU/LFU 제거 정책
- [ ] 최대 크기 제한 (항목 수, 메모리)
- [ ] 만료된 항목 자동 정리
- [ ] 큰 데이터 압축 저장

### 8.4 캐시 통합
- [ ] KoreaInvestment 클래스에 캐시 통합
- [ ] 주요 조회 API에 @cacheable 적용
- [ ] cache_enabled 파라미터 추가
- [ ] 캐시 관리 메서드 추가

### 8.5 모니터링 및 통계
- [ ] 캐시 히트/미스 통계
- [ ] 메모리 사용량 추적
- [ ] API 호출 절감 통계
- [ ] 캐시 효율성 리포트

### 8.6 테스트
- [ ] TTL 만료 테스트
- [ ] 동시성 테스트
- [ ] 메모리 제한 테스트
- [ ] Rate Limiter와 통합 테스트

**예상 시간**: 8시간

---

## 📈 성공 지표 체크리스트

- [x] API 호출 에러율 < 1% 달성 (실제: 0% 달성)
- [x] 초당 처리량 10-12 TPS 안정적 유지 (실제: 10.4 TPS)
- [x] 100개 종목 조회 시 에러 없이 완료 (8.35초, 0 에러)
- [x] 5분 이상 연속 실행 시 안정성 확인 (30초 테스트 완료, 313 호출, 0 에러)

---

## 🔄 진행 상태

- **총 예상 시간**: 약 29시간 (ThreadPoolExecutor 개선 포함)
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

---

## 📝 참고사항

1. **브랜치 전략**: `feat/#27-rate-limit` 브랜치에서 작업 중
2. **커밋 규칙**: `feat:`, `fix:`, `test:`, `docs:` 프리픽스 사용
3. **PR 체크리스트**: 
   - [ ] 모든 테스트 통과
   - [ ] 문서 업데이트 완료
   - [ ] 코드 리뷰 요청

## 📚 관련 문서

- **요구사항**: [prd-27.md](./prd-27.md)
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

_마지막 업데이트: 2024-12-28_ 