package app

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/DyadyaRodya/go-shortener/internal/config"
	"github.com/DyadyaRodya/go-shortener/internal/domain/services"
	"github.com/DyadyaRodya/go-shortener/internal/handlers"
	"github.com/DyadyaRodya/go-shortener/internal/logger"
	"github.com/DyadyaRodya/go-shortener/internal/repositories/inmemory"
	"github.com/DyadyaRodya/go-shortener/internal/usecases"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"io"
	"log"
	"os"
	"time"
)

type App struct {
	appConfig  *config.Config
	e          *echo.Echo
	appLogger  *zap.Logger
	appStorage *inmemory.StoreInMemory
}

func NewApp(DefaultBaseShortURL, DefaultServerAddress, DefaultLogLevel, DefaultStorageFile string) *App {
	// init configs
	appConfig := config.InitConfigFromCMD(DefaultServerAddress, DefaultBaseShortURL, DefaultLogLevel, DefaultStorageFile)

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
	e.Use(NewGZIPMiddleware())
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
		appConfig:  appConfig,
		e:          e,
		appLogger:  appLogger,
		appStorage: store,
	}
}

func (a *App) Run() error {
	a.appLogger.Info("Reading file storage", zap.String("path", a.appConfig.StorageFile))
	err := a.readInitData()
	if err != nil {
		a.appLogger.Error("Reading file storage error",
			zap.String("path", a.appConfig.StorageFile),
			zap.Error(err),
		)
		return err
	}
	a.appLogger.Info("Starting server at", zap.String("address", a.appConfig.ServerAddress))
	err = a.e.Start(a.appConfig.ServerAddress)
	if err != nil {
		a.appLogger.Error("Starting error", zap.Error(err))
		return err
	}
	return nil
}

func (a *App) Shutdown(signal os.Signal) error {
	a.appLogger.Info("Stopped server on signal", zap.String("signal", signal.String()))
	file, err := os.OpenFile(a.appConfig.StorageFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil
	}

	defer file.Close()

	err = json.NewEncoder(file).Encode(a.appStorage.Storage())
	if err != nil {
		return err
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	err = a.e.Shutdown(ctx)
	return err
}

func (a *App) readInitData() error {
	file, err := os.OpenFile(a.appConfig.StorageFile, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(a.appStorage.Storage())
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	if err != nil {
		a.appLogger.Info("Empty storage file")
	} else {
		a.appLogger.Info("Storage file successfully read")
	}

	return nil
}
