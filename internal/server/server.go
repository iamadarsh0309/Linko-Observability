package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"boot.dev/linko/internal/auth"
	"boot.dev/linko/internal/store"
)

type Server struct {
	httpServer *http.Server
	store      *store.Store
	cancel     context.CancelFunc
	logger     *slog.Logger
}

func New(st *store.Store, port int, cancel context.CancelFunc, logger *slog.Logger) *Server {
	mux := http.NewServeMux()

	s := &Server{
		store:  st,
		cancel: cancel,
		logger: logger,
	}
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: requestLogger(logger)(mux),
	}

	mux.HandleFunc("GET /", s.handlerIndex)
	mux.Handle("POST /api/login", auth.Middleware(logger, http.HandlerFunc(s.handlerLogin)))
	mux.Handle("POST /api/shorten", auth.Middleware(logger, http.HandlerFunc(s.handlerShortenLink)))
	mux.Handle("GET /api/stats", auth.Middleware(logger, http.HandlerFunc(s.handlerStats)))
	mux.Handle("GET /api/urls", auth.Middleware(logger, http.HandlerFunc(s.handlerListURLs)))
	mux.HandleFunc("GET /{shortCode}", s.handlerRedirect)
	mux.HandleFunc("POST /admin/shutdown", s.handlerShutdown)

	return s
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return err
	}
	s.logger.Info("Linko is running on http://localhost:%d\n", ln.Addr().(*net.TCPAddr).Port)
	if err := s.httpServer.Serve(ln); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) handlerShutdown(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("ENV") == "production" {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	go s.cancel()
}
