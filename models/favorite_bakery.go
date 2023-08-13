package models

type FavoriteBakery struct {
	UserID   int    `gorm:"primaryKey"`
	BakeryID int    `gorm:"primaryKey"`
	User     User   `gorm:"constraint:OnDelete:CASCADE"`
	Bakery   Bakery `gorm:"constraint:OnDelete:CASCADE"`
}
