package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int          `gorm:"primaryKey" json:"id"`
	UID       string       `gorm:"size:255;not null;uniqueIndex:ux_users_uid_provider" json:"-"`
	Provider  string       `gorm:"size:50;not null;uniqueIndex:ux_users_uid_provider" json:"provider"`
	Email     string       `gorm:"size:255;not null" json:"email"`
	AvatarURL string       `gorm:"size:255" json:"avatarUrl"`
	CreatedAt time.Time    `json:"-"`
	UpdatedAt time.Time    `json:"-"`
	DeletedAt sql.NullTime `gorm:"index" json:"-"`
}
