// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: admin/v1/problem.proto

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
	// ProblemServiceName is the fully-qualified name of the ProblemService service.
	ProblemServiceName = "admin.v1.ProblemService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// ProblemServiceGetProblemProcedure is the fully-qualified name of the ProblemService's GetProblem
	// RPC.
	ProblemServiceGetProblemProcedure = "/admin.v1.ProblemService/GetProblem"
	// ProblemServiceGetProblemsProcedure is the fully-qualified name of the ProblemService's
	// GetProblems RPC.
	ProblemServiceGetProblemsProcedure = "/admin.v1.ProblemService/GetProblems"
	// ProblemServicePatchProblemProcedure is the fully-qualified name of the ProblemService's
	// PatchProblem RPC.
	ProblemServicePatchProblemProcedure = "/admin.v1.ProblemService/PatchProblem"
	// ProblemServicePostProblemProcedure is the fully-qualified name of the ProblemService's
	// PostProblem RPC.
	ProblemServicePostProblemProcedure = "/admin.v1.ProblemService/PostProblem"
	// ProblemServiceDeleteProblemProcedure is the fully-qualified name of the ProblemService's
	// DeleteProblem RPC.
	ProblemServiceDeleteProblemProcedure = "/admin.v1.ProblemService/DeleteProblem"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	problemServiceServiceDescriptor             = v1.File_admin_v1_problem_proto.Services().ByName("ProblemService")
	problemServiceGetProblemMethodDescriptor    = problemServiceServiceDescriptor.Methods().ByName("GetProblem")
	problemServiceGetProblemsMethodDescriptor   = problemServiceServiceDescriptor.Methods().ByName("GetProblems")
	problemServicePatchProblemMethodDescriptor  = problemServiceServiceDescriptor.Methods().ByName("PatchProblem")
	problemServicePostProblemMethodDescriptor   = problemServiceServiceDescriptor.Methods().ByName("PostProblem")
	problemServiceDeleteProblemMethodDescriptor = problemServiceServiceDescriptor.Methods().ByName("DeleteProblem")
)

// ProblemServiceClient is a client for the admin.v1.ProblemService service.
type ProblemServiceClient interface {
	GetProblem(context.Context, *connect.Request[v1.GetProblemRequest]) (*connect.Response[v1.GetProblemResponse], error)
	GetProblems(context.Context, *connect.Request[v1.GetProblemsRequest]) (*connect.Response[v1.GetProblemsResponse], error)
	PatchProblem(context.Context, *connect.Request[v1.PatchProblemRequest]) (*connect.Response[v1.PatchProblemResponse], error)
	PostProblem(context.Context, *connect.Request[v1.PostProblemRequest]) (*connect.Response[v1.PostProblemResponse], error)
	DeleteProblem(context.Context, *connect.Request[v1.DeleteProblemRequest]) (*connect.Response[v1.DeleteProblemResponse], error)
}

