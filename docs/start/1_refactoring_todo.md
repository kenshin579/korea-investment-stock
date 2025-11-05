# TODO: 프로젝트 구조 리팩토링

> **관련 문서**: [1_refactoring_prd.md](1_refactoring_prd.md) | [1_refactoring_implementation.md](1_refactoring_implementation.md)

---

## Phase 1: 디렉토리 구조 생성

### 준비 작업
- [ ] 현재 작업 브랜치 확인 (`git branch`)
- [ ] 작업 디렉토리 확인 (`pwd` → `.conductor/karachi`)
- [ ] 깨끗한 작업 환경 확인 (`git status`)

### 디렉토리 생성
- [ ] `mkdir -p korea_investment_stock/cache`
- [ ] `mkdir -p korea_investment_stock/token_storage`

### 초기화 파일 생성
- [ ] `touch korea_investment_stock/cache/__init__.py`
- [ ] `touch korea_investment_stock/token_storage/__init__.py`

### 검증
- [ ] `tree korea_investment_stock -L 2` 실행하여 구조 확인
- [ ] 생성된 디렉토리가 정상적으로 보이는지 확인

---

## Phase 2: 캐시 모듈 파일 이동

### 구현 파일 이동
- [ ] `mv korea_investment_stock/cache_manager.py korea_investment_stock/cache/`
- [ ] `mv korea_investment_stock/cached_korea_investment.py korea_investment_stock/cache/`

### 테스트 파일 이동
- [ ] `mv korea_investment_stock/tests/test_cache_manager.py korea_investment_stock/cache/`
- [ ] `mv korea_investment_stock/tests/test_cached_integration.py korea_investment_stock/cache/`

### 검증
- [ ] `ls korea_investment_stock/cache/` 실행하여 4개 파일 확인
  - `__init__.py`
  - `cache_manager.py`
  - `cached_korea_investment.py`
  - `test_cache_manager.py`
  - `test_cached_integration.py`

---

## Phase 3: 토큰 저장소 모듈 이동

### 구현 파일 이동
- [ ] `mv korea_investment_stock/token_storage.py korea_investment_stock/token_storage/token_storage.py`

### 테스트 파일 이동
- [ ] `mv korea_investment_stock/test_token_storage.py korea_investment_stock/token_storage/`

### 빈 디렉토리 정리
- [ ] `rmdir korea_investment_stock/tests/` (tests 폴더가 비어있으면 제거)

### 검증
- [ ] `ls korea_investment_stock/token_storage/` 실행하여 3개 파일 확인
  - `__init__.py`
  - `token_storage.py`
  - `test_token_storage.py`
- [ ] `tests/` 디렉토리가 제거되었는지 확인

---

## Phase 4: __init__.py 파일 작성

### cache/__init__.py
- [ ] `cache/__init__.py` 파일 열기
- [ ] 모듈 docstring 추가
- [ ] `CacheManager`, `CacheEntry`, `CachedKoreaInvestment` import
- [ ] `__all__` 리스트 정의

### token_storage/__init__.py
- [ ] `token_storage/__init__.py` 파일 열기
- [ ] 모듈 docstring 추가
- [ ] `TokenStorage`, `FileTokenStorage`, `RedisTokenStorage` import
- [ ] `__all__` 리스트 정의

### korea_investment_stock/__init__.py
- [ ] 기존 `__init__.py` 파일 백업 (복사)
- [ ] 캐시 모듈 import 추가
- [ ] 토큰 저장소 모듈 import 추가
- [ ] `__all__` 리스트 업데이트

---

## Phase 5: 내부 Import 경로 수정

### cache/cached_korea_investment.py
- [ ] 파일 열기
- [ ] `from korea_investment_stock.cache_manager import CacheManager` 찾기
- [ ] `from .cache_manager import CacheManager`로 변경
- [ ] 파일 저장

### cache/test_cache_manager.py
- [ ] 파일 열기
- [ ] `from korea_investment_stock.cache_manager import` 찾기
- [ ] `from .cache_manager import` 또는 `from korea_investment_stock.cache import`로 변경
- [ ] 파일 저장

