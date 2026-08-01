package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"urlshortener/internal/models"
	"urlshortener/pkg/utils"
)

// Contains dependencies for handling HTTP requests
type ShortenerHandler struct {
	store *utils.Store
}

// Function to initialize a new handler
func NewShortenerHandler(store *utils.Store) *ShortenerHandler {
	return &ShortenerHandler{store: store}
}

// Handles POST requests to shorten the URL
func (h *ShortenerHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	// Parse the request body
	r.ParseForm()
	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "Missing URL parameter", http.StatusBadRequest)
		return
	}

	// Generate short URL
	shortURL, err := utils.GenerateShortURL()
	if err != nil {
		http.Error(w, "Error generarting the short URL", http.StatusInternalServerError)
		return
	}

	// Save mapping into our in-memory store
	h.store.Save(shortURL, originalURL)
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

	// Retrieve the long URL from the store
	originalURL, exists := h.store.Find(shortURL)
	if !exists {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	// Redirect to the original URL
	http.Redirect(w, r, originalURL, http.StatusFound)
}
