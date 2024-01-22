package middlewares

import (
	"context"
	"net/http"

	"github.com/atanurdemir/gatekeeper/src/types"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretKey := "your-secret-key"
		accessToken := r.Header.Get("Authorization")
		if accessToken == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		accessToken = accessToken[len("Bearer "):]

		claims := &types.Claims{}
		_, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil || claims.UserID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), types.UserKey, claims.UserID))

		next.ServeHTTP(w, r)
	})
}
