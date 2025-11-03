package models

import (
	"time"

	"gorm.io/gorm"
)

type ClassSubject struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ClassID   uint           `gorm:"not null" json:"class_id"`
	SubjectID uint           `gorm:"not null" json:"subject_id"`
	TeacherID uint           `gorm:"not null" json:"teacher_id"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	Class     Class          `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	Subject   Subject        `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
	Teacher   Teacher        `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ClassSubject) TableName() string {
	return "class_subjects"
}

// Unique constraint for class-subject-teacher combination
func (cs *ClassSubject) BeforeCreate(tx *gorm.DB) error {
	// Add unique index constraint
	return nil
}