package models

import (
	"time"
)

type RetrospectiveItem struct {
	ID        string    `gorm:"primaryKey;type:text"`
	SprintID  string    `gorm:"index"`
	Type      string
	Content   string
	UserID    string
	CreatedAt time.Time

	Sprint Sprint `gorm:"foreignKey:SprintID;constraint:OnDelete:CASCADE"`
	User   User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
