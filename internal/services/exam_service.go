package services

import (
	"errors"
	"school-exam-system/internal/models"
	"time"

	"gorm.io/gorm"
)

type ExamService struct {
	db *gorm.DB
}

type CreateExamRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	SubjectID   uint      `json:"subject_id" binding:"required"`
	ClassID     uint      `json:"class_id" binding:"required"`
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time" binding:"required"`
	Duration    int       `json:"duration" binding:"required"` // minutes
	MaxAttempts int       `json:"max_attempts"`
}

type UpdateExamRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Duration    int       `json:"duration"`
	MaxAttempts int       `json:"max_attempts"`
	Status      string    `json:"status"`
}

func NewExamService(db *gorm.DB) *ExamService {
	return &ExamService{db: db}
}

func (s *ExamService) CreateExam(teacherID uint, req CreateExamRequest) (*models.Exam, error) {
	// Validate teacher has access to this class and subject
	var assignment models.ClassSubject
	if err := s.db.Where("teacher_id = ? AND class_id = ? AND subject_id = ? AND is_active = ?", 
		teacherID, req.ClassID, req.SubjectID, true).First(&assignment).Error; err != nil {
		return nil, errors.New("teacher not assigned to this class and subject")
	}

	// Validate time
	if req.EndTime.Before(req.StartTime) {
		return nil, errors.New("end time must be after start time")
	}

	if req.MaxAttempts == 0 {
		req.MaxAttempts = 1
	}

	exam := models.Exam{
		Title:       req.Title,
		Description: req.Description,
		SubjectID:   req.SubjectID,
		ClassID:     req.ClassID,
		TeacherID:   teacherID,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Duration:    req.Duration,
		MaxAttempts: req.MaxAttempts,
		Status:      models.ExamStatusDraft,
	}

	if err := s.db.Create(&exam).Error; err != nil {
		return nil, err
	}

	// Load relationships
	if err := s.db.Preload("Subject").Preload("Class").Preload("Teacher").First(&exam, exam.ID).Error; err != nil {
		return nil, err
	}

	return &exam, nil
}

func (s *ExamService) GetExamByID(id uint) (*models.Exam, error) {
	var exam models.Exam
	if err := s.db.Preload("Subject").Preload("Class").Preload("Teacher").
		Preload("Questions").First(&exam, id).Error; err != nil {
		return nil, err
	}
	return &exam, nil
}

func (s *ExamService) UpdateExam(id uint, teacherID uint, req UpdateExamRequest) (*models.Exam, error) {
	var exam models.Exam
	if err := s.db.First(&exam, id).Error; err != nil {
		return nil, err
	}

	// Check if teacher owns this exam
	if exam.TeacherID != teacherID {
		return nil, errors.New("unauthorized to update this exam")
	}

	// Update fields
	if req.Title != "" {
		exam.Title = req.Title
	}
	if req.Description != "" {
		exam.Description = req.Description
	}
	if !req.StartTime.IsZero() {
		exam.StartTime = req.StartTime
	}
	if !req.EndTime.IsZero() {
		exam.EndTime = req.EndTime
	}
	if req.Duration > 0 {
		exam.Duration = req.Duration
	}
	if req.MaxAttempts > 0 {
		exam.MaxAttempts = req.MaxAttempts
	}
	if req.Status != "" {
		exam.Status = models.ExamStatus(req.Status)
	}

	if err := s.db.Save(&exam).Error; err != nil {
		return nil, err
	}

	return s.GetExamByID(exam.ID)
}

func (s *ExamService) DeleteExam(id uint, teacherID uint) error {
	var exam models.Exam
	if err := s.db.First(&exam, id).Error; err != nil {
		return err
	}

	// Check if teacher owns this exam
	if exam.TeacherID != teacherID {
		return errors.New("unauthorized to delete this exam")
	}

	// Check if exam has submissions
	var submissionCount int64
	s.db.Model(&models.ExamSubmission{}).Where("exam_id = ?", id).Count(&submissionCount)
	if submissionCount > 0 {
		return errors.New("cannot delete exam with existing submissions")
	}

	return s.db.Delete(&exam).Error
}

func (s *ExamService) ListExamsByTeacher(teacherID uint) ([]models.Exam, error) {
	var exams []models.Exam
	if err := s.db.Preload("Subject").Preload("Class").
		Where("teacher_id = ?", teacherID).
		Order("created_at DESC").Find(&exams).Error; err != nil {
		return nil, err
	}
	return exams, nil
}

func (s *ExamService) ListExamsByStudent(studentID uint) ([]models.Exam, error) {
	// Get student's class
	var student models.Student
	if err := s.db.First(&student, studentID).Error; err != nil {
		return nil, err
	}

	var exams []models.Exam
	if err := s.db.Preload("Subject").Preload("Class").Preload("Teacher").
		Where("class_id = ? AND status IN ?", student.ClassID, 
			[]models.ExamStatus{models.ExamStatusPublished, models.ExamStatusActive}).
		Order("start_time ASC").Find(&exams).Error; err != nil {
		return nil, err
	}

	return exams, nil
}

func (s *ExamService) PublishExam(id uint, teacherID uint) error {
	var exam models.Exam
	if err := s.db.First(&exam, id).Error; err != nil {
		return err
	}

	if exam.TeacherID != teacherID {
		return errors.New("unauthorized to publish this exam")
	}

	if exam.Status != models.ExamStatusDraft {
		return errors.New("only draft exams can be published")
	}

	// Check if exam has questions
	var questionCount int64
	s.db.Model(&models.Question{}).Where("exam_id = ?", id).Count(&questionCount)
	if questionCount == 0 {
		return errors.New("exam must have at least one question before publishing")
	}

	exam.Status = models.ExamStatusPublished
	return s.db.Save(&exam).Error
}

func (s *ExamService) ActivateExam(id uint) error {
	var exam models.Exam
	if err := s.db.First(&exam, id).Error; err != nil {
		return err
	}

	if exam.Status != models.ExamStatusPublished {
		return errors.New("only published exams can be activated")
	}

	if time.Now().Before(exam.StartTime) {
		return errors.New("exam cannot be activated before start time")
	}

	exam.Status = models.ExamStatusActive
	return s.db.Save(&exam).Error
}

func (s *ExamService) CloseExam(id uint) error {
	var exam models.Exam
	if err := s.db.First(&exam, id).Error; err != nil {
		return err
	}

	exam.Status = models.ExamStatusClosed
	return s.db.Save(&exam).Error
}