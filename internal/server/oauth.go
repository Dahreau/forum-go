package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"forum-go/internal/models"
	"forum-go/internal/shared"
	"golang.org/x/crypto/bcrypt"
)

var (
	clientID     = "SECRET"
	clientSecret = "SECRET"
	redirectURI  = "http://localhost:8080/auth/github/callback"
)

// GithubLoginHandler initiates the GitHub OAuth flow.
func (s *Server) GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	authURL := "https://github.com/login/oauth/authorize?client_id=" + clientID +
		"&redirect_uri=" + url.QueryEscape(redirectURI) +
		"&scope=user:email"

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// GithubCallbackHandler handles the callback from GitHub.
func (s *Server) GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Authorization code is missing", http.StatusBadRequest)
		return
	}

	// Exchange the authorization code for an access token
	tokenResp, err := http.PostForm("https://github.com/login/oauth/access_token", url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"redirect_uri":  {redirectURI},
		"code":          {code},
	})
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}
	defer tokenResp.Body.Close()

	if tokenResp.StatusCode != http.StatusOK {
		http.Error(w, "Token exchange failed", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(tokenResp.Body)
	if err != nil {
		http.Error(w, "Error reading token response", http.StatusInternalServerError)
		return
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		http.Error(w, "Error parsing token response", http.StatusInternalServerError)
		return
	}

	accessToken := values.Get("access_token")
	if accessToken == "" {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect) // Restart the OAuth flow
		return
	}

	// Fetch user information
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Error fetching user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Error parsing user info", http.StatusInternalServerError)
		return
	}

	// Extract user details
	email, emailOk := userInfo["email"].(string)
	username, usernameOk := userInfo["login"].(string)

	if !usernameOk || username == "" {
		http.Error(w, "Failed to retrieve username", http.StatusInternalServerError)
		return
	}
	email, errMail := getMail(accessToken)
	log.Println(email, emailOk, errMail)
	if !emailOk || email == "" {
		// If no email exists, use the GitHub username as a fallback for account uniqueness
		email = fmt.Sprintf("%s@github.local", username) // Fake email to ensure unique account creation
	}

	// Check if the email already exists in the database
	IsUnique, err := s.db.FindEmailUser(email)
	if err != nil {
		http.Error(w, "Error checking user existence", http.StatusInternalServerError)
		return
	}

	if !IsUnique {
		user, err := s.db.FindUserByEmail(email)
		if err != nil {
			http.Error(w, "Error fetching user", http.StatusInternalServerError)
			return
		}

		if user.Role == "ban" {
			render(w, r, "login", map[string]interface{}{"Error": "You are banned", "email": email})
			return
		}

		// Automatically log the user in by creating a session
		userID := shared.ParseUUID(shared.GenerateUUID())
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
		return
	}

	// Create a new user if the email is not found
	password := shared.GenerateUUID().String()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	user := models.User{
		Username:     username,
		Email:        email,
		Password:     string(passwordHash),
		Role:         "user",
		CreationDate: time.Now(),
		UserId:       shared.ParseUUID(shared.GenerateUUID()),
	}
	err = s.db.CreateUser(user)
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// Automatically log the new user in
	userID := user.UserId
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

// getMail retrieves the user's primary email from GitHub.
func getMail(accessToken string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch emails")
	}

	var emails []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

	for _, email := range emails {
		if primary, ok := email["primary"].(bool); ok && primary {
			if emailAddr, ok := email["email"].(string); ok {
				return emailAddr, nil
			}
		}
	}

	return "", fmt.Errorf("no primary email found")
}
