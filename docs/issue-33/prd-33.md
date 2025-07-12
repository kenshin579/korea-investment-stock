# PRD-33: 미국 주식 조회 기능 개선

## 1. 개요

### 1.1 프로젝트 배경
현재 `korea-investment-stock` 라이브러리는 한국투자증권 API를 활용한 주식 거래 시스템을 제공하고 있으나, 미국 주식(NASDAQ, NYSE, AMEX) 조회 기능이 부분적으로만 구현되어 있어 사용자가 미국 주식 정보를 완전하게 조회할 수 없는 상황입니다.

### 1.2 목표
- 미국 주식 ticker(예: AAPL, TSLA)로 종목 정보, 현재가, 상세 정보를 완벽하게 조회할 수 있도록 기능 보완
- 국내 주식과 동일한 수준의 API 인터페이스 제공
- 기존 코드와의 일관성 유지 및 backward compatibility 보장

## 2. 현재 상태 분석

### 2.1 구현 완료된 기능
| 기능 | 메서드 | 상태 | 비고 |
|------|--------|------|------|
| 해외주식 상세시세 조회 | `fetch_price_detail_oversea_list()` | ✅ 정상 | NASDAQ, NYSE, AMEX 지원 |
| 해외주식 상세시세 조회 (내부) | `__fetch_price_detail_oversea()` | ✅ 정상 | 캐시 및 재시도 지원 |

### 2.2 미구현/오류 기능
| 기능 | 메서드 | 문제점 | 영향도 |
|------|--------|--------|--------|
| 미국 주식 현재가 조회 | `fetch_price_list()` | `fetch_oversea_price()` 메서드 없음 | 높음 |
| 미국 주식 종목정보 조회 | `fetch_stock_info_list()` | 국내 전용 API 사용 | 중간 |
| 미국 주식 검색정보 조회 | `fetch_search_stock_info_list()` | 국내 전용 API 사용 | 낮음 |

### 2.3 코드 분석 결과
```python
# __fetch_price() 메서드 내부
elif market == "US":
    resp_json = self.fetch_oversea_price(symbol)  # ❌ 메서드 없음
```

## 3. 요구사항

### 3.1 기능 요구사항

#### FR-001: 미국 주식 현재가 조회
- **설명**: 미국 주식 ticker로 현재가를 조회할 수 있어야 함
- **입력**: 종목 코드(예: "AAPL"), 시장 코드("US")
- **출력**: 현재가, 등락률, 거래량 등 시세 정보
- **우선순위**: P0 (Critical)

#### FR-002: 미국 주식 종목정보 조회
- **설명**: 미국 주식의 기본 정보를 조회할 수 있어야 함
- **입력**: 종목 코드, 시장 코드
- **출력**: 종목명, 시가총액, 섹터 정보 등
- **우선순위**: P1 (High)

#### FR-003: 통합 인터페이스 일관성
- **설명**: 국내/해외 주식 조회 시 동일한 메서드 사용
- **입력**: stock_list = [("005930", "KR"), ("AAPL", "US")]
- **출력**: 통일된 형식의 응답
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

### 4.1 단기 해결책 (Quick Fix)

#### 옵션 1: `__fetch_price()` 메서드 수정
```python
def __fetch_price(self, symbol: str, market: str = "KR") -> dict:
    if market == "KR" or market == "KRX":
        # 기존 로직 유지
    elif market == "US":
        # fetch_oversea_price() 대신 __fetch_price_detail_oversea() 활용
        detail_result = self.__fetch_price_detail_oversea(symbol, market)
        # 응답 형식을 fetch_price와 일치하도록 변환
        return self._convert_oversea_detail_to_price(detail_result)
```

**장점**: 
- 빠른 구현 가능
- 기존 코드 최소 변경

**단점**:
- 상세 API 호출로 인한 오버헤드
- 응답 변환 로직 필요

### 4.2 장기 해결책 (Complete Implementation)

