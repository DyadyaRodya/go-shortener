package config

import "github.com/DyadyaRodya/go-shortener/internal/handlers"

type Config struct {
	HandlersConfig *handlers.Config
	ServerAddress  string
}

func NewConfig(handlerConfig *handlers.Config, serverAddress string) *Config {
	return &Config{
		HandlersConfig: handlerConfig,
		ServerAddress:  serverAddress,
	}
}
