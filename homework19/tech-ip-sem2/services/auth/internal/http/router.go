package authhttp

import (
	"net/http"

	"go.uber.org/zap"

	"tech-ip-sem2/shared/middleware"
)

func NewRouter(h *Handler, logger *zap.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/auth/login", h.Login)
	mux.HandleFunc("/v1/auth/verify", h.Verify)

	var handler http.Handler = mux
	handler = middleware.AccessLog(logger)(handler)
	handler = middleware.RequestID(handler)
	return handler
}
