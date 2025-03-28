// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: contestant/v1/problem.proto

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
	// ProblemServiceName is the fully-qualified name of the ProblemService service.
	ProblemServiceName = "contestant.v1.ProblemService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// ProblemServiceListProblemsProcedure is the fully-qualified name of the ProblemService's
	// ListProblems RPC.
	ProblemServiceListProblemsProcedure = "/contestant.v1.ProblemService/ListProblems"
	// ProblemServiceGetProblemProcedure is the fully-qualified name of the ProblemService's GetProblem
	// RPC.
	ProblemServiceGetProblemProcedure = "/contestant.v1.ProblemService/GetProblem"
	// ProblemServiceListDeploymentsProcedure is the fully-qualified name of the ProblemService's
	// ListDeployments RPC.
	ProblemServiceListDeploymentsProcedure = "/contestant.v1.ProblemService/ListDeployments"
	// ProblemServiceDeployProcedure is the fully-qualified name of the ProblemService's Deploy RPC.
	ProblemServiceDeployProcedure = "/contestant.v1.ProblemService/Deploy"
)

// ProblemServiceClient is a client for the contestant.v1.ProblemService service.
type ProblemServiceClient interface {
	ListProblems(context.Context, *connect.Request[v1.ListProblemsRequest]) (*connect.Response[v1.ListProblemsResponse], error)
	GetProblem(context.Context, *connect.Request[v1.GetProblemRequest]) (*connect.Response[v1.GetProblemResponse], error)
	ListDeployments(context.Context, *connect.Request[v1.ListDeploymentsRequest]) (*connect.Response[v1.ListDeploymentsResponse], error)
	Deploy(context.Context, *connect.Request[v1.DeployRequest]) (*connect.Response[v1.DeployResponse], error)
}

// NewProblemServiceClient constructs a client for the contestant.v1.ProblemService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewProblemServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) ProblemServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	problemServiceMethods := v1.File_contestant_v1_problem_proto.Services().ByName("ProblemService").Methods()
	return &problemServiceClient{
		listProblems: connect.NewClient[v1.ListProblemsRequest, v1.ListProblemsResponse](
			httpClient,
			baseURL+ProblemServiceListProblemsProcedure,
			connect.WithSchema(problemServiceMethods.ByName("ListProblems")),
			connect.WithClientOptions(opts...),
		),
		getProblem: connect.NewClient[v1.GetProblemRequest, v1.GetProblemResponse](
			httpClient,
			baseURL+ProblemServiceGetProblemProcedure,
			connect.WithSchema(problemServiceMethods.ByName("GetProblem")),
			connect.WithClientOptions(opts...),
		),
		listDeployments: connect.NewClient[v1.ListDeploymentsRequest, v1.ListDeploymentsResponse](
			httpClient,
			baseURL+ProblemServiceListDeploymentsProcedure,
			connect.WithSchema(problemServiceMethods.ByName("ListDeployments")),
			connect.WithClientOptions(opts...),
		),
		deploy: connect.NewClient[v1.DeployRequest, v1.DeployResponse](
			httpClient,
			baseURL+ProblemServiceDeployProcedure,
			connect.WithSchema(problemServiceMethods.ByName("Deploy")),
			connect.WithClientOptions(opts...),
		),
	}
}

// problemServiceClient implements ProblemServiceClient.
type problemServiceClient struct {
	listProblems    *connect.Client[v1.ListProblemsRequest, v1.ListProblemsResponse]
	getProblem      *connect.Client[v1.GetProblemRequest, v1.GetProblemResponse]
	listDeployments *connect.Client[v1.ListDeploymentsRequest, v1.ListDeploymentsResponse]
	deploy          *connect.Client[v1.DeployRequest, v1.DeployResponse]
}

// ListProblems calls contestant.v1.ProblemService.ListProblems.
func (c *problemServiceClient) ListProblems(ctx context.Context, req *connect.Request[v1.ListProblemsRequest]) (*connect.Response[v1.ListProblemsResponse], error) {
	return c.listProblems.CallUnary(ctx, req)
}

// GetProblem calls contestant.v1.ProblemService.GetProblem.
func (c *problemServiceClient) GetProblem(ctx context.Context, req *connect.Request[v1.GetProblemRequest]) (*connect.Response[v1.GetProblemResponse], error) {
	return c.getProblem.CallUnary(ctx, req)
}

