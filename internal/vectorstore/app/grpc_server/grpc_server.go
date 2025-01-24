package grpc_server

import (
	"context"
	"log/slog"

	vectorstorev1 "github.com/aria3ppp/rag-server/gen/go/vectorstore/v1"
	internal_error "github.com/aria3ppp/rag-server/internal/pkg/error"
	"github.com/aria3ppp/rag-server/internal/vectorstore/domain"
	"github.com/aria3ppp/rag-server/internal/vectorstore/usecase"

	"github.com/samber/lo"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	grpc_codes "google.golang.org/grpc/codes"
	grpc_status "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type grpcServer struct {
	vectorstorev1.UnimplementedVectorStoreServiceServer
	uc     usecase.UseCase
	tracer trace.Tracer
	logger *slog.Logger
}

var _ vectorstorev1.VectorStoreServiceServer = (*grpcServer)(nil)

func NewGRPCServer(
	uc usecase.UseCase,
	tracer trace.Tracer,
	logger *slog.Logger,
) *grpcServer {
	return &grpcServer{
		uc:     uc,
		tracer: tracer,
		logger: logger,
	}
}

func (grpcServer *grpcServer) InsertTexts(ctx context.Context, req *vectorstorev1.VectorStoreServiceInsertTextsRequest) (_ *vectorstorev1.VectorStoreServiceInsertTextsResponse, err error) {
	ctx, span := grpcServer.tracer.Start(ctx, "grpcServer.InsertTexts")
	defer func() {
		defer span.End()
		if err != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	texts := lo.Map(req.Texts, func(item *vectorstorev1.VectorStoreServiceInsertTextsRequestText, _ int) *domain.InsertTextsInputText {
		return &domain.InsertTextsInputText{
			Text:     item.Text,
			Metadata: item.Metadata.AsMap(),
		}
	})

	insertTextsInput := &domain.InsertTextsInput{
		Texts: texts,
	}

	if err := grpcServer.uc.InsertTexts(ctx, insertTextsInput); err != nil {
		grpcServer.logger.ErrorContext(ctx, "failed to usecase insert texts", slog.String("error", err.Error()))
		if _, ok := err.(*internal_error.ValidationError); ok {
			return nil, grpc_status.New(grpc_codes.InvalidArgument, err.Error()).Err()
		}
		return nil, err
	}

	vectorStoreServiceInsertTextsResponse := &vectorstorev1.VectorStoreServiceInsertTextsResponse{}

	return vectorStoreServiceInsertTextsResponse, nil
}

func (grpcServer *grpcServer) SearchText(ctx context.Context, req *vectorstorev1.VectorStoreServiceSearchTextRequest) (_ *vectorstorev1.VectorStoreServiceSearchTextResponse, err error) {
	ctx, span := grpcServer.tracer.Start(ctx, "grpcServer.SearchText")
	defer func() {
		defer span.End()
		if err != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	searchTextInput := &domain.SearchTextInput{
		Text:     req.Text,
		TopK:     int(req.TopK),
		MinScore: req.MinScore,
		Filter:   req.Filter.AsMap(),
	}

	searchTextResults, err := grpcServer.uc.SearchText(ctx, searchTextInput)
	if err != nil {
		grpcServer.logger.ErrorContext(ctx, "failed to usecase search text", slog.String("error", err.Error()))
		if _, ok := err.(*internal_error.ValidationError); ok {
			return nil, grpc_status.New(grpc_codes.InvalidArgument, err.Error()).Err()
		}
		return nil, err
	}

	similarTexts := make([]*vectorstorev1.VectorStoreServiceSearchTextResponseSimilarText, 0, len(searchTextResults.SimilarTexts))

	for _, similarText := range searchTextResults.SimilarTexts {
		metadata, err := structpb.NewStruct(similarText.Metadata)
		if err != nil {
			grpcServer.logger.ErrorContext(ctx, "failed to structpb new struct", slog.String("error", err.Error()))
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
