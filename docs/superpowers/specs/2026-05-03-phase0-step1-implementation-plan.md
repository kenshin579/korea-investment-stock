# Phase 0 Step 1 — Python 정리 + Go 스켈레톤 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** `korea-investment-stock` repo 의 Python 코드를 deprecation notice 부착 후 정리하고, Go 라이브러리 빈 스켈레톤을 main 에 마련한다.

**Architecture:** 외부 PyPI 사용자에게 충격이 가지 않도록 (1) deprecation notice 가 들어간 마지막 Python release (`v0.19.0`) 를 먼저 PyPI 에 공개 → (2) 마지막 Python 커밋에 별칭 태그 (`python-final`) 박기 → (3) main 에서 Python 코드 제거 + Go 스켈레톤 추가 → (4) PR.

**Tech Stack:** Python 3.11+ (정리 대상), Go 1.23+ (신규), git, gh CLI

**참고 spec:** [`2026-05-03-korea-investment-go-migration-design.md`](./2026-05-03-korea-investment-go-migration-design.md) (커밋 `ebfa2e7`)

---

## 사전 정보 / 결정 사항

| 항목 | 값 |
|------|---|
| 현재 최신 태그 | `v0.18.0` |
| 마지막 Python release (deprecation notice 포함) | **`v0.19.0`** (이 plan 에서 새로 만듦) |
| 별칭 태그 (영구 참조용) | **`python-final`** → `v0.19.0` 와 동일 commit |
| Python 코드 제거 commit | `python-final` 태그 다음 commit |
| Go module path | `github.com/kenshin579/korea-investment-stock` |
| Go 패키지명 | `kis` |
| 작업 브랜치 | `docs/golang-migration-spec` (이미 존재, spec 커밋 완료) |
| PR 베이스 | `main` |

## 파일 구조 (변경 대상)

### 수정 (Group A — deprecation notice)
- `pyproject.toml` — description 변경, classifiers 에 `Development Status :: 7 - Inactive` 추가
- `README.md` — 상단에 deprecation 경고 박스 추가 (Group A 단계, 임시. Group F 에서 Go 용으로 완전 교체)
- `CHANGELOG.md` — `[Unreleased]` 섹션을 `## [v0.19.0]` 로 변환 + deprecation notice entry 추가

### 삭제 (Group D — Python 코드 정리)
- `korea_investment_stock/` (전체 디렉터리)
- `korea_investment_stock.egg-info/` (전체)
- `examples/` (Python 예시)
- `scripts/` (Python 스크립트)
- `upload.sh` (PyPI 업로드 스크립트)
- `pyproject.toml`
- `Makefile` (Python 작업용)
- `TEST_REPORT.md`

### 신규 (Group E — Go 스켈레톤)
- `go.mod`
- `client.go` — `Client` 구조체와 `NewClient` 시그니처 (스텁)
- `domestic/doc.go` — 패키지 doc comment 만 (구현은 Phase 1)
- `overseas/doc.go` — 동일
- `internal/httpclient/doc.go` — 패키지 doc comment
- `internal/ratelimit/doc.go` — 동일
- `internal/token/doc.go` — 동일
- `.gitignore` — Go 용 (`*.exe`, `vendor/`, etc.)

### 신규/교체 (Group F — Go README)
- `README.md` — 완전 교체. Go 사용법 안내 + Python v1.x 사용자 마이그레이션 안내 + WIP 표기

### 보존
- `LICENSE`
- `CHANGELOG.md` (Python 히스토리 보존, Go 항목은 향후 추가)
- `docs/` (api 문서, superpowers spec 모두 보존)
- `.github/` (워크플로우는 추후 Go 용으로 갱신, 본 plan 에서는 그대로)

---

## Task 1: pyproject.toml 에 deprecation notice 부착

**Files:**
- Modify: `pyproject.toml`

- [ ] **Step 1: 현재 pyproject.toml 의 description 라인 확인**

Run:
```bash
grep -n 'description' pyproject.toml
```
Expected: `description = "Pure Python wrapper for Korea Investment Securities OpenAPI"` 라인이 있음

- [ ] **Step 2: description 을 deprecation 명시 문구로 교체**

