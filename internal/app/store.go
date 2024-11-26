package app

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/DyadyaRodya/go-shortener/internal/repositories/inmemory"
	"github.com/DyadyaRodya/go-shortener/internal/repositories/pgx"
	"github.com/DyadyaRodya/go-shortener/pkg/jsonfile"
)

// StoreBuf buffer for reading data from file
type StoreBuf struct {
	URLs           *map[string]string   `json:"urls"`
	UsersShortUrls *map[string][]string `json:"usersShortUrls"`
}

func (a *App) readInitData() error {
	switch v := (*a.appStorage).(type) {
	case *inmemory.StoreInMemory:
		return a.loadToStoreInMemory(v)

	case *pgx.StorePGX:
		ctx := context.Background()
		ctx, cancelFunc := context.WithTimeout(ctx, 3*time.Second)
		defer cancelFunc()
		err := v.InitSchema(ctx)
		if err != nil {
			return err
		}

		return a.loadToStorePGX(ctx, v)
	default:
		a.appLogger.Warn("Not using init file for current storage type")
		return nil

	}
}

func (a *App) loadToStoreInMemory(store *inmemory.StoreInMemory) error {
	if a.appConfig.StorageFile == "" {
		return nil
	}

	buf := &StoreBuf{
		URLs:           store.URLS(),
		UsersShortUrls: store.UsersShortUrls(),
	}

	err := jsonfile.ReadFileToAny(a.appConfig.StorageFile, buf)

	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	if err != nil {
		a.appLogger.Info("Empty storage file")
	} else {
		a.appLogger.Info("Storage file successfully read")
	}

	return nil
}

func (a *App) saveFromStoreInMemory(store *inmemory.StoreInMemory) error {
	buf := &StoreBuf{
		URLs:           store.URLS(),
		UsersShortUrls: store.UsersShortUrls(),
	}

	return jsonfile.WriteAnyToFile(a.appConfig.StorageFile, buf)
}

func (a *App) loadToStorePGX(ctx context.Context, store *pgx.StorePGX) error {
	if a.appConfig.StorageFile == "" {
		return nil
	}

	buf := &StoreBuf{}
	err := jsonfile.ReadFileToAny(a.appConfig.StorageFile, buf)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	if err != nil {
		a.appLogger.Info("Empty storage file")
		return nil
	} else {
		a.appLogger.Info("Storage file successfully read, loading data to DB")
	}

	err = store.LoadURLs(ctx, *buf.URLs)
	if err != nil {
		return err
	}

	return store.LoadUsersURLs(ctx, *buf.UsersShortUrls)
}

func (a *App) saveFromStorePGX(ctx context.Context, store *pgx.StorePGX) error {
	urls, err := store.URLs(ctx)
	if err != nil {
		return err
	}
	usersURLs, err := store.UsersURLs(ctx)
	if err != nil {
		return err
	}
	buf := &StoreBuf{
		URLs:           urls,
		UsersShortUrls: usersURLs,
	}
	return jsonfile.WriteAnyToFile(a.appConfig.StorageFile, buf)
}
