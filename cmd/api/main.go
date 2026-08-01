package main

import (
	"fmt"
	"net/http"
	"urlshortener/internal/api/handlers"
	"urlshortener/pkg/utils"
)

func main() {

	// Initialize a new store
	store := utils.NewStore()

	// Initialize the handler
	ShortenerHandler := handlers.NewShortenerHandler(store)

	// Define routes
	http.HandleFunc("/shorten", ShortenerHandler.ShortenURL)
	http.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "web/index.html")
	})
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	http.HandleFunc("/", ShortenerHandler.RedirectURL)

	// Start the server
	port := ":8080"
	fmt.Println("Server starting on port:", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Println("Error starting server:", err.Error())
	}

}
