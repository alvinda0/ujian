package models

import (
	"time"

	"gorm.io/gorm"
)

type SubmissionAnswer struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	SubmissionID uint           `gorm:"not null" json:"submission_id"`
	QuestionID   uint           `gorm:"not null" json:"question_id"`
	Answer       string         `gorm:"type:text" json:"answer"`
	IsCorrect    *bool          `json:"is_correct"`
	Points       *float64       `json:"points"`
	Submission   ExamSubmission `gorm:"foreignKey:SubmissionID" json:"submission,omitempty"`
	Question     Question       `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (SubmissionAnswer) TableName() string {
	return "submission_answers"
}