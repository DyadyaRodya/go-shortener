package usecases

import (
	"context"
	"errors"

	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
)

// CreateShortURL creates new short ID for URL and links it to user.
//
// Creates short ID for new URL. Existing URL are only linked to user.
// Existing URL already linked to user are also presented in return, but error marks that conflict.
func (u *Usecases) CreateShortURL(ctx context.Context, url string, UserUUID string) (*entity.ShortURL, error) {
	var id string
	var err error
	for {
		id, err = u.idGenerator.Generate()
		if err != nil {
			return nil, err
		}

		shortURL := &entity.ShortURL{ID: id, URL: url}
		if err = u.urlStorage.AddURL(ctx, shortURL, UserUUID); err != nil {
			switch {
			case errors.Is(err, entity.ErrUUIDTaken):
				continue // should generate another uuid
			case errors.Is(err, entity.ErrShortURLExists):
				return u.ensureUserLinkedToURL(ctx, url, UserUUID)
			default:
				return nil, err
			}
		}
		return shortURL, nil
	}
}

// returns *entity.ShortURL and entity.ErrShortURLExists if was not deleted
func (u *Usecases) ensureUserLinkedToURL(ctx context.Context, url, userUUID string) (*entity.ShortURL, error) {
	tx, err := u.urlStorage.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return nil, err
	}

	addCalled := false

	retErr := entity.ErrShortURLExists
	shortURL, err := tx.GetShortByURL(ctx, url)
	if err != nil {
		if errors.Is(err, entity.ErrShortURLDeleted) {
			err = tx.AddURL(ctx, shortURL, true) // allows uuid rewriting
			if err != nil {
				return nil, err
			}
			retErr = nil     // we recreated it
			addCalled = true // to not forget tx.Commit()
		} else {
			return nil, err
		}
	}

	userShortURLs, err := tx.GetUserUrls(ctx, userUUID)
	if err != nil {
		return nil, err
	}
	// check if user had url linked or userUUID not set (which means working without auth)
	if _, ok := userShortURLs[shortURL.ID]; ok || userUUID == "" {
		if addCalled {
			err = tx.Commit(ctx)
			if err != nil {
				return nil, err
			}
		}
		return shortURL, retErr
	}

	err = tx.AddUserIfNotExists(ctx, userUUID)
	if err != nil {
		return nil, err
	}
	err = tx.AddUserURL(ctx, shortURL.ID, userUUID)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return shortURL, retErr
}
