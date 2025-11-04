# Redis Token Storage - 구현 가이드

> **프로젝트**: Korea Investment Stock - Token Storage Enhancement
> **버전**: 1.2
> **작성일**: 2025-01-04
> **최종 수정**: 2025-01-04 - 테스트 전략 (fakeredis 사용), 가상환경 설정 추가

---

## 개요

파일 기반 토큰 저장소에 Redis 저장소를 추가하여 분산 환경과 멀티 프로세스 환경을 지원합니다.

**핵심 변경**:
- `TokenStorage` 추상 클래스 추가 (저장소 인터페이스)
- `FileTokenStorage` 클래스 (기존 구현 래핑)
- `RedisTokenStorage` 클래스 (신규, Redis 인증 지원)
- `KoreaInvestment` 클래스에 저장소 주입 기능 추가

---

## 1. 저장소 클래스 구현

### 1.1 TokenStorage 추상 클래스

```python
from abc import ABC, abstractmethod
from typing import Optional, Dict, Any

class TokenStorage(ABC):
    """토큰 저장소 추상 클래스"""

    @abstractmethod
    def save_token(self, token_data: Dict[str, Any]) -> bool:
        """토큰 저장"""
        pass

    @abstractmethod
    def load_token(self, api_key: str, api_secret: str) -> Optional[Dict[str, Any]]:
        """토큰 로드"""
        pass

    @abstractmethod
    def check_token_valid(self, api_key: str, api_secret: str) -> bool:
        """토큰 유효성 확인"""
        pass

    @abstractmethod
    def delete_token(self, api_key: str, api_secret: str) -> bool:
        """토큰 삭제"""
        pass
```

### 1.2 FileTokenStorage 클래스

```python
import pickle
from pathlib import Path
from datetime import datetime

class FileTokenStorage(TokenStorage):
    """파일 기반 토큰 저장소"""

    def __init__(self, file_path: Path = None):
        """
        Args:
            file_path: 토큰 파일 경로 (기본값: ~/.cache/kis/token.key)
        """
        self.token_file = file_path or Path("~/.cache/kis/token.key").expanduser()

    def save_token(self, token_data: Dict[str, Any]) -> bool:
        """토큰을 Pickle 파일로 저장"""
        try:
            self.token_file.parent.mkdir(parents=True, exist_ok=True)
            with self.token_file.open("wb") as f:
                pickle.dump(token_data, f)
            return True
        except Exception as e:
            logger.error(f"토큰 파일 저장 실패: {e}")
            return False

    def load_token(self, api_key: str, api_secret: str) -> Optional[Dict[str, Any]]:
        """파일에서 토큰 로드 및 검증"""
        if not self.check_token_valid(api_key, api_secret):
            return None

        try:
            with self.token_file.open("rb") as f:
                data = pickle.load(f)
            return data
        except Exception as e:
            logger.error(f"토큰 파일 로드 실패: {e}")
            return None

    def check_token_valid(self, api_key: str, api_secret: str) -> bool:
        """파일 기반 토큰 유효성 확인"""
        if not self.token_file.exists():
            return False

        try:
            with self.token_file.open("rb") as f:
                data = pickle.load(f)
        except Exception:
            return False

        # API Key/Secret 확인
        if (data.get('api_key') != api_key) or (data.get('api_secret') != api_secret):
            return False

        # 만료 시각 확인
        ts_now = int(datetime.now().timestamp())
        return ts_now < data.get('timestamp', 0)

    def delete_token(self, api_key: str, api_secret: str) -> bool:
        """토큰 파일 삭제"""
        try:
            if self.token_file.exists():
                self.token_file.unlink()
            return True
        except Exception as e:
            logger.error(f"토큰 파일 삭제 실패: {e}")
            return False
```

### 1.3 RedisTokenStorage 클래스

