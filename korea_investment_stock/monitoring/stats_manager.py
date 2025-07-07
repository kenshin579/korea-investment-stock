#!/usr/bin/env python3
"""
통합 통계 관리자
Date: 2024-12-28
Issue: #27 - Phase 5.1

다양한 모듈의 통계를 통합 관리하고 여러 형식으로 저장
"""

import json
import csv
import gzip
import time
import logging
from pathlib import Path
from datetime import datetime, timedelta
from typing import Dict, Any, Optional, List, Union
from collections import defaultdict

logger = logging.getLogger(__name__)


class StatsManager:
    """통합 통계 관리자
    
    여러 모듈의 통계를 수집하고 다양한 형식으로 저장/관리합니다.
    """
    
    def __init__(self, 
                 base_dir: str = "logs",
                 enable_rotation: bool = True,
                 retention_days: int = 7):
        """
        Args:
            base_dir: 통계 파일 저장 기본 디렉토리
            enable_rotation: 오래된 파일 자동 삭제 활성화
            retention_days: 통계 파일 보관 기간 (일)
        """
        self.base_dir = Path(base_dir)
        self.enable_rotation = enable_rotation
        self.retention_days = retention_days
        
        # 통계 저장 디렉토리 생성
        self.stats_dir = self.base_dir / "integrated_stats"
        self.stats_dir.mkdir(parents=True, exist_ok=True)
        
        logger.info(f"StatsManager 초기화: base_dir={base_dir}, retention_days={retention_days}")
    
    def collect_all_stats(self, 
                         rate_limiter=None,
                         backoff_strategy=None,
                         error_recovery=None,
                         batch_controller=None,
                         cache=None) -> Dict[str, Any]:
        """모든 모듈의 통계 수집
        
        Args:
            rate_limiter: EnhancedRateLimiter 인스턴스
            backoff_strategy: EnhancedBackoffStrategy 인스턴스
            error_recovery: ErrorRecoverySystem 인스턴스
            batch_controller: DynamicBatchController 인스턴스
            cache: TTLCache 인스턴스
            
        Returns:
            통합된 통계 딕셔너리
        """
        stats = {
            'timestamp': datetime.now().isoformat(),
            'timestamp_epoch': time.time(),
            'modules': {}
        }
        
        # Rate Limiter 통계
        if rate_limiter and hasattr(rate_limiter, 'get_stats'):
            stats['modules']['rate_limiter'] = rate_limiter.get_stats()
        
        # Backoff Strategy 통계
        if backoff_strategy and hasattr(backoff_strategy, 'get_stats'):
            stats['modules']['backoff_strategy'] = backoff_strategy.get_stats()
        
        # Error Recovery 통계
        if error_recovery and hasattr(error_recovery, 'get_error_summary'):
            stats['modules']['error_recovery'] = error_recovery.get_error_summary(hours=24)
        
        # Batch Controller 통계
        if batch_controller and hasattr(batch_controller, 'get_stats'):
            stats['modules']['batch_controller'] = batch_controller.get_stats()
        
        # Cache 통계
        if cache and hasattr(cache, 'get_stats'):
            stats['modules']['cache'] = cache.get_stats()
        
        # 전체 요약
        stats['summary'] = self._generate_summary(stats['modules'])
        
        return stats
    
    def _generate_summary(self, modules: Dict[str, Any]) -> Dict[str, Any]:
        """전체 통계 요약 생성"""
        summary = {
            'total_api_calls': 0,
            'total_errors': 0,
            'overall_error_rate': 0,
            'avg_tps': 0,
            'system_health': 'UNKNOWN',
            'cache_hit_rate': 0,
            'api_calls_saved': 0
        }
        
        # Rate Limiter 정보
        if 'rate_limiter' in modules:
            rl = modules['rate_limiter']
            summary['total_api_calls'] = rl.get('total_calls', 0)
            summary['total_errors'] += rl.get('error_count', 0)
            summary['max_tps'] = rl.get('max_calls_per_second', 0)
        
        # Error Recovery 정보
        if 'error_recovery' in modules:
            er = modules['error_recovery']
            summary['total_errors'] += er.get('total_errors', 0)
        
        # Cache 정보
        if 'cache' in modules:
            cache = modules['cache']
            summary['cache_hit_rate'] = cache.get('hit_rate', 0)
            summary['api_calls_saved'] = cache.get('hit_count', 0)
            summary['cache_entries'] = cache.get('size', 0)
        
        # 전체 에러율 계산
        if summary['total_api_calls'] > 0:
            summary['overall_error_rate'] = summary['total_errors'] / summary['total_api_calls']
        
        # 시스템 상태 판단
        if summary['overall_error_rate'] < 0.01:
            summary['system_health'] = 'HEALTHY'
        elif summary['overall_error_rate'] < 0.05:
            summary['system_health'] = 'WARNING'
        else:
            summary['system_health'] = 'CRITICAL'
        
        return summary
    
    def save_stats(self,
                  stats: Dict[str, Any],
                  format: str = 'json',
                  filename: Optional[str] = None,
                  compress: bool = False,
                  include_timestamp: bool = True) -> str:
        """통계를 파일로 저장
        
        Args:
            stats: 저장할 통계 데이터
            format: 저장 형식 ('json', 'csv', 'jsonl')
            filename: 파일명 (None이면 자동 생성)
            compress: gzip 압축 여부
            include_timestamp: 파일명에 타임스탬프 포함 여부
            
        Returns:
            저장된 파일 경로
        """
        # 파일명 생성
        if filename is None:
            if include_timestamp:
                timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
                filename = f"stats_{timestamp}"
            else:
                filename = "stats_latest"
        
        # 형식별 처리
        if format == 'json':
            filepath = self._save_json(stats, filename, compress)
        elif format == 'csv':
            filepath = self._save_csv(stats, filename, compress)
        elif format == 'jsonl':
            filepath = self._save_jsonl(stats, filename, compress)
        else:
            raise ValueError(f"지원하지 않는 형식: {format}")
        
        # 파일 로테이션
        if self.enable_rotation:
            self._rotate_old_files()
        
        logger.info(f"통계 저장 완료: {filepath}")
        return str(filepath)
    
    def _save_json(self, stats: Dict[str, Any], filename: str, compress: bool) -> Path:
        """JSON 형식으로 저장"""
        ext = '.json.gz' if compress else '.json'
        filepath = self.stats_dir / f"{filename}{ext}"
        
        content = json.dumps(stats, indent=2, ensure_ascii=False)
        
        if compress:
            with gzip.open(filepath, 'wt', encoding='utf-8') as f:
                f.write(content)
        else:
            with open(filepath, 'w', encoding='utf-8') as f:
                f.write(content)
        
        return filepath
    
    def _save_csv(self, stats: Dict[str, Any], filename: str, compress: bool) -> Path:
        """CSV 형식으로 저장 (요약 정보만)"""
        ext = '.csv.gz' if compress else '.csv'
        filepath = self.stats_dir / f"{filename}{ext}"
        
        # CSV로 저장할 데이터 평탄화
        flat_data = self._flatten_stats(stats)
        
        if compress:
            with gzip.open(filepath, 'wt', encoding='utf-8', newline='') as f:
                writer = csv.DictWriter(f, fieldnames=flat_data.keys())
                writer.writeheader()
                writer.writerow(flat_data)
        else:
            with open(filepath, 'w', encoding='utf-8', newline='') as f:
                writer = csv.DictWriter(f, fieldnames=flat_data.keys())
                writer.writeheader()
                writer.writerow(flat_data)
        
        return filepath
    
    def _save_jsonl(self, stats: Dict[str, Any], filename: str, compress: bool) -> Path:
        """JSON Lines 형식으로 저장 (추가 가능)"""
        ext = '.jsonl.gz' if compress else '.jsonl'
        filepath = self.stats_dir / f"{filename}{ext}"
        
        # 한 줄로 저장
        line = json.dumps(stats, ensure_ascii=False) + '\n'
        
        if compress:
            mode = 'at' if filepath.exists() else 'wt'
            with gzip.open(filepath, mode, encoding='utf-8') as f:
                f.write(line)
        else:
            mode = 'a' if filepath.exists() else 'w'
            with open(filepath, mode, encoding='utf-8') as f:
                f.write(line)
        
        return filepath
    
    def _flatten_stats(self, stats: Dict[str, Any], prefix: str = '') -> Dict[str, Any]:
        """중첩된 딕셔너리를 평탄화"""
        flat = {}
        
        for key, value in stats.items():
            new_key = f"{prefix}_{key}" if prefix else key
            
            if isinstance(value, dict):
                flat.update(self._flatten_stats(value, new_key))
            elif isinstance(value, (list, tuple)):
                flat[new_key] = str(value)
            else:
                flat[new_key] = value
        
        return flat
    
    def _rotate_old_files(self):
        """오래된 통계 파일 삭제"""
        if not self.enable_rotation:
            return
        
        cutoff_time = datetime.now() - timedelta(days=self.retention_days)
        deleted_count = 0
        
        for file_path in self.stats_dir.glob("stats_*"):
            # 파일 수정 시간 확인
            mtime = datetime.fromtimestamp(file_path.stat().st_mtime)
            
            if mtime < cutoff_time:
                try:
                    file_path.unlink()
                    deleted_count += 1
                except Exception as e:
                    logger.error(f"파일 삭제 실패 {file_path}: {e}")
        
        if deleted_count > 0:
            logger.info(f"{deleted_count}개의 오래된 통계 파일 삭제됨")
    
    def load_stats(self, filepath: Union[str, Path], format: str = 'auto') -> Dict[str, Any]:
        """저장된 통계 파일 로드
        
        Args:
            filepath: 파일 경로
            format: 파일 형식 ('auto'면 확장자로 판단)
            
        Returns:
            로드된 통계 데이터
        """
        filepath = Path(filepath)
        
        # 형식 자동 감지
        if format == 'auto':
            if filepath.suffix == '.gz':
                # 압축 파일
                base_name = filepath.stem
                if base_name.endswith('.json'):
                    format = 'json'
                elif base_name.endswith('.csv'):
                    format = 'csv'
                elif base_name.endswith('.jsonl'):
                    format = 'jsonl'
            else:
                # 비압축 파일
                format = filepath.suffix[1:]  # .json -> json
        
        # 압축 여부 확인
        compressed = filepath.suffix == '.gz'
        
        # 형식별 로드
        if format == 'json':
            return self._load_json(filepath, compressed)
        elif format == 'csv':
            return self._load_csv(filepath, compressed)
        elif format == 'jsonl':
            return self._load_jsonl(filepath, compressed)
        else:
            raise ValueError(f"지원하지 않는 형식: {format}")
    
    def _load_json(self, filepath: Path, compressed: bool) -> Dict[str, Any]:
        """JSON 파일 로드"""
        if compressed:
            with gzip.open(filepath, 'rt', encoding='utf-8') as f:
                return json.load(f)
        else:
            with open(filepath, 'r', encoding='utf-8') as f:
                return json.load(f)
    
    def _load_csv(self, filepath: Path, compressed: bool) -> List[Dict[str, Any]]:
        """CSV 파일 로드"""
        data = []
        
        if compressed:
            with gzip.open(filepath, 'rt', encoding='utf-8', newline='') as f:
                reader = csv.DictReader(f)
                for row in reader:
                    data.append(row)
        else:
            with open(filepath, 'r', encoding='utf-8', newline='') as f:
                reader = csv.DictReader(f)
                for row in reader:
                    data.append(row)
        
        return data
    
    def _load_jsonl(self, filepath: Path, compressed: bool) -> List[Dict[str, Any]]:
        """JSON Lines 파일 로드"""
        data = []
        
        if compressed:
            with gzip.open(filepath, 'rt', encoding='utf-8') as f:
                for line in f:
                    data.append(json.loads(line.strip()))
        else:
            with open(filepath, 'r', encoding='utf-8') as f:
                for line in f:
                    data.append(json.loads(line.strip()))
        
        return data
    
    def get_latest_stats_file(self, format: str = 'json') -> Optional[Path]:
        """가장 최근 통계 파일 경로 반환"""
        pattern = f"stats_*{format}*"
        files = list(self.stats_dir.glob(pattern))
        
        if not files:
            return None
        
        # 수정 시간 기준 정렬
        files.sort(key=lambda f: f.stat().st_mtime, reverse=True)
        return files[0]


# 싱글톤 인스턴스
_stats_manager = None

def get_stats_manager() -> StatsManager:
    """StatsManager 싱글톤 인스턴스 반환"""
    global _stats_manager
    if _stats_manager is None:
        _stats_manager = StatsManager()
    return _stats_manager 