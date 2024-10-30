package server

import (
	"database/sql"
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
	render(w, r, "../categories", map[string]interface{}{"Categories": categories})
}

func (s *Server) PostCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	category := r.FormValue("categoryName")
	if !IsUniqueCategory(s.categories, category) {
		render(w, r, "../categories", map[string]interface{}{"Categories": s.categories, "Error": "Category already exists"})
		return
	}
	err := s.db.AddCategory(category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.categories = append(s.categories, models.Category{Name: sql.NullString{String: category, Valid: true}})
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}

func (s *Server) DeleteCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	categoryID := r.FormValue("categoryId")
	err := s.db.DeleteCategory(categoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for i, category := range s.categories {
		if category.CategoryId.String == categoryID {
			s.categories = append(s.categories[:i], s.categories[i+1:]...)
			break
		}
	}
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}

func (s *Server) EditCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	categoryID := sql.NullString{String: r.FormValue("categoryId"), Valid: true}
	categoryName := sql.NullString{String: r.FormValue("newCategoryName"), Valid: true}
	if !IsUniqueCategory(s.categories, categoryName.String) {
		render(w, r, "../categories", map[string]interface{}{"Categories": s.categories, "Error": "Category already exists"})
		return
	}
	err := s.db.EditCategory(categoryID.String, categoryName.String)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for i, category := range s.categories {
		if category.CategoryId == categoryID {
			s.categories[i].Name = categoryName
			break
		}
	}
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}

func IsUniqueCategory(categories []models.Category, category string) bool {
	for _, existingCategory := range categories {
		if strings.EqualFold(existingCategory.Name.String, category) {
			return false
		}
	}
	return true
}
