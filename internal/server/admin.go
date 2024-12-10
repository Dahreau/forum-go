package server

import (
	"forum-go/internal/models"
	"net/http"
)

func (s *Server) ModRequestsHandler(w http.ResponseWriter, r *http.Request) {
	// ModRequestsHandler handles the moderator requests page
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
	// GetModRequestHandler handles the moderator request page
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
	// PostModRequestHandler handles the moderator request form submission
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
	// AcceptRequestHandler handles the moderator request acceptance
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
	//	RejectRequestHandler handles the moderator request rejection
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

func (s *Server) GetReportsHandler(w http.ResponseWriter, r *http.Request) {
	// GetReportsHandler handles the reports page
	if !IsAdmin(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	Reports, err := s.db.GetReports()
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	for i, report := range Reports {
		for _, post := range s.posts {
			if post.PostId == report.PostId {
				Reports[i].Post = post
				break
			}
		}
	}
	render(w, r, "admin/reports", map[string]interface{}{"Reports": Reports})
}

func (s *Server) AcceptReportHandler(w http.ResponseWriter, r *http.Request) {
	// AcceptReportHandler handles the report acceptance
	if !IsAdmin(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	postid := r.FormValue("postid")
	err := s.db.DeletePost(postid)
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, "../adminPanel/reports", http.StatusSeeOther)
}

func (s *Server) RejectReportHandler(w http.ResponseWriter, r *http.Request) {
	// RejectReportHandler handles the report rejection
	if !IsAdmin(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	reportId := r.FormValue("reportid")
	err := s.db.UpdateReportStatus(reportId, "rejected")
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, "../adminPanel/reports", http.StatusSeeOther)
}