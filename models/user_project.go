package models

import (
	"time"

	"gorm.io/gorm"
)

type UserProject struct {
	ID                int64     `gorm:"primaryKey" json:"id"`
	UserID            uint      // Foreign key
	ProjectPlatformID uint      // Foregin Key
	Slug              string    `gorm:"varchar;not null;size:126" json:"slug"`
	ProjectCreatedAt  time.Time `gorm:"default:current_timestamp;type:timestamp(0);" json:"project_created_at"`
	ProjectUpdatedAt  time.Time `gorm:"default:current_timestamp;type:timestamp(0);" json:"project_updated_at"`
	ImageUrl          string    `gorm:"varchar;size:300" json:"image_url"`
	IsActive          bool      `gorm:"bool;default:1" json:"is_active"`
	CreatedAt         time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoUpdateTime" json:"updated_at"`

	UserProjectTranslations []UserProjectTranslation `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserProjectAttachments  []UserProjectAttachment  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (UserProject) BeforeUpdate(db *gorm.DB) error {
	// manually updated at
	db.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
}
