package app

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"github.com/DyadyaRodya/go-shortener/internal/auth"
	"github.com/DyadyaRodya/go-shortener/internal/config"
	"github.com/DyadyaRodya/go-shortener/internal/domain/services"
	"github.com/DyadyaRodya/go-shortener/internal/handlers"
	"github.com/DyadyaRodya/go-shortener/internal/logger"
	"github.com/DyadyaRodya/go-shortener/internal/repositories/inmemory"
	pgxrepo "github.com/DyadyaRodya/go-shortener/internal/repositories/pgx"
	"github.com/DyadyaRodya/go-shortener/internal/usecases"
	"github.com/DyadyaRodya/go-shortener/internal/usecases/dto"
)

// App Main structure for shortener service
type App struct {
	appConfig  *config.Config
	e          *echo.Echo
	appLogger  *zap.Logger
	appStorage *usecases.URLStorage
	close      func()
}

// NewApp Creates new *App for shortener service.
//
// Reads config, generates keys if needed,
// initialize structures, setup handlers, middlewares,
// routes, and starts deleter goroutine.
//
// Does not start server for handling requests
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

	secretKeyString := os.Getenv("SECRET_KEY")
	var secretKey []byte
	if secretKeyString == "" {
		secretKey = newSecretKey(32)

		base64Text := make([]byte, base64.URLEncoding.EncodedLen(len(secretKey)))
		base64.URLEncoding.Encode(base64Text, secretKey)
		appLogger.Debug("New secret key", zap.ByteString("SECRET_KEY", base64Text))
	} else {
		appLogger.Debug("old secret key", zap.String("SECRET_KEY", secretKeyString))

		secretKey = make([]byte, base64.URLEncoding.DecodedLen(len(secretKeyString))-1)
		n, err := base64.URLEncoding.Decode(secretKey, []byte(secretKeyString))
		appLogger.Debug("after decoding secret key", zap.Int("n", n), zap.Error(err))
	}

	// init Echo
	e := echo.New()

	var uuid4Generator auth.UUIDGenerator = services.NewUUID4Generator()

	// Middleware
	e.Use(loggerMW)
	e.Use(NewGZIPMiddleware())
	e.Use(middleware.Recover())
	e.Use(auth.NewAuthJWTMiddleware(&uuid4Generator, secretKey))

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

	// separate deleter into other routine
	ctxDeleter := context.Background()
	ctxDeleter, stopDeleter := context.WithCancel(ctxDeleter)
	delChan := make(chan *dto.DeleteUserShortURLsRequest, 1024)
	go u.UsersShortURLsDeleter(ctxDeleter, func(msg string) {
		appLogger.Error(msg) // pass error logger function to log some errors
	}, delChan)

	// init handlers
	h := handlers.NewHandlers(u, appConfig.HandlersConfig, delChan)

	// setup handlers for routes
	setupRoutes(e, h)

	// need to close all connections, channels and stop goroutines
	closer := func() {
		closeStorage()
		stopDeleter()
		for range delChan {

		}
		close(delChan)
	}

	return &App{
		appConfig:  appConfig,
		e:          e,
		appLogger:  appLogger,
		appStorage: &store,
		close:      closer,
	}
}

// Run Starts App server for handling requests.
//
// Reads init data from file if provided in config, loads it to storage
func (a *App) Run() error {
	a.appLogger.Info("Initializing storage and reading file", zap.String("path", a.appConfig.StorageFile))
	err := a.readInitData()
	if err != nil {
		a.appLogger.Error("Initializing storage or reading file error",
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

// Shutdown Graceful shutdown of running server.
//
// Stops serving requests, saves data from storage to file if provided. Stops deleter goroutine, closes chan.
func (a *App) Shutdown(signal os.Signal) error {
	defer a.close()

	ctx := context.Background()
	a.appLogger.Info("Stopped server on signal", zap.String("signal", signal.String()))

	if a.appConfig.StorageFile != "" {
		switch v := (*a.appStorage).(type) {
		case *inmemory.StoreInMemory:
			err := a.saveFromStoreInMemory(v)
			if err != nil {
				return err
			}
		case *pgxrepo.StorePGX:
			err := a.saveFromStorePGX(ctx, v)
			if err != nil {
				return err
			}
		default:

		}
	}

	ctx, cancelFunc := context.WithTimeout(ctx, 10*time.Second)
	defer cancelFunc()
	err := a.e.Shutdown(ctx)
	return err
}

func newSecretKey(size int) []byte {
	secretKey := make([]byte, size)
	_, err := rand.Read(secretKey)
	if err != nil {
		panic(err)
	}
	return secretKey
}
