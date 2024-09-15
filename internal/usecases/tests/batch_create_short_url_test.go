package tests

import (
	"context"
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	"github.com/DyadyaRodya/go-shortener/internal/usecases/dto"
	"github.com/DyadyaRodya/go-shortener/internal/usecases/mocks"
)

func (u *usecasesSuite) TestUsecases_BatchCreateShortURLs_Success() {
	ctx := context.Background()
	fullURL := "http://test.url/blabla"
	existedShortURL := &entity.ShortURL{
		ID:  "idexists",
		URL: "https://exists.url/blabla",
	}
	foundIDs := []string{
		"teststr1", "teststr2", "teststr3",
	}
	generatedID := "teststring"
	shortURL := &entity.ShortURL{ID: generatedID, URL: fullURL}

	requests := []*dto.BatchCreateRequest{
		&dto.BatchCreateRequest{
			CorrelationID: "1",
			OriginalURL:   existedShortURL.URL,
		},
		&dto.BatchCreateRequest{
			CorrelationID: "2",
			OriginalURL:   fullURL,
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
	}

	tx := mocks.NewTransaction(u.T())

	// start transaction
	u.urlStorage.EXPECT().Begin(ctx).Return(tx, nil).Once()

	// check existing
	tx.EXPECT().GetByURLs(ctx, []string{existedShortURL.URL, fullURL}).Return(map[string]*entity.ShortURL{
		existedShortURL.URL: existedShortURL,
	}, nil).Once()

	// ensure generated unique uuids
	for _, id := range foundIDs {
		u.idGenerator.EXPECT().Generate().Return(id, nil).Once()
		tx.EXPECT().CheckIDs(ctx, []string{id}).Return([]string{id}, nil).Once()
	}
	u.idGenerator.EXPECT().Generate().Return(generatedID, nil).Once()
	tx.EXPECT().CheckIDs(ctx, []string{generatedID}).Return([]string{}, nil).Once()

	// created new short url
	tx.EXPECT().AddURL(ctx, shortURL).Return(nil).Once()

	// finish transaction
	tx.EXPECT().Commit(ctx).Return(nil).Once()

	// from defer
	tx.EXPECT().Rollback(ctx).Return(nil).Once()

	// asserts
	result, err := u.usecases.BatchCreateShortURLs(ctx, requests)
	u.Nil(err)
	u.Equal(expected, result)
}

func (u *usecasesSuite) TestUsecases_BatchCreateShortURLs_EmptyRequest() {
	ctx := context.Background()
	requests := make([]*dto.BatchCreateRequest, 0)
	expected := make([]*dto.BatchCreateResponse, 0)

	result, err := u.usecases.BatchCreateShortURLs(ctx, requests)
	u.Nil(err)
	u.Equal(expected, result)
}
