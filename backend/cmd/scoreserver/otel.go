package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/cockroachdb/errors"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
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

const (
	exporterTypeNone = "none"
	exporterTypeOTLP = "otlp"

	otlpProtocolGRPC         = "grpc"
	otlpProtocolHTTPProtobuf = "http/protobuf"
	otlpProtocolHTTPJSON     = "http/json"
)

type noopSpanExporter struct{}

var _ trace.SpanExporter = &noopSpanExporter{}

func (n *noopSpanExporter) ExportSpans(context.Context, []trace.ReadOnlySpan) error {
	return nil
}

func (n *noopSpanExporter) Shutdown(context.Context) error {
	return nil
}

func newSpanExporter(ctx context.Context) (trace.SpanExporter, error) {
	expType := exporterTypeOTLP
	if t := os.Getenv("OTEL_TRACES_EXPORTER"); t != "" {
		expType = t
	}
	switch expType {
	case exporterTypeNone:
		return &noopSpanExporter{}, nil
	case exporterTypeOTLP:
		proto := otlpProtocolHTTPProtobuf
		if p := os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL"); p != "" {
			proto = p
		}
		if p := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_PROTOCOL"); p != "" {
			proto = p
		}
		switch proto {
		case otlpProtocolGRPC:
			exp, err := otlptracegrpc.New(ctx)
			return exp, errors.WithStack(err)
		case otlpProtocolHTTPProtobuf:
			exp, err := otlptracehttp.New(ctx)
			return exp, errors.WithStack(err)
		default:
			return nil, errors.Newf("unsupported OTLP protocol: %s", proto)
		}
	default:
		return nil, errors.Newf("unsupported trace exporter type: %s", expType)
	}
}

func newMericReader(ctx context.Context) (metric.Reader, error) {
	expTyp := exporterTypeOTLP
	if t := os.Getenv("OTEL_METRICS_EXPORTER"); t != "" {
		expTyp = t
	}
	switch expTyp {
	case exporterTypeNone:
		// ManualReader は単体では何もしない
		return metric.NewManualReader(), nil
	case exporterTypeOTLP:
		proto := otlpProtocolHTTPProtobuf
		if p := os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL"); p != "" {
			proto = p
		}
		if p := os.Getenv("OTEL_EXPORTER_OTLP_METRICS_PROTOCOL"); p != "" {
			proto = p
		}
		switch proto {
		case otlpProtocolHTTPProtobuf:
			exporter, err := otlpmetrichttp.New(ctx)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			return metric.NewPeriodicReader(exporter), nil
		case otlpProtocolGRPC:
			exporter, err := otlpmetricgrpc.New(ctx)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			return metric.NewPeriodicReader(exporter), nil
		default:
			return nil, errors.Newf("unsupported OTLP protocol: %s", proto)
		}
	default:
		return nil, errors.Newf("unsupported metric exporter type: %s", expTyp)
	}
}
