package inmemory

import (
	"context"
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	"slices"
)

type TransactionInMemory struct {
	*StoreInMemory
	buf               map[string]string
	bufUsersShortURLs map[string][]string
}

func (t *TransactionInMemory) Commit(_ context.Context) error {
	t.StoreInMemory.urls = t.buf
	t.StoreInMemory.usersShortUrls = t.bufUsersShortURLs
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

func (t *TransactionInMemory) AddUserIfNotExists(_ context.Context, UserUUID string) error {
	if _, ok := t.bufUsersShortURLs[UserUUID]; !ok {
		t.bufUsersShortURLs[UserUUID] = []string{}
	}
	return nil
}

func (t *TransactionInMemory) AddUserURL(_ context.Context, ShortURLUUID, UserUUID string) error {
	if owns, ok := t.bufUsersShortURLs[UserUUID]; !ok {
		t.bufUsersShortURLs[UserUUID] = []string{ShortURLUUID}
	} else if !slices.Contains(owns, ShortURLUUID) {
		t.bufUsersShortURLs[UserUUID] = append(owns, ShortURLUUID)
	}
	return nil
}

func (t *TransactionInMemory) GetUserUrls(_ context.Context, UserUUID string) (map[string]*entity.ShortURL, error) {
	result := make(map[string]*entity.ShortURL)

	owns, ok := t.bufUsersShortURLs[UserUUID]
	if !ok {
		return result, nil
	}

	for _, uuid := range owns {
		result[uuid] = &entity.ShortURL{
			ID:  uuid,
			URL: t.buf[uuid],
		}
	}
	return result, nil
}

func (t *TransactionInMemory) GetShortByURL(_ context.Context, URL string) (*entity.ShortURL, error) {
	for uuid, url := range t.buf {
		if url == URL {
			shortURL := &entity.ShortURL{
				ID:  uuid,
				URL: url,
			}
			return shortURL, nil
		}
	}

	return nil, entity.ErrShortURLNotFound
}
