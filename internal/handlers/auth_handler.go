package handlers

import (
	"school-exam-system/internal/services"
	"school-exam-system/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	response, err := h.authService.Login(req)
	if err != nil {
		utils.UnauthorizedResponse(c, err.Error())
		return
	}

	utils.OKResponse(c, "Login successful", response)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// For JWT, logout is typically handled client-side by removing the token
	// Here we just return success
	utils.OKResponseWithoutData(c, "Logged out successfully")
}

func (h *AuthHandler) Profile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	profileData := gin.H{
		"user_id":  userID,
		"username": username,
		"role":     role,
	}

	utils.OKResponse(c, "Profile retrieved successfully", profileData)
}