package grpchandlers

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	pb "github.com/DyadyaRodya/go-shortener/proto/v1"
)

// GetFullByID Get full URL for short URL
func (h *GrpcHandlers) GetFullByID(ctx context.Context, in *pb.GetFullByIDRequest) (*pb.GetFullByIDResponse, error) {
	shortURL, err := h.Usecases.GetShortURL(ctx, in.GetId())
	if err != nil {
		if errors.Is(err, entity.ErrShortURLDeleted) {
			return &pb.GetFullByIDResponse{Status: pb.GetFullByIDResponse_DELETED}, status.Error(codes.NotFound, err.Error())
		}

		if errors.Is(err, entity.ErrShortURLNotFound) {
			return &pb.GetFullByIDResponse{Status: pb.GetFullByIDResponse_NOT_FOUND}, status.Error(codes.NotFound, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetFullByIDResponse{FullUrl: shortURL.URL, Status: pb.GetFullByIDResponse_OK}, nil
}
