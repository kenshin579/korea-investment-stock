# PRD-33: 미국 주식 현재가 조회 통합 인터페이스 구현

## 1. 개요

### 1.1 프로젝트 배경
현재 `korea-investment-stock` 라이브러리는 한국투자증권 API를 활용한 주식 거래 시스템을 제공하고 있습니다. 사용자들은 `fetch_price_list()` 하나의 메서드로 국내와 미국 주식의 현재가를 모두 조회하기를 원하지만, 현재 미국 주식 조회 시 필요한 `fetch_oversea_price()` 메서드가 누락되어 있어 통합 인터페이스가 제대로 작동하지 않고 있습니다.

### 1.2 목표
- **단일 인터페이스**: `fetch_price_list()` 하나로 국내/미국 주식 현재가 조회
- **일관된 사용 경험**: 시장 구분 없이 동일한 메서드로 조회
- **기존 코드 호환성**: 현재 국내 주식 조회 코드는 그대로 유지

## 2. 현재 상태 분석

### 2.1 통합 인터페이스 관점에서의 현황

#### 현재 `fetch_price_list()` 호출 플로우
```
fetch_price_list(stock_list)
    ↓
__fetch_price(symbol, market)
    ↓
    if market == "KR": 
        → fetch_domestic_price() ✅ (작동)
    elif market == "US":
        → fetch_oversea_price() ❌ (메서드 없음)
        
제안된 개선:
        → __fetch_price_detail_oversea() ✅ (이미 구현됨)
```

#### 메서드별 지원 현황
| 메서드 | 국내 주식 | 미국 주식 | 비고 |
|--------|-----------|-----------|------|
| `fetch_price_list()` | ✅ | ❌ | 미국 주식 시 에러 발생 |
| `fetch_stock_info_list()` | ✅ | ✅ | 통합 인터페이스 정상 작동 |
| `fetch_search_stock_info_list()` | ✅ | ❌ | 국내 전용 (의도된 제한) |
| `fetch_price_detail_oversea_list()` | ❌ | ✅ | 해외 전용 별도 메서드 |

### 2.4 코드 분석 결과
```python
# __fetch_price() 메서드 내부
elif market == "US":
    resp_json = self.fetch_oversea_price(symbol)  # ❌ 메서드 없음
```

## 3. 요구사항

### 3.1 기능 요구사항

#### FR-001: 통합 현재가 조회 인터페이스
- **설명**: `fetch_price_list()` 하나로 국내/미국 주식 현재가 모두 조회
- **예시 코드**:
  ```python
  stock_list = [
      ("005930", "KR"),  # 삼성전자
      ("AAPL", "US"),    # 애플
      ("TSLA", "US"),    # 테슬라
      ("035720", "KR")   # 카카오
  ]
  results = broker.fetch_price_list(stock_list)  # 모두 정상 조회
  ```
- **우선순위**: P0 (Critical)

#### FR-002: `__fetch_price()` 메서드 수정
- **설명**: 미국 주식 조회 시 기존 `__fetch_price_detail_oversea()` 활용하도록 수정
- **구현 방안**: `market == "US"`일 때 `__fetch_price_detail_oversea()` 호출
- **장점**: 이미 검증된 메서드 활용으로 안정성 확보
- **우선순위**: P0 (Critical)

#### FR-003: 일관된 응답 형식
- **설명**: 국내/해외 주식 조회 결과가 동일한 구조를 가져야 함
- **요구사항**: 필요시 응답 변환 로직 추가
- **우선순위**: P1 (High)

#### FR-004: API 캡슐화
- **설명**: 내부 구현 메서드들을 private으로 변경하여 통합 인터페이스 사용 유도
- **변경 대상**:
  - `fetch_etf_domestic_price()` → `__fetch_etf_domestic_price()`
  - `fetch_domestic_price()` → `__fetch_domestic_price()`
- **목적**: 사용자는 `fetch_price_list()` 하나만 사용하도록 API 단순화
- **우선순위**: P1 (High)

### 3.2 비기능 요구사항

#### NFR-001: 성능
- API 응답 시간: 개별 조회 < 500ms
- 배치 조회: 100개 종목 < 10초
- 캐시 활용으로 중복 조회 최소화

#### NFR-002: 안정성
- Rate Limit 준수 (15 calls/sec)
- 재시도 로직 구현 (exponential backoff)
- 에러 복구 시스템 통합

#### NFR-003: 호환성
- 기존 API와 backward compatible
- Python 3.9+ 지원
- 기존 캐시/모니터링 시스템과 통합

