package server

import (
	"encoding/base64"
	"forum-go/internal/models"
	"html/template"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) GetLoginHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./assets/login.tmpl.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)

}

func (s *Server) PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	//Simulates login
	userID := generateToken(32)
	sessionID := generateToken(32)

	//Creates cookie session
	expiration := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{
		Name:     sessionID,
		Value:    userID,
		Expires:  expiration,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
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
	PasswordHash, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := models.User{Username: r.FormValue("username"), Email: r.FormValue("email"), Password: string(PasswordHash), Role: "user", CreationDate: time.Now(), UserId: strconv.Itoa(rand.Intn(math.MaxInt32))}
	err = s.db.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func generateToken(lenght int) string {
	bytes := make([]byte, lenght)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}
