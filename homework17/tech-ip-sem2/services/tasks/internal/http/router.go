package taskshttp

import (
	"log/slog"
	"net/http"

	"tech-ip-sem2/shared/middleware"
)

func NewRouter(h *Handler, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/tasks", h.Tasks)

	mux.HandleFunc("/v1/tasks/", h.Task)

	var handler http.Handler = mux
	handler = middleware.Logging(logger)(handler)
	handler = middleware.RequestID(handler)
	return handler
}
