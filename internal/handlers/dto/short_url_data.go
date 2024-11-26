package dto

// ShortURLData Structure of items of JSON response with short URLs information
type ShortURLData struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
