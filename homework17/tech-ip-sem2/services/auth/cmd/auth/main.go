package main

import (
	"log/slog"
	"net/http"
	"os"

	authhttp "tech-ip-sem2/services/auth/internal/http"
	"tech-ip-sem2/services/auth/internal/service"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	port := os.Getenv("AUTH_PORT")
	if port == "" {
		port = "8081"
	}

	svc := service.New()
	handler := authhttp.NewHandler(svc)
	router := authhttp.NewRouter(handler, logger)

	logger.Info("auth service starting", "port", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		logger.Error("server error", "err", err)
	}
}
