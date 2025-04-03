package routes

import (
	"Ev-Charge-Hub/Server/internal/delivery/http"
	"Ev-Charge-Hub/Server/middleware"

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
		stationGroup.Use(middleware.AuthMiddleware())
		stationGroup.GET("/filter", stationHandler.FilterStations)
		stationGroup.GET("/:id", stationHandler.GetStationByID)
		stationGroup.PUT("/set-booking", stationHandler.SetBooking)
		stationGroup.GET("", stationHandler.ShowAllStations)
		stationGroup.POST("/create", stationHandler.CreateStation)
		stationGroup.PUT("/:id", stationHandler.EditStation)
		stationGroup.DELETE("/:id", stationHandler.RemoveStation)
		stationGroup.GET("/booking/:username", stationHandler.GetBookingByUserName)
		stationGroup.GET("/bookings/:username", stationHandler.GetBookingsByUserName)	
		stationGroup.GET("/connector/:connector_id", stationHandler.GetStationByConnectorID)
		stationGroup.GET("/username/:username", stationHandler.GetStationByUserName)
	}
}