// ListDeployments calls contestant.v1.ProblemService.ListDeployments.
func (c *problemServiceClient) ListDeployments(ctx context.Context, req *connect.Request[v1.ListDeploymentsRequest]) (*connect.Response[v1.ListDeploymentsResponse], error) {
	return c.listDeployments.CallUnary(ctx, req)
}

// Deploy calls contestant.v1.ProblemService.Deploy.
func (c *problemServiceClient) Deploy(ctx context.Context, req *connect.Request[v1.DeployRequest]) (*connect.Response[v1.DeployResponse], error) {
	return c.deploy.CallUnary(ctx, req)
}

// ProblemServiceHandler is an implementation of the contestant.v1.ProblemService service.
type ProblemServiceHandler interface {
	ListProblems(context.Context, *connect.Request[v1.ListProblemsRequest]) (*connect.Response[v1.ListProblemsResponse], error)
	GetProblem(context.Context, *connect.Request[v1.GetProblemRequest]) (*connect.Response[v1.GetProblemResponse], error)
	ListDeployments(context.Context, *connect.Request[v1.ListDeploymentsRequest]) (*connect.Response[v1.ListDeploymentsResponse], error)
	Deploy(context.Context, *connect.Request[v1.DeployRequest]) (*connect.Response[v1.DeployResponse], error)
}

// NewProblemServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewProblemServiceHandler(svc ProblemServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	problemServiceMethods := v1.File_contestant_v1_problem_proto.Services().ByName("ProblemService").Methods()
	problemServiceListProblemsHandler := connect.NewUnaryHandler(
		ProblemServiceListProblemsProcedure,
		svc.ListProblems,
		connect.WithSchema(problemServiceMethods.ByName("ListProblems")),
		connect.WithHandlerOptions(opts...),
	)
	problemServiceGetProblemHandler := connect.NewUnaryHandler(
		ProblemServiceGetProblemProcedure,
		svc.GetProblem,
		connect.WithSchema(problemServiceMethods.ByName("GetProblem")),
		connect.WithHandlerOptions(opts...),
	)
	problemServiceListDeploymentsHandler := connect.NewUnaryHandler(
		ProblemServiceListDeploymentsProcedure,
		svc.ListDeployments,
		connect.WithSchema(problemServiceMethods.ByName("ListDeployments")),
		connect.WithHandlerOptions(opts...),
	)
	problemServiceDeployHandler := connect.NewUnaryHandler(
		ProblemServiceDeployProcedure,
		svc.Deploy,
		connect.WithSchema(problemServiceMethods.ByName("Deploy")),
		connect.WithHandlerOptions(opts...),
	)
	return "/contestant.v1.ProblemService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ProblemServiceListProblemsProcedure:
			problemServiceListProblemsHandler.ServeHTTP(w, r)
		case ProblemServiceGetProblemProcedure:
			problemServiceGetProblemHandler.ServeHTTP(w, r)
		case ProblemServiceListDeploymentsProcedure:
			problemServiceListDeploymentsHandler.ServeHTTP(w, r)
		case ProblemServiceDeployProcedure:
			problemServiceDeployHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedProblemServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedProblemServiceHandler struct{}

func (UnimplementedProblemServiceHandler) ListProblems(context.Context, *connect.Request[v1.ListProblemsRequest]) (*connect.Response[v1.ListProblemsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("contestant.v1.ProblemService.ListProblems is not implemented"))
}

func (UnimplementedProblemServiceHandler) GetProblem(context.Context, *connect.Request[v1.GetProblemRequest]) (*connect.Response[v1.GetProblemResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("contestant.v1.ProblemService.GetProblem is not implemented"))
}

func (UnimplementedProblemServiceHandler) ListDeployments(context.Context, *connect.Request[v1.ListDeploymentsRequest]) (*connect.Response[v1.ListDeploymentsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("contestant.v1.ProblemService.ListDeployments is not implemented"))
}

func (UnimplementedProblemServiceHandler) Deploy(context.Context, *connect.Request[v1.DeployRequest]) (*connect.Response[v1.DeployResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("contestant.v1.ProblemService.Deploy is not implemented"))
}
