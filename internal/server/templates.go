package server

import (
	"html/template"
	"net/http"
)

func render(w http.ResponseWriter, page string, data map[string]interface{}) {
	t, err := template.ParseFiles("./assets/" + page + ".tmpl.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}
