#!/usr/bin/env python3
"""
AI 기반 웹 to Markdown 변환 스크립트
LLM을 활용하여 더 스마트한 콘텐츠 추출 및 구조화
"""

import asyncio
import argparse
import os
import json
from datetime import datetime
from typing import Dict, Any, List
from crawl4ai import AsyncWebCrawler
from crawl4ai.extraction_strategy import LLMExtractionStrategy


class AIWebToMarkdown:
    """AI를 활용한 웹 to Markdown 변환기"""
    
    def __init__(self, output_dir: str = "ai_markdown_output"):
        self.output_dir = output_dir
        os.makedirs(output_dir, exist_ok=True)
        
    async def extract_with_ai(self, url: str, extraction_type: str = "article") -> Dict[str, Any]:
        """AI를 사용하여 웹페이지에서 구조화된 콘텐츠 추출"""
        
        # 추출 타입별 스키마 정의
        schemas = {
            "article": {
                "title": "string",
                "author": "string",
                "date": "string",
                "summary": "string",
                "main_content": "string",
                "key_points": ["string"],
                "tags": ["string"]
            },
            "product": {
                "name": "string",
                "price": "string",
                "description": "string",
                "features": ["string"],
                "specifications": {"key": "value"},
                "reviews_summary": "string"
            },
            "documentation": {
                "title": "string",
                "overview": "string",
                "sections": [
                    {
                        "heading": "string",
                        "content": "string",
                        "code_examples": ["string"]
                    }
                ],
                "prerequisites": ["string"],
                "related_topics": ["string"]
            },
            "general": {
                "title": "string",
                "main_topic": "string",
                "summary": "string",
                "key_sections": [
                    {
                        "heading": "string",
                        "content": "string"
                    }
                ],
                "important_links": ["string"],
                "metadata": {"key": "value"}
            }
        }
        
        # LLM 추출 전략 설정
        extraction_strategy = LLMExtractionStrategy(
            provider="openai/gpt-4",  # 실제 사용시 API 키 필요
            api_token=os.getenv("OPENAI_API_KEY"),
            schema=schemas.get(extraction_type, schemas["general"]),
            extraction_type="schema",
            instruction=f"""
            웹페이지에서 {extraction_type} 형식의 콘텐츠를 추출하세요.
            가능한 한 자세하고 구조화된 정보를 제공하되,
            원본 콘텐츠의 의미를 정확히 보존하세요.
            코드 예제가 있다면 반드시 포함시키세요.
            """
        )
        
        async with AsyncWebCrawler(verbose=True) as crawler:
            result = await crawler.arun(
                url=url,
                extraction_strategy=extraction_strategy,
                bypass_cache=True
            )
            
            return result
    
    async def smart_crawl_and_convert(self, url: str, 
                                    extraction_type: str = "general",
                                    include_raw: bool = False) -> str:
        """AI를 활용한 스마트 크롤링 및 변환"""
        
        print(f"🤖 AI 기반 콘텐츠 추출 중... (타입: {extraction_type})")
        
        # 기본 크롤링 (백업용)
        async with AsyncWebCrawler(verbose=True) as crawler:
            basic_result = await crawler.arun(
                url=url,
                word_count_threshold=10,
                exclude_external_links=True,
                remove_overlay=True
            )
        
        # AI 추출 시도
        ai_extracted = None
        try:
            if os.getenv("OPENAI_API_KEY"):
                ai_result = await self.extract_with_ai(url, extraction_type)
                if ai_result.success and ai_result.extracted_content:
                    ai_extracted = json.loads(ai_result.extracted_content)
        except Exception as e:
            print(f"⚠️ AI 추출 실패, 기본 추출 사용: {e}")
        
        # Markdown 생성
        markdown_content = self._create_ai_enhanced_markdown(
            basic_result,
            ai_extracted,
            extraction_type,
            include_raw
        )
        
        # 파일 저장
        filename = self._save_markdown(url, markdown_content, extraction_type)
        
        return filename
    
    def _create_ai_enhanced_markdown(self, basic_result, ai_data, 
                                   extraction_type: str, include_raw: bool) -> str:
        """AI 추출 데이터를 활용한 향상된 Markdown 생성"""
        
        parts = []
        
        # 헤더
        parts.append("---")
        parts.append(f"title: {basic_result.metadata.get('title', 'Untitled')}")
        parts.append(f"url: {basic_result.url}")
        parts.append(f"extracted_at: {datetime.now().isoformat()}")
        parts.append(f"extraction_type: {extraction_type}")
        parts.append(f"ai_enhanced: {'yes' if ai_data else 'no'}")
        parts.append("---")
        parts.append("")
        
        # AI 추출 데이터가 있는 경우
        if ai_data:
            parts.append("# " + ai_data.get('title', basic_result.metadata.get('title', 'Untitled')))
            parts.append("")
            
            # 타입별 특수 처리
            if extraction_type == "article":
                self._format_article(parts, ai_data)
            elif extraction_type == "product":
                self._format_product(parts, ai_data)
            elif extraction_type == "documentation":
                self._format_documentation(parts, ai_data)
            else:
                self._format_general(parts, ai_data)
        else:
            # AI 데이터가 없는 경우 기본 포맷
            parts.append("# " + basic_result.metadata.get('title', 'Untitled'))
            parts.append("")
            if basic_result.markdown:
                parts.append(basic_result.markdown)
            else:
                parts.append(basic_result.cleaned_text)
        
        # 원본 콘텐츠 포함 옵션
        if include_raw and not ai_data:
            parts.append("\n---\n")
            parts.append("## 원본 콘텐츠")
            parts.append("")
            parts.append("```")
            parts.append(basic_result.cleaned_text[:5000])  # 처음 5000자만
            parts.append("```")
        
        # 메타데이터
        parts.append("\n---\n")
        parts.append("## 페이지 정보")
        parts.append("")
        parts.append(f"- **URL**: {basic_result.url}")
        parts.append(f"- **제목**: {basic_result.metadata.get('title', 'N/A')}")
        parts.append(f"- **설명**: {basic_result.metadata.get('description', 'N/A')}")
        
        return "\n".join(parts)
    
    def _format_article(self, parts: List[str], data: Dict[str, Any]):
        """기사 형식 포맷팅"""
        if data.get('author'):
            parts.append(f"**작성자**: {data['author']}")
        if data.get('date'):
            parts.append(f"**날짜**: {data['date']}")
        parts.append("")
        
        if data.get('summary'):
            parts.append("## 요약")
            parts.append("")
            parts.append(data['summary'])
            parts.append("")
        
        if data.get('main_content'):
            parts.append("## 본문")
            parts.append("")
            parts.append(data['main_content'])
            parts.append("")
        
        if data.get('key_points'):
            parts.append("## 핵심 포인트")
            parts.append("")
            for point in data['key_points']:
                parts.append(f"- {point}")
            parts.append("")
        
        if data.get('tags'):
            parts.append("**태그**: " + ", ".join(f"`{tag}`" for tag in data['tags']))
    
    def _format_product(self, parts: List[str], data: Dict[str, Any]):
        """제품 정보 포맷팅"""
        if data.get('price'):
            parts.append(f"**가격**: {data['price']}")
        parts.append("")
        
        if data.get('description'):
            parts.append("## 제품 설명")
            parts.append("")
            parts.append(data['description'])
            parts.append("")
        
        if data.get('features'):
            parts.append("## 주요 특징")
            parts.append("")
            for feature in data['features']:
                parts.append(f"- {feature}")
            parts.append("")
        
        if data.get('specifications'):
            parts.append("## 사양")
            parts.append("")
            parts.append("| 항목 | 내용 |")
            parts.append("|------|------|")
            for key, value in data['specifications'].items():
                parts.append(f"| {key} | {value} |")
            parts.append("")
    
    def _format_documentation(self, parts: List[str], data: Dict[str, Any]):
        """문서 형식 포맷팅"""
        if data.get('overview'):
            parts.append("## 개요")
            parts.append("")
            parts.append(data['overview'])
            parts.append("")
        
        if data.get('prerequisites'):
            parts.append("## 사전 요구사항")
            parts.append("")
            for prereq in data['prerequisites']:
                parts.append(f"- {prereq}")
            parts.append("")
        
        if data.get('sections'):
            for section in data['sections']:
                parts.append(f"## {section['heading']}")
                parts.append("")
                parts.append(section['content'])
                
                if section.get('code_examples'):
                    parts.append("")
                    for code in section['code_examples']:
                        parts.append("```")
                        parts.append(code)
                        parts.append("```")
                        parts.append("")
    
    def _format_general(self, parts: List[str], data: Dict[str, Any]):
        """일반 형식 포맷팅"""
        if data.get('main_topic'):
            parts.append(f"**주제**: {data['main_topic']}")
            parts.append("")
        
        if data.get('summary'):
            parts.append("## 요약")
            parts.append("")
            parts.append(data['summary'])
            parts.append("")
        
        if data.get('key_sections'):
            for section in data['key_sections']:
                parts.append(f"## {section['heading']}")
                parts.append("")
                parts.append(section['content'])
                parts.append("")
    
    def _save_markdown(self, url: str, content: str, extraction_type: str) -> str:
        """Markdown 파일 저장"""
        from urllib.parse import urlparse
        
        parsed_url = urlparse(url)
        domain = parsed_url.netloc.replace('.', '_')
        path = parsed_url.path.strip('/').replace('/', '_') or 'index'
        
        filename = f"{domain}_{path}_{extraction_type}.md"
        filepath = os.path.join(self.output_dir, filename)
        
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        
        print(f"✅ 저장 완료: {filepath}")
        return filepath


