package models

type BreadAvailability struct {
	BakeryID  int    `gorm:"primaryKey"`
	BreadID   int    `gorm:"primaryKey"`
	Available bool   `gorm:"default:true"`
	Bakery    Bakery `gorm:"constraint:OnDelete:CASCADE"`
	Bread     Bread  `gorm:"constraint:OnDelete:CASCADE"`
}
