# 캐싱 기능 추가 PRD (Product Requirements Document)

## 📋 Executive Summary

### 목적
한국투자증권 API 호출 결과를 메모리에 캐싱하여 불필요한 API 호출을 줄이고 응답 속도를 개선합니다.

### 선택된 아키텍처
**🎯 Option B: Wrapper 클래스 패턴 (최종 확정)**

```python
# 기존 broker를 래핑하는 방식
cached_broker = CachedKoreaInvestment(broker, price_ttl=5)
result = cached_broker.fetch_price("005930", "KR")
```

### 핵심 요구사항
- **메모리 기반 캐싱**: 외부 의존성 없이 Python 내장 기능 사용
- **데이터별 TTL 설정**: 실시간성에 따라 차등적인 만료 시간 적용
- **투명한 사용**: 기존 API 인터페이스 유지 (옵션으로 제공)
- **철학 준수**: "Simple, transparent, flexible - let users implement features their way"

### 성공 지표
- API 호출 횟수 30-50% 감소
- 반복 조회 시 응답 속도 90% 이상 개선
- 기존 코드 100% 하위 호환성 유지

---

## 🔍 Current State Analysis

### 현재 구조 분석

#### 1. 주요 API 메서드 (10개)
```python
# 가격 조회 (고빈도, 실시간)
fetch_price(symbol, market)               # 통합 가격 조회
fetch_domestic_price(market_code, symbol) # 국내 주식
fetch_etf_domestic_price(market_code, symbol) # 국내 ETF
fetch_price_detail_oversea(symbol, market) # 해외 주식

# 종목 정보 (중빈도, 준실시간)
fetch_stock_info(symbol, market)          # 종목 정보
fetch_search_stock_info(symbol, market)   # 종목 검색

# 종목 리스트 (저빈도, 정적)
fetch_kospi_symbols()                     # KOSPI 종목 리스트
fetch_kosdaq_symbols()                    # KOSDAQ 종목 리스트

# IPO 정보 (저빈도, 정적)
fetch_ipo_schedule(from_date, to_date, symbol) # 공모주 일정
```

#### 2. 현재 API 호출 패턴
- **직접 호출**: 모든 요청이 즉시 API로 전달됨
- **토큰 관리**: `access_token`만 캐싱 (파일 기반, `~/.cache/mojito2/token.dat`)
- **Rate Limiting**: 사용자가 직접 구현해야 함
- **중복 조회**: 동일한 데이터 반복 조회 시 불필요한 API 호출 발생

#### 3. 문제점
- **불필요한 API 호출**: 짧은 시간 내 동일 데이터 반복 조회
- **Rate Limit 부담**: API 호출 제한(20 req/sec)에 쉽게 도달
- **응답 지연**: 매번 네트워크 요청으로 인한 레이턴시
- **비용**: 불필요한 네트워크 대역폭 사용

---

## 📊 Data Types & TTL Strategy

### 데이터 분류 및 캐시 전략

| 데이터 유형 | API 메서드 | 실시간성 | 권장 TTL | 이유 |
|------------|-----------|---------|---------|------|
| **실시간 가격** | `fetch_price()`<br>`fetch_domestic_price()`<br>`fetch_etf_domestic_price()`<br>`fetch_price_detail_oversea()` | 매우 높음 | **3-5초** | 주식 가격은 실시간으로 변동<br>너무 짧으면 캐싱 효과 없음<br>너무 길면 정확도 하락 |
| **종목 정보** | `fetch_stock_info()`<br>`fetch_search_stock_info()` | 중간 | **5-10분** | 종목 기본정보는 자주 변경되지 않음<br>시가총액, 업종 등은 일중 변경 가능 |
| **종목 리스트** | `fetch_kospi_symbols()`<br>`fetch_kosdaq_symbols()` | 낮음 | **1-24시간** | 상장/폐지는 드물게 발생<br>일 1회 갱신으로 충분 |
| **IPO 일정** | `fetch_ipo_schedule()` | 낮음 | **30분-1시간** | 공모주 일정은 거의 변경되지 않음<br>당일 청약 진행 시에만 빈번한 조회 |

### TTL 설정 근거

