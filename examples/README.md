# Korea Investment Stock 예제 모음

이 폴더에는 korea-investment-stock 패키지의 주요 기능을 보여주는 예제들이 있습니다.

## 📚 예제 파일 설명

### 1. **visualization_integrated_example.py** ⭐ 권장
- **목적**: 패키지에 통합된 시각화 기능 사용법
- **특징**: 
  - 가장 간단하고 권장되는 사용법
  - `broker.create_monitoring_dashboard()` 등의 메서드 활용
  - 실시간 대시보드, 차트 생성
- **실행**: `python visualization_integrated_example.py`

### 2. **rate_limiting_example.py**
- **목적**: Rate Limiting 기능 시연
- **특징**:
  - 자동 속도 제한
  - 배치 처리
  - 에러 모니터링
  - 통계 수집
- **실행**: `python rate_limiting_example.py`

### 3. **stats_management_example.py**
- **목적**: 통계 수집 및 관리 기능
- **특징**:
  - 다양한 형식으로 통계 저장 (JSON, CSV, JSONL)
  - 통계 분석
  - 파일 로테이션
- **실행**: `python stats_management_example.py`

### 4. **stats_visualization_plotly.py** (고급)
- **목적**: 독립형 Plotly 시각화 예제
- **특징**:
  - 고급 사용자를 위한 상세 구현
  - 커스터마이징 참고용
  - 패키지 없이 독립 실행 가능
- **실행**: `python stats_visualization_plotly.py`

## 🚀 빠른 시작

### 1. 환경 설정
```bash
# API 키 설정 (환경 변수)
export KOREA_INVESTMENT_API_KEY='your_api_key'
export KOREA_INVESTMENT_API_SECRET='your_api_secret'
export KOREA_INVESTMENT_ACC_NO='12345678-01'
```

### 2. 의존성 설치
```bash
pip install plotly pandas numpy
```

### 3. 예제 실행
```bash
# 통합 시각화 예제 (권장)
python visualization_integrated_example.py

# Rate Limiting 예제
python rate_limiting_example.py
```

## 📊 생성되는 파일들

예제 실행 시 다음과 같은 파일들이 생성됩니다:

- **HTML 파일**: 인터랙티브 차트 및 대시보드
  - `system_health.html` - 시스템 상태
  - `realtime_dashboard.html` - 실시간 대시보드
  - `api_monitoring_report_*.html` - 각종 리포트

- **통계 파일**: logs 폴더에 저장
  - `logs/integrated_stats/` - 통합 통계
  - `logs/rate_limiter_stats/` - Rate Limiter 통계

## 💡 팁

1. **시각화 기능 사용**: `visualization_integrated_example.py`를 먼저 실행해보세요
2. **API 키 없이 테스트**: `mock=True` 옵션으로 모의 서버 사용 가능
3. **커스터마이징**: `stats_visualization_plotly.py`의 코드를 참고하여 커스텀 차트 생성

## 📌 주의사항

- 실제 API 사용 시 유효한 API 키가 필요합니다
- 일부 기능은 plotly 설치가 필요합니다
- PNG/PDF 내보내기는 kaleido 추가 설치가 필요합니다 