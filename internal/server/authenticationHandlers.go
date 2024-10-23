package server

import (
	"database/sql"
	"encoding/base64"
	"forum-go/internal/models"
	"html/template"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) GetLoginHandler(w http.ResponseWriter, r *http.Request) {
	if s.isLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	render(w, "login", nil)
}
func (s *Server) PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := s.db.GetUser(email, password)
	if user.UserId == "" || err != nil {

		render(w, "login", map[string]interface{}{"Error": "Invalid username or password. Please try again.", "email": email})
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
	user.SessionId = sql.NullString{String: userID, Valid: true}
	user.SessionExpired = sql.NullTime{Time: expiration, Valid: true}
	err = s.db.UpdateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	//Creates a cookie with the same name
	cookie := http.Cookie{
		Name:     s.SESSION_ID,    // Cookie name
		Value:    "",              // EMpty value to delete it
		Expires:  time.Unix(0, 0), // Set expiration date in the past
		HttpOnly: true,
		Path:     "/", // Cookie path
	}

	// Deletes cookie
	http.SetCookie(w, &cookie)

	// Redirect to home
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) GetRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if s.isLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	render(w, "register", nil)
}

func (s *Server) PostRegisterHandler(w http.ResponseWriter, r *http.Request) {
	PasswordHash, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	IsUnique, _ := s.db.FindEmailUser(r.FormValue("email"))
	if !IsUnique {
		t, err := template.ParseFiles("./assets/register.tmpl.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, map[string]string{"email_used": "Email already used, change it"})
		return
	}

	IsUniqueUsername, _ := s.db.FindUsername(r.FormValue("email"))
	if !IsUniqueUsername {
		t, err := template.ParseFiles("./assets/register.tmpl.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, map[string]string{"username_used": "Username already used, change it"})
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
	users, err := s.db.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render(w, "../users", map[string]interface{}{"users": users})
}

func generateToken(lenght int) string {
	bytes := make([]byte, lenght)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)

}

func (s *Server) DeleteUsersHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	pathParts := strings.Split(path, "/")
	// Check if the path matches the structure
	id := ""
	if len(pathParts) >= 4 && pathParts[2] == "users" {
		id = pathParts[3] // Extract user ID from the path
	}
	err := s.db.DeleteUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
