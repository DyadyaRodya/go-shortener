package handlers

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (h *Handlers) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if contentType := r.Header.Get("Content-Type"); !strings.Contains(contentType, "text/plain") {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	defer func() {
		_ = r.Body.Close()
	}()

	sourceURL := strings.TrimSpace(string(body))
	_, err = url.ParseRequestURI(sourceURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL, err := h.Usecases.CreateShortURL(sourceURL)
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
