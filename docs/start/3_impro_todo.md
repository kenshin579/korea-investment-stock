# korea_investment_stock.py 리팩토링 TODO

> PRD: `3_impro_prd.md` | 구현 가이드: `3_impro_implementation.md`

---

## Phase 1: 즉시 정리 (삭제만) - 1~2시간

### 1.1 사용 안 하는 import 제거
- [ ] `import pickle` 삭제 (라인 7)
- [ ] `from typing import List` 삭제 (라인 14에서 List만 제거)

### 1.2 DEPRECATED 메서드 제거
- [ ] `__handle_rate_limit_error` 메서드 삭제 (라인 507-524)

### 1.3 `__main__` 테스트 코드 제거
- [ ] `if __name__ == "__main__":` 블록 전체 삭제 (라인 1293-1342)

### 1.4 죽은 코드 제거
- [ ] `fetch_symbols` 메서드 삭제 (라인 749-766)
  - `self.exchange` 속성이 존재하지 않음

### 1.5 디버그 print 문을 logger.debug로 변경
- [ ] `fetch_price_detail_oversea`의 `print(...)` → `logger.debug(...)` 변경 (라인 1014)
- [ ] `fetch_stock_info`의 `print(e)` → `logger.debug(...)` 변경 (라인 1055)

### 1.6 Phase 1 검증
- [ ] `pytest` 실행 - 모든 테스트 통과
- [ ] `python -c "from korea_investment_stock import KoreaInvestment"` 성공

---

## Phase 2: 상수 분리 - 1~2시간

### 2.1 constants.py 생성
- [ ] `korea_investment_stock/constants.py` 파일 생성
- [ ] `EXCHANGE_CODE` → `EXCHANGE_CODE_QUOTE` 이름 변경 후 이동
- [ ] `EXCHANGE_CODE2` → `EXCHANGE_CODE_ORDER` 이름 변경 후 이동
- [ ] `EXCHANGE_CODE3` → `EXCHANGE_CODE_BALANCE` 이름 변경 후 이동
- [ ] `EXCHANGE_CODE4` → `EXCHANGE_CODE_DETAIL` 이름 변경 후 이동
- [ ] `CURRENCY_CODE` 이동
- [ ] `MARKET_TYPE_MAP` 이동
- [ ] `MARKET_TYPE`, `EXCHANGE_TYPE` 타입 정의 이동
- [ ] `MARKET_CODE_MAP`, `EXCHANGE_CODE_MAP` 이동
- [ ] `API_RETURN_CODE` 이동
- [ ] 하위 호환성 alias 추가 (기존 이름 유지)

### 2.2 메인 파일 수정
- [ ] `korea_investment_stock.py`에서 상수 정의 삭제 (라인 27-160)
- [ ] `from .constants import ...` 추가

### 2.3 Phase 2 검증
- [ ] `pytest` 실행 - 모든 테스트 통과
- [ ] `python -c "from korea_investment_stock.constants import MARKET_TYPE_MAP"` 성공

---

## Phase 3: 설정 로직 분리 - 2~3시간

### 3.1 ConfigResolver 클래스 생성
- [ ] `korea_investment_stock/config_resolver.py` 파일 생성
- [ ] `ConfigResolver` 클래스 구현
- [ ] `resolve()` 메서드 구현
- [ ] `_merge_config()` 메서드 구현
- [ ] `_load_default_config_file()` 메서드 구현
- [ ] `_load_config_file()` 메서드 구현
- [ ] `_load_from_env()` 메서드 구현

### 3.2 메인 파일 수정
- [ ] `KoreaInvestment`에서 `_resolve_config` 관련 메서드 삭제
- [ ] `KoreaInvestment`에서 `_merge_config` 삭제
- [ ] `KoreaInvestment`에서 `_load_default_config_file` 삭제
- [ ] `KoreaInvestment`에서 `_load_config_file` 삭제
- [ ] `KoreaInvestment`에서 `_load_from_env` 삭제
- [ ] `__init__`에서 `ConfigResolver` 사용하도록 수정

### 3.3 Phase 3 검증
- [ ] `pytest` 실행 - 모든 테스트 통과
- [ ] Config 관련 테스트 통과 확인

---

## Phase 4: 파서 분리 - 2~3시간

### 4.1 parsers 모듈 생성
- [ ] `korea_investment_stock/parsers/` 디렉토리 생성
- [ ] `korea_investment_stock/parsers/__init__.py` 생성
- [ ] `korea_investment_stock/parsers/master_parser.py` 생성

