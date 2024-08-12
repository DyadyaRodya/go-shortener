package app

import (
	"fmt"
	"github.com/DyadyaRodya/go-shortener/internal/config"
	"github.com/DyadyaRodya/go-shortener/internal/domain/services"
	"github.com/DyadyaRodya/go-shortener/internal/handlers"
	"github.com/DyadyaRodya/go-shortener/internal/repositories/inmemory"
	"github.com/DyadyaRodya/go-shortener/internal/usecases"
	"log"
	"net/http"
)

type App struct {
	appConfig *config.Config
	mux       *http.ServeMux
}

func NewApp(BaseShortURL, ServerAddress string) *App {
	log.SetFlags(log.Ldate | log.Ltime)
	mux := http.NewServeMux()

	handlersConfig := &handlers.Config{
		BaseShortURL: BaseShortURL,
	}
	appConfig := config.NewConfig(handlersConfig, ServerAddress)
	log.Printf("%+v\n", *handlersConfig)

	idGenerator := services.NewIDGenerator()
	store := inmemory.NewStoreInMemory()
	u := usecases.NewUsecases(store, idGenerator)
	h := handlers.NewHandlers(u, appConfig.HandlersConfig)

	setupRoutes(mux, h)
	return &App{
		appConfig: appConfig,
		mux:       mux,
	}
}

func (a *App) Run() error {
	log.Println("Starting server at", a.appConfig.ServerAddress)
	err := http.ListenAndServe(a.appConfig.ServerAddress, a.mux)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
