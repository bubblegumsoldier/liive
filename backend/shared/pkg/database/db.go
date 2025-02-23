package database

import (
    "fmt"
    "log"
    "os"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/liive/backend/shared/pkg/models"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
    host := os.Getenv("DB_HOST")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")
    
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=UTC",
        host, user, password, dbname)
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %v", err)
    }

    // Auto Migrate the schema
    err = db.AutoMigrate(&models.User{}, &models.Role{})
    if err != nil {
        return nil, fmt.Errorf("failed to migrate database: %v", err)
    }

    DB = db
    log.Println("Database connection established")
    return db, nil
}

func GetDB() *gorm.DB {
    return DB
}
