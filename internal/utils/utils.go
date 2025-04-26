package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Utsavch189/logview/internal/models/request"
	"github.com/xuri/excelize/v2"
)

func GenerateSqlQueryForFilterSearch(filter request.LogFilterSearch, project request.ProjectEntry, page int, pageSize int) (string, string, string, string, string, string, string) {

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
	var sqlPaginateCount string

	sqlCount = fmt.Sprintf(`SELECT COUNT(*) FROM logs WHERE source_token = '%s'`, project.SourceToken)
	sqlInfoLogCount = fmt.Sprintf(`SELECT COUNT(*) FROM logs WHERE LOWER(level) = '%s' AND source_token = '%s'`, "info", project.SourceToken)
	sqlWarnLogCount = fmt.Sprintf(`SELECT COUNT(*) FROM logs WHERE LOWER(level) = '%s' AND source_token = '%s'`, "warning", project.SourceToken)
	sqlErrorLogCount = fmt.Sprintf(`SELECT COUNT(*) FROM logs WHERE LOWER(level) = '%s' AND source_token = '%s'`, "error", project.SourceToken)
	sqlDebugLogCount = fmt.Sprintf(`SELECT COUNT(*) FROM logs WHERE LOWER(level) = '%s' AND source_token = '%s'`, "debug", project.SourceToken)
	sqlPaginateCount = fmt.Sprintf(`SELECT COUNT(*) FROM logs WHERE source_token = '%s' AND LOWER(level) IN (%s)`, project.SourceToken, levels)

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
		sqlPaginateCount += fmt.Sprintf(` AND created_at BETWEEN '%s' AND '%s'`, filter.LogDates.From.UTC().Format("2006-01-02 15:04:05"), endOfDay.Format("2006-01-02 15:04:05"))
	}

	sql += fmt.Sprintf(` ORDER BY ID DESC LIMIT %d OFFSET %d`, pageSize, offset)

	return sql, sqlCount, sqlInfoLogCount, sqlWarnLogCount, sqlErrorLogCount, sqlDebugLogCount, sqlPaginateCount
}

func GenerateSqlQueryForLogDownload(filter request.LogFilterSearch, project request.ProjectEntry) string {
	activeLevels := []string{}

	for _, level := range filter.LogLevels {
		if level.Checked {
			activeLevels = append(activeLevels, strings.ToLower(fmt.Sprintf(`'%s'`, level.Level)))
		}
	}

	levels := strings.Join(activeLevels, ",")

	sql := fmt.Sprintf(`SELECT id,
		time, level, logger, message, hostname, source_token,
		pathname, filename, func_name, lineno, thread,
		process, module, created, exception, created_at
	FROM logs Where LOWER(level) IN (%s) AND source_token = '%s'`, levels, project.SourceToken)

	if !filter.LogDates.From.IsZero() && !filter.LogDates.To.IsZero() {
		endOfDay := filter.LogDates.To.Add(24 * time.Hour).Truncate(24 * time.Hour).Add(-time.Second)
		sql += fmt.Sprintf(` AND created_at BETWEEN '%s' AND '%s'`, filter.LogDates.From.UTC().Format("2006-01-02 15:04:05"), endOfDay.Format("2006-01-02 15:04:05"))
	}

	sql += ` ORDER BY ID`

	return sql
}

func GenerateXlLogs(logs []request.LogEntry) *excelize.File {
	f := excelize.NewFile()
	sheet := f.GetSheetName(0)

	headers := []string{
		"ID", "Time", "Level", "Logger", "Message", "Hostname", "SourceToken",
		"Pathname", "Filename", "FuncName", "Lineno", "Thread", "Process",
		"Module", "Created", "Exception", "CreatedAt",
	}

	for colIdx, header := range headers {
		col, _ := excelize.ColumnNumberToName(colIdx + 1)
		cell := col + "1"
		f.SetCellValue(sheet, cell, header)
	}

	for rowIdx, log := range logs {
		row := rowIdx + 2

		createdReadable := time.Unix(
			int64(log.Created),
			int64((log.Created-float64(int64(log.Created)))*1e9),
		).Format(time.RFC3339)

		createdAtFormatted := log.CreatedAt.Format(time.RFC3339)

		values := []interface{}{
			log.ID, log.Time, log.Level, log.Logger, log.Message, log.Hostname, log.SourceToken,
			log.Pathname, log.Filename, log.FuncName, log.Lineno, log.Thread, log.Process,
			log.Module, createdReadable, log.Exception, createdAtFormatted,
		}

		for colIdx, val := range values {
			col, _ := excelize.ColumnNumberToName(colIdx + 1)
			cell := col + strconv.Itoa(row)
			f.SetCellValue(sheet, cell, val)
		}
	}

	return f
}

func SanitizeLogMessage(msg string) string {
	msg = strings.ReplaceAll(msg, `"`, `\"`)   // Escape double quotes
	msg = strings.ReplaceAll(msg, `'`, `\'`)   // Escape single quotes
	msg = strings.ReplaceAll(msg, "\n", "\\n") // Escape newlines
	msg = strings.ReplaceAll(msg, "\r", "\\r") // Escape carriage returns
	return msg
}
