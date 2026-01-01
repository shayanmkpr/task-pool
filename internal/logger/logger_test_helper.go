package logger

import (
	"io"
	"log/slog"
)

// NewTestLogger creates a logger that writes to io.Discard for testing
func NewTestLogger() *Logger {
	handler := slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return &Logger{
		file:   nil, // No file for tests
		Logger: slog.New(handler),
	}
}

// NewTestLoggerWithOutput creates a logger that writes to the provided writer for testing
func NewTestLoggerWithOutput(w io.Writer) *Logger {
	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return &Logger{
		file:   nil,
		Logger: slog.New(handler),
	}
}
