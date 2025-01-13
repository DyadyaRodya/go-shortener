package usecases

import (
	"context"

	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	"github.com/DyadyaRodya/go-shortener/internal/usecases/dto"
)

// Usecases interfaces
type (
	// Transaction interface of storage transaction (session) for Usecases.
	Transaction interface {
		// Commit finishes storage Transaction by commiting changes.
		Commit(ctx context.Context) error
		// Rollback finishes storage Transaction by changes cancellation.
		// Should be called in defer part in case of errors. It is safe to call it after Commit.
		Rollback(ctx context.Context) error
		// GetShortByURL performs query as part of Transaction. Reads *entity.ShortURL if exists for full URL.
		GetShortByURL(ctx context.Context, URL string) (*entity.ShortURL, error)
		// GetByURLs performs query as part of Transaction.
		//
		// Reads map[string]*entity.ShortURL where key is full URL and value is *entity.ShortURL.
		// If full URL not found it won't be presented in result.
		GetByURLs(ctx context.Context, URLs []string) (map[string]*entity.ShortURL, error)
		// CheckIDs performs query as part of Transaction. Allows to check short IDs, returns taken IDs.
		CheckIDs(ctx context.Context, IDs []string) ([]string, error)
		// AddURL performs query as part of Transaction. Adds new *entity.ShortURL.
		AddURL(ctx context.Context, ShortURL *entity.ShortURL, force bool) error
		// AddUserIfNotExists performs query as part of Transaction. Ensures user with UserUUID exists in storage.
		AddUserIfNotExists(ctx context.Context, UserUUID string) error
		// AddUserURL performs query as part of Transaction. Links short URL ID to user
		AddUserURL(ctx context.Context, ShortURLUUID, UserUUID string) error
		// GetUserUrls performs query as part of Transaction.
		// Returns map[string]*entity.ShortURL where key is full URL and value is *entity.ShortURL linked to user UUID.
		GetUserUrls(ctx context.Context, UserUUID string) (map[string]*entity.ShortURL, error)
	}
	// URLStorage interface of storage for Usecases.
	// Allows to start new Transaction or performs a transaction in one database call.
	URLStorage interface {
		// TestConnection allows to ping DB end check the readiness of the service
		TestConnection(ctx context.Context) error
		// GetURLByID performs a query transaction in one database call. Returns *entity.ShortURL for ID
		GetURLByID(ctx context.Context, ID string) (*entity.ShortURL, error)
		// AddURL performs a query transaction in one database call. Adds new user *entity.ShortURL and links it to user
		AddURL(ctx context.Context, ShortURL *entity.ShortURL, OwnerUUID string) error
		// GetUserUrls performs a query transaction in one database call.
		// Returns map[string]*entity.ShortURL where key is full URL and value is *entity.ShortURL linked to user UUID.
		GetUserUrls(ctx context.Context, UserUUID string) (map[string]*entity.ShortURL, error)
		// DeleteUserURLs performs a query transaction in one database call. Unlinks URLs and users.
		//
		// Removes URLs that are not linked to any other users.
		DeleteUserURLs(ctx context.Context, requests ...*dto.DeleteUserShortURLsRequest) error

		// GetStats returns summary *dto.StatsResponse with total numbers of shortened URLs and users stored
		GetStats(ctx context.Context) (*dto.StatsResponse, error)

		// Begin starts new Transaction session for Storage
		Begin(ctx context.Context) (Transaction, error)
	}
	// IDGenerator interface for Usecases to remove the dependence of business logic
	// on libraries and simplify testing by generating mocks.
	IDGenerator interface {
		Generate() (string, error)
	}
)

// Usecases contains main business scripts of the service.
// Does not depend on storages, UI, API, libraries.
type Usecases struct {
	urlStorage  URLStorage
	idGenerator IDGenerator
}

// NewUsecases constructor for Usecases
func NewUsecases(URLStorage URLStorage, IDGenerator IDGenerator) *Usecases {
	return &Usecases{
		urlStorage:  URLStorage,
		idGenerator: IDGenerator,
	}
}
