package services

import (
	"errors"
	"school-exam-system/internal/models"

	"gorm.io/gorm"
)

type GradingService struct {
	db *gorm.DB
}

type GradeSubmissionRequest struct {
	Score float64 `json:"score" binding:"required"`
}

type GradeAnswerRequest struct {
	Points    float64 `json:"points" binding:"required"`
	IsCorrect bool    `json:"is_correct"`
	Feedback  string  `json:"feedback"`
}

type GradeStatistics struct {
	TotalSubmissions int     `json:"total_submissions"`
	AverageScore     float64 `json:"average_score"`
	HighestScore     float64 `json:"highest_score"`
	LowestScore      float64 `json:"lowest_score"`
	PassRate         float64 `json:"pass_rate"`
	ScoreDistribution map[string]int `json:"score_distribution"`
}

func NewGradingService(db *gorm.DB) *GradingService {
	return &GradingService{db: db}
}

func (s *GradingService) GradeSubmission(submissionID uint, teacherID uint, req GradeSubmissionRequest) error {
	var submission models.ExamSubmission
	if err := s.db.Preload("Exam").First(&submission, submissionID).Error; err != nil {
		return err
	}

	// Verify teacher owns this exam
	if submission.Exam.TeacherID != teacherID {
		return errors.New("unauthorized to grade this submission")
	}

	// Update submission score
	submission.Score = &req.Score
	submission.Status = models.SubmissionStatusGraded

	return s.db.Save(&submission).Error
}

func (s *GradingService) GradeAnswer(answerID uint, teacherID uint, req GradeAnswerRequest) error {
	var answer models.SubmissionAnswer
	if err := s.db.Preload("Submission").Preload("Submission.Exam").Preload("Question").
		First(&answer, answerID).Error; err != nil {
		return err
	}

	// Verify teacher owns this exam
	if answer.Submission.Exam.TeacherID != teacherID {
		return errors.New("unauthorized to grade this answer")
	}

	// Validate points don't exceed question points
	if req.Points > float64(answer.Question.Points) {
		return errors.New("points cannot exceed question maximum points")
	}

	// Update answer
	answer.Points = &req.Points
	answer.IsCorrect = &req.IsCorrect

	if err := s.db.Save(&answer).Error; err != nil {
		return err
	}

	// Recalculate submission total score
	return s.recalculateSubmissionScore(answer.SubmissionID)
}

func (s *GradingService) recalculateSubmissionScore(submissionID uint) error {
	var submission models.ExamSubmission
	if err := s.db.First(&submission, submissionID).Error; err != nil {
		return err
	}

	// Get all answers for this submission
	var answers []models.SubmissionAnswer
	if err := s.db.Preload("Question").Where("submission_id = ?", submissionID).Find(&answers).Error; err != nil {
		return err
	}

	totalScore := 0.0
	maxScore := 0.0
	gradedAnswers := 0

	for _, answer := range answers {
		maxScore += float64(answer.Question.Points)
		
		if answer.Points != nil {
			totalScore += *answer.Points
			gradedAnswers++
		}
	}

	// Update submission
	submission.Score = &totalScore
	submission.MaxScore = &maxScore

	// If all answers are graded, mark submission as graded
	if gradedAnswers == len(answers) {
		submission.Status = models.SubmissionStatusGraded
	}

	return s.db.Save(&submission).Error
}

func (s *GradingService) GetSubmissionGrades(examID uint, teacherID uint) ([]models.ExamSubmission, error) {
	// Verify teacher owns this exam
	var exam models.Exam
	if err := s.db.First(&exam, examID).Error; err != nil {
		return nil, err
	}

	if exam.TeacherID != teacherID {
		return nil, errors.New("unauthorized to view grades for this exam")
	}

	var submissions []models.ExamSubmission
	if err := s.db.Preload("Student").Preload("Student.User").
		Where("exam_id = ?", examID).
		Order("score DESC").Find(&submissions).Error; err != nil {
		return nil, err
	}

	return submissions, nil
}

func (s *GradingService) GetStudentGrades(studentID uint) ([]models.ExamSubmission, error) {
	var submissions []models.ExamSubmission
	if err := s.db.Preload("Exam").Preload("Exam.Subject").
		Where("student_id = ? AND status = ?", studentID, models.SubmissionStatusGraded).
		Order("created_at DESC").Find(&submissions).Error; err != nil {
		return nil, err
	}

	return submissions, nil
}