### 4.2 MasterParser 클래스 구현
- [ ] `MasterParser` 클래스 생성
- [ ] KOSPI 설정값 정의 (OFFSET, FIELD_SPECS, COLUMNS)
- [ ] KOSDAQ 설정값 정의 (OFFSET, FIELD_SPECS, COLUMNS)
- [ ] `parse_kospi()` 메서드 구현
- [ ] `parse_kosdaq()` 메서드 구현
- [ ] `_parse_master()` 공통 메서드 구현 (중복 제거)

### 4.3 메인 파일 수정
- [ ] `parse_kospi_master` 메서드 삭제
- [ ] `parse_kosdaq_master` 메서드 삭제
- [ ] `MasterParser` import 추가
- [ ] `fetch_kospi_symbols`에서 `MasterParser` 사용
- [ ] `fetch_kosdaq_symbols`에서 `MasterParser` 사용

### 4.4 Phase 4 검증
- [ ] `pytest` 실행 - 모든 테스트 통과
- [ ] `fetch_kospi_symbols()` 동작 확인
- [ ] `fetch_kosdaq_symbols()` 동작 확인

---

## Phase 5: IPO 헬퍼 분리 - 1~2시간

### 5.1 ipo 모듈 생성
- [ ] `korea_investment_stock/ipo/` 디렉토리 생성
- [ ] `korea_investment_stock/ipo/__init__.py` 생성
- [ ] `korea_investment_stock/ipo/ipo_helpers.py` 생성

### 5.2 IPO 헬퍼 함수 구현
- [ ] `parse_ipo_date_range()` 함수 이동
- [ ] `format_ipo_date()` 함수 이동
- [ ] `calculate_ipo_d_day()` 함수 이동
- [ ] `get_ipo_status()` 함수 이동
- [ ] `format_number()` 함수 이동

### 5.3 메인 파일 수정
- [ ] IPO 관련 정적 메서드 → 위임 패턴으로 변경
- [ ] `_validate_date_format`, `_validate_date_range` 유지 (fetch_ipo_schedule에서 사용)

### 5.4 Phase 5 검증
- [ ] `pytest` 실행 - 모든 테스트 통과
- [ ] IPO 관련 테스트 통과 확인

---

## Phase 6: 최종 검증 - 2~3시간

### 6.1 __init__.py 업데이트
- [ ] 새로운 모듈 export 추가
- [ ] 하위 호환성 확인

### 6.2 전체 테스트
- [ ] `pytest` 전체 테스트 통과
- [ ] `pytest korea_investment_stock/tests/test_integration_us_stocks.py -v` 통과
- [ ] `pytest korea_investment_stock/tests/test_ipo_integration.py -v` 통과

### 6.3 Import 테스트
- [ ] `python -c "from korea_investment_stock import KoreaInvestment"` 성공
- [ ] `python -c "from korea_investment_stock.constants import MARKET_TYPE_MAP"` 성공
- [ ] `python -c "from korea_investment_stock.config_resolver import ConfigResolver"` 성공
- [ ] `python -c "from korea_investment_stock.parsers import MasterParser"` 성공
- [ ] `python -c "from korea_investment_stock.ipo import ipo_helpers"` 성공

### 6.4 예제 실행
- [ ] `python examples/basic_example.py` 성공
- [ ] `python examples/ipo_schedule_example.py` 성공
- [ ] `python examples/us_stock_price_example.py` 성공

### 6.5 코드 품질 확인
- [ ] `korea_investment_stock.py` 라인 수 ≤ 400줄
- [ ] 중복 코드 0줄
- [ ] 사용 안 하는 코드 0줄

---

## 완료 체크리스트

- [ ] Phase 1 완료 (즉시 정리)
- [ ] Phase 2 완료 (상수 분리)
- [ ] Phase 3 완료 (설정 로직 분리)
- [ ] Phase 4 완료 (파서 분리)
- [ ] Phase 5 완료 (IPO 헬퍼 분리)
- [ ] Phase 6 완료 (최종 검증)
- [ ] CLAUDE.md 업데이트 (새로운 모듈 구조 문서화)
- [ ] CHANGELOG.md 업데이트

---

**예상 총 소요 시간**: 10-15시간

**권장 진행 방식**:
1. Phase 1만 먼저 완료하여 즉각적인 개선 효과
2. 나머지 Phase는 필요에 따라 점진적 진행
3. 각 Phase 완료 후 반드시 테스트 검증
