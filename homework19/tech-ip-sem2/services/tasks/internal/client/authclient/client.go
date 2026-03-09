package authclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"tech-ip-sem2/shared/httpx"
	"tech-ip-sem2/shared/middleware"
)

const defaultTimeout = 3 * time.Second

type VerifyResult struct {
	Valid   bool
	Subject string
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: httpx.NewClient(defaultTimeout),
	}
}

var (
	ErrUnauthorized = fmt.Errorf("unauthorized")
	ErrUpstream     = fmt.Errorf("auth service unavailable")
)

func (c *Client) Verify(ctx context.Context, authHeader string) (VerifyResult, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	url := c.baseURL + "/v1/auth/verify"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return VerifyResult{}, ErrUpstream
	}

	req.Header.Set("Authorization", authHeader)

	if rid := middleware.GetRequestID(ctx); rid != "" {
		req.Header.Set("X-Request-ID", rid)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return VerifyResult{}, ErrUpstream
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode == http.StatusOK:
		var body struct {
			Valid   bool   `json:"valid"`
			Subject string `json:"subject"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			return VerifyResult{}, ErrUpstream
		}
		return VerifyResult{Valid: body.Valid, Subject: body.Subject}, nil

	case resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden:
		return VerifyResult{}, ErrUnauthorized

	default:
		return VerifyResult{}, ErrUpstream
	}
}
