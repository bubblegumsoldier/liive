package main

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/liive/backend/liive-user-manager/docs" // Import generated Swagger docs
	"github.com/liive/backend/liive-user-manager/internal/handlers"
	"github.com/liive/backend/liive-user-manager/internal/service"
	"github.com/liive/backend/shared/pkg/database"
	"github.com/liive/backend/shared/pkg/auth"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title User Manager API
// @version 1.0
// @description This is the user management service API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8083
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

	// Initialize services and handlers
	userService := service.NewUserService(db)
	userHandler := handlers.NewUserHandler(userService)

	// Create Echo instance
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// JWT middleware
	jwtMiddleware := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte(os.Getenv("JWT_SECRET_KEY")),
		Claims:      &auth.Claims{},
		TokenLookup: "header:Authorization",
		AuthScheme:  "Bearer",
	})

	// Swagger documentation
	e.GET("/swagger", echoSwagger.EchoWrapHandler())
	e.GET("/swagger/*", echoSwagger.EchoWrapHandler())

	// Public routes
	e.POST("/register", userHandler.Register)
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "healthy",
		})
	})

	// Protected routes
	api := e.Group("")
	api.Use(jwtMiddleware)
	api.GET("/profile", userHandler.GetProfile)
	api.PUT("/profile", userHandler.UpdateProfile)
	api.POST("/change-password", userHandler.UpdatePassword)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
