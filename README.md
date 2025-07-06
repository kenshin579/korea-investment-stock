# Korea Investment Stock

대한민국 증권사의 Rest API 기반의 Open API에 대한 통합 파이썬 레퍼 모듈입니다. 
통합 모듈이라 칵테일 이름인 모히토를 프로젝트명으로 사용하고 있으며, 돈 벌어서 몰디브가서 모히토 한 잔 하자는 의미도 있습니다. 


# 설치 

```sh
$ pip install korea-investment-stock
```

# 지원 API 

| 카테고리 | 기능 | 함수 |
|--------|-----|-----|
| OAuth 인증 | Hasheky | `issue_hashkey()` |
| OAuth 인증 | 접근토근발급(P) | `issue_access_token()` |
| OAuth 인증 | 접근토근폐기(P) | 미지원 |
| 국내주식주문 | 주식주문(현금) |  |
| 국내주식주문 | 주식잔고조회 | `fetch_balance()` |

# 사용법
## 한국투자증권

https://wikidocs.net/book/7845  

현재가 조회

```py
import korea_investment_stock
import pprint

key = "발급받은 API KEY"
secret = "발급받은 API SECRET"
acc_no = "12345678-01"

broker = korea_investment_stock.KoreaInvestment(api_key=key, api_secret=secret, acc_no=acc_no)
resp = broker.fetch_price("005930")
pprint.pprint(resp)

```

일봉 데이터 조회

```py
import korea_investment_stock
import pprint

key = "발급받은 API KEY"
secret = "발급받은 API SECRET"
acc_no = "12345678-01"

broker = korea_investment_stock.KoreaInvestment(api_key=key, api_secret=secret, acc_no=acc_no)
resp = broker.fetch_daily_price("005930")
pprint.pprint(resp)
```

잔고 조회 

```py
resp = broker.fetch_balance()
pprint.pprint(resp)
```

주문 

```py
resp = broker.create_market_buy_order("005930", 10) # 삼성전자, 10주, 시장가
pprint.pprint(resp)
```

```yaml
{
 'rt_cd': '0',
 'msg_cd': 'APBK0013',
 'msg1': '주문 전송 완료 되었습니다.',
 'output': {'KRX_FWDG_ORD_ORGNO': '91252',
  'ODNO': '0000117057',
  'ORD_TMD': '121052'}
}
```

주문 취소

```py
resp = broker.cancel_order("91252", "0000117057", "00", 60000, 5, "Y") # KRX_FWDG_ORD_ORGNO, ODNO, 지정가 주문, 가격, 수량, 모두 
print(resp)
```

미국주식 주문

```py
broker = KoreaInvestment(key, secret, acc_no=acc_no, exchange="NASD")
resp = broker.create_limit_buy_order("TQQQ", 35, 1)
print(resp)
```

웹소켓

```py
import pprint
import korea_investment_stock

with open("../../koreainvestment.key", encoding="utf-8") as f:
    lines = f.readlines()
key = lines[0].strip()
secret = lines[1].strip()

if __name__ == "__main__":
    broker_ws = korea_investment_stock.KoreaInvestmentWS(key, secret, ["H0STCNT0", "H0STASP0"], ["005930", "000660"],
                                                         user_id="idjhh82")
    broker_ws.start()
    while True:
        data_ = broker_ws.get()
        if data_[0] == '체결':
            print(data_[1])
        elif data_[0] == '호가':
            print(data_[1])
        elif data_[0] == '체잔':
            print(data_[1])
```        

# Rate Limiting (속도 제한 관리)

한국투자증권 API는 초당 20회 호출 제한이 있습니다. 이 라이브러리는 자동으로 Rate Limit을 관리하여 `EGW00201` 에러를 방지합니다.

## 특징

- **자동 속도 제어**: Token Bucket + Sliding Window 하이브리드 방식
- **보수적 설정**: API 한계의 60%(12 TPS)만 사용하여 안정성 확보
- **자동 재시도**: Rate Limit 에러 발생 시 자동 재시도
- **통계 및 모니터링**: 실시간 성능 추적

## 기본 사용법

Rate Limiting은 자동으로 작동합니다:

```python
# 일반 사용 - Rate Limiting이 자동으로 적용됨
broker = korea_investment_stock.KoreaInvestment(api_key=key, api_secret=secret, acc_no=acc_no)

# 여러 종목 조회 시에도 자동으로 속도 제어
stock_list = [("005930", "KR"), ("000660", "KR"), ("035720", "KR")]
results = broker.fetch_price_list(stock_list)
```

## 통계 확인

### 실시간 통계 출력
```python
# 통계 출력
broker.rate_limiter.print_stats()

# 통계 데이터 가져오기
stats = broker.rate_limiter.get_stats()
print(f"총 호출 수: {stats['total_calls']}")
print(f"에러율: {stats['error_rate']:.1%}")
print(f"평균 대기 시간: {stats['avg_wait_time']:.3f}초")
```

