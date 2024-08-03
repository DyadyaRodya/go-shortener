package handlers

import (
	"encoding/hex"
	"net/http"
	"strings"
)

func (h *Handlers) GetByShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/")

	_, err := hex.DecodeString(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL, err := h.Usecases.GetShortURL(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", shortURL.URL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
