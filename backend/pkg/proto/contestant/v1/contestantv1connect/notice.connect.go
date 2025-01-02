// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: contestant/v1/notice.proto

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
	// NoticeServiceName is the fully-qualified name of the NoticeService service.
	NoticeServiceName = "contestant.v1.NoticeService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// NoticeServiceListNoticesProcedure is the fully-qualified name of the NoticeService's ListNotices
	// RPC.
	NoticeServiceListNoticesProcedure = "/contestant.v1.NoticeService/ListNotices"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	noticeServiceServiceDescriptor           = v1.File_contestant_v1_notice_proto.Services().ByName("NoticeService")
	noticeServiceListNoticesMethodDescriptor = noticeServiceServiceDescriptor.Methods().ByName("ListNotices")
)

// NoticeServiceClient is a client for the contestant.v1.NoticeService service.
type NoticeServiceClient interface {
	ListNotices(context.Context, *connect.Request[v1.ListNoticesRequest]) (*connect.Response[v1.ListNoticesResponse], error)
}

// NewNoticeServiceClient constructs a client for the contestant.v1.NoticeService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewNoticeServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) NoticeServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &noticeServiceClient{
		listNotices: connect.NewClient[v1.ListNoticesRequest, v1.ListNoticesResponse](
			httpClient,
			baseURL+NoticeServiceListNoticesProcedure,
			connect.WithSchema(noticeServiceListNoticesMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// noticeServiceClient implements NoticeServiceClient.
type noticeServiceClient struct {
	listNotices *connect.Client[v1.ListNoticesRequest, v1.ListNoticesResponse]
}

// ListNotices calls contestant.v1.NoticeService.ListNotices.
func (c *noticeServiceClient) ListNotices(ctx context.Context, req *connect.Request[v1.ListNoticesRequest]) (*connect.Response[v1.ListNoticesResponse], error) {
	return c.listNotices.CallUnary(ctx, req)
}

// NoticeServiceHandler is an implementation of the contestant.v1.NoticeService service.
type NoticeServiceHandler interface {
	ListNotices(context.Context, *connect.Request[v1.ListNoticesRequest]) (*connect.Response[v1.ListNoticesResponse], error)
}

// NewNoticeServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewNoticeServiceHandler(svc NoticeServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	noticeServiceListNoticesHandler := connect.NewUnaryHandler(
		NoticeServiceListNoticesProcedure,
		svc.ListNotices,
		connect.WithSchema(noticeServiceListNoticesMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/contestant.v1.NoticeService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case NoticeServiceListNoticesProcedure:
			noticeServiceListNoticesHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedNoticeServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedNoticeServiceHandler struct{}

func (UnimplementedNoticeServiceHandler) ListNotices(context.Context, *connect.Request[v1.ListNoticesRequest]) (*connect.Response[v1.ListNoticesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("contestant.v1.NoticeService.ListNotices is not implemented"))
}