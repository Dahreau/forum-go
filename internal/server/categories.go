package server

import (
	"forum-go/internal/models"
	"net/http"
	"strings"
)

func (s *Server) GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := s.db.GetCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render(w, "../categories", map[string]interface{}{"Categories": categories})
}

func (s *Server) PostCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	category := r.FormValue("categoryName")
	if !IsUniqueCategory(s.categories, category) {
		render(w, "../categories", map[string]interface{}{"Categories": s.categories, "Error": "Category already exists"})
		return
	}
	err := s.db.AddCategory(category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}

func (s *Server) DeleteCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	categoryID := r.FormValue("categoryId")
	err := s.db.DeleteCategory(categoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}

func (s *Server) EditCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	categoryID := r.FormValue("categoryId")
	categoryName := r.FormValue("newCategoryName")
	if !IsUniqueCategory(s.categories, categoryName) {
		render(w, "../categories", map[string]interface{}{"Categories": s.categories, "Error": "Category already exists"})
		return
	}
	err := s.db.EditCategory(categoryID, categoryName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}

func IsUniqueCategory(categories []models.Category, category string) bool {
	for _, existingCategory := range categories {
		if strings.EqualFold(existingCategory.Name, category) {
			return false
		}
	}
	return true
}
