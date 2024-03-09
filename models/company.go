package models

import (
	"time"

	"gorm.io/gorm"
)

type CompanyAllResponse struct {
	ID        int64     `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Company struct {
	ID                 int64     `gorm:"primaryKey" json:"id"`
	Code               string    `gorm:"varchar;unique;not null;size:5" json:"code"`
	Name               string    `gorm:"varchar;unique;not null;size:48" json:"name"`
	ImageUrl           string    `gorm:"varchar;size:300" json:"image_url"`
	Url                string    `gorm:"varchar;size:300" json:"url"`
	IsExternalUrl      bool      `gorm:"boolean" json:"is_external_url"`
	IsExternalImageUrl bool      `gorm:"boolean" json:"is_external_image_url"`
	Address            string    `gorm:"varchar;not null;size:300" json:"address"`
	CreatedAt          time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoUpdateTime" json:"updated_at"`

	UserExperiences []UserExperience `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Company) BeforeUpdate(db *gorm.DB) error {
	// manually updated at
	db.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
}
