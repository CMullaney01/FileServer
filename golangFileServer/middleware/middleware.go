package middleware

import (
	"net/http"
	"log"

	"github.com/CMullaney01/FileServer/handlers"
)

// AuthCORSHandler adds CORS headers and performs authentication.
func AuthCORSHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("Middleware")
		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Expose-Headers", "Set-Cookie")
			w.WriteHeader(http.StatusOK)
			return
		}
		// log.Printf("Middleware1")
		// Get session token from cookie
		c, err := r.Cookie("session_token")
		if err != nil {
			// log.Printf("Middleware2")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		sessionToken := c.Value
		// Check if the session token is valid
		userSession, exists := handlers.Sessions[sessionToken]  // Update this line to use the correct variable name
		if !exists || userSession.IsExpired() {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		log.Printf("Middleware3")
		
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Expose-Headers", "Set-Cookie")
		// Continue to the next handler
		h.ServeHTTP(w, r)
	})
}

// CORSMiddleware adds CORS headers for all requests.
func CORSHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Expose-Headers", "Set-Cookie")
			w.WriteHeader(http.StatusOK)
			return
		}

		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Expose-Headers", "Set-Cookie")
		h.ServeHTTP(w, r)
	})
}