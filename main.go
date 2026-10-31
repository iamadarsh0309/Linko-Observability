package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"boot.dev/linko/internal/logging"
	"boot.dev/linko/internal/server"
	"boot.dev/linko/internal/store"
)

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
	logger, closeLogger, err := logging.New(os.Getenv("LINKO_LOG_FILE"))
	if err != nil {
		log.Printf("failed to open access log: %v\n", err)
		return 1
	}
	defer func() {
		if err := closeLogger(); err != nil {
			log.Printf("failed to close logger: %v\n", err)
		}
	}()

	st, err := store.New(dataDir, logger)
	if err != nil {
		logger.Printf("failed to create store: %v\n", err)
		return 1
	}

	s := server.New(st, httpPort, cancel, logger)

	var serverErr error
	go func() {
		serverErr = s.Start()
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Println("Linko is shutting down")

	if err := s.Shutdown(shutdownCtx); err != nil {
		logger.Printf("failed to shutdown server: %v\n", err)
		return 1
	}
	if serverErr != nil {
		logger.Printf("server error: %v\n", serverErr)
		return 1
	}
	return 0
}
