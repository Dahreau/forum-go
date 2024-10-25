package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	mux.HandleFunc("/", s.HomePageHandler)

	mux.HandleFunc("GET /login", s.GetLoginHandler)
	mux.HandleFunc("POST /login", s.PostLoginHandler)

	mux.HandleFunc("POST /logout", s.LogoutHandler)

	mux.HandleFunc("GET /register", s.GetRegisterHandler)
	mux.HandleFunc("POST /register", s.PostRegisterHandler)

	mux.HandleFunc("GET /users", s.GetUsersHandler)
	mux.HandleFunc("GET /delete/users/{id}", s.DeleteUsersHandler)

	mux.HandleFunc("GET /posts/create", s.GetNewPostsHandler)
	mux.HandleFunc("POST /posts/create", s.PostNewPostsHandler)

	mux.HandleFunc("GET /categories", s.GetCategoriesHandler)
	mux.HandleFunc("POST /categories/add", s.PostCategoriesHandler)
	mux.HandleFunc("POST /categories/delete/{id}", s.DeleteCategoriesHandler)
	mux.HandleFunc("POST /categories/edit/{id}", s.EditCategoriesHandler)

	mux.HandleFunc("GET /post/{id}", s.GetPostHandler)
	mux.HandleFunc("POST /post/comment", s.PostCommentHandler)
	mux.HandleFunc("/health", s.healthHandler)

	return s.authenticate(mux)
}

func (s *Server) HomePageHandler(w http.ResponseWriter, r *http.Request) {
	render(w, "../index", map[string]interface{}{"User": r.Context().Value(contextKeyUser)})
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
