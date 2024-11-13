package server

import (
	"forum-go/internal/models"
	"net/http"
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
		CommentId:    ParseUUID(GenerateUUID()),
		Content:      r.FormValue("comment"),
		CreationDate: time.Now(),
		UserID:       r.FormValue("UserId"),
		PostID:       r.FormValue("PostId"),
		Likes:        0,
		Dislikes:     0,
	}
	err := s.db.AddComment(newComment)
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, "/post/"+newComment.PostID, http.StatusSeeOther)
}

func (s *Server) GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := s.db.GetPosts()
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render(w, r, "../posts", map[string]interface{}{"Posts": posts})
}

func (s *Server) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	PostID := r.FormValue("PostId")
	CommentID := r.FormValue("CommentId")
	UserID := r.FormValue("UserId")
	if !s.isLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if UserID != s.getUser(r).UserId && !IsAdmin(r) {
		http.Error(w, "You are not allowed to delete this comment", http.StatusForbidden)
		return
	}
	if CommentID == "" {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}
	if PostID == "" {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	err := s.db.DeleteComment(CommentID)
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, "/post/"+PostID, http.StatusSeeOther)
}

func (s *Server) GetNewCommentHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := s.db.GetCategories()
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	if !s.isLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	render(w, r, "createPost", map[string]interface{}{"Categories": categories})
}

const MaxCharComment = 400

func ValidateCommentChar(content string) bool {
	if len(content) > MaxCharComment || len(content) == 0 {
		return true
	}
	return false
}

func (s *Server) EditCommentHandler(w http.ResponseWriter, r *http.Request) {
	CommentID := r.FormValue("CommentId")
	PostId := r.FormValue("PostId")
	UpdatedContent := r.FormValue("UpdatedContent")

	err := s.db.EditComment(CommentID, UpdatedContent)
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, "/post/"+PostId, http.StatusSeeOther)
}
