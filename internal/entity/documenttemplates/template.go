package documenttemplates

import "time"

type DocumentTemplate struct {
	ID          int64      `json:"id"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Engine      string     `json:"engine"`
	OutputFormat string    `json:"output_format"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

