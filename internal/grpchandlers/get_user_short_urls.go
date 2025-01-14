package grpchandlers

import (
	"context"
	"net/url"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/DyadyaRodya/go-shortener/proto/v1"
)

// GetUserShortURLs Get user short URLs for full URLs
func (h *GrpcHandlers) GetUserShortURLs(ctx context.Context, _ *pb.GetUserShortURLsRequest) (*pb.GetUserShortURLsResponse, error) {
	userUUID, authorized, err := h.ProcessAuth(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if !authorized {
		return nil, status.Error(codes.Unauthenticated, "You are not authorized to perform this operation")
	}

	shortURLs, err := h.Usecases.GetUserShortURLs(ctx, userUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var fullShortURL string
	resp := make([]*pb.GetUserShortURLsResponseItem, 0, len(shortURLs))
	for _, shortURL := range shortURLs {
		fullShortURL, err = url.JoinPath(h.Config.BaseShortURL, shortURL.ID)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		resp = append(resp, &pb.GetUserShortURLsResponseItem{ShortUrl: fullShortURL, OriginalUrl: shortURL.URL})
	}

	token, err := h.JWTService.GenerateToken(userUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.GetUserShortURLsResponse{
		NewJwtToken: token,
		Urls:        resp,
	}, nil
}
