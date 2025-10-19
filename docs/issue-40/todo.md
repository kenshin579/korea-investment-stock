# TODO: Korea Investment Stock 단순화 체크리스트

> 이 문서는 [PRD](prd.md)와 [Implementation Guide](implementation.md)의 구현 체크리스트입니다.

**진행 상태 범례**:
- [ ] 미완료
- [x] 완료
- [~] 진행중

---

## Phase 1: 모듈 삭제 (우선순위: HIGH) ✅ 완료

### 1.1 rate_limiting/ 디렉토리 전체 삭제 ✅

- [x] `enhanced_rate_limiter.py` (~400 lines)
- [x] `enhanced_backoff_strategy.py` (~300 lines)
- [x] `enhanced_retry_decorator.py` (~200 lines)
- [x] `__init__.py` (~50 lines)

```bash
rm -rf korea_investment_stock/rate_limiting/
```

### 1.2 caching/ 디렉토리 전체 삭제 ✅

- [x] `ttl_cache.py` (~500 lines)
- [x] `market_hours.py` (~100 lines)
- [x] `__init__.py` (~50 lines)

```bash
rm -rf korea_investment_stock/caching/
```

### 1.3 visualization/ 디렉토리 전체 삭제 ✅

- [x] `plotly_visualizer.py` (~400 lines)
- [x] `dashboard.py` (~350 lines)
- [x] `charts.py` (~250 lines)
- [x] `__init__.py` (~50 lines)

```bash
rm -rf korea_investment_stock/visualization/
```

### 1.4 batch_processing/ 디렉토리 전체 삭제 ✅

- [x] `dynamic_batch_controller.py` (~300 lines)
- [x] `__init__.py` (~30 lines)

```bash
rm -rf korea_investment_stock/batch_processing/
```

### 1.5 monitoring/ 디렉토리 전체 삭제 ✅

- [x] `stats_manager.py` (~600 lines)
- [x] `__init__.py` (~30 lines)

```bash
rm -rf korea_investment_stock/monitoring/
```

### 1.6 error_handling/ 디렉토리 전체 삭제 ✅

- [x] `error_recovery_system.py` (~500 lines)
- [x] `__init__.py` (~30 lines)

```bash
rm -rf korea_investment_stock/error_handling/
```

### 1.7 legacy/ 디렉토리 전체 삭제 (선택사항) ✅

- [x] `rate_limiter_v1.py`

```bash
rm -rf korea_investment_stock/legacy/
```

**예상 결과**: ~4,090 lines 삭제

---

## Phase 2: 메인 모듈 수정 (우선순위: HIGH) ✅ 완료

**파일**: `korea_investment_stock/korea_investment_stock.py`

### 2.1 Import 문 제거 ✅

- [x] Rate limiting imports (4줄)
  ```python
  from .rate_limiting.enhanced_rate_limiter import EnhancedRateLimiter
  from .rate_limiting.enhanced_backoff_strategy import get_backoff_strategy
  from .rate_limiting.enhanced_retry_decorator import retry_on_rate_limit, retry_on_network_error
  ```

- [x] Error handling imports (1줄)
  ```python
  from .error_handling.error_recovery_system import get_error_recovery_system
  ```

- [x] Monitoring imports (1줄)
  ```python
  from .monitoring.stats_manager import get_stats_manager
  ```

- [x] Caching imports (1줄)
  ```python
  from .caching import TTLCache, cacheable, CACHE_TTL_CONFIG
  ```

- [x] Visualization imports (3줄)
  ```python
  try:
      from .visualization import PlotlyVisualizer, DashboardManager
      VISUALIZATION_AVAILABLE = True
  except ImportError:
      VISUALIZATION_AVAILABLE = False
  ```

### 2.2 __init__() 메서드 간소화 ✅

- [x] Rate limiter 초기화 제거
- [x] Cache 초기화 제거
- [x] ThreadPoolExecutor 초기화 제거
- [x] Semaphore 초기화 제거
- [x] Visualizer 초기화 제거
- [x] Stats manager 초기화 제거
- [x] Error recovery 초기화 제거
- [x] atexit.register() 제거
- [x] max_workers, cache_enabled 파라미터 제거
- [x] Docstring 업데이트

### 2.3 List 기반 메서드 제거 (7개) ✅

