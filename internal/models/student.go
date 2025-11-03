package models

import (
	"time"

	"gorm.io/gorm"
)

type Student struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"unique;not null" json:"user_id"`
	StudentID string         `gorm:"unique;not null;size:20" json:"student_id"` // NIS
	FullName  string         `gorm:"not null;size:100" json:"full_name"`
	ClassID   uint           `gorm:"not null" json:"class_id"`
	Phone     string         `gorm:"size:20" json:"phone"`
	Address   string         `gorm:"size:255" json:"address"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Class     Class          `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Student) TableName() string {
	return "students"
}