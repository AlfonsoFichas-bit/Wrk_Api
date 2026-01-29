package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"primaryKey;type:text"`
	Email     string         `gorm:"uniqueIndex;not null"`
	Name      string         `gorm:"not null"`
	Password  string         `gorm:"not null"`
	Role      string         `gorm:"default:'TEAM_DEVELOPER'"`
	Avatar    *string        `gorm:"type:text"`
	Active    bool           `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relations
	Projects           []Project           `gorm:"foreignKey:OwnerID"`
	Tasks              []Task              `gorm:"foreignKey:AssigneeID"`
	Evaluations        []Evaluation        `gorm:"foreignKey:EvaluatorID"`
	UserStories        []UserStory         `gorm:"foreignKey:AssigneeID"`
	Messages           []Message           `gorm:"foreignKey:UserID"`
	Chats              []ChatParticipant   `gorm:"foreignKey:UserID"`
	ProjectMemberships []ProjectMember     `gorm:"foreignKey:UserID"`
	Notifications      []Notification      `gorm:"foreignKey:UserID"`
	RetrospectiveItems []RetrospectiveItem `gorm:"foreignKey:UserID"`
}
