# Korea Investment Stock 폴더 구조 재정리

Date: 2024-12-28

## 변경 개요

korea_investment_stock 패키지의 파일들을 기능별로 그룹화하여 폴더 구조를 재정리했습니다.

## 이전 구조 (평면적)

```
korea_investment_stock/
├── __init__.py
├── korea_investment_stock.py
├── enhanced_rate_limiter.py
├── enhanced_retry_decorator.py
├── enhanced_backoff_strategy.py
├── error_recovery_system.py
├── dynamic_batch_controller.py
├── stats_manager.py
├── test_*.py (13개 테스트 파일)
├── legacy/
└── logs/
```

## 새로운 구조 (기능별 그룹화)

```
korea_investment_stock/
├── __init__.py                    # 패키지 메인 엔트리 포인트
├── korea_investment_stock.py      # 메인 클래스 (패키지 루트로 이동)
├── utils/                         # 헬퍼 함수 및 내부 유틸리티
│   └── __init__.py
├── rate_limiting/                 # Rate Limiting 관련
│   ├── __init__.py
│   ├── enhanced_rate_limiter.py
│   ├── enhanced_retry_decorator.py
│   └── enhanced_backoff_strategy.py
├── error_handling/                # 에러 처리 관련
│   ├── __init__.py
│   └── error_recovery_system.py
├── batch_processing/              # 배치 처리 관련
│   ├── __init__.py
│   └── dynamic_batch_controller.py
├── monitoring/                    # 모니터링 및 통계
│   ├── __init__.py
│   └── stats_manager.py
├── tests/                         # 모든 테스트 파일
│   ├── __init__.py
│   ├── test_batch_processing.py
│   ├── test_enhanced_backoff.py
│   ├── test_enhanced_integration.py
│   ├── test_error_handling.py
│   ├── test_error_recovery.py
│   ├── test_integration.py
│   ├── test_korea_investment_stock.py
│   ├── test_load.py
│   ├── test_rate_limit_error_detection.py
│   ├── test_rate_limit_simulation.py
│   ├── test_rate_limiter.py
│   ├── test_stats_save.py
│   └── test_threadpool_improvement.py
├── legacy/                        # 레거시 코드
│   └── rate_limiter_v1.py
└── logs/                          # 로그 파일
    └── rate_limiter_stats/
```

## 변경 사항

### 1. 파일 이름 변경
- `koreainvestmentstock.py` → `korea_investment_stock.py` (파일명 일관성)
- `test_koreainvestmentstock.py` → `test_korea_investment_stock.py`

### 2. 폴더별 역할
- **korea_investment_stock.py**: 메인 클래스 (KoreaInvestment) - 패키지 루트에 위치
- **utils/**: 헬퍼 함수와 내부 유틸리티들을 위한 공간 (현재 비어있음)
- **rate_limiting/**: API 호출 제한 관리 관련 모듈들
- **error_handling/**: 에러 처리 및 복구 시스템
- **batch_processing/**: 대량 요청 처리를 위한 배치 처리 기능
- **monitoring/**: 통계 수집 및 모니터링 기능
- **tests/**: 모든 단위 테스트 및 통합 테스트

### 3. Import 경로 업데이트
모든 파일의 import 경로를 새로운 폴더 구조에 맞게 업데이트했습니다.

#### 예시: 이전
```python
from .enhanced_rate_limiter import EnhancedRateLimiter
```

#### 예시: 이후
```python
from ..rate_limiting import EnhancedRateLimiter
```

### 4. 패키지 __init__.py 업데이트
메인 __init__.py 파일을 업데이트하여 각 모듈의 주요 클래스와 함수를 export합니다:

```python
# Core imports
from .core.korea_investment_stock import KoreaInvestment, MARKET_CODE_MAP, EXCHANGE_CODE_MAP, API_RETURN_CODE

# Rate limiting imports
from .rate_limiting import EnhancedRateLimiter, retry_on_rate_limit, retry_on_network_error, get_backoff_strategy

# Error handling imports
from .error_handling import ErrorRecoverySystem, get_error_recovery_system

# Batch processing imports
from .batch_processing import DynamicBatchController

# Monitoring imports
from .monitoring import StatsManager, get_stats_manager
```

## 장점

1. **모듈화**: 기능별로 명확하게 분리되어 코드 관리가 용이
2. **확장성**: 새로운 기능 추가 시 적절한 폴더에 배치 가능
3. **가독성**: 프로젝트 구조를 한눈에 파악 가능
4. **테스트 격리**: 테스트 파일들이 별도 폴더에 정리되어 프로덕션 코드와 분리

## 마이그레이션 가이드

기존 코드에서 import를 사용하는 경우:
```python
# 이전
from korea_investment_stock import KoreaInvestment

# 이후 (변경 없음, 패키지 __init__.py에서 export하므로)
from korea_investment_stock import KoreaInvestment
```

내부 모듈을 직접 import하는 경우:
```python
# 이전
from korea_investment_stock.enhanced_rate_limiter import EnhancedRateLimiter

# 이후
from korea_investment_stock.rate_limiting import EnhancedRateLimiter
```

## 호환성
- 패키지의 공개 API는 변경되지 않았으므로 기존 사용자 코드와 완전히 호환됩니다.
- 내부 모듈을 직접 import하는 경우에만 경로 수정이 필요합니다. 