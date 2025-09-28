package middleware

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/Soumil-2007/file-sharing-webApp/services/auth"
	"github.com/Soumil-2007/file-sharing-webApp/types"
)

// AuthMiddleware returns a mux.MiddlewareFunc
func AuthMiddleware(userStore types.UserStore) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(auth.WithJWTAuth(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		}, userStore))
	}
}
