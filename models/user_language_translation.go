package models

import (
	"time"

	"gorm.io/gorm"
)

type UserLanguageTranslationResponse struct {
	ID          int64     `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	LogoUrl     string    `json:"logo_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserLanguageTranslationCreateForm struct {
	SelectLanguageID int64  `gorm:"int64;not null" json:"select_language_id"`
	LanguageID       int64  `gorm:"int64;not null" json:"language_id"`
	Title            string `gorm:"varchar;not null;size:30" json:"title"`
	Description      string `gorm:"varchar;not null;size:300" json:"description"`
}

type UserLanguageTranslation struct {
	ID             int64     `gorm:"primaryKey" json:"id"`
	LanguageID     uint      // Foreign key to link user experience translation to language
	UserLanguageID uint      // Foreign key to link user skill translation to skill
	Title          string    `gorm:"varchar;not null;size:30" json:"title"`
	Description    string    `gorm:"varchar;not null;size:300" json:"description"`
	CreatedAt      time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoUpdateTime" json:"updated_at"`
}

func (UserLanguageTranslation) BeforeUpdate(db *gorm.DB) error {
	// manually updated at
	db.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
}
