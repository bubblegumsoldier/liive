package main

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/liive/backend/liive-rest-api/docs" // Import generated Swagger docs
	"github.com/liive/backend/liive-rest-api/internal/handlers"
	"github.com/liive/backend/liive-rest-api/internal/service"
	"github.com/liive/backend/shared/pkg/database"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title REST API
// @version 1.0
// @description This is the REST API service
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /
// @schemes http
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize services
	chatService := service.NewChatService(db)

	// Initialize handlers
	chatHandler := handlers.NewChatHandler(chatService)

	// Create Echo instance
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Swagger documentation
	e.GET("/swagger", echoSwagger.EchoWrapHandler())
	e.GET("/swagger/*", echoSwagger.EchoWrapHandler())

	// Chat routes
	api := e.Group("/api")
	// TODO: Add JWT middleware
	// api.Use(middleware.JWT([]byte(os.Getenv("JWT_SECRET"))))

	chats := api.Group("/chats")
	chats.POST("", chatHandler.CreateChat)
	chats.GET("", chatHandler.GetUserChats)
	chats.GET("/:id", chatHandler.GetChat)
	chats.PUT("/:id", chatHandler.UpdateChatTitle)
	chats.POST("/:id/leave", chatHandler.LeaveChat)
	chats.POST("/:id/members", chatHandler.AddMembers)
	chats.DELETE("/:id/members/:userId", chatHandler.RemoveMember)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "healthy",
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
