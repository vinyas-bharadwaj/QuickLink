package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"urlshortener/internal/models"
	"urlshortener/pkg/utils"

	"github.com/redis/go-redis/v9"
)

// Contains dependencies for handling HTTP requests
type ShortenerHandler struct {
	db    *sql.DB
	cache *redis.Client
}

// Function to initialize a new handler
func NewShortenerHandler(db *sql.DB, cache *redis.Client) *ShortenerHandler {
	return &ShortenerHandler{db: db, cache: cache}
}

// Handles POST requests to shorten the URL
func (h *ShortenerHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	r.ParseForm()
	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "Missing URL parameter", http.StatusBadRequest)
		return
	}

	// Check if original URL already exists in database
	var existingShortURL string
	err := h.db.QueryRow("SELECT short_url FROM short_long_mapping WHERE original_url = ?", originalURL).Scan(&existingShortURL)
	if err == nil {
		response := models.URLResponse{
			Message:     "URL already shortened",
			ShortURL:    existingShortURL,
			OriginalURL: originalURL,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Maximum retries in case of a collision occuring
	const maxRetries = 5
	var shortURL string
	for i := 0; i < maxRetries; i++ {
		// Generate short URL
		shortURL, err = utils.GenerateShortURL()
		if err != nil {
			http.Error(w, "Error generarting the short URL", http.StatusInternalServerError)
			return
		}

		// Save mapping into our in-memory store
		_, err = h.db.Exec("INSERT INTO short_long_mapping VALUES(?, ?)", shortURL, originalURL)
		// Break out of the loop if no error occurred
		if err == nil {
			break
		}

		// Continue the next iteration of the loop and retry generating another shortURL
	}

	// In case we fail all the retry attempts
	if err != nil {
		http.Error(w, "Could not save short url after retries", http.StatusInternalServerError)
		return
	}

	response := models.URLResponse{
		Message:     "Successfully created a shortened URL",
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Handles GET requests to redirect to the original URL
func (h *ShortenerHandler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the short URL from the request URL
	shortURL := r.URL.Path[1:] // Remove the leading slash
	fmt.Println("Shortened URL:", shortURL)

	// Check if URL is present in the cache
	ctx := context.Background()
	cachedURL, err := h.cache.Get(ctx, shortURL).Result()
	if err == nil {
		// Print statement to check whether the cache was accessed or not
		fmt.Println("Cache accessed")
		// Key is already present in cache
		// We can directly redirect the user without having to make a call to the DB
		http.Redirect(w, r, cachedURL, http.StatusFound)
		return

	}

	// Retrieve the long URL from the database
	var originalURL string
	err = h.db.QueryRow("SELECT original_url FROM short_long_mapping WHERE short_url = ?", shortURL).Scan(&originalURL)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	// The key is not present in the cache
	// Therefore we set the key
	_, err = h.cache.Set(ctx, shortURL, originalURL, 0).Result()
	if err != nil {
		http.Error(w, "Error inserting value to the cache", http.StatusInternalServerError)
		return
	}

	// Redirect to the original URL
	http.Redirect(w, r, originalURL, http.StatusFound)
}
