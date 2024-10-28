package server

import (
	"forum-go/internal/models"
	"net/http"
	"strings"
)

func (s *Server) GetNewPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := s.db.GetPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render(w, r, "../posts", map[string]interface{}{"Posts": posts})
}

func (s *Server) PostNewPostsHandler(w http.ResponseWriter, r *http.Request) {
	post := r.FormValue("Postid")
	if !IsUniquePost(s.posts, post) {
		render(w, r, "../posts", map[string]interface{}{"Posts": s.posts, "Error": "Post already exists"})
		return
	}
	err := s.db.AddPost(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}

func (s *Server) DeletePostsHandler(w http.ResponseWriter, r *http.Request) {
	PostID := r.FormValue("postID")
	err := s.db.DeletePost(PostID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}

func (s *Server) GetPostHandler(w http.ResponseWriter, r *http.Request) {
}
func (s *Server) PostCommentHandler(w http.ResponseWriter, r *http.Request) {
}

func IsUniquePost(posts []models.Post, post string) bool {
	for _, existingPost := range posts {
		if strings.EqualFold(existingPost.PostId, post) {
			return false
		}
	}
	return true
}
