package taskshttp

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type metricsWriter struct {
	http.ResponseWriter
	statusCode int
}

func newMetricsWriter(w http.ResponseWriter) *metricsWriter {
	return &metricsWriter{w, http.StatusOK}
}

func (mw *metricsWriter) WriteHeader(code int) {
	mw.statusCode = code
	mw.ResponseWriter.WriteHeader(code)
}

func MetricsMiddleware(m *Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			route := normalizeRoute(r.URL.Path)

			m.InFlightRequests.Inc()
			defer m.InFlightRequests.Dec()

			start := time.Now()
			mw := newMetricsWriter(w)
			next.ServeHTTP(mw, r)
			elapsed := time.Since(start).Seconds()

			status := fmt.Sprintf("%d", mw.statusCode)
			m.RequestsTotal.WithLabelValues(r.Method, route, status).Inc()
			m.RequestDuration.WithLabelValues(r.Method, route).Observe(elapsed)
		})
	}
}

func normalizeRoute(path string) string {
	if path == "/v1/tasks" {
		return "/v1/tasks"
	}
	if strings.HasPrefix(path, "/v1/tasks/") {
		return "/v1/tasks/:id"
	}
	return "other"
}
