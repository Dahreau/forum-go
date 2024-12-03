package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"forum-go/internal/models"
	"forum-go/internal/shared"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()

	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	mux.HandleFunc("/", s.HomePageHandler)
	mux.HandleFunc("/about", s.AboutPageHandler)

	mux.HandleFunc("/activity", s.ActivityPageHandler)

	mux.HandleFunc("GET /login", s.GetLoginHandler)
	mux.HandleFunc("POST /login", s.PostLoginHandler)

	mux.HandleFunc("POST /logout", s.LogoutHandler)

	mux.HandleFunc("GET /register", s.GetRegisterHandler)
	mux.HandleFunc("POST /register", s.PostRegisterHandler)

	mux.HandleFunc("GET /delete/users/{id}", s.DeleteUsersHandler)
	mux.HandleFunc("GET /ban/users/{id}", s.BanUserHandler)
	mux.HandleFunc("GET /promote/users/{id}", s.PromoteUserHandler)
	mux.HandleFunc("GET /demote/users/{id}", s.DemoteUserHandler)

	mux.HandleFunc("GET /posts/create", s.GetNewPostHandler)
	mux.HandleFunc("POST /posts/create", s.PostNewPostsHandler)
	mux.HandleFunc("POST /posts/delete/{id}", s.DeletePostsHandler)
	mux.HandleFunc("POST /posts/edit/{id}", s.EditPostHandler)

	mux.HandleFunc("GET /categories", s.GetCategoriesHandler)
	mux.HandleFunc("POST /categories/add", s.PostCategoriesHandler)
	mux.HandleFunc("POST /categories/delete/{id}", s.DeleteCategoriesHandler)
	mux.HandleFunc("POST /categories/edit/{id}", s.EditCategoriesHandler)

	mux.HandleFunc("GET /post/{id}", s.GetPostHandler)
	mux.HandleFunc("POST /comment/delete/{id}", s.DeleteCommentHandler)
	mux.HandleFunc("POST /comment/edit/{id}", s.EditCommentHandler)
	mux.HandleFunc("POST /post/comment", s.PostCommentHandler)

	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("GET /adminPanel", s.AdminPanelHandler)
	mux.HandleFunc("GET /report/{id}", s.GetReportHandler)
	mux.HandleFunc("POST /report", s.PostReportHandler)
	mux.HandleFunc("GET /adminPanel/modrequests", s.ModRequestsHandler)
	mux.HandleFunc("POST /vote", s.VoteHandler)

	mux.HandleFunc("POST /modRequest", s.PostModRequestHandler)
	mux.HandleFunc("GET /modRequest", s.GetModRequestHandler)
	mux.HandleFunc("POST /modRequest/accepted", s.AcceptRequestHandler)
	mux.HandleFunc("POST /modRequest/rejected", s.RejectRequestHandler)

	mux.HandleFunc("GET /adminPanel/reports", s.GetReportsHandler)
	mux.HandleFunc("POST /reports/accepted", s.AcceptReportHandler)
	mux.HandleFunc("POST /reports/rejected", s.RejectReportHandler)

	// AUTH ROUTES
	mux.HandleFunc("/auth/google", s.GoogleLoginHandler)
	mux.HandleFunc("/auth/google/callback", s.GoogleCallbackHandler)
	// mux.HandleFunc("/auth/github", s.GithubLoginHandler)
	// mux.HandleFunc("/auth/github/callback", s.GithubCallbackHandler)

	return s.authenticate(mux)
}

