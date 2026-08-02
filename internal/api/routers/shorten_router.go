package routers

import (
	"log"
	"net/http"
	"urlshortener/internal/api/handlers"
	"urlshortener/internal/repository/sqlconnect"
)

func shortenRourter(mux *http.ServeMux) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		log.Fatal("Error: ", err.Error())
		return
	}

	shortenHandler := handlers.NewShortenerHandler(db)

	// URL shortening routes
	mux.HandleFunc("/shorten", shortenHandler.ShortenURL)
	mux.HandleFunc("/", shortenHandler.RedirectURL)
}
