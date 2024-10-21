package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.HomePageHandler)
	mux.HandleFunc("GET /login", s.GetLoginHandler)
	mux.HandleFunc("POST /login", s.PostLoginHandler)
	mux.HandleFunc("POST /logout", s.LogoutHandler)
	mux.HandleFunc("GET /register", s.GetRegisterHandler)
	mux.HandleFunc("POST /register", s.PostRegisterHandler)
	mux.HandleFunc("GET /posts/create", s.GetNewPostsHandler)
	mux.HandleFunc("POST /posts/create", s.PostNewPostsHandler)
	mux.HandleFunc("GET /post/{id}", s.GetPostHandler)
	mux.HandleFunc("POST /post/{id}/comments", s.PostCommentHandler)
	mux.HandleFunc("/health", s.healthHandler)

	return mux
}

func (s *Server) HomePageHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

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

func (s *Server) GetLoginHandler(w http.ResponseWriter, r *http.Request) {
}
func (s *Server) PostLoginHandler(w http.ResponseWriter, r *http.Request) {
}
func (s *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
}
func (s *Server) GetRegisterHandler(w http.ResponseWriter, r *http.Request) {
}
func (s *Server) PostRegisterHandler(w http.ResponseWriter, r *http.Request) {
}
func (s *Server) GetNewPostsHandler(w http.ResponseWriter, r *http.Request) {
}
func (s *Server) PostNewPostsHandler(w http.ResponseWriter, r *http.Request) {
}
func (s *Server) GetPostHandler(w http.ResponseWriter, r *http.Request) {
}
func (s *Server) PostCommentHandler(w http.ResponseWriter, r *http.Request) {
}
