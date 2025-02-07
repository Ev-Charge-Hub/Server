package main

import (
	"Ev-Charge-Hub/Server/configs"
	"Ev-Charge-Hub/Server/internal/delivery/http"
	"Ev-Charge-Hub/Server/internal/repository"
	"Ev-Charge-Hub/Server/internal/usecase"
	"Ev-Charge-Hub/Server/routes"
	"fmt"
	"log"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	db := configs.ConnectDB()

	// Init -> Repository, Use Case and Handler
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := http.NewUserHandler(userUsecase)

	stationRepo := repository.NewEVStationRepository(db)
	stationUsecase := usecase.NewEVStationUsecase(stationRepo)
	stationHandler := http.NewEVStationHandler(stationUsecase)


	// Set -> Routing
	router := gin.Default()
	err := router.SetTrustedProxies(nil) 
	if err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	routes.SetupRoutes(router, userHandler, stationHandler)

	port := ":8080"
	fmt.Printf("Server is running on http://localhost%s\n", port)
	router.Run(port)
	printRegisteredRoutes(router)

}

func printRegisteredRoutes(router *gin.Engine) {
	for _, route := range router.Routes() {
		fmt.Printf("Registered route: %s %s\n", route.Method, route.Path)
	}
}