package http

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
	logger Logger
	app    Application
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	DebugKV(msg string, keysAndValues ...interface{})
	InfoKV(msg string, keysAndValues ...interface{})
	WarnKV(msg string, keysAndValues ...interface{})
	ErrorKV(msg string, keysAndValues ...interface{})
}

type Application interface {
	GetResizedImage(w http.ResponseWriter, r *http.Request)
	ClearCache(w http.ResponseWriter, r *http.Request)
}

func NewServer(logger Logger, app Application, port int) *Server {
	srv := &Server{
		logger: logger,
		app:    app,
	}

	srv.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           srv.setupRoutes(),
		ReadHeaderTimeout: time.Second * 10,
	}

	return srv
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error(fmt.Sprintf("Server failed: %s", err))
		}
	}()

	<-ctx.Done()

	return s.Stop(ctx)
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	return nil
}

func (s *Server) setupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", s.loggingMiddleware(http.HandlerFunc(s.app.GetResizedImage)))
	mux.Handle("/clear", s.loggingMiddleware(http.HandlerFunc(s.app.ClearCache)))
	return mux
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body = b
	return rw.ResponseWriter.Write(b)
}
