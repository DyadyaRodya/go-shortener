package grpchandlers

import (
	"context"
	"errors"
	"net/url"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/DyadyaRodya/go-shortener/proto/v1"

	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
)

// CreateShortURL Create short URL for full URL
func (h *GrpcHandlers) CreateShortURL(ctx context.Context, in *pb.CreateShortURLRequest) (*pb.CreateShortURLResponse, error) {
	userUUID, _, err := h.ProcessAuth(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	shortURL, err := h.Usecases.CreateShortURL(ctx, in.GetUrl(), userUUID)
	if err != nil && !errors.Is(err, entity.ErrShortURLExists) {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var statusCode pb.CreateShortURLResponse_Status
	if err != nil {
		statusCode = pb.CreateShortURLResponse_ALREADY_EXISTS
	} else {
		statusCode = pb.CreateShortURLResponse_CREATED
	}

	fullShortURL, err := url.JoinPath(h.Config.BaseShortURL, shortURL.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	token, err := h.JWTService.GenerateToken(userUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.CreateShortURLResponse{
		NewJwtToken: token,
		Result:      fullShortURL,
		Status:      statusCode,
	}, nil
}
