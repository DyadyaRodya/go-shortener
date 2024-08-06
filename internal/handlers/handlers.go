package handlers

import (
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	"log"
	"net/http"
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

func (h *Handlers) RootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Request at", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
	switch r.Method {
	case http.MethodGet:
		h.GetByShortURL(w, r)
	case http.MethodPost:
		h.CreateShortURL(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	log.Println("Response at", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
}
