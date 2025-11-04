# PRD: Redis Token Storage Integration

> **프로젝트**: Korea Investment Stock - Token Storage Enhancement
> **작성일**: 2025-01-04
> **버전**: 1.4
> **관련 이슈**: Token Storage Improvement

---

## 📚 관련 문서

- **[구현 가이드](1_redis_token_implementation.md)** - 코드 구조, 클래스 설계, 환경 변수, 테스트 전략
- **[TODO 체크리스트](1_redis_token_todo.md)** - 단계별 구현 작업 목록 및 일정

---

## 📋 Executive Summary

### 프로젝트 목표
현재 파일 기반(`token.key`) 토큰 저장 방식에 Redis 저장소를 추가하여, 분산 환경과 멀티 프로세스 환경에서의 토큰 관리를 개선합니다.

### 핵심 변경사항
- **추가**: Redis 기반 토큰 저장/조회 기능
- **유지**: 기존 파일 기반 저장 방식 (하위 호환성)
- **개선**: 설정 기반 저장소 선택 (file/redis)

### 기대효과
- 분산 환경에서 토큰 공유 가능 (여러 서버/프로세스)
- Redis TTL 기능으로 자동 만료 관리
- 파일 I/O 부하 감소
- 멀티 프로세스 환경에서 토큰 동기화 문제 해결

---

## 🎯 Background & Context

### 현재 토큰 저장 방식 분석

**1. 파일 기반 저장 (token.key)**

```python
# 현재 구현 (korea_investment_stock.py:207-320)

class KoreaInvestment:
    def __init__(self, ...):
        # 토큰 파일 경로 설정
        self.token_file = Path("~/.cache/kis/token.key").expanduser()

        # 토큰 로딩
        if self.check_access_token():
            self.load_access_token()
        else:
            self.issue_access_token()
```

**저장 위치**: `~/.cache/kis/token.key`

**저장 형식**: Python Pickle 직렬화

**저장 데이터**:
```python
{
    'access_token': 'Bearer eyJ0eXAiOiJKV1Q...',  # JWT 토큰
    'access_token_token_expired': '2025-01-05 09:30:00',  # 만료 시각 (KST)
    'timestamp': 1736036400,  # Unix epoch (만료 시각)
    'api_key': 'PSxxxxxxxxxx',  # API Key
    'api_secret': 'xxxxxxxxxxxx'  # API Secret
}
```

**토큰 발급 흐름**:
```
1. issue_access_token() 호출
   ↓
2. Korea Investment API에 POST 요청
   POST /oauth2/tokenP
   Body: {
     "grant_type": "client_credentials",
     "appkey": api_key,
     "appsecret": api_secret
   }
   ↓
3. API 응답 수신
   {
     "access_token": "eyJ0eXAiOiJKV1Q...",
     "access_token_token_expired": "2025-01-05 09:30:00",
     "expires_in": 86400
   }
   ↓
4. 타임존 변환 (Asia/Seoul)
   dt = datetime.strptime(resp_data['access_token_token_expired'], '%Y-%m-%d %H:%M:%S')
   dt = dt.replace(tzinfo=ZoneInfo('Asia/Seoul'))
   resp_data['timestamp'] = int(dt.timestamp())
   ↓
5. Pickle 직렬화 후 파일 저장
   with token_file.open("wb") as f:
       pickle.dump(resp_data, f)
```

**토큰 검증 흐름**:
```
1. check_access_token() 호출
   ↓
2. 파일 존재 확인
   if not token_file.exists():
       return False
   ↓
3. Pickle 역직렬화
   with token_file.open("rb") as f:
       data = pickle.load(f)
   ↓
4. API Key/Secret 일치 확인
   if (data['api_key'] != self.api_key) or
      (data['api_secret'] != self.api_secret):
       return False
   ↓
5. 만료 시각 확인
   ts_now = int(datetime.now().timestamp())
   return ts_now < data['timestamp']
```

**토큰 로딩 흐름**:
```
1. load_access_token() 호출
   ↓
2. Pickle 역직렬화
   with token_file.open("rb") as f:
       data = pickle.load(f)
   ↓
3. 메모리에 토큰 설정
   self.access_token = f'Bearer {data["access_token"]}'
```

