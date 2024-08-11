package handlers

import (
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
)

type (
	Usecases interface {
		GetShortURL(ID string) (*entity.ShortURL, error)
		CreateShortURL(URL string) (*entity.ShortURL, error)
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
