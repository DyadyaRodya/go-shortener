package dto

import "github.com/DyadyaRodya/go-shortener/internal/domain/entity"

// BatchCreateRequest DTO for passing batch create short URL requests to usecases.Usecases
type BatchCreateRequest struct {
	CorrelationID string
	OriginalURL   string
}

// BatchCreateResponse DTO for returning batch create short URL results from usecases.Usecases
type BatchCreateResponse struct {
	CorrelationID string
	ShortURL      *entity.ShortURL
}
