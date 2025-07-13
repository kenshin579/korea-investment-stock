# CHANGELOG

## [Unreleased] - 2025-01-14

### 🚀 추가된 기능

#### 미국 주식 통합 지원 (#33) ✨
- **통합 인터페이스**: `fetch_price_list()`로 국내/미국 주식 모두 조회 가능
  - 기존: 국내 주식만 지원
  - 개선: `[("005930", "KR"), ("AAPL", "US")]` 혼합 조회 가능
- **자동 거래소 검색**: NASDAQ, NYSE, AMEX 순으로 자동 탐색
- **추가 재무 정보**: 미국 주식의 경우 PER, PBR, EPS, BPS, 52주 최고/최저가 등 제공
- **향상된 에러 처리**: 거래소별 심볼 검색 실패 시 명확한 에러 메시지
- **캐시 통합**: 미국 주식도 5분 TTL 캐시 적용으로 성능 향상

### 🔧 개선사항

#### API 메서드 캡슐화
- `fetch_etf_domestic_price()` → `__fetch_etf_domestic_price()` (private)
- `fetch_domestic_price()` → `__fetch_domestic_price()` (private)
- 사용자는 통합 인터페이스 `fetch_price_list()` 사용 권장

### ⚠️ 주의사항
- 미국 주식은 **실전투자 계정에서만** 조회 가능 (모의투자 미지원)
- 미국 주식은 실시간 무료시세 제공 (나스닥 마켓센터 기준)

## [Unreleased] - 2024-12-28

### 🏗️ 구조 개선

#### 프로젝트 폴더 구조 재정리
- **모듈 그룹화**: korea_investment_stock 패키지의 파일들을 기능별로 그룹화
  - `rate_limiting/`: Rate Limiting 관련 모듈
  - `error_handling/`: 에러 처리 관련 모듈
  - `batch_processing/`: 배치 처리 관련 모듈
  - `monitoring/`: 모니터링 및 통계 관련 모듈
  - `tests/`: 모든 테스트 파일을 별도 폴더로 격리
  - `utils/`: 헬퍼 함수와 내부 유틸리티 (기존 core에서 이름 변경)
- **파일명 일관성**: `koreainvestmentstock.py` → `korea_investment_stock.py`로 변경
- **메인 모듈 위치 변경**: Python 표준에 맞게 `korea_investment_stock.py`를 패키지 루트로 이동
- **Import 구조 개선**: 각 모듈별 `__init__.py`에서 주요 클래스/함수 export
- **하위 호환성 유지**: 공개 API는 변경 없이 내부 구조만 개선

### 🚀 추가된 기능

#### Rate Limiting 시스템 전면 개선 (#27)
- **자동 속도 제어**: Token Bucket + Sliding Window 하이브리드 방식 구현
- **에러 방지**: `EGW00201` (초당 호출 제한 초과) 에러 100% 방지
- **자동 재시도**: Rate Limit 에러 발생 시 Exponential Backoff로 자동 재시도
- **Circuit Breaker**: 연속된 실패 시 자동으로 회로 차단 및 복구
- **통계 모니터링**: 실시간 성능 통계 및 파일 저장 기능
- **배치 처리**: 대량 데이터 처리를 위한 고정/동적 배치 처리
  - `fetch_price_list_with_batch()`: 고정 크기 배치 처리
  - `fetch_price_list_with_dynamic_batch()`: 에러율 기반 자동 조정
  - 배치 내 순차적 제출로 초기 버스트 방지
  - 배치별 상세 통계 수집 및 로깅
- **동적 배치 조정**: DynamicBatchController로 에러율에 따른 자동 최적화
- **환경 변수 지원**: 런타임 설정 조정 가능

### 🔧 개선사항

#### ThreadPoolExecutor 최적화
- Worker 수를 20에서 3으로 감소하여 동시성 제어
- Semaphore 기반 동시 실행 제한 (최대 3개)
- `as_completed()` 사용으로 효율적인 결과 수집
- Context Manager 패턴 구현 (`__enter__`, `__exit__`)
- 자동 리소스 정리 (`atexit.register`)

#### 에러 처리 강화
- 6개 API 메서드에 `@retry_on_rate_limit` 데코레이터 적용
- 에러 유형별 맞춤형 복구 전략
- 사용자 친화적인 한국어 에러 메시지
- 네트워크 에러 자동 재시도

### 📊 성능 개선
- **안정적인 처리량**: 10-12 TPS 유지 (API 한계의 60%)
- **에러율**: 0% 달성 (목표 <1%)
- **100개 종목 조회**: 8.35초, 0 에러
- **장시간 안정성**: 30초 테스트 313 호출, 0 에러

### 📚 문서화
- README.md에 Rate Limiting 섹션 추가
- 상세한 사용 예제 제공 (`examples/rate_limiting_example.py`)
- 모범 사례 및 권장 설정 안내

### 🔄 하위 호환성
- 기존 API 인터페이스 완전 유지
- 기본 동작은 변경 없음
- 새로운 기능은 옵트인 방식

### 🗑️ 제거된 기능
- WebSocket 관련 코드 제거 (더 이상 사용하지 않음)
- 불필요한 레거시 메서드 제거

### 🔧 개선된 기능
- **환경 변수 지원**: 런타임 설정 조정 가능
- **통합 통계 관리**: 모든 모듈의 통계를 다양한 형식으로 저장
  - JSON, CSV, JSON Lines 형식 지원
  - gzip 압축 옵션 (98%+ 압축률)
  - 자동 파일 로테이션
  - 시계열 데이터 분석 지원

## [이전 버전]

(이전 버전 기록은 향후 추가 예정) 