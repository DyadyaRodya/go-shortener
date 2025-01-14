package grpchandlers

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/DyadyaRodya/go-shortener/proto/v1"
)

// Ping Check server readiness
func (h *GrpcHandlers) Ping(ctx context.Context, _ *pb.PingRequest) (*pb.PingResponse, error) {
	err := h.Usecases.CheckConnection(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.PingResponse{}, nil
}
