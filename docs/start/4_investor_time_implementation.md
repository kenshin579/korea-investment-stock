# 시장별 투자자매매동향(시세) API 구현

## 개요

| 항목 | 값 |
|------|-----|
| API명 | 시장별 투자자매매동향(시세) |
| API ID | v1_국내주식-074 |
| API 경로 | `/uapi/domestic-stock/v1/quotations/inquire-investor-time-by-market` |
| TR ID | FHPTJ04030000 |
| 모의투자 | 미지원 |

---

## Phase 1: korea_investment_stock 라이브러리

### 1.1 API 메서드 추가

**파일**: `korea_investment_stock/korea_investment_stock.py`

```python
def fetch_investor_trend_by_market(
    self,
    market_code: str = "KSP",
    sector_code: str = "0001"
) -> dict:
    """
    시장별 투자자매매동향(시세) API

    Args:
        market_code: 시장구분 (KSP=코스피, KSQ=코스닥, ETF, ELW, ETN 등)
        sector_code: 업종구분 (0001=코스피종합, 1001=코스닥종합 등)

    Returns:
        dict: API 응답 (투자자별 매매 현황)
    """
    path = "uapi/domestic-stock/v1/quotations/inquire-investor-time-by-market"
    url = f"{self.base_url}/{path}"

    headers = {
        "content-type": "application/json",
        "authorization": self.access_token,
        "appKey": self.api_key,
        "appSecret": self.api_secret,
        "tr_id": "FHPTJ04030000"
    }

    params = {
        "fid_input_iscd": market_code,
        "fid_input_iscd_2": sector_code
    }

    return self._request_with_token_refresh("GET", url, headers, params)
```

### 1.2 상수 추가

**파일**: `korea_investment_stock/constants.py`

```python
# 시장별 투자자동향 - 시장 코드
MARKET_INVESTOR_TREND_CODE = {
    "KOSPI": "KSP",
    "KOSDAQ": "KSQ",
    "ETF": "ETF",
    "ELW": "ELW",
    "ETN": "ETN",
    "FUTURES": "K2I",
    "STOCK_FUTURES": "999",
    "MINI": "MKI",
    "WEEKLY_MONTH": "WKM",
    "WEEKLY_THUR": "WKI",
    "KOSDAQ150": "KQI"
}

# 시장별 투자자동향 - 업종 코드 (주요 항목만)
SECTOR_CODE = {
    "KOSPI_TOTAL": "0001",
    "KOSDAQ_TOTAL": "1001",
    "FUTURES": "F001",
    "CALL_OPTION": "OC01",
    "PUT_OPTION": "OP01",
    "ETF_TOTAL": "T000",
    "ELW_TOTAL": "W000",
    "ETN_TOTAL": "E199"
}
```

### 1.3 테스트 코드

**파일**: `tests/test_investor_trend.py`

```python
import pytest
from korea_investment_stock import KoreaInvestment

class TestInvestorTrendByMarket:

    @pytest.fixture
    def client(self):
        return KoreaInvestment()

    def test_fetch_kospi_investor_trend(self, client):
        """코스피 종합 투자자 동향 조회"""
        result = client.fetch_investor_trend_by_market("KSP", "0001")

        assert result['rt_cd'] == '0'
        assert 'output' in result

    def test_fetch_kosdaq_investor_trend(self, client):
        """코스닥 종합 투자자 동향 조회"""
        result = client.fetch_investor_trend_by_market("KSQ", "1001")

        assert result['rt_cd'] == '0'
```

### 1.4 버전 업데이트

`pyproject.toml` 버전을 `0.16.0`으로 업데이트

---

## Phase 2: DB 스키마 (moneyflow.advenoh.pe.kr)

DB 스키마는 `moneyflow.advenoh.pe.kr` 레포지토리에서 Liquibase로 관리합니다.

### 2.1 Liquibase Changelog

**파일**: `moneyflow.advenoh.pe.kr/backend/db/changelog/2025-12/8_create_market_investor_trend.sql`

