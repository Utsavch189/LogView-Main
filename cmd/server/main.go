package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/Utsavch189/logview/internal/api"
	"github.com/Utsavch189/logview/internal/configs"
	"github.com/Utsavch189/logview/internal/scripts"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	db, err := configs.Connect()
	if err != nil {
		log.Fatalf("DB connect issue: %v", err)
	}
	defer db.Close()

	err = scripts.CreateTables(db)
	if err != nil {
		log.Fatalf("Table creation issue: %v", err)
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Accept"},
		AllowCredentials: true,
	})

	r := mux.NewRouter()
	handler := c.Handler(r)

	api.LogIngestHandler(r)
	api.ProjectHandler(r)
	api.TemplateHandler(r)

	staticDir := http.Dir(filepath.Join("internal", "static"))
	fs := http.FileServer(staticDir)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	http.ListenAndServe(":53423", handler)

}
