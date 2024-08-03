package usecases

import "github.com/DyadyaRodya/go-shortener/internal/domain/entity"

type (
	URLStorage interface {
		GetURLByID(string) (*entity.ShortURL, error)
		AddURL(*entity.ShortURL) error
	}
	Usecases struct {
		urlStorage URLStorage
	}
)

func NewUsecases(URLStorage URLStorage) *Usecases {
	return &Usecases{
		urlStorage: URLStorage,
	}
}