- [x] `fetch_price_list()` 삭제
- [x] `fetch_price_list_with_batch()` 삭제
- [x] `fetch_price_list_with_dynamic_batch()` 삭제
- [x] `fetch_stock_info_list()` 삭제
- [x] `fetch_search_stock_info_list()` 삭제 - 첫 번째 정의
- [x] `fetch_search_stock_info_list()` 삭제 - 두 번째 정의
- [x] `fetch_price_detail_oversea_list()` 삭제

### 2.4 내부 실행 메서드 제거 (2개) ✅

- [x] `__execute_concurrent_requests()` 삭제 (~292 lines)
- [x] `__execute_concurrent_requests_with_cache()` 삭제 (~101 lines)

### 2.5 Private 메서드 → Public 전환 (7개) ✅

#### __fetch_price() → fetch_price() ✅

- [x] 메서드명 변경: `__fetch_price` → `fetch_price`
- [x] Docstring 업데이트
- [x] `__get_symbol_type` 호출을 `get_symbol_type`으로 변경
- [x] `__fetch_etf_domestic_price` 호출을 `fetch_etf_domestic_price`로 변경
- [x] `__fetch_domestic_price` 호출을 `fetch_domestic_price`로 변경
- [x] `__fetch_price_detail_oversea` 호출을 `fetch_price_detail_oversea`로 변경

#### __get_symbol_type() → get_symbol_type() ✅

- [x] 메서드명 변경: `__get_symbol_type` → `get_symbol_type`

#### __fetch_etf_domestic_price() → fetch_etf_domestic_price() ✅

- [x] 메서드명 변경: `__fetch_etf_domestic_price` → `fetch_etf_domestic_price`
- [x] Docstring 업데이트
- [x] `@cacheable` 데코레이터 제거
- [x] `@retry_on_rate_limit` 데코레이터 제거

#### __fetch_domestic_price() → fetch_domestic_price() ✅

- [x] 메서드명 변경: `__fetch_domestic_price` → `fetch_domestic_price`
- [x] Docstring 업데이트
- [x] `@cacheable` 데코레이터 제거
- [x] `@retry_on_rate_limit` 데코레이터 제거

#### __fetch_price_detail_oversea() → fetch_price_detail_oversea() ✅

- [x] 메서드명 변경: `__fetch_price_detail_oversea` → `fetch_price_detail_oversea`
- [x] `@cacheable` 데코레이터 제거
- [x] `@retry_on_rate_limit` 데코레이터 제거

#### __fetch_stock_info() → fetch_stock_info() ✅

- [x] 메서드명 변경: `__fetch_stock_info` → `fetch_stock_info`
- [x] `@cacheable` 데코레이터 제거
- [x] `@retry_on_rate_limit` 데코레이터 제거

#### __fetch_search_stock_info() → fetch_search_stock_info() ✅

- [x] 메서드명 변경: `__fetch_search_stock_info` → `fetch_search_stock_info`
- [x] `@cacheable` 데코레이터 제거
- [x] `@retry_on_rate_limit` 데코레이터 제거

#### __handle_rate_limit_error() 유지 (DEPRECATED 마킹됨)

- [x] `__handle_rate_limit_error()` 메서드는 하위 호환성을 위해 유지됨

### 2.6 Cache 관련 메서드 제거 (4개) ✅

- [x] `clear_cache()` 삭제
- [x] `get_cache_stats()` 삭제
- [x] `set_cache_enabled()` 삭제
- [x] `preload_cache()` 삭제

### 2.7 Monitoring 관련 메서드 제거 (6개) ✅

- [x] `create_monitoring_dashboard()` 삭제
- [x] `save_monitoring_dashboard()` 삭제
- [x] `create_stats_report()` 삭제
- [x] `get_system_health_chart()` 삭제
- [x] `get_api_usage_chart()` 삭제
- [x] `show_monitoring_dashboard()` 삭제

### 2.8 나머지 메서드 데코레이터 제거 ✅

- [x] `issue_access_token()` - `@retry_on_rate_limit` 제거
- [x] `fetch_kospi_symbols()` - `@cacheable` 제거
- [x] `fetch_kosdaq_symbols()` - `@cacheable` 제거
- [x] `fetch_ipo_schedule()` - `@cacheable` 제거
- [x] `fetch_ipo_schedule()` - `@retry_on_rate_limit` 제거

### 2.9 shutdown() 메서드 간소화 ✅

- [x] ThreadPoolExecutor shutdown 코드 제거
- [x] Event 처리 제거
- [x] Stats 저장 코드 제거
- [x] Rate limiter, backoff, cache, error recovery cleanup 제거
- [x] 간소화된 버전으로 교체 (단순 pass 문)

