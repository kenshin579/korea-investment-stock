# Issue #27: Rate Limiting 개선 프로젝트 문서

이 폴더는 Korea Investment Stock API의 Rate Limiting 개선 프로젝트(Issue #27)와 관련된 모든 문서를 포함합니다.

## 📋 주요 문서

### 요구사항 및 계획
- [`prd-27.md`](./prd-27.md) - Product Requirements Document
- [`prd-27-cache.md`](./prd-27-cache.md) - TTL 캐시 기능 요구사항
- [`todo-27.md`](./todo-27.md) - 작업 TODO 리스트 및 진행 상황

### 분석 문서
- [`rate_limiter_analysis.md`](./rate_limiter_analysis.md) - 기존 RateLimiter 분석
- [`thread_pool_executor_analysis.md`](./thread_pool_executor_analysis.md) - ThreadPoolExecutor 분석
- [`api_methods_analysis.md`](./api_methods_analysis.md) - API 메서드 분석
- [`error_pattern_analysis.md`](./error_pattern_analysis.md) - 에러 패턴 분석

### 구현 가이드
- [`rate_limit_implementation.md`](./rate_limit_implementation.md) - Rate Limiting 구현 상세
- [`rate_limiter_defense_mechanisms.md`](./rate_limiter_defense_mechanisms.md) - 방어 메커니즘
- [`thread_pool_executor_improvement.md`](./thread_pool_executor_improvement.md) - ThreadPool 개선안
- [`improved_threadpool_pattern.py`](./improved_threadpool_pattern.py) - 개선된 패턴 예제

### Phase별 완료 보고서

#### Phase 2: Enhanced RateLimiter
- [`threadpool_executor_phase2_4_completion.md`](./threadpool_executor_phase2_4_completion.md) - ThreadPoolExecutor 개선

#### Phase 3: 에러 핸들링
- [`phase3_1_error_detection_completion.md`](./phase3_1_error_detection_completion.md) - 에러 감지
- [`phase3_2_exponential_backoff_completion.md`](./phase3_2_exponential_backoff_completion.md) - Exponential Backoff
- [`phase3_3_error_recovery_completion.md`](./phase3_3_error_recovery_completion.md) - 에러 복구
- [`phase3_4_completion_report.md`](./phase3_4_completion_report.md) - Phase 3.4 완료
- [`phase3_completion_summary.md`](./phase3_completion_summary.md) - Phase 3 전체 요약

#### Phase 4: 배치 처리
- [`phase4_1_batch_params_completion.md`](./phase4_1_batch_params_completion.md) - 배치 파라미터화
- [`phase4_completion_report.md`](./phase4_completion_report.md) - Phase 4 완료 보고서

#### Phase 5: 모니터링 및 통계
- [`phase5_1_stats_save_completion.md`](./phase5_1_stats_save_completion.md) - 통계 저장 기능
- [`phase5_1_advanced_stats_completion.md`](./phase5_1_advanced_stats_completion.md) - 고급 통계 관리

#### Phase 6: 테스트
- [`phase6_1_unit_tests_completion.md`](./phase6_1_unit_tests_completion.md) - 단위 테스트
- [`phase6_completion_summary.md`](./phase6_completion_summary.md) - Phase 6 요약

#### Phase 7: 문서화
- [`phase7_1_documentation_completion.md`](./phase7_1_documentation_completion.md) - 문서화 완료

### 프로젝트 요약
- [`