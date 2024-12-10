package server

import (
	"forum-go/internal/models"
	"net/http"
)

func (s *Server) ActivityPageHandler(w http.ResponseWriter, r *http.Request) {
	// ActivityPageHandler handles the activity page
	if !s.isLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	user, ok := r.Context().Value(contextKeyUser).(models.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}
	render(w, r, "activity", nil)
	for i := range user.Activities {
		user.Activities[i].IsRead = true
	}
	s.db.ReadActivites(user.UserId)
}
