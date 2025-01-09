// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: contestant/v1/contest.proto

package contestantv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/contestant/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// ContestServiceName is the fully-qualified name of the ContestService service.
	ContestServiceName = "contestant.v1.ContestService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// ContestServiceGetScheduleProcedure is the fully-qualified name of the ContestService's
	// GetSchedule RPC.
	ContestServiceGetScheduleProcedure = "/contestant.v1.ContestService/GetSchedule"
)

// ContestServiceClient is a client for the contestant.v1.ContestService service.
type ContestServiceClient interface {
	GetSchedule(context.Context, *connect.Request[v1.GetScheduleRequest]) (*connect.Response[v1.GetScheduleResponse], error)
}

// NewContestServiceClient constructs a client for the contestant.v1.ContestService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewContestServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) ContestServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	contestServiceMethods := v1.File_contestant_v1_contest_proto.Services().ByName("ContestService").Methods()
	return &contestServiceClient{
		getSchedule: connect.NewClient[v1.GetScheduleRequest, v1.GetScheduleResponse](
			httpClient,
			baseURL+ContestServiceGetScheduleProcedure,
			connect.WithSchema(contestServiceMethods.ByName("GetSchedule")),
			connect.WithClientOptions(opts...),
		),
	}
}

// contestServiceClient implements ContestServiceClient.
type contestServiceClient struct {
	getSchedule *connect.Client[v1.GetScheduleRequest, v1.GetScheduleResponse]
}

// GetSchedule calls contestant.v1.ContestService.GetSchedule.
func (c *contestServiceClient) GetSchedule(ctx context.Context, req *connect.Request[v1.GetScheduleRequest]) (*connect.Response[v1.GetScheduleResponse], error) {
	return c.getSchedule.CallUnary(ctx, req)
}

// ContestServiceHandler is an implementation of the contestant.v1.ContestService service.
type ContestServiceHandler interface {
	GetSchedule(context.Context, *connect.Request[v1.GetScheduleRequest]) (*connect.Response[v1.GetScheduleResponse], error)
}

// NewContestServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewContestServiceHandler(svc ContestServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	contestServiceMethods := v1.File_contestant_v1_contest_proto.Services().ByName("ContestService").Methods()
	contestServiceGetScheduleHandler := connect.NewUnaryHandler(
		ContestServiceGetScheduleProcedure,
		svc.GetSchedule,
		connect.WithSchema(contestServiceMethods.ByName("GetSchedule")),
		connect.WithHandlerOptions(opts...),
	)
	return "/contestant.v1.ContestService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ContestServiceGetScheduleProcedure:
			contestServiceGetScheduleHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedContestServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedContestServiceHandler struct{}

func (UnimplementedContestServiceHandler) GetSchedule(context.Context, *connect.Request[v1.GetScheduleRequest]) (*connect.Response[v1.GetScheduleResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("contestant.v1.ContestService.GetSchedule is not implemented"))
}
