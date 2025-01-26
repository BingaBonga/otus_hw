package internalhttp

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func (s *Server) loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

func (s *Server) contentTypeJSONMiddleware(fn func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) { //nolint: lll
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fn(w, r)
	}
}
