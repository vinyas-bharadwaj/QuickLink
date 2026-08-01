package models

type URLResponse struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
	Message     string `json:"message"`
}
