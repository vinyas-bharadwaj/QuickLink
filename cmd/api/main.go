package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"urlshortener/internal/api/middlewares"
	"urlshortener/internal/api/routers"
	"urlshortener/internal/repository/sqlconnect"

	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}

	// Initialize the SQL database
	if err := sqlconnect.InitDB(); err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Check if the connection to the SQL database works
	if _, err := sqlconnect.ConnectDB(); err != nil {
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
