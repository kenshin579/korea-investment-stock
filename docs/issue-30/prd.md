# 공모주청약일정 API 추가 개발 PRD

## 1. 개요

### 1.1 프로젝트 배경
- 한국투자증권 Open API에서 제공하는 공모주청약일정 정보를 korea-investment-stock 라이브러리에 추가
- 투자자들이 공모주 청약 일정을 쉽게 조회하고 관리할 수 있도록 지원

### 1.2 목표
- 공모주청약일정 조회 API를 라이브러리에 통합
- 기존 Rate Limiting, 캐싱, 에러 처리 시스템과 완벽한 통합
- 사용하기 쉬운 인터페이스 제공

## 2. API 분석

### 2.1 엔드포인트 정보
- **URL Path**: `/uapi/domestic-stock/v1/ksdinfo/pub-offer`
- **Method**: GET
- **Category**: [국내주식] 종목정보
- **Description**: 예탁원정보(공모주청약일정) API
- **TR_ID**: `HHKDB669108C0`
- **모의투자**: 미지원

### 2.2 요청 파라미터 (Query Parameters)
```json
{
  "SHT_CD": "",              // 종목코드 (공백: 전체, 특정종목 조회시: 종목코드)
  "CTS": "",                 // CTS (공백)
  "F_DT": "20240101",       // 조회일자From (YYYYMMDD)
  "T_DT": "20240131"        // 조회일자To (YYYYMMDD)
}
```

### 2.3 응답 형식
```json
{
  "rt_cd": "0",
  "msg_cd": "KIOK0000",
  "msg1": "정상처리되었습니다",
  "output1": [
    {
      "record_date": "20240115",        // 기준일
      "sht_cd": "123456",               // 종목코드
      "isin_name": "테스트회사",         // 종목명
      "fix_subscr_pri": "15000",        // 공모가
      "face_value": "5000",             // 액면가
      "subscr_dt": "2024.01.15~2024.01.16",  // 청약기간
      "pay_dt": "2024.01.18",           // 납입일
      "refund_dt": "2024.01.19",        // 환불일
      "list_dt": "2024.01.25",          // 상장/등록일
      "lead_mgr": "한국투자증권",        // 주간사
      "pub_bf_cap": "1000000000",       // 공모전자본금
      "pub_af_cap": "1500000000",       // 공모후자본금
      "assign_stk_qty": "100000"        // 당사배정물량
    }
  ]
}
```

### 2.4 응답 필드 상세 설명
| 필드명 | 한글명 | 타입 | 길이 | 설명 |
|--------|--------|------|------|------|
| record_date | 기준일 | String | 8 | YYYYMMDD 형식 |
| sht_cd | 종목코드 | String | 9 | 종목코드 |
| isin_name | 종목명 | String | 40 | 회사명 |
| fix_subscr_pri | 공모가 | String | 12 | 공모가격 |
| face_value | 액면가 | String | 9 | 액면가격 |
| subscr_dt | 청약기간 | String | 23 | 청약 시작일~종료일 |
| pay_dt | 납입일 | String | 10 | 납입일자 |
| refund_dt | 환불일 | String | 10 | 환불일자 |
| list_dt | 상장/등록일 | String | 10 | 상장예정일 |
| lead_mgr | 주간사 | String | 41 | 주간사명 |
| pub_bf_cap | 공모전자본금 | String | 12 | 공모 전 자본금 |
| pub_af_cap | 공모후자본금 | String | 12 | 공모 후 자본금 |
| assign_stk_qty | 당사배정물량 | String | 12 | 한국투자증권 배정물량 |

## 3. 기능 요구사항

### 3.1 기본 기능
1. **공모주 일정 조회**
   - 전체 공모주 일정 조회
   - 기간별 공모주 일정 조회
   - 특정 종목 공모주 정보 조회

2. **데이터 처리 및 변환**
   - 날짜 형식 파싱 (청약기간 문자열 처리)
   - 숫자 데이터 타입 변환
   - 청약 상태 판단 (날짜 기준)

### 3.2 캐싱 전략
- TTL: 1시간 (3600초) - 공모주 정보는 자주 변경되지 않음
- 캐시 키: `fetch_ipo_schedule:{from_date}:{to_date}:{symbol}`
- 전체 조회 시: `fetch_ipo_schedule:{from_date}:{to_date}:ALL`

