package tests

import (
	"context"
	"errors"
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	"github.com/stretchr/testify/mock"
)

func (u *usecasesSuite) TestUsecases_CreateShortURL_Success() {
	ctx := context.Background()
	fullURL := "http://test.url/blabla"
	foundIDs := []string{
		"teststr1", "teststr2", "teststr3",
	}
	generatedID := "teststring"
	shortURL := &entity.ShortURL{ID: generatedID, URL: fullURL}

	for _, id := range foundIDs {
		u.idGenerator.EXPECT().Generate().Return(id, nil).Once()
		u.urlStorage.EXPECT().AddURL(ctx, mock.Anything).RunAndReturn(func(callCTX context.Context, callURL *entity.ShortURL) error {
			u.Equal(id, callURL.ID)
			u.Equal(fullURL, callURL.URL)
			return entity.ErrUUIDTaken
		}).Once()
	}
	u.idGenerator.EXPECT().Generate().Return(generatedID, nil).Once()
	u.urlStorage.EXPECT().AddURL(ctx, mock.Anything).RunAndReturn(func(ctx context.Context, shortURLParam *entity.ShortURL) error {
		u.Equal(*shortURLParam, *shortURL)
		return nil
	}).Once()

	result, err := u.usecases.CreateShortURL(ctx, fullURL)

	u.NoError(err)
	u.Equal(result, shortURL)
}

func (u *usecasesSuite) TestUsecases_CreateShortURL_GenerateError() {
	ctx := context.Background()
	fullURL := "http://test.url/blabla"
	expectedError := errors.New("error")

	u.idGenerator.EXPECT().Generate().Return("", expectedError).Once()

	result, err := u.usecases.CreateShortURL(ctx, fullURL)

	u.ErrorIs(err, expectedError)
	u.Nil(result)
}

func (u *usecasesSuite) TestUsecases_CreateShortURL_AddURLError() {
	ctx := context.Background()
	fullURL := "http://test.url/blabla"
	generatedID := "teststring"
	shortURL := &entity.ShortURL{ID: generatedID, URL: fullURL}
	expectedError := errors.New("error")

	u.idGenerator.EXPECT().Generate().Return(generatedID, nil).Once()
	u.urlStorage.EXPECT().AddURL(ctx, mock.Anything).RunAndReturn(func(ctx context.Context, shortURLParam *entity.ShortURL) error {
		u.Equal(*shortURLParam, *shortURL)
		return expectedError
	}).Once()

	result, err := u.usecases.CreateShortURL(ctx, fullURL)

	u.ErrorIs(err, expectedError)
	u.Nil(result)
}

func (u *usecasesSuite) TestUsecases_CreateShortURL_AddURL_ExistsError() {
	ctx := context.Background()
	fullURL := "http://test.url/blabla"
	generatedID := "teststring"
	shortURL := &entity.ShortURL{ID: generatedID, URL: fullURL}

	expectedShortURL := &entity.ShortURL{ID: "expected", URL: fullURL}
	expectedError := entity.ErrShortURLExists

	u.idGenerator.EXPECT().Generate().Return(generatedID, nil).Once()
	u.urlStorage.EXPECT().AddURL(ctx, mock.Anything).RunAndReturn(func(ctx context.Context, shortURLParam *entity.ShortURL) error {
		u.Equal(*shortURLParam, *shortURL)
		return expectedError
	}).Once()

	u.urlStorage.EXPECT().GetShortByURL(ctx, fullURL).Return(expectedShortURL, nil).Once()

	result, err := u.usecases.CreateShortURL(ctx, fullURL)

	u.ErrorIs(err, expectedError)
	u.Equal(*expectedShortURL, *result)
}
