package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"urlshortener/internal/models"
	"urlshortener/pkg/utils"
)

// Contains dependencies for handling HTTP requests
type ShortenerHandler struct {
	db *sql.DB
}

// Function to initialize a new handler
func NewShortenerHandler(db *sql.DB) *ShortenerHandler {
	return &ShortenerHandler{db: db}
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

	// Generate short URL
	shortURL, err := utils.GenerateShortURL()
	if err != nil {
		http.Error(w, "Error generarting the short URL", http.StatusInternalServerError)
		return
	}

	// Save mapping into our in-memory store
	_, err = h.db.Exec("INSERT INTO short_long_mapping VALUES(?, ?)", shortURL, originalURL)
	if err != nil {
		http.Error(w, "Error adding values to the database", http.StatusInternalServerError)
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

	// Retrieve the long URL from the database
	var originalURL string
	err := h.db.QueryRow("SELECT original_url FROM short_long_mapping WHERE short_url = ?", shortURL).Scan(&originalURL)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	// Redirect to the original URL
	http.Redirect(w, r, originalURL, http.StatusFound)
}
