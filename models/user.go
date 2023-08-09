package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int    `gorm:"primaryKey"`
	Email     string `gorm:"size:255;not null;uniqueIndex:ux_users_email_provider" json:"email"`
	Provider  string `gorm:"size:50;not null;uniqueIndex:ux_users_email_provider" json:"provider"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `gorm:"index"`
}
