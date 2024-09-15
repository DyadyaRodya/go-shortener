package pgx

import (
	"context"
	"errors"
	"fmt"
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	"github.com/DyadyaRodya/go-shortener/internal/usecases"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type StorePGX struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewStorePGX(pool *pgxpool.Pool, logger *zap.Logger) *StorePGX {
	return &StorePGX{pool: pool, logger: logger}
}

func (s *StorePGX) InitSchema(ctx context.Context) error {
	s.logger.Info("Initializing schema")
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.logger.Error("Failed to start transaction", zap.Error(err))
		return err
	}
	defer tx.Rollback(ctx)

	s.logger.Info("Creating table `short_urls`")
	_, err = tx.Exec(ctx, `CREATE TABLE IF NOT EXISTS short_urls(uuid CHAR(8) NOT NULL PRIMARY KEY, url TEXT UNIQUE NOT NULL)`)
	if err != nil {
		s.logger.Error("Failed to create table `short_urls`", zap.Error(err))
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		s.logger.Error("Failed to commit transaction", zap.Error(err))
		return err
	}
	s.logger.Info("Initializing schema done")
	return nil
}

func (s *StorePGX) Load(ctx context.Context, src map[string]string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.logger.Error("Failed to start transaction", zap.Error(err))
		return err
	}
	defer tx.Rollback(ctx)

	sqlQuery := `INSERT INTO short_urls (uuid, url) VALUES (@uuid, @url) ON CONFLICT (uuid) DO NOTHING`

	for uuid, url := range src {
		ct, err := tx.Exec(ctx, sqlQuery, pgx.NamedArgs{"uuid": uuid, "url": url})

		if err != nil {
			s.logger.Error("Failed to insert into short_urls",
				zap.String("uuid", uuid),
				zap.String("url", url),
				zap.Error(err))
			return fmt.Errorf("error in StorePGX.Load for url: %s uuid %s: %w", url, uuid, err)
		}
		if !ct.Insert() {
			s.logger.Error("Failed to insert into short_urls",
				zap.String("uuid", uuid),
				zap.String("url", url),
				zap.Any("commandTag", ct))
			return errors.New("error in StorePGX.Load: not inserted url: " + url + " uuid: " + uuid)
		}
	}
	err = tx.Commit(ctx)
	return err
}

func (s *StorePGX) Save(ctx context.Context) (*map[string]string, error) {
	dst := make(map[string]string)

	rows, err := s.pool.Query(ctx, `SELECT uuid, url FROM short_urls`)
	if err != nil {
		s.logger.Error("Failed to query database", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var uuid, url string
		if err := rows.Scan(&uuid, &url); err != nil {
			s.logger.Error("Failed to scan", zap.Error(err))
			return nil, err
		}
		dst[uuid] = url
	}
	if err := rows.Err(); err != nil {
		s.logger.Error("Failed to query database", zap.Error(err))
		return nil, err
	}
	s.logger.Info("Successfully read database")
	return &dst, nil
}

func (s *StorePGX) TestConnection(ctx context.Context) error {
	s.logger.Debug("Pinging database")
	return s.pool.Ping(ctx)
}

func (s *StorePGX) AddURL(ctx context.Context, ShortURL *entity.ShortURL) error {
	s.logger.Debug("Adding URL", zap.Any("ShortURL", ShortURL))
	ct, err := s.pool.Exec(ctx, `INSERT INTO short_urls (uuid, url) VALUES (@uuid, @url)`, pgx.NamedArgs{"uuid": ShortURL.ID, "url": ShortURL.URL})

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		switch pgErr.ColumnName {
		case "uuid":
			return entity.ErrUUIDTaken
		case "url":
			return entity.ErrShortURLExists
		}
	}

	if err != nil {
		s.logger.Info("Failed to insert into short_urls",
			zap.String("uuid", ShortURL.ID),
			zap.String("url", ShortURL.URL),
			zap.Error(err))
		return err
	}
	if !ct.Insert() {
		s.logger.Error("Failed to insert into short_urls",
			zap.String("uuid", ShortURL.ID),
			zap.String("url", ShortURL.URL),
			zap.Any("commandTag", ct))
		return errors.New("error in StorePGX.AddURL: not inserted url: " + ShortURL.URL + " uuid: " + ShortURL.ID)
	}
	return nil
}

func (s *StorePGX) GetURLByID(ctx context.Context, ID string) (*entity.ShortURL, error) {
	s.logger.Debug("Getting URL by ID", zap.String("ID", ID))

	var url string
	err := s.pool.QueryRow(ctx, `SELECT url FROM short_urls WHERE uuid = @uuid`, pgx.NamedArgs{"uuid": ID}).
		Scan(&url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrShortURLNotFound
		}
		s.logger.Error("Failed to get URL by ID", zap.String("ID", ID), zap.Error(err))
		return nil, err
	}
	return &entity.ShortURL{URL: url, ID: ID}, nil
}

func (s *StorePGX) Begin(ctx context.Context) (usecases.Transaction, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.logger.Error("Failed to start transaction", zap.Error(err))
		return nil, fmt.Errorf("error in StorePGX.Begin: %w", err)
	}
	txPGX := &TransactionPGX{tx: tx, logger: s.logger}
	return txPGX, nil
}
