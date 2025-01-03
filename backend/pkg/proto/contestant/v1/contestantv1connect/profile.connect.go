// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: contestant/v1/profile.proto

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
	// ProfileServiceName is the fully-qualified name of the ProfileService service.
	ProfileServiceName = "contestant.v1.ProfileService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// ProfileServiceListTeamsProcedure is the fully-qualified name of the ProfileService's ListTeams
	// RPC.
	ProfileServiceListTeamsProcedure = "/contestant.v1.ProfileService/ListTeams"
	// ProfileServiceUpdateProfileProcedure is the fully-qualified name of the ProfileService's
	// UpdateProfile RPC.
	ProfileServiceUpdateProfileProcedure = "/contestant.v1.ProfileService/UpdateProfile"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	profileServiceServiceDescriptor             = v1.File_contestant_v1_profile_proto.Services().ByName("ProfileService")
	profileServiceListTeamsMethodDescriptor     = profileServiceServiceDescriptor.Methods().ByName("ListTeams")
	profileServiceUpdateProfileMethodDescriptor = profileServiceServiceDescriptor.Methods().ByName("UpdateProfile")
)

// ProfileServiceClient is a client for the contestant.v1.ProfileService service.
type ProfileServiceClient interface {
	ListTeams(context.Context, *connect.Request[v1.ListTeamsRequest]) (*connect.Response[v1.ListTeamsResponse], error)
	UpdateProfile(context.Context, *connect.Request[v1.UpdateProfileRequest]) (*connect.Response[v1.UpdateProfileResponse], error)
}

// NewProfileServiceClient constructs a client for the contestant.v1.ProfileService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewProfileServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) ProfileServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &profileServiceClient{
		listTeams: connect.NewClient[v1.ListTeamsRequest, v1.ListTeamsResponse](
			httpClient,
			baseURL+ProfileServiceListTeamsProcedure,
			connect.WithSchema(profileServiceListTeamsMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		updateProfile: connect.NewClient[v1.UpdateProfileRequest, v1.UpdateProfileResponse](
			httpClient,
			baseURL+ProfileServiceUpdateProfileProcedure,
			connect.WithSchema(profileServiceUpdateProfileMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// profileServiceClient implements ProfileServiceClient.
type profileServiceClient struct {
	listTeams     *connect.Client[v1.ListTeamsRequest, v1.ListTeamsResponse]
	updateProfile *connect.Client[v1.UpdateProfileRequest, v1.UpdateProfileResponse]
}

// ListTeams calls contestant.v1.ProfileService.ListTeams.
func (c *profileServiceClient) ListTeams(ctx context.Context, req *connect.Request[v1.ListTeamsRequest]) (*connect.Response[v1.ListTeamsResponse], error) {
	return c.listTeams.CallUnary(ctx, req)
}

// UpdateProfile calls contestant.v1.ProfileService.UpdateProfile.
func (c *profileServiceClient) UpdateProfile(ctx context.Context, req *connect.Request[v1.UpdateProfileRequest]) (*connect.Response[v1.UpdateProfileResponse], error) {
	return c.updateProfile.CallUnary(ctx, req)
}

// ProfileServiceHandler is an implementation of the contestant.v1.ProfileService service.
type ProfileServiceHandler interface {
	ListTeams(context.Context, *connect.Request[v1.ListTeamsRequest]) (*connect.Response[v1.ListTeamsResponse], error)
	UpdateProfile(context.Context, *connect.Request[v1.UpdateProfileRequest]) (*connect.Response[v1.UpdateProfileResponse], error)
}

// NewProfileServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewProfileServiceHandler(svc ProfileServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	profileServiceListTeamsHandler := connect.NewUnaryHandler(
		ProfileServiceListTeamsProcedure,
		svc.ListTeams,
		connect.WithSchema(profileServiceListTeamsMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	profileServiceUpdateProfileHandler := connect.NewUnaryHandler(
		ProfileServiceUpdateProfileProcedure,
		svc.UpdateProfile,
		connect.WithSchema(profileServiceUpdateProfileMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/contestant.v1.ProfileService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ProfileServiceListTeamsProcedure:
			profileServiceListTeamsHandler.ServeHTTP(w, r)
		case ProfileServiceUpdateProfileProcedure:
			profileServiceUpdateProfileHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedProfileServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedProfileServiceHandler struct{}

func (UnimplementedProfileServiceHandler) ListTeams(context.Context, *connect.Request[v1.ListTeamsRequest]) (*connect.Response[v1.ListTeamsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("contestant.v1.ProfileService.ListTeams is not implemented"))
}

func (UnimplementedProfileServiceHandler) UpdateProfile(context.Context, *connect.Request[v1.UpdateProfileRequest]) (*connect.Response[v1.UpdateProfileResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("contestant.v1.ProfileService.UpdateProfile is not implemented"))
}