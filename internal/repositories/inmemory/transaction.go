package inmemory

import (
	"context"
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	"slices"
)

type TransactionInMemory struct {
	*StoreInMemory
	buf map[string]string
}

func (t *TransactionInMemory) Commit(_ context.Context) error {
	t.StoreInMemory.storage = t.buf
	t.lock.Unlock()
	return nil
}

func (t *TransactionInMemory) Rollback(_ context.Context) error {
	t.lock.Unlock()
	return nil
}

func (t *TransactionInMemory) GetByURLs(_ context.Context, URLs []string) (map[string]*entity.ShortURL, error) {
	result := make(map[string]*entity.ShortURL)
	for id, url := range t.buf {
		if slices.Contains(URLs, url) {
			result[url] = &entity.ShortURL{
				ID:  id,
				URL: url,
			}
		}
	}
	return result, nil
}

func (t *TransactionInMemory) CheckIDs(_ context.Context, IDs []string) ([]string, error) {
	result := make([]string, 0, len(IDs))
	for _, id := range IDs {
		if _, ok := t.buf[id]; ok {
			result = append(result, id)
		}
	}
	return result, nil
}

func (t *TransactionInMemory) AddURL(_ context.Context, ShortURL *entity.ShortURL) error {
	_, ok := t.buf[ShortURL.ID]
	if ok {
		return entity.ErrUUIDTaken
	}
	for _, url := range t.buf {
		if url == ShortURL.URL {
			return entity.ErrShortURLExists
		}
	}
	t.buf[ShortURL.ID] = ShortURL.URL
	return nil
}
