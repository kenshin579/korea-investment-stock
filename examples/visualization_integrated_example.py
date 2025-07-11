#!/usr/bin/env python3
"""
통합 Visualization 예제
Korea Investment Stock 패키지에 통합된 시각화 기능 사용 예제
"""

import korea_investment_stock
import time
from datetime import datetime


def main():
    """메인 함수"""
    print("한국투자증권 통합 시각화 예제")
    print("=" * 50)
    
    # API 키 설정 (실제 사용 시 본인의 키로 교체)
    key = "YOUR_API_KEY"
    secret = "YOUR_API_SECRET"
    acc_no = "12345678-01"
    
    # 환경 변수에서 API 키 가져오기 (선택적)
    import os
    key = os.getenv("KOREA_INVESTMENT_API_KEY", key)
    secret = os.getenv("KOREA_INVESTMENT_API_SECRET", secret)
    acc_no = os.getenv("KOREA_INVESTMENT_ACC_NO", acc_no)
    
    # 브로커 객체 생성
    broker = korea_investment_stock.KoreaInvestment(
        api_key=key, 
        api_secret=secret, 
        acc_no=acc_no,
        mock=True  # 모의투자 서버 사용
    )
    
    print("\n=== 시각화 기능 사용 예제 ===\n")
    
    try:
        # 1. 시스템 헬스 차트 생성
        print("1. 시스템 헬스 차트 생성")
        health_chart = broker.get_system_health_chart()
        if health_chart:
            # 파일로 저장
            broker.visualizer.save_chart(
                health_chart, 
                "system_health.html", 
                format='html'
            )
            print("   ✅ system_health.html 생성 완료")
        else:
            print("   ⚠️ 시스템 헬스 차트 생성 실패 (통계 데이터 없음)")
        
        # 2. API 사용량 차트 생성 (최근 24시간)
        print("\n2. API 사용량 차트 생성 (최근 24시간)")
        api_chart = broker.get_api_usage_chart(hours=24)
        if api_chart:
            # 파일로 저장
            broker.visualizer.save_chart(
                api_chart,
                "api_usage_24h.html",
                format='html'
            )
            print("   ✅ api_usage_24h.html 생성 완료")
        else:
            print("   ⚠️ API 사용량 차트 생성 실패 (통계 데이터 없음)")
        
        # 3. 실시간 모니터링 대시보드 생성
        print("\n3. 실시간 모니터링 대시보드 생성")
        dashboard = broker.create_monitoring_dashboard(
            update_interval=5000  # 5초마다 업데이트
        )
        
        if dashboard:
            # HTML 파일로 저장
            saved = broker.save_monitoring_dashboard("realtime_dashboard.html")
            if saved:
                print("   ✅ realtime_dashboard.html 생성 완료")
                
                # 브라우저에서 열기 (선택적)
                # broker.show_monitoring_dashboard()
            else:
                print("   ❌ 대시보드 저장 실패")
        else:
            print("   ⚠️ 대시보드 생성 실패 (통계 데이터 없음)")
        
        # 4. 종합 리포트 생성
        print("\n4. 종합 리포트 생성")
        report_paths = broker.create_stats_report("api_monitoring_report")
        
        if report_paths:
            print("   ✅ 리포트 파일 생성:")
            for name, path in report_paths.items():
                print(f"      - {name}: {path}")
        else:
            print("   ⚠️ 리포트 생성 실패 (통계 데이터 없음)")
        
        # 5. API 호출하여 통계 생성
        print("\n5. API 호출하여 통계 데이터 생성")
        print("   몇 가지 API를 호출하여 통계를 생성합니다...")
        
        # 샘플 종목 리스트
        sample_stocks = [
            ("005930", "KR"),  # 삼성전자
            ("000660", "KR"),  # SK하이닉스
            ("035720", "KR"),  # 카카오
        ]
        
        # API 호출 (통계 생성을 위해)
        results = broker.fetch_price_list(sample_stocks)
        
        print(f"   ✅ {len(results)}개 종목 조회 완료")
        
        # 잠시 대기 (통계 저장을 위해)
        time.sleep(2)
        
        # 6. 업데이트된 통계로 다시 차트 생성
        print("\n6. 업데이트된 통계로 차트 재생성")
        
        # 시스템 헬스 차트
        health_chart = broker.get_system_health_chart()
        if health_chart:
            broker.visualizer.save_chart(
                health_chart, 
                "system_health_updated.html", 
                format='html'
            )
            print("   ✅ system_health_updated.html 생성 완료")
        
        # 대시보드
        dashboard = broker.create_monitoring_dashboard()
        if dashboard:
            broker.save_monitoring_dashboard("realtime_dashboard_updated.html")
            print("   ✅ realtime_dashboard_updated.html 생성 완료")
        
        # 7. 차트 팩토리 사용 예제
        if broker.visualizer:
            print("\n7. 차트 팩토리 사용 예제")
            from korea_investment_stock.visualization import ChartFactory
            
            chart_factory = ChartFactory()
            
            # 에러율 추이 시계열 차트
            df = broker.visualizer.prepare_dataframe()
            if not df.empty:
                time_series = chart_factory.create_time_series(
                    df,
                    x_col='timestamp',
                    y_cols=['error_rate'],
                    title='에러율 추이',
                    y_title='에러율 (%)'
                )
                broker.visualizer.save_chart(
                    time_series,
                    "error_rate_trend.html",
                    format='html'
                )
                print("   ✅ error_rate_trend.html 생성 완료")
        
        print("\n✅ 모든 시각화 예제 완료!")
        print("\n생성된 HTML 파일들을 브라우저에서 열어보세요.")
        
    except Exception as e:
        print(f"\n❌ 예제 실행 중 오류 발생: {e}")
        import traceback
        traceback.print_exc()
    
    finally:
        # 리소스 정리
        broker.shutdown()


if __name__ == "__main__":
    main() 