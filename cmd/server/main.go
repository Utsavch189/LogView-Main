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

	configs.InitRedisClient(
		configs.GetEnv("REDIS_ADDR"),
		"",
		0,
	)
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
	api.CoreSystemHandler(r)

	staticDir := http.Dir(filepath.Join("internal", "static"))
	templateDir := http.Dir(filepath.Join("internal", "templates"))
	fs := http.FileServer(staticDir)
	tfs := http.FileServer(templateDir)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.PathPrefix("/templates/").Handler(http.StripPrefix("/templates/", tfs))

	http.ListenAndServe(":53423", handler)

}
