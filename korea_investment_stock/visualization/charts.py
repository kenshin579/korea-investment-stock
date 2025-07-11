#!/usr/bin/env python3
"""
차트 팩토리
다양한 유형의 차트를 쉽게 생성할 수 있는 팩토리 클래스
"""

import logging
from typing import List, Dict, Any, Optional, Union
import pandas as pd

try:
    import plotly.graph_objects as go
    import plotly.express as px
    from plotly.subplots import make_subplots
    PLOTLY_AVAILABLE = True
except ImportError:
    PLOTLY_AVAILABLE = False
    logging.warning("plotly가 설치되지 않았습니다.")

logger = logging.getLogger(__name__)


class ChartFactory:
    """차트 생성 팩토리 클래스
    
    다양한 유형의 차트를 쉽게 생성할 수 있습니다.
    """
    
    # 기본 색상 팔레트
    DEFAULT_COLORS = {
        'primary': '#3498DB',
        'success': '#2ECC71',
        'warning': '#F39C12',
        'danger': '#E74C3C',
        'info': '#16A085',
        'secondary': '#9B59B6',
        'light': '#95A5A6',
        'dark': '#2C3E50'
    }
    
    def __init__(self):
        """ChartFactory 초기화"""
        if not PLOTLY_AVAILABLE:
            raise ImportError("plotly가 설치되지 않았습니다.")
        
        self.template = 'plotly_white'
        logger.info("ChartFactory 초기화 완료")
    
    def create_time_series(self, 
                          data: pd.DataFrame,
                          x_col: str,
                          y_cols: Union[str, List[str]],
                          title: str = "시계열 차트",
                          y_title: str = "값",
                          secondary_y_cols: Optional[List[str]] = None) -> go.Figure:
        """시계열 차트 생성
        
        Args:
            data: 데이터프레임
            x_col: X축 컬럼명 (시간)
            y_cols: Y축 컬럼명(들)
            title: 차트 제목
            y_title: Y축 제목
            secondary_y_cols: 보조 Y축에 표시할 컬럼들
            
        Returns:
            시계열 차트
        """
        if isinstance(y_cols, str):
            y_cols = [y_cols]
        
        secondary_y_cols = secondary_y_cols or []
        
        # 보조 Y축 사용 여부
        has_secondary = len(secondary_y_cols) > 0
        
        if has_secondary:
            fig = make_subplots(specs=[[{"secondary_y": True}]])
        else:
            fig = go.Figure()
        
        # 색상 설정
        colors = list(self.DEFAULT_COLORS.values())
        
        # 각 컬럼에 대한 트레이스 추가
        for i, col in enumerate(y_cols):
            is_secondary = col in secondary_y_cols
            
            trace = go.Scatter(
                x=data[x_col],
                y=data[col],
                name=col,
                line=dict(color=colors[i % len(colors)], width=2),
                mode='lines+markers',
                marker=dict(size=6),
                hovertemplate=f'{col}: %{{y}}<extra></extra>'
            )
            
            if has_secondary:
                fig.add_trace(trace, secondary_y=is_secondary)
            else:
                fig.add_trace(trace)
        
        # 레이아웃 설정
        fig.update_layout(
            title=title,
            xaxis_title="시간",
            template=self.template,
            hovermode='x unified'
        )
        
        if has_secondary:
            fig.update_yaxes(title_text=y_title, secondary_y=False)
            fig.update_yaxes(title_text="보조 Y축", secondary_y=True)
        else:
            fig.update_yaxes(title_text=y_title)
        
        return fig
    
    def create_bar_chart(self,
                        data: pd.DataFrame,
                        x_col: str,
                        y_col: str,
                        title: str = "막대 차트",
                        orientation: str = 'v',
                        color_col: Optional[str] = None) -> go.Figure:
        """막대 차트 생성
        
        Args:
            data: 데이터프레임
            x_col: X축 컬럼명
            y_col: Y축 컬럼명
            title: 차트 제목
            orientation: 'v' (수직) 또는 'h' (수평)
            color_col: 색상 구분을 위한 컬럼
            
        Returns:
            막대 차트
        """
        if color_col:
            fig = px.bar(
                data,
                x=x_col if orientation == 'v' else y_col,
                y=y_col if orientation == 'v' else x_col,
                color=color_col,
                title=title,
                orientation=orientation,
                template=self.template
            )
        else:
            fig = go.Figure([
                go.Bar(
                    x=data[x_col] if orientation == 'v' else data[y_col],
                    y=data[y_col] if orientation == 'v' else data[x_col],
                    orientation=orientation,
                    marker_color=self.DEFAULT_COLORS['primary']
                )
            ])
            
            fig.update_layout(
                title=title,
                xaxis_title=x_col if orientation == 'v' else y_col,
                yaxis_title=y_col if orientation == 'v' else x_col,
                template=self.template
            )
        
        return fig
    
    def create_pie_chart(self,
                        labels: List[str],
                        values: List[float],
                        title: str = "파이 차트",
                        hole: float = 0) -> go.Figure:
        """파이 차트 생성
        
        Args:
            labels: 라벨 리스트
            values: 값 리스트
            title: 차트 제목
            hole: 도넛 차트를 위한 구멍 크기 (0-1)
            
        Returns:
            파이 차트
        """
        fig = go.Figure([
            go.Pie(
                labels=labels,
                values=values,
                hole=hole,
                marker=dict(
                    colors=px.colors.qualitative.Set3[:len(labels)]
                ),
                textinfo='label+percent',
                hovertemplate='%{label}: %{value}<br>%{percent}<extra></extra>'
            )
        ])
        
        fig.update_layout(
            title=title,
            template=self.template
        )
        
        return fig
    
    def create_gauge_chart(self,
                          value: float,
                          title: str = "게이지 차트",
                          range_min: float = 0,
                          range_max: float = 100,
                          thresholds: Optional[List[Dict]] = None) -> go.Figure:
        """게이지 차트 생성
        
        Args:
            value: 현재 값
            title: 차트 제목
            range_min: 최소값
            range_max: 최대값
            thresholds: 임계값 설정 [{'value': 80, 'color': 'red'}, ...]
            
        Returns:
            게이지 차트
        """
        # 기본 임계값 설정
        if thresholds is None:
            thresholds = [
                {'value': range_max * 0.9, 'color': 'red'},
                {'value': range_max * 0.7, 'color': 'orange'}
            ]
        
        # 색상 결정
        color = self.DEFAULT_COLORS['success']
        for threshold in sorted(thresholds, key=lambda x: x['value'], reverse=True):
            if value >= threshold['value']:
                color = threshold['color']
                break
        
        fig = go.Figure([
            go.Indicator(
                mode="gauge+number+delta",
                value=value,
                title={'text': title},
                delta={'reference': range_max * 0.8},
                gauge={
                    'axis': {'range': [range_min, range_max]},
                    'bar': {'color': color},
                    'steps': [
                        {'range': [range_min, range_max * 0.5], 'color': "lightgray"},
                        {'range': [range_max * 0.5, range_max * 0.8], 'color': "gray"}
                    ],
                    'threshold': {
                        'line': {'color': "red", 'width': 4},
                        'thickness': 0.75,
                        'value': range_max * 0.9
                    }
                }
            )
        ])
        
        fig.update_layout(
            height=400,
            template=self.template
        )
        
        return fig
    
    def create_heatmap(self,
                      data: pd.DataFrame,
                      title: str = "히트맵",
                      colorscale: str = 'RdBu') -> go.Figure:
        """히트맵 생성
        
        Args:
            data: 2D 데이터프레임
            title: 차트 제목
            colorscale: 색상 스케일
            
        Returns:
            히트맵
        """
        fig = go.Figure([
            go.Heatmap(
                z=data.values,
                x=data.columns,
                y=data.index,
                colorscale=colorscale,
                hovertemplate='%{x}<br>%{y}<br>값: %{z}<extra></extra>'
            )
        ])
        
        fig.update_layout(
            title=title,
            template=self.template
        )
        
        return fig
    
    def create_scatter_plot(self,
                           data: pd.DataFrame,
                           x_col: str,
                           y_col: str,
                           title: str = "산점도",
                           size_col: Optional[str] = None,
                           color_col: Optional[str] = None) -> go.Figure:
        """산점도 생성
        
        Args:
            data: 데이터프레임
            x_col: X축 컬럼명
            y_col: Y축 컬럼명
            title: 차트 제목
            size_col: 크기를 결정할 컬럼
            color_col: 색상을 결정할 컬럼
            
        Returns:
            산점도
        """
        if size_col or color_col:
            fig = px.scatter(
                data,
                x=x_col,
                y=y_col,
                size=size_col,
                color=color_col,
                title=title,
                template=self.template
            )
        else:
            fig = go.Figure([
                go.Scatter(
                    x=data[x_col],
                    y=data[y_col],
                    mode='markers',
                    marker=dict(
                        size=10,
                        color=self.DEFAULT_COLORS['primary']
                    ),
                    hovertemplate=f'{x_col}: %{{x}}<br>{y_col}: %{{y}}<extra></extra>'
                )
            ])
            
            fig.update_layout(
                title=title,
                xaxis_title=x_col,
                yaxis_title=y_col,
                template=self.template
            )
        
        return fig
    
    def create_box_plot(self,
                       data: pd.DataFrame,
                       y_col: str,
                       x_col: Optional[str] = None,
                       title: str = "박스 플롯") -> go.Figure:
        """박스 플롯 생성
        
        Args:
            data: 데이터프레임
            y_col: Y축 컬럼명
            x_col: X축 컬럼명 (그룹화용, 선택사항)
            title: 차트 제목
            
        Returns:
            박스 플롯
        """
        if x_col:
            fig = px.box(
                data,
                x=x_col,
                y=y_col,
                title=title,
                template=self.template
            )
        else:
            fig = go.Figure([
                go.Box(
                    y=data[y_col],
                    name=y_col,
                    marker_color=self.DEFAULT_COLORS['primary']
                )
            ])
            
            fig.update_layout(
                title=title,
                yaxis_title=y_col,
                template=self.template
            )
        
        return fig
    
    def create_multi_axis_chart(self,
                               data: pd.DataFrame,
                               x_col: str,
                               y_configs: List[Dict[str, Any]],
                               title: str = "다중 축 차트") -> go.Figure:
        """다중 축 차트 생성
        
        Args:
            data: 데이터프레임
            x_col: X축 컬럼명
            y_configs: Y축 설정 리스트
                      [{'col': 'col1', 'type': 'line', 'yaxis': 'y1', 'name': '라인1'},
                       {'col': 'col2', 'type': 'bar', 'yaxis': 'y2', 'name': '막대1'}]
            title: 차트 제목
            
        Returns:
            다중 축 차트
        """
        fig = go.Figure()
        
        # Y축 개수 확인
        y_axes = list(set([cfg.get('yaxis', 'y1') for cfg in y_configs]))
        
        # 각 설정에 따라 트레이스 추가
        for i, cfg in enumerate(y_configs):
            col = cfg['col']
            chart_type = cfg.get('type', 'line')
            yaxis = cfg.get('yaxis', 'y1')
            name = cfg.get('name', col)
            color = cfg.get('color', list(self.DEFAULT_COLORS.values())[i % len(self.DEFAULT_COLORS)])
            
            if chart_type == 'line':
                trace = go.Scatter(
                    x=data[x_col],
                    y=data[col],
                    name=name,
                    line=dict(color=color, width=2),
                    yaxis=yaxis.replace('y', 'y')
                )
            elif chart_type == 'bar':
                trace = go.Bar(
                    x=data[x_col],
                    y=data[col],
                    name=name,
                    marker_color=color,
                    yaxis=yaxis.replace('y', 'y'),
                    opacity=0.7
                )
            else:
                continue
            
            fig.add_trace(trace)
        
        # 레이아웃 설정
        layout_dict = {
            'title': title,
            'xaxis': {'title': x_col},
            'template': self.template,
            'hovermode': 'x unified'
        }
        
        # Y축 설정
        for i, yaxis in enumerate(y_axes):
            if i == 0:
                layout_dict['yaxis'] = {'title': yaxis}
            else:
                layout_dict[f'yaxis{i+1}'] = {
                    'title': yaxis,
                    'overlaying': 'y',
                    'side': 'right' if i == 1 else 'left',
                    'position': 0.85 if i > 1 else None
                }
        
        fig.update_layout(**layout_dict)
        
        return fig
    
    def apply_custom_theme(self, fig: go.Figure, theme: Dict[str, Any]) -> go.Figure:
        """커스텀 테마 적용
        
        Args:
            fig: Figure 객체
            theme: 테마 설정 딕셔너리
            
        Returns:
            테마가 적용된 Figure
        """
        # 기본 테마 설정
        default_theme = {
            'font': {'family': 'Arial, sans-serif', 'size': 12},
            'plot_bgcolor': 'white',
            'paper_bgcolor': 'white',
            'margin': {'l': 60, 'r': 60, 't': 80, 'b': 60}
        }
        
        # 사용자 테마로 업데이트
        default_theme.update(theme)
        
        # 테마 적용
        fig.update_layout(**default_theme)
        
        return fig 