# PRD: 프로젝트 구조 리팩토링

> **프로젝트**: Korea Investment Stock - Project Structure Refactoring
> **작성일**: 2025-11-05
> **버전**: 1.0
> **관련 이슈**: Code Organization & Test Structure Improvement

---

## 📚 관련 문서

- **[구현 가이드](1_refactoring_implementation.md)** - 상세 구현 절차, 마이그레이션 가이드
- **[TODO 체크리스트](1_refactoring_todo.md)** - 단계별 작업 목록 및 일정

---

## 📋 Executive Summary

### 프로젝트 목표
현재 평면적인 패키지 구조를 기능별로 재구성하고, 테스트 코드를 구현 파일과 함께 배치하여 코드 가독성과 유지보수성을 향상시킵니다.

### 핵심 변경사항
- **테스트 구조**: 구현 파일과 테스트 파일을 기능별로 같은 디렉토리에 배치
- **기능별 모듈화**: 캐시, 토큰 저장소 등을 독립된 서브패키지로 분리
- **테스트 안정성**: 현재 실패하는 테스트 원인 분석 및 수정

### 기대효과
- 관련 코드 찾기 용이 (co-location)
- 기능별 독립성 향상 (모듈화)
- 테스트 유지보수성 개선 (co-located tests)
- 코드 리뷰 효율성 증대

---

## 🔍 Current State Analysis

### 1. 현재 디렉토리 구조

```
korea_investment_stock/
├── __init__.py                      # 패키지 초기화
├── korea_investment_stock.py        # 메인 클래스 (1,011 lines)
├── cache_manager.py                 # 캐시 관리자
├── cached_korea_investment.py       # 캐싱 래퍼
├── token_storage.py                 # 토큰 저장소
├── test_korea_investment_stock.py   # ❌ 메인 테스트 (잘못된 위치)
├── test_integration_us_stocks.py    # ❌ 통합 테스트 (잘못된 위치)
├── test_token_storage.py            # ❌ 토큰 저장소 테스트 (잘못된 위치)
└── tests/                           # 일부 테스트만 있는 폴더
    ├── test_cache_manager.py        # 캐시 매니저 테스트
    └── test_cached_integration.py   # 캐시 통합 테스트
```

### 2. 문제점 분석

#### 🔴 문제 1: 테스트 파일 위치 불일치
**현상:**
- 일부 테스트는 `korea_investment_stock/tests/`에 위치
- 일부 테스트는 `korea_investment_stock/` 루트에 위치
- 구현 파일과 테스트 파일이 분리되어 있음

**영향:**
- 코드 수정 시 관련 테스트 찾기 어려움
- 테스트 파일 위치가 일관되지 않아 혼란
- 새로운 기여자의 진입 장벽 증가

**Python 커뮤니티 권장사항:**
```
✅ 권장: Co-located tests (구현과 테스트가 같은 위치)

# Django, Flask, FastAPI 등 대부분의 Python 프로젝트 구조
myproject/
├── feature_a/
│   ├── __init__.py
│   ├── models.py
│   └── test_models.py          # ✅ 같은 폴더
└── feature_b/
    ├── __init__.py
    ├── views.py
    └── test_views.py            # ✅ 같은 폴더

# pytest는 자동으로 test_*.py 또는 *_test.py를 찾음
```

#### 🔴 문제 2: 기능별 모듈화 부족
**현상:**
- 캐시 관련 파일: `cache_manager.py`, `cached_korea_investment.py`가 루트에 산재
- 토큰 저장소: `token_storage.py`가 루트에 단독 존재
- 기능별 그룹화가 되어 있지 않음

**영향:**
- 기능 확장 시 파일 수 급증
- 관련 파일 찾기 어려움
- 의존성 관리 복잡도 증가

**예시:**
```python
# 현재: 어떤 파일이 캐시 관련인지 불명확
from korea_investment_stock.cache_manager import CacheManager
from korea_investment_stock.cached_korea_investment import CachedKoreaInvestment

# 개선: 캐시 기능이 명확히 그룹화됨
from korea_investment_stock.cache import CacheManager, CachedKoreaInvestment
```

#### 🔴 문제 3: 실패하는 테스트 존재
**현상:**
- 일부 테스트가 실패하거나 스킵되고 있음
- 실패 원인이 명확하지 않음

**확인 필요:**
```python
# test_korea_investment_stock.py:72
@skip("Skipping test_fetch_kospi_symbols")
def test_fetch_kospi_symbols(self):
    # 왜 스킵되었는지?

# test_korea_investment_stock.py:78
# todo: 이 unit test는 정리가 필요하다
def test_fetch_price_detail_oversea(self):
    # 무엇을 정리해야 하는지?
```

---

## 🎯 Proposed Solution

### 1. 목표 디렉토리 구조

