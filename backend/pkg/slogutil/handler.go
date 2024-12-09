package slogutil

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/errors/errbase"
	"github.com/go-slog/otelslog"
	"github.com/phsym/console-slog"
	slogformatter "github.com/samber/slog-formatter"
)

func NewHandler(dev bool, level slog.Leveler) slog.Handler {
	var handler slog.Handler
	if dev {
		handler = console.NewHandler(os.Stderr, &console.HandlerOptions{
			AddSource:  true,
			Level:      level,
			NoColor:    false,
			TimeFormat: time.RFC3339,
			Theme:      console.NewDefaultTheme(),
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   true,
			Level:       level,
			ReplaceAttr: nil,
		})
	}

	handler = slogformatter.NewFormatterHandler(
		slogformatter.FormatByType(formatError),
	)(handler)

	handler = otelslog.NewHandler(handler)

	return handler
}

type withStackError interface {
	SafeFormatError(p errbase.Printer) error
}

func formatError(err error) slog.Value {
	attrs := make([]slog.Attr, 0, 2) //nolint:mnd // 全てのエラーはスタックトレースを持つはずなので属性は2つある
	attrs = append(attrs, slog.String("message", err.Error()))

	var stackErr withStackError
	if errors.As(err, &stackErr) {
		trace := fmt.Sprintf("%+v", stackErr)
		_, trace, _ = strings.Cut(trace, "\n")
		attrs = append(attrs, slog.String("stacktrace", trace))
	}

	return slog.GroupValue(attrs...)
}
