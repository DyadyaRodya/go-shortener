package usecases

import (
	"context"
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	"github.com/DyadyaRodya/go-shortener/internal/usecases/dto"
	"slices"
)

func (u *Usecases) BatchCreateShortURLs(
	ctx context.Context,
	createRequests []*dto.BatchCreateRequest,
) ([]*dto.BatchCreateResponse, error) {
	if len(createRequests) == 0 {
		return []*dto.BatchCreateResponse{}, nil
	}

	originalURLs := make([]string, 0, len(createRequests))
	for _, createRequest := range createRequests {
		originalURLs = append(originalURLs, createRequest.OriginalURL)
	}

	// start transaction
	tx, err := u.urlStorage.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// check existing
	existingShortURLs, err := tx.GetByURLs(ctx, originalURLs)
	if err != nil {
		return nil, err
	}

	// generate new uuids for not existing urls
	newIDs := make([]string, 0, len(createRequests)-len(existingShortURLs))
	for _, createRequest := range createRequests {
		if _, ok := existingShortURLs[createRequest.OriginalURL]; !ok {
			id, err := u.idGenerator.Generate()
			if err != nil {
				return nil, err
			}
			newIDs = append(newIDs, id)
		}
	}

	// ensure generated unique uuids
	newIDs, err = u.ensureUniqueIDs(ctx, tx, newIDs)
	if err != nil {
		return nil, err
	}

	// creating new short urls if not existed
	resp := make([]*dto.BatchCreateResponse, 0, len(createRequests))
	for _, createRequest := range createRequests {
		shortURL, ok := existingShortURLs[createRequest.OriginalURL]
		if !ok {
			var id string
			id, newIDs = newIDs[0], newIDs[1:]
			shortURL = &entity.ShortURL{
				ID: id, URL: createRequest.OriginalURL,
			}
			err = tx.AddURL(ctx, shortURL)
			if err != nil {
				return nil, err
			}
		}
		resp = append(resp, &dto.BatchCreateResponse{
			CorrelationID: createRequest.CorrelationID,
			ShortURL:      shortURL,
		})
	}

	// finish transaction
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (u *Usecases) ensureUniqueIDs(ctx context.Context, tx Transaction, newIDs []string) ([]string, error) {
	for {
		existingIDs, err := tx.CheckIDs(ctx, newIDs)
		if err != nil {
			return nil, err
		}
		if len(existingIDs) == 0 {
			break
		}
		for i, id := range newIDs {
			if slices.Contains(existingIDs, id) {
				newIDs[i], err = u.idGenerator.Generate()
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return newIDs, nil
}
