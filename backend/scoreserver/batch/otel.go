package batch

import "go.opentelemetry.io/otel"

const otelName = "github.com/ictsc/ictsc-regalia/backend/scoreserver/batch"

var (
	tracer = otel.Tracer(otelName)
	meter  = otel.Meter(otelName)
)
