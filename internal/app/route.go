package app

import (
	"net/http/pprof"
	"strings"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"

	"github.com/DyadyaRodya/go-shortener/internal/handlers"
	"github.com/DyadyaRodya/go-shortener/internal/handlers/dto"

	//nolint:blank-imports // for swagger
	_ "github.com/DyadyaRodya/go-shortener/docs"
)

const (
	rootURL            = "/"
	apiShortenURL      = "/api/shorten"
	apiShortenURLBatch = "/api/shorten/batch"
	apiGetUserURLs     = "/api/user/urls"
	apiDeleteUserURLs  = "/api/user/urls"
	idURL              = "/:" + dto.IDParamName
	pingURL            = "/ping"
)

func setupRoutes(e *echo.Echo, handlers *handlers.Handlers) {
	e.GET(idURL, handlers.GetByShortURL)
	e.POST(rootURL, handlers.CreateShortURL)
	e.POST(apiShortenURL, handlers.CreateShortURLJSON)
	e.POST(apiShortenURLBatch, handlers.BatchCreateShortURLJSON)
	e.GET(apiGetUserURLs, handlers.GetUserShortURLs)
	e.DELETE(apiDeleteUserURLs, handlers.DeleteUserShortURLs)
	e.GET(pingURL, handlers.PingHandler)

	e.GET("/docs/swagger/*", echoSwagger.WrapHandler)

	Wrap(e)
}

// echo-pprof from https://github.com/sevennt/echo-pprof/blob/master/pprof.go

// Wrap adds several routes from package `net/http/pprof` to *echo.Echo object.
func Wrap(e *echo.Echo) {
	WrapGroup("", e.Group("/debug/pprof"))
}

// Wrapper make sure we are backward compatible.
var Wrapper = Wrap

// WrapGroup adds several routes from package `net/http/pprof` to *echo.Group object.
func WrapGroup(prefix string, g *echo.Group) {
	routers := []struct {
		Handler echo.HandlerFunc
		Method  string
		Path    string
	}{
		{IndexHandler(), "GET", ""},
		{IndexHandler(), "GET", "/"},
		{HeapHandler(), "GET", "/heap"},
		{GoroutineHandler(), "GET", "/goroutine"},
		{BlockHandler(), "GET", "/block"},
		{ThreadCreateHandler(), "GET", "/threadcreate"},
		{CmdlineHandler(), "GET", "/cmdline"},
		{ProfileHandler(), "GET", "/profile"},
		{SymbolHandler(), "GET", "/symbol"},
		{SymbolHandler(), "POST", "/symbol"},
		{TraceHandler(), "GET", "/trace"},
		{MutexHandler(), "GET", "/mutex"},
	}

	for _, r := range routers {
		switch r.Method {
		case "GET":
			g.GET(strings.TrimPrefix(r.Path, prefix), r.Handler)
		case "POST":
			g.POST(strings.TrimPrefix(r.Path, prefix), r.Handler)
		}
	}
}

// IndexHandler will pass the call from /debug/pprof to pprof.
func IndexHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Index(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

// HeapHandler will pass the call from /debug/pprof/heap to pprof.
func HeapHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Handler("heap").ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	}
}

// GoroutineHandler will pass the call from /debug/pprof/goroutine to pprof.
func GoroutineHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Handler("goroutine").ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

// BlockHandler will pass the call from /debug/pprof/block to pprof.
func BlockHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Handler("block").ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

// ThreadCreateHandler will pass the call from /debug/pprof/threadcreate to pprof.
func ThreadCreateHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Handler("threadcreate").ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

// CmdlineHandler will pass the call from /debug/pprof/cmdline to pprof.
func CmdlineHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Cmdline(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

// ProfileHandler will pass the call from /debug/pprof/profile to pprof.
func ProfileHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Profile(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

// SymbolHandler will pass the call from /debug/pprof/symbol to pprof.
func SymbolHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Symbol(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

// TraceHandler will pass the call from /debug/pprof/trace to pprof.
func TraceHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Trace(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

// MutexHandler will pass the call from /debug/pprof/mutex to pprof.
func MutexHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Handler("mutex").ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	}
}
