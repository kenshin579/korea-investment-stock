#!/usr/bin/env python3
"""
Load Tests for Rate Limiting System
Date: 2024-12-28

Phase 6.3: 부하 테스트 작성
- 최대 처리량 측정 스크립트
- 에러율 측정 및 리포트
- 성능 프로파일링
"""

import time
import threading
import statistics
from concurrent.futures import ThreadPoolExecutor, as_completed
from dataclasses import dataclass
from typing import List, Dict, Any
import json

from ..rate_limiting import EnhancedRateLimiter
from .test_integration import MockAPIServer


@dataclass
class LoadTestResult:
    """부하 테스트 결과"""
    duration: float
    total_requests: int
    successful_requests: int
    failed_requests: int
    rate_limit_errors: int
    avg_tps: float
    max_tps: float
    min_tps: float
    error_rate: float
    response_times: List[float]


class LoadTester:
    """부하 테스트 실행기"""
    
    def __init__(self, rate_limiter: EnhancedRateLimiter, mock_server: MockAPIServer):
        self.rate_limiter = rate_limiter
        self.mock_server = mock_server
        self.results = []
        self.lock = threading.Lock()
        
    def single_request(self):
        """단일 요청 실행"""
        start_time = time.time()
        
        if not self.rate_limiter.acquire():
            return {
                "success": False,
                "error": "rate_limiter_blocked",
                "duration": time.time() - start_time
            }
        
        try:
            response = self.mock_server.api_call()
            duration = time.time() - start_time
            
            if response["rt_cd"] == "0":
                return {
                    "success": True,
                    "duration": duration
                }
            else:
                return {
                    "success": False,
                    "error": response["rt_cd"],
                    "duration": duration
                }
        except Exception as e:
            return {
                "success": False,
                "error": str(e),
                "duration": time.time() - start_time
            }
    
    def run_load_test(self, duration_seconds: int, max_workers: int) -> LoadTestResult:
        """부하 테스트 실행"""
        print(f"\n부하 테스트 시작: {duration_seconds}초, {max_workers} 워커")
        
        start_time = time.time()
        end_time = start_time + duration_seconds
        
        results = []
        tps_per_second = {}
        
        with ThreadPoolExecutor(max_workers=max_workers) as executor:
            futures = []
            
            while time.time() < end_time:
                # 지속적으로 요청 제출
                future = executor.submit(self.single_request)
                futures.append((future, time.time()))
                
                # CPU 과부하 방지
                time.sleep(0.001)
            
            # 모든 결과 수집
            for future, submit_time in futures:
                try:
                    result = future.result(timeout=5)
                    result["submit_time"] = submit_time
                    results.append(result)
                    
                    # TPS 계산용 데이터 수집
                    second = int(submit_time - start_time)
                    if second not in tps_per_second:
                        tps_per_second[second] = 0
                    tps_per_second[second] += 1
                    
                except Exception as e:
                    results.append({
                        "success": False,
                        "error": str(e),
                        "duration": 0,
                        "submit_time": submit_time
                    })
        
        # 결과 분석
        total_duration = time.time() - start_time
        successful = [r for r in results if r.get("success", False)]
        failed = [r for r in results if not r.get("success", False)]
        rate_limit_errors = [r for r in failed if r.get("error") == "EGW00201"]
        
        response_times = [r["duration"] for r in successful if r["duration"] > 0]
        
        tps_values = list(tps_per_second.values())
        
        return LoadTestResult(
            duration=total_duration,
            total_requests=len(results),
            successful_requests=len(successful),
            failed_requests=len(failed),
            rate_limit_errors=len(rate_limit_errors),
            avg_tps=len(results) / total_duration,
            max_tps=max(tps_values) if tps_values else 0,
            min_tps=min(tps_values) if tps_values else 0,
            error_rate=len(failed) / len(results) if results else 0,
            response_times=response_times
        )


def test_maximum_throughput():
    """최대 처리량 측정"""
    print("\n" + "="*60)
    print("최대 처리량 측정 테스트")
    print("="*60)
    
    # 다양한 설정으로 테스트
    configurations = [
        {"max_calls": 10, "safety_margin": 0.8, "workers": 3},
        {"max_calls": 15, "safety_margin": 0.8, "workers": 3},
        {"max_calls": 15, "safety_margin": 0.9, "workers": 5},
        {"max_calls": 20, "safety_margin": 0.7, "workers": 5},
    ]
    
    results_summary = []
    
    for config in configurations:
        print(f"\n\n테스트 설정: max_calls={config['max_calls']}, "
              f"safety_margin={config['safety_margin']}, workers={config['workers']}")
        
        rate_limiter = EnhancedRateLimiter(
            max_calls=config["max_calls"],
            per_seconds=1.0,
            safety_margin=config["safety_margin"]
        )
        
        mock_server = MockAPIServer(rate_limit=20)
        tester = LoadTester(rate_limiter, mock_server)
        
        # 30초 부하 테스트
        result = tester.run_load_test(duration_seconds=30, max_workers=config["workers"])
        
        # 결과 출력
        print(f"\n결과:")
        print(f"- 총 요청: {result.total_requests}")
        print(f"- 성공: {result.successful_requests}")
        print(f"- 실패: {result.failed_requests}")
        print(f"- Rate Limit 에러: {result.rate_limit_errors}")
        print(f"- 평균 TPS: {result.avg_tps:.2f}")
        print(f"- 최대 TPS: {result.max_tps}")
        print(f"- 최소 TPS: {result.min_tps}")
        print(f"- 에러율: {result.error_rate:.2%}")
        
        if result.response_times:
            print(f"- 평균 응답 시간: {statistics.mean(result.response_times)*1000:.1f}ms")
            print(f"- P95 응답 시간: {statistics.quantiles(result.response_times, n=20)[18]*1000:.1f}ms")
        
        # 서버 통계
        server_stats = mock_server.get_stats()
        print(f"- 서버 Rate Limit 에러: {server_stats['error_count']}")
        
        results_summary.append({
            "config": config,
            "avg_tps": result.avg_tps,
            "error_rate": result.error_rate,
            "server_errors": server_stats['error_count']
        })
    
    # 최적 설정 찾기
    print("\n\n" + "="*60)
    print("최적 설정 분석")
    print("="*60)
    
    # 서버 에러가 0인 설정 중 최대 TPS
    valid_configs = [r for r in results_summary if r["server_errors"] == 0]
    if valid_configs:
        best_config = max(valid_configs, key=lambda x: x["avg_tps"])
        print(f"\n최적 설정: {best_config['config']}")
        print(f"- 평균 TPS: {best_config['avg_tps']:.2f}")
        print(f"- 에러율: {best_config['error_rate']:.2%}")
    
    return results_summary


