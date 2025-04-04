// internal/delivery/http/auth_handler.go
package http

import (
	"Ev-Charge-Hub/Server/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func TokenValidationHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing", "valid": false})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format", "valid": false})
		return
	}

	token := parts[1]
	claims, err := utils.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "valid": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
		"role":  claims.Role,
	})
}

