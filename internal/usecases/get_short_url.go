package usecases

import (
	"context"

	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
)

func (u *Usecases) GetShortURL(ctx context.Context, ID string) (*entity.ShortURL, error) {
	shortURL, err := u.urlStorage.GetURLByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	return shortURL, nil
}