```python
import hashlib
from datetime import datetime
from typing import Optional, Dict, Any

class RedisTokenStorage(TokenStorage):
    """Redis 기반 토큰 저장소"""

    def __init__(self, redis_url: str = "redis://localhost:6379/0",
                 password: Optional[str] = None,
                 key_prefix: str = "korea_investment:token"):
        """
        Args:
            redis_url: Redis 연결 URL
            password: Redis 인증 비밀번호 (선택)
            key_prefix: Redis 키 프리픽스
        """
        import redis

        # Redis URL에 비밀번호가 없고 password 파라미터가 제공된 경우
        if password and ':' not in redis_url.split('@')[0].split('//')[1]:
            # redis://host:port/db → redis://:password@host:port/db
            parts = redis_url.split('//')
            protocol = parts[0]
            rest = parts[1]
            redis_url = f"{protocol}//:{password}@{rest}"

        self.redis_client = redis.from_url(redis_url, decode_responses=True)
        self.key_prefix = key_prefix

    def _get_redis_key(self, api_key: str) -> str:
        """Redis 키 생성 (API Key 해시 사용)"""
        key_hash = hashlib.sha256(api_key.encode()).hexdigest()[:12]
        return f"{self.key_prefix}:{key_hash}"

    def save_token(self, token_data: Dict[str, Any]) -> bool:
        """토큰을 Redis Hash로 저장 (TTL 자동 설정)"""
        try:
            redis_key = self._get_redis_key(token_data['api_key'])

            # Hash 저장
            self.redis_client.hset(
                redis_key,
                mapping={
                    'access_token': token_data['access_token'],
                    'access_token_token_expired': token_data['access_token_token_expired'],
                    'timestamp': str(token_data['timestamp']),
                    'api_key': token_data['api_key'],
                    'api_secret': token_data['api_secret']
                }
            )

            # TTL 설정 (만료 시각까지 남은 시간)
            ts_now = int(datetime.now().timestamp())
            ttl = token_data['timestamp'] - ts_now
            if ttl > 0:
                self.redis_client.expire(redis_key, ttl)

            return True
        except Exception as e:
            logger.error(f"Redis 토큰 저장 실패: {e}")
            return False

    def load_token(self, api_key: str, api_secret: str) -> Optional[Dict[str, Any]]:
        """Redis에서 토큰 로드"""
        if not self.check_token_valid(api_key, api_secret):
            return None

        try:
            redis_key = self._get_redis_key(api_key)
            data = self.redis_client.hgetall(redis_key)

            if not data:
                return None

            return {
                'access_token': data['access_token'],
                'access_token_token_expired': data['access_token_token_expired'],
                'timestamp': int(data['timestamp']),
                'api_key': data['api_key'],
                'api_secret': data['api_secret']
            }
        except Exception as e:
            logger.error(f"Redis 토큰 로드 실패: {e}")
            return None

    def check_token_valid(self, api_key: str, api_secret: str) -> bool:
        """Redis 토큰 유효성 확인"""
        try:
            redis_key = self._get_redis_key(api_key)

            if not self.redis_client.exists(redis_key):
                return False

            data = self.redis_client.hgetall(redis_key)

            # API Secret 확인
            if data.get('api_secret') != api_secret:
                return False

            # 만료 시각 확인
            ts_now = int(datetime.now().timestamp())
            timestamp = int(data.get('timestamp', 0))
            return ts_now < timestamp

        except Exception as e:
            logger.error(f"Redis 토큰 확인 실패: {e}")
            return False

    def delete_token(self, api_key: str, api_secret: str) -> bool:
        """Redis 토큰 삭제"""
        try:
            redis_key = self._get_redis_key(api_key)
            self.redis_client.delete(redis_key)
            return True
        except Exception as e:
            logger.error(f"Redis 토큰 삭제 실패: {e}")
            return False
```

---

## 2. KoreaInvestment 클래스 통합

### 2.1 `__init__()` 메서드 수정

```python
def __init__(self, api_key: str, api_secret: str, acc_no: str, mock: bool = False,
             token_storage: Optional[TokenStorage] = None):
    """
    Args:
        token_storage: 토큰 저장소 인스턴스 (None이면 환경 변수로 결정)
    """
    # ... (기존 코드)

    # 토큰 저장소 초기화
    if token_storage:
        self.token_storage = token_storage
    else:
        self.token_storage = self._create_token_storage()

    # access token
    self.access_token = None
    if self.token_storage.check_token_valid(self.api_key, self.api_secret):
        token_data = self.token_storage.load_token(self.api_key, self.api_secret)
        self.access_token = f'Bearer {token_data["access_token"]}'
    else:
        self.issue_access_token()
```

### 2.2 `_create_token_storage()` 메서드 추가

```python
def _create_token_storage(self) -> TokenStorage:
    """환경 변수 기반 토큰 저장소 생성"""
    storage_type = os.getenv("KOREA_INVESTMENT_TOKEN_STORAGE", "file").lower()

    if storage_type == "file":
        file_path = os.getenv("KOREA_INVESTMENT_TOKEN_FILE")
        if file_path:
            file_path = Path(file_path).expanduser()
        return FileTokenStorage(file_path)

    elif storage_type == "redis":
        redis_url = os.getenv("KOREA_INVESTMENT_REDIS_URL", "redis://localhost:6379/0")
        redis_password = os.getenv("KOREA_INVESTMENT_REDIS_PASSWORD")
        return RedisTokenStorage(redis_url, password=redis_password)

    else:
        raise ValueError(f"지원하지 않는 저장소 타입: {storage_type}")
```

