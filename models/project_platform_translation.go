package models

import (
	"time"

	"gorm.io/gorm"
)

type ProjectPlatformTranslationResponse struct {
	ID                 int64     `json:"id"`
	LanguageID         int64     `json:"language_id"`
	Code               string    `json:"code"`
	Name               string    `json:"name"`
	Title              string    `json:"title"`
	ImageUrl           string    `json:"image_url"`
	Url                string    `json:"url"`
	IsExternalUrl      bool      `json:"is_external_url"`
	IsExternalImageUrl bool      `json:"is_external_image_url"`
	Description        string    `json:"description"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type ProjectPlatformTranslationCreateForm struct {
	LanguageID        int64  `gorm:"int64;not null" json:"language_id"`
	ProjectPlatformID int64  `gorm:"int64;not null" json:"ProjectPlatform_id"`
	Title             string `gorm:"varchar;not null;size:30" json:"title"`
	Description       string `gorm:"varchar;not null;size:300" json:"description"`
}
type ProjectPlatformTranslation struct {
	ID                int64     `gorm:"primaryKey" json:"id"`
	ProjectPlatformID uint      // Foreign key to link user ProjectPlatform translation to ProjectPlatform
	LanguageID        uint      // Foreign key to link user ProjectPlatform translation to language
	Title             string    `gorm:"varchar;not null;size:30" json:"title"`
	Description       string    `gorm:"varchar;not null;size:300" json:"description"`
	CreatedAt         time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoUpdateTime" json:"updated_at"`
}

func (ProjectPlatformTranslation) BeforeUpdate(db *gorm.DB) error {
	// manually updated at
	db.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
}
