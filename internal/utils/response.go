package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success responses
func SuccessResponse(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

func CreatedResponse(c *gin.Context, message string, data interface{}) {
	SuccessResponse(c, http.StatusCreated, message, data)
}

func OKResponse(c *gin.Context, message string, data interface{}) {
	SuccessResponse(c, http.StatusOK, message, data)
}

func OKResponseWithoutData(c *gin.Context, message string) {
	c.JSON(http.StatusOK, APIResponse{
		Status:  http.StatusOK,
		Message: message,
	})
}

// Error responses
func ErrorResponse(c *gin.Context, status int, message string, error string) {
	c.JSON(status, APIResponse{
		Status:  status,
		Message: message,
		Error:   error,
	})
}

func BadRequestResponse(c *gin.Context, error string) {
	ErrorResponse(c, http.StatusBadRequest, "Bad Request", error)
}

func UnauthorizedResponse(c *gin.Context, error string) {
	ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", error)
}

func ForbiddenResponse(c *gin.Context, error string) {
	ErrorResponse(c, http.StatusForbidden, "Forbidden", error)
}

func NotFoundResponse(c *gin.Context, error string) {
	ErrorResponse(c, http.StatusNotFound, "Not Found", error)
}

func InternalServerErrorResponse(c *gin.Context, error string) {
	ErrorResponse(c, http.StatusInternalServerError, "Internal Server Error", error)
}