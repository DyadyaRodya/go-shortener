package tests

import (
	"context"
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	"github.com/brianvoe/gofakeit/v6"
)

func (u *usecasesSuite) TestUsecases_GetShortURL_Success() {
	ctx := context.Background()
	id := gofakeit.Word()
	shortURL := &entity.ShortURL{ID: id, URL: gofakeit.URL()}

	u.urlStorage.EXPECT().GetURLByID(ctx, id).Return(shortURL, nil).Once()

	result, err := u.usecases.GetShortURL(ctx, id)

	u.NoError(err)
	u.Equal(shortURL, result)
}

func (u *usecasesSuite) TestUsecases_GetShortURL_NotFound() {
	ctx := context.Background()
	id := gofakeit.Word()

	u.urlStorage.EXPECT().GetURLByID(ctx, id).Return(nil, entity.ErrShortURLNotFound).Once()

	result, err := u.usecases.GetShortURL(ctx, id)

	u.ErrorIs(err, entity.ErrShortURLNotFound)
	u.Empty(result)
}

func (u *usecasesSuite) TestUsecases_GetShortURL_AnyError() {
	ctx := context.Background()
	id := gofakeit.Word()
	otherError := gofakeit.Error()

	u.urlStorage.EXPECT().GetURLByID(ctx, id).Return(nil, otherError).Once()

	result, err := u.usecases.GetShortURL(ctx, id)

	u.ErrorIs(err, otherError)
	u.Empty(result)
}
