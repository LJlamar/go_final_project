package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func authenticationHandler(w http.ResponseWriter, r *http.Request) {

	pswd := map[string]string{
		//Structure prior to JSON:
		"password": "",
	}

	err := json.NewDecoder(r.Body).Decode(&pswd)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": "Некорректный формат JSON"})
		return
	}

	password := os.Getenv("TODO_PASSWORD")
	if pswd["password"] == password {
		preHash := sha256.Sum256([]byte(password))
		hashString := hex.EncodeToString(preHash[:])

		claims := jwt.MapClaims{
			"data": hashString,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		secretKey := os.Getenv("JWT_SECRET")
		if secretKey == "" {
			log.Fatal("internal error of generating token: secret key (JWT_SECRET) for JWT isn't found or have empty value")
			w.WriteHeader(http.StatusInternalServerError)
			writeJson(w, map[string]string{
				"error": "Внутренняя ошибка сервера. Попробуйте позже.",
			})
			return
		}

		tokenString, err := token.SignedString([]byte(secretKey))
		if err != nil {
			log.Printf("Internal error generating token: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			writeJson(w, map[string]string{
				"error": "Внутренняя ошибка сервера. Попробуйте позже.",
			})
			return
		}
		//Setting up variable for production:
		isSecure := os.Getenv("APP_ENV") == "production"

		cookie := &http.Cookie{
			Name:     "token",
			Value:    tokenString,
			Path:     "/",
			HttpOnly: true,
			Secure:   isSecure,
			MaxAge:   3600,
		}

		http.SetCookie(w, cookie)

		writeJson(w, map[string]string{"token": tokenString})
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": "Неверный пароль"})
		return
	}
}
