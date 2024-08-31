package app

import (
	"github.com/DyadyaRodya/go-shortener/internal/config"
	"github.com/DyadyaRodya/go-shortener/internal/domain/services"
	"github.com/DyadyaRodya/go-shortener/internal/handlers"
	"github.com/DyadyaRodya/go-shortener/internal/logger"
	"github.com/DyadyaRodya/go-shortener/internal/repositories/inmemory"
	"github.com/DyadyaRodya/go-shortener/internal/usecases"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"log"
)

type App struct {
	appConfig *config.Config
	e         *echo.Echo
	appLogger *zap.Logger
}

func NewApp(DefaultBaseShortURL, DefaultServerAddress, DefaultLogLevel string) *App {
	// init configs
	appConfig := config.InitConfigFromCMD(DefaultServerAddress, DefaultBaseShortURL, DefaultLogLevel)

	var appLogger *zap.Logger
	var loggerMW echo.MiddlewareFunc
	var err error
	// init logger middleware
	if appLogger, loggerMW, err = logger.Initialize(appConfig.LogLevel); err != nil {
		log.Printf("Config %+v\n", *appConfig.HandlersConfig)
		log.Fatalf("Cannot initialize logger %+v\n", err)
	}

	appLogger.Info("Config", zap.Any("config", appConfig))

	// init Echo
	e := echo.New()

	// Middleware
	e.Use(loggerMW)
	e.Use(middleware.Recover())

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
		appLogger: appLogger,
	}
}

func (a *App) Run() error {
	a.appLogger.Info("Starting server at", zap.String("address", a.appConfig.ServerAddress))
	err := a.e.Start(a.appConfig.ServerAddress)
	if err != nil {
		a.appLogger.Error("Starting error", zap.Error(err))
		return err
	}
	return nil
}
