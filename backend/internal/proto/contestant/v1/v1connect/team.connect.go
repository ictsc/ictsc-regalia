// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: contestant/v1/team.proto

package v1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "github.com/ictsc/ictsc-outlands/backend/internal/proto/contestant/v1"
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
	// TeamServiceName is the fully-qualified name of the TeamService service.
	TeamServiceName = "contestant.v1.TeamService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// TeamServiceGetTeamsProcedure is the fully-qualified name of the TeamService's GetTeams RPC.
	TeamServiceGetTeamsProcedure = "/contestant.v1.TeamService/GetTeams"
	// TeamServiceGetTeamProcedure is the fully-qualified name of the TeamService's GetTeam RPC.
	TeamServiceGetTeamProcedure = "/contestant.v1.TeamService/GetTeam"
	// TeamServiceGetMembersProcedure is the fully-qualified name of the TeamService's GetMembers RPC.
	TeamServiceGetMembersProcedure = "/contestant.v1.TeamService/GetMembers"
	// TeamServiceGetConnectionInfoProcedure is the fully-qualified name of the TeamService's
	// GetConnectionInfo RPC.
	TeamServiceGetConnectionInfoProcedure = "/contestant.v1.TeamService/GetConnectionInfo"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	teamServiceServiceDescriptor                 = v1.File_contestant_v1_team_proto.Services().ByName("TeamService")
	teamServiceGetTeamsMethodDescriptor          = teamServiceServiceDescriptor.Methods().ByName("GetTeams")
	teamServiceGetTeamMethodDescriptor           = teamServiceServiceDescriptor.Methods().ByName("GetTeam")
	teamServiceGetMembersMethodDescriptor        = teamServiceServiceDescriptor.Methods().ByName("GetMembers")
	teamServiceGetConnectionInfoMethodDescriptor = teamServiceServiceDescriptor.Methods().ByName("GetConnectionInfo")
)

// TeamServiceClient is a client for the contestant.v1.TeamService service.
type TeamServiceClient interface {
	GetTeams(context.Context, *connect.Request[v1.GetTeamsRequest]) (*connect.Response[v1.GetTeamsResponse], error)
	GetTeam(context.Context, *connect.Request[v1.GetTeamRequest]) (*connect.Response[v1.GetTeamResponse], error)
	GetMembers(context.Context, *connect.Request[v1.GetMembersRequest]) (*connect.Response[v1.GetMembersResponse], error)
	GetConnectionInfo(context.Context, *connect.Request[v1.GetConnectionInfoRequest]) (*connect.Response[v1.GetConnectionInfoResponse], error)
}