## 4. 솔루션 설계

### 4.1 구현 방향
통합 인터페이스를 위해 `__fetch_price()` 메서드를 수정하여, 미국 주식 조회 시 이미 구현된 `__fetch_price_detail_oversea()` 메서드를 활용합니다. 이는 새로운 메서드를 구현하는 것보다 안정적이고 빠른 해결책입니다.

### 4.2 이 접근법의 장점
1. **빠른 구현**: 새로운 API 메서드 구현 불필요
2. **검증된 안정성**: 이미 작동 중인 `__fetch_price_detail_oversea()` 활용
3. **일관된 동작**: 기존 해외 상세 시세 조회와 동일한 로직 사용
4. **유지보수 용이**: 하나의 해외 API 호출 로직만 관리
5. **추가 정보 제공**: 현재가 외에도 PER, PBR, EPS, BPS, 매매단위, 호가단위 등 풍부한 정보 제공

### 4.3 핵심 구현 사항

#### 1. `__fetch_price()` 메서드 수정
```python
def __fetch_price(self, symbol: str, market: str = "KR") -> dict:
    """국내/해외 주식 현재가 통합 조회
    
    Args:
        symbol (str): 종목코드
        market (str): 시장 구분 ("KR", "US" 등)
        
    Returns:
        dict: 현재가 정보
    """
    if market == "KR" or market == "KRX":
        stock_info = self.__fetch_stock_info(symbol, market)
        symbol_type = self.__get_symbol_type(stock_info)
        if symbol_type == "ETF":
            resp_json = self.__fetch_etf_domestic_price("J", symbol)
        else:
            resp_json = self.__fetch_domestic_price("J", symbol)
    elif market == "US":
        # 기존: resp_json = self.fetch_oversea_price(symbol)  # 메서드 없음
        # 개선: 이미 구현된 __fetch_price_detail_oversea() 활용
        resp_json = self.__fetch_price_detail_oversea(symbol, market)
        # 참고: 이 API는 현재가 외에도 PER, PBR, EPS 등 추가 정보 제공
    else:
        raise ValueError("Unsupported market type")
    
    return resp_json
```

**참고: `__fetch_price_detail_oversea()` 구현 세부사항**
- API 엔드포인트: `/uapi/overseas-price/v1/quotations/price-detail`
- TR ID: `HHDFS76200200`
- 지원 거래소 (EXCD 파라미터):
  - `NAS`: 나스닥
  - `NYS`: 뉴욕
  - `AMS`: 아멕스
- MARKET_TYPE_MAP["US"] = ["512", "513", "529"]를 순회하며 조회

#### 2. 관련 메서드 리팩토링
국내 주식 가격 조회 메서드들을 private으로 변경하여 내부 구현 세부사항을 캡슐화합니다:
- `fetch_etf_domestic_price()` → `__fetch_etf_domestic_price()`
- `fetch_domestic_price()` → `__fetch_domestic_price()`

이를 통해 사용자는 `fetch_price_list()` 통합 인터페이스만 사용하도록 유도합니다.

#### 3. 응답 형식 통일 (선택사항)
필요한 경우 국내/해외 응답 형식을 통일하는 변환 로직을 추가할 수 있습니다:
```python
def _normalize_response(self, resp_json: dict, market: str) -> dict:
    """국내/해외 응답 형식 통일"""
    if market == "US" and resp_json.get('rt_cd') == '0':
        # 해외 주식 응답을 국내 형식과 유사하게 변환
        # (실제 필드명은 API 응답에 따라 조정 필요)
        pass
    return resp_json
```



## 5. 구현 계획

구현 작업은 3개의 Phase로 나누어 진행됩니다. 상세한 작업 항목과 체크리스트는 [TODO-33](./todo-33.md) 문서를 참조하세요.

### 5.1 Phase 1: 핵심 구현 (2-3일)
`__fetch_price()` 메서드 수정과 국내 가격 조회 메서드들의 캡슐화를 통해 통합 인터페이스의 기반을 구축합니다.

**주요 작업**:
- US market 처리 로직을 `__fetch_price_detail_oversea()` 호출로 변경
- 국내 전용 메서드들을 private으로 변경하여 API 단순화
- 기본 동작 검증을 위한 초기 테스트

### 5.2 Phase 2: 안정화 및 최적화 (3-4일)
응답 형식 통일성 검토와 종합적인 테스트를 통해 안정성을 확보합니다.

**주요 작업**:
- 국내/해외 응답 형식 분석 및 필요시 변환 로직 구현
- 통합 테스트 케이스 작성 및 에러 처리 검증
- 성능 최적화 및 모니터링 통합

