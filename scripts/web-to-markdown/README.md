# Web to Markdown Converter Scripts

Crawl4AI를 활용하여 웹사이트를 AI 친화적인 Markdown으로 변환하는 스크립트 모음입니다.

## 🚀 설치

### 빠른 설치 (권장)

```bash
# 설치 스크립트 실행
bash scripts/web-to-markdown/install.sh
```

### 수동 설치

#### 1. 필요한 패키지 설치

```bash
cd scripts/web-to-markdown
pip install -e .
```

AI 기능 포함 설치:

```bash
cd scripts/web-to-markdown
pip install -e ".[ai]"
```

개발 도구 포함 설치:

```bash
cd scripts/web-to-markdown
pip install -e ".[dev]"
```

#### 2. Playwright 브라우저 설치 (Crawl4AI 필수 요소)

```bash
playwright install chromium
```

#### 3. (선택) AI 버전 사용 시 OpenAI API 키 설정

```bash
export OPENAI_API_KEY='your-api-key'
```

## 📁 스크립트 설명

### 1. `web_to_markdown.py` - 고급 기능 버전

다양한 옵션과 기능을 제공하는 메인 스크립트입니다.

**주요 기능:**
- 멀티 페이지 크롤링 (깊이 설정 가능)
- 메타데이터 포함 (제목, 설명, 키워드 등)
- 이미지 및 링크 목록 추출
- JSON 형태의 크롤링 요약 생성
- 중복 파일명 자동 처리

**사용법:**

```bash
# 기본 사용 (단일 페이지)
python scripts/web-to-markdown/web_to_markdown.py https://example.com

# 출력 디렉토리 지정
python scripts/web-to-markdown/web_to_markdown.py https://example.com -o my_output

# 깊이 2로 크롤링 (메인 페이지 + 링크된 페이지들)
python scripts/web-to-markdown/web_to_markdown.py https://example.com -d 2

# 메타데이터 제외
python scripts/web-to-markdown/web_to_markdown.py https://example.com --no-metadata

# 링크 목록 제외
python scripts/web-to-markdown/web_to_markdown.py https://example.com --no-links

# 이미지 목록 제외
python scripts/web-to-markdown/web_to_markdown.py https://example.com --no-images

# 모든 옵션 조합
python scripts/web-to-markdown/web_to_markdown.py https://example.com -o output -d 2 --no-images
```

### 2. `simple_web_to_markdown.py` - 간단한 버전

빠른 사용을 위한 단순화된 버전입니다.

**주요 기능:**
- 단일 페이지 크롤링
- 기본 메타데이터 포함
- 간단한 사용법

**사용법:**

```bash
# 기본 사용 (자동 파일명 생성)
python scripts/web-to-markdown/simple_web_to_markdown.py https://example.com

# 출력 파일명 지정
python scripts/web-to-markdown/simple_web_to_markdown.py https://example.com my_output.md
```

### 3. `ai_web_to_markdown.py` - AI 강화 버전

OpenAI GPT-4를 활용하여 콘텐츠를 더 스마트하게 추출하고 구조화합니다.

**주요 기능:**
- AI 기반 콘텐츠 구조화
- 타입별 특화 추출 (기사, 제품, 문서, 일반)
- 자동 요약 및 핵심 포인트 추출
- 스키마 기반 정보 추출

**사용법:**

```bash
# OpenAI API 키 설정 (필수)
export OPENAI_API_KEY='your-api-key'

# 기본 사용 (일반 타입)
python scripts/web-to-markdown/ai_web_to_markdown.py https://example.com

# 기사 타입으로 추출
python scripts/web-to-markdown/ai_web_to_markdown.py https://news.example.com/article -t article

# 제품 페이지 추출
python scripts/web-to-markdown/ai_web_to_markdown.py https://shop.example.com/product -t product

# 문서 페이지 추출
python scripts/web-to-markdown/ai_web_to_markdown.py https://docs.example.com -t documentation

# 원본 콘텐츠 포함
python scripts/web-to-markdown/ai_web_to_markdown.py https://example.com --include-raw
```

