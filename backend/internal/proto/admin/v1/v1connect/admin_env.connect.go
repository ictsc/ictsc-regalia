// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: admin/v1/admin_env.proto

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
	// AdminEnvServiceName is the fully-qualified name of the AdminEnvService service.
	AdminEnvServiceName = "admin.v1.AdminEnvService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// AdminEnvServiceGetAdminConnectionInfoProcedure is the fully-qualified name of the
	// AdminEnvService's GetAdminConnectionInfo RPC.
	AdminEnvServiceGetAdminConnectionInfoProcedure = "/admin.v1.AdminEnvService/GetAdminConnectionInfo"
	// AdminEnvServicePutAdminConnectionInfoProcedure is the fully-qualified name of the
	// AdminEnvService's PutAdminConnectionInfo RPC.
	AdminEnvServicePutAdminConnectionInfoProcedure = "/admin.v1.AdminEnvService/PutAdminConnectionInfo"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	adminEnvServiceServiceDescriptor                      = v1.File_admin_v1_admin_env_proto.Services().ByName("AdminEnvService")
	adminEnvServiceGetAdminConnectionInfoMethodDescriptor = adminEnvServiceServiceDescriptor.Methods().ByName("GetAdminConnectionInfo")
	adminEnvServicePutAdminConnectionInfoMethodDescriptor = adminEnvServiceServiceDescriptor.Methods().ByName("PutAdminConnectionInfo")
)

// AdminEnvServiceClient is a client for the admin.v1.AdminEnvService service.
type AdminEnvServiceClient interface {
	GetAdminConnectionInfo(context.Context, *connect.Request[v1.GetAdminConnectionInfoRequest]) (*connect.Response[v1.GetAdminConnectionInfoResponse], error)
	PutAdminConnectionInfo(context.Context, *connect.Request[v1.PutAdminConnectionInfoRequest]) (*connect.Response[v1.PutAdminConnectionInfoResponse], error)
}

// NewAdminEnvServiceClient constructs a client for the admin.v1.AdminEnvService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewAdminEnvServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) AdminEnvServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &adminEnvServiceClient{
		getAdminConnectionInfo: connect.NewClient[v1.GetAdminConnectionInfoRequest, v1.GetAdminConnectionInfoResponse](
			httpClient,
			baseURL+AdminEnvServiceGetAdminConnectionInfoProcedure,
			connect.WithSchema(adminEnvServiceGetAdminConnectionInfoMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		putAdminConnectionInfo: connect.NewClient[v1.PutAdminConnectionInfoRequest, v1.PutAdminConnectionInfoResponse](
			httpClient,
			baseURL+AdminEnvServicePutAdminConnectionInfoProcedure,
			connect.WithSchema(adminEnvServicePutAdminConnectionInfoMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// adminEnvServiceClient implements AdminEnvServiceClient.
type adminEnvServiceClient struct {
	getAdminConnectionInfo *connect.Client[v1.GetAdminConnectionInfoRequest, v1.GetAdminConnectionInfoResponse]
	putAdminConnectionInfo *connect.Client[v1.PutAdminConnectionInfoRequest, v1.PutAdminConnectionInfoResponse]
}

// GetAdminConnectionInfo calls admin.v1.AdminEnvService.GetAdminConnectionInfo.
func (c *adminEnvServiceClient) GetAdminConnectionInfo(ctx context.Context, req *connect.Request[v1.GetAdminConnectionInfoRequest]) (*connect.Response[v1.GetAdminConnectionInfoResponse], error) {
	return c.getAdminConnectionInfo.CallUnary(ctx, req)
}

// PutAdminConnectionInfo calls admin.v1.AdminEnvService.PutAdminConnectionInfo.
func (c *adminEnvServiceClient) PutAdminConnectionInfo(ctx context.Context, req *connect.Request[v1.PutAdminConnectionInfoRequest]) (*connect.Response[v1.PutAdminConnectionInfoResponse], error) {
	return c.putAdminConnectionInfo.CallUnary(ctx, req)
}

// AdminEnvServiceHandler is an implementation of the admin.v1.AdminEnvService service.
type AdminEnvServiceHandler interface {
	GetAdminConnectionInfo(context.Context, *connect.Request[v1.GetAdminConnectionInfoRequest]) (*connect.Response[v1.GetAdminConnectionInfoResponse], error)
	PutAdminConnectionInfo(context.Context, *connect.Request[v1.PutAdminConnectionInfoRequest]) (*connect.Response[v1.PutAdminConnectionInfoResponse], error)
}

// NewAdminEnvServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewAdminEnvServiceHandler(svc AdminEnvServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	adminEnvServiceGetAdminConnectionInfoHandler := connect.NewUnaryHandler(
		AdminEnvServiceGetAdminConnectionInfoProcedure,
		svc.GetAdminConnectionInfo,
		connect.WithSchema(adminEnvServiceGetAdminConnectionInfoMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	adminEnvServicePutAdminConnectionInfoHandler := connect.NewUnaryHandler(
		AdminEnvServicePutAdminConnectionInfoProcedure,
		svc.PutAdminConnectionInfo,
		connect.WithSchema(adminEnvServicePutAdminConnectionInfoMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/admin.v1.AdminEnvService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case AdminEnvServiceGetAdminConnectionInfoProcedure:
			adminEnvServiceGetAdminConnectionInfoHandler.ServeHTTP(w, r)
		case AdminEnvServicePutAdminConnectionInfoProcedure:
			adminEnvServicePutAdminConnectionInfoHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedAdminEnvServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedAdminEnvServiceHandler struct{}

func (UnimplementedAdminEnvServiceHandler) GetAdminConnectionInfo(context.Context, *connect.Request[v1.GetAdminConnectionInfoRequest]) (*connect.Response[v1.GetAdminConnectionInfoResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("admin.v1.AdminEnvService.GetAdminConnectionInfo is not implemented"))
}

func (UnimplementedAdminEnvServiceHandler) PutAdminConnectionInfo(context.Context, *connect.Request[v1.PutAdminConnectionInfoRequest]) (*connect.Response[v1.PutAdminConnectionInfoResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("admin.v1.AdminEnvService.PutAdminConnectionInfo is not implemented"))
}
