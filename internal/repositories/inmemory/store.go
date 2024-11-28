package inmemory

import (
	"context"
	"maps"
	"slices"
	"sync"

	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	"github.com/DyadyaRodya/go-shortener/internal/usecases"
	"github.com/DyadyaRodya/go-shortener/internal/usecases/dto"
	"github.com/DyadyaRodya/go-shortener/pkg/itemsets"
)

// StoreInMemory stores urls and users in process RAM
type StoreInMemory struct {
	urls           map[string]string
	usersShortUrls map[string][]string
	lock           sync.RWMutex
}

// NewStoreInMemory constructor for StoreInMemory
func NewStoreInMemory() *StoreInMemory {
	return &StoreInMemory{
		urls:           make(map[string]string),
		usersShortUrls: make(map[string][]string),
		lock:           sync.RWMutex{},
	}
}

// URLS get all pairs of ID and full URL
func (s *StoreInMemory) URLS() *map[string]string {
	return &s.urls
}

// UsersShortUrls get lists of urls by users
func (s *StoreInMemory) UsersShortUrls() *map[string][]string {
	return &s.usersShortUrls
}

// AddURL performs a query transaction in one database call. Adds new user *entity.ShortURL and links it to user
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

// GetURLByID performs a query transaction in one database call. Returns *entity.ShortURL for ID
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

// TestConnection allows to ping DB end check the readiness of the service
func (s *StoreInMemory) TestConnection(_ context.Context) error {
	return nil
}

// Begin starts new TransactionInMemory session for StoreInMemory
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

// GetUserUrls performs a query transaction in one database call.
// Returns map[string]*entity.ShortURL where key is full URL and value is *entity.ShortURL linked to user UUID.
func (s *StoreInMemory) GetUserUrls(_ context.Context, UserUUID string) (map[string]*entity.ShortURL, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

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

// DeleteUserURLs performs a query transaction in one database call. Unlinks URLs and users.
//
// Removes URLs that are not linked to any other users.
func (s *StoreInMemory) DeleteUserURLs(_ context.Context, requests ...*dto.DeleteUserShortURLsRequest) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	toDelete := make([]string, 0)

	for _, request := range requests {
		if request.UserUUID == "" {
			continue
		}
		if owns, ok := s.usersShortUrls[request.UserUUID]; ok {
			inter := itemsets.Intersection(owns, request.ShortURLUUIDs)
			s.usersShortUrls[request.UserUUID] = itemsets.RemoveItems(owns, inter)
			toDelete = itemsets.AddItems(toDelete, inter)
		}
	}

	for _, uuid := range toDelete {
		found := false
		for _, owns := range s.usersShortUrls {
			if slices.Contains(owns, uuid) {
				found = true
				break
			}
		}
		if !found {
			delete(s.urls, uuid) // do not support deletion mark for memory and file
		}
	}
	return nil
}
