package models

import (
	"time"

	"gorm.io/gorm"
)

type UserLanguageCreateForm struct {
	LanguageID int64 `gorm:"int64;not null" json:"language_id"`
}

type UserLanguage struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	UserID     uint      // Foreign key to link user experience to user
	LanguageID uint      // Foreign key to link user experience to skill
	CreatedAt  time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoUpdateTime" json:"updated_at"`

	UserLanguageTranslations []UserLanguageTranslation `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (UserLanguage) BeforeUpdate(db *gorm.DB) error {
	// manually updated at
	db.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
}
