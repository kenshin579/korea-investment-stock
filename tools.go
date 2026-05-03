//go:build tools

// tools.go pins build-time and test-time dependencies so they appear in go.mod
// and are not removed by go mod tidy.
package kis

import (
	_ "github.com/alicebob/miniredis/v2"
	_ "github.com/go-resty/resty/v2"
	_ "github.com/jarcoal/httpmock"
	_ "github.com/redis/go-redis/v9"
	_ "github.com/shopspring/decimal"
	_ "github.com/stretchr/testify/assert"
	_ "golang.org/x/sync/singleflight"
	_ "gopkg.in/yaml.v3"
)
