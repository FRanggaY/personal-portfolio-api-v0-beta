package models

import (
	"time"

	"gorm.io/gorm"
)

type UserEducationCreateForm struct {
	UserID   int64 `gorm:"int64;not null" json:"user_id"`
	SchoolID int64 `gorm:"int64;not null" json:"school_id"`
}

type UserEducation struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    uint      // Foreign key to link user education to user
	SchoolId  uint      // Foreign key to link user education to school
	CreatedAt time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoUpdateTime" json:"updated_at"`

	UserEducationTranslations []UserEducationTranslation `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (UserEducation) BeforeUpdate(db *gorm.DB) error {
	// manually updated at
	db.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
}