// NewProblemServiceClient constructs a client for the admin.v1.ProblemService service. By default,
// it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and
// sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC()
// or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewProblemServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) ProblemServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &problemServiceClient{
		getProblem: connect.NewClient[v1.GetProblemRequest, v1.GetProblemResponse](
			httpClient,
			baseURL+ProblemServiceGetProblemProcedure,
			connect.WithSchema(problemServiceGetProblemMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		getProblems: connect.NewClient[v1.GetProblemsRequest, v1.GetProblemsResponse](
			httpClient,
			baseURL+ProblemServiceGetProblemsProcedure,
			connect.WithSchema(problemServiceGetProblemsMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		patchProblem: connect.NewClient[v1.PatchProblemRequest, v1.PatchProblemResponse](
			httpClient,
			baseURL+ProblemServicePatchProblemProcedure,
			connect.WithSchema(problemServicePatchProblemMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		postProblem: connect.NewClient[v1.PostProblemRequest, v1.PostProblemResponse](
			httpClient,
			baseURL+ProblemServicePostProblemProcedure,
			connect.WithSchema(problemServicePostProblemMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		deleteProblem: connect.NewClient[v1.DeleteProblemRequest, v1.DeleteProblemResponse](
			httpClient,
			baseURL+ProblemServiceDeleteProblemProcedure,
			connect.WithSchema(problemServiceDeleteProblemMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// problemServiceClient implements ProblemServiceClient.
type problemServiceClient struct {
	getProblem    *connect.Client[v1.GetProblemRequest, v1.GetProblemResponse]
	getProblems   *connect.Client[v1.GetProblemsRequest, v1.GetProblemsResponse]
	patchProblem  *connect.Client[v1.PatchProblemRequest, v1.PatchProblemResponse]
	postProblem   *connect.Client[v1.PostProblemRequest, v1.PostProblemResponse]
	deleteProblem *connect.Client[v1.DeleteProblemRequest, v1.DeleteProblemResponse]
}

// GetProblem calls admin.v1.ProblemService.GetProblem.
func (c *problemServiceClient) GetProblem(ctx context.Context, req *connect.Request[v1.GetProblemRequest]) (*connect.Response[v1.GetProblemResponse], error) {
	return c.getProblem.CallUnary(ctx, req)
}

// GetProblems calls admin.v1.ProblemService.GetProblems.
func (c *problemServiceClient) GetProblems(ctx context.Context, req *connect.Request[v1.GetProblemsRequest]) (*connect.Response[v1.GetProblemsResponse], error) {
	return c.getProblems.CallUnary(ctx, req)
}

// PatchProblem calls admin.v1.ProblemService.PatchProblem.
func (c *problemServiceClient) PatchProblem(ctx context.Context, req *connect.Request[v1.PatchProblemRequest]) (*connect.Response[v1.PatchProblemResponse], error) {
	return c.patchProblem.CallUnary(ctx, req)
}

// PostProblem calls admin.v1.ProblemService.PostProblem.
func (c *problemServiceClient) PostProblem(ctx context.Context, req *connect.Request[v1.PostProblemRequest]) (*connect.Response[v1.PostProblemResponse], error) {
	return c.postProblem.CallUnary(ctx, req)
}

// DeleteProblem calls admin.v1.ProblemService.DeleteProblem.
func (c *problemServiceClient) DeleteProblem(ctx context.Context, req *connect.Request[v1.DeleteProblemRequest]) (*connect.Response[v1.DeleteProblemResponse], error) {
	return c.deleteProblem.CallUnary(ctx, req)
}

// ProblemServiceHandler is an implementation of the admin.v1.ProblemService service.
type ProblemServiceHandler interface {
	GetProblem(context.Context, *connect.Request[v1.GetProblemRequest]) (*connect.Response[v1.GetProblemResponse], error)
	GetProblems(context.Context, *connect.Request[v1.GetProblemsRequest]) (*connect.Response[v1.GetProblemsResponse], error)
	PatchProblem(context.Context, *connect.Request[v1.PatchProblemRequest]) (*connect.Response[v1.PatchProblemResponse], error)
	PostProblem(context.Context, *connect.Request[v1.PostProblemRequest]) (*connect.Response[v1.PostProblemResponse], error)
	DeleteProblem(context.Context, *connect.Request[v1.DeleteProblemRequest]) (*connect.Response[v1.DeleteProblemResponse], error)
}

// NewProblemServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewProblemServiceHandler(svc ProblemServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	problemServiceGetProblemHandler := connect.NewUnaryHandler(
		ProblemServiceGetProblemProcedure,
		svc.GetProblem,
		connect.WithSchema(problemServiceGetProblemMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	problemServiceGetProblemsHandler := connect.NewUnaryHandler(
		ProblemServiceGetProblemsProcedure,
		svc.GetProblems,
		connect.WithSchema(problemServiceGetProblemsMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	problemServicePatchProblemHandler := connect.NewUnaryHandler(
		ProblemServicePatchProblemProcedure,
		svc.PatchProblem,
		connect.WithSchema(problemServicePatchProblemMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	problemServicePostProblemHandler := connect.NewUnaryHandler(
		ProblemServicePostProblemProcedure,
		svc.PostProblem,
		connect.WithSchema(problemServicePostProblemMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	problemServiceDeleteProblemHandler := connect.NewUnaryHandler(
		ProblemServiceDeleteProblemProcedure,
		svc.DeleteProblem,
		connect.WithSchema(problemServiceDeleteProblemMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/admin.v1.ProblemService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ProblemServiceGetProblemProcedure:
			problemServiceGetProblemHandler.ServeHTTP(w, r)
		case ProblemServiceGetProblemsProcedure:
			problemServiceGetProblemsHandler.ServeHTTP(w, r)
		case ProblemServicePatchProblemProcedure:
			problemServicePatchProblemHandler.ServeHTTP(w, r)
		case ProblemServicePostProblemProcedure:
			problemServicePostProblemHandler.ServeHTTP(w, r)
		case ProblemServiceDeleteProblemProcedure:
			problemServiceDeleteProblemHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedProblemServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedProblemServiceHandler struct{}

func (UnimplementedProblemServiceHandler) GetProblem(context.Context, *connect.Request[v1.GetProblemRequest]) (*connect.Response[v1.GetProblemResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("admin.v1.ProblemService.GetProblem is not implemented"))
}

func (UnimplementedProblemServiceHandler) GetProblems(context.Context, *connect.Request[v1.GetProblemsRequest]) (*connect.Response[v1.GetProblemsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("admin.v1.ProblemService.GetProblems is not implemented"))
}

func (UnimplementedProblemServiceHandler) PatchProblem(context.Context, *connect.Request[v1.PatchProblemRequest]) (*connect.Response[v1.PatchProblemResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("admin.v1.ProblemService.PatchProblem is not implemented"))
}

func (UnimplementedProblemServiceHandler) PostProblem(context.Context, *connect.Request[v1.PostProblemRequest]) (*connect.Response[v1.PostProblemResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("admin.v1.ProblemService.PostProblem is not implemented"))
}

func (UnimplementedProblemServiceHandler) DeleteProblem(context.Context, *connect.Request[v1.DeleteProblemRequest]) (*connect.Response[v1.DeleteProblemResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("admin.v1.ProblemService.DeleteProblem is not implemented"))
}