### 2.3 기존 메서드 수정

**`issue_access_token()` 수정**:
```python
def issue_access_token(self):
    # ... API 호출 ...

    # BEFORE: 파일 직접 저장
    # with self.token_file.open("wb") as f:
    #     pickle.dump(resp_data, f)

    # AFTER: 저장소 추상화
    self.token_storage.save_token(resp_data)
```

**`check_access_token()` 수정**:
```python
def check_access_token(self) -> bool:
    # BEFORE: 파일 직접 확인
    # if not self.token_file.exists():
    #     return False
    # ...

    # AFTER: 저장소 추상화
    return self.token_storage.check_token_valid(self.api_key, self.api_secret)
```

**`load_access_token()` 수정**:
```python
def load_access_token(self):
    # BEFORE: 파일 직접 로드
    # with self.token_file.open("rb") as f:
    #     data = pickle.load(f)

    # AFTER: 저장소 추상화
    token_data = self.token_storage.load_token(self.api_key, self.api_secret)
    if token_data:
        self.access_token = f'Bearer {token_data["access_token"]}'
```

---

## 3. 환경 변수

| 환경 변수 | 설명 | 기본값 |
|-----------|------|--------|
| `KOREA_INVESTMENT_TOKEN_STORAGE` | 저장소 타입 ("file" or "redis") | `"file"` |
| `KOREA_INVESTMENT_REDIS_URL` | Redis 연결 URL | `"redis://localhost:6379/0"` |
| `KOREA_INVESTMENT_REDIS_PASSWORD` | Redis 비밀번호 | `None` |
| `KOREA_INVESTMENT_TOKEN_FILE` | 토큰 파일 경로 | `"~/.cache/kis/token.key"` |

**사용 예시**:
```bash
# File 저장소 (기본값)
export KOREA_INVESTMENT_TOKEN_STORAGE="file"

# Redis 저장소
export KOREA_INVESTMENT_TOKEN_STORAGE="redis"
export KOREA_INVESTMENT_REDIS_URL="redis://localhost:6379/0"

# Redis with 인증
export KOREA_INVESTMENT_TOKEN_STORAGE="redis"
export KOREA_INVESTMENT_REDIS_URL="redis://redis-server:6379/1"
export KOREA_INVESTMENT_REDIS_PASSWORD="mypassword"
```

---

## 4. Redis 키 스키마

```
KEY: korea_investment:token:{api_key_hash}
TYPE: Hash
TTL: 자동 (만료 시각까지 남은 시간)

FIELDS:
- access_token: "Bearer eyJ0eXAiOiJKV1Q..."
- access_token_token_expired: "2025-01-05 09:30:00"
- timestamp: "1736036400"
- api_key: "PSxxxxxxxxxx"
- api_secret: "xxxxxxxxxxxx"
```

---

## 5. 의존성 관리

**pyproject.toml 수정**:
```toml
[project]
name = "korea-investment-stock"
version = "0.6.1"
dependencies = [
    "requests",
    "pandas",
    "websockets",
    "pycryptodome",
    "crypto>=1.4.1",
]

[project.optional-dependencies]
redis = [
    "redis>=4.5.0",
]

dev = [
    "pytest>=7.0.0",
    "pytest-mock>=3.10.0",
    "fakeredis>=2.10.0",  # In-memory Redis (테스트용)
]
```

**설치**:
```bash
# 가상환경 생성 및 활성화
python -m venv .venv
source .venv/bin/activate  # macOS/Linux
# .venv\Scripts\activate  # Windows

# 기본 설치 (File 저장소만)
pip install korea-investment-stock

# Redis 지원 포함
pip install korea-investment-stock[redis]

# 개발 의존성 포함 (테스트용)
pip install -e ".[dev,redis]"
```

---

## 6. 테스트 전략

### 6.1 In-memory Redis 라이브러리 사용

Redis 테스트는 **`fakeredis`** 라이브러리를 사용하여 Docker 없이 순수 Python으로 수행합니다.

**`fakeredis` 특징**:
- Pure Python으로 구현된 Redis
- 외부 의존성 없음 (Docker 불필요)
- 실제 Redis 명령어 지원
- 메모리에서 동작 (빠른 테스트)
- TTL, Hash, String 등 대부분의 Redis 기능 지원

