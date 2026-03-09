package main

import (
	"log/slog"
	"net/http"
	"os"

	"tech-ip-sem2/services/tasks/internal/client/authgrpc"
	taskshttp "tech-ip-sem2/services/tasks/internal/http"
	"tech-ip-sem2/services/tasks/internal/service"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	port := os.Getenv("TASKS_PORT")
	if port == "" {
		port = "8082"
	}

	grpcAddr := os.Getenv("AUTH_GRPC_ADDR")
	if grpcAddr == "" {
		grpcAddr = "localhost:50051"
	}

	authClient, err := authgrpc.New(grpcAddr)
	if err != nil {
		logger.Error("failed to connect to auth gRPC", "err", err)
		os.Exit(1)
	}

	svc := service.New()
	handler := taskshttp.NewHandler(svc, authClient)
	router := taskshttp.NewRouter(handler, logger)

	logger.Info("tasks service starting", "port", port, "auth_grpc", grpcAddr)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		logger.Error("server error", "err", err)
	}
}
