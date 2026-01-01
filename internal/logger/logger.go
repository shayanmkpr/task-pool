package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	file *os.File
	*slog.Logger
}

// filename is ignored if toStdout is true
func New(filename string, toStdout bool) (*Logger, error) {
	var (
		handler slog.Handler
		f       *os.File
	)

	if toStdout {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		var err error
		f, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			return nil, err
		}

		handler = slog.NewJSONHandler(f, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

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
