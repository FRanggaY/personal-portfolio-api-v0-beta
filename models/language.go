package models

import (
	"time"

	"gorm.io/gorm"
)

type LanguageResponse struct {
	ID        int64     `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	LogoUrl   string    `json:"logo_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Language struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Code      string    `gorm:"varchar;size:5;unique;not null" json:"code"`
	Name      string    `gorm:"varchar;size:32;unique;not null" json:"name"`
	LogoUrl   string    `gorm:"varchar;size:300" json:"logo_url"`
	CreatedAt time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoUpdateTime" json:"updated_at"`

	UserExperienceTranslations  []UserExperienceTranslation  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserEducationTranslations   []UserEducationTranslation   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SkillTranslations           []SkillTranslation           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProjectPlatformTranslations []ProjectPlatformTranslation `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserLanguages               []UserLanguage               `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserLanguageTranslation     []UserLanguageTranslation    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserProjectTranslation      []UserProjectTranslation     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Language) BeforeUpdate(db *gorm.DB) error {
	// manually updated at
	db.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
}
