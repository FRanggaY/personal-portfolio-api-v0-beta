package models

import (
	"time"

	"gorm.io/gorm"
)

type ExperienceTranslationResponse struct {
	ID                        int64     `json:"id"`
	LanguageID                int64     `json:"language_id"`
	CompanyID                 int64     `json:"company_id"`
	Title                     string    `gorm:"varchar"`
	Description               string    `gorm:"varchar"`
	Category                  string    `gorm:"varchar"`
	Location                  string    `gorm:"varchar"`
	LocationType              string    `gorm:"varchar"`
	Industry                  string    `gorm:"varchar"`
	MonthStart                int       `gorm:"int"`
	MonthEnd                  int       `gorm:"int"`
	YearStart                 uint      `gorm:"uint"`
	YearEnd                   uint      `gorm:"uint"`
	CompanyCode               string    `gorm:"varchar"`
	CompanyName               string    `gorm:"varchar"`
	CompanyImageUrl           string    `gorm:"varchar"`
	CompanyUrl                string    `gorm:"varchar"`
	CompanyIsExternalUrl      bool      `gorm:"boolean"`
	CompanyIsExternalImageUrl bool      `gorm:"boolean"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

type UserExperienceTranslationCreateForm struct {
	LanguageID       int64  `gorm:"int64;not null" json:"language_id"`
	UserExperienceID int64  `gorm:"int64;not null" json:"user_experience_id"`
	Title            string `gorm:"varchar;not null;size:64" json:"title"`
	Description      string `gorm:"varchar;not null;size:300" json:"description"`
	Category         string `gorm:"varchar;not null;size:14" json:"category"`
	Location         string `gorm:"varchar;not null;size:128" json:"location"`
	LocationType     string `gorm:"varchar;not null;size:36" json:"location_type"`
	Industry         string `gorm:"varchar;not null;size:300" json:"industry"`
	MonthStart       int    `gorm:"int;not null;size:2" json:"month_start"`
	MonthEnd         int    `gorm:"int;size:2" json:"month_end"`
	YearStart        int64  `gorm:"int64;not null;size:4" json:"year_start"`
	YearEnd          int64  `gorm:"int64;size:4" json:"year_end"`
}

type UserExperienceTranslation struct {
	ID               int64     `gorm:"primaryKey" json:"id"`
	LanguageID       uint      // Foreign key to link user experience translation to language
	UserExperienceID uint      // Foreign key to link user experience translation to user experience
	Title            string    `gorm:"varchar;not null;size:64" json:"title"`
	Description      string    `gorm:"varchar;not null;size:300" json:"description"`
	Category         string    `gorm:"varchar;not null;size:14" json:"category"`
	Location         string    `gorm:"varchar;not null;size:128" json:"location"`
	LocationType     string    `gorm:"varchar;not null;size:36" json:"location_type"`
	Industry         string    `gorm:"varchar;not null;size:300" json:"industry"`
	MonthStart       int       `gorm:"int;not null;size:2" json:"month_start"`
	MonthEnd         int       `gorm:"int;size:2" json:"month_end"`
	YearStart        uint      `gorm:"uint;not null" json:"year_start"`
	YearEnd          uint      `gorm:"uint" json:"year_end"`
	CreatedAt        time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoUpdateTime" json:"updated_at"`
}

func (UserExperienceTranslation) BeforeUpdate(db *gorm.DB) error {
	// manually updated at
	db.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
}
