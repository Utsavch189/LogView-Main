package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"context"

	"github.com/Utsavch189/logview/internal/configs"
	"github.com/Utsavch189/logview/internal/controller"
	"github.com/Utsavch189/logview/internal/models/request"
	"github.com/Utsavch189/logview/internal/models/response"
	"github.com/Utsavch189/logview/internal/utils"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

var ctx = context.Background()
var rdb *redis.Client

func LogIngestService(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rdb = redis.NewClient(&redis.Options{
		Addr:     configs.GetEnv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})

	logData := &request.LogEntry{}

	err := json.NewDecoder(r.Body).Decode(&logData)
	if err != nil {
		http.Error(w, "Invalid log format", http.StatusBadRequest)
		return
	}

	_, perr := controller.GetProjectBySourceToken(logData.SourceToken)

	if perr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.ErrorResponse(perr, "log is denied!"))
		return
	}

	jsonLog, _ := json.Marshal(logData)
	// print(logData.Exception)
	rdb.RPush(ctx, "logs", jsonLog)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Log is published",
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

	sql, sqlCount, sqlInfoLogCount, sqlWarnLogCount, sqlErrorLogCount, sqlDebugLogCount := utils.GenerateSqlQueryForFilterSearch(logFilterSearch, *_project, page, pageSize)
	// print(sql)

	logs, count, infoCount, warnCount, errorCount, debugCount, err := controller.GetFilteredLogs(sql, sqlCount, sqlInfoLogCount, sqlWarnLogCount, sqlErrorLogCount, sqlDebugLogCount)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.ErrorResponse(err, "logs fetching failed!"))
		return
	}

	response := response.LogFilteredResponse{
		Logs:       logs,
		Count:      count,
		InfoCount:  infoCount,
		WarnCount:  warnCount,
		ErrorCount: errorCount,
		DebugCount: debugCount,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&response)
}
