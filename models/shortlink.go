package models

import (
	"gorm.io/gorm"
)

type Shortlink struct {
	gorm.Model
	UserID       uint
	OriginalUrl  *string `gorm:"type:text"`
	ShortCode    string
	CustomDomain string

	User      User        `gorm:"constraint:OnDelete:CASCADE;"`
	Clicklogs []Clicklogs `gorm:"foreignKey:ShortlinkID"`
}
