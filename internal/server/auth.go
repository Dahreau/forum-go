package server

import (
	"database/sql"
	"forum-go/internal/models"
	"forum-go/internal/shared"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) GetLoginHandler(w http.ResponseWriter, r *http.Request) {
	if s.isLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	render(w, r, "login", nil)
}
func (s *Server) PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := s.db.GetUser(email, password)
	if user.UserId == "" || err != nil {
		render(w, r, "login", map[string]interface{}{"Error": "Invalid username or password. Please try again.", "email": email})
		return
	}
	if user.Role == "ban" {
		render(w, r, "login", map[string]interface{}{"Error": "You are banned", "email": email})
		return
	}
	userID := generateToken(32)

	//Creates cookie session
	expiration := time.Now().Add(time.Hour)
	cookie := http.Cookie{
		Name:    s.SESSION_ID,
		Value:   userID,
		Expires: expiration,
		Path:    "/",
	}
	user.SessionId = sql.NullString{String: userID, Valid: true}
	user.SessionExpire = sql.NullTime{Time: expiration, Valid: true}
	err = s.db.UpdateUser(user)
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
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
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
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
	render(w, r, "register", nil)
}

func (s *Server) PostRegisterHandler(w http.ResponseWriter, r *http.Request) {
	type FormData struct {
		Username string
		Email    string
		Errors   map[string]string
	}
	formData := FormData{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Errors:   make(map[string]string),
	}
	PasswordHash, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	IsUnique, _ := s.db.FindEmailUser(formData.Email)
	if !IsUnique {
		formData.Errors["email_used"] = "Email already used, change it"
	}

	IsUniqueUsername, _ := s.db.FindUsername(r.FormValue("username"))
	if !IsUniqueUsername {
		formData.Errors["username_used"] = "Username already used, change it"
	}
	if len(formData.Username) < 3 {
		formData.Errors["username_len"] = "Username must be at least 3 characters long"
	} else if len(formData.Username) > 20 {
		formData.Errors["username_len"] = "Username must be at most 20 characters long"
	}
	if strings.Contains(formData.Username, " ") {
		formData.Errors["username_spaces"] = "Username must not contain spaces"
	}
	if !IsAlphanumeric(formData.Username) {
		formData.Errors["username_alpha"] = "Username must contain only alphanumeric characters"
	}

	r.FormValue("Confirmpassword")
	if r.FormValue("password") != r.FormValue("Confirmpassword") {
		formData.Errors["password"] = "Passwords don't match"
	}
	if len(formData.Errors) > 0 {
		render(w, r, "register", map[string]interface{}{"FormData": formData})
		return
	}
	user := models.User{Username: r.FormValue("username"), Email: r.FormValue("email"), Password: string(PasswordHash), Role: "user", CreationDate: time.Now(), UserId: shared.ParseUUID(shared.GenerateUUID())}
	err = s.db.CreateUser(user)
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	s.users = append(s.users, user)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (s *Server) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	if !IsAdmin(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	users, err := s.db.GetUsers()
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render(w, r, "../users", map[string]interface{}{"users": users})
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
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	for i, user := range s.users {
		if user.UserId == id {
			s.users = append(s.users[:i], s.users[i+1:]...)
			break
		}
	}
	http.Redirect(w, r, "/adminPanel", http.StatusSeeOther)
}

func (s *Server) BanUserHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	pathParts := strings.Split(path, "/")
	// Check if the path matches the structure
	id := ""
	if len(pathParts) >= 4 && pathParts[2] == "users" {
		id = pathParts[3] // Extract user ID from the path
	}
	userToUpdate := models.User{}
	for _, user := range s.users {
		if user.UserId == id {
			userToUpdate = user
			break
		}
	}
	if userToUpdate.Role == "ban" {
		userToUpdate.Role = "user"
	} else {
		userToUpdate.Role = "ban"
	}
	for i, user := range s.users {
		if user.UserId == id {
			s.users[i] = userToUpdate
			break
		}
	}
	err := s.db.UpdateUser(userToUpdate)
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, "/adminPanel", http.StatusSeeOther)
}

func (s *Server) PromoteUserHandler(w http.ResponseWriter, r *http.Request) {
	if !IsAdmin(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	path := r.URL.Path
	pathParts := strings.Split(path, "/")
	// Check if the path matches the structure
	id := ""
	if len(pathParts) >= 4 && pathParts[2] == "users" {
		id = pathParts[3] // Extract user ID from the path
	}
	userToUpdate := models.User{}
	for _, user := range s.users {
		if user.UserId == id {
			userToUpdate = user
			break
		}
	}
	if userToUpdate.Role == "admin" {
		s.errorHandler(w, r, http.StatusInternalServerError, "User is already an admin")
		return
	} else if userToUpdate.Role == "user" {
		userToUpdate.Role = "moderator"
	} else if userToUpdate.Role == "moderator" {
		userToUpdate.Role = "admin"
	} else {
		s.errorHandler(w, r, http.StatusInternalServerError, "User is banned")
		return
	}
	for i, user := range s.users {
		if user.UserId == id {
			s.users[i] = userToUpdate
			break
		}
	}
	err := s.db.UpdateUser(userToUpdate)
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, "/adminPanel", http.StatusSeeOther)
}

func (s *Server) DemoteUserHandler(w http.ResponseWriter, r *http.Request) {
	if !IsAdmin(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	path := r.URL.Path
	pathParts := strings.Split(path, "/")
	// Check if the path matches the structure
	id := ""
	if len(pathParts) >= 4 && pathParts[2] == "users" {
		id = pathParts[3] // Extract user ID from the path
	}
	userToUpdate := models.User{}
	for _, user := range s.users {
		if user.UserId == id {
			userToUpdate = user
			break
		}
	}
	if userToUpdate.Role == "user" {
		s.errorHandler(w, r, http.StatusInternalServerError, "User is already a user")
		return
	} else if userToUpdate.Role == "moderator" {
		userToUpdate.Role = "user"
	} else if userToUpdate.Role == "admin" {
		s.errorHandler(w, r, http.StatusForbidden, "You can't demote an admin, contact the big boss \"dahreau\" on discord")
		return
	} else {
		s.errorHandler(w, r, http.StatusInternalServerError, "User is banned")
		return
	}
	for i, user := range s.users {
		if user.UserId == id {
			s.users[i] = userToUpdate
			break
		}
	}
	err := s.db.UpdateUser(userToUpdate)
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, "/adminPanel", http.StatusSeeOther)
}
