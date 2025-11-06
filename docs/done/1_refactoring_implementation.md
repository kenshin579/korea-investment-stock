# 구현 가이드: 프로젝트 구조 리팩토링

> **관련 문서**: [1_refactoring_prd.md](1_refactoring_prd.md)
> **작업 목록**: [1_refactoring_todo.md](1_refactoring_todo.md)

---

## 📦 Phase 1: 디렉토리 구조 생성

### 1.1 서브패키지 디렉토리 생성

```bash
# 캐시 모듈 디렉토리
mkdir -p korea_investment_stock/cache

# 토큰 저장소 모듈 디렉토리
mkdir -p korea_investment_stock/token_storage
```

### 1.2 __init__.py 파일 생성

```bash
# 캐시 모듈 초기화 파일
touch korea_investment_stock/cache/__init__.py

# 토큰 저장소 모듈 초기화 파일
touch korea_investment_stock/token_storage/__init__.py
```

---

## 🔄 Phase 2: 캐시 모듈 파일 이동

### 2.1 구현 파일 이동

```bash
# 캐시 관리자
mv korea_investment_stock/cache_manager.py \
   korea_investment_stock/cache/cache_manager.py

# 캐싱 래퍼
mv korea_investment_stock/cached_korea_investment.py \
   korea_investment_stock/cache/cached_korea_investment.py
```

### 2.2 테스트 파일 이동

```bash
# tests/ 폴더에서 cache/ 폴더로 이동
mv korea_investment_stock/tests/test_cache_manager.py \
   korea_investment_stock/cache/test_cache_manager.py

mv korea_investment_stock/tests/test_cached_integration.py \
   korea_investment_stock/cache/test_cached_integration.py
```

---

## 🔐 Phase 3: 토큰 저장소 모듈 이동

### 3.1 구현 파일 이동

```bash
# 토큰 저장소 (파일명 유지 위해 서브디렉토리 생성)
mv korea_investment_stock/token_storage.py \
   korea_investment_stock/token_storage/token_storage.py
```

### 3.2 테스트 파일 이동

```bash
# 루트에서 token_storage/ 폴더로 이동
mv korea_investment_stock/test_token_storage.py \
   korea_investment_stock/token_storage/test_token_storage.py
```

### 3.3 빈 디렉토리 제거

```bash
# 모든 파일이 이동되었으면 tests/ 디렉토리 제거
rmdir korea_investment_stock/tests/
```

---

## 📝 Phase 4: __init__.py 작성

### 4.1 cache/__init__.py

```python
"""캐시 기능 모듈

Memory-based caching for Korea Investment API responses.
"""

from .cache_manager import CacheManager, CacheEntry
from .cached_korea_investment import CachedKoreaInvestment

__all__ = [
    'CacheManager',
    'CacheEntry',
    'CachedKoreaInvestment',
]
```

### 4.2 token_storage/__init__.py

```python
"""토큰 저장소 모듈

File-based and Redis-based token storage implementations.
"""

from .token_storage import TokenStorage, FileTokenStorage, RedisTokenStorage

__all__ = [
    'TokenStorage',
    'FileTokenStorage',
    'RedisTokenStorage',
]
```

### 4.3 korea_investment_stock/__init__.py (수정)

```python
"""Korea Investment Stock API Wrapper

Simple, transparent, and flexible Python wrapper for Korea Investment Securities OpenAPI.
"""

# 메인 클래스
from .korea_investment_stock import KoreaInvestment, API_RETURN_CODE

# 캐시 기능 (서브패키지)
from .cache import CacheManager, CachedKoreaInvestment

# 토큰 저장소 (서브패키지)
from .token_storage import FileTokenStorage, RedisTokenStorage

__all__ = [
    # 메인 API
    'KoreaInvestment',
    'API_RETURN_CODE',

    # 캐시 기능
    'CacheManager',
    'CachedKoreaInvestment',

    # 토큰 저장소
    'FileTokenStorage',
    'RedisTokenStorage',
]

__version__ = "0.7.0"
```

---

## 🔗 Phase 5: 내부 Import 수정

### 5.1 cached_korea_investment.py

**파일**: `korea_investment_stock/cache/cached_korea_investment.py`

```python
# Before
from korea_investment_stock.cache_manager import CacheManager

# After
from .cache_manager import CacheManager
```

### 5.2 테스트 파일 Import 수정

**cache/test_cache_manager.py**:
```python
# Before
from korea_investment_stock.cache_manager import CacheManager, CacheEntry

# After
from .cache_manager import CacheManager, CacheEntry
# 또는
from korea_investment_stock.cache import CacheManager, CacheEntry
```

**cache/test_cached_integration.py**:
```python
# Before
from korea_investment_stock import KoreaInvestment, CachedKoreaInvestment

# After (변경 없음 - 메인 __init__.py에서 export됨)
from korea_investment_stock import KoreaInvestment, CachedKoreaInvestment
```