**추출 타입:**
- `article`: 뉴스, 블로그 포스트 (작성자, 날짜, 요약, 핵심 포인트)
- `product`: 제품 페이지 (가격, 특징, 사양, 리뷰 요약)
- `documentation`: 기술 문서 (개요, 섹션, 코드 예제, 사전 요구사항)
- `general`: 일반 웹페이지 (기본값)

## 📄 출력 형식

### Markdown 파일 구조

```markdown
---
title: 페이지 제목
url: https://example.com
crawled_at: 2024-01-01T12:00:00
description: 페이지 설명
keywords: 키워드1, 키워드2
---

# 페이지 제목

[페이지 내용이 여기에 들어갑니다]

## 이미지 목록

- ![이미지 설명](이미지_URL)

## 참조 링크

- https://example.com/link1
- https://example.com/link2
```

### JSON 요약 파일 (crawl_summary.json)

```json
{
  "https://example.com": {
    "success": true,
    "filename": "example_com_index.md",
    "content_length": 12345,
    "title": "Example Domain",
    "description": "Example Domain for illustrative examples"
  }
}
```

## 🎯 사용 예시

### 패키지 설치 후 사용 (권장)

설치 후에는 어디서든 명령어를 사용할 수 있습니다:

```bash
# 설치
cd scripts/web-to-markdown && pip install -e . && cd ../..

# 사용
simple-web-to-markdown https://example.com
web-to-markdown https://example.com -d 2
ai-web-to-markdown https://example.com -t article
```

### 직접 스크립트 실행

```bash
python scripts/web-to-markdown/simple_web_to_markdown.py https://example.com
python scripts/web-to-markdown/web_to_markdown.py https://example.com -d 2
python scripts/web-to-markdown/ai_web_to_markdown.py https://example.com -t article
```

### 1. 블로그 백업

```bash
# 패키지 설치된 경우
web-to-markdown https://myblog.com -d 2 -o blog_backup

# 직접 실행
python scripts/web-to-markdown/web_to_markdown.py https://myblog.com -d 2 -o blog_backup
```

### 2. 문서 사이트 크롤링

```bash
python scripts/web_to_markdown.py https://docs.example.com -d 3 -o documentation
```

### 3. 단일 페이지 빠른 변환

```bash
python scripts/simple_web_to_markdown.py https://news.example.com/article
```

## ⚠️ 주의사항

1. **크롤링 예절**: 대상 사이트의 robots.txt를 확인하고 과도한 요청을 피하세요.
2. **저작권**: 크롤링한 콘텐츠의 저작권을 확인하고 적절히 사용하세요.
3. **리소스 사용**: 깊이가 깊을수록 많은 페이지를 크롤링하므로 시간과 리소스가 많이 소요됩니다.
4. **동적 콘텐츠**: JavaScript로 렌더링되는 콘텐츠도 Crawl4AI가 처리하지만, 일부 복잡한 SPA는 완벽히 캡처되지 않을 수 있습니다.

## 🔧 문제 해결

### 1. "No module named 'crawl4ai'" 오류

```bash
pip install crawl4ai --upgrade
```

### 2. Playwright 관련 오류

```bash
playwright install chromium
playwright install-deps  # Linux에서 필요한 경우
```

### 3. 인코딩 오류

스크립트는 UTF-8 인코딩을 사용합니다. 다른 인코딩이 필요한 경우 코드에서 `encoding` 파라미터를 수정하세요.

## 📚 참고 자료

- [Crawl4AI GitHub](https://github.com/unclecode/crawl4ai)
- [Crawl4AI Documentation](https://crawl4ai.com/mkdocs/)
- [Markdownify Documentation](https://github.com/matthewwithanm/python-markdownify)

## 🤝 기여

이슈나 개선사항이 있다면 언제든 PR을 보내주세요! 