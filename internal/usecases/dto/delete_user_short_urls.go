package dto

// DeleteUserShortURLsRequest DTO for passing delete short URLs requests to usecases.Usecases deleter throw the chan
type DeleteUserShortURLsRequest struct {
	UserUUID      string
	ShortURLUUIDs []string
}
