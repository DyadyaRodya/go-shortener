package inmemory

import (
	"context"
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	"github.com/DyadyaRodya/go-shortener/internal/usecases"
	"maps"
	"slices"
	"sync"
)

type StoreInMemory struct {
	urls           map[string]string
	usersShortUrls map[string][]string
	lock           sync.RWMutex
}

func NewStoreInMemory() *StoreInMemory {
	return &StoreInMemory{
		urls:           make(map[string]string),
		usersShortUrls: make(map[string][]string),
		lock:           sync.RWMutex{},
	}
}

func (s *StoreInMemory) URLS() *map[string]string {
	return &s.urls
}

func (s *StoreInMemory) UsersShortUrls() *map[string][]string {
	return &s.usersShortUrls
}

func (s *StoreInMemory) AddURL(_ context.Context, ShortURL *entity.ShortURL, OwnerUUID string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if owns, ok := s.usersShortUrls[OwnerUUID]; !ok {
		s.usersShortUrls[OwnerUUID] = []string{ShortURL.ID}
	} else if !slices.Contains(owns, ShortURL.ID) {
		s.usersShortUrls[OwnerUUID] = append(owns, ShortURL.ID)
	}

	_, ok := s.urls[ShortURL.ID]
	if ok {
		return entity.ErrUUIDTaken
	}
	for _, url := range s.urls {
		if url == ShortURL.URL {
			return entity.ErrShortURLExists
		}
	}
	s.urls[ShortURL.ID] = ShortURL.URL
	return nil
}

func (s *StoreInMemory) GetURLByID(_ context.Context, ID string) (*entity.ShortURL, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	url, ok := s.urls[ID]
	if !ok {
		return nil, entity.ErrShortURLNotFound
	}

	shortURL := &entity.ShortURL{
		ID:  ID,
		URL: url,
	}
	return shortURL, nil
}

func (s *StoreInMemory) TestConnection(_ context.Context) error {
	return nil
}

func (s *StoreInMemory) Begin(_ context.Context) (usecases.Transaction, error) {
	buf := make(map[string]string)
	bufOwns := make(map[string][]string)
	s.lock.Lock()
	maps.Copy(buf, s.urls)
	maps.Copy(bufOwns, s.usersShortUrls)
	tx := &TransactionInMemory{
		StoreInMemory: s, buf: buf, bufUsersShortURLs: bufOwns,
	}
	return tx, nil
}

func (s *StoreInMemory) GetUserUrls(_ context.Context, UserUUID string) (map[string]*entity.ShortURL, error) {
	result := make(map[string]*entity.ShortURL)

	owns, ok := s.usersShortUrls[UserUUID]
	if !ok {
		return result, nil
	}

	for _, uuid := range owns {
		result[uuid] = &entity.ShortURL{
			ID:  uuid,
			URL: s.urls[uuid],
		}
	}
	return result, nil
}
