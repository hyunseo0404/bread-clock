package models

type BakeryPhoto struct {
	ID       int    `gorm:"primaryKey"`
	BakeryID int    `gorm:"index:ux_bakery_photos_bakery_id_url,unique"`
	URL      string `gorm:"size:255;index:ux_bakery_photos_bakery_id_url,unique"`
	Bakery   Bakery `gorm:"constraint:OnDelete:CASCADE"`
}