#### 1. 실시간 가격 (3-5초)
```python
# 사용 시나리오
for _ in range(100):  # 100번 반복 조회
    price = broker.fetch_price("005930", "KR")
    time.sleep(1)  # 1초 간격 조회

# 캐싱 효과
# - TTL 3초: 100회 → 34회 API 호출 (66% 감소)
# - TTL 5초: 100회 → 20회 API 호출 (80% 감소)
```

**근거:**
- 일반 투자자가 체결가를 확인하는 최소 간격: 3-5초
- 알고리즘 트레이딩: 1초 단위 → 캐시 비활성화 옵션 제공
- 실시간 호가는 별도 WebSocket API 사용 권장

#### 2. 종목 정보 (5-10분)
```python
# 종목 기본정보는 자주 변경되지 않음
# - 종목명, 업종, 상장주식수 등
# - 시가총액은 가격에 따라 변동하지만 참고용으로 충분
```

#### 3. 종목 리스트 (1-24시간)
```python
# KOSPI/KOSDAQ 종목 리스트
# - 상장: 월 1-2회 정도
# - 폐지: 분기 1-2회 정도
# - 장 시작 전 1회 조회로 하루 사용 가능
```

#### 4. IPO 일정 (30분-1시간)
```python
# 공모주 청약 일정
# - 청약일 당일: 30분 (진행 상태 확인)
# - 평상시: 1시간 (일정 확인용)
```

---

## 🏗️ Architecture Design

### 1. 설계 원칙

#### 핵심 철학 준수
> "Simple, transparent, flexible - let users implement features their way"

**구현 방침:**
1. **Optional Feature**: 캐싱은 선택 사항 (기본 비활성화)
2. **No Magic**: 캐싱 동작은 명시적이고 투명하게
3. **User Control**: 사용자가 TTL, 활성화 여부 등을 제어
4. **Backward Compatible**: 기존 코드 100% 호환

### 2. 아키텍처 옵션 비교 및 최종 결정

---

#### ❌ Option A: Decorator 패턴 (거부됨)
```python
@cached(ttl=5)
def fetch_price(self, symbol: str, market: str = "KR") -> dict:
    # 기존 코드
```

**장점:**
- 간결한 코드
- 메서드별 TTL 설정 용이

**단점:**
- ❌ 기존 철학(No decorators) 위반
- ❌ v0.6.0에서 데코레이터 모두 제거한 배경과 상충
- ❌ 투명성 저하

**결정:** 프로젝트 철학에 맞지 않아 **거부**

---

#### ✅ Option B: Wrapper 클래스 패턴 (✨ 최종 채택 ✨)
```python
class CachedKoreaInvestment:
    def __init__(self, broker: KoreaInvestment,
                 enable_cache: bool = True,
                 price_ttl: int = 5):
        self.broker = broker
        self.cache = {}
        self.ttl = {
            'price': price_ttl,
            'stock_info': 300,
            'symbols': 3600,
            'ipo': 1800
        }
```

**장점:**
- ✅ 기존 `KoreaInvestment` 클래스 불변
- ✅ 사용자 선택으로 캐싱 활성화
- ✅ 투명하고 명시적인 동작
- ✅ 철학 완벽 준수

**단점:**
- 추가 클래스 필요 (하지만 분리가 더 명확함)
- 약간의 추가 코드 (하지만 유지보수성 향상)

**결정:** 프로젝트 철학에 완벽히 부합하여 **최종 채택** 🎯

---

#### ❌ Option C: 내장 캐싱 옵션 (거부됨)
```python
class KoreaInvestment:
    def __init__(self, ..., enable_cache: bool = False):
        self.cache_enabled = enable_cache
```

**장점:**
- 단일 클래스 유지
- 간편한 사용

**단점:**
- ❌ 기존 철학(No built-in features) 위반
- ❌ 클래스 복잡도 증가
- ❌ v0.6.0 simplification과 상충

**결정:** v0.6.0 단순화 정신에 반하여 **거부**

---

### 3. 🎯 최종 결정: Option B - Wrapper 클래스 패턴

#### 선택 이유:
1. **철학 100% 준수**: "Simple, transparent, flexible - let users implement features their way"
2. **완벽한 투명성**: 캐싱 로직이 분리되어 이해하기 쉬움
3. **최대 유연성**: 사용자가 자체 캐싱 전략 구현 가능
4. **완전한 하위호환**: 기존 코드에 전혀 영향 없음
5. **v0.6.0 정신 유지**: 핵심 클래스는 단순하게, 기능은 외부로

