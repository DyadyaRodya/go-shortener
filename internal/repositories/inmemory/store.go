package inmemory

import (
	"context"
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	"github.com/DyadyaRodya/go-shortener/internal/usecases"
	"maps"
	"sync"
)

type StoreInMemory struct {
	storage map[string]string
	lock    sync.RWMutex
}

func NewStoreInMemory() *StoreInMemory {
	return &StoreInMemory{
		storage: make(map[string]string),
		lock:    sync.RWMutex{},
	}
}

func (s *StoreInMemory) Storage() *map[string]string {
	return &s.storage
}

func (s *StoreInMemory) AddURL(_ context.Context, ShortURL *entity.ShortURL) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.storage[ShortURL.ID]
	if ok {
		return entity.ErrUUIDTaken
	}
	for _, url := range s.storage {
		if url == ShortURL.URL {
			return entity.ErrShortURLExists
		}
	}
	s.storage[ShortURL.ID] = ShortURL.URL
	return nil
}

func (s *StoreInMemory) GetURLByID(_ context.Context, ID string) (*entity.ShortURL, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	url, ok := s.storage[ID]
	if !ok {
		return nil, entity.ErrShortURLNotFound
	}

	shortURL := &entity.ShortURL{
		ID:  ID,
		URL: url,
	}
	return shortURL, nil
}

func (s *StoreInMemory) GetShortByURL(_ context.Context, URL string) (*entity.ShortURL, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	for uuid, url := range s.storage {
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

func (s *StoreInMemory) Load(src map[string]string) {
	maps.Copy(s.storage, src)
}

func (s *StoreInMemory) Save(dst map[string]string) {
	maps.Copy(dst, s.storage)
}

func (s *StoreInMemory) TestConnection(_ context.Context) error {
	return nil
}

func (s *StoreInMemory) Begin(_ context.Context) (usecases.Transaction, error) {
	buf := make(map[string]string)
	s.lock.Lock()
	maps.Copy(buf, s.storage)
	tx := &TransactionInMemory{
		StoreInMemory: s, buf: buf,
	}
	return tx, nil
}