**pytest fixture 예시**:
```python
import pytest
import fakeredis

@pytest.fixture
def fake_redis():
    """In-memory Redis 인스턴스 (fakeredis)"""
    return fakeredis.FakeStrictRedis(decode_responses=True)

@pytest.fixture
def redis_storage(fake_redis, monkeypatch):
    """fakeredis를 사용하는 RedisTokenStorage"""
    def mock_from_url(*args, **kwargs):
        return fake_redis

    monkeypatch.setattr('redis.from_url', mock_from_url)
    return RedisTokenStorage("redis://localhost:6379/0")

def test_save_and_load(redis_storage):
    """토큰 저장 및 로드 테스트"""
    token_data = {
        'access_token': 'Bearer test',
        'access_token_token_expired': '2025-12-31 23:59:59',
        'timestamp': int((datetime.now() + timedelta(days=1)).timestamp()),
        'api_key': 'test_key',
        'api_secret': 'test_secret'
    }

    assert redis_storage.save_token(token_data)
    loaded = redis_storage.load_token('test_key', 'test_secret')
    assert loaded['access_token'] == 'Bearer test'

def test_ttl_auto_expire(redis_storage, fake_redis):
    """TTL 자동 만료 테스트"""
    token_data = {
        'access_token': 'Bearer ttl_test',
        'timestamp': int((datetime.now() + timedelta(seconds=10)).timestamp()),
        'api_key': 'ttl_key',
        'api_secret': 'ttl_secret'
    }

    redis_storage.save_token(token_data)

    # TTL 확인
    redis_key = redis_storage._get_redis_key('ttl_key')
    ttl = fake_redis.ttl(redis_key)
    assert 5 < ttl <= 10  # TTL이 설정되어 있음
```

**테스트 실행**:
```bash
# 가상환경 활성화
source .venv/bin/activate

# 테스트 실행 (Docker 불필요)
pytest korea_investment_stock/tests/test_token_storage.py -v

# 특정 테스트만 실행
pytest korea_investment_stock/tests/test_token_storage.py::TestRedisTokenStorage -v
```

### 6.2 테스트 파일 구조

```
korea_investment_stock/tests/
├── test_token_storage.py (신규)
│   ├── TestFileTokenStorage
│   ├── TestRedisTokenStorage (fakeredis - in-memory)
│   └── TestTokenStorageIntegration
└── test_korea_investment_stock.py (수정)
    └── test_token_storage_integration
```

---

## 7. 사용 예시

**사전 준비**: 모든 예시는 가상환경 활성화 후 실행해야 합니다.

```bash
# 가상환경 생성 (최초 1회)
python -m venv .venv

# 가상환경 활성화
source .venv/bin/activate  # macOS/Linux
# .venv\Scripts\activate  # Windows

# 패키지 설치
pip install korea-investment-stock[redis]
```

### 7.1 기본 사용 (변경 없음)
```python
from korea_investment_stock import KoreaInvestment

# 가상환경 활성화 상태에서 실행
broker = KoreaInvestment(api_key, secret, acc_no)
# 내부적으로 File 저장소 사용 (~/.cache/kis/token.key)
```

### 7.2 Redis 사용
```python
import os
from korea_investment_stock import KoreaInvestment

# 환경 변수 설정
os.environ["KOREA_INVESTMENT_TOKEN_STORAGE"] = "redis"
os.environ["KOREA_INVESTMENT_REDIS_URL"] = "redis://localhost:6379/0"

# 가상환경 활성화 상태에서 실행
broker = KoreaInvestment(api_key, secret, acc_no)
# Redis에서 토큰 로드
```

### 7.3 커스텀 저장소
```python
from korea_investment_stock import KoreaInvestment, RedisTokenStorage

# 가상환경 활성화 상태에서 실행
storage = RedisTokenStorage("redis://custom:6379/1", password="secret")
broker = KoreaInvestment(api_key, secret, acc_no, token_storage=storage)
```

---

## 문서 히스토리

| 버전 | 날짜 | 변경사항 |
|------|------|---------|
| 1.0 | 2025-01-04 | 초안 작성 - 구현 가이드 및 테스트 전략 |
| 1.1 | 2025-01-04 | 테스트 전략 변경 (FakeRedis → Docker) |
| 1.2 | 2025-01-04 | 테스트 전략 재변경 (Docker → fakeredis), 가상환경 설정 추가 |

---

**작성일**: 2025-01-04
**최종 수정**: 2025-01-04
