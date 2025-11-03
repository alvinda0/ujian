package models

import (
	"time"

	"gorm.io/gorm"
)

type Class struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null;size:50" json:"name"` // e.g., "X IPA 1"
	Grade       int            `gorm:"not null" json:"grade"`        // 10, 11, 12
	Stream      string         `gorm:"not null;size:20" json:"stream"` // IPA, IPS, etc.
	Section     string         `gorm:"not null;size:10" json:"section"` // 1, 2, 3, etc.
	Description string         `gorm:"size:255" json:"description"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	Students    []Student      `gorm:"foreignKey:ClassID" json:"students,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Class) TableName() string {
	return "classes"
}