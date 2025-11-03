package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleSystemAdmin UserRole = "system_admin"
	RoleTeacher     UserRole = "teacher"
	RoleStudent     UserRole = "student"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"unique;not null;size:50" json:"username"`
	Email     string         `gorm:"unique;not null;size:100" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	Role      UserRole       `gorm:"not null;type:enum('system_admin','teacher','student')" json:"role"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}