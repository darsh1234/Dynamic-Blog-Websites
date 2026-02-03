package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func writeError(c *gin.Context, statusCode int, code, message string, details any) {
	errorObj := gin.H{
		"code":    code,
		"message": message,
	}
	if details != nil {
		errorObj["details"] = details
	}

	c.JSON(statusCode, gin.H{"error": errorObj})
}

func writeValidationError(c *gin.Context, err error) {
	writeError(c, http.StatusBadRequest, "validation_error", "Request validation failed", gin.H{"reason": err.Error()})
}
