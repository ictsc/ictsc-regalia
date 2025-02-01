package pgxutil

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"

	"github.com/XSAM/otelsql"
	"github.com/jackc/pgx/v5"
	pgxstdlib "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

const (
	driverName = "pgx"
)

type openConfig struct {
	enableOTel  bool
	otelOptions []otelsql.Option
}

func NewDB(pgcfg pgx.ConnConfig, opts ...OpenOption) *sql.DB {
	var cfg openConfig
	for _, opt := range opts {
		opt.applyOpenOption(&cfg)
	}

	connector := newConnector(pgcfg)

	if cfg.enableOTel {
		opts := defaultOTelOptions(pgcfg)
		opts = append(opts, cfg.otelOptions...)
		return otelsql.OpenDB(connector, opts...)
	}

	return sql.OpenDB(connector)
}

func NewDBx(pgcfg pgx.ConnConfig, opts ...OpenOption) *sqlx.DB {
	return sqlx.NewDb(NewDB(pgcfg, opts...), driverName)
}

func newConnector(pgcfg pgx.ConnConfig) driver.Connector {
	return pgxstdlib.GetConnector(pgcfg)
}

func defaultOTelOptions(pgcfg pgx.ConnConfig) []otelsql.Option {
	return []otelsql.Option{
		otelsql.WithAttributes(
			semconv.DBSystemPostgreSQL,
			semconv.DBNamespace(pgcfg.Database),
			semconv.ServerAddress(pgcfg.Host),
			semconv.ServerPort(int(pgcfg.Port)),
		),
		otelsql.WithSpanNameFormatter(defaultSpanNameFormatter),
		otelsql.WithSpanOptions(otelsql.SpanOptions{
			DisableQuery: true,
		}),
		otelsql.WithAttributesGetter(defaultAttributeGetter),
	}
}

func defaultSpanNameFormatter(_ context.Context, method otelsql.Method, query string) string {
	summary := summarizeQuery(query)
	if summary == "" {
		return string(method)
	}
	return summary
}

func defaultAttributeGetter(_ context.Context, _ otelsql.Method, query string, _ []driver.NamedValue) []attribute.KeyValue {
	if query == "" {
		return []attribute.KeyValue{}
	}
	return []attribute.KeyValue{
		semconv.DBQueryText(query),
	}
}

var (
	queryRegexMap = map[string]*regexp.Regexp{
		"SELECT": regexp.MustCompile(`(?ims)^\s*SELECT\s+[\s\S]*?\sFROM\s+([^\s]+)`),
		"UPDATE": regexp.MustCompile(`(?ims)^\s*UPDATE\s+([^\s]+)`),
		"DELETE": regexp.MustCompile(`(?ims)^\s*DELETE\s+FROM\s+([^\s]+)`),
		"INSERT": regexp.MustCompile(`(?ims)^\s*INSERT\s+INTO\s+([^\s]+)`),
	}
)

func summarizeQuery(query string) string {
	for key, re := range queryRegexMap {
		if match := re.FindStringSubmatch(query); match != nil {
			table := strings.Trim(match[1], `"`)
			return fmt.Sprintf("%s %s", key, table)
		}
	}
	return ""
}
