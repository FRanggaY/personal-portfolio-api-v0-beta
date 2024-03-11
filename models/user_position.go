package models

import (
	"time"

	"gorm.io/gorm"
)

type UserPositionCreateForm struct {
	Title string `gorm:"varchar;not null;size:64" json:"title"`
}

type UserPosition struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    uint      // Foreign key to link user position to user
	Title     string    `gorm:"varchar;not null;size:64" json:"title"`
	IsActive  bool      `gorm:"bool;default:1" json:"is_active"`
	CreatedAt time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:current_timestamp;type:timestamp(0);autoUpdateTime" json:"updated_at"`
}

func (UserPosition) BeforeUpdate(db *gorm.DB) error {
	// manually updated at
	db.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
}
