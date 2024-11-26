package tests

import (
	"context"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	"github.com/DyadyaRodya/go-shortener/internal/usecases/dto"
	"github.com/DyadyaRodya/go-shortener/internal/usecases/mocks"
)

func (u *usecasesSuite) TestUsecases_BatchCreateShortURLs_Success() {
	ctx := context.Background()
	fullURL := gofakeit.URL()
	existedShortURL := &entity.ShortURL{
		ID:  gofakeit.Word(),
		URL: gofakeit.URL(),
	}
	existedUserShortURL := &entity.ShortURL{
		ID:  gofakeit.Word(),
		URL: gofakeit.URL(),
	}
	foundIDs := []string{
		gofakeit.Word(), gofakeit.Word(), gofakeit.Word(),
	}
	generatedID := gofakeit.Word()
	shortURL := &entity.ShortURL{ID: generatedID, URL: fullURL}

	userUUID := gofakeit.UUID()

	requests := []*dto.BatchCreateRequest{
		&dto.BatchCreateRequest{
			CorrelationID: "1",
			OriginalURL:   existedShortURL.URL,
		},
		&dto.BatchCreateRequest{
			CorrelationID: "2",
			OriginalURL:   fullURL,
		},
		&dto.BatchCreateRequest{
			CorrelationID: "3",
			OriginalURL:   existedUserShortURL.URL,
		},
	}

	expected := []*dto.BatchCreateResponse{
		&dto.BatchCreateResponse{
			CorrelationID: "1",
			ShortURL:      existedShortURL,
		},
		&dto.BatchCreateResponse{
			CorrelationID: "2",
			ShortURL:      shortURL,
		},
		&dto.BatchCreateResponse{
			CorrelationID: "3",
			ShortURL:      existedUserShortURL,
		},
	}

	tx := mocks.NewTransaction(u.T())

	// start transaction
	u.urlStorage.EXPECT().Begin(ctx).Return(tx, nil).Once()

	// ensure user exists
	tx.EXPECT().AddUserIfNotExists(ctx, userUUID).Return(nil).Once()

	// check existing
	tx.EXPECT().GetByURLs(ctx, []string{existedShortURL.URL, fullURL, existedUserShortURL.URL}).Return(map[string]*entity.ShortURL{
		existedShortURL.URL: existedShortURL, existedUserShortURL.URL: existedUserShortURL,
	}, nil).Once()

	// get urls owned by user
	tx.EXPECT().GetUserUrls(ctx, userUUID).Return(map[string]*entity.ShortURL{
		existedUserShortURL.URL: existedUserShortURL,
	}, nil).Once()

	// ensure generated unique uuids
	for _, id := range foundIDs {
		u.idGenerator.EXPECT().Generate().Return(id, nil).Once()
		tx.EXPECT().CheckIDs(ctx, []string{id}).Return([]string{id}, nil).Once()
	}
	u.idGenerator.EXPECT().Generate().Return(generatedID, nil).Once()
	tx.EXPECT().CheckIDs(ctx, []string{generatedID}).Return([]string{}, nil).Once()

	// link existed to user
	tx.EXPECT().AddUserURL(ctx, existedShortURL.ID, userUUID).Return(nil).Once()

	// created new short url
	tx.EXPECT().AddURL(ctx, shortURL, true).Return(nil).Once()

	// link to user
	tx.EXPECT().AddUserURL(ctx, shortURL.ID, userUUID).Return(nil).Once()

	// finish transaction
	tx.EXPECT().Commit(ctx).Return(nil).Once()

	// from defer
	tx.EXPECT().Rollback(ctx).Return(nil).Once()

	// asserts
	result, err := u.usecases.BatchCreateShortURLs(ctx, requests, userUUID)
	u.Nil(err)
	u.Equal(expected, result)
}

func (u *usecasesSuite) TestUsecases_BatchCreateShortURLs_EmptyRequest() {
	ctx := context.Background()
	userUUID := gofakeit.UUID()
	requests := make([]*dto.BatchCreateRequest, 0)
	expected := make([]*dto.BatchCreateResponse, 0)

	result, err := u.usecases.BatchCreateShortURLs(ctx, requests, userUUID)
	u.Nil(err)
	u.Equal(expected, result)
}
