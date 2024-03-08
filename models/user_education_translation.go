package models

import (
	"time"

	"gorm.io/gorm"
)

type UserEducationTranslationCreateForm struct {
	LanguageID      int64  `gorm:"int64;not null" json:"language_id"`
	UserEducationID int64  `gorm:"int64;not null" json:"user_education_id"`
	Title           string `gorm:"varchar;not null;size:64" json:"title"`
	Description     string `gorm:"varchar;not null;size:300" json:"description"`
	Category        string `gorm:"varchar;not null;size:14" json:"category"`
	Location        string `gorm:"varchar;not null;size:128" json:"location"`
	LocationType    string `gorm:"varchar;not null;size:36" json:"location_type"`
	MonthStart      int    `gorm:"int;not null;size:2" json:"month_start"`
	MonthEnd        int    `gorm:"int;size:2" json:"month_end"`
	YearStart       int64  `gorm:"int64;not null;size:4" json:"year_start"`
	YearEnd         int64  `gorm:"int64;size:4" json:"year_end"`
}

type UserEducationTranslation struct {
	Id              int64     `gorm:"primaryKey" json:"id"`
	LanguageID      uint      // Foreign key to link user education translation to language
	UserEducationID uint      // Foreign key to link user education translation to user education
	Title           string    `gorm:"varchar;not null;size:64" json:"title"`
	Description     string    `gorm:"varchar;not null;size:300" json:"description"`
	Category        string    `gorm:"varchar;not null;size:14" json:"category"`
	Location        string    `gorm:"varchar;not null;size:128" json:"location"`
	LocationType    string    `gorm:"varchar;not null;size:36" json:"location_type"`
	MonthStart      int       `gorm:"int;not null;size:2" json:"month_start"`
	MonthEnd        int       `gorm:"int;size:2" json:"month_end"`
	YearStart       uint      `gorm:"uint;not null" json:"year_start"`
	YearEnd         uint      `gorm:"uint" json:"year_end"`
	CreatedAt       time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoUpdateTime" json:"updated_at"`
}

func (UserEducationTranslation) BeforeUpdate(db *gorm.DB) error {
	// manually updated at
	db.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
}
