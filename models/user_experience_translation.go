package models

import (
	"time"

	"gorm.io/gorm"
)

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
	YearStart        int    `gorm:"int;not null;size:4" json:"year_start"`
	YearEnd          int    `gorm:"int;size:4" json:"year_end"`
}

type UserExperienceTranslation struct {
	Id               int64     `gorm:"primaryKey" json:"id"`
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
	YearStart        int       `gorm:"int;not null;size:4" json:"year_start"`
	YearEnd          int       `gorm:"int;size:4" json:"year_end"`
	CreatedAt        time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoUpdateTime" json:"updated_at"`
}

func (UserExperienceTranslation) BeforeUpdate(db *gorm.DB) error {
	// manually updated at
	db.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
}
