# TODO: Korea Investment Stock 단순화 체크리스트

> 이 문서는 [PRD](prd.md)와 [Implementation Guide](implementation.md)의 구현 체크리스트입니다.

**진행 상태 범례**:
- [ ] 미완료
- [x] 완료
- [~] 진행중

---

## Phase 1: 모듈 삭제 (우선순위: HIGH)

### 1.1 rate_limiting/ 디렉토리 전체 삭제

- [ ] `enhanced_rate_limiter.py` (~400 lines)
- [ ] `enhanced_backoff_strategy.py` (~300 lines)
- [ ] `enhanced_retry_decorator.py` (~200 lines)
- [ ] `__init__.py` (~50 lines)

```bash
rm -rf korea_investment_stock/rate_limiting/
```

### 1.2 caching/ 디렉토리 전체 삭제

- [ ] `ttl_cache.py` (~500 lines)
- [ ] `market_hours.py` (~100 lines)
- [ ] `__init__.py` (~50 lines)

```bash
rm -rf korea_investment_stock/caching/
```

### 1.3 visualization/ 디렉토리 전체 삭제

- [ ] `plotly_visualizer.py` (~400 lines)
- [ ] `dashboard.py` (~350 lines)
- [ ] `charts.py` (~250 lines)
- [ ] `__init__.py` (~50 lines)

```bash
rm -rf korea_investment_stock/visualization/
```

### 1.4 batch_processing/ 디렉토리 전체 삭제

- [ ] `dynamic_batch_controller.py` (~300 lines)
- [ ] `__init__.py` (~30 lines)

```bash
rm -rf korea_investment_stock/batch_processing/
```

### 1.5 monitoring/ 디렉토리 전체 삭제

- [ ] `stats_manager.py` (~600 lines)
- [ ] `__init__.py` (~30 lines)

```bash
rm -rf korea_investment_stock/monitoring/
```

### 1.6 error_handling/ 디렉토리 전체 삭제

- [ ] `error_recovery_system.py` (~500 lines)
- [ ] `__init__.py` (~30 lines)

```bash
rm -rf korea_investment_stock/error_handling/
```

### 1.7 legacy/ 디렉토리 전체 삭제 (선택사항)

- [ ] `rate_limiter_v1.py`

```bash
rm -rf korea_investment_stock/legacy/
```

**예상 결과**: ~4,090 lines 삭제

---

## Phase 2: 메인 모듈 수정 (우선순위: HIGH)

**파일**: `korea_investment_stock/korea_investment_stock.py`

### 2.1 Import 문 제거

- [ ] Rate limiting imports (4줄)
  ```python
  from .rate_limiting.enhanced_rate_limiter import EnhancedRateLimiter
  from .rate_limiting.enhanced_backoff_strategy import get_backoff_strategy
  from .rate_limiting.enhanced_retry_decorator import retry_on_rate_limit, retry_on_network_error
  ```

- [ ] Error handling imports (1줄)
  ```python
  from .error_handling.error_recovery_system import get_error_recovery_system
  ```

- [ ] Monitoring imports (1줄)
  ```python
  from .monitoring.stats_manager import get_stats_manager
  ```

- [ ] Caching imports (1줄)
  ```python
  from .caching import TTLCache, cacheable, CACHE_TTL_CONFIG
  ```

- [ ] Visualization imports (3줄)
  ```python
  try:
      from .visualization import PlotlyVisualizer, DashboardManager
      VISUALIZATION_AVAILABLE = True
  except ImportError:
      VISUALIZATION_AVAILABLE = False
  ```

### 2.2 __init__() 메서드 간소화

- [ ] Rate limiter 초기화 제거
  ```python
  self.rate_limiter = EnhancedRateLimiter(...)
  self.backoff_strategy = get_backoff_strategy()
  self._rate_limit_semaphore = threading.Semaphore(max_workers)
  ```

- [ ] Cache 초기화 제거
  ```python
  self.cache = TTLCache(...)
  self.cache_enabled = cache_enabled
  ```

- [ ] ThreadPoolExecutor 초기화 제거
  ```python
  self.executor = ThreadPoolExecutor(max_workers=max_workers)
  self._shutdown_event = threading.Event()
  atexit.register(self.shutdown)
  ```

