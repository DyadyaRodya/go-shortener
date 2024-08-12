package app

import (
	"github.com/DyadyaRodya/go-shortener/internal/handlers"
	"net/http"
)

const (
	rootURL = "/"
	// idURL   = "/{id}"
)

func setupRoutes(mux *http.ServeMux, handlers *handlers.Handlers) {
	mux.HandleFunc(rootURL, handlers.RootHandler)
}
