package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Getting secret key
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			log.Fatal("Ошибка конфигурации: JWT_SECRET не установлен!")
		}

		// Getting data from cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Unauthorized: missing token cookie", http.StatusUnauthorized)
			return
		}

		//Getting token
		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			//Cheking signing method:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrInvalidKey
			}
			return []byte(secret), nil
		})
		//Validating token
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: invalid or expired token", http.StatusUnauthorized)
			return
		}

		_, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized: invalid claims format", http.StatusUnauthorized)
			return
		}

		next(w, r)
	})
}
