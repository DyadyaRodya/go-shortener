package dto

type DeleteUserShortURLsRequest struct {
	UserUUID      string
	ShortURLUUIDs []string
}
