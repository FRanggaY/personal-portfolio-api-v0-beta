package models

import (
	"time"

	"gorm.io/gorm"
)

type UserSkillCreateForm struct {
	UserId  int64 `gorm:"int64;not null" json:"user_id"`
	SkillId int64 `gorm:"int64;not null" json:"skill_id"`
}

type UserSkill struct {
	Id        int64     `gorm:"primaryKey" json:"id"`
	UserID    uint      // Foreign key to link user experience to user
	SkillId   uint      // Foreign key to link user experience to skill
	CreatedAt time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoUpdateTime" json:"updated_at"`
}

func (UserSkill) BeforeUpdate(db *gorm.DB) error {
	// manually updated at
	db.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
}