### 3.3 에러 처리
- 날짜 형식 검증 (YYYYMMDD)
- 날짜 범위 검증 (시작일 <= 종료일)
- API 응답 에러 처리
- Rate Limit 에러 자동 재시도
- 모의투자 미지원 안내

## 4. 구현 계획

### 4.1 메서드 구현

```python
class KoreaInvestment:
    
    @cacheable(ttl=3600, key_generator=lambda self, from_date, to_date, symbol: f"fetch_ipo_schedule:{from_date}:{to_date}:{symbol or 'ALL'}")
    @retry_on_rate_limit()
    def fetch_ipo_schedule(self, from_date: str = None, to_date: str = None, symbol: str = "") -> dict:
        """공모주 청약 일정 조회
        
        예탁원정보(공모주청약일정) API를 통해 공모주 정보를 조회합니다.
        한국투자 HTS(eFriend Plus) > [0667] 공모주청약 화면과 동일한 기능입니다.
        
        Args:
            from_date: 조회 시작일 (YYYYMMDD, 기본값: 오늘)
            to_date: 조회 종료일 (YYYYMMDD, 기본값: 30일 후)
            symbol: 종목코드 (선택, 공백시 전체 조회)
            
        Returns:
            dict: 공모주 청약 일정 정보
                {
                    "rt_cd": "0",  # 성공여부
                    "msg_cd": "응답코드",
                    "msg1": "응답메시지",
                    "output1": [
                        {
                            "record_date": "기준일",
                            "sht_cd": "종목코드",
                            "isin_name": "종목명",
                            "fix_subscr_pri": "공모가",
                            "face_value": "액면가",
                            "subscr_dt": "청약기간",  # "2024.01.15~2024.01.16"
                            "pay_dt": "납입일",
                            "refund_dt": "환불일",
                            "list_dt": "상장/등록일",
                            "lead_mgr": "주간사",
                            "pub_bf_cap": "공모전자본금",
                            "pub_af_cap": "공모후자본금",
                            "assign_stk_qty": "당사배정물량"
                        }
                    ]
                }
                
        Raises:
            ValueError: 모의투자 사용시 또는 날짜 형식 오류시
            
        Note:
            - 모의투자는 지원하지 않습니다.
            - 예탁원에서 제공한 자료이므로 정보용으로만 사용하시기 바랍니다.
            - 실제 청약시에는 반드시 공식 공모주 청약 공고문을 확인하세요.
            
        Examples:
            >>> # 전체 공모주 조회 (오늘부터 30일)
            >>> ipos = broker.fetch_ipo_schedule()
            
            >>> # 특정 기간 조회
            >>> ipos = broker.fetch_ipo_schedule(
            ...     from_date="20240101",
            ...     to_date="20240131"
            ... )
            
            >>> # 특정 종목 조회
            >>> ipo = broker.fetch_ipo_schedule(symbol="123456")
        """
        # 모의투자 체크
        if self.is_mock:
            raise ValueError("공모주청약일정 조회는 모의투자를 지원하지 않습니다.")
            
        self.rate_limiter.acquire()
        
        # 날짜 기본값 설정
        if not from_date:
            from_date = datetime.now().strftime("%Y%m%d")
        if not to_date:
            to_date = (datetime.now() + timedelta(days=30)).strftime("%Y%m%d")
        
        # 날짜 유효성 검증
        if not validate_date_format(from_date) or not validate_date_format(to_date):
            raise ValueError("날짜 형식은 YYYYMMDD 이어야 합니다.")
        
        if not validate_date_range(from_date, to_date):
            raise ValueError("시작일은 종료일보다 이전이어야 합니다.")
        
        path = "uapi/domestic-stock/v1/ksdinfo/pub-offer"
        url = f"{self.base_url}/{path}"
        headers = self._get_base_headers()
        headers["tr_id"] = "HHKDB669108C0"
        headers["custtype"] = "P"  # 개인
        
        params = {
            "SHT_CD": symbol,
            "CTS": "",
            "F_DT": from_date,
            "T_DT": to_date
        }
        
        resp = self._execute_request("GET", url, headers=headers, params=params)
        return self._process_response(resp)
```

### 4.2 헬퍼 함수

