# Testcontainers 도입 TODO

## Phase 1: 기본 인프라 구축

- [ ] `pyproject.toml`에 testcontainers 의존성 추가
  ```toml
  dev = [
      ...
      "testcontainers>=4.0.0",
  ]
  ```

- [ ] `pyproject.toml`에 pytest markers 설정 추가
  ```toml
  [tool.pytest.ini_options]
  markers = [
      "unit: 단위 테스트 (fakeredis, Docker 불필요)",
      "integration: 통합 테스트 (Docker 필요)",
  ]
  ```

- [ ] `korea_investment_stock/tests/` 디렉토리 생성

- [ ] `korea_investment_stock/tests/__init__.py` 생성

- [ ] `korea_investment_stock/tests/conftest.py` 생성
  - [ ] `pytest_configure()` 함수 구현 (marker 등록)
  - [ ] `redis_container` fixture 구현 (scope="session")
  - [ ] `redis_url` fixture 구현
  - [ ] `redis_client` fixture 구현 (FLUSHDB 포함)

## Phase 2: Redis 통합 테스트 구현

- [ ] `korea_investment_stock/tests/test_redis_integration.py` 생성

- [ ] `TestRedisTokenStorageIntegration` 클래스 구현
  - [ ] `test_save_and_load` - 토큰 저장/로드
  - [ ] `test_check_token_valid` - 토큰 유효성 확인
  - [ ] `test_delete_token` - 토큰 삭제
  - [ ] `test_connection_pool` - 다중 스레드 연결 풀
  - [ ] `test_ttl_actual_expiry` - 실제 TTL 만료 확인

- [ ] `TestRedisConnectionHandling` 클래스 구현
  - [ ] `test_reconnection_after_flushdb` - FLUSHDB 후 동작
  - [ ] `test_multiple_databases` - 다중 DB 격리

## Phase 3: 테스트 검증

- [ ] 단위 테스트 실행 확인 (Docker 없이)
  ```bash
  pytest -m "not integration" -v
  ```

- [ ] 통합 테스트 실행 확인 (Docker 필요)
  ```bash
  pytest -m integration -v
  ```

- [ ] 전체 테스트 실행 확인
  ```bash
  pytest -v
  ```

- [ ] Docker 미설치 환경에서 graceful skip 확인

## Phase 4: CI/CD 통합

- [ ] `.github/workflows/test.yml` 수정 (또는 생성)
  - [ ] `unit-tests` job 추가
  - [ ] `integration-tests` job 추가

- [ ] GitHub Actions에서 테스트 실행 확인

## Phase 5: 문서화

- [ ] `README.md` 테스트 섹션 업데이트
- [ ] `CLAUDE.md` 테스트 관련 내용 업데이트
- [ ] `CHANGELOG.md`에 변경 사항 추가