`pyproject.toml` 의 `description` 을 다음으로 교체:
```toml
description = "[DEPRECATED — use the Go version] Python wrapper for Korea Investment Securities OpenAPI. This package is no longer maintained. New users should use the Go module at github.com/kenshin579/korea-investment-stock."
```

- [ ] **Step 3: classifiers 추가**

`[project]` 섹션 (description 바로 아래) 에 다음 키 추가 (이미 `classifiers` 가 있으면 항목만 추가):

```toml
classifiers = [
    "Development Status :: 7 - Inactive",
    "Intended Audience :: Developers",
    "License :: OSI Approved :: MIT License",
    "Programming Language :: Python :: 3",
]
```

- [ ] **Step 4: 검증**

Run:
```bash
python -c "import tomllib; data = tomllib.loads(open('pyproject.toml').read()); print(data['project']['description'][:50]); print(data['project']['classifiers'])"
```
Expected: 출력에 `[DEPRECATED` 시작 문자열과 `Development Status :: 7 - Inactive` 가 들어감

- [ ] **Step 5: Commit**

```bash
git add pyproject.toml
git commit -m "[chore] pyproject.toml deprecation notice 부착

description 변경 + classifiers 에 Inactive 추가.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 2: README.md 상단에 deprecation 경고 추가 (임시)

> 이 단계의 README 변경은 **마지막 Python release(v0.19.0) 시점의 README** 가 됨. Group F 에서 Go 용으로 완전 교체할 예정.

**Files:**
- Modify: `README.md` (최상단)

- [ ] **Step 1: README.md 상단 (첫 번째 `# 🚀 Korea Investment Stock` 헤더 바로 위) 에 다음 블록 삽입**

```markdown
> ## ⚠️ DEPRECATED — 이 라이브러리는 더 이상 유지보수되지 않습니다
>
> 이 Python 패키지는 **`v0.19.0` 을 마지막으로 신규 기능 추가가 중단**되었습니다.
> 신규 사용자는 같은 repo 의 **Go 모듈** 을 사용해주세요.
>
> - **Go 모듈**: `github.com/kenshin579/korea-investment-stock` (개발 중)
> - **기존 Python v1.x 사용자**: 이 코드는 [`python-final`](https://github.com/kenshin579/korea-investment-stock/tree/python-final) 태그로 영구 보존됩니다. critical security fix 를 제외한 신규 PR 은 받지 않습니다.
> - **마이그레이션 가이드**: 본 repo 의 [Phase 0 design spec](docs/superpowers/specs/2026-05-03-korea-investment-go-migration-design.md) 참고
>
> ---
```

- [ ] **Step 2: 인코딩 검증**

Run:
```bash
file -I README.md
```
Expected: `text/plain; charset=utf-8`

- [ ] **Step 3: Commit**

```bash
git add README.md
git commit -m "[chore] README.md 상단 deprecation 경고 추가

신규 사용자에게 Go 모듈 안내, 기존 Python 사용자에게 python-final 태그 안내.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 3: CHANGELOG.md 에 v0.19.0 release entry 추가

**Files:**
- Modify: `CHANGELOG.md`

- [ ] **Step 1: 현재 CHANGELOG 구조 확인**

Run:
```bash
head -10 CHANGELOG.md
```
Expected: 첫 줄 `# CHANGELOG`, 그 다음에 `## [Unreleased]` 섹션이 있음

- [ ] **Step 2: `## [Unreleased]` 헤더를 `## [v0.19.0] — 2026-05-03` 으로 교체하고, 그 위에 새 [Unreleased] 와 deprecation notice 추가**

`# CHANGELOG` 바로 아래에 다음 블록을 삽입 (기존 `## [Unreleased]` 헤더는 `## [v0.19.0] — 2026-05-03` 으로 변경):

