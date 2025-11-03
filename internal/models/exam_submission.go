package models

import (
	"time"

	"gorm.io/gorm"
)

type SubmissionStatus string

const (
	SubmissionStatusInProgress SubmissionStatus = "in_progress"
	SubmissionStatusSubmitted  SubmissionStatus = "submitted"
	SubmissionStatusGraded     SubmissionStatus = "graded"
	SubmissionStatusExpired    SubmissionStatus = "expired"
)

type ExamSubmission struct {
	ID          uint               `gorm:"primaryKey" json:"id"`
	ExamID      uint               `gorm:"not null" json:"exam_id"`
	StudentID   uint               `gorm:"not null" json:"student_id"`
	StartedAt   time.Time          `gorm:"not null" json:"started_at"`
	SubmittedAt *time.Time         `json:"submitted_at"`
	Score       *float64           `json:"score"`
	MaxScore    *float64           `json:"max_score"`
	Status      SubmissionStatus   `gorm:"default:'in_progress';type:enum('in_progress','submitted','graded','expired')" json:"status"`
	AttemptNum  int                `gorm:"default:1" json:"attempt_num"`
	Answers     []SubmissionAnswer `gorm:"foreignKey:SubmissionID" json:"answers,omitempty"`
	Exam        Exam               `gorm:"foreignKey:ExamID" json:"exam,omitempty"`
	Student     Student            `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	DeletedAt   gorm.DeletedAt     `gorm:"index" json:"-"`
}

func (ExamSubmission) TableName() string {
	return "exam_submissions"
}