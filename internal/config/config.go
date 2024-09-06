package config

import (
	"flag"
	"github.com/DyadyaRodya/go-shortener/internal/handlers"
	"os"
)

type Config struct {
	HandlersConfig *handlers.Config
	ServerAddress  string
	LogLevel       string
	StorageFile    string
}

func InitConfigFromCMD(defaultServerAddress, defaultBaseURL, defaultLogLevel, defaultStorageFile string) *Config {
	serverAddress := flag.String("a", defaultServerAddress, "server address to bind")
	baseURL := flag.String("b", defaultBaseURL, "base url for short url")
	logLevel := flag.String("l", defaultLogLevel, "log level")
	storageFile := flag.String("f", defaultStorageFile, "file storage path")
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

	return &Config{
		HandlersConfig: &handlers.Config{
			BaseShortURL: *baseURL,
		},
		ServerAddress: *serverAddress,
		LogLevel:      *logLevel,
		StorageFile:   *storageFile,
	}
}
