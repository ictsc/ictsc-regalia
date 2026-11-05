// Package telemetry wires OpenTelemetry traces and metrics for the API.
//
// Export is configured only through the standard OTEL_* environment
// variables, e.g. OTEL_EXPORTER_OTLP_ENDPOINT. When none of them is set the
// global providers stay no-op so local development and tests do not try to
// reach a collector.
package telemetry

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
)

// DefaultServiceName is used unless OTEL_SERVICE_NAME or
// OTEL_RESOURCE_ATTRIBUTES overrides service.name.
const DefaultServiceName = "ictsc-regalia-api"

// exportEnvs enable export when any of them is non-empty.
var exportEnvs = []string{
	"OTEL_EXPORTER_OTLP_ENDPOINT",
	"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT",
	"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT",
	"OTEL_TRACES_EXPORTER",
	"OTEL_METRICS_EXPORTER",
}

// Enabled reports whether the environment asks for telemetry export.
func Enabled() bool {
	if strings.EqualFold(strings.TrimSpace(os.Getenv("OTEL_SDK_DISABLED")), "true") {
		return false
	}
	for _, key := range exportEnvs {
		if strings.TrimSpace(os.Getenv(key)) != "" {
			return true
		}
	}
	return false
}

// Setup installs the global propagator and, when Enabled, the global tracer
// and meter providers. The returned function flushes and stops them.
func Setup(ctx context.Context) (func(context.Context) error, error) {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{},
	))
	if !Enabled() {
		return func(context.Context) error { return nil }, nil
	}

	res, err := resource.New(ctx,
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithProcessRuntimeName(),
		resource.WithProcessRuntimeVersion(),
		resource.WithAttributes(semconv.ServiceName(DefaultServiceName)),
		// Environment wins over the defaults above.
		resource.WithFromEnv(),
	)
	if err != nil && !errors.Is(err, resource.ErrPartialResource) {
		return nil, fmt.Errorf("build telemetry resource: %w", err)
	}

	var shutdowns []func(context.Context) error
	shutdown := func(ctx context.Context) error {
		var errs error
		for i := len(shutdowns) - 1; i >= 0; i-- {
			errs = errors.Join(errs, shutdowns[i](ctx))
		}
		return errs
	}

	spanExporter, err := autoexport.NewSpanExporter(ctx)
	if err != nil {
		return nil, fmt.Errorf("create span exporter: %w", err)
	}
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(spanExporter),
		sdktrace.WithResource(res),
	)
	shutdowns = append(shutdowns, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	metricReader, err := autoexport.NewMetricReader(ctx)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("create metric reader: %w", err), shutdown(ctx))
	}
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(metricReader),
		sdkmetric.WithResource(res),
	)
	shutdowns = append(shutdowns, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	if err := runtime.Start(runtime.WithMeterProvider(meterProvider)); err != nil {
		return nil, errors.Join(fmt.Errorf("start runtime metrics: %w", err), shutdown(ctx))
	}
	return shutdown, nil
}
