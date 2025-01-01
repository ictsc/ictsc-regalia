// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: admin/v1/contestant.proto

package adminv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "github.com/ictsc/ictsc-regalia/backend/pkg/proto/admin/v1"
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
	// ContestantServiceName is the fully-qualified name of the ContestantService service.
	ContestantServiceName = "admin.v1.ContestantService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// ContestantServiceListContestantsProcedure is the fully-qualified name of the ContestantService's
	// ListContestants RPC.
	ContestantServiceListContestantsProcedure = "/admin.v1.ContestantService/ListContestants"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	contestantServiceServiceDescriptor               = v1.File_admin_v1_contestant_proto.Services().ByName("ContestantService")
	contestantServiceListContestantsMethodDescriptor = contestantServiceServiceDescriptor.Methods().ByName("ListContestants")
)

// ContestantServiceClient is a client for the admin.v1.ContestantService service.
type ContestantServiceClient interface {
	ListContestants(context.Context, *connect.Request[v1.ListContestantsRequest]) (*connect.Response[v1.ListContestantsResponse], error)
}

// NewContestantServiceClient constructs a client for the admin.v1.ContestantService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewContestantServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) ContestantServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &contestantServiceClient{
		listContestants: connect.NewClient[v1.ListContestantsRequest, v1.ListContestantsResponse](
			httpClient,
			baseURL+ContestantServiceListContestantsProcedure,
			connect.WithSchema(contestantServiceListContestantsMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// contestantServiceClient implements ContestantServiceClient.
type contestantServiceClient struct {
	listContestants *connect.Client[v1.ListContestantsRequest, v1.ListContestantsResponse]
}

// ListContestants calls admin.v1.ContestantService.ListContestants.
func (c *contestantServiceClient) ListContestants(ctx context.Context, req *connect.Request[v1.ListContestantsRequest]) (*connect.Response[v1.ListContestantsResponse], error) {
	return c.listContestants.CallUnary(ctx, req)
}

// ContestantServiceHandler is an implementation of the admin.v1.ContestantService service.
type ContestantServiceHandler interface {
	ListContestants(context.Context, *connect.Request[v1.ListContestantsRequest]) (*connect.Response[v1.ListContestantsResponse], error)
}

// NewContestantServiceHandler builds an HTTP handler from the service implementation. It returns
// the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewContestantServiceHandler(svc ContestantServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	contestantServiceListContestantsHandler := connect.NewUnaryHandler(
		ContestantServiceListContestantsProcedure,
		svc.ListContestants,
		connect.WithSchema(contestantServiceListContestantsMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/admin.v1.ContestantService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ContestantServiceListContestantsProcedure:
			contestantServiceListContestantsHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedContestantServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedContestantServiceHandler struct{}

func (UnimplementedContestantServiceHandler) ListContestants(context.Context, *connect.Request[v1.ListContestantsRequest]) (*connect.Response[v1.ListContestantsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("admin.v1.ContestantService.ListContestants is not implemented"))
}