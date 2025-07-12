# Scripts 디렉토리

이 디렉토리는 프로젝트에서 사용되는 다양한 유틸리티 스크립트들을 포함합니다.

## 📁 하위 디렉토리

### web-to-markdown/
웹사이트를 AI 친화적인 Markdown으로 변환하는 도구 모음입니다.

- **주요 기능**: 
  - 웹 페이지를 구조화된 Markdown으로 변환
  - AI 기반 콘텐츠 추출 및 요약
  - 배치 처리 및 멀티 페이지 크롤링 지원

- **자세한 사용법**: `scripts/web-to-markdown/README.md` 참조

- **빠른 시작**:
  ```bash
  # 설치
  bash scripts/web-to-markdown/install.sh
  
  # 사용
  python scripts/web-to-markdown/simple_web_to_markdown.py https://example.com
  ``` 