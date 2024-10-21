package server

import (
	"forum-go/internal/models"
	"html/template"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func (s *Server) GetLoginHandler(w http.ResponseWriter, r *http.Request) {
}
func (s *Server) PostLoginHandler(w http.ResponseWriter, r *http.Request) {
}
func (s *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
}
func (s *Server) GetRegisterHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./assets/register.tmpl.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}
func (s *Server) PostRegisterHandler(w http.ResponseWriter, r *http.Request) {
	user := models.User{Username: r.FormValue("username"), Email: r.FormValue("email"), Password: r.FormValue("password"), Role: "user", CreationDate: time.Now(), UserId: strconv.Itoa(rand.Intn(math.MaxInt32))}
	err := s.db.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
