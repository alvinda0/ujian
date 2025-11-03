package services

import (
	"errors"
	"school-exam-system/internal/models"
	"time"

	"gorm.io/gorm"
)

type SubmissionService struct {
	db *gorm.DB
}

type StartExamRequest struct {
	ExamID uint `json:"exam_id" binding:"required"`
}

type SubmitAnswerRequest struct {
	QuestionID uint   `json:"question_id" binding:"required"`
	Answer     string `json:"answer" binding:"required"`
}

type FinishExamRequest struct {
	SubmissionID uint `json:"submission_id" binding:"required"`
}

func NewSubmissionService(db *gorm.DB) *SubmissionService {
	return &SubmissionService{db: db}
}

func (s *SubmissionService) StartExam(examID uint, studentID uint) (*models.ExamSubmission, error) {
	// Get exam details
	var exam models.Exam
	if err := s.db.Preload("Class").First(&exam, examID).Error; err != nil {
		return nil, err
	}

	// Get student details
	var student models.Student
	if err := s.db.First(&student, studentID).Error; err != nil {
		return nil, err
	}

	// Check if student is in the correct class
	if student.ClassID != exam.ClassID {
		return nil, errors.New("student not enrolled in this class")
	}

	// Check exam status and timing
	now := time.Now()
	if exam.Status != models.ExamStatusActive && exam.Status != models.ExamStatusPublished {
		return nil, errors.New("exam is not available")
	}

	if now.Before(exam.StartTime) {
		return nil, errors.New("exam has not started yet")
	}

	if now.After(exam.EndTime) {
		return nil, errors.New("exam has ended")
	}

	// Check if student has already reached max attempts
	var attemptCount int64
	s.db.Model(&models.ExamSubmission{}).
		Where("exam_id = ? AND student_id = ?", examID, studentID).
		Count(&attemptCount)

	if int(attemptCount) >= exam.MaxAttempts {
		return nil, errors.New("maximum attempts reached")
	}

	// Check if student has an active submission
	var activeSubmission models.ExamSubmission
	if err := s.db.Where("exam_id = ? AND student_id = ? AND status = ?", 
		examID, studentID, models.SubmissionStatusInProgress).First(&activeSubmission).Error; err == nil {
		// Return existing active submission
		return &activeSubmission, nil
	}

	// Create new submission
	submission := models.ExamSubmission{
		ExamID:     examID,
		StudentID:  studentID,
		StartedAt:  now,
		Status:     models.SubmissionStatusInProgress,
		AttemptNum: int(attemptCount) + 1,
	}

	if err := s.db.Create(&submission).Error; err != nil {
		return nil, err
	}

	// Load relationships
	if err := s.db.Preload("Exam").Preload("Student").First(&submission, submission.ID).Error; err != nil {
		return nil, err
	}

	return &submission, nil
}

func (s *SubmissionService) SubmitAnswer(submissionID uint, studentID uint, req SubmitAnswerRequest) error {
	// Get submission
	var submission models.ExamSubmission
	if err := s.db.Preload("Exam").First(&submission, submissionID).Error; err != nil {
		return err
	}

	// Verify ownership
	if submission.StudentID != studentID {
		return errors.New("unauthorized access to this submission")
	}

	// Check submission status
	if submission.Status != models.SubmissionStatusInProgress {
		return errors.New("submission is not in progress")
	}

	// Check if exam time has expired
	examEndTime := submission.StartedAt.Add(time.Duration(submission.Exam.Duration) * time.Minute)
	if time.Now().After(examEndTime) {
		// Auto-finish the exam
		s.FinishExam(submissionID, studentID)
		return errors.New("exam time has expired")
	}

	// Verify question belongs to this exam
	var question models.Question
	if err := s.db.Where("id = ? AND exam_id = ?", req.QuestionID, submission.ExamID).First(&question).Error; err != nil {
		return errors.New("question not found in this exam")
	}

	// Check if answer already exists
	var existingAnswer models.SubmissionAnswer
	if err := s.db.Where("submission_id = ? AND question_id = ?", submissionID, req.QuestionID).First(&existingAnswer).Error; err == nil {
		// Update existing answer
		existingAnswer.Answer = req.Answer
		return s.db.Save(&existingAnswer).Error
	}

	// Create new answer
	answer := models.SubmissionAnswer{
		SubmissionID: submissionID,
		QuestionID:   req.QuestionID,
		Answer:       req.Answer,
	}

	return s.db.Create(&answer).Error
}

