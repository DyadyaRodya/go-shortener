package tests

import (
	"context"

	"github.com/DyadyaRodya/go-shortener/internal/usecases/dto"

	"github.com/brianvoe/gofakeit/v6"
)

func (u *usecasesSuite) TestUsecases_GetStats_Success() {
	ctx := context.Background()
	stats := &dto.StatsResponse{}

	u.urlStorage.EXPECT().GetStats(ctx).Return(stats, nil).Once()

	result, err := u.usecases.GetStats(ctx)

	u.NoError(err)
	u.Equal(stats, result)
}

func (u *usecasesSuite) TestUsecases_GetStats_AnyError() {
	ctx := context.Background()
	otherError := gofakeit.Error()

	u.urlStorage.EXPECT().GetStats(ctx).Return(nil, otherError).Once()

	result, err := u.usecases.GetStats(ctx)

	u.ErrorIs(err, otherError)
	u.Empty(result)
}
