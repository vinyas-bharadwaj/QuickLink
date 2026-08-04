package main

import (
	"fmt"
	"net/http"
	"time"
	"urlshortener/internal/api/middlewares"
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

	// Set a rate limiter of 10 requests every 30 seconds
	rlMiddleware := middlewares.NewRateLimiter(10, 30*time.Second)
	handler := rlMiddleware.Middleware(router)

	// Start the server
	port := ":8080"
	fmt.Println("Server starting on port:", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		fmt.Println("Error starting server:", err.Error())
	}

}