**예상 결과**: 1,941 lines → ~800 lines

---

## Phase 3: Package 설정 수정 (우선순위: HIGH) ✅ 완료

**파일**: `korea_investment_stock/__init__.py`

### 3.1 Import 문 정리 ✅

- [x] Rate limiting imports 제거 (없음 - 이미 정리됨)
- [x] Error handling imports 제거 (없음 - 이미 정리됨)
- [x] Batch processing imports 제거 (없음 - 이미 정리됨)
- [x] Monitoring imports 제거 (없음 - 이미 정리됨)

### 3.2 __all__ 리스트 업데이트 ✅

- [x] 제거된 모듈 exports 삭제 (이미 완료)
- [x] 핵심 4개만 유지: `KoreaInvestment`, `MARKET_CODE_MAP`, `EXCHANGE_CODE_MAP`, `API_RETURN_CODE`
- [x] Docstring 업데이트 (Pure wrapper 설명 추가)

**실제 결과**: 15 lines → 18 lines (docstring 개선으로 약간 증가)

---

## Phase 4: 테스트 수정 (우선순위: MEDIUM) ✅ 완료

### 4.1 테스트 파일 삭제 (12개) ✅

- [x] `test_rate_limiter.py` 삭제 (Phase 2에서 완료)
- [x] `test_enhanced_backoff.py` 삭제 (Phase 2에서 완료)
- [x] `test_rate_limit_error_detection.py` 삭제 (Phase 2에서 완료)
- [x] `test_rate_limit_simulation.py` 삭제 (Phase 2에서 완료)
- [x] `test_ttl_cache.py` 삭제 (Phase 2에서 완료)
- [x] `test_cache_integration.py` 삭제 (Phase 2에서 완료)
- [x] `test_batch_processing.py` 삭제 (Phase 2에서 완료)
- [x] `test_error_recovery.py` 삭제 (Phase 2에서 완료)
- [x] `test_error_handling.py` 삭제 (Phase 2에서 완료)
- [x] `test_stats_save.py` 삭제 (Phase 2에서 완료)
- [x] `test_enhanced_integration.py` 삭제 (Phase 2에서 완료)
- [x] `test_threadpool_improvement.py` 삭제 (Phase 2에서 완료)

### 4.2 남은 테스트 파일 확인 ✅

- [x] `test_korea_investment_stock.py` - 이미 Public API 사용 (수정 불필요)
- [x] `test_integration_us_stocks.py` - 이미 fetch_price() loop 사용 (수정 불필요)
- [x] `test_ipo_schedule.py` - 이미 단순화된 API 사용 (수정 불필요)
- [x] `test_ipo_integration.py` - 이미 단순화된 API 사용 (수정 불필요)

**실제 결과**: Phase 2에서 테스트 파일 삭제 및 정리가 이미 완료됨. 남은 4개 테스트 파일은 모두 단순화된 Public API를 사용하고 있어 추가 수정 불필요.

---

## Phase 5: Example 파일 수정 (우선순위: MEDIUM) ✅ 완료

### 5.1 Example 파일 삭제 (4개) ✅

- [x] `examples/rate_limiting_example.py` 삭제 (Phase 2에서 완료)
- [x] `examples/stats_management_example.py` 삭제 (Phase 2에서 완료)
- [x] `examples/stats_visualization_plotly.py` 삭제 (Phase 2에서 완료)
- [x] `examples/visualization_integrated_example.py` 삭제 (Phase 2에서 완료)

### 5.2 남은 Example 파일 확인 ✅

- [x] `basic_example.py` (208 lines) - Phase 2에서 생성, fetch_price() loop 패턴 (수정 불필요)
- [x] `ipo_schedule_example.py` - 이미 fetch_ipo_schedule() 단일 메서드 사용 (수정 불필요)
- [x] `us_stock_price_example.py` - 이미 fetch_price() loop 패턴 사용 (수정 불필요)

**실제 결과**: Phase 2에서 example 파일 삭제 및 생성이 이미 완료됨. 모든 example 파일이 단순화된 API를 사용하고 있어 추가 수정 불필요.

---

## Phase 6: 문서 업데이트 (우선순위: HIGH) ✅ 완료

### 6.1 README.md 업데이트
#### Features 섹션 ✅
#### Features 섹션
- [x] Rate limiting 항목 제거
- [x] Cache 항목 제거
- [x] Visualization 항목 제거
- [x] Batch processing 항목 제거
- [x] 단순화된 기능 목록으로 교체

