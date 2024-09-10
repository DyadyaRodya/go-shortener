package handlers

import (
	"context"
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
)

type (
	Usecases interface {
		CheckConnection(ctx context.Context) error
		GetShortURL(ctx context.Context, ID string) (*entity.ShortURL, error)
		CreateShortURL(ctx context.Context, URL string) (*entity.ShortURL, error)
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
