package models

import (
	"time"

	"gorm.io/gorm"
)

type ExamStatus string

const (
	ExamStatusDraft     ExamStatus = "draft"
	ExamStatusPublished ExamStatus = "published"
	ExamStatusActive    ExamStatus = "active"
	ExamStatusClosed    ExamStatus = "closed"
)

type Exam struct {
	ID          uint             `gorm:"primaryKey" json:"id"`
	Title       string           `gorm:"not null;size:200" json:"title"`
	Description string           `gorm:"size:500" json:"description"`
	SubjectID   uint             `gorm:"not null" json:"subject_id"`
	ClassID     uint             `gorm:"not null" json:"class_id"`
	TeacherID   uint             `gorm:"not null" json:"teacher_id"`
	StartTime   time.Time        `gorm:"not null" json:"start_time"`
	EndTime     time.Time        `gorm:"not null" json:"end_time"`
	Duration    int              `gorm:"not null" json:"duration"` // minutes
	MaxAttempts int              `gorm:"default:1" json:"max_attempts"`
	Status      ExamStatus       `gorm:"default:'draft';type:enum('draft','published','active','closed')" json:"status"`
	Questions   []Question       `gorm:"foreignKey:ExamID" json:"questions,omitempty"`
	Submissions []ExamSubmission `gorm:"foreignKey:ExamID" json:"submissions,omitempty"`
	Subject     Subject          `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
	Class       Class            `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	Teacher     Teacher          `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	DeletedAt   gorm.DeletedAt   `gorm:"index" json:"-"`
}

func (Exam) TableName() string {
	return "exams"
}