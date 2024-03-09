package models

import (
	"time"

	"gorm.io/gorm"
)

type UserExperienceCreateForm struct {
	UserID    int64 `gorm:"int64;not null" json:"user_id"`
	CompanyID int64 `gorm:"int64;not null" json:"company_id"`
}
type UserExperience struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    uint      // Foreign key to link user experience to user
	CompanyID uint      // Foreign key to link user experience to company
	CreatedAt time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoUpdateTime" json:"updated_at"`

	UserExperienceTranslations []UserExperienceTranslation `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (UserExperience) BeforeUpdate(db *gorm.DB) error {
	// manually updated at
	db.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
}
