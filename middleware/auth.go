package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
	"personal_blog/config"
)

type Authenticator struct {
	config *config.Config
}

func NewAuthenticator(config *config.Config) *Authenticator {
	return &Authenticator{
		config: config,
	}
}

// BasicAuth middleware
func (a *Authenticator) BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || !a.validateCredentials(username, password) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (a *Authenticator) validateCredentials(username, password string) bool {
	usernameHash := sha256.Sum256([]byte(username))
	passwordHash := sha256.Sum256([]byte(password))
	expectedUsernameHash := sha256.Sum256([]byte(a.config.Auth.Username))
	expectedPasswordHash := sha256.Sum256([]byte(a.config.Auth.Password))

	usernameMatch := subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1
	passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1

	return usernameMatch && passwordMatch
}