func (s *Server) VoteHandler(w http.ResponseWriter, r *http.Request) {
	if !s.isLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	postID := r.FormValue("post_id")
	userID := r.FormValue("user_id")
	vote := r.FormValue("vote")
	commentID := r.FormValue("comment_id")
	var isLike bool
	if vote == "like" {
		isLike = true
	} else {
		isLike = false
	}

	err := s.db.Vote(postID, commentID, userID, isLike)
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	post, err := s.db.GetPost(postID)
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
	}
	if commentID == "" {
		if isLike {
			if post.UserID != s.getUser(r).UserId {
				newActivity := models.NewActivity(post.UserID, userID, string(models.GET_POST_LIKED), postID, "", post.Title)
				err = s.db.CreateActivity(newActivity)
				if err != nil {
					s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
				}
			}
			newActivity := models.NewActivity(s.getUser(r).UserId, userID, string(models.POST_LIKED), postID, "", post.Title)
			err = s.db.CreateActivity(newActivity)
			if err != nil {
				s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
			}
		} else {
			if post.UserID != s.getUser(r).UserId {
				newActivity := models.NewActivity(post.UserID, userID, string(models.GET_POST_DISLIKED), postID, "", post.Title)
				err = s.db.CreateActivity(newActivity)
				if err != nil {
					s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
				}
			}
			newActivity := models.NewActivity(s.getUser(r).UserId, userID, string(models.POST_DISLIKED), postID, "", post.Title)
			err = s.db.CreateActivity(newActivity)
			if err != nil {
				s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
			}
		}
	} else {
		ActualComment := models.Comment{}
		for _, comment := range post.Comments {
			if comment.CommentId == commentID {
				ActualComment = comment
			}
		}
		if isLike {
			if ActualComment.UserID != s.getUser(r).UserId {
				newActivity := models.NewActivity(ActualComment.UserID, userID, string(models.GET_COMMENT_LIKED), postID, ActualComment.CommentId, ActualComment.Content)
				err = s.db.CreateActivity(newActivity)
				if err != nil {
					s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
				}
			}
			newActivity := models.NewActivity(s.getUser(r).UserId, userID, string(models.COMMENT_LIKED), postID, ActualComment.CommentId, ActualComment.Content)
			err = s.db.CreateActivity(newActivity)
			if err != nil {
				s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
			}
		} else {
			if ActualComment.UserID != s.getUser(r).UserId {
				newActivity := models.NewActivity(ActualComment.UserID, userID, string(models.GET_COMMENT_DISLIKED), postID, ActualComment.CommentId, ActualComment.Content)
				err = s.db.CreateActivity(newActivity)
				if err != nil {
					s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
				}
			}
			newActivity := models.NewActivity(s.getUser(r).UserId, userID, string(models.COMMENT_DISLIKED), postID, ActualComment.CommentId, ActualComment.Content)
			err = s.db.CreateActivity(newActivity)
			if err != nil {
				s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
			}
		}
	}
	referer := r.Header.Get("Referer")
	if referer != "" {
		http.Redirect(w, r, referer, http.StatusSeeOther)
		return
	}
}

func (s *Server) GetReportHandler(w http.ResponseWriter, r *http.Request) {
	if !s.isLoggedIn(r) || (!IsAdmin(r) && !IsModerator(r)) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	id := r.URL.Path[len("/report/"):]
	ReportPost := models.Post{}
	for _, post := range s.posts {
		if post.PostId == id {
			ReportPost = post
			break
		}
	}
	if ReportPost.PostId == "" {
		s.errorHandler(w, r, http.StatusNotFound, "Post not found")
		return
	}

	render(w, r, "report", map[string]interface{}{"post": ReportPost})

}

func (s *Server) PostReportHandler(w http.ResponseWriter, r *http.Request) {
	if !s.isLoggedIn(r) || (!IsAdmin(r) && !IsModerator(r)) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	postID := r.FormValue("postid")
	content := r.FormValue("content")
	reason := r.FormValue("reason")
	username := r.FormValue("username")
	userid := r.FormValue("userid")
	report := models.NewReport(userid, username, postID, content, reason)
	err := s.db.CreateReport(report)
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) AdminPanelHandler(w http.ResponseWriter, r *http.Request) {
	if !s.isLoggedIn(r) || !IsAdmin(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	users, err := s.db.GetUsers()
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render(w, r, "admin/adminPanel", map[string]interface{}{"users": users})
}

func (s *Server) HomePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/created" && r.URL.Path != "/liked" {
		s.errorHandler(w, r, http.StatusNotFound, "Page not found")
		return
	}
	if !s.isLoggedIn(r) && (r.URL.Path == "/created" || r.URL.Path == "/liked") {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	err := error(nil)
	s.categories, err = s.db.GetCategories()
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	s.posts, err = s.db.GetPosts()
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	for i, post := range s.posts {
		s.posts[i].HasVoted = GetUserVote(post, s.getUser(r).UserId)
	}
	postsToRender := []models.Post{}
	if r.URL.Path == "/created" {
		for _, post := range s.posts {
			if post.UserID == s.getUser(r).UserId {
				postsToRender = append(postsToRender, post)
			}
		}
	} else if r.URL.Path == "/liked" {
		for _, post := range s.posts {
			if post.HasVoted == 1 {
				postsToRender = append(postsToRender, post)
			}
		}
	} else {
		postsToRender = s.posts
	}

	render(w, r, "home", map[string]interface{}{"Categories": s.categories, "Posts": postsToRender})
}
func (s *Server) AboutPageHandler(w http.ResponseWriter, r *http.Request) {
	render(w, r, "about", nil)
}
func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"
	for k, v := range r.Header {
		resp[k] = fmt.Sprintf("%v", v)
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(s.db.Health())

	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) errorHandler(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.WriteHeader(status)
	error := models.Error{Message: message, StatusCode: status}
	render(w, r, "error", map[string]interface{}{"Error": error})
}

// Auth
// Google

var googleClientID = "googleClientID"         // TODO Put googleClientID
var googleClientSecret = "googleClientSecret" // TODO Put googleClientSecret
var googleRedirectURL = "http://localhost:8080/auth/google/callback"

