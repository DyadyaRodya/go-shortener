package tests

import (
	"context"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/DyadyaRodya/go-shortener/internal/config"
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	"github.com/DyadyaRodya/go-shortener/internal/domain/services"
	"github.com/DyadyaRodya/go-shortener/internal/repositories/inmemory"
	pgxrepo "github.com/DyadyaRodya/go-shortener/internal/repositories/pgx"
	"github.com/DyadyaRodya/go-shortener/internal/usecases"
	"github.com/DyadyaRodya/go-shortener/internal/usecases/dto"
)

func BenchmarkBatchCreateShortURLsWithDB(b *testing.B) {
	b.StopTimer()
	appConfig := config.InitConfigFromCMD(`:8080`, `http://localhost:8080/`, "info", "")
	idGenerator := services.NewIDGenerator()

	lvl, err := zap.ParseAtomicLevel(appConfig.LogLevel)
	if err != nil {
		b.Fatal(err)
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	cfg.EncoderConfig.CallerKey = zapcore.OmitKey
	appLogger, err := cfg.Build()
	if err != nil {
		b.Fatal(err)
	}

	var store usecases.URLStorage
	if appConfig.DSN != "" {
		pool, err := pgxpool.New(context.Background(), appConfig.DSN)
		if err != nil {
			b.Fatal("Cannot create connection pool to database", appConfig.DSN, err)
		}
		defer pool.Close()
		s := pgxrepo.NewStorePGX(pool, appLogger)
		store = s
	} else {
		s := inmemory.NewStoreInMemory()
		store = s
	}

	u := usecases.NewUsecases(store, idGenerator)

	sampleSize := 1000

	ctx := context.Background()
	userUUID := gofakeit.UUID()
	tx, _ := store.Begin(ctx)
	defer tx.Rollback(ctx)
	tx.AddUserIfNotExists(ctx, userUUID)

	requests := make([]*dto.BatchCreateRequest, 3*sampleSize)
	for i := 0; i < sampleSize; i++ {
		requests[i] = &dto.BatchCreateRequest{
			CorrelationID: strconv.Itoa(i + 1),
			OriginalURL:   gofakeit.URL(),
		}
	}

	for i := 0; i < sampleSize; i++ {
		id, _ := idGenerator.Generate()
		foundShortURL := &entity.ShortURL{
			ID:  id,
			URL: gofakeit.URL(),
		}
		tx.AddURL(ctx, foundShortURL, true)
		requests[i+sampleSize] = &dto.BatchCreateRequest{
			CorrelationID: strconv.Itoa(i + sampleSize + 1),
			OriginalURL:   foundShortURL.URL,
		}
	}
	for i := 0; i < sampleSize; i++ {
		id, _ := idGenerator.Generate()
		foundShortURL := &entity.ShortURL{
			ID:  id,
			URL: gofakeit.URL(),
		}
		tx.AddURL(ctx, foundShortURL, true)
		tx.AddUserURL(ctx, foundShortURL.ID, userUUID)
		requests[i+2*sampleSize] = &dto.BatchCreateRequest{
			CorrelationID: strconv.Itoa(i + 2*sampleSize + 1),
			OriginalURL:   foundShortURL.URL,
		}
	}
	// additional user urls
	for i := 0; i < sampleSize; i++ {
		id, _ := idGenerator.Generate()
		url := &entity.ShortURL{
			ID:  id,
			URL: gofakeit.URL(),
		}
		tx.AddURL(ctx, url, true)
		tx.AddUserURL(ctx, url.ID, userUUID)
	}
	tx.Commit(ctx)

	//
	b.StartTimer()
	b.ResetTimer()
	b.Run("BatchCreateShortURLs", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			u.BatchCreateShortURLs(ctx, requests, userUUID)
		}
	})
	b.StopTimer()
}
