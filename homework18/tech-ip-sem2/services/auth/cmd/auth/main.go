package main

import (
	"log/slog"
	"net"
	"os"

	"google.golang.org/grpc"

	authpb "tech-ip-sem2/gen/auth"
	authgrpcserver "tech-ip-sem2/services/auth/internal/grpc"
	"tech-ip-sem2/services/auth/internal/service"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	grpcPort := os.Getenv("AUTH_GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	svc := service.New()

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		logger.Error("gRPC listen failed", "err", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(grpcServer, authgrpcserver.New(svc))

	logger.Info("auth gRPC server starting", "port", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("gRPC server error", "err", err)
	}
}