- [ ] Semaphore 초기화 제거
  ```python
  self._rate_limit_semaphore = threading.Semaphore(max_workers)
  ```

- [ ] Visualizer 초기화 제거
  ```python
  if VISUALIZATION_AVAILABLE:
      self.visualizer = PlotlyVisualizer()
      self.dashboard_manager = DashboardManager()
  ```

- [ ] Stats manager 초기화 제거
  ```python
  self.stats_manager = get_stats_manager()
  ```

- [ ] Error recovery 초기화 제거
  ```python
  self.error_recovery = get_error_recovery_system()
  ```

- [ ] atexit.register() 제거 (또는 간소화)

- [ ] max_workers, cache_enabled 파라미터 제거
  ```python
  # BEFORE
  def __init__(self, api_key, api_secret, acc_no, mock=True, max_workers=3, cache_enabled=True)
  
  # AFTER
  def __init__(self, api_key, api_secret, acc_no, mock=True)
  ```

- [ ] Docstring 업데이트

### 2.3 List 기반 메서드 제거 (6개)

- [ ] `fetch_price_list()` 삭제 (Line ~817)
- [ ] `fetch_price_list_with_batch()` 삭제 (Line ~820)
- [ ] `fetch_price_list_with_dynamic_batch()` 삭제 (Line ~840)
- [ ] `fetch_stock_info_list()` 삭제 (Line ~1262)
- [ ] `fetch_search_stock_info_list()` 삭제 - 첫 번째 정의 (Line ~814)
- [ ] `fetch_search_stock_info_list()` 삭제 - 두 번째 정의 (Line ~1302)
- [ ] `fetch_price_detail_oversea_list()` 삭제 (Line ~1212)

### 2.4 내부 실행 메서드 제거 (2개)

- [ ] `__execute_concurrent_requests()` 삭제 (Line ~290, ~150 lines)
- [ ] `__execute_concurrent_requests_with_cache()` 삭제 (Line ~1349, ~80 lines)

### 2.5 Private 메서드 → Public 전환 (8개)

#### __fetch_price() → fetch_price()

- [ ] 메서드명 변경: `__fetch_price` → `fetch_price` (Line ~865)
- [ ] Docstring 추가 (Args, Returns, Example 포함)
- [ ] `__get_symbol_type` 호출을 `get_symbol_type`으로 변경
- [ ] `__fetch_etf_domestic_price` 호출을 `fetch_etf_domestic_price`로 변경
- [ ] `__fetch_domestic_price` 호출을 `fetch_domestic_price`로 변경
- [ ] `__fetch_price_detail_oversea` 호출을 `fetch_price_detail_oversea`로 변경

#### __get_symbol_type() → get_symbol_type()

- [ ] 메서드명 변경: `__get_symbol_type` → `get_symbol_type` (Line ~893)
- [ ] Docstring 추가

#### __fetch_etf_domestic_price() → fetch_etf_domestic_price()

- [ ] 메서드명 변경: `__fetch_etf_domestic_price` → `fetch_etf_domestic_price` (Line ~907)
- [ ] Docstring 추가
- [ ] `@cacheable` 데코레이터 제거 (Line ~902)
- [ ] `@retry_on_rate_limit` 데코레이터 제거 (Line ~906)
- [ ] Rate limiter 코드 제거: `with self.rate_limiter.acquire():`

#### __fetch_domestic_price() → fetch_domestic_price()

- [ ] 메서드명 변경: `__fetch_domestic_price` → `fetch_domestic_price` (Line ~940)
- [ ] Docstring 추가
- [ ] `@cacheable` 데코레이터 제거 (Line ~935)
- [ ] `@retry_on_rate_limit` 데코레이터 제거 (Line ~939)
- [ ] Rate limiter 코드 제거: `with self.rate_limiter.acquire():`

#### __fetch_price_detail_oversea() → fetch_price_detail_oversea()

- [ ] 메서드명 변경: `__fetch_price_detail_oversea` → `fetch_price_detail_oversea` (Line ~1220)
- [ ] Docstring 추가
- [ ] `@cacheable` 데코레이터 제거 (Line ~1215)
- [ ] `@retry_on_rate_limit` 데코레이터 제거 (Line ~1219)
- [ ] Rate limiter 코드 제거

