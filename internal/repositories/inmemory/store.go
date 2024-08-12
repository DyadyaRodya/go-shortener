package inmemory

import (
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
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

func (s *StoreInMemory) AddURL(ShortURL *entity.ShortURL) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.storage[ShortURL.ID] = ShortURL.URL
	return nil
}

func (s *StoreInMemory) GetURLByID(ID string) (*entity.ShortURL, error) {
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
