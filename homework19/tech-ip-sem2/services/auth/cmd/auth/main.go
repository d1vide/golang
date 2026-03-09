package main

import (
	"net/http"
	"os"

	authhttp "tech-ip-sem2/services/auth/internal/http"
	"tech-ip-sem2/services/auth/internal/service"
	"tech-ip-sem2/shared/logger"
)

func main() {
	log := logger.New()
	defer log.Sync()

	port := os.Getenv("AUTH_PORT")
	if port == "" {
		port = "8081"
	}

	svc := service.New()
	handler := authhttp.NewHandler(svc)
	router := authhttp.NewRouter(handler, log)

	log.Info("auth service starting", logger.Port(port))
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Error("server error", logger.Err(err))
	}
}
