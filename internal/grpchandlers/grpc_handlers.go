package grpchandlers

import (
	"context"
	"net"

	"google.golang.org/grpc/metadata"

	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	usecasesdto "github.com/DyadyaRodya/go-shortener/internal/usecases/dto"
	pb "github.com/DyadyaRodya/go-shortener/proto/v1"
)

// GRPC Handlers interfaces
type (
	Usecases interface {
		CheckConnection(ctx context.Context) error
		GetShortURL(ctx context.Context, ID string) (*entity.ShortURL, error)
		CreateShortURL(ctx context.Context, URL, UserUUID string) (*entity.ShortURL, error)
		BatchCreateShortURLs(
			ctx context.Context,
			createRequests []*usecasesdto.BatchCreateRequest,
			UserUUID string,
		) ([]*usecasesdto.BatchCreateResponse, error)
		GetUserShortURLs(ctx context.Context, userUUID string) ([]*entity.ShortURL, error)
		GetStats(ctx context.Context) (*usecasesdto.StatsResponse, error)
	}

	JWTService interface {
		ProcessToken(token string) (userUUID string, authenticated bool, err error)
		GenerateToken(userUUID string) (string, error)
	}

	// Config config for grpc handlers
	Config struct {
		BaseShortURL  string
		TrustedSubnet *net.IPNet
		SecretKey     []byte
	}
)

// GrpcHandlers supports all required server calls
type GrpcHandlers struct {
	pb.UnimplementedGoShortenerServiceServer
	Usecases   Usecases
	JWTService JWTService
	Config     *Config
	DelChan    chan *usecasesdto.DeleteUserShortURLsRequest
}

// NewGrpcHandlers constructor for Handlers
func NewGrpcHandlers(
	usecases Usecases,
	jwtService JWTService,
	config *Config,
	DelChan chan *usecasesdto.DeleteUserShortURLsRequest,
) *GrpcHandlers {
	return &GrpcHandlers{
		Usecases:   usecases,
		JWTService: jwtService,
		Config:     config,
		DelChan:    DelChan,
	}
}

// ProcessAuth Reads jwt token from metadata and processes it with JWTService
func (h *GrpcHandlers) ProcessAuth(ctx context.Context) (userUUID string, authenticated bool, err error) {
	var token string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("auth")
		if len(values) > 0 {
			token = values[0]
		}
	}
	return h.JWTService.ProcessToken(token)
}
