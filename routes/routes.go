package routes

import (
	"Ev-Charge-Hub/Server/internal/delivery/http"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userHandler http.UserHandlerInterface) {
	r := router.Group("/users")
	r.POST("/register", userHandler.RegisterUser)
	r.POST("/login", userHandler.LoginUser)
}