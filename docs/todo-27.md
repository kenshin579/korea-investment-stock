# Rate Limiting 개선 구현 TODO List

> Issue #27: 한국투자증권 API Rate Limiting 개선
> 
> 관련 문서: [PRD-27.md](./PRD-27.md)

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
- [ ] 컨텍스트 매니저 패턴 구현 (`__enter__`, `__exit__`)
- [ ] 세마포어로 동시 실행 제한 (최대 2-3개)
- [ ] `as_completed()` 사용으로 효율적 결과 수집
- [ ] 에러 처리 강화 (개별 future 예외 처리)
- [ ] `atexit.register()` 자동 정리 추가
- [ ] 워커 수 감소 (max_workers=3)

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
- [ ] API 응답에서 `msg_cd` 확인 로직 추가
- [ ] 에러 발생 시 별도 예외 발생
- [ ] 각 API 호출 메서드에 에러 체크 추가

### 3.2 Exponential Backoff 구현
- [x] `__handle_rate_limit_error` 메서드 추가 완료
- [ ] 재시도 로직 구현
  - [ ] 재시도 횟수 관리 (최대 5회)
  - [ ] 백오프 시간 계산 (1, 2, 4, 8, 16, 32초)
  - [ ] Jitter 추가 (0-10%)
- [ ] API 호출 래퍼 함수 생성

### 3.3 에러 복구 흐름
- [ ] 재시도 가능한 에러와 불가능한 에러 구분
- [ ] 실패 시 사용자에게 명확한 에러 메시지 전달
- [ ] 에러 통계 수집

### 3.4 ThreadPoolExecutor 에러 처리 통합
- [ ] `__execute_concurrent_requests`에 에러 처리 래퍼 추가
- [ ] Future 타임아웃 설정 (30초)
- [ ] 에러 발생 시 결과에 에러 정보 포함
- [ ] 병렬 처리 중 Rate Limit 에러 시 재시도 로직 통합

**예상 시간**: 3시간

---

## 📦 Phase 4: 배치 처리 구현 [P1]

### 4.1 배치 처리 로직
- [x] `__execute_concurrent_requests` 메서드에 배치 처리 추가 완료
- [ ] 배치 크기 설정 가능하도록 파라미터화
- [ ] 배치 간 대기 시간 조정 가능하도록 개선
- [ ] 배치 내 순차적 제출로 초기 버스트 방지
- [ ] 배치별 결과 통계 수집 및 로깅

### 4.2 동적 배치 크기 조정 (선택사항)
- [ ] 에러율에 따른 배치 크기 자동 조정
- [ ] 서버 응답 시간 기반 조정

**예상 시간**: 2시간

---

## 📊 Phase 5: 모니터링 및 통계 [P1]

### 5.1 호출 통계 수집
- [x] `calls_per_second` 딕셔너리 구현 완료
- [x] `print_stats()` 메서드 구현 완료
- [ ] 통계를 파일로 저장하는 옵션 추가

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
- [ ] `test_rate_limiter.py` 작성
  - [ ] Token Bucket 리필 테스트
  - [ ] 동시성 테스트
  - [ ] 최소 간격 보장 테스트
- [ ] `test_error_handling.py` 작성
  - [ ] Exponential Backoff 테스트
  - [ ] 재시도 로직 테스트

### 6.2 통합 테스트
- [ ] Mock 서버를 이용한 Rate Limit 시나리오 테스트
- [ ] 100개 종목 동시 조회 테스트
- [ ] 장시간 실행 안정성 테스트

### 6.3 부하 테스트
- [ ] 최대 처리량 측정 스크립트 작성
- [ ] 에러율 측정 및 리포트
- [ ] 성능 프로파일링

**예상 시간**: 4시간

---

## 📚 Phase 7: 문서화 및 배포 [P1]

### 7.1 문서 업데이트
- [ ] README.md에 Rate Limiting 섹션 추가
- [ ] CHANGELOG.md 업데이트
- [ ] API 문서에 에러 핸들링 가이드 추가

### 7.2 예제 코드
- [ ] Rate Limit 설정 커스터마이징 예제
- [ ] 에러 핸들링 예제
- [ ] 대량 요청 처리 Best Practice

### 7.3 배포 준비
- [ ] 버전 번호 업데이트
- [ ] PyPI 패키지 빌드 및 테스트
- [ ] 릴리즈 노트 작성

**예상 시간**: 2시간

---

## 📈 성공 지표 체크리스트

- [ ] API 호출 에러율 < 1% 달성
- [ ] 초당 처리량 10-12 TPS 안정적 유지
- [ ] 100개 종목 조회 시 에러 없이 완료
- [ ] 5분 이상 연속 실행 시 안정성 확인

---

## 🔄 진행 상태

- **총 예상 시간**: 약 29시간 (ThreadPoolExecutor 개선 포함)
- **우선순위**:
  - P0 (필수): Phase 1, 2, 3, 6
  - P1 (권장): Phase 4, 5, 7
  - P2 (선택): 각 Phase 내 선택사항 표시된 항목

### 일일 진행 체크
- [ ] Day 1: Phase 1 완료, Phase 2.1-2.3 진행
- [ ] Day 2: Phase 2.4 (ThreadPoolExecutor 개선) 완료
- [ ] Day 3: Phase 3 완료 (에러 처리 통합)
- [ ] Day 4: Phase 6 (테스트 작성)
- [ ] Day 5: Phase 4, 5 진행
- [ ] Day 6: Phase 7 및 최종 테스트

---

## 📝 참고사항

1. **브랜치 전략**: `feat/#27-rate-limit` 브랜치에서 작업 중
2. **커밋 규칙**: `feat:`, `fix:`, `test:`, `docs:` 프리픽스 사용
3. **PR 체크리스트**: 
   - [ ] 모든 테스트 통과
   - [ ] 문서 업데이트 완료
   - [ ] 코드 리뷰 요청

## 📚 관련 문서

- **요구사항**: [PRD-27.md](./PRD-27.md)
- **기술 문서**:
  - [RateLimiter_Analysis.md](./RateLimiter_Analysis.md) - 현재 RateLimiter 분석
  - [ThreadPoolExecutor_Analysis.md](./ThreadPoolExecutor_Analysis.md) - 현재 병렬 처리 분석
  - [ThreadPoolExecutor_Improvement.md](./ThreadPoolExecutor_Improvement.md) - 병렬 처리 개선안
  - [API_Methods_Analysis.md](./API_Methods_Analysis.md) - API 호출 메서드 목록 및 분석
  - [Error_Pattern_Analysis.md](./Error_Pattern_Analysis.md) - Rate Limit 에러 발생 패턴
  - [RateLimiter_Defense_Mechanisms.md](./RateLimiter_Defense_Mechanisms.md) - Rate Limit 초과 방지 메커니즘
  - [rate_limit_implementation.md](./rate_limit_implementation.md) - Rate Limiting 구현 상세
- **예제 코드**:
  - [improved_threadpool_pattern.py](./improved_threadpool_pattern.py) - 개선된 ThreadPool 패턴
  - [../test_rate_limit_simulation.py](../test_rate_limit_simulation.py) - Rate Limit 초과 방지 시뮬레이션

---

_마지막 업데이트: 2024-12-28_ 