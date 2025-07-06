package models

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	UserID    uint
	Token     string `gorm:"unique"`
	UserAgent string
	IPAddress string
	ExpiredAt time.Time

	User User `gorm:"constraint:OnDelete:CASCADE;"`
}