```python
from datetime import datetime, timedelta
import re

def validate_date_format(date_str: str) -> bool:
    """날짜 형식 검증 (YYYYMMDD)"""
    if len(date_str) != 8:
        return False
    try:
        datetime.strptime(date_str, "%Y%m%d")
        return True
    except ValueError:
        return False

def validate_date_range(from_date: str, to_date: str) -> bool:
    """날짜 범위 유효성 검증"""
    try:
        start = datetime.strptime(from_date, "%Y%m%d")
        end = datetime.strptime(to_date, "%Y%m%d")
        return start <= end
    except ValueError:
        return False

def parse_ipo_date_range(date_range_str: str) -> tuple:
    """청약기간 문자열 파싱
    
    Args:
        date_range_str: "2024.01.15~2024.01.16" 형식의 문자열
        
    Returns:
        tuple: (시작일 datetime, 종료일 datetime) 또는 (None, None)
    """
    if not date_range_str:
        return (None, None)
    
    # "2024.01.15~2024.01.16" 형식 파싱
    pattern = r'(\d{4}\.\d{2}\.\d{2})~(\d{4}\.\d{2}\.\d{2})'
    match = re.match(pattern, date_range_str)
    
    if match:
        try:
            start_str = match.group(1).replace('.', '')
            end_str = match.group(2).replace('.', '')
            start_date = datetime.strptime(start_str, "%Y%m%d")
            end_date = datetime.strptime(end_str, "%Y%m%d")
            return (start_date, end_date)
        except ValueError:
            pass
    
    return (None, None)

def format_ipo_date(date_str: str) -> str:
    """날짜 형식 변환 (YYYYMMDD -> YYYY-MM-DD)"""
    if len(date_str) == 8:
        return f"{date_str[:4]}-{date_str[4:6]}-{date_str[6:8]}"
    elif '.' in date_str:
        return date_str.replace('.', '-')
    return date_str

def calculate_ipo_d_day(ipo_date_str: str) -> int:
    """청약일까지 남은 일수 계산"""
    if '~' in ipo_date_str:
        start_date, _ = parse_ipo_date_range(ipo_date_str)
        if start_date:
            today = datetime.now()
            return (start_date - today).days
    return -999

def get_ipo_status(subscr_dt: str) -> str:
    """청약 상태 판단
    
    Returns:
        str: "예정", "진행중", "마감", "알수없음"
    """
    start_date, end_date = parse_ipo_date_range(subscr_dt)
    if not start_date or not end_date:
        return "알수없음"
    
    today = datetime.now()
    if today < start_date:
        return "예정"
    elif start_date <= today <= end_date:
        return "진행중"
    else:
        return "마감"

def format_number(num_str: str) -> str:
    """숫자 문자열에 천단위 콤마 추가"""
    try:
        return f"{int(num_str):,}"
    except (ValueError, TypeError):
        return num_str
```

## 5. 사용 예제

