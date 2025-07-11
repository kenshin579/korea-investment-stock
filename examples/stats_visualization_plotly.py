#!/usr/bin/env python3
"""
통계 시각화 예제 (Plotly 독립형 버전)
저장된 통계 파일을 읽어서 인터랙티브한 그래프를 생성합니다.

이 예제는 고급 사용자를 위한 상세한 구현을 보여줍니다.
일반적인 사용은 visualization_integrated_example.py를 참고하세요.

특징:
- 실시간 업데이트 지원
- 인터랙티브 대시보드
- 다양한 통계 지표 시각화
- HTML 내보내기 지원
- 패키지 없이 독립적으로 실행 가능
"""

import json
import gzip
import time
from pathlib import Path
from datetime import datetime, timedelta
from typing import List, Dict, Any, Optional
import pandas as pd
import numpy as np

# Plotly imports
import plotly.graph_objects as go
from plotly.subplots import make_subplots
import plotly.express as px
import plotly.io as pio

# Korea Investment Stock imports
import sys
import os
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from korea_investment_stock.monitoring import get_stats_manager


class PlotlyStatsVisualizer:
    """Plotly를 사용한 통계 시각화 클래스"""
    
    def __init__(self, stats_dir: str = "logs/integrated_stats"):
        self.stats_dir = Path(stats_dir)
        self.history_data = []
        self.latest_stats = None
        self.dashboard = None
        self.stats_manager = get_stats_manager()
        
    def load_history_data(self, filename: str = "stats_history.jsonl.gz") -> List[Dict]:
        """압축된 시계열 데이터 로드"""
        filepath = self.stats_dir / filename
        
        if not filepath.exists():
            print(f"파일을 찾을 수 없습니다: {filepath}")
            return []
        
        # StatsManager를 사용하여 로드
        data = self.stats_manager.load_stats(filepath, format='jsonl')
        
        self.history_data = data
        print(f"✅ {len(data)}개의 시계열 데이터 로드 완료")
        return data
    
    def load_latest_stats(self) -> Dict:
        """가장 최근 통계 파일 로드"""
        # JSON 파일 직접 검색
        json_files = list(self.stats_dir.glob("stats_*.json"))
        
        if not json_files:
            print("JSON 통계 파일을 찾을 수 없습니다.")
            # JSONL 파일에서 최신 데이터 사용
            if self.history_data:
                self.latest_stats = self.history_data[-1]
                print(f"✅ JSONL에서 최신 통계 로드")
                return self.latest_stats
            return {}
        
        # 가장 최근 파일 선택
        latest_file = max(json_files, key=lambda f: f.stat().st_mtime)
        
        with open(latest_file, 'r', encoding='utf-8') as f:
            self.latest_stats = json.load(f)
            
        print(f"✅ 최신 통계 파일 로드: {latest_file.name}")
        return self.latest_stats
    
    def create_realtime_dashboard(self, update_interval: int = 5000):
        """실시간 업데이트가 가능한 대시보드 생성
        
        Args:
            update_interval: 업데이트 간격 (밀리초)
        """
        if not self.history_data:
            print("시계열 데이터가 없습니다.")
            return None
        
        # 데이터 준비
        df = self._prepare_dataframe()
        
        # 레이아웃 생성 (3x2 그리드)
        fig = make_subplots(
            rows=3, cols=2,
            subplot_titles=(
                '📈 API 호출 및 에러 추이',
                '📊 시스템 헬스 상태',
                '💾 캐시 효율성',
                '⚡ Rate Limiter 성능',
                '🔄 배치 처리 효율',
                '🚨 에러 타입 분포'
            ),
            specs=[
                [{"secondary_y": True}, {"type": "indicator"}],
                [{"secondary_y": False}, {"secondary_y": True}],
                [{"secondary_y": False}, {"type": "pie"}]
            ],
            vertical_spacing=0.1,
            horizontal_spacing=0.12,
            row_heights=[0.35, 0.35, 0.3]
        )
        
        # 1. API 호출 추이
        self._add_api_calls_trace(fig, df, row=1, col=1)
        
        # 2. 시스템 헬스 인디케이터
        self._add_health_indicator(fig, row=1, col=2)
        
        # 3. 캐시 효율성
        self._add_cache_efficiency_trace(fig, df, row=2, col=1)
        
        # 4. Rate Limiter 성능
        self._add_rate_limiter_performance(fig, df, row=2, col=2)
        
        # 5. 배치 처리 효율
        self._add_batch_processing_trace(fig, df, row=3, col=1)
        
        # 6. 에러 타입 분포
        self._add_error_distribution(fig, row=3, col=2)
        
        # 레이아웃 설정
        fig.update_layout(
            title={
                'text': '한국투자증권 API 실시간 모니터링 대시보드',
                'font': {'size': 26, 'color': '#2C3E50', 'family': 'Arial Black'},
                'x': 0.5,
                'xanchor': 'center'
            },
            height=1000,
            showlegend=True,
            template='plotly_white',
            hovermode='x unified',
            updatemenus=[{
                'type': 'buttons',
                'direction': 'left',
                'buttons': [
                    {
                        'label': '▶ 재생',
                        'method': 'animate',
                        'args': [None, {
                            'frame': {'duration': update_interval, 'redraw': True},
                            'fromcurrent': True
                        }]
                    },
                    {
                        'label': '⏸ 일시정지',
                        'method': 'animate',
                        'args': [[None], {
                            'frame': {'duration': 0, 'redraw': False},
                            'mode': 'immediate',
                            'transition': {'duration': 0}
                        }]
                    }
                ],
                'pad': {'r': 10, 't': 60},
                'showactive': False,
                'x': 0.1,
                'xanchor': 'right',
                'y': 1.15,
                'yanchor': 'top'
            }]
        )
        
        # 축 포맷 업데이트
        for row in [2, 3]:
            for col in [1, 2]:
                if row == 3 and col == 2:
                    continue  # pie chart
                fig.update_xaxes(
                    title_text="시간",
                    tickformat="%H:%M:%S",
                    row=row, col=col
                )
        
        self.dashboard = fig
        return fig
    
    def _prepare_dataframe(self) -> pd.DataFrame:
        """히스토리 데이터를 DataFrame으로 변환"""
        records = []
        
        for stat in self.history_data:
            # 기본 값들로 record 초기화
            record = {
                'timestamp': datetime.now(),
                'total_api_calls': 0,
                'total_errors': 0,
                'error_rate': 0,
                'cache_hit_rate': 0,
                'api_calls_saved': 0,
                'system_health': 'UNKNOWN',
                'max_tps': 0
            }
            
            # timestamp 처리
            if 'timestamp' in stat:
                try:
                    if isinstance(stat['timestamp'], str):
                        record['timestamp'] = datetime.fromisoformat(stat['timestamp'])
                    else:
                        record['timestamp'] = stat['timestamp']
                except:
                    pass
            
            # summary 데이터 처리
            if 'summary' in stat and isinstance(stat['summary'], dict):
                summary = stat['summary']
                record.update({
                    'total_api_calls': summary.get('total_api_calls', 0),
                    'total_errors': summary.get('total_errors', 0),
                    'error_rate': summary.get('overall_error_rate', 0) * 100,
                    'cache_hit_rate': summary.get('cache_hit_rate', 0) * 100,
                    'api_calls_saved': summary.get('api_calls_saved', 0),
                    'system_health': summary.get('system_health', 'UNKNOWN'),
                    'max_tps': summary.get('max_tps', 0)
                })
            
            # modules 데이터가 있을 경우에만 처리
            if 'modules' in stat and isinstance(stat['modules'], dict):
                # Rate Limiter 데이터
                if 'rate_limiter' in stat['modules']:
                    rl = stat['modules']['rate_limiter']
                    record['rl_total_calls'] = rl.get('total_calls', 0)
                    record['rl_error_count'] = rl.get('error_count', 0)
                    record['rl_avg_wait_time'] = rl.get('avg_wait_time', 0)
                    record['rl_max_calls_per_second'] = rl.get('max_calls_per_second', 0)
                    
                    # config 정보
                    if 'config' in rl:
                        record['rl_nominal_max_calls'] = rl['config'].get('nominal_max_calls', 15)
                        record['rl_effective_max_calls'] = rl['config'].get('effective_max_calls', 12)
                
                # Batch Controller 데이터
                if 'batch_controller' in stat['modules']:
                    bc = stat['modules']['batch_controller']
                    record['bc_current_batch_size'] = bc.get('current_batch_size', 0)
                    record['bc_total_batches'] = bc.get('total_batches', 0)
                
                # Error Recovery 데이터
                if 'error_recovery' in stat['modules']:
                    er = stat['modules']['error_recovery']
                    record['er_total_errors'] = er.get('total_errors', 0)
                    if 'by_type' in er:
                        record['error_types'] = er['by_type']
            
            # 기존 통계 데이터 형식 (modules가 없는 경우) 처리
            elif 'rate_limiter' in stat:
                # 직접 rate_limiter 데이터에 접근
                rl = stat.get('rate_limiter', {})
                if isinstance(rl, dict):
                    record['rl_total_calls'] = rl.get('total_calls', 0)
                    record['rl_error_count'] = rl.get('error_count', 0)
                    record['rl_avg_wait_time'] = rl.get('avg_wait_time', 0)
                    record['rl_max_calls_per_second'] = rl.get('max_calls_per_second', 0)
                    record['total_api_calls'] = rl.get('total_calls', 0)
                    record['total_errors'] = rl.get('error_count', 0)
                    
                    if rl.get('total_calls', 0) > 0:
                        record['error_rate'] = (rl.get('error_count', 0) / rl['total_calls']) * 100
            
            records.append(record)
        
        # DataFrame 생성 및 정렬
        df = pd.DataFrame(records)
        if not df.empty and 'timestamp' in df.columns:
            df = df.sort_values('timestamp').reset_index(drop=True)
        
        return df
    
    def _add_api_calls_trace(self, fig, df, row, col):
        """API 호출 추이 추가"""
        # API 호출 수 (Area chart)
        fig.add_trace(
            go.Scatter(
                x=df['timestamp'],
                y=df['total_api_calls'],
                name='API 호출',
                line=dict(color='#3498DB', width=3),
                fill='tozeroy',
                fillcolor='rgba(52, 152, 219, 0.2)',
                mode='lines',
                hovertemplate='API 호출: %{y:,}<extra></extra>'
            ),
            row=row, col=col, secondary_y=False
        )
        
        # 에러 수 (보조 Y축)
        fig.add_trace(
            go.Scatter(
                x=df['timestamp'],
                y=df['total_errors'],
                name='에러',
                line=dict(color='#E74C3C', width=2),
                mode='lines+markers',
                marker=dict(size=6, symbol='x'),
                hovertemplate='에러: %{y}<extra></extra>'
            ),
            row=row, col=col, secondary_y=True
        )
        
        # 축 레이블
        fig.update_yaxes(title_text="API 호출 수", secondary_y=False, row=row, col=col)
        fig.update_yaxes(title_text="에러 수", secondary_y=True, row=row, col=col)
    
    def _add_health_indicator(self, fig, row, col):
        """시스템 헬스 인디케이터 추가"""
        if not self.latest_stats:
            return
        
        summary = self.latest_stats.get('summary', {})
        health = summary.get('system_health', 'UNKNOWN')
        error_rate = summary.get('overall_error_rate', 0) * 100
        
        # 색상 매핑
        color_map = {
            'HEALTHY': '#2ECC71',
            'WARNING': '#F39C12',
            'CRITICAL': '#E74C3C',
            'UNKNOWN': '#95A5A6'
        }
        
        fig.add_trace(
            go.Indicator(
                mode="gauge+number+delta",
                value=100 - error_rate,  # 건강도 (100 - 에러율)
                title={'text': f"시스템 상태: {health}"},
                delta={'reference': 95},  # 95% 이상이 목표
                gauge={
                    'axis': {'range': [0, 100]},
                    'bar': {'color': color_map.get(health, '#95A5A6')},
                    'steps': [
                        {'range': [0, 95], 'color': "lightgray"},
                        {'range': [95, 99], 'color': "gray"}
                    ],
                    'threshold': {
                        'line': {'color': "red", 'width': 4},
                        'thickness': 0.75,
                        'value': 99
                    }
                },
                domain={'row': row-1, 'column': col-1}
            ),
            row=row, col=col
        )
    
    def _add_cache_efficiency_trace(self, fig, df, row, col):
        """캐시 효율성 시각화"""
        # 캐시 적중률과 절감된 API 호출 수를 함께 표시
        fig.add_trace(
            go.Scatter(
                x=df['timestamp'],
                y=df['cache_hit_rate'],
                name='캐시 적중률 (%)',
                line=dict(color='#2ECC71', width=3),
                mode='lines+markers',
                marker=dict(size=8),
                yaxis='y',
                hovertemplate='캐시 적중률: %{y:.1f}%<extra></extra>'
            ),
            row=row, col=col
        )
        
        # 절감된 API 호출 수 (누적)
        df['api_calls_saved_cumsum'] = df['api_calls_saved'].cumsum()
        fig.add_trace(
            go.Scatter(
                x=df['timestamp'],
                y=df['api_calls_saved_cumsum'],
                name='누적 절감 호출',
                line=dict(color='#16A085', width=2, dash='dash'),
                mode='lines',
                yaxis='y2',
                hovertemplate='누적 절감: %{y:,}<extra></extra>'
            ),
            row=row, col=col
        )
        
        fig.update_yaxes(title_text="캐시 적중률 (%)", row=row, col=col)
    
    def _add_rate_limiter_performance(self, fig, df, row, col):
        """Rate Limiter 성능 지표"""
        if 'rl_max_calls_per_second' in df.columns:
            # 실제 TPS vs 제한값
            fig.add_trace(
                go.Scatter(
                    x=df['timestamp'],
                    y=df['rl_max_calls_per_second'],
                    name='실제 TPS',
                    line=dict(color='#F39C12', width=3),
                    mode='lines+markers',
                    hovertemplate='실제 TPS: %{y}<extra></extra>'
                ),
                row=row, col=col, secondary_y=False
            )
            
            # 제한값 라인
            if 'rl_effective_max_calls' in df.columns:
                fig.add_trace(
                    go.Scatter(
                        x=df['timestamp'],
                        y=df['rl_effective_max_calls'],
                        name='TPS 제한',
                        line=dict(color='red', width=2, dash='dash'),
                        mode='lines',
                        hovertemplate='제한: %{y}<extra></extra>'
                    ),
                    row=row, col=col, secondary_y=False
                )
            
            # 평균 대기 시간
            fig.add_trace(
                go.Bar(
                    x=df['timestamp'],
                    y=df['rl_avg_wait_time'] * 1000,  # ms로 변환
                    name='평균 대기 시간 (ms)',
                    marker=dict(color='#E67E22', opacity=0.5),
                    yaxis='y2',
                    hovertemplate='대기: %{y:.1f}ms<extra></extra>'
                ),
                row=row, col=col, secondary_y=True
            )
            
            fig.update_yaxes(title_text="TPS", secondary_y=False, row=row, col=col)
            fig.update_yaxes(title_text="대기 시간 (ms)", secondary_y=True, row=row, col=col)
    
    def _add_batch_processing_trace(self, fig, df, row, col):
        """배치 처리 효율성"""
        if 'bc_current_batch_size' in df.columns:
            # 배치 크기 변화
            fig.add_trace(
                go.Scatter(
                    x=df['timestamp'],
                    y=df['bc_current_batch_size'],
                    name='배치 크기',
                    line=dict(color='#9B59B6', width=3),
                    mode='lines+markers+text',
                    text=df['bc_current_batch_size'],
                    textposition='top center',
                    hovertemplate='배치 크기: %{y}<extra></extra>'
                ),
                row=row, col=col
            )
            
            fig.update_yaxes(title_text="배치 크기", row=row, col=col)
            
            # 에러율에 따른 배경색 추가
            for i in range(len(df)-1):
                if df['error_rate'].iloc[i] > 5:
                    color = 'rgba(231, 76, 60, 0.1)'  # 빨강
                elif df['error_rate'].iloc[i] > 1:
                    color = 'rgba(243, 156, 18, 0.1)'  # 주황
                else:
                    color = 'rgba(46, 204, 113, 0.1)'  # 초록
                
                fig.add_vrect(
                    x0=df['timestamp'].iloc[i],
                    x1=df['timestamp'].iloc[i+1] if i+1 < len(df) else df['timestamp'].iloc[i],
                    fillcolor=color,
                    layer="below",
                    line_width=0,
                    row=row, col=col
                )
    
    def _add_error_distribution(self, fig, row, col):
        """에러 타입 분포 (파이 차트)"""
        if self.latest_stats and 'modules' in self.latest_stats:
            if 'error_recovery' in self.latest_stats['modules']:
                er = self.latest_stats['modules']['error_recovery']
                if 'by_type' in er:
                    error_types = er['by_type']
                    
                    labels = list(error_types.keys())
                    values = list(error_types.values())
                    
                    fig.add_trace(
                        go.Pie(
                            labels=labels,
                            values=values,
                            hole=0.4,
                            marker=dict(
                                colors=px.colors.qualitative.Set3[:len(labels)]
                            ),
                            textinfo='label+percent',
                            hovertemplate='%{label}: %{value}건<br>%{percent}<extra></extra>'
                        ),
                        row=row, col=col
                    )
    
    def create_summary_card(self) -> go.Figure:
        """향상된 요약 정보 카드"""
        if not self.latest_stats:
            return None
        
        summary = self.latest_stats.get('summary', {})
        modules = self.latest_stats.get('modules', {})
        
        # 지표 카드 생성
        fig = make_subplots(
            rows=2, cols=4,
            specs=[
                [{"type": "indicator"}, {"type": "indicator"}, {"type": "indicator"}, {"type": "indicator"}],
                [{"type": "indicator"}, {"type": "indicator"}, {"type": "indicator"}, {"type": "indicator"}]
            ]
        )
        
        # 첫 번째 행
        indicators = [
            {
                'value': 100 - summary.get('overall_error_rate', 0) * 100,  # 건강도 점수
                'mode': 'number+gauge',
                'title': {'text': f"시스템 상태: {summary.get('system_health', 'UNKNOWN')}"},
                'number': {'suffix': '%', 'font': {'color': self._get_health_color(summary.get('system_health', 'UNKNOWN'))}},
                'gauge': {
                    'axis': {'range': [0, 100]},
                    'bar': {'color': self._get_health_color(summary.get('system_health', 'UNKNOWN'))},
                    'threshold': {
                        'line': {'color': "red", 'width': 4},
                        'thickness': 0.75,
                        'value': 95
                    }
                },
                'row': 1, 'col': 1
            },
            {
                'value': summary.get('total_api_calls', 0),
                'mode': 'number',
                'title': {'text': 'API 호출'},
                'number': {'font': {'color': '#3498DB'}},
                'row': 1, 'col': 2
            },
            {
                'value': summary.get('overall_error_rate', 0) * 100,
                'mode': 'number+gauge',
                'title': {'text': '에러율'},
                'number': {'suffix': '%', 'font': {'color': '#E74C3C'}},
                'gauge': {'axis': {'range': [0, 10]}},
                'row': 1, 'col': 3
            },
            {
                'value': summary.get('cache_hit_rate', 0) * 100,
                'mode': 'number+gauge',
                'title': {'text': '캐시 적중률'},
                'number': {'suffix': '%', 'font': {'color': '#2ECC71'}},
                'gauge': {'axis': {'range': [0, 100]}},
                'row': 1, 'col': 4
            }
        ]
        
        # 두 번째 행 (Rate Limiter 통계)
        if 'rate_limiter' in modules:
            rl = modules['rate_limiter']
            indicators.extend([
                {
                    'value': rl.get('max_calls_per_second', 0),
                    'mode': 'number',
                    'title': {'text': '최대 TPS'},
                    'number': {'font': {'color': '#F39C12'}},
                    'row': 2, 'col': 1
                },
                {
                    'value': rl.get('avg_wait_time', 0) * 1000,
                    'mode': 'number',
                    'title': {'text': '평균 대기시간'},
                    'number': {'suffix': 'ms', 'font': {'color': '#E67E22'}},
                    'row': 2, 'col': 2
                },
                {
                    'value': summary.get('api_calls_saved', 0),
                    'mode': 'number',
                    'title': {'text': '절감된 호출'},
                    'number': {'font': {'color': '#16A085'}},
                    'row': 2, 'col': 3
                },
                {
                    'value': modules.get('batch_controller', {}).get('current_batch_size', 0),
                    'mode': 'number',
                    'title': {'text': '활성 배치'},
                    'number': {'font': {'color': '#9B59B6'}},
                    'row': 2, 'col': 4
                }
            ])
        
        # 인디케이터 추가
        for ind in indicators:
            fig.add_trace(
                go.Indicator(
                    value=ind['value'],
                    mode=ind.get('mode', 'number'),
                    title=ind.get('title', {}),
                    number=ind.get('number', {}),
                    gauge=ind.get('gauge', {}),
                ),
                row=ind['row'], col=ind['col']
            )
        
        fig.update_layout(
            title={
                'text': "시스템 현황 요약 대시보드",
                'font': {'size': 20},
                'x': 0.5,
                'xanchor': 'center'
            },
            height=400,
            showlegend=False,
            margin=dict(l=20, r=20, t=60, b=20)
        )
        
        return fig
    
    def _get_health_color(self, health: str) -> str:
        """시스템 상태에 따른 색상 반환"""
        colors = {
            'HEALTHY': '#2ECC71',
            'WARNING': '#F39C12',
            'CRITICAL': '#E74C3C',
            'UNKNOWN': '#95A5A6'
        }
        return colors.get(health, '#95A5A6')
    
    def save_dashboard(self, filename: str = "api_monitoring_dashboard.html", 
                      include_plotlyjs: str = 'cdn'):
        """대시보드를 HTML 파일로 저장
        
        Args:
            filename: 저장할 파일명
            include_plotlyjs: 'cdn', 'inline', 'directory' 중 선택
        """
        if self.dashboard:
            self.dashboard.write_html(
                filename, 
                include_plotlyjs=include_plotlyjs,
                config={'displayModeBar': True, 'displaylogo': False}
            )
            print(f"✅ 대시보드가 저장되었습니다: {filename}")
        else:
            print("저장할 대시보드가 없습니다.")
    
    def show_dashboard(self):
        """대시보드 표시"""
        if self.dashboard:
            self.dashboard.show()
        else:
            print("표시할 대시보드가 없습니다.")
    
    def create_report(self, save_as: str = "monitoring_report.pdf"):
        """PDF 리포트 생성 (plotly + kaleido 필요)"""
        try:
            import kaleido
            
            # 모든 차트를 이미지로 저장
            if self.dashboard:
                self.dashboard.write_image(f"{save_as.replace('.pdf', '')}_dashboard.png", 
                                         width=1600, height=1000, scale=2)
            
            summary_card = self.create_summary_card()
            if summary_card:
                summary_card.write_image(f"{save_as.replace('.pdf', '')}_summary.png",
                                       width=1200, height=400, scale=2)
            
            print(f"✅ 리포트 이미지가 생성되었습니다.")
            
        except ImportError:
            print("PDF 리포트 생성을 위해서는 kaleido가 필요합니다: pip install kaleido")


