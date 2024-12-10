package server

import (
	"database/sql"
	"encoding/json"
	"forum-go/internal/models"
	"forum-go/internal/shared"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// DiscordLoginHandler initiates the Discord OAuth flow.
func (s *Server) DiscordLoginHandler(w http.ResponseWriter, r *http.Request) {
	discordClientID := shared.GetEnv("DiscordClientID")
	discordRedirectURI := shared.GetEnv("DiscordRedirectURI")

	// Construct the authorization URL
	authURL := "https://discord.com/oauth2/authorize?client_id=" + discordClientID +
		"&redirect_uri=" + url.QueryEscape(discordRedirectURI) +
		"&response_type=code&scope=identify email"

	// Redirect the user to Discord's authorization page
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// DiscordCallbackHandler handles the callback from Discord after authentication.
func (s *Server) DiscordCallbackHandler(w http.ResponseWriter, r *http.Request) {
	discordClientID := shared.GetEnv("DiscordClientID")
	discordClientSecret := shared.GetEnv("DiscordClientSecret")
	discordRedirectURI := shared.GetEnv("DiscordRedirectURI")

	// Get the authorization code from the query string
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Authorization code is missing", http.StatusBadRequest)
		return
	}

	// Exchange the authorization code for an access token
	tokenResp, err := http.PostForm("https://discord.com/api/oauth2/token", url.Values{
		"client_id":     {discordClientID},
		"client_secret": {discordClientSecret},
		"redirect_uri":  {discordRedirectURI},
		"grant_type":    {"authorization_code"},
		"code":          {code},
	})
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}
	defer tokenResp.Body.Close()

	// Check if the response status is OK
	if tokenResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(tokenResp.Body)
		http.Error(w, "Token exchange failed: "+string(body), http.StatusInternalServerError)
		return
	}

	// Parse the token response
	var tokenData map[string]interface{}
	if err := json.NewDecoder(tokenResp.Body).Decode(&tokenData); err != nil {
		http.Error(w, "Error parsing token response", http.StatusInternalServerError)
		return
	}

	// Get the access token from the response
	accessToken := tokenData["access_token"].(string)
	if accessToken == "" {
		http.Error(w, "Access token is missing", http.StatusInternalServerError)
		return
	}

	// Fetch user information from Discord
	req, err := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
	if err != nil {
		http.Error(w, "Failed to create user info request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, "Failed to fetch user info: "+string(body), http.StatusInternalServerError)
		return
	}

	// Parse the user information response
	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Error parsing user info", http.StatusInternalServerError)
		return
	}

	// Extract user details
	email, emailOk := userInfo["email"].(string)
	username, usernameOk := userInfo["username"].(string)
	if !emailOk || !usernameOk || email == "" || username == "" {
		http.Error(w, "Email or username is missing", http.StatusInternalServerError)
		return
	}

	// Check if the email exists in the database
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
		if user.Provider != "discord" {
			render(w, r, "login", map[string]interface{}{"Error": "Email already used by another provider", "email": email})
			return
		}

		// Log the user in by creating a session
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

	// Create a new user
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
		Provider:     "discord",
	}

	err = s.db.CreateUser(user)
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// Log the new user in
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
