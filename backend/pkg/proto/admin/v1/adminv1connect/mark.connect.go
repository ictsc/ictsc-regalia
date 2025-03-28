// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: admin/v1/mark.proto

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
	// MarkServiceName is the fully-qualified name of the MarkService service.
	MarkServiceName = "admin.v1.MarkService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// MarkServiceListAnswersProcedure is the fully-qualified name of the MarkService's ListAnswers RPC.
	MarkServiceListAnswersProcedure = "/admin.v1.MarkService/ListAnswers"
	// MarkServiceGetAnswerProcedure is the fully-qualified name of the MarkService's GetAnswer RPC.
	MarkServiceGetAnswerProcedure = "/admin.v1.MarkService/GetAnswer"
	// MarkServiceListMarkingResultsProcedure is the fully-qualified name of the MarkService's
	// ListMarkingResults RPC.
	MarkServiceListMarkingResultsProcedure = "/admin.v1.MarkService/ListMarkingResults"
	// MarkServiceCreateMarkingResultProcedure is the fully-qualified name of the MarkService's
	// CreateMarkingResult RPC.
	MarkServiceCreateMarkingResultProcedure = "/admin.v1.MarkService/CreateMarkingResult"
	// MarkServiceUpdateMarkingResultVisibilitiesProcedure is the fully-qualified name of the
	// MarkService's UpdateMarkingResultVisibilities RPC.
	MarkServiceUpdateMarkingResultVisibilitiesProcedure = "/admin.v1.MarkService/UpdateMarkingResultVisibilities"
	// MarkServiceUpdateScoresProcedure is the fully-qualified name of the MarkService's UpdateScores
	// RPC.
	MarkServiceUpdateScoresProcedure = "/admin.v1.MarkService/UpdateScores"
)

// MarkServiceClient is a client for the admin.v1.MarkService service.
type MarkServiceClient interface {
	ListAnswers(context.Context, *connect.Request[v1.ListAnswersRequest]) (*connect.Response[v1.ListAnswersResponse], error)
	GetAnswer(context.Context, *connect.Request[v1.GetAnswerRequest]) (*connect.Response[v1.GetAnswerResponse], error)
	ListMarkingResults(context.Context, *connect.Request[v1.ListMarkingResultsRequest]) (*connect.Response[v1.ListMarkingResultsResponse], error)
	CreateMarkingResult(context.Context, *connect.Request[v1.CreateMarkingResultRequest]) (*connect.Response[v1.CreateMarkingResultResponse], error)
	UpdateMarkingResultVisibilities(context.Context, *connect.Request[v1.UpdateMarkingResultVisibilitiesRequest]) (*connect.Response[v1.UpdateMarkingResultVisibilitiesResponse], error)
	UpdateScores(context.Context, *connect.Request[v1.UpdateScoresRequest]) (*connect.Response[v1.UpdateScoresResponse], error)
}

// NewMarkServiceClient constructs a client for the admin.v1.MarkService service. By default, it
// uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewMarkServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) MarkServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	markServiceMethods := v1.File_admin_v1_mark_proto.Services().ByName("MarkService").Methods()
	return &markServiceClient{
		listAnswers: connect.NewClient[v1.ListAnswersRequest, v1.ListAnswersResponse](
			httpClient,
			baseURL+MarkServiceListAnswersProcedure,
			connect.WithSchema(markServiceMethods.ByName("ListAnswers")),
			connect.WithClientOptions(opts...),
		),
		getAnswer: connect.NewClient[v1.GetAnswerRequest, v1.GetAnswerResponse](
			httpClient,
			baseURL+MarkServiceGetAnswerProcedure,
			connect.WithSchema(markServiceMethods.ByName("GetAnswer")),
			connect.WithClientOptions(opts...),
		),
		listMarkingResults: connect.NewClient[v1.ListMarkingResultsRequest, v1.ListMarkingResultsResponse](
			httpClient,
			baseURL+MarkServiceListMarkingResultsProcedure,
			connect.WithSchema(markServiceMethods.ByName("ListMarkingResults")),
			connect.WithClientOptions(opts...),
		),
		createMarkingResult: connect.NewClient[v1.CreateMarkingResultRequest, v1.CreateMarkingResultResponse](
			httpClient,
			baseURL+MarkServiceCreateMarkingResultProcedure,
			connect.WithSchema(markServiceMethods.ByName("CreateMarkingResult")),
			connect.WithClientOptions(opts...),
		),
		updateMarkingResultVisibilities: connect.NewClient[v1.UpdateMarkingResultVisibilitiesRequest, v1.UpdateMarkingResultVisibilitiesResponse](
			httpClient,
			baseURL+MarkServiceUpdateMarkingResultVisibilitiesProcedure,
			connect.WithSchema(markServiceMethods.ByName("UpdateMarkingResultVisibilities")),
			connect.WithClientOptions(opts...),
		),
		updateScores: connect.NewClient[v1.UpdateScoresRequest, v1.UpdateScoresResponse](
			httpClient,
			baseURL+MarkServiceUpdateScoresProcedure,
			connect.WithSchema(markServiceMethods.ByName("UpdateScores")),
			connect.WithClientOptions(opts...),
		),
	}
}

