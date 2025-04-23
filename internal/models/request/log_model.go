package request

import "time"

type LogEntry struct {
	ID          int    `json:"id,omitempty"`
	Time        string `json:"time"`
	Level       string `json:"level"`
	Logger      string `json:"logger"`
	Message     string `json:"message"`
	Hostname    string `json:"hostname"`
	SourceToken string `json:"source_token"`

	Pathname string  `json:"pathname"`
	Filename string  `json:"filename"`
	FuncName string  `json:"func_name"`
	Lineno   int     `json:"lineno"`
	Thread   string  `json:"thread"`
	Process  string  `json:"process"`
	Module   string  `json:"module"`
	Created  float64 `json:"created"` // float64 to match SQLite REAL

	Exception string    `json:"exception,omitempty"`  // optional stack trace
	CreatedAt time.Time `json:"created_at,omitempty"` // mapped to SQLite TIMESTAMP
}

type LogLevel struct {
	Level   string `json:"level"`
	ID      string `json:"id"`
	Checked bool   `json:"checked"`
}

type LogDate struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type LogFilterSearch struct {
	LogLevels []LogLevel `json:"loglevels"`
	LogDates  LogDate    `json:"dates,omitempty"`
}
