package server

import (
	"context"

	ragv1 "github.com/aria3ppp/rag-server/gen/go/rag/v1"
	"github.com/aria3ppp/rag-server/internal/rag/domain"
	"github.com/aria3ppp/rag-server/internal/rag/usecase"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ragGRPCService struct {
	ragv1.UnimplementedRAGServiceServer
	uc usecase.UseCase
}

var _ ragv1.RAGServiceServer = (*ragGRPCService)(nil)

func NewRagGRPCService(uc usecase.UseCase) *ragGRPCService {
	return &ragGRPCService{
		uc: uc,
	}
}

func (grpcService *ragGRPCService) QuerySync(ctx context.Context, request *ragv1.RAGServiceQuerySyncRequest) (*ragv1.RAGServiceQuerySyncResponse, error) {
	result, err := grpcService.uc.QuerySync(ctx, request.Query)
	if err != nil {
		return nil, err
	}

	response := &ragv1.RAGServiceQuerySyncResponse{
		Content: result,
	}

	return response, nil
}

func (grpcService *ragGRPCService) QueryAsync(request *ragv1.RAGServiceQueryAsyncRequest, stream grpc.ServerStreamingServer[ragv1.RAGServiceQueryAsyncResponse]) error {
	if err := grpcService.uc.QueryAsync(stream.Context(), request.Query, func(ctx context.Context, event *domain.QueryAsyncEvent) error {
		item := &ragv1.RAGServiceQueryAsyncResponse{
			Done:      event.Done,
			Content:   event.Content,
			CreatedAt: timestamppb.New(event.CreatedAt),
			Error:     "",
		}

		if err := stream.Send(item); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
