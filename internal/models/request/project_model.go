package request

import (
	"time"

	"github.com/google/uuid"
)

type ProjectEntry struct {
	ID          int       `json:"id,omitempty"`
	SourceToken string    `json:"source_token,omitempty"`
	ProjectName string    `json:"project_name"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewProjectEntry(project_name string) *ProjectEntry {
	return &ProjectEntry{
		SourceToken: uuid.New().String(),
		ProjectName: project_name,
		CreatedAt:   time.Now(),
	}
}
