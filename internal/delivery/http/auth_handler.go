// internal/delivery/http/auth_handler.go
package http

import (
	"Ev-Charge-Hub/Server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TokenValidationHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"valid": false,
			"error": "Authorization header missing",
		})
		return
	}

	tokenString := authHeader[len("Bearer "):]

	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"valid": false,
			"error": "Invalid token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
		"role":  claims.Role,
	})
}
