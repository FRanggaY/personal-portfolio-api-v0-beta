package models

import (
	"time"

	"gorm.io/gorm"
)

type SkillTranslationResponse struct {
	ID                 int64     `json:"id"`
	LanguageID         int64     `json:"language_id"`
	Code               string    `json:"code"`
	Name               string    `json:"name"`
	ImageUrl           string    `json:"image_url"`
	Url                string    `json:"url"`
	IsExternalUrl      bool      `json:"is_external_url"`
	IsExternalImageUrl bool      `json:"is_external_image_url"`
	Description        string    `json:"description"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type SkillTranslation struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	SkillID     uint      // Foreign key to link user skill translation to skill
	LanguageID  uint      // Foreign key to link user skill translation to language
	Description string    `gorm:"varchar;not null;size:300" json:"description"`
	CreatedAt   time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoUpdateTime" json:"updated_at"`
}

func (SkillTranslation) BeforeUpdate(db *gorm.DB) error {
	// manually updated at
	db.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
}
