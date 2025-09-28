package middleware

import (
	"net/http"
)

// AuthMiddleware returns a mux.MiddlewareFunc
func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: validate JWT / session here
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			// If valid → call the next handler
			next.ServeHTTP(w, r)
		})
	}
}
