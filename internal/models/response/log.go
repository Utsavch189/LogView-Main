package response

import "github.com/Utsavch189/logview/internal/models/request"

type LogFilteredResponse struct {
	Logs          []request.LogEntry `json:"logs"`
	Count         int                `json:"count"`
	InfoCount     int                `json:"info_count"`
	WarnCount     int                `json:"warn_count"`
	ErrorCount    int                `json:"error_count"`
	DebugCount    int                `json:"debug_count"`
	PaginateCount int                `json:"paginate_count"`
}
