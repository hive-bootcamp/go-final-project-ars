package api

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwtTokenValue string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				jwtTokenValue = cookie.Value
			}

			parsedToken, err := jwt.Parse(jwtTokenValue, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return []byte(os.Getenv("SECRET_KEY")), nil
			})

			if err != nil || !parsedToken.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// здесь код для валидации и проверки JWT-токена
			hashOfBytes := sha256.Sum256([]byte(pass))
			currentHash := hex.EncodeToString(hashOfBytes[:])

			jwtClaims, ok := parsedToken.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid claims", http.StatusUnauthorized)
				return
			}
			if jwtClaims["hash"].(string) != currentHash {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
