package models

import "time"

type UserForm struct {
	Name     string `gorm:"varchar(48),unique,not null" json:"name"`
	Username string `gorm:"varchar(12),unique,not null" json:"username"`
	Password string `gorm:"varchar(300),unique,not null" json:"password"`
}

type User struct {
	Id        int64     `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"varchar(48),unique,not null" json:"name"`
	Username  string    `gorm:"varchar(12),unique,not null" json:"username"`
	Password  string    `gorm:"varchar(300),unique,not null" json:"password"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