### 통계 파일 저장
```python
# 수동 저장
filepath = broker.rate_limiter.save_stats()
print(f"통계 저장됨: {filepath}")

# 자동 저장 활성화 (5분마다)
broker.rate_limiter.enable_auto_save(interval_seconds=300)

# 프로그램 종료 시 자동 저장
broker.shutdown()  # 통계가 자동으로 저장됨
```

## 대량 데이터 처리 (배치 처리)

많은 종목을 조회할 때 배치 처리를 사용하여 서버 부하를 분산할 수 있습니다:

### 고정 배치 처리
```python
# 100개 종목을 20개씩 나누어 처리
large_stock_list = [(f"{code:06d}", "KR") for code in range(1, 101)]

# 새로운 배치 처리 API 사용 (권장)
results = broker.fetch_price_list_with_batch(
    large_stock_list,
    batch_size=20,      # 20개씩 처리
    batch_delay=1.0,    # 배치 간 1초 대기
    progress_interval=10 # 10개마다 진행상황 출력
)
```

### 동적 배치 처리
```python
# 에러율에 따라 배치 크기를 자동 조정
from korea_investment_stock.dynamic_batch_controller import DynamicBatchController

# 컨트롤러 생성 (선택사항, 없으면 자동 생성)
controller = DynamicBatchController(
    initial_batch_size=50,
    target_error_rate=0.01  # 목표 에러율 1%
)

# 동적 배치 처리 실행
results = broker.fetch_price_list_with_dynamic_batch(
    large_stock_list,
    dynamic_batch_controller=controller
)

# 처리 결과 확인
stats = controller.get_stats()
print(f"최종 배치 크기: {stats['current_batch_size']}")
print(f"전체 에러율: {stats['overall_error_rate']:.1%}")
```

### 기존 방식 (하위 호환성)
```python
# 내부 메서드 직접 사용 (비권장)
results = broker._KoreaInvestment__execute_concurrent_requests(
    broker._KoreaInvestment__fetch_price,
    large_stock_list,
    batch_size=20,
    batch_delay=1.0
)
```

## 환경 변수 설정

런타임에 Rate Limiting 동작을 조정할 수 있습니다:

```bash
# 최대 호출 수 조정 (기본값: 15)
export RATE_LIMIT_MAX_CALLS=10

# 안전 마진 조정 (기본값: 0.8)
export RATE_LIMIT_SAFETY_MARGIN=0.7
```

## 에러 처리

Rate Limit 관련 에러는 자동으로 처리됩니다:

```python
# @retry_on_rate_limit 데코레이터가 자동으로 적용됨
# EGW00201 에러 발생 시 최대 5회 재시도
resp = broker.fetch_price("005930")

# 수동으로 에러 처리하기
try:
    resp = broker.fetch_price("005930")
except Exception as e:
    if "EGW00201" in str(e):
        print("Rate limit 초과! 잠시 후 다시 시도하세요.")
```

## Circuit Breaker

연속된 실패 시 자동으로 Circuit Breaker가 작동합니다:

```python
# Circuit Breaker 상태 확인
from korea_investment_stock.enhanced_backoff_strategy import get_backoff_strategy

backoff = get_backoff_strategy()
stats = backoff.get_stats()
print(f"Circuit 상태: {stats['state']}")  # CLOSED, OPEN, HALF_OPEN
print(f"성공률: {stats['success_rate']:.1%}")
```

## 모범 사례

1. **대량 조회 시 배치 처리 사용**
   ```python
   # 좋은 예: 배치 처리
   results = broker._KoreaInvestment__execute_concurrent_requests(
       method, stock_list, batch_size=50, batch_delay=1.0
   )
   
   # 피하세요: 한 번에 너무 많은 요청
   results = broker.fetch_price_list(huge_list)  # 1000개 이상
   ```

2. **피크 시간대 고려**
   ```python
   import datetime
   
   # 장 시작/종료 시간대는 보수적으로 처리
   current_hour = datetime.datetime.now().hour
   if 9 <= current_hour <= 10 or 15 <= current_hour <= 16:
       batch_size = 20  # 작은 배치
       batch_delay = 2.0  # 긴 대기
   else:
       batch_size = 50
       batch_delay = 0.5
   ```

3. **통계 모니터링**
   ```python
   # 주기적으로 성능 확인
   if broker.rate_limiter.get_stats()['error_rate'] > 0.01:
       print("경고: 에러율이 1%를 초과했습니다!")
   ```

4. **리소스 정리**
   ```python
   # 프로그램 종료 시 반드시 호출
   broker.shutdown()
   ```

## 성능 지표

- **안정적인 처리량**: 10-12 TPS 유지
- **에러율**: 0% (목표 <1%)
- **100개 종목 조회**: 약 8-10초
- **자동 복구**: Rate Limit 에러 시 지수 백오프로 재시도

이 Rate Limiting 시스템은 한국투자증권 API의 제한사항을 자동으로 처리하여, 개발자가 비즈니스 로직에만 집중할 수 있도록 합니다.
