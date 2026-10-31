package logging

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

// New returns a logger that writes to stderr and, when logFile is non-empty,
// appends the same output to that file through a buffered writer. The returned
// cleanup function flushes the buffer and closes the file; callers must invoke it
// before exiting.
func New(logFile string) (*log.Logger, func() error, error) {
	if logFile == "" {
		return log.New(os.Stderr, "", log.LstdFlags), func() error { return nil }, nil
	}

	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open log file: %w", err)
	}

	bufferedFile := bufio.NewWriterSize(file, 8192)
	logger := log.New(io.MultiWriter(os.Stderr, bufferedFile), "", log.LstdFlags)

	cleanup := func() error {
		if err := bufferedFile.Flush(); err != nil {
			_ = file.Close()
			return err
		}
		return file.Close()
	}

	return logger, cleanup, nil
}