func (s *Server) GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := "https://accounts.google.com/o/oauth2/auth?client_id=" + googleClientID +
		"&redirect_uri=" + googleRedirectURL +
		"&response_type=code&scope=email%20profile&state=state"
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (s *Server) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère le code d'autorisation
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code d'autorisation manquant", http.StatusBadRequest)
		log.Println("Code d'autorisation manquant")
		return
	}

	// Échange le code d'autorisation contre un token d'accès
	tokenResp, err := http.PostForm("https://oauth2.googleapis.com/token", url.Values{
		"client_id":     {googleClientID},
		"client_secret": {googleClientSecret},
		"redirect_uri":  {googleRedirectURL},
		"grant_type":    {"authorization_code"},
		"code":          {code},
	})
	if err != nil {
		http.Error(w, "Échec lors de l'échange du token : "+err.Error(), http.StatusInternalServerError)
		log.Println("Erreur HTTP POST lors de l'échange du token :", err)
		return
	}
	defer tokenResp.Body.Close()

	// Vérifie le statut HTTP
	if tokenResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(tokenResp.Body)
		log.Printf("Erreur lors de l'échange du token: %s", body)
		http.Error(w, "Erreur lors de l'échange du token avec Google", http.StatusInternalServerError)
		return
	}

	// Décoder la réponse JSON contenant le token
	var tokenData map[string]interface{}
	if err := json.NewDecoder(tokenResp.Body).Decode(&tokenData); err != nil {
		http.Error(w, "Erreur lors du parsing du token : "+err.Error(), http.StatusInternalServerError)
		log.Println("Erreur JSON lors du parsing :", err)
		return
	}
	log.Printf("Réponse JSON token: %+v\n", tokenData)

	// Vérifie et extrait le token d'accès
	accessToken, ok := tokenData["access_token"]
	if !ok || accessToken == nil {
		http.Error(w, "Token d'accès manquant ou invalide", http.StatusInternalServerError)
		log.Println("Token d'accès non trouvé ou null dans la réponse JSON")
		return
	}

	accessTokenStr, ok := accessToken.(string)
	if !ok {
		http.Error(w, "Token d'accès n'est pas une chaîne valide", http.StatusInternalServerError)
		log.Println("Token d'accès non convertible en string")
		return
	}

	// Utilise le token d'accès pour récupérer les infos utilisateur
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		http.Error(w, "Erreur lors de la création de la requête utilisateur : "+err.Error(), http.StatusInternalServerError)
		log.Println("Erreur HTTP GET lors de la récupération des infos utilisateur :", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+accessTokenStr)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Erreur lors de la requête utilisateur : "+err.Error(), http.StatusInternalServerError)
		log.Println("Erreur HTTP GET lors de la récupération des infos utilisateur :", err)
		return
	}
	defer resp.Body.Close()

	// Vérifie le statut HTTP
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Erreur lors de la récupération des infos utilisateur: %s", body)
		http.Error(w, "Erreur lors de la récupération des infos utilisateur", http.StatusInternalServerError)
		return
	}

	// Décoder les infos utilisateur
	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Erreur lors du parsing des infos utilisateur : "+err.Error(), http.StatusInternalServerError)
		log.Println("Erreur JSON lors du parsing des infos utilisateur :", err)
		return
	}
	log.Printf("Infos utilisateur: %+v\n", userInfo)

	// Récupère les informations utilisateur
	email := userInfo["email"].(string)
	name := userInfo["name"].(string)

	// Vérifiez si l'adresse email existe déjà dans la base de données
	IsUnique, err := s.db.FindEmailUser(email)
	if err != nil {
		http.Error(w, "Erreur lors de la vérification de l'utilisateur : "+err.Error(), http.StatusInternalServerError)
		log.Println("Erreur lors de la vérification de l'utilisateur :", err)
		return
	}

	if !IsUnique {
		// L'adresse email existe déjà dans la base de données
		log.Println("L'adresse email existe déjà dans la base de données")

		// Récupérer l'utilisateur depuis la base de données
		user, err := s.db.FindUserByEmail(email)
		if err != nil {
			http.Error(w, "Erreur lors de la récupération de l'utilisateur : "+err.Error(), http.StatusInternalServerError)
			log.Println("Erreur lors de la récupération de l'utilisateur :", err)
			return
		}

		if user.Role == "ban" {
			render(w, r, "login", map[string]interface{}{"Error": "You are banned", "email": email})
			return
		}

		userID := shared.ParseUUID(shared.GenerateUUID())

		// Créez une session pour l'utilisateur
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

		// Redirigez l'utilisateur vers la page d'accueil ou une autre page appropriée
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Si l'adresse email n'existe pas, continuez avec le reste de la logique
	log.Println("L'adresse email n'existe pas dans la base de données")

	// Générer un mot de passe unique
	password := shared.GenerateUUID().String()

	// Hacher le mot de passe
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// Créez un nouvel utilisateur
	user := models.User{
		Username:     name,
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
	s.users = append(s.users, user)

	// Redirigez l'utilisateur vers la page de connexion ou une autre page appropriée
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
