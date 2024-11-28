package tests

import (
	"context"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
)

func (u *usecasesSuite) TestUsecases_GetUserShortURLs_Success() {
	ctx := context.Background()
	userUUID := gofakeit.UUID()
	id := gofakeit.Word()
	url := gofakeit.URL()
	shortURL := &entity.ShortURL{ID: id, URL: url}

	u.urlStorage.EXPECT().GetUserUrls(ctx, userUUID).Return(map[string]*entity.ShortURL{
		shortURL.ID: shortURL,
	}, nil).Once()

	result, err := u.usecases.GetUserShortURLs(ctx, userUUID)

	u.NoError(err)
	u.Equal([]*entity.ShortURL{shortURL}, result)
}

func (u *usecasesSuite) TestUsecases_GetUserShortURLs_Error() {
	ctx := context.Background()
	userUUID := gofakeit.UUID()
	expectedError := gofakeit.Error()

	u.urlStorage.EXPECT().GetUserUrls(ctx, userUUID).Return(nil, expectedError).Once()

	result, err := u.usecases.GetUserShortURLs(ctx, userUUID)

	u.ErrorIs(err, expectedError)
	u.Nil(result)
}
