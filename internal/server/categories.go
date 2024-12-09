package server

import (
	"forum-go/internal/models"
	"net/http"
	"strings"
)

func (s *Server) GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	if !s.isLoggedIn(r) || !IsAdmin(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	categories, err := s.db.GetCategories()
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render(w, r, "admin/categories", map[string]interface{}{"Categories": categories})
}

func (s *Server) PostCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	if !s.isLoggedIn(r) || !IsAdmin(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	category := r.FormValue("categoryName")
	if !IsUniqueCategory(s.categories, category) {
		render(w, r, "admin/categories", map[string]interface{}{"Categories": s.categories, "Error": "Category already exists"})
		return
	}
	err := s.db.AddCategory(category)
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	s.categories = append(s.categories, models.Category{Name: category})
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}

func (s *Server) DeleteCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	if !s.isLoggedIn(r) || !IsAdmin(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	categoryID := r.FormValue("categoryId")
	err := s.db.DeleteCategory(categoryID)
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	for i, category := range s.categories {
		if category.CategoryId == categoryID {
			s.categories = append(s.categories[:i], s.categories[i+1:]...)
			break
		}
	}
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}

func (s *Server) EditCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	if !s.isLoggedIn(r) || !IsAdmin(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	categoryID := r.FormValue("categoryId")
	categoryName := r.FormValue("newCategoryName")
	if !IsUniqueCategory(s.categories, categoryName) {
		render(w, r, "admin/categories", map[string]interface{}{"Categories": s.categories, "Error": "Category already exists"})
		return
	}
	err := s.db.EditCategory(categoryID, categoryName)
	if err != nil {
		s.errorHandler(w, r, http.StatusInternalServerError, err.Error())
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
		if strings.EqualFold(existingCategory.Name, category) {
			return false
		}
	}
	return true
}