```markdown
## [Unreleased]

> 본 repo 는 Go 로 마이그레이션 중입니다. Python 신규 기능 추가는 중단되었습니다.

## [v0.19.0] — 2026-05-03

### Deprecation Notice

- **이 버전이 Python 라이브러리의 마지막 기능 release 입니다.** 이후 Go 모듈로 대체됩니다.
- 마지막 Python 커밋은 `python-final` 태그로 영구 보존됩니다.
- PyPI 패키지 자체는 archive 하지 않으며, critical security fix 만 v0.19.x patch 로 받을 수 있습니다.
- 신규 사용자는 Go 모듈 (`github.com/kenshin579/korea-investment-stock`) 을 사용해주세요.

상세 내용: [Phase 0 design spec](docs/superpowers/specs/2026-05-03-korea-investment-go-migration-design.md)

### Added (이 버전에 포함된 기능 — 기존 [Unreleased] 항목 유지)

(여기 아래는 기존 [Unreleased] 섹션 내용을 그대로 유지)
```

> 주의: 기존 `[Unreleased]` 섹션의 본문(Phase 1 API 확장 등) 은 모두 `## [v0.19.0]` 의 `### Added` 아래로 이동.

- [ ] **Step 3: 인코딩 검증**

Run:
```bash
file -I CHANGELOG.md
```
Expected: `charset=utf-8`

- [ ] **Step 4: Commit**

```bash
git add CHANGELOG.md
git commit -m "[chore] CHANGELOG v0.19.0 deprecation release 항목 추가

기존 [Unreleased] 섹션을 v0.19.0 으로 확정하고 deprecation notice 부착.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 4: v0.19.0 + python-final 태그 박기

> 이 task 는 **PyPI release 와 별개로** git 태그만 박는다. PyPI 업로드는 Task 5 에서 사용자가 직접 실행.

**Files:** (git 태그만)

- [ ] **Step 1: 현재 HEAD 가 Task 1~3 의 모든 commit 을 포함하는지 확인**

Run:
```bash
git log --oneline -5
```
Expected: 최근 3개 커밋이 (1) pyproject.toml deprecation, (2) README deprecation 경고, (3) CHANGELOG v0.19.0 — 순서로 보임

- [ ] **Step 2: 현재 HEAD 에 두 개의 태그 박기**

Run:
```bash
git tag -a v0.19.0 -m "v0.19.0: Last Python release with deprecation notice. Future development moves to Go."
git tag -a python-final -m "python-final: Permanent reference tag for the last Python codebase. Same commit as v0.19.0."
```

- [ ] **Step 3: 태그 검증**

Run:
```bash
git tag -l 'v0.19.0' 'python-final' && git show v0.19.0 --no-patch --format='%H' && git show python-final --no-patch --format='%H'
```
Expected: 두 줄 다 같은 commit hash 출력

> 태그는 push 하지 않는다 — Task 11 의 PR 검토 후 사용자가 직접 push.

---

## Task 5: PyPI v0.19.0 release (사용자 직접 실행)

> **이 task 는 Claude 가 직접 실행할 수 없습니다.** PyPI 업로드 권한과 토큰이 필요하므로 사용자가 직접 실행합니다.

**작업자:** kenshin579 (사용자)

**Files:**
- Use: `upload.sh` (기존)

- [ ] **Step 1: 사용자에게 다음 절차 안내**

```bash
# 1. Python 가상환경 활성화
source .venv/bin/activate

# 2. v0.19.0 태그가 박힌 commit 으로 빌드
python -m build

# 3. dist/ 확인
ls dist/
# 예: korea_investment_stock-0.19.0-py3-none-any.whl
#     korea_investment_stock-0.19.0.tar.gz