// NewTeamServiceClient constructs a client for the contestant.v1.TeamService service. By default,
// it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and
// sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC()
// or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewTeamServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) TeamServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &teamServiceClient{
		getTeams: connect.NewClient[v1.GetTeamsRequest, v1.GetTeamsResponse](
			httpClient,
			baseURL+TeamServiceGetTeamsProcedure,
			connect.WithSchema(teamServiceGetTeamsMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		getTeam: connect.NewClient[v1.GetTeamRequest, v1.GetTeamResponse](
			httpClient,
			baseURL+TeamServiceGetTeamProcedure,
			connect.WithSchema(teamServiceGetTeamMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		getMembers: connect.NewClient[v1.GetMembersRequest, v1.GetMembersResponse](
			httpClient,
			baseURL+TeamServiceGetMembersProcedure,
			connect.WithSchema(teamServiceGetMembersMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		getConnectionInfo: connect.NewClient[v1.GetConnectionInfoRequest, v1.GetConnectionInfoResponse](
			httpClient,
			baseURL+TeamServiceGetConnectionInfoProcedure,
			connect.WithSchema(teamServiceGetConnectionInfoMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// teamServiceClient implements TeamServiceClient.
type teamServiceClient struct {
	getTeams          *connect.Client[v1.GetTeamsRequest, v1.GetTeamsResponse]
	getTeam           *connect.Client[v1.GetTeamRequest, v1.GetTeamResponse]
	getMembers        *connect.Client[v1.GetMembersRequest, v1.GetMembersResponse]
	getConnectionInfo *connect.Client[v1.GetConnectionInfoRequest, v1.GetConnectionInfoResponse]
}

// GetTeams calls contestant.v1.TeamService.GetTeams.
func (c *teamServiceClient) GetTeams(ctx context.Context, req *connect.Request[v1.GetTeamsRequest]) (*connect.Response[v1.GetTeamsResponse], error) {
	return c.getTeams.CallUnary(ctx, req)
}

// GetTeam calls contestant.v1.TeamService.GetTeam.
func (c *teamServiceClient) GetTeam(ctx context.Context, req *connect.Request[v1.GetTeamRequest]) (*connect.Response[v1.GetTeamResponse], error) {
	return c.getTeam.CallUnary(ctx, req)
}

// GetMembers calls contestant.v1.TeamService.GetMembers.
func (c *teamServiceClient) GetMembers(ctx context.Context, req *connect.Request[v1.GetMembersRequest]) (*connect.Response[v1.GetMembersResponse], error) {
	return c.getMembers.CallUnary(ctx, req)
}

// GetConnectionInfo calls contestant.v1.TeamService.GetConnectionInfo.
func (c *teamServiceClient) GetConnectionInfo(ctx context.Context, req *connect.Request[v1.GetConnectionInfoRequest]) (*connect.Response[v1.GetConnectionInfoResponse], error) {
	return c.getConnectionInfo.CallUnary(ctx, req)
}

// TeamServiceHandler is an implementation of the contestant.v1.TeamService service.
type TeamServiceHandler interface {
	GetTeams(context.Context, *connect.Request[v1.GetTeamsRequest]) (*connect.Response[v1.GetTeamsResponse], error)
	GetTeam(context.Context, *connect.Request[v1.GetTeamRequest]) (*connect.Response[v1.GetTeamResponse], error)
	GetMembers(context.Context, *connect.Request[v1.GetMembersRequest]) (*connect.Response[v1.GetMembersResponse], error)
	GetConnectionInfo(context.Context, *connect.Request[v1.GetConnectionInfoRequest]) (*connect.Response[v1.GetConnectionInfoResponse], error)
}

// NewTeamServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewTeamServiceHandler(svc TeamServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	teamServiceGetTeamsHandler := connect.NewUnaryHandler(
		TeamServiceGetTeamsProcedure,
		svc.GetTeams,
		connect.WithSchema(teamServiceGetTeamsMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	teamServiceGetTeamHandler := connect.NewUnaryHandler(
		TeamServiceGetTeamProcedure,
		svc.GetTeam,
		connect.WithSchema(teamServiceGetTeamMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	teamServiceGetMembersHandler := connect.NewUnaryHandler(
		TeamServiceGetMembersProcedure,
		svc.GetMembers,
		connect.WithSchema(teamServiceGetMembersMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	teamServiceGetConnectionInfoHandler := connect.NewUnaryHandler(
		TeamServiceGetConnectionInfoProcedure,
		svc.GetConnectionInfo,
		connect.WithSchema(teamServiceGetConnectionInfoMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/contestant.v1.TeamService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case TeamServiceGetTeamsProcedure:
			teamServiceGetTeamsHandler.ServeHTTP(w, r)
		case TeamServiceGetTeamProcedure:
			teamServiceGetTeamHandler.ServeHTTP(w, r)
		case TeamServiceGetMembersProcedure:
			teamServiceGetMembersHandler.ServeHTTP(w, r)
		case TeamServiceGetConnectionInfoProcedure:
			teamServiceGetConnectionInfoHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedTeamServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedTeamServiceHandler struct{}

func (UnimplementedTeamServiceHandler) GetTeams(context.Context, *connect.Request[v1.GetTeamsRequest]) (*connect.Response[v1.GetTeamsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("contestant.v1.TeamService.GetTeams is not implemented"))
}

func (UnimplementedTeamServiceHandler) GetTeam(context.Context, *connect.Request[v1.GetTeamRequest]) (*connect.Response[v1.GetTeamResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("contestant.v1.TeamService.GetTeam is not implemented"))
}

func (UnimplementedTeamServiceHandler) GetMembers(context.Context, *connect.Request[v1.GetMembersRequest]) (*connect.Response[v1.GetMembersResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("contestant.v1.TeamService.GetMembers is not implemented"))
}

func (UnimplementedTeamServiceHandler) GetConnectionInfo(context.Context, *connect.Request[v1.GetConnectionInfoRequest]) (*connect.Response[v1.GetConnectionInfoResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("contestant.v1.TeamService.GetConnectionInfo is not implemented"))
}
