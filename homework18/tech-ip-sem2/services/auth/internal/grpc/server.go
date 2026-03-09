package authgrpcserver

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authpb "tech-ip-sem2/gen/auth"
	"tech-ip-sem2/services/auth/internal/service"
)

type Server struct {
	authpb.UnimplementedAuthServiceServer
	svc *service.AuthService
}

func New(svc *service.AuthService) *Server {
	return &Server{svc: svc}
}

func (s *Server) Verify(ctx context.Context, req *authpb.VerifyRequest) (*authpb.VerifyResponse, error) {
	if req.Token == "" {
		return nil, status.Error(codes.Unauthenticated, "token is required")
	}

	subject, valid := s.svc.Verify("Bearer " + strings.TrimPrefix(req.Token, "Bearer "))
	if !valid {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	return &authpb.VerifyResponse{Valid: true, Subject: subject}, nil
}
