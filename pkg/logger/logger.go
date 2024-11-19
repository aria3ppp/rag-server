package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/aria3ppp/rag-server/pkg/profile"
)

func NewLogger(w io.Writer) *slog.Logger {
	if w == nil {
		w = os.Stdout
	}

	addSource := false
	lvl := slog.LevelInfo

	if profile.IsDebug {
		addSource = true
		lvl = slog.LevelDebug
		if w != os.Stdout {
			w = io.MultiWriter(w, os.Stdout)
		}
	}

	return slog.New(slog.NewJSONHandler(
		w,
		&slog.HandlerOptions{
			AddSource: addSource,
			Level:     lvl,
		},
	))
}
