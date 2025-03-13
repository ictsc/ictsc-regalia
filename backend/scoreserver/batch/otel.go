package batch

import "go.opentelemetry.io/otel"

var tracer = otel.Tracer("github.com/ictsc/ictsc-regalia/backend/scoreserver/batch")
