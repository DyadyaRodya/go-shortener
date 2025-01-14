package config

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"os"

	"github.com/DyadyaRodya/go-shortener/internal/grpchandlers"

	"github.com/DyadyaRodya/go-shortener/internal/handlers"
)

// Config stores all config data for app.App
type Config struct {
	HandlersConfig     *handlers.Config
	GrpcHandlersConfig *grpchandlers.Config
	TrustedSubnet      *net.IPNet
	ServerAddress      string
	GrpcAddress        string
	LogLevel           string
	StorageFile        string
	DSN                string
	EnableHTTPS        bool
}

type config struct {
	ServerAddress string `json:"server_address"`
	GrpcAddress   string `json:"grpc_address"`
	BaseURL       string `json:"base_url"`
	StorageFile   string `json:"file_storage_path"`
	DSN           string `json:"database_dsn"`
	TrustedSubnet string `json:"trusted_subnet"`
	EnableHTTPS   bool   `json:"enable_https"`
}

// InitConfigFromCMD reads CMD line and env arguments to Config
func InitConfigFromCMD(
	defaultServerAddress,
	defaultGrpcAddress,
	defaultBaseURL,
	defaultLogLevel,
	defaultStorageFile string,
) *Config {
	configFile := flag.String("c", "", "path to config file")
	flag.StringVar(configFile, "config", "", "path to config file")

	serverAddress := flag.String("a", "", "server address to bind")
	grpcAddress := flag.String("g", "", "grpc address to bind")
	baseURL := flag.String("b", "", "base url for short url")
	logLevel := flag.String("l", defaultLogLevel, "log level")
	storageFile := flag.String("f", "", "file storage path")
	dsn := flag.String("d", "", "database connection string")
	enableHTTPS := flag.Bool("s", false, "enable https")
	trustedSubnet := flag.String("t", "", "trusted subnet")
	flag.Parse()

	if envConfigFile := os.Getenv("CONFIG"); envConfigFile != "" {
		configFile = &envConfigFile
	}

	if envServerAddress := os.Getenv("SERVER_ADDRESS"); envServerAddress != "" {
		serverAddress = &envServerAddress
	}
	if envGrpcAddress := os.Getenv("GRPC_ADDRESS"); envGrpcAddress != "" {
		grpcAddress = &envGrpcAddress
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
	if envTrustedSubnet := os.Getenv("TRUSTED_SUBNET"); envTrustedSubnet != "" {
		trustedSubnet = &envTrustedSubnet
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
	if *grpcAddress == "" {
		if conf.GrpcAddress != "" {
			*grpcAddress = conf.GrpcAddress
		} else {
			*grpcAddress = defaultGrpcAddress
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
	if *trustedSubnet == "" {
		*trustedSubnet = conf.TrustedSubnet
	}

	var confNet *net.IPNet
	if *trustedSubnet != "" {
		_, ipNet, err := net.ParseCIDR(*trustedSubnet)
		if err != nil {
			log.Fatalf("failed to parse trusted subnet %s: %s", *trustedSubnet, err.Error())
		}
		confNet = ipNet
	}

	return &Config{
		HandlersConfig: &handlers.Config{
			BaseShortURL: *baseURL,
		},
		GrpcHandlersConfig: &grpchandlers.Config{
			BaseShortURL:  *baseURL,
			TrustedSubnet: confNet,
		},
		ServerAddress: *serverAddress,
		GrpcAddress:   *grpcAddress,
		LogLevel:      *logLevel,
		StorageFile:   *storageFile,
		DSN:           *dsn,
		EnableHTTPS:   *enableHTTPS,
		TrustedSubnet: confNet,
	}
}
