package services

import (
	"errors"
	"school-exam-system/internal/models"

	"gorm.io/gorm"
)

type QuestionService struct {
	db *gorm.DB
}

type AddQuestionRequest struct {
	Text       string                  `json:"text" binding:"required"`
	Type       models.QuestionType     `json:"type" binding:"required"`
	Options    models.QuestionOptions  `json:"options"`
	Answer     string                  `json:"answer" binding:"required"`
	Points     int                     `json:"points"`
	OrderIndex int                     `json:"order_index"`
}

type UpdateQuestionRequest struct {
	Text       string                  `json:"text"`
	Type       models.QuestionType     `json:"type"`
	Options    models.QuestionOptions  `json:"options"`
	Answer     string                  `json:"answer"`
	Points     int                     `json:"points"`
	OrderIndex int                     `json:"order_index"`
}

func NewQuestionService(db *gorm.DB) *QuestionService {
	return &QuestionService{db: db}
}

func (s *QuestionService) AddQuestion(examID uint, teacherID uint, req AddQuestionRequest) (*models.Question, error) {
	// Verify teacher owns the exam
	var exam models.Exam
	if err := s.db.First(&exam, examID).Error; err != nil {
		return nil, err
	}

	if exam.TeacherID != teacherID {
		return nil, errors.New("unauthorized to add questions to this exam")
	}

	if exam.Status != models.ExamStatusDraft {
		return nil, errors.New("can only add questions to draft exams")
	}

	// Validate question type and options
	if err := s.validateQuestion(req.Type, req.Options, req.Answer); err != nil {
		return nil, err
	}

	// Set default values
	if req.Points == 0 {
		req.Points = 1
	}

	// Auto-set order index if not provided
	if req.OrderIndex == 0 {
		var maxOrder int
		s.db.Model(&models.Question{}).Where("exam_id = ?", examID).
			Select("COALESCE(MAX(order_index), 0)").Scan(&maxOrder)
		req.OrderIndex = maxOrder + 1
	}

	question := models.Question{
		ExamID:     examID,
		Text:       req.Text,
		Type:       req.Type,
		Options:    req.Options,
		Answer:     req.Answer,
		Points:     req.Points,
		OrderIndex: req.OrderIndex,
	}

	if err := s.db.Create(&question).Error; err != nil {
		return nil, err
	}

	return &question, nil
}

func (s *QuestionService) UpdateQuestion(id uint, teacherID uint, req UpdateQuestionRequest) (*models.Question, error) {
	var question models.Question
	if err := s.db.Preload("Exam").First(&question, id).Error; err != nil {
		return nil, err
	}

	// Verify teacher owns the exam
	if question.Exam.TeacherID != teacherID {
		return nil, errors.New("unauthorized to update this question")
	}

	if question.Exam.Status != models.ExamStatusDraft {
		return nil, errors.New("can only update questions in draft exams")
	}

	// Update fields
	if req.Text != "" {
		question.Text = req.Text
	}
	if req.Type != "" {
		question.Type = req.Type
	}
	if req.Options != nil {
		question.Options = req.Options
	}
	if req.Answer != "" {
		question.Answer = req.Answer
	}
	if req.Points > 0 {
		question.Points = req.Points
	}
	if req.OrderIndex > 0 {
		question.OrderIndex = req.OrderIndex
	}

	// Validate updated question
	if err := s.validateQuestion(question.Type, question.Options, question.Answer); err != nil {
		return nil, err
	}

	if err := s.db.Save(&question).Error; err != nil {
		return nil, err
	}

	return &question, nil
}

func (s *QuestionService) DeleteQuestion(id uint, teacherID uint) error {
	var question models.Question
	if err := s.db.Preload("Exam").First(&question, id).Error; err != nil {
		return err
	}

	// Verify teacher owns the exam
	if question.Exam.TeacherID != teacherID {
		return errors.New("unauthorized to delete this question")
	}

	if question.Exam.Status != models.ExamStatusDraft {
		return errors.New("can only delete questions from draft exams")
	}

	return s.db.Delete(&question).Error
}

func (s *QuestionService) GetQuestionsByExam(examID uint) ([]models.Question, error) {
	var questions []models.Question
	if err := s.db.Where("exam_id = ?", examID).
		Order("order_index ASC").Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

func (s *QuestionService) GetQuestionByID(id uint) (*models.Question, error) {
	var question models.Question
	if err := s.db.First(&question, id).Error; err != nil {
		return nil, err
	}
	return &question, nil
}

func (s *QuestionService) ReorderQuestions(examID uint, teacherID uint, questionIDs []uint) error {
	// Verify teacher owns the exam
	var exam models.Exam
	if err := s.db.First(&exam, examID).Error; err != nil {
		return err
	}

	if exam.TeacherID != teacherID {
		return errors.New("unauthorized to reorder questions in this exam")
	}

	if exam.Status != models.ExamStatusDraft {
		return errors.New("can only reorder questions in draft exams")
	}

	// Update order indexes
	tx := s.db.Begin()
	for i, questionID := range questionIDs {
		if err := tx.Model(&models.Question{}).
			Where("id = ? AND exam_id = ?", questionID, examID).
			Update("order_index", i+1).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()

	return nil
}

func (s *QuestionService) validateQuestion(qType models.QuestionType, options models.QuestionOptions, answer string) error {
	switch qType {
	case models.QuestionTypeMultipleChoice:
		if len(options) < 2 {
			return errors.New("multiple choice questions must have at least 2 options")
		}
		// Check if answer is one of the options
		found := false
		for _, option := range options {
			if option == answer {
				found = true
				break
			}
		}
		if !found {
			return errors.New("answer must be one of the provided options")
		}

	case models.QuestionTypeTrueFalse:
		if answer != "true" && answer != "false" {
			return errors.New("true/false questions must have 'true' or 'false' as answer")
		}

	case models.QuestionTypeEssay, models.QuestionTypeShortAnswer:
		if answer == "" {
			return errors.New("essay and short answer questions must have a sample answer")
		}
	}

	return nil
}