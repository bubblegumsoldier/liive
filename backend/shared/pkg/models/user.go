package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"uniqueIndex;not null"`
	Username     string `gorm:"uniqueIndex;not null"`
	Password     string `gorm:"not null"`
	FirstName    string
	LastName     string
	LastLogin    time.Time
	IsActive     bool      `gorm:"default:true"`
	Roles        []Role    `gorm:"many2many:user_roles;"`
}

type Role struct {
	gorm.Model
	Name        string    `gorm:"uniqueIndex;not null"`
	Description string
	Users       []User    `gorm:"many2many:user_roles;"`
}
