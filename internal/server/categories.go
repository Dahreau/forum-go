package server

import "net/http"

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
	err := s.db.AddCategory(category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}