# 4. PyPI 업로드
./upload.sh
# 또는: twine upload dist/*
```

- [ ] **Step 2: PyPI 페이지 확인 안내**

업로드 완료 후 사용자가 https://pypi.org/project/korea-investment-stock/0.19.0/ 에서 다음 확인:
- 페이지 상단에 "[DEPRECATED — use the Go version] ..." description 표시
- Classifiers 에 "Development Status :: 7 - Inactive" 표시
- README 의 deprecation 경고 박스가 보임

> 이 task 가 끝나기 전까지 Task 6 (Python 코드 제거) 으로 넘어가지 않는다 — PyPI 업로드 시점에 main 에 Python 코드가 살아있어야 빌드 가능.

---

## Task 6: main 에서 Python 코드 제거

> 사전 조건: Task 5 (PyPI v0.19.0 업로드) 가 완료되었거나, 사용자가 "PyPI 업로드는 나중에 하겠다" 고 명시한 상태.
>
> Task 4 의 `python-final` 태그가 박혀 있으므로, git history 에는 Python 코드가 영구 보존됨.

**Files:**
- Delete: `korea_investment_stock/`
- Delete: `korea_investment_stock.egg-info/`
- Delete: `examples/`
- Delete: `scripts/`
- Delete: `upload.sh`
- Delete: `pyproject.toml`
- Delete: `Makefile`
- Delete: `TEST_REPORT.md`

- [ ] **Step 1: 삭제 대상 목록 검증**

Run:
```bash
ls -d korea_investment_stock/ korea_investment_stock.egg-info/ examples/ scripts/ upload.sh pyproject.toml Makefile TEST_REPORT.md 2>/dev/null
```
Expected: 위 8개 항목이 모두 출력됨 (없는 게 있다면 plan 의 삭제 목록에서 제외)

- [ ] **Step 2: git rm 으로 일괄 삭제**

Run:
```bash
git rm -r korea_investment_stock/
git rm -r korea_investment_stock.egg-info/ 2>/dev/null || true   # gitignore 일 수 있음
git rm -r examples/
git rm -r scripts/
git rm upload.sh
git rm pyproject.toml
git rm Makefile
git rm TEST_REPORT.md 2>/dev/null || true
```

- [ ] **Step 3: 삭제 결과 검증**

Run:
```bash
ls korea_investment_stock 2>/dev/null && echo "ERROR: still exists" || echo "OK: deleted"
git status --short | head -20
```
Expected: "OK: deleted" + `git status` 에서 `D` (deleted) 마커가 다수 표시됨

- [ ] **Step 4: 남아있는 root 파일 점검**

Run:
```bash
ls -1 | grep -v '^\.' | sort
```
Expected: 다음 항목만 남아있어야 함
```
CHANGELOG.md
LICENSE
README.md
docs
```

- [ ] **Step 5: Commit**

```bash
git commit -m "[chore] Python 코드 제거 — Go 마이그레이션 시작

python-final 태그(v0.19.0 동일 commit) 에 Python 코드 영구 보존.
이후 main 은 Go-only 로 진행.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 7: README.md 를 Go 용으로 완전 교체

**Files:**
- Replace: `README.md` (기존 내용 전부 폐기)

- [ ] **Step 1: 새 README.md 내용을 작성** (아래 그대로 사용)

````markdown
# korea-investment-stock (Go)

