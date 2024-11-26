package services

import (
	"crypto/rand"
	"encoding/hex"

	uuidPkg "github.com/google/uuid"
)

type IDGenerator struct {
}

func NewIDGenerator() *IDGenerator {
	return &IDGenerator{}
}

func (i *IDGenerator) Generate() (string, error) {
	id := make([]byte, 4)
	_, err := rand.Read(id)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(id), nil
}

type UUID4Generator struct{}

func NewUUID4Generator() *UUID4Generator {
	return &UUID4Generator{}
}

func (u *UUID4Generator) Generate() (string, error) {
	uuid, err := uuidPkg.NewRandom()
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}
