package logging

import (
	"fmt"
	"io"
	"log"
	"os"
)

// New returns a logger that writes to stderr and, when logFile is non-empty,
// appends the same output to that file.
func New(logFile string) (*log.Logger, error) {
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		return log.New(io.MultiWriter(os.Stderr, file), "", log.LstdFlags), nil
	}
	return log.New(os.Stderr, "", log.LstdFlags), nil
}
