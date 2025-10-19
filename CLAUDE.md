# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Korea Investment Stock is a Python library providing a production-ready wrapper for the Korea Investment Securities OpenAPI. It features advanced rate limiting, automatic retry, batch processing, TTL caching, error recovery, and real-time monitoring.

**Key Capabilities:**
- Domestic (KR) and US stock price/info queries
- IPO schedule lookup
- Unified interface for mixed KR/US stock queries
- Production-tested reliability with 0% error rate under normal operation

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
pytest korea_investment_stock/tests/test_rate_limiter.py

# Run with verbose output
pytest -v

# Run tests matching a pattern
pytest -k "test_cache"

# Run integration tests (requires API credentials)
pytest korea_investment_stock/tests/test_integration.py
```

### Running Examples

```bash
# Always activate virtual environment first
source .venv/bin/activate

# Run examples
python examples/rate_limiting_example.py
python examples/ipo_schedule_example.py
python examples/us_stock_price_example.py

# Run stress test
python scripts/stress_runner.py
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

### Core Component Flow

```
User API Call
  ↓
@retry_on_rate_limit decorator (5 retries)
  ↓
@cacheable decorator (TTL-based)
  ↓
EnhancedRateLimiter.acquire() (Token Bucket + Sliding Window)
  ↓
HTTP Request to Korea Investment API
  ↓
Error Recovery System (pattern matching)
  ↓
Circuit Breaker (CLOSED → OPEN → HALF_OPEN)
  ↓
Exponential Backoff with Jitter
```

### Key Modules

1. **`korea_investment_stock/korea_investment_stock.py`** - Main class
   - `KoreaInvestment`: Primary API interface (context manager pattern)
   - Thread pool executor (max 3 concurrent workers)
   - Semaphore-based concurrency control
   - Unified KR/US stock query methods

2. **`rate_limiting/`** - Rate limiting system
   - `EnhancedRateLimiter`: Hybrid Token Bucket + Sliding Window (15 calls/sec, 80% safety margin)
   - `EnhancedBackoffStrategy`: Circuit breaker with exponential backoff + jitter (singleton)
   - `enhanced_retry_decorator`: `@retry_on_rate_limit`, `@retry_on_network_error`, `@auto_refresh_token`

3. **`caching/`** - TTL cache system
   - `TTLCache`: Thread-safe LRU/LFU cache with auto-expiry and compression
   - Background cleanup thread (60-second intervals)
   - `@cacheable`: Method-level caching decorator with custom key generation
   - Market-aware TTL (price: 5min, stock info: 5hrs, symbols: 3 days)

4. **`error_handling/`** - Error recovery
   - `ErrorRecoverySystem`: Pattern-based error matching (singleton)
   - Error severity levels: LOW, MEDIUM, HIGH, CRITICAL
   - Recovery actions: RETRY, WAIT, REFRESH_TOKEN, FAIL_FAST, CIRCUIT_BREAK
   - Tracks last 1000 errors for diagnostics

5. **`batch_processing/`** - Dynamic batch control
   - `DynamicBatchController`: Auto-adjusts batch size (5-100) and delay (0.5-5.0s)
   - Target error rate: 1%
   - Adjustment factor: 20% per cycle (min 10-second intervals)

6. **`monitoring/`** - Statistics management
   - `StatsManager`: Aggregates metrics from all subsystems (singleton)
   - Export formats: JSON, CSV, JSONL (with optional gzip)
   - 7-day auto-rotation
   - System health classification: HEALTHY (<1%), WARNING (<5%), CRITICAL (≥5%)

7. **`visualization/`** - Monitoring dashboards
   - `PlotlyVisualizer`: Interactive charts (requires plotly)
   - `DashboardManager`: Real-time monitoring HTML dashboards
   - Error rate trends, API usage graphs, system health metrics

### Singleton Pattern Usage

These components are global singletons shared across all `KoreaInvestment` instances:

```python
from korea_investment_stock.rate_limiting import get_backoff_strategy
from korea_investment_stock.error_handling import get_error_recovery_system
from korea_investment_stock.monitoring import get_stats_manager

backoff = get_backoff_strategy()       # Circuit breaker state
recovery = get_error_recovery_system() # Error pattern matching
stats = get_stats_manager()            # Metrics aggregation
```

### Threading & Concurrency

- **ThreadPoolExecutor**: 3 worker threads (configurable)
- **Semaphore**: Max 3 concurrent API calls
- **Locks**:
  - `threading.Lock` on RateLimiter, BackoffStrategy, ErrorRecoverySystem
  - `threading.RLock` on TTLCache (re-entrant)
- **Background Threads**:
  - Cache cleanup (daemon, 60s interval)
  - Optional rate limiter auto-save

All shared state is protected by locks. Thread-safe by design.

### Batch Processing Workflow

When calling `fetch_price_list()` or `fetch_price_list_with_dynamic_batch()`:

1. Split stock list into batches (default 50 items)
2. Submit items to ThreadPoolExecutor with 10ms delay between submissions (prevents burst)
3. Each item: rate_limit acquire → cache check → API call
4. DynamicBatchController monitors error rate
5. Auto-adjust batch size/delay if error rate > 1%
6. Next batch uses adjusted parameters

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
2. Apply decorators in order:
   ```python
   @retry_on_rate_limit(max_retries=5)
   @cacheable(ttl=300, key_generator=lambda self, symbol: f"price:{symbol}")
   def fetch_new_data(self, symbol: str) -> Dict[str, Any]:
       with self.rate_limiter.acquire():
           response = self._call("/endpoint", params={...})
           return self._parse_response(response)
   ```
