package server

import (
	"forum-go/internal/models"
	"html/template"
	"net/http"
)

func render(w http.ResponseWriter, r *http.Request, page string, data map[string]interface{}) {
	// render renders the template with the given data
	t, err := template.ParseFiles("./assets/templates/" + page + ".tmpl.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if data == nil {
		data = map[string]interface{}{}
	}
	user, ok := r.Context().Value(contextKeyUser).(models.User)
	if ok {
		data["User"] = user
	}
	t.Execute(w, data)
}
