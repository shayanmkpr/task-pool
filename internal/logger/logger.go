package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	file *os.File
	*slog.Logger
}

func New(filename string) (*Logger, error) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	handler := slog.NewJSONHandler(f, &slog.HandlerOptions{
		Level: slog.LevelInfo, // adjust as needed
	})

	return &Logger{
		file:   f,
		Logger: slog.New(handler),
	}, nil
}

func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}
