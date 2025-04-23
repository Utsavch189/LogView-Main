package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/Utsavch189/logview/internal/models/request"
)

func GenerateSqlQueryForFilterSearch(filter request.LogFilterSearch, project request.ProjectEntry, page int, pageSize int) (string, string, string, string, string, string) {

	activeLevels := []string{}

	offset := (page - 1) * pageSize

	for _, level := range filter.LogLevels {
		if level.Checked {
			activeLevels = append(activeLevels, strings.ToLower(fmt.Sprintf(`'%s'`, level.Level)))
		}
	}

	levels := strings.Join(activeLevels, ",")

	var sqlCount string
	var sqlInfoLogCount string
	var sqlWarnLogCount string
	var sqlErrorLogCount string
	var sqlDebugLogCount string

	sqlCount = fmt.Sprintf(`SELECT COUNT(*) FROM logs WHERE source_token = '%s'`, project.SourceToken)
	sqlInfoLogCount = fmt.Sprintf(`SELECT COUNT(*) FROM logs WHERE LOWER(level) = '%s' AND source_token = '%s'`, "info", project.SourceToken)
	sqlWarnLogCount = fmt.Sprintf(`SELECT COUNT(*) FROM logs WHERE LOWER(level) = '%s' AND source_token = '%s'`, "warning", project.SourceToken)
	sqlErrorLogCount = fmt.Sprintf(`SELECT COUNT(*) FROM logs WHERE LOWER(level) = '%s' AND source_token = '%s'`, "error", project.SourceToken)
	sqlDebugLogCount = fmt.Sprintf(`SELECT COUNT(*) FROM logs WHERE LOWER(level) = '%s' AND source_token = '%s'`, "debug", project.SourceToken)

	sql := fmt.Sprintf(`SELECT id,
		time, level, logger, message, hostname, source_token,
		pathname, filename, func_name, lineno, thread,
		process, module, created, exception, created_at
	FROM logs Where LOWER(level) IN (%s) AND source_token = '%s'`, levels, project.SourceToken)

	if !filter.LogDates.From.IsZero() && !filter.LogDates.To.IsZero() {
		endOfDay := filter.LogDates.To.Add(24 * time.Hour).Truncate(24 * time.Hour).Add(-time.Second)
		sql += fmt.Sprintf(` AND created_at BETWEEN '%s' AND '%s'`, filter.LogDates.From.UTC().Format("2006-01-02 15:04:05"), endOfDay.Format("2006-01-02 15:04:05"))
		sqlCount += fmt.Sprintf(` AND created_at BETWEEN '%s' AND '%s'`, filter.LogDates.From.UTC().Format("2006-01-02 15:04:05"), endOfDay.Format("2006-01-02 15:04:05"))
		sqlInfoLogCount += fmt.Sprintf(` AND created_at BETWEEN '%s' AND '%s'`, filter.LogDates.From.UTC().Format("2006-01-02 15:04:05"), endOfDay.Format("2006-01-02 15:04:05"))
		sqlWarnLogCount += fmt.Sprintf(` AND created_at BETWEEN '%s' AND '%s'`, filter.LogDates.From.UTC().Format("2006-01-02 15:04:05"), endOfDay.Format("2006-01-02 15:04:05"))
		sqlErrorLogCount += fmt.Sprintf(` AND created_at BETWEEN '%s' AND '%s'`, filter.LogDates.From.UTC().Format("2006-01-02 15:04:05"), endOfDay.Format("2006-01-02 15:04:05"))
		sqlDebugLogCount += fmt.Sprintf(` AND created_at BETWEEN '%s' AND '%s'`, filter.LogDates.From.UTC().Format("2006-01-02 15:04:05"), endOfDay.Format("2006-01-02 15:04:05"))
	}

	sql += fmt.Sprintf(` ORDER BY ID DESC LIMIT %d OFFSET %d`, pageSize, offset)

	return sql, sqlCount, sqlInfoLogCount, sqlWarnLogCount, sqlErrorLogCount, sqlDebugLogCount
}