3. Add tests in `korea_investment_stock/tests/test_korea_investment_stock.py`
4. Update cache TTL config in `caching/__init__.py` if needed

### Error Handling Pattern

Custom exceptions hierarchy:
- `RateLimitError`: API rate limit hit (triggers retry)
- `TokenExpiredError`: Auth token expired (triggers refresh)
- `APIError`: General API error

Always use `ErrorRecoverySystem` for consistent handling:

```python
try:
    result = api_call()
except Exception as e:
    recovery_system = get_error_recovery_system()
    action = recovery_system.handle_error(e, context={"method": "fetch_price"})
    # Recovery action executed automatically
```

### Cache Key Generation

Cache keys should include all parameters that affect the result:

```python
def _generate_cache_key(self, market: str, symbol: str, detail: bool) -> str:
    return f"price:{market}:{symbol}:{detail}"
```

Avoid caching on error responses (use `cache_condition` parameter).

### Working with US Stocks

US stock queries require real account (mock not supported):
- Exchange auto-detection: NASDAQ → NYSE → AMEX
- Symbol format: Plain ticker (e.g., "AAPL", not "AAPL.US")
- Additional fields: PER, PBR, EPS, BPS included in response

## Testing Strategy

### Test Organization

- **Unit tests**: `test_rate_limiter.py`, `test_ttl_cache.py`, `test_batch_processing.py`
- **Integration tests**: `test_integration.py`, `test_integration_us_stocks.py`
- **Load tests**: `test_load.py`, `test_rate_limit_simulation.py`
- **Feature tests**: `test_ipo_schedule.py`, `test_error_recovery.py`

### Running Integration Tests

Integration tests require valid API credentials:

```bash
# Set credentials in environment
export KOREA_INVESTMENT_API_KEY="..."
export KOREA_INVESTMENT_API_SECRET="..."
export KOREA_INVESTMENT_ACCOUNT_NO="..."

# Run integration tests
pytest korea_investment_stock/tests/test_integration.py -v
```

### Stress Testing

The stress test framework (`scripts/stress_runner.py`) validates:
- Rate limit compliance under high load
- Circuit breaker activation/recovery
- Cache effectiveness
- Error recovery behavior

Configuration: `docs/issue-35/config.yaml`

## Important Files

- **`pyproject.toml`**: Package metadata, dependencies, build config
- **`.cursorrules`**: Development conventions (env vars, virtual env, naming)
- **`CHANGELOG.md`**: Version history and release notes
- **`examples/`**: Usage examples for all major features
- **`docs/`**: Design documents and API exploration notes

## API Rate Limiting

Korea Investment Securities API limits:
- **Official**: 20 requests/second
- **Library default**: 15 requests/second (conservative, 80% safety margin)
- **Minimum interval**: 83ms between calls (when enabled)

Adjust via environment variables:
```bash
export RATE_LIMIT_MAX_CALLS=15
export RATE_LIMIT_SAFETY_MARGIN=0.8
```

## Monitoring & Debugging

### Statistics Collection

```python
with KoreaInvestment(api_key, api_secret, account_no) as broker:
    # ... perform operations ...

    # View rate limiter stats
    broker.rate_limiter.print_stats()

    # View cache stats
    cache_stats = broker.get_cache_stats()
    print(f"Hit rate: {cache_stats['hit_rate']:.1%}")

    # Save all statistics
    from korea_investment_stock.monitoring import get_stats_manager
    stats_mgr = get_stats_manager()
    stats_mgr.save_all_stats(format='json', compress=True)
```

### Dashboard Generation

```python
# Create HTML monitoring dashboard
broker.save_monitoring_dashboard("dashboard.html")

# Generate system health chart
health_chart = broker.get_system_health_chart()

# API usage chart (last 24 hours)
usage_chart = broker.get_api_usage_chart(hours=24)
```

Dashboards saved to project root as HTML files.

## Known Limitations

1. **US stocks**: Real account only (mock not supported by API)
2. **IPO data**: Reference only (from 예탁원), mock not supported
3. **Order functionality**: Not yet implemented (planned)
4. **WebSocket**: Not included in this library

## Git Commit Message Guidelines

Follow conventional commit format:
- `[feat]`: New feature
- `[fix]`: Bug fix
- `[chore]`: Maintenance (docs, CI, etc.)
- `[refactor]`: Code restructuring
- `[test]`: Test additions/changes

Example: `[feat] Add US stock PER/PBR data to fetch_price_list`

## GitHub Actions Workflows

- **`label-merge-conflict.yml`**: Auto-labels PRs with merge conflicts
- Additional workflows may be present for CI/CD

## Context Manager Pattern

Always use context manager for resource cleanup:

```python
with KoreaInvestment(api_key, api_secret, account_no) as broker:
    # Auto-shutdown thread pool and save stats on exit
    prices = broker.fetch_price_list(stock_list)
```

Manual cleanup if not using context manager:
```python
broker = KoreaInvestment(api_key, api_secret, account_no)
try:
    # operations
finally:
    broker.shutdown()  # Essential for thread pool cleanup
```

## Performance Characteristics

**Benchmarks** (from README):
- Throughput: 10-12 TPS (stable)
- Error rate: <0.1%
- 100 stocks query: ~8.5 seconds
- Memory usage: <100MB
- CPU usage: <5%
- Cache hit rate: >80% (typical usage patterns)

## Additional Resources

- **Main docs**: https://wikidocs.net/book/7845
- **Issues**: https://github.com/kenshin579/korea-investment-stock/issues
- **PyPI**: https://pypi.org/project/korea-investment-stock/