#### 사용 예시:
```python
# 기존 코드 (변경 없음)
broker = KoreaInvestment(api_key, api_secret, acc_no)
result = broker.fetch_price("005930", "KR")  # ✅ 그대로 동작

# 캐싱 원하는 경우만 래퍼 사용 (Opt-in)
cached_broker = CachedKoreaInvestment(broker, price_ttl=5)
result = cached_broker.fetch_price("005930", "KR")  # ✅ 캐싱 적용
```

---

## ⚡ Performance Considerations

### 1. 메모리 사용량

| 데이터 유형 | 평균 응답 크기 | 1000개 캐시 시 메모리 |
|------------|--------------|-------------------|
| 실시간 가격 | ~2KB | ~2MB |
| 종목 정보 | ~5KB | ~5MB |
| 종목 리스트 | ~500KB | ~500MB (1회만) |
| IPO 일정 | ~10KB | ~10MB |

**예상 메모리 사용량:**
- 일반 사용: 10-50MB
- 대량 조회: 100-200MB
- 최대 (모든 종목): ~1GB

### 2. API 호출 감소 예상

| 시나리오 | 캐싱 전 | 캐싱 후 | 감소율 |
|---------|--------|--------|-------|
| 동일 종목 반복 조회 (1분) | 60회 | 12회 | 80% |
| 종목 리스트 조회 (하루) | 10회 | 1회 | 90% |
| IPO 일정 조회 (하루) | 20회 | 1회 | 95% |

### 3. 응답 시간 개선

| 작업 | 캐싱 전 | 캐싱 후 | 개선율 |
|-----|--------|--------|-------|
| 가격 조회 | 100-300ms | <1ms | 99% |
| 종목 정보 | 150-400ms | <1ms | 99% |
| 종목 리스트 | 2-5초 | <1ms | 99% |

---

## ⚠️ Limitations & Considerations

### 1. 실시간성 Trade-off
- **짧은 TTL**: API 호출 많음, 실시간성 높음
- **긴 TTL**: API 호출 적음, 실시간성 낮음
- **권장**: 사용 시나리오에 맞게 TTL 조정

### 2. 메모리 제약
- **대량 데이터**: 종목 리스트 전체 캐싱 시 메모리 사용량 큼
- **해결책**: 필요 시 외부 캐시 사용 (사용자 구현)

### 3. Multi-process 환경
- **현재**: 프로세스별 독립 캐시
- **해결책**: 필요 시 공유 캐시 구현 (사용자 선택)

### 4. 캐시 Warm-up 권장
```python
# 장 시작 전 주요 종목 캐싱
cached_broker = CachedKoreaInvestment(broker)
symbols = ["005930", "000660", "035720"]
for symbol in symbols:
    cached_broker.fetch_price(symbol, "KR")
```

---

## ✅ 최종 결정 요약

### 채택된 아키텍처
**🎯 Option B: Wrapper 클래스 패턴**

### 핵심 구성요소
1. **CacheManager**: Thread-safe 메모리 캐시 관리
2. **CachedKoreaInvestment**: 래퍼 클래스로 기존 broker 래핑
3. **데이터별 TTL**: 실시간 가격(5초), 종목정보(5분), 종목리스트(1시간), IPO(30분)

### 철학 준수 체크리스트
- ✅ Simple: 추가 클래스 하나로 명확한 구조
- ✅ Transparent: 캐싱 동작이 명시적이고 투명
- ✅ Flexible: 사용자가 TTL과 활성화 여부 제어
- ✅ Optional: 기존 코드는 그대로, 원하는 경우만 사용 (Opt-in)
- ✅ Backward Compatible: 기존 코드 100% 동작 보장

---

## 📚 References

### 관련 이슈
- v0.6.0: 모든 캐싱/장식자 기능 제거 (#40)
- 철학: "Simple, transparent, flexible"

### Korea Investment API
- Rate Limit: 20 req/sec (공식)
- 권장: 15 req/sec (보수적)
- 토큰 만료: 24시간

---

## 📂 관련 문서

- **구현 가이드**: [2_cache_implementation.md](2_cache_implementation.md)
- **Todo 체크리스트**: [2_cache_todo.md](2_cache_todo.md)

---

**작성일**: 2025-11-04
**버전**: 1.1
**상태**: Final
