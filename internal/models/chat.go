package models

import (
	"time"
)

type Chat struct {
	ID        string    `gorm:"primaryKey;type:text"`
	ProjectID *string   `gorm:"index"`
	Title     *string
	Type      string    `gorm:"default:'PROJECT'"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Project      *Project          `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	Messages     []Message         `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE"`
	Participants []ChatParticipant `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE"`
}

type ChatParticipant struct {
	ChatID string `gorm:"primaryKey;type:text"` // Composite PK part 1
	UserID string `gorm:"primaryKey;type:text"` // Composite PK part 2

	Chat Chat `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE"`
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type Message struct {
	ID        string    `gorm:"primaryKey;type:text"`
	ChatID    string    `gorm:"index"`
	UserID    string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time

	Chat Chat `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE"`
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
