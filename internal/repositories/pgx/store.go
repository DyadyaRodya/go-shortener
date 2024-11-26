package pgx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	"github.com/DyadyaRodya/go-shortener/internal/usecases"
	"github.com/DyadyaRodya/go-shortener/internal/usecases/dto"
	"github.com/DyadyaRodya/go-shortener/pkg/itemsets"
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
	_, err = tx.Exec(ctx, `CREATE TABLE IF NOT EXISTS short_urls (
        uuid CHAR(8) NOT NULL PRIMARY KEY, 
        url TEXT UNIQUE NOT NULL,
        deleted_at TIMESTAMPTZ NULL DEFAULT NULL
    )`)
	if err != nil {
		s.logger.Error("Failed to create table `short_urls`", zap.Error(err))
		return err
	}

	s.logger.Info("Creating table `users`")
	_, err = tx.Exec(ctx, `CREATE TABLE IF NOT EXISTS users (
        uuid UUID NOT NULL PRIMARY KEY
    )`)
	if err != nil {
		s.logger.Error("Failed to create table `users`", zap.Error(err))
		return err
	}

	s.logger.Info("Creating table `users_short_urls`")
	_, err = tx.Exec(ctx, `CREATE TABLE IF NOT EXISTS users_short_urls (
        user_uuid UUID REFERENCES users(uuid) ON DELETE CASCADE, 
        short_url_uuid CHAR(8) REFERENCES short_urls(uuid) ON DELETE CASCADE,
        CONSTRAINT users_short_urls_pkey PRIMARY KEY (user_uuid, short_url_uuid)                                           
	)`)
	if err != nil {
		s.logger.Error("Failed to create table `users_short_urls`", zap.Error(err))
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

func (s *StorePGX) LoadURLs(ctx context.Context, src map[string]string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.logger.Error("Failed to start transaction", zap.Error(err))
		return err
	}
	defer tx.Rollback(ctx)

	sqlQuery := `INSERT INTO short_urls (uuid, url) VALUES (@uuid, @url) ON CONFLICT DO NOTHING`

	for uuid, url := range src {
		ct, err := tx.Exec(ctx, sqlQuery, pgx.NamedArgs{"uuid": uuid, "url": url})

		if err != nil {
			s.logger.Error("Failed to insert into short_urls",
				zap.String("uuid", uuid),
				zap.String("url", url),
				zap.Error(err))
			return fmt.Errorf("error in StorePGX.LoadURLs for url: %s uuid %s: %w", url, uuid, err)
		}
		if !ct.Insert() {
			s.logger.Error("Failed to insert into short_urls",
				zap.String("uuid", uuid),
				zap.String("url", url),
				zap.Any("commandTag", ct))
			return errors.New("error in StorePGX.LoadURLs: not inserted url: " + url + " uuid: " + uuid)
		}
	}
	err = tx.Commit(ctx)
	return err
}

func (s *StorePGX) LoadUsersURLs(ctx context.Context, src map[string][]string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.logger.Error("Failed to start transaction", zap.Error(err))
		return err
	}
	defer tx.Rollback(ctx)

	sqlQueryUsers := `INSERT INTO users (uuid) VALUES (@uuid) ON CONFLICT DO NOTHING`
	sqlQueryUsersURLs := `INSERT INTO users_short_urls (user_uuid, short_url_uuid) 
							VALUES (@user_uuid, @short_url_uuid) ON CONFLICT DO NOTHING`

	for userUUID, shortUUIDs := range src {
		ct, err := tx.Exec(ctx, sqlQueryUsers, pgx.NamedArgs{"uuid": userUUID})

		if err != nil {
			s.logger.Error("Failed to insert into users", zap.String("user_uuid", userUUID), zap.Error(err))
			return fmt.Errorf("error in StorePGX.UsersURLs for uuid: %s: %w", userUUID, err)
		}
		if !ct.Insert() {
			s.logger.Error("Failed to insert into users",
				zap.String("user_uuid", userUUID),
				zap.Any("commandTag", ct))
			return fmt.Errorf("error in StorePGX.UsersURLs: not inserted user_uuid: %s", userUUID)
		}

		for _, shortUUID := range shortUUIDs {
			ct, err := tx.Exec(ctx, sqlQueryUsersURLs, pgx.NamedArgs{"user_uuid": userUUID, "short_url_uuid": shortUUID})

			if err != nil {
				s.logger.Error("Failed to insert into short_urls",
					zap.String("user_uuid", userUUID),
					zap.String("short_url_uuid", shortUUID),
					zap.Error(err))
				return fmt.Errorf("error in StorePGX.UsersURLs for user_uuid: %s short_url_uuid %s: %w",
					shortUUID,
					userUUID,
					err,
				)
			}
			if !ct.Insert() {
				s.logger.Error("Failed to insert into short_urls",
					zap.String("user_uuid", userUUID),
					zap.String("short_url_uuid", shortUUID),
					zap.Any("commandTag", ct))
				return fmt.Errorf("error in StorePGX.UsersURLs: not inserted user_uuid: %s short_url_uuid %s",
					shortUUID,
					userUUID,
				)
			}
		}
	}
	err = tx.Commit(ctx)
	return err
}

