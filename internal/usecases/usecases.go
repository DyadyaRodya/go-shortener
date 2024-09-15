package usecases

import (
	"context"
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
)

type (
	Transaction interface {
		Commit(ctx context.Context) error
		Rollback(ctx context.Context) error
		GetByURLs(ctx context.Context, URLs []string) (map[string]*entity.ShortURL, error)
		CheckIDs(ctx context.Context, IDs []string) ([]string, error)
		AddURL(ctx context.Context, ShortURL *entity.ShortURL) error
	}

	URLStorage interface {
		TestConnection(ctx context.Context) error
		GetURLByID(ctx context.Context, ID string) (*entity.ShortURL, error)
		AddURL(ctx context.Context, ShortURL *entity.ShortURL) error

		Begin(ctx context.Context) (Transaction, error)
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
