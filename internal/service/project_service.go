package service

import (
	"encoding/json"
	"net/http"

	"github.com/Utsavch189/logview/internal/controller"
	"github.com/Utsavch189/logview/internal/models/request"
	"github.com/Utsavch189/logview/internal/models/response"
	"github.com/gorilla/mux"
)

func CreateProjectService(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newProject request.ProjectEntry
	err := json.NewDecoder(r.Body).Decode(&newProject)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.ErrorResponse(err, "project creation failed!"))
		return
	}

	project := request.NewProjectEntry(newProject.ProjectName)

	createdProject, cerr := controller.CreateProject(project)
	if cerr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.ErrorResponse(cerr, "project creation failed!"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&createdProject)
}

func GetAllProjectSerive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	projects, err := controller.GetAllProject()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.ErrorResponse(err, "projects fetching failed!"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&projects)
}

func DeleteProjectService(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	source_token := params["source_token"]

	err := controller.DeleteProject(source_token)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.ErrorResponse(err, "project deletion failed!"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "deleted",
	})
}
