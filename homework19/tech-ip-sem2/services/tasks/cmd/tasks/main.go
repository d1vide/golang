package main

import (
	"net/http"
	"os"

	"go.uber.org/zap"

	"tech-ip-sem2/services/tasks/internal/client/authclient"
	taskshttp "tech-ip-sem2/services/tasks/internal/http"
	"tech-ip-sem2/services/tasks/internal/service"
	"tech-ip-sem2/shared/logger"
)

func main() {
	log := logger.New()
	defer log.Sync()

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
	router := taskshttp.NewRouter(handler, log)

	log.Info("tasks service starting",
		logger.Port(port),
		zap.String("auth_url", authBaseURL),
	)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Error("server error", logger.Err(err))
	}
}
