package models

import (
	"time"
)

// Goal represents a project goal
type Goal struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"not null" json:"title"`
	Priority  int       `gorm:"default:100" json:"priority"`
	Status    string    `gorm:"check:status IN ('active','paused','done');default:active" json:"status"`
	Notes     string    `json:"notes"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// ADR represents an Architecture Decision Record
type ADR struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"not null" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// CIRun represents a CI test run
type CIRun struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	Scope      string     `json:"scope"`
	Status     string     `gorm:"check:status IN ('pass','fail','error');not null" json:"status"`
	StartedAt  time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
}

// MarkdownTemplate represents a markdown template
type MarkdownTemplate struct {
	ID          string             `gorm:"primaryKey" json:"id"`
	Name        string             `gorm:"not null" json:"name"`
	Description string             `json:"description"`
	Category    string             `json:"category"`
	Content     string             `gorm:"not null" json:"content"`
	CreatedAt   time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
	Variables   []TemplateVariable `gorm:"foreignKey:TemplateID;constraint:OnDelete:CASCADE" json:"variables"`
}

// TemplateVariable represents a variable in a markdown template
type TemplateVariable struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	TemplateID   string `gorm:"not null" json:"template_id"`
	Name         string `gorm:"not null" json:"name"`
	Type         string `gorm:"check:type IN ('string','date','list','number','boolean');default:string" json:"type"`
	Required     bool   `gorm:"default:false" json:"required"`
	DefaultValue string `json:"default_value"`
	Description  string `json:"description"`
}

// PreferredTool represents a preferred tool for specific use cases
type PreferredTool struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	Category    string    `gorm:"not null" json:"category"`
	Description string    `gorm:"default:''" json:"description"`
	Language    string    `gorm:"default:''" json:"language"`
	UseCase     string    `gorm:"default:''" json:"use_case"`
	Priority    int       `gorm:"default:0" json:"priority"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// CursorRule represents a Cursor IDE rule file
type CursorRule struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	Category    string    `gorm:"not null" json:"category"`
	Description string    `gorm:"default:''" json:"description"`
	Content     string    `gorm:"not null" json:"content"`
	Tags        string    `gorm:"default:''" json:"tags"`        // comma-separated tags
	Source      string    `gorm:"default:'local'" json:"source"` // local, community, custom
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// ChangelogEntry represents a changelog entry
type ChangelogEntry struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Summary   string    `gorm:"not null" json:"summary"`
	Files     string    `json:"files"` // comma-separated list of files
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