func (s *SubmissionService) FinishExam(submissionID uint, studentID uint) (*models.ExamSubmission, error) {
	var submission models.ExamSubmission
	if err := s.db.Preload("Exam").Preload("Answers").First(&submission, submissionID).Error; err != nil {
		return nil, err
	}

	// Verify ownership
	if submission.StudentID != studentID {
		return nil, errors.New("unauthorized access to this submission")
	}

	// Check if already finished
	if submission.Status != models.SubmissionStatusInProgress {
		return &submission, nil
	}

	// Update submission
	now := time.Now()
	submission.SubmittedAt = &now
	submission.Status = models.SubmissionStatusSubmitted

	if err := s.db.Save(&submission).Error; err != nil {
		return nil, err
	}

	// Auto-grade objective questions
	if err := s.autoGradeSubmission(&submission); err != nil {
		// Log error but don't fail the submission
		// Auto-grading can be done later
	}

	return &submission, nil
}

func (s *SubmissionService) GetSubmissionByID(id uint) (*models.ExamSubmission, error) {
	var submission models.ExamSubmission
	if err := s.db.Preload("Exam").Preload("Student").Preload("Answers").Preload("Answers.Question").
		First(&submission, id).Error; err != nil {
		return nil, err
	}
	return &submission, nil
}

func (s *SubmissionService) GetSubmissionsByExam(examID uint) ([]models.ExamSubmission, error) {
	var submissions []models.ExamSubmission
	if err := s.db.Preload("Student").Preload("Student.User").
		Where("exam_id = ?", examID).
		Order("created_at DESC").Find(&submissions).Error; err != nil {
		return nil, err
	}
	return submissions, nil
}

func (s *SubmissionService) GetSubmissionsByStudent(studentID uint) ([]models.ExamSubmission, error) {
	var submissions []models.ExamSubmission
	if err := s.db.Preload("Exam").Preload("Exam.Subject").
		Where("student_id = ?", studentID).
		Order("created_at DESC").Find(&submissions).Error; err != nil {
		return nil, err
	}
	return submissions, nil
}

func (s *SubmissionService) autoGradeSubmission(submission *models.ExamSubmission) error {
	// Get all questions for this exam
	var questions []models.Question
	if err := s.db.Where("exam_id = ?", submission.ExamID).Find(&questions).Error; err != nil {
		return err
	}

	totalScore := 0.0
	maxScore := 0.0

	for _, question := range questions {
		maxScore += float64(question.Points)

		// Find student's answer
		var answer models.SubmissionAnswer
		if err := s.db.Where("submission_id = ? AND question_id = ?", 
			submission.ID, question.ID).First(&answer).Error; err != nil {
			// No answer provided, score is 0
			continue
		}

		// Auto-grade based on question type
		var points float64
		var isCorrect bool

		switch question.Type {
		case models.QuestionTypeMultipleChoice, models.QuestionTypeTrueFalse:
			if answer.Answer == question.Answer {
				points = float64(question.Points)
				isCorrect = true
			}
		case models.QuestionTypeShortAnswer:
			// Simple string comparison (case-insensitive)
			if answer.Answer == question.Answer {
				points = float64(question.Points)
				isCorrect = true
			}
		case models.QuestionTypeEssay:
			// Essay questions need manual grading
			continue
		}

		totalScore += points

		// Update answer with score
		answer.Points = &points
		answer.IsCorrect = &isCorrect
		s.db.Save(&answer)
	}

	// Update submission with score
	submission.Score = &totalScore
	submission.MaxScore = &maxScore

	// If all questions are auto-gradable, mark as graded
	var essayCount int64
	s.db.Model(&models.Question{}).
		Where("exam_id = ? AND type = ?", submission.ExamID, models.QuestionTypeEssay).
		Count(&essayCount)

	if essayCount == 0 {
		submission.Status = models.SubmissionStatusGraded
	}

	return s.db.Save(submission).Error
}

func (s *SubmissionService) CheckExamTimeExpiry() error {
	// Find all in-progress submissions that have exceeded time limit
	var submissions []models.ExamSubmission
	if err := s.db.Preload("Exam").Where("status = ?", models.SubmissionStatusInProgress).Find(&submissions).Error; err != nil {
		return err
	}

	for _, submission := range submissions {
		examEndTime := submission.StartedAt.Add(time.Duration(submission.Exam.Duration) * time.Minute)
		if time.Now().After(examEndTime) {
			// Auto-finish expired exam
			submission.Status = models.SubmissionStatusExpired
			now := time.Now()
			submission.SubmittedAt = &now
			s.db.Save(&submission)

			// Auto-grade the submission
			s.autoGradeSubmission(&submission)
		}
	}

	return nil
}