func (s *GradingService) GetExamStatistics(examID uint, teacherID uint) (*GradeStatistics, error) {
	// Verify teacher owns this exam
	var exam models.Exam
	if err := s.db.First(&exam, examID).Error; err != nil {
		return nil, err
	}

	if exam.TeacherID != teacherID {
		return nil, errors.New("unauthorized to view statistics for this exam")
	}

	var submissions []models.ExamSubmission
	if err := s.db.Where("exam_id = ? AND status = ?", examID, models.SubmissionStatusGraded).
		Find(&submissions).Error; err != nil {
		return nil, err
	}

	if len(submissions) == 0 {
		return &GradeStatistics{
			TotalSubmissions: 0,
			ScoreDistribution: make(map[string]int),
		}, nil
	}

	// Calculate statistics
	var totalScore float64
	var highestScore float64
	var lowestScore float64 = 100.0 // Assuming percentage
	var passCount int
	scoreDistribution := make(map[string]int)

	for i, submission := range submissions {
		if submission.Score == nil || submission.MaxScore == nil {
			continue
		}

		// Calculate percentage
		percentage := (*submission.Score / *submission.MaxScore) * 100

		totalScore += percentage

		if i == 0 {
			highestScore = percentage
			lowestScore = percentage
		} else {
			if percentage > highestScore {
				highestScore = percentage
			}
			if percentage < lowestScore {
				lowestScore = percentage
			}
		}

		// Count passes (assuming 60% is passing)
		if percentage >= 60 {
			passCount++
		}

		// Score distribution
		var range_ string
		switch {
		case percentage >= 90:
			range_ = "A (90-100)"
		case percentage >= 80:
			range_ = "B (80-89)"
		case percentage >= 70:
			range_ = "C (70-79)"
		case percentage >= 60:
			range_ = "D (60-69)"
		default:
			range_ = "F (0-59)"
		}
		scoreDistribution[range_]++
	}

	averageScore := totalScore / float64(len(submissions))
	passRate := (float64(passCount) / float64(len(submissions))) * 100

	return &GradeStatistics{
		TotalSubmissions:  len(submissions),
		AverageScore:      averageScore,
		HighestScore:      highestScore,
		LowestScore:       lowestScore,
		PassRate:          passRate,
		ScoreDistribution: scoreDistribution,
	}, nil
}

func (s *GradingService) GetPendingGrades(teacherID uint) ([]models.ExamSubmission, error) {
	var submissions []models.ExamSubmission
	if err := s.db.Preload("Exam").Preload("Student").Preload("Student.User").
		Joins("JOIN exams ON exam_submissions.exam_id = exams.id").
		Where("exams.teacher_id = ? AND exam_submissions.status IN ?", 
			teacherID, []models.SubmissionStatus{models.SubmissionStatusSubmitted}).
		Order("exam_submissions.created_at ASC").Find(&submissions).Error; err != nil {
		return nil, err
	}

	return submissions, nil
}

func (s *GradingService) GetAnswersForGrading(submissionID uint, teacherID uint) ([]models.SubmissionAnswer, error) {
	var submission models.ExamSubmission
	if err := s.db.Preload("Exam").First(&submission, submissionID).Error; err != nil {
		return nil, err
	}

	// Verify teacher owns this exam
	if submission.Exam.TeacherID != teacherID {
		return nil, errors.New("unauthorized to view answers for this submission")
	}

	var answers []models.SubmissionAnswer
	if err := s.db.Preload("Question").Where("submission_id = ?", submissionID).
		Order("question_id ASC").Find(&answers).Error; err != nil {
		return nil, err
	}

	return answers, nil
}

func (s *GradingService) AutoGradeExam(examID uint, teacherID uint) error {
	// Verify teacher owns this exam
	var exam models.Exam
	if err := s.db.First(&exam, examID).Error; err != nil {
		return err
	}

	if exam.TeacherID != teacherID {
		return errors.New("unauthorized to auto-grade this exam")
	}

	// Get all submitted submissions for this exam
	var submissions []models.ExamSubmission
	if err := s.db.Where("exam_id = ? AND status = ?", examID, models.SubmissionStatusSubmitted).
		Find(&submissions).Error; err != nil {
		return err
	}

	// Auto-grade each submission
	for _, submission := range submissions {
		if err := s.autoGradeSubmission(&submission); err != nil {
			// Log error but continue with other submissions
			continue
		}
	}

	return nil
}

func (s *GradingService) autoGradeSubmission(submission *models.ExamSubmission) error {
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