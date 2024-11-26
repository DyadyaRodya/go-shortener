package pgx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"

	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
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
	rows, err := t.tx.Query(ctx, `SELECT uuid, url FROM short_urls WHERE url = ANY($1) AND deleted_at IS NULL`, URLs)
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
	rows, err := t.tx.Query(ctx, `SELECT uuid FROM short_urls WHERE uuid = ANY($1) AND deleted_at IS NULL`, IDs)
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

func (t *TransactionPGX) AddURL(ctx context.Context, ShortURL *entity.ShortURL, force bool) error {
	t.logger.Debug("Adding URL", zap.Any("ShortURL", ShortURL))
	var query string
	if force {
		query = `INSERT INTO short_urls (uuid, url) VALUES (@uuid, @url) 
				     ON CONFLICT (uuid) DO UPDATE SET deleted_at = NULL, url = @url` // allow uuid rewriting
	} else {
		query = `INSERT INTO short_urls (uuid, url) VALUES (@uuid, @url)`
	}
	ct, err := t.tx.Exec(ctx, query,
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

func (t *TransactionPGX) AddUserIfNotExists(ctx context.Context, UserUUID string) error {
	t.logger.Debug("Adding User", zap.String("UserUUID", UserUUID))
	ct, err := t.tx.Exec(ctx, `INSERT INTO users (uuid) VALUES (@uuid) ON CONFLICT DO NOTHING`,
		pgx.NamedArgs{"uuid": UserUUID})

	if err != nil {
		t.logger.Info("Failed to insert into users",
			zap.String("uuid", UserUUID),
			zap.Error(err))
		return fmt.Errorf("TransactionPGX.AddUserIfNotExists: %w", err)
	}
	if !ct.Insert() {
		t.logger.Error("Failed to insert into users",
			zap.String("uuid", UserUUID),
			zap.Any("commandTag", ct))
		return errors.New("error in TransactionPGX.AddUserIfNotExists: not inserted uuid: " + UserUUID)
	}
	return nil
}

func (t *TransactionPGX) AddUserURL(ctx context.Context, ShortURLUUID, UserUUID string) error {
	t.logger.Debug("Adding User",
		zap.String("UserUUID", UserUUID),
		zap.String("ShortURLUUID", ShortURLUUID),
	)
	ct, err := t.tx.Exec(ctx, `INSERT INTO users_short_urls (user_uuid, short_url_uuid) VALUES 
                                                             (@user_uuid, @short_url_uuid) ON CONFLICT DO NOTHING`,
		pgx.NamedArgs{"user_uuid": UserUUID, "short_url_uuid": ShortURLUUID})

	if err != nil {
		t.logger.Info("Failed to insert into users_short_urls",
			zap.String("user_uuid", UserUUID),
			zap.String("short_url_uuid", ShortURLUUID),
			zap.Error(err))
		return fmt.Errorf("TransactionPGX.AddUserURL: %w", err)
	}
	if !ct.Insert() {
		t.logger.Error("Failed to insert into users_short_urls",
			zap.String("user_uuid", UserUUID),
			zap.String("short_url_uuid", ShortURLUUID),
			zap.Any("commandTag", ct))
		return errors.New("error in TransactionPGX.AddUserURL: not inserted user_uuid: " + UserUUID +
			" short_url_uuid: " + ShortURLUUID)
	}
	return nil
}

func (t *TransactionPGX) GetUserUrls(ctx context.Context, UserUUID string) (map[string]*entity.ShortURL, error) {
	rows, err := t.tx.Query(ctx, `SELECT uuid, url FROM short_urls AS su 
		JOIN users_short_urls AS usu ON su.uuid=usu.short_url_uuid 				
		WHERE usu.user_uuid = @user_uuid`,
		pgx.NamedArgs{"user_uuid": UserUUID})
	if err != nil {
		t.logger.Error("Failed to query database", zap.Error(err))
		return nil, fmt.Errorf("error in TransactionPGX.GetUserUrls: %w", err)
	}
	defer rows.Close()
	shortURLs := make(map[string]*entity.ShortURL)
	for rows.Next() {
		var uuid, url string
		if err := rows.Scan(&uuid, &url); err != nil {
			t.logger.Error("Failed to scan", zap.Error(err))
			return nil, fmt.Errorf("error in TransactionPGX.GetUserUrls: %w", err)
		}
		shortURLs[url] = &entity.ShortURL{URL: url, ID: uuid}
	}
	if err := rows.Err(); err != nil {
		t.logger.Error("Failed to query database", zap.Error(err))
		return nil, fmt.Errorf("error in TransactionPGX.GetUserUrls: %w", err)
	}
	return shortURLs, nil
}

func (t *TransactionPGX) GetShortByURL(ctx context.Context, URL string) (*entity.ShortURL, error) {
	t.logger.Debug("Getting Short URL by full URL", zap.String("URL", URL))

	var uuid string
	var deletedAt sql.NullTime
	err := t.tx.QueryRow(ctx, `SELECT uuid, deleted_at FROM short_urls WHERE url = @url`,
		pgx.NamedArgs{"url": URL},
	).Scan(&uuid, &deletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrShortURLNotFound
		}
		t.logger.Error("Failed to get Short URL by full URL", zap.String("URL", URL), zap.Error(err))
		return nil, err
	}

	shortURL := &entity.ShortURL{URL: URL, ID: uuid}
	if !deletedAt.Valid {
		return shortURL, nil
	}
	return shortURL, entity.ErrShortURLDeleted
}
