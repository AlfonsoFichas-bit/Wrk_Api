package models

import (
	"time"
)

type Project struct {
	ID          string    `gorm:"primaryKey;type:text"`
	Name        string    `gorm:"not null"`
	Description *string
	Status      string    `gorm:"default:'ACTIVE'"`
	StartDate   *time.Time
	EndDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Relations
	OwnerID     string
	Owner       User      `gorm:"foreignKey:OwnerID"`
	Members     []ProjectMember `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	Sprints     []Sprint        `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	UserStories []UserStory     `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	Tasks       []Task          `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	Evaluations []Evaluation    `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	Rubrics     []Rubric        `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	Chats       []Chat          `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	Documents   []Document      `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
}

type ProjectMember struct {
	ID        string    `gorm:"primaryKey;type:text"`
	ProjectID string    `gorm:"index:idx_project_user,unique"`
	UserID    string    `gorm:"index:idx_project_user,unique"`
	Role      string    `gorm:"not null"`
	JoinedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`

	Project Project `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	User    User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
