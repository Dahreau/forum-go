package server

import (
	"encoding/base64"
	"fmt"
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
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := s.db.GetUser(email, password)
	if user.UserId == "" || err != nil {
		t, err := template.ParseFiles("./assets/login.tmpl.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		t.Execute(w, map[string]string{"Error": "Invalid username or password. Please try again."})
		return
	}
	//Simulates login
	userID := generateToken(32)

	//Creates cookie session
	expiration := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{
		Name:    s.SESSION_ID,
		Value:   userID,
		Expires: expiration,
		Path:    "/",
	}
	fmt.Println("Cookie: ", cookie)
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

func (s *Server) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	users, err := s.db.GetUsers()
	fmt.Println("Exec time for GetUsersHandler: ", time.Since(start))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t, err := template.ParseFiles("./assets/users.tmpl.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, users)
}

func generateToken(lenght int) string {
	bytes := make([]byte, lenght)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)

}
