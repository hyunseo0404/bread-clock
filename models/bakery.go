package models

type Bakery struct {
	ID           int    `gorm:"primaryKey"`
	Name         string `gorm:"size:128;not null"`
	Address      string `gorm:"size:255"`
	Coordinates  string `gorm:"size:128"`
	OpeningHours string `gorm:"size:70"`
}
