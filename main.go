package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"boot.dev/linko/internal/store"
)

func requestLogger(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			logger.Printf("Served request: %s %s", r.Method, r.URL.Path)
		})
	}
}

func initializeLogger(logFile string) (*log.Logger, error) {
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		multiwriter := io.MultiWriter(os.Stderr, file)
		return log.New(multiwriter, "", log.LstdFlags), nil
	}
	return log.New(os.Stderr, "", log.LstdFlags), nil
}

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	httpPort := flag.Int("port", 8899, "port to listen on")
	dataDir := flag.String("data", "./data", "directory to store data")
	flag.Parse()

	status := run(ctx, cancel, *httpPort, *dataDir)
	cancel()
	os.Exit(status)
}

func run(ctx context.Context, cancel context.CancelFunc, httpPort int, dataDir string) int {
	logger, err := initializeLogger(os.Getenv("LINKO_LOG_FILE"))
	if err != nil {
		log.Printf("failed to open access log: %v\n", err)
		return 1
	}

	st, err := store.New(dataDir, logger)
	if err != nil {
		logger.Printf("failed to create store: %v\n", err)
		return 1
	}
	s := newServer(*st, httpPort, cancel, logger)

	var serverErr error
	go func() {
		serverErr = s.start()
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Println("Linko is shutting down")

	if err := s.shutdown(shutdownCtx); err != nil {
		logger.Printf("failed to shutdown server: %v\n", err)
		return 1
	}
	if serverErr != nil {
		logger.Printf("server error: %v\n", serverErr)
		return 1
	}
	return 0
}