---

### 현재 구현의 문제점

**1. 분산 환경 미지원 (심각도: HIGH)**
- 여러 서버에서 동일한 API 계정 사용 시 각자 토큰 발급 필요
- 토큰 발급 API 호출 증가 → 불필요한 부하
- 토큰 동기화 불가능

```python
# 문제 시나리오
# Server 1
broker1 = KoreaInvestment(api_key, secret, acc_no)  # 토큰 발급 1

# Server 2 (동시 실행)
broker2 = KoreaInvestment(api_key, secret, acc_no)  # 토큰 발급 2 (중복!)

# → 2번의 토큰 발급 API 호출 (불필요)
```

**2. 멀티 프로세스 동시성 문제 (심각도: MEDIUM)**
- 동일 서버에서 여러 프로세스 실행 시 파일 동시 쓰기 위험
- Race condition 가능성

```python
# 문제 시나리오
# Process 1
broker1.issue_access_token()  # token.key 쓰기 시작
                              # (아직 완료 안됨)

# Process 2 (거의 동시)
broker2.check_access_token()  # token.key 읽기 시도
                              # → 불완전한 데이터 읽기 가능
```

**3. 파일 I/O 오버헤드 (심각도: LOW)**
- 매 인스턴스 생성 시 파일 읽기 필요
- 고빈도 인스턴스 생성 시 I/O 부하

**4. 자동 만료 관리 부재 (심각도: LOW)**
- 만료된 토큰 파일이 디스크에 계속 남음
- 수동 삭제 필요

---

### 목표 상태 (Target State)

**Redis 통합 후**:
```python
# 1. 환경 변수로 저장소 선택
export KOREA_INVESTMENT_TOKEN_STORAGE="redis"  # or "file"
export KOREA_INVESTMENT_REDIS_URL="redis://localhost:6379/0"
export KOREA_INVESTMENT_REDIS_PASSWORD="your-password"  # 인증 필요 시

# 2. 코드 변경 없음 (투명한 전환)
broker = KoreaInvestment(api_key, secret, acc_no)

# 3. 내부적으로 Redis 사용
# - 분산 환경에서 토큰 공유
# - TTL 자동 만료
# - 동시성 안전
```

**저장소 비교**:

| 특성 | File 저장소 | Redis 저장소 |
|------|-------------|--------------|
| **분산 환경 지원** | ❌ 각 서버 독립 | ✅ 모든 서버 공유 |
| **멀티 프로세스** | ⚠️ Race condition | ✅ Atomic 연산 |
| **자동 만료** | ❌ 수동 삭제 | ✅ TTL 자동 삭제 |
| **성능** | 🐢 파일 I/O | 🚀 In-memory |
| **설정 복잡도** | ✅ 설정 불필요 | ⚠️ Redis 서버 필요 |
| **의존성** | ✅ 없음 | ⚠️ redis-py |

---

## 📝 Requirements Summary

### R1: 저장소 추상화

**구현 대상**:
- `TokenStorage` 추상 클래스 (저장소 인터페이스)
- `FileTokenStorage` 클래스 (기존 파일 저장 래핑)
- `RedisTokenStorage` 클래스 (Redis 저장, 인증 지원)

**핵심 메서드**:
- `save_token()` - 토큰 저장
- `load_token()` - 토큰 로드
- `check_token_valid()` - 유효성 확인
- `delete_token()` - 토큰 삭제

**상세 구현**: [구현 가이드](1_redis_token_implementation.md) 참조

---

### R2: KoreaInvestment 클래스 통합

**변경사항**:
- `__init__()` 메서드에 `token_storage` 파라미터 추가
- `_create_token_storage()` 메서드로 환경 변수 기반 저장소 생성
- 기존 메서드 (`issue_access_token`, `check_access_token`, `load_access_token`) 수정

**하위 호환성**: 기본 동작 유지 (파일 저장소)

---

### R3: 환경 변수