#### __fetch_stock_info() → fetch_stock_info()

- [ ] 메서드명 변경: `__fetch_stock_info` → `fetch_stock_info` (Line ~1270)
- [ ] Docstring 추가
- [ ] `@cacheable` 데코레이터 제거 (Line ~1265)
- [ ] `@retry_on_rate_limit` 데코레이터 제거 (Line ~1269)
- [ ] Rate limiter 코드 제거

#### __fetch_search_stock_info() → fetch_search_stock_info()

- [ ] 메서드명 변경: `__fetch_search_stock_info` → `fetch_search_stock_info` (Line ~1310)
- [ ] Docstring 추가
- [ ] `@cacheable` 데코레이터 제거 (Line ~1305)
- [ ] `@retry_on_rate_limit` 데코레이터 제거 (Line ~1309)
- [ ] Rate limiter 코드 제거

#### __handle_rate_limit_error() 삭제

- [ ] `__handle_rate_limit_error()` 메서드 완전 삭제 (Line ~583, DEPRECATED)

### 2.6 Cache 관련 메서드 제거 (5개)

- [ ] `clear_cache()` 삭제 (Line ~1452)
- [ ] `get_cache_stats()` 삭제 (Line ~1471)
- [ ] `set_cache_enabled()` 삭제 (Line ~1498)
- [ ] `preload_cache()` 삭제 (Line ~1507)

### 2.7 Monitoring 관련 메서드 제거 (7개)

- [ ] `create_monitoring_dashboard()` 삭제 (Line ~1536)
- [ ] `save_monitoring_dashboard()` 삭제 (Line ~1570)
- [ ] `create_stats_report()` 삭제 (Line ~1591)
- [ ] `get_system_health_chart()` 삭제 (Line ~1612)
- [ ] `get_api_usage_chart()` 삭제 (Line ~1634)
- [ ] `show_monitoring_dashboard()` 삭제 (Line ~1668)

### 2.8 나머지 메서드 데코레이터 제거

- [ ] `issue_access_token()` - `@retry_on_rate_limit` 제거 (Line ~731)
- [ ] `fetch_kospi_symbols()` - `@cacheable` 제거 (Line ~968)
- [ ] `fetch_kosdaq_symbols()` - `@cacheable` 제거 (Line ~1002)
- [ ] `fetch_ipo_schedule()` - `@cacheable` 제거 (Line ~1780)
- [ ] `fetch_ipo_schedule()` - `@retry_on_rate_limit` 제거 (Line ~1784)

### 2.9 shutdown() 메서드 간소화

- [ ] ThreadPoolExecutor shutdown 코드 제거
  ```python
  if hasattr(self, 'executor'):
      self.executor.shutdown(wait=True)
  ```

- [ ] Event 처리 제거
  ```python
  if self._shutdown_event.is_set():
      return
  self._shutdown_event.set()
  ```

- [ ] Stats 저장 코드 제거
  ```python
  if hasattr(self, 'stats_manager'):
      self.stats_manager.save_all_stats()
  ```

- [ ] 간소화된 버전으로 교체 또는 완전히 제거 검토

**예상 결과**: 1,941 lines → ~800 lines

---

## Phase 3: Package 설정 수정 (우선순위: HIGH)

**파일**: `korea_investment_stock/__init__.py`

### 3.1 Import 문 정리

- [ ] Rate limiting imports 제거
  ```python
  from .rate_limiting.enhanced_rate_limiter import EnhancedRateLimiter
  from .rate_limiting.enhanced_retry_decorator import retry_on_rate_limit, retry_on_network_error
  from .rate_limiting.enhanced_backoff_strategy import EnhancedBackoffStrategy, get_backoff_strategy
  ```

- [ ] Error handling imports 제거
  ```python
  from .error_handling.error_recovery_system import ErrorRecoverySystem, get_error_recovery_system
  ```

- [ ] Batch processing imports 제거
  ```python
  from .batch_processing.dynamic_batch_controller import DynamicBatchController
  ```

- [ ] Monitoring imports 제거
  ```python
  from .monitoring.stats_manager import StatsManager, get_stats_manager
  ```

### 3.2 __all__ 리스트 업데이트