```
korea_investment_stock/
├── __init__.py                          # 패키지 초기화, public API 노출
├── korea_investment_stock.py            # 메인 KoreaInvestment 클래스
├── test_korea_investment_stock.py       # ✅ 메인 클래스 테스트 (co-located)
├── test_integration_us_stocks.py        # ✅ 통합 테스트 (co-located)
│
├── cache/                               # 🆕 캐시 기능 모듈
│   ├── __init__.py                      # cache_manager, cached_korea_investment 노출
│   ├── cache_manager.py                 # CacheManager, CacheEntry
│   ├── test_cache_manager.py            # ✅ 캐시 매니저 테스트 (co-located)
│   ├── cached_korea_investment.py       # CachedKoreaInvestment
│   └── test_cached_integration.py       # ✅ 캐시 통합 테스트 (co-located)
│
└── token_storage/                       # 🆕 토큰 저장소 모듈
    ├── __init__.py                      # TokenStorage, FileTokenStorage, RedisTokenStorage 노출
    ├── token_storage.py                 # 토큰 저장소 구현
    └── test_token_storage.py            # ✅ 토큰 저장소 테스트 (co-located)
```

### 2. Import 경로 변경 (하위 호환성 유지)

#### 변경 전 (현재)
```python
from korea_investment_stock import KoreaInvestment, CachedKoreaInvestment
from korea_investment_stock.cache_manager import CacheManager
from korea_investment_stock.token_storage import FileTokenStorage, RedisTokenStorage
```

#### 변경 후 (신규)
```python
# 메인 API (변경 없음)
from korea_investment_stock import KoreaInvestment, CachedKoreaInvestment

# 캐시 모듈 (명확한 그룹화)
from korea_investment_stock.cache import CacheManager, CachedKoreaInvestment

# 토큰 저장소 모듈 (명확한 그룹화)
from korea_investment_stock.token_storage import FileTokenStorage, RedisTokenStorage
```

#### 하위 호환성 유지 전략
```python
# korea_investment_stock/__init__.py

# 메인 클래스
from .korea_investment_stock import KoreaInvestment

# 하위 호환성: 기존 import 경로 유지
from .cache.cache_manager import CacheManager
from .cache.cached_korea_investment import CachedKoreaInvestment
from .token_storage.token_storage import FileTokenStorage, RedisTokenStorage

__all__ = [
    'KoreaInvestment',
    'CachedKoreaInvestment',
    'CacheManager',
    'FileTokenStorage',
    'RedisTokenStorage',
]
```

### 3. 테스트 파일 배치 원칙

#### ✅ Co-location 원칙
**기본 규칙:**
- 테스트 파일은 구현 파일과 같은 디렉토리에 위치
- 파일명: `test_<module_name>.py` 또는 `<module_name>_test.py`
- pytest는 자동으로 `test_*.py` 패턴을 찾음

**예시:**
```
cache/
├── __init__.py
├── cache_manager.py
├── test_cache_manager.py      # cache_manager.py 테스트
├── cached_korea_investment.py
└── test_cached_integration.py # cached_korea_investment.py 테스트
```

#### ✅ 통합 테스트 배치
**규칙:**
- 통합 테스트는 메인 기능과 같은 레벨에 위치
- 명확한 네이밍: `test_integration_*.py`

**예시:**
```
korea_investment_stock/
├── korea_investment_stock.py
├── test_korea_investment_stock.py    # 단위 테스트
└── test_integration_us_stocks.py     # 통합 테스트 (US stocks 특화)
```

---

## 🧪 Test Analysis & Fixes

### 1. 스킵된 테스트 분석

#### test_fetch_kospi_symbols (test_korea_investment_stock.py:72)
```python
@skip("Skipping test_fetch_kospi_symbols")
def test_fetch_kospi_symbols(self):
    resp = self.kis.fetch_kospi_symbols()
    print(resp)
    self.assertEqual(resp['rt_cd'], API_RETURN_CODE["SUCCESS"])
```

**분석 필요:**
- [ ] 왜 스킵되었는지 확인
- [ ] API 변경으로 인한 실패인지?
- [ ] Mock 데이터 사용 시 문제인지?
- [ ] 실제 API 호출 시 성공하는지 확인

**조치:**
- 원인 파악 후 수정 또는 문서화
- Mock 환경에서 재현 가능하도록 개선

#### test_fetch_price_detail_oversea (test_korea_investment_stock.py:78)
```python
# todo: 이 unit test는 정리가 필요하다
def test_fetch_price_detail_oversea(self):
    stock_market_list = [
        # ("AAPL", "US"),  # 주석 처리됨
        ("QQQM", "US"), # ETF
    ]
```

**문제점:**
- AAPL 테스트가 주석 처리됨
- "정리가 필요하다"는 주석만 있고 구체적 내용 없음

