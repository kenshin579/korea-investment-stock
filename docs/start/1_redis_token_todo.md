# Redis Token Storage - TODO

> **프로젝트**: Korea Investment Stock - Token Storage Enhancement
> **버전**: 1.2
> **작성일**: 2025-01-04
> **최종 수정**: 2025-01-04 - 테스트 전략 (fakeredis), 가상환경 설정 추가
> **예상 소요 시간**: 2-3일

---

## Phase 1: 코어 구현 (7시간) ✅ COMPLETED

### 1.1 저장소 클래스 설계 ✅

- [x] `TokenStorage` 추상 클래스 구현 (1h)
  - [x] `save_token()` 추상 메서드
  - [x] `load_token()` 추상 메서드
  - [x] `check_token_valid()` 추상 메서드
  - [x] `delete_token()` 추상 메서드

### 1.2 FileTokenStorage 구현 ✅

- [x] 기존 파일 저장 로직 래핑
  - [x] `save_token()` - Pickle 저장
  - [x] `load_token()` - Pickle 로드
  - [x] `check_token_valid()` - 파일 존재 및 만료 확인
  - [x] `delete_token()` - 파일 삭제
- [x] 기본 경로 설정: `~/.cache/kis/token.key`
- [x] 디렉토리 자동 생성 로직

### 1.3 RedisTokenStorage 구현 ✅

- [x] Redis 클라이언트 초기화
  - [x] `redis_url` 파라미터 처리
  - [x] `password` 파라미터 처리 (URL에 주입)
  - [x] Redis 연결 설정
- [x] Redis 키 생성 로직
  - [x] `_get_redis_key()` 메서드 (SHA-256 해시)
- [x] 토큰 저장/로드 구현
  - [x] `save_token()` - Hash 저장 + TTL 설정
  - [x] `load_token()` - Hash 로드 + 타입 변환
  - [x] `check_token_valid()` - 존재 확인 + 만료 확인
  - [x] `delete_token()` - Redis 키 삭제

### 1.4 단위 테스트 ✅

- [x] `test_token_storage.py` 파일 생성
- [x] `TestFileTokenStorage` 클래스
  - [x] `test_save_and_load()`
  - [x] `test_expired_token()`
  - [x] `test_wrong_credentials()`
- [x] `TestRedisTokenStorage` 클래스 (fakeredis)
  - [x] `fakeredis` fixture 설정
  - [x] `test_save_and_load()`
  - [x] `test_redis_with_password()`
  - [x] `test_ttl_auto_expire()`
  - [x] `test_concurrent_access()`

---

## Phase 2: 통합 및 테스트 (7시간)

### 2.1 KoreaInvestment 클래스 수정 (3h)

- [ ] `__init__()` 메서드 수정 (1h)
  - [ ] `token_storage` 파라미터 추가
  - [ ] 토큰 저장소 초기화 로직
  - [ ] 기존 `self.token_file` 제거
  - [ ] 토큰 로드 로직 변경

- [ ] `_create_token_storage()` 메서드 구현 (1h)
  - [ ] 환경 변수 읽기 (`KOREA_INVESTMENT_TOKEN_STORAGE`)
  - [ ] `"file"` → `FileTokenStorage` 생성
  - [ ] `"redis"` → `RedisTokenStorage` 생성
  - [ ] 환경 변수 검증 및 에러 처리

- [ ] 기존 메서드 수정 (1h)
  - [ ] `issue_access_token()` - `save_token()` 호출로 변경
  - [ ] `check_access_token()` - `check_token_valid()` 호출로 변경
  - [ ] `load_access_token()` - `load_token()` 호출로 변경

### 2.2 통합 테스트 (2h)

- [ ] `test_korea_investment_stock.py` 수정
  - [ ] `test_file_storage_default()` - 기본 파일 저장소 테스트
  - [ ] `test_redis_storage_via_env()` - 환경 변수로 Redis 사용
  - [ ] `test_custom_storage()` - 커스텀 저장소 주입
  - [ ] 기존 테스트 호환성 확인

### 2.3 예제 코드 작성 (1.5h)

- [ ] `examples/redis_token_example.py` 생성
  - [ ] File 저장소 예제
  - [ ] Redis 저장소 예제 (인증 없음)
  - [ ] Redis 저장소 예제 (인증 포함)
  - [ ] 커스텀 저장소 예제