```python
from korea_investment_stock import KoreaInvestment
from datetime import datetime, timedelta

# 클라이언트 생성
broker = KoreaInvestment(api_key, api_secret, acc_no, is_mock=False)  # 실전투자만 지원

# 1. 전체 공모주 일정 조회 (기본값: 오늘부터 30일)
result = broker.fetch_ipo_schedule()
if result.get('rt_cd') == '0':
    ipos = result.get('output1', [])
    print(f"조회된 공모주: {len(ipos)}개")

# 2. 특정 기간 조회
result = broker.fetch_ipo_schedule(
    from_date="20240101",
    to_date="20240131"
)

# 3. 특정 종목 공모주 정보 조회
result = broker.fetch_ipo_schedule(symbol="123456")
if result.get('rt_cd') == '0' and result.get('output1'):
    ipo = result['output1'][0]
    print(f"종목명: {ipo['isin_name']}")
    print(f"공모가: {format_number(ipo['fix_subscr_pri'])}원")

# 4. 활용 예제: 현재 청약 진행중인 공모주 필터링
today = datetime.now()
week_ago = (today - timedelta(days=7)).strftime("%Y%m%d")
week_later = (today + timedelta(days=7)).strftime("%Y%m%d")

result = broker.fetch_ipo_schedule(from_date=week_ago, to_date=week_later)
if result.get('rt_cd') == '0':
    for ipo in result.get('output1', []):
        status = get_ipo_status(ipo['subscr_dt'])
        if status == "진행중":
            print(f"{ipo['isin_name']} - 청약진행중 ({ipo['subscr_dt']})")

# 5. 활용 예제: 공모주 상세 정보 출력
def print_ipo_details(ipo):
    """공모주 상세 정보 출력"""
    print(f"종목명: {ipo['isin_name']} ({ipo['sht_cd']})")
    print(f"공모가: {format_number(ipo['fix_subscr_pri'])}원 (액면가: {format_number(ipo['face_value'])}원)")
    print(f"청약기간: {ipo['subscr_dt']}")
    print(f"주간사: {ipo['lead_mgr']}")
    print(f"상장예정일: {ipo['list_dt']}")
    print(f"당사배정물량: {format_number(ipo['assign_stk_qty'])}주")
    print(f"공모전 자본금: {format_number(ipo['pub_bf_cap'])}원")
    print(f"공모후 자본금: {format_number(ipo['pub_af_cap'])}원")
    print(f"상태: {get_ipo_status(ipo['subscr_dt'])}")
    print("-" * 50)

# 전체 조회 후 상세 출력
result = broker.fetch_ipo_schedule()
if result.get('rt_cd') == '0':
    for ipo in result.get('output1', []):
        print_ipo_details(ipo)

# 6. 활용 예제: 청약 D-Day 계산
result = broker.fetch_ipo_schedule()
if result.get('rt_cd') == '0':
    upcoming_ipos = []
    for ipo in result.get('output1', []):
        d_day = calculate_ipo_d_day(ipo['subscr_dt'])
        if 0 <= d_day <= 7:  # 7일 이내 청약 예정
            upcoming_ipos.append({
                'name': ipo['isin_name'],
                'subscr_dt': ipo['subscr_dt'],
                'd_day': d_day
            })
    
    # D-Day 순으로 정렬
    upcoming_ipos.sort(key=lambda x: x['d_day'])
    for ipo in upcoming_ipos:
        print(f"D-{ipo['d_day']}: {ipo['name']} ({ipo['subscr_dt']})")
```

## 6. 테스트 계획

### 6.1 단위 테스트
- 날짜 형식 검증 테스트
- 날짜 범위 검증 테스트
- 청약기간 파싱 테스트
- 청약 상태 판단 테스트
- API 응답 파싱 테스트
- 캐시 동작 테스트

### 6.2 통합 테스트
- 실제 API 호출 테스트 (실전투자만)
- Rate Limiting 테스트
- 에러 복구 테스트
- 다양한 날짜 범위 조회 테스트
- 특정 종목 조회 테스트

## 7. 문서화

### 7.1 README 업데이트
- 공모주 조회 기능 추가
- 사용 예제 추가
- 모의투자 미지원 명시
- API 제한사항 명시

### 7.2 API 문서
- 메서드별 상세 문서
- 파라미터 설명
- 응답 형식 설명
- 에러 코드 설명
- 예탁원 정보 기반 데이터임을 명시

## 8. 보안 고려사항
- 민감한 공모 정보는 로깅하지 않음
- 캐시된 데이터 암호화 고려
- API 키 노출 방지

## 9. 성능 목표
- 단일 조회: < 100ms (캐시 히트 시)
- 전체 조회: < 500ms
- 캐시 적중률: > 70%

## 10. 향후 확장 가능성
- 공모주 청약 알림 기능
- 공모주 수익률 분석 (상장 후 가격 추적)
- 청약 경쟁률 예측
- 공모주 캘린더 뷰 제공
- 주간사별 공모주 통계

## 11. 마일스톤
1. **Phase 1**: API 구현 및 기본 기능 (2일)
2. **Phase 2**: 헬퍼 함수 및 데이터 처리 (1일)
3. **Phase 3**: 캐싱 및 에러 처리 통합 (1일)
4. **Phase 4**: 테스트 작성 및 디버깅 (2일)
5. **Phase 5**: 문서화 및 예제 작성 (1일)

## 12. 리스크 및 대응방안
- **모의투자 미지원**: README 및 메서드 docstring에 명확히 표시
- **날짜 형식 불일치**: 다양한 날짜 형식 파싱 함수 구현
- **Rate Limit 제약**: 배치 처리 시 적절한 딜레이 추가
- **데이터 정확성**: 예탁원 제공 정보임을 명시, 참고용으로만 사용 권고

## 13. 주의사항
- 이 API는 한국투자 HTS(eFriend Plus) > [0667] 공모주청약 화면과 동일한 기능
- 예탁원에서 제공한 자료이므로 정보용으로만 사용
- 실제 청약 시에는 반드시 공식 공모주 청약 공고문 확인 필요
