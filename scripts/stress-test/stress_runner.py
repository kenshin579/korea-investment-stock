#!/usr/bin/env python3
"""
Stress Test Runner for korea-investment-stock (issue-35)

- Reads docs/issue-35/config.yaml if PyYAML is available; otherwise uses conservative defaults.
- Implements scenarios S1–S4 described in PRD.
- Applies a global token-bucket rate limiter (default cap: 18 rps) independent of library internals.
- Supports dry-run mode (no external API calls) if credentials are missing or --dry-run is passed.
- Collects metrics (latency, outcomes, retries estimated) and writes CSV + a simple summary report to ./stress.

Usage examples:
  python scripts/stress_runner.py --dry-run --duration 15
  python scripts/stress_runner.py --config docs/issue-35/config.yaml --scenarios S1,S2 --target-rps 12 --out-dir ./stress

Environment variables for real run:
  KI_API_KEY, KI_API_SECRET, KI_ACC_NO
"""
from __future__ import annotations

import argparse
import csv
import os
import random
import threading
import time
from collections import deque, defaultdict
from concurrent.futures import ThreadPoolExecutor, as_completed
from dataclasses import dataclass
from pathlib import Path
from statistics import quantiles
from typing import Any, Dict, List, Optional, Tuple

try:
    import yaml  # type: ignore
except Exception:  # pragma: no cover
    yaml = None

# Try import the library
try:
    from korea_investment_stock.korea_investment_stock import KoreaInvestment
    # Import EnhancedRateLimiter from the library to configure client-level throttling
    try:
        from korea_investment_stock.rate_limiting.enhanced_rate_limiter import EnhancedRateLimiter
    except Exception:
        # Fallback path if package structure differs
        from rate_limiting.enhanced_rate_limiter import EnhancedRateLimiter  # type: ignore
except Exception:
    KoreaInvestment = None  # type: ignore
    EnhancedRateLimiter = None  # type: ignore


@dataclass
class LoadProfile:
    target_rps: float = 15.0
    ramp_up_sec: int = 30
    sustain_sec: int = 240
    ramp_down_sec: int = 30
    max_concurrency: int = 100


@dataclass
class RetryConfig:
    max_attempts: int = 3
    backoff: str = "exponential"
    base_delay_ms: int = 100
    jitter_ms: Tuple[int, int] = (50, 150)


@dataclass
class RateLimitGuard:
    global_rps_cap: float = 18.0
    burst_size: int = 5
    retry: RetryConfig = None  # will be set in __post_init__

    def __post_init__(self):
        if self.retry is None:
            self.retry = RetryConfig()


@dataclass
class ScenarioConfig:
    enable: List[str] = None
    s1_tickers: List[str] = None
    s1_duration_sec: int = 300
    s2_use_all: bool = True
    s2_duration_sec: int = 300
    s3_burst_rps: float = 30
    s3_burst_sec: int = 30
    s3_cool_rps: float = 12
    s3_cool_sec: int = 180
    s3_cycles: int = 3
    s4_cache_hit_ratio_target: float = 0.8


@dataclass
class ReportingConfig:
    out_dir: str = "./stress"
    formats: List[str] = None


@dataclass
class Config:
    stock_list: List[Tuple[str, str]]
    load_profile: LoadProfile
    rate_limit_guard: RateLimitGuard
    scenario: ScenarioConfig
    reporting: ReportingConfig


