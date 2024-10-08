// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: admin/v1/contestant.proto

package v1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "github.com/ictsc/ictsc-outlands/backend/internal/proto/admin/v1"
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
	// ContestantServiceGetContestantProcedure is the fully-qualified name of the ContestantService's
	// GetContestant RPC.
	ContestantServiceGetContestantProcedure = "/admin.v1.ContestantService/GetContestant"
	// ContestantServiceGetContestantsProcedure is the fully-qualified name of the ContestantService's
	// GetContestants RPC.
	ContestantServiceGetContestantsProcedure = "/admin.v1.ContestantService/GetContestants"
	// ContestantServiceDeleteContestantProcedure is the fully-qualified name of the ContestantService's
	// DeleteContestant RPC.
	ContestantServiceDeleteContestantProcedure = "/admin.v1.ContestantService/DeleteContestant"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	contestantServiceServiceDescriptor                = v1.File_admin_v1_contestant_proto.Services().ByName("ContestantService")
	contestantServiceGetContestantMethodDescriptor    = contestantServiceServiceDescriptor.Methods().ByName("GetContestant")
	contestantServiceGetContestantsMethodDescriptor   = contestantServiceServiceDescriptor.Methods().ByName("GetContestants")
	contestantServiceDeleteContestantMethodDescriptor = contestantServiceServiceDescriptor.Methods().ByName("DeleteContestant")
)

// ContestantServiceClient is a client for the admin.v1.ContestantService service.
type ContestantServiceClient interface {
	GetContestant(context.Context, *connect.Request[v1.GetContestantRequest]) (*connect.Response[v1.GetContestantResponse], error)
	GetContestants(context.Context, *connect.Request[v1.GetContestantsRequest]) (*connect.Response[v1.GetContestantsResponse], error)
	DeleteContestant(context.Context, *connect.Request[v1.DeleteContestantRequest]) (*connect.Response[v1.DeleteContestantResponse], error)
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
		getContestant: connect.NewClient[v1.GetContestantRequest, v1.GetContestantResponse](
			httpClient,
			baseURL+ContestantServiceGetContestantProcedure,
			connect.WithSchema(contestantServiceGetContestantMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		getContestants: connect.NewClient[v1.GetContestantsRequest, v1.GetContestantsResponse](
			httpClient,
			baseURL+ContestantServiceGetContestantsProcedure,
			connect.WithSchema(contestantServiceGetContestantsMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		deleteContestant: connect.NewClient[v1.DeleteContestantRequest, v1.DeleteContestantResponse](
			httpClient,
			baseURL+ContestantServiceDeleteContestantProcedure,
			connect.WithSchema(contestantServiceDeleteContestantMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// contestantServiceClient implements ContestantServiceClient.
type contestantServiceClient struct {
	getContestant    *connect.Client[v1.GetContestantRequest, v1.GetContestantResponse]
	getContestants   *connect.Client[v1.GetContestantsRequest, v1.GetContestantsResponse]
	deleteContestant *connect.Client[v1.DeleteContestantRequest, v1.DeleteContestantResponse]
}

// GetContestant calls admin.v1.ContestantService.GetContestant.
func (c *contestantServiceClient) GetContestant(ctx context.Context, req *connect.Request[v1.GetContestantRequest]) (*connect.Response[v1.GetContestantResponse], error) {
	return c.getContestant.CallUnary(ctx, req)
}

// GetContestants calls admin.v1.ContestantService.GetContestants.
func (c *contestantServiceClient) GetContestants(ctx context.Context, req *connect.Request[v1.GetContestantsRequest]) (*connect.Response[v1.GetContestantsResponse], error) {
	return c.getContestants.CallUnary(ctx, req)
}

// DeleteContestant calls admin.v1.ContestantService.DeleteContestant.
func (c *contestantServiceClient) DeleteContestant(ctx context.Context, req *connect.Request[v1.DeleteContestantRequest]) (*connect.Response[v1.DeleteContestantResponse], error) {
	return c.deleteContestant.CallUnary(ctx, req)
}

// ContestantServiceHandler is an implementation of the admin.v1.ContestantService service.
type ContestantServiceHandler interface {
	GetContestant(context.Context, *connect.Request[v1.GetContestantRequest]) (*connect.Response[v1.GetContestantResponse], error)
	GetContestants(context.Context, *connect.Request[v1.GetContestantsRequest]) (*connect.Response[v1.GetContestantsResponse], error)
	DeleteContestant(context.Context, *connect.Request[v1.DeleteContestantRequest]) (*connect.Response[v1.DeleteContestantResponse], error)
}

// NewContestantServiceHandler builds an HTTP handler from the service implementation. It returns
// the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewContestantServiceHandler(svc ContestantServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	contestantServiceGetContestantHandler := connect.NewUnaryHandler(
		ContestantServiceGetContestantProcedure,
		svc.GetContestant,
		connect.WithSchema(contestantServiceGetContestantMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	contestantServiceGetContestantsHandler := connect.NewUnaryHandler(
		ContestantServiceGetContestantsProcedure,
		svc.GetContestants,
		connect.WithSchema(contestantServiceGetContestantsMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	contestantServiceDeleteContestantHandler := connect.NewUnaryHandler(
		ContestantServiceDeleteContestantProcedure,
		svc.DeleteContestant,
		connect.WithSchema(contestantServiceDeleteContestantMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/admin.v1.ContestantService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ContestantServiceGetContestantProcedure:
			contestantServiceGetContestantHandler.ServeHTTP(w, r)
		case ContestantServiceGetContestantsProcedure:
			contestantServiceGetContestantsHandler.ServeHTTP(w, r)
		case ContestantServiceDeleteContestantProcedure:
			contestantServiceDeleteContestantHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedContestantServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedContestantServiceHandler struct{}

func (UnimplementedContestantServiceHandler) GetContestant(context.Context, *connect.Request[v1.GetContestantRequest]) (*connect.Response[v1.GetContestantResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("admin.v1.ContestantService.GetContestant is not implemented"))
}

func (UnimplementedContestantServiceHandler) GetContestants(context.Context, *connect.Request[v1.GetContestantsRequest]) (*connect.Response[v1.GetContestantsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("admin.v1.ContestantService.GetContestants is not implemented"))
}

func (UnimplementedContestantServiceHandler) DeleteContestant(context.Context, *connect.Request[v1.DeleteContestantRequest]) (*connect.Response[v1.DeleteContestantResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("admin.v1.ContestantService.DeleteContestant is not implemented"))
}
