#!/bin/bash

echo "🚀 Web to Markdown 변환기 설치 시작..."
echo ""

# Python 버전 확인
if ! command -v python3 &> /dev/null; then
    echo "❌ Python 3가 설치되어 있지 않습니다."
    echo "먼저 Python 3를 설치해주세요."
    exit 1
fi

echo "✅ Python 3 확인됨: $(python3 --version)"

# pip 확인
if ! command -v pip3 &> /dev/null; then
    echo "❌ pip3가 설치되어 있지 않습니다."
    echo "Python pip를 설치해주세요."
    exit 1
fi

echo "✅ pip3 확인됨"
echo ""

# 가상환경 생성 옵션
echo "가상환경을 생성하시겠습니까? (권장) [Y/n]"
read -r response
response=${response:-Y}

if [[ "$response" =~ ^[Yy]$ ]]; then
    echo "📦 가상환경 생성 중..."
    python3 -m venv venv
    
    # 가상환경 활성화
    if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
        # Windows
        source venv/Scripts/activate
    else
        # Linux/Mac
        source venv/bin/activate
    fi
    
    echo "✅ 가상환경 생성 및 활성화 완료"
    echo ""
fi

# 패키지 설치
echo "📥 필요한 패키지 설치 중..."
# 현재 스크립트 위치 기준으로 작업
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"
pip3 install -e .

if [ $? -eq 0 ]; then
    echo "✅ 패키지 설치 완료"
    
    # AI 기능 사용 옵션
    echo ""
    echo "AI 기능을 위한 추가 패키지를 설치하시겠습니까? (OpenAI API 필요) [y/N]"
    read -r ai_response
    ai_response=${ai_response:-N}
    
    if [[ "$ai_response" =~ ^[Yy]$ ]]; then
        echo "📥 AI 패키지 설치 중..."
        pip3 install -e ".[ai]"
        echo "✅ AI 패키지 설치 완료"
    fi
    
else
    echo "❌ 패키지 설치 실패"
    exit 1
fi

echo ""
echo "🎭 Playwright 브라우저 설치 중..."
playwright install chromium

if [ $? -eq 0 ]; then
    echo "✅ Playwright 브라우저 설치 완료"
else
    echo "❌ Playwright 브라우저 설치 실패"
    echo "수동으로 'playwright install chromium' 명령을 실행해주세요."
fi

echo ""
echo "🎉 설치 완료!"
echo ""
echo "사용 방법:"
echo "1. 간단한 버전: python scripts/simple_web_to_markdown.py <URL>"
echo "2. 고급 버전: python scripts/web_to_markdown.py <URL> [옵션]"
echo "3. AI 버전: python scripts/ai_web_to_markdown.py <URL> [옵션]"
echo ""
echo "자세한 사용법은 scripts/README.md 파일을 참조하세요."

if [[ "$response" =~ ^[Yy]$ ]]; then
    echo ""
    echo "⚠️  가상환경 사용 시 매번 다음 명령으로 활성화하세요:"
    if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
        echo "    source venv/Scripts/activate"
    else
        echo "    source venv/bin/activate"
    fi
fi 