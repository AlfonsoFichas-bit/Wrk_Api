package models

import (
	"time"
)

type Sprint struct {
	ID          string    `gorm:"primaryKey;type:text"`
	ProjectID   string    `gorm:"index"`
	Name        string    `gorm:"not null"`
	Description *string
	StartDate   time.Time
	EndDate     time.Time
	Status      string    `gorm:"default:'PLANNING'"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Project            Project             `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	UserStories        []UserStory         `gorm:"foreignKey:SprintID"`
	Tasks              []Task              `gorm:"foreignKey:SprintID"`
	RetrospectiveItems []RetrospectiveItem `gorm:"foreignKey:SprintID"`
	Evaluations        []Evaluation        `gorm:"foreignKey:SprintID"`
}

type UserStory struct {
	ID          string    `gorm:"primaryKey;type:text"`
	ProjectID   string    `gorm:"index"`
	Title       string    `gorm:"not null"`
	Description string    `gorm:"not null"`
	Acceptance  *string
	Priority    string    `gorm:"default:'MEDIUM'"`
	StoryPoints *int
	Status      string    `gorm:"default:'BACKLOG'"`
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time

	AssigneeID *string
	Assignee   *User   `gorm:"foreignKey:AssigneeID"`
	SprintID   *string
	Sprint     *Sprint `gorm:"foreignKey:SprintID"`
	Project    Project `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	Tasks      []Task  `gorm:"foreignKey:UserStoryID"`
}

type Task struct {
	ID          string    `gorm:"primaryKey;type:text"`
	ProjectID   string    `gorm:"index"`
	UserStoryID *string   `gorm:"index"`
	SprintID    *string   `gorm:"index"`
	Title       string    `gorm:"not null"`
	Description *string
	Priority    string    `gorm:"default:'MEDIUM'"`
	Status      string    `gorm:"default:'TODO'"`
	Deadline    *time.Time
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time

	AssigneeID *string
	Assignee   *User      `gorm:"foreignKey:AssigneeID"`
	Project    Project    `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	UserStory  *UserStory `gorm:"foreignKey:UserStoryID"`
	Sprint     *Sprint    `gorm:"foreignKey:SprintID"`
	Evaluations []Evaluation `gorm:"foreignKey:TaskID"`
}
