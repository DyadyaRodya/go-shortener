package pgx

import (
	"context"
	"errors"
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StorePGX struct {
	pool *pgxpool.Pool
}

func NewStorePGX(pool *pgxpool.Pool) *StorePGX {
	return &StorePGX{pool: pool}
}

func (s *StorePGX) TestConnection(ctx context.Context) error {
	return s.pool.Ping(ctx)
}

func (s *StorePGX) AddURL(ctx context.Context, ShortURL *entity.ShortURL) error {
	return errors.New("not implemented")
}

func (s *StorePGX) GetURLByID(ctx context.Context, ID string) (*entity.ShortURL, error) {
	return nil, errors.New("not implemented")
}
