package loggers

import (
	"app/news-parser/internal/response"
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
	*DataLog
}
type DataLog struct {
	UserUUID string
	Errors   []response.Error
	MapLog   map[string]any
}

func NewLogger() *Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return &Logger{
		Logger: logger,
		DataLog: &DataLog{
			MapLog: make(map[string]any),
		},
	}
}

func (l *Logger) HandlerLogger(path, userUUID string, status int, errors []response.Error, data map[string]any) {
	if status < 500 {
		if status >= 200 && status < 300 {
			l.Info("request successful",
				slog.String("path", path),
				slog.Int("status", status),
				slog.String("user_uuid", userUUID),
			)
		} else {
			l.Info("request failed",
				slog.String("path", path),
				slog.Int("status", status),
				slog.String("user_uuid", userUUID),
				slog.Any("errors", errors),
				slog.Any("data", data),
			)
		}
	} else {
		l.Error("request failed",
			slog.String("path", path),
			slog.Int("status", status),
			slog.String("user_uuid", userUUID),
			slog.Any("errors", errors),
			slog.Any("data", data),
		)
	}
}
func (l *Logger) SystemLogger(level slog.Level, msg string) {
	switch level {
	case slog.LevelInfo:
		l.Info(msg)
	case slog.LevelError:
		l.Error(msg)
	}
}
