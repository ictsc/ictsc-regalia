package connectutil

import (
	"context"
	"net/http"
	"strings"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/protobuf/proto"
)

const (
	instrumentationName = "github.com/ictsc/ictsc-regalia/backend/pkg/connectutil"
	rpcMessageEvent     = "rpc.message"
)

var meter = otel.Meter(instrumentationName)

// NewOtelInterceptor returns an interceptor that instruments OpenTelemetry metrics and traces.
//
// otelconnect のメトリック計装がおかしいので，メトリックは自分で計装する
func NewOtelInterceptor() connect.Interceptor {
	serverDurationMetric, err := meter.Int64Histogram(
		semconv.RPCServerDurationName,
		metric.WithUnit(semconv.RPCServerDurationUnit),
		metric.WithDescription(semconv.RPCServerDurationDescription),
	)
	if err != nil {
		otel.Handle(err)
		serverDurationMetric = noop.Int64Histogram{}
	}

	serverRequestSizeMetric, err := meter.Int64Histogram(
		semconv.RPCServerRequestSizeName,
		metric.WithUnit(semconv.RPCServerRequestSizeUnit),
		metric.WithDescription(semconv.RPCServerRequestSizeDescription),
	)
	if err != nil {
		otel.Handle(err)
		serverRequestSizeMetric = noop.Int64Histogram{}
	}

	serverResponseSizeMetric, err := meter.Int64Histogram(
		semconv.RPCServerResponseSizeName,
		metric.WithUnit(semconv.RPCServerResponseSizeName),
		metric.WithDescription(semconv.RPCServerResponseSizeDescription),
	)
	if err != nil {
		otel.Handle(err)
		serverResponseSizeMetric = noop.Int64Histogram{}
	}

	traceInterceptor, err := otelconnect.NewInterceptor(
		otelconnect.WithoutMetrics(),
		otelconnect.WithTrustRemote(),
	)
	if err != nil {
		otel.Handle(err)
		traceInterceptor = nil
	}

	return &otelInterceptor{
		Interceptor:              traceInterceptor,
		serverDurationMetric:     serverDurationMetric,
		serverRequestSizeMetric:  serverRequestSizeMetric,
		serverResponseSizeMetric: serverResponseSizeMetric,
	}
}

type otelInterceptor struct {
	*otelconnect.Interceptor
	serverDurationMetric     metric.Int64Histogram
	serverRequestSizeMetric  metric.Int64Histogram
	serverResponseSizeMetric metric.Int64Histogram
}

var now = time.Now

func (o *otelInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	if o.Interceptor != nil {
		next = o.Interceptor.WrapUnary(next)
	}
	unaryFunc := func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
		requestStartTime := now()

		response, err := next(ctx, request)

		requestEndTime := now()

		var requestSize int
		if request != nil {
			if msg, ok := request.Any().(proto.Message); ok {
				requestSize = proto.Size(msg)
			}
		}

		var responseSize int
		if err == nil {
			if msg, ok := response.Any().(proto.Message); ok {
				responseSize = proto.Size(msg)
			}
		}

		duration := requestEndTime.Sub(requestStartTime)

		service, method := splitProceduce(request.Spec().Procedure)

		protocol := protocolToSemConv(request.Peer().Protocol)

		attrs := []attribute.KeyValue{
			semconv.RPCSystemKey.String(protocol),
			semconv.RPCService(service),
			semconv.RPCMethod(method),
		}
		if attr, ok := statusCodeAttribute(protocol, err); ok {
			attrs = append(attrs, attr)
		}

		o.serverDurationMetric.Record(ctx, duration.Milliseconds(), metric.WithAttributes(attrs...))
		o.serverRequestSizeMetric.Record(ctx, int64(requestSize), metric.WithAttributes(attrs...))
		o.serverResponseSizeMetric.Record(ctx, int64(responseSize), metric.WithAttributes(attrs...))

		return response, err
	}
	return connect.UnaryFunc(unaryFunc)
}

func statusCodeAttribute(protocol string, err error) (attribute.KeyValue, bool) {
	switch protocol {
	case grpcProtocol, grpcwebProtocol:
		if err != nil {
			return semconv.RPCGRPCStatusCodeKey.Int64(int64(connect.CodeOf(err))), true
		}
		return semconv.RPCGRPCStatusCodeOk, true
	case connectProtocol:
		if connect.IsNotModifiedError(err) {
			return semconv.HTTPResponseStatusCode(http.StatusNotModified), true
		}
		if err != nil {
			return semconv.RPCConnectRPCErrorCodeKey.String(connect.CodeOf(err).String()), true
		}
	}
	return attribute.KeyValue{}, false
}

const (
	grpcwebString   = "grpcweb"
	grpcwebProtocol = "grpc_web"
	grpcString      = "grpc"
	grpcProtocol    = "grpc"
	connectString   = "connect"
	connectProtocol = "connect_rpc"
)

func protocolToSemConv(protocol string) string {
	switch protocol {
	case grpcwebString:
		return grpcwebProtocol
	case grpcString:
		return grpcProtocol
	case connectString:
		return connectProtocol
	default:
		return protocol
	}
}

//nolint:nonamedreturns // 名前が付いたほうが適切
func splitProceduce(proc string) (service, method string) {
	proc = strings.TrimLeft(proc, "/")
	service, method, _ = strings.Cut(proc, "/")
	return
}
