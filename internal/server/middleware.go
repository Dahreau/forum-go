package server

import (
	"context"
	"net/http"
)

type contextKey string

const contextKeyUser = contextKey("user")

func (s *Server) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(s.SESSION_ID)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		// Use the cookie value as needed
		user, err := s.db.FindUserCookie(cookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		// Set the user in the request context
		ctx := context.WithValue(r.Context(), contextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) isLoggedIn(r *http.Request) bool {
	user := r.Context().Value(contextKeyUser)
	return user != nil
}
