package authgrpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	authpb "tech-ip-sem2/gen/auth"
	"tech-ip-sem2/shared/middleware"
)

const defaultDeadline = 2 * time.Second

var (
	ErrUnauthorized = fmt.Errorf("unauthorized")
	ErrUpstream     = fmt.Errorf("auth service unavailable")
)

type Client struct {
	stub authpb.AuthServiceClient
}

func New(addr string) (*Client, error) {
	//nolint:staticcheck
	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("authgrpc: dial %s: %w", addr, err)
	}
	return &Client{stub: authpb.NewAuthServiceClient(conn)}, nil
}

func (c *Client) Verify(ctx context.Context, authHeader string) error {
	token := extractToken(authHeader)
	if token == "" {
		return ErrUnauthorized
	}

	_ = middleware.GetRequestID(ctx)

	ctx, cancel := context.WithTimeout(ctx, defaultDeadline)
	defer cancel()

	_, err := c.stub.Verify(ctx, &authpb.VerifyRequest{Token: token})
	if err == nil {
		return nil
	}

	switch status.Code(err) {
	case codes.Unauthenticated, codes.PermissionDenied:
		return ErrUnauthorized
	default:
		return ErrUpstream
	}
}

func extractToken(authHeader string) string {
	const prefix = "Bearer "
	if len(authHeader) > len(prefix) && authHeader[:len(prefix)] == prefix {
		return authHeader[len(prefix):]
	}
	return ""
}
