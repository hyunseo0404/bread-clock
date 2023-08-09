package models

type Bread struct {
	ID   int    `gorm:"primaryKey"`
	Name string `gorm:"size:20;not null;uniqueIndex"`
}
