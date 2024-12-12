package server

import (
	"fmt"
	"forum-go/internal/models"
	"forum-go/internal/shared"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

func (s *Server) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := s.db.GetPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Sorting posts by creation date in descending order
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreationDate.After(posts[j].CreationDate)
	})

	// Render sorted posts
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
	erri := r.ParseMultipartForm(20 * 1024 * 1024) // 20MB limit
	formData := FormData{
		Title:      r.FormValue("title"),
		Content:    r.FormValue("content"),
		Categories: r.Form["categories"],
		Errors:     make(map[string]string),
	}
	if erri != nil {
		log.Println(erri)
	}
	fmt.Println(formData)

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

	// Handle image upload
	var imageURL string
	if r.MultipartForm != nil && r.MultipartForm.File["file"] != nil {
		imageURL, erri = UploadImageHandler(w, r)
		if erri != nil {
			formData.Errors["Image"] = erri.Error()
		}
	}

	// Check for errors
	if len(formData.Errors) > 0 {
		render(w, r, "createPost", map[string]interface{}{"FormData": formData, "Categories": s.categories})
		return
	}

	// Create new post
	newPost := models.Post{
		PostId:                shared.ParseUUID(shared.GenerateUUID()),
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
	newActivity := models.NewActivity(newPost.UserID, newPost.UserID, string(models.POST_CREATED), newPost.PostId, "", newPost.Title)
	s.db.CreateActivity(newActivity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/post/"+newPost.PostId, http.StatusSeeOther)
}

func (s *Server) DeletePostsHandler(w http.ResponseWriter, r *http.Request) {
	PostID := r.FormValue("postId")
	// Fetch the post to get the image path
	post, err := s.db.GetPost(PostID)
	if err != nil {
		http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
		return
	}

	// Delete the image file if it exists
	if post.ImageURL != "" {
		err = os.Remove("assets/img/uploads/" + post.ImageURL)
		fmt.Println(post.ImageURL)
		if err != nil && !os.IsNotExist(err) { // Ignore errors if the file doesn't exist
			log.Printf("Failed to delete image file: %v\n", err)
		}
	}

	// Delete the post from the database
	err = s.db.DeletePost(PostID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) EditPostHandler(w http.ResponseWriter, r *http.Request) {
	// Edit a post
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
	// GetNewPostHandler handles the new post page
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
	// GetPostHandler handles the post page
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
	post.HasVoted = GetUserVote(post, s.getUser(r).UserId)
	for i, comment := range post.Comments {
		post.Comments[i].HasVoted = GetUserVote(comment, s.getUser(r).UserId)
	}
	if post.ImageURL != "" {
		render(w, r, "detailsPost", map[string]interface{}{"Post": post, "ImageURL": post.ImageURL})
		return
	}
	render(w, r, "detailsPost", map[string]interface{}{"Post": post})
}

func IsUniquePost(posts []models.Post, post string) bool {
	// Check if the post is unique
	for _, existingPost := range posts {
		if strings.EqualFold(existingPost.PostId, post) {
			return false
		}
	}
	return true
}

const MaxChar = 1000

func ValidatePostChar(content string) bool {
	// Validate post character length
	if len(content) > MaxChar || len(content) == 0 {
		return true
	}
	return false
}

func ValidateTitle(title string) bool {
	// Validate title character length
	return len(title) == 0
}

func ValidateCategory(categories []string) bool {
	// Validate categories
	return len(categories) < 1
}
