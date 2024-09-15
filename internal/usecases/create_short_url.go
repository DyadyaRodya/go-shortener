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

		shortURL := &entity.ShortURL{ID: id, URL: url}
		if err = u.urlStorage.AddURL(ctx, shortURL); err != nil {
			switch {
			case errors.Is(err, entity.ErrUUIDTaken):
				continue // should generate another uuid
			case errors.Is(err, entity.ErrShortURLExists):
				shortURL, err = u.urlStorage.GetShortByURL(ctx, url)
				if err != nil {
					return nil, err
				}
				return shortURL, entity.ErrShortURLExists
			default:
				return nil, err
			}
		}
		return shortURL, nil
	}
}
