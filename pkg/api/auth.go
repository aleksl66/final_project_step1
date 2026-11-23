package api

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret — секрет для подписи JWT.
var jwtSecret = []byte("todo-scheduler-secret")

// passwordHash возвращает SHA-256 хэш пароля в hex.
func passwordHash(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

// makeToken создаёт JWT с хэшем пароля в полезной нагрузке.
func makeToken(password string) (string, error) {
	claims := jwt.MapClaims{
		"hash": passwordHash(password),
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// validateToken проверяет JWT и соответствие хэша текущему паролю.
func validateToken(tokenStr, password string) bool {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}
	hash, ok := claims["hash"].(string)
	if !ok {
		return false
	}
	return hash == passwordHash(password)
}

// auth — middleware для проверки аутентификации.
func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
			if !validateToken(cookie.Value, pass) {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	}
}
