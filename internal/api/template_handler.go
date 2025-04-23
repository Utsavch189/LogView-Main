package api

import (
	"github.com/Utsavch189/logview/internal/service"
	"github.com/gorilla/mux"
)

func TemplateHandler(r *mux.Router) {
	r.HandleFunc("/", service.HomeService).
		Methods("GET")
}
