package handlers

import (
	"net/http"
	"school-exam-system/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ExamHandler struct {
	examService       *services.ExamService
	questionService   *services.QuestionService
	submissionService *services.SubmissionService
	gradingService    *services.GradingService
	teacherService    *services.TeacherService
	studentService    *services.StudentService
}

func NewExamHandler(
	examService *services.ExamService,
	questionService *services.QuestionService,
	submissionService *services.SubmissionService,
	gradingService *services.GradingService,
	teacherService *services.TeacherService,
	studentService *services.StudentService,
) *ExamHandler {
	return &ExamHandler{
		examService:       examService,
		questionService:   questionService,
		submissionService: submissionService,
		gradingService:    gradingService,
		teacherService:    teacherService,
		studentService:    studentService,
	}
}

// Teacher Exam Management
func (h *ExamHandler) CreateExam(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	teacher, err := h.teacherService.GetTeacherByUserID(uint(userID.(uint)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher profile not found"})
		return
	}

	var req services.CreateExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exam, err := h.examService.CreateExam(teacher.ID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    exam,
	})
}

func (h *ExamHandler) GetMyExams(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	teacher, err := h.teacherService.GetTeacherByUserID(uint(userID.(uint)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher profile not found"})
		return
	}

	exams, err := h.examService.ListExamsByTeacher(teacher.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    exams,
	})
}

func (h *ExamHandler) GetExam(c *gin.Context) {
	examIDStr := c.Param("id")
	examID, err := strconv.ParseUint(examIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exam ID"})
		return
	}

	exam, err := h.examService.GetExamByID(uint(examID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exam not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    exam,
	})
}

func (h *ExamHandler) UpdateExam(c *gin.Context) {
	userID, _ := c.Get("user_id")
	examIDStr := c.Param("id")
	examID, err := strconv.ParseUint(examIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exam ID"})
		return
	}

	teacher, err := h.teacherService.GetTeacherByUserID(uint(userID.(uint)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher profile not found"})
		return
	}

	var req services.UpdateExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exam, err := h.examService.UpdateExam(uint(examID), teacher.ID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    exam,
	})
}

func (h *ExamHandler) DeleteExam(c *gin.Context) {
	userID, _ := c.Get("user_id")
	examIDStr := c.Param("id")
	examID, err := strconv.ParseUint(examIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exam ID"})
		return
	}

	teacher, err := h.teacherService.GetTeacherByUserID(uint(userID.(uint)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher profile not found"})
		return
	}

	if err := h.examService.DeleteExam(uint(examID), teacher.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Exam deleted successfully",
	})
}

func (h *ExamHandler) PublishExam(c *gin.Context) {
	userID, _ := c.Get("user_id")
	examIDStr := c.Param("id")
	examID, err := strconv.ParseUint(examIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exam ID"})
		return
	}

	teacher, err := h.teacherService.GetTeacherByUserID(uint(userID.(uint)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher profile not found"})
		return
	}

	if err := h.examService.PublishExam(uint(examID), teacher.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Exam published successfully",
	})
}

// Question Management
func (h *ExamHandler) AddQuestion(c *gin.Context) {
	userID, _ := c.Get("user_id")
	examIDStr := c.Param("id")
	examID, err := strconv.ParseUint(examIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exam ID"})
		return
	}

	teacher, err := h.teacherService.GetTeacherByUserID(uint(userID.(uint)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher profile not found"})
		return
	}

	var req services.AddQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question, err := h.questionService.AddQuestion(uint(examID), teacher.ID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    question,
	})
}

func (h *ExamHandler) GetQuestions(c *gin.Context) {
	examIDStr := c.Param("id")
	examID, err := strconv.ParseUint(examIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exam ID"})
		return
	}

	questions, err := h.questionService.GetQuestionsByExam(uint(examID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    questions,
	})
}

// Student Exam Taking
func (h *ExamHandler) GetAvailableExams(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	student, err := h.studentService.GetStudentByUserID(uint(userID.(uint)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student profile not found"})
		return
	}

	exams, err := h.examService.ListExamsByStudent(student.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    exams,
	})
}

func (h *ExamHandler) StartExam(c *gin.Context) {
	userID, _ := c.Get("user_id")
	examIDStr := c.Param("id")
	examID, err := strconv.ParseUint(examIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exam ID"})
		return
	}

	student, err := h.studentService.GetStudentByUserID(uint(userID.(uint)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student profile not found"})
		return
	}

	submission, err := h.submissionService.StartExam(uint(examID), student.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    submission,
	})
}

func (h *ExamHandler) SubmitAnswer(c *gin.Context) {
	userID, _ := c.Get("user_id")
	submissionIDStr := c.Param("submissionId")
	submissionID, err := strconv.ParseUint(submissionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
		return
	}

	student, err := h.studentService.GetStudentByUserID(uint(userID.(uint)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student profile not found"})
		return
	}

	var req services.SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.submissionService.SubmitAnswer(uint(submissionID), student.ID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Answer submitted successfully",
	})
}

func (h *ExamHandler) FinishExam(c *gin.Context) {
	userID, _ := c.Get("user_id")
	submissionIDStr := c.Param("submissionId")
	submissionID, err := strconv.ParseUint(submissionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
		return
	}

	student, err := h.studentService.GetStudentByUserID(uint(userID.(uint)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student profile not found"})
		return
	}

	submission, err := h.submissionService.FinishExam(uint(submissionID), student.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    submission,
	})
}

// Grading
func (h *ExamHandler) GetSubmissions(c *gin.Context) {
	userID, _ := c.Get("user_id")
	examIDStr := c.Param("id")
	examID, err := strconv.ParseUint(examIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exam ID"})
		return
	}

	teacher, err := h.teacherService.GetTeacherByUserID(uint(userID.(uint)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher profile not found"})
		return
	}

	submissions, err := h.gradingService.GetSubmissionGrades(uint(examID), teacher.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    submissions,
	})
}

func (h *ExamHandler) GetMyGrades(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	student, err := h.studentService.GetStudentByUserID(uint(userID.(uint)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student profile not found"})
		return
	}

	grades, err := h.gradingService.GetStudentGrades(student.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    grades,
	})
}