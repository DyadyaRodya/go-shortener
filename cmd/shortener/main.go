package main

import "github.com/DyadyaRodya/go-shortener/internal/app"

func main() {
	server := app.NewApp("http://localhost:8080/", `:8080`)

	err := server.Run()
	if err != nil {
		panic(err)
	}
}