- [ ] 제거된 모듈 exports 삭제
- [ ] 핵심 4개만 유지: `KoreaInvestment`, `MARKET_CODE_MAP`, `EXCHANGE_CODE_MAP`, `API_RETURN_CODE`

**예상 결과**: 36 lines → ~10 lines

---

## Phase 4: 테스트 수정 (우선순위: MEDIUM)

### 4.1 테스트 파일 삭제 (12개)

- [ ] `test_rate_limiter.py` 삭제
- [ ] `test_enhanced_backoff.py` 삭제
- [ ] `test_rate_limit_error_detection.py` 삭제
- [ ] `test_rate_limit_simulation.py` 삭제
- [ ] `test_ttl_cache.py` 삭제
- [ ] `test_cache_integration.py` 삭제
- [ ] `test_batch_processing.py` 삭제
- [ ] `test_error_recovery.py` 삭제
- [ ] `test_error_handling.py` 삭제
- [ ] `test_stats_save.py` 삭제
- [ ] `test_enhanced_integration.py` 삭제
- [ ] `test_threadpool_improvement.py` 삭제

```bash
cd korea_investment_stock/tests/
rm test_rate_limiter.py test_enhanced_backoff.py test_rate_limit_error_detection.py
rm test_rate_limit_simulation.py test_ttl_cache.py test_cache_integration.py
rm test_batch_processing.py test_error_recovery.py test_error_handling.py
rm test_stats_save.py test_enhanced_integration.py test_threadpool_improvement.py
```

### 4.2 test_korea_investment_stock.py 업데이트

- [ ] `TestBatchProcessing` 클래스 삭제
  - [ ] `test_fetch_price_list()` 삭제
  - [ ] `test_fetch_price_list_with_batch()` 삭제
  - [ ] `test_concurrent_requests()` 삭제

- [ ] `TestCaching` 클래스 삭제
  - [ ] `test_cache_hit()` 삭제
  - [ ] `test_cache_miss()` 삭제
  - [ ] `test_clear_cache()` 삭제

- [ ] `TestSingleFetch` 클래스 추가
  - [ ] `test_fetch_price()` 추가 (KR)
  - [ ] `test_fetch_price()` 추가 (US)
  - [ ] `test_fetch_domestic_price()` 추가
  - [ ] `test_fetch_etf_domestic_price()` 추가
  - [ ] `test_fetch_stock_info()` 추가
  - [ ] `test_get_symbol_type()` 추가

### 4.3 test_integration.py 업데이트

- [ ] `fetch_price_list()` 호출을 `fetch_price()` loop로 변경
- [ ] Rate limiter 검증 코드 제거
- [ ] Context manager 패턴 제거 (선택사항)

### 4.4 test_integration_us_stocks.py 업데이트

- [ ] `fetch_price_detail_oversea()` Public 호출 테스트 추가
- [ ] List 기반 메서드 호출 제거

### 4.5 test_ipo_schedule.py 업데이트

- [ ] 데코레이터 제거 반영 (동작 변경 없음)
- [ ] 테스트 케이스 검증

### 4.6 test_ipo_integration.py 업데이트

- [ ] 변경사항 검증
- [ ] 필요시 테스트 케이스 수정

### 4.7 test_load.py 업데이트

- [ ] ThreadPoolExecutor 사용 제거
- [ ] 단순 for loop 기반으로 변경
- [ ] Rate limiting 검증 제거
- [ ] 순차 부하 테스트로 변경

### 4.8 test_public_api.py 생성 (신규)

- [ ] 파일 생성: `korea_investment_stock/tests/test_public_api.py`
- [ ] `TestPublicAPI` 클래스 작성
  - [ ] `test_fetch_price_kr()` 추가
  - [ ] `test_fetch_price_us()` 추가
  - [ ] `test_user_controlled_batch()` 추가
  - [ ] `test_user_controlled_retry()` 추가
- [ ] `TestUserImplementation` 클래스 작성
  - [ ] `test_user_caching()` 추가

---

## Phase 5: Example 파일 수정 (우선순위: MEDIUM)

### 5.1 Example 파일 삭제 (4개)

