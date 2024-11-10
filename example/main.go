package main

import (
	"log/slog"

	humanslog "github.com/cschleiden/humanslog"
)

func main() {
	logger := slog.New(humanslog.New(slog.Default().Handler()))

	logger.Info("info message", slog.String("key", "value"), slog.Int("key2", 42))
	logger.Error("error message", slog.String("key", "value"), slog.Int("key2", 42))
	logger.Debug("debug message", slog.String("key", "value"), slog.Int("key2", 42))
	logger.Warn("warn message", slog.String("key", "value"), slog.Int("key2", 42))

	lattr := logger.With(slog.String("attr-key", "value"))

	lgroup := lattr.WithGroup("group")

	lgroup.Info("info message", slog.String("key", "value"), slog.Int("key2", 42))
	lgroup.Error("error message", slog.String("key", "value"), slog.Int("key2", 42))
	lgroup.Debug("debug message", slog.String("key", "value"), slog.Int("key2", 42))
	lgroup.Warn("warn message", slog.String("key", "value"), slog.Int("key2", 42))
}