**조치 필요:**
- [ ] AAPL 테스트가 실패하는 원인 확인
- [ ] ETF만 테스트하는 이유 문서화
- [ ] 일반 주식 테스트 추가 또는 제거 이유 명시

### 2. 테스트 실행 계획

#### 실행 전 준비
```bash
# 가상환경 활성화
source .venv/bin/activate

# 의존성 설치 확인
pip install -e ".[dev]"

# 환경 변수 확인
echo $KOREA_INVESTMENT_API_KEY
echo $KOREA_INVESTMENT_API_SECRET
echo $KOREA_INVESTMENT_ACCOUNT_NO
```

#### 테스트 실행 순서
```bash
# 1. 전체 테스트 실행 (현재 상태 확인)
pytest korea_investment_stock -v --tb=short > test_results_before.txt 2>&1

# 2. 실패 테스트만 확인
pytest korea_investment_stock --lf -v

# 3. 리팩토링 후 테스트 실행
pytest korea_investment_stock -v --tb=short > test_results_after.txt 2>&1

# 4. 결과 비교
diff test_results_before.txt test_results_after.txt
```

#### 예상 테스트 결과
- ✅ `test_cache_manager.py`: 모두 통과 예상
- ✅ `test_cached_integration.py`: 모두 통과 예상
- ✅ `test_token_storage.py`: fakeredis 설치 시 통과
- ⚠️ `test_korea_investment_stock.py`: 일부 스킵/실패 예상
- ⚠️ `test_integration_us_stocks.py`: 실제 API 필요

---

## ✅ Success Criteria

### 1. 구조 개선
- [x] 캐시 관련 파일이 `cache/` 디렉토리에 그룹화됨
- [x] 토큰 저장소 파일이 `token_storage/` 디렉토리에 그룹화됨
- [x] 모든 테스트가 구현 파일과 같은 디렉토리에 위치
- [x] `tests/` 디렉토리가 제거되고 파일이 적절히 재배치됨

### 2. 하위 호환성
- [x] 기존 import 경로가 모두 동작함
- [x] 외부 사용자 코드 수정 불필요
- [x] 예제 코드가 그대로 동작함

### 3. 테스트 안정성
- [x] 모든 테스트가 새 구조에서 실행됨
- [x] 스킵된 테스트 원인이 문서화됨
- [x] 실패하는 테스트가 수정되거나 이유가 명확함

### 4. 문서화
- [x] 새 구조가 CLAUDE.md에 반영됨
- [x] 마이그레이션 가이드 작성됨
- [x] 변경 사항이 CHANGELOG.md에 기록됨

---

## ⚠️ Risks & Mitigation

### Risk 1: Import 경로 깨짐
**위험도**: 🟡 중간
**내용**: 내부 import가 깨져서 패키지가 동작하지 않음

**완화 전략:**
- Phase별 점진적 마이그레이션
- 각 Phase 후 테스트 실행
- `__init__.py`에서 하위 호환 경로 유지

### Risk 2: 테스트 발견 실패
**위험도**: 🟢 낮음
**내용**: pytest가 새 위치의 테스트를 찾지 못함

**완화 전략:**
- pytest는 기본적으로 `test_*.py` 패턴을 모든 디렉토리에서 찾음
- `pytest.ini` 또는 `pyproject.toml`에서 testpaths 확인

### Risk 3: 외부 패키지 호환성
**위험도**: 🟢 낮음 (사용자 거의 없음)
**내용**: 외부에서 내부 모듈을 직접 import하는 경우

**완화 전략:**
- Public API만 `__init__.py`에 노출
- 내부 구조는 private으로 간주
- 문서에 권장 import 방법 명시

---

## 📚 References

### Python 프로젝트 구조 참고
- **Django**: 앱별 tests.py 또는 tests/ 디렉토리
- **Flask**: 각 모듈과 함께 test_*.py
- **FastAPI**: 기능별 디렉토리 + co-located tests
- **pytest 공식**: test discovery 패턴 문서

### 관련 이슈
- v0.6.0: 프로젝트 단순화 (#40)
- 철학: "Simple, transparent, flexible"

### 프로젝트 철학
> "Simple, transparent, flexible - let users implement features their way"

**리팩토링 원칙:**
- ✅ Simple: 명확한 디렉토리 구조
- ✅ Transparent: 기능별 그룹화로 찾기 쉬움
- ✅ Flexible: 사용자는 필요한 기능만 import

---

## 📂 관련 문서

- **[구현 가이드](1_refactoring_implementation.md)** - 상세 구현 절차, 마이그레이션 가이드
- **[TODO 체크리스트](1_refactoring_todo.md)** - 단계별 작업 목록 및 검증

---

**작성일**: 2025-11-05
**버전**: 1.0
**상태**: Ready
