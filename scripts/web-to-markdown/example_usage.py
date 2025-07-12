#!/usr/bin/env python3
"""
Web to Markdown 변환기 사용 예제
프로그래밍 방식으로 변환기를 사용하는 방법을 보여줍니다.
"""

import asyncio
import sys
import os

# 현재 디렉토리를 Python 경로에 추가
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from web_to_markdown import WebToMarkdownConverter
from simple_web_to_markdown import convert_to_markdown
from ai_web_to_markdown import AIWebToMarkdown


async def example_simple_converter():
    """간단한 변환기 사용 예제"""
    print("=== 간단한 변환기 예제 ===")
    
    # 단일 페이지 변환
    url = "https://example.com"
    await convert_to_markdown(url, "example_output.md")
    print()


async def example_advanced_converter():
    """고급 변환기 사용 예제"""
    print("=== 고급 변환기 예제 ===")
    
    # 변환기 인스턴스 생성
    converter = WebToMarkdownConverter(output_dir="advanced_output")
    
    # 단일 페이지 변환
    results = await converter.crawl_and_convert(
        url="https://example.com",
        max_depth=1,
        include_metadata=True,
        include_links=True,
        include_images=True
    )
    
    print(f"변환 결과: {len(results)} 페이지")
    for url, result in results.items():
        if result['success']:
            print(f"✅ {url}: {result['filename']}")
        else:
            print(f"❌ {url}: {result.get('error', 'Unknown error')}")
    print()


async def example_ai_converter():
    """AI 변환기 사용 예제"""
    print("=== AI 변환기 예제 ===")
    
    # AI 변환기 인스턴스 생성
    ai_converter = AIWebToMarkdown(output_dir="ai_output")
    
    # API 키 확인
    if not os.getenv("OPENAI_API_KEY"):
        print("⚠️ OPENAI_API_KEY가 설정되지 않아 AI 기능을 사용할 수 없습니다.")
        print("export OPENAI_API_KEY='your-api-key' 명령으로 설정하세요.")
        return
    
    # 기사 형식으로 추출
    filename = await ai_converter.smart_crawl_and_convert(
        url="https://example.com",
        extraction_type="article",
        include_raw=False
    )
    
    print(f"AI 변환 완료: {filename}")
    print()


async def example_batch_processing():
    """배치 처리 예제"""
    print("=== 배치 처리 예제 ===")
    
    urls = [
        "https://example.com",
        "https://example.org",
        "https://example.net"
    ]
    
    converter = WebToMarkdownConverter(output_dir="batch_output")
    
    # 모든 URL을 순차적으로 처리
    all_results = {}
    for url in urls:
        print(f"처리 중: {url}")
        results = await converter.crawl_and_convert(url, max_depth=1)
        all_results.update(results)
        
        # 서버 부하 방지를 위한 대기
        await asyncio.sleep(2)
    
    # 결과 요약
    success_count = sum(1 for r in all_results.values() if r['success'])
    print(f"\n배치 처리 완료: {success_count}/{len(all_results)} 성공")
    print()


async def example_custom_processing():
    """커스텀 처리 예제"""
    print("=== 커스텀 처리 예제 ===")
    
    from crawl4ai import AsyncWebCrawler
    
    async with AsyncWebCrawler(verbose=False) as crawler:
        # 커스텀 옵션으로 크롤링
        result = await crawler.arun(
            url="https://example.com",
            word_count_threshold=50,  # 최소 단어 수
            exclude_external_links=True,  # 외부 링크 제외
            remove_overlay=True,  # 오버레이 제거
            bypass_cache=True,  # 캐시 무시
            css_selector="main, article",  # 특정 CSS 선택자만
        )
        
        if result.success:
            # 커스텀 Markdown 생성
            custom_md = f"""# 커스텀 추출 결과

**URL**: {result.url}
**추출 시간**: {result.metadata.get('crawl_date', 'Unknown')}

## 콘텐츠

{result.markdown[:1000]}...

## 통계
- 단어 수: {len(result.cleaned_text.split())}
- 링크 수: {len(result.links)}
- 이미지 수: {len(result.media.get('images', []))}
"""
            
            # 파일로 저장
            with open("custom_output.md", "w", encoding="utf-8") as f:
                f.write(custom_md)
            
            print("✅ 커스텀 변환 완료: custom_output.md")
        else:
            print(f"❌ 크롤링 실패: {result.error_message}")


async def main():
    """모든 예제 실행"""
    print("Web to Markdown 변환기 사용 예제")
    print("=" * 50)
    print()
    
    # 각 예제 실행
    await example_simple_converter()
    await example_advanced_converter()
    await example_ai_converter()
    await example_batch_processing()
    await example_custom_processing()
    
    print("모든 예제 실행 완료!")


if __name__ == "__main__":
    # 비동기 함수 실행
    asyncio.run(main()) 