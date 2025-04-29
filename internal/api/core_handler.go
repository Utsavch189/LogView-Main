package api

import (
	"github.com/Utsavch189/logview/internal/service"
	"github.com/gorilla/mux"
)

func CoreSystemHandler(r *mux.Router) {
	r.HandleFunc("/api/core/settings", service.UpdateCoreSettingsService).
		Methods("POST")

	r.HandleFunc("/api/core/settings", service.GetCoreSettingsService).
		Methods("GET")
}