func (s *StorePGX) URLs(ctx context.Context) (*map[string]string, error) {
	dst := make(map[string]string)

	rows, err := s.pool.Query(ctx, `SELECT uuid, url FROM short_urls WHERE deleted_at IS NULL`)
	if err != nil {
		s.logger.Error("Failed to query database StorePGX.URLs", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var uuid, url string
		if err := rows.Scan(&uuid, &url); err != nil {
			s.logger.Error("Failed to scan StorePGX.URLs", zap.Error(err))
			return nil, err
		}
		dst[uuid] = url
	}
	if err := rows.Err(); err != nil {
		s.logger.Error("Failed to query database StorePGX.URLs", zap.Error(err))
		return nil, err
	}
	s.logger.Info("Successfully read database StorePGX.URLs")
	return &dst, nil
}

func (s *StorePGX) UsersURLs(ctx context.Context) (*map[string][]string, error) {
	dst := make(map[string][]string)

	rows, err := s.pool.Query(ctx, `SELECT user_uuid, short_url_uuid FROM users_short_urls`)
	if err != nil {
		s.logger.Error("Failed to query database StorePGX.UsersURLs", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userUUID, ShortUUID string
		if err := rows.Scan(&userUUID, &ShortUUID); err != nil {
			s.logger.Error("Failed to scan StorePGX.UsersURLs", zap.Error(err))
			return nil, err
		}
		dst[userUUID] = append(dst[userUUID], ShortUUID)
	}
	if err := rows.Err(); err != nil {
		s.logger.Error("Failed to query database StorePGX.UsersURLs", zap.Error(err))
		return nil, err
	}
	s.logger.Info("Successfully read database StorePGX.UsersURLs")
	return &dst, nil
}

func (s *StorePGX) TestConnection(ctx context.Context) error {
	s.logger.Debug("Pinging database")
	return s.pool.Ping(ctx)
}

func (s *StorePGX) AddURL(ctx context.Context, ShortURL *entity.ShortURL, OwnerUUID string) error {
	tx, err := s.Begin(ctx)
	if err != nil {
		err = fmt.Errorf("StorePGX.AddURL: %w", err)
		s.logger.Error("Failed to start transaction in StorePGX.AddURL", zap.Error(err))
		return err
	}
	defer tx.Rollback(ctx)

	if err := tx.AddURL(ctx, ShortURL, false); err != nil {
		err = fmt.Errorf("StorePGX.AddURL: %w", err)
		s.logger.Info("Failed to add URL in StorePGX.AddURL", zap.Error(err))
		return err
	}

	if OwnerUUID != "" {
		if err := tx.AddUserIfNotExists(ctx, OwnerUUID); err != nil {
			err = fmt.Errorf("StorePGX.AddURL: %w", err)
			s.logger.Error("Failed to add user in StorePGX.AddURL", zap.Error(err))
			return err
		}

		if err := tx.AddUserURL(ctx, ShortURL.ID, OwnerUUID); err != nil {
			err = fmt.Errorf("StorePGX.AddUserURL: %w", err)
			s.logger.Error("Failed to add user in StorePGX.AddUserURL", zap.Error(err))
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		err = fmt.Errorf("StorePGX.AddURL: %w", err)
		s.logger.Error("Failed to commit transaction in StorePGX.AddURL", zap.Error(err))
		return err
	}
	return nil
}

func (s *StorePGX) GetURLByID(ctx context.Context, ID string) (*entity.ShortURL, error) {
	s.logger.Debug("Getting URL by ID", zap.String("ID", ID))

	var url string
	var deletedAt sql.NullTime
	err := s.pool.QueryRow(ctx, `SELECT url, deleted_at FROM short_urls WHERE uuid = @uuid`, pgx.NamedArgs{"uuid": ID}).
		Scan(&url, &deletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrShortURLNotFound
		}
		s.logger.Error("Failed to get URL by ID", zap.String("ID", ID), zap.Error(err))
		return nil, err
	}
	shortURL := &entity.ShortURL{URL: url, ID: ID}
	if !deletedAt.Valid {
		return shortURL, nil
	}
	return shortURL, entity.ErrShortURLDeleted
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

func (s *StorePGX) GetUserUrls(ctx context.Context, UserUUID string) (map[string]*entity.ShortURL, error) {
	tx, err := s.Begin(ctx)
	if err != nil {
		err = fmt.Errorf("StorePGX.GetUserUrls: %w", err)
		s.logger.Error("Failed to start transaction in StorePGX.GetUserUrls", zap.Error(err))
		return nil, err
	}
	defer tx.Rollback(ctx)

	res, err := tx.GetUserUrls(ctx, UserUUID)
	if err != nil {
		err = fmt.Errorf("StorePGX.GetUserUrls: %w", err)
		s.logger.Error("Failed to get user URLs in StorePGX.GetUserUrls", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (s *StorePGX) DeleteUserURLs(ctx context.Context, requests ...*dto.DeleteUserShortURLsRequest) error {
	s.logger.Debug("StorePGX.DeleteUserURLs", zap.Any("requests", requests))

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.logger.Error("Failed to begin transaction StorePGX.DeleteUserURLs", zap.Error(err))
		return fmt.Errorf("StorePGX.DeleteUserURLs: %w", err)
	}
	defer tx.Rollback(ctx)

	toDelete := make([]string, 0)

	deleteQuery := `DELETE FROM users_short_urls WHERE user_uuid = $1 AND short_url_uuid = ANY($2)`
	for _, request := range requests {
		if request.UserUUID == "" {
			continue
		}
		ct, err := tx.Exec(ctx, deleteQuery, request.UserUUID, request.ShortURLUUIDs)

		if err != nil {
			s.logger.Info("Failed to delete from users_short_urls",
				zap.String("user_uuid", request.UserUUID),
				zap.Any("short_url_uuids", request.ShortURLUUIDs),
				zap.Error(err))
			return fmt.Errorf("StorePGX.DeleteUserURLs: %w", err)
		}

		if ct.Delete() { // deleted some links - can be last owner
			toDelete = itemsets.AddItems(toDelete, request.ShortURLUUIDs)
		}
	}
	_, err = tx.Exec(ctx, `UPDATE short_urls AS su SET deleted_at=now() 
                        				WHERE su.uuid = ANY($1) AND NOT EXISTS(
                        				SELECT 1 FROM users_short_urls AS usu WHERE usu.short_url_uuid=su.uuid
                        				                                      )`,
		toDelete)

	if err != nil {
		s.logger.Info("Failed to update short_urls with deleted_at = now()",
			zap.Any("short_url_uuids", toDelete),
			zap.Error(err))
		return fmt.Errorf("StorePGX.DeleteUserURLs: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.logger.Error("Failed to commit transaction in StorePGX.DeleteUserURLs", zap.Error(err))
		return fmt.Errorf("StorePGX.DeleteUserURLs: %w", err)
	}
	return nil
}
