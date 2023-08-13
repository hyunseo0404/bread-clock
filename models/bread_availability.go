package models

type BreadAvailability struct {
	BakeryID       int    `gorm:"primaryKey"`
	BreadID        int    `gorm:"primaryKey"`
	Available      bool   `gorm:"default:true"`
	AvailableHours string `gorm:"size:255"`
	Bakery         Bakery `gorm:"constraint:OnDelete:CASCADE"`
	Bread          Bread  `gorm:"constraint:OnDelete:CASCADE"`
}
