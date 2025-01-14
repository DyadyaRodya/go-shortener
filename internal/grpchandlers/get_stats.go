package grpchandlers

import (
	"context"
	"net"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/DyadyaRodya/go-shortener/proto/v1"
)

// GetStats Get total numbers  of users and shortened URLs
func (h *GrpcHandlers) GetStats(ctx context.Context, in *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	if h.Config.TrustedSubnet == nil {
		return nil, status.Error(codes.PermissionDenied, "Forbidden")
	}

	if ip := net.ParseIP(in.GetRealIp()); ip == nil || !h.Config.TrustedSubnet.Contains(ip) {
		return nil, status.Error(codes.PermissionDenied, "Forbidden")
	}

	stats, err := h.Usecases.GetStats(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetStatsResponse{
		Urls:  uint32(stats.URLs),
		Users: uint32(stats.Users),
	}, nil
}
