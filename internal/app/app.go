package app

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/DyadyaRodya/go-shortener/internal/grpchandlers"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"

	"github.com/DyadyaRodya/go-shortener/internal/auth"
	"github.com/DyadyaRodya/go-shortener/internal/config"
	"github.com/DyadyaRodya/go-shortener/internal/domain/services"
	"github.com/DyadyaRodya/go-shortener/internal/handlers"
	"github.com/DyadyaRodya/go-shortener/internal/logger"
	"github.com/DyadyaRodya/go-shortener/internal/repositories/inmemory"
	pgxrepo "github.com/DyadyaRodya/go-shortener/internal/repositories/pgx"
	"github.com/DyadyaRodya/go-shortener/internal/usecases"
	"github.com/DyadyaRodya/go-shortener/internal/usecases/dto"
	pb "github.com/DyadyaRodya/go-shortener/proto/v1"
)

// App Main structure for shortener service
type App struct {
	appConfig  *config.Config
	e          *echo.Echo
	gh         *grpchandlers.GrpcHandlers
	grpcServer *grpc.Server
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
func NewApp(
	DefaultBaseShortURL,
	DefaultServerAddress,
	DefaultGrpcAddress,
	DefaultLogLevel,
	DefaultStorageFile string,
) *App {
	// init configs
	appConfig := config.InitConfigFromCMD(
		DefaultServerAddress,
		DefaultGrpcAddress,
		DefaultBaseShortURL,
		DefaultLogLevel,
		DefaultStorageFile,
	)

	var appLogger *zap.Logger
	var err error
	// init logger middleware
	if appLogger, err = logger.BuildLogger(appConfig.LogLevel); err != nil {
		log.Printf("Config %+v\n", *appConfig.HandlersConfig)
		log.Fatalf("Cannot initialize logger %+v\n", err)
	}
	loggerMW := logger.InitializeEcho(appLogger)

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
	appConfig.GrpcHandlersConfig.SecretKey = secretKey

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
	jwtService := auth.NewJWTService(secretKey, uuid4Generator)

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
	gh := grpchandlers.NewGrpcHandlers(u, jwtService, appConfig.GrpcHandlersConfig, delChan)

	// setup handlers for routes
	setupRoutes(e, h, auth.NewTrustedSubnetMiddleware(appConfig.TrustedSubnet))

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
		gh:         gh,
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

	g, _ := errgroup.WithContext(context.Background())

	if a.appConfig.EnableHTTPS {
		ensureCertAndKeyExist(a.appLogger)
	}

	g.Go(func() error {
		if a.appConfig.EnableHTTPS {
			a.appLogger.Info("Ensure certificate and private key exist")
			a.appLogger.Info("Starting HTTPS server", zap.String("address", a.appConfig.ServerAddress))
			err = a.e.StartTLS(a.appConfig.ServerAddress, "goshortener.cert.pem", "goshortener.key.pem")
		} else {
			a.appLogger.Info("Starting HTTP server", zap.String("address", a.appConfig.ServerAddress))
			err = a.e.Start(a.appConfig.ServerAddress)

		}
		if err != nil && !strings.Contains(err.Error(), "http: Server closed") {
			a.appLogger.Error("Starting error", zap.Error(err))
			return err
		}
		return nil
	})
	g.Go(func() error {
		listen, err := net.Listen("tcp", a.appConfig.GrpcAddress)
		if err != nil {
			a.appLogger.Error("Starting gRPC error", zap.Error(err))
			return err
		}

		var s *grpc.Server
		if a.appConfig.EnableHTTPS {
			creds, _ := credentials.NewServerTLSFromFile("goshortener.cert.pem", "goshortener.key.pem")
			s = grpc.NewServer(grpc.Creds(creds), grpc.ChainUnaryInterceptor(logging.UnaryServerInterceptor(logger.InterceptorLogger(a.appLogger))))
		} else {
			s = grpc.NewServer(grpc.ChainUnaryInterceptor(logging.UnaryServerInterceptor(logger.InterceptorLogger(a.appLogger))))
		}

		pb.RegisterGoShortenerServiceServer(s, a.gh)
		a.grpcServer = s

		a.appLogger.Info("gRPC started at", zap.String("grpc_address", a.appConfig.GrpcAddress))
		if err := s.Serve(listen); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			a.appLogger.Error("Serving gRPC error", zap.Error(err))
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
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

	a.grpcServer.GracefulStop()

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
