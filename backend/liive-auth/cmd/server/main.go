package main

import (
    "log"
    "os"

    "github.com/go-playground/validator/v10"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    _ "github.com/liive/backend/liive-auth/docs" // Import generated Swagger docs
    "github.com/liive/backend/liive-auth/internal/handlers"
    "github.com/liive/backend/liive-auth/internal/service"
    "github.com/liive/backend/shared/pkg/database"
    echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Auth API
// @version 1.0
// @description This is the authentication service API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8082
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
    authService := service.NewAuthService(db)
    authHandler := handlers.NewAuthHandler(authService)

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

    // Routes
    e.POST("/login", authHandler.Login)

    // Health check
    e.GET("/health", func(c echo.Context) error {
        return c.JSON(200, map[string]string{
            "status": "healthy",
        })
    })

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8082"
    }
    e.Logger.Fatal(e.Start(":" + port))
}
