package handlers

import (
	"github.com/DyadyaRodya/go-shortener/internal/handlers/dto"
	"net/http"
	"net/url"
)

func (h *Handlers) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	createShortURLData, errorResponse := dto.CreateShortURLDataFromRequest(r)
	if errorResponse != nil {
		w.WriteHeader(errorResponse.Code)
		return
	}

	shortURL, err := h.Usecases.CreateShortURL(createShortURLData.URL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	fullShortURL, err := url.JoinPath(h.Config.BaseShortURL, shortURL.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(fullShortURL))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
