package grpc_server

import (
	"context"
	"log/slog"

	ragv1 "github.com/aria3ppp/rag-server/gen/go/rag/v1"
	"github.com/aria3ppp/rag-server/internal/rag/domain"
	"github.com/aria3ppp/rag-server/internal/rag/usecase"
	"github.com/samber/lo"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type ragGRPCServer struct {
	ragv1.UnimplementedRAGServiceServer
	uc     usecase.UseCase
	tracer trace.Tracer
	logger *slog.Logger
}

var _ ragv1.RAGServiceServer = (*ragGRPCServer)(nil)

func NewGRPCServer(
	uc usecase.UseCase,
	tracer trace.Tracer,
	logger *slog.Logger,
) *ragGRPCServer {
	return &ragGRPCServer{
		uc:     uc,
		tracer: tracer,
		logger: logger,
	}
}

func (grpcServer *ragGRPCServer) Query(ctx context.Context, request *ragv1.RAGServiceQueryRequest) (_ *ragv1.RAGServiceQueryResponse, err error) {
	ctx, span := grpcServer.tracer.Start(ctx, "grpcServer.Query")
	defer func() {
		defer span.End()
		if err != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	messages := lo.Map(request.GetMessages(), func(m *ragv1.Message, _ int) *domain.Message {
		return &domain.Message{
			Role:    domain.Role(m.GetRole()),
			Content: m.GetContent(),
		}
	})

	input := &domain.QueryInput{
		Query:    request.GetQuery(),
		Messages: messages,
	}

	result, err := grpcServer.uc.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	response := &ragv1.RAGServiceQueryResponse{
		Content:     result.Content,
		CreatedInMs: result.CreatedInMS,
	}

	return response, nil
}

func (grpcServer *ragGRPCServer) QueryStream(request *ragv1.RAGServiceQueryStreamRequest, stream grpc.ServerStreamingServer[ragv1.RAGServiceQueryStreamResponse]) (err error) {
	ctx := stream.Context()

	ctx, span := grpcServer.tracer.Start(ctx, "grpcServer.QueryStream")
	defer func() {
		defer span.End()
		if err != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	messages := lo.Map(request.GetMessages(), func(m *ragv1.Message, _ int) *domain.Message {
		return &domain.Message{
			Role:    domain.Role(m.GetRole()),
			Content: m.GetContent(),
		}
	})

	input := &domain.QueryStreamInput{
		Query:    request.GetQuery(),
		Messages: messages,
	}

	grpcServer.uc.QueryStream(ctx, input, func(event *domain.QueryStreamResultEvent) (continueRunning bool) {
		err = event.Error

		var responseError string
		if err != nil {
			responseError = err.Error()
		}

		item := &ragv1.RAGServiceQueryStreamResponse{
			Content:     event.Content,
			CreatedAtMs: event.CreatedAtMS,
			StopReason:  ragv1.StopReason(event.StopReason),
			Error:       responseError,
		}

		if err = stream.Send(item); err != nil {
			return false
		}

		return true
	})

	return nil
}
