package main

import (
	"Ev-Charge-Hub/Server/configs"
	"Ev-Charge-Hub/Server/internal/delivery/http"
	"Ev-Charge-Hub/Server/internal/repository"
	"Ev-Charge-Hub/Server/internal/usecase"
	"Ev-Charge-Hub/Server/routes"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func main() {
	err := godotenv.Load()
	
	if err != nil {
		log.Printf("⚠️ .env not loaded: %v\n", err)
	}

	clientPort := os.Getenv("CLIENT_PORT")
	if clientPort == "" {
		clientPort = "8080"
	}
	port := ":" + clientPort
	

	// 🔧 Release mode
	gin.SetMode(gin.ReleaseMode)

	// ✅ Connect to MongoDB
	db := configs.ConnectDB()

	// ✅ Initialize Dependencies
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := http.NewUserHandler(userUsecase)

	stationRepo := repository.NewEVStationRepository(db)
	stationUsecase := usecase.NewEVStationUsecase(stationRepo)
	stationHandler := http.NewEVStationHandler(stationUsecase)

	// ✅ Set up Router
	router := gin.New()                    // ❌ No default logger
	router.Use(gin.Recovery())             // ✅ Add panic recovery
	router.Use(requestPerformanceLogger()) // ✅ Log API duration

	// ✅ Prometheus /metrics
	p := ginprometheus.NewPrometheus("ev_station")
	p.Use(router)

	// ✅ CORS (can adjust for production)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// ✅ Trusted proxy
	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	// ✅ Register Routes
	routes.SetupRoutes(router, userHandler, stationHandler)
	printRegisteredRoutes(router)

	fmt.Printf("🚀 Server is running on http://localhost%s\n", port)
	router.Run(port)
}

// ✅ Log API performance for each request
func requestPerformanceLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		// format duration to ms
		ms := float64(duration.Microseconds()) / 1000.0

		// align output log
		log.Printf(
			"[PERF] %s | %3d |  %-6s %-25s | %6.2f ms",
			start.Format("2006/01/02 15:04:05"),
			c.Writer.Status(),
			c.Request.Method,
			c.FullPath(),
			ms,
		)
	}
}

// ✅ Optional: Print all registered routes
func printRegisteredRoutes(router *gin.Engine) {
	for _, route := range router.Routes() {
		fmt.Printf("🔗 %s %s\n", route.Method, route.Path)
	}
}
