package config

import (
	"flag"
	"os"

	"github.com/DyadyaRodya/go-shortener/internal/handlers"
)

// Config stores all config data for app.App
type Config struct {
	HandlersConfig *handlers.Config
	ServerAddress  string
	LogLevel       string
	StorageFile    string
	DSN            string
	EnableHTTPS    bool
}

// InitConfigFromCMD reads CMD line and env arguments to Config
func InitConfigFromCMD(defaultServerAddress, defaultBaseURL, defaultLogLevel, defaultStorageFile string) *Config {
	serverAddress := flag.String("a", defaultServerAddress, "server address to bind")
	baseURL := flag.String("b", defaultBaseURL, "base url for short url")
	logLevel := flag.String("l", defaultLogLevel, "log level")
	storageFile := flag.String("f", defaultStorageFile, "file storage path")
	dsn := flag.String("d", "", "database connection string")
	enableHTTPS := flag.Bool("s", false, "enable https")
	flag.Parse()

	if envServerAddress := os.Getenv("SERVER_ADDRESS"); envServerAddress != "" {
		serverAddress = &envServerAddress
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		baseURL = &envBaseURL
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		logLevel = &envLogLevel
	}
	if envStorageFile := os.Getenv("FILE_STORAGE_PATH"); envStorageFile != "" {
		storageFile = &envStorageFile
	}
	if envDSN := os.Getenv("DATABASE_DSN"); envDSN != "" {
		dsn = &envDSN
	}
	if envEnableHTTPS := os.Getenv("ENABLE_HTTPS"); envEnableHTTPS != "" {
		*enableHTTPS = true
	}

	return &Config{
		HandlersConfig: &handlers.Config{
			BaseShortURL: *baseURL,
		},
		ServerAddress: *serverAddress,
		LogLevel:      *logLevel,
		StorageFile:   *storageFile,
		DSN:           *dsn,
		EnableHTTPS:   *enableHTTPS,
	}
}
