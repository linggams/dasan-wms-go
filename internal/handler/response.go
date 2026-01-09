package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
type Response struct {
	Status  string      `json:"status"`
	Success bool        `json:"success,omitempty"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   bool        `json:"error,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// SuccessResponse sends a success response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Status:  "success",
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, message string, errors interface{}) {
	c.JSON(statusCode, Response{
		Status:  "error",
		Error:   true,
		Message: message,
		Errors:  errors,
	})
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *gin.Context, message string, errors map[string][]string) {
	c.JSON(http.StatusUnprocessableEntity, gin.H{
		"message": message,
		"errors":  errors,
	})
}
