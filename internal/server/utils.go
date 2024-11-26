package server

import (
	"forum-go/internal/models"
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

func IsAlphanumeric(s string) bool {
	for _, char := range s {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}

func GetUserVote(post_comment models.Post_Comment, userId string) int {
	for _, like := range post_comment.GetUserLikes() {
		if like.UserId == userId {
			if like.IsLike {
				return 1
			}
			return -1
		}
	}
	return 0
}
