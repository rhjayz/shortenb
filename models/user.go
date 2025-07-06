package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Email    string `gorm:"unique"`
	Password string
	Image    *string `gorm:"type:text"`

	Session   []Session   `gorm:"foreignKey:UserID"`
	Shortlink []Shortlink `gorm:"foreignKey:UserID"`
}
