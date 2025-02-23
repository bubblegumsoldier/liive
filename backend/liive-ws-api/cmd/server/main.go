package main

import (
    "log"
    "os"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "github.com/liive/backend/shared/pkg/database"
)

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

    // Routes
    e.GET("/health", func(c echo.Context) error {
        return c.JSON(200, map[string]string{
            "status": "healthy",
        })
    })

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    e.Logger.Fatal(e.Start(":" + port))
}
