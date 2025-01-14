package grpchandlers

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	usecasesdto "github.com/DyadyaRodya/go-shortener/internal/usecases/dto"
	pb "github.com/DyadyaRodya/go-shortener/proto/v1"
)

// DeleteShortURLs Delete user short URLs
func (h *GrpcHandlers) DeleteShortURLs(ctx context.Context, in *pb.DeleteShortURLsRequest) (*pb.DeleteShortURLsResponse, error) {
	userUUID, authorized, err := h.ProcessAuth(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if !authorized {
		return nil, status.Error(codes.Unauthenticated, "You are not authorized to perform this operation")
	}

	req := &usecasesdto.DeleteUserShortURLsRequest{
		UserUUID:      userUUID,
		ShortURLUUIDs: in.GetIds(),
	}
	go func() { h.DelChan <- req }()

	token, err := h.JWTService.GenerateToken(userUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.DeleteShortURLsResponse{NewJwtToken: token}, status.Error(codes.OK, "Request accepted")
}
