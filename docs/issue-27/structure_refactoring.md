# 메인 모듈 위치 리팩토링

Date: 2024-12-28

## 변경 개요

Python 라이브러리의 표준 구조에 맞게 메인 모듈인 `korea_investment_stock.py`를 `core/` 폴더에서 패키지 루트로 이동했습니다.

## 변경 전 구조

```
korea_investment_stock/
├── __init__.py
├── core/
│   ├── __init__.py
│   └── korea_investment_stock.py  # 메인 클래스가 여기 있었음
├── rate_limiting/
├── error_handling/
└── ...
```

## 변경 후 구조

```
korea_investment_stock/
├── __init__.py
├── korea_investment_stock.py  # 메인 클래스를 루트로 이동
├── utils/                     # 헬퍼 함수 및 내부 유틸리티
│   └── __init__.py
├── rate_limiting/
├── error_handling/
└── ...
```

## 변경 이유

1. **Python 커뮤니티 표준 준수**
   - requests, flask, pandas 등 대부분의 라이브러리가 메인 모듈을 패키지 루트에 배치
   
2. **Import 경로 단순화**
   ```python
   # 이전 (깊은 중첩)
   from korea_investment_stock.core.korea_investment_stock import KoreaInvestment
   
   # 현재 (간단하고 직관적)
   from korea_investment_stock import KoreaInvestment
   ```

3. **사용자 경험 개선**
   - 더 짧고 기억하기 쉬운 import 경로
   - IDE 자동완성이 더 효과적으로 작동

## 변경 사항

### 1. 파일 이동
- `korea_investment_stock/core/korea_investment_stock.py` → `korea_investment_stock/korea_investment_stock.py`

### 2. Import 경로 수정
- 메인 모듈 내부의 상대 import 경로 수정 (`..` → `.`)
- `__init__.py`의 import 경로 수정
- 테스트 파일들의 import 경로 수정 (`..core` → `..`)

### 3. utils 폴더 생성
- core 폴더를 utils로 이름 변경 (더 명확한 용도 표현)
- 향후 헬퍼 함수와 내부 유틸리티들을 위한 공간
- 현재는 비어있는 상태

## 호환성

- **하위 호환성 완전 유지**: 공개 API 변경 없음
- 사용자는 기존과 동일하게 `from korea_investment_stock import KoreaInvestment` 사용 가능
- 내부 구조만 개선되었으므로 사용자 코드 수정 불필요

## 참고: 다른 라이브러리들의 구조

### Requests
```
requests/
├── __init__.py
├── api.py       # 메인 API 함수들
├── models.py    # Request, Response 클래스
└── ...
```

### Flask
```
flask/
├── __init__.py
├── app.py       # Flask 클래스
└── ...
```

### Pandas
```
pandas/
├── __init__.py
├── core/        # 내부 구현
└── ...
# DataFrame 등은 __init__.py에서 export
```

이처럼 대부분의 Python 라이브러리는 메인 클래스를 패키지 루트에 두고, 
내부 유틸리티는 `utils/`, `core/`, `_internal/` 등의 하위 폴더에 배치합니다. 