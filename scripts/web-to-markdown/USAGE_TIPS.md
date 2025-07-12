# Web to Markdown 변환기 사용 팁 및 모범 사례

## 🎯 각 스크립트 선택 가이드

### 언제 `simple_web_to_markdown.py`를 사용할까?
- ✅ 빠른 테스트가 필요할 때
- ✅ 단일 페이지만 변환하면 될 때
- ✅ 기본적인 텍스트 추출만 필요할 때
- ✅ 설정이나 옵션 없이 바로 사용하고 싶을 때

### 언제 `web_to_markdown.py`를 사용할까?
- ✅ 여러 페이지를 한 번에 크롤링해야 할 때
- ✅ 이미지와 링크 목록이 필요할 때
- ✅ 메타데이터를 포함/제외하고 싶을 때
- ✅ 크롤링 결과를 JSON으로도 저장하고 싶을 때
- ✅ 블로그나 문서 사이트 전체를 백업할 때

### 언제 `ai_web_to_markdown.py`를 사용할까?
- ✅ 콘텐츠를 구조화된 형태로 추출하고 싶을 때
- ✅ 자동 요약이 필요할 때
- ✅ 기사, 제품, 문서 등 특정 형식에 맞춰 추출할 때
- ✅ 핵심 포인트만 뽑아내고 싶을 때
- ❌ OpenAI API 키가 없거나 비용이 부담될 때는 피하세요

## 💡 실전 사용 예시

### 1. 뉴스 기사 스크랩
```bash
# AI 버전 (추천)
python scripts/web-to-markdown/ai_web_to_markdown.py https://news.site.com/article-123 -t article

# 일반 버전
python scripts/web-to-markdown/web_to_markdown.py https://news.site.com/article-123 --no-images
```

### 2. 제품 정보 수집
```bash
# AI 버전으로 구조화된 정보 추출
python scripts/web-to-markdown/ai_web_to_markdown.py https://shop.com/product/laptop -t product -o products

# 이미지 포함한 전체 페이지 저장
python scripts/web-to-markdown/web_to_markdown.py https://shop.com/product/laptop -o products
```

### 3. 기술 문서 백업
```bash
# 문서 사이트 전체 크롤링 (깊이 3)
python scripts/web-to-markdown/web_to_markdown.py https://docs.example.com -d 3 -o tech_docs

# 특정 페이지만 AI로 구조화
python scripts/web-to-markdown/ai_web_to_markdown.py https://docs.example.com/api-guide -t documentation
```

### 4. 블로그 백업
```bash
# 블로그 전체 백업 (링크된 페이지 포함)
python scripts/web-to-markdown/web_to_markdown.py https://myblog.com -d 2 -o blog_backup

# 메타데이터 없이 본문만
python scripts/web-to-markdown/web_to_markdown.py https://myblog.com --no-metadata --no-links
```

## 🛠️ 고급 사용법

### 배치 처리 스크립트 예시

```bash
#!/bin/bash
# batch_convert.sh

# URL 목록 파일에서 일괄 변환
while IFS= read -r url; do
    echo "Processing: $url"
    python scripts/web-to-markdown/simple_web_to_markdown.py "$url"
    sleep 2  # 서버 부하 방지를 위한 대기
done < urls.txt
```

### Python에서 직접 사용하기

```python
# 다른 Python 스크립트에서 import하여 사용
import asyncio
import sys
sys.path.append('scripts/web-to-markdown')
from web_to_markdown import WebToMarkdownConverter

async def batch_convert():
    converter = WebToMarkdownConverter(output_dir="my_output")
    urls = [
        "https://example1.com",
        "https://example2.com",
        "https://example3.com"
    ]
    
    for url in urls:
        results = await converter.crawl_and_convert(url)
        print(f"Converted {url}: {results}")

asyncio.run(batch_convert())
```

## ⚡ 성능 최적화 팁

1. **깊이 제한**: 크롤링 깊이(`-d`)를 필요한 만큼만 설정하세요
   ```bash
   # 메인 페이지만
   python scripts/web-to-markdown/web_to_markdown.py https://example.com -d 1
   
   # 링크된 페이지까지
   python scripts/web-to-markdown/web_to_markdown.py https://example.com -d 2
   ```

2. **불필요한 요소 제외**: 필요없는 요소는 제외하여 속도 향상
   ```bash
   python scripts/web-to-markdown/web_to_markdown.py https://example.com --no-images --no-links
   ```

3. **동시 처리**: 여러 URL을 처리할 때는 병렬 처리 고려
   ```bash
   # GNU Parallel 사용 예시
   cat urls.txt | parallel -j 4 python scripts/web-to-markdown/simple_web_to_markdown.py {}
   ```

## 🐛 문제 해결

### 일반적인 오류와 해결법

1. **"Connection refused" 오류**
   - 원인: 사이트가 봇을 차단
   - 해결: User-Agent 헤더 추가 또는 속도 제한

2. **"Timeout" 오류**
   - 원인: 페이지 로딩이 너무 오래 걸림
   - 해결: JavaScript가 많은 사이트는 로딩 시간 증가 필요

3. **인코딩 오류**
   - 원인: 특수 문자 처리 문제
   - 해결: UTF-8 인코딩 확인

4. **메모리 부족**
   - 원인: 너무 많은 페이지를 한 번에 크롤링
   - 해결: 깊이 제한 또는 배치 처리

## 📋 체크리스트

크롤링 전 확인사항:
- [ ] robots.txt 확인했나요?
- [ ] 사이트 이용약관을 확인했나요?
- [ ] 적절한 딜레이를 설정했나요?
- [ ] 저작권을 고려했나요?
- [ ] 출력 디렉토리에 충분한 공간이 있나요?

## 🔍 디버깅 모드

더 자세한 로그를 보려면:
```bash
# Crawl4AI verbose 모드는 기본적으로 켜져 있음
# 추가 디버깅이 필요한 경우 코드 수정 필요
```

## 📚 추가 리소스

- [Crawl4AI 공식 문서](https://crawl4ai.com/mkdocs/)
- [Markdown 문법 가이드](https://www.markdownguide.org/)
- [Python asyncio 가이드](https://docs.python.org/3/library/asyncio.html)

---

💡 **Pro Tip**: 대량의 페이지를 크롤링할 때는 먼저 작은 샘플로 테스트하세요! 