package usecases

import (
	"context"
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
)

type (
	URLStorage interface {
		TestConnection(ctx context.Context) error
		GetURLByID(ctx context.Context, ID string) (*entity.ShortURL, error)
		AddURL(ctx context.Context, ShortURL *entity.ShortURL) error
	}
	IDGenerator interface {
		Generate() (string, error)
	}
	Usecases struct {
		urlStorage  URLStorage
		idGenerator IDGenerator
	}
)

func NewUsecases(URLStorage URLStorage, IDGenerator IDGenerator) *Usecases {
	return &Usecases{
		urlStorage:  URLStorage,
		idGenerator: IDGenerator,
	}
}
