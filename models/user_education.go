package models

import (
	"time"

	"gorm.io/gorm"
)

type UserEducationCreateForm struct {
	SchoolID   int64 `gorm:"int64;not null" json:"school_id"`
	MonthStart int   `gorm:"int;not null;size:2" json:"month_start"`
	MonthEnd   int   `gorm:"int;size:2" json:"month_end"`
	YearStart  int64 `gorm:"int64;not null;size:4" json:"year_start"`
	YearEnd    int64 `gorm:"int64;size:4" json:"year_end"`
}

type UserEducation struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	UserID     uint      // Foreign key to link user education to user
	SchoolId   uint      // Foreign key to link user education to school
	MonthStart int       `gorm:"int;not null;size:2" json:"month_start"`
	MonthEnd   int       `gorm:"int;size:2" json:"month_end"`
	YearStart  uint      `gorm:"uint;not null" json:"year_start"`
	YearEnd    uint      `gorm:"uint" json:"year_end"`
	IsActive   bool      `gorm:"bool;default:1" json:"is_active"`
	CreatedAt  time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoUpdateTime" json:"updated_at"`

	UserEducationTranslations []UserEducationTranslation `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (UserEducation) BeforeUpdate(db *gorm.DB) error {
	// manually updated at
	db.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
}
