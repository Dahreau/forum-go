package server

import (
	"encoding/base64"
	"forum-go/internal/models"
	"log"
	"math/rand"
	"net/http"
	"unicode"
)

func (s *Server) isLoggedIn(r *http.Request) bool {
	user := r.Context().Value(contextKeyUser)
	return user != nil
}
func (s *Server) getUser(r *http.Request) models.User {
	user := r.Context().Value(contextKeyUser)
	if user == nil {
		return models.User{}
	}
	return user.(models.User)
}
func IsAdmin(r *http.Request) bool {
	user := r.Context().Value(contextKeyUser)
	if user == nil {
		return false
	}
	return user.(models.User).Role == "admin"
}

func generateToken(lenght int) string {
	bytes := make([]byte, lenght)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)

}
func IsAlphanumeric(s string) bool {
	for _, char := range s {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}
