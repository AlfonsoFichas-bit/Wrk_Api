package models

import (
	"time"
)

type Rubric struct {
	ID          string    `gorm:"primaryKey;type:text"`
	ProjectID   *string   `gorm:"index"`
	Name        string    `gorm:"not null"`
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Project  *Project   `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	Criteria []Criteria `gorm:"foreignKey:RubricID;constraint:OnDelete:CASCADE"`
}

type Criteria struct {
	ID          string `gorm:"primaryKey;type:text"`
	RubricID    string `gorm:"index"`
	Name        string `gorm:"not null"`
	Description *string
	MaxScore    int `gorm:"default:100"`
	Weight      int `gorm:"default:1"`

	Rubric      Rubric               `gorm:"foreignKey:RubricID;constraint:OnDelete:CASCADE"`
	Evaluations []EvaluationCriteria `gorm:"foreignKey:CriteriaID"`
}

type Evaluation struct {
	ID          string    `gorm:"primaryKey;type:text"`
	ProjectID   string    `gorm:"index"`
	TaskID      *string   `gorm:"index"`
	SprintID    *string   `gorm:"index"`
	EvaluatorID string    `gorm:"index"`
	Status      string    `gorm:"default:'PENDING'"`
	Feedback    *string
	Score       *int
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Project   Project              `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	Task      *Task                `gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE"`
	Sprint    *Sprint              `gorm:"foreignKey:SprintID;constraint:OnDelete:CASCADE"`
	Evaluator User                 `gorm:"foreignKey:EvaluatorID"`
	Criteria  []EvaluationCriteria `gorm:"foreignKey:EvaluationID;constraint:OnDelete:CASCADE"`
}

type EvaluationCriteria struct {
	ID           string `gorm:"primaryKey;type:text"`
	EvaluationID string `gorm:"index:idx_eval_criteria,unique"`
	CriteriaID   string `gorm:"index:idx_eval_criteria,unique"`
	Score        int    `gorm:"default:0"`
	Comment      *string

	Evaluation Evaluation `gorm:"foreignKey:EvaluationID;constraint:OnDelete:CASCADE"`
	Criteria   Criteria   `gorm:"foreignKey:CriteriaID"`
}