async def main():
    parser = argparse.ArgumentParser(
        description="AI를 활용한 웹사이트 to Markdown 변환기"
    )
    parser.add_argument("url", help="변환할 웹사이트 URL")
    parser.add_argument("-t", "--type", 
                       choices=["article", "product", "documentation", "general"],
                       default="general",
                       help="추출 타입 (기본값: general)")
    parser.add_argument("-o", "--output", default="ai_markdown_output",
                       help="출력 디렉토리")
    parser.add_argument("--include-raw", action="store_true",
                       help="원본 콘텐츠 포함")
    
    args = parser.parse_args()
    
    # URL 검증
    if not args.url.startswith(('http://', 'https://')):
        args.url = 'https://' + args.url
    
    print(f"🌐 URL: {args.url}")
    print(f"🤖 추출 타입: {args.type}")
    print(f"📁 출력 디렉토리: {args.output}")
    
    # API 키 확인
    if not os.getenv("OPENAI_API_KEY"):
        print("⚠️ 주의: OPENAI_API_KEY가 설정되지 않아 기본 추출만 사용합니다.")
        print("AI 기반 추출을 사용하려면 환경 변수를 설정하세요:")
        print("export OPENAI_API_KEY='your-api-key'")
    
    converter = AIWebToMarkdown(output_dir=args.output)
    
    try:
        filename = await converter.smart_crawl_and_convert(
            url=args.url,
            extraction_type=args.type,
            include_raw=args.include_raw
        )
        print(f"\n✨ 변환 완료!")
        
    except Exception as e:
        print(f"❌ 오류 발생: {e}")
        import traceback
        traceback.print_exc()


def run():
    """명령줄 엔트리포인트를 위한 동기 래퍼"""
    asyncio.run(main())

if __name__ == "__main__":
    run() 