[![Go Reference](https://pkg.go.dev/badge/github.com/kenshin579/korea-investment-stock.svg)](https://pkg.go.dev/github.com/kenshin579/korea-investment-stock)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**한국투자증권 OpenAPI Go 클라이언트** — typed struct, context-first, functional options, 자동 토큰 갱신/rate limit 내장.

> ⚠️ **Work in progress.** 이 라이브러리는 현재 초기 개발 단계입니다 (`v0.x`). 안정화 시 `v1.0.0` 으로 올릴 예정입니다.

## Python 사용자에게

이 repo 는 **2026-05-03 부로 Python → Go 로 전환**되었습니다.

- 기존 Python 코드 (`v0.18.0` 까지 + `v0.19.0` deprecation release): [`python-final`](https://github.com/kenshin579/korea-investment-stock/tree/python-final) 태그로 영구 보존
- PyPI 패키지 (`korea-investment-stock`): `v0.19.0` 까지 그대로 유지. critical security fix 외 신규 기능 없음.
- 마이그레이션 배경: [Phase 0 design spec](docs/superpowers/specs/2026-05-03-korea-investment-go-migration-design.md)

## Install

```bash
go get github.com/kenshin579/korea-investment-stock
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    kis "github.com/kenshin579/korea-investment-stock"
)

func main() {
    client, err := kis.NewClient(
        os.Getenv("KOREA_INVESTMENT_API_KEY"),
        os.Getenv("KOREA_INVESTMENT_API_SECRET"),
        os.Getenv("KOREA_INVESTMENT_ACCOUNT_NO"),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    price, err := client.Domestic.FetchPrice(ctx, "005930")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Samsung Electronics: %s KRW\n", price.Current.String())
}
```

> Phase 1 에서 실제 메서드들이 추가됩니다. 현재 main 의 `client.Domestic.FetchPrice` 는 정의만 되어 있고 구현은 Phase 1 에서 채워집니다.

## Design

- **호출 스타일**: `client.Domestic.<Method>(ctx, ...)` 1단계 그룹화 (go-github / stripe-go 패턴)
- **응답**: typed struct, 한투 API 의 한글 약어 필드는 JSON 태그로 매핑하고 영문 필드명으로 노출
- **에러**: `*kis.APIError` (rt_cd / msg_cd / msg1) + sentinel 에러 (`ErrTokenExpired`, `ErrRateLimited`, `ErrNotFound`)
- **자동 처리**: 토큰 갱신, rate limit (token bucket, 기본 15 req/sec), 429/5xx 재시도
- **HTTP**: 내부적으로 [resty](https://github.com/go-resty/resty) 사용 (사용자는 표준 `*http.Client` 만 알면 됨)
- **금융 정밀도**: 가격 필드는 [shopspring/decimal](https://github.com/shopspring/decimal)

상세 설계: [Phase 0 design spec](docs/superpowers/specs/2026-05-03-korea-investment-go-migration-design.md)

## Scope

- ✅ 국내주식 (시세, 차트, 순위, 재무, 투자자 동향, IPO/예탁원, 심볼)
- ✅ 해외주식 (시세, 차트, 순위)
- ❌ 선물옵션 — 영구 제외
- ❌ 장내채권 — 영구 제외
- ❌ 실시간 WebSocket — 추후 별도 spec
- ❌ 주식 주문/잔고/예약주문 — 본 spec 에서 다루지 않음

## License

MIT — 기존 Python 라이브러리와 동일.
````

- [ ] **Step 2: README 작성**

`README.md` 를 위 내용으로 완전 교체.

- [ ] **Step 3: 인코딩 검증**

Run:
```bash
file -I README.md
```
Expected: `charset=utf-8`

- [ ] **Step 4: Commit**

```bash
git add README.md
git commit -m "[chore] README Go 버전으로 완전 교체

Python 사용자 안내, Go install/quick start, 설계 요약, scope 정리.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 8: Go module 초기화 + .gitignore

**Files:**
- Create: `go.mod`
- Create: `.gitignore`

- [ ] **Step 1: go.mod 초기화**

Run:
```bash
go mod init github.com/kenshin579/korea-investment-stock
```

- [ ] **Step 2: Go 버전 1.23 으로 설정**

`go.mod` 내용 확인:
```bash
cat go.mod
```

Expected: 다음과 비슷
```go
module github.com/kenshin579/korea-investment-stock

go 1.23
```

`go 1.23` 이 아니라면 `go.mod` 의 첫 번째 두 줄을 위와 동일하게 편집.

- [ ] **Step 3: .gitignore 작성**

`.gitignore` 파일을 다음 내용으로 작성:

```
# Go
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out
vendor/
go.work
go.work.sum

# IDE
.idea/
.vscode/
*.iml

# OS
.DS_Store
Thumbs.db

# Logs
*.log

# Coverage
coverage.txt
coverage.html
```

- [ ] **Step 4: 검증**

Run:
```bash
go env GOMOD && go mod verify
```
Expected: 첫 번째 라인은 `go.mod` 의 절대경로, 두 번째 라인은 `all modules verified` 또는 의존성이 없으니 `verifying ...` 없이 그냥 종료

- [ ] **Step 5: Commit**

```bash
git add go.mod .gitignore
git commit -m "[chore] go.mod 초기화 + .gitignore 추가

Go 1.23, module path = github.com/kenshin579/korea-investment-stock.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 9: Go 패키지 스켈레톤 — `client.go` 와 sub-package doc 파일

> 빈 패키지가 아니라 **컴파일 가능한 최소 골격**. `kis.NewClient(...)` 시그니처 + sub-package 의 `doc.go` 만.
>
> 실제 메서드 구현 (`FetchPrice` 등) 은 Phase 1 의 영역. 이 task 는 **스켈레톤이 컴파일 되는지** 까지만 책임진다.

**Files:**
- Create: `client.go`
- Create: `domestic/doc.go`
- Create: `overseas/doc.go`
- Create: `internal/httpclient/doc.go`
- Create: `internal/ratelimit/doc.go`
- Create: `internal/token/doc.go`

- [ ] **Step 1: 디렉터리 생성**

Run:
```bash
mkdir -p domestic overseas internal/httpclient internal/ratelimit internal/token
```

- [ ] **Step 2: `client.go` 작성**

```go
// Package kis is a Go client for the Korea Investment Securities OpenAPI.
//
// This is a Phase 0 skeleton. Methods are added in Phase 1.
//
// See docs/superpowers/specs/2026-05-03-korea-investment-go-migration-design.md
// for the design rationale.
package kis

import (
	"errors"
	"net/http"
)

// Client is the entry point of the kis library. Domain-specific operations
// are grouped under sub-clients (Client.Domestic, Client.Overseas).
//
// Construct a Client with NewClient, passing functional options to override
// defaults.
type Client struct {
	apiKey    string
	apiSecret string
	accountNo string

	// Sub-clients. Initialized in NewClient.
	Domestic *DomesticClient
	Overseas *OverseasClient
}

// DomesticClient groups methods for the domestic Korean stock market.
// Implementations are added in Phase 1.
type DomesticClient struct {
	parent *Client
}

// OverseasClient groups methods for overseas (US/HK/JP/CN/VN) markets.
// Implementations are added in Phase 1.
type OverseasClient struct {
	parent *Client
}

// Option configures a Client. See WithBaseURL, WithRetries, etc.
type Option func(*clientOptions)

type clientOptions struct {
	baseURL    string
	retries    int
	rateLimit  int
	httpClient *http.Client
}

// NewClient constructs a kis Client.
//
// apiKey, apiSecret, accountNo are required. Phase 1 will add functional
// options (WithBaseURL, WithRetries, WithRateLimit, WithHTTPClient,
// WithTokenStorage, WithLogger).
func NewClient(apiKey, apiSecret, accountNo string, opts ...Option) (*Client, error) {
	if apiKey == "" || apiSecret == "" || accountNo == "" {
		return nil, errors.New("kis: apiKey, apiSecret, accountNo are all required")
	}
	c := &Client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		accountNo: accountNo,
	}
	c.Domestic = &DomesticClient{parent: c}
	c.Overseas = &OverseasClient{parent: c}
	return c, nil
}
```

- [ ] **Step 3: sub-package `doc.go` 파일 작성**

`domestic/doc.go`:
```go
// Package domestic contains response types for domestic (Korean) stock APIs.
//
// Phase 0 skeleton — implementations are added in Phase 1.
package domestic
```

`overseas/doc.go`:
```go
// Package overseas contains response types for overseas (US/HK/JP/CN/VN) stock APIs.
//
// Phase 0 skeleton — implementations are added in Phase 1.
package overseas
```

`internal/httpclient/doc.go`:
```go
// Package httpclient wraps resty for transparent token refresh, rate limiting,
// and 429/5xx retries. Not exposed to library users.
//
// Phase 0 skeleton.
package httpclient
```

`internal/ratelimit/doc.go`:
```go
// Package ratelimit implements a token-bucket rate limiter shared across
// goroutines. Default 15 req/sec.
//
// Phase 0 skeleton.
package ratelimit
```

`internal/token/doc.go`:
```go
// Package token manages the KIS access token: issuance, automatic refresh
// (5 minutes before expiry), and pluggable storage backends (file/redis).
//
// Phase 0 skeleton.
package token
```

- [ ] **Step 4: 컴파일 검증**

Run:
```bash
go build ./...
```
Expected: 출력 없음 (성공). 에러가 나면 Step 2 의 `client.go` 또는 doc.go 의 패키지 선언 점검.

- [ ] **Step 5: vet 검증**

Run:
```bash
go vet ./...
```
Expected: 출력 없음

- [ ] **Step 6: 테스트 (이 task 는 테스트 없음 확인)**

Run:
```bash
go test ./... 2>&1 | head -10
```
Expected: `?       github.com/kenshin579/korea-investment-stock        [no test files]` 같은 메시지가 6개 (각 패키지마다)

- [ ] **Step 7: Commit**

```bash
git add client.go domestic overseas internal
git commit -m "[chore] Go 패키지 스켈레톤 — client.go + sub-package doc

NewClient 시그니처와 Domestic/Overseas sub-client 정의 (구현은 Phase 1).
internal/{httpclient,ratelimit,token} 디렉터리 placeholder.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 10: 최종 점검

**Files:** (점검만, 파일 변경 없음)

- [ ] **Step 1: 디렉터리 구조 검증**

Run:
```bash
ls -1 | grep -v '^\.' | sort
```
Expected:
```
CHANGELOG.md
LICENSE
README.md
client.go
docs
domestic
go.mod
internal
overseas
```

- [ ] **Step 2: Go 컴파일 + vet 재검증**

Run:
```bash
go build ./... && go vet ./... && echo "OK: build + vet"
```
Expected: `OK: build + vet`

- [ ] **Step 3: 한글 콘텐츠 인코딩 검증**

Run:
```bash
file -I README.md CHANGELOG.md docs/superpowers/specs/2026-05-03-korea-investment-go-migration-design.md docs/superpowers/specs/2026-05-03-phase0-step1-implementation-plan.md
```
Expected: 모두 `charset=utf-8`

- [ ] **Step 4: 태그 검증**

Run:
```bash
git tag -l 'v0.19.0' 'python-final' && git show v0.19.0 --no-patch --format='%H' && git show python-final --no-patch --format='%H'
```
Expected: 두 태그 모두 존재, 같은 commit hash

- [ ] **Step 5: 브랜치 commit history 검증**

Run:
```bash
git log main..HEAD --oneline
```
Expected: 다음 8개 commit 이 시간순으로 보임 (가장 위가 최신):
1. Go 패키지 스켈레톤
2. go.mod 초기화 + .gitignore
3. README Go 버전 교체
4. Python 코드 제거
5. CHANGELOG v0.19.0
6. README deprecation 경고 (※ 2번 task — Python 시점)
7. pyproject.toml deprecation
8. Phase 0 Go 마이그레이션 설계 문서 (※ spec 커밋 — 이미 ebfa2e7)

(7→1 순서. 4와 5 사이의 순서는 plan 의 task 순서 기준)

> commit 6, 7 은 main 에서는 제거된 파일을 수정한 commit 이지만, history 에는 남아 PyPI v0.19.0 빌드 base 가 됨.

---

## Task 11: PR 생성 및 push (사용자 승인 후)

> Claude 는 push / PR 생성을 사용자 명시적 승인 후에만 실행한다 (글로벌 정책).

**Files:** (git 원격 작업)

- [ ] **Step 1: 사용자에게 승인 요청**

Claude 가 사용자에게 다음 중 하나의 응답을 요청:
- "지금 push + PR 생성하라"
- "PyPI 업로드 (Task 5) 먼저 하고 그 다음에 push 하겠다"
- "수정 사항 있다"

- [ ] **Step 2: 승인 받으면 push + 태그 push**

Run:
```bash
git push -u origin docs/golang-migration-spec
git push origin v0.19.0 python-final
```

- [ ] **Step 3: PR 생성 (gh CLI + HEREDOC)**

Run:
```bash
gh pr create --title "Phase 0: Python → Go 마이그레이션 (spec + Step 1 정리)" --base main --body "$(cat <<'EOF'
## Summary

`korea-investment-stock` repo 를 Python → Go 로 전환하기 위한 Phase 0 작업.

### 포함 내용

- **Phase 0 design spec**: \`docs/superpowers/specs/2026-05-03-korea-investment-go-migration-design.md\` (결정 기록)
- **Phase 0 Step 1 implementation plan**: \`docs/superpowers/specs/2026-05-03-phase0-step1-implementation-plan.md\`
- **PyPI deprecation 절차**: \`v0.19.0\` 마지막 release (description / classifiers / README / CHANGELOG 에 deprecation notice)
- **태그**: \`v0.19.0\` (마지막 Python release), \`python-final\` (영구 보존 별칭)
- **main 정리**: Python 코드 전부 제거 (git history 보존)
- **Go 스켈레톤**: \`go.mod\` (Go 1.23), \`client.go\` (Client/NewClient 골격), \`domestic\`, \`overseas\`, \`internal/{httpclient,ratelimit,token}\` 패키지

### 외부 사용자 영향

- PyPI \`korea-investment-stock\` 패키지: \`v0.19.0\` 까지 그대로 유지, 새 버전부터 Inactive 상태
- \`python-final\` 태그로 옛 코드 영구 참조 가능
- 신규 기능 추가 없음, critical security fix 만 v0.19.x patch 허용

### 다음 단계

- Phase 1 spec: 어떤 한투 API 를 어떤 우선순위로 Go 로 옮길지 (별도 PR)

## Test plan

- [ ] \`go build ./...\` 성공
- [ ] \`go vet ./...\` 출력 없음
- [ ] 모든 한글 .md 파일 \`charset=utf-8\`
- [ ] \`v0.19.0\` 와 \`python-final\` 태그가 같은 commit
- [ ] (Task 5) PyPI \`v0.19.0\` release 페이지에 deprecation notice 노출 확인
EOF
)"
```

- [ ] **Step 4: PR URL 사용자에게 전달**

Run:
```bash
gh pr view --json url --jq '.url'
```

PR URL 을 사용자에게 보여주면서 review 안내.

---

## Self-Review 결과

**1. Spec coverage:**

| Spec 요구사항 | 구현 task |
|---|---|
| §1: Python 코드 즉시 제거 | Task 6 |
| §1: v1.x 태그 보존 | Task 4 (`python-final` + `v0.19.0`) |
| §1: PyPI deprecation notice | Task 1 (description), Task 2 (README), Task 3 (CHANGELOG), Task 5 (release) |
| §1: 신규 추가 없음 정책 | Task 1 의 `Development Status :: 7 - Inactive`, Task 3 의 CHANGELOG entry |
| §2 Step 1: Python freeze | Task 1~3 으로 deprecation 메타데이터 박고 Task 4 로 태그 |
| §2 Step 1: README 마이그레이션 안내 | Task 2 (Python 시점 경고), Task 7 (Go 시점 풀 README) |
| §2 Step 1: main 에서 Python 제거 | Task 6 |
| §2 Step 1: go.mod 초기화 + 빈 패키지 스켈레톤 | Task 8 (go.mod), Task 9 (client.go + 5개 sub-package) |
| §3: 패키지 이름 `kis`, module path `github.com/kenshin579/korea-investment-stock` | Task 8, Task 9 의 `client.go` 첫 줄 |
| §3: 디렉터리 구조 (domestic/, overseas/, internal/{httpclient,ratelimit,token}) | Task 9 의 mkdir |
| §5: 인코딩 검증 (한글 UTF-8) | Task 2 Step 2, Task 3 Step 3, Task 7 Step 3, Task 10 Step 3 |

> §3 의 의존성 (`resty`, `decimal`) 은 본 plan 에서는 추가하지 않음. Phase 1 에서 첫 메서드 구현할 때 함께 import. Phase 0 의 스켈레톤은 표준 라이브러리만 사용해 컴파일 가능해야 한다.

**2. Placeholder scan:** 모든 step 에 실제 명령어 / 코드 / 본문 포함. "TBD", "TODO", "implement later" 사용 안 함.

**3. Type consistency:** Task 9 의 `Client.Domestic`, `Client.Overseas` 가 spec §3 의 호출 스타일 (`client.Domestic.FetchPrice(ctx, ...)`) 과 일치.

**4. Risk:** Task 5 (PyPI 업로드) 가 사용자 직접 작업이라 plan 이 그 단계에서 멈춤. Task 6~11 은 Task 5 완료 (또는 명시적 보류) 후에 진행.

---

## Execution Handoff

Plan 작성 완료. 두 가지 실행 옵션:

1. **Subagent-Driven (recommended)** — task 별로 fresh subagent 가 실행 + task 사이 review. 빠른 iteration.
2. **Inline Execution** — 본 세션에서 task 들 batch 실행 + checkpoint 마다 review.

어느 쪽으로 진행할까요?