```sql
--liquibase formatted sql
--changeset kenshin579:create_market_investor_trend
-- comment: 시장별 투자자 매매동향(시세) 테이블 생성

CREATE TABLE IF NOT EXISTS market_investor_trend (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT 'PK',

    -- 시장 정보
    market_code VARCHAR(10) NOT NULL COMMENT '시장코드 (KSP=코스피, KSQ=코스닥 등)',
    sector_code VARCHAR(10) NOT NULL COMMENT '업종코드 (0001=코스피종합, 1001=코스닥종합 등)',
    collected_at DATETIME NOT NULL COMMENT '수집일시',
    trade_date DATE NOT NULL COMMENT '거래일',

    -- 외국인
    foreign_net_qty BIGINT NULL COMMENT '외국인 순매수량',
    foreign_net_amount BIGINT NULL COMMENT '외국인 순매수금액',

    -- 개인
    individual_net_qty BIGINT NULL COMMENT '개인 순매수량',
    individual_net_amount BIGINT NULL COMMENT '개인 순매수금액',

    -- 기관
    institution_net_qty BIGINT NULL COMMENT '기관 순매수량',
    institution_net_amount BIGINT NULL COMMENT '기관 순매수금액',

    -- 증권
    securities_net_qty BIGINT NULL COMMENT '증권 순매수량',
    securities_net_amount BIGINT NULL COMMENT '증권 순매수금액',

    -- 투자신탁
    investment_trust_net_qty BIGINT NULL COMMENT '투자신탁 순매수량',
    investment_trust_net_amount BIGINT NULL COMMENT '투자신탁 순매수금액',

    -- 사모펀드
    pe_fund_net_qty BIGINT NULL COMMENT '사모펀드 순매수량',
    pe_fund_net_amount BIGINT NULL COMMENT '사모펀드 순매수금액',

    -- 은행
    bank_net_qty BIGINT NULL COMMENT '은행 순매수량',
    bank_net_amount BIGINT NULL COMMENT '은행 순매수금액',

    -- 보험
    insurance_net_qty BIGINT NULL COMMENT '보험 순매수량',
    insurance_net_amount BIGINT NULL COMMENT '보험 순매수금액',

    -- 기금
    pension_fund_net_qty BIGINT NULL COMMENT '기금 순매수량',
    pension_fund_net_amount BIGINT NULL COMMENT '기금 순매수금액',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '생성일',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '수정일',

    UNIQUE KEY uk_market_sector_date (market_code, sector_code, trade_date),
    INDEX idx_trade_date (trade_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='시장별 투자자 매매동향';

--rollback DROP TABLE IF EXISTS market_investor_trend;
```

### 2.2 마이그레이션 실행

```bash
cd moneyflow.advenoh.pe.kr/backend/db
./liquibase.sh update
```

---

## Phase 3: stock-data-batch

### 3.1 DB 모델 추가

**파일**: `database/models.py`

```python
class MarketInvestorTrend(BaseModel):
    """시장별 투자자 매매동향 (시세)"""

    id = AutoField()
    market_code = CharField(max_length=10)
    sector_code = CharField(max_length=10)
    collected_at = DateTimeField()
    trade_date = DateField()

    # 외국인/개인/기관 (순매수량, 순매수금액)
    foreign_net_qty = BigIntegerField(null=True)
    foreign_net_amount = BigIntegerField(null=True)
    individual_net_qty = BigIntegerField(null=True)
    individual_net_amount = BigIntegerField(null=True)
    institution_net_qty = BigIntegerField(null=True)
    institution_net_amount = BigIntegerField(null=True)

    # 기타 투자자 (순매수량, 순매수금액만)
    securities_net_qty = BigIntegerField(null=True)
    securities_net_amount = BigIntegerField(null=True)
    investment_trust_net_qty = BigIntegerField(null=True)
    investment_trust_net_amount = BigIntegerField(null=True)
    pe_fund_net_qty = BigIntegerField(null=True)
    pe_fund_net_amount = BigIntegerField(null=True)
    bank_net_qty = BigIntegerField(null=True)
    bank_net_amount = BigIntegerField(null=True)
    insurance_net_qty = BigIntegerField(null=True)
    insurance_net_amount = BigIntegerField(null=True)
    pension_fund_net_qty = BigIntegerField(null=True)
    pension_fund_net_amount = BigIntegerField(null=True)

    created_at = DateTimeField(default=datetime.now)
    updated_at = DateTimeField(default=datetime.now)

    class Meta:
        table_name = 'market_investor_trend'
        indexes = (
            (('market_code', 'sector_code', 'trade_date'), True),
        )
```