#### Usage 섹션
#### Usage 섹션 ✅- [x] 복잡한 사용 예시 제거
- [x] 단순한 사용 예시로 교체
- [x] 사용자 제어 패턴 강조

#### Migration Guide ✅
- [x] Migration 섹션 추가
- [x] PRD 링크 추가

### 6.2 CLAUDE.md 업데이트
#### Architecture 섹션 ✅
#### Architecture 섹션
- [x] 제거된 모듈 설명 삭제
  - [x] rate_limiting/ 섹션 삭제
  - [x] caching/ 섹션 삭제
  - [x] visualization/ 섹션 삭제
  - [x] batch_processing/ 섹션 삭제
  - [x] monitoring/ 섹션 삭제
  - [x] error_handling/ 섹션 삭제

- [x] 단순화된 아키텍처 다이어그램 추가
- [x] Singleton Patterns 섹션 삭제
- [x] Threading & Concurrency 섹션 삭제

#### API Methods 섹션
- [x] Private → Public 전환 문서화
- [x] 제거된 메서드 목록 삭제
- [x] 새로운 Public API 목록 작성

#### Performance Characteristics
- [x] Benchmark 섹션 업데이트 (변경 반영)

### 6.3 CHANGELOG.md 업데이트

- [x] 0.6.0 버전 추가
- [x] Breaking Changes 명시
  - [x] Removed Modules 목록
  - [x] Removed Methods 목록
  - [x] Changed Methods 목록
- [x] Migration guide 링크 추가

### 6.4 MIGRATION.md 생성 (신규) - 선택사항

- [x] 파일 생성: `docs/MIGRATION.md`
- [x] Breaking changes 상세 ���명
- [x] Before/After 코드 예시 (5가지)
  - [x] 단일 조회
  - [x] 배치 조회
  - [x] 캐싱
  - [x] 모니터링
  - [x] IPO 조회
- [x] 권장 migration 전략 (Phase 1-4)

---

## Phase 7: 버전 관리 (우선순위: HIGH)

### 7.1 pyproject.toml 업데이트

- [ ] version: `0.5.0` → `0.6.0`
- [ ] dependencies 검토 (plotly 관련)

### 7.2 Git 작업

- [ ] Feature branch 생성
  ```bash
  git checkout -b feat/issue-40-simplify
  ```

- [ ] 단계별 커밋 (각 Phase별로)
  - [ ] Phase 1 커밋: `[feat] #40 - Remove rate_limiting, caching, visualization modules`
  - [ ] Phase 2 커밋: `[feat] #40 - Simplify main module and convert private to public methods`
  - [ ] Phase 3 커밋: `[feat] #40 - Update package exports`
  - [ ] Phase 4 커밋: `[feat] #40 - Update tests for simplified API`
  - [ ] Phase 5 커밋: `[feat] #40 - Update examples`
  - [ ] Phase 6 커밋: `[feat] #40 - Update documentation`
  - [ ] Phase 7 커밋: `[feat] #40 - Bump version to 0.6.0`

- [ ] PR 생성
  ```bash
  git push origin feat/issue-40-simplify
  gh pr create --title "[feat] #40 - Simplify library to pure API wrapper" \
    --body "$(cat docs/issue-40/prd.md)"
  ```

---

## Phase 8: 검증 & 배포 (우선순위: HIGH)

### 8.1 로컬 테스트

- [ ] 전체 테스트 실행
  ```bash
  pytest korea_investment_stock/tests/ -v
  ```

- [ ] 커버리지 확인
  ```bash
  pytest --cov=korea_investment_stock --cov-report=html
  ```

- [ ] Examples 실행 검증
  ```bash
  python examples/ipo_schedule_example.py
  python examples/us_stock_price_example.py
  python examples/basic_usage_example.py
  ```

- [ ] Integration 테스트 (실제 API 필요)
  ```bash
  pytest korea_investment_stock/tests/test_integration.py -v
  ```

### 8.2 코드 리뷰

- [ ] API surface 검증
  - [ ] Public 메서드 18개 확인
  - [ ] Private 메서드 없음 확인
  - [ ] 제거된 메서드 호출 없음 확인

- [ ] Breaking changes 확인
  - [ ] fetch_price_list() 제거 확인
  - [ ] Cache 관련 메서드 제거 확인
  - [ ] Monitoring 관련 메서드 제거 확인