def load_config(path: Optional[str]) -> Config:
    # Defaults aligned with PRD
    default_stock_list = [("005930", "KR"), ("AAPL", "US"), ("091170", "KR"), ("NVDA", "US")]
    lp = LoadProfile()
    rl = RateLimitGuard()
    sc = ScenarioConfig(
        enable=["S1", "S2", "S3", "S4"],
        s1_tickers=["005930:KR", "AAPL:US"],
    )
    rp = ReportingConfig(out_dir="./stress", formats=["csv", "html"])

    if path and yaml is not None and Path(path).exists():
        with open(path, "r", encoding="utf-8") as f:
            data = yaml.safe_load(f)
        stock_list = [(s[0], s[1]) for s in data.get("stock_list", default_stock_list)]
        lp_data = data.get("load_profile", {})
        lp = LoadProfile(
            target_rps=float(lp_data.get("target_rps", lp.target_rps)),
            ramp_up_sec=int(lp_data.get("ramp_up_sec", lp.ramp_up_sec)),
            sustain_sec=int(lp_data.get("sustain_sec", lp.sustain_sec)),
            ramp_down_sec=int(lp_data.get("ramp_down_sec", lp.ramp_down_sec)),
            max_concurrency=int(lp_data.get("max_concurrency", lp.max_concurrency)),
        )
        rl_data = data.get("rate_limit_guard", {})
        r_retry = rl_data.get("retry", {})
        rl = RateLimitGuard(
            global_rps_cap=float(rl_data.get("global_rps_cap", rl.global_rps_cap)),
            burst_size=int(rl_data.get("burst_size", rl.burst_size)),
            retry=RetryConfig(
                max_attempts=int(r_retry.get("max_attempts", 3)),
                backoff=str(r_retry.get("backoff", "exponential")),
                base_delay_ms=int(r_retry.get("base_delay_ms", 100)),
                jitter_ms=tuple(r_retry.get("jitter_ms", [50, 150])),
            ),
        )
        sc_data = data.get("scenario", {})
        enable = sc_data.get("enable", sc.enable)
        s1 = sc_data.get("s1", {})
        s2 = sc_data.get("s2", {})
        s3 = sc_data.get("s3", {})
        s4 = sc_data.get("s4", {})
        sc = ScenarioConfig(
            enable=enable,
            s1_tickers=s1.get("tickers", ["005930:KR", "AAPL:US"]),
            s1_duration_sec=int(s1.get("duration_sec", 300)),
            s2_use_all=bool(s2.get("use_all", True)),
            s2_duration_sec=int(s2.get("duration_sec", 300)),
            s3_burst_rps=float(s3.get("burst_rps", 30)),
            s3_burst_sec=int(s3.get("burst_sec", 30)),
            s3_cool_rps=float(s3.get("cool_rps", 12)),
            s3_cool_sec=int(s3.get("cool_sec", 180)),
            s3_cycles=int(s3.get("cycles", 3)),
            s4_cache_hit_ratio_target=float(s4.get("cache_hit_ratio_target", 0.8)),
        )
        rp_data = data.get("reporting", {})
        rp = ReportingConfig(
            out_dir=str(rp_data.get("out_dir", "./stress")),
            formats=list(rp_data.get("formats", ["csv", "html"])),
        )
    else:
        stock_list = default_stock_list

    return Config(stock_list=stock_list, load_profile=lp, rate_limit_guard=rl, scenario=sc, reporting=rp)


class TokenBucket:
    def __init__(self, rate_per_sec: float, capacity: int):
        self.rate = rate_per_sec
        self.capacity = max(1, capacity)
        self.tokens = capacity
        self.lock = threading.Lock()
        self.last_refill = time.perf_counter()

    def take(self, n: int = 1):
        while True:
            with self.lock:
                now = time.perf_counter()
                elapsed = now - self.last_refill
                add = elapsed * self.rate
                if add > 0:
                    self.tokens = min(self.capacity, self.tokens + add)
                    self.last_refill = now
                if self.tokens >= n:
                    self.tokens -= n
                    return True
            time.sleep(0.001)


def percentile(data: List[float], p: float) -> float:
    if not data:
        return 0.0
    if len(data) == 1:
        return data[0]
    # statistics.quantiles provides quartiles; for generic percentiles we can sort manually
    k = max(0, min(len(data) - 1, int(round((p / 100.0) * (len(data) - 1)))))
    return sorted(data)[k]


def make_client(dry_run: bool) -> Optional[KoreaInvestment]:
    if dry_run:
        return None
    if KoreaInvestment is None:
        raise RuntimeError("Library import failed. Cannot run real test.")
    api_key = os.getenv("KI_API_KEY")
    api_secret = os.getenv("KI_API_SECRET")
    acc_no = os.getenv("KI_ACC_NO")
    if not api_key or not api_secret or not acc_no:
        raise RuntimeError("Missing credentials. Set KI_API_KEY, KI_API_SECRET, KI_ACC_NO or use --dry-run")
    return KoreaInvestment(api_key=api_key, api_secret=api_secret, acc_no=acc_no, mock=False, cache_enabled=True)


