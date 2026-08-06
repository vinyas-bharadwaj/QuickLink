package routers

import (
	"log"
	"net/http"
	"urlshortener/internal/api/handlers"
	"urlshortener/internal/repository/redisconnect"
	"urlshortener/internal/repository/sqlconnect"
)

func shortenRourter(mux *http.ServeMux) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		log.Fatal("Error: ", err.Error())
		return
	}

	cache := redisconnect.ConnectRedis()
	shortenHandler := handlers.NewShortenerHandler(db, cache)

	// URL shortening routes
	mux.HandleFunc("/shorten", shortenHandler.ShortenURL)
	mux.HandleFunc("/", shortenHandler.RedirectURL)
}
