package usecases

import (
	"context"
	"errors"
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
)

func (u *Usecases) CreateShortURL(ctx context.Context, url string) (*entity.ShortURL, error) {
	var id string
	var err error
	for {
		id, err = u.idGenerator.Generate()
		if err != nil {
			return nil, err
		}
		_, err = u.urlStorage.GetURLByID(ctx, id)
		if err != nil {
			if errors.Is(err, entity.ErrShortURLNotFound) {
				break
			}
			return nil, err
		}
	}

	shortURL := &entity.ShortURL{ID: id, URL: url}
	if err = u.urlStorage.AddURL(ctx, shortURL); err != nil {
		return nil, err
	}
	return shortURL, nil
}