#### 1. `fetch_oversea_price()` 메서드 구현
```python
@cacheable(
    ttl=300,  # 5분
    key_generator=lambda self, symbol: f"fetch_oversea_price:{symbol}"
)
@retry_on_rate_limit()
def fetch_oversea_price(self, symbol: str) -> dict:
    """해외주식 현재가 조회
    
    Args:
        symbol (str): 종목코드 (예: AAPL, TSLA)
        
    Returns:
        dict: 현재가 정보
    """
    self.rate_limiter.acquire()
    
    path = "uapi/overseas-price/v1/quotations/price"
    url = f"{self.base_url}/{path}"
    
    headers = {
        "content-type": "application/json",
        "authorization": self.access_token,
        "appKey": self.api_key,
        "appSecret": self.api_secret,
        "tr_id": "HHDFS00000300"
    }
    
    # 미국 시장 거래소 순회 (NASDAQ -> NYSE -> AMEX)
    for market_code in MARKET_TYPE_MAP["US"]:
        market_type = MARKET_CODE_MAP[market_code]
        exchange_code = EXCHANGE_CODE_MAP[market_type]
        
        params = {
            "AUTH": "",
            "EXCD": exchange_code,
            "SYMB": symbol
        }
        
        resp = requests.get(url, headers=headers, params=params)
        resp_json = resp.json()
        
        if resp_json['rt_cd'] == API_RETURN_CODE["SUCCESS"] and resp_json['output'].get('rsym'):
            return resp_json
            
    # 모든 거래소에서 찾지 못한 경우
    return {
        'rt_cd': API_RETURN_CODE["NO_DATA"],
        'msg1': f'Symbol {symbol} not found in US markets'
    }
```

#### 2. `fetch_oversea_stock_info()` 메서드 구현
```python
@cacheable(
    ttl=18000,  # 5시간
    key_generator=lambda self, symbol: f"fetch_oversea_stock_info:{symbol}"
)
@retry_on_rate_limit()
def fetch_oversea_stock_info(self, symbol: str) -> dict:
    """해외주식 종목정보 조회"""
    # 해외주식 종목정보 API 엔드포인트 사용
    pass
```

## 5. 구현 계획

### 5.1 Phase 1: Quick Fix (1주)
- [ ] `__fetch_price()` 메서드 수정
- [ ] 응답 변환 유틸리티 구현
- [ ] 단위 테스트 작성
- [ ] 통합 테스트 실행

### 5.2 Phase 2: Complete Implementation (2-3주)
- [ ] `fetch_oversea_price()` 메서드 구현
- [ ] `fetch_oversea_stock_info()` 메서드 구현
- [ ] 캐시 통합
- [ ] Rate Limiter 통합
- [ ] 에러 처리 강화
- [ ] 문서화

### 5.3 Phase 3: Testing & Optimization (1주)
- [ ] 성능 테스트
- [ ] 부하 테스트
- [ ] 실제 시장 데이터 검증
- [ ] 모니터링 대시보드 업데이트

## 6. 테스트 계획

### 6.1 단위 테스트
```python
def test_fetch_oversea_price():
    """미국 주식 가격 조회 테스트"""
    result = broker.fetch_oversea_price("AAPL")
    assert result['rt_cd'] == '0'
    assert 'output' in result
    assert 'last' in result['output']  # 현재가
```

### 6.2 통합 테스트
```python
def test_mixed_market_price_list():
    """국내/해외 혼합 조회 테스트"""
    stock_list = [
        ("005930", "KR"),  # 삼성전자
        ("AAPL", "US"),   # 애플
        ("TSLA", "US")    # 테슬라
    ]
    results = broker.fetch_price_list(stock_list)
    assert len(results) == 3
    assert all(r['rt_cd'] == '0' for r in results)
```

## 7. 위험 요소 및 대응 방안

### 7.1 위험 요소
| 위험 | 확률 | 영향도 | 대응 방안 |
|------|------|--------|-----------|
| API 스펙 변경 | 낮음 | 높음 | API 버전 관리, 변경사항 모니터링 |
| Rate Limit 초과 | 중간 | 중간 | 백오프 전략 강화, 캐시 활용 극대화 |
| 거래소별 심볼 차이 | 높음 | 낮음 | 심볼 매핑 테이블 구축 |

### 7.2 의존성
- 한국투자증권 API 해외주식 엔드포인트 가용성
- 기존 rate limiting 시스템
- 캐시 시스템

## 8. 성공 지표

### 8.1 기능적 지표
- 미국 주식 조회 성공률 > 99%
- API 응답 시간 < 500ms (캐시 미적중 시)
- 모든 미국 거래소(NASDAQ, NYSE, AMEX) 지원

### 8.2 품질 지표
- 테스트 커버리지 > 90%
- 에러 발생률 < 0.1%
- 문서화 완성도 100%

## 9. 참고 자료

### 9.1 관련 문서
- 한국투자증권 OpenAPI 개발 가이드
- [Issue #27: Rate Limiting 구현](../issue-27/)
- [Issue #30: IPO 조회 기능](../issue-30/)

### 9.2 API 엔드포인트
- 해외주식 현재가: `/uapi/overseas-price/v1/quotations/price`
- 해외주식 현재가상세: `/uapi/overseas-price/v1/quotations/price-detail`
- 해외주식 종목정보: TBD (API 문서 확인 필요)

---

**작성일**: 2025-01-13  
**작성자**: AI Assistant  
**버전**: 1.0.0  
**상태**: 초안 