def call_endpoint(client: Optional[KoreaInvestment], endpoint: str, batch: List[Tuple[str, str]], retry: RetryConfig, dry_run: bool) -> Tuple[bool, float, str, int]:
    attempts = 0
    start = time.perf_counter()
    last_err = ""
    while attempts < max(1, retry.max_attempts):
        attempts += 1
        try:
            if dry_run:
                # Simulate network latency and success/error for a batch request
                base = random.uniform(0.05, 0.2)
                # Slightly increase latency with batch size to be realistic
                base += 0.01 * max(0, len(batch) - 1)
                # 2% simulated error rate per request
                err = random.random() < 0.02
                time.sleep(base)
                if err:
                    raise RuntimeError("SimulatedError: upstream 5xx")
                latency = (time.perf_counter() - start) * 1000.0
                return True, latency, "200", attempts
            else:
                # Real calls: pass the whole batch to list endpoints
                try:
                    print(f"[API][CALL] ep={endpoint} batch={len(batch)} first={batch[0][0]}:{batch[0][1]} attempt={attempts}")
                except Exception:
                    print(f"[API][CALL] ep={endpoint} batch={len(batch)} attempt={attempts}")
                if endpoint == "fetch_price_list":
                    _ = client.fetch_price_list(batch)
                elif endpoint == "fetch_stock_info_list":
                    _ = client.fetch_stock_info_list(batch)
                elif endpoint == "fetch_search_stock_info_list":
                    # KR-only: filter non-KR items
                    kr_batch = [(s, m) for (s, m) in batch if m in ("KR", "KRX")]
                    if not kr_batch:
                        raise ValueError("search_stock_info supports KR only: empty KR batch")
                    _ = client.fetch_search_stock_info_list(kr_batch)
                else:
                    raise ValueError(f"Unknown endpoint {endpoint}")
                latency = (time.perf_counter() - start) * 1000.0
                print(f"[API][OK] ep={endpoint} batch={len(batch)} latency_ms={latency:.2f} attempts={attempts}")
                return True, latency, "200", attempts
        except Exception as e:  # handle retry
            last_err = type(e).__name__ + ": " + str(e)
            try:
                print(f"[API][ERR] ep={endpoint} attempt={attempts} error={last_err}")
            except Exception:
                print(f"[API][ERR] attempt={attempts} error")
            if attempts >= retry.max_attempts:
                latency = (time.perf_counter() - start) * 1000.0
                return False, latency, last_err, attempts
            # backoff
            delay = retry.base_delay_ms / 1000.0
            if retry.backoff == "exponential":
                delay = delay * (2 ** (attempts - 1))
            elif retry.backoff == "linear":
                delay = delay * attempts
            jmin, jmax = retry.jitter_ms
            jitter = random.uniform(jmin / 1000.0, jmax / 1000.0)
            time.sleep(delay + jitter)
    latency = (time.perf_counter() - start) * 1000.0
    return False, latency, last_err or "UnknownError", attempts


def schedule_requests(total_duration: int, target_rps: float):
    # Generator yielding launch timestamps (relative)
    if target_rps <= 0:
        return
    period = 1.0 / target_rps
    t = 0.0
    while t < total_duration:
        # add +-10% jitter to mitigate sync peaks
        jitter = random.uniform(-0.1 * period, 0.1 * period)
        yield max(0.0, t + jitter)
        t += period


