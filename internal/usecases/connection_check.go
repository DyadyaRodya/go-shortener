package usecases

import (
	"context"
	"fmt"
)

// CheckConnection checks readiness of storage connection.
func (u *Usecases) CheckConnection(ctx context.Context) error {
	err := u.urlStorage.TestConnection(ctx)
	if err != nil {
		return fmt.Errorf("Usecases.CheckConnection: %w", err)
	}
	return nil
}
