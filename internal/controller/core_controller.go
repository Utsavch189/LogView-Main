package controller

import (
	"database/sql"

	"github.com/Utsavch189/logview/internal/configs"
	"github.com/Utsavch189/logview/internal/models/request"
)

func GetCoreSystemSettings() (*request.CoreSettings, error) {
	db, err := configs.Connect()

	var settings request.CoreSettings

	if err != nil {
		return &settings, err
	}
	defer db.Close()

	query := `Select * from core_settings`
	row := db.QueryRow(query)

	errs := row.Scan(
		&settings.AutoLogDeleteDays,
	)

	if errs == sql.ErrNoRows {
		return nil, nil
	}

	if errs != nil {
		return &settings, errs
	}

	return &settings, nil
}

func UpdateCoreSettings(settings *request.CoreSettings) error {
	db, err := configs.Connect()

	if err != nil {
		return err
	}
	defer db.Close()

	_settings, errs := GetCoreSystemSettings()

	if errs != nil {
		return errs
	}

	var query string

	if _settings != nil {
		query = `
			Update core_settings set autolog_delete_days = ?
		`
	} else {
		query = `
			Insert into core_settings(autolog_delete_days) Values(?)
		`
	}
	// print(query)
	_, err1 := db.Exec(query,
		settings.AutoLogDeleteDays,
	)

	if err1 != nil {
		return err1
	}

	return nil
}