def test_stress_conditions():
    """스트레스 조건 테스트"""
    print("\n" + "="*60)
    print("스트레스 조건 테스트")
    print("="*60)
    
    rate_limiter = EnhancedRateLimiter(
        max_calls=15,
        per_seconds=1.0,
        safety_margin=0.8
    )
    
    mock_server = MockAPIServer(rate_limit=20)
    
    # 1. 버스트 트래픽 테스트
    print("\n1. 버스트 트래픽 테스트 (100개 동시 요청)")
    burst_results = []
    
    with ThreadPoolExecutor(max_workers=100) as executor:
        futures = []
        start_time = time.time()
        
        # 100개 요청 동시 제출
        for i in range(100):
            future = executor.submit(lambda: (
                rate_limiter.acquire(),
                mock_server.api_call() if rate_limiter.acquire() else {"rt_cd": "BLOCKED"}
            ))
            futures.append(future)
        
        for future in as_completed(futures):
            result = future.result()
            burst_results.append(result)
    
    burst_duration = time.time() - start_time
    burst_success = sum(1 for _, r in burst_results if r.get("rt_cd") == "0")
    
    print(f"- 실행 시간: {burst_duration:.2f}초")
    print(f"- 성공: {burst_success}/100")
    print(f"- 서버 에러: {mock_server.get_stats()['error_count']}")
    
    # 2. 지속적 고부하 테스트
    print("\n2. 지속적 고부하 테스트 (60초)")
    tester = LoadTester(rate_limiter, MockAPIServer(rate_limit=20))
    high_load_result = tester.run_load_test(duration_seconds=60, max_workers=10)
    
    print(f"- 총 요청: {high_load_result.total_requests}")
    print(f"- 평균 TPS: {high_load_result.avg_tps:.2f}")
    print(f"- 에러율: {high_load_result.error_rate:.2%}")
    
    # 3. 성능 저하 감지
    print("\n3. 성능 저하 패턴 분석")
    if high_load_result.response_times:
        # 시간대별 응답 시간 분석
        time_windows = []
        window_size = 10  # 10초 단위
        
        for i in range(0, 60, window_size):
            window_times = [
                r["duration"] for r in tester.results 
                if i <= (r["submit_time"] - start_time) < i + window_size
                and r.get("success", False)
            ]
            
            if window_times:
                avg_response = statistics.mean(window_times) * 1000
                time_windows.append({
                    "period": f"{i}-{i+window_size}초",
                    "avg_response_ms": avg_response
                })
        
        # 성능 저하 확인
        if time_windows:
            first_window_avg = time_windows[0]["avg_response_ms"]
            last_window_avg = time_windows[-1]["avg_response_ms"]
            
            print(f"- 초기 평균 응답: {first_window_avg:.1f}ms")
            print(f"- 후기 평균 응답: {last_window_avg:.1f}ms")
            
            if last_window_avg > first_window_avg * 1.5:
                print("⚠️  성능 저하 감지: 응답 시간이 50% 이상 증가")
            else:
                print("✅ 안정적인 성능 유지")


def generate_performance_report():
    """성능 리포트 생성"""
    print("\n" + "="*60)
    print("종합 성능 리포트 생성")
    print("="*60)
    
    # 최대 처리량 테스트
    throughput_results = test_maximum_throughput()
    
    # 스트레스 테스트
    test_stress_conditions()
    
    # 리포트 생성
    report = {
        "test_date": time.strftime("%Y-%m-%d %H:%M:%S"),
        "summary": {
            "optimal_tps": max(r["avg_tps"] for r in throughput_results if r["server_errors"] == 0),
            "zero_error_configs": len([r for r in throughput_results if r["server_errors"] == 0]),
            "tested_configurations": len(throughput_results)
        },
        "configurations": throughput_results,
        "recommendations": []
    }
    
    # 권장사항 추가
    if report["summary"]["optimal_tps"] >= 10:
        report["recommendations"].append("현재 설정으로 초당 10회 이상의 안정적인 처리 가능")
    
    if all(r["server_errors"] == 0 for r in throughput_results):
        report["recommendations"].append("모든 테스트 설정에서 서버 Rate Limit 에러 0건 달성")
    
    # JSON 파일로 저장
    with open("performance_report.json", "w", encoding="utf-8") as f:
        json.dump(report, f, indent=2, ensure_ascii=False)
    
    print("\n✅ 성능 리포트가 'performance_report.json'에 저장되었습니다.")
    
    return report


if __name__ == "__main__":
    # 전체 부하 테스트 실행
    report = generate_performance_report()
    
    print("\n" + "="*60)
    print("부하 테스트 완료")
    print(f"최적 TPS: {report['summary']['optimal_tps']:.2f}")
    print("="*60) 