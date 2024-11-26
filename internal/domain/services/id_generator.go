package services

import (
	"crypto/rand"
	"encoding/hex"

	uuidPkg "github.com/google/uuid"
)

// IDGenerator Provides method for generating hex string
type IDGenerator struct {
}

// NewIDGenerator constructor for generator IDGenerator
func NewIDGenerator() *IDGenerator {
	return &IDGenerator{}
}

// Generate Generates random 4-byte hex and returns as 8 chars long hex string
func (i *IDGenerator) Generate() (string, error) {
	id := make([]byte, 4)
	_, err := rand.Read(id)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(id), nil
}

// UUID4Generator Provides method for generating UUID string
type UUID4Generator struct{}

// NewUUID4Generator constructor for generator UUID4Generator
func NewUUID4Generator() *UUID4Generator {
	return &UUID4Generator{}
}

// Generate Generates UUID4 string
func (u *UUID4Generator) Generate() (string, error) {
	uuid, err := uuidPkg.NewRandom()
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}