- [ ] `examples/rate_limiting_example.py` 삭제
- [ ] `examples/stats_management_example.py` 삭제
- [ ] `examples/stats_visualization_plotly.py` 삭제
- [ ] `examples/visualization_integrated_example.py` 삭제

```bash
cd examples/
rm rate_limiting_example.py stats_management_example.py
rm stats_visualization_plotly.py visualization_integrated_example.py
```

### 5.2 ipo_schedule_example.py 업데이트

- [ ] Context manager 제거 (선택사항)
- [ ] 단순화된 사용법으로 변경
- [ ] 주석 업데이트

### 5.3 us_stock_price_example.py 업데이트

- [ ] `fetch_price_list()` → `fetch_price()` loop로 변경
- [ ] 사용자 제어 배치 조회 예시 추가
- [ ] Rate limiting 코드 예시 추가 (time.sleep)

### 5.4 basic_usage_example.py 생성 (신규)

- [ ] 파일 생성: `examples/basic_usage_example.py`
- [ ] 기본 초기화 예시
- [ ] 단일 조회 예시
- [ ] 배치 조회 (사용자 제어) 예시
- [ ] 재시도 로직 구현 예시
- [ ] 캐싱 구현 예시
- [ ] IPO 조회 예시

---

## Phase 6: 문서 업데이트 (우선순위: HIGH)

### 6.1 README.md 업데이트

#### Features 섹션
- [ ] Rate limiting 항목 제거
- [ ] Cache 항목 제거
- [ ] Visualization 항목 제거
- [ ] Batch processing 항목 제거
- [ ] 단순화된 기능 목록으로 교체

#### Usage 섹션
- [ ] 복잡한 사용 예시 제거
- [ ] 단순한 사용 예시로 교체
- [ ] 사용자 제어 패턴 강조

#### Migration Guide 링크
- [ ] Migration 섹션 추가
- [ ] PRD 링크 추가

### 6.2 CLAUDE.md 업데이트

#### Architecture 섹션
- [ ] 제거된 모듈 설명 삭제
  - [ ] rate_limiting/ 섹션 삭제
  - [ ] caching/ 섹션 삭제
  - [ ] visualization/ 섹션 삭제
  - [ ] batch_processing/ 섹션 삭제
  - [ ] monitoring/ 섹션 삭제
  - [ ] error_handling/ 섹션 삭제

- [ ] 단순화된 아키텍처 다이어그램 추가
- [ ] Singleton Patterns 섹션 삭제
- [ ] Threading & Concurrency 섹션 삭제

#### API Methods 섹션
- [ ] Private → Public 전환 문서화
- [ ] 제거된 메서드 목록 삭제
- [ ] 새로운 Public API 목록 작성

#### Performance Characteristics
- [ ] Benchmark 섹션 업데이트 (변경 반영)

### 6.3 CHANGELOG.md 업데이트

- [ ] 0.6.0 버전 추가
- [ ] Breaking Changes 명시
  - [ ] Removed Modules 목록
  - [ ] Removed Methods 목록
  - [ ] Changed Methods 목록
- [ ] Migration guide 링크 추가

### 6.4 MIGRATION.md 생성 (신규) - 선택사항

- [ ] 파일 생성: `docs/MIGRATION.md`
- [ ] Breaking changes 상세 ���명
- [ ] Before/After 코드 예시 (5가지)
  - [ ] 단일 조회
  - [ ] 배치 조회
  - [ ] 캐싱
  - [ ] 모니터링
  - [ ] IPO 조회
- [ ] 권장 migration 전략 (Phase 1-4)

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

**전체 진행률**: 0/100 (0%)

| Phase | 작업 | 완료 | 진행률 |
|-------|------|------|--------|
| Phase 1 | 모듈 삭제 (16개 파일) | 0/16 | 0% |
| Phase 2 | 메인 모듈 수정 | 0/50+ | 0% |
| Phase 3 | Package 설정 | 0/2 | 0% |
| Phase 4 | 테스트 수정 | 0/20+ | 0% |
| Phase 5 | Example 수정 | 0/8 | 0% |
| Phase 6 | 문서 업데이트 | 0/15+ | 0% |
| Phase 7 | 버전 관리 | 0/5 | 0% |
| Phase 8 | 검증 & 배포 | 0/15+ | 0% |

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
