package taskshttp

import (
	"log/slog"
	"net/http"
	"tech-ip-sem2/shared/middleware"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(h *Handler, logger *slog.Logger) http.Handler {
	m := NewMetrics()
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/v1/tasks", h.Tasks)
	mux.HandleFunc("/v1/tasks/", h.Task)

	var handler http.Handler = mux
	handler = MetricsMiddleware(m)(handler)
	handler = middleware.Logging(logger)(handler)
	handler = middleware.RequestID(handler)
	return handler

}
