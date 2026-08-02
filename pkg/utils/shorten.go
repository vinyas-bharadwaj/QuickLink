package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// Creates a unique 6 character URL
func GenerateShortURL() (string, error) {
	randomBytes := make([]byte, 6)
	_, err := rand.Read(randomBytes) // Create a byte slice array of length 6

	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(randomBytes)[:6], nil
}
