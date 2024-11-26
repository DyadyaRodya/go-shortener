package logger

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Initialize(level string) (*zap.Logger, echo.MiddlewareFunc, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, nil, err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	cfg.EncoderConfig.CallerKey = zapcore.OmitKey
	Log, err := cfg.Build()
	if err != nil {
		return nil, nil, err
	}

	loggerContextMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			res := next(c)
			duration := time.Since(start)

			request := c.Request()
			response := c.Response()
			Log.Info("HTTP response for",
				zap.String("method", request.Method),
				zap.String("URI", request.RequestURI),
				zap.Int("status", response.Status),
				zap.Int64("size", response.Size),
				zap.Duration("duration", duration),
			)

			return res
		}
	}
	return Log, loggerContextMiddleware, nil
}
