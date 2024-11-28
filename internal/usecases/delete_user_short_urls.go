package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/DyadyaRodya/go-shortener/internal/usecases/dto"
)

// DeleterErrorLogger function for error logging
type DeleterErrorLogger func(msg string)

// UsersShortURLsDeleter reads provided delChan chan *dto.DeleteUserShortURLsRequest once per 10 seconds
// and executes batch delete for URLs.
//
// Unlinks URL from user in request and removes all URLs that are not linked to any users.
//
// In case of failure it writes URLs back to channel to retry later.
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
