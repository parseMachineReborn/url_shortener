package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func AuthenticationMiddleware(next http.Handler, secretKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		if bearer == "" {
			http.Error(w, "Отсутствует токен", http.StatusUnauthorized)
			return
		}
		bearer = strings.TrimPrefix(bearer, "Bearer ")

		token, err := jwt.Parse(bearer, func(t *jwt.Token) (any, error) { return []byte(secretKey), nil })
		if err != nil {
			http.Error(w, "Не валидный токен", http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := claims["user_id"]
		ctx := context.WithValue(r.Context(), "user_id", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