def main():
    parser = argparse.ArgumentParser(description="Stress Test Runner (issue-35)")
    parser.add_argument("--config", default="docs/issue-35/config.yaml")
    parser.add_argument("--out-dir", default=None)
    parser.add_argument("--dry-run", action="store_true")
    parser.add_argument("--scenarios", default=None, help="Comma-separated subset of S1,S2,S3,S4")
    parser.add_argument("--target-rps", type=float, default=None)
    parser.add_argument("--duration", type=int, default=None, help="Override sustain duration for simple run")
    parser.add_argument("--bulk-size", type=int, default=1, help="Number of symbols per request (batch size)")
    args = parser.parse_args()

    cfg = load_config(args.config)
    if args.out_dir:
        cfg.reporting.out_dir = args.out_dir
    if args.target_rps is not None:
        cfg.load_profile.target_rps = args.target_rps
    if args.duration is not None:
        # set a simple uniform profile for quick run
        cfg.load_profile.ramp_up_sec = 0
        cfg.load_profile.sustain_sec = max(1, args.duration)
        cfg.load_profile.ramp_down_sec = 0
        # also override scenario-specific durations for brevity
        cfg.scenario.s1_duration_sec = max(1, args.duration)
        cfg.scenario.s2_duration_sec = max(1, args.duration)
        cfg.scenario.s3_burst_sec = min(cfg.scenario.s3_burst_sec, max(1, args.duration))
        cfg.scenario.s3_cool_sec = min(cfg.scenario.s3_cool_sec, max(1, args.duration))
    enabled = cfg.scenario.enable or ["S1", "S2", "S3", "S4"]
    if args.scenarios:
        enabled = [s.strip() for s in args.scenarios.split(",") if s.strip()]

    dry_run = bool(args.dry_run)

    # Prepare output
    out_dir = Path(cfg.reporting.out_dir)
    out_dir.mkdir(parents=True, exist_ok=True)
    raw_path = out_dir / "stress_raw.csv"
    summary_path = out_dir / "stress_summary.html"

    # Prepare client
    client = None
    if not dry_run:
        try:
            client = make_client(dry_run=False)
            # Override the client's internal rate limiter with safer stress settings
            if client is not None and EnhancedRateLimiter is not None:
                client.rate_limiter = EnhancedRateLimiter(
                    max_calls=7,
                    per_seconds=1.0,
                    safety_margin=0.9,
                    enable_stats=True
                )
        except Exception as e:
            print(f"[WARN] Falling back to dry-run: {e}")
            dry_run = True

    # Build working sets
    # Convert scenario s1 tickers "CODE:MARKET" to tuples
    s1_tickers = []
    if cfg.scenario.s1_tickers:
        for item in cfg.scenario.s1_tickers:
            if ":" in item:
                code, mk = item.split(":", 1)
                s1_tickers.append((code, mk))
    if not s1_tickers:
        s1_tickers = [(cfg.stock_list[0][0], cfg.stock_list[0][1])]

    all_list = list(cfg.stock_list)
    if not all_list:
        all_list = [("005930", "KR"), ("AAPL", "US")]

    # Endpoints under test
    endpoints = [
        "fetch_price_list",
        "fetch_stock_info_list",
        "fetch_search_stock_info_list",
    ]

    # Use library-provided EnhancedRateLimiter via client instead of custom TokenBucket.
    # Global throttling will be handled by client.rate_limiter configured below.

    # ThreadPoolExecutor based on max_concurrency
    max_workers = max(1, min(cfg.load_profile.max_concurrency, 200))
    executor = ThreadPoolExecutor(max_workers=max_workers)

    # Metrics storage
    rows: List[Dict[str, Any]] = []

    def submit_task(ts_rel: float, scenario: str, endpoint: str, batch: List[Tuple[str, str]]):
        def task():
            # Throttling is handled by the library's EnhancedRateLimiter inside the client
            ok, latency_ms, status, attempts = call_endpoint(client, endpoint, batch, cfg.rate_limit_guard.retry, dry_run)
            # Use first symbol/market for compact logging
            first_sym, first_mk = batch[0]
            rows.append({
                "time_s": time.perf_counter() - t0,
                "scenario": scenario,
                "endpoint": endpoint,
                "symbol": first_sym,
                "market": first_mk,
                "batch_size": len(batch),
                "success": 1 if ok else 0,
                "status": status,
                "attempts": attempts,
                "latency_ms": round(latency_ms, 2),
            })
        # Delay until scheduled time
        delay = max(0.0, ts_rel - (time.perf_counter() - t0))
        if delay > 0:
            time.sleep(delay)
        return executor.submit(task)

    futures = []
    t0 = time.perf_counter()

    def run_uniform(duration_sec: int, rps: float, symbol_iter):
        for ts in schedule_requests(duration_sec, rps):
            # Build a batch of symbols
            batch: List[Tuple[str, str]] = [next(symbol_iter) for _ in range(max(1, args.bulk_size))]
            ep = random.choice(endpoints)
            futures.append(submit_task(ts, current_scenario, ep, batch))

    def rr_iter(pool: List[Tuple[str, str]]):
        dq = deque(pool)
        while True:
            s = dq[0]
            dq.rotate(-1)
            yield s

    for current_scenario in enabled:
        if current_scenario == "S1":
            # High frequency single/few tickers
            iter_src = rr_iter(s1_tickers)
            total = cfg.scenario.s1_duration_sec
            rps = min(cfg.load_profile.target_rps, cfg.rate_limit_guard.global_rps_cap)
            run_uniform(total, rps, iter_src)
        elif current_scenario == "S2":
            iter_src = rr_iter(all_list)
            total = cfg.scenario.s2_duration_sec
            rps = min(cfg.load_profile.target_rps, cfg.rate_limit_guard.global_rps_cap)
            run_uniform(total, rps, iter_src)
        elif current_scenario == "S3":
            cycles = cfg.scenario.s3_cycles
            for i in range(cycles):
                iter_src = rr_iter(all_list)
                run_uniform(cfg.scenario.s3_burst_sec, min(cfg.scenario.s3_burst_rps, cfg.rate_limit_guard.global_rps_cap), iter_src)
                iter_src = rr_iter(all_list)
                run_uniform(cfg.scenario.s3_cool_sec, min(cfg.scenario.s3_cool_rps, cfg.rate_limit_guard.global_rps_cap), iter_src)
        elif current_scenario == "S4":
            # Cache hit vs miss: assemble batches with a desired hit ratio
            base_sym = s1_tickers[0]
            alt_iter = rr_iter(all_list)
            duration = max(30, int(cfg.load_profile.sustain_sec / 4))
            rps = min(cfg.load_profile.target_rps, cfg.rate_limit_guard.global_rps_cap)
            for ts in schedule_requests(duration, rps):
                batch: List[Tuple[str, str]] = []
                for _ in range(max(1, args.bulk_size)):
                    if random.random() < cfg.scenario.s4_cache_hit_ratio_target:
                        batch.append(base_sym)
                    else:
                        batch.append(next(alt_iter))
                ep = random.choice(["fetch_price_list", "fetch_stock_info_list"])  # KR-only API excluded here
                futures.append(submit_task(ts, current_scenario, ep, batch))
        else:
            print(f"[WARN] Unknown scenario {current_scenario}, skipping")

    # Wait for completion
    for fut in as_completed(futures):
        try:
            fut.result()
        except Exception as e:
            print("[WorkerError]", e)

    executor.shutdown(wait=True)

    # Write CSV
    with open(raw_path, "w", newline="", encoding="utf-8") as f:
        writer = csv.DictWriter(f, fieldnames=["time_s", "scenario", "endpoint", "symbol", "market", "batch_size", "success", "status", "attempts", "latency_ms"])
        writer.writeheader()
        for r in rows:
            writer.writerow(r)

    # Build summary
    total = len(rows)
    ok = sum(r["success"] for r in rows)
    err = total - ok
    latencies = [r["latency_ms"] for r in rows if r["success"] == 1]
    p50 = percentile(latencies, 50)
    p90 = percentile(latencies, 90)
    p95 = percentile(latencies, 95)
    p99 = percentile(latencies, 99)
    duration_total = max((max([r["time_s"] for r in rows]) if rows else 0.0), 0.001)
    observed_rps = total / max(duration_total, 1e-6)

    html = f"""
    <html><head><meta charset='utf-8'><title>Stress Summary</title></head>
    <body>
    <h2>Stress Test Summary</h2>
    <ul>
      <li>Total requests: {total}</li>
      <li>Success: {ok} | Error: {err} | Error rate: { (err/total*100 if total else 0):.2f}%</li>
      <li>Observed RPS: {observed_rps:.2f}</li>
      <li>Latency p50: {p50:.1f} ms, p90: {p90:.1f} ms, p95: {p95:.1f} ms, p99: {p99:.1f} ms</li>
      <li>Scenarios run: {', '.join(enabled)}</li>
      <li>Dry-run: {dry_run}</li>
    </ul>
    <p>Raw metrics: {raw_path}</p>
    </body></html>
    """
    with open(summary_path, "w", encoding="utf-8") as f:
        f.write(html)

    print(f"Saved raw metrics to: {raw_path}")
    print(f"Saved summary to: {summary_path}")


if __name__ == "__main__":
    main()