def main():
    """메인 함수"""
    print("=" * 60)
    print("한국투자증권 API 통계 시각화 (Plotly Enhanced)")
    print("=" * 60)
    
    # 시각화 객체 생성
    visualizer = PlotlyStatsVisualizer()
    
    # 데이터 로드
    print("\n1. 데이터 로드 중...")
    visualizer.load_history_data()
    visualizer.load_latest_stats()
    
    if not visualizer.history_data and not visualizer.latest_stats:
        print("시각화할 데이터가 없습니다.")
        return
    
    # 대시보드 생성
    print("\n2. 인터랙티브 대시보드 생성 중...")
    
    try:
        # 실시간 대시보드
        dashboard = visualizer.create_realtime_dashboard(update_interval=5000)
        
        # 요약 카드
        summary_card = visualizer.create_summary_card()
        
        if dashboard:
            print("\n📊 대시보드 생성 완료!")
            
            # HTML로 저장
            visualizer.save_dashboard("api_monitoring_realtime.html")
            
            # 대시보드 표시
            print("\n브라우저에서 대시보드를 여는 중...")
            visualizer.show_dashboard()
            
            # 요약 카드도 표시
            if summary_card:
                summary_card.show()
            
            # PDF 리포트 생성 (선택사항)
            print("\n3. 리포트 생성 시도 중...")
            visualizer.create_report()
        
    except Exception as e:
        print(f"대시보드 생성 중 오류: {e}")
        import traceback
        traceback.print_exc()
    
    print("\n✅ 완료!")
    print("=" * 60)


if __name__ == "__main__":
    main() 