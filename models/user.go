package models

import (
	"time"

	"gorm.io/gorm"
)

type UserCreateForm struct {
	Name     string `gorm:"varchar;not null;size:48" json:"name"`
	Username string `gorm:"varchar;unique;not null;size:48" json:"username"`
	Password string `gorm:"varchar;unique;not null;size:300" json:"password"`
}

type UserEditForm struct {
	Name     string `gorm:"varchar;not null;size:48" json:"name"`
	Username string `gorm:"varchar;unique;not null;size:48" json:"username"`
}

type UserLoginForm struct {
	Username string `gorm:"varchar;unique;not null;size:48" json:"username"`
	Password string `gorm:"varchar;unique;not null;size:300" json:"password"`
}

type User struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"varchar;not null;size:48" json:"name"`
	Username  string    `gorm:"varchar;unique;not null;size:48" json:"username"`
	Password  string    `gorm:"varchar;unique;not null;size:300" json:"password"`
	CreatedAt time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoUpdateTime" json:"updated_at"`

	Positions   []UserPosition   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Attachments []UserAttachment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Experiences []UserExperience `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Educations  []UserEducation  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Languages   []UserLanguage   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (User) BeforeUpdate(db *gorm.DB) error {
	// manually updated at
	db.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
}