### 3.2 API 응답 매핑

**파일**: `services/stock_collector.py`

```python
def _map_market_investor_trend(market_code: str, sector_code: str, data: Dict) -> Dict:
    """API 응답 -> DB 모델 매핑"""
    return {
        'market_code': market_code,
        'sector_code': sector_code,
        'collected_at': datetime.now(),
        'trade_date': date.today(),

        'foreign_net_qty': safe_int(data.get('frgn_ntby_qty')),
        'foreign_net_amount': safe_int(data.get('frgn_ntby_tr_pbmn')),
        'individual_net_qty': safe_int(data.get('prsn_ntby_qty')),
        'individual_net_amount': safe_int(data.get('prsn_ntby_tr_pbmn')),
        'institution_net_qty': safe_int(data.get('orgn_ntby_qty')),
        'institution_net_amount': safe_int(data.get('orgn_ntby_tr_pbmn')),

        'securities_net_qty': safe_int(data.get('scrt_ntby_qty')),
        'securities_net_amount': safe_int(data.get('scrt_ntby_tr_pbmn')),
        'investment_trust_net_qty': safe_int(data.get('ivtr_ntby_qty')),
        'investment_trust_net_amount': safe_int(data.get('ivtr_ntby_tr_pbmn')),
        'pe_fund_net_qty': safe_int(data.get('pe_fund_ntby_vol')),
        'pe_fund_net_amount': safe_int(data.get('pe_fund_ntby_tr_pbmn')),
        'bank_net_qty': safe_int(data.get('bank_ntby_qty')),
        'bank_net_amount': safe_int(data.get('bank_ntby_tr_pbmn')),
        'insurance_net_qty': safe_int(data.get('insu_ntby_qty')),
        'insurance_net_amount': safe_int(data.get('insu_ntby_tr_pbmn')),
        'pension_fund_net_qty': safe_int(data.get('fund_ntby_qty')),
        'pension_fund_net_amount': safe_int(data.get('fund_ntby_tr_pbmn')),
    }
```

### 3.3 배치 처리

**파일**: `main.py`

```python
def process_market_investor_trend():
    """시장별 투자자 매매동향 수집 배치"""
    targets = [
        ("KSP", "0001"),  # 코스피 종합
        ("KSQ", "1001"),  # 코스닥 종합
        ("ETF", "T000"),  # ETF 전체
    ]

    for market_code, sector_code in targets:
        data = collector.fetch_market_investor_trend(market_code, sector_code)
        with mysql_db.atomic():
            saver.upsert_market_investor_trend(data)
```

### 3.4 CLI 옵션

```python
parser.add_argument(
    '--market-investor-trend',
    action='store_true',
    help='시장별 투자자 매매동향(시세) 수집'
)
```

---

## Phase 4: 배포 및 통합 테스트

### 4.1 배포 순서

1. **korea_investment_stock** PR merge 후 PyPI 배포
   - PR을 master 브랜치에 merge
   - GitHub Actions에서 `release.yml` workflow 수동 실행
   - PyPI에 자동 배포됨
2. **moneyflow.advenoh.pe.kr** DB 마이그레이션 실행
   ```bash
   cd moneyflow.advenoh.pe.kr/backend/db
   ./liquibase.sh update
   ```
3. **stock-data-batch** 의존성 업데이트 및 배치 테스트
   ```bash
   python main.py --market-investor-trend
   ```
