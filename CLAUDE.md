# CLAUDE.md

Go client library for the Korea Investment Securities OpenAPI.

> **Phase 0 — skeleton.** API methods are added in Phase 1.

- Design spec: [`docs/superpowers/specs/2026-05-03-korea-investment-go-migration-design.md`](docs/superpowers/specs/2026-05-03-korea-investment-go-migration-design.md)
- Phase 0 implementation plan: [`docs/superpowers/specs/2026-05-03-phase0-step1-implementation-plan.md`](docs/superpowers/specs/2026-05-03-phase0-step1-implementation-plan.md)
- Legacy Python: see `python-final` tag (commit `e3fc52f`); PyPI `korea-investment-stock` v0.19.0 deprecated.

## Stack

Go 1.23+ · module `github.com/kenshin579/korea-investment-stock` · package `kis` · `go-resty/resty/v2` (HTTP, internal) · `shopspring/decimal` (가격) · `stretchr/testify` (test)

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

- Call style: `client.Domestic.FetchPrice(ctx, ...)` (1-level service grouping)
- Errors: `*kis.APIError` + sentinel (`ErrTokenExpired`, `ErrRateLimited`, `ErrNotFound`)
- Korean comments preferred for domain-specific code

## Out of Scope (Phase 0)

선물옵션 · 장내채권 · 실시간 WebSocket · 주문/잔고/예약주문

## Git Policy

- Never commit directly to `main`. Always feature branch.
- Tags: `v0.x.y` for Go releases. `python-final` for legacy Python reference.
