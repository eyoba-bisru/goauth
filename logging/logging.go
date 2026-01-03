package logging

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

// NewLogger creates a logger that writes to stdout and the given file path.
// It returns the logger, a closer function to close the underlying file, and an error if any.
func NewLogger(path string) (*log.Logger, func() error, error) {
	dir := filepath.Dir(path)
	if dir == "" || dir == "." {
		dir = "logs"
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, nil, err
	}
	lf, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}
	mw := io.MultiWriter(os.Stdout, lf)
	logger := log.New(mw, "", log.LstdFlags|log.Lmicroseconds)
	closer := func() error { return lf.Close() }
	return logger, closer, nil
}
