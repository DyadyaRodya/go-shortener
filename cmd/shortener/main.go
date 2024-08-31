package main

import "github.com/DyadyaRodya/go-shortener/internal/app"

const (
	defaultBaseShortURL  = "http://localhost:8080/"
	defaultServerAddress = `:8080`
	defaultLogLevel      = "info"
)

func main() {
	server := app.NewApp(defaultBaseShortURL, defaultServerAddress, defaultLogLevel)

	err := server.Run()
	if err != nil {
		panic(err)
	}
}
