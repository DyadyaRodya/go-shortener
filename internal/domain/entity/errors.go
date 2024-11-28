package entity

import "errors"

// All domain errors
var (
	// ErrShortURLNotFound indicates if short URL not found
	ErrShortURLNotFound = errors.New("short url not found")
	// ErrUUIDTaken indicates if short URL ID taken
	ErrUUIDTaken = errors.New("uuid taken")
	// ErrShortURLExists indicates if short URL exists for full URL
	ErrShortURLExists = errors.New("short url exists")
	// ErrShortURLDeleted indicates if short URL marked as deleted
	ErrShortURLDeleted = errors.New("short url deleted")
)
