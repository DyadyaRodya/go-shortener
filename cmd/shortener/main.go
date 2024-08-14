package main

import "github.com/DyadyaRodya/go-shortener/internal/app"

const (
	defaultBaseShortURL  = "http://localhost:8080/"
	defaultServerAddress = `:8080`
)

func main() {
	server := app.NewApp(defaultBaseShortURL, defaultServerAddress)

	err := server.Run()
	if err != nil {
		panic(err)
	}
}
