package middlewares

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	mu        sync.Mutex
	visitors  map[string]int
	limit     int
	resetTime time.Duration
}

func NewRateLimiter(limit int, resetTime time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors:  make(map[string]int),
		limit:     limit,
		resetTime: resetTime,
	}

	// Initialize a goroutine which runs in the background
	// This clears the visitors list after every 'resetTime' interval
	go rl.resetVisitorCount()
	return rl
}

func (rl *RateLimiter) resetVisitorCount() {
	// Infinite loop
	for {
		// Refresh the visitors list after every 'resetTime' interval
		time.Sleep(rl.resetTime)

		// Lock the mutex before accessing the shared resource
		rl.mu.Lock()
		rl.visitors = make(map[string]int)
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.mu.Lock()
		defer rl.mu.Unlock()

		visitorIP := r.RemoteAddr
		rl.visitors[visitorIP]++

		if rl.visitors[visitorIP] > rl.limit {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
