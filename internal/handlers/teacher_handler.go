package handlers

import (
	"net/http"
	"school-exam-system/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TeacherHandler struct {
	academicService *services.AcademicService
	teacherService  *services.TeacherService
}

func NewTeacherHandler(academicService *services.AcademicService, teacherService *services.TeacherService) *TeacherHandler {
	return &TeacherHandler{
		academicService: academicService,
		teacherService:  teacherService,
	}
}

func (h *TeacherHandler) GetMyClasses(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	// Get teacher ID from user ID
	teacherID, err := h.getTeacherIDFromUserID(uint(userID.(uint)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher profile not found"})
		return
	}

	assignments, err := h.academicService.GetTeacherAssignments(teacherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    assignments,
	})
}

func (h *TeacherHandler) GetMyStudents(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	// Get teacher ID from user ID
	teacherID, err := h.getTeacherIDFromUserID(uint(userID.(uint)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher profile not found"})
		return
	}

	students, err := h.academicService.GetStudentsByTeacher(teacherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    students,
	})
}

func (h *TeacherHandler) GetStudentsByClass(c *gin.Context) {
	classIDStr := c.Param("classId")
	classID, err := strconv.ParseUint(classIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	// Get teacher ID from user ID
	teacherID, err := h.getTeacherIDFromUserID(uint(userID.(uint)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher profile not found"})
		return
	}

	// Verify teacher has access to this class
	assignments, err := h.academicService.GetTeacherAssignments(teacherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	hasAccess := false
	for _, assignment := range assignments {
		if assignment.ClassID == uint(classID) {
			hasAccess = true
			break
		}
	}

	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this class"})
		return
	}

	students, err := h.academicService.GetStudentsByClass(uint(classID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    students,
	})
}

func (h *TeacherHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	teacher, err := h.teacherService.GetTeacherProfile(uint(userID.(uint)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher profile not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    teacher,
	})
}

// Helper function to get teacher ID from user ID
func (h *TeacherHandler) getTeacherIDFromUserID(userID uint) (uint, error) {
	teacher, err := h.teacherService.GetTeacherByUserID(userID)
	if err != nil {
		return 0, err
	}
	return teacher.ID, nil
}