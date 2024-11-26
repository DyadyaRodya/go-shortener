package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/DyadyaRodya/go-shortener/internal/usecases/dto"
)

type DeleterErrorLogger func(msg string)

func (u *Usecases) UsersShortURLsDeleter(
	ctx context.Context,
	errorLogger DeleterErrorLogger,
	delChan chan *dto.DeleteUserShortURLsRequest,
) {
	// process requests collected for 10 sec
	ticker := time.NewTicker(10 * time.Second)

	var requests []*dto.DeleteUserShortURLsRequest
	for {
		select {
		case request := <-delChan:
			requests = append(requests, request)
		case <-ctx.Done():
			return
		case <-ticker.C:
			if len(requests) == 0 {
				continue
			}
			err := u.urlStorage.DeleteUserURLs(ctx, requests...)
			if err != nil {
				errorLogger(fmt.Sprintf("Error deleting urls: %v", err))
				// not clearing requests - we will try again
				continue
			}
			requests = nil // successfully processed - clear requests
		}
	}
}
