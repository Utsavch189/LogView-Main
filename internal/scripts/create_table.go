package scripts

import "database/sql"

func CreateTables(db *sql.DB) error {
	schema := `

    CREATE TABLE IF NOT EXISTS logs (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	time TEXT,
    	level TEXT,
    	logger TEXT,
    	message TEXT,
    	hostname TEXT,
    	source_token TEXT,

    	pathname TEXT,
    	filename TEXT,
    	func_name TEXT,
    	lineno INTEGER,
    	thread TEXT,
    	process TEXT,
    	module TEXT,
    	created REAL,

    	exception TEXT, -- optional, stack trace if exists

    	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS projects (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	source_token TEXT NOT NULL UNIQUE,
		project_name TEXT NOT NULL UNIQUE,
    	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS core_settings (
    	autolog_delete_days INTEGER DEFAULT 60
	);

    `

	_, err := db.Exec(schema)
	return err
}
