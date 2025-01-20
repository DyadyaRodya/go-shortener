package grpchandlers

import (
	"context"
	"net/url"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	usecasesdto "github.com/DyadyaRodya/go-shortener/internal/usecases/dto"
	pb "github.com/DyadyaRodya/go-shortener/proto/v1"
)

// BatchCreateShortURL Batch create short URLs for full URLs
func (h *GrpcHandlers) BatchCreateShortURL(ctx context.Context, in *pb.BatchCreateShortURLRequest) (*pb.BatchCreateShortURLResponse, error) {
	req := make([]*usecasesdto.BatchCreateRequest, 0, len(in.GetUrls()))
	for _, d := range in.GetUrls() {
		req = append(req, &usecasesdto.BatchCreateRequest{
			CorrelationID: d.GetCorrelationId(),
			OriginalURL:   d.GetOriginalUrl(),
		})
	}

	userUUID, _, err := h.ProcessAuth(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp, err := h.Usecases.BatchCreateShortURLs(ctx, req, userUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var fullShortURL string
	result := make([]*pb.BatchCreateShortURLResponseItem, 0, len(resp))
	for _, d := range resp {
		fullShortURL, err = url.JoinPath(h.Config.BaseShortURL, d.ShortURL.ID)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		result = append(result, &pb.BatchCreateShortURLResponseItem{
			CorrelationId: d.CorrelationID,
			ShortUrl:      fullShortURL,
		})
	}

	token, err := h.JWTService.GenerateToken(userUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.BatchCreateShortURLResponse{
		NewJwtToken: token,
		Results:     result,
	}, nil
}
