package main

import (
    "log"
    "os"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    _ "github.com/liive/backend/liive-message-storer/docs" // Import generated Swagger docs
    "github.com/liive/backend/shared/pkg/database"
    echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Message Storer API
// @version 1.0
// @description This is the message storage service API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8084
// @BasePath /
// @schemes http
func main() {
    // Initialize database
    _, err := database.InitDB()
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }

    // Create Echo instance
    e := echo.New()

    // Middleware
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.Use(middleware.CORS())

    // Swagger documentation
    e.GET("/swagger", echoSwagger.EchoWrapHandler())
    e.GET("/swagger/*", echoSwagger.EchoWrapHandler())

    // Routes
    e.GET("/health", func(c echo.Context) error {
        return c.JSON(200, map[string]string{
            "status": "healthy",
        })
    })

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8084"
    }
    e.Logger.Fatal(e.Start(":" + port))
}
