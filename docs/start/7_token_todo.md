# Token 관리 리팩토링 TODO

## Phase 1: 폴더 구조 변경

- [x] `token_storage/` 폴더를 `token/`으로 리네이밍
- [x] `token_storage.py`를 `storage.py`로 리네이밍
- [x] `test_token_storage.py`를 `test_storage.py`로 리네이밍
- [x] `token/__init__.py` import 경로 업데이트
- [x] `korea_investment_stock.py` import 경로 업데이트
- [x] 기존 테스트 실행하여 리네이밍 확인
  ```bash
  pytest korea_investment_stock/token/test_storage.py -v
  ```

## Phase 2: TokenManager 클래스 생성

- [x] `token/manager.py` 파일 생성
- [x] `TokenManager` 클래스 구현
  - [x] `__init__` 메서드 (storage, base_url, api_key, api_secret)
  - [x] `access_token` 프로퍼티
  - [x] `get_valid_token()` 메서드
  - [x] `is_token_valid()` 메서드
  - [x] `_load_token()` 메서드
  - [x] `_issue_token()` 메서드
  - [x] `_parse_token_response()` 메서드
  - [x] `issue_hashkey()` 메서드
  - [x] `invalidate()` 메서드
- [x] `token/test_manager.py` 테스트 파일 생성
  - [x] `test_get_valid_token_when_valid` 테스트
  - [x] `test_get_valid_token_when_invalid` 테스트
  - [x] `test_is_token_valid` 테스트
  - [x] `test_invalidate` 테스트
- [x] TokenManager 단위 테스트 통과 확인 (16 passed)
  ```bash
  pytest korea_investment_stock/token/test_manager.py -v
  ```

## Phase 3: TokenStorageFactory 분리

- [ ] `token/factory.py` 파일 생성
- [ ] `create_token_storage()` 함수 구현
- [ ] `_get_config_value()` 헬퍼 함수 구현
- [ ] `_create_file_storage()` 함수 구현
- [ ] `_create_redis_storage()` 함수 구현
- [ ] `token/test_factory.py` 테스트 파일 생성
  - [ ] `test_default_file_storage` 테스트
  - [ ] `test_config_file_storage` 테스트
  - [ ] `test_config_redis_storage` 테스트
  - [ ] `test_invalid_storage_type` 테스트
  - [ ] `test_env_var_storage_type` 테스트
- [ ] Factory 테스트 통과 확인
  ```bash
  pytest korea_investment_stock/token/test_factory.py -v
  ```

## Phase 4: KoreaInvestment 수정

- [ ] `token` 모듈 import 추가
  ```python
  from .token import TokenManager, create_token_storage
  ```
- [ ] `__init__`에서 `TokenManager` 초기화
- [ ] `issue_access_token()` 위임 패턴으로 변경
- [ ] `check_access_token()` 위임 패턴으로 변경
- [ ] `load_access_token()` 위임 패턴으로 변경
- [ ] `issue_hashkey()` 위임 패턴으로 변경
- [ ] `_create_token_storage()` 메서드 삭제
- [ ] 메인 클래스 테스트 통과 확인
  ```bash
  pytest korea_investment_stock/tests/test_korea_investment_stock.py -v
  ```

## Phase 5: __init__.py 업데이트

- [ ] `token/__init__.py`에 TokenManager export 추가
- [ ] `token/__init__.py`에 create_token_storage export 추가
- [ ] `__all__` 리스트 업데이트
- [ ] import 테스트
  ```bash
  python -c "from korea_investment_stock import KoreaInvestment"
  python -c "from korea_investment_stock.token import TokenManager"
  python -c "from korea_investment_stock.token import create_token_storage"
  ```

## Phase 6: 최종 검증

- [ ] 전체 테스트 통과 확인
  ```bash
  pytest -v
  ```
- [ ] 통합 테스트 통과 확인 (API 자격 증명 필요)
  ```bash
  pytest korea_investment_stock/tests/test_integration_us_stocks.py -v
  ```
- [ ] 예제 실행 확인
  ```bash
  python examples/basic_example.py
  ```
- [ ] 하위 호환성 검증 (기존 코드 동작 확인)

---

## 완료 기준

- [ ] `korea_investment_stock.py` 토큰 관련 코드 ~100줄 감소
- [ ] `token/` 폴더에 5개 파일 (storage.py, manager.py, factory.py, __init__.py, 테스트 파일들)
- [ ] 기존 API 시그니처 변경 없음 (Breaking Change 없음)
- [ ] TokenManager 테스트 커버리지 80% 이상

---

**작성일**: 2025-12-06