- [ ] Documentation completeness
  - [ ] README.md 업데이트 확인
  - [ ] CLAUDE.md 업데이트 확인
  - [ ] CHANGELOG.md 작성 확인
  - [ ] Docstring 추가 확인

### 8.3 검증 스크립트 실행

- [ ] 삭제된 모듈 확인
  ```bash
  ! test -d korea_investment_stock/rate_limiting
  ! test -d korea_investment_stock/caching
  ! test -d korea_investment_stock/visualization
  ! test -d korea_investment_stock/batch_processing
  ! test -d korea_investment_stock/monitoring
  ! test -d korea_investment_stock/error_handling
  ```

- [ ] 라인 수 확인
  ```bash
  lines=$(wc -l < korea_investment_stock/korea_investment_stock.py)
  [ $lines -lt 1000 ] && echo "✓ Line count acceptable"
  ```

- [ ] Public 메서드 확인
  ```bash
  grep -c "^    def fetch_price(" korea_investment_stock/korea_investment_stock.py
  grep -c "^    def fetch_domestic_price(" korea_investment_stock/korea_investment_stock.py
  ```

- [ ] 데코레이터 제거 확인
  ```bash
  ! grep -q "@retry_on_rate_limit" korea_investment_stock/korea_investment_stock.py
  ! grep -q "@cacheable" korea_investment_stock/korea_investment_stock.py
  ```

### 8.4 PyPI 배포 준비

- [ ] 빌드 실행
  ```bash
  python -m build
  ```

- [ ] 패키지 검증
  ```bash
  twine check dist/*
  ```

- [ ] TestPyPI 배포 (테스트)
  ```bash
  twine upload --repository testpypi dist/*
  ```

- [ ] TestPyPI에서 설치 테스트
  ```bash
  pip install --index-url https://test.pypi.org/simple/ korea-investment-stock==0.6.0
  ```

### 8.5 배포

- [ ] PyPI 업로드
  ```bash
  twine upload dist/*
  ```

- [ ] GitHub Release 생성
  - [ ] Tag: v0.6.0
  - [ ] Title: "v0.6.0 - Simplification Release"
  - [ ] Description: CHANGELOG.md 내용 포함

- [ ] Migration guide 공지
  - [ ] GitHub Discussion 작성
  - [ ] README에 배너 추가

---

## 📊 진행 상황 요약

**전체 진행률**: 5/8 Phases 완료 (62.5%)

| Phase | 작업 | 완료 | 진행률 | 상태 |
|-------|------|------|--------|------|
| Phase 1 | 모듈 삭제 (16개 파일) | 16/16 | 100% | ✅ 완료 |
| Phase 2 | 메인 모듈 수정 | 50+/50+ | 100% | ✅ 완료 |
| Phase 3 | Package 설정 | 2/2 | 100% | ✅ 완료 |
| Phase 4 | 테스트 수정 | 16/16 | 100% | ✅ 완료 (Phase 2에서 처리됨) |
| Phase 5 | Example 수정 | 7/7 | 100% | ✅ 완료 (Phase 2에서 처리됨) |
| Phase 6 | 문서 업데이트 | 15+/15+ | 100% | ✅ 완료 |
| Phase 7 | 버전 관리 | 1/5 | 20% | 🔄 부분완료 (version 0.6.0) |
| Phase 8 | 검증 & 배포 | 0/15+ | 0% | ⏸️ 대기 |

### 상세 성과
- **삭제된 코드**: ~5,364 lines (4,434 lines from modules + 930 lines from main file)
- **파일 크기**: 1,941 lines → 1,011 lines (48% 감소)
- **완료된 커밋**: 6개 (Phase 1-3: 6개)
- **브랜치**: `feat/#40-simplify`
- **남은 테스트 파일**: 4개 (모두 Public API 사용)
- **남은 Example 파일**: 3개 (모두 단순화된 API 사용)
- **문서 업데이트**: README.md, CHANGELOG.md, __init__.py

---

## ⚠️ 주의사항

1. **순서 준수**: Phase 순서대로 진행 (Phase 1 → Phase 8)
2. **단위 커밋**: 각 Phase별로 커밋하여 롤백 가능하도록 유지
3. **테스트 우선**: 코드 변경 후 반드시 테스트 실행
4. **문서 동기화**: 코드 변경과 문서 업데이트 동시 진행
5. **Breaking Changes**: 모든 변경사항을 CHANGELOG에 명시

---

**작성**: Claude Code  
**시작일**: (To be filled)  
**완료일**: (To be filled)
