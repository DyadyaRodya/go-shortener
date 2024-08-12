package config

import (
	"flag"
	"github.com/DyadyaRodya/go-shortener/internal/handlers"
)

type Config struct {
	HandlersConfig *handlers.Config
	ServerAddress  string
}

func InitConfigFromCMD(defaultServerAddress, defaultBaseURL string) *Config {
	serverAddress := flag.String("a", defaultServerAddress, "server address to bind")
	baseURL := flag.String("b", defaultBaseURL, "base url for short url")
	flag.Parse()

	return &Config{
		HandlersConfig: &handlers.Config{
			BaseShortURL: *baseURL,
		},
		ServerAddress: *serverAddress,
	}
}
