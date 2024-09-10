package inmemory

import (
	"context"
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
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

func (s *StoreInMemory) Load(src map[string]string) {
	maps.Copy(s.storage, src)
}

func (s *StoreInMemory) Save(dst map[string]string) {
	maps.Copy(dst, s.storage)
}

func (s *StoreInMemory) TestConnection(_ context.Context) error {
	return nil
}
