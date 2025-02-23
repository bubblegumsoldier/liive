package models

import (
	"time"
	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	Title    string       `gorm:"type:varchar(255)"`
	IsGroup  bool         `gorm:"not null;default:false"`
	Members  []ChatMember `gorm:"foreignKey:ChatID"`
}

type ChatMember struct {
	gorm.Model
	ChatID      uint      `gorm:"not null"`
	UserID      uint      `gorm:"not null"`
	CurrentMessage string  `gorm:"type:text"`
	LastUpdated   time.Time
	JoinedAt    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	LeftAt      *time.Time
	Chat        Chat      `gorm:"foreignKey:ChatID"`
	User        User      `gorm:"foreignKey:UserID"`
} 