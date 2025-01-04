package config

import (
	"encoding/json"
	"flag"
	"log"
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

type config struct {
	ServerAddress string `json:"server_address"`
	BaseURL       string `json:"base_url"`
	StorageFile   string `json:"file_storage_path"`
	DSN           string `json:"database_dsn"`
	EnableHTTPS   bool   `json:"enable_https"`
}

// InitConfigFromCMD reads CMD line and env arguments to Config
func InitConfigFromCMD(defaultServerAddress, defaultBaseURL, defaultLogLevel, defaultStorageFile string) *Config {
	configFile := flag.String("c", "", "path to config file")
	flag.StringVar(configFile, "config", "", "path to config file")

	serverAddress := flag.String("a", "", "server address to bind")
	baseURL := flag.String("b", "", "base url for short url")
	logLevel := flag.String("l", defaultLogLevel, "log level")
	storageFile := flag.String("f", "", "file storage path")
	dsn := flag.String("d", "", "database connection string")
	enableHTTPS := flag.Bool("s", false, "enable https")
	flag.Parse()

	if envConfigFile := os.Getenv("CONFIG"); envConfigFile != "" {
		configFile = &envConfigFile
	}

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

	conf := &config{}
	if *configFile != "" {
		f, err := os.Open(*configFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		err = json.NewDecoder(f).Decode(conf)
		if err != nil {
			log.Fatal(err)
		}
	}

	if *serverAddress == "" {
		if conf.ServerAddress != "" {
			*serverAddress = conf.ServerAddress
		} else {
			*serverAddress = defaultServerAddress
		}
	}
	if *baseURL == "" {
		if conf.BaseURL != "" {
			*baseURL = conf.BaseURL
		} else {
			*baseURL = defaultBaseURL
		}
	}
	if *storageFile == "" {
		if conf.StorageFile != "" {
			*storageFile = conf.StorageFile
		} else {
			*storageFile = defaultStorageFile
		}
	}
	if *dsn == "" {
		*dsn = conf.DSN
	}
	if !*enableHTTPS {
		*enableHTTPS = conf.EnableHTTPS
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
