# 코드 스타일 및 문서화 규칙

## Python 코드 스타일

### 1. 기본 규칙
- Python 3.11+ 문법 사용
- PEP 8 준수
- 한 줄은 최대 120자 (PEP 8의 79자보다 완화)

### 2. Type Hints
```python
from typing import Dict, List, Optional, Union

def get_stock_price(
    symbol: str,
    start_date: Optional[str] = None,
    end_date: Optional[str] = None
) -> Dict[str, Union[str, float]]:
    """주식 가격 조회"""
    pass
```

### 3. Docstring
```python
def calculate_profit(buy_price: float, sell_price: float, quantity: int) -> float:
    """
    수익금 계산
    
    Args:
        buy_price: 매수 단가
        sell_price: 매도 단가  
        quantity: 수량
        
    Returns:
        수익금 (수수료 제외)
        
    Raises:
        ValueError: 가격이나 수량이 음수인 경우
    """
    pass
```

### 4. 한글 사용
- 주석과 docstring은 한글 사용 권장
- 변수명과 함수명은 영어 사용
- 사용자 대면 메시지는 한글 사용

### 5. 네이밍 컨벤션
- 함수/변수: snake_case
- 클래스: PascalCase
- 상수: UPPER_SNAKE_CASE
- Private: _leading_underscore

### 6. Import 순서
```python
# 1. 표준 라이브러리
import os
import sys
from datetime import datetime

# 2. 서드파티 라이브러리
import requests
import pandas as pd

# 3. 로컬 모듈
from korea_investment_stock import KoreaInvestmentStock
from .utils import validate_stock_code
``` 