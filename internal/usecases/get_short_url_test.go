package usecases

import (
	"errors"
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
)

func (u *usecasesSuite) TestUsecases_GetShortURL_Success() {
	id := "teststring"
	shortURL := &entity.ShortURL{ID: id, URL: "http://test.url/blabla"}

	u.urlStorage.EXPECT().GetURLByID(id).Return(shortURL, nil).Once()

	result, err := u.usecases.GetShortURL(id)

	u.NoError(err)
	u.Equal(shortURL, result)
}

func (u *usecasesSuite) TestUsecases_GetShortURL_NotFound() {
	id := "teststring"

	u.urlStorage.EXPECT().GetURLByID(id).Return(nil, entity.ErrShortURLNotFound).Once()

	result, err := u.usecases.GetShortURL(id)

	u.ErrorIs(err, entity.ErrShortURLNotFound)
	u.Empty(result)
}

func (u *usecasesSuite) TestUsecases_GetShortURL_AnyError() {
	id := "teststring"
	otherError := errors.New("some other error")

	u.urlStorage.EXPECT().GetURLByID(id).Return(nil, otherError).Once()

	result, err := u.usecases.GetShortURL(id)

	u.ErrorIs(err, otherError)
	u.Empty(result)
}
