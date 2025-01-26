package internalhttp

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func (s *Server) loggingMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			next.ServeHTTP(w, r)
			s.logger.Info("Rest Request INFO",
				zap.String("IP", r.RemoteAddr),
				zap.Time("Time", startTime),
				zap.String("Version", r.Proto),
				zap.String("Path", r.URL.Path),
				zap.String("Method", r.Method),
				zap.Duration("Duration", time.Since(startTime)),
				zap.String("UserAgent", r.UserAgent()))
		})
	}
}

func (s *Server) contentTypeJSONMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			next.ServeHTTP(w, r)
		})
	}
}
