package app

import (
	"context"
	"errors"
	"github.com/DyadyaRodya/go-shortener/internal/config"
	"github.com/DyadyaRodya/go-shortener/internal/domain/services"
	"github.com/DyadyaRodya/go-shortener/internal/handlers"
	"github.com/DyadyaRodya/go-shortener/internal/logger"
	"github.com/DyadyaRodya/go-shortener/internal/repositories/inmemory"
	pgxrepo "github.com/DyadyaRodya/go-shortener/internal/repositories/pgx"
	"github.com/DyadyaRodya/go-shortener/internal/usecases"
	"github.com/DyadyaRodya/go-shortener/pkg/jsonfile"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type App struct {
	appConfig  *config.Config
	e          *echo.Echo
	appLogger  *zap.Logger
	appStorage *usecases.URLStorage
	close      func()
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
	var closeStorage func()
	var store usecases.URLStorage
	if appConfig.DSN == "" {
		s := inmemory.NewStoreInMemory()
		store = s
		closeStorage = func() {}
	} else {
		pool, err := pgxpool.New(context.Background(), appConfig.DSN)
		if err != nil {
			appLogger.Fatal("Cannot create connection pool to database", zap.String("DSN", appConfig.DSN), zap.Error(err))
		}
		closeStorage = pool.Close
		s := pgxrepo.NewStorePGX(pool, appLogger)
		store = s
	}
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
		appStorage: &store,
		close:      closeStorage,
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
	if err != nil && !strings.Contains(err.Error(), "http: Server closed") {
		a.appLogger.Error("Starting error", zap.Error(err))
		return err
	}
	return nil
}

func (a *App) Shutdown(signal os.Signal) error {
	defer a.close()

	ctx := context.Background()
	a.appLogger.Info("Stopped server on signal", zap.String("signal", signal.String()))

	if a.appConfig.StorageFile != "" {
		var err error
		var data *map[string]string
		switch v := (*a.appStorage).(type) {
		case *inmemory.StoreInMemory:
			data = v.Storage()

		case *pgxrepo.StorePGX:
			data, err = v.Save(ctx)
			if err != nil {
				return err
			}
		default:

		}
		err = jsonfile.WriteMapToFile(a.appConfig.StorageFile, data)
		if err != nil {
			return err
		}
	}

	ctx, cancelFunc := context.WithTimeout(ctx, 10*time.Second)
	defer cancelFunc()
	err := a.e.Shutdown(ctx)
	return err
}

func (a *App) readInitData() error {
	switch v := (*a.appStorage).(type) {
	case *inmemory.StoreInMemory:
		if a.appConfig.StorageFile == "" {
			return nil
		}

		err := jsonfile.ReadFileToMap(a.appConfig.StorageFile, v.Storage())

		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}

		if err != nil {
			a.appLogger.Info("Empty storage file")
		} else {
			a.appLogger.Info("Storage file successfully read")
		}

		return nil
	case *pgxrepo.StorePGX:
		ctx := context.Background()
		ctx, cancelFunc := context.WithTimeout(ctx, 3*time.Second)
		defer cancelFunc()
		err := v.InitSchema(ctx)
		if err != nil {
			return err
		}

		if a.appConfig.StorageFile == "" {
			return nil
		}

		src := make(map[string]string)
		err = jsonfile.ReadFileToMap(a.appConfig.StorageFile, &src)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}

		if err != nil {
			a.appLogger.Info("Empty storage file")
			return nil
		} else {
			a.appLogger.Info("Storage file successfully read, loading data to DB")
		}

		return v.Load(ctx, src)
	default:
		a.appLogger.Warn("Not using init file for current storage type")
		return nil

	}
}
