package handlers

import (
	"context"

	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	usecasesdto "github.com/DyadyaRodya/go-shortener/internal/usecases/dto"
)

// Handlers interfaces
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
	}

	Config struct {
		BaseShortURL string
	}
)

// Handlers
// @title URL Shortener API
// @version 0.1
// @description URL Shortener API OpenAPI
// @termsOfService http://swagger.io/terms/
type Handlers struct {
	Usecases Usecases
	Config   *Config
	DelChan  chan *usecasesdto.DeleteUserShortURLsRequest
}

// NewHandlers constructor for Handlers
func NewHandlers(usecases Usecases, config *Config, DelChan chan *usecasesdto.DeleteUserShortURLsRequest) *Handlers {
	return &Handlers{
		Usecases: usecases,
		Config:   config,
		DelChan:  DelChan,
	}
}
