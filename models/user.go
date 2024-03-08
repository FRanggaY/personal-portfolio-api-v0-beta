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

type UserSkillCreateForm struct {
	UserId  int64 `gorm:"int64;not null" json:"user_id"`
	SkillId int64 `gorm:"int64;not null" json:"skill_id"`
}

type UserExperienceCreateForm struct {
	UserId    int64 `gorm:"int64;not null" json:"user_id"`
	CompanyId int64 `gorm:"int64;not null" json:"company_id"`
}

type UserEducationCreateForm struct {
	UserId   int64 `gorm:"int64;not null" json:"user_id"`
	SchoolId int64 `gorm:"int64;not null" json:"school_id"`
}

type UserPositionCreateForm struct {
	UserId int64  `gorm:"int64;not null" json:"user_id"`
	Title  string `gorm:"varchar;not null;size:64" json:"title"`
}

type UserLoginForm struct {
	Username string `gorm:"varchar;unique;not null;size:48" json:"username"`
	Password string `gorm:"varchar;unique;not null;size:300" json:"password"`
}

type User struct {
	Id        int64     `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"varchar;not null;size:48" json:"name"`
	Username  string    `gorm:"varchar;unique;not null;size:48" json:"username"`
	Password  string    `gorm:"varchar;unique;not null;size:300" json:"password"`
	CreatedAt time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoUpdateTime" json:"updated_at"`

	Positions   []UserPosition   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Attachments []UserAttachment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Experiences []UserExperience `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Educations  []UserEducation  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (User) BeforeUpdate(db *gorm.DB) error {
	// manually updated at
	db.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
}
