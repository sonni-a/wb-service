package metrics

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func normalizePath(path string) string {
	if strings.HasPrefix(path, "/order/") {
		return "/order/{id}"
	}
	return path
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		rec := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rec, r)

		duration := time.Since(start).Seconds()

		path := normalizePath(r.URL.Path)

		HttpRequestsTotal.
			WithLabelValues(r.Method, path, strconv.Itoa(rec.statusCode)).
			Inc()

		HttpRequestDuration.
			WithLabelValues(r.Method, path).
			Observe(duration)
	})
}
