package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int    `gorm:"primaryKey"`
	UID       string `gorm:"size:255;not null;uniqueIndex:ux_users_uid_provider"`
	Provider  string `gorm:"size:50;not null;uniqueIndex:ux_users_uid_provider" json:"provider"`
	Email     string `gorm:"size:255;not null" json:"email"`
	AvatarURL string `gorm:"size:255" json:"avatarUrl"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `gorm:"index"`
}
