package router

import (
	"context"

	vectorstorev1 "github.com/aria3ppp/rag-server/gen/go/vectorstore/v1"
	"github.com/aria3ppp/rag-server/internal/vectorstore/domain"
	vectorstore_error "github.com/aria3ppp/rag-server/internal/vectorstore/error"
	"github.com/aria3ppp/rag-server/internal/vectorstore/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/structpb"
)

type vectorStoreGRPCService struct {
	vectorstorev1.UnimplementedVectorStoreServiceServer
	uc usecase.UseCase
}

var _ vectorstorev1.VectorStoreServiceServer = (*vectorStoreGRPCService)(nil)

func NewVectorStoreGRPCService(uc usecase.UseCase) *vectorStoreGRPCService {
	return &vectorStoreGRPCService{
		uc: uc,
	}
}

func (grpcService *vectorStoreGRPCService) InsertTexts(ctx context.Context, req *vectorstorev1.VectorStoreServiceInsertTextsRequest) (*vectorstorev1.VectorStoreServiceInsertTextsResponse, error) {
	texts := lo.Map(req.Texts, func(item *vectorstorev1.VectorStoreServiceInsertTextsRequestText, _ int) *domain.InsertTextsInputText {
		return &domain.InsertTextsInputText{
			Text:     item.Text,
			Metadata: item.Metadata.AsMap(),
		}
	})

	insertTextsInput := &domain.InsertTextsInput{
		Texts: texts,
	}

	if err := grpcService.uc.InsertTexts(ctx, insertTextsInput); err != nil {
		if _, ok := err.(*vectorstore_error.Error); ok {
			return nil, status.New(codes.InvalidArgument, err.Error()).Err()
		}
		return nil, err
	}

	vectorStoreServiceInsertTextsResponse := &vectorstorev1.VectorStoreServiceInsertTextsResponse{}

	return vectorStoreServiceInsertTextsResponse, nil
}

func (grpcService *vectorStoreGRPCService) SearchText(ctx context.Context, req *vectorstorev1.VectorStoreServiceSearchTextRequest) (*vectorstorev1.VectorStoreServiceSearchTextResponse, error) {
	searchTextInput := &domain.SearchTextInput{
		Text:   req.Text,
		TopK:   int(req.TopK),
		Filter: req.Filter.AsMap(),
	}

	searchTextResults, err := grpcService.uc.SearchText(ctx, searchTextInput)
	if err != nil {
		if _, ok := err.(*vectorstore_error.Error); ok {
			return nil, status.New(codes.InvalidArgument, err.Error()).Err()
		}
		return nil, err
	}

	similarTexts := make([]*vectorstorev1.VectorStoreServiceSearchTextResponseSimilarText, 0, len(searchTextResults.SimilarTexts))

	for _, similarText := range searchTextResults.SimilarTexts {
		metadata, err := structpb.NewStruct(similarText.Metadata)
		if err != nil {
			return nil, err
		}

		similarTexts = append(
			similarTexts,
			&vectorstorev1.VectorStoreServiceSearchTextResponseSimilarText{
				Text:     similarText.Text,
				Score:    similarText.Score,
				Metadata: metadata,
			},
		)
	}

	vectorStoreServiceSearchTextResponse := &vectorstorev1.VectorStoreServiceSearchTextResponse{
		SimilarTexts: similarTexts,
	}

	return vectorStoreServiceSearchTextResponse, nil
}
