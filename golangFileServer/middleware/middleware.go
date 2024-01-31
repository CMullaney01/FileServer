package middleware

import (
	"net/http"

	"github.com/CMullaney01/FileServer/handlers"
)

// AuthCORSHandler adds CORS headers and performs authentication.
func AuthCORSHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Get session token from cookie
		c, err := r.Cookie("session_token")
		if err != nil {
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

		// Continue to authentication
		if !handlers.AuthenticateUser(userSession.Username, "") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

		// Continue to the next handler
		h.ServeHTTP(w, r)
	})
}