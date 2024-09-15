package handlers

import (
	"context"
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	usecasesdto "github.com/DyadyaRodya/go-shortener/internal/usecases/dto"
)

type (
	Usecases interface {
		CheckConnection(ctx context.Context) error
		GetShortURL(ctx context.Context, ID string) (*entity.ShortURL, error)
		CreateShortURL(ctx context.Context, URL string) (*entity.ShortURL, error)
		BatchCreateShortURLs(
			ctx context.Context,
			createRequests []*usecasesdto.BatchCreateRequest,
		) ([]*usecasesdto.BatchCreateResponse, error)
	}

	Config struct {
		BaseShortURL string
	}

	Handlers struct {
		Usecases Usecases
		Config   *Config
	}
)

func NewHandlers(usecases Usecases, config *Config) *Handlers {
	return &Handlers{
		Usecases: usecases,
		Config:   config,
	}
}
