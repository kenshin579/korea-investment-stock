# 🚀 Korea Investment Stock

[![PyPI version](https://badge.fury.io/py/korea-investment-stock.svg)](https://badge.fury.io/py/korea-investment-stock)
[![Python Versions](https://img.shields.io/pypi/pyversions/korea-investment-stock.svg)](https://pypi.org/project/korea-investment-stock/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

한국투자증권 OpenAPI를 위한 파이썬 라이브러리입니다. Rate Limiting, 자동 재시도, 배치 처리 등 프로덕션 환경에 필요한 기능들을 포함하고 있습니다.

## 🌟 주요 특징

### ✨ 핵심 기능
- **한국투자증권 API 시세 조회**: 국내/해외 주식 시세 및 종목 정보 조회
- **공모주 청약 일정 조회**: 공모주 정보, 청약 일정, 상태 확인
- **자동 Rate Limiting**: API 호출 제한(초당 20회)을 자동으로 관리
- **스마트 재시도**: Exponential Backoff와 Circuit Breaker 패턴 구현
- **배치 처리**: 대량 데이터 조회를 위한 최적화된 배치 처리
- **실시간 모니터링**: 상세한 통계 및 성능 추적
- **TTL 캐싱**: 자동 만료되는 캐시로 API 호출 최소화

### 🛡️ 안정성
- **에러율 0%**: 프로덕션 환경에서 검증된 안정성
- **자동 에러 복구**: 일시적 오류 자동 처리
- **Thread-Safe**: 멀티스레드 환경 지원

## 📦 설치

```bash
pip install korea-investment-stock
```

### 요구사항
- Python 3.9 이상
- 한국투자증권 API 계정

## 🚀 빠른 시작

### 기본 사용법

```python
from korea_investment_stock import KoreaInvestment

# API 인증 정보
api_key = "YOUR_API_KEY"
api_secret = "YOUR_API_SECRET"
account_no = "12345678-01"

# 클라이언트 생성
broker = KoreaInvestment(
    api_key=api_key,
    api_secret=api_secret,
    acc_no=account_no,
    mock=False  # 실거래: False, 모의투자: True
)

# 현재가 조회
price_info = broker.fetch_price("005930")  # 삼성전자
print(f"현재가: {price_info['output']['stck_prpr']}원")

# 종목 정보 조회
stock_info = broker.fetch_stock_info_list([("005930", "KR")])
print(stock_info)
```

### Context Manager 사용 (권장)

```python
from korea_investment_stock import KoreaInvestment

with KoreaInvestment(api_key, api_secret, account_no) as broker:
    # 자동으로 리소스 정리
    price = broker.fetch_price("005930")
    # ... 작업 수행
# 자동으로 broker.shutdown() 호출됨
```

## 📊 주요 기능

### 1. 주식 정보 조회

```python
# 국내 주식 현재가
price = broker.fetch_domestic_price("J", "005930")

# ETF 현재가
etf_price = broker.fetch_etf_domestic_price("J", "069500")  # KODEX 200

# 종목 정보 조회
stock_info = broker.fetch_stock_info_list([("005930", "KR")])

# 여러 종목 동시 조회
stock_list = [("005930", "KR"), ("000660", "KR"), ("035720", "KR")]
prices = broker.fetch_price_list(stock_list)

# 코스피/코스닥 전체 종목 조회
kospi_symbols = broker.fetch_kospi_symbols()
kosdaq_symbols = broker.fetch_kosdaq_symbols()
```

### 2. 공모주 청약 일정 조회

```python
# 전체 공모주 일정 조회 (오늘부터 30일)
ipo_schedule = broker.fetch_ipo_schedule()
if ipo_schedule['rt_cd'] == '0':
    for ipo in ipo_schedule['output1']:
        print(f"{ipo['isin_name']} - 청약기간: {ipo['subscr_dt']}")
        print(f"공모가: {broker.format_number(ipo['fix_subscr_pri'])}원")
        print(f"주간사: {ipo['lead_mgr']}")

# 특정 기간 공모주 조회
ipos = broker.fetch_ipo_schedule(
    from_date="20240101",
    to_date="20240131"
)

# 특정 종목 공모주 정보 조회
ipo_info = broker.fetch_ipo_schedule(symbol="123456")

# 공모주 상태 확인 (예정/진행중/마감)
for ipo in ipo_schedule['output1']:
    status = broker.get_ipo_status(ipo['subscr_dt'])
    d_day = broker.calculate_ipo_d_day(ipo['subscr_dt'])
    print(f"{ipo['isin_name']}: {status}, D-{d_day}")

# 주의: 모의투자는 지원하지 않습니다
# 예탁원 제공 정보이므로 참고용으로만 사용하세요
```

### 3. 배치 처리 (대량 데이터)

```python
# 100개 종목 조회
large_stock_list = [(f"{i:06d}", "KR") for i in range(1, 101)]

# 고정 크기 배치 처리
results = broker.fetch_price_list_with_batch(
    large_stock_list,
    batch_size=20,       # 20개씩 처리
    batch_delay=1.0,     # 배치 간 1초 대기
    progress_interval=10 # 진행상황 출력
)

# 동적 배치 처리 (에러율에 따라 자동 조정)
results = broker.fetch_price_list_with_dynamic_batch(large_stock_list)
```

### 4. Rate Limiting 관리

```python
# Rate Limiter 통계 확인
broker.rate_limiter.print_stats()

# 통계 데이터 가져오기
stats = broker.rate_limiter.get_stats()
print(f"총 호출: {stats['total_calls']}")
print(f"에러율: {stats['error_rate']:.1%}")
print(f"평균 대기시간: {stats['avg_wait_time']:.3f}초")

# 통계 저장
broker.rate_limiter.save_stats()
```

### 5. 캐시 관리

```python
# 캐시 상태 확인
cache_stats = broker.get_cache_stats()
print(f"캐시 적중률: {cache_stats['hit_rate']:.1%}")
print(f"메모리 사용량: {cache_stats['memory_usage']:.1f}MB")

# 캐시 비우기
broker.clear_cache()  # 전체 캐시 삭제
broker.clear_cache("fetch_domestic_price:J:005930")  # 특정 항목만 삭제

# 자주 사용하는 종목 미리 캐싱
popular_stocks = ["005930", "000660", "035720", "051910", "005380"]
broker.preload_cache(popular_stocks, market="KR")
```

### 6. 고급 기능

```python
# Circuit Breaker 상태 확인
from korea_investment_stock.rate_limiting import get_backoff_strategy

backoff = get_backoff_strategy()
state = backoff.get_stats()['state']  # CLOSED, OPEN, HALF_OPEN

# 에러 복구 시스템
from korea_investment_stock.error_handling import get_error_recovery_system

recovery = get_error_recovery_system()
summary = recovery.get_error_summary(hours=1)
print(f"최근 1시간 에러: {summary['total_errors']}건")

# 통계 매니저
from korea_investment_stock.monitoring import get_stats_manager

stats_mgr = get_stats_manager()
stats_mgr.save_all_stats()  # 모든 통계 저장
```

### 7. 모니터링 및 시각화

```python
# 실시간 모니터링 대시보드 생성
dashboard = broker.create_monitoring_dashboard()
broker.save_monitoring_dashboard("monitoring.html")

# 시스템 헬스 체크
health_chart = broker.get_system_health_chart()

# API 사용량 차트
usage_chart = broker.get_api_usage_chart(hours=24)

# 통계 리포트 생성
report_files = broker.create_stats_report("weekly_report")
```

## 🚧 개발 중인 기능

다음 기능들은 향후 버전에서 추가될 예정입니다:

- **주문 기능**: 시장가/지정가 주문, 주문 취소
- **잔고 조회**: 보유 종목 및 예수금 조회
- **차트 데이터**: 일봉, 분봉 등 OHLCV 데이터
- **해외 주식**: 미국 주식 시세 조회 기능 확대

## 📁 프로젝트 구조

```
korea_investment_stock/
├── __init__.py
├── korea_investment_stock.py      # 메인 클래스
├── rate_limiting/                 # Rate Limiting 모듈
│   ├── enhanced_rate_limiter.py
│   ├── enhanced_backoff_strategy.py
│   └── enhanced_retry_decorator.py
├── error_handling/                # 에러 처리 모듈
│   └── error_recovery_system.py
├── batch_processing/              # 배치 처리 모듈
│   └── dynamic_batch_controller.py
├── caching/                       # 캐싱 모듈
│   ├── ttl_cache.py
│   └── market_hours.py
├── monitoring/                    # 모니터링 및 통계
│   └── stats_manager.py
├── visualization/                 # 시각화 모듈
│   ├── charts.py
│   ├── dashboard.py
│   └── plotly_visualizer.py
└── utils/                        # 유틸리티 함수
```

## 🔧 환경 설정

### 환경 변수

```bash
# Rate Limiting 설정
export RATE_LIMIT_MAX_CALLS=15        # 최대 호출 수 (기본: 15)
export RATE_LIMIT_SAFETY_MARGIN=0.8   # 안전 마진 (기본: 0.8)

# Backoff 전략 설정
export BACKOFF_BASE_DELAY=1.0         # 기본 대기 시간
export BACKOFF_MAX_DELAY=60.0         # 최대 대기 시간
export CIRCUIT_FAILURE_THRESHOLD=5    # Circuit Breaker 임계값

# 캐시 설정
export CACHE_DEFAULT_TTL=300          # 기본 캐시 TTL (초)
export CACHE_MAX_SIZE=10000          # 최대 캐시 크기
```

## 📈 성능 지표

- **처리량**: 10-12 TPS (안정적)
- **에러율**: < 0.1%
- **100종목 조회**: ~8.5초
- **메모리 사용**: < 100MB
- **CPU 사용률**: < 5%
- **캐시 적중률**: > 80% (일반적인 사용 패턴)

## 🤝 기여하기

프로젝트에 기여를 환영합니다! 

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### 개발 환경 설정

```bash
# 저장소 클론
git clone https://github.com/softyoungha/korea-investment-stock.git
cd korea-investment-stock

# 가상환경 생성 및 활성화
python -m venv venv
source venv/bin/activate  # Windows: venv\Scripts\activate

# 개발 의존성 설치
pip install -e ".[dev]"

# 테스트 실행
pytest
```

## 📚 문서

- [상세 문서](https://wikidocs.net/book/7845)
- [API 레퍼런스](docs/api_reference.md)
- [예제 코드](examples/)
- [변경 이력](CHANGELOG.md)

## ⚠️ 주의사항

- 이 라이브러리는 한국투자증권 OpenAPI의 공식 라이브러리가 아닙니다
- 실거래 사용 시 충분한 테스트를 거쳐 사용하세요
- API 호출 제한을 준수하여 사용하세요
- 현재는 시세 조회 기능만 지원합니다

## 📄 라이선스

이 프로젝트는 MIT 라이선스 하에 배포됩니다. 자세한 내용은 [LICENSE](LICENSE) 파일을 참조하세요.

## 🙏 감사의 글

- 한국투자증권 OpenAPI 팀
- 모든 기여자들
- 이슈 리포트와 피드백을 주신 사용자분들

## 📞 지원

- **이슈 트래커**: [GitHub Issues](https://github.com/softyoungha/korea-investment-stock/issues)
- **토론**: [GitHub Discussions](https://github.com/softyoungha/korea-investment-stock/discussions)

---

<p align="center">
  Made with ❤️ by the Korea Investment Stock community
</p>
