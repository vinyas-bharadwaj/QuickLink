package utils

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
)

// In-memory store of all the URLs
type Store struct {
	urls map[string]string // Key-Value pairs (ShortURL -> LongURL)
	mu   sync.Mutex        // Mutex for concurrency control
}

// Function to generate a new in-memory store
func NewStore() *Store {
	return &Store{
		urls: make(map[string]string),
	}
}

func (s *Store) Save(shortURL, originalURL string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.urls[shortURL] = originalURL
}

func (s *Store) Find(shortURL string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	originalURL, exists := s.urls[shortURL]
	return originalURL, exists
}

// Creates a unique 6 character URL
func GenerateShortURL() (string, error) {
	randomBytes := make([]byte, 6)
	_, err := rand.Read(randomBytes) // Create a byte slice array of length 6

	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(randomBytes)[:6], nil
}
