package controller

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Utsavch189/logview/internal/configs"
	"github.com/Utsavch189/logview/internal/models/request"
)

func CreateProject(project *request.ProjectEntry) (*request.ProjectEntry, error) {
	db, err := configs.Connect()

	if err != nil {
		return &request.ProjectEntry{}, err
	}
	defer db.Close()

	query := `Insert into projects(source_token,project_name,created_at) Values(?,?,?)`

	_, cerr := db.Exec(query, project.SourceToken, project.ProjectName, project.CreatedAt)

	if cerr != nil {
		return &request.ProjectEntry{}, cerr
	}

	return project, nil
}

func GetProjectBySourceToken(source_token string) (*request.ProjectEntry, error) {
	db, err := configs.Connect()

	var project request.ProjectEntry

	if err != nil {
		return &project, err
	}
	defer db.Close()

	query := `Select * from projects Where source_token = ?`
	row := db.QueryRow(query, source_token)

	errs := row.Scan(
		&project.ID,
		&project.SourceToken,
		&project.ProjectName,
		&project.CreatedAt,
	)

	if errs == sql.ErrNoRows {
		return &project, fmt.Errorf("project with source token %s not found", source_token)
	}

	if errs != nil {
		return &project, errs
	}

	return &project, nil
}

func GetProjectByName(project_name string) (*request.ProjectEntry, error) {
	db, err := configs.Connect()

	var project request.ProjectEntry

	if err != nil {
		return &project, err
	}
	defer db.Close()

	query := `Select * from projects Where project_name = ?`
	row := db.QueryRow(query, project_name)

	errs := row.Scan(
		&project.ID,
		&project.SourceToken,
		&project.ProjectName,
		&project.CreatedAt,
	)

	if errs == sql.ErrNoRows {
		return &project, fmt.Errorf("project with name %s not found", project_name)
	}

	if errs != nil {
		return &project, errs
	}

	return &project, nil
}

func GetAllProject() ([]request.ProjectEntry, error) {
	db, err := configs.Connect()

	var projects []request.ProjectEntry

	if err != nil {
		return projects, err
	}
	defer db.Close()

	query := `Select * from projects order by created_at;`
	rows, rerr := db.Query(query)

	if rerr != nil {
		return projects, err
	}

	for rows.Next() {
		var project request.ProjectEntry
		var createdAt string

		errs := rows.Scan(
			&project.ID,
			&project.SourceToken,
			&project.ProjectName,
			&project.CreatedAt,
		)

		if errs != nil {
			return projects, errs
		}

		project.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)

		projects = append(projects, project)
	}

	return projects, nil
}
