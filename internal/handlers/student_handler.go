package handlers

import (
	"net/http"
	"school-exam-system/internal/services"

	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	studentService  *services.StudentService
	academicService *services.AcademicService
}

func NewStudentHandler(studentService *services.StudentService, academicService *services.AcademicService) *StudentHandler {
	return &StudentHandler{
		studentService:  studentService,
		academicService: academicService,
	}
}

func (h *StudentHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	student, err := h.studentService.GetStudentProfile(uint(userID.(uint)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student profile not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    student,
	})
}

func (h *StudentHandler) GetMySubjects(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	subjects, err := h.studentService.GetStudentSubjects(uint(userID.(uint)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    subjects,
	})
}

func (h *StudentHandler) GetMyClass(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	student, err := h.studentService.GetStudentByUserID(uint(userID.(uint)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	class, err := h.academicService.GetClassByID(student.ClassID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    class,
	})
}

func (h *StudentHandler) GetClassmates(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	classmates, err := h.studentService.GetClassmates(uint(userID.(uint)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    classmates,
	})
}