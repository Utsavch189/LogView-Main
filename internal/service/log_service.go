package service

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Utsavch189/logview/internal/configs"
	"github.com/Utsavch189/logview/internal/controller"
	"github.com/Utsavch189/logview/internal/models/request"
	"github.com/Utsavch189/logview/internal/models/response"
	"github.com/Utsavch189/logview/internal/utils"
	"github.com/gorilla/mux"
)

func LogIngestService(w http.ResponseWriter, r *http.Request) {

	var logDatas []request.LogEntry

	err := json.NewDecoder(r.Body).Decode(&logDatas)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.ErrorResponse(err, "Invalid log format!"))
		return
	}

	// fmt.Printf("logDatas[0]: %+v\n", logDatas[0])

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	_, perr := controller.GetProjectBySourceToken(token)
	if perr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response.ErrorResponse(perr, "Unauthorized"))
		return
	}

	for i := range logDatas {
		logData := &logDatas[i]
		jsonLog, err := json.Marshal(logData)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Failed to marshal log entry",
			})
			return
		}

		err = configs.Rdb.RPush(configs.Ctx, "logs", jsonLog).Err()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Failed to store log entry in Redis",
			})
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logs have been successfully ingested",
	})
}

func GetAllLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	project := params["project"]

	_project, perr := controller.GetProjectByName(project)

	if perr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.ErrorResponse(perr, "logs fetching failed!"))
		return
	}

	logs, err := controller.GetAllLogs(_project.SourceToken)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.ErrorResponse(err, "logs fetching failed!"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&logs)
}

func GetAllLogsWithFilters(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var logFilterSearch request.LogFilterSearch

	err := json.NewDecoder(r.Body).Decode(&logFilterSearch)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.ErrorResponse(err, "wrong payload!"))
		return
	}

	params := mux.Vars(r)
	project := params["project"]

	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	page, err := strconv.Atoi(pageStr)
	if err != nil || pageStr == "" {
		page = 0
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSizeStr == "" {
		pageSize = 10
	}

	_project, perr := controller.GetProjectByName(project)

	if perr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.ErrorResponse(perr, "logs fetching failed!"))
		return
	}

	sql, sqlCount, sqlInfoLogCount, sqlWarnLogCount, sqlErrorLogCount, sqlDebugLogCount, sqlPaginateCount := utils.GenerateSqlQueryForFilterSearch(logFilterSearch, *_project, page, pageSize)
	// print(sql)

	logs, count, infoCount, warnCount, errorCount, debugCount, paginateCount, err := controller.GetFilteredLogs(sql, sqlCount, sqlInfoLogCount, sqlWarnLogCount, sqlErrorLogCount, sqlDebugLogCount, sqlPaginateCount)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.ErrorResponse(err, "logs fetching failed!"))
		return
	}

	response := response.LogFilteredResponse{
		Logs:          logs,
		Count:         count,
		InfoCount:     infoCount,
		WarnCount:     warnCount,
		ErrorCount:    errorCount,
		DebugCount:    debugCount,
		PaginateCount: paginateCount,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&response)
}

func DownloadLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var logFilterSearch request.LogFilterSearch

	err := json.NewDecoder(r.Body).Decode(&logFilterSearch)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.ErrorResponse(err, "wrong payload!"))
		return
	}

	params := mux.Vars(r)
	project := params["project"]

	_project, perr := controller.GetProjectByName(project)

	if perr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.ErrorResponse(perr, "logs fetching failed!"))
		return
	}

	sql := utils.GenerateSqlQueryForLogDownload(logFilterSearch, *_project)

	logs, lerr := controller.GetLogsForDownload(sql)

	if lerr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.ErrorResponse(lerr, "logs fetching failed!"))
		return
	}

	f := utils.GenerateXlLogs(logs)

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=logs.xlsx")
	w.Header().Set("File-Name", "logs.xlsx")
	w.Header().Set("Access-Control-Expose-Headers", "File-Name")

	if err := f.Write(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
