# CLAUDE.md

Go client library for the Korea Investment Securities OpenAPI.

> **Phase 1.3 — domestic 순위/재무 (v1.1.0).** Phase 1.4+ 메서드는 추후 sub-plan 으로.

- Design spec: [`docs/superpowers/specs/2026-05-03-korea-investment-go-migration-design.md`](docs/superpowers/specs/2026-05-03-korea-investment-go-migration-design.md)
- Phase 0 implementation plan: [`docs/superpowers/specs/2026-05-03-phase0-step1-implementation-plan.md`](docs/superpowers/specs/2026-05-03-phase0-step1-implementation-plan.md)
- Phase 1.2 implementation plan: [`docs/superpowers/specs/2026-05-03-phase1-2-domestic-quotes-implementation-plan.md`](docs/superpowers/specs/2026-05-03-phase1-2-domestic-quotes-implementation-plan.md)
- Phase 1.3 implementation plan: [`docs/superpowers/specs/2026-05-04-phase1-3-ranking-financial-implementation-plan.md`](docs/superpowers/specs/2026-05-04-phase1-3-ranking-financial-implementation-plan.md)
- Legacy Python: see `python-final` tag (commit `e3fc52f`); PyPI `korea-investment-stock` v0.19.0 deprecated.

## Stack

Go 1.25+ (`golang.org/x/sync` v0.20+ 이 1.25 요구) · module `github.com/kenshin579/korea-investment-stock` · package `kis` · `go-resty/resty/v2` (HTTP, internal) · `shopspring/decimal` (가격) · `stretchr/testify` (test)

## Layout

```
client.go           # Client + NewClient + sub-clients
domestic/           # 국내주식 API (Phase 1)
overseas/           # 해외주식 API (Phase 1)
internal/
  httpclient/       # resty wrapper, token refresh, retries
  ratelimit/        # token bucket
  token/            # token storage (file/redis)
docs/api/           # KIS API specs (Korean)
docs/superpowers/   # design specs / implementation plans
```

## Common Commands

```bash
go build ./...
go vet ./...
go test ./...
go mod tidy
```

## Conventions

- Call style: `client.Domestic.InquirePrice(ctx, ...)` (1-level service grouping)
- Errors: `error` 반환. `error.Error()` 에 `msg_cd`/`msg1` 포함. typed error 는 추후 재도입 예정.
- Korean comments preferred for domain-specific code

## Out of Scope (Phase 0)

선물옵션 · 장내채권 · 실시간 WebSocket · 주문/잔고/예약주문

## Git Policy

- Never commit directly to `main`. Always feature branch.
- Tags: `v0.x.y` for Go releases. `python-final` for legacy Python reference.
