#!/usr/bin/env python3
"""
간단한 웹 to Markdown 변환 스크립트
빠른 테스트와 간단한 사용을 위한 버전
"""

import asyncio
from crawl4ai import AsyncWebCrawler
import sys
import os


async def convert_to_markdown(url: str, output_file: str = None):
    """웹페이지를 Markdown으로 변환하는 간단한 함수"""
    
    async with AsyncWebCrawler(verbose=True) as crawler:
        # 크롤링 실행
        result = await crawler.arun(
            url=url,
            word_count_threshold=10,
            exclude_external_links=True,
            remove_overlay=True,
            bypass_cache=True
        )
        
        if result.success:
            # 출력 파일명 결정
            if not output_file:
                # URL에서 파일명 생성
                safe_filename = url.replace('https://', '').replace('http://', '')
                safe_filename = safe_filename.replace('/', '_').replace('?', '_')
                output_file = f"{safe_filename}.md"
            
            # Markdown 내용 구성
            markdown_content = f"""---
title: {result.metadata.get('title', 'Untitled')}
url: {url}
---

# {result.metadata.get('title', 'Untitled')}

{result.markdown if result.markdown else result.cleaned_text}

---

## 메타 정보

- **URL**: {url}
- **설명**: {result.metadata.get('description', 'N/A')}
- **크롤링 시간**: {result.metadata.get('crawl_date', 'N/A')}
"""
            
            # 파일 저장
            with open(output_file, 'w', encoding='utf-8') as f:
                f.write(markdown_content)
            
            print(f"✅ 성공적으로 저장됨: {output_file}")
            print(f"📄 파일 크기: {len(markdown_content):,} 바이트")
            
            return True
        else:
            print(f"❌ 크롤링 실패: {result.error_message}")
            return False


def main():
    if len(sys.argv) < 2:
        print("사용법: python simple_web_to_markdown.py <URL> [출력파일명]")
        print("예시: python simple_web_to_markdown.py https://example.com")
        print("      python simple_web_to_markdown.py https://example.com output.md")
        sys.exit(1)
    
    url = sys.argv[1]
    output_file = sys.argv[2] if len(sys.argv) > 2 else None
    
    # URL 검증
    if not url.startswith(('http://', 'https://')):
        url = 'https://' + url
    
    print(f"🌐 크롤링 중: {url}")
    
    # 비동기 함수 실행
    success = asyncio.run(convert_to_markdown(url, output_file))
    
    if not success:
        sys.exit(1)


def run():
    """명령줄 엔트리포인트를 위한 동기 래퍼"""
    main()

if __name__ == "__main__":
    run() 