package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

const defaultServiceName = "scoreserver"

// nolint:nonamedreturns
func setupOpenTelemetry(ctx context.Context) (shutdown func(context.Context) error, err error) {
	var shutdownFns []func(context.Context) error
	shutdownFn := func(ctx context.Context) error {
		errs := make([]error, 0, len(shutdownFns))
		for _, fn := range shutdownFns {
			errs = append(errs, fn(ctx))
		}
		return errors.Join(errs...)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, shutdownFn(ctx))
		}
	}()

	res, err := newResource(ctx)
	if err != nil {
		return nil, err
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(cause error) {
		slog.Error("OpenTelemetry error", "error", cause)
	}))

	spanExporter, err := newSpanExporter(ctx)
	if err != nil {
		return nil, err
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithBatcher(spanExporter),
	)
	otel.SetTracerProvider(tracerProvider)
	shutdownFns = append(shutdownFns, tracerProvider.Shutdown)

	metricReader, err := newMericReader(ctx)
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metricReader),
	)
	otel.SetMeterProvider(meterProvider)
	shutdownFns = append(shutdownFns, meterProvider.Shutdown)

	if err := runtime.Start(); err != nil {
		return nil, errors.Wrap(err, "failed to start runtime instrumentation")
	}

	return shutdownFn, nil
}

func newResource(ctx context.Context) (*resource.Resource, error) {
	res, err := resource.New(
		ctx,
		resource.WithAttributes(semconv.ServiceName(defaultServiceName)),
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
	)
	return res, errors.WithStack(err)
}

type noopSpanExporter struct{}

var _ trace.SpanExporter = &noopSpanExporter{}

func (n *noopSpanExporter) ExportSpans(context.Context, []trace.ReadOnlySpan) error {
	return nil
}

func (n *noopSpanExporter) Shutdown(context.Context) error {
	return nil
}

func newSpanExporter(ctx context.Context) (trace.SpanExporter, error) {
	switch os.Getenv("OTEL_TRACES_EXPORTER") {
	case "none":
		return &noopSpanExporter{}, nil
	case "console":
		exp, err := stdouttrace.New()
		return exp, errors.WithStack(err)
	default /* otlp */ :
		exp, err := otlptracegrpc.New(ctx)
		return exp, errors.WithStack(err)
	}
}

func newMericReader(ctx context.Context) (metric.Reader, error) {
	switch os.Getenv("OTEL_METRICS_EXPORTER") {
	case "none":
		// ManualReader は単体では何もしない
		return metric.NewManualReader(), nil
	case "prometheus":
		return newPrometheusMetricExporter()
	default: /* otlp */
		exporter, err := otlpmetricgrpc.New(ctx)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return metric.NewPeriodicReader(exporter), nil
	}
}

type prometheusReader struct {
	otelprom.Exporter
	s *http.Server
}

func newPrometheusMetricExporter(opts ...otelprom.Option) (*prometheusReader, error) {
	registry := prometheus.NewRegistry()

	opts = append(opts, otelprom.WithRegisterer(registry), otelprom.WithoutScopeInfo())
	exporter, err := otelprom.New(opts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create prometheus exporter")
	}

	host := "localhost"
	if h, ok := os.LookupEnv("OTEL_EXPORTER_PROMETHEUS_HOST"); ok {
		host = h
	}
	port := "9464"
	if p, ok := os.LookupEnv("OTEL_EXPORTER_PROMETHEUS_PORT"); ok {
		port = p
	}
	// nolint:gosec
	srv := &http.Server{
		Addr: net.JoinHostPort(host, port),
		Handler: promhttp.HandlerFor(registry, promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		}),
	}

	go func() {
		slog.Info("Starting prometheus server", "address", srv.Addr)
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start prometheus server", "error", err)
		}
	}()

	return &prometheusReader{
		Exporter: *exporter,
		s:        srv,
	}, nil
}

func (r *prometheusReader) Shutdown(ctx context.Context) error {
	if err := r.Exporter.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "failed to shutdown prometheus exporter")
	}
	if err := r.s.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "failed to shutdown prometheus server")
	}
	return nil
}
