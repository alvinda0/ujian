package models

import (
	"time"

	"gorm.io/gorm"
)

type Teacher struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"unique;not null" json:"user_id"`
	EmployeeID  string         `gorm:"unique;not null;size:20" json:"employee_id"`
	FullName    string         `gorm:"not null;size:100" json:"full_name"`
	Phone       string         `gorm:"size:20" json:"phone"`
	Address     string         `gorm:"size:255" json:"address"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Assignments []ClassSubject `gorm:"foreignKey:TeacherID" json:"assignments,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Teacher) TableName() string {
	return "teachers"
}