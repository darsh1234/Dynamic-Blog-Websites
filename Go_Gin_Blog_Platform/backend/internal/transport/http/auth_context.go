package http

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func currentUserFromContext(c *gin.Context) (userID string, role string, ok bool) {
	userValue, userExists := c.Get(ContextKeyUserID)
	roleValue, roleExists := c.Get(ContextKeyRole)
	if !userExists || !roleExists {
		return "", "", false
	}

	userID, okUser := userValue.(string)
	role, okRole := roleValue.(string)
	if !okUser || !okRole {
		return "", "", false
	}

	userID = strings.TrimSpace(userID)
	role = strings.TrimSpace(role)
	if userID == "" || role == "" {
		return "", "", false
	}

	return userID, role, true
}
