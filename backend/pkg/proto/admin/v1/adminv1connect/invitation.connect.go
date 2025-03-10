// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: admin/v1/invitation.proto

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
	// InvitationServiceName is the fully-qualified name of the InvitationService service.
	InvitationServiceName = "admin.v1.InvitationService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// InvitationServiceListInvitationCodesProcedure is the fully-qualified name of the
	// InvitationService's ListInvitationCodes RPC.
	InvitationServiceListInvitationCodesProcedure = "/admin.v1.InvitationService/ListInvitationCodes"
	// InvitationServiceCreateInvitationCodeProcedure is the fully-qualified name of the
	// InvitationService's CreateInvitationCode RPC.
	InvitationServiceCreateInvitationCodeProcedure = "/admin.v1.InvitationService/CreateInvitationCode"
)

// InvitationServiceClient is a client for the admin.v1.InvitationService service.
type InvitationServiceClient interface {
	ListInvitationCodes(context.Context, *connect.Request[v1.ListInvitationCodesRequest]) (*connect.Response[v1.ListInvitationCodesResponse], error)
	CreateInvitationCode(context.Context, *connect.Request[v1.CreateInvitationCodeRequest]) (*connect.Response[v1.CreateInvitationCodeResponse], error)
}

// NewInvitationServiceClient constructs a client for the admin.v1.InvitationService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewInvitationServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) InvitationServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	invitationServiceMethods := v1.File_admin_v1_invitation_proto.Services().ByName("InvitationService").Methods()
	return &invitationServiceClient{
		listInvitationCodes: connect.NewClient[v1.ListInvitationCodesRequest, v1.ListInvitationCodesResponse](
			httpClient,
			baseURL+InvitationServiceListInvitationCodesProcedure,
			connect.WithSchema(invitationServiceMethods.ByName("ListInvitationCodes")),
			connect.WithClientOptions(opts...),
		),
		createInvitationCode: connect.NewClient[v1.CreateInvitationCodeRequest, v1.CreateInvitationCodeResponse](
			httpClient,
			baseURL+InvitationServiceCreateInvitationCodeProcedure,
			connect.WithSchema(invitationServiceMethods.ByName("CreateInvitationCode")),
			connect.WithClientOptions(opts...),
		),
	}
}

// invitationServiceClient implements InvitationServiceClient.
type invitationServiceClient struct {
	listInvitationCodes  *connect.Client[v1.ListInvitationCodesRequest, v1.ListInvitationCodesResponse]
	createInvitationCode *connect.Client[v1.CreateInvitationCodeRequest, v1.CreateInvitationCodeResponse]
}

// ListInvitationCodes calls admin.v1.InvitationService.ListInvitationCodes.
func (c *invitationServiceClient) ListInvitationCodes(ctx context.Context, req *connect.Request[v1.ListInvitationCodesRequest]) (*connect.Response[v1.ListInvitationCodesResponse], error) {
	return c.listInvitationCodes.CallUnary(ctx, req)
}

// CreateInvitationCode calls admin.v1.InvitationService.CreateInvitationCode.
func (c *invitationServiceClient) CreateInvitationCode(ctx context.Context, req *connect.Request[v1.CreateInvitationCodeRequest]) (*connect.Response[v1.CreateInvitationCodeResponse], error) {
	return c.createInvitationCode.CallUnary(ctx, req)
}

// InvitationServiceHandler is an implementation of the admin.v1.InvitationService service.
type InvitationServiceHandler interface {
	ListInvitationCodes(context.Context, *connect.Request[v1.ListInvitationCodesRequest]) (*connect.Response[v1.ListInvitationCodesResponse], error)
	CreateInvitationCode(context.Context, *connect.Request[v1.CreateInvitationCodeRequest]) (*connect.Response[v1.CreateInvitationCodeResponse], error)
}

// NewInvitationServiceHandler builds an HTTP handler from the service implementation. It returns
// the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewInvitationServiceHandler(svc InvitationServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	invitationServiceMethods := v1.File_admin_v1_invitation_proto.Services().ByName("InvitationService").Methods()
	invitationServiceListInvitationCodesHandler := connect.NewUnaryHandler(
		InvitationServiceListInvitationCodesProcedure,
		svc.ListInvitationCodes,
		connect.WithSchema(invitationServiceMethods.ByName("ListInvitationCodes")),
		connect.WithHandlerOptions(opts...),
	)
	invitationServiceCreateInvitationCodeHandler := connect.NewUnaryHandler(
		InvitationServiceCreateInvitationCodeProcedure,
		svc.CreateInvitationCode,
		connect.WithSchema(invitationServiceMethods.ByName("CreateInvitationCode")),
		connect.WithHandlerOptions(opts...),
	)
	return "/admin.v1.InvitationService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case InvitationServiceListInvitationCodesProcedure:
			invitationServiceListInvitationCodesHandler.ServeHTTP(w, r)
		case InvitationServiceCreateInvitationCodeProcedure:
			invitationServiceCreateInvitationCodeHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedInvitationServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedInvitationServiceHandler struct{}

func (UnimplementedInvitationServiceHandler) ListInvitationCodes(context.Context, *connect.Request[v1.ListInvitationCodesRequest]) (*connect.Response[v1.ListInvitationCodesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("admin.v1.InvitationService.ListInvitationCodes is not implemented"))
}

func (UnimplementedInvitationServiceHandler) CreateInvitationCode(context.Context, *connect.Request[v1.CreateInvitationCodeRequest]) (*connect.Response[v1.CreateInvitationCodeResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("admin.v1.InvitationService.CreateInvitationCode is not implemented"))
}
