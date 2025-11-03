package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type QuestionType string

const (
	QuestionTypeMultipleChoice QuestionType = "multiple_choice"
	QuestionTypeEssay          QuestionType = "essay"
	QuestionTypeTrueFalse      QuestionType = "true_false"
	QuestionTypeShortAnswer    QuestionType = "short_answer"
)

type QuestionOptions []string

// Implement driver.Valuer interface for database storage
func (qo QuestionOptions) Value() (driver.Value, error) {
	return json.Marshal(qo)
}

// Implement sql.Scanner interface for database retrieval
func (qo *QuestionOptions) Scan(value interface{}) error {
	if value == nil {
		*qo = nil
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	
	return json.Unmarshal(bytes, qo)
}

type Question struct {
	ID         uint            `gorm:"primaryKey" json:"id"`
	ExamID     uint            `gorm:"not null" json:"exam_id"`
	Text       string          `gorm:"not null;type:text" json:"text"`
	Type       QuestionType    `gorm:"not null;type:enum('multiple_choice','essay','true_false','short_answer')" json:"type"`
	Options    QuestionOptions `gorm:"type:json" json:"options"`
	Answer     string          `gorm:"not null;type:text" json:"answer"`
	Points     int             `gorm:"default:1" json:"points"`
	OrderIndex int             `gorm:"not null" json:"order_index"`
	Exam       Exam            `gorm:"foreignKey:ExamID" json:"exam,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	DeletedAt  gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (Question) TableName() string {
	return "questions"
}