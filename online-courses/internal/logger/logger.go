package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func MustInitLogger(LEVEL_LOG string) *slog.Logger {

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("ошибки инциализации логгера - запись в файл: %v", err))
	}

	level := slog.LevelInfo
	LEVEL_LOG = strings.ToLower(LEVEL_LOG)

	switch LEVEL_LOG {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: level,
	})

	return slog.New(handler)

}
