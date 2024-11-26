package server

import (
	"forum-go/internal/models"
	"net/http"
)

func (s *Server) ModRequestsHandler(w http.ResponseWriter, r *http.Request) {
	if !IsAdmin(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	modRequests, err := s.db.GetRequests()
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render(w, r, "admin/requests", map[string]interface{}{"modRequests": modRequests})
}

func (s *Server) GetModRequestHandler(w http.ResponseWriter, r *http.Request) {
	if !s.isLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	Requests, err := s.db.GetRequests()
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	UserRequests := []models.Request{}
	HasPendingRequest := false
	for _, request := range Requests {
		if request.UserId == s.getUser(r).UserId {
			UserRequests = append(UserRequests, request)
			if request.Status == "pending" {
				HasPendingRequest = true
			}
		}
	}
	if len(UserRequests) == 0 {
		render(w, r, "modRequest", map[string]interface{}{"HasPendingRequest": false})
		return
	}
	render(w, r, "modRequest", map[string]interface{}{"UserRequests": UserRequests, "HasPendingRequest": HasPendingRequest})
}

func (s *Server) PostModRequestHandler(w http.ResponseWriter, r *http.Request) {
	if !s.isLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	content := r.FormValue("content")
	userid := r.FormValue("userid")
	username := r.FormValue("username")
	request := models.NewRequest(userid, username, content)
	err := s.db.CreateRequest(request)
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, "modRequest", http.StatusSeeOther)
}

func (s *Server) AcceptRequestHandler(w http.ResponseWriter, r *http.Request) {
	if !IsAdmin(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	requestId := r.FormValue("request_id")
	userid := r.FormValue("user_id")
	err := s.db.UpdateRequestStatus(requestId, "accepted")
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	for i, user := range s.users {
		if user.UserId == userid {
			s.users[i].Role = "moderator"
			err = s.db.UpdateUser(s.users[i])
			if err != nil {
				s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
				return
			}
			break
		}
	}
	http.Redirect(w, r, "../adminPanel/modrequests", http.StatusSeeOther)
}

func (s *Server) RejectRequestHandler(w http.ResponseWriter, r *http.Request) {
	if !IsAdmin(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	requestId := r.FormValue("request_id")
	err := s.db.UpdateRequestStatus(requestId, "rejected")
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, "../adminPanel/modrequests", http.StatusSeeOther)
}
