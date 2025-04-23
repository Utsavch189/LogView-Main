package service

import (
	"net/http"
	"text/template"
)

func HomeService(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("internal/templates/index.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := struct {
		Name string
	}{
		Name: "Gopher",
	}

	tmpl.Execute(w, data)
}
