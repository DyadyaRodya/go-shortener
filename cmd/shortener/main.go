package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/DyadyaRodya/go-shortener/internal/app"
)

const (
	defaultBaseShortURL  = "http://localhost:8080/"
	defaultServerAddress = `:8080`
	defaultGrpcAddress   = `:50051`
	defaultLogLevel      = "info"
	defaultStorageFile   = ""
)

var buildVersion = "N/A" //nolint: gochecknoglobals // This var could be global
var buildDate = "N/A"    //nolint: gochecknoglobals // This var could be global
var buildCommit = "N/A"  //nolint: gochecknoglobals // This var could be global

func main() {
	fmt.Printf(
		"Build version: %s\nBuild date: %s\nBuild commit: %s\n",
		buildVersion,
		buildDate,
		buildCommit,
	)
	server := app.NewApp(
		defaultBaseShortURL,
		defaultServerAddress,
		defaultGrpcAddress,
		defaultLogLevel,
		defaultStorageFile,
	)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
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
