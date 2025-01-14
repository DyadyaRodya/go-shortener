package usecases

import (
	"context"

	"github.com/DyadyaRodya/go-shortener/internal/usecases/dto"
)

// GetStats Reads total *dto.StatsData from storage
func (u *Usecases) GetStats(ctx context.Context) (*dto.StatsResponse, error) {
	stats, err := u.urlStorage.GetStats(ctx)
	if err != nil {
		return nil, err
	}
	return stats, nil
}