### cache/test_cached_integration.py
- [ ] 파일 열기
- [ ] import 경로 확인 (메인 API는 변경 불필요)
- [ ] 필요시 수정

### token_storage/test_token_storage.py
- [ ] 파일 열기
- [ ] `from korea_investment_stock.token_storage import` 찾기
- [ ] `from .token_storage import`로 변경
- [ ] 파일 저장

---

## Phase 6: 검증 및 테스트

### Import 검증
- [ ] Python 인터프리터 테스트 실행:
  ```bash
  python3 -c "
  from korea_investment_stock import KoreaInvestment, CachedKoreaInvestment
  from korea_investment_stock.cache import CacheManager
  from korea_investment_stock.token_storage import FileTokenStorage, RedisTokenStorage
  print('✅ All imports successful')
  "
  ```

### 디렉토리 구조 확인
- [ ] `tree korea_investment_stock -I "__pycache__|*.pyc" -L 2` 실행
- [ ] 목표 구조와 일치하는지 확인

### 테스트 발견 확인
- [ ] `pytest korea_investment_stock --collect-only` 실행
- [ ] 모든 테스트 파일이 발견되는지 확인

### 단위 테스트 실행
- [ ] `pytest korea_investment_stock/cache/test_cache_manager.py -v`
- [ ] `pytest korea_investment_stock/cache/test_cached_integration.py -v`
- [ ] `pytest korea_investment_stock/token_storage/test_token_storage.py -v`
- [ ] `pytest korea_investment_stock/test_korea_investment_stock.py -v`

### 전체 테스트 실행
- [ ] `pytest korea_investment_stock -v`
- [ ] 실패한 테스트 기록 (스킵된 테스트 제외)
- [ ] 기존 실패 테스트와 비교하여 새로운 실패 없는지 확인

---

## Phase 7: 문서 업데이트

### CLAUDE.md
- [ ] `CLAUDE.md` 파일 열기
- [ ] "Package Structure" 섹션 찾기
- [ ] 새로운 디렉토리 구조로 업데이트
- [ ] 변경사항 저장

### CHANGELOG.md
- [ ] `CHANGELOG.md` 파일 열기
- [ ] `[Unreleased]` 섹션 찾기 (없으면 생성)
- [ ] `### Changed` 항목 추가
- [ ] 프로젝트 구조 리팩토링 내용 기록
- [ ] 변경사항 저장

---

## Phase 8: 커밋 및 정리

### Git 작업
- [ ] `git status` 실행하여 변경사항 확인
- [ ] `git add korea_investment_stock/` 실행
- [ ] `git add CLAUDE.md CHANGELOG.md` 실행
- [ ] Commit 메시지 작성:
  ```bash
  git commit -m "[refactor] Reorganize package structure into feature modules

  - Created cache/ module for caching functionality
  - Created token_storage/ module for token storage
  - Moved test files to co-locate with implementation
  - Removed tests/ directory
  - Maintained backward compatibility for all imports
  "
  ```

### 최종 검증
- [ ] `git diff HEAD~1` 실행하여 변경사항 리뷰
- [ ] 테스트 한 번 더 실행: `pytest korea_investment_stock -v`
- [ ] 예제 코드 실행 확인 (optional)

---

## 완료 기준

### 구조
- [x] `cache/` 모듈 생성 완료
- [x] `token_storage/` 모듈 생성 완료
- [x] `tests/` 디렉토리 제거 완료
- [x] 모든 테스트 파일이 co-located

### 기능
- [x] 모든 import 경로 정상 동작
- [x] 기존 테스트 모두 발견됨
- [x] 새로운 실패 테스트 없음

### 문서
- [x] CLAUDE.md 업데이트 완료
- [x] CHANGELOG.md 업데이트 완료
- [x] Git commit 완료

---

**작성일**: 2025-11-05
**버전**: 1.0
**상태**: Ready
