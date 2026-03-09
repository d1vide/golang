package authhttp

import (
	"log/slog"
	"net/http"

	"tech-ip-sem2/shared/middleware"
)

func NewRouter(h *Handler, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/auth/login", h.Login)
	mux.HandleFunc("/v1/auth/verify", h.Verify)

	var handler http.Handler = mux
	handler = middleware.Logging(logger)(handler)
	handler = middleware.RequestID(handler)
	return handler
}
