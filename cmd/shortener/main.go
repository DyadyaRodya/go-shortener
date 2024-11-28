package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/DyadyaRodya/go-shortener/internal/app"
)

const (
	defaultBaseShortURL  = "http://localhost:8080/"
	defaultServerAddress = `:8080`
	defaultLogLevel      = "info"
	defaultStorageFile   = ""
)

func main() {
	server := app.NewApp(defaultBaseShortURL, defaultServerAddress, defaultLogLevel, defaultStorageFile)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		s := <-c
		err := server.Shutdown(s)
		if err != nil {
			panic(err)
		}
	}()
	err := server.Run()
	if err != nil {
		panic(err)
	}
}
