package main

import (
	"log/slog"
	"net/http"
	"os"

	"tech-ip-sem2/services/tasks/internal/client/authclient"
	taskshttp "tech-ip-sem2/services/tasks/internal/http"
	"tech-ip-sem2/services/tasks/internal/service"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	port := os.Getenv("TASKS_PORT")
	if port == "" {
		port = "8082"
	}

	authBaseURL := os.Getenv("AUTH_BASE_URL")
	if authBaseURL == "" {
		authBaseURL = "http://localhost:8081"
	}

	svc := service.New()
	authClient := authclient.New(authBaseURL)
	handler := taskshttp.NewHandler(svc, authClient)
	router := taskshttp.NewRouter(handler, logger)

	logger.Info("tasks service starting", "port", port, "auth_url", authBaseURL)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		logger.Error("server error", "err", err)
	}
}
