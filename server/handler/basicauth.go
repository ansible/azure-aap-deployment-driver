package handler

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
	"server/config"
)

const USERNAME string = "installer"

func BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Provided username/pass
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))

			expectedUsernameHash := sha256.Sum256([]byte(USERNAME))
			expectedPasswordHash := sha256.Sum256([]byte(config.GetEnvironment().PASSWORD))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		respondError(w, 401, "Unauthorized")
	})
}
