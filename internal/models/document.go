package models

import (
	"time"
)

type Document struct {
	ID        string    `gorm:"primaryKey;type:text"`
	ProjectID string    `gorm:"index"`
	Name      string
	URL       string
	Type      string
	Size      *int
	Version   int       `gorm:"default:1"`
	ParentID  *string   `gorm:"index"`
	UploadedAt time.Time

	Project Project   `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	Parent  *Document `gorm:"foreignKey:ParentID;constraint:OnDelete:SET NULL"`
	Versions []Document `gorm:"foreignKey:ParentID"`
}
