package routes

import (
	"Ev-Charge-Hub/Server/internal/delivery/http"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userHandler http.UserHandlerInterface, stationHandler *http.EVStationHandler) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("/register", userHandler.RegisterUser)
		userGroup.POST("/login", userHandler.LoginUser)
	}

	stationGroup := router.Group("/stations")
	{
		stationGroup.GET("/filter", stationHandler.FilterStations)
		stationGroup.GET("/:id", stationHandler.GetStationByID)
		stationGroup.PUT("/set-booking", stationHandler.SetBooking)
		stationGroup.GET("", stationHandler.ShowAllStations)
		stationGroup.POST("/create", stationHandler.CreateStation)
		stationGroup.PUT("/:id", stationHandler.EditStation)
		stationGroup.DELETE("/:id", stationHandler.RemoveStation)
	}
}
