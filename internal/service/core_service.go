package service

import (
	"encoding/json"
	"net/http"

	"github.com/Utsavch189/logview/internal/controller"
	"github.com/Utsavch189/logview/internal/models/request"
	"github.com/Utsavch189/logview/internal/models/response"
)

func GetCoreSettingsService(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	settings, err := controller.GetCoreSystemSettings()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.ErrorResponse(err, "settings fetching failed!"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(settings)
}

func UpdateCoreSettingsService(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var settings request.CoreSettings

	err := json.NewDecoder(r.Body).Decode(&settings)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.ErrorResponse(err, "wrong payload!"))
		return
	}

	err = controller.UpdateCoreSettings(&settings)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.ErrorResponse(err, "settings update failed!"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(settings)

}
