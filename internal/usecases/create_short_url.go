package usecases

import (
	"errors"
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	"github.com/DyadyaRodya/go-shortener/internal/domain/services"
)

func (u *Usecases) CreateShortURL(url string) (*entity.ShortURL, error) {
	var id string
	var err error
	for {
		id, err = services.GenerateID()
		if err != nil {
			return nil, err
		}
		_, err = u.urlStorage.GetURLByID(id)
		if err != nil {
			if errors.Is(err, entity.ErrShortURLNotFound) {
				break
			}
			return nil, err
		}
	}

	shortURL := &entity.ShortURL{ID: id, URL: url}
	if err = u.urlStorage.AddURL(shortURL); err != nil {
		return nil, err
	}
	return shortURL, nil
}
