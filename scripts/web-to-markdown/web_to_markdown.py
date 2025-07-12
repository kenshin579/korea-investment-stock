#!/usr/bin/env python3
"""
웹사이트를 AI 친화적인 Markdown으로 변환하는 스크립트
Crawl4AI를 사용하여 웹페이지의 구조와 내용을 보존하면서 Markdown으로 변환합니다.
"""

import asyncio
import argparse
import os
import sys
from datetime import datetime
from urllib.parse import urlparse, urljoin
import json
from typing import Optional, Dict, Any

from crawl4ai import AsyncWebCrawler
from crawl4ai.extraction_strategy import LLMExtractionStrategy, JsonCssExtractionStrategy
import markdownify


class WebToMarkdownConverter:
    """웹사이트를 AI 친화적인 Markdown으로 변환하는 클래스"""
    
    def __init__(self, output_dir: str = "markdown_output"):
        self.output_dir = output_dir
        os.makedirs(output_dir, exist_ok=True)
        
    async def crawl_and_convert(self, url: str, max_depth: int = 1, 
                               include_metadata: bool = True,
                               include_links: bool = True,
                               include_images: bool = True) -> Dict[str, Any]:
        """
        웹사이트를 크롤링하고 Markdown으로 변환
        
        Args:
            url: 크롤링할 URL
            max_depth: 크롤링 깊이 (기본값: 1)
            include_metadata: 메타데이터 포함 여부
            include_links: 링크 포함 여부
            include_images: 이미지 포함 여부
            
        Returns:
            크롤링 결과 딕셔너리
        """
        results = {}
        
        async with AsyncWebCrawler(verbose=True) as crawler:
            # 메인 페이지 크롤링
            result = await crawler.arun(
                url=url,
                word_count_threshold=10,
                exclude_external_links=True,
                remove_overlay=True,
                process_iframes=True
            )
            
            if result.success:
                # Markdown 변환
                markdown_content = self._convert_to_markdown(
                    result,
                    include_metadata=include_metadata,
                    include_links=include_links,
                    include_images=include_images
                )
                
                # 파일 저장
                filename = self._save_markdown(url, markdown_content)
                
                results[url] = {
                    'success': True,
                    'filename': filename,
                    'content_length': len(markdown_content),
                    'title': result.metadata.get('title', 'No title'),
                    'description': result.metadata.get('description', 'No description')
                }
                
                # 링크된 페이지들도 크롤링 (max_depth에 따라)
                if max_depth > 1 and result.links:
                    for link in result.links[:10]:  # 최대 10개 링크만
                        if link.startswith('http'):
                            sub_result = await self._crawl_subpage(crawler, link, url)
                            if sub_result:
                                results[link] = sub_result
            else:
                results[url] = {
                    'success': False,
                    'error': result.error_message
                }
                
        return results
    
    async def _crawl_subpage(self, crawler, url: str, base_url: str) -> Optional[Dict[str, Any]]:
        """서브페이지 크롤링"""
        try:
            result = await crawler.arun(
                url=url,
                word_count_threshold=10,
                exclude_external_links=True,
                remove_overlay=True
            )
            
            if result.success:
                markdown_content = self._convert_to_markdown(result)
                filename = self._save_markdown(url, markdown_content)
                
                return {
                    'success': True,
                    'filename': filename,
                    'content_length': len(markdown_content),
                    'title': result.metadata.get('title', 'No title')
                }
        except Exception as e:
            return {
                'success': False,
                'error': str(e)
            }
        
        return None
    
    def _convert_to_markdown(self, crawl_result, 
                           include_metadata: bool = True,
                           include_links: bool = True,
                           include_images: bool = True) -> str:
        """크롤링 결과를 AI 친화적인 Markdown으로 변환"""
        markdown_parts = []
        
        # 메타데이터 헤더
        if include_metadata:
            markdown_parts.append("---")
            markdown_parts.append(f"title: {crawl_result.metadata.get('title', 'Untitled')}")
            markdown_parts.append(f"url: {crawl_result.url}")
            markdown_parts.append(f"crawled_at: {datetime.now().isoformat()}")
            if crawl_result.metadata.get('description'):
                markdown_parts.append(f"description: {crawl_result.metadata.get('description')}")
            if crawl_result.metadata.get('keywords'):
                markdown_parts.append(f"keywords: {crawl_result.metadata.get('keywords')}")
            markdown_parts.append("---")
            markdown_parts.append("")
        
        # 제목
        title = crawl_result.metadata.get('title', 'Untitled')
        markdown_parts.append(f"# {title}")
        markdown_parts.append("")
        
        # 메인 콘텐츠
        if crawl_result.markdown:
            # Crawl4AI가 이미 markdown으로 변환한 경우
            content = crawl_result.markdown
        else:
            # HTML을 markdown으로 변환
            content = markdownify.markdownify(
                crawl_result.html,
                heading_style="ATX",
                bullets="-",
                code_language="python"
            )
        
        # 콘텐츠 정리
        content = self._clean_markdown(content)
        markdown_parts.append(content)
        
        # 이미지 목록
        if include_images and crawl_result.media.get('images'):
            markdown_parts.append("\n## 이미지 목록")
            markdown_parts.append("")
            for img in crawl_result.media['images']:
                img_url = img.get('src', '')
                alt_text = img.get('alt', 'Image')
                if img_url:
                    markdown_parts.append(f"- ![{alt_text}]({img_url})")
        
        # 링크 목록
        if include_links and crawl_result.links:
            markdown_parts.append("\n## 참조 링크")
            markdown_parts.append("")
            unique_links = list(set(crawl_result.links))[:20]  # 최대 20개
            for link in unique_links:
                if link.startswith('http'):
                    markdown_parts.append(f"- {link}")
        
        return "\n".join(markdown_parts)
    
    def _clean_markdown(self, content: str) -> str:
        """Markdown 콘텐츠 정리"""
        # 중복된 줄바꿈 제거
        lines = content.split('\n')
        cleaned_lines = []
        empty_count = 0
        
        for line in lines:
            if line.strip() == '':
                empty_count += 1
                if empty_count <= 2:  # 최대 2개의 빈 줄만 허용
                    cleaned_lines.append(line)
            else:
                empty_count = 0
                cleaned_lines.append(line)
        
        return '\n'.join(cleaned_lines)
    
    def _save_markdown(self, url: str, content: str) -> str:
        """Markdown 파일로 저장"""
        # URL에서 파일명 생성
        parsed_url = urlparse(url)
        domain = parsed_url.netloc.replace('.', '_')
        path = parsed_url.path.strip('/').replace('/', '_')
        
        if not path:
            path = 'index'
        
        filename = f"{domain}_{path}.md"
        filepath = os.path.join(self.output_dir, filename)
        
        # 파일명 중복 처리
        counter = 1
        while os.path.exists(filepath):
            filename = f"{domain}_{path}_{counter}.md"
            filepath = os.path.join(self.output_dir, filename)
            counter += 1
        
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        
        print(f"✅ Saved: {filepath}")
        return filepath


