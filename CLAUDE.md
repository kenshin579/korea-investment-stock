# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Korea Investment Stock is a **pure Python wrapper** for the Korea Investment Securities OpenAPI. This library focuses on providing direct, transparent access to the API without abstraction layers.

**Philosophy**: Simple, transparent, and flexible - let users implement features their way.

**Key Capabilities:**
- Domestic (KR) and US stock price/info queries
- Stock information and search
- IPO schedule lookup
- Unified interface for mixed KR/US stock queries via `fetch_price(symbol, market)`
- Context manager support for resource cleanup

## Development Commands

### Environment Setup

```bash
# Create and activate virtual environment (.venv is required)
python -m venv .venv
source .venv/bin/activate  # macOS/Linux
# .venv\Scripts\activate  # Windows

# Install package in editable mode (uses pyproject.toml, NOT requirements.txt)
pip install -e .

# Install with development dependencies
pip install -e ".[dev]"
```

### Testing

```bash
# Run all tests
pytest

# Run specific test file
pytest korea_investment_stock/tests/test_korea_investment_stock.py

# Run with verbose output
pytest -v

# Run integration tests (requires API credentials)
pytest korea_investment_stock/tests/test_integration_us_stocks.py -v
pytest korea_investment_stock/tests/test_ipo_integration.py -v
```

### Running Examples

```bash
# Always activate virtual environment first
source .venv/bin/activate

# Run examples
python examples/basic_example.py
python examples/ipo_schedule_example.py
python examples/us_stock_price_example.py
```

### Package Management

```bash
# Build distribution packages
python -m build

# Upload to PyPI (maintainers only)
./upload.sh
```

## Required Environment Variables

**CRITICAL:** This project uses OS environment variables only. Never use `.env` files or `python-dotenv`.

Set these in your shell profile (`~/.zshrc` or `~/.bashrc`):

```bash
export KOREA_INVESTMENT_API_KEY="your-api-key"
export KOREA_INVESTMENT_API_SECRET="your-api-secret"
export KOREA_INVESTMENT_ACCOUNT_NO="12345678-01"
```

**Naming Convention:**
- Always use `KOREA_INVESTMENT_` prefix
- Account number is `ACCOUNT_NO` (not `ACC_NO`)
- All uppercase with underscore separators

## Architecture Overview

### Simplified Component Flow

```
User API Call
  ↓
KoreaInvestment.fetch_price(symbol, market)
  ↓
HTTP Request to Korea Investment API
  ↓
Return raw API response
```

**That's it.** No decorators, no caching, no rate limiting, no magic.

### Core Module

**`korea_investment_stock/korea_investment_stock.py`** - Main class (1,011 lines)
- `KoreaInvestment`: Primary API interface
- Context manager pattern (`__enter__`, `__exit__`)
- Token management (`issue_access_token()`)
- Public API methods (18 total):
  - `fetch_price(symbol, market)` - Unified KR/US price query
  - `fetch_domestic_price(market_code, symbol)` - KR stocks
  - `fetch_etf_domestic_price(market_code, symbol)` - KR ETFs
  - `fetch_price_detail_oversea(symbol, market)` - US stocks
  - `fetch_stock_info(symbol, market)` - Stock information
  - `fetch_search_stock_info(symbol, market)` - Stock search
  - `fetch_kospi_symbols()` - KOSPI symbol list
  - `fetch_kosdaq_symbols()` - KOSDAQ symbol list
  - `fetch_ipo_schedule()` - IPO schedule
  - IPO helper methods (9 total)

### Package Structure

```
korea_investment_stock/
├── __init__.py (18 lines)
├── korea_investment_stock.py (1,011 lines)
└── tests/
    ├── test_korea_investment_stock.py
    ├── test_integration_us_stocks.py
    ├── test_ipo_schedule.py
    └── test_ipo_integration.py
```

**Dependencies:** `requests`, `pandas` (minimal)

## Code Style & Conventions

**From `.cursorrules`:**

1. **Python Version**: 3.11+ (uses `zoneinfo`, modern type hints)
2. **Type Hints**: Required for all public methods
3. **Comments**: Prefer Korean for domain-specific comments
4. **Error Messages**: Korean for user-facing messages
5. **Dependency Management**: Use `pyproject.toml` only (no `requirements.txt`)

## Common Development Patterns

### Adding a New API Method

1. Define method in `KoreaInvestment` class
2. **No decorators** - just pure API calls
3. Add tests in `korea_investment_stock/tests/test_korea_investment_stock.py`

```python
def fetch_new_data(self, symbol: str) -> Dict[str, Any]:
    """
    새로운 데이터 조회

    Args:
        symbol: 종목 코드

    Returns:
        API 응답 딕셔너리
    """
    url = f"{self.base_url}/endpoint"
    headers = {
        "authorization": self.access_token,
        "appkey": self.api_key,
        "appsecret": self.api_secret,
        "tr_id": "TRANSACTION_ID"
    }
    params = {"symbol": symbol}

    response = requests.get(url, headers=headers, params=params)
    return response.json()
```

### Error Handling Pattern

API returns error codes in response:
```python
result = broker.fetch_price("005930", "KR")

if result['rt_cd'] == '0':
    # Success
    price = result['output1']['stck_prpr']
else:
    # Error
    print(f"Error: {result['msg1']}")
```

Users should implement their own retry logic:
```python
import time

def fetch_with_retry(broker, symbol, market, max_retries=3):
    for attempt in range(max_retries):
        try:
            result = broker.fetch_price(symbol, market)
            if result['rt_cd'] == '0':
                return result
            time.sleep(2 ** attempt)  # Exponential backoff
        except Exception as e:
            if attempt == max_retries - 1:
                raise
            time.sleep(2 ** attempt)
```

### Working with US Stocks

US stock queries require real account (mock not supported):
- Exchange auto-detection: NASDAQ → NYSE → AMEX
- Symbol format: Plain ticker (e.g., "AAPL", not "AAPL.US")
- Additional fields: PER, PBR, EPS, BPS included in response

Example:
```python
# Requires mock=False (real account)
with KoreaInvestment(api_key, api_secret, acc_no, mock=False) as broker:
    result = broker.fetch_price("AAPL", "US")

    if result['rt_cd'] == '0':
        output = result['output']
        print(f"Price: ${output['last']}")
        print(f"PER: {output['perx']}")
```

## Testing Strategy

### Test Organization

- **Unit tests**: `test_korea_investment_stock.py` (94 lines)
- **Integration tests**:
  - `test_integration_us_stocks.py` (225 lines) - US stock queries
  - `test_ipo_integration.py` (186 lines) - IPO schedule
- **Feature tests**: `test_ipo_schedule.py` (297 lines) - IPO helpers

### Running Integration Tests

Integration tests require valid API credentials:

```bash
# Set credentials in environment
export KOREA_INVESTMENT_API_KEY="..."
export KOREA_INVESTMENT_API_SECRET="..."
export KOREA_INVESTMENT_ACCOUNT_NO="..."

# Run integration tests
pytest korea_investment_stock/tests/test_integration_us_stocks.py -v
pytest korea_investment_stock/tests/test_ipo_integration.py -v
```

## Important Files

- **`pyproject.toml`**: Package metadata, dependencies, build config
- **`.cursorrules`**: Development conventions (env vars, virtual env, naming)
- **`CHANGELOG.md`**: Version history and release notes (v0.6.0 = breaking changes)
- **`examples/`**: Usage examples:
  - `basic_example.py` - Getting started
  - `ipo_schedule_example.py` - IPO queries
  - `us_stock_price_example.py` - US stock queries

## API Rate Limiting

**User Responsibility**: You must implement your own rate limiting.

Korea Investment Securities API limits:
- **Official**: 20 requests/second
- **Recommended**: 15 requests/second (conservative)

Simple rate limiting example:
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

## Known Limitations

1. **US stocks**: Real account only (mock not supported by API)
2. **IPO data**: Reference only (from 예탁원), mock not supported
3. **Order functionality**: Not yet implemented (planned)
4. **WebSocket**: Not included in this library
5. **No built-in rate limiting**: Users must implement
6. **No built-in caching**: Users must implement
7. **No automatic retry**: Users must implement

## Git Commit Message Guidelines

Follow conventional commit format:
- `[feat]`: New feature
- `[fix]`: Bug fix
- `[chore]`: Maintenance (docs, CI, etc.)
- `[refactor]`: Code restructuring
- `[test]`: Test additions/changes

Example: `[feat] Add US stock PER/PBR data to fetch_price`

## GitHub Actions Workflows

- **`label-merge-conflict.yml`**: Auto-labels PRs with merge conflicts
- Additional workflows may be present for CI/CD

## Context Manager Pattern

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

## User Implementation Examples

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

### Batch Processing Example

```python
def fetch_multiple_stocks(broker, stock_list, rate_limit=15):
    """Fetch multiple stocks with rate limiting"""
    import time

    min_interval = 1.0 / rate_limit
    last_call = 0
    results = []

    for symbol, market in stock_list:
        # Rate limiting
        elapsed = time.time() - last_call
        if elapsed < min_interval:
            time.sleep(min_interval - elapsed)
        last_call = time.time()

        # API call
        result = broker.fetch_price(symbol, market)
        results.append(result)

    return results
```

## Migration from v0.5.0

**Breaking Changes in v0.6.0:**

1. **Removed Methods:**
   - `fetch_price_list()` → Use loop with `fetch_price()`
   - `fetch_stock_info_list()` → Use loop with `fetch_stock_info()`
   - All batch processing methods
   - All caching methods
   - All monitoring methods

2. **Removed Features:**
   - Rate limiting system
   - TTL caching
   - Statistics collection
   - Visualization tools
   - Automatic retry decorators

3. **Migration Example:**
```python
# v0.5.0 (Old)
results = broker.fetch_price_list([("005930", "KR"), ("AAPL", "US")])

# v0.6.0 (New)
results = []
for symbol, market in [("005930", "KR"), ("AAPL", "US")]:
    result = broker.fetch_price(symbol, market)
    results.append(result)
    # Add your own rate limiting here if needed
```

See [CHANGELOG.md](CHANGELOG.md) for complete details.

## Additional Resources

- **Official API Docs**: https://wikidocs.net/book/7845
- **GitHub Issues**: https://github.com/kenshin579/korea-investment-stock/issues
- **PyPI**: https://pypi.org/project/korea-investment-stock/
- **CHANGELOG**: [CHANGELOG.md](CHANGELOG.md) - v0.6.0 breaking changes

---

**Remember**: This is a pure wrapper. You control rate limiting, caching, error handling, and monitoring according to your needs.
