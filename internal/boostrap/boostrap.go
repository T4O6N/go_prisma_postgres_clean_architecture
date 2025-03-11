package boostrap

// import (
// 	"log"
// 	"log/slog"
// 	"os"
// 	"time"

// 	_ "sample-project/docs"
// 	"sample-project/internal/config"
// 	http "sample-project/internal/delivery/http"
// 	"sample-project/internal/repository"
// 	"sample-project/internal/usecase"

// 	"github.com/gin-contrib/cors"
// 	"github.com/gin-gonic/gin"
// 	swaggerFiles "github.com/swaggo/files"
// 	ginSwagger "github.com/swaggo/gin-swagger"
// )

// type App struct {
// 	Router *gin.Engine
// }

// func (a *App) Initialize() {
// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "8080"
// 	}

// 	client, err := config.ConnectDB()
// 	if err != nil {
// 		slog.Error("Failed to connect to database", "error", err)
// 		os.Exit(1)
// 	}

// 	a.Router = gin.Default()

// 	// Configure CORS
// 	a.Router.Use(cors.New(cors.Config{
// 		AllowOrigins:     []string{"*"},                                       // Allows all origins (adjust as needed)
// 		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allowed methods
// 		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Allowed headers
// 		ExposeHeaders:    []string{"Content-Length"},                          // Expose custom headers
// 		AllowCredentials: true,                                                // Whether to allow cookies
// 		MaxAge:           12 * time.Hour,                                      // Caching for preflight request
// 	}))

// 	// Initialize Repositories
// 	userRepo := repository.NewUserRepository(client)
// 	subjectRepo := repository.NewSubjectRepository(client)

// 	// Initialize Usecases
// 	userUsecase := usecase.NewUserUsecase(userRepo)
// 	subjectUsecase := usecase.NewSubjectUseCase(subjectRepo)

// 	// Initialize Handlers
// 	http.NewUserHandler(a.Router, userUsecase)
// 	http.NewSubjectHandler(a.Router, subjectUsecase)

// 	// Swagger
// 	a.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

// 	log.Println("Server running on :8080")
// 	a.Router.Run(":" + port)
// }
