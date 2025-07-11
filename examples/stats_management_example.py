#!/usr/bin/env python3
"""
통합 통계 관리 예제
Date: 2024-12-28
Issue: #27 - Phase 5.1

통합 통계 관리자를 사용하여 다양한 형식으로 통계를 저장하고 분석하는 예제
"""

import sys
import os
import time
import json
from pathlib import Path

# 모듈 경로 추가
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from korea_investment_stock import KoreaInvestment
from korea_investment_stock.monitoring import StatsManager, get_stats_manager
from korea_investment_stock.rate_limiting import EnhancedRateLimiter, get_backoff_strategy
from korea_investment_stock.error_handling import get_error_recovery_system
from korea_investment_stock.batch_processing import DynamicBatchController


def example_basic_stats_collection():
    """기본 통계 수집 및 저장 예제"""
    print("=== 1. 기본 통계 수집 예제 ===\n")
    
    # StatsManager 생성
    stats_mgr = StatsManager(base_dir="example_stats", retention_days=30)
    
    # 각 모듈 생성 및 활동
    rate_limiter = EnhancedRateLimiter(max_calls=15)
    
    # 몇 가지 API 호출 시뮬레이션
    print("API 호출 시뮬레이션...")
    for i in range(10):
        if rate_limiter.acquire():
            print(f"  API 호출 {i+1} 성공")
            time.sleep(0.1)
    
    # 에러 시뮬레이션
    rate_limiter.record_error()
    
    # 통계 수집
    all_stats = stats_mgr.collect_all_stats(rate_limiter=rate_limiter)
    
    # JSON으로 저장
    json_path = stats_mgr.save_stats(all_stats, format='json')
    print(f"\n통계 저장됨: {json_path}")
    
    # 저장된 내용 출력
    with open(json_path, 'r', encoding='utf-8') as f:
        saved = json.load(f)
    
    print("\n저장된 통계 요약:")
    summary = saved['summary']
    print(f"- 시스템 상태: {summary['system_health']}")
    print(f"- 총 API 호출: {summary['total_api_calls']}")
    print(f"- 에러율: {summary['overall_error_rate']:.1%}")


def example_multi_format_export():
    """다양한 형식으로 내보내기 예제"""
    print("\n\n=== 2. 다양한 형식으로 내보내기 예제 ===\n")
    
    stats_mgr = get_stats_manager()
    
    # 테스트 데이터
    test_stats = {
        'timestamp': '2024-12-28T10:00:00',
        'modules': {
            'rate_limiter': {
                'total_calls': 1000,
                'error_count': 5,
                'max_calls_per_second': 12
            }
        },
        'summary': {
            'system_health': 'HEALTHY',
            'overall_error_rate': 0.005
        }
    }
    
    # 각 형식으로 저장
    formats = ['json', 'csv', 'jsonl']
    for fmt in formats:
        path = stats_mgr.save_stats(test_stats, format=fmt)
        file_size = Path(path).stat().st_size
        print(f"{fmt.upper()} 형식: {path} ({file_size} bytes)")
    
    # 압축 저장
    print("\n압축 저장:")
    for fmt in formats:
        path = stats_mgr.save_stats(test_stats, format=fmt, compress=True)
        file_size = Path(path).stat().st_size
        print(f"{fmt.upper()}.GZ 형식: {path} ({file_size} bytes)")


def example_real_world_usage():
    """실제 사용 시나리오 예제"""
    print("\n\n=== 3. 실제 사용 시나리오 예제 ===\n")
    
    # KoreaInvestment 인스턴스 생성 (Mock 모드)
    broker = KoreaInvestment(
        api_key="test_key",
        api_secret="test_secret", 
        acc_no="12345678-01",
        mock=True
    )
    
    # 통계 관리자
    stats_mgr = get_stats_manager()
    
    # 주기적 통계 저장 함수
    def save_periodic_stats():
        """주기적으로 통계 저장"""
        all_stats = stats_mgr.collect_all_stats(
            rate_limiter=broker.rate_limiter,
            backoff_strategy=get_backoff_strategy(),
            error_recovery=get_error_recovery_system()
        )
        
        # JSON Lines 형식으로 추가 (시계열 분석용)
        jsonl_path = stats_mgr.save_stats(
            all_stats,
            format='jsonl',
            filename='stats_timeline',
            include_timestamp=False
        )
        
        return all_stats['summary']
    
    # 시뮬레이션
    print("API 사용 시뮬레이션 중...")
    
    # 5초마다 통계 저장
    for i in range(3):
        # API 호출 시뮬레이션
        for j in range(5):
            broker.rate_limiter.acquire()
            time.sleep(0.1)
        
        # 통계 저장
        summary = save_periodic_stats()
        print(f"\n[{i+1}/3] 통계 저장 완료")
        print(f"  - 상태: {summary['system_health']}")
        print(f"  - API 호출: {summary['total_api_calls']}")
        
        if i < 2:
            time.sleep(2)
    
    # 최종 통계 분석
    print("\n최종 통계 분석:")
    timeline_file = Path("logs/integrated_stats/stats_timeline.jsonl")
    
    if timeline_file.exists():
        # 시계열 데이터 로드
        timeline_data = stats_mgr.load_stats(timeline_file, format='jsonl')
        
        print(f"총 {len(timeline_data)}개의 스냅샷 저장됨")
        
        # API 호출 추이
        total_calls = [d['summary']['total_api_calls'] for d in timeline_data]
        print(f"API 호출 추이: {' → '.join(map(str, total_calls))}")


def example_stats_analysis():
    """저장된 통계 분석 예제"""
    print("\n\n=== 4. 저장된 통계 분석 예제 ===\n")
    
    stats_mgr = get_stats_manager()
    
    # 최신 통계 파일 찾기
    latest_json = stats_mgr.get_latest_stats_file(format='json')
    
    if latest_json:
        print(f"최신 통계 파일: {latest_json}")
        
        # 통계 로드
        stats = stats_mgr.load_stats(latest_json)
        
        # 분석
        print("\n통계 분석:")
        
        # Rate Limiter 분석
        if 'rate_limiter' in stats.get('modules', {}):
            rl = stats['modules']['rate_limiter']
            print(f"\nRate Limiter:")
            print(f"  - 총 호출: {rl.get('total_calls', 0)}")
            print(f"  - 최대 TPS: {rl.get('max_calls_per_second', 0)}")
            print(f"  - 평균 대기: {rl.get('avg_wait_time', 0):.3f}초")
        
        # Error Recovery 분석
        if 'error_recovery' in stats.get('modules', {}):
            er = stats['modules']['error_recovery']
            print(f"\nError Recovery:")
            print(f"  - 총 에러: {er.get('total_errors', 0)}")
            
            if 'by_type' in er:
                print("  - 에러 타입별:")
                for err_type, count in er['by_type'].items():
                    print(f"    - {err_type}: {count}회")
    else:
        print("저장된 통계 파일이 없습니다.")


def main():
    """모든 예제 실행"""
    print("=" * 60)
    print("통합 통계 관리 예제")
    print("=" * 60)
    
    try:
        # 1. 기본 통계 수집
        example_basic_stats_collection()
        
        # 2. 다양한 형식 내보내기
        example_multi_format_export()
        
        # 3. 실제 사용 시나리오
        example_real_world_usage()
        
        # 4. 통계 분석
        example_stats_analysis()
        
        print("\n" + "=" * 60)
        print("모든 예제 실행 완료!")
        print("=" * 60)
        
    except Exception as e:
        print(f"\n예제 실행 중 오류: {e}")
        import traceback
        traceback.print_exc()


if __name__ == "__main__":
    main() 