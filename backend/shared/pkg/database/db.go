package database

import (
    "fmt"
    "log"
    "os"
    "time"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
    host := os.Getenv("DB_HOST")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")
    
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=UTC",
        host, user, password, dbname)
    
    // Retry connection a few times
    var db *gorm.DB
    var err error
    maxRetries := 5
    for i := 0; i < maxRetries; i++ {
        db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err == nil {
            break
        }
        log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
        time.Sleep(time.Second * 3)
    }
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database after %d attempts: %v", maxRetries, err)
    }

    DB = db
    log.Println("Database connection established")
    return db, nil
}

func GetDB() *gorm.DB {
    return DB
}