| 환경 변수 | 기본값 | 설명 |
|-----------|--------|------|
| `KOREA_INVESTMENT_TOKEN_STORAGE` | `"file"` | 저장소 타입 (`"file"` or `"redis"`) |
| `KOREA_INVESTMENT_REDIS_URL` | `"redis://localhost:6379/0"` | Redis 연결 URL |
| `KOREA_INVESTMENT_REDIS_PASSWORD` | `None` | Redis 인증 비밀번호 |
| `KOREA_INVESTMENT_TOKEN_FILE` | `"~/.cache/kis/token.key"` | 토큰 파일 경로 |

---

### R4: Redis 키 스키마

```
KEY: korea_investment:token:{api_key_hash}
TYPE: Hash
TTL: 자동 (만료 시각까지 남은 시간)

FIELDS:
- access_token
- access_token_token_expired
- timestamp
- api_key
- api_secret
```

---

### R5: 의존성 관리

**pyproject.toml**:
- `version`: `0.6.1`
- `optional-dependencies`: `redis = ["redis>=4.5.0"]`
- `dev`: `pytest`, `pytest-mock`, `fakeredis>=2.10.0`

**설치**:
```bash
# 가상환경 생성 및 활성화
python -m venv .venv
source .venv/bin/activate  # macOS/Linux

# Redis 지원 포함
pip install korea-investment-stock[redis]
```

**테스트 환경**:
- Redis 테스트는 `fakeredis` 라이브러리 사용 (Pure Python, Docker 불필요)
- In-memory에서 동작하여 빠른 테스트 가능
- 실제 Redis 명령어 대부분 지원

---

## ⚠️ Risk Assessment

### High Risk Areas

**1. Redis 연결 실패 처리 (심각도: HIGH)**
- Redis 서버 다운 시 애플리케이션 전체 중단 위험
- 네트워크 지연 시 토큰 로딩 타임아웃

**완화 전략**:
- Redis 연결 타임아웃 설정 (3초)
- 예외 처리 후 명확한 에러 메시지
- Health check 엔드포인트 추가
- 문서에 Redis 장애 시 대응 방안 명시

**2. 기존 사용자 영향 (심각도: MEDIUM)**
- 환경 변수 설정 실수 시 토큰 로딩 실패
- Redis 의존성 미설치 시 import 에러

**완화 전략**:
- 기본값은 "file" 유지 (변경 없음)
- Redis 사용 시 명시적 에러 메시지
- Optional dependency로 redis-py 분리

**3. Redis 인증 보안 (심각도: MEDIUM)**
- 비밀번호가 환경 변수에 노출
- URL에 비밀번호 포함 시 로그 노출 위험

**완화 전략**:
- 환경 변수 분리 옵션 제공 (KOREA_INVESTMENT_REDIS_PASSWORD)
- 문서에 보안 권장사항 명시
- Kubernetes Secrets, AWS Secrets Manager 사용 권장

---

## 🎯 Success Criteria

### 정량적 지표
1. **하위 호환성**: 기존 코드 100% 동작 (환경 변수 없이)
2. **테스트 커버리지**: 90% 이상 (token_storage.py)
3. **성능**: Redis 조회 < 10ms, File 조회 < 50ms
4. **동시성**: 100 스레드 동시 접근 시 에러 0건

### 정성적 지표
5. **투명성**: 사용자가 저장소 변경 시 코드 수정 불필요
6. **유연성**: 커스텀 저장소 구현 가능
7. **보안성**: Redis 인증 지원으로 프로덕션 환경 안전성 확보

---

## ✍️ Document History

| 버전 | 날짜 | 작성자 | 변경사항 |
|-----|------|--------|---------|
| 1.0 | 2025-01-04 | Claude Code | 초안 작성 - 현재 구현 분석 및 Redis 통합 설계 |
| 1.1 | 2025-01-04 | Claude Code | 파일 경로 변경 (.cache/kis/token.key), Redis 인증 추가, DualTokenStorage 제거 |
| 1.2 | 2025-01-04 | Claude Code | 구현/TODO 별도 파일 분리, PRD 간소화 |
| 1.3 | 2025-01-04 | Claude Code | 테스트 전략 변경 (→ Docker) |
| 1.4 | 2025-01-04 | Claude Code | 테스트 전략 재변경 (→ fakeredis), 가상환경 설정 추가 |

---

**작성**: Claude Code
**검토**: (To be reviewed)
**승인**: (To be approved)
