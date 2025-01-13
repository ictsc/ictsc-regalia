package slogutil

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/errors/errbase"
	"github.com/go-slog/otelslog"
	"github.com/phsym/console-slog"
	slogformatter "github.com/samber/slog-formatter"
)

func NewHandler(out io.Writer, format Format, level slog.Leveler) slog.Handler {
	var handler slog.Handler
	switch format {
	case FormatJSON:
		handler = slog.NewJSONHandler(out, &slog.HandlerOptions{
			AddSource: true,
			Level:     level,
		})
	case FormatConsole:
		handler = slog.NewTextHandler(out, &slog.HandlerOptions{
			AddSource: true,
			Level:     level,
		})
	case FormatPretty:
		handler = console.NewHandler(out, &console.HandlerOptions{
			AddSource:  true,
			Level:      level,
			NoColor:    false,
			TimeFormat: time.RFC3339,
			Theme:      console.NewDefaultTheme(),
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
