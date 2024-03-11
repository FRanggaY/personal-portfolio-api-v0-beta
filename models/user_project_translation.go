package models

import (
	"time"

	"gorm.io/gorm"
)

type ProjectTranslationResponse struct {
	ID                int64     `json:"id"`
	LanguageID        int64     `json:"language_id"`
	ProjectPlatformID int64     `json:"project_platform_id"`
	Name              string    `gorm:"varchar" json:"name"`
	Slug              string    `gorm:"varchar" json:"slug"`
	Description       string    `gorm:"varchar" json:"description"`
	ImageUrl          string    `gorm:"varchar" json:"image_url"`
	ProjectCreatedAt  time.Time `json:"project_created_at"`
	ProjectUpdatedAt  time.Time `json:"project_updated_at"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type UserProjectTranslationCreateForm struct {
	LanguageID    int64  `gorm:"int64;not null" json:"language_id"`
	UserProjectID int64  `gorm:"int64;not null" json:"user_project_id"`
	Name          string `gorm:"varchar;not null;size:48" json:"name"`
	Description   string `gorm:"varchar;not null;size:300" json:"description"`
}

type UserProjectTranslation struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	LanguageID    uint      // Foreign key
	UserProjectID uint      // Foreign key
	Name          string    `gorm:"varchar;not null;size:48" json:"name"`
	Description   string    `gorm:"varchar;not null;size:300" json:"description"`
	CreatedAt     time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoUpdateTime" json:"updated_at"`
}

func (UserProjectTranslation) BeforeUpdate(db *gorm.DB) error {
	// manually updated at
	db.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
}
