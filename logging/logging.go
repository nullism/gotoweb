package logging

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/fatih/color"
)

type Logger struct {
	*slog.Logger
}

func GetLogger() *Logger {
	return &Logger{
		slog.New(
			&CLIHandler{slog.NewTextHandler(os.Stdout, nil)},
		),
	}
}

type CLIHandler struct {
	slog.Handler
}

func (h *CLIHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String()

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	timeStr := r.Time.Format("[15:05:05.000]")
	msg := color.WhiteString(r.Message)

	println(timeStr, level, msg, string(b))

	return nil
}
