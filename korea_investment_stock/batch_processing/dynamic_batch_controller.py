#!/usr/bin/env python3
"""
Dynamic Batch Controller
Date: 2024-12-28
Issue: #27 - Rate Limiting 개선

Phase 4.2: 에러율에 따른 배치 크기 자동 조정
- 에러율이 높으면 배치 크기 감소
- 안정적이면 배치 크기 증가
- 배치 간 대기 시간도 동적 조정
"""

import time
from typing import Tuple, List, Dict, Any
from collections import deque
import logging

logger = logging.getLogger(__name__)


class DynamicBatchController:
    """동적 배치 조정 컨트롤러
    
    에러율과 성능 지표를 기반으로 배치 크기와 대기 시간을 자동 조정합니다.
    """
    
    def __init__(
        self,
        initial_batch_size: int = 50,
        min_batch_size: int = 5,
        max_batch_size: int = 100,
        initial_batch_delay: float = 1.0,
        min_batch_delay: float = 0.5,
        max_batch_delay: float = 5.0,
        target_error_rate: float = 0.01,  # 목표 에러율 1%
        adjustment_factor: float = 0.2,    # 조정 비율 20%
        history_size: int = 10             # 최근 N개 배치의 성과 추적
    ):
        """
        Args:
            initial_batch_size: 초기 배치 크기
            min_batch_size: 최소 배치 크기
            max_batch_size: 최대 배치 크기
            initial_batch_delay: 초기 배치 간 대기 시간
            min_batch_delay: 최소 대기 시간
            max_batch_delay: 최대 대기 시간
            target_error_rate: 목표 에러율
            adjustment_factor: 조정 비율 (0.2 = 20%씩 조정)
            history_size: 성과 추적을 위한 히스토리 크기
        """
        # 배치 크기 설정
        self.batch_size = initial_batch_size
        self.min_batch_size = min_batch_size
        self.max_batch_size = max_batch_size
        
        # 배치 대기 시간 설정
        self.batch_delay = initial_batch_delay
        self.min_batch_delay = min_batch_delay
        self.max_batch_delay = max_batch_delay
        
        # 목표 및 조정 설정
        self.target_error_rate = target_error_rate
        self.adjustment_factor = adjustment_factor
        
        # 성과 추적
        self.history_size = history_size
        self.batch_history = deque(maxlen=history_size)
        
        # 통계
        self.total_batches = 0
        self.total_items = 0
        self.total_errors = 0
        self.adjustment_count = 0
        self.last_adjustment_time = time.time()
        
        logger.info(
            f"DynamicBatchController 초기화: "
            f"batch_size={initial_batch_size}, batch_delay={initial_batch_delay:.1f}s, "
            f"target_error_rate={target_error_rate:.1%}"
        )
    
    def record_batch_result(
        self, 
        batch_size: int, 
        success_count: int, 
        error_count: int,
        elapsed_time: float
    ) -> None:
        """배치 처리 결과 기록
        
        Args:
            batch_size: 처리한 배치 크기
            success_count: 성공한 요청 수
            error_count: 실패한 요청 수
            elapsed_time: 처리 시간 (초)
        """
        total_count = success_count + error_count
        error_rate = error_count / total_count if total_count > 0 else 0
        throughput = total_count / elapsed_time if elapsed_time > 0 else 0
        
        # 히스토리에 추가
        batch_result = {
            'timestamp': time.time(),
            'batch_size': batch_size,
            'success_count': success_count,
            'error_count': error_count,
            'error_rate': error_rate,
            'elapsed_time': elapsed_time,
            'throughput': throughput
        }
        self.batch_history.append(batch_result)
        
        # 전체 통계 업데이트
        self.total_batches += 1
        self.total_items += total_count
        self.total_errors += error_count
        
        logger.info(
            f"배치 결과: 크기={batch_size}, 성공={success_count}, "
            f"실패={error_count}, 에러율={error_rate:.1%}, "
            f"처리시간={elapsed_time:.2f}s, TPS={throughput:.1f}"
        )
        
        # 동적 조정 수행
        self._adjust_parameters()
    
    def _adjust_parameters(self) -> None:
        """에러율과 성능을 기반으로 파라미터 조정"""
        
        # 충분한 히스토리가 없으면 조정하지 않음
        if len(self.batch_history) < 3:
            return
        
        # 너무 자주 조정하지 않도록 제한 (최소 10초 간격)
        if time.time() - self.last_adjustment_time < 10:
            return
        
        # 최근 성과 분석
        recent_errors = sum(b['error_count'] for b in self.batch_history)
        recent_total = sum(b['success_count'] + b['error_count'] for b in self.batch_history)
        recent_error_rate = recent_errors / recent_total if recent_total > 0 else 0
        
        # 평균 처리량 계산
        avg_throughput = sum(b['throughput'] for b in self.batch_history) / len(self.batch_history)
        
        # 이전 값 저장
        old_batch_size = self.batch_size
        old_batch_delay = self.batch_delay
        
        # 에러율 기반 조정
        if recent_error_rate > self.target_error_rate * 2:
            # 에러율이 목표의 2배 이상: 적극적으로 감소
            self._decrease_batch_size(aggressive=True)
            self._increase_batch_delay(aggressive=True)
            
        elif recent_error_rate > self.target_error_rate:
            # 에러율이 목표 초과: 보수적으로 감소
            self._decrease_batch_size(aggressive=False)
            self._increase_batch_delay(aggressive=False)
            
        elif recent_error_rate < self.target_error_rate * 0.5:
            # 에러율이 목표의 절반 미만: 점진적으로 증가
            self._increase_batch_size()
            self._decrease_batch_delay()
        
        # 조정이 발생했는지 확인
        if self.batch_size != old_batch_size or self.batch_delay != old_batch_delay:
            self.adjustment_count += 1
            self.last_adjustment_time = time.time()
            
            logger.info(
                f"배치 파라미터 조정 #{self.adjustment_count}: "
                f"batch_size: {old_batch_size} → {self.batch_size}, "
                f"batch_delay: {old_batch_delay:.1f}s → {self.batch_delay:.1f}s "
                f"(에러율: {recent_error_rate:.1%}, TPS: {avg_throughput:.1f})"
            )
    
    def _decrease_batch_size(self, aggressive: bool = False) -> None:
        """배치 크기 감소"""
        factor = self.adjustment_factor * (2 if aggressive else 1)
        new_size = int(self.batch_size * (1 - factor))
        self.batch_size = max(new_size, self.min_batch_size)
    
    def _increase_batch_size(self) -> None:
        """배치 크기 증가"""
        new_size = int(self.batch_size * (1 + self.adjustment_factor))
        self.batch_size = min(new_size, self.max_batch_size)
    
    def _increase_batch_delay(self, aggressive: bool = False) -> None:
        """배치 대기 시간 증가"""
        factor = self.adjustment_factor * (2 if aggressive else 1)
        new_delay = self.batch_delay * (1 + factor)
        self.batch_delay = min(new_delay, self.max_batch_delay)
    
    def _decrease_batch_delay(self) -> None:
        """배치 대기 시간 감소"""
        new_delay = self.batch_delay * (1 - self.adjustment_factor)
        self.batch_delay = max(new_delay, self.min_batch_delay)
    
    def get_current_parameters(self) -> Tuple[int, float]:
        """현재 배치 파라미터 반환
        
        Returns:
            (batch_size, batch_delay) 튜플
        """
        return self.batch_size, self.batch_delay
    
    def get_stats(self) -> Dict[str, Any]:
        """통계 정보 반환"""
        overall_error_rate = self.total_errors / self.total_items if self.total_items > 0 else 0
        
        return {
            'current_batch_size': self.batch_size,
            'current_batch_delay': self.batch_delay,
            'total_batches': self.total_batches,
            'total_items': self.total_items,
            'total_errors': self.total_errors,
            'overall_error_rate': overall_error_rate,
            'adjustment_count': self.adjustment_count,
            'target_error_rate': self.target_error_rate,
            'recent_history': list(self.batch_history)
        }
    
    def reset(self, batch_size: int = None, batch_delay: float = None) -> None:
        """파라미터 리셋
        
        Args:
            batch_size: 새로운 배치 크기 (None이면 초기값 유지)
            batch_delay: 새로운 대기 시간 (None이면 초기값 유지)
        """
        if batch_size is not None:
            self.batch_size = max(min(batch_size, self.max_batch_size), self.min_batch_size)
        
        if batch_delay is not None:
            self.batch_delay = max(min(batch_delay, self.max_batch_delay), self.min_batch_delay)
        
        self.batch_history.clear()
        self.last_adjustment_time = time.time()
        
        logger.info(f"배치 파라미터 리셋: batch_size={self.batch_size}, batch_delay={self.batch_delay:.1f}s")
    
    def suggest_initial_parameters(self, total_items: int) -> Tuple[int, float]:
        """전체 항목 수를 기반으로 초기 파라미터 제안
        
        Args:
            total_items: 처리할 전체 항목 수
            
        Returns:
            (suggested_batch_size, suggested_batch_delay) 튜플
        """
        if total_items <= 20:
            # 소량: 한 번에 처리
            return total_items, 0.0
        elif total_items <= 100:
            # 중간: 20개씩
            return 20, 0.5
        elif total_items <= 500:
            # 대량: 50개씩
            return 50, 1.0
        else:
            # 초대량: 100개씩
            return 100, 2.0


