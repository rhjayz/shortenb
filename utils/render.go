package utils

import (
	"html/template"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, filename string, data any) {
	tmpl := template.Must(template.ParseFiles(
		"views/layout.html",
		"views/templates/templates.html",
		filename,
	))

	err := tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
