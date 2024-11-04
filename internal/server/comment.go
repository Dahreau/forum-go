package server

import (
	"fmt"
	"forum-go/internal/models"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Implement function : retrieve form values and call AddCommet function in database/comment.go
// Add model instance
func (s *Server) PostCommentHandler(w http.ResponseWriter, r *http.Request) {
	type CommentData struct {
		Content string
		UserID  string
		PostID  string
		Errors  map[string]string
	}

	commentData := CommentData{
		Content: r.FormValue("comment"),
		UserID:  "UserId",
		PostID:  "PostId",
		Errors:  make(map[string]string),
	}

	if ValidateCommentChar(commentData.Content) {
		commentData.Errors["Comment"] = "Comments must have a maximum of 400 characters"
	}
	if len(commentData.Errors) > 0 {
		render(w, r, "detailsPost", map[string]interface{}{"FormData": commentData, "Categories": s.categories})
		return
	}

	newComment := models.Comment{
		CommentId:    strconv.Itoa(rand.Intn(math.MaxInt32)),
		Content:      r.FormValue("comment"),
		CreationDate: time.Now(),
		UserID:       r.FormValue("UserId"),
		PostID:       r.FormValue("PostId"),
		Likes:        0,
		Dislikes:     0,
	}
	err := s.db.AddComment(newComment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/post/"+newComment.PostID, http.StatusSeeOther)
}

func validComment(content string) bool {
	return len(content) < 401
}

func (s *Server) GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := s.db.GetPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render(w, r, "../posts", map[string]interface{}{"Posts": posts})
}

func (s *Server) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	PostID := r.FormValue("postId")
	fmt.Println(PostID)
	err := s.db.DeletePost(PostID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}

func (s *Server) GetNewCommentHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := s.db.GetCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !s.isLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	render(w, r, "createPost", map[string]interface{}{"Categories": categories})
}

func (s *Server) GetCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := strings.Split(r.URL.Path, "/")
	if len(vars) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	postID := vars[2]
	post, err := s.db.GetPost(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if post.PostId == "" {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	render(w, r, "detailsPost", map[string]interface{}{"Post": post})
}

// Handler in the post.go file line 130

const MaxCharComment = 400

func ValidateCommentChar(content string) bool {
	if len(content) > MaxCharComment || len(content) == 0 {
		return true
	}
	return false
}

func (s *Server) EditCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentID := r.FormValue("commentId")
	commentContent := r.FormValue("newCommentContent")

	err := s.db.EditComment(commentID, commentContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}
