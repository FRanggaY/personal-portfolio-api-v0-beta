package models

import (
	"time"

	"gorm.io/gorm"
)

type SkillTranslation struct {
	Id          int64     `gorm:"primaryKey" json:"id"`
	SkillId     uint      // Foreign key to link user skill translation to skill
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
