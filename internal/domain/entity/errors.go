package entity

import "errors"

var (
	ErrShortURLNotFound = errors.New("short url not found")
	ErrUUIDTaken        = errors.New("uuid taken")
	ErrShortURLExists   = errors.New("short url exists")
)
