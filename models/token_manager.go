package models

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/silasburger/lenslocked/rand"
)

type TokenManager struct {
	BytesPerToken int
}

const (
	MinBytesPerToken = 32
)

func (tm TokenManager) New() (token, tokenHash string, err error) {
	bytesPerToken := tm.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err = rand.String(bytesPerToken)
	if err != nil {
		return "", "", fmt.Errorf("new: %w", err)
	}
	return token, tm.Hash(token), err
}

func (tm TokenManager) Hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	// base64 encode the data into a string
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
