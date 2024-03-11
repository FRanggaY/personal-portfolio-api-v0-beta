package models

import (
	"time"

	"gorm.io/gorm"
)

type UserExperienceCreateForm struct {
	CompanyID  int64 `gorm:"int64;not null" json:"company_id"`
	MonthStart int   `gorm:"int;not null;size:2" json:"month_start"`
	MonthEnd   int   `gorm:"int;size:2" json:"month_end"`
	YearStart  int64 `gorm:"int64;not null;size:4" json:"year_start"`
	YearEnd    int64 `gorm:"int64;size:4" json:"year_end"`
}
type UserExperience struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	UserID     uint      // Foreign key to link user experience to user
	CompanyID  uint      // Foreign key to link user experience to company
	MonthStart int       `gorm:"int;not null;size:2" json:"month_start"`
	MonthEnd   int       `gorm:"int;size:2" json:"month_end"`
	YearStart  uint      `gorm:"uint;not null" json:"year_start"`
	YearEnd    uint      `gorm:"uint" json:"year_end"`
	IsActive   bool      `gorm:"bool;default:1" json:"is_active"`
	CreatedAt  time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoUpdateTime" json:"updated_at"`

	UserExperienceTranslations []UserExperienceTranslation `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (UserExperience) BeforeUpdate(db *gorm.DB) error {
	// manually updated at
	db.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
}
