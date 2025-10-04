package middleware

import (
	"context"
	"net/http"
	"strings"

	"blog-api/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ambil token dari Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		// Format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			respondError(w, http.StatusUnauthorized, "Invalid authorization format")
			return
		}

		tokenString := parts[1]

		// Parse dan validasi token
		cfg := config.LoadConfig()
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validasi signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			respondError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// Ekstrak user_id dari claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			respondError(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			respondError(w, http.StatusUnauthorized, "Invalid user_id in token")
			return
		}

		// Simpan user_id ke context
		ctx := context.WithValue(r.Context(), UserIDKey, uint(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper untuk mengambil user_id dari context
func GetUserID(r *http.Request) (uint, bool) {
	userID, ok := r.Context().Value(UserIDKey).(uint)
	return userID, ok
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(`{"error":"` + message + `"}`))
}