### 5.3 Phase 3: 문서화 및 배포 (1-2일)
사용자를 위한 문서 업데이트와 릴리즈 준비를 진행합니다.

**주요 작업**:
- README 및 API 문서 업데이트
- 예제 코드 작성 및 기존 예제 수정
- CHANGELOG 작성 및 릴리즈 노트 준비

### 5.4 예상 구현 시간
- **전체 소요 시간**: 약 1주 (5-7 영업일)
- **난이도**: 낮음 (기존 메서드 재사용)
- **위험도**: 낮음 (검증된 코드 활용)

> 💡 **참고**: 구체적인 작업 항목과 진행 상황은 [TODO-33](./todo-33.md)에서 관리됩니다.

## 6. 테스트 계획

### 6.1 통합 인터페이스 테스트
```python
def test_unified_price_interface():
    """통합 가격 조회 인터페이스 테스트"""
    stock_list = [
        ("005930", "KR"),  # 삼성전자
        ("AAPL", "US"),    # 애플
        ("035720", "KR"),  # 카카오
        ("TSLA", "US"),    # 테슬라
        ("NVDA", "US"),    # 엔비디아
    ]
    
    # 한 번의 호출로 모든 주식 조회
    results = broker.fetch_price_list(stock_list)
    
    # 검증
    assert len(results) == 5
    assert all(r['rt_cd'] == '0' for r in results)
    
    # 국내/미국 주식 모두 정상 조회 확인
    kr_stocks = [r for r, (_, m) in zip(results, stock_list) if m == "KR"]
    us_stocks = [r for r, (_, m) in zip(results, stock_list) if m == "US"]
    
    assert len(kr_stocks) == 2
    assert len(us_stocks) == 3
```

### 6.2 내부 메서드 통합 테스트
```python
def test_fetch_price_internal_routing():
    """__fetch_price 메서드의 라우팅 로직 테스트"""
    # 국내 주식 -> __fetch_domestic_price (private)
    result = broker._KoreaInvestment__fetch_price("005930", "KR")
    assert result['rt_cd'] == '0'
    
    # 미국 주식 -> __fetch_price_detail_oversea
    result = broker._KoreaInvestment__fetch_price("AAPL", "US")
    assert result['rt_cd'] == '0'
    assert 'output' in result
    
    # 상세 시세 API 응답 확인
    assert result.get('output', {}).get('rsym') == 'AAPL'
```

### 6.3 성능 및 캐시 테스트
```python
def test_cache_performance():
    """캐시 적중률 및 성능 테스트"""
    stock_list = [("AAPL", "US"), ("TSLA", "US")] * 10  # 동일 종목 반복
    
    start_time = time.time()
    results1 = broker.fetch_price_list(stock_list)
    first_call_time = time.time() - start_time
    
    # 두 번째 호출은 캐시에서
    start_time = time.time()
    results2 = broker.fetch_price_list(stock_list)
    second_call_time = time.time() - start_time
    
    # 캐시 적중으로 두 번째가 훨씬 빨라야 함
    assert second_call_time < first_call_time * 0.1
```

## 7. 위험 요소 및 대응 방안

### 7.1 위험 요소
| 위험 | 확률 | 영향도 | 대응 방안 |
|------|------|--------|-----------|
| API 스펙 변경 | 낮음 | 높음 | API 버전 관리, 변경사항 모니터링 |
| Rate Limit 초과 | 중간 | 중간 | 백오프 전략 강화, 캐시 활용 극대화 |
| 거래소별 심볼 차이 | 높음 | 낮음 | 심볼 매핑 테이블 구축 |
| 모의투자 미지원 | 확실 | 낮음 | 문서에 명시, 실전 계정 사용 필수 |

### 7.2 의존성
- 한국투자증권 API 해외주식 엔드포인트 가용성
- 기존 rate limiting 시스템
- 캐시 시스템

## 8. 성공 지표

### 8.1 기능적 지표
- 미국 주식 조회 성공률 > 99%
- API 응답 시간 < 500ms (캐시 미적중 시)
- 모든 미국 거래소(NASDAQ, NYSE, AMEX) 지원
- 통합 인터페이스로 국내/미국 주식 혼합 조회 가능
- 미국 주식 조회 시 PER, PBR 등 추가 정보도 함께 제공

### 8.2 품질 지표
- 테스트 커버리지 > 90%
- 에러 발생률 < 0.1%
- 문서화 완성도 100%

## 9. 사용 예시

