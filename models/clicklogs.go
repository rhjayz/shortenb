package models

import (
	"time"

	"gorm.io/gorm"
)

type Clicklogs struct {
	gorm.Model
	ShortlinkID uint
	ClickedAt   time.Time
	IPAddress   string
	Useragent   *string `gorm:"type:text"`
	Location    *string `gorm:"type:text"`

	Shortlink Shortlink `gorm:"constraint:OnDelete:CASCADE;"`
}
