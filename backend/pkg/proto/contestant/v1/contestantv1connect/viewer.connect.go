// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: contestant/v1/viewer.proto

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
	// ViewerServiceName is the fully-qualified name of the ViewerService service.
	ViewerServiceName = "contestant.v1.ViewerService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// ViewerServiceGetViewerProcedure is the fully-qualified name of the ViewerService's GetViewer RPC.
	ViewerServiceGetViewerProcedure = "/contestant.v1.ViewerService/GetViewer"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	viewerServiceServiceDescriptor         = v1.File_contestant_v1_viewer_proto.Services().ByName("ViewerService")
	viewerServiceGetViewerMethodDescriptor = viewerServiceServiceDescriptor.Methods().ByName("GetViewer")
)

// ViewerServiceClient is a client for the contestant.v1.ViewerService service.
type ViewerServiceClient interface {
	GetViewer(context.Context, *connect.Request[v1.GetViewerRequest]) (*connect.Response[v1.GetViewerResponse], error)
}

// NewViewerServiceClient constructs a client for the contestant.v1.ViewerService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewViewerServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) ViewerServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &viewerServiceClient{
		getViewer: connect.NewClient[v1.GetViewerRequest, v1.GetViewerResponse](
			httpClient,
			baseURL+ViewerServiceGetViewerProcedure,
			connect.WithSchema(viewerServiceGetViewerMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// viewerServiceClient implements ViewerServiceClient.
type viewerServiceClient struct {
	getViewer *connect.Client[v1.GetViewerRequest, v1.GetViewerResponse]
}

// GetViewer calls contestant.v1.ViewerService.GetViewer.
func (c *viewerServiceClient) GetViewer(ctx context.Context, req *connect.Request[v1.GetViewerRequest]) (*connect.Response[v1.GetViewerResponse], error) {
	return c.getViewer.CallUnary(ctx, req)
}

// ViewerServiceHandler is an implementation of the contestant.v1.ViewerService service.
type ViewerServiceHandler interface {
	GetViewer(context.Context, *connect.Request[v1.GetViewerRequest]) (*connect.Response[v1.GetViewerResponse], error)
}

// NewViewerServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewViewerServiceHandler(svc ViewerServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	viewerServiceGetViewerHandler := connect.NewUnaryHandler(
		ViewerServiceGetViewerProcedure,
		svc.GetViewer,
		connect.WithSchema(viewerServiceGetViewerMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/contestant.v1.ViewerService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ViewerServiceGetViewerProcedure:
			viewerServiceGetViewerHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedViewerServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedViewerServiceHandler struct{}

func (UnimplementedViewerServiceHandler) GetViewer(context.Context, *connect.Request[v1.GetViewerRequest]) (*connect.Response[v1.GetViewerResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("contestant.v1.ViewerService.GetViewer is not implemented"))
}
