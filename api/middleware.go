package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "userID"

func (cfg *ApiConfig) MiddlewareAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			RespondWithError(w, http.StatusUnauthorized, "Missing Authorization Header")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			RespondWithError(w, http.StatusUnauthorized, "Malformed Token")
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			RespondWithError(w, http.StatusUnauthorized, "Invalid Token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			RespondWithError(w, http.StatusInternalServerError, "Invalid Token Claims")
			return
		}

		userIDStr, ok := claims["sub"].(string)
		if !ok {
			RespondWithError(w, http.StatusInternalServerError, "Token missing User ID")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userIDStr)
		next(w, r.WithContext(ctx))
	}
}
