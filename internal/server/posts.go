package server

import (
	"fmt"
	"forum-go/internal/models"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sort" // Import pour trier les posts
	"strconv"
	"strings"
	"time"
)

func (s *Server) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := s.db.GetPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Tri des posts par date de création dans l'ordre décroissant
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreationDate.After(posts[j].CreationDate)
	})

	// Rendu des posts triés
	render(w, r, "../posts", map[string]interface{}{"Posts": posts})
}

func (s *Server) PostNewPostsHandler(w http.ResponseWriter, r *http.Request) {
	type FormData struct {
		Title      string
		Content    string
		Image      string
		Categories []string
		Errors     map[string]string
	}
	erri := r.ParseForm()
	formData := FormData{
		Title:      r.FormValue("title"),
		Content:    r.FormValue("content"),
		Categories: r.Form["categories"],
		Errors:     make(map[string]string),
	}
	if erri != nil {
		log.Println(erri)
	}

	// Validate title
	if ValidateTitle(formData.Title) {
		formData.Errors["Title"] = "Title cannot be empty"
	}

	// Validate content
	if ValidatePostChar(formData.Content) {
		formData.Errors["Content"] = "Content cannot be empty or more than 1000 characters"
	}

	// Validate Categories
	if ValidateCategory(formData.Categories) {
		formData.Errors["Categories"] = "Please select at least one category"
	}

	var imageURL string
	if r.MultipartForm != nil && r.MultipartForm.File["image"] != nil {
		imageURL, erri = UploadImageHandler(w, r)
		if erri != nil {
			formData.Errors["Image"] = "Failed to upload image"
		}
	}

	if len(formData.Errors) > 0 {
		render(w, r, "createPost", map[string]interface{}{"FormData": formData, "Categories": s.categories})
		return
	}

	newPost := models.Post{
		PostId:                strconv.Itoa(rand.Intn(math.MaxInt32)),
		Title:                 formData.Title,
		Content:               formData.Content,
		UserID:                r.FormValue("UserId"),
		ImageURL:              imageURL,
		CreationDate:          time.Now(),
		FormattedCreationDate: time.Now().Format("Jan 02, 2006 - 15:04:05"),
	}

	categories := []models.Category{}
	for _, categoryID := range formData.Categories {
		for _, category := range s.categories {
			if category.CategoryId == categoryID {
				categories = append(categories, category)
			}
		}
	}
	newPost.Categories = categories
	err := s.db.AddPost(newPost, categories)
	s.posts = append(s.posts, newPost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}

func (s *Server) DeletePostsHandler(w http.ResponseWriter, r *http.Request) {
	PostID := r.FormValue("postId")
	fmt.Println(PostID)
	err := s.db.DeletePost(PostID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) EditPostHandler(w http.ResponseWriter, r *http.Request) {
	PostId := r.FormValue("PostId")
	UpdatedContent := r.FormValue("UpdatedContent")

	err := s.db.EditPost(PostId, UpdatedContent)
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, "/post/"+PostId, http.StatusSeeOther)
}

func (s *Server) GetNewPostHandler(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := strings.Split(r.URL.Path, "/")
	if len(vars) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	postID := vars[2]
	post, err := s.db.GetPost(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if post.PostId == "" {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	render(w, r, "detailsPost", map[string]interface{}{"Post": post})
}

func IsUniquePost(posts []models.Post, post string) bool {
	for _, existingPost := range posts {
		if strings.EqualFold(existingPost.PostId, post) {
			return false
		}
	}
	return true
}

const MaxChar = 1000

func ValidatePostChar(content string) bool {
	if len(content) > MaxChar || len(content) == 0 {
		return true
	}
	return false
}

func ValidateTitle(title string) bool {
	return len(title) == 0
}

func ValidateCategory(categories []string) bool {
	return len(categories) < 1
}
