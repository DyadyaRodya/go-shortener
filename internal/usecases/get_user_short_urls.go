package usecases

import (
	"context"
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
)

func (u *Usecases) GetUserShortURLs(ctx context.Context, userUUID string) ([]*entity.ShortURL, error) {
	shortURLs, err := u.urlStorage.GetUserUrls(ctx, userUUID)
	if err != nil {
		return nil, err
	}
	res := make([]*entity.ShortURL, 0, len(shortURLs))
	for _, shortURL := range shortURLs {
		res = append(res, shortURL)
	}
	return res, nil
}