### 9.1 현재 (문제 상황)
```python
# 국내 주식만 조회 가능
kr_stocks = [("005930", "KR"), ("035720", "KR")]
kr_results = broker.fetch_price_list(kr_stocks)  # ✅ 정상

# 미국 주식 포함 시 에러
mixed_stocks = [("005930", "KR"), ("AAPL", "US")]
mixed_results = broker.fetch_price_list(mixed_stocks)  # ❌ AttributeError: fetch_oversea_price
```

### 9.2 개선 후 (목표)
```python
# 통합 인터페이스로 모든 주식 조회
mixed_stocks = [
    ("005930", "KR"),    # 삼성전자
    ("AAPL", "US"),      # 애플
    ("035720", "KR"),    # 카카오
    ("TSLA", "US"),      # 테슬라
    ("000660", "KR"),    # SK하이닉스
    ("MSFT", "US"),      # 마이크로소프트
]

# 한 번의 호출로 모든 주식 조회
results = broker.fetch_price_list(mixed_stocks)  # ✅ 모두 정상

# 결과 활용
for (symbol, market), result in zip(mixed_stocks, results):
    if result['rt_cd'] == '0':
        price = result['output']['last']  # 또는 적절한 필드
        print(f"{symbol} ({market}): {price}")
```

## 10. 참고 자료

### 10.1 관련 문서
- 한국투자증권 OpenAPI 개발 가이드
- [Issue #27: Rate Limiting 구현](../issue-27/)
- [Issue #30: IPO 조회 기능](../issue-30/)

### 10.2 API 엔드포인트
- 국내주식 현재가: `/uapi/domestic-stock/v1/quotations/inquire-price`
- 해외주식 현재가상세: `/uapi/overseas-price/v1/quotations/price-detail` (실제 사용)
  - TR ID: `HHDFS76200200`
  - Method: `GET`
  - 모의투자: **미지원**
  - 무료시세(지연시세) 제공 (미국: 실시간, 기타: 15분 지연)
- ~~해외주식 현재가: `/uapi/overseas-price/v1/quotations/price`~~ (미구현)
- 주식 종목정보: `/uapi/domestic-stock/v1/quotations/search-info`

### 10.3 참고 사항
- 해외주식 현재가상세 API는 현재가 외에도 매매단위(vnit), 호가단위(e_hogau), PER, PBR, EPS, BPS 등의 추가 정보 제공
- 모의투자에서는 해외주식 시세 조회 불가 (실전 계정 필수)
- 한국투자증권 API 상세 문서: `/docs/api/` 디렉토리 참조
- 미국 주식은 실시간 무료시세 제공 (나스닥 마켓센터 기준)
- 응답 형식이 국내/해외 간 다를 수 있으므로 필요시 변환 로직 구현 권장

---

**작성일**: 2025-01-13  
**작성자**: AI Assistant  
**버전**: 2.3.1  
**상태**: 초안

### 변경 이력
- v2.3.1 (2025-01-13): 구현 계획 개선
  - 구현 계획의 TODO 항목들을 별도 [TODO-33](./todo-33.md) 파일로 분리
  - 더 상세한 작업 항목과 체크리스트 제공
  - Phase별 주요 작업 설명 추가
- v2.3.0 (2025-01-13): 한국투자 API 문서 반영
  - `/docs/api/` 디렉토리의 공식 API 문서 정보 추가
  - 해외주식 현재가상세 API의 상세 파라미터 명시
  - 모의투자 미지원 및 지연시세 정보 추가
- v2.2.0 (2025-01-13): 메서드 캡슐화 추가
  - `fetch_etf_domestic_price()`, `fetch_domestic_price()`를 private 메서드로 변경
  - 사용자가 통합 인터페이스 `fetch_price_list()`만 사용하도록 유도
- v2.1.0 (2025-01-13): 구현 방식 변경
  - `fetch_oversea_price()` 신규 구현 대신 기존 `__fetch_price_detail_oversea()` 활용
  - 구현 기간 단축 (2-3주 → 1주)
  - 이미 검증된 메서드 재사용으로 안정성 향상
- v2.0.0 (2025-01-13): 통합 인터페이스 중심으로 전면 재작성
  - 목표 변경: 별도 해외 메서드 대신 `fetch_price_list()` 통합 인터페이스 구현
  - `fetch_price_detail_oversea_list()` 사용 방안 제거
  - 테스트 계획을 통합 인터페이스 관점으로 수정
- v1.0.1 (2025-01-13): `fetch_stock_info_list()` 국내/미국 지원 확인
- v1.0.0 (2025-01-13): 초안 작성 