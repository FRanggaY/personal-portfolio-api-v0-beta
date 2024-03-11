package models

import (
	"time"

	"gorm.io/gorm"
)

type EducationTranslationResponse struct {
	ID                       int64     `json:"id"`
	LanguageID               int64     `json:"language_id"`
	SchoolID                 int64     `json:"school_id"`
	Title                    string    `gorm:"varchar" json:"title"`
	Description              string    `gorm:"varchar" json:"description"`
	Category                 string    `gorm:"varchar" json:"category"`
	Location                 string    `gorm:"varchar" json:"location"`
	LocationType             string    `gorm:"varchar" json:"location_type"`
	MonthStart               int       `gorm:"int" json:"month_start"`
	MonthEnd                 int       `gorm:"int" json:"month_end"`
	YearStart                uint      `gorm:"uint" json:"year_start"`
	YearEnd                  uint      `gorm:"uint" json:"year_end"`
	SchoolCode               string    `gorm:"varchar" json:"school_code"`
	SchoolName               string    `gorm:"varchar" json:"school_name"`
	SchoolImageUrl           string    `gorm:"varchar" json:"school_image_url"`
	SchoolUrl                string    `gorm:"varchar" json:"school_url"`
	SchoolIsExternalUrl      bool      `gorm:"boolean" json:"school_is_external_url"`
	SchoolIsExternalImageUrl bool      `gorm:"boolean" json:"school_is_external_image_url"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

type UserEducationTranslationCreateForm struct {
	LanguageID   int64  `gorm:"int64;not null" json:"language_id"`
	SchoolID     int64  `gorm:"int64;not null" json:"school_id"`
	Title        string `gorm:"varchar;not null;size:64" json:"title"`
	Description  string `gorm:"varchar;not null;size:300" json:"description"`
	Category     string `gorm:"varchar;not null;size:14" json:"category"`
	Location     string `gorm:"varchar;not null;size:128" json:"location"`
	LocationType string `gorm:"varchar;not null;size:36" json:"location_type"`
}

type UserEducationTranslation struct {
	ID              int64     `gorm:"primaryKey" json:"id"`
	LanguageID      uint      // Foreign key to link user education translation to language
	UserEducationID uint      // Foreign key to link user education translation to user education
	Title           string    `gorm:"varchar;not null;size:64" json:"title"`
	Description     string    `gorm:"varchar;not null;size:300" json:"description"`
	Category        string    `gorm:"varchar;not null;size:14" json:"category"`
	Location        string    `gorm:"varchar;not null;size:128" json:"location"`
	LocationType    string    `gorm:"varchar;not null;size:36" json:"location_type"`
	CreatedAt       time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoUpdateTime" json:"updated_at"`
}

func (UserEducationTranslation) BeforeUpdate(db *gorm.DB) error {
	// manually updated at
	db.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
}