### 2.4 의존성 업데이트 (0.5h)

- [ ] `pyproject.toml` 수정
  - [ ] `version` 업데이트: `0.6.1`
  - [ ] `[project.optional-dependencies]` 섹션 추가
  - [ ] `redis = ["redis>=4.5.0"]` 추가
  - [ ] `dev`에 `fakeredis>=2.10.0` 추가

---

## Phase 3: 문서화 및 배포 (5시간)

### 3.1 README.md 업데이트 (1h)

- [ ] Redis 저장소 섹션 추가
  - [ ] 설치 방법 (`pip install korea-investment-stock[redis]`)
  - [ ] 환경 변수 설정 가이드
  - [ ] 사용 예시 (File vs Redis)
  - [ ] 분산 환경 사용 예시

### 3.2 CLAUDE.md 업데이트 (1h)

- [ ] 아키텍처 섹션 수정
  - [ ] 저장소 클래스 구조 설명
  - [ ] 환경 변수 목록 추가
- [ ] 개발 패턴 섹션
  - [ ] 새로운 저장소 구현 방법
  - [ ] FakeRedis 테스트 패턴

### 3.3 CHANGELOG.md 작성 (0.5h)

- [ ] v0.6.1 섹션 추가
  - [ ] Added: Redis token storage support
  - [ ] Added: Token storage abstraction (TokenStorage, FileTokenStorage, RedisTokenStorage)
  - [ ] Added: Environment variables for storage configuration
  - [ ] Changed: Token storage is now pluggable

### 3.4 환경 변수 가이드 작성 (0.5h)

- [ ] `docs/environment_variables.md` 생성 (선택)
  - [ ] 환경 변수 목록 및 설명
  - [ ] 보안 권장사항 (Redis 비밀번호 관리)
  - [ ] Docker/Kubernetes 환경 설정 예시

### 3.5 실제 Redis 환경 테스트 (1h)

- [ ] 가상환경 설정
  - [ ] `python -m venv .venv`
  - [ ] `source .venv/bin/activate`
  - [ ] `pip install -e ".[dev,redis]"`
- [ ] 단위 테스트 실행 (fakeredis)
  - [ ] `pytest korea_investment_stock/tests/test_token_storage.py -v`
- [ ] 실제 Redis 서버 테스트 (선택)
  - [ ] Docker로 Redis 실행: `docker run -d -p 6379:6379 redis:7-alpine`
  - [ ] 환경 변수 설정 및 수동 테스트
  - [ ] Redis CLI로 데이터 확인
  - [ ] 테스트 완료 후 컨테이너 정리

### 3.6 PyPI 배포 (1h)

- [ ] 빌드 및 검증
  - [ ] `python -m build`
  - [ ] `twine check dist/*`
- [ ] 배포
  - [ ] `twine upload dist/*`
- [ ] 설치 확인
  - [ ] `pip install korea-investment-stock[redis]==0.6.1`
  - [ ] 간단한 동작 테스트

---

## 체크리스트 요약

### 필수 작업
- [ ] Phase 1 완료 (7시간) - 코어 구현 및 단위 테스트
- [ ] Phase 2 완료 (7시간) - 통합 및 예제 코드
- [ ] Phase 3 완료 (5시간) - 문서화 및 배포

### 검증 항목
- [ ] 모든 단위 테스트 통과 (`pytest`)
- [ ] 하위 호환성 확인 (기존 코드 동작)
- [ ] Redis 실제 환경 테스트 통과
- [ ] 문서화 완료 (README, CLAUDE, CHANGELOG)
- [ ] PyPI 배포 성공

---

## 문서 히스토리

| 버전 | 날짜 | 변경사항 |
|------|------|---------|
| 1.0 | 2025-01-04 | 초안 작성 - 단계별 TODO 체크리스트 |
| 1.1 | 2025-01-04 | 테스트 전략 변경 (→ Docker) |
| 1.2 | 2025-01-04 | 테스트 전략 재변경 (→ fakeredis), 가상환경 설정 추가 |

---

**작성일**: 2025-01-04
**최종 수정**: 2025-01-04
