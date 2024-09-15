package pgx

import (
	"context"
	"errors"
	"fmt"
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"strings"
)

type TransactionPGX struct {
	tx     pgx.Tx
	logger *zap.Logger
}

func (t *TransactionPGX) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *TransactionPGX) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

func (t *TransactionPGX) GetByURLs(ctx context.Context, URLs []string) (map[string]*entity.ShortURL, error) {
	rows, err := t.tx.Query(ctx, `SELECT uuid, url FROM short_urls WHERE url = ANY($1)`, URLs)
	if err != nil {
		t.logger.Error("Failed to query database", zap.Error(err))
		return nil, fmt.Errorf("error in TransactionPGX.GetByURLs: %w", err)
	}
	defer rows.Close()
	shortURLs := make(map[string]*entity.ShortURL)
	for rows.Next() {
		var uuid, url string
		if err := rows.Scan(&uuid, &url); err != nil {
			t.logger.Error("Failed to scan", zap.Error(err))
			return nil, fmt.Errorf("error in TransactionPGX.GetByURLs: %w", err)
		}
		shortURLs[url] = &entity.ShortURL{URL: url, ID: uuid}
	}
	if err := rows.Err(); err != nil {
		t.logger.Error("Failed to query database", zap.Error(err))
		return nil, fmt.Errorf("error in TransactionPGX.GetByURLs: %w", err)
	}
	return shortURLs, nil
}

func (t *TransactionPGX) CheckIDs(ctx context.Context, IDs []string) ([]string, error) {
	rows, err := t.tx.Query(ctx, `SELECT uuid FROM short_urls WHERE uuid = ANY($1)`, IDs)
	if err != nil {
		t.logger.Error("Failed to query database", zap.Error(err))
		return nil, fmt.Errorf("error in TransactionPGX.GetByIDs: %w", err)
	}
	defer rows.Close()
	var takenIDs []string
	for rows.Next() {
		var uuid string
		if err := rows.Scan(&uuid); err != nil {
			t.logger.Error("Failed to scan", zap.Error(err))
			return nil, fmt.Errorf("error in TransactionPGX.GetByIDs: %w", err)
		}
		takenIDs = append(takenIDs, uuid)
	}
	if err := rows.Err(); err != nil {
		t.logger.Error("Failed to query database", zap.Error(err))
		return nil, fmt.Errorf("error in TransactionPGX.GetByIDs: %w", err)
	}
	return takenIDs, nil
}

func (t *TransactionPGX) AddURL(ctx context.Context, ShortURL *entity.ShortURL) error {
	t.logger.Debug("Adding URL", zap.Any("ShortURL", ShortURL))
	ct, err := t.tx.Exec(ctx, `INSERT INTO short_urls (uuid, url) VALUES (@uuid, @url)`,
		pgx.NamedArgs{"uuid": ShortURL.ID, "url": ShortURL.URL})

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		switch {
		case strings.Contains(pgErr.ConstraintName, "short_urls_pkey"):
			return entity.ErrUUIDTaken
		case strings.Contains(pgErr.ConstraintName, "short_urls_url_key"):
			return entity.ErrShortURLExists
		}
	}

	if err != nil {
		t.logger.Info("Failed to insert into short_urls",
			zap.String("uuid", ShortURL.ID),
			zap.String("url", ShortURL.URL),
			zap.Error(err))
		return fmt.Errorf("TransactionPGX.AddURL: %w", err)
	}
	if !ct.Insert() {
		t.logger.Error("Failed to insert into short_urls",
			zap.String("uuid", ShortURL.ID),
			zap.String("url", ShortURL.URL),
			zap.Any("commandTag", ct))
		return errors.New("error in TransactionPGX.AddURL: not inserted url: " +
			ShortURL.URL + " uuid: " + ShortURL.ID)
	}
	return nil
}
