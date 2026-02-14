from typing import Optional, Dict, Any
from ..korea_investment_stock import KoreaInvestment
from .cache_manager import CacheManager


class CachedKoreaInvestment:
    """캐싱 기능이 추가된 KoreaInvestment 래퍼"""

    DEFAULT_TTL = {
        'price': 5,           # 실시간 가격: 5초
        'stock_info': 300,    # 종목 정보: 5분
        'symbols': 3600,      # 종목 리스트: 1시간
        'ipo': 1800           # IPO 일정: 30분
    }

    def __init__(
        self,
        broker: KoreaInvestment,
        enable_cache: bool = True,
        price_ttl: Optional[int] = None,
        stock_info_ttl: Optional[int] = None,
        symbols_ttl: Optional[int] = None,
        ipo_ttl: Optional[int] = None
    ):
        """
        Args:
            broker: KoreaInvestment 인스턴스
            enable_cache: 캐싱 활성화 여부
            price_ttl: 실시간 가격 TTL (초)
            stock_info_ttl: 종목정보 TTL (초)
            symbols_ttl: 종목리스트 TTL (초)
            ipo_ttl: IPO 일정 TTL (초)
        """
        self.broker = broker
        self.enable_cache = enable_cache
        self.cache = CacheManager() if enable_cache else None

        # TTL 설정
        self.ttl = {
            'price': price_ttl or self.DEFAULT_TTL['price'],
            'stock_info': stock_info_ttl or self.DEFAULT_TTL['stock_info'],
            'symbols': symbols_ttl or self.DEFAULT_TTL['symbols'],
            'ipo': ipo_ttl or self.DEFAULT_TTL['ipo']
        }

    def _make_cache_key(self, method: str, *args, **kwargs) -> str:
        """캐시 키 생성"""
        args_str = "_".join(str(arg) for arg in args)
        kwargs_str = "_".join(f"{k}={v}" for k, v in sorted(kwargs.items()))
        return f"{method}:{args_str}:{kwargs_str}"

    def fetch_price(self, symbol: str, country_code: str = "KR") -> dict:
        """가격 조회 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_price(symbol, country_code)

        cache_key = self._make_cache_key("fetch_price", symbol, country_code)
        cached_data = self.cache.get(cache_key)

        if cached_data is not None:
            return cached_data

        result = self.broker.fetch_price(symbol, country_code)

        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['price'])

        return result

    def fetch_domestic_price(self, symbol: str, symbol_type: str = "Stock") -> dict:
        """국내 주식/ETF 가격 조회 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_domestic_price(symbol, symbol_type)

        cache_key = self._make_cache_key("fetch_domestic_price", symbol, symbol_type)
        cached_data = self.cache.get(cache_key)

        if cached_data is not None:
            return cached_data

        result = self.broker.fetch_domestic_price(symbol, symbol_type)

        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['price'])

        return result

    def fetch_price_detail_oversea(self, symbol: str, country_code: str = "US") -> dict:
        """해외 주식 가격 조회 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_price_detail_oversea(symbol, country_code)

        cache_key = self._make_cache_key("fetch_price_detail_oversea", symbol, country_code)
        cached_data = self.cache.get(cache_key)

        if cached_data is not None:
            return cached_data

        result = self.broker.fetch_price_detail_oversea(symbol, country_code)

        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['price'])

        return result

    def fetch_stock_info(self, symbol: str, country_code: str = "KR") -> dict:
        """종목 정보 조회 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_stock_info(symbol, country_code)

        cache_key = self._make_cache_key("fetch_stock_info", symbol, country_code)
        cached_data = self.cache.get(cache_key)

        if cached_data is not None:
            return cached_data

        result = self.broker.fetch_stock_info(symbol, country_code)

        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['stock_info'])

        return result

    def fetch_search_stock_info(self, symbol: str, country_code: str = "KR") -> dict:
        """종목 검색 (캐싱 지원) - 국내주식 전용 (KR만 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_search_stock_info(symbol, country_code)

        cache_key = self._make_cache_key("fetch_search_stock_info", symbol, country_code)
        cached_data = self.cache.get(cache_key)

        if cached_data is not None:
            return cached_data

        result = self.broker.fetch_search_stock_info(symbol, country_code)

        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['stock_info'])

        return result

    def fetch_ipo_schedule(
        self,
        from_date: Optional[str] = None,
        to_date: Optional[str] = None,
        symbol: str = ""
    ) -> dict:
        """IPO 일정 조회 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_ipo_schedule(from_date, to_date, symbol)

        cache_key = self._make_cache_key("fetch_ipo_schedule", from_date, to_date, symbol)
        cached_data = self.cache.get(cache_key)

        if cached_data is not None:
            return cached_data

        result = self.broker.fetch_ipo_schedule(from_date, to_date, symbol)

        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['ipo'])

        return result

    def fetch_investor_trading_by_stock_daily(
        self,
        symbol: str,
        date: str,
        market_code: str = "J"
    ) -> dict:
        """투자자 매매동향 조회 (캐싱 지원)

        과거 날짜 데이터는 1시간 캐시, 당일 데이터는 가격 TTL 적용
        """
        if not self.enable_cache:
            return self.broker.fetch_investor_trading_by_stock_daily(symbol, date, market_code)

        cache_key = self._make_cache_key("fetch_investor_trading_by_stock_daily", symbol, date, market_code)
        cached_data = self.cache.get(cache_key)

        if cached_data is not None:
            return cached_data

        result = self.broker.fetch_investor_trading_by_stock_daily(symbol, date, market_code)

        if result.get('rt_cd') == '0':
            # 과거 데이터는 더 긴 TTL 적용 (1시간)
            from datetime import datetime
            today = datetime.now().strftime("%Y%m%d")
            ttl = 3600 if date < today else self.ttl['price']
            self.cache.set(cache_key, result, ttl)

        return result

    def fetch_domestic_chart(
        self,
        symbol: str,
        period: str = "D",
        start_date: str = "",
        end_date: str = "",
        adjusted: bool = True,
        market_code: str = "J"
    ) -> dict:
        """국내주식기간별시세 조회 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_domestic_chart(symbol, period, start_date, end_date, adjusted, market_code)

        cache_key = self._make_cache_key("fetch_domestic_chart", symbol, period, start_date, end_date, adjusted, market_code)
        cached_data = self.cache.get(cache_key)

        if cached_data is not None:
            return cached_data

        result = self.broker.fetch_domestic_chart(symbol, period, start_date, end_date, adjusted, market_code)

        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['price'])

        return result

    def fetch_domestic_minute_chart(
        self,
        symbol: str,
        time_from: str = "",
        market_code: str = "J"
    ) -> dict:
        """주식당일분봉조회 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_domestic_minute_chart(symbol, time_from, market_code)

        cache_key = self._make_cache_key("fetch_domestic_minute_chart", symbol, time_from, market_code)
        cached_data = self.cache.get(cache_key)

        if cached_data is not None:
            return cached_data

        result = self.broker.fetch_domestic_minute_chart(symbol, time_from, market_code)

        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['price'])

        return result

    def fetch_overseas_chart(
        self,
        symbol: str,
        country_code: str = "US",
        period: str = "D",
        end_date: str = "",
        adjusted: bool = True
    ) -> dict:
        """해외주식 기간별시세 조회 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_overseas_chart(symbol, country_code, period, end_date, adjusted)

        cache_key = self._make_cache_key("fetch_overseas_chart", symbol, country_code, period, end_date, adjusted)
        cached_data = self.cache.get(cache_key)

        if cached_data is not None:
            return cached_data

        result = self.broker.fetch_overseas_chart(symbol, country_code, period, end_date, adjusted)

        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['price'])

        return result

    def fetch_volume_ranking(self, market_code: str = "J", sort_by: str = "0") -> dict:
        """거래량순위 조회 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_volume_ranking(market_code, sort_by)
        cache_key = self._make_cache_key("fetch_volume_ranking", market_code, sort_by)
        cached_data = self.cache.get(cache_key)
        if cached_data is not None:
            return cached_data
        result = self.broker.fetch_volume_ranking(market_code, sort_by)
        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['price'])
        return result

    def fetch_change_rate_ranking(self, market_code: str = "J", sort_order: str = "0", period_days: str = "0") -> dict:
        """등락률 순위 조회 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_change_rate_ranking(market_code, sort_order, period_days)
        cache_key = self._make_cache_key("fetch_change_rate_ranking", market_code, sort_order, period_days)
        cached_data = self.cache.get(cache_key)
        if cached_data is not None:
            return cached_data
        result = self.broker.fetch_change_rate_ranking(market_code, sort_order, period_days)
        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['price'])
        return result

    def fetch_market_cap_ranking(self, market_code: str = "J", target_market: str = "0000") -> dict:
        """시가총액 상위 조회 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_market_cap_ranking(market_code, target_market)
        cache_key = self._make_cache_key("fetch_market_cap_ranking", market_code, target_market)
        cached_data = self.cache.get(cache_key)
        if cached_data is not None:
            return cached_data
        result = self.broker.fetch_market_cap_ranking(market_code, target_market)
        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['price'])
        return result

    def fetch_overseas_change_rate_ranking(self, country_code: str = "US", sort_order: str = "1", period: str = "0", volume_filter: str = "0") -> dict:
        """해외주식 상승율/하락율 조회 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_overseas_change_rate_ranking(country_code, sort_order, period, volume_filter)
        cache_key = self._make_cache_key("fetch_overseas_change_rate_ranking", country_code, sort_order, period, volume_filter)
        cached_data = self.cache.get(cache_key)
        if cached_data is not None:
            return cached_data
        result = self.broker.fetch_overseas_change_rate_ranking(country_code, sort_order, period, volume_filter)
        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['price'])
        return result

    def fetch_financial_ratio(self, symbol: str, period_type: str = "0", market_code: str = "J") -> dict:
        """재무비율 조회 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_financial_ratio(symbol, period_type, market_code)
        cache_key = self._make_cache_key("fetch_financial_ratio", symbol, period_type, market_code)
        cached_data = self.cache.get(cache_key)
        if cached_data is not None:
            return cached_data
        result = self.broker.fetch_financial_ratio(symbol, period_type, market_code)
        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['stock_info'])
        return result

    def fetch_income_statement(self, symbol: str, period_type: str = "0", market_code: str = "J") -> dict:
        """손익계산서 조회 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_income_statement(symbol, period_type, market_code)
        cache_key = self._make_cache_key("fetch_income_statement", symbol, period_type, market_code)
        cached_data = self.cache.get(cache_key)
        if cached_data is not None:
            return cached_data
        result = self.broker.fetch_income_statement(symbol, period_type, market_code)
        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['stock_info'])
        return result

    def fetch_balance_sheet(self, symbol: str, period_type: str = "0", market_code: str = "J") -> dict:
        """대차대조표 조회 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_balance_sheet(symbol, period_type, market_code)
        cache_key = self._make_cache_key("fetch_balance_sheet", symbol, period_type, market_code)
        cached_data = self.cache.get(cache_key)
        if cached_data is not None:
            return cached_data
        result = self.broker.fetch_balance_sheet(symbol, period_type, market_code)
        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['stock_info'])
        return result

    def fetch_profitability_ratio(self, symbol: str, period_type: str = "0", market_code: str = "J") -> dict:
        """수익성비율 조회 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_profitability_ratio(symbol, period_type, market_code)
        cache_key = self._make_cache_key("fetch_profitability_ratio", symbol, period_type, market_code)
        cached_data = self.cache.get(cache_key)
        if cached_data is not None:
            return cached_data
        result = self.broker.fetch_profitability_ratio(symbol, period_type, market_code)
        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['stock_info'])
        return result

    def fetch_growth_ratio(self, symbol: str, period_type: str = "0", market_code: str = "J") -> dict:
        """성장성비율 조회 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_growth_ratio(symbol, period_type, market_code)
        cache_key = self._make_cache_key("fetch_growth_ratio", symbol, period_type, market_code)
        cached_data = self.cache.get(cache_key)
        if cached_data is not None:
            return cached_data
        result = self.broker.fetch_growth_ratio(symbol, period_type, market_code)
        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['stock_info'])
        return result

    def fetch_dividend_ranking(self, market_type: str = "0", dividend_type: str = "2", start_date: str = "", end_date: str = "", settlement_type: str = "0") -> dict:
        """배당률 상위 조회 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_dividend_ranking(market_type, dividend_type, start_date, end_date, settlement_type)
        cache_key = self._make_cache_key("fetch_dividend_ranking", market_type, dividend_type, start_date, end_date, settlement_type)
        cached_data = self.cache.get(cache_key)
        if cached_data is not None:
            return cached_data
        result = self.broker.fetch_dividend_ranking(market_type, dividend_type, start_date, end_date, settlement_type)
        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['stock_info'])
        return result

    def fetch_industry_index(self, industry_code: str = "0001") -> dict:
        """업종 현재지수 조회 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_industry_index(industry_code)
        cache_key = self._make_cache_key("fetch_industry_index", industry_code)
        cached_data = self.cache.get(cache_key)
        if cached_data is not None:
            return cached_data
        result = self.broker.fetch_industry_index(industry_code)
        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['stock_info'])
        return result

    def fetch_industry_category_price(self, market_type: str = "K", category_code: str = "0") -> dict:
        """업종 구분별전체시세 조회 (캐싱 지원)"""
        if not self.enable_cache:
            return self.broker.fetch_industry_category_price(market_type, category_code)
        cache_key = self._make_cache_key("fetch_industry_category_price", market_type, category_code)
        cached_data = self.cache.get(cache_key)
        if cached_data is not None:
            return cached_data
        result = self.broker.fetch_industry_category_price(market_type, category_code)
        if result.get('rt_cd') == '0':
            self.cache.set(cache_key, result, self.ttl['stock_info'])
        return result

    def invalidate_cache(self, method: Optional[str] = None):
        """캐시 무효화"""
        if not self.enable_cache:
            return

        self.cache.clear()

    def get_cache_stats(self) -> Dict[str, Any]:
        """캐시 통계 반환"""
        if not self.enable_cache:
            return {'cache_enabled': False}

        stats = self.cache.get_stats()
        stats['cache_enabled'] = True
        stats['ttl_config'] = self.ttl
        return stats

    def __enter__(self):
        """컨텍스트 매니저 진입"""
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        """컨텍스트 매니저 종료"""
        if self.enable_cache:
            self.cache.clear()
        return False
