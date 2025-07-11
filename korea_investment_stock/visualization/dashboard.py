#!/usr/bin/env python3
"""
대시보드 관리자
실시간 모니터링 대시보드와 종합 리포트 생성
"""

import logging
from pathlib import Path
from datetime import datetime
from typing import Dict, Any, Optional, List
import pandas as pd

try:
    import plotly.graph_objects as go
    from plotly.subplots import make_subplots
    import plotly.express as px
    PLOTLY_AVAILABLE = True
except ImportError:
    PLOTLY_AVAILABLE = False
    logging.warning("plotly가 설치되지 않았습니다.")

from .plotly_visualizer import PlotlyVisualizer

logger = logging.getLogger(__name__)


class DashboardManager:
    """대시보드 관리자 클래스
    
    실시간 모니터링 대시보드와 종합 리포트를 생성하고 관리합니다.
    """
    
    def __init__(self, visualizer: Optional[PlotlyVisualizer] = None):
        """
        Args:
            visualizer: PlotlyVisualizer 인스턴스 (None이면 새로 생성)
        """
        if not PLOTLY_AVAILABLE:
            raise ImportError("plotly가 설치되지 않았습니다.")
            
        self.visualizer = visualizer or PlotlyVisualizer()
        self.dashboard = None
        self.summary_card = None
        
        logger.info("DashboardManager 초기화 완료")
    
    def create_realtime_dashboard(self, update_interval: int = 5000) -> go.Figure:
        """실시간 업데이트가 가능한 대시보드 생성
        
        Args:
            update_interval: 업데이트 간격 (밀리초)
            
        Returns:
            대시보드 Figure 객체
        """
        # 데이터 로드
        if not self.visualizer.history_data:
            logger.warning("시계열 데이터가 없습니다.")
            return None
        
        # 데이터 준비
        df = self.visualizer.prepare_dataframe()
        
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
                'buttons': self._create_update_buttons(update_interval),
                'pad': {'r': 10, 't': 60},
                'showactive': False,
                'x': 0.1,
                'xanchor': 'right',
                'y': 1.15,
                'yanchor': 'top'
            }]
        )
        
        # 축 포맷 업데이트
        self._update_axes_format(fig)
        
        self.dashboard = fig
        logger.info("실시간 대시보드 생성 완료")
        return fig
    
    def _create_update_buttons(self, update_interval: int) -> List[Dict]:
        """업데이트 버튼 생성"""
        return [
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
        ]
    
    def _update_axes_format(self, fig):
        """축 포맷 업데이트"""
        for row in [2, 3]:
            for col in [1, 2]:
                if row == 3 and col == 2:
                    continue  # pie chart
                fig.update_xaxes(
                    title_text="시간",
                    tickformat="%H:%M:%S",
                    row=row, col=col
                )
    
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
        if not self.visualizer.latest_stats:
            return
        
        summary = self.visualizer.latest_stats.get('summary', {})
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
        if self.visualizer.latest_stats and 'modules' in self.visualizer.latest_stats:
            if 'error_recovery' in self.visualizer.latest_stats['modules']:
                er = self.visualizer.latest_stats['modules']['error_recovery']
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
        if not self.visualizer.latest_stats:
            return None
        
        summary = self.visualizer.latest_stats.get('summary', {})
        modules = self.visualizer.latest_stats.get('modules', {})
        
        # 지표 카드 생성
        fig = make_subplots(
            rows=2, cols=4,
            specs=[
                [{"type": "indicator"}, {"type": "indicator"}, {"type": "indicator"}, {"type": "indicator"}],
                [{"type": "indicator"}, {"type": "indicator"}, {"type": "indicator"}, {"type": "indicator"}]
            ]
        )
        
        # 지표 정의
        indicators = self._create_summary_indicators(summary, modules)
        
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
        
        self.summary_card = fig
        logger.info("요약 카드 생성 완료")
        return fig
    
    def _create_summary_indicators(self, summary: Dict, modules: Dict) -> List[Dict]:
        """요약 지표 생성"""
        indicators = [
            {
                'value': 100 - summary.get('overall_error_rate', 0) * 100,  # 건강도 점수
                'mode': 'number+gauge',
                'title': {'text': f"시스템 상태: {summary.get('system_health', 'UNKNOWN')}"},
                'number': {'suffix': '%', 'font': {'color': self.visualizer.get_health_color(summary.get('system_health', 'UNKNOWN'))}},
                'gauge': {
                    'axis': {'range': [0, 100]},
                    'bar': {'color': self.visualizer.get_health_color(summary.get('system_health', 'UNKNOWN'))},
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
        
        return indicators
    
    def save_dashboard(self, filename: str = "api_monitoring_dashboard.html", 
                      include_plotlyjs: str = 'cdn') -> str:
        """대시보드를 HTML 파일로 저장
        
        Args:
            filename: 저장할 파일명
            include_plotlyjs: 'cdn', 'inline', 'directory' 중 선택
            
        Returns:
            저장된 파일 경로
        """
        if self.dashboard:
            path = self.visualizer.save_chart(
                self.dashboard, 
                filename, 
                format='html',
                include_plotlyjs=include_plotlyjs
            )
            logger.info(f"대시보드 저장 완료: {path}")
            return path
        else:
            logger.warning("저장할 대시보드가 없습니다.")
            return ""
    
    def show_dashboard(self):
        """대시보드 표시"""
        if self.dashboard:
            self.dashboard.show()
        else:
            logger.warning("표시할 대시보드가 없습니다.")
    
    def create_report(self, save_as: str = "monitoring_report") -> Dict[str, str]:
        """종합 리포트 생성
        
        Args:
            save_as: 저장할 파일명 (확장자 제외)
            
        Returns:
            생성된 파일 경로들
        """
        paths = {}
        
        # 대시보드 이미지
        if self.dashboard:
            try:
                path = self.visualizer.save_chart(
                    self.dashboard,
                    f"{save_as}_dashboard.png",
                    format='png',
                    width=1600,
                    height=1000,
                    scale=2
                )
                paths['dashboard'] = path
            except Exception as e:
                logger.error(f"대시보드 이미지 생성 실패: {e}")
        
        # 요약 카드 이미지
        if self.summary_card:
            try:
                path = self.visualizer.save_chart(
                    self.summary_card,
                    f"{save_as}_summary.png",
                    format='png',
                    width=1200,
                    height=400,
                    scale=2
                )
                paths['summary'] = path
            except Exception as e:
                logger.error(f"요약 카드 이미지 생성 실패: {e}")
        
        # 개별 차트들
        if self.visualizer.history_data:
            df = self.visualizer.prepare_dataframe()
            
            # API 호출 차트
            fig = self.visualizer.create_api_calls_chart(df)
            path = self.visualizer.save_chart(
                fig, f"{save_as}_api_calls.html", format='html'
            )
            paths['api_calls'] = path
            
            # 캐시 효율성 차트
            fig = self.visualizer.create_cache_efficiency_chart(df)
            path = self.visualizer.save_chart(
                fig, f"{save_as}_cache.html", format='html'
            )
            paths['cache'] = path
            
            # Rate Limiter 차트
            fig = self.visualizer.create_rate_limiter_chart(df)
            path = self.visualizer.save_chart(
                fig, f"{save_as}_rate_limiter.html", format='html'
            )
            paths['rate_limiter'] = path
        
        logger.info(f"리포트 생성 완료: {len(paths)}개 파일")
        return paths 