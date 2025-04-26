package api

import (
	"github.com/Utsavch189/logview/internal/service"
	"github.com/gorilla/mux"
)

func ProjectHandler(r *mux.Router) {
	r.HandleFunc("/api/project/create", service.CreateProjectService).
		Methods("POST")

	r.HandleFunc("/api/project/get-all", service.GetAllProjectSerive).
		Methods("GET")

	r.HandleFunc("/api/project/delete/{source_token}", service.DeleteProjectService).
		Methods("DELETE")
}
