#!/usr/bin/env python3
"""
Plotly 기반 시각화 클래스
실시간 대시보드 및 인터랙티브 차트 생성
"""

import json
import gzip
from pathlib import Path
from datetime import datetime, timedelta
from typing import List, Dict, Any, Optional, Union
import pandas as pd
import numpy as np
import logging

# Plotly imports
try:
    import plotly.graph_objects as go
    from plotly.subplots import make_subplots
    import plotly.express as px
    import plotly.io as pio
    PLOTLY_AVAILABLE = True
except ImportError:
    PLOTLY_AVAILABLE = False
    logging.warning("plotly가 설치되지 않았습니다. 시각화 기능이 제한됩니다.")

from ..monitoring import get_stats_manager

logger = logging.getLogger(__name__)


class PlotlyVisualizer:
    """Plotly를 사용한 통계 시각화 클래스
    
    실시간 모니터링 대시보드와 다양한 차트를 생성합니다.
    """
    
    def __init__(self, stats_dir: str = "logs/integrated_stats"):
        """
        Args:
            stats_dir: 통계 파일이 저장된 디렉토리 경로
        """
        if not PLOTLY_AVAILABLE:
            raise ImportError("plotly가 설치되지 않았습니다. pip install plotly를 실행하세요.")
            
        self.stats_dir = Path(stats_dir)
        self.history_data = []
        self.latest_stats = None
        self.dashboard = None
        self.stats_manager = get_stats_manager()
        
        logger.info(f"PlotlyVisualizer 초기화: stats_dir={stats_dir}")
    
    def load_history_data(self, filename: str = "stats_history.jsonl.gz") -> List[Dict]:
        """압축된 시계열 데이터 로드
        
        Args:
            filename: 로드할 파일명
            
        Returns:
            시계열 통계 데이터 리스트
        """
        filepath = self.stats_dir / filename
        
        if not filepath.exists():
            logger.warning(f"파일을 찾을 수 없습니다: {filepath}")
            return []
        
        # StatsManager를 사용하여 로드
        data = self.stats_manager.load_stats(filepath, format='jsonl')
        
        self.history_data = data
        logger.info(f"{len(data)}개의 시계열 데이터 로드 완료")
        return data
    
    def load_latest_stats(self) -> Dict:
        """가장 최근 통계 파일 로드
        
        Returns:
            최신 통계 데이터
        """
        # JSON 파일 직접 검색
        json_files = list(self.stats_dir.glob("stats_*.json"))
        
        if not json_files:
            logger.warning("JSON 통계 파일을 찾을 수 없습니다.")
            # JSONL 파일에서 최신 데이터 사용
            if self.history_data:
                self.latest_stats = self.history_data[-1]
                logger.info("JSONL에서 최신 통계 로드")
                return self.latest_stats
            return {}
        
        # 가장 최근 파일 선택
        latest_file = max(json_files, key=lambda f: f.stat().st_mtime)
        
        with open(latest_file, 'r', encoding='utf-8') as f:
            self.latest_stats = json.load(f)
            
        logger.info(f"최신 통계 파일 로드: {latest_file.name}")
        return self.latest_stats
    
    def prepare_dataframe(self, data: Optional[List[Dict]] = None) -> pd.DataFrame:
        """히스토리 데이터를 DataFrame으로 변환
        
        Args:
            data: 변환할 데이터 (None이면 self.history_data 사용)
            
        Returns:
            통계 데이터프레임
        """
        if data is None:
            data = self.history_data
            
        records = []
        
        for stat in data:
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
                except Exception as e:
                    logger.error(f"타임스탬프 변환 실패: {e}")
            
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
            
            # modules 데이터 처리
            self._process_modules_data(stat, record)
            
            records.append(record)
        
        # DataFrame 생성 및 정렬
        df = pd.DataFrame(records)
        if not df.empty and 'timestamp' in df.columns:
            df = df.sort_values('timestamp').reset_index(drop=True)
        
        return df
    
    def _process_modules_data(self, stat: Dict, record: Dict):
        """모듈별 데이터 처리"""
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
    
    def create_api_calls_chart(self, df: pd.DataFrame) -> go.Figure:
        """API 호출 추이 차트 생성"""
        fig = go.Figure()
        
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
            )
        )
        
        # 에러 수
        fig.add_trace(
            go.Scatter(
                x=df['timestamp'],
                y=df['total_errors'],
                name='에러',
                line=dict(color='#E74C3C', width=2),
                mode='lines+markers',
                marker=dict(size=6, symbol='x'),
                hovertemplate='에러: %{y}<extra></extra>',
                yaxis='y2'
            )
        )
        
        fig.update_layout(
            title='API 호출 및 에러 추이',
            xaxis_title='시간',
            yaxis_title='API 호출 수',
            yaxis2=dict(
                title='에러 수',
                overlaying='y',
                side='right'
            ),
            hovermode='x unified',
            template='plotly_white'
        )
        
        return fig
    
    def create_cache_efficiency_chart(self, df: pd.DataFrame) -> go.Figure:
        """캐시 효율성 차트 생성"""
        fig = go.Figure()
        
        # 캐시 적중률
        fig.add_trace(
            go.Scatter(
                x=df['timestamp'],
                y=df['cache_hit_rate'],
                name='캐시 적중률 (%)',
                line=dict(color='#2ECC71', width=3),
                mode='lines+markers',
                marker=dict(size=8),
                hovertemplate='캐시 적중률: %{y:.1f}%<extra></extra>'
            )
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
                hovertemplate='누적 절감: %{y:,}<extra></extra>',
                yaxis='y2'
            )
        )
        
        fig.update_layout(
            title='캐시 효율성',
            xaxis_title='시간',
            yaxis_title='캐시 적중률 (%)',
            yaxis2=dict(
                title='누적 절감 호출',
                overlaying='y',
                side='right'
            ),
            hovermode='x unified',
            template='plotly_white'
        )
        
        return fig
    
    def create_rate_limiter_chart(self, df: pd.DataFrame) -> go.Figure:
        """Rate Limiter 성능 차트 생성"""
        fig = make_subplots(
            rows=1, cols=1,
            specs=[[{"secondary_y": True}]]
        )
        
        if 'rl_max_calls_per_second' in df.columns:
            # 실제 TPS
            fig.add_trace(
                go.Scatter(
                    x=df['timestamp'],
                    y=df['rl_max_calls_per_second'],
                    name='실제 TPS',
                    line=dict(color='#F39C12', width=3),
                    mode='lines+markers',
                    hovertemplate='실제 TPS: %{y}<extra></extra>'
                ),
                secondary_y=False
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
                    secondary_y=False
                )
            
            # 평균 대기 시간
            fig.add_trace(
                go.Bar(
                    x=df['timestamp'],
                    y=df['rl_avg_wait_time'] * 1000,  # ms로 변환
                    name='평균 대기 시간 (ms)',
                    marker=dict(color='#E67E22', opacity=0.5),
                    hovertemplate='대기: %{y:.1f}ms<extra></extra>'
                ),
                secondary_y=True
            )
        
        fig.update_xaxes(title_text="시간")
        fig.update_yaxes(title_text="TPS", secondary_y=False)
        fig.update_yaxes(title_text="대기 시간 (ms)", secondary_y=True)
        
        fig.update_layout(
            title='Rate Limiter 성능',
            hovermode='x unified',
            template='plotly_white'
        )
        
        return fig
    
    def create_error_distribution_chart(self) -> go.Figure:
        """에러 타입 분포 차트 생성"""
        if not self.latest_stats:
            return go.Figure()
            
        fig = go.Figure()
        
        if 'modules' in self.latest_stats:
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
                        )
                    )
        
        fig.update_layout(
            title='에러 타입 분포',
            template='plotly_white'
        )
        
        return fig
    
    def create_system_health_indicator(self) -> go.Figure:
        """시스템 헬스 인디케이터 생성"""
        if not self.latest_stats:
            return go.Figure()
            
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
        
        fig = go.Figure()
        
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
                }
            )
        )
        
        fig.update_layout(
            height=400,
            template='plotly_white'
        )
        
        return fig
    
    def save_chart(self, fig: go.Figure, filename: str, 
                  format: str = 'html', **kwargs) -> str:
        """차트를 파일로 저장
        
        Args:
            fig: 저장할 Figure 객체
            filename: 파일명
            format: 저장 형식 ('html', 'png', 'pdf', 'svg')
            **kwargs: 추가 옵션
            
        Returns:
            저장된 파일 경로
        """
        output_path = Path(filename)
        
        if format == 'html':
            fig.write_html(
                output_path,
                include_plotlyjs=kwargs.get('include_plotlyjs', 'cdn'),
                config=kwargs.get('config', {'displayModeBar': True, 'displaylogo': False})
            )
        elif format in ['png', 'pdf', 'svg']:
            try:
                fig.write_image(
                    output_path,
                    format=format,
                    width=kwargs.get('width', 1200),
                    height=kwargs.get('height', 800),
                    scale=kwargs.get('scale', 2)
                )
            except ImportError:
                logger.error(f"{format} 형식으로 저장하려면 kaleido가 필요합니다: pip install kaleido")
                return ""
        else:
            logger.error(f"지원하지 않는 형식: {format}")
            return ""
        
        logger.info(f"차트 저장 완료: {output_path}")
        return str(output_path)
    
    def get_health_color(self, health: str) -> str:
        """시스템 상태에 따른 색상 반환"""
        colors = {
            'HEALTHY': '#2ECC71',
            'WARNING': '#F39C12',
            'CRITICAL': '#E74C3C',
            'UNKNOWN': '#95A5A6'
        }
        return colors.get(health, '#95A5A6') 