package models

import (
	"time"
)

type Notification struct {
	ID        string    `gorm:"primaryKey;type:text"`
	UserID    string    `gorm:"index"`
	Title     string
	Message   string
	Type      string
	Read      bool      `gorm:"default:false"`
	CreatedAt time.Time

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
