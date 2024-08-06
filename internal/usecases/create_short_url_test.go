package usecases

import (
	"errors"
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	"github.com/stretchr/testify/mock"
)

func (u *usecasesSuite) TestUsecases_CreateShortURL_Success() {
	fullURL := "http://test.url/blabla"
	foundIDs := []string{
		"teststr1", "teststr2", "teststr3",
	}
	generatedID := "teststring"
	shortURL := &entity.ShortURL{ID: generatedID, URL: fullURL}

	for _, id := range foundIDs {
		u.idGenerator.EXPECT().Generate().Return(id, nil).Once()
		u.urlStorage.EXPECT().GetURLByID(id).Return(&entity.ShortURL{}, nil).Once()
	}
	u.idGenerator.EXPECT().Generate().Return(generatedID, nil).Once()
	u.urlStorage.EXPECT().GetURLByID(generatedID).Return(nil, entity.ErrShortURLNotFound).Once()

	u.urlStorage.EXPECT().AddURL(mock.Anything).RunAndReturn(func(shortURLParam *entity.ShortURL) error {
		u.Equal(*shortURLParam, *shortURL)
		return nil
	}).Once()

	result, err := u.usecases.CreateShortURL(fullURL)

	u.NoError(err)
	u.Equal(result, shortURL)
}

func (u *usecasesSuite) TestUsecases_CreateShortURL_GenerateError() {
	fullURL := "http://test.url/blabla"
	expectedError := errors.New("error")

	u.idGenerator.EXPECT().Generate().Return("", expectedError).Once()

	result, err := u.usecases.CreateShortURL(fullURL)

	u.ErrorIs(err, expectedError)
	u.Nil(result)
}

func (u *usecasesSuite) TestUsecases_CreateShortURL_GetURLByIDError() {
	fullURL := "http://test.url/blabla"
	generatedID := "teststring"
	expectedError := errors.New("error")

	u.idGenerator.EXPECT().Generate().Return(generatedID, nil).Once()
	u.urlStorage.EXPECT().GetURLByID(generatedID).Return(nil, expectedError).Once()

	result, err := u.usecases.CreateShortURL(fullURL)

	u.ErrorIs(err, expectedError)
	u.Nil(result)
}

func (u *usecasesSuite) TestUsecases_CreateShortURL_AddURLError() {
	fullURL := "http://test.url/blabla"
	generatedID := "teststring"
	shortURL := &entity.ShortURL{ID: generatedID, URL: fullURL}
	expectedError := errors.New("error")

	u.idGenerator.EXPECT().Generate().Return(generatedID, nil).Once()
	u.urlStorage.EXPECT().GetURLByID(generatedID).Return(nil, entity.ErrShortURLNotFound).Once()

	u.urlStorage.EXPECT().AddURL(mock.Anything).RunAndReturn(func(shortURLParam *entity.ShortURL) error {
		u.Equal(*shortURLParam, *shortURL)
		return expectedError
	}).Once()

	result, err := u.usecases.CreateShortURL(fullURL)

	u.ErrorIs(err, expectedError)
	u.Nil(result)
}
