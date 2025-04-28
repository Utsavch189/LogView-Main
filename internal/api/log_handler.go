package api

import (
	"github.com/Utsavch189/logview/internal/service"
	"github.com/gorilla/mux"
)

func LogIngestHandler(r *mux.Router) {
	r.HandleFunc("/api/logs/ingest", service.LogIngestService).
		Methods("POST")

	r.HandleFunc("/api/logs/{project}/get-all", service.GetAllLogs).
		Methods("GET")

	r.HandleFunc("/api/logs/{project}/apply-filters/get-all", service.GetAllLogsWithFilters).
		Methods("POST")

	r.HandleFunc("/api/logs/{project}/download-logs", service.DownloadLogs).
		Methods("POST")

	r.HandleFunc("/api/logs/{project}/delete-logs", service.DeleteLogsService).
		Methods("DELETE")
}
