package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	_ "sample-project/docs"
	"sample-project/internal/config"
	"sample-project/internal/config/cache"
	http "sample-project/internal/delivery/http"
	"sample-project/internal/repository"
	"sample-project/internal/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found or could not be loaded")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Connect to Postgres
	client, err := config.ConnectDB()
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	// Connect to Redis
	cache.ConnectRedis()
	redisClient := cache.GetRedisClient()

	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // Allows all origins (adjust as needed)
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allowed methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Allowed headers
		ExposeHeaders:    []string{"Content-Length"},                          // Expose custom headers
		AllowCredentials: true,                                                // Whether to allow cookies
		MaxAge:           12 * time.Hour,                                      // Caching for preflight request
	}))

	// Initialize Repositories
	userRepo := repository.NewUserRepository(client, redisClient)
	subjectRepo := repository.NewSubjectRepository(client, redisClient)

	// Initialize Usecases
	userUsecase := usecase.NewUserUsecase(userRepo)
	subjectUsecase := usecase.NewSubjectUseCase(subjectRepo)

	// Initialize Handlers
	http.NewUserHandler(router, userUsecase)
	http.NewSubjectHandler(router, subjectUsecase)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	fmt.Printf("☄️  Swagger UI available at: http://localhost:%s/swagger/index.html ☄️\n\n", port)

	router.Run(":" + port)
}