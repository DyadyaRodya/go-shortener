package services

import (
	"crypto/rand"
	"encoding/hex"
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
