// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: contestant/v1/ranking.proto

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
	// RankingServiceName is the fully-qualified name of the RankingService service.
	RankingServiceName = "contestant.v1.RankingService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// RankingServiceGetRankingProcedure is the fully-qualified name of the RankingService's GetRanking
	// RPC.
	RankingServiceGetRankingProcedure = "/contestant.v1.RankingService/GetRanking"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	rankingServiceServiceDescriptor          = v1.File_contestant_v1_ranking_proto.Services().ByName("RankingService")
	rankingServiceGetRankingMethodDescriptor = rankingServiceServiceDescriptor.Methods().ByName("GetRanking")
)

// RankingServiceClient is a client for the contestant.v1.RankingService service.
type RankingServiceClient interface {
	GetRanking(context.Context, *connect.Request[v1.GetRankingRequest]) (*connect.Response[v1.GetRankingResponse], error)
}

// NewRankingServiceClient constructs a client for the contestant.v1.RankingService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewRankingServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) RankingServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &rankingServiceClient{
		getRanking: connect.NewClient[v1.GetRankingRequest, v1.GetRankingResponse](
			httpClient,
			baseURL+RankingServiceGetRankingProcedure,
			connect.WithSchema(rankingServiceGetRankingMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// rankingServiceClient implements RankingServiceClient.
type rankingServiceClient struct {
	getRanking *connect.Client[v1.GetRankingRequest, v1.GetRankingResponse]
}

// GetRanking calls contestant.v1.RankingService.GetRanking.
func (c *rankingServiceClient) GetRanking(ctx context.Context, req *connect.Request[v1.GetRankingRequest]) (*connect.Response[v1.GetRankingResponse], error) {
	return c.getRanking.CallUnary(ctx, req)
}

// RankingServiceHandler is an implementation of the contestant.v1.RankingService service.
type RankingServiceHandler interface {
	GetRanking(context.Context, *connect.Request[v1.GetRankingRequest]) (*connect.Response[v1.GetRankingResponse], error)
}

// NewRankingServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewRankingServiceHandler(svc RankingServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	rankingServiceGetRankingHandler := connect.NewUnaryHandler(
		RankingServiceGetRankingProcedure,
		svc.GetRanking,
		connect.WithSchema(rankingServiceGetRankingMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/contestant.v1.RankingService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case RankingServiceGetRankingProcedure:
			rankingServiceGetRankingHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedRankingServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedRankingServiceHandler struct{}

func (UnimplementedRankingServiceHandler) GetRanking(context.Context, *connect.Request[v1.GetRankingRequest]) (*connect.Response[v1.GetRankingResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("contestant.v1.RankingService.GetRanking is not implemented"))
}
