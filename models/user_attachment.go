package models

import (
	"time"

	"gorm.io/gorm"
)

type UserAttachment struct {
	Id                 int64     `gorm:"primaryKey" json:"id"`
	UserID             uint      // Foreign key to link user attachment to user
	Title              string    `gorm:"varchar;not null;size:64" json:"title"`
	Category           string    `gorm:"varchar;size:36" json:"category"`
	ImageUrl           string    `gorm:"varchar;size:300" json:"image_url"`
	Url                string    `gorm:"varchar;size:300" json:"url"`
	IsExternalUrl      bool      `gorm:"boolean" json:"is_external_url"`
	IsExternalImageUrl bool      `gorm:"boolean" json:"is_external_image_url"`
	CreatedAt          time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoUpdateTime" json:"updated_at"`
}

func (UserAttachment) BeforeUpdate(db *gorm.DB) error {
	// manually updated at
	db.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
}
