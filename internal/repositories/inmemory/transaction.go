package inmemory

import (
	"context"
	"slices"

	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
)

// TransactionInMemory imitate storage transaction (session)
type TransactionInMemory struct {
	*StoreInMemory
	buf               map[string]string
	bufUsersShortURLs map[string][]string
}

// Commit finishes storage Transaction by commiting changes.
func (t *TransactionInMemory) Commit(_ context.Context) error {
	t.StoreInMemory.urls = t.buf
	t.StoreInMemory.usersShortUrls = t.bufUsersShortURLs
	t.lock.Unlock()
	return nil
}

// Rollback finishes storage Transaction by changes cancellation.
// Should be called in defer part in case of errors. It is safe to call it after Commit.
func (t *TransactionInMemory) Rollback(_ context.Context) error {
	t.lock.TryLock() // after commit it can be unlocked
	t.lock.Unlock()
	return nil
}

// GetByURLs performs query as part of Transaction.
//
// Reads map[string]*entity.ShortURL where key is full URL and value is *entity.ShortURL.
// If full URL not found it won't be presented in result.
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

// CheckIDs performs query as part of Transaction. Allows to check short IDs, returns taken IDs.
func (t *TransactionInMemory) CheckIDs(_ context.Context, IDs []string) ([]string, error) {
	result := make([]string, 0, len(IDs))
	for _, id := range IDs {
		if _, ok := t.buf[id]; ok {
			result = append(result, id)
		}
	}
	return result, nil
}

// AddURL performs query as part of Transaction. Adds new *entity.ShortURL.
func (t *TransactionInMemory) AddURL(_ context.Context, ShortURL *entity.ShortURL, force bool) error {
	oldURL, ok := t.buf[ShortURL.ID]
	if ok && !force {
		return entity.ErrUUIDTaken
	}
	for _, url := range t.buf {
		if url == ShortURL.URL && url != oldURL {
			return entity.ErrShortURLExists
		}
	}
	t.buf[ShortURL.ID] = ShortURL.URL
	return nil
}

// AddUserIfNotExists performs query as part of Transaction. Ensures user with UserUUID exists in storage.
func (t *TransactionInMemory) AddUserIfNotExists(_ context.Context, UserUUID string) error {
	if _, ok := t.bufUsersShortURLs[UserUUID]; !ok {
		t.bufUsersShortURLs[UserUUID] = []string{}
	}
	return nil
}

// AddUserURL performs query as part of Transaction. Links short URL ID to user
func (t *TransactionInMemory) AddUserURL(_ context.Context, ShortURLUUID, UserUUID string) error {
	if owns, ok := t.bufUsersShortURLs[UserUUID]; !ok {
		t.bufUsersShortURLs[UserUUID] = []string{ShortURLUUID}
	} else if !slices.Contains(owns, ShortURLUUID) {
		t.bufUsersShortURLs[UserUUID] = append(owns, ShortURLUUID)
	}
	return nil
}

// GetUserUrls performs query as part of Transaction.
// Returns map[string]*entity.ShortURL where key is full URL and value is *entity.ShortURL linked to user UUID.
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

// GetShortByURL performs query as part of Transaction. Reads *entity.ShortURL if exists for full URL.
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
