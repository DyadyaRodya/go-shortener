package config

import (
	"flag"
	"github.com/DyadyaRodya/go-shortener/internal/handlers"
	"os"
)

type Config struct {
	HandlersConfig *handlers.Config
	ServerAddress  string
}

func InitConfigFromCMD(defaultServerAddress, defaultBaseURL string) *Config {
	serverAddress := flag.String("a", defaultServerAddress, "server address to bind")
	baseURL := flag.String("b", defaultBaseURL, "base url for short url")
	flag.Parse()

	if envServerAddress := os.Getenv("SERVER_ADDRESS"); envServerAddress != "" {
		serverAddress = &envServerAddress
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		baseURL = &envBaseURL
	}

	return &Config{
		HandlersConfig: &handlers.Config{
			BaseShortURL: *baseURL,
		},
		ServerAddress: *serverAddress,
	}
}
