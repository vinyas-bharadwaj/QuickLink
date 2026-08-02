package main

import (
	"fmt"
	"net/http"
	"urlshortener/internal/api/routers"
	"urlshortener/internal/repository/sqlconnect"
)

func init() {
	err := sqlconnect.InitDB()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	_, err = sqlconnect.ConnectDB()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func main() {

	// Initialize the router
	router := routers.Router()

	// Define routes
	http.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "web/index.html")
	})
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))

	// Start the server
	port := ":8080"
	fmt.Println("Server starting on port:", port)
	if err := http.ListenAndServe(port, router); err != nil {
		fmt.Println("Error starting server:", err.Error())
	}

}
