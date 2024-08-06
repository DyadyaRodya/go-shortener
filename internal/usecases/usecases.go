package usecases

import "github.com/DyadyaRodya/go-shortener/internal/domain/entity"

type (
	URLStorage interface {
		GetURLByID(string) (*entity.ShortURL, error)
		AddURL(*entity.ShortURL) error
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
