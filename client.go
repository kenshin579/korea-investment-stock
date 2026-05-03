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

// KIS OpenAPI base URLs.
const (
	// RealEnv is the production endpoint.
	RealEnv = "https://openapi.koreainvestment.com:9443"
	// PaperEnv is the paper-trading (mock) endpoint.
	PaperEnv = "https://openapivts.koreainvestment.com:29443"
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

	opts clientOptions

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

// Option configures a Client. See WithBaseURL, WithRetries, etc. (added in Phase 1).
type Option func(*clientOptions)

type clientOptions struct {
	baseURL    string
	retries    int
	rateLimit  int
	httpClient *http.Client
}

// NewClient constructs a kis Client.
//
// apiKey, apiSecret, accountNo are required and must not be empty. Phase 1 will
// add functional options (WithBaseURL, WithRetries, WithRateLimit, WithHTTPClient,
// WithTokenStorage, WithLogger).
func NewClient(apiKey, apiSecret, accountNo string, opts ...Option) (*Client, error) {
	if apiKey == "" || apiSecret == "" || accountNo == "" {
		return nil, errors.New("kis: apiKey, apiSecret, and accountNo are required and must not be empty")
	}
	cfg := clientOptions{
		baseURL: RealEnv,
		// retries / rateLimit defaults are set in Phase 1 alongside the
		// httpclient / ratelimit packages.
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	c := &Client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		accountNo: accountNo,
		opts:      cfg,
	}
	c.Domestic = &DomesticClient{parent: c}
	c.Overseas = &OverseasClient{parent: c}
	return c, nil
}