# 테스트 코드
if __name__ == "__main__":
    # 로깅 설정
    logging.basicConfig(
        level=logging.INFO,
        format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
    )
    
    # 컨트롤러 생성
    controller = DynamicBatchController(
        initial_batch_size=30,
        target_error_rate=0.02  # 2% 목표
    )
    
    print("Dynamic Batch Controller 테스트")
    print("=" * 50)
    
    # 시뮬레이션: 다양한 에러율로 배치 처리
    scenarios = [
        # (성공, 실패, 처리시간)
        (28, 2, 3.5),   # 정상
        (27, 3, 3.6),   # 정상
        (25, 5, 3.8),   # 에러율 상승
        (20, 10, 4.0),  # 높은 에러율
        (18, 12, 4.2),  # 매우 높은 에러율
        (28, 2, 3.0),   # 회복
        (29, 1, 2.8),   # 안정화
        (30, 0, 2.5),   # 매우 안정
    ]
    
    for i, (success, error, elapsed) in enumerate(scenarios):
        print(f"\n배치 #{i+1}")
        
        # 현재 파라미터 가져오기
        batch_size, batch_delay = controller.get_current_parameters()
        print(f"파라미터: batch_size={batch_size}, batch_delay={batch_delay:.1f}s")
        
        # 배치 처리 시뮬레이션
        time.sleep(0.5)  # 실제 처리 시뮬레이션
        
        # 결과 기록
        controller.record_batch_result(
            batch_size=batch_size,
            success_count=success,
            error_count=error,
            elapsed_time=elapsed
        )
    
    # 최종 통계
    print("\n" + "=" * 50)
    print("최종 통계:")
    stats = controller.get_stats()
    print(f"- 총 배치 수: {stats['total_batches']}")
    print(f"- 총 처리 항목: {stats['total_items']}")
    print(f"- 전체 에러율: {stats['overall_error_rate']:.1%}")
    print(f"- 파라미터 조정 횟수: {stats['adjustment_count']}")
    print(f"- 최종 batch_size: {stats['current_batch_size']}")
    print(f"- 최종 batch_delay: {stats['current_batch_delay']:.1f}s") 