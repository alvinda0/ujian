package handlers

import (
	"school-exam-system/internal/services"
	"school-exam-system/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	userService     *services.UserService
	academicService *services.AcademicService
}

func NewAdminHandler(userService *services.UserService, academicService *services.AcademicService) *AdminHandler {
	return &AdminHandler{
		userService:     userService,
		academicService: academicService,
	}
}

func (h *AdminHandler) CreateUser(c *gin.Context) {
	var req services.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	user, err := h.userService.CreateUser(req)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.CreatedResponse(c, "User created successfully", user)
}

func (h *AdminHandler) GetUsers(c *gin.Context) {
	role := c.Query("role")
	
	users, err := h.userService.ListUsers(role)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.OKResponse(c, "Users retrieved successfully", users)
}

func (h *AdminHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "User not found")
		return
	}

	utils.OKResponse(c, "User retrieved successfully", user)
}

func (h *AdminHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}

	var req services.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	user, err := h.userService.UpdateUser(uint(id), req)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.OKResponse(c, "User updated successfully", user)
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID")
		return
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.OKResponseWithoutData(c, "User deleted successfully")
}

// Class Management
func (h *AdminHandler) CreateClass(c *gin.Context) {
	var req services.CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	class, err := h.academicService.CreateClass(req)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.CreatedResponse(c, "Class created successfully", class)
}

func (h *AdminHandler) GetClasses(c *gin.Context) {
	classes, err := h.academicService.ListClasses()
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.OKResponse(c, "Classes retrieved successfully", classes)
}

func (h *AdminHandler) CreateSubject(c *gin.Context) {
	var req services.CreateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	subject, err := h.academicService.CreateSubject(req)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.CreatedResponse(c, "Subject created successfully", subject)
}

func (h *AdminHandler) GetSubjects(c *gin.Context) {
	subjects, err := h.academicService.ListSubjects()
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.OKResponse(c, "Subjects retrieved successfully", subjects)
}

func (h *AdminHandler) AssignTeacher(c *gin.Context) {
	var req services.AssignTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	if err := h.academicService.AssignTeacherToClass(req); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.OKResponseWithoutData(c, "Teacher assigned successfully")
}

// Student Enrollment Management
func (h *AdminHandler) EnrollStudent(c *gin.Context) {
	studentIDStr := c.Param("studentId")
	studentID, err := strconv.ParseUint(studentIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid student ID")
		return
	}

	var req struct {
		ClassID uint `json:"class_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	if err := h.academicService.EnrollStudent(uint(studentID), req.ClassID); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.OKResponseWithoutData(c, "Student enrolled successfully")
}

func (h *AdminHandler) TransferStudent(c *gin.Context) {
	studentIDStr := c.Param("studentId")
	studentID, err := strconv.ParseUint(studentIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid student ID")
		return
	}

	var req struct {
		NewClassID uint `json:"new_class_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	if err := h.academicService.TransferStudent(uint(studentID), req.NewClassID); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.OKResponseWithoutData(c, "Student transferred successfully")
}

func (h *AdminHandler) GetClassRoster(c *gin.Context) {
	classIDStr := c.Param("classId")
	classID, err := strconv.ParseUint(classIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid class ID")
		return
	}

	students, err := h.academicService.GetClassRoster(uint(classID))
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.OKResponse(c, "Class roster retrieved successfully", students)
}

func (h *AdminHandler) GetTeacherWorkload(c *gin.Context) {
	teacherIDStr := c.Param("teacherId")
	teacherID, err := strconv.ParseUint(teacherIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid teacher ID")
		return
	}

	workload, err := h.academicService.GetTeacherWorkload(uint(teacherID))
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.OKResponse(c, "Teacher workload retrieved successfully", workload)
}