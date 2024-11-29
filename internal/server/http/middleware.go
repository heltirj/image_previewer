package http

import (
	"net/http"
	"time"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ip := r.RemoteAddr
		userAgent := r.UserAgent()

		rw := &responseWriter{ResponseWriter: w}
		next.ServeHTTP(rw, r)

		w.WriteHeader(rw.statusCode)

		latency := time.Since(start)

		if rw.statusCode != http.StatusOK {
			s.logger.ErrorKV("error to handle request", "ip", ip, "time", time.Now(), "method", r.Method, "path",
				r.URL.Path, "httpVersion", r.Proto, "statusCode", rw.statusCode, "headers", rw.Header(), "latency",
				latency,
				"userAgent",
				userAgent)
			return
		}

		s.logger.InfoKV("handled request", "ip", ip, "time", time.Now(), "method", r.Method, "path", r.URL.Path,
			"httpVersion", r.Proto, "statusCode", rw.statusCode, "latency", latency, "userAgent", userAgent)
	})
}
