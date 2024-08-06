package handlers

import (
	"github.com/DyadyaRodya/go-shortener/internal/handlers/dto"
	"net/http"
)

func (h *Handlers) GetByShortURL(w http.ResponseWriter, r *http.Request) {
	getShortURLData, errorResponse := dto.GetShortURLDataFromRequest(r)
	if errorResponse != nil {
		w.WriteHeader(errorResponse.Code)
		return
	}

	shortURL, err := h.Usecases.GetShortURL(getShortURLData.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", shortURL.URL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