**token_storage/test_token_storage.py**:
```python
# Before
from korea_investment_stock.token_storage import (
    TokenStorage,
    FileTokenStorage,
    RedisTokenStorage,
)

# After
from .token_storage import (
    TokenStorage,
    FileTokenStorage,
    RedisTokenStorage,
)
# 또는
from korea_investment_stock.token_storage import (
    TokenStorage,
    FileTokenStorage,
    RedisTokenStorage,
)
```

---

## ✅ Phase 6: 검증

### 6.1 Import 테스트

```python
# Python 인터프리터에서 확인
python3 -c "
from korea_investment_stock import KoreaInvestment, CachedKoreaInvestment
from korea_investment_stock.cache import CacheManager
from korea_investment_stock.token_storage import FileTokenStorage, RedisTokenStorage
print('✅ All imports successful')
"
```

### 6.2 테스트 실행

```bash
# 전체 테스트 실행
pytest korea_investment_stock -v

# 캐시 모듈만 테스트
pytest korea_investment_stock/cache -v

# 토큰 저장소 모듈만 테스트
pytest korea_investment_stock/token_storage -v
```

### 6.3 디렉토리 구조 확인

```bash
tree korea_investment_stock -I "__pycache__|*.pyc" -L 2
```

**예상 결과**:
```
korea_investment_stock/
├── __init__.py
├── korea_investment_stock.py
├── test_korea_investment_stock.py
├── test_integration_us_stocks.py
├── cache/
│   ├── __init__.py
│   ├── cache_manager.py
│   ├── test_cache_manager.py
│   ├── cached_korea_investment.py
│   └── test_cached_integration.py
└── token_storage/
    ├── __init__.py
    ├── token_storage.py
    └── test_token_storage.py
```

---

## 📚 Phase 7: 문서 업데이트

### 7.1 CLAUDE.md 수정

**Package Structure 섹션 업데이트**:

```markdown
### Package Structure

```
korea_investment_stock/
├── __init__.py                      # Module exports
├── korea_investment_stock.py        # Main KoreaInvestment class
├── test_korea_investment_stock.py   # Main class tests
├── test_integration_us_stocks.py    # Integration tests
│
├── cache/                           # Cache module
│   ├── __init__.py
│   ├── cache_manager.py             # CacheManager, CacheEntry
│   ├── test_cache_manager.py        # Cache manager tests
│   ├── cached_korea_investment.py   # CachedKoreaInvestment wrapper
│   └── test_cached_integration.py   # Cache integration tests
│
└── token_storage/                   # Token storage module
    ├── __init__.py
    ├── token_storage.py             # FileTokenStorage, RedisTokenStorage
    └── test_token_storage.py        # Token storage tests
```

**Dependencies:** `requests`, `pandas` (minimal)
\```
```

### 7.2 CHANGELOG.md 추가

```markdown
## [Unreleased]

### Changed
- **Project Structure**: Reorganized package into feature-based modules
  - Created `cache/` module for caching functionality
  - Created `token_storage/` module for token storage implementations
  - Moved test files to co-locate with implementation files (co-located tests)
  - Removed `tests/` directory in favor of feature-specific test files
  - All existing import paths remain compatible (backward compatible)
```

---

## 🔍 검증 체크리스트

### 구조 검증
- [ ] `cache/` 디렉토리 생성됨
- [ ] `token_storage/` 디렉토리 생성됨
- [ ] `tests/` 디렉토리 제거됨
- [ ] 모든 파일이 올바른 위치로 이동됨

### Import 검증
- [ ] `from korea_investment_stock import KoreaInvestment` 동작
- [ ] `from korea_investment_stock import CachedKoreaInvestment` 동작
- [ ] `from korea_investment_stock.cache import CacheManager` 동작
- [ ] `from korea_investment_stock.token_storage import FileTokenStorage` 동작

### 테스트 검증
- [ ] `pytest korea_investment_stock -v` 모든 테스트 발견
- [ ] 캐시 모듈 테스트 통과
- [ ] 토큰 저장소 테스트 통과
- [ ] 메인 클래스 테스트 통과

### 문서 검증
- [ ] CLAUDE.md 업데이트됨
- [ ] CHANGELOG.md 업데이트됨
- [ ] 예제 코드 동작 확인

---

## ⚠️ 주의사항

### Import 수정 시
- **상대 import 사용**: 같은 패키지 내에서는 `.` 사용
- **절대 import 유지**: 외부에서 사용하는 공개 API는 `from korea_investment_stock import ...`

### 테스트 실행
- **pytest discovery**: pytest는 자동으로 `test_*.py` 패턴을 찾음
- **상대 경로 주의**: 테스트에서 같은 패키지 모듈 import 시 `.` 사용 가능

### 하위 호환성
- **Public API 유지**: `__init__.py`에서 기존 클래스 모두 export
- **기존 코드 동작**: 외부 사용자 코드 수정 불필요

---

**작성일**: 2025-11-05
**버전**: 1.0
