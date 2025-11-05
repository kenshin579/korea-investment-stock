# 🚀 Korea Investment Stock

[![PyPI version](https://badge.fury.io/py/korea-investment-stock.svg)](https://badge.fury.io/py/korea-investment-stock)
[![Python Versions](https://img.shields.io/pypi/pyversions/korea-investment-stock.svg)](https://pypi.org/project/korea-investment-stock/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Pure Python wrapper** for Korea Investment Securities OpenAPI

## 🎯 Purpose

A simple, transparent wrapper around the Korea Investment Securities OpenAPI. This library handles API authentication and request formatting, giving you direct access to the API responses without abstraction layers.

### Philosophy: Keep It Simple

- **Pure wrapper**: Direct API access without magic
- **Minimal dependencies**: Only `requests` and `pandas`
- **No abstraction**: You get exactly what the API returns
- **Implement your way**: Add rate limiting, caching, retries as you need them

## 🌟 Features

### Core API Support
- ✅ **Stock Price Queries**: Domestic (KR) and US stocks
- ✅ **Stock Information**: Company details and market data
- ✅ **IPO Schedule**: Public offering information and schedules
- ✅ **Unified Interface**: Query KR/US stocks with `fetch_price(symbol, market)`
- ✅ **Search Functions**: Stock search and lookup

### Technical Features
- 🔧 **Context Manager**: Automatic resource cleanup
- 🔧 **Thread Pool**: Basic concurrent execution support
- 🔧 **Mock Trading**: Test mode for development
- 🔧 **Environment Variables**: API credentials via env vars

## 📦 Installation

```bash
pip install korea-investment-stock
```

### Requirements
- Python 3.11 or higher
- Korea Investment Securities API account

## 🚀 Quick Start

### 1. Set Up Credentials

```bash
# Add to your ~/.zshrc or ~/.bashrc
export KOREA_INVESTMENT_API_KEY="your-api-key"
export KOREA_INVESTMENT_API_SECRET="your-api-secret"
export KOREA_INVESTMENT_ACCOUNT_NO="12345678-01"
```

### 2. Basic Usage

```python
from korea_investment_stock import KoreaInvestment
import os

# Create client (using environment variables)
with KoreaInvestment(
    api_key=os.getenv('KOREA_INVESTMENT_API_KEY'),
    api_secret=os.getenv('KOREA_INVESTMENT_API_SECRET'),
    acc_no=os.getenv('KOREA_INVESTMENT_ACCOUNT_NO'),
    mock=True  # False for real trading
) as broker:
    # Query Samsung Electronics
    result = broker.fetch_price("005930", "KR")

    if result['rt_cd'] == '0':
        price = result['output1']['stck_prpr']
        print(f"Price: {price}원")
```

## 📖 API Methods

### Stock Price Queries

```python
# Domestic stock price
result = broker.fetch_price("005930", "KR")  # Samsung Electronics

# US stock price (requires real account)
result = broker.fetch_price("AAPL", "US")   # Apple

# Direct methods
result = broker.fetch_domestic_price("J", "005930")
result = broker.fetch_etf_domestic_price("J", "069500")  # KODEX 200
result = broker.fetch_price_detail_oversea("AAPL", "US")
```

### Stock Information

```python
# Stock info
result = broker.fetch_stock_info("005930", "KR")

# Search stock
result = broker.fetch_search_stock_info("005930", "KR")
```

### IPO Schedule

```python
# All IPOs (today + 30 days)
result = broker.fetch_ipo_schedule()

# Specific period
result = broker.fetch_ipo_schedule(
    from_date="20250101",
    to_date="20250131"
)

# Specific symbol
result = broker.fetch_ipo_schedule(symbol="123456")

# Helper methods
status = broker.get_ipo_status(ipo['subscr_dt'])  # "예정", "진행중", "마감"
d_day = broker.calculate_ipo_d_day(ipo['subscr_dt'])  # Days until subscription
```

### Symbol Lists

```python
# KOSPI symbols
result = broker.fetch_kospi_symbols()

# KOSDAQ symbols
result = broker.fetch_kosdaq_symbols()
```

## 🔧 Advanced Usage

### Multiple Stock Queries

```python
# Query multiple stocks
stocks = [
    ("005930", "KR"),  # Samsung
    ("000660", "KR"),  # SK Hynix
    ("035720", "KR"),  # Kakao
]

results = []
for symbol, market in stocks:
    result = broker.fetch_price(symbol, market)
    results.append(result)
    # Add your own rate limiting here if needed
```

### Mixed KR/US Portfolio

```python
portfolio = [
    ("005930", "KR"),  # Samsung Electronics
    ("AAPL", "US"),    # Apple
    ("035720", "KR"),  # Kakao
    ("MSFT", "US"),    # Microsoft
]

for symbol, market in portfolio:
    result = broker.fetch_price(symbol, market)

    if result['rt_cd'] == '0':
        if market == "KR":
            output = result['output1']
            price = output['stck_prpr']
            print(f"{symbol}: ₩{int(price):,}")
        else:
            output = result['output']
            price = output['last']
            print(f"{symbol}: ${price}")
```

### Error Handling

```python
try:
    result = broker.fetch_price("INVALID", "US")

    if result['rt_cd'] != '0':
        # API returned error
        print(f"API Error: {result['msg1']}")
except ValueError as e:
    # Invalid parameters
    print(f"Invalid request: {e}")
except Exception as e:
    # Network or other errors
    print(f"Error: {e}")
```

### Memory Caching (Optional)

Reduce API calls and improve response times with built-in memory caching:

```python
from korea_investment_stock import KoreaInvestment, CachedKoreaInvestment

# Create base broker
broker = KoreaInvestment(api_key, api_secret, acc_no, mock=True)

# Wrap with caching (opt-in)
cached_broker = CachedKoreaInvestment(broker, price_ttl=5)

# First call: API request (cache miss)
result1 = cached_broker.fetch_price("005930", "KR")  # ~200ms

# Second call: from cache (cache hit)
result2 = cached_broker.fetch_price("005930", "KR")  # <1ms

# Cache statistics
stats = cached_broker.get_cache_stats()
print(f"Hit rate: {stats['hit_rate']}")  # "50.00%"
```

**TTL Configuration** (in seconds):

```python
cached_broker = CachedKoreaInvestment(
    broker,
    price_ttl=5,        # Real-time price: 5 seconds
    stock_info_ttl=300, # Stock info: 5 minutes
    symbols_ttl=3600,   # Symbol lists: 1 hour
    ipo_ttl=1800        # IPO schedule: 30 minutes
)
```

**Performance Benefits**:
- 📉 API calls reduced by 30-50%
- ⚡ Response time improved by 90%+ (cached queries)
- 🔒 Thread-safe with automatic expiration
- 💾 No external dependencies (memory-only)

**See**: `examples/cached_basic_example.py` for comprehensive examples

## 📊 Response Format

### Domestic Stock (KR)

```python
{
    'rt_cd': '0',               # Return code ('0' = success)
    'msg1': '정상처리되었습니다',   # Message
    'output1': {
        'stck_prpr': '62600',   # Current price
        'prdy_vrss': '1600',    # Change from previous day
        'prdy_ctrt': '2.62',    # Change rate (%)
        'stck_oprc': '61000',   # Opening price
        'stck_hgpr': '63000',   # High price
        'stck_lwpr': '60500',   # Low price
        'acml_vol': '15234567'  # Volume
        # ... more fields
    }
}
```

### US Stock (US)

```python
{
    'rt_cd': '0',
    'msg1': '정상처리되었습니다',
    'output': {
        'rsym': 'DNASAAPL',     # Exchange + Symbol
        'last': '211.16',       # Current price
        'open': '210.56',       # Opening price
        'high': '212.13',       # High price
        'low': '209.86',        # Low price
        'tvol': '39765812',     # Volume
        't_xdif': '1.72',       # Change
        't_xrat': '-0.59',      # Change rate (%)
        'perx': '32.95',        # PER
        'pbrx': '47.23',        # PBR
        'epsx': '6.41',         # EPS
        'bpsx': '4.47'          # BPS
        # ... more fields
    }
}
```

## ⚠️ Important Notes

### API Rate Limits
- Korea Investment API: **20 requests/second**
- **You are responsible** for implementing rate limiting
- Exceeding limits will cause API errors

### US Stocks
- **Real account only** (mock trading not supported)
- Auto-detects exchange (NASDAQ → NYSE → AMEX)
- Includes financial ratios (PER, PBR, EPS, BPS)

### Context Manager
Always use context manager for proper resource cleanup:

```python
# ✅ Good: Automatic cleanup
with KoreaInvestment(api_key, api_secret, acc_no) as broker:
    result = broker.fetch_price("005930", "KR")

# ❌ Bad: Manual cleanup required
broker = KoreaInvestment(api_key, api_secret, acc_no)
result = broker.fetch_price("005930", "KR")
broker.shutdown()  # Must call manually
```

## 🔨 Implementing Your Own Features

### Rate Limiting Example

```python
import time

class RateLimiter:
    def __init__(self, calls_per_second=15):
        self.min_interval = 1.0 / calls_per_second
        self.last_call = 0

    def wait(self):
        elapsed = time.time() - self.last_call
        if elapsed < self.min_interval:
            time.sleep(self.min_interval - elapsed)
        self.last_call = time.time()

# Usage
limiter = RateLimiter(calls_per_second=15)

for symbol, market in stocks:
    limiter.wait()
    result = broker.fetch_price(symbol, market)
```

### Caching Example

```python
from functools import lru_cache
from datetime import datetime, timedelta

class CachedBroker:
    def __init__(self, broker):
        self.broker = broker
        self.cache = {}
        self.ttl = timedelta(minutes=5)

    def fetch_price_cached(self, symbol, market):
        key = f"{symbol}:{market}"
        now = datetime.now()

        if key in self.cache:
            cached_time, cached_result = self.cache[key]
            if now - cached_time < self.ttl:
                return cached_result

        result = self.broker.fetch_price(symbol, market)
        self.cache[key] = (now, result)
        return result
```

### Retry Example

```python
import time

def fetch_with_retry(broker, symbol, market, max_retries=3):
    for attempt in range(max_retries):
        try:
            result = broker.fetch_price(symbol, market)
            if result['rt_cd'] == '0':
                return result
        except Exception as e:
            if attempt == max_retries - 1:
                raise
            time.sleep(2 ** attempt)  # Exponential backoff
```

## 📚 Examples

See the `examples/` directory:
- `basic_example.py`: Simple usage patterns
- `ipo_schedule_example.py`: IPO queries and helpers
- `us_stock_price_example.py`: US stock queries

## 🔄 Migration from v0.5.0

### Breaking Changes in v0.6.0

**Removed features** (~6,000 lines of code):
- Rate limiting system
- TTL caching
- Batch processing methods
- Monitoring and statistics
- Visualization tools
- Automatic retry decorators

**API changes**:
```python
# v0.5.0 (Old)
results = broker.fetch_price_list([("005930", "KR"), ("AAPL", "US")])

# v0.6.0 (New)
results = []
for symbol, market in [("005930", "KR"), ("AAPL", "US")]:
    result = broker.fetch_price(symbol, market)
    results.append(result)
```

See [CHANGELOG.md](CHANGELOG.md) for complete migration guide.

## 📖 Documentation

- [Official API Docs](https://wikidocs.net/book/7845)
- [GitHub Issues](https://github.com/kenshin579/korea-investment-stock/issues)

## 🤝 Contributing

Contributions welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Submit a pull request

## 📝 License

MIT License - see [LICENSE](LICENSE) file

## ⚡ Performance Tips

1. **Implement rate limiting**: Prevent API errors
2. **Use caching**: Reduce redundant API calls
3. **Batch requests wisely**: Don't overwhelm the API
4. **Handle errors gracefully**: Check `rt_cd` and `msg1`
5. **Use context manager**: Ensure proper cleanup

## 🙏 Credits

- Korea Investment Securities for providing the OpenAPI
- Original contributors: Jonghun Yoo, Brayden Jo, Frank Oh

---

**Remember**: This is a pure wrapper. You control rate limiting, caching, error handling, and monitoring according to your needs.
