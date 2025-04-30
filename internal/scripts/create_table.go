package scripts

import "database/sql"

func CreateTables(db *sql.DB) error {
	statements := []string{

		`CREATE TABLE IF NOT EXISTS logs (
    	id INT AUTO_INCREMENT PRIMARY KEY,
    	time TEXT,
    	level TEXT,
    	logger TEXT,
    	message TEXT,
    	hostname TEXT,
    	source_token TEXT,

    	pathname TEXT,
    	filename TEXT,
    	func_name TEXT,
    	lineno INT,
    	thread TEXT,
    	process TEXT,
    	module TEXT,
    	created DOUBLE,

    	exception TEXT, -- optional, stack trace if exists

    	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`,

		`CREATE TABLE IF NOT EXISTS projects (
	    	id INT AUTO_INCREMENT PRIMARY KEY,
	    	source_token VARCHAR(50) NOT NULL UNIQUE,
	    	project_name VARCHAR(50) NOT NULL UNIQUE,
	    	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS core_settings (
	    	autolog_delete_days INT DEFAULT 60
		);`,
	}

	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}
