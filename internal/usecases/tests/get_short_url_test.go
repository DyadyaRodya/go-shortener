package tests

import (
	"context"
	"errors"
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
)

func (u *usecasesSuite) TestUsecases_GetShortURL_Success() {
	ctx := context.Background()
	id := "teststring"
	shortURL := &entity.ShortURL{ID: id, URL: "http://test.url/blabla"}

	u.urlStorage.EXPECT().GetURLByID(ctx, id).Return(shortURL, nil).Once()

	result, err := u.usecases.GetShortURL(ctx, id)

	u.NoError(err)
	u.Equal(shortURL, result)
}

func (u *usecasesSuite) TestUsecases_GetShortURL_NotFound() {
	ctx := context.Background()
	id := "teststring"

	u.urlStorage.EXPECT().GetURLByID(ctx, id).Return(nil, entity.ErrShortURLNotFound).Once()

	result, err := u.usecases.GetShortURL(ctx, id)

	u.ErrorIs(err, entity.ErrShortURLNotFound)
	u.Empty(result)
}

func (u *usecasesSuite) TestUsecases_GetShortURL_AnyError() {
	ctx := context.Background()
	id := "teststring"
	otherError := errors.New("some other error")

	u.urlStorage.EXPECT().GetURLByID(ctx, id).Return(nil, otherError).Once()

	result, err := u.usecases.GetShortURL(ctx, id)

	u.ErrorIs(err, otherError)
	u.Empty(result)
}
