package controller

import (
	"time"

	"github.com/Utsavch189/logview/internal/configs"
	"github.com/Utsavch189/logview/internal/models/request"
)

func SaveLogToDB(log *request.LogEntry) error {
	db, err := configs.Connect()

	if err != nil {
		return err
	}
	defer db.Close()

	query := `
	INSERT INTO logs (
		time, level, logger, message, hostname, source_token,
		pathname, filename, func_name, lineno, thread,
		process, module, created, exception
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err1 := db.Exec(query,
		log.Time,
		log.Level,
		log.Logger,
		log.Message,
		log.Hostname,
		log.SourceToken,
		log.Pathname,
		log.Filename,
		log.FuncName,
		log.Lineno,
		log.Thread,
		log.Process,
		log.Module,
		log.Created,
		log.Exception,
	)

	if err1 != nil {
		return err1
	}

	return nil
}

func GetAllLogs(source_token string) ([]request.LogEntry, error) {

	db, err := configs.Connect()

	var logs []request.LogEntry

	if err != nil {
		return logs, err
	}

	query := `SELECT id,
		time, level, logger, message, hostname, source_token,
		pathname, filename, func_name, lineno, thread,
		process, module, created, exception, created_at
	FROM logs Where source_token = ? ORDER BY id DESC`

	rows, err := db.Query(query, source_token)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var log request.LogEntry
		var createdAt string

		err := rows.Scan(
			&log.ID,
			&log.Time,
			&log.Level,
			&log.Logger,
			&log.Message,
			&log.Hostname,
			&log.SourceToken,
			&log.Pathname,
			&log.Filename,
			&log.FuncName,
			&log.Lineno,
			&log.Thread,
			&log.Process,
			&log.Module,
			&log.Created,
			&log.Exception,
			&createdAt,
		)

		if err != nil {
			return nil, err
		}

		log.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)

		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

func GetFilteredLogs(sql string, sqlCount string, sqlInfoLogCount string, sqlWarnLogCount string, sqlErrorLogCount string, sqlDebugLogCount string) ([]request.LogEntry, int, int, int, int, int, error) {

	db, err := configs.Connect()

	var logs []request.LogEntry

	if err != nil {
		return logs, 0, 0, 0, 0, 0, err
	}

	rows, err := db.Query(sql)
	if err != nil {
		return nil, 0, 0, 0, 0, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var log request.LogEntry
		var createdAt string

		err := rows.Scan(
			&log.ID,
			&log.Time,
			&log.Level,
			&log.Logger,
			&log.Message,
			&log.Hostname,
			&log.SourceToken,
			&log.Pathname,
			&log.Filename,
			&log.FuncName,
			&log.Lineno,
			&log.Thread,
			&log.Process,
			&log.Module,
			&log.Created,
			&log.Exception,
			&createdAt,
		)

		if err != nil {
			return nil, 0, 0, 0, 0, 0, err
		}

		log.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)

		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, 0, 0, 0, 0, err
	}

	var totalCount int
	err = db.QueryRow(sqlCount).Scan(&totalCount)
	if err != nil {
		return nil, 0, 0, 0, 0, 0, err
	}

	var totalInfoCount int
	err = db.QueryRow(sqlInfoLogCount).Scan(&totalInfoCount)
	if err != nil {
		return nil, 0, 0, 0, 0, 0, err
	}

	var totalWarningCount int
	err = db.QueryRow(sqlWarnLogCount).Scan(&totalWarningCount)
	if err != nil {
		return nil, 0, 0, 0, 0, 0, err
	}

	var totalErrorCount int
	err = db.QueryRow(sqlErrorLogCount).Scan(&totalErrorCount)
	if err != nil {
		return nil, 0, 0, 0, 0, 0, err
	}

	var totalDebugCount int
	err = db.QueryRow(sqlDebugLogCount).Scan(&totalDebugCount)
	if err != nil {
		return nil, 0, 0, 0, 0, 0, err
	}

	return logs, totalCount, totalInfoCount, totalWarningCount, totalErrorCount, totalDebugCount, nil
}

func GetLogsForDownload(sql string) ([]request.LogEntry, error) {

	db, err := configs.Connect()

	var logs []request.LogEntry

	if err != nil {
		return logs, err
	}

	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var log request.LogEntry
		var createdAt string

		err := rows.Scan(
			&log.ID,
			&log.Time,
			&log.Level,
			&log.Logger,
			&log.Message,
			&log.Hostname,
			&log.SourceToken,
			&log.Pathname,
			&log.Filename,
			&log.FuncName,
			&log.Lineno,
			&log.Thread,
			&log.Process,
			&log.Module,
			&log.Created,
			&log.Exception,
			&createdAt,
		)

		if err != nil {
			return nil, err
		}

		log.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)

		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}
