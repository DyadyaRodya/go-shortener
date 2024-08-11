package app

import (
	"github.com/DyadyaRodya/go-shortener/internal/config"
	"github.com/DyadyaRodya/go-shortener/internal/domain/services"
	"github.com/DyadyaRodya/go-shortener/internal/handlers"
	"github.com/DyadyaRodya/go-shortener/internal/repositories/inmemory"
	"github.com/DyadyaRodya/go-shortener/internal/usecases"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
)

type App struct {
	appConfig *config.Config
	e         *echo.Echo
}

func NewApp(BaseShortURL, ServerAddress string) *App {
	log.SetFlags(log.Ldate | log.Ltime)

	// init Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// init configs
	handlersConfig := &handlers.Config{
		BaseShortURL: BaseShortURL,
	}
	appConfig := config.NewConfig(handlersConfig, ServerAddress)
	log.Printf("%+v\n", *handlersConfig)

	// init services
	idGenerator := services.NewIDGenerator()

	// init storage
	store := inmemory.NewStoreInMemory()

	// init usecases
	u := usecases.NewUsecases(store, idGenerator)

	// init handlers
	h := handlers.NewHandlers(u, appConfig.HandlersConfig)

	// setup handlers for routes
	setupRoutes(e, h)

	return &App{
		appConfig: appConfig,
		e:         e,
	}
}

func (a *App) Run() error {
	log.Println("Starting server at", a.appConfig.ServerAddress)
	err := a.e.Start(a.appConfig.ServerAddress)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
