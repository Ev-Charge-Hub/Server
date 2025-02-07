package routes

import (
	"Ev-Charge-Hub/Server/internal/delivery/http/user"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userHandler user.UserHandlerInterface) {
	// กลุ่มเส้นทาง API สำหรับ User
	r := router.Group("/users")
	r.POST("/register", userHandler.RegisterUser)
	r.POST("/login", userHandler.LoginUser)
}