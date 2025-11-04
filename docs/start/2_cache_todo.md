# 캐싱 기능 구현 Todo

## Phase 1: 핵심 구현

### 1.1 CacheManager 구현
- [ ] `korea_investment_stock/cache_manager.py` 파일 생성
- [ ] `CacheEntry` 클래스 구현
  - [ ] `__init__()`: 데이터 저장, TTL 설정
  - [ ] `is_expired()`: 만료 여부 확인
  - [ ] `age_seconds()`: 캐시 생성 후 경과 시간
- [ ] `CacheManager` 클래스 구현
  - [ ] `__init__()`: 캐시 딕셔너리, lock, 통계 초기화
  - [ ] `get()`: Thread-safe 캐시 조회
  - [ ] `set()`: Thread-safe 캐시 저장
  - [ ] `invalidate()`: 특정 캐시 무효화
  - [ ] `clear()`: 전체 캐시 삭제
  - [ ] `get_stats()`: 캐시 통계 반환
  - [ ] `get_cache_info()`: 캐시 엔트리 정보

### 1.2 CachedKoreaInvestment 래퍼 구현
- [ ] `korea_investment_stock/cached_korea_investment.py` 파일 생성
- [ ] `CachedKoreaInvestment` 클래스 구현
  - [ ] `__init__()`: broker, TTL 설정, CacheManager 초기화
  - [ ] `_make_cache_key()`: 캐시 키 생성 로직
  - [ ] `fetch_price()`: 캐싱 지원 추가
  - [ ] `fetch_domestic_price()`: 캐싱 지원 추가
  - [ ] `fetch_etf_domestic_price()`: 캐싱 지원 추가
  - [ ] `fetch_price_detail_oversea()`: 캐싱 지원 추가
  - [ ] `fetch_stock_info()`: 캐싱 지원 추가
  - [ ] `fetch_search_stock_info()`: 캐싱 지원 추가
  - [ ] `fetch_kospi_symbols()`: 캐싱 지원 추가
  - [ ] `fetch_kosdaq_symbols()`: 캐싱 지원 추가
  - [ ] `fetch_ipo_schedule()`: 캐싱 지원 추가
  - [ ] `invalidate_cache()`: 캐시 무효화 메서드
  - [ ] `get_cache_stats()`: 캐시 통계 메서드
  - [ ] `__enter__()`, `__exit__()`: 컨텍스트 매니저 지원

### 1.3 모듈 export 설정
- [ ] `korea_investment_stock/__init__.py` 업데이트
  - [ ] `CacheManager` import 추가
  - [ ] `CacheEntry` import 추가
  - [ ] `CachedKoreaInvestment` import 추가
  - [ ] `__all__` 리스트 업데이트

---

## Phase 2: 테스트 작성

### 2.1 단위 테스트
- [ ] `korea_investment_stock/tests/test_cache_manager.py` 파일 생성
- [ ] `TestCacheEntry` 클래스
  - [ ] `test_cache_entry_creation`: 엔트리 생성 테스트
  - [ ] `test_cache_entry_expiration`: 만료 동작 테스트
- [ ] `TestCacheManager` 클래스
  - [ ] `test_cache_set_get`: 저장/조회 테스트
  - [ ] `test_cache_miss`: 캐시 미스 테스트
  - [ ] `test_cache_expiration`: 만료 후 삭제 테스트
  - [ ] `test_cache_invalidation`: 무효화 테스트
  - [ ] `test_cache_clear`: 전체 삭제 테스트
  - [ ] `test_cache_stats`: 통계 테스트
  - [ ] `test_thread_safety`: Thread-safe 테스트

### 2.2 통합 테스트
- [ ] `korea_investment_stock/tests/test_cached_integration.py` 파일 생성
- [ ] `TestCachedKoreaInvestment` 클래스
  - [ ] `test_cached_fetch_price`: 가격 조회 캐싱 테스트
  - [ ] `test_cached_expiration`: 캐시 만료 테스트
  - [ ] `test_cache_disabled`: 캐시 비활성화 테스트
  - [ ] `test_cache_invalidation`: 캐시 무효화 테스트
  - [ ] `test_multiple_symbols`: 여러 종목 캐싱 테스트
  - [ ] `test_context_manager`: 컨텍스트 매니저 테스트

### 2.3 기존 테스트 확인
- [ ] 전체 테스트 실행: `pytest`
- [ ] 기존 테스트 100% 통과 확인
- [ ] 캐싱 기능 추가가 기존 기능에 영향 없음 확인

---

## Phase 3: 문서화

### 3.1 사용 예제 작성
- [ ] `examples/cached_basic_example.py` 파일 생성
  - [ ] 환경 설정 가이드 (가상환경 생성 및 활성화)
  - [ ] 기본 사용법 예제
  - [ ] TTL 커스터마이징 예제
  - [ ] 컨텍스트 매니저 예제
  - [ ] 캐시 제어 예제

### 3.2 README.md 업데이트
- [ ] "캐싱 기능" 섹션 추가
- [ ] 기본 사용법 설명
- [ ] TTL 설정 가이드
- [ ] 성능 개선 예상치

### 3.3 CLAUDE.md 업데이트
- [ ] `CachedKoreaInvestment` 클래스 설명 추가
- [ ] 캐싱 패턴 가이드 추가
- [ ] 주의사항 및 Best Practices 추가

---

## 성공 기준

### 기능 요구사항
- [ ] 메모리 기반 캐싱 동작
- [ ] 데이터별 TTL 차등 적용
- [ ] Thread-safe 동작
- [ ] 캐시 통계 제공
- [ ] 컨텍스트 매니저 지원

### 품질 요구사항
- [ ] 테스트 커버리지 90% 이상
- [ ] 기존 테스트 100% 통과
- [ ] 문서화 완료

### 철학 준수
- [ ] 기존 코드 100% 하위 호환
- [ ] 옵트인 방식 (기본 비활성화)
- [ ] 투명하고 명시적인 동작
- [ ] 사용자 제어 가능

---

## 검증 항목

### 기능 검증
```bash
# 단위 테스트
pytest korea_investment_stock/tests/test_cache_manager.py -v

# 통합 테스트
pytest korea_investment_stock/tests/test_cached_integration.py -v

# 전체 테스트
pytest -v
```

### 성능 검증
```bash
# 캐싱 전후 성능 비교
python examples/cached_basic_example.py
```

### 사용성 검증
```python
# 기존 코드 동작 확인
broker = KoreaInvestment(api_key, api_secret, acc_no)
result = broker.fetch_price("005930", "KR")  # ✅ 정상 동작

# 캐싱 적용 확인
cached_broker = CachedKoreaInvestment(broker)
result = cached_broker.fetch_price("005930", "KR")  # ✅ 캐싱 동작
```