async def main():
    parser = argparse.ArgumentParser(
        description="웹사이트를 AI 친화적인 Markdown으로 변환합니다."
    )
    parser.add_argument("url", help="변환할 웹사이트 URL")
    parser.add_argument("-o", "--output", default="markdown_output", 
                       help="출력 디렉토리 (기본값: markdown_output)")
    parser.add_argument("-d", "--depth", type=int, default=1,
                       help="크롤링 깊이 (기본값: 1)")
    parser.add_argument("--no-metadata", action="store_true",
                       help="메타데이터 제외")
    parser.add_argument("--no-links", action="store_true",
                       help="링크 목록 제외")
    parser.add_argument("--no-images", action="store_true",
                       help="이미지 목록 제외")
    
    args = parser.parse_args()
    
    # URL 검증
    if not args.url.startswith(('http://', 'https://')):
        args.url = 'https://' + args.url
    
    print(f"🌐 크롤링 시작: {args.url}")
    print(f"📁 출력 디렉토리: {args.output}")
    print(f"🔍 크롤링 깊이: {args.depth}")
    
    converter = WebToMarkdownConverter(output_dir=args.output)
    
    try:
        results = await converter.crawl_and_convert(
            url=args.url,
            max_depth=args.depth,
            include_metadata=not args.no_metadata,
            include_links=not args.no_links,
            include_images=not args.no_images
        )
        
        # 결과 요약
        print("\n📊 크롤링 결과:")
        success_count = sum(1 for r in results.values() if r['success'])
        print(f"✅ 성공: {success_count}개 페이지")
        
        if len(results) > success_count:
            print(f"❌ 실패: {len(results) - success_count}개 페이지")
        
        # 결과를 JSON 파일로도 저장
        summary_file = os.path.join(args.output, "crawl_summary.json")
        with open(summary_file, 'w', encoding='utf-8') as f:
            json.dump(results, f, ensure_ascii=False, indent=2)
        print(f"\n📄 요약 파일: {summary_file}")
        
    except Exception as e:
        print(f"❌ 오류 발생: {e}")
        sys.exit(1)


def run():
    """명령줄 엔트리포인트를 위한 동기 래퍼"""
    asyncio.run(main())

if __name__ == "__main__":
    run() 