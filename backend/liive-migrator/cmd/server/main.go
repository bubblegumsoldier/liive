package main

import (
    "fmt"
    "log"
    "os"
    "time"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/liive/backend/shared/pkg/models"
)

func main() {
    log.Println("Starting database migration service...")

    // Database connection parameters
    host := os.Getenv("DB_HOST")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")
    
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=UTC",
        host, user, password, dbname)
    
    // Retry connection until database is ready
    var db *gorm.DB
    var err error
    maxRetries := 30 // More retries for initial setup
    for i := 0; i < maxRetries; i++ {
        db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err == nil {
            break
        }
        log.Printf("Waiting for database to be ready (attempt %d/%d): %v", i+1, maxRetries, err)
        time.Sleep(time.Second * 2)
    }
    if err != nil {
        log.Fatalf("Failed to connect to database after %d attempts: %v", maxRetries, err)
    }

    log.Println("Connected to database, starting migrations...")

    // Define all models to migrate
    models := []interface{}{
        &models.User{},
        &models.Role{},
        &models.Chat{},
        &models.ChatMember{},
    }

    // Perform migrations
    for _, model := range models {
        log.Printf("Migrating %T...", model)
        if err := db.AutoMigrate(model); err != nil {
            log.Fatalf("Failed to migrate %T: %v", model, err)
        }
    }

    // Create indexes
    log.Println("Creating indexes...")
    indexSQL := "CREATE UNIQUE INDEX IF NOT EXISTS idx_chat_members_chat_user ON chat_members (chat_id, user_id) WHERE deleted_at IS NULL"
    if err := db.Exec(indexSQL).Error; err != nil {
        log.Fatalf("Failed to create chat members index: %v", err)
    }

    log.Println("All migrations completed successfully!")
} 