package security

import (
	"net/http"
	"sync"
	"time"
)

// Structure to manage a client's requests
type client struct {
	requests []time.Time // Request history
}

// Map to store clients and a mutex for thread safety
var clients = make(map[string]*client)
var mu sync.Mutex

// Rate limiting parameters
const maxRequests = 12              // Maximum number of allowed requests
const timeWindow = 10 * time.Second // Time interval to limit requests

// Main function to handle requests with rate limiting
func RateLimitedHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr // Identify the client by its IP

		// Lock access to the clients map
		mu.Lock()
		defer mu.Unlock()

		// Retrieve client information or create a new record
		c, exists := clients[ip]
		if !exists {
			c = &client{requests: []time.Time{}}
			clients[ip] = c
		}

		// Clean up requests outside the time window
		now := time.Now()
		validRequests := []time.Time{}
		for _, t := range c.requests {
			if now.Sub(t) <= timeWindow {
				validRequests = append(validRequests, t)
			}
		}
		c.requests = validRequests

		// Check if the request limit is reached
		if len(c.requests) >= maxRequests {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		// Add the new request to the history
		c.requests = append(c.requests, now)

		next(w, r)
	}
}

// Periodic cleanup of inactive clients
func CleanupInactiveClients() {
	for {
		time.Sleep(1 * time.Minute) // Clean up every minute

		mu.Lock()
		for ip, c := range clients {
			// If a client has not sent requests for a certain time, delete it
			if len(c.requests) == 0 || time.Since(c.requests[len(c.requests)-1]) > 2*timeWindow {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}
