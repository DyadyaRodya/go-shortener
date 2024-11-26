package tests

import (
	"context"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/mock"

	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	"github.com/DyadyaRodya/go-shortener/internal/usecases/mocks"
)

func (u *usecasesSuite) TestUsecases_CreateShortURL_Success() {
	ctx := context.Background()
	fullURL := gofakeit.URL()
	foundIDs := []string{
		gofakeit.Word(), gofakeit.Word(), gofakeit.Word(),
	}
	generatedID := gofakeit.Word()
	shortURL := &entity.ShortURL{ID: generatedID, URL: fullURL}
	userUUID := gofakeit.UUID()

	for _, id := range foundIDs {
		u.idGenerator.EXPECT().Generate().Return(id, nil).Once()
		u.urlStorage.EXPECT().AddURL(ctx, mock.Anything, userUUID).RunAndReturn(func(callCTX context.Context, callURL *entity.ShortURL, userUUIDParam string) error {
			u.Equal(id, callURL.ID)
			u.Equal(userUUID, userUUIDParam)
			u.Equal(fullURL, callURL.URL)
			return entity.ErrUUIDTaken
		}).Once()
	}
	u.idGenerator.EXPECT().Generate().Return(generatedID, nil).Once()
	u.urlStorage.EXPECT().AddURL(ctx, mock.Anything, userUUID).RunAndReturn(func(ctx context.Context, shortURLParam *entity.ShortURL, userUUIDParam string) error {
		u.Equal(*shortURLParam, *shortURL)
		u.Equal(userUUID, userUUIDParam)
		return nil
	}).Once()

	result, err := u.usecases.CreateShortURL(ctx, fullURL, userUUID)

	u.NoError(err)
	u.Equal(result, shortURL)
}

func (u *usecasesSuite) TestUsecases_CreateShortURL_GenerateError() {
	ctx := context.Background()
	fullURL := gofakeit.URL()
	expectedError := gofakeit.Error()
	userUUID := gofakeit.UUID()

	u.idGenerator.EXPECT().Generate().Return("", expectedError).Once()

	result, err := u.usecases.CreateShortURL(ctx, fullURL, userUUID)

	u.ErrorIs(err, expectedError)
	u.Nil(result)
}

func (u *usecasesSuite) TestUsecases_CreateShortURL_AddURLError() {
	ctx := context.Background()
	fullURL := gofakeit.URL()
	generatedID := gofakeit.Word()
	shortURL := &entity.ShortURL{ID: generatedID, URL: fullURL}
	userUUID := gofakeit.UUID()
	expectedError := gofakeit.Error()

	u.idGenerator.EXPECT().Generate().Return(generatedID, nil).Once()
	u.urlStorage.EXPECT().AddURL(ctx, mock.Anything, userUUID).RunAndReturn(func(ctx context.Context, shortURLParam *entity.ShortURL, userUUIDParam string) error {
		u.Equal(*shortURLParam, *shortURL)
		u.Equal(userUUID, userUUIDParam)
		return expectedError
	}).Once()

	result, err := u.usecases.CreateShortURL(ctx, fullURL, userUUID)

	u.ErrorIs(err, expectedError)
	u.Nil(result)
}

func (u *usecasesSuite) TestUsecases_CreateShortURL_AddURL_ExistsError() {
	ctx := context.Background()
	fullURL := gofakeit.URL()
	generatedID := gofakeit.Word()
	shortURL := &entity.ShortURL{ID: generatedID, URL: fullURL}

	existedShortURL := &entity.ShortURL{ID: gofakeit.Word(), URL: fullURL}
	userUUID := gofakeit.UUID()

	expectedError := entity.ErrShortURLExists

	tx := mocks.NewTransaction(u.T())

	u.idGenerator.EXPECT().Generate().Return(generatedID, nil).Once()
	u.urlStorage.EXPECT().AddURL(ctx, mock.Anything, userUUID).RunAndReturn(func(ctx context.Context, shortURLParam *entity.ShortURL, userUUIDParam string) error {
		u.Equal(*shortURLParam, *shortURL)
		u.Equal(userUUID, userUUIDParam)
		return expectedError
	}).Once()

	u.urlStorage.EXPECT().Begin(ctx).Return(tx, nil).Once()

	tx.EXPECT().GetShortByURL(ctx, fullURL).Return(existedShortURL, nil).Once()
	tx.EXPECT().GetUserUrls(ctx, userUUID).Return(map[string]*entity.ShortURL{
		existedShortURL.ID: existedShortURL,
	}, nil).Once()

	// from defer
	tx.EXPECT().Rollback(ctx).Return(nil).Once()
	result, err := u.usecases.CreateShortURL(ctx, fullURL, userUUID)

	u.ErrorIs(err, expectedError)
	u.Equal(*existedShortURL, *result)
}

func (u *usecasesSuite) TestUsecases_CreateShortURL_AddURL_ExistsButNotOwns() {
	ctx := context.Background()
	fullURL := gofakeit.URL()
	generatedID := gofakeit.Word()
	shortURL := &entity.ShortURL{ID: generatedID, URL: fullURL}

	existedShortURL := &entity.ShortURL{ID: gofakeit.Word(), URL: fullURL}
	userUUID := gofakeit.UUID()

	expectedError := entity.ErrShortURLExists

	tx := mocks.NewTransaction(u.T())

	u.idGenerator.EXPECT().Generate().Return(generatedID, nil).Once()
	u.urlStorage.EXPECT().AddURL(ctx, mock.Anything, userUUID).RunAndReturn(func(ctx context.Context, shortURLParam *entity.ShortURL, userUUIDParam string) error {
		u.Equal(*shortURLParam, *shortURL)
		u.Equal(userUUID, userUUIDParam)
		return expectedError
	}).Once()

	u.urlStorage.EXPECT().Begin(ctx).Return(tx, nil).Once()

	tx.EXPECT().GetShortByURL(ctx, fullURL).Return(existedShortURL, nil).Once()
	tx.EXPECT().GetUserUrls(ctx, userUUID).Return(map[string]*entity.ShortURL{}, nil).Once()

	tx.EXPECT().AddUserIfNotExists(ctx, userUUID).Return(nil).Once()
	tx.EXPECT().AddUserURL(ctx, existedShortURL.ID, userUUID).Return(nil).Once()

	// finish transaction
	tx.EXPECT().Commit(ctx).Return(nil).Once()

	// from defer
	tx.EXPECT().Rollback(ctx).Return(nil).Once()
	result, err := u.usecases.CreateShortURL(ctx, fullURL, userUUID)

	u.ErrorIs(err, expectedError)
	u.Equal(*existedShortURL, *result)
}