// markServiceClient implements MarkServiceClient.
type markServiceClient struct {
	listAnswers                     *connect.Client[v1.ListAnswersRequest, v1.ListAnswersResponse]
	getAnswer                       *connect.Client[v1.GetAnswerRequest, v1.GetAnswerResponse]
	listMarkingResults              *connect.Client[v1.ListMarkingResultsRequest, v1.ListMarkingResultsResponse]
	createMarkingResult             *connect.Client[v1.CreateMarkingResultRequest, v1.CreateMarkingResultResponse]
	updateMarkingResultVisibilities *connect.Client[v1.UpdateMarkingResultVisibilitiesRequest, v1.UpdateMarkingResultVisibilitiesResponse]
	updateScores                    *connect.Client[v1.UpdateScoresRequest, v1.UpdateScoresResponse]
}

// ListAnswers calls admin.v1.MarkService.ListAnswers.
func (c *markServiceClient) ListAnswers(ctx context.Context, req *connect.Request[v1.ListAnswersRequest]) (*connect.Response[v1.ListAnswersResponse], error) {
	return c.listAnswers.CallUnary(ctx, req)
}

// GetAnswer calls admin.v1.MarkService.GetAnswer.
func (c *markServiceClient) GetAnswer(ctx context.Context, req *connect.Request[v1.GetAnswerRequest]) (*connect.Response[v1.GetAnswerResponse], error) {
	return c.getAnswer.CallUnary(ctx, req)
}

// ListMarkingResults calls admin.v1.MarkService.ListMarkingResults.
func (c *markServiceClient) ListMarkingResults(ctx context.Context, req *connect.Request[v1.ListMarkingResultsRequest]) (*connect.Response[v1.ListMarkingResultsResponse], error) {
	return c.listMarkingResults.CallUnary(ctx, req)
}

// CreateMarkingResult calls admin.v1.MarkService.CreateMarkingResult.
func (c *markServiceClient) CreateMarkingResult(ctx context.Context, req *connect.Request[v1.CreateMarkingResultRequest]) (*connect.Response[v1.CreateMarkingResultResponse], error) {
	return c.createMarkingResult.CallUnary(ctx, req)
}

// UpdateMarkingResultVisibilities calls admin.v1.MarkService.UpdateMarkingResultVisibilities.
func (c *markServiceClient) UpdateMarkingResultVisibilities(ctx context.Context, req *connect.Request[v1.UpdateMarkingResultVisibilitiesRequest]) (*connect.Response[v1.UpdateMarkingResultVisibilitiesResponse], error) {
	return c.updateMarkingResultVisibilities.CallUnary(ctx, req)
}

// UpdateScores calls admin.v1.MarkService.UpdateScores.
func (c *markServiceClient) UpdateScores(ctx context.Context, req *connect.Request[v1.UpdateScoresRequest]) (*connect.Response[v1.UpdateScoresResponse], error) {
	return c.updateScores.CallUnary(ctx, req)
}

// MarkServiceHandler is an implementation of the admin.v1.MarkService service.
type MarkServiceHandler interface {
	ListAnswers(context.Context, *connect.Request[v1.ListAnswersRequest]) (*connect.Response[v1.ListAnswersResponse], error)
	GetAnswer(context.Context, *connect.Request[v1.GetAnswerRequest]) (*connect.Response[v1.GetAnswerResponse], error)
	ListMarkingResults(context.Context, *connect.Request[v1.ListMarkingResultsRequest]) (*connect.Response[v1.ListMarkingResultsResponse], error)
	CreateMarkingResult(context.Context, *connect.Request[v1.CreateMarkingResultRequest]) (*connect.Response[v1.CreateMarkingResultResponse], error)
	UpdateMarkingResultVisibilities(context.Context, *connect.Request[v1.UpdateMarkingResultVisibilitiesRequest]) (*connect.Response[v1.UpdateMarkingResultVisibilitiesResponse], error)
	UpdateScores(context.Context, *connect.Request[v1.UpdateScoresRequest]) (*connect.Response[v1.UpdateScoresResponse], error)
}

// NewMarkServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewMarkServiceHandler(svc MarkServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	markServiceMethods := v1.File_admin_v1_mark_proto.Services().ByName("MarkService").Methods()
	markServiceListAnswersHandler := connect.NewUnaryHandler(
		MarkServiceListAnswersProcedure,
		svc.ListAnswers,
		connect.WithSchema(markServiceMethods.ByName("ListAnswers")),
		connect.WithHandlerOptions(opts...),
	)
	markServiceGetAnswerHandler := connect.NewUnaryHandler(
		MarkServiceGetAnswerProcedure,
		svc.GetAnswer,
		connect.WithSchema(markServiceMethods.ByName("GetAnswer")),
		connect.WithHandlerOptions(opts...),
	)
	markServiceListMarkingResultsHandler := connect.NewUnaryHandler(
		MarkServiceListMarkingResultsProcedure,
		svc.ListMarkingResults,
		connect.WithSchema(markServiceMethods.ByName("ListMarkingResults")),
		connect.WithHandlerOptions(opts...),
	)
	markServiceCreateMarkingResultHandler := connect.NewUnaryHandler(
		MarkServiceCreateMarkingResultProcedure,
		svc.CreateMarkingResult,
		connect.WithSchema(markServiceMethods.ByName("CreateMarkingResult")),
		connect.WithHandlerOptions(opts...),
	)
	markServiceUpdateMarkingResultVisibilitiesHandler := connect.NewUnaryHandler(
		MarkServiceUpdateMarkingResultVisibilitiesProcedure,
		svc.UpdateMarkingResultVisibilities,
		connect.WithSchema(markServiceMethods.ByName("UpdateMarkingResultVisibilities")),
		connect.WithHandlerOptions(opts...),
	)
	markServiceUpdateScoresHandler := connect.NewUnaryHandler(
		MarkServiceUpdateScoresProcedure,
		svc.UpdateScores,
		connect.WithSchema(markServiceMethods.ByName("UpdateScores")),
		connect.WithHandlerOptions(opts...),
	)
	return "/admin.v1.MarkService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case MarkServiceListAnswersProcedure:
			markServiceListAnswersHandler.ServeHTTP(w, r)
		case MarkServiceGetAnswerProcedure:
			markServiceGetAnswerHandler.ServeHTTP(w, r)
		case MarkServiceListMarkingResultsProcedure:
			markServiceListMarkingResultsHandler.ServeHTTP(w, r)
		case MarkServiceCreateMarkingResultProcedure:
			markServiceCreateMarkingResultHandler.ServeHTTP(w, r)
		case MarkServiceUpdateMarkingResultVisibilitiesProcedure:
			markServiceUpdateMarkingResultVisibilitiesHandler.ServeHTTP(w, r)
		case MarkServiceUpdateScoresProcedure:
			markServiceUpdateScoresHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedMarkServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedMarkServiceHandler struct{}

func (UnimplementedMarkServiceHandler) ListAnswers(context.Context, *connect.Request[v1.ListAnswersRequest]) (*connect.Response[v1.ListAnswersResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("admin.v1.MarkService.ListAnswers is not implemented"))
}

func (UnimplementedMarkServiceHandler) GetAnswer(context.Context, *connect.Request[v1.GetAnswerRequest]) (*connect.Response[v1.GetAnswerResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("admin.v1.MarkService.GetAnswer is not implemented"))
}

func (UnimplementedMarkServiceHandler) ListMarkingResults(context.Context, *connect.Request[v1.ListMarkingResultsRequest]) (*connect.Response[v1.ListMarkingResultsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("admin.v1.MarkService.ListMarkingResults is not implemented"))
}

func (UnimplementedMarkServiceHandler) CreateMarkingResult(context.Context, *connect.Request[v1.CreateMarkingResultRequest]) (*connect.Response[v1.CreateMarkingResultResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("admin.v1.MarkService.CreateMarkingResult is not implemented"))
}

func (UnimplementedMarkServiceHandler) UpdateMarkingResultVisibilities(context.Context, *connect.Request[v1.UpdateMarkingResultVisibilitiesRequest]) (*connect.Response[v1.UpdateMarkingResultVisibilitiesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("admin.v1.MarkService.UpdateMarkingResultVisibilities is not implemented"))
}

func (UnimplementedMarkServiceHandler) UpdateScores(context.Context, *connect.Request[v1.UpdateScoresRequest]) (*connect.Response[v1.UpdateScoresResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("admin.v1.MarkService.UpdateScores is not implemented"